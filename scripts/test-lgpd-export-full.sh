#!/usr/bin/env bash

################################################################################
# test-lgpd-export-full.sh - Validação Completa de Exportação LGPD
#
# Descrição:
#   Valida que a exportação retorna JSON completo sem campos vazios,
#   com todas as seções necessárias para portabilidade (Art. 18, V).
#
# Uso:
#   ./scripts/test-lgpd-export-full.sh [API_URL] [AUTH_TOKEN]
#
# Exemplo:
#   TOKEN=$(curl -s ... | jq -r '.token')
#   ./scripts/test-lgpd-export-full.sh http://localhost:8080 "$TOKEN"
#
# Requisitos:
#   - curl
#   - jq
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
TOKEN="${2:-}"
PASSED=0
FAILED=0

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  📦 Validação Completa - Exportação LGPD (Portabilidade)"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Validar token
if [ -z "$TOKEN" ]; then
  log_error "Token de autenticação não fornecido"
  echo ""
  echo "Uso: $0 [API_URL] [AUTH_TOKEN]"
  echo ""
  echo "Exemplo:"
  echo "  TOKEN=\$(curl -s -X POST http://localhost:8080/api/v1/auth/login \\"
  echo "    -H 'Content-Type: application/json' \\"
  echo "    -d '{\"email\":\"user@example.com\",\"password\":\"senha123\"}' \\"
  echo "    | jq -r '.token')"
  echo "  $0 http://localhost:8080 \"\$TOKEN\""
  exit 1
fi

log_info "API URL: $API_URL"
log_info "Token: ${TOKEN:0:20}..."
echo ""

# ──────────────────────────────────────────────────────────────────
# Executar exportação
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🚀 Executando Exportação"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "GET /api/v1/me/export"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/api/v1/me/export" \
  -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" != "200" ]; then
  log_error "Status $HTTP_CODE (esperado: 200)"
  echo ""
  echo "Resposta:"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
  exit 1
fi

log_success "Status 200 OK - Exportação realizada"

# Salvar JSON em arquivo temporário para análise
EXPORT_FILE="/tmp/lgpd-export-$(date +%s).json"
echo "$BODY" > "$EXPORT_FILE"
log_info "JSON salvo em: $EXPORT_FILE"

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar Estrutura do JSON
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔍 Validando Estrutura do JSON"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Verificar se é JSON válido
if jq empty "$EXPORT_FILE" 2>/dev/null; then
  log_success "JSON válido e bem formado"
else
  log_error "JSON inválido ou malformado"
  exit 1
fi

# ──────────────────────────────────────────────────────────────────
# Validar Seções Obrigatórias
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📋 Validando Seções Obrigatórias"
echo "───────────────────────────────────────────────────────────────"
echo ""

# 1. Seção: user
if jq -e '.user' "$EXPORT_FILE" > /dev/null 2>&1; then
  log_success "Seção 'user' presente"
else
  log_error "Seção 'user' ausente"
fi

# 2. Seção: tenant
if jq -e '.tenant' "$EXPORT_FILE" > /dev/null 2>&1; then
  log_success "Seção 'tenant' presente"
else
  log_error "Seção 'tenant' ausente"
fi

# 3. Seção: preferences
if jq -e '.preferences' "$EXPORT_FILE" > /dev/null 2>&1; then
  log_success "Seção 'preferences' presente"
else
  log_error "Seção 'preferences' ausente"
fi

# 4. Seção: audit_logs (opcional, mas recomendado)
if jq -e '.audit_logs' "$EXPORT_FILE" > /dev/null 2>&1; then
  log_success "Seção 'audit_logs' presente"
else
  log_warning "Seção 'audit_logs' ausente (recomendado incluir)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar Campos Obrigatórios: user
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  👤 Validando Campos: user"
echo "───────────────────────────────────────────────────────────────"
echo ""

USER_FIELDS=("id" "email" "nome" "role" "criado_em" "atualizado_em")

