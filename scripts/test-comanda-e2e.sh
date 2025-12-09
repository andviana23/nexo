#!/bin/bash

# Script de teste E2E para Módulo de Comanda
# Testa fluxo completo: Appointment → Finalizar → Criar Comanda → Pagamento → Fechar → Done

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
TENANT_ID="${TENANT_ID:-00000000-0000-0000-0000-000000000001}"

echo "🧪 Teste E2E - Módulo de Comanda (Estilo Trinks)"
echo "================================================"
echo "Base URL: $BASE_URL"
echo "Tenant ID: $TENANT_ID"
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para log de sucesso
log_success() {
  echo -e "${GREEN}✅ $1${NC}"
}

# Função para log de erro
log_error() {
  echo -e "${RED}❌ $1${NC}"
}

# Função para log de info
log_info() {
  echo -e "${YELLOW}📝 $1${NC}"
}

# ============================================================================
# SETUP: Criar dados necessários (Professional, Customer, Service)
# ============================================================================

log_info "Fase 0: Setup inicial - Criando dados necessários"
echo ""

# Criar Professional
log_info "Criando profissional..."
PROFESSIONAL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/professionals" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "name": "João Barbeiro E2E",
    "email": "joao.e2e@test.com",
    "phone": "11999999999",
    "specialty": "Corte masculino",
    "commission_percentage": 40.0,
    "active": true
  }')

PROFESSIONAL_ID=$(echo "$PROFESSIONAL_RESPONSE" | jq -r '.id')
if [ "$PROFESSIONAL_ID" = "null" ] || [ -z "$PROFESSIONAL_ID" ]; then
  log_error "Erro ao criar profissional"
  exit 1
fi
log_success "Profissional criado: $PROFESSIONAL_ID"
echo ""

# Criar Customer
log_info "Criando cliente..."
CUSTOMER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/customers" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "nome": "Carlos Cliente E2E",
    "telefone": "11988888888",
    "email": "carlos.e2e@test.com"
  }')

CUSTOMER_ID=$(echo "$CUSTOMER_RESPONSE" | jq -r '.id')
if [ "$CUSTOMER_ID" = "null" ] || [ -z "$CUSTOMER_ID" ]; then
  log_error "Erro ao criar cliente"
  exit 1
fi
log_success "Cliente criado: $CUSTOMER_ID"
echo ""

# Criar Service
log_info "Criando serviço..."
SERVICE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/services" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "name": "Corte E2E",
    "description": "Corte de teste",
    "price": "50.00",
    "duration_minutes": 30,
    "active": true
  }')

SERVICE_ID=$(echo "$SERVICE_RESPONSE" | jq -r '.id')
if [ "$SERVICE_ID" = "null" ] || [ -z "$SERVICE_ID" ]; then
  log_error "Erro ao criar serviço"
  exit 1
fi
log_success "Serviço criado: $SERVICE_ID"
echo ""

# Criar Meio de Pagamento (PIX com taxa)
log_info "Criando meio de pagamento (PIX)..."
PAYMENT_METHOD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/meios-pagamento" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "nome": "PIX E2E",
    "tipo": "PIX",
    "taxa_percentual": 2.0,
    "taxa_fixa": "0.50",
    "ativo": true,
    "ordem_exibicao": 1
  }')

PAYMENT_METHOD_ID=$(echo "$PAYMENT_METHOD_RESPONSE" | jq -r '.id')
if [ "$PAYMENT_METHOD_ID" = "null" ] || [ -z "$PAYMENT_METHOD_ID" ]; then
  log_error "Erro ao criar meio de pagamento"
  exit 1
fi
log_success "Meio de pagamento criado: $PAYMENT_METHOD_ID"
echo ""

# ============================================================================
# TESTE 1: Criar Appointment
# ============================================================================

log_info "Teste 1: Criando agendamento"
echo ""

START_TIME=$(date -u -d '+1 hour' '+%Y-%m-%dT%H:%M:%SZ')
END_TIME=$(date -u -d '+1 hour 30 minutes' '+%Y-%m-%dT%H:%M:%SZ')

APPOINTMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/appointments" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d "{
    \"professional_id\": \"$PROFESSIONAL_ID\",
    \"customer_id\": \"$CUSTOMER_ID\",
    \"start_time\": \"$START_TIME\",
    \"service_ids\": [\"$SERVICE_ID\"],
    \"notes\": \"Teste E2E Comanda\"
  }")

echo "$APPOINTMENT_RESPONSE" | jq .
echo ""

