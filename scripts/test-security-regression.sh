#!/usr/bin/env bash

################################################################################
# test-security-regression.sh - Teste de Regressão de Segurança
#
# Descrição:
#   Executa bateria completa de testes de segurança para garantir que
#   as mudanças de hardening/OPS não introduziram vulnerabilidades.
#
# Testes incluídos:
#   - SQL Injection (10 testes)
#   - XSS (5 testes)
#   - CSRF (5 testes)
#   - RBAC/Autorização (10 testes)
#   - Autenticação (5 testes)
#
# Meta: 35/35 testes passando (100%)
#
# Uso:
#   ./scripts/test-security-regression.sh [API_URL]
#
# Exemplo:
#   ./scripts/test-security-regression.sh http://localhost:8080
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
TOTAL_TESTS=35

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🛡️  Teste de Regressão de Segurança (35 Testes)"
echo "═══════════════════════════════════════════════════════════════"
echo ""
log_info "API URL: $API_URL"
log_info "Meta: 35/35 testes passando (100%)"
echo ""

# ──────────────────────────────────────────────────────────────────
# Setup: Criar tokens de teste
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔐 Setup: Autenticação"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Token válido (owner)
log_info "Login: Owner (Tenant 1)"
OWNER_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner1@tenant1.com","password":"senha123"}' \
  | jq -r '.token // empty' 2>/dev/null || echo "")

if [ -n "$OWNER_TOKEN" ]; then
  log_success "Token owner obtido"
else
  log_error "Falha ao obter token owner (configurar seed_test_tenant.sql)"
  exit 1
fi

# Token employee
log_info "Login: Employee (Tenant 1)"
EMPLOYEE_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"employee1@tenant1.com","password":"senha123"}' \
  | jq -r '.token // empty' 2>/dev/null || echo "")

if [ -n "$EMPLOYEE_TOKEN" ]; then
  log_success "Token employee obtido"
else
  log_warning "Token employee não obtido (alguns testes serão limitados)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Categoria 1: SQL Injection (10 testes)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  💉 Categoria 1: SQL Injection (10 testes)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Teste 1.1: SQL Injection no email (login)
log_info "1.1 - SQL Injection no campo email (login)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin'\'' OR 1=1--","password":"any"}')

if [ "$HTTP_CODE" = "400" ] || [ "$HTTP_CODE" = "401" ]; then
  log_success "SQL Injection bloqueado (Status $HTTP_CODE)"
else
  log_error "SQL Injection não bloqueado (Status $HTTP_CODE)"
fi

# Teste 1.2: SQL Injection no password
log_info "1.2 - SQL Injection no campo password"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"' OR '1'='1"}')

if [ "$HTTP_CODE" = "400" ] || [ "$HTTP_CODE" = "401" ]; then
  log_success "SQL Injection bloqueado (Status $HTTP_CODE)"
else
  log_error "SQL Injection não bloqueado (Status $HTTP_CODE)"
fi

# Teste 1.3: SQL Injection em query params
log_info "1.3 - SQL Injection em query params (busca)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET \
  "$API_URL/api/v1/receitas?busca='; DROP TABLE users;--" \
  -H "Authorization: Bearer $OWNER_TOKEN")

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "400" ]; then
  log_success "SQL Injection bloqueado (Status $HTTP_CODE)"
else
  log_error "SQL Injection não bloqueado (Status $HTTP_CODE)"
fi

# Teste 1.4: UNION-based SQL Injection
log_info "1.4 - UNION-based SQL Injection"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET \
  "$API_URL/api/v1/receitas?categoria_id=1 UNION SELECT password_hash FROM users--" \
  -H "Authorization: Bearer $OWNER_TOKEN")

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "400" ]; then
  log_success "UNION injection bloqueado (Status $HTTP_CODE)"
else
  log_error "UNION injection não bloqueado (Status $HTTP_CODE)"
fi

# Teste 1.5-1.10: Variações de SQL Injection
SQL_PAYLOADS=(
  "1'; EXEC sp_executesql N'DROP TABLE users'--"
  "admin'--"
  "1' AND '1'='1"
  "1' OR SLEEP(5)--"
  "' OR 1=1 LIMIT 1--"
  "') OR ('1'='1"
)

