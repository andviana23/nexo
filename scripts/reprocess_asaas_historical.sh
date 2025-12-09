#!/bin/bash
# ============================================================
# Script de Reprocessamento Histórico Asaas → NEXO
# ============================================================
# 
# Este script busca pagamentos confirmados do Asaas que não
# possuem conta a receber correspondente no NEXO e cria os
# registros necessários.
#
# Uso: ./reprocess_asaas_historical.sh [TENANT_ID] [DATA_INICIO] [DATA_FIM]
#
# Exemplo:
#   ./reprocess_asaas_historical.sh e2e00000-0000-0000-0000-000000000001 2024-01-01 2024-12-31
#
# Variáveis de ambiente necessárias:
#   - DATABASE_URL: URL de conexão PostgreSQL
#   - ASAAS_API_KEY: Chave de API do Asaas (opcional, para sync com API)
# ============================================================

set -euo pipefail

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funções de log
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Verificar parâmetros
TENANT_ID="${1:-}"
DATA_INICIO="${2:-$(date -d '30 days ago' '+%Y-%m-%d')}"
DATA_FIM="${3:-$(date '+%Y-%m-%d')}"

if [ -z "$TENANT_ID" ]; then
    log_error "TENANT_ID não informado"
    echo ""
    echo "Uso: $0 <TENANT_ID> [DATA_INICIO] [DATA_FIM]"
    echo ""
    echo "Exemplo:"
    echo "  $0 e2e00000-0000-0000-0000-000000000001 2024-01-01 2024-12-31"
    exit 1
fi

# Verificar DATABASE_URL
if [ -z "${DATABASE_URL:-}" ]; then
    log_error "DATABASE_URL não configurada"
    exit 1
fi

log_info "=============================================="
log_info "Reprocessamento Histórico Asaas → NEXO"
log_info "=============================================="
log_info "Tenant:     $TENANT_ID"
log_info "Período:    $DATA_INICIO a $DATA_FIM"
log_info "=============================================="
echo ""

# 1. Listar pagamentos CONFIRMED/RECEIVED sem conta a receber
log_info "Buscando pagamentos sem conta a receber..."

QUERY_ORPHAN_PAYMENTS=$(cat <<EOF
SELECT 
    sp.id,
    sp.asaas_payment_id,
    sp.subscription_id,
    sp.valor,
    sp.status,
    sp.confirmed_date,
    sp.credit_date,
    s.plano_id,
    s.cliente_id,
    p.nome as plano_nome
FROM subscription_payments sp
JOIN subscriptions s ON sp.subscription_id = s.id AND sp.tenant_id = s.tenant_id
JOIN plans p ON s.plano_id = p.id
WHERE sp.tenant_id = '$TENANT_ID'
  AND sp.status IN ('CONFIRMADO', 'RECEBIDO')
  AND sp.asaas_payment_id IS NOT NULL
  AND sp.created_at >= '$DATA_INICIO'
  AND sp.created_at <= '$DATA_FIM'
  AND NOT EXISTS (
      SELECT 1 FROM contas_a_receber cr
      WHERE cr.asaas_payment_id = sp.asaas_payment_id
        AND cr.tenant_id = sp.tenant_id
  )
ORDER BY sp.confirmed_date ASC;
EOF
)

ORPHAN_PAYMENTS=$(psql "$DATABASE_URL" -t -c "$QUERY_ORPHAN_PAYMENTS" 2>/dev/null || echo "")

if [ -z "$ORPHAN_PAYMENTS" ]; then
    log_success "Nenhum pagamento órfão encontrado!"
    exit 0
fi

# Contar pagamentos
ORPHAN_COUNT=$(echo "$ORPHAN_PAYMENTS" | grep -v '^$' | wc -l)
log_warn "Encontrados $ORPHAN_COUNT pagamentos sem conta a receber"
echo ""

# 2. Confirmar com usuário
read -p "Deseja criar as contas a receber? (s/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Ss]$ ]]; then
    log_info "Operação cancelada pelo usuário"
    exit 0
fi

# 3. Processar cada pagamento
log_info "Processando pagamentos..."
SUCCESS_COUNT=0
ERROR_COUNT=0

