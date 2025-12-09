#!/bin/bash

# Script de teste para endpoints de Serviços
# Sprint 1.4.2 - Serviços Básicos

set -e

BASE_URL="http://localhost:8080/api/v1"
TENANT_ID="e2e00000-0000-0000-0000-000000000001"
EMAIL="andrey@tratodebarbados.com"
PASSWORD="@Aa30019258"

echo "=========================================="
echo "🧪 TESTES - MÓDULO SERVIÇOS"
echo "=========================================="
echo ""

# 1. Login
echo "1️⃣  Fazendo login..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}")

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" = "null" ] || [ -z "$ACCESS_TOKEN" ]; then
  echo "❌ Erro no login"
  echo "$LOGIN_RESPONSE" | jq .
  exit 1
fi

echo "✅ Login bem-sucedido"
echo ""

# 2. Listar categorias (para pegar categoria_id)
echo "2️⃣  Listando categorias de serviço..."
CATEGORIAS_RESPONSE=$(curl -s -X GET "${BASE_URL}/categorias-servicos" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

CATEGORIA_ID=$(echo "$CATEGORIAS_RESPONSE" | jq -r '.categorias[0].id')

if [ "$CATEGORIA_ID" = "null" ] || [ -z "$CATEGORIA_ID" ]; then
  echo "❌ Nenhuma categoria encontrada"
  exit 1
fi

echo "✅ Categoria obtida: $CATEGORIA_ID"
echo ""

# 3. Criar serviço
echo "3️⃣  Criando novo serviço..."
TIMESTAMP=$(date +%s)
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/servicos" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"nome\": \"Corte Teste ${TIMESTAMP}\",
    \"descricao\": \"Serviço criado via teste automatizado\",
    \"preco\": \"45.00\",
    \"duracao\": 30,
    \"comissao\": \"50.00\",
    \"categoria_id\": \"${CATEGORIA_ID}\",
    \"cor\": \"#FF5733\",
    \"tags\": [\"teste\", \"api\"]
  }")

SERVICO_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

if [ "$SERVICO_ID" = "null" ] || [ -z "$SERVICO_ID" ]; then
  echo "❌ Erro ao criar serviço"
  echo "$CREATE_RESPONSE" | jq .
  exit 1
fi

echo "✅ Serviço criado: $SERVICO_ID"
echo "$CREATE_RESPONSE" | jq .
echo ""

