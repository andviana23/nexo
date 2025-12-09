#!/bin/bash
################################################################################
# sync_asaas_to_nexo.sh
# 
# Script para sincronizar assinaturas do Asaas para o banco local do NEXO
# Insere diretamente no banco de dados PostgreSQL
# 
# Uso: ./scripts/sync_asaas_to_nexo.sh
################################################################################

set -o pipefail

# Configurações
API_URL="${API_URL:-http://localhost:8080}"
ASAAS_API_KEY='$aact_prod_000MzkwODA2MWY2OGM3MWRlMDU2NWM3MzJlNzZmNGZhZGY6OjEzMjMzMWIyLThmYzItNDUyZi04NzEwLWRjOWE2MjcwNjdhZTo6JGFhY2hfNWJiMjY4NTMtNGU0Yy00ODVjLTg1Y2ItZWM5M2IwMjhiMzY1'
TENANT_ID="00000000-0000-0000-0000-000000000001"
DATABASE_URL="postgresql://neondb_owner:npg_83COkAjHMotv@ep-winter-leaf-adhqz08p-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require"

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🔄 NEXO - Sincronização de Assinaturas do Asaas"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Autenticação no NEXO
echo -e "${BLUE}[1/6]${NC} Autenticando no NEXO..."
JWT_TOKEN=$(curl -s -X POST "${API_URL}/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d '{"email":"admin@teste.com","password":"Admin123!"}' | jq -r '.access_token')

if [ -z "$JWT_TOKEN" ] || [ "$JWT_TOKEN" == "null" ]; then
    echo -e "${RED}✗ Erro: Falha na autenticação${NC}"
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

# 2. Criar planos baseados nos valores únicos
echo ""
echo -e "${BLUE}[2/6]${NC} Criando planos no NEXO..."

PLANS_CREATED=0

# Definir os planos padrão
create_plan() {
    local nome="$1"
    local valor="$2"
    local qtd="$3"
    
    RESULT=$(nexo_request "POST" "/api/v1/plans" "{
        \"nome\": \"$nome\",
        \"descricao\": \"Plano importado do Asaas\",
        \"valor\": \"$valor\",
        \"qtd_servicos\": $qtd
    }" 2>/dev/null)
    
    PLAN_ID=$(echo "$RESULT" | jq -r '.id // empty')
    if [ -n "$PLAN_ID" ]; then
        echo -e "  ${GREEN}✓${NC} $nome (R\$ $valor) - ID: $PLAN_ID"
        ((PLANS_CREATED++))
    else
        echo -e "  ${YELLOW}⏭${NC} $nome já existe ou erro"
    fi
}

create_plan "Clube do Trato One - Corte" "67.50" 1
create_plan "Clube do Trato Corte 2x" "114.90" 2
create_plan "Clube do Trato Seletivo" "119.90" 2
create_plan "Clube do Trato One - Corte e Barba" "120.00" 1
create_plan "Clube do Trato Corte Exclusivo" "149.90" 4
create_plan "Clube do Trato Barba" "179.90" 4
create_plan "Clube do Trato Corte e Barba 2x" "199.90" 2
create_plan "Clube do Trato Gold Corte" "210.00" 4
create_plan "Clube do Trato Corte e Barba Exclusivo" "249.90" 4
create_plan "Clube do Trato Mult Premium" "299.90" 6
create_plan "Clube do Trato Gold Corte + Barba" "349.90" 8

echo -e "${GREEN}✓ Planos processados${NC}"

# Recarregar planos
EXISTING_PLANS=$(nexo_request "GET" "/api/v1/plans")
echo ""
echo "Planos disponíveis:"
echo "$EXISTING_PLANS" | jq -r '.[] | select(.ativo == true) | "  - \(.nome): R$ \(.valor) (ID: \(.id))"'

# 3. Buscar todas assinaturas ativas do Asaas
echo ""
echo -e "${BLUE}[3/6]${NC} Buscando assinaturas ativas do Asaas..."

OFFSET=0
LIMIT=100
ALL_SUBS="[]"

while true; do
    PAGE=$(curl -s -X GET "https://api.asaas.com/v3/subscriptions?status=ACTIVE&limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    
    DATA=$(echo "$PAGE" | jq '.data')
    COUNT=$(echo "$DATA" | jq 'length')
    
    ALL_SUBS=$(echo "$ALL_SUBS" "$DATA" | jq -s 'add')
    
    HAS_MORE=$(echo "$PAGE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ] || [ "$COUNT" -eq 0 ]; then
        break
    fi
    
    OFFSET=$((OFFSET + LIMIT))
done

TOTAL_SUBS=$(echo "$ALL_SUBS" | jq 'length')
echo -e "${GREEN}✓ Encontradas $TOTAL_SUBS assinaturas ativas${NC}"

# 4. Buscar todos clientes do Asaas
echo ""
echo -e "${BLUE}[4/6]${NC} Buscando clientes do Asaas..."

OFFSET=0
ALL_CUSTOMERS="[]"

while true; do
    PAGE=$(curl -s -X GET "https://api.asaas.com/v3/customers?limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    
    DATA=$(echo "$PAGE" | jq '.data')
    COUNT=$(echo "$DATA" | jq 'length')
    
    ALL_CUSTOMERS=$(echo "$ALL_CUSTOMERS" "$DATA" | jq -s 'add')
    
    HAS_MORE=$(echo "$PAGE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ] || [ "$COUNT" -eq 0 ]; then
        break
    fi
    
    OFFSET=$((OFFSET + LIMIT))
done

TOTAL_CUSTOMERS=$(echo "$ALL_CUSTOMERS" | jq 'length')
echo -e "${GREEN}✓ Encontrados $TOTAL_CUSTOMERS clientes${NC}"

# 5. Processar e importar
echo ""
echo -e "${BLUE}[5/6]${NC} Importando assinaturas..."

IMPORTED=0
SKIPPED=0
ERRORS=0

# Criar arquivo temporário com os dados
TEMP_FILE="/tmp/nexo_import_$$"

echo "$ALL_SUBS" | jq -c '.[]' | while read -r sub; do
    ASAAS_SUB_ID=$(echo "$sub" | jq -r '.id')
    ASAAS_CUSTOMER_ID=$(echo "$sub" | jq -r '.customer')
    VALOR=$(echo "$sub" | jq -r '.value')
    DESCRICAO=$(echo "$sub" | jq -r '.description // "Assinatura"' | tr -d '\r\n' | head -c 100)
    NEXT_DUE_DATE=$(echo "$sub" | jq -r '.nextDueDate')
    DATE_CREATED=$(echo "$sub" | jq -r '.dateCreated')
    BILLING_TYPE=$(echo "$sub" | jq -r '.billingType')
    
    # Buscar dados do cliente
    CUSTOMER=$(echo "$ALL_CUSTOMERS" | jq -r --arg id "$ASAAS_CUSTOMER_ID" '.[] | select(.id == $id)')
    CUSTOMER_NAME=$(echo "$CUSTOMER" | jq -r '.name // "Cliente Asaas"')
    CUSTOMER_PHONE=$(echo "$CUSTOMER" | jq -r '.mobilePhone // .phone // "00000000000"')
    
    # Limpar telefone
    CUSTOMER_PHONE=$(echo "$CUSTOMER_PHONE" | tr -dc '0-9')
    if [ -z "$CUSTOMER_PHONE" ] || [ ${#CUSTOMER_PHONE} -lt 10 ]; then
        CUSTOMER_PHONE="00000000000"
    fi
    
    # Verificar se cliente já existe no NEXO
    NEXO_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")
    NEXO_CUSTOMER_ID=$(echo "$NEXO_CUSTOMERS" | jq -r --arg phone "$CUSTOMER_PHONE" '.data[]? | select(.telefone == $phone) | .id' 2>/dev/null | head -1)
    
    # Se cliente não existe, criar
    if [ -z "$NEXO_CUSTOMER_ID" ]; then
        CREATE_RESULT=$(nexo_request "POST" "/api/v1/customers" "{
            \"nome\": \"$CUSTOMER_NAME\",
            \"telefone\": \"$CUSTOMER_PHONE\",
            \"tags\": [\"asaas-import\"]
        }")
        NEXO_CUSTOMER_ID=$(echo "$CREATE_RESULT" | jq -r '.id // empty')
        
        if [ -n "$NEXO_CUSTOMER_ID" ]; then
            echo -e "  ${GREEN}+${NC} Cliente criado: $CUSTOMER_NAME ($CUSTOMER_PHONE)"
        else
            echo -e "  ${RED}✗${NC} Erro ao criar cliente: $CUSTOMER_NAME"
            continue
        fi
    fi
    
    # Encontrar plano pelo valor - usar awk para formatar (locale-safe)
    VALOR_FMT=$(echo "$VALOR" | LC_ALL=C awk '{printf "%.2f", $1}')
    NEXO_PLAN_ID=$(echo "$EXISTING_PLANS" | jq -r --arg val "$VALOR_FMT" '.[] | select(.valor == $val and .ativo == true) | .id' | head -1)
    
    if [ -z "$NEXO_PLAN_ID" ]; then
        echo -e "  ${YELLOW}⏭${NC} Plano não encontrado para R\$ $VALOR - $CUSTOMER_NAME"
        continue
    fi
    
    # Determinar forma de pagamento
    FORMA_PAG="CARTAO"
    if [ "$BILLING_TYPE" == "PIX" ]; then
        FORMA_PAG="PIX"
    elif [ "$BILLING_TYPE" == "BOLETO" ]; then
        FORMA_PAG="DINHEIRO"
    fi
    
    # Criar assinatura
    CREATE_SUB=$(nexo_request "POST" "/api/v1/subscriptions" "{
        \"cliente_id\": \"$NEXO_CUSTOMER_ID\",
        \"plano_id\": \"$NEXO_PLAN_ID\",
        \"forma_pagamento\": \"$FORMA_PAG\"
    }")
    
    NEW_ID=$(echo "$CREATE_SUB" | jq -r '.id // empty')
    if [ -n "$NEW_ID" ]; then
        echo -e "  ${GREEN}✓${NC} Importada: $CUSTOMER_NAME - R\$ $VALOR ($FORMA_PAG)"
    else
        ERROR=$(echo "$CREATE_SUB" | jq -r '.message // .error // "Erro"')
        echo -e "  ${YELLOW}⏭${NC} $CUSTOMER_NAME: $ERROR"
    fi
done

# 6. Resumo final
echo ""
echo -e "${BLUE}[6/6]${NC} Verificando resultado..."

FINAL_SUBS=$(nexo_request "GET" "/api/v1/subscriptions")
FINAL_COUNT=$(echo "$FINAL_SUBS" | jq 'length')
FINAL_METRICS=$(nexo_request "GET" "/api/v1/subscriptions/metrics")

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 RESULTADO FINAL"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  Total de assinaturas no NEXO: ${GREEN}$FINAL_COUNT${NC}"
echo ""
echo "  Métricas:"
echo "$FINAL_METRICS" | jq '.'
echo ""
echo -e "${GREEN}✓ Sincronização concluída!${NC}"
