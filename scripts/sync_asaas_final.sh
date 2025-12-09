#!/bin/bash
################################################################################
# sync_asaas_final.sh
# 
# Script FINAL para sincronizar assinaturas do Asaas para o NEXO
# Usa jq para fazer match de valores com tolerância decimal
# 
# Uso: ./scripts/sync_asaas_final.sh
################################################################################

set -o pipefail

# Configurações
API_URL="${API_URL:-http://localhost:8080}"
ASAAS_API_KEY='$aact_prod_000MzkwODA2MWY2OGM3MWRlMDU2NWM3MzJlNzZmNGZhZGY6OjEzMjMzMWIyLThmYzItNDUyZi04NzEwLWRjOWE2MjcwNjdhZTo6JGFhY2hfNWJiMjY4NTMtNGU0Yy00ODVjLTg1Y2ItZWM5M2IwMjhiMzY1'
TENANT_ID="00000000-0000-0000-0000-000000000001"

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🔄 NEXO - Sincronização FINAL de Assinaturas do Asaas"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Autenticação
echo -e "${BLUE}[1/5]${NC} Autenticando no NEXO..."
LOGIN_RESP=$(curl -s -X POST "${API_URL}/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d '{"email":"admin@teste.com","password":"Admin123!"}')

JWT_TOKEN=$(echo "$LOGIN_RESP" | jq -r '.access_token')

if [ -z "$JWT_TOKEN" ] || [ "$JWT_TOKEN" == "null" ]; then
    echo -e "${RED}✗ Erro: Falha na autenticação${NC}"
    echo "$LOGIN_RESP"
    exit 1
fi
echo -e "${GREEN}✓ Autenticação OK${NC}"

# Helper para fazer requests ao NEXO
nexo_request() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    
    if [ -n "$data" ]; then
        curl -s -X "$method" "${API_URL}${endpoint}" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -X "$method" "${API_URL}${endpoint}" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "X-Tenant-ID: $TENANT_ID"
    fi
}

# 2. Buscar dados do Asaas
echo ""
echo -e "${BLUE}[2/5]${NC} Buscando dados do Asaas..."