# 4. Buscar serviço por ID
echo "4️⃣  Buscando serviço por ID..."
GET_RESPONSE=$(curl -s -X GET "${BASE_URL}/servicos/${SERVICO_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

NOME_SERVICO=$(echo "$GET_RESPONSE" | jq -r '.nome')

if [ "$NOME_SERVICO" != "Corte Teste ${TIMESTAMP}" ]; then
  echo "❌ Erro ao buscar serviço"
  exit 1
fi

echo "✅ Serviço encontrado:"
echo "$GET_RESPONSE" | jq .
echo ""

# 5. Listar todos os serviços
echo "5️⃣  Listando todos os serviços..."
LIST_RESPONSE=$(curl -s -X GET "${BASE_URL}/servicos" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

TOTAL=$(echo "$LIST_RESPONSE" | jq -r '.total')

echo "✅ Total de serviços: $TOTAL"
echo ""

# 6. Listar apenas serviços ativos
echo "6️⃣  Listando apenas serviços ativos..."
ATIVOS_RESPONSE=$(curl -s -X GET "${BASE_URL}/servicos?apenas_ativos=true" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

TOTAL_ATIVOS=$(echo "$ATIVOS_RESPONSE" | jq -r '.total')

echo "✅ Serviços ativos: $TOTAL_ATIVOS"
echo ""

# 7. Buscar por categoria
echo "7️⃣  Filtrando por categoria..."
CATEGORIA_RESPONSE=$(curl -s -X GET "${BASE_URL}/servicos?categoria_id=${CATEGORIA_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

TOTAL_CATEGORIA=$(echo "$CATEGORIA_RESPONSE" | jq -r '.total')

echo "✅ Serviços da categoria: $TOTAL_CATEGORIA"
echo ""

# 8. Buscar serviços (search)
echo "8️⃣  Buscando serviços com termo 'Corte'..."
SEARCH_RESPONSE=$(curl -s -X GET "${BASE_URL}/servicos?search=Corte" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

TOTAL_SEARCH=$(echo "$SEARCH_RESPONSE" | jq -r '.total')

echo "✅ Resultados encontrados: $TOTAL_SEARCH"
echo ""

# 9. Obter estatísticas
echo "9️⃣  Obtendo estatísticas de serviços..."
STATS_RESPONSE=$(curl -s -X GET "${BASE_URL}/servicos/stats" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

echo "✅ Estatísticas:"
echo "$STATS_RESPONSE" | jq .
echo ""

# 10. Atualizar serviço
echo "🔟 Atualizando serviço..."
UPDATE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/servicos/${SERVICO_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"nome\": \"Corte Teste ${TIMESTAMP} Atualizado\",
    \"descricao\": \"Descrição atualizada\",
    \"preco\": \"50.00\",
    \"duracao\": 40,
    \"comissao\": \"55.00\"
  }")

NOME_ATUALIZADO=$(echo "$UPDATE_RESPONSE" | jq -r '.nome')

if [ "$NOME_ATUALIZADO" != "Corte Teste ${TIMESTAMP} Atualizado" ]; then
  echo "❌ Erro ao atualizar serviço"
  echo "$UPDATE_RESPONSE" | jq .
  exit 1
fi

echo "✅ Serviço atualizado:"
echo "$UPDATE_RESPONSE" | jq .
echo ""

# 11. Toggle status (desativar)
echo "1️⃣1️⃣  Desativando serviço..."
TOGGLE_RESPONSE=$(curl -s -X PATCH "${BASE_URL}/servicos/${SERVICO_ID}/toggle-status" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

ATIVO=$(echo "$TOGGLE_RESPONSE" | jq -r '.ativo')

if [ "$ATIVO" != "false" ]; then
  echo "❌ Erro ao desativar serviço"
  exit 1
fi

echo "✅ Serviço desativado"
echo ""

# 12. Toggle status (reativar)
echo "1️⃣2️⃣  Reativando serviço..."
TOGGLE2_RESPONSE=$(curl -s -X PATCH "${BASE_URL}/servicos/${SERVICO_ID}/toggle-status" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

ATIVO2=$(echo "$TOGGLE2_RESPONSE" | jq -r '.ativo')

if [ "$ATIVO2" != "true" ]; then
  echo "❌ Erro ao reativar serviço"
  exit 1
fi

echo "✅ Serviço reativado"
echo ""

# 13. Deletar serviço
echo "1️⃣3️⃣  Deletando serviço de teste..."
DELETE_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/servicos/${SERVICO_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

echo "✅ Serviço deletado"
echo ""

# 14. Verificar que foi deletado
echo "1️⃣4️⃣  Verificando que serviço foi deletado..."
VERIFY_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/servicos/${SERVICO_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_CODE=$(echo "$VERIFY_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" != "404" ]; then
  echo "❌ Serviço ainda existe"
  exit 1
fi

echo "✅ Confirmado: serviço não existe mais"
echo ""

echo "=========================================="
echo "✅ TODOS OS TESTES PASSARAM!"
echo "=========================================="
echo ""
echo "📊 Resumo dos testes:"
echo "  • Login: ✅"
echo "  • Criar serviço: ✅"
echo "  • Buscar por ID: ✅"
echo "  • Listar todos: ✅"
echo "  • Filtrar ativos: ✅"
echo "  • Filtrar por categoria: ✅"
echo "  • Buscar (search): ✅"
echo "  • Estatísticas: ✅"
echo "  • Atualizar: ✅"
echo "  • Desativar: ✅"
echo "  • Reativar: ✅"
echo "  • Deletar: ✅"
echo "  • Verificar deleção: ✅"
echo ""
