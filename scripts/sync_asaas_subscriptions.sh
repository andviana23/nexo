#!/bin/bash
################################################################################
# sync_asaas_subscriptions.sh
# 
# Script para sincronizar assinaturas do Asaas para o banco local do NEXO
# 
# Uso: ./scripts/sync_asaas_subscriptions.sh
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
NC='\033[0m'

# Contadores
PLANS_CREATED=0
CUSTOMERS_CREATED=0
SUBSCRIPTIONS_IMPORTED=0
SUBSCRIPTIONS_SKIPPED=0
ERRORS=0

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🔄 NEXO - Sincronização de Assinaturas do Asaas"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Autenticação no NEXO
echo -e "${BLUE}[1/5]${NC} Autenticando no NEXO..."
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
echo -e "${BLUE}[2/5]${NC} Criando planos no NEXO..."

# Definir os planos padrão baseados nos valores mais comuns
declare -A PLAN_MAPPING
PLAN_MAPPING["67.50"]="Clube do Trato One - Corte"
PLAN_MAPPING["114.90"]="Clube do Trato Corte 2x"
PLAN_MAPPING["119.90"]="Clube do Trato Seletivo"
PLAN_MAPPING["120.00"]="Clube do Trato One - Corte e Barba"
PLAN_MAPPING["149.90"]="Clube do Trato Corte Exclusivo"
PLAN_MAPPING["179.90"]="Clube do Trato Barba"
PLAN_MAPPING["199.90"]="Clube do Trato Corte e Barba 2x"
PLAN_MAPPING["210.00"]="Clube do Trato Gold Corte"
PLAN_MAPPING["249.90"]="Clube do Trato Corte e Barba Exclusivo"
PLAN_MAPPING["299.90"]="Clube do Trato Mult Premium"
PLAN_MAPPING["349.90"]="Clube do Trato Gold Corte + Barba"

# Buscar planos existentes
EXISTING_PLANS=$(nexo_request "GET" "/api/v1/plans")

# Criar cada plano
for valor in "${!PLAN_MAPPING[@]}"; do
    nome="${PLAN_MAPPING[$valor]}"
    
    # Verificar se já existe
    EXISTS=$(echo "$EXISTING_PLANS" | jq -r --arg nome "$nome" '.[] | select(.nome == $nome) | .id')
    
    if [ -n "$EXISTS" ]; then
        echo -e "  ${YELLOW}⏭${NC} Plano '$nome' já existe (ID: $EXISTS)"
    else
        RESULT=$(nexo_request "POST" "/api/v1/plans" "{
            \"nome\": \"$nome\",
            \"descricao\": \"Plano importado do Asaas\",
            \"valor\": \"$valor\",
            \"qtd_servicos\": 4
        }")
        
        PLAN_ID=$(echo "$RESULT" | jq -r '.id // empty')
        if [ -n "$PLAN_ID" ]; then
            echo -e "  ${GREEN}✓${NC} Plano criado: $nome (R$ $valor)"
            ((PLANS_CREATED++))
        else
            echo -e "  ${RED}✗${NC} Erro ao criar plano: $nome"
            ((ERRORS++))
        fi
    fi
done

# Recarregar planos
EXISTING_PLANS=$(nexo_request "GET" "/api/v1/plans")
echo -e "${GREEN}✓ $PLANS_CREATED planos criados${NC}"

# 3. Buscar clientes existentes e criar os que faltam
echo ""
echo -e "${BLUE}[3/5]${NC} Sincronizando clientes do Asaas..."

# Buscar todos os clientes do Asaas (paginado)
OFFSET=0
LIMIT=100
ALL_ASAAS_CUSTOMERS="[]"

