#!/bin/bash
################################################################################
# sync_asaas_v2.sh
# 
# Script SIMPLIFICADO para sincronizar assinaturas do Asaas para o NEXO
# Usa mapeamento direto de valores Asaas -> Planos NEXO
# 
# Uso: ./scripts/sync_asaas_v2.sh
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
echo "  🔄 NEXO - Sincronização de Assinaturas (v2)"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Autenticação
echo -e "${BLUE}[1/5]${NC} Autenticando..."
JWT_TOKEN=$(curl -s -X POST "${API_URL}/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d '{"email":"admin@teste.com","password":"Admin123!"}' | jq -r '.access_token')

if [ -z "$JWT_TOKEN" ] || [ "$JWT_TOKEN" == "null" ]; then
    echo -e "${RED}✗ Erro: Falha na autenticação${NC}"
    exit 1
fi
echo -e "${GREEN}✓ OK${NC}"

# Helper
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

# 2. Buscar planos e criar mapa
echo ""
echo -e "${BLUE}[2/5]${NC} Carregando planos..."

NEXO_PLANS=$(nexo_request "GET" "/api/v1/plans")

# Criar mapa de valor -> plano_id (manualmente por clareza)
# Formato Asaas -> Formato NEXO -> Plan ID
get_plan_id() {
    local asaas_value="$1"
    
    # Converter para formato NEXO (sempre 2 casas decimais)
    local nexo_value
    case "$asaas_value" in
        "67.5"|"67.50")   nexo_value="67.50" ;;
        "114.9"|"114.90") nexo_value="114.90" ;;
        "119.9"|"119.90") nexo_value="119.90" ;;
        "120"|"120.0"|"120.00") nexo_value="120.00" ;;
        "149.9"|"149.90") nexo_value="149.90" ;;
        "179.9"|"179.90") nexo_value="179.90" ;;
        "199.9"|"199.90") nexo_value="199.90" ;;
        "210"|"210.0"|"210.00") nexo_value="210.00" ;;
        "249.9"|"249.90") nexo_value="249.90" ;;
        "299.9"|"299.90") nexo_value="299.90" ;;
        "349.9"|"349.90") nexo_value="349.90" ;;
        *) nexo_value="" ;;
    esac
    
    if [ -n "$nexo_value" ]; then
        echo "$NEXO_PLANS" | jq -r --arg val "$nexo_value" '.[] | select(.valor == $val and .ativo == true) | .id' | head -1
    fi
}

echo -e "${GREEN}✓ Planos carregados${NC}"

# 3. Buscar dados Asaas
echo ""
echo -e "${BLUE}[3/5]${NC} Buscando dados do Asaas..."