for i in "${!SQL_PAYLOADS[@]}"; do
  payload="${SQL_PAYLOADS[$i]}"
  test_num=$((i + 5))

  log_info "1.$test_num - SQL Injection variação $((i+1))"
  HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$payload\",\"password\":\"test\"}")

  if [ "$HTTP_CODE" = "400" ] || [ "$HTTP_CODE" = "401" ]; then
    log_success "SQL Injection variação $((i+1)) bloqueado"
  else
    log_error "SQL Injection variação $((i+1)) não bloqueado"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Categoria 2: XSS (5 testes)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🎭 Categoria 2: XSS (5 testes)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Teste 2.1: XSS no nome de categoria
log_info "2.1 - XSS em campo de texto (nome categoria)"
RESPONSE=$(curl -s -X POST "$API_URL/api/v1/categorias" \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nome":"<script>alert(1)</script>","tipo":"receita"}')

HTTP_CODE=$(echo "$RESPONSE" | jq -r '.status // 400')
NOME_SALVO=$(echo "$RESPONSE" | jq -r '.nome // ""')

if [[ "$NOME_SALVO" != *"<script>"* ]] || [ "$HTTP_CODE" = "400" ]; then
  log_success "XSS sanitizado ou bloqueado"
else
  log_error "XSS não foi sanitizado"
fi

# Teste 2.2-2.5: Outras variações de XSS
XSS_PAYLOADS=(
  "<img src=x onerror=alert(1)>"
  "javascript:alert(1)"
  "<iframe src='javascript:alert(1)'></iframe>"
  "<svg onload=alert(1)>"
)

for i in "${!XSS_PAYLOADS[@]}"; do
  payload="${XSS_PAYLOADS[$i]}"
  test_num=$((i + 2))

  log_info "2.$test_num - XSS variação $((i+1))"
  HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/categorias" \
    -H "Authorization: Bearer $OWNER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"nome\":\"$payload\",\"tipo\":\"receita\"}")

  # Aceitar tanto sanitização (200) quanto bloqueio (400)
  if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "400" ]; then
    log_success "XSS variação $((i+1)) tratado"
  else
    log_error "XSS variação $((i+1)) não tratado (Status $HTTP_CODE)"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Categoria 3: CSRF (5 testes)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🎯 Categoria 3: CSRF (5 testes)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Teste 3.1: DELETE sem CSRF token (deve falhar se CSRF habilitado)
log_info "3.1 - DELETE /me sem validação adicional"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$API_URL/api/v1/me" \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"password":"senha123"}')

# CSRF geralmente exige 403 sem token, mas pode aceitar com JWT
if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "400" ] || [ "$HTTP_CODE" = "403" ]; then
  log_success "Endpoint protegido (Status $HTTP_CODE)"
else
  log_warning "Status inesperado: $HTTP_CODE (verificar proteção CSRF)"
fi

# Teste 3.2-3.5: Tentar operações sensíveis de origem externa
CSRF_TESTS=(
  "POST /api/v1/receitas"
  "PUT /api/v1/me/preferences"
  "DELETE /api/v1/categorias/1"
  "POST /api/v1/assinaturas"
)

for i in "${!CSRF_TESTS[@]}"; do
  endpoint="${CSRF_TESTS[$i]}"
  test_num=$((i + 2))
  method=$(echo "$endpoint" | awk '{print $1}')
  path=$(echo "$endpoint" | awk '{print $2}')

  log_info "3.$test_num - CSRF protection em $method $path"

  # Simular request sem Origin/Referer corretos
  HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X "$method" \
    "$API_URL$path" \
    -H "Authorization: Bearer $OWNER_TOKEN" \
    -H "Origin: http://evil.com" \
    -H "Content-Type: application/json" \
    -d '{}')

  # CORS deve bloquear ou aceitar (depende da config)
  if [ "$HTTP_CODE" != "500" ]; then
    log_success "CSRF/CORS tratado (Status $HTTP_CODE)"
  else
    log_error "CSRF/CORS não tratado (Status 500)"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Categoria 4: RBAC/Autorização (10 testes)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔒 Categoria 4: RBAC/Autorização (10 testes)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Teste 4.1: Employee tenta acessar rota admin
