#!/bin/bash
# ============================================================================
# NEXO MVP v1.0.0 - Teste de APIs REST
# ============================================================================
# Testa todas as APIs listadas no TAREFAS_MVP_V1.0.0.md
# Uso: ./scripts/test-mvp-apis.sh
# ============================================================================

set -e

BASE_URL="http://localhost:8080/api/v1"
TOKEN=$(cat /tmp/nexo_token.txt 2>/dev/null || echo "")
UNIT_ID=$(cat /tmp/nexo_unit.txt 2>/dev/null || echo "")

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Contadores
PASSED=0
FAILED=0
TOTAL=0

# Função para testar endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_codes=$3  # Códigos esperados separados por vírgula (ex: "200,201")
    local data=$4
    local description=$5
    
    TOTAL=$((TOTAL + 1))
    
    local headers="-H 'Authorization: Bearer $TOKEN'"
    if [ -n "$UNIT_ID" ]; then
        headers="$headers -H 'X-Unit-ID: $UNIT_ID'"
    fi
    headers="$headers -H 'Content-Type: application/json'"
    
    local cmd="curl -s -o /tmp/api_response.json -w '%{http_code}' -X $method"
    
    if [ -n "$data" ]; then
        cmd="$cmd -d '$data'"
    fi
    
    cmd="$cmd $headers '$BASE_URL$endpoint'"
    
    local http_code=$(eval $cmd)
    local response=$(cat /tmp/api_response.json 2>/dev/null | head -c 200)
    
    # Verificar se o código está entre os esperados
    local is_expected=false
    IFS=',' read -ra EXPECTED_ARRAY <<< "$expected_codes"
    for expected in "${EXPECTED_ARRAY[@]}"; do
        if [ "$http_code" = "$expected" ]; then
            is_expected=true
            break
        fi
    done
    
    if [ "$is_expected" = true ]; then
        echo -e "${GREEN}✅ PASS${NC} [$method] $endpoint → HTTP $http_code"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}❌ FAIL${NC} [$method] $endpoint → HTTP $http_code (esperado: $expected_codes)"
        echo "   Response: $response"
        FAILED=$((FAILED + 1))
    fi
}

# Verificar se temos token
if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Token não encontrado. Execute o login primeiro.${NC}"
    exit 1
fi

echo ""
echo "============================================================================"
echo "🧪 NEXO MVP v1.0.0 - Teste de APIs REST"
echo "============================================================================"
echo "Token: ${TOKEN:0:20}..."
echo "Unit ID: $UNIT_ID"
echo "============================================================================"
echo ""

# ============================================================================
# MÓDULO 1: AUTENTICAÇÃO
# ============================================================================
echo -e "${YELLOW}🔐 MÓDULO 1: AUTENTICAÇÃO${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/auth/me" "200" "" "Dados do usuário logado"
echo ""

# ============================================================================
# MÓDULO 2: CATEGORIAS DE SERVIÇOS
# ============================================================================
echo -e "${YELLOW}🏷️  MÓDULO 2: CATEGORIAS DE SERVIÇOS${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/categorias-servicos" "200" "" "Listar categorias"
echo ""

# ============================================================================
# MÓDULO 3: SERVIÇOS
# ============================================================================
echo -e "${YELLOW}✂️  MÓDULO 3: SERVIÇOS${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/servicos" "200" "" "Listar serviços"
test_endpoint "GET" "/servicos/stats" "200" "" "Estatísticas de serviços"
echo ""

# ============================================================================
# MÓDULO 4: PROFISSIONAIS
# ============================================================================
echo -e "${YELLOW}👨‍💼 MÓDULO 4: PROFISSIONAIS${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/professionals" "200" "" "Listar profissionais"
echo ""

# ============================================================================
# MÓDULO 5: LISTA DA VEZ (Barber Turn)
# ============================================================================
echo -e "${YELLOW}🔄 MÓDULO 5: LISTA DA VEZ${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/barber-turn/list" "200" "" "Listar barbeiros na vez"
test_endpoint "GET" "/barber-turn/available" "200" "" "Barbeiros disponíveis"
test_endpoint "GET" "/barber-turn/history" "200" "" "Histórico de turnos"
test_endpoint "GET" "/barber-turn/history/summary" "200" "" "Resumo do histórico"
echo ""

# ============================================================================
# MÓDULO 6: CLIENTES (CRM)
# ============================================================================
echo -e "${YELLOW}👥 MÓDULO 6: CLIENTES (CRM)${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/customers" "200" "" "Listar clientes"
test_endpoint "GET" "/customers/stats" "200" "" "Estatísticas de clientes"
test_endpoint "GET" "/customers/search?q=test" "200" "" "Buscar clientes"
echo ""

# ============================================================================
# MÓDULO 7: AGENDAMENTO
# ============================================================================
echo -e "${YELLOW}📅 MÓDULO 7: AGENDAMENTO${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/appointments" "200" "" "Listar agendamentos"
test_endpoint "GET" "/blocked-times" "200" "" "Listar horários bloqueados"
echo ""

