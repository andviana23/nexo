#!/usr/bin/env bash

################################################################################
# smoke_tests_subscription.sh - Smoke Tests para Módulo de Assinaturas
#
# Descrição:
#   Executa testes de ponta-a-ponta para os endpoints de Planos e Assinaturas.
#   Valida operações CRUD, validações de negócio, e integração Asaas.
#
# Uso:
#   ./scripts/smoke_tests_subscription.sh [API_URL]
#
# Exemplo:
#   ./scripts/smoke_tests_subscription.sh http://localhost:8080
#
# Fluxo Testado:
#   PLANOS:
#   1. Criar plano válido
#   2. Criar plano com valor inválido (espera erro)
#   3. Listar planos (verifica plano criado)
#   4. Atualizar plano
#   5. Desativar plano
#
#   ASSINATURAS:
#   6. Criar assinatura PIX
#   7. Listar assinaturas
#   8. Buscar assinatura por ID
#   9. Renovar assinatura
#   10. Cancelar assinatura
#   11. Verificar métricas
#
# Requisitos:
#   - curl
#   - jq
#   - Backend rodando com JWT válido
#
# Autor: Copilot/Codex
# Versão: 1.0.0
################################################################################

set -euo pipefail

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configurações
API_URL="${1:-http://localhost:8080}"
TIMEOUT=15

# Contadores
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Variáveis globais
TENANT_ID=""
JWT_TOKEN=""
PLAN_ID=""
SUBSCRIPTION_ID=""
CUSTOMER_ID=""

# Timestamp único
TIMESTAMP=$(date +%s)
UNIQUE_SUFFIX="test_${TIMESTAMP}"

# Função de log
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_test() {
    echo -e "${CYAN}[TEST]${NC} $1"
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    TESTS_PASSED=$((TESTS_PASSED + 1))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    TESTS_FAILED=$((TESTS_FAILED + 1))
}

log_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
}

# Função para fazer requisições HTTP
http_request() {
    local method="$1"
    local endpoint="$2"
    local data="${3:-}"
    local full_url="${API_URL}${endpoint}"

    if [ -n "$data" ]; then
        curl -s -w "\n%{http_code}" -X "$method" "$full_url" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "$data" \
            --max-time "$TIMEOUT" 2>/dev/null || echo -e "\n000"
    else
        curl -s -w "\n%{http_code}" -X "$method" "$full_url" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "X-Tenant-ID: $TENANT_ID" \
            --max-time "$TIMEOUT" 2>/dev/null || echo -e "\n000"
    fi
}

# Função para extrair código HTTP da resposta
get_http_code() {
    echo "$1" | tail -n1
}

# Função para extrair body da resposta
get_response_body() {
    echo "$1" | sed '$d'
}

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🧪 Barber Analytics Pro - Smoke Tests: Assinaturas"
echo "═══════════════════════════════════════════════════════════════"
echo ""
log_info "API URL: $API_URL"
echo ""

# Verificar dependências
log_info "Verificando dependências..."
if ! command -v curl &> /dev/null; then
    log_error "curl não está instalado"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    log_warning "jq não está instalado (parsing JSON será limitado)"
    JQ_AVAILABLE=false
else
    JQ_AVAILABLE=true
fi

log_success "Dependências verificadas"
echo ""

# ============================================================================
# SETUP: Health Check e Autenticação
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "0. Setup: Health Check"
echo "───────────────────────────────────────────────────────────────"

RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_URL}/health" --max-time "$TIMEOUT" 2>/dev/null || echo -e "\n000")
HTTP_CODE=$(get_http_code "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    log_success "Health check passou (HTTP 200)"
else
    log_error "Health check falhou (HTTP $HTTP_CODE)"
    exit 1
fi

# Autenticação - usar credenciais de teste
echo ""
log_info "Autenticando para obter JWT..."

# Credenciais do seed (002_seed_test_user.sql)
# Email: admin@teste.com / Senha: Admin123!
AUTH_PAYLOAD='{"email":"admin@teste.com","password":"Admin123!"}'
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "$AUTH_PAYLOAD" \
    --max-time "$TIMEOUT" 2>/dev/null || echo -e "\n000")

HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        JWT_TOKEN=$(echo "$BODY" | jq -r '.access_token // .token // empty')
        TENANT_ID=$(echo "$BODY" | jq -r '.user.tenant_id // .tenant_id // empty')
    else
        # Fallback: extração simples
        JWT_TOKEN=$(echo "$BODY" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        TENANT_ID=$(echo "$BODY" | grep -o '"tenant_id":"[^"]*"' | cut -d'"' -f4)
    fi
    
    # Se tenant_id estiver vazio, usar o valor padrão de teste
    if [ -z "$TENANT_ID" ]; then
        TENANT_ID="00000000-0000-0000-0000-000000000001"
    fi
    
    if [ -n "$JWT_TOKEN" ]; then
        log_success "Login bem-sucedido, JWT obtido"
        log_info "Tenant ID: $TENANT_ID"
    else
        log_error "JWT não encontrado na resposta"
        log_info "Response: $BODY"
        exit 1
    fi
else
    log_error "Login falhou (HTTP $HTTP_CODE)"
    log_info "Response: $BODY"
    exit 1
fi

echo ""

# ============================================================================
# TESTE 1: Criar Plano Válido
# ============================================================================
echo "═══════════════════════════════════════════════════════════════"
echo "  📋 TESTES DE PLANOS"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "───────────────────────────────────────────────────────────────"
log_test "1. Criar Plano Válido"
echo "───────────────────────────────────────────────────────────────"

PLAN_NAME="Plano Teste ${UNIQUE_SUFFIX}"
CREATE_PLAN_PAYLOAD=$(cat <<EOF
{
    "nome": "${PLAN_NAME}",
    "descricao": "Plano de teste criado pelo smoke test",
    "valor": "99.90",
    "qtd_servicos": 10,
    "limite_uso_mensal": 30
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/plans" "$CREATE_PLAN_PAYLOAD")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 201 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        PLAN_ID=$(echo "$BODY" | jq -r '.id // empty')
        PLAN_VALOR=$(echo "$BODY" | jq -r '.valor // empty')
    else
        PLAN_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        PLAN_VALOR=$(echo "$BODY" | grep -o '"valor":"[^"]*"' | cut -d'"' -f4)
    fi
    
    if [ -n "$PLAN_ID" ]; then
        log_success "Plano criado com sucesso (ID: $PLAN_ID)"
        log_info "Valor: R$ $PLAN_VALOR"
    else
        log_error "ID do plano não encontrado na resposta"
        log_info "Response: $BODY"
    fi
else
    log_error "Falha ao criar plano (HTTP $HTTP_CODE)"
    log_info "Response: $BODY"
fi

echo ""

# ============================================================================
# TESTE 2: Criar Plano com Valor Inválido (Espera Erro)
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "2. Criar Plano com Valor Inválido (espera 400)"
echo "───────────────────────────────────────────────────────────────"

INVALID_PLAN_PAYLOAD=$(cat <<EOF
{
    "nome": "Plano Inválido",
    "descricao": "Teste com valor negativo",
    "valor": "-50.00",
    "qtd_servicos": 5
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/plans" "$INVALID_PLAN_PAYLOAD")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 400 ]; then
    log_success "Validação funcionou corretamente (HTTP 400 para valor inválido)"
else
    log_error "Deveria retornar 400, mas retornou HTTP $HTTP_CODE"
    log_info "Response: $BODY"
fi

echo ""

# ============================================================================
# TESTE 3: Listar Planos
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "3. Listar Planos"
echo "───────────────────────────────────────────────────────────────"

RESPONSE=$(http_request "GET" "/api/v1/plans")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        PLAN_COUNT=$(echo "$BODY" | jq 'if type == "array" then length else 0 end')
        FOUND_PLAN=$(echo "$BODY" | jq -r --arg id "$PLAN_ID" '.[] | select(.id == $id) | .nome // empty')
    else
        PLAN_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l)
        FOUND_PLAN=$(echo "$BODY" | grep -o "$PLAN_ID" || echo "")
    fi
    
    log_success "Listagem de planos OK (HTTP 200)"
    log_info "Total de planos: $PLAN_COUNT"
    
    if [ -n "$FOUND_PLAN" ]; then
        log_success "Plano criado encontrado na lista"
    else
        log_warning "Plano criado não encontrado na lista"
    fi
else
    log_error "Falha ao listar planos (HTTP $HTTP_CODE)"
    log_info "Response: $BODY"
fi

echo ""

# ============================================================================
# TESTE 4: Atualizar Plano
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "4. Atualizar Plano"
echo "───────────────────────────────────────────────────────────────"

if [ -n "$PLAN_ID" ]; then
    UPDATE_PLAN_PAYLOAD=$(cat <<EOF
{
    "nome": "${PLAN_NAME} - Atualizado",
    "descricao": "Plano atualizado pelo smoke test",
    "valor": "149.90",
    "qtd_servicos": 15,
    "limite_uso_mensal": 50,
    "ativo": true
}
EOF
)

    RESPONSE=$(http_request "PUT" "/api/v1/plans/${PLAN_ID}" "$UPDATE_PLAN_PAYLOAD")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    BODY=$(get_response_body "$RESPONSE")

    if [ "$HTTP_CODE" -eq 200 ]; then
        if [ "$JQ_AVAILABLE" = true ]; then
            NEW_VALOR=$(echo "$BODY" | jq -r '.valor // empty')
        else
            NEW_VALOR=$(echo "$BODY" | grep -o '"valor":"[^"]*"' | cut -d'"' -f4)
        fi
        log_success "Plano atualizado com sucesso"
        log_info "Novo valor: R$ $NEW_VALOR"
    else
        log_error "Falha ao atualizar plano (HTTP $HTTP_CODE)"
        log_info "Response: $BODY"
    fi
else
    log_warning "Pulando teste - PLAN_ID não disponível"
fi

echo ""

# ============================================================================
# TESTE 5: Desativar Plano
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "5. Desativar Plano"
echo "───────────────────────────────────────────────────────────────"

if [ -n "$PLAN_ID" ]; then
    RESPONSE=$(http_request "DELETE" "/api/v1/plans/${PLAN_ID}")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    BODY=$(get_response_body "$RESPONSE")

    if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 204 ]; then
        log_success "Plano desativado com sucesso"
    else
        log_error "Falha ao desativar plano (HTTP $HTTP_CODE)"
        log_info "Response: $BODY"
    fi
else
    log_warning "Pulando teste - PLAN_ID não disponível"
fi

echo ""

# ============================================================================
# TESTES DE ASSINATURAS
# ============================================================================
echo "═══════════════════════════════════════════════════════════════"
echo "  📄 TESTES DE ASSINATURAS"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Primeiro, criar um plano ativo para usar nas assinaturas
echo "───────────────────────────────────────────────────────────────"
log_info "Setup: Criando plano para assinaturas..."
echo "───────────────────────────────────────────────────────────────"

PLAN_FOR_SUB_NAME="Plano Assinatura ${UNIQUE_SUFFIX}"
CREATE_PLAN_PAYLOAD=$(cat <<EOF
{
    "nome": "${PLAN_FOR_SUB_NAME}",
    "descricao": "Plano para testar assinaturas",
    "valor": "79.90",
    "qtd_servicos": 8
}
EOF
)

RESPONSE=$(http_request "POST" "/api/v1/plans" "$CREATE_PLAN_PAYLOAD")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 201 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        ACTIVE_PLAN_ID=$(echo "$BODY" | jq -r '.id // empty')
    else
        ACTIVE_PLAN_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    fi
    log_success "Plano para assinaturas criado (ID: $ACTIVE_PLAN_ID)"
else
    log_warning "Não foi possível criar plano para assinaturas (HTTP $HTTP_CODE)"
    ACTIVE_PLAN_ID=""
fi

# Buscar um cliente existente para usar nas assinaturas
echo ""
log_info "Buscando cliente existente..."

RESPONSE=$(http_request "GET" "/api/v1/customers?limit=1")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        # Resposta paginada tem formato: {"data": [...], "page": 1, "total": N}
        CUSTOMER_ID=$(echo "$BODY" | jq -r '.data[0].id // .[0].id // empty' 2>/dev/null)
        CUSTOMER_NAME=$(echo "$BODY" | jq -r '.data[0].nome // .[0].nome // "Cliente"' 2>/dev/null)
    else
        CUSTOMER_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        CUSTOMER_NAME="Cliente"
    fi
    
    if [ -n "$CUSTOMER_ID" ] && [ "$CUSTOMER_ID" != "null" ]; then
        log_success "Cliente encontrado: $CUSTOMER_NAME (ID: $CUSTOMER_ID)"
    else
        log_warning "Nenhum cliente encontrado - alguns testes serão pulados"
        CUSTOMER_ID=""
    fi
else
    log_warning "Não foi possível buscar clientes (HTTP $HTTP_CODE)"
    CUSTOMER_ID=""
fi

echo ""

# ============================================================================
# TESTE 6: Criar Assinatura PIX
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "6. Criar Assinatura PIX"
echo "───────────────────────────────────────────────────────────────"

if [ -n "$ACTIVE_PLAN_ID" ] && [ -n "$CUSTOMER_ID" ]; then
    CREATE_SUB_PAYLOAD=$(cat <<EOF
{
    "cliente_id": "${CUSTOMER_ID}",
    "plano_id": "${ACTIVE_PLAN_ID}",
    "forma_pagamento": "PIX"
}
EOF
)

    RESPONSE=$(http_request "POST" "/api/v1/subscriptions" "$CREATE_SUB_PAYLOAD")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    BODY=$(get_response_body "$RESPONSE")

    if [ "$HTTP_CODE" -eq 201 ]; then
        if [ "$JQ_AVAILABLE" = true ]; then
            SUBSCRIPTION_ID=$(echo "$BODY" | jq -r '.id // empty')
            SUB_STATUS=$(echo "$BODY" | jq -r '.status // empty')
            LINK_PAGAMENTO=$(echo "$BODY" | jq -r '.link_pagamento // "N/A"')
        else
            SUBSCRIPTION_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
            SUB_STATUS=$(echo "$BODY" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
            LINK_PAGAMENTO="N/A"
        fi
        log_success "Assinatura criada com sucesso (ID: $SUBSCRIPTION_ID)"
        log_info "Status: $SUB_STATUS"
        log_info "Link Pagamento: $LINK_PAGAMENTO"
    else
        log_error "Falha ao criar assinatura (HTTP $HTTP_CODE)"
        log_info "Response: $BODY"
        SUBSCRIPTION_ID=""
    fi
else
    log_warning "Pulando teste - ACTIVE_PLAN_ID ou CUSTOMER_ID não disponível"
    SUBSCRIPTION_ID=""
fi

echo ""

# ============================================================================
# TESTE 7: Listar Assinaturas
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "7. Listar Assinaturas"
echo "───────────────────────────────────────────────────────────────"

RESPONSE=$(http_request "GET" "/api/v1/subscriptions")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        SUB_COUNT=$(echo "$BODY" | jq 'if type == "array" then length else 0 end')
    else
        SUB_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l)
    fi
    log_success "Listagem de assinaturas OK (HTTP 200)"
    log_info "Total de assinaturas: $SUB_COUNT"
else
    log_error "Falha ao listar assinaturas (HTTP $HTTP_CODE)"
    log_info "Response: $BODY"
fi

echo ""

# ============================================================================
# TESTE 8: Buscar Assinatura por ID
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "8. Buscar Assinatura por ID"
echo "───────────────────────────────────────────────────────────────"

if [ -n "$SUBSCRIPTION_ID" ]; then
    RESPONSE=$(http_request "GET" "/api/v1/subscriptions/${SUBSCRIPTION_ID}")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    BODY=$(get_response_body "$RESPONSE")

    if [ "$HTTP_CODE" -eq 200 ]; then
        if [ "$JQ_AVAILABLE" = true ]; then
            SUB_CLIENTE=$(echo "$BODY" | jq -r '.cliente_nome // empty')
            SUB_PLANO=$(echo "$BODY" | jq -r '.plano_nome // empty')
        else
            SUB_CLIENTE=$(echo "$BODY" | grep -o '"cliente_nome":"[^"]*"' | cut -d'"' -f4)
            SUB_PLANO=$(echo "$BODY" | grep -o '"plano_nome":"[^"]*"' | cut -d'"' -f4)
        fi
        log_success "Assinatura encontrada"
        log_info "Cliente: $SUB_CLIENTE"
        log_info "Plano: $SUB_PLANO"
    else
        log_error "Falha ao buscar assinatura (HTTP $HTTP_CODE)"
        log_info "Response: $BODY"
    fi
else
    log_warning "Pulando teste - SUBSCRIPTION_ID não disponível"
fi

echo ""

# ============================================================================
# TESTE 9: Renovar Assinatura
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "9. Renovar Assinatura"
echo "───────────────────────────────────────────────────────────────"

if [ -n "$SUBSCRIPTION_ID" ]; then
    RENEW_PAYLOAD=$(cat <<EOF
{
    "forma_pagamento": "PIX"
}
EOF
)

    RESPONSE=$(http_request "POST" "/api/v1/subscriptions/${SUBSCRIPTION_ID}/renew" "$RENEW_PAYLOAD")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    BODY=$(get_response_body "$RESPONSE")

    # Renovação pode falhar se assinatura não estiver no estado correto
    if [ "$HTTP_CODE" -eq 200 ]; then
        log_success "Assinatura renovada com sucesso"
    elif [ "$HTTP_CODE" -eq 400 ]; then
        log_warning "Renovação rejeitada (assinatura pode já estar ativa)"
        log_info "Response: $BODY"
    else
        log_error "Falha ao renovar assinatura (HTTP $HTTP_CODE)"
        log_info "Response: $BODY"
    fi
else
    log_warning "Pulando teste - SUBSCRIPTION_ID não disponível"
fi

echo ""

# ============================================================================
# TESTE 10: Cancelar Assinatura
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "10. Cancelar Assinatura"
echo "───────────────────────────────────────────────────────────────"

if [ -n "$SUBSCRIPTION_ID" ]; then
    RESPONSE=$(http_request "DELETE" "/api/v1/subscriptions/${SUBSCRIPTION_ID}")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    BODY=$(get_response_body "$RESPONSE")

    if [ "$HTTP_CODE" -eq 200 ]; then
        log_success "Assinatura cancelada com sucesso"
    elif [ "$HTTP_CODE" -eq 409 ]; then
        log_warning "Assinatura já estava cancelada"
    else
        log_error "Falha ao cancelar assinatura (HTTP $HTTP_CODE)"
        log_info "Response: $BODY"
    fi
else
    log_warning "Pulando teste - SUBSCRIPTION_ID não disponível"
fi

echo ""

# ============================================================================
# TESTE 11: Verificar Métricas
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_test "11. Verificar Métricas de Assinaturas"
echo "───────────────────────────────────────────────────────────────"

RESPONSE=$(http_request "GET" "/api/v1/subscriptions/metrics")
HTTP_CODE=$(get_http_code "$RESPONSE")
BODY=$(get_response_body "$RESPONSE")

if [ "$HTTP_CODE" -eq 200 ]; then
    if [ "$JQ_AVAILABLE" = true ]; then
        TOTAL_ATIVAS=$(echo "$BODY" | jq -r '.total_ativas // 0')
        TOTAL_INATIVAS=$(echo "$BODY" | jq -r '.total_inativas // 0')
        TOTAL_INADIMPLENTES=$(echo "$BODY" | jq -r '.total_inadimplentes // 0')
        RECEITA_MENSAL=$(echo "$BODY" | jq -r '.receita_mensal // "0.00"')
    else
        TOTAL_ATIVAS=$(echo "$BODY" | grep -o '"total_ativas":[0-9]*' | cut -d':' -f2)
        RECEITA_MENSAL=$(echo "$BODY" | grep -o '"receita_mensal":"[^"]*"' | cut -d'"' -f4)
    fi
    log_success "Métricas obtidas com sucesso"
    log_info "Assinaturas Ativas: $TOTAL_ATIVAS"
    log_info "Assinaturas Inativas: ${TOTAL_INATIVAS:-N/A}"
    log_info "Inadimplentes: ${TOTAL_INADIMPLENTES:-N/A}"
    log_info "Receita Mensal: R$ $RECEITA_MENSAL"
else
    log_error "Falha ao obter métricas (HTTP $HTTP_CODE)"
    log_info "Response: $BODY"
fi

echo ""

# ============================================================================
# CLEANUP
# ============================================================================
echo "───────────────────────────────────────────────────────────────"
log_info "Cleanup: Desativando plano criado para testes..."
echo "───────────────────────────────────────────────────────────────"

if [ -n "$ACTIVE_PLAN_ID" ]; then
    RESPONSE=$(http_request "DELETE" "/api/v1/plans/${ACTIVE_PLAN_ID}")
    HTTP_CODE=$(get_http_code "$RESPONSE")
    
    if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 204 ]; then
        log_success "Plano de teste desativado"
    else
        log_warning "Não foi possível desativar plano de teste"
    fi
fi

echo ""

# ============================================================================
# RESUMO
# ============================================================================
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 RESUMO DOS TESTES"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "Total de testes:   ${CYAN}${TESTS_TOTAL}${NC}"
echo -e "Testes passaram:   ${GREEN}${TESTS_PASSED}${NC}"
echo -e "Testes falharam:   ${RED}${TESTS_FAILED}${NC}"
echo ""

# Cálculo de taxa de sucesso
if [ "$TESTS_TOTAL" -gt 0 ]; then
    SUCCESS_RATE=$((TESTS_PASSED * 100 / TESTS_TOTAL))
    echo -e "Taxa de sucesso:   ${CYAN}${SUCCESS_RATE}%${NC}"
else
    SUCCESS_RATE=0
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"

# Exit code
if [ "$TESTS_FAILED" -gt 0 ]; then
    echo -e "${RED}❌ ALGUNS TESTES FALHARAM${NC}"
    exit 1
else
    echo -e "${GREEN}✅ TODOS OS TESTES PASSARAM${NC}"
    exit 0
fi
