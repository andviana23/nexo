#!/usr/bin/env bash

################################################################################
# test-cookie-consent-banner.sh - Teste E2E de Banner de Consentimento
#
# Descrição:
#   Valida funcionamento do cookie consent banner:
#   - Preferências persistem após reload
#   - Sincronização com backend
#   - LocalStorage + API
#   - Botões "Aceitar tudo" / "Rejeitar tudo"
#
# Uso:
#   ./scripts/test-cookie-consent-banner.sh [FRONTEND_URL]
#
# Exemplo:
#   ./scripts/test-cookie-consent-banner.sh http://localhost:3000
#
# Requisitos:
#   - Playwright ou Cypress instalado
#   - Node.js
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
FRONTEND_URL="${1:-http://localhost:3000}"
PASSED=0
FAILED=0

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🍪 Teste E2E - Cookie Consent Banner"
echo "═══════════════════════════════════════════════════════════════"
echo ""

log_warning "Este teste requer ambiente de teste E2E (Playwright/Cypress)"
log_info "Frontend URL: $FRONTEND_URL"
echo ""

# ──────────────────────────────────────────────────────────────────
# Checklist Manual (Requer E2E real)
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📝 Checklist Manual de Teste"
echo "───────────────────────────────────────────────────────────────"
echo ""

echo "Execute os seguintes testes manualmente no navegador:"
echo ""

echo "1. EXIBIÇÃO INICIAL"
echo "   [ ] Acessar $FRONTEND_URL (sem estar logado)"
echo "   [ ] Banner aparece na primeira visita"
echo "   [ ] Banner não aparece se já foi respondido"
echo ""

echo "2. BOTÃO 'ACEITAR TUDO'"
echo "   [ ] Clicar em 'Aceitar tudo'"
echo "   [ ] Banner fecha"
echo "   [ ] localStorage['cookie-consent'] = JSON com todos true"
echo "   [ ] Recarregar página → banner não aparece mais"
echo ""

echo "3. BOTÃO 'REJEITAR TUDO'"
echo "   [ ] Limpar localStorage"
echo "   [ ] Recarregar página"
echo "   [ ] Clicar em 'Rejeitar tudo'"
echo "   [ ] Banner fecha"
echo "   [ ] localStorage['cookie-consent'] = JSON com todos false"
echo "   [ ] Recarregar página → banner não aparece mais"
echo ""

echo "4. PERSONALIZAÇÃO"
echo "   [ ] Limpar localStorage"
echo "   [ ] Recarregar página"
echo "   [ ] Abrir 'Personalizar cookies'"
echo "   [ ] Habilitar: Analytics ✓, Marketing ✗"
echo "   [ ] Salvar"
echo "   [ ] Verificar localStorage com valores corretos"
echo "   [ ] Recarregar → preferências mantidas"
echo ""

echo "5. PERSISTÊNCIA APÓS LOGIN"
echo "   [ ] Fazer login como usuário autenticado"
echo "   [ ] Verificar que preferências locais foram sincronizadas"
echo "   [ ] GET /api/v1/me/preferences retorna valores corretos"
echo ""

echo "6. SINCRONIZAÇÃO BACKEND"
echo "   [ ] Mudar preferências no banner (logado)"
echo "   [ ] Verificar chamada PUT /api/v1/me/preferences"
echo "   [ ] Status 200 OK"
echo "   [ ] Fazer logout e login novamente"
echo "   [ ] Preferências mantidas (vieram do backend)"
echo ""

echo "7. REVOGAÇÃO"
echo "   [ ] Ir para /privacy"
echo "   [ ] Clicar em 'Revogar consentimento'"
echo "   [ ] Banner reaparece"
echo "   [ ] Novas escolhas podem ser feitas"
echo ""

echo "8. INTEGRAÇÃO COM ANALYTICS"
echo "   [ ] Aceitar apenas 'Essenciais'"
echo "   [ ] Google Analytics NÃO deve carregar"
echo "   [ ] Aceitar 'Analytics'"
echo "   [ ] Google Analytics DEVE carregar"
echo ""

echo "───────────────────────────────────────────────────────────────"
echo ""

# ──────────────────────────────────────────────────────────────────
# Teste Básico: Verificar se componente existe
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔍 Teste Básico: Componente Existe"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Verificando se /privacy está acessível..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$FRONTEND_URL/privacy")

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Página /privacy acessível (Status 200)"
else
  log_error "Página /privacy não acessível (Status $HTTP_CODE)"
fi

log_info "Verificando se componente cookie-consent-banner.tsx existe..."
if [ -f "frontend/components/cookie-consent-banner.tsx" ]; then
  log_success "Componente cookie-consent-banner.tsx existe"
else
  log_error "Componente cookie-consent-banner.tsx não encontrado"
fi

log_info "Verificando se hook use-user-preferences.ts existe..."
if [ -f "frontend/hooks/use-user-preferences.ts" ]; then
  log_success "Hook use-user-preferences.ts existe"
else
  log_error "Hook use-user-preferences.ts não encontrado"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Teste: localStorage schema
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📦 Validar Schema do localStorage"
echo "───────────────────────────────────────────────────────────────"
echo ""

cat <<'EOF' > /tmp/test-localstorage-schema.js
// Exemplo de schema esperado no localStorage
const expectedSchema = {
  "cookie-consent": {
    "analytics": true,
    "marketing": false,
    "third_party": false,
    "timestamp": "2025-11-24T10:30:00Z"
  }
};

console.log("Schema esperado:");
console.log(JSON.stringify(expectedSchema, null, 2));

// Validação
const stored = localStorage.getItem('cookie-consent');
if (!stored) {
  console.error("❌ 'cookie-consent' não encontrado no localStorage");
  process.exit(1);
}

try {
  const parsed = JSON.parse(stored);

  // Verificar campos obrigatórios
  const required = ['analytics', 'marketing', 'third_party'];
  const missing = required.filter(f => !(f in parsed));

  if (missing.length > 0) {
    console.error(`❌ Campos faltando: ${missing.join(', ')}`);
    process.exit(1);
  }

  console.log("✅ Schema válido");
  process.exit(0);

} catch (err) {
  console.error("❌ JSON inválido:", err.message);
  process.exit(1);
}
EOF

log_info "Schema de validação criado em: /tmp/test-localstorage-schema.js"
log_warning "Execute no console do navegador para validar localStorage"

echo ""

# ──────────────────────────────────────────────────────────────────
# Resumo
# ──────────────────────────────────────────────────────────────────
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo ""
echo -e "${YELLOW}⚠️  Este teste requer validação manual ou E2E automatizado${NC}"
echo ""
echo "Para automação completa, criar:"
echo "  - tests/e2e/cookie-consent.spec.ts (Playwright)"
echo "  - cypress/e2e/cookie-consent.cy.ts (Cypress)"
echo ""

if [ $FAILED -eq 0 ]; then
  exit 0
else
  exit 1
fi
