#!/bin/bash

# ============================================================================
# 🧪 Barber Analytics Pro V2 — Test Categories API
# Testa endpoints de Categorias de Serviço
# ============================================================================

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="http://localhost:8080/api/v1"
EMAIL="andrey@tratodebarbados.com"
PASSWORD="@Aa30019258"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  🧪 Testando Categorias API — Barber Analytics Pro V2"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo

# ============================================================================
# 1. LOGIN
# ============================================================================

echo -e "${YELLOW}🔑 1. Login ($EMAIL)${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo -e "   ${RED}❌ Login falhou${NC}"
    echo "$LOGIN_RESPONSE"
    exit 1
fi

echo -e "   ${GREEN}✅ Login OK${NC}"
echo

# ============================================================================
# 2. CREATE CATEGORY
# ============================================================================

echo -e "${YELLOW}➕ 2. Criar Categoria${NC}"
CREATE_PAYLOAD='{
  "nome": "Categoria Teste API",
  "descricao": "Criada via script de teste",
  "cor": "#FF5733",
  "icone": "content_cut"
}'

CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/categorias-servicos" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "$CREATE_PAYLOAD")

CATEGORY_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

if [ "$CATEGORY_ID" == "null" ] || [ -z "$CATEGORY_ID" ]; then
    echo -e "   ${RED}❌ Falha ao criar categoria${NC}"
    echo "$CREATE_RESPONSE"
    exit 1
fi

echo -e "   ${GREEN}✅ Categoria criada (ID: $CATEGORY_ID)${NC}"
echo

# ============================================================================
# 3. LIST CATEGORIES
# ============================================================================

echo -e "${YELLOW}📋 3. Listar Categorias${NC}"
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/categorias-servicos" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

COUNT=$(echo "$LIST_RESPONSE" | jq '. | length')

if [ "$COUNT" -gt 0 ]; then
    echo -e "   ${GREEN}✅ Listagem OK ($COUNT categorias encontradas)${NC}"
else
    echo -e "   ${RED}❌ Listagem vazia ou falhou${NC}"
    echo "$LIST_RESPONSE"
    exit 1
fi

echo

# ============================================================================
# 4. UPDATE CATEGORY
# ============================================================================

echo -e "${YELLOW}✏️  4. Atualizar Categoria${NC}"
UPDATE_PAYLOAD='{
  "nome": "Categoria Teste API (Editada)",
  "descricao": "Editada via script",
  "cor": "#00FF00"
}'

UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/categorias-servicos/$CATEGORY_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "$UPDATE_PAYLOAD")

# Check status code or response
# Assuming 200 OK or 204 No Content
# curl output is the body. If empty, we might need -w "%{http_code}"

# Let's verify by getting it again
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/categorias-servicos/$CATEGORY_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

UPDATED_NAME=$(echo "$GET_RESPONSE" | jq -r '.nome')

if [ "$UPDATED_NAME" == "Categoria Teste API (Editada)" ]; then
    echo -e "   ${GREEN}✅ Atualização confirmada${NC}"
else
    echo -e "   ${RED}❌ Atualização falhou${NC}"
    echo "Esperado: Categoria Teste API (Editada)"
    echo "Obtido: $UPDATED_NAME"
    exit 1
fi

echo

# ============================================================================
# 5. DELETE CATEGORY
# ============================================================================

echo -e "${YELLOW}🗑️  5. Deletar Categoria${NC}"
DELETE_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/categorias-servicos/$CATEGORY_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if [ "$DELETE_RESPONSE" -eq 204 ] || [ "$DELETE_RESPONSE" -eq 200 ]; then
    echo -e "   ${GREEN}✅ Deleção OK (HTTP $DELETE_RESPONSE)${NC}"
else
    echo -e "   ${RED}❌ Deleção falhou (HTTP $DELETE_RESPONSE)${NC}"
    exit 1
fi

# Verify deletion
VERIFY_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/categorias-servicos/$CATEGORY_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if [ "$VERIFY_RESPONSE" -eq 404 ]; then
    echo -e "   ${GREEN}✅ Verificação: Categoria não encontrada (404)${NC}"
else
    echo -e "   ${RED}❌ Verificação falhou: Categoria ainda existe (HTTP $VERIFY_RESPONSE)${NC}"
    exit 1
fi

echo

echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ TESTES DE CATEGORIA CONCLUÍDOS COM SUCESSO!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════${NC}"
echo
