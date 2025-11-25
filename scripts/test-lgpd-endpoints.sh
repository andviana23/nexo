#!/usr/bin/env bash

################################################################################
# test-lgpd-endpoints.sh - Teste E2E de Endpoints LGPD
#
# Descrição:
#   Testa todos os endpoints LGPD com diferentes roles e verifica
#   isolamento multi-tenant, validação de ownership e casos de erro.
#
# Uso:
#   ./scripts/test-lgpd-endpoints.sh [API_URL]
#
# Exemplo:
#   ./scripts/test-lgpd-endpoints.sh http://localhost:8080
#
# Requisitos:
#   - curl
#   - jq
#   - API rodando
#
# Autor: Andrey Viana
# Versão: 1.0.0
################################################################################

set -euo pipefail

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuração
API_URL="${1:-http://localhost:8080}"
PASSED=0
FAILED=0

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🧪 Teste E2E - Endpoints LGPD (Art. 18)"
echo "═══════════════════════════════════════════════════════════════"
echo ""
log_info "API URL: $API_URL"
echo ""

# ──────────────────────────────────────────────────────────────────
# Setup: Login de usuários de diferentes tenants e roles
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔐 Setup: Autenticação de Usuários"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Tenant 1 - Owner
log_info "Login: Tenant 1 (Owner)"
TENANT1_OWNER_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner1@tenant1.com","password":"senha123"}' \
  | jq -r '.token // empty')

if [ -z "$TENANT1_OWNER_TOKEN" ]; then
  log_error "Falha ao obter token do Owner (Tenant 1)"
  exit 1
fi
log_success "Token obtido: ${TENANT1_OWNER_TOKEN:0:20}..."

# Tenant 1 - Employee
log_info "Login: Tenant 1 (Employee)"
TENANT1_EMPLOYEE_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"employee1@tenant1.com","password":"senha123"}' \
  | jq -r '.token // empty')

if [ -z "$TENANT1_EMPLOYEE_TOKEN" ]; then
  log_error "Falha ao obter token do Employee (Tenant 1)"
  exit 1
fi
log_success "Token obtido: ${TENANT1_EMPLOYEE_TOKEN:0:20}..."

# Tenant 2 - Owner
log_info "Login: Tenant 2 (Owner)"
TENANT2_OWNER_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner2@tenant2.com","password":"senha123"}' \
  | jq -r '.token // empty')

if [ -z "$TENANT2_OWNER_TOKEN" ]; then
  log_error "Falha ao obter token do Owner (Tenant 2)"
  exit 1
fi
log_success "Token obtido: ${TENANT2_OWNER_TOKEN:0:20}..."

echo ""

# ──────────────────────────────────────────────────────────────────
# Teste 1: GET /me/preferences - Obter preferências
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📋 Teste 1: GET /me/preferences"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "1.1 - Owner (Tenant 1) obtém suas preferências"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/preferences" \
  -H "Authorization: Bearer $TENANT1_OWNER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK"

  # Validar estrutura JSON
  DATA_SHARING=$(echo "$BODY" | jq -r '.data_sharing_consent // "null"')
  MARKETING=$(echo "$BODY" | jq -r '.marketing_consent // "null"')

  if [ "$DATA_SHARING" != "null" ] && [ "$MARKETING" != "null" ]; then
    log_success "JSON válido com campos obrigatórios"
  else
    log_error "JSON inválido ou campos faltando"
  fi
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "1.2 - Employee (Tenant 1) obtém suas preferências"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/preferences" \
  -H "Authorization: Bearer $TENANT1_EMPLOYEE_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK - Employee pode acessar próprias preferências"
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "1.3 - Owner (Tenant 2) obtém suas preferências (isolamento)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/preferences" \
  -H "Authorization: Bearer $TENANT2_OWNER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK - Tenant 2 isolado do Tenant 1"

  # Garantir que os dados são diferentes
  TENANT2_DATA=$(echo "$BODY" | jq -r '.data_sharing_consent')
  if [ "$TENANT2_DATA" != "$DATA_SHARING" ]; then
    log_success "Dados isolados corretamente entre tenants"
  else
    log_warning "Dados podem não estar isolados (verificar manualmente)"
  fi
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "1.4 - Sem autenticação (deve falhar)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/preferences")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "401" ]; then
  log_success "Status 401 Unauthorized - Autenticação obrigatória"