APPOINTMENT_ID=$(echo "$APPOINTMENT_RESPONSE" | jq -r '.id')
if [ "$APPOINTMENT_ID" = "null" ] || [ -z "$APPOINTMENT_ID" ]; then
  log_error "Erro ao criar agendamento"
  exit 1
fi

APPOINTMENT_STATUS=$(echo "$APPOINTMENT_RESPONSE" | jq -r '.status')
if [ "$APPOINTMENT_STATUS" != "CREATED" ]; then
  log_error "Status inicial esperado: CREATED, recebido: $APPOINTMENT_STATUS"
  exit 1
fi

log_success "Agendamento criado: $APPOINTMENT_ID (status: $APPOINTMENT_STATUS)"
echo ""

# ============================================================================
# TESTE 2: Confirmar Appointment
# ============================================================================

log_info "Teste 2: Confirmando agendamento"
echo ""

CONFIRM_RESPONSE=$(curl -s -X PATCH "$BASE_URL/api/v1/appointments/$APPOINTMENT_ID/status" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "status": "CONFIRMED"
  }')

echo "$CONFIRM_RESPONSE" | jq .
echo ""

CONFIRMED_STATUS=$(echo "$CONFIRM_RESPONSE" | jq -r '.status')
if [ "$CONFIRMED_STATUS" != "CONFIRMED" ]; then
  log_error "Status esperado: CONFIRMED, recebido: $CONFIRMED_STATUS"
  exit 1
fi

log_success "Agendamento confirmado (status: $CONFIRMED_STATUS)"
echo ""

# ============================================================================
# TESTE 3: Check-in do Cliente
# ============================================================================

log_info "Teste 3: Check-in do cliente"
echo ""

CHECKIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/appointments/$APPOINTMENT_ID/check-in" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$CHECKIN_RESPONSE" | jq .
echo ""

CHECKIN_STATUS=$(echo "$CHECKIN_RESPONSE" | jq -r '.status')
if [ "$CHECKIN_STATUS" != "CHECKED_IN" ]; then
  log_error "Status esperado: CHECKED_IN, recebido: $CHECKIN_STATUS"
  exit 1
fi

log_success "Cliente fez check-in (status: $CHECKIN_STATUS)"
echo ""

# ============================================================================
# TESTE 4: Iniciar Atendimento
# ============================================================================

log_info "Teste 4: Iniciando atendimento"
echo ""

START_SERVICE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/appointments/$APPOINTMENT_ID/start" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$START_SERVICE_RESPONSE" | jq .
echo ""

IN_SERVICE_STATUS=$(echo "$START_SERVICE_RESPONSE" | jq -r '.status')
if [ "$IN_SERVICE_STATUS" != "IN_SERVICE" ]; then
  log_error "Status esperado: IN_SERVICE, recebido: $IN_SERVICE_STATUS"
  exit 1
fi

log_success "Atendimento iniciado (status: $IN_SERVICE_STATUS)"
echo ""

# ============================================================================
# TESTE 5: Finalizar Atendimento → AWAITING_PAYMENT
# ============================================================================

log_info "Teste 5: Finalizando atendimento (→ AWAITING_PAYMENT)"
echo ""

FINISH_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/appointments/$APPOINTMENT_ID/finish" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$FINISH_RESPONSE" | jq .
echo ""

AWAITING_STATUS=$(echo "$FINISH_RESPONSE" | jq -r '.status')
if [ "$AWAITING_STATUS" != "AWAITING_PAYMENT" ]; then
  log_error "Status esperado: AWAITING_PAYMENT, recebido: $AWAITING_STATUS"
  exit 1
fi

log_success "Atendimento finalizado (status: $AWAITING_STATUS)"
log_info "✨ Appointment agora aguarda pagamento - pronto para criar comanda"
echo ""

# ============================================================================
# TESTE 6: Criar Comanda a partir do Appointment
# ============================================================================

log_info "Teste 6: Criando comanda a partir do appointment"
echo ""

COMMAND_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/commands" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d "{
    \"appointment_id\": \"$APPOINTMENT_ID\",
    \"customer_id\": \"$CUSTOMER_ID\"
  }")

echo "$COMMAND_RESPONSE" | jq .
echo ""

COMMAND_ID=$(echo "$COMMAND_RESPONSE" | jq -r '.id')
if [ "$COMMAND_ID" = "null" ] || [ -z "$COMMAND_ID" ]; then
  log_error "Erro ao criar comanda"
  exit 1
fi

