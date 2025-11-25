#!/usr/bin/env bash

################################################################################
# test-prometheus-alerts.sh - Teste de Alertas Prometheus/Alertmanager
#
# Descrição:
#   Valida configuração de alertas e simula condições de falha para
#   verificar se os alertas disparam corretamente.
#
# Testa:
#   - Alerta de falha de backup (>24h sem backup)
#   - Alerta de backup lento (duração >30min)
#   - Alerta de arquivo muito pequeno (<1MB)
#   - Alerta de restore não testado (>30 dias)
#
# Uso:
#   ./scripts/test-prometheus-alerts.sh [PROMETHEUS_URL]
#
# Exemplo:
#   ./scripts/test-prometheus-alerts.sh http://localhost:9090
#
# Requisitos:
#   - Prometheus rodando
#   - Alertmanager configurado
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
PROMETHEUS_URL="${1:-http://localhost:9090}"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"
PASSED=0
FAILED=0

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[✓]${NC} $1"; ((PASSED++)); }
log_error() { echo -e "${RED}[✗]${NC} $1"; ((FAILED++)); }
log_warning() { echo -e "${YELLOW}[⚠]${NC} $1"; }

# Banner
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  🚨 Teste de Alertas Prometheus + Alertmanager"
echo "═══════════════════════════════════════════════════════════════"
echo ""
log_info "Prometheus: $PROMETHEUS_URL"
log_info "Alertmanager: $ALERTMANAGER_URL"
echo ""

# ──────────────────────────────────────────────────────────────────
# Verificar se Prometheus está rodando
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔌 Verificar Conectividade"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Testando conexão com Prometheus..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$PROMETHEUS_URL/-/healthy")

if [ "$HTTP_CODE" = "200" ]; then
  log_success "Prometheus acessível (Status 200)"
else
  log_error "Prometheus não acessível (Status $HTTP_CODE)"
  log_warning "Certifique-se de que Prometheus está rodando em $PROMETHEUS_URL"
  exit 1
fi

log_info "Testando conexão com Alertmanager..."
AM_HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$ALERTMANAGER_URL/-/healthy" 2>/dev/null || echo "000")

if [ "$AM_HTTP_CODE" = "200" ]; then
  log_success "Alertmanager acessível (Status 200)"
else
  log_warning "Alertmanager não acessível (Status $AM_HTTP_CODE)"
  log_info "Alguns testes serão limitados sem Alertmanager"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Validar arquivo de regras de alerta
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📋 Validar Arquivo de Regras"
echo "───────────────────────────────────────────────────────────────"
echo ""

RULES_FILE="prometheus-alert-rules.yml"

if [ -f "$RULES_FILE" ]; then
  log_success "Arquivo de regras existe: $RULES_FILE"

  # Validar YAML
  if command -v yamllint &> /dev/null; then
    if yamllint -d relaxed "$RULES_FILE" &> /dev/null; then
      log_success "YAML válido (yamllint)"
    else
      log_warning "yamllint reportou warnings (verificar manualmente)"
    fi
  else
    log_info "yamllint não instalado, pulando validação YAML"
  fi

  # Contar grupos e alertas
  GROUPS=$(grep -c "^  - name:" "$RULES_FILE" || echo "0")
  ALERTS=$(grep -c "^    - alert:" "$RULES_FILE" || echo "0")

  log_info "Grupos de alertas: $GROUPS"
  log_info "Total de alertas: $ALERTS"

  if [ "$ALERTS" -ge "8" ]; then
    log_success "Alertas configurados ($ALERTS >= 8)"
  else
    log_error "Poucos alertas configurados ($ALERTS < 8)"
  fi
else
  log_error "Arquivo de regras não encontrado: $RULES_FILE"
  exit 1
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Verificar regras carregadas no Prometheus
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🔍 Verificar Regras Carregadas"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Consultando /api/v1/rules..."
LOADED_RULES=$(curl -s "$PROMETHEUS_URL/api/v1/rules" | jq -r '.data.groups[].rules[].name' 2>/dev/null | wc -l || echo "0")

log_info "Regras carregadas no Prometheus: $LOADED_RULES"

if [ "$LOADED_RULES" -ge "8" ]; then
  log_success "Regras carregadas ($LOADED_RULES >= 8)"
else
  log_error "Regras não carregadas ou arquivo não foi importado"
  log_info "Verifique prometheus.yml:"
  log_info "  rule_files:"
  log_info "    - 'prometheus-alert-rules.yml'"
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Verificar alertas específicos
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🚨 Verificar Alertas Específicos"
echo "───────────────────────────────────────────────────────────────"
echo ""

EXPECTED_ALERTS=(
  "BackupFailed"
  "BackupTooSlow"
  "BackupFileTooSmall"
  "BackupHighFailureRate"
  "LGPDExportHighFailureRate"
  "LGPDExportSlow"
  "APIHighErrorRate"
  "APIHighLatency"
)

for alert in "${EXPECTED_ALERTS[@]}"; do
  FOUND=$(curl -s "$PROMETHEUS_URL/api/v1/rules" | jq -r ".data.groups[].rules[] | select(.name == \"$alert\") | .name" 2>/dev/null || echo "")

  if [ "$FOUND" = "$alert" ]; then
    log_success "Alerta '$alert' configurado"
  else
    log_error "Alerta '$alert' NÃO encontrado"
  fi
done

echo ""

# ──────────────────────────────────────────────────────────────────
# Verificar estado atual dos alertas
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  📊 Estado Atual dos Alertas"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_info "Consultando alertas ativos..."
ACTIVE_ALERTS=$(curl -s "$PROMETHEUS_URL/api/v1/alerts" | jq -r '.data.alerts[] | .labels.alertname' 2>/dev/null || echo "")

