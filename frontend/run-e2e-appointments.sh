#!/bin/bash

# NEXO - Testes E2E do Módulo de Agendamento
# Execute este script para rodar os testes end-to-end do módulo de agendamento

set -e

echo "🧪 NEXO - Testes E2E: Módulo de Agendamento"
echo "=========================================="
echo ""

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar se estamos no diretório correto
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ Erro: Execute este script no diretório frontend/${NC}"
    exit 1
fi

# Verificar se o Playwright está instalado
if ! command -v npx &> /dev/null; then
    echo -e "${RED}❌ Erro: Node.js não encontrado${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Pré-requisitos:${NC}"
echo "  1. Backend rodando em http://localhost:8080"
echo "  2. Frontend rodando em http://localhost:3000"
echo "  3. Banco de dados com dados de teste"
echo ""

read -p "Os serviços estão rodando? (s/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Ss]$ ]]; then
    echo -e "${YELLOW}ℹ️  Inicie os serviços antes de executar os testes${NC}"
    echo ""
    echo "Backend:"
    echo "  cd ../backend && make dev"
    echo ""
    echo "Frontend:"
    echo "  pnpm dev"
    echo ""
    exit 1
fi

echo ""
echo -e "${GREEN}▶️  Instalando dependências do Playwright (se necessário)...${NC}"
npx playwright install --with-deps chromium firefox

echo ""
echo -e "${GREEN}▶️  Executando testes E2E do módulo de agendamento...${NC}"
echo ""

# Executar apenas os testes de agendamento
npx playwright test tests/e2e/appointments.spec.ts --project=chromium

EXIT_CODE=$?

echo ""
if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ Todos os testes E2E passaram!${NC}"
    echo ""
    echo -e "${YELLOW}📊 Para ver o relatório HTML:${NC}"
    echo "  npx playwright show-report"
else
    echo -e "${RED}❌ Alguns testes falharam${NC}"
    echo ""
    echo -e "${YELLOW}🔍 Para ver detalhes:${NC}"
    echo "  npx playwright show-report"
    echo ""
    echo -e "${YELLOW}💡 Dicas de debug:${NC}"
    echo "  - Execute com modo UI: npx playwright test --ui"
    echo "  - Execute com debug: npx playwright test --debug"
    echo "  - Veja screenshots: ls -la test-results/"
fi

exit $EXIT_CODE
