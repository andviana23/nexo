#!/bin/bash

# Script temporário para rodar backend sem módulos com erro (financial)
# Apenas para testar autenticação

cd /home/andrey/Projetos/barber-analytics-proV2/backend

echo "🔧 Compilando backend (modo auth-only)..."

# Tentar compilar
if go build -tags "auth_only" -o bin/api-auth ./cmd/api 2>&1 | grep -q "undefined:"; then
    echo "❌ Ainda há erros de compilação"
    echo ""
    echo "Vou iniciar o backend ignorando erros de módulos antigos..."
    echo ""
fi

# Source do .env
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Iniciar backend
echo "🚀 Iniciando backend na porta ${PORT:-8080}..."
echo ""
echo "📍 Endpoints disponíveis:"
echo "   - GET  http://localhost:${PORT:-8080}/health"
echo "   - POST http://localhost:${PORT:-8080}/api/v1/auth/login"
echo "   - POST http://localhost:${PORT:-8080}/api/v1/auth/refresh"
echo "   - GET  http://localhost:${PORT:-8080}/api/v1/auth/me"
echo "   - POST http://localhost:${PORT:-8080}/api/v1/auth/logout"
echo ""
echo "🔑 Credenciais de teste:"
echo "   Email: admin@teste.com"
echo "   Senha: Admin123!"
echo ""
echo "Press Ctrl+C to stop"
echo "======================================"
echo ""

go run ./cmd/api
