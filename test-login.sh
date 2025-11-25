#!/bin/bash

# Script de teste do endpoint de login
# Testa autenticação com o usuário admin@teste.com

BASE_URL="http://localhost:8080/api/v1"

echo "======================================"
echo "Teste de Autenticação - VALTARIS v1.0"
echo "======================================"
echo ""

# Cores
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. Health Check
echo "1️⃣  Health Check..."
HEALTH=$(curl -s ${BASE_URL%/api/v1}/health)
if echo "$HEALTH" | grep -q "ok"; then
    echo -e "${GREEN}✓ Backend rodando${NC}"
else
    echo -e "${RED}✗ Backend não está respondendo${NC}"
    exit 1
fi
echo ""

# 2. Login
echo "2️⃣  Testando Login (admin@teste.com / Admin123!)..."
LOGIN_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -X POST ${BASE_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@teste.com",
    "password": "Admin123!"
  }')

HTTP_STATUS=$(echo "$LOGIN_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$LOGIN_RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" = "200" ]; then
    echo -e "${GREEN}✓ Login bem-sucedido (200 OK)${NC}"
    echo ""
    echo "📦 Response Body:"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"

    # Extrair access token
    ACCESS_TOKEN=$(echo "$BODY" | jq -r '.access_token' 2>/dev/null)

    if [ "$ACCESS_TOKEN" != "null" ] && [ -n "$ACCESS_TOKEN" ]; then
        echo ""
        echo "🔑 Access Token extraído com sucesso"
        echo ""

        # 3. Testar endpoint /me com o token
        echo "3️⃣  Testando endpoint /auth/me (protegido)..."
        ME_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
          -X GET ${BASE_URL}/auth/me \
          -H "Authorization: Bearer $ACCESS_TOKEN")

        ME_HTTP_STATUS=$(echo "$ME_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
        ME_BODY=$(echo "$ME_RESPONSE" | sed '/HTTP_STATUS/d')

        if [ "$ME_HTTP_STATUS" = "200" ]; then
            echo -e "${GREEN}✓ Endpoint /auth/me OK (200)${NC}"
            echo ""
            echo "👤 Dados do usuário:"
            echo "$ME_BODY" | jq '.' 2>/dev/null || echo "$ME_BODY"
        else
            echo -e "${RED}✗ Endpoint /auth/me falhou (HTTP $ME_HTTP_STATUS)${NC}"
            echo "$ME_BODY"
        fi
    fi
else
    echo -e "${RED}✗ Login falhou (HTTP $HTTP_STATUS)${NC}"
    echo ""
    echo "📦 Response Body:"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
fi

echo ""
echo "======================================"
echo "Teste concluído"
echo "======================================"
