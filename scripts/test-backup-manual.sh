#!/usr/bin/env bash

################################################################################
# test-backup-manual.sh - Teste Manual de Backup
#
# Descrição:
#   Executa backup manual do PostgreSQL, faz upload para S3 e valida:
#   - Arquivo .sql.gz gerado
#   - Checksum válido
#   - Upload S3 bem-sucedido
#   - Tamanho adequado
#   - Manifesto JSON criado
#
# Uso:
#   ./scripts/test-backup-manual.sh
#
# Requisitos:
#   - psql
#   - pg_dump
#   - gzip
#   - aws-cli (configurado)
#   - Variáveis de ambiente:
#       - NEON_DB_HOST
#       - NEON_DB_USER
#       - NEON_DB_PASSWORD
#       - NEON_DB_NAME
#       - S3_BACKUP_BUCKET
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
PASSED=0
FAILED=0
BACKUP_DIR="${BACKUP_DIR:-/tmp/backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="barber-backup-${TIMESTAMP}.sql"
COMPRESSED_FILE="${BACKUP_FILE}.gz"
MANIFEST_FILE="barber-backup-${TIMESTAMP}-manifest.json"

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  💾 Teste Manual de Backup PostgreSQL → S3"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# ──────────────────────────────────────────────────────────────────
# Validar variáveis de ambiente
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔐 Validar Credenciais"
echo "───────────────────────────────────────────────────────────────"
echo ""

REQUIRED_VARS=(
  "NEON_DB_HOST"
  "NEON_DB_USER"
  "NEON_DB_PASSWORD"
  "NEON_DB_NAME"
  "S3_BACKUP_BUCKET"
)

for var in "${REQUIRED_VARS[@]}"; do
  if [ -z "${!var:-}" ]; then
    log_error "Variável $var não configurada"
    ((FAILED++))
  else
    log_success "Variável $var configurada"
  fi
done

if [ $FAILED -gt 0 ]; then
  echo ""
  log_error "Configure as variáveis de ambiente antes de continuar"
  exit 1
fi

log_info "Database: ${NEON_DB_USER}@${NEON_DB_HOST}/${NEON_DB_NAME}"
log_info "S3 Bucket: ${S3_BACKUP_BUCKET}"
echo ""

# ──────────────────────────────────────────────────────────────────
# Preparar diretório de backup
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📁 Preparar Diretório"
echo "───────────────────────────────────────────────────────────────"
echo ""

mkdir -p "$BACKUP_DIR"
log_success "Diretório criado: $BACKUP_DIR"
log_info "Backup será salvo em: $BACKUP_DIR/$COMPRESSED_FILE"
echo ""

# ──────────────────────────────────────────────────────────────────
# Executar pg_dump
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔄 Executar pg_dump"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Iniciando backup..."
START_TIME=$(date +%s)

export PGPASSWORD="$NEON_DB_PASSWORD"

if pg_dump \
  -h "$NEON_DB_HOST" \
  -U "$NEON_DB_USER" \
  -d "$NEON_DB_NAME" \
  -F p \
  --no-owner \
  --no-acl \
  -f "$BACKUP_DIR/$BACKUP_FILE"; then

  END_TIME=$(date +%s)
  DURATION=$((END_TIME - START_TIME))
  log_success "pg_dump concluído (${DURATION}s)"
else
  log_error "pg_dump falhou"
  exit 1
fi

unset PGPASSWORD

# Validar arquivo gerado
if [ -f "$BACKUP_DIR/$BACKUP_FILE" ]; then
  FILE_SIZE=$(wc -c < "$BACKUP_DIR/$BACKUP_FILE")
  FILE_SIZE_MB=$((FILE_SIZE / 1024 / 1024))
  log_success "Arquivo gerado: ${FILE_SIZE_MB}MB"

  if [ "$FILE_SIZE" -lt 100 ]; then
    log_error "Arquivo muito pequeno (< 100 bytes) - backup pode estar vazio"
  else
    log_success "Tamanho adequado"
  fi
else
  log_error "Arquivo não foi gerado"
  exit 1
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Comprimir arquivo
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🗜️  Comprimir Backup"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Comprimindo com gzip..."
if gzip -f "$BACKUP_DIR/$BACKUP_FILE"; then
  log_success "Compressão concluída"
else
  log_error "Falha ao comprimir"
  exit 1
fi

if [ -f "$BACKUP_DIR/$COMPRESSED_FILE" ]; then
  COMPRESSED_SIZE=$(wc -c < "$BACKUP_DIR/$COMPRESSED_FILE")
  COMPRESSED_SIZE_MB=$((COMPRESSED_SIZE / 1024 / 1024))
  COMPRESSION_RATIO=$(awk "BEGIN {printf \"%.1f\", ($FILE_SIZE / $COMPRESSED_SIZE)}")

  log_success "Arquivo comprimido: ${COMPRESSED_SIZE_MB}MB (ratio: ${COMPRESSION_RATIO}x)"
else
  log_error "Arquivo .gz não foi gerado"
  exit 1
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Calcular checksum
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔐 Calcular Checksum"
echo "───────────────────────────────────────────────────────────────"
echo ""

CHECKSUM=$(sha256sum "$BACKUP_DIR/$COMPRESSED_FILE" | awk '{print $1}')
log_success "SHA256: $CHECKSUM"
echo ""