# Buscar assinaturas
ALL_SUBS="[]"
OFFSET=0
while true; do
    PAGE=$(curl -s "https://api.asaas.com/v3/subscriptions?status=ACTIVE&limit=100&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    DATA=$(echo "$PAGE" | jq '.data')
    ALL_SUBS=$(echo "$ALL_SUBS" "$DATA" | jq -s 'add')
    [ "$(echo "$PAGE" | jq '.hasMore')" != "true" ] && break
    OFFSET=$((OFFSET + 100))
done
TOTAL_SUBS=$(echo "$ALL_SUBS" | jq 'length')
echo -e "  Assinaturas: ${GREEN}$TOTAL_SUBS${NC}"

# Buscar clientes
ALL_CUSTOMERS="[]"
OFFSET=0
while true; do
    PAGE=$(curl -s "https://api.asaas.com/v3/customers?limit=100&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    DATA=$(echo "$PAGE" | jq '.data')
    ALL_CUSTOMERS=$(echo "$ALL_CUSTOMERS" "$DATA" | jq -s 'add')
    [ "$(echo "$PAGE" | jq '.hasMore')" != "true" ] && break
    OFFSET=$((OFFSET + 100))
done
TOTAL_CUSTOMERS=$(echo "$ALL_CUSTOMERS" | jq 'length')
echo -e "  Clientes: ${GREEN}$TOTAL_CUSTOMERS${NC}"

# 4. Importar
echo ""
echo -e "${BLUE}[4/5]${NC} Importando assinaturas..."
echo ""

IMPORTED=0
SKIPPED=0
ERRORS=0

# Processar cada assinatura
echo "$ALL_SUBS" | jq -c '.[]' | while read -r sub; do
    ASAAS_CUSTOMER_ID=$(echo "$sub" | jq -r '.customer')
    VALOR=$(echo "$sub" | jq -r '.value')
    BILLING_TYPE=$(echo "$sub" | jq -r '.billingType')
    
    # Buscar dados do cliente
    CUSTOMER=$(echo "$ALL_CUSTOMERS" | jq -r --arg id "$ASAAS_CUSTOMER_ID" '.[] | select(.id == $id)')
    CUSTOMER_NAME=$(echo "$CUSTOMER" | jq -r '.name // "Cliente"')
    CUSTOMER_PHONE=$(echo "$CUSTOMER" | jq -r '.mobilePhone // .phone // ""' | tr -dc '0-9')
    
    # Garantir telefone único usando cpfCnpj como fallback
    if [ -z "$CUSTOMER_PHONE" ] || [ ${#CUSTOMER_PHONE} -lt 10 ]; then
        # Usar ID do cliente Asaas para gerar telefone único
        UNIQUE_ID=$(echo "$ASAAS_CUSTOMER_ID" | tr -dc '0-9' | tail -c 11)
        CUSTOMER_PHONE="31$UNIQUE_ID"
    fi
    
    # Buscar plano
    PLAN_ID=$(get_plan_id "$VALOR")
    
    if [ -z "$PLAN_ID" ]; then
        echo -e "  ${YELLOW}⏭${NC} Plano não encontrado: R\$ $VALOR - $CUSTOMER_NAME"
        ((SKIPPED++))
        continue
    fi
    
    # Verificar/criar cliente
    NEXO_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")
    NEXO_CUSTOMER_ID=$(echo "$NEXO_CUSTOMERS" | jq -r --arg phone "$CUSTOMER_PHONE" '.data[]? | select(.telefone == $phone) | .id' 2>/dev/null | head -1)
    
    if [ -z "$NEXO_CUSTOMER_ID" ] || [ "$NEXO_CUSTOMER_ID" == "null" ]; then
        CREATE_RESULT=$(nexo_request "POST" "/api/v1/customers" "{
            \"nome\": \"$CUSTOMER_NAME\",
            \"telefone\": \"$CUSTOMER_PHONE\",
            \"tags\": [\"asaas-import\"]
        }")
        NEXO_CUSTOMER_ID=$(echo "$CREATE_RESULT" | jq -r '.id // empty')
        
        if [ -z "$NEXO_CUSTOMER_ID" ] || [ "$NEXO_CUSTOMER_ID" == "null" ]; then
            ERROR_MSG=$(echo "$CREATE_RESULT" | jq -r '.message // .error // "?"')
            echo -e "  ${RED}✗${NC} Erro cliente: $CUSTOMER_NAME - $ERROR_MSG"
            ((ERRORS++))
            continue
        fi
        echo -e "  ${CYAN}+${NC} Cliente: $CUSTOMER_NAME"
    fi
    
    # Forma de pagamento
    FORMA_PAG="CARTAO"
    case "$BILLING_TYPE" in
        "PIX") FORMA_PAG="PIX" ;;
        "BOLETO") FORMA_PAG="DINHEIRO" ;;
    esac
    
    # Criar assinatura
    CREATE_SUB=$(nexo_request "POST" "/api/v1/subscriptions" "{
        \"cliente_id\": \"$NEXO_CUSTOMER_ID\",
        \"plano_id\": \"$PLAN_ID\",
        \"forma_pagamento\": \"$FORMA_PAG\"
    }")
    
    NEW_ID=$(echo "$CREATE_SUB" | jq -r '.id // empty')
    if [ -n "$NEW_ID" ] && [ "$NEW_ID" != "null" ]; then
        echo -e "  ${GREEN}✓${NC} $CUSTOMER_NAME - R\$ $VALOR"
        ((IMPORTED++))
    else
        ERROR_MSG=$(echo "$CREATE_SUB" | jq -r '.message // .error // "?"')
        if [[ "$ERROR_MSG" == *"já possui"* ]]; then
            echo -e "  ${YELLOW}⏭${NC} $CUSTOMER_NAME já tem assinatura"
            ((SKIPPED++))
        else
            echo -e "  ${RED}✗${NC} $CUSTOMER_NAME: $ERROR_MSG"
            ((ERRORS++))
        fi
    fi
done

# 5. Resumo
echo ""
echo -e "${BLUE}[5/5]${NC} Resultado..."

FINAL_SUBS=$(nexo_request "GET" "/api/v1/subscriptions")
FINAL_METRICS=$(nexo_request "GET" "/api/v1/subscriptions/metrics")
FINAL_COUNT=$(echo "$FINAL_SUBS" | jq 'length')

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 RESULTADO"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  Assinaturas no NEXO: ${GREEN}$FINAL_COUNT${NC}"
echo ""
echo "  Métricas:"
echo "$FINAL_METRICS" | jq -r '
  "    Ativas: \(.assinaturas_ativas // 0)",
  "    MRR: R$ \(.mrr // "0.00")"
'
echo ""
echo -e "${GREEN}✓ Concluído!${NC}"