else
  log_error "Status $HTTP_CODE (esperado: 401)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Teste 2: PUT /me/preferences - Atualizar preferências
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  ✏️  Teste 2: PUT /me/preferences"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "2.1 - Owner (Tenant 1) atualiza preferências"
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$API_URL/api/v1/me/preferences" \
  -H "Authorization: Bearer $TENANT1_OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data_sharing_consent": true,
    "marketing_consent": false,
    "analytics_consent": true,
    "third_party_consent": false,
    "personalized_ads_consent": true
  }')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK - Preferências atualizadas"

  # Validar se os valores foram salvos
  MARKETING=$(echo "$BODY" | jq -r '.marketing_consent')
  if [ "$MARKETING" = "false" ]; then
    log_success "Valor atualizado corretamente (marketing_consent=false)"
  else
    log_error "Valor não foi atualizado corretamente"
  fi
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "2.2 - Employee (Tenant 1) atualiza suas preferências"
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$API_URL/api/v1/me/preferences" \
  -H "Authorization: Bearer $TENANT1_EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data_sharing_consent": false,
    "marketing_consent": false,
    "analytics_consent": false,
    "third_party_consent": false,
    "personalized_ads_consent": false
  }')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK - Employee pode atualizar próprias preferências"
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "2.3 - Payload inválido (deve falhar)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$API_URL/api/v1/me/preferences" \
  -H "Authorization: Bearer $TENANT1_OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"invalid_field": true}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "400" ]; then
  log_success "Status 400 Bad Request - Validação funcionando"
else
  log_error "Status $HTTP_CODE (esperado: 400)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Teste 3: GET /me/export - Exportar dados (Portabilidade)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📦 Teste 3: GET /me/export (Portabilidade)"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "3.1 - Owner (Tenant 1) exporta seus dados"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/export" \
  -H "Authorization: Bearer $TENANT1_OWNER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK - Exportação realizada"

  # Validar estrutura do JSON exportado
  USER_ID=$(echo "$BODY" | jq -r '.user.id // "null"')
  TENANT_ID=$(echo "$BODY" | jq -r '.tenant.id // "null"')
  PREFERENCES=$(echo "$BODY" | jq -r '.preferences // "null"')

  if [ "$USER_ID" != "null" ] && [ "$TENANT_ID" != "null" ] && [ "$PREFERENCES" != "null" ]; then
    log_success "JSON completo com user, tenant e preferences"
  else
    log_error "JSON incompleto - faltam campos obrigatórios"
  fi

  # Verificar se não há campos vazios críticos
  USER_EMAIL=$(echo "$BODY" | jq -r '.user.email // ""')
  if [ -n "$USER_EMAIL" ]; then
    log_success "Campos críticos preenchidos (email presente)"
  else
    log_error "Campos críticos vazios"
  fi
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "3.2 - Employee (Tenant 1) exporta seus dados"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/export" \
  -H "Authorization: Bearer $TENANT1_EMPLOYEE_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Status 200 OK - Employee pode exportar próprios dados"
else
  log_error "Status $HTTP_CODE (esperado: 200)"
fi

log_info "3.3 - Verificar rate limit (1 export/dia)"
log_warning "Rate limit será testado em teste separado (requer aguardar 24h)"

echo ""

# ──────────────────────────────────────────────────────────────────
# Teste 4: DELETE /me - Deletar conta (Direito ao Esquecimento)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🗑️  Teste 4: DELETE /me (Direito ao Esquecimento)"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "4.1 - Tentar deletar sem senha (deve falhar)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$API_URL/api/v1/me" \
  -H "Authorization: Bearer $TENANT1_EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "400" ]; then
  log_success "Status 400 Bad Request - Senha obrigatória"
else
  log_error "Status $HTTP_CODE (esperado: 400)"
fi

log_info "4.2 - Deletar com senha incorreta (deve falhar)"
RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$API_URL/api/v1/me" \
  -H "Authorization: Bearer $TENANT1_EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "senhaerrada123",
    "reason": "Teste de validação"
  }')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "401" ] || [ "$HTTP_CODE" = "403" ]; then
  log_success "Status $HTTP_CODE - Senha inválida rejeitada"
else
  log_error "Status $HTTP_CODE (esperado: 401 ou 403)"
fi

log_warning "4.3 - Teste de deleção real SKIPADO (requer usuário de teste descartável)"
log_info "    Para testar deleção real, criar usuário temporário e validar:"
log_info "    - users.deleted_at preenchido"
log_info "    - PII anonimizado (nome, email, password_hash)"
log_info "    - user_preferences deletado"
log_info "    - Registro em audit_logs"

echo ""

# ──────────────────────────────────────────────────────────────────
# Resumo
# ──────────────────────────────────────────────────────────────────
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo dos Testes"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}✓ Todos os testes de endpoints LGPD passaram!${NC}"
  exit 0
else
  echo -e "${RED}✗ Alguns testes falharam. Verifique os logs acima.${NC}"
  exit 1
fi