while IFS='|' read -r id asaas_payment_id subscription_id valor status confirmed_date credit_date plano_id cliente_id plano_nome; do
    # Limpar espaços
    id=$(echo "$id" | xargs)
    asaas_payment_id=$(echo "$asaas_payment_id" | xargs)
    subscription_id=$(echo "$subscription_id" | xargs)
    valor=$(echo "$valor" | xargs)
    status=$(echo "$status" | xargs)
    confirmed_date=$(echo "$confirmed_date" | xargs)
    credit_date=$(echo "$credit_date" | xargs)
    plano_nome=$(echo "$plano_nome" | xargs)
    
    if [ -z "$asaas_payment_id" ]; then
        continue
    fi
    
    # Determinar competência (mês do confirmed_date ou do due_date)
    if [ -n "$confirmed_date" ]; then
        COMPETENCIA=$(echo "$confirmed_date" | cut -d'-' -f1-2)
    else
        COMPETENCIA=$(date '+%Y-%m')
    fi
    
    # Determinar status da conta
    if [ "$status" = "RECEBIDO" ]; then
        CONTA_STATUS="RECEBIDO"
        DATA_RECEBIMENTO="'$credit_date'::date"
        VALOR_PAGO="$valor"
    else
        CONTA_STATUS="CONFIRMADO"
        DATA_RECEBIMENTO="NULL"
        VALOR_PAGO="0"
    fi
    
    # Descrição
    DESCRICAO="Assinatura: $plano_nome (Reprocessado)"
    
    # Inserir conta a receber
    INSERT_QUERY=$(cat <<EOF
INSERT INTO contas_a_receber (
    tenant_id,
    origem,
    assinatura_id,
    subscription_id,
    descricao,
    valor,
    valor_pago,
    data_vencimento,
    data_recebimento,
    status,
    observacoes,
    asaas_payment_id,
    competencia_mes,
    confirmed_at,
    received_at
) VALUES (
    '$TENANT_ID',
    'ASSINATURA',
    NULL,
    '$subscription_id',
    '$DESCRICAO',
    $valor,
    $VALOR_PAGO,
    COALESCE('$confirmed_date'::date, CURRENT_DATE),
    $DATA_RECEBIMENTO,
    '$CONTA_STATUS',
    'Conta criada via reprocessamento histórico',
    '$asaas_payment_id',
    '$COMPETENCIA',
    $([ -n "$confirmed_date" ] && echo "'$confirmed_date'::timestamp" || echo "NULL"),
    $([ "$status" = "RECEBIDO" ] && [ -n "$credit_date" ] && echo "'$credit_date'::timestamp" || echo "NULL")
)
ON CONFLICT (tenant_id, asaas_payment_id) WHERE asaas_payment_id IS NOT NULL 
DO UPDATE SET 
    status = EXCLUDED.status,
    valor_pago = EXCLUDED.valor_pago,
    data_recebimento = EXCLUDED.data_recebimento,
    received_at = EXCLUDED.received_at
RETURNING id;
EOF
)

    if RESULT=$(psql "$DATABASE_URL" -t -c "$INSERT_QUERY" 2>&1); then
        ((SUCCESS_COUNT++))
        log_success "Conta criada: $asaas_payment_id → $RESULT"
    else
        ((ERROR_COUNT++))
        log_error "Falha ao criar conta para $asaas_payment_id: $RESULT"
    fi

done <<< "$ORPHAN_PAYMENTS"

echo ""
log_info "=============================================="
log_info "Resumo do Reprocessamento"
log_info "=============================================="
log_success "Sucesso: $SUCCESS_COUNT contas criadas"
if [ $ERROR_COUNT -gt 0 ]; then
    log_error "Erros:   $ERROR_COUNT falhas"
fi
log_info "=============================================="

# 4. Gerar log de conciliação
LOG_QUERY=$(cat <<EOF
INSERT INTO asaas_reconciliation_logs (
    tenant_id,
    period_start,
    period_end,
    total_asaas,
    total_nexo,
    divergences,
    auto_fixed,
    pending_review,
    details
) VALUES (
    '$TENANT_ID',
    '$DATA_INICIO'::date,
    '$DATA_FIM'::date,
    $ORPHAN_COUNT,
    $SUCCESS_COUNT,
    $ERROR_COUNT,
    $SUCCESS_COUNT,
    $ERROR_COUNT,
    '{"type": "historical_reprocess", "script": "reprocess_asaas_historical.sh"}'::jsonb
)
RETURNING id;
EOF
)

if LOG_ID=$(psql "$DATABASE_URL" -t -c "$LOG_QUERY" 2>/dev/null); then
    log_info "Log de conciliação criado: $LOG_ID"
fi

exit $ERROR_COUNT