# ============================================================================
# MÓDULO 8: ESTOQUE
# ============================================================================
echo -e "${YELLOW}📦 MÓDULO 8: ESTOQUE${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/stock/items" "200" "" "Listar produtos do estoque"
test_endpoint "GET" "/stock/alerts" "200" "" "Alertas de estoque baixo"
test_endpoint "GET" "/fornecedores" "200" "" "Listar fornecedores"
test_endpoint "GET" "/categorias-produtos" "200" "" "Listar categorias de produtos"
echo ""

# ============================================================================
# MÓDULO 9: FINANCEIRO
# ============================================================================
echo -e "${YELLOW}💰 MÓDULO 9: FINANCEIRO${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/financial/payables" "200" "" "Contas a pagar"
test_endpoint "GET" "/financial/receivables" "200" "" "Contas a receber"
test_endpoint "GET" "/financial/compensations" "200" "" "Compensações"
test_endpoint "GET" "/financial/cashflow" "200" "" "Fluxo de caixa"
test_endpoint "GET" "/financial/dre" "200" "" "DRE"
test_endpoint "GET" "/financial/dashboard" "200" "" "Dashboard financeiro"
test_endpoint "GET" "/financial/projections" "200" "" "Projeções financeiras"
test_endpoint "GET" "/financial/despesas-fixas" "200" "" "Despesas fixas"
echo ""

# ============================================================================
# MÓDULO 10: METAS
# ============================================================================
echo -e "${YELLOW}🎯 MÓDULO 10: METAS${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/metas/monthly" "200" "" "Metas mensais"
test_endpoint "GET" "/metas/barbers" "200" "" "Metas por barbeiro"
test_endpoint "GET" "/metas/ticket" "200" "" "Metas de ticket médio"
echo ""

# ============================================================================
# MÓDULO 11: COMISSÕES
# ============================================================================
echo -e "${YELLOW}💵 MÓDULO 11: COMISSÕES${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/commissions/rules" "200" "" "Regras de comissão"
test_endpoint "GET" "/commissions/periods" "200" "" "Períodos de comissão"
test_endpoint "GET" "/commissions/items" "200" "" "Itens de comissão"
test_endpoint "GET" "/commissions/advances" "200" "" "Adiantamentos"
test_endpoint "GET" "/commissions/pending" "200" "" "Comissões pendentes"
echo ""

# ============================================================================
# MÓDULO 12: COMANDAS
# ============================================================================
echo -e "${YELLOW}🧾 MÓDULO 12: COMANDAS${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/commands" "200" "" "Listar comandas"
echo ""

# ============================================================================
# MÓDULO 13: CAIXA DIÁRIO
# ============================================================================
echo -e "${YELLOW}🏧 MÓDULO 13: CAIXA DIÁRIO${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/caixa/aberto" "200,404" "" "Caixa aberto"
test_endpoint "GET" "/caixa/historico" "200" "" "Histórico de caixa"
echo ""

# ============================================================================
# MÓDULO 14: MEIOS DE PAGAMENTO
# ============================================================================
echo -e "${YELLOW}💳 MÓDULO 14: MEIOS DE PAGAMENTO${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/meios-pagamento" "200" "" "Listar meios de pagamento"
echo ""

# ============================================================================
# MÓDULO 15: PRECIFICAÇÃO
# ============================================================================
echo -e "${YELLOW}💲 MÓDULO 15: PRECIFICAÇÃO${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/pricing/config" "200,404" "" "Configuração de preços"
test_endpoint "GET" "/pricing/simulations" "200" "" "Simulações de preços"
echo ""

# ============================================================================
# MÓDULO 16: UNIDADES
# ============================================================================
echo -e "${YELLOW}🏢 MÓDULO 16: UNIDADES${NC}"
echo "----------------------------------------------------------------------------"
test_endpoint "GET" "/units/me" "200" "" "Minhas unidades"
echo ""

# ============================================================================
# RELATÓRIO FINAL
# ============================================================================
echo ""
echo "============================================================================"
echo "📊 RELATÓRIO FINAL"
echo "============================================================================"
echo -e "Total de testes: ${TOTAL}"
echo -e "Passou: ${GREEN}${PASSED}${NC}"
echo -e "Falhou: ${RED}${FAILED}${NC}"
echo ""

PERCENT=$((PASSED * 100 / TOTAL))
if [ $PERCENT -ge 90 ]; then
    echo -e "${GREEN}✅ APIs MVP: ${PERCENT}% funcionando${NC}"
elif [ $PERCENT -ge 70 ]; then
    echo -e "${YELLOW}⚠️  APIs MVP: ${PERCENT}% funcionando${NC}"
else
    echo -e "${RED}❌ APIs MVP: ${PERCENT}% funcionando${NC}"
fi

echo "============================================================================"
echo ""

# Exit com código de erro se houver falhas
if [ $FAILED -gt 0 ]; then
    exit 1
fi