COMMAND_STATUS=$(echo "$COMMAND_RESPONSE" | jq -r '.status')
if [ "$COMMAND_STATUS" != "OPEN" ]; then
  log_error "Status da comanda esperado: OPEN, recebido: $COMMAND_STATUS"
  exit 1
fi

COMMAND_TOTAL=$(echo "$COMMAND_RESPONSE" | jq -r '.total')
log_success "Comanda criada: $COMMAND_ID (status: $COMMAND_STATUS, total: R$ $COMMAND_TOTAL)"
echo ""

# ============================================================================
# TESTE 7: Buscar Comanda Completa (com items e payments)
# ============================================================================

log_info "Teste 7: Buscando comanda completa"
echo ""

COMMAND_GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/commands/$COMMAND_ID" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$COMMAND_GET_RESPONSE" | jq .
echo ""

ITEMS_COUNT=$(echo "$COMMAND_GET_RESPONSE" | jq '.items | length')
if [ "$ITEMS_COUNT" -lt 1 ]; then
  log_error "Comanda deveria ter pelo menos 1 item (serviço do appointment)"
  exit 1
fi

log_success "Comanda possui $ITEMS_COUNT item(s)"
echo ""

# ============================================================================
# TESTE 8: Adicionar Pagamento (PIX com taxa)
# ============================================================================

log_info "Teste 8: Adicionando pagamento (PIX R$ 50,00)"
echo ""

PAYMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/commands/$COMMAND_ID/payments" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d "{
    \"meio_pagamento_id\": \"$PAYMENT_METHOD_ID\",
    \"valor_recebido\": \"50.00\"
  }")

echo "$PAYMENT_RESPONSE" | jq .
echo ""

PAYMENT_ID=$(echo "$PAYMENT_RESPONSE" | jq -r '.id')
if [ "$PAYMENT_ID" = "null" ] || [ -z "$PAYMENT_ID" ]; then
  log_error "Erro ao adicionar pagamento"
  exit 1
fi

VALOR_RECEBIDO=$(echo "$PAYMENT_RESPONSE" | jq -r '.valor_recebido')
TAXA_PERCENTUAL=$(echo "$PAYMENT_RESPONSE" | jq -r '.taxa_percentual')
TAXA_FIXA=$(echo "$PAYMENT_RESPONSE" | jq -r '.taxa_fixa')
VALOR_LIQUIDO=$(echo "$PAYMENT_RESPONSE" | jq -r '.valor_liquido')

log_success "Pagamento adicionado: $PAYMENT_ID"
log_info "   Valor recebido: R$ $VALOR_RECEBIDO"
log_info "   Taxa: $TAXA_PERCENTUAL% + R$ $TAXA_FIXA"
log_info "   Valor líquido: R$ $VALOR_LIQUIDO"
echo ""

# Validar cálculo de taxa
# Taxa = (50.00 * 2/100) + 0.50 = 1.00 + 0.50 = 1.50
# Líquido = 50.00 - 1.50 = 48.50

EXPECTED_LIQUIDO="48.50"
if [ "$VALOR_LIQUIDO" != "$EXPECTED_LIQUIDO" ]; then
  log_error "Cálculo de taxa incorreto! Esperado: $EXPECTED_LIQUIDO, Recebido: $VALOR_LIQUIDO"
  exit 1
fi

log_success "Cálculo de taxa validado: R$ 50.00 - R$ 1.50 = R$ 48.50 ✓"
echo ""

# ============================================================================
# TESTE 9: Fechar Comanda
# ============================================================================

log_info "Teste 9: Fechando comanda"
echo ""

CLOSE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/commands/$COMMAND_ID/close" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "deixar_troco_gorjeta": false,
    "deixar_saldo_divida": false,
    "observacoes": "Teste E2E concluído"
  }')

echo "$CLOSE_RESPONSE" | jq .
echo ""

CLOSED_STATUS=$(echo "$CLOSE_RESPONSE" | jq -r '.status')
if [ "$CLOSED_STATUS" != "CLOSED" ]; then
  log_error "Status da comanda esperado: CLOSED, recebido: $CLOSED_STATUS"
  exit 1
fi

FECHADO_EM=$(echo "$CLOSE_RESPONSE" | jq -r '.fechado_em')
if [ "$FECHADO_EM" = "null" ]; then
  log_error "Campo fechado_em deveria estar preenchido"
  exit 1
fi

log_success "Comanda fechada com sucesso (status: $CLOSED_STATUS)"
log_info "   Fechado em: $FECHADO_EM"
echo ""

# ============================================================================
# TESTE 10: Verificar que Appointment foi para DONE
# ============================================================================