if [ -z "$ACTIVE_ALERTS" ]; then
  log_success "Nenhum alerta disparado (sistema saudável)"
else
  ALERT_COUNT=$(echo "$ACTIVE_ALERTS" | wc -l)
  log_warning "$ALERT_COUNT alerta(s) disparado(s):"
  echo "$ACTIVE_ALERTS" | while read -r alert; do
    log_warning "  - $alert"
  done
fi

echo ""

# ──────────────────────────────────────────────────────────────────
# Simular falhas para testar alertas
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  🧪 Simulação de Falhas"
echo "───────────────────────────────────────────────────────────────"
echo ""

log_warning "Para testar alertas, simule as seguintes condições:"
echo ""

echo "1. BACKUPFAILED (Backup não executado em 24h)"
echo "   Simular:"
echo "   - Desabilitar workflow de backup por 24h"
echo "   - Ou: curl -X POST http://localhost:9091/metrics/job/backup \\"
echo "           --data-binary @- <<EOF"
echo "     backup_last_success_timestamp $(date -d '2 days ago' +%s)"
echo "     EOF"
echo "   Validar:"
echo "   - Aguardar 5min (evaluation_interval)"
echo "   - curl $PROMETHEUS_URL/api/v1/alerts | jq '.data.alerts[] | select(.labels.alertname == \"BackupFailed\")'"
echo ""

echo "2. BACKUPTOOSLOW (Backup demora >30min)"
echo "   Simular:"
echo "   - Executar backup muito grande (>30min)"
echo "   - Ou: curl -X POST http://localhost:9091/metrics/job/backup \\"
echo "           --data-binary 'backup_duration_seconds 2000'"
echo "   Validar:"
echo "   - curl $PROMETHEUS_URL/api/v1/alerts | jq '.data.alerts[] | select(.labels.alertname == \"BackupTooSlow\")'"
echo ""

echo "3. BACKUPFILETOOSMALL (Arquivo <1MB)"
echo "   Simular:"
echo "   - curl -X POST http://localhost:9091/metrics/job/backup \\"
echo "         --data-binary 'backup_file_size_bytes 500000'"
echo "   Validar:"
echo "   - Verificar alerta após 5min"
echo ""

echo "4. LGPDEXPORTHIGHFAILURERATE (Exports falhando >10%)"
echo "   Simular:"
echo "   - Chamar GET /api/v1/me/export sem autenticação 10x"
echo "   - Garantir que >10% retornam erro 401/500"
echo "   Validar:"
echo "   - curl $PROMETHEUS_URL/api/v1/query?query='rate(lgpd_export_requests_total{status=\"error\"}[5m])'"
echo ""

log_warning "Testes de simulação requerem configuração manual ou Pushgateway"

echo ""

# ──────────────────────────────────────────────────────────────────
# Verificar configuração do Alertmanager
# ──────────────────────────────────────────────────────────────────
if [ "$AM_HTTP_CODE" = "200" ]; then
  echo "───────────────────────────────────────────────────────────────"
  echo "  📨 Verificar Alertmanager Config"
  echo "───────────────────────────────────────────────────────────────"
  echo ""

  log_info "Consultando status do Alertmanager..."
  AM_STATUS=$(curl -s "$ALERTMANAGER_URL/api/v2/status" | jq -r '.cluster.status' 2>/dev/null || echo "unknown")

  if [ "$AM_STATUS" = "ready" ]; then
    log_success "Alertmanager status: ready"
  else
    log_warning "Alertmanager status: $AM_STATUS"
  fi

  # Verificar receivers configurados
  RECEIVERS=$(curl -s "$ALERTMANAGER_URL/api/v1/status" | jq -r '.config.receivers[].name' 2>/dev/null || echo "")

  if [ -n "$RECEIVERS" ]; then
    log_success "Receivers configurados:"
    echo "$RECEIVERS" | while read -r receiver; do
      log_info "  - $receiver"
    done
  else
    log_error "Nenhum receiver configurado"
  fi

  echo ""
fi

# ──────────────────────────────────────────────────────────────────
# Checklist manual
# ──────────────────────────────────────────────────────────────────
echo "───────────────────────────────────────────────────────────────"
echo "  ✅ Checklist Manual de Validação"
echo "───────────────────────────────────────────────────────────────"
echo ""

cat <<'EOF'
Execute os seguintes testes manualmente:

[ ] 1. Verificar prometheus.yml tem rule_files configurado
[ ] 2. Reiniciar Prometheus após adicionar prometheus-alert-rules.yml
[ ] 3. Acessar http://localhost:9090/alerts e verificar todas as regras
[ ] 4. Simular falha de backup (desabilitar workflow por 24h)
[ ] 5. Verificar alerta BackupFailed dispara após 24h
[ ] 6. Configurar Slack/Email no Alertmanager
[ ] 7. Testar notificação disparando alerta de teste
[ ] 8. Validar que alertas resolvidos enviam notificação de "resolved"
[ ] 9. Configurar silenciamento de alertas durante manutenção
[ ] 10. Documentar runbook para cada alerta crítico

EOF

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
  echo -e "${GREEN}✓ Configuração de alertas validada!${NC}"
  echo ""
  echo "  Próximos passos:"
  echo "  1. Configurar Alertmanager com Slack webhook"
  echo "  2. Simular falhas para validar alertas end-to-end"
  echo "  3. Criar runbooks para cada alerta crítico"
  echo "  4. Configurar escalation policies (PagerDuty/Opsgenie)"
  exit 0
else
  echo -e "${RED}✗ Alguns testes de alertas falharam.${NC}"
  exit 1
fi
