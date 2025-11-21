# 🛠️ Configuração de Ambiente de Desenvolvimento — Barber Analytics Pro v2.0

**Sistema Operacional:** Pop!\_OS 22.04 LTS (Ubuntu/Debian)
**Última Atualização:** 21/11/2025
**Propósito:** Setup completo "Zero to Hero" após formatação

---

## 📋 Índice

1. [Stack Tecnológica](#1-stack-tecnológica)
2. [Dependências de Sistema](#2-dependências-de-sistema)
3. [Configuração do Projeto](#3-configuração-do-projeto)
4. [Extensões Obrigatórias (VS Code)](#4-extensões-obrigatórias-vs-code)
5. [Verificação do Ambiente](#5-verificação-do-ambiente)

---

## 1. Stack Tecnológica

| Tecnologia         | Versão                      | Função no Projeto                                     |
| ------------------ | --------------------------- | ----------------------------------------------------- |
| **Go**             | 1.24.0+ (toolchain 1.24.10) | Backend API (Clean Architecture)                      |
| **Node.js**        | 20.x LTS                    | Runtime para Next.js frontend                         |
| **pnpm**           | 9.x                         | Gerenciador de pacotes frontend (mais rápido que npm) |
| **PostgreSQL**     | 14+                         | Banco de dados principal (Neon recomendado)           |
| **Next.js**        | 16.0.3                      | Framework frontend (App Router)                       |
| **React**          | 19.0.0                      | Biblioteca UI                                         |
| **TypeScript**     | 5.x                         | Linguagem frontend                                    |
| **Echo**           | v4.13.4                     | Framework HTTP para Go                                |
| **SQLC**           | Latest                      | Type-safe SQL code generator                          |
| **golang-migrate** | Latest                      | Gerenciamento de migrations                           |
| **Docker**         | 24.x+                       | Containerização (opcional)                            |
| **Git**            | 2.x+                        | Controle de versão                                    |
| **Make**           | 4.x+                        | Automação de tarefas                                  |

---

## 2. Dependências de Sistema

### 2.1 Atualizar Sistema

```bash
sudo apt update && sudo apt upgrade -y
```

### 2.2 Instalar Dependências Essenciais

```bash
# Build essentials, Git, Make, curl, wget
sudo apt install -y build-essential git make curl wget \
  ca-certificates gnupg lsb-release software-properties-common \
  apt-transport-https
```

### 2.3 Instalar Go 1.24+

```bash
# Remover versões antigas (se houver)
sudo rm -rf /usr/local/go

# Baixar e instalar Go 1.24.0
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
rm go1.24.0.linux-amd64.tar.gz

# Adicionar ao PATH (adicione ao ~/.bashrc ou ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Verificar instalação
go version
```

### 2.4 Instalar Node.js 20 LTS (via NodeSource)

```bash
# Adicionar repositório NodeSource
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -

# Instalar Node.js
sudo apt install -y nodejs

# Verificar instalação
node -v
npm -v
```

### 2.5 Instalar pnpm

```bash
# Instalar pnpm globalmente
curl -fsSL https://get.pnpm.io/install.sh | sh -

# Recarregar shell
source ~/.bashrc

# Verificar instalação
pnpm -v
```

### 2.6 Instalar PostgreSQL 14+ (Cliente)

```bash
# Adicionar repositório oficial PostgreSQL
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update

# Instalar cliente PostgreSQL
sudo apt install -y postgresql-client-14

# Verificar instalação
psql --version
```

**⚠️ IMPORTANTE:** Para desenvolvimento, recomenda-se usar **Neon** (https://neon.tech) em vez de PostgreSQL local.

### 2.7 Instalar golang-migrate

```bash
# Via apt (pode estar desatualizado)
# sudo apt install -y golang-migrate

# OU via Go (recomendado - versão mais recente)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verificar instalação
migrate -version
```

### 2.8 Instalar SQLC

```bash
# Via Go
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verificar instalação
sqlc version
```

### 2.9 Instalar golangci-lint

```bash
# Via Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verificar instalação
golangci-lint version
```

### 2.10 Instalar goimports

```bash
go install golang.org/x/tools/cmd/goimports@latest

# Verificar instalação
goimports -h
```

### 2.11 Instalar Docker (Opcional)

```bash
# Adicionar repositório Docker
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update

# Instalar Docker Engine
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# IMPORTANTE: Fazer logout e login novamente para aplicar

# Verificar instalação
docker --version
docker compose version
```

### 2.12 Instalar VS Code

```bash
# Via snap (recomendado)
sudo snap install code --classic

# OU via apt
wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
rm -f packages.microsoft.gpg
sudo apt update
sudo apt install -y code
```

---

## 3. Configuração do Projeto

### 3.1 Clonar Repositório

```bash
cd ~/projetos
git clone https://github.com/andviana23/barber-analytics-proV2.git
cd barber-analytics-proV2
```

### 3.2 Configurar Backend (Go)

```bash
cd backend

# Baixar dependências Go
go mod download

# Instalar ferramentas de desenvolvimento (se ainda não instalou)
make install-tools

# Copiar arquivo de ambiente
cp .env.example .env

# Editar .env com suas credenciais
# Altere DATABASE_URL, JWT keys, ASAAS keys, etc.
nano .env
```

**⚠️ Gerar Chaves JWT:**

```bash
# Criar diretório de chaves
mkdir -p keys

# Gerar chave privada RSA
openssl genrsa -out keys/private.pem 2048

# Extrair chave pública
openssl rsa -in keys/private.pem -pubout -out keys/public.pem

# Verificar permissões
chmod 600 keys/private.pem
chmod 644 keys/public.pem
```

**🗄️ Configurar Banco de Dados:**

Opção A: **Neon (Recomendado)**

```bash
# 1. Criar conta em https://neon.tech
# 2. Criar projeto "barber-analytics-dev"
# 3. Copiar connection string
# 4. Colar no .env em DATABASE_URL
```

Opção B: **PostgreSQL Local**

```bash
# Instalar servidor PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Criar banco de dados
sudo -u postgres psql
CREATE DATABASE barber_analytics_dev;
CREATE USER barber_user WITH PASSWORD 'sua_senha_forte';
GRANT ALL PRIVILEGES ON DATABASE barber_analytics_dev TO barber_user;
\q

# Atualizar .env
DATABASE_URL=postgresql://barber_user:sua_senha_forte@localhost:5432/barber_analytics_dev?sslmode=disable
```

**🔄 Rodar Migrations:**

```bash
# Executar migrations
make migrate-up

# Verificar status
psql $DATABASE_URL -c "SELECT version, dirty FROM schema_migrations;"
```

**🌱 Seed de Dados Demo (Opcional):**

```bash
# Criar tenant demo + dados de exemplo
make seed-demo

# Verificar seed
make seed-verify
```

**▶️ Rodar Backend:**

```bash
# Desenvolvimento (hot reload com Air - se instalado)
make run

# OU rodar diretamente
go run cmd/api/main.go

# Testar
curl http://localhost:8080/health
```

### 3.3 Configurar Frontend (Next.js)

```bash
cd ../frontend

# Instalar dependências com pnpm
pnpm install

# Copiar arquivo de ambiente (se houver)
# cp .env.example .env.local

# Rodar em desenvolvimento
pnpm dev

# Abrir navegador em http://localhost:3000
```

**🧪 Rodar Testes:**

```bash
# Testes unitários
pnpm test:unit

# Testes E2E (Playwright)
pnpm test:e2e

# Testes com coverage
pnpm test:coverage
```

---

## 4. Extensões Obrigatórias (VS Code)

### 4.1 Backend (Go)

| Extensão            | ID                            | Descrição                                                    |
| ------------------- | ----------------------------- | ------------------------------------------------------------ |
| **Go**              | `golang.go`                   | Suporte completo para Go (IntelliSense, debug, test, format) |
| **Go Nightly**      | `golang.go-nightly`           | Versão nightly do Go extension (opcional)                    |
| **Error Lens**      | `usernamehw.errorlens`        | Mostra erros inline                                          |
| **Better Comments** | `aaron-bond.better-comments`  | Highlight de comentários                                     |
| **Docker**          | `ms-azuretools.vscode-docker` | Suporte Docker/Docker Compose                                |
| **PostgreSQL**      | `ckolkman.vscode-postgres`    | Cliente PostgreSQL no VS Code                                |
| **REST Client**     | `humao.rest-client`           | Testar APIs HTTP direto no VS Code                           |
| **GitLens**         | `eamodio.gitlens`             | Git supercharged                                             |

**Instalação via CLI:**

```bash
code --install-extension golang.go
code --install-extension usernamehw.errorlens
code --install-extension aaron-bond.better-comments
code --install-extension ms-azuretools.vscode-docker
code --install-extension ckolkman.vscode-postgres
code --install-extension humao.rest-client
code --install-extension eamodio.gitlens
```

### 4.2 Frontend (Next.js/React/TypeScript)

| Extensão                      | ID                                     | Descrição                          |
| ----------------------------- | -------------------------------------- | ---------------------------------- |
| **ESLint**                    | `dbaeumer.vscode-eslint`               | Linter JavaScript/TypeScript       |
| **Prettier**                  | `esbenp.prettier-vscode`               | Formatador de código               |
| **TypeScript Next.js**        | `willstakayama.vscode-nextjs-snippets` | Snippets Next.js                   |
| **Tailwind CSS IntelliSense** | `bradlc.vscode-tailwindcss`            | Autocomplete Tailwind              |
| **ES7+ React Snippets**       | `dsznajder.es7-react-js-snippets`      | Snippets React                     |
| **Auto Rename Tag**           | `formulahendry.auto-rename-tag`        | Renomear tags HTML automaticamente |
| **Import Cost**               | `wix.vscode-import-cost`               | Mostrar tamanho de imports         |
| **Path Intellisense**         | `christian-kohler.path-intellisense`   | Autocomplete de paths              |

**Instalação via CLI:**

```bash
code --install-extension dbaeumer.vscode-eslint
code --install-extension esbenp.prettier-vscode
code --install-extension willstakayama.vscode-nextjs-snippets
code --install-extension bradlc.vscode-tailwindcss
code --install-extension dsznajder.es7-react-js-snippets
code --install-extension formulahendry.auto-rename-tag
code --install-extension wix.vscode-import-cost
code --install-extension christian-kohler.path-intellisense
```

### 4.3 Geral (Produtividade)

| Extensão                | ID                             | Descrição              |
| ----------------------- | ------------------------------ | ---------------------- |
| **GitHub Copilot**      | `github.copilot`               | AI pair programming    |
| **GitHub Copilot Chat** | `github.copilot-chat`          | Chat com Copilot       |
| **Todo Tree**           | `gruntfuggly.todo-tree`        | Highlight de TODOs     |
| **Markdown All in One** | `yzhang.markdown-all-in-one`   | Suporte Markdown       |
| **Thunder Client**      | `rangav.vscode-thunder-client` | Cliente HTTP/REST      |
| **EditorConfig**        | `editorconfig.editorconfig`    | Configuração de editor |

**Instalação via CLI:**

```bash
code --install-extension github.copilot
code --install-extension github.copilot-chat
code --install-extension gruntfuggly.todo-tree
code --install-extension yzhang.markdown-all-in-one
code --install-extension rangav.vscode-thunder-client
code --install-extension editorconfig.editorconfig
```

### 4.4 Configurações Recomendadas (settings.json)

Adicione ao `.vscode/settings.json` do projeto:

```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true,
    "source.organizeImports": true
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.formatOnSave": true
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[javascript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[json]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "typescript.tsdk": "node_modules/typescript/lib",
  "typescript.enablePromptUseWorkspaceTsdk": true,
  "files.associations": {
    "*.sql": "sql"
  }
}
```

---

## 5. Verificação do Ambiente

### 5.1 Checklist Completo

Execute os comandos abaixo para verificar se tudo está instalado corretamente:

```bash
# Go
go version                    # Deve retornar 1.24.0+

# Node.js
node -v                       # Deve retornar v20.x.x

# pnpm
pnpm -v                       # Deve retornar 9.x.x

# PostgreSQL Client
psql --version               # Deve retornar 14.x+

# migrate
migrate -version             # Deve retornar versão

# sqlc
sqlc version                 # Deve retornar versão

# golangci-lint
golangci-lint version        # Deve retornar versão

# goimports
goimports -h                 # Deve mostrar help

# Docker (opcional)
docker --version             # Deve retornar versão
docker compose version       # Deve retornar versão

# Git
git --version                # Deve retornar versão

# Make
make --version               # Deve retornar versão

# VS Code
code --version               # Deve retornar versão
```

### 5.2 Testar Backend

```bash
cd ~/projetos/barber-analytics-proV2/backend

# Rodar testes
make test

# Rodar linter
make lint

# Build
make build

# Verificar binário
./bin/barber-analytics-backend --help
```

### 5.3 Testar Frontend

```bash
cd ~/projetos/barber-analytics-proV2/frontend

# Rodar testes unitários
pnpm test:unit

# Build de produção
pnpm build

# Verificar se build passou
```

### 5.4 Script de Verificação Automática

Crie um arquivo `scripts/check-env.sh`:

```bash
#!/bin/bash

echo "🔍 Verificando ambiente de desenvolvimento..."
echo ""

check_command() {
    if command -v $1 &> /dev/null; then
        echo "✅ $1: $(command -v $1)"
    else
        echo "❌ $1: NÃO INSTALADO"
    fi
}

check_command go
check_command node
check_command pnpm
check_command psql
check_command migrate
check_command sqlc
check_command golangci-lint
check_command goimports
check_command docker
check_command git
check_command make
check_command code

echo ""
echo "🎯 Versões:"
echo "Go: $(go version)"
echo "Node: $(node -v)"
echo "pnpm: $(pnpm -v)"

echo ""
echo "✅ Verificação concluída!"
```

Execute:

```bash
chmod +x scripts/check-env.sh
./scripts/check-env.sh
```

---

## 6. Troubleshooting

### 6.1 Erro: "go: command not found"

```bash
# Adicionar Go ao PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 6.2 Erro: "pnpm: command not found"

```bash
# Reinstalar pnpm
curl -fsSL https://get.pnpm.io/install.sh | sh -
source ~/.bashrc
```

### 6.3 Erro: "migrate: command not found"

```bash
# Instalar via Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verificar se $GOPATH/bin está no PATH
echo $PATH | grep $GOPATH/bin
```

### 6.4 Erro: "permission denied" ao rodar Docker

```bash
# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# Fazer logout e login novamente
```

### 6.5 Erro de conexão com PostgreSQL

```bash
# Verificar se PostgreSQL está rodando (se local)
sudo systemctl status postgresql

# Testar conexão
psql "postgresql://user:password@localhost:5432/barber_analytics_dev?sslmode=disable"

# Verificar variável de ambiente
echo $DATABASE_URL
```

---

## 7. Próximos Passos

Após configurar o ambiente:

1. ✅ Ler `docs/ARQUITETURA.md` para entender a estrutura
2. ✅ Ler `docs/GUIA_DEV_BACKEND.md` para convenções de código Go
3. ✅ Ler `docs/GUIA_FRONTEND.md` para padrões do frontend
4. ✅ Executar `make seed-demo` para ter dados de teste
5. ✅ Explorar endpoints em `docs/API_REFERENCE.md`
6. ✅ Começar a codar! 🚀

---

**Última Atualização:** 21/11/2025
**Mantenedor:** Equipe Barber Analytics Pro
**Suporte:** Abra uma issue no GitHub