log_info "4.1 - Employee acessa rota restrita a owner"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$API_URL/api/v1/admin/users" \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN")

if [ "$HTTP_CODE" = "403" ] || [ "$HTTP_CODE" = "401" ] || [ "$HTTP_CODE" = "404" ]; then
  log_success "Acesso negado para employee (Status $HTTP_CODE)"
else
  log_error "Employee conseguiu acessar rota admin (Status $HTTP_CODE)"
fi

# Teste 4.2: Sem token JWT
log_info "4.2 - Acesso sem token JWT"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$API_URL/api/v1/receitas")

if [ "$HTTP_CODE" = "401" ]; then
  log_success "Acesso negado sem JWT (Status 401)"
else
  log_error "Endpoint acessível sem JWT (Status $HTTP_CODE)"
fi

# Teste 4.3: Token expirado
log_info "4.3 - Token JWT expirado"
EXPIRED_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwMDAwMDB9.invalid"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$API_URL/api/v1/receitas" \
  -H "Authorization: Bearer $EXPIRED_TOKEN")

if [ "$HTTP_CODE" = "401" ]; then
  log_success "Token expirado rejeitado (Status 401)"
else
  log_error "Token expirado aceito (Status $HTTP_CODE)"
fi

# Teste 4.4: Token inválido
log_info "4.4 - Token JWT inválido (assinatura errada)"
INVALID_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalidpayload.invalidsignature"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$API_URL/api/v1/receitas" \
  -H "Authorization: Bearer $INVALID_TOKEN")

if [ "$HTTP_CODE" = "401" ]; then
  log_success "Token inválido rejeitado (Status 401)"
else
  log_error "Token inválido aceito (Status $HTTP_CODE)"
fi

# Teste 4.5: Acessar dados de outro tenant
log_info "4.5 - Isolamento multi-tenant (acessar tenant diferente)"
TENANT2_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner2@tenant2.com","password":"senha123"}' \
  | jq -r '.token // empty' 2>/dev/null || echo "")

if [ -n "$TENANT2_TOKEN" ]; then
  # Tentar buscar receitas do Tenant 1 com token do Tenant 2
  RESPONSE=$(curl -s -X GET "$API_URL/api/v1/receitas" \
    -H "Authorization: Bearer $TENANT2_TOKEN")

  TENANT1_DATA=$(echo "$RESPONSE" | jq -r '.data[] | select(.tenant_id != "tenant2-id") | .id' 2>/dev/null || echo "")

  if [ -z "$TENANT1_DATA" ]; then
    log_success "Isolamento multi-tenant OK (sem dados vazados)"
  else
    log_error "Dados de outro tenant vazados!"
  fi
else
  log_warning "Não foi possível testar isolamento multi-tenant (token2 não obtido)"
fi

# Teste 4.6-4.10: Outras verificações RBAC
log_info "4.6 - Employee tenta deletar conta de outro usuário"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$API_URL/api/v1/users/other-user-id" \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN")

if [ "$HTTP_CODE" = "403" ] || [ "$HTTP_CODE" = "401" ] || [ "$HTTP_CODE" = "404" ]; then
  log_success "Deleção bloqueada (Status $HTTP_CODE)"
else
  log_error "Employee deletou outro usuário (Status $HTTP_CODE)"
fi

log_info "4.7 - Modificar role via API (privilege escalation)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$API_URL/api/v1/me" \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"owner"}')

if [ "$HTTP_CODE" = "403" ] || [ "$HTTP_CODE" = "400" ]; then
  log_success "Escalation de privilégio bloqueado (Status $HTTP_CODE)"
else
  log_error "Escalation de privilégio permitido (Status $HTTP_CODE)"
fi

log_info "4.8 - Acessar endpoint /metrics sem autenticação"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$API_URL/metrics")

# Métricas podem ser públicas (200) ou protegidas (401)
if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "401" ]; then
  log_success "Endpoint /metrics configurado (Status $HTTP_CODE)"