while true; do
    PAGE=$(curl -s -X GET "https://api.asaas.com/v3/customers?limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    
    DATA=$(echo "$PAGE" | jq '.data')
    COUNT=$(echo "$DATA" | jq 'length')
    
    ALL_ASAAS_CUSTOMERS=$(echo "$ALL_ASAAS_CUSTOMERS" "$DATA" | jq -s 'add')
    
    HAS_MORE=$(echo "$PAGE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ] || [ "$COUNT" -eq 0 ]; then
        break
    fi
    
    OFFSET=$((OFFSET + LIMIT))
done

TOTAL_ASAAS_CUSTOMERS=$(echo "$ALL_ASAAS_CUSTOMERS" | jq 'length')
echo -e "  Encontrados $TOTAL_ASAAS_CUSTOMERS clientes no Asaas"

# Buscar clientes existentes no NEXO
EXISTING_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")

# 4. Buscar assinaturas ativas do Asaas
echo ""
echo -e "${BLUE}[4/5]${NC} Buscando assinaturas ativas do Asaas..."

OFFSET=0
ALL_ASAAS_SUBSCRIPTIONS="[]"

while true; do
    PAGE=$(curl -s -X GET "https://api.asaas.com/v3/subscriptions?status=ACTIVE&limit=$LIMIT&offset=$OFFSET" \
        -H "access_token: $ASAAS_API_KEY")
    
    DATA=$(echo "$PAGE" | jq '.data')
    COUNT=$(echo "$DATA" | jq 'length')
    
    ALL_ASAAS_SUBSCRIPTIONS=$(echo "$ALL_ASAAS_SUBSCRIPTIONS" "$DATA" | jq -s 'add')
    
    HAS_MORE=$(echo "$PAGE" | jq '.hasMore')
    if [ "$HAS_MORE" != "true" ] || [ "$COUNT" -eq 0 ]; then
        break
    fi
    
    OFFSET=$((OFFSET + LIMIT))
done

TOTAL_ASAAS_SUBS=$(echo "$ALL_ASAAS_SUBSCRIPTIONS" | jq 'length')
echo -e "  Encontradas $TOTAL_ASAAS_SUBS assinaturas ativas no Asaas"

# 5. Importar assinaturas
echo ""
echo -e "${BLUE}[5/5]${NC} Importando assinaturas para o NEXO..."

# Processar cada assinatura
echo "$ALL_ASAAS_SUBSCRIPTIONS" | jq -c '.[]' | while read -r sub; do
    ASAAS_SUB_ID=$(echo "$sub" | jq -r '.id')
    ASAAS_CUSTOMER_ID=$(echo "$sub" | jq -r '.customer')
    VALOR=$(echo "$sub" | jq -r '.value')
    DESCRICAO=$(echo "$sub" | jq -r '.description // "Assinatura"')
    NEXT_DUE_DATE=$(echo "$sub" | jq -r '.nextDueDate')
    DATE_CREATED=$(echo "$sub" | jq -r '.dateCreated')
    
    # Verificar se já existe no NEXO
    EXISTING=$(nexo_request "GET" "/api/v1/subscriptions" | jq -r --arg id "$ASAAS_SUB_ID" '.[] | select(.asaas_subscription_id == $id) | .id // empty' 2>/dev/null)
    
    if [ -n "$EXISTING" ]; then
        echo -e "  ${YELLOW}⏭${NC} Assinatura $ASAAS_SUB_ID já existe"
        continue
    fi
    
    # Buscar dados do cliente no Asaas
    CUSTOMER_DATA=$(echo "$ALL_ASAAS_CUSTOMERS" | jq -r --arg id "$ASAAS_CUSTOMER_ID" '.[] | select(.id == $id)')
    CUSTOMER_NAME=$(echo "$CUSTOMER_DATA" | jq -r '.name // "Cliente Asaas"')
    CUSTOMER_PHONE=$(echo "$CUSTOMER_DATA" | jq -r '.mobilePhone // .phone // ""')
    
    # Verificar se cliente existe no NEXO (por telefone)
    if [ -n "$CUSTOMER_PHONE" ]; then
        NEXO_CUSTOMER_ID=$(echo "$EXISTING_CUSTOMERS" | jq -r --arg phone "$CUSTOMER_PHONE" '.data[]? | select(.telefone == $phone) | .id // empty' 2>/dev/null)
    fi
    
    # Se cliente não existe, criar
    if [ -z "$NEXO_CUSTOMER_ID" ]; then
        CREATE_CUSTOMER=$(nexo_request "POST" "/api/v1/customers" "{
            \"nome\": \"$CUSTOMER_NAME\",
            \"telefone\": \"${CUSTOMER_PHONE:-11999999999}\",
            \"tags\": [\"asaas-import\"]
        }")
        NEXO_CUSTOMER_ID=$(echo "$CREATE_CUSTOMER" | jq -r '.id // empty')
        
        if [ -n "$NEXO_CUSTOMER_ID" ]; then
            echo -e "  ${GREEN}+${NC} Cliente criado: $CUSTOMER_NAME"
            ((CUSTOMERS_CREATED++))
            # Recarregar lista de clientes
            EXISTING_CUSTOMERS=$(nexo_request "GET" "/api/v1/customers?limit=1000")
        else
            echo -e "  ${RED}✗${NC} Erro ao criar cliente: $CUSTOMER_NAME"
            ((ERRORS++))
            continue
        fi
    fi
    
    # Encontrar plano correspondente pelo valor
    VALOR_FORMATTED=$(printf "%.2f" "$VALOR")
    NEXO_PLAN_ID=$(echo "$EXISTING_PLANS" | jq -r --arg val "$VALOR_FORMATTED" '.[] | select(.valor == $val and .ativo == true) | .id' | head -1)
    
    if [ -z "$NEXO_PLAN_ID" ]; then
        echo -e "  ${YELLOW}⚠${NC} Plano não encontrado para valor R$ $VALOR (Sub: $ASAAS_SUB_ID)"
        ((SUBSCRIPTIONS_SKIPPED++))
        continue
    fi
    
    # Criar assinatura no NEXO
    CREATE_SUB=$(nexo_request "POST" "/api/v1/subscriptions" "{
        \"cliente_id\": \"$NEXO_CUSTOMER_ID\",
        \"plano_id\": \"$NEXO_PLAN_ID\",
        \"forma_pagamento\": \"CARTAO\",
        \"asaas_subscription_id\": \"$ASAAS_SUB_ID\",
        \"asaas_customer_id\": \"$ASAAS_CUSTOMER_ID\"
    }")
    
    NEW_SUB_ID=$(echo "$CREATE_SUB" | jq -r '.id // empty')
    if [ -n "$NEW_SUB_ID" ]; then
        echo -e "  ${GREEN}✓${NC} Assinatura importada: $CUSTOMER_NAME - R$ $VALOR"
        ((SUBSCRIPTIONS_IMPORTED++))
    else
        ERROR_MSG=$(echo "$CREATE_SUB" | jq -r '.message // .error // "Erro desconhecido"')
        echo -e "  ${RED}✗${NC} Erro: $ERROR_MSG (Cliente: $CUSTOMER_NAME)"
        ((ERRORS++))
    fi
done

# Resumo final
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 RESUMO DA SINCRONIZAÇÃO"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  Planos criados:       ${GREEN}$PLANS_CREATED${NC}"
echo -e "  Clientes criados:     ${GREEN}$CUSTOMERS_CREATED${NC}"
echo -e "  Assinaturas importadas: ${GREEN}$SUBSCRIPTIONS_IMPORTED${NC}"
echo -e "  Assinaturas puladas:  ${YELLOW}$SUBSCRIPTIONS_SKIPPED${NC}"
echo -e "  Erros:                ${RED}$ERRORS${NC}"
echo ""

# Verificar resultado final
echo "=== Verificando resultado ==="
FINAL_SUBS=$(nexo_request "GET" "/api/v1/subscriptions" | jq 'length')
FINAL_METRICS=$(nexo_request "GET" "/api/v1/subscriptions/metrics")
echo -e "Total de assinaturas no NEXO: ${GREEN}$FINAL_SUBS${NC}"
echo "$FINAL_METRICS" | jq

if [ "$ERRORS" -gt 0 ]; then
    echo -e "\n${YELLOW}⚠ Sincronização concluída com alguns erros${NC}"
    exit 1
else
    echo -e "\n${GREEN}✓ Sincronização concluída com sucesso!${NC}"
    exit 0
fi
