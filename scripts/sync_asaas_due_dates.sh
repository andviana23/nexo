#!/bin/bash
# =============================================================================
# NEXO - Sincronização de Datas de Vencimento com Asaas
# 
# Este script busca a data do próximo pagamento (nextDueDate) de cada 
# assinatura ativa no Asaas e atualiza o campo data_vencimento no NEXO.
# =============================================================================

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "═══════════════════════════════════════════════════════════════"
echo "  🔄 NEXO - Sincronização de Datas de Vencimento com Asaas"
echo "═══════════════════════════════════════════════════════════════"
echo -e "${NC}"

# Carregar variáveis de ambiente
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../backend/.env"

if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
    echo -e "${GREEN}✓ Variáveis de ambiente carregadas${NC}"
else
    echo -e "${RED}✗ Arquivo .env não encontrado em $ENV_FILE${NC}"
    exit 1
fi

# Verificar variáveis obrigatórias
if [ -z "$ASAAS_API_KEY" ]; then
    echo -e "${RED}✗ ASAAS_API_KEY não configurada${NC}"
    exit 1
fi

if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}✗ DATABASE_URL não configurada${NC}"
    exit 1
fi

echo -e "${BLUE}Conectando ao Asaas...${NC}"

# Contadores
TOTAL_ASAAS=0
ATUALIZADOS=0
NAO_ENCONTRADOS=0
ERROS=0

# Função para atualizar data de vencimento no banco
update_due_date() {
    local asaas_sub_id=$1
    local next_due_date=$2
    
    # Atualizar no banco usando psql
    RESULT=$(psql "$DATABASE_URL" -t -c "
        UPDATE subscriptions 
        SET data_vencimento = '$next_due_date'::timestamptz,
            updated_at = NOW()
        WHERE asaas_subscription_id = '$asaas_sub_id'
          AND tenant_id = 'e2e00000-0000-0000-0000-000000000001'
        RETURNING id;
    " 2>/dev/null | tr -d ' ')
    
    if [ -n "$RESULT" ]; then
        return 0
    else
        return 1
    fi
}

# Buscar todas as assinaturas ativas do Asaas (paginado)
OFFSET=0
LIMIT=100

while true; do
    echo -e "\n${BLUE}Buscando assinaturas do Asaas (offset: $OFFSET)...${NC}"
    
    RESPONSE=$(curl -s "https://api.asaas.com/v3/subscriptions?status=ACTIVE&limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    
    # Verificar se há dados
    COUNT=$(echo "$RESPONSE" | jq '.data | length')
    
    if [ "$COUNT" -eq 0 ]; then
        echo -e "${YELLOW}Nenhuma assinatura encontrada nesta página.${NC}"
        break
    fi
    
    echo -e "${GREEN}Encontradas $COUNT assinaturas nesta página${NC}"
    
    # Processar cada assinatura
    echo "$RESPONSE" | jq -c '.data[]' | while read -r SUB; do
        ASAAS_SUB_ID=$(echo "$SUB" | jq -r '.id')
        NEXT_DUE_DATE=$(echo "$SUB" | jq -r '.nextDueDate')
        VALUE=$(echo "$SUB" | jq -r '.value')
        CUSTOMER_ID=$(echo "$SUB" | jq -r '.customer')
        
        if [ "$NEXT_DUE_DATE" != "null" ] && [ -n "$NEXT_DUE_DATE" ]; then
            echo -n "  📅 $ASAAS_SUB_ID (R$ $VALUE) -> Vencimento: $NEXT_DUE_DATE ... "
            
            if update_due_date "$ASAAS_SUB_ID" "$NEXT_DUE_DATE"; then
                echo -e "${GREEN}✓${NC}"
            else
                echo -e "${YELLOW}Não encontrado no NEXO${NC}"
            fi
        fi
    done
    
    TOTAL_ASAAS=$((TOTAL_ASAAS + COUNT))
    
    # Próxima página
    HAS_MORE=$(echo "$RESPONSE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ]; then
        break
    fi
    
    OFFSET=$((OFFSET + LIMIT))
done

echo -e "\n${BLUE}"
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 RESULTADO DA SINCRONIZAÇÃO"
echo "═══════════════════════════════════════════════════════════════"
echo -e "${NC}"

# Verificar resultado no banco
echo -e "${BLUE}Verificando assinaturas atualizadas no NEXO...${NC}"

STATS=$(psql "$DATABASE_URL" -t -c "
    SELECT 
        COUNT(*) as total,
        COUNT(*) FILTER (WHERE data_vencimento IS NOT NULL) as com_vencimento,
        COUNT(*) FILTER (WHERE data_vencimento IS NULL) as sem_vencimento,
        COUNT(*) FILTER (WHERE data_vencimento > NOW()) as futuras,
        COUNT(*) FILTER (WHERE data_vencimento <= NOW() AND data_vencimento IS NOT NULL) as vencidas
    FROM subscriptions
    WHERE tenant_id = 'e2e00000-0000-0000-0000-000000000001'
      AND status = 'ATIVO';
")

echo "$STATS" | while IFS='|' read -r TOTAL COM_VENC SEM_VENC FUTURAS VENCIDAS; do
    TOTAL=$(echo "$TOTAL" | tr -d ' ')
    COM_VENC=$(echo "$COM_VENC" | tr -d ' ')
    SEM_VENC=$(echo "$SEM_VENC" | tr -d ' ')
    FUTURAS=$(echo "$FUTURAS" | tr -d ' ')
    VENCIDAS=$(echo "$VENCIDAS" | tr -d ' ')
    
    echo -e "  Total de assinaturas ativas: ${GREEN}$TOTAL${NC}"
    echo -e "  Com data de vencimento:      ${GREEN}$COM_VENC${NC}"
    echo -e "  Sem data de vencimento:      ${YELLOW}$SEM_VENC${NC}"
    echo -e "  Vencimentos futuros:         ${GREEN}$FUTURAS${NC}"
    echo -e "  Vencimentos passados:        ${RED}$VENCIDAS${NC}"
done

# Mostrar próximas renovações
echo -e "\n${BLUE}📅 Próximas Renovações (7 dias):${NC}"
psql "$DATABASE_URL" -c "
    SELECT 
        c.nome as cliente,
        s.valor,
        TO_CHAR(s.data_vencimento, 'DD/MM/YYYY') as vencimento
    FROM subscriptions s
    JOIN clientes c ON s.cliente_id = c.id
    WHERE s.tenant_id = 'e2e00000-0000-0000-0000-000000000001'
      AND s.status = 'ATIVO'
      AND s.data_vencimento IS NOT NULL
      AND s.data_vencimento BETWEEN NOW() AND NOW() + INTERVAL '7 days'
    ORDER BY s.data_vencimento
    LIMIT 10;
"

echo -e "\n${GREEN}✓ Sincronização concluída!${NC}"