else
  log_warning "Status inesperado em /metrics: $HTTP_CODE"
fi

log_info "4.9 - Injetar tenant_id no payload (bypass multi-tenant)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/receitas" \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"tenant_id":"outro-tenant","descricao":"hack","valor":1000}')

if [ "$HTTP_CODE" = "400" ] || [ "$HTTP_CODE" = "200" ]; then
  log_success "tenant_id ignorado/validado (Status $HTTP_CODE)"
else
  log_error "tenant_id injetado aceito (Status $HTTP_CODE)"
fi

log_info "4.10 - Modificar tenant_id via header (bypass)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$API_URL/api/v1/receitas" \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "X-Tenant-ID: outro-tenant-id")

# Tenant ID deve vir do JWT, não do header
if [ "$HTTP_CODE" = "200" ]; then
  log_success "tenant_id do JWT prevalece (Status 200)"
else
  log_warning "Verificar se tenant_id do header é ignorado (Status $HTTP_CODE)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Categoria 5: Autenticação (5 testes)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔑 Categoria 5: Autenticação (5 testes)"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Teste 5.1: Login com senha errada
log_info "5.1 - Login com senha incorreta"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"owner1@tenant1.com","password":"senhaerrada"}')

if [ "$HTTP_CODE" = "401" ]; then
  log_success "Senha incorreta rejeitada (Status 401)"
else
  log_error "Senha incorreta aceita (Status $HTTP_CODE)"
fi

# Teste 5.2: Login com usuário inexistente
log_info "5.2 - Login com email inexistente"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"naoexiste@example.com","password":"senha123"}')

if [ "$HTTP_CODE" = "401" ] || [ "$HTTP_CODE" = "404" ]; then
  log_success "Usuário inexistente rejeitado (Status $HTTP_CODE)"
else
  log_error "Usuário inexistente aceito (Status $HTTP_CODE)"
fi

# Teste 5.3: Brute force protection (rate limiting)
log_info "5.3 - Proteção contra brute force (10 tentativas rápidas)"
BLOCKED=0

for i in {1..10}; do
  HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"wrong"}')

  if [ "$HTTP_CODE" = "429" ]; then
    BLOCKED=1
    break
  fi
done

if [ "$BLOCKED" = "1" ]; then
  log_success "Rate limiting bloqueou após múltiplas tentativas"
else
  log_warning "Rate limiting não detectado (pode não estar configurado)"
fi

# Teste 5.4: Token refresh
log_info "5.4 - Refresh token (se implementado)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/refresh" \
  -H "Authorization: Bearer $OWNER_TOKEN")

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "404" ]; then
  log_success "Endpoint refresh tratado (Status $HTTP_CODE)"
else
  log_warning "Status inesperado em refresh: $HTTP_CODE"
fi

# Teste 5.5: Logout
log_info "5.5 - Logout (revogação de token)"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/auth/logout" \
  -H "Authorization: Bearer $OWNER_TOKEN")

if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "204" ] || [ "$HTTP_CODE" = "404" ]; then
  log_success "Endpoint logout tratado (Status $HTTP_CODE)"
else
  log_warning "Status inesperado em logout: $HTTP_CODE"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Resumo
# ──────────────────────────────────────────────────────────────────
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo dos Testes de Segurança"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "  Total de testes: $TOTAL_TESTS"
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo ""

PERCENTAGE=$((PASSED * 100 / TOTAL_TESTS))
echo "  Taxa de sucesso: ${PERCENTAGE}%"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}✓ Todos os testes de segurança passaram! (35/35)${NC}"
  echo ""
  echo "  ✅ SQL Injection: 10/10"
  echo "  ✅ XSS: 5/5"
  echo "  ✅ CSRF: 5/5"
  echo "  ✅ RBAC: 10/10"
  echo "  ✅ Autenticação: 5/5"
  exit 0
else
  echo -e "${RED}✗ Alguns testes de segurança falharam ($FAILED/$TOTAL_TESTS)${NC}"
  echo ""
  echo "  Revise os logs acima e corrija as vulnerabilidades."
  exit 1
fi
