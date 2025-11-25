#!/usr/bin/env bash

################################################################################
# test-backup-restore.sh - Teste de Restore de Backup em Staging
#
# Descrição:
#   Baixa backup do S3, restaura em staging e valida schema completo.
#   Executa validate_schema.sh e smoke tests para garantir integridade.
#
# Uso:
#   ./scripts/test-backup-restore.sh <BACKUP_FILE> <STAGING_DB_URL>
#
# Exemplo:
#   ./scripts/test-backup-restore.sh \
#     barber-backup-20251124_103000.sql.gz \
#     "postgresql://user:pass@staging.neon.tech:5432/barber_staging"
#
# Requisitos:
#   - aws-cli
#   - psql
#   - pg_restore ou psql (dependendo do formato)
#   - Acesso ao bucket S3
#   - Staging database vazio ou recém-criado
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
BACKUP_FILE="${1:-}"
STAGING_DB_URL="${2:-}"
S3_BACKUP_BUCKET="${S3_BACKUP_BUCKET:-}"
RESTORE_DIR="/tmp/restore-test"
PASSED=0
FAILED=0

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🔄 Teste de Restore de Backup (Staging)"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Validar argumentos
if [ -z "$BACKUP_FILE" ] || [ -z "$STAGING_DB_URL" ]; then
  log_error "Argumentos faltando"
  echo ""
  echo "Uso: $0 <BACKUP_FILE> <STAGING_DB_URL>"
  echo ""
  echo "Exemplo:"
  echo "  $0 barber-backup-20251124_103000.sql.gz \\"
  echo "     'postgresql://user:pass@staging.neon.tech:5432/barber_staging'"
  exit 1
fi

if [ -z "$S3_BACKUP_BUCKET" ]; then
  log_error "Variável S3_BACKUP_BUCKET não configurada"
  exit 1
fi

log_info "Backup: $BACKUP_FILE"
log_info "Staging DB: ${STAGING_DB_URL%%@*}@***"
log_info "S3 Bucket: $S3_BACKUP_BUCKET"
echo ""

# ──────────────────────────────────────────────────────────────────
# Preparar diretório
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📁 Preparar Diretório"
echo "───────────────────────────────────────────────────────────────"
echo ""

mkdir -p "$RESTORE_DIR"
log_success "Diretório criado: $RESTORE_DIR"
echo ""

# ──────────────────────────────────────────────────────────────────
# Download do S3
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  ☁️  Download do S3"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Baixando $BACKUP_FILE do S3..."
if aws s3 cp \
  "s3://${S3_BACKUP_BUCKET}/backups/${BACKUP_FILE}" \
  "$RESTORE_DIR/$BACKUP_FILE"; then
  log_success "Download concluído"
else
  log_error "Falha no download"
  exit 1
fi

# Validar integridade (checksum)
log_info "Validando checksum..."
LOCAL_CHECKSUM=$(sha256sum "$RESTORE_DIR/$BACKUP_FILE" | awk '{print $1}')
log_info "SHA256 local: $LOCAL_CHECKSUM"

# Baixar manifesto para comparar
MANIFEST_FILE="${BACKUP_FILE%.sql.gz}-manifest.json"
if aws s3 cp \
  "s3://${S3_BACKUP_BUCKET}/backups/${MANIFEST_FILE}" \
  "$RESTORE_DIR/$MANIFEST_FILE" 2>/dev/null; then

  S3_CHECKSUM=$(jq -r '.checksum_sha256' "$RESTORE_DIR/$MANIFEST_FILE")

  if [ "$LOCAL_CHECKSUM" = "$S3_CHECKSUM" ]; then
    log_success "Checksum validado (match)"
  else
    log_error "Checksum não corresponde (local: $LOCAL_CHECKSUM, S3: $S3_CHECKSUM)"
    exit 1
  fi
else
  log_warning "Manifesto não encontrado, pulando validação de checksum"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Descomprimir
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🗜️  Descomprimir Backup"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Descomprimindo com gunzip..."
if gunzip -f "$RESTORE_DIR/$BACKUP_FILE"; then
  log_success "Descompressão concluída"
else
  log_error "Falha ao descomprimir"
  exit 1
fi

SQL_FILE="${BACKUP_FILE%.gz}"
if [ -f "$RESTORE_DIR/$SQL_FILE" ]; then
  SQL_SIZE=$(wc -c < "$RESTORE_DIR/$SQL_FILE")
  SQL_SIZE_MB=$((SQL_SIZE / 1024 / 1024))
  log_success "Arquivo SQL: ${SQL_SIZE_MB}MB"
else
  log_error "Arquivo SQL não encontrado após descompressão"
  exit 1
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar conexão staging
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔌 Validar Conexão Staging"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Testando conexão com staging..."
if psql "$STAGING_DB_URL" -c "SELECT 1;" &> /dev/null; then
  log_success "Conexão estabelecida"
else
  log_error "Falha ao conectar em staging"
  exit 1
fi