# ──────────────────────────────────────────────────────────────────
# Criar manifesto JSON
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📝 Criar Manifesto"
echo "───────────────────────────────────────────────────────────────"
echo ""

cat > "$BACKUP_DIR/$MANIFEST_FILE" <<EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "database": "${NEON_DB_NAME}",
  "host": "${NEON_DB_HOST}",
  "backup_file": "${COMPRESSED_FILE}",
  "size_bytes": ${COMPRESSED_SIZE},
  "size_mb": ${COMPRESSED_SIZE_MB},
  "original_size_bytes": ${FILE_SIZE},
  "compression_ratio": "${COMPRESSION_RATIO}",
  "checksum_sha256": "${CHECKSUM}",
  "duration_seconds": ${DURATION},
  "pg_version": "$(psql --version | awk '{print $3}')",
  "backup_type": "full"
}
EOF

log_success "Manifesto criado: $MANIFEST_FILE"
cat "$BACKUP_DIR/$MANIFEST_FILE" | jq '.'
echo ""

# ──────────────────────────────────────────────────────────────────
# Upload para S3
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  ☁️  Upload para S3"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Verificando credenciais AWS..."
if aws sts get-caller-identity &> /dev/null; then
  log_success "Credenciais AWS válidas"
else
  log_error "Credenciais AWS inválidas. Execute: aws configure"
  exit 1
fi

log_info "Fazendo upload do backup..."
if aws s3 cp \
  "$BACKUP_DIR/$COMPRESSED_FILE" \
  "s3://${S3_BACKUP_BUCKET}/backups/${COMPRESSED_FILE}" \
  --storage-class STANDARD_IA \
  --server-side-encryption AES256; then
  log_success "Backup enviado para S3"
else
  log_error "Falha no upload do backup"
  exit 1
fi

log_info "Fazendo upload do manifesto..."
if aws s3 cp \
  "$BACKUP_DIR/$MANIFEST_FILE" \
  "s3://${S3_BACKUP_BUCKET}/backups/${MANIFEST_FILE}" \
  --content-type "application/json" \
  --server-side-encryption AES256; then
  log_success "Manifesto enviado para S3"
else
  log_error "Falha no upload do manifesto"
  exit 1
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar artefatos no S3
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  ✅ Validar Artefatos no S3"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Verificando arquivo no S3..."
S3_SIZE=$(aws s3api head-object \
  --bucket "$S3_BACKUP_BUCKET" \
  --key "backups/$COMPRESSED_FILE" \
  --query ContentLength \
  --output text 2>/dev/null || echo "0")

if [ "$S3_SIZE" -eq "$COMPRESSED_SIZE" ]; then
  log_success "Tamanho no S3 corresponde ($S3_SIZE bytes)"
else
  log_error "Tamanho no S3 não corresponde (local: $COMPRESSED_SIZE, S3: $S3_SIZE)"
fi

log_info "Verificando criptografia..."
ENCRYPTION=$(aws s3api head-object \
  --bucket "$S3_BACKUP_BUCKET" \
  --key "backups/$COMPRESSED_FILE" \
  --query ServerSideEncryption \
  --output text 2>/dev/null || echo "NONE")

if [ "$ENCRYPTION" = "AES256" ]; then
  log_success "Criptografia AES256 habilitada"
else
  log_error "Criptografia não habilitada (encontrado: $ENCRYPTION)"
fi

log_info "Verificando versionamento..."
VERSIONS=$(aws s3api list-object-versions \
  --bucket "$S3_BACKUP_BUCKET" \
  --prefix "backups/$COMPRESSED_FILE" \
  --query 'length(Versions)' \
  --output text 2>/dev/null || echo "0")

if [ "$VERSIONS" -ge "1" ]; then
  log_success "Versionamento habilitado ($VERSIONS versões)"
else
  log_warning "Versionamento não habilitado ou não detectado"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Cleanup local
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🧹 Cleanup Local"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Removendo arquivos locais..."
rm -f "$BACKUP_DIR/$COMPRESSED_FILE"
rm -f "$BACKUP_DIR/$MANIFEST_FILE"
log_success "Arquivos locais removidos"

echo ""

# ──────────────────────────────────────────────────────────────────
# Resumo
# ──────────────────────────────────────────────────────────────────
echo "═══════════════════════════════════════════════════════════════"
echo "  📊 Resumo do Backup"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "  📁 Arquivo: $COMPRESSED_FILE"
echo "  💾 Tamanho: ${COMPRESSED_SIZE_MB}MB (comprimido)"
echo "  📦 Original: ${FILE_SIZE_MB}MB"
echo "  🗜️  Compressão: ${COMPRESSION_RATIO}x"
echo "  🔐 Checksum: $CHECKSUM"
echo "  ⏱️  Duração: ${DURATION}s"
echo "  ☁️  S3: s3://${S3_BACKUP_BUCKET}/backups/${COMPRESSED_FILE}"
echo "  🔒 Criptografia: $ENCRYPTION"
echo ""
echo -e "  ${GREEN}Aprovados:${NC} $PASSED"
echo -e "  ${RED}Falhados:${NC}  $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}✓ Backup concluído com sucesso!${NC}"
  exit 0
else
  echo -e "${RED}✗ Backup teve problemas. Verifique os logs acima.${NC}"
  exit 1
fi
