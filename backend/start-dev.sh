#!/bin/bash
# Script para iniciar o backend com variáveis de ambiente

# Carregar .env se existir
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

# Executar air
exec air
