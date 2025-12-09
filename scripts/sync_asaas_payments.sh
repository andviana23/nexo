#!/bin/bash
# ==============================================================================
# NEXO - Sincronização de Pagamentos Asaas
# ==============================================================================
# Este script sincroniza pagamentos do Asaas para o NEXO:
# - CONFIRMED → subscription_payments (métricas MRR)
# - RECEIVED → operacoes_caixa (fluxo de caixa)
# ==============================================================================

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuração
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_DIR/backend"

# Carregar variáveis de ambiente
if [ -f "$BACKEND_DIR/.env" ]; then
    export $(grep -v '^#' "$BACKEND_DIR/.env" | xargs)
fi

# Verificar variáveis
if [ -z "$ASAAS_API_KEY" ]; then
    echo -e "${RED}Erro: ASAAS_API_KEY não configurada${NC}"
    exit 1
fi

if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}Erro: DATABASE_URL não configurada${NC}"
    exit 1
fi

# Tenant ID (Trato de Barbados)
TENANT_ID="e2e00000-0000-0000-0000-000000000001"

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  NEXO - Sincronização Asaas         ${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# ==============================================================================
# Função: Buscar pagamentos CONFIRMED do Asaas
# ==============================================================================
fetch_confirmed_payments() {
    echo -e "${YELLOW}Buscando pagamentos CONFIRMED...${NC}"
    
    curl -s "https://api.asaas.com/v3/payments?status=CONFIRMED&limit=100" \
        -H "access_token: $ASAAS_API_KEY" > /tmp/asaas_confirmed.json
    
    TOTAL=$(jq '.totalCount // 0' /tmp/asaas_confirmed.json)
    echo -e "${GREEN}Encontrados: $TOTAL pagamentos CONFIRMED${NC}"
}

# ==============================================================================
# Função: Buscar pagamentos RECEIVED do Asaas
# ==============================================================================
fetch_received_payments() {
    echo -e "${YELLOW}Buscando pagamentos RECEIVED...${NC}"
    
    curl -s "https://api.asaas.com/v3/payments?status=RECEIVED&limit=100" \
        -H "access_token: $ASAAS_API_KEY" > /tmp/asaas_received.json
    
    TOTAL=$(jq '.totalCount // 0' /tmp/asaas_received.json)
    echo -e "${GREEN}Encontrados: $TOTAL pagamentos RECEIVED${NC}"
}

# ==============================================================================
# Função: Gerar SQL para CONFIRMED (subscription_payments)
# ==============================================================================
generate_confirmed_sql() {
    echo -e "${YELLOW}Gerando SQL para subscription_payments...${NC}"
    
    # Extrair pagamentos CONFIRMED com subscription vinculada
    jq -r --arg tenant "$TENANT_ID" '
        .data[] | 
        select(.subscription != null) |
        "INSERT INTO subscription_payments (id, tenant_id, subscription_id, asaas_payment_id, valor, forma_pagamento, status, data_pagamento, codigo_transacao, observacao, created_at)
SELECT 
    gen_random_uuid(),
    '"'"'\($tenant)'"'"'::uuid,
    s.id,
    '"'"'\(.id)'"'"',
    \(.value),
    '"'"'\(.billingType)'"'"',
    '"'"'CONFIRMED'"'"',
    '"'"'\(.confirmedDate)'"'"'::timestamp,
    '"'"'\(.invoiceNumber // "")'"'"',
    '"'"'\(.description // "")'"'"',
    NOW()
FROM subscriptions s
WHERE s.asaas_subscription_id = '"'"'\(.subscription)'"'"'
  AND s.tenant_id = '"'"'\($tenant)'"'"'::uuid
  AND NOT EXISTS (
    SELECT 1 FROM subscription_payments sp 
    WHERE sp.asaas_payment_id = '"'"'\(.id)'"'"' 
      AND sp.tenant_id = '"'"'\($tenant)'"'"'::uuid
  );"
    ' /tmp/asaas_confirmed.json > /tmp/confirmed_sql.sql
    
    LINES=$(wc -l < /tmp/confirmed_sql.sql)
    echo -e "${GREEN}Gerados $LINES registros para subscription_payments${NC}"
}

# ==============================================================================
# Função: Gerar SQL para RECEIVED (operacoes_caixa via API)
# Para operacoes_caixa precisa de caixa_id, então vamos usar subscription_payments
# com status RECEIVED para rastreamento
# ==============================================================================
generate_received_sql() {
    echo -e "${YELLOW}Gerando SQL para pagamentos RECEIVED...${NC}"
    
    # Pagamentos RECEIVED também vão para subscription_payments com status RECEIVED
    # O fluxo de caixa será alimentado separadamente
    jq -r --arg tenant "$TENANT_ID" '
        .data[] | 
        select(.subscription != null) |
        "INSERT INTO subscription_payments (id, tenant_id, subscription_id, asaas_payment_id, valor, forma_pagamento, status, data_pagamento, codigo_transacao, observacao, created_at)
SELECT 
    gen_random_uuid(),
    '"'"'\($tenant)'"'"'::uuid,
    s.id,
    '"'"'\(.id)'"'"',
    \(.value),
    '"'"'\(.billingType)'"'"',
    '"'"'RECEIVED'"'"',
    '"'"'\(.paymentDate)'"'"'::timestamp,
    '"'"'\(.invoiceNumber // "")'"'"',
    '"'"'\(.description // "")'"'"',
    NOW()
FROM subscriptions s
WHERE s.asaas_subscription_id = '"'"'\(.subscription)'"'"'
  AND s.tenant_id = '"'"'\($tenant)'"'"'::uuid
  AND NOT EXISTS (
    SELECT 1 FROM subscription_payments sp 
    WHERE sp.asaas_payment_id = '"'"'\(.id)'"'"' 
      AND sp.tenant_id = '"'"'\($tenant)'"'"'::uuid
  );"
    ' /tmp/asaas_received.json > /tmp/received_sql.sql
    
    LINES=$(wc -l < /tmp/received_sql.sql)
    echo -e "${GREEN}Gerados $LINES registros para RECEIVED${NC}"
}

# ==============================================================================
# Função: Mostrar resumo dos pagamentos
# ==============================================================================
show_summary() {
    echo ""
    echo -e "${BLUE}=== RESUMO ===${NC}"
    
    # CONFIRMED no mês atual
    CONFIRMED_MES=$(jq '[.data[] | select(.confirmedDate >= "2025-12-01")]' /tmp/asaas_confirmed.json)
    CONFIRMED_COUNT=$(echo "$CONFIRMED_MES" | jq 'length')
    CONFIRMED_TOTAL=$(echo "$CONFIRMED_MES" | jq '[.[].value] | add // 0')
    
    echo -e "${GREEN}CONFIRMED (Dezembro/2025):${NC}"
    echo -e "  Quantidade: $CONFIRMED_COUNT"
    echo -e "  Total: R$ $CONFIRMED_TOTAL"
    
    # RECEIVED no mês atual  
    RECEIVED_MES=$(jq '[.data[] | select(.paymentDate >= "2025-12-01")]' /tmp/asaas_received.json)
    RECEIVED_COUNT=$(echo "$RECEIVED_MES" | jq 'length')
    RECEIVED_TOTAL=$(echo "$RECEIVED_MES" | jq '[.[].value] | add // 0')
    
    echo ""
    echo -e "${GREEN}RECEIVED (Dezembro/2025):${NC}"
    echo -e "  Quantidade: $RECEIVED_COUNT"
    echo -e "  Total: R$ $RECEIVED_TOTAL"
    
    echo ""
    echo -e "${BLUE}=== ARQUIVOS GERADOS ===${NC}"
    echo -e "  /tmp/asaas_confirmed.json"
    echo -e "  /tmp/asaas_received.json"
    echo -e "  /tmp/confirmed_sql.sql"
    echo -e "  /tmp/received_sql.sql"
}

# ==============================================================================
# Main
# ==============================================================================
main() {
    fetch_confirmed_payments
    echo ""
    fetch_received_payments
    echo ""
    generate_confirmed_sql
    echo ""
    generate_received_sql
    echo ""
    show_summary
    
    echo ""
    echo -e "${YELLOW}Para aplicar as alterações, execute:${NC}"
    echo -e "  cat /tmp/confirmed_sql.sql | psql \"\$DATABASE_URL\""
    echo -e "  cat /tmp/received_sql.sql | psql \"\$DATABASE_URL\""
    echo ""
    echo -e "${BLUE}Ou use a ferramenta pgsql do VS Code para executar os SQLs${NC}"
}

main "$@"
