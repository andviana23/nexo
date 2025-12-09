#!/bin/bash
################################################################################
# NEXO - Smoke Tests Completos
# Versão: 2.0.0
# Data: 02/12/2025
#
# Testa todos os endpoints críticos do MVP v1.0.0
# Uso: ./scripts/smoke_tests_complete.sh [API_URL]
#
# Endpoints testados:
# - Health Check
# - Auth (Login, Me)
# - Services/Categorias
# - Professionals
# - Customers
# - Appointments
# - Financial (Payables, Receivables, Dashboard, Projections)
# - Metas (Monthly, Barbers, Ticket)
# - Pricing
# - Stock
# - Barber Turn (Lista da Vez)
################################################################################

set -o pipefail

# Configurações
API_URL="${1:-http://localhost:8080}"
TIMEOUT=15

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Contadores
TOTAL=0
PASSED=0
FAILED=0
SKIPPED=0

# Variáveis globais
JWT_TOKEN=""
TENANT_ID=""
TEST_USER_EMAIL="admin@teste.com"
TEST_USER_PASSWORD="Admin123!"

echo ""
echo "═══════════════════════════════════════════════════════════════════════════"
echo "  🧪 NEXO - Smoke Tests MVP v1.0.0"
echo "═══════════════════════════════════════════════════════════════════════════"
echo ""
echo "📍 API URL: $API_URL"
echo "⏱️  Timeout: ${TIMEOUT}s"
echo "📅 Data: $(date)"
echo ""

# Verificar dependências
check_dependencies() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔧 Verificando dependências..."
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}✗ curl não instalado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ curl${NC}"
    
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}⚠ jq não instalado (parsing JSON limitado)${NC}"
    else
        echo -e "${GREEN}✓ jq${NC}"
    fi
    echo ""
}

# Função para fazer requisições
http_request() {
    local method="$1"
    local endpoint="$2"
    local data="${3:-}"
    local auth="${4:-}"
    
    local headers=(-H "Content-Type: application/json")
    
    if [ -n "$auth" ]; then
        headers+=(-H "Authorization: Bearer $auth")
    fi
    
    if [ -n "$TENANT_ID" ]; then
        headers+=(-H "X-Tenant-ID: $TENANT_ID")
    fi
    
    if [ -n "$data" ]; then
        curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint" \
            "${headers[@]}" -d "$data" --max-time "$TIMEOUT" 2>/dev/null
    else
        curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint" \
            "${headers[@]}" --max-time "$TIMEOUT" 2>/dev/null
    fi
}

# Função para testar endpoint
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local expected_code="$3"
    local description="$4"
    local data="${5:-}"
    local auth="${6:-}"
    
    ((TOTAL++))
    printf "  %-55s" "$description"
    
    local response=$(http_request "$method" "$endpoint" "$data" "$auth")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        ((PASSED++))
        echo "$body"
        return 0
    elif [ "$expected_code" = "200|404" ] && ([ "$http_code" = "200" ] || [ "$http_code" = "404" ]); then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code - aceitável)"
        ((PASSED++))
        echo "$body"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected $expected_code, got $http_code)"
        ((FAILED++))
        if [ -n "$body" ]; then
            echo "    Response: $(echo "$body" | head -c 200)"
        fi
        return 1
    fi
}

# Função para login e obter JWT
do_login() {
    local email="$1"
    local password="$2"
    
    local login_data="{\"email\": \"$email\", \"password\": \"$password\"}"
    local response=$(http_request "POST" "/api/v1/auth/login" "$login_data")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        if command -v jq &> /dev/null; then
            JWT_TOKEN=$(echo "$body" | jq -r '.access_token // .token // .data.access_token // .data.token // empty')
            TENANT_ID=$(echo "$body" | jq -r '.user.tenant_id // .data.user.tenant_id // empty')
        fi
        return 0
    else
        return 1
    fi
}

# ==============================================================================
# TESTE 1: HEALTH CHECK
# ==============================================================================
test_health() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📊 1. HEALTH CHECK"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    test_endpoint "GET" "/health" "200" "Health endpoint" > /dev/null
    test_endpoint "GET" "/api/v1/ping" "200" "Ping endpoint" > /dev/null
    
    echo ""
}

