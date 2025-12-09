#!/bin/bash

# ============================================================================
# 🧪 Barber Analytics Pro V2 — Test API
# Testa endpoints críticos da API com autenticação JWT
# ============================================================================

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

BASE_URL="http://localhost:8080/api/v1"

# Credenciais de teste
TEST_EMAIL="andrey@tratodebarbados.com"
TEST_PASSWORD="@Aa30019258"

# Variáveis globais
ACCESS_TOKEN=""
TENANT_ID=""
TESTS_PASSED=0
TESTS_FAILED=0

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  🧪 Testando API — Barber Analytics Pro V2"
echo -e "${BLUE}║${NC}  📧 Usuário: $TEST_EMAIL"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

# ============================================================================
# TEST PING (Público)
# ============================================================================

echo -e "${YELLOW}📡 Test 1: PING (público)${NC}"
response=$(curl -s "$BASE_URL/ping")
if [[ $response == *"pong"* ]]; then
    echo -e "   ${GREEN}✅ Backend respondendo${NC}"
    ((TESTS_PASSED++))
else
    echo -e "   ${RED}❌ Backend não respondeu${NC}"
    ((TESTS_FAILED++))
    echo -e "   ${RED}⚠️  Backend precisa estar rodando para continuar${NC}"
    exit 1
fi

echo

# ============================================================================
# LOGIN - Obter JWT Token
# ============================================================================

echo -e "${YELLOW}🔐 Test 2: LOGIN${NC}"
login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"$TEST_PASSWORD\"}")