# Assinaturas
echo -n "  Assinaturas: "
OFFSET=0
LIMIT=100
ALL_SUBS="[]"
while true; do
    PAGE=$(curl -s -X GET "https://api.asaas.com/v3/subscriptions?status=ACTIVE&limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    DATA=$(echo "$PAGE" | jq '.data')
    ALL_SUBS=$(echo "$ALL_SUBS" "$DATA" | jq -s 'add')
    HAS_MORE=$(echo "$PAGE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ]; then break; fi
    OFFSET=$((OFFSET + LIMIT))
done
TOTAL_SUBS=$(echo "$ALL_SUBS" | jq 'length')
echo -e "${GREEN}$TOTAL_SUBS${NC}"

# Clientes
echo -n "  Clientes: "
OFFSET=0
ALL_CUSTOMERS="[]"
while true; do
    PAGE=$(curl -s -X GET "https://api.asaas.com/v3/customers?limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    DATA=$(echo "$PAGE" | jq '.data')
    ALL_CUSTOMERS=$(echo "$ALL_CUSTOMERS" "$DATA" | jq -s 'add')
    HAS_MORE=$(echo "$PAGE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ]; then break; fi
    OFFSET=$((OFFSET + LIMIT))
done
TOTAL_CUSTOMERS=$(echo "$ALL_CUSTOMERS" | jq 'length')
echo -e "${GREEN}$TOTAL_CUSTOMERS${NC}"

# 3. Carregar dados do NEXO
echo ""
echo -e "${BLUE}[3/5]${NC} Carregando dados do NEXO..."

NEXO_PLANS=$(nexo_request "GET" "/api/v1/plans")
NEXO_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")
NEXO_SUBS=$(nexo_request "GET" "/api/v1/subscriptions")

echo -e "  Planos: $(echo "$NEXO_PLANS" | jq '[.[] | select(.ativo == true)] | length')"
echo -e "  Clientes: $(echo "$NEXO_CUSTOMERS" | jq '.data | length // 0' 2>/dev/null || echo "0")"
echo -e "  Assinaturas: $(echo "$NEXO_SUBS" | jq 'length')"

# Criar mapa de planos por valor (como número para comparação)
# Formato: valor_int -> plan_id
# Ex: 6750 -> uuid (67.50)
create_plan_map() {
    echo "$NEXO_PLANS" | jq -r '.[] | select(.ativo == true) | "\((.valor | gsub("\\."; "") | tonumber))=\(.id)"'
}
PLAN_MAP=$(create_plan_map)

# 4. Importar assinaturas
echo ""
echo -e "${BLUE}[4/5]${NC} Importando assinaturas..."
echo ""

IMPORTED=0
SKIPPED=0
ERRORS=0

# Processar cada assinatura
echo "$ALL_SUBS" | jq -c '.[]' | while read -r sub; do
    ASAAS_SUB_ID=$(echo "$sub" | jq -r '.id')
    ASAAS_CUSTOMER_ID=$(echo "$sub" | jq -r '.customer')
    VALOR=$(echo "$sub" | jq -r '.value')
    BILLING_TYPE=$(echo "$sub" | jq -r '.billingType')
    
    # Converter valor para inteiro (ex: 67.5 -> 6750, 119.9 -> 11990)
    VALOR_INT=$(echo "$VALOR" | awk '{printf "%.0f", $1 * 100}')
    
    # Buscar dados do cliente no Asaas
    CUSTOMER=$(echo "$ALL_CUSTOMERS" | jq -r --arg id "$ASAAS_CUSTOMER_ID" '.[] | select(.id == $id)')
    CUSTOMER_NAME=$(echo "$CUSTOMER" | jq -r '.name // "Cliente Asaas"')
    CUSTOMER_PHONE=$(echo "$CUSTOMER" | jq -r '.mobilePhone // .phone // ""' | tr -dc '0-9')
    
    # Limpar telefone
    if [ -z "$CUSTOMER_PHONE" ] || [ ${#CUSTOMER_PHONE} -lt 10 ]; then
        CUSTOMER_PHONE="31900000000"
    fi
    
    # Verificar se cliente já existe no NEXO (por telefone)
    NEXO_CUSTOMER_ID=$(echo "$NEXO_CUSTOMERS" | jq -r --arg phone "$CUSTOMER_PHONE" '.data[]? | select(.telefone == $phone) | .id' 2>/dev/null | head -1)
    
    # Se não existe, criar cliente
    if [ -z "$NEXO_CUSTOMER_ID" ] || [ "$NEXO_CUSTOMER_ID" == "null" ]; then
        CREATE_RESULT=$(nexo_request "POST" "/api/v1/customers" "{
            \"nome\": \"$CUSTOMER_NAME\",
            \"telefone\": \"$CUSTOMER_PHONE\",
            \"tags\": [\"asaas-import\"]
        }")
        NEXO_CUSTOMER_ID=$(echo "$CREATE_RESULT" | jq -r '.id // empty')
        
        if [ -n "$NEXO_CUSTOMER_ID" ] && [ "$NEXO_CUSTOMER_ID" != "null" ]; then
            echo -e "  ${CYAN}+${NC} Cliente: $CUSTOMER_NAME"
            # Atualizar cache
            NEXO_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")
        else
            echo -e "  ${RED}✗${NC} Erro ao criar cliente: $CUSTOMER_NAME"
            ((ERRORS++))
            continue
        fi
    fi
    
    # Encontrar plano pelo valor (comparando como inteiro)
    # Procurar por valor exato ou valor com .00 adicionado
    NEXO_PLAN_ID=""
    
    # Tentar match exato
    for plan_entry in $PLAN_MAP; do
        PLAN_VALOR_INT=$(echo "$plan_entry" | cut -d= -f1)
        PLAN_ID=$(echo "$plan_entry" | cut -d= -f2)
        
        if [ "$VALOR_INT" == "$PLAN_VALOR_INT" ]; then
            NEXO_PLAN_ID="$PLAN_ID"
            break
        fi
    done
    
    if [ -z "$NEXO_PLAN_ID" ]; then
        echo -e "  ${YELLOW}⏭${NC} Plano não encontrado: R\$ $VALOR - $CUSTOMER_NAME"
        ((SKIPPED++))
        continue
    fi
    
    # Determinar forma de pagamento
    FORMA_PAG="CARTAO"
    case "$BILLING_TYPE" in
        "PIX") FORMA_PAG="PIX" ;;
        "BOLETO") FORMA_PAG="DINHEIRO" ;;
        "CREDIT_CARD") FORMA_PAG="CARTAO" ;;
    esac
    
    # Criar assinatura
    CREATE_SUB=$(nexo_request "POST" "/api/v1/subscriptions" "{
        \"cliente_id\": \"$NEXO_CUSTOMER_ID\",
        \"plano_id\": \"$NEXO_PLAN_ID\",
        \"forma_pagamento\": \"$FORMA_PAG\"
    }")
    
    NEW_ID=$(echo "$CREATE_SUB" | jq -r '.id // empty')
    if [ -n "$NEW_ID" ] && [ "$NEW_ID" != "null" ]; then
        echo -e "  ${GREEN}✓${NC} Assinatura: $CUSTOMER_NAME - R\$ $VALOR ($FORMA_PAG)"
        ((IMPORTED++))
    else
        ERROR=$(echo "$CREATE_SUB" | jq -r '.message // .error // "Erro desconhecido"')
        # Se já existe, é OK
        if [[ "$ERROR" == *"já possui"* ]] || [[ "$ERROR" == *"already"* ]]; then
            echo -e "  ${YELLOW}⏭${NC} $CUSTOMER_NAME já tem assinatura"
            ((SKIPPED++))
        else
            echo -e "  ${RED}✗${NC} $CUSTOMER_NAME: $ERROR"
            ((ERRORS++))
        fi
    fi
done

# 5. Resumo final
echo ""
echo -e "${BLUE}[5/5]${NC} Resultado final..."

FINAL_SUBS=$(nexo_request "GET" "/api/v1/subscriptions")
FINAL_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")
FINAL_METRICS=$(nexo_request "GET" "/api/v1/subscriptions/metrics")

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 RESUMO DA IMPORTAÇÃO"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  Clientes no NEXO: ${GREEN}$(echo "$FINAL_CUSTOMERS" | jq '.meta.total // (.data | length)' 2>/dev/null || echo "?")${NC}"
echo -e "  Assinaturas no NEXO: ${GREEN}$(echo "$FINAL_SUBS" | jq 'length')${NC}"
echo ""
echo "  Métricas de Assinaturas:"
echo "$FINAL_METRICS" | jq -r '
  "    Ativas: \(.assinaturas_ativas // 0)",
  "    MRR: R$ \(.mrr // "0.00")",
  "    Churn Rate: \(.churn_rate // "0.00")%"
'
echo ""
echo -e "${GREEN}✓ Sincronização concluída!${NC}"
echo ""