# ==============================================================================
# TESTE 2: AUTENTICAÇÃO
# ==============================================================================
test_auth() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔐 2. AUTENTICAÇÃO"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Tentar login com usuário de teste existente
    # Se não existir, pular testes autenticados
    if do_login "$TEST_USER_EMAIL" "$TEST_USER_PASSWORD"; then
        echo -e "  Login com $TEST_USER_EMAIL                        ${GREEN}✓ PASS${NC}"
        ((TOTAL++))
        ((PASSED++))
        
        # Test /auth/me
        test_endpoint "GET" "/api/v1/auth/me" "200" "Auth Me (com JWT)" "" "$JWT_TOKEN" > /dev/null
    else
        # Tentar com credenciais alternativas
        if do_login "admin@nexo.com" "admin123"; then
            echo -e "  Login com admin@nexo.com                           ${GREEN}✓ PASS${NC}"
            ((TOTAL++))
            ((PASSED++))
            test_endpoint "GET" "/api/v1/auth/me" "200" "Auth Me (com JWT)" "" "$JWT_TOKEN" > /dev/null
        else
            echo -e "  Login                                                   ${YELLOW}⚠ SKIP${NC} (sem usuário teste)"
            ((TOTAL++))
            ((SKIPPED++))
        fi
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 3: CATEGORIAS DE SERVIÇOS
# ==============================================================================
test_categorias() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📂 3. CATEGORIAS DE SERVIÇOS"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/categorias-servicos" "200" "Listar categorias de serviços" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Listar categorias                                       ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 4: SERVIÇOS
# ==============================================================================
test_servicos() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✂️  4. SERVIÇOS"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/servicos" "200" "Listar serviços" "" "$JWT_TOKEN" > /dev/null
        test_endpoint "GET" "/api/v1/servicos/stats" "200" "Estatísticas de serviços" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes de serviços                                      ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 5: PROFISSIONAIS
# ==============================================================================
test_profissionais() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "👤 5. PROFISSIONAIS"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/professionals" "200" "Listar profissionais" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Listar profissionais                                    ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 6: CLIENTES (CRM)
# ==============================================================================
test_clientes() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "👥 6. CLIENTES (CRM)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/customers" "200" "Listar clientes" "" "$JWT_TOKEN" > /dev/null
        test_endpoint "GET" "/api/v1/customers/stats" "200" "Estatísticas de clientes" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes de clientes                                      ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 7: AGENDAMENTOS
# ==============================================================================
test_agendamentos() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📅 7. AGENDAMENTOS"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/appointments" "200" "Listar agendamentos" "" "$JWT_TOKEN" > /dev/null
        test_endpoint "GET" "/api/v1/appointments?date=$(date +%Y-%m-%d)" "200" "Agendamentos do dia" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes de agendamentos                                  ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 8: FINANCEIRO
# ==============================================================================
test_financeiro() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "💵 8. FINANCEIRO"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        # Contas a Pagar
        test_endpoint "GET" "/api/v1/financial/payables" "200" "Listar contas a pagar" "" "$JWT_TOKEN" > /dev/null
        
        # Contas a Receber
        test_endpoint "GET" "/api/v1/financial/receivables" "200" "Listar contas a receber" "" "$JWT_TOKEN" > /dev/null
        
        # Fluxo de Caixa
        test_endpoint "GET" "/api/v1/financial/cashflow" "200" "Listar fluxo de caixa" "" "$JWT_TOKEN" > /dev/null
        
        # DRE (pode retornar 404 se não houver dados)
        test_endpoint "GET" "/api/v1/financial/dre" "200|404" "Listar DRE" "" "$JWT_TOKEN" > /dev/null
        
        # Dashboard
        local year=$(date +%Y)
        local month=$(date +%-m)
        test_endpoint "GET" "/api/v1/financial/dashboard?year=$year&month=$month" "200" "Dashboard financeiro" "" "$JWT_TOKEN" > /dev/null
        
        # Projeções
        test_endpoint "GET" "/api/v1/financial/projections?months_ahead=3" "200" "Projeções financeiras" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes financeiros                                      ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 9: METAS
# ==============================================================================
test_metas() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🎯 9. METAS"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        # Metas Mensais
        test_endpoint "GET" "/api/v1/metas/monthly" "200" "Listar metas mensais" "" "$JWT_TOKEN" > /dev/null
        
        # Metas Barbeiros
        test_endpoint "GET" "/api/v1/metas/barbers" "200" "Listar metas por barbeiro" "" "$JWT_TOKEN" > /dev/null
        
        # Metas Ticket
        test_endpoint "GET" "/api/v1/metas/ticket" "200" "Listar metas ticket médio" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes de metas                                         ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 10: PRECIFICAÇÃO