for field in "${USER_FIELDS[@]}"; do
  VALUE=$(jq -r ".user.$field // \"\"" "$EXPORT_FILE")

  if [ -n "$VALUE" ] && [ "$VALUE" != "null" ]; then
    log_success "user.$field = $VALUE"
  else
    log_error "user.$field está vazio ou null"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar Campos Obrigatórios: tenant
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🏢 Validando Campos: tenant"
echo "───────────────────────────────────────────────────────────────"
echo ""

TENANT_FIELDS=("id" "nome" "cnpj")

for field in "${TENANT_FIELDS[@]}"; do
  VALUE=$(jq -r ".tenant.$field // \"\"" "$EXPORT_FILE")

  if [ -n "$VALUE" ] && [ "$VALUE" != "null" ]; then
    log_success "tenant.$field = $VALUE"
  else
    log_error "tenant.$field está vazio ou null"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar Campos Obrigatórios: preferences
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔒 Validando Campos: preferences"
echo "───────────────────────────────────────────────────────────────"
echo ""

PREF_FIELDS=(
  "data_sharing_consent"
  "marketing_consent"
  "analytics_consent"
  "third_party_consent"
  "personalized_ads_consent"
)

for field in "${PREF_FIELDS[@]}"; do
  VALUE=$(jq -r ".preferences.$field // \"\"" "$EXPORT_FILE")

  if [ "$VALUE" = "true" ] || [ "$VALUE" = "false" ]; then
    log_success "preferences.$field = $VALUE"
  else
    log_error "preferences.$field inválido (esperado: true/false, recebido: '$VALUE')"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Verificar Campos Vazios/Corrompidos
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔎 Verificando Campos Vazios/Corrompidos"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Contar campos null
NULL_COUNT=$(jq '[.. | select(type == "null")] | length' "$EXPORT_FILE")
log_info "Total de campos null encontrados: $NULL_COUNT"

if [ "$NULL_COUNT" -eq 0 ]; then
  log_success "Nenhum campo null encontrado"
else
  log_warning "$NULL_COUNT campos null detectados (verificar se são opcionais)"
fi

# Contar campos vazios (strings "")
EMPTY_COUNT=$(jq '[.. | select(type == "string" and . == "")] | length' "$EXPORT_FILE")
log_info "Total de strings vazias encontradas: $EMPTY_COUNT"

if [ "$EMPTY_COUNT" -eq 0 ]; then
  log_success "Nenhuma string vazia encontrada"
else
  log_warning "$EMPTY_COUNT strings vazias detectadas"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar Tamanho do JSON
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📏 Validando Tamanho do JSON"
echo "───────────────────────────────────────────────────────────────"
echo ""

FILE_SIZE=$(wc -c < "$EXPORT_FILE")
log_info "Tamanho do arquivo: $FILE_SIZE bytes"

if [ "$FILE_SIZE" -lt 100 ]; then
  log_error "JSON muito pequeno (< 100 bytes) - pode estar incompleto"
elif [ "$FILE_SIZE" -gt 10485760 ]; then
  log_warning "JSON muito grande (> 10MB) - considerar streaming"
else
  log_success "Tamanho adequado"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar Metadados
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🏷️  Validando Metadados"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Verificar se há timestamp de exportação
EXPORTED_AT=$(jq -r '.exported_at // ""' "$EXPORT_FILE")
if [ -n "$EXPORTED_AT" ]; then
  log_success "Timestamp de exportação presente: $EXPORTED_AT"
else
  log_warning "Timestamp de exportação ausente (recomendado incluir)"
fi

# Verificar versão do formato
EXPORT_VERSION=$(jq -r '.version // ""' "$EXPORT_FILE")
if [ -n "$EXPORT_VERSION" ]; then
  log_success "Versão do formato presente: $EXPORT_VERSION"
else
  log_warning "Versão do formato ausente (recomendado incluir)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Resumo
# ──────────────────────────────────────────────────────────────────
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo da Validação"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo ""
echo "  Arquivo exportado: $EXPORT_FILE"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}✓ Exportação LGPD válida e completa!${NC}"
  exit 0
else
  echo -e "${RED}✗ Exportação tem campos ausentes ou inválidos.${NC}"
  exit 1
fi