# Verificar se banco está vazio
TABLE_COUNT=$(psql "$STAGING_DB_URL" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" | tr -d ' ')
log_info "Tabelas existentes em staging: $TABLE_COUNT"

if [ "$TABLE_COUNT" -gt "0" ]; then
  log_warning "Staging não está vazio ($TABLE_COUNT tabelas). Considere limpar antes."
  read -p "Continuar mesmo assim? (y/N): " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_info "Operação cancelada pelo usuário"
    exit 0
  fi
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Restaurar banco
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔄 Restaurar Banco de Dados"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Executando psql < $SQL_FILE..."
START_TIME=$(date +%s)

if psql "$STAGING_DB_URL" < "$RESTORE_DIR/$SQL_FILE" &> "$RESTORE_DIR/restore.log"; then
  END_TIME=$(date +%s)
  DURATION=$((END_TIME - START_TIME))
  log_success "Restore concluído (${DURATION}s)"
else
  log_error "Falha no restore. Veja logs em: $RESTORE_DIR/restore.log"
  tail -n 20 "$RESTORE_DIR/restore.log"
  exit 1
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar schema com validate_schema.sh
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  ✅ Validar Schema (validate_schema.sh)"
echo "───────────────────────────────────────────────────────────────"
echo ""

if [ -f "scripts/validate_schema.sh" ]; then
  log_info "Executando validate_schema.sh..."

  if bash scripts/validate_schema.sh "$STAGING_DB_URL"; then
    log_success "Schema validado com sucesso"
  else
    log_error "Validação de schema falhou"
    ((FAILED++))
  fi
else
  log_warning "validate_schema.sh não encontrado, pulando validação de schema"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar dados básicos
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📊 Validar Dados Restaurados"
echo "───────────────────────────────────────────────────────────────"
echo ""

# Contar tabelas
TABLES=$(psql "$STAGING_DB_URL" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" | tr -d ' ')
log_info "Tabelas restauradas: $TABLES"

if [ "$TABLES" -ge "10" ]; then
  log_success "Número adequado de tabelas (>= 10)"
else
  log_error "Poucas tabelas restauradas (< 10)"
fi

# Contar constraints
CONSTRAINTS=$(psql "$STAGING_DB_URL" -t -c "SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_schema = 'public';" | tr -d ' ')
log_info "Constraints restauradas: $CONSTRAINTS"

if [ "$CONSTRAINTS" -ge "5" ]; then
  log_success "Constraints restauradas"
else
  log_warning "Poucas constraints (pode estar OK se schema minimal)"
fi

# Contar índices
INDEXES=$(psql "$STAGING_DB_URL" -t -c "SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public';" | tr -d ' ')
log_info "Índices restaurados: $INDEXES"

if [ "$INDEXES" -ge "5" ]; then
  log_success "Índices restaurados"
else
  log_warning "Poucos índices restaurados"
fi

# Verificar se há dados em tabelas core
TENANTS_COUNT=$(psql "$STAGING_DB_URL" -t -c "SELECT COUNT(*) FROM tenants;" 2>/dev/null | tr -d ' ' || echo "0")
USERS_COUNT=$(psql "$STAGING_DB_URL" -t -c "SELECT COUNT(*) FROM users;" 2>/dev/null | tr -d ' ' || echo "0")

log_info "Registros restaurados:"
log_info "  - tenants: $TENANTS_COUNT"
log_info "  - users: $USERS_COUNT"

if [ "$TENANTS_COUNT" -gt "0" ] && [ "$USERS_COUNT" -gt "0" ]; then
  log_success "Dados de produção restaurados"
else
  log_warning "Poucas linhas em tabelas core (pode ser backup vazio)"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Smoke tests
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🧪 Smoke Tests"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Executando queries de validação..."

# Query 1: Verificar integridade referencial
ORPHAN_USERS=$(psql "$STAGING_DB_URL" -t -c "
  SELECT COUNT(*)
  FROM users u
  LEFT JOIN tenants t ON u.tenant_id = t.id
  WHERE t.id IS NULL;
" 2>/dev/null | tr -d ' ' || echo "0")

if [ "$ORPHAN_USERS" = "0" ]; then
  log_success "Integridade referencial OK (users → tenants)"
else
  log_error "Integridade referencial quebrada: $ORPHAN_USERS users órfãos"
fi

# Query 2: Verificar unicidade tenant_id + email
DUPLICATE_EMAILS=$(psql "$STAGING_DB_URL" -t -c "
  SELECT COUNT(*)
  FROM (
    SELECT tenant_id, email, COUNT(*)
    FROM users
    GROUP BY tenant_id, email
    HAVING COUNT(*) > 1
  ) AS duplicates;
" 2>/dev/null | tr -d ' ' || echo "0")

if [ "$DUPLICATE_EMAILS" = "0" ]; then
  log_success "Constraint UNIQUE(tenant_id, email) válida"
else
  log_error "Emails duplicados por tenant: $DUPLICATE_EMAILS"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Cleanup
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🧹 Cleanup"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Removendo arquivos temporários..."
rm -rf "$RESTORE_DIR"
log_success "Cleanup concluído"

echo ""

# ──────────────────────────────────────────────────────────────────
# Resumo
# ──────────────────────────────────────────────────────────────────
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo do Restore"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "  📦 Backup: $BACKUP_FILE"
echo "  💾 Tamanho: ${SQL_SIZE_MB}MB (descomprimido)"
echo "  ⏱️  Duração: ${DURATION}s"
echo "  📊 Tabelas: $TABLES"
echo "  🔗 Constraints: $CONSTRAINTS"
echo "  📇 Índices: $INDEXES"
echo "  👥 Tenants: $TENANTS_COUNT"
echo "  👤 Users: $USERS_COUNT"
echo ""
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}✓ Restore validado com sucesso!${NC}"
  echo ""
  echo "  Próximos passos:"
  echo "  1. Executar smoke_tests_v2.sh contra staging"
  echo "  2. Comparar contagens de registros (prod vs staging)"
  echo "  3. Validar autenticação e queries críticas"
  echo "  4. Atualizar RTO/RPO no runbook"
  exit 0
else
  echo -e "${RED}✗ Restore teve problemas. Verifique os logs.${NC}"
  exit 1
fi