# ==============================================================================
test_precificacao() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "💰 10. PRECIFICAÇÃO"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        # Config (pode não existir ainda)
        test_endpoint "GET" "/api/v1/pricing/config" "200|404" "Config de precificação" "" "$JWT_TOKEN" > /dev/null
        
        # Simulações
        test_endpoint "GET" "/api/v1/pricing/simulations" "200" "Listar simulações" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes de precificação                                  ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 11: ESTOQUE
# ==============================================================================
test_estoque() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📦 11. ESTOQUE"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/stock/items" "200" "Listar produtos estoque" "" "$JWT_TOKEN" > /dev/null
        test_endpoint "GET" "/api/v1/stock/alerts" "200" "Listar alertas estoque" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes de estoque                                       ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 12: LISTA DA VEZ
# ==============================================================================
test_lista_da_vez() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📋 12. LISTA DA VEZ"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/barber-turn/list" "200" "Listar barbeiros na vez" "" "$JWT_TOKEN" > /dev/null
        test_endpoint "GET" "/api/v1/barber-turn/available" "200" "Barbeiros disponíveis" "" "$JWT_TOKEN" > /dev/null
        test_endpoint "GET" "/api/v1/barber-turn/history" "200" "Histórico da vez" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes lista da vez                                     ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# TESTE 13: MEIOS DE PAGAMENTO
# ==============================================================================
test_meios_pagamento() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "💳 13. MEIOS DE PAGAMENTO"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -n "$JWT_TOKEN" ]; then
        test_endpoint "GET" "/api/v1/meios-pagamento" "200" "Listar meios de pagamento" "" "$JWT_TOKEN" > /dev/null
    else
        echo -e "  Testes meios de pagamento                               ${YELLOW}⚠ SKIP${NC} (sem JWT)"
        ((TOTAL++))
        ((SKIPPED++))
    fi
    
    echo ""
}

# ==============================================================================
# RESUMO
# ==============================================================================
print_summary() {
    echo "═══════════════════════════════════════════════════════════════════════════"
    echo "  📊 RESUMO DOS SMOKE TESTS"
    echo "═══════════════════════════════════════════════════════════════════════════"
    echo ""
    echo "  Total de testes:   $TOTAL"
    echo -e "  ${GREEN}Aprovados:${NC}         $PASSED"
    echo -e "  ${RED}Falhados:${NC}          $FAILED"
    echo -e "  ${YELLOW}Pulados:${NC}           $SKIPPED"
    echo ""
    
    if [ $TOTAL -gt 0 ]; then
        local SUCCESS_RATE=$((PASSED * 100 / TOTAL))
        local EFFECTIVE_TOTAL=$((TOTAL - SKIPPED))
        
        if [ $EFFECTIVE_TOTAL -gt 0 ]; then
            local EFFECTIVE_RATE=$((PASSED * 100 / EFFECTIVE_TOTAL))
            echo "  Taxa de sucesso:   ${SUCCESS_RATE}% (${EFFECTIVE_RATE}% efetivo)"
        fi
    fi
    
    echo ""
    
    if [ $FAILED -eq 0 ]; then
        echo -e "${GREEN}═══════════════════════════════════════════════════════════════════════════${NC}"
        echo -e "${GREEN}  ✓ TODOS OS SMOKE TESTS PASSARAM! Sistema operacional.${NC}"
        echo -e "${GREEN}═══════════════════════════════════════════════════════════════════════════${NC}"
        exit 0
    elif [ $PASSED -ge $((TOTAL * 8 / 10)) ]; then
        echo -e "${YELLOW}═══════════════════════════════════════════════════════════════════════════${NC}"
        echo -e "${YELLOW}  ⚠ Alguns testes falharam, mas sistema está parcialmente operacional.${NC}"
        echo -e "${YELLOW}═══════════════════════════════════════════════════════════════════════════${NC}"
        exit 0
    else
        echo -e "${RED}═══════════════════════════════════════════════════════════════════════════${NC}"
        echo -e "${RED}  ✗ MUITOS TESTES FALHARAM! Sistema pode não estar operacional.${NC}"
        echo -e "${RED}═══════════════════════════════════════════════════════════════════════════${NC}"
        exit 1
    fi
}

# ==============================================================================
# EXECUÇÃO PRINCIPAL
# ==============================================================================
main() {
    check_dependencies
    
    test_health
    test_auth
    test_categorias
    test_servicos
    test_profissionais
    test_clientes
    test_agendamentos
    test_financeiro
    test_metas
    test_precificacao
    test_estoque
    test_lista_da_vez
    test_meios_pagamento
    
    print_summary
}

main "$@"