# Extrair access_token do JSON
ACCESS_TOKEN=$(echo "$login_response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    # Extrair tenant_id do response
    TENANT_ID=$(echo "$login_response" | grep -o '"tenant_id":"[^"]*"' | cut -d'"' -f4)
    USER_NAME=$(echo "$login_response" | grep -o '"nome":"[^"]*"' | cut -d'"' -f4)
    USER_ROLE=$(echo "$login_response" | grep -o '"role":"[^"]*"' | cut -d'"' -f4)
    
    echo -e "   ${GREEN}✅ Login bem-sucedido${NC}"
    echo -e "   ${CYAN}👤 Usuário: $USER_NAME${NC}"
    echo -e "   ${CYAN}🏢 Tenant: $TENANT_ID${NC}"
    echo -e "   ${CYAN}🎭 Role: $USER_ROLE${NC}"
    echo -e "   ${CYAN}🔑 Token: ${ACCESS_TOKEN:0:50}...${NC}"
    ((TESTS_PASSED++))
else
    echo -e "   ${RED}❌ Login falhou${NC}"
    echo -e "   ${RED}Resposta: $login_response${NC}"
    ((TESTS_FAILED++))
    echo -e "   ${RED}⚠️  Sem token, testes autenticados serão ignorados${NC}"
fi

echo

# ============================================================================
# TEST AUTH/ME (Autenticado)
# ============================================================================

echo -e "${YELLOW}👤 Test 3: AUTH/ME${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    me_response=$(curl -s -X GET "$BASE_URL/auth/me" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    
    if [[ $me_response == *"email"* ]] && [[ $me_response == *"$TEST_EMAIL"* ]]; then
        echo -e "   ${GREEN}✅ Endpoint /auth/me OK${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Endpoint /auth/me falhou${NC}"
        echo -e "   ${RED}Resposta: $me_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST HEALTH (Autenticado)
# ============================================================================

echo -e "${YELLOW}💚 Test 4: HEALTH CHECK${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    if [ "$http_code" -eq 200 ]; then
        echo -e "   ${GREEN}✅ Health check OK${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Health check falhou (HTTP $http_code)${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST RECEITAS (Autenticado)
# ============================================================================

echo -e "${YELLOW}💰 Test 5: LIST RECEITAS${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    receitas_response=$(curl -s -X GET "$BASE_URL/receitas" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/receitas" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "   ${GREEN}✅ Receitas endpoint OK (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Receitas falhou (HTTP $http_code)${NC}"
        echo -e "   ${RED}Resposta: $receitas_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST MEIOS DE PAGAMENTO (Autenticado)
# ============================================================================

echo -e "${YELLOW}💳 Test 6: LIST MEIOS DE PAGAMENTO${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    meios_response=$(curl -s -X GET "$BASE_URL/meios-pagamento" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/meios-pagamento" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if [ "$http_code" -eq 200 ]; then
        count=$(echo "$meios_response" | grep -o '"total":[0-9]*' | cut -d':' -f2)
        echo -e "   ${GREEN}✅ Meios de Pagamento OK (HTTP $http_code)${NC}"
        echo -e "   ${CYAN}📊 Total de registros: ${count:-0}${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Meios de Pagamento falhou (HTTP $http_code)${NC}"
        echo -e "   ${RED}Resposta: $meios_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST AGENDAMENTOS (Autenticado)
# ============================================================================

echo -e "${YELLOW}📅 Test 7: LIST AGENDAMENTOS${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    today=$(date +%Y-%m-%d)
    agendamentos_response=$(curl -s -X GET "$BASE_URL/agendamentos?data=$today" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/agendamentos?data=$today" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "   ${GREEN}✅ Agendamentos endpoint OK (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Agendamentos falhou (HTTP $http_code)${NC}"
        echo -e "   ${RED}Resposta: $agendamentos_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST METAS (Autenticado)
# ============================================================================

echo -e "${YELLOW}🎯 Test 8: LIST METAS${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    current_month=$(date +%Y-%m)
    metas_response=$(curl -s -X GET "$BASE_URL/metas?mes=$current_month" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/metas?mes=$current_month" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "   ${GREEN}✅ Metas endpoint OK (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Metas falhou (HTTP $http_code)${NC}"
        echo -e "   ${RED}Resposta: $metas_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST CLIENTES (Autenticado)
# ============================================================================

echo -e "${YELLOW}👥 Test 9: LIST CLIENTES${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    clientes_response=$(curl -s -X GET "$BASE_URL/clientes" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/clientes" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "   ${GREEN}✅ Clientes endpoint OK (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Clientes falhou (HTTP $http_code)${NC}"
        echo -e "   ${RED}Resposta: $clientes_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST PROFISSIONAIS (Autenticado)
# ============================================================================

echo -e "${YELLOW}✂️  Test 10: LIST PROFISSIONAIS${NC}"
if [ -n "$ACCESS_TOKEN" ]; then
    profissionais_response=$(curl -s -X GET "$BASE_URL/profissionais" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json")
    http_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/profissionais" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "   ${GREEN}✅ Profissionais endpoint OK (HTTP $http_code)${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "   ${RED}❌ Profissionais falhou (HTTP $http_code)${NC}"
        echo -e "   ${RED}Resposta: $profissionais_response${NC}"
        ((TESTS_FAILED++))
    fi
else
    echo -e "   ${YELLOW}⏭️  Ignorado (sem token)${NC}"
fi

echo

# ============================================================================
# TEST FRONTEND
# ============================================================================

echo -e "${YELLOW}🌐 Test 11: FRONTEND${NC}"
http_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:3000")
if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 307 ]; then
    echo -e "   ${GREEN}✅ Frontend respondendo (HTTP $http_code)${NC}"
    ((TESTS_PASSED++))
else
    echo -e "   ${RED}❌ Frontend não respondeu (HTTP $http_code)${NC}"
    ((TESTS_FAILED++))
fi

echo

# ============================================================================
# SUMMARY
# ============================================================================

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))

echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ TODOS OS TESTES PASSARAM! ($TESTS_PASSED/$TOTAL_TESTS)${NC}"
else
    echo -e "${YELLOW}⚠️  TESTES COMPLETOS: $TESTS_PASSED passed, $TESTS_FAILED failed${NC}"
fi
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo

echo -e "🌐 Frontend: ${BLUE}http://localhost:3000${NC}"
echo -e "📡 API: ${BLUE}http://localhost:8080/api/v1${NC}"
echo -e "🔐 Token válido: ${CYAN}$([ -n \"$ACCESS_TOKEN\" ] && echo 'Sim' || echo 'Não')${NC}"
echo

exit $TESTS_FAILED