log_info "Teste 10: Verificando status final do appointment"
echo ""

FINAL_APPOINTMENT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/appointments/$APPOINTMENT_ID" \
  -H "X-Tenant-ID: $TENANT_ID")

echo "$FINAL_APPOINTMENT_RESPONSE" | jq .
echo ""

FINAL_STATUS=$(echo "$FINAL_APPOINTMENT_RESPONSE" | jq -r '.status')
if [ "$FINAL_STATUS" != "DONE" ]; then
  log_error "Status final do appointment esperado: DONE, recebido: $FINAL_STATUS"
  exit 1
fi

log_success "Appointment marcado como DONE (status: $FINAL_STATUS)"
echo ""

# ============================================================================
# TESTE 11: Edição de Item (Adicional)
# ============================================================================

log_info "Teste 11: Testando edição de item da comanda"
echo ""

# Criar nova comanda para teste de edição
NEW_APPOINTMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/appointments" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d "{
    \"professional_id\": \"$PROFESSIONAL_ID\",
    \"customer_id\": \"$CUSTOMER_ID\",
    \"start_time\": \"$(date -u -d '+2 hour' '+%Y-%m-%dT%H:%M:%SZ')\",
    \"service_ids\": [\"$SERVICE_ID\"]
  }")

NEW_APPOINTMENT_ID=$(echo "$NEW_APPOINTMENT_RESPONSE" | jq -r '.id')

# Finalizar para AWAITING_PAYMENT
curl -s -X POST "$BASE_URL/api/v1/appointments/$NEW_APPOINTMENT_ID/finish" \
  -H "X-Tenant-ID: $TENANT_ID" > /dev/null

# Criar comanda
NEW_COMMAND_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/commands" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d "{
    \"appointment_id\": \"$NEW_APPOINTMENT_ID\",
    \"customer_id\": \"$CUSTOMER_ID\"
  }")

NEW_COMMAND_ID=$(echo "$NEW_COMMAND_RESPONSE" | jq -r '.id')

# Buscar item ID
ITEM_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/commands/$NEW_COMMAND_ID" \
  -H "X-Tenant-ID: $TENANT_ID")

ITEM_ID=$(echo "$ITEM_RESPONSE" | jq -r '.items[0].id')

# Atualizar item (aplicar desconto)
UPDATE_ITEM_RESPONSE=$(curl -s -X PATCH "$BASE_URL/api/v1/commands/$NEW_COMMAND_ID/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "desconto_valor": "5.00"
  }')

echo "$UPDATE_ITEM_RESPONSE" | jq .
echo ""

UPDATED_DESCONTO=$(echo "$UPDATE_ITEM_RESPONSE" | jq -r '.desconto_valor')
if [ "$UPDATED_DESCONTO" != "5.00" ]; then
  log_error "Desconto não aplicado corretamente"
  exit 1
fi

log_success "Item editado com sucesso (desconto: R$ $UPDATED_DESCONTO)"
echo ""

# ============================================================================
# CLEANUP
# ============================================================================

log_info "Limpando dados de teste..."
echo ""

# Deletar appointment, comanda, etc seria ideal aqui
# Por simplicidade, apenas logamos o sucesso

# ============================================================================
# RESUMO FINAL
# ============================================================================

echo ""
echo "═══════════════════════════════════════════════════════════"
log_success "TODOS OS TESTES PASSARAM! 🎉"
echo "═══════════════════════════════════════════════════════════"
echo ""
log_info "Fluxo completo validado:"
echo "  1. ✅ Appointment criado (CREATED)"
echo "  2. ✅ Appointment confirmado (CONFIRMED)"
echo "  3. ✅ Cliente check-in (CHECKED_IN)"
echo "  4. ✅ Atendimento iniciado (IN_SERVICE)"
echo "  5. ✅ Atendimento finalizado (AWAITING_PAYMENT)"
echo "  6. ✅ Comanda criada com items do appointment"
echo "  7. ✅ Pagamento adicionado com cálculo de taxas correto"
echo "  8. ✅ Comanda fechada (CLOSED)"
echo "  9. ✅ Appointment marcado como concluído (DONE)"
echo " 10. ✅ Edição de item da comanda funcional"
echo ""
log_info "Cálculos financeiros validados:"
echo "  • Taxa PIX: 2% + R$ 0,50 = R$ 1,50 sobre R$ 50,00"
echo "  • Valor líquido: R$ 48,50 ✓"
echo ""
log_success "Sistema de Comanda 100% funcional!"
echo ""
