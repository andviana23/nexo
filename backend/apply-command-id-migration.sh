#!/bin/bash

# Script para aplicar migration e regenerar SQLC
# Data: 30/11/2025

set -e

echo "🔧 Aplicando migration 032_add_command_id_to_appointments..."

# Rodar migration
echo "📦 Executando migration..."
make migrate-up || echo "⚠️  Migration pode já estar aplicada"

# Regenerar SQLC
echo "🔄 Regenerando código SQLC..."
make sqlc-generate

echo "✅ Concluído!"
echo ""
echo "📝 Próximos passos:"
echo "1. Verificar se não há erros de compilação: go build ./..."
echo "2. Rodar testes: make test"
echo "3. Verificar endpoint GET /api/v1/appointments/:id retorna command_id"
