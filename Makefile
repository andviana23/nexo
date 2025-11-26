# ============================================================================
# Barber Analytics Pro - Makefile
# ============================================================================
# Descrição: Automação para desenvolvimento local (Backend Go + Frontend Next.js)
# Backend: Air (hot-reload Go)
# Frontend: Next.js 14.2.4 (compatível com React 18.2.0 + MUI 5.15.21/Emotion 11.11)
# Database: Neon PostgreSQL (remoto)
# ============================================================================

# ============================================================================
# Variáveis
# ============================================================================
PROJECT_ROOT := $(shell pwd)
BACKEND_DIR := $(PROJECT_ROOT)/backend
FRONTEND_DIR := $(PROJECT_ROOT)/frontend

BACKEND_PID := $(BACKEND_DIR)/.backend.pid
FRONTEND_PID := $(FRONTEND_DIR)/.frontend.pid

BACKEND_LOG := $(BACKEND_DIR)/tmp/backend.log
FRONTEND_LOG := $(FRONTEND_DIR)/tmp/frontend.log

# DATABASE_URL deve ser definida via variável de ambiente (.env)
# Exemplo: export DATABASE_URL="postgresql://user:pass@host/db?sslmode=require"
DATABASE_URL ?= $(shell echo $$DATABASE_URL)
API_URL ?= http://localhost:8080

GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
RED := \033[0;31m
NC := \033[0m

# ============================================================================
# Targets Principais
# ============================================================================

.PHONY: help
help: ## Exibir esta mensagem de ajuda
	@echo ""
	@echo "$(BLUE)╔════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(BLUE)║$(NC)  🚀 Barber Analytics Pro - Makefile Commands"
	@echo "$(BLUE)╚════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

.PHONY: dev
dev: ## Iniciar backend + frontend em paralelo
	@echo "$(BLUE)╔════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(BLUE)║$(NC)  🚀 Iniciando Backend + Frontend..."
	@echo "$(BLUE)╚════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@$(MAKE) -j2 backend frontend
	@echo ""
	@echo "$(GREEN)✅ Sistema iniciado!$(NC)"
	@echo ""
	@echo "$(YELLOW)📡 Backend:$(NC)  http://localhost:8080"
	@echo "$(YELLOW)🌐 Frontend:$(NC) http://localhost:3000"
	@echo ""
	@echo "$(BLUE)Para parar:$(NC) make stop"
	@echo ""

.PHONY: backend
backend: ## Iniciar apenas o backend (Air + Go)
	@echo "$(YELLOW)🔧 Backend (Go + Air)...$(NC)"
	@mkdir -p "$(BACKEND_DIR)/tmp"
	@if [ -f "$(BACKEND_PID)" ]; then \
		echo "$(RED)❌ Backend já está rodando (PID: $$(cat "$(BACKEND_PID)"))$(NC)"; \
		exit 1; \
	fi
	@cd "$(BACKEND_DIR)" && \
		nohup ./start-dev.sh > "$(BACKEND_LOG)" 2>&1 & echo $$! > "$(BACKEND_PID)"
	@sleep 2
	@if ps -p $$(cat "$(BACKEND_PID)") > /dev/null 2>&1; then \
		echo "$(GREEN)✅ Backend rodando (PID: $$(cat "$(BACKEND_PID)"))$(NC)"; \
		echo "$(BLUE)   Logs:$(NC) tail -f $(BACKEND_LOG)"; \
		echo "$(BLUE)   URL:$(NC)  http://localhost:8080"; \
	else \
		echo "$(RED)❌ Falha ao iniciar backend$(NC)"; \
		cat "$(BACKEND_LOG)"; \
		rm -f "$(BACKEND_PID)"; \
		exit 1; \
	fi

.PHONY: frontend
frontend: ## Iniciar apenas o frontend (Next.js)
	@echo "$(YELLOW)⚛️  Frontend (Next.js)...$(NC)"
	@mkdir -p "$(FRONTEND_DIR)/tmp"
	@if [ -f "$(FRONTEND_PID)" ]; then \
		echo "$(RED)❌ Frontend já está rodando (PID: $$(cat "$(FRONTEND_PID)"))$(NC)"; \
		exit 1; \
	fi
	@echo "   Verificando porta 3000..."
	@if lsof -ti:3000 >/dev/null 2>&1; then \
		echo "$(RED)   ❌ Porta 3000 em uso. Finalizando processo...$(NC)"; \
		lsof -ti:3000 | xargs kill -9 2>/dev/null || true; \
		sleep 1; \
	fi
	@echo "   Removendo locks anteriores..."
	@rm -rf "$(FRONTEND_DIR)/.next/dev/lock" 2>/dev/null || true
	@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
		echo "$(YELLOW)📦 Instalando dependências...$(NC)"; \
		cd "$(FRONTEND_DIR)" && pnpm install; \
	fi
	@cd "$(FRONTEND_DIR)" && \
		nohup pnpm dev > "$(FRONTEND_LOG)" 2>&1 & echo $$! > "$(FRONTEND_PID)"
	@sleep 3
	@if ps -p $$(cat "$(FRONTEND_PID)") > /dev/null 2>&1; then \
		echo "$(GREEN)✅ Frontend rodando (PID: $$(cat "$(FRONTEND_PID)"))$(NC)"; \
		echo "$(BLUE)   Logs:$(NC) tail -f $(FRONTEND_LOG)"; \
		echo "$(BLUE)   URL:$(NC)  http://localhost:3000"; \
	else \
		echo "$(RED)❌ Falha ao iniciar frontend$(NC)"; \
		cat "$(FRONTEND_LOG)"; \
		rm -f "$(FRONTEND_PID)"; \
		exit 1; \
	fi

.PHONY: stop
stop: ## Parar backend + frontend
	@echo "$(RED)🛑 Parando serviços...$(NC)"
	@if [ -f "$(BACKEND_PID)" ]; then \
		echo "   Parando backend (PID: $$(cat "$(BACKEND_PID)"))..."; \
		kill $$(cat "$(BACKEND_PID)") 2>/dev/null || true; \
		pkill -P $$(cat "$(BACKEND_PID)") 2>/dev/null || true; \
		rm -f "$(BACKEND_PID)"; \
		echo "   $(GREEN)✅ Backend parado$(NC)"; \
	else \
		echo "   $(YELLOW)⚠️  Backend não estava rodando$(NC)"; \
	fi
	@if [ -f "$(FRONTEND_PID)" ]; then \
		echo "   Parando frontend (PID: $$(cat "$(FRONTEND_PID)"))..."; \
		kill -TERM $$(cat "$(FRONTEND_PID)") 2>/dev/null || true; \
		sleep 1; \
		pkill -P $$(cat "$(FRONTEND_PID)") 2>/dev/null || true; \
		rm -f "$(FRONTEND_PID)"; \
		echo "   $(GREEN)✅ Frontend parado$(NC)"; \
	else \
		echo "   $(YELLOW)⚠️  Frontend não estava rodando$(NC)"; \
	fi
	@echo "   Finalizando processos remanescentes..."
	@pkill -9 -f "air" 2>/dev/null || true
	@pkill -9 -f "next dev" 2>/dev/null || true
	@pkill -9 -f "next-server" 2>/dev/null || true
	@pkill -9 -f "node.*next" 2>/dev/null || true
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "   Removendo arquivos de lock..."
	@rm -rf "$(FRONTEND_DIR)/.next/dev/lock" 2>/dev/null || true
	@rm -rf "$(FRONTEND_DIR)/.next/cache/webpack" 2>/dev/null || true
	@sleep 1
	@echo ""
	@echo "$(GREEN)✅ Todos os serviços foram parados$(NC)"

.PHONY: restart
restart: ## Reiniciar backend + frontend
	@echo "$(YELLOW)🔄 Reiniciando sistema...$(NC)"
	@echo ""
	@$(MAKE) stop
	@echo ""
	@sleep 2
	@$(MAKE) dev

.PHONY: force-stop
force-stop: ## Parar TODOS os processos (emergência - mata tudo brutalmente)
	@echo "$(RED)⚠️  FORÇA BRUTA: Matando TODOS os processos...$(NC)"
	@pkill -9 -f "air" 2>/dev/null || true
	@pkill -9 -f "next" 2>/dev/null || true
	@pkill -9 -f "node" 2>/dev/null || true
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@rm -rf $(BACKEND_DIR)/.backend.pid 2>/dev/null || true
	@rm -rf $(FRONTEND_DIR)/.frontend.pid 2>/dev/null || true
	@rm -rf $(FRONTEND_DIR)/.next/dev 2>/dev/null || true
	@echo "$(GREEN)✅ Limpeza forçada concluída$(NC)"
	@echo "$(YELLOW)⚠️  AVISO: Este comando matou TODOS os processos Node.js e Go$(NC)"

.PHONY: status
status: ## Verificar status dos serviços
	@echo "$(BLUE)╔════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(BLUE)║$(NC)  📊 Status dos Serviços"
	@echo "$(BLUE)╚════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@if [ -f "$(BACKEND_PID)" ] && ps -p $$(cat "$(BACKEND_PID)") > /dev/null 2>&1; then \
		echo "$(GREEN)✅ Backend:$(NC)  Rodando (PID: $$(cat "$(BACKEND_PID)"))"; \
		echo "   $(BLUE)URL:$(NC) http://localhost:8080"; \
	else \
		echo "$(RED)❌ Backend:$(NC)  Parado"; \
		rm -f "$(BACKEND_PID)" 2>/dev/null || true; \
	fi
	@if [ -f "$(FRONTEND_PID)" ] && ps -p $$(cat "$(FRONTEND_PID)") > /dev/null 2>&1; then \
		echo "$(GREEN)✅ Frontend:$(NC) Rodando (PID: $$(cat "$(FRONTEND_PID)"))"; \
		echo "   $(BLUE)URL:$(NC) http://localhost:3000"; \
	else \
		echo "$(RED)❌ Frontend:$(NC) Parado"; \
		rm -f "$(FRONTEND_PID)" 2>/dev/null || true; \
	fi
	@echo ""

.PHONY: logs-backend
logs-backend: ## Ver logs do backend (tail -f)
	@if [ -f $(BACKEND_LOG) ]; then \
		echo "$(BLUE)📋 Logs do Backend (Ctrl+C para sair):$(NC)"; \
		tail -f $(BACKEND_LOG); \
	else \
		echo "$(RED)❌ Arquivo de log não encontrado$(NC)"; \
	fi

.PHONY: logs-frontend
logs-frontend: ## Ver logs do frontend (tail -f)
	@if [ -f $(FRONTEND_LOG) ]; then \
		echo "$(BLUE)📋 Logs do Frontend (Ctrl+C para sair):$(NC)"; \
		tail -f $(FRONTEND_LOG); \
	else \
		echo "$(RED)❌ Arquivo de log não encontrado$(NC)"; \
	fi

.PHONY: logs
logs: ## Ver logs de ambos em paralelo
	@echo "$(BLUE)📋 Logs do Sistema (Ctrl+C para sair):$(NC)"
	@echo ""
	@if [ -f $(BACKEND_LOG) ] && [ -f $(FRONTEND_LOG) ]; then \
		tail -f $(BACKEND_LOG) $(FRONTEND_LOG); \
	else \
		echo "$(RED)❌ Arquivos de log não encontrados$(NC)"; \
	fi

.PHONY: clean
clean: stop ## Limpar arquivos temporários e logs
	@echo "$(YELLOW)🧹 Limpando arquivos temporários...$(NC)"
	@rm -rf $(BACKEND_DIR)/tmp/*.log
	@rm -rf $(BACKEND_DIR)/.backend.pid
	@rm -rf $(FRONTEND_DIR)/tmp/*.log
	@rm -rf $(FRONTEND_DIR)/.frontend.pid
	@rm -rf $(BACKEND_DIR)/.air.toml.lock 2>/dev/null || true
	@rm -rf $(FRONTEND_DIR)/.next/dev 2>/dev/null || true
	@rm -rf $(FRONTEND_DIR)/.next/cache 2>/dev/null || true
	@echo "$(GREEN)✅ Limpeza concluída$(NC)"

.PHONY: test-backend
test-backend: ## Testar se backend está respondendo
	@echo "$(BLUE)🧪 Testando backend...$(NC)"
	@curl -s http://localhost:8080/api/v1/ping || echo "$(RED)❌ Backend não está respondendo$(NC)"

.PHONY: test-frontend
test-frontend: ## Testar se frontend está respondendo
	@echo "$(BLUE)🧪 Testando frontend...$(NC)"
	@curl -s http://localhost:3000 > /dev/null && echo "$(GREEN)✅ Frontend OK$(NC)" || echo "$(RED)❌ Frontend não está respondendo$(NC)"

.PHONY: install
install: ## Instalar dependências (backend + frontend)
	@echo "$(YELLOW)📦 Instalando dependências...$(NC)"
	@echo ""
	@echo "$(BLUE)🔧 Backend (Go modules)...$(NC)"
	@cd $(BACKEND_DIR) && go mod download
	@echo "$(GREEN)✅ Backend OK$(NC)"
	@echo ""
	@echo "$(BLUE)⚛️  Frontend (pnpm)...$(NC)"
	@cd $(FRONTEND_DIR) && pnpm install
	@echo "$(GREEN)✅ Frontend OK$(NC)"
	@echo ""
	@echo "$(GREEN)✅ Todas as dependências instaladas$(NC)"

.PHONY: build-backend
build-backend: ## Build do backend (produção)
	@echo "$(BLUE)🏗️  Building backend...$(NC)"
	@cd $(BACKEND_DIR) && go build -o bin/barber-api ./cmd/api
	@echo "$(GREEN)✅ Backend compilado: $(BACKEND_DIR)/bin/barber-api$(NC)"

.PHONY: build-frontend
build-frontend: ## Build do frontend (produção)
	@echo "$(BLUE)🏗️  Building frontend...$(NC)"
	@cd $(FRONTEND_DIR) && pnpm build
	@echo "$(GREEN)✅ Frontend compilado$(NC)"

.PHONY: build
build: build-backend build-frontend ## Build completo (backend + frontend)

.PHONY: validate-schema
validate-schema: ## Validar schema do banco com scripts/validate_schema.sh (usa DATABASE_URL)
	@echo "$(BLUE)🔍 Validando schema do banco...$(NC)"
	@./scripts/validate_schema.sh "$(DATABASE_URL)"

.PHONY: smoke-tests
smoke-tests: ## Executar smoke tests E2E contra API (ajuste API_URL se necessário)
	@echo "$(BLUE)🧪 Smoke tests na API: $(API_URL)$(NC)"
	@./scripts/smoke_tests.sh "$(API_URL)"

# ============================================================================
# Default Target
# ============================================================================
.DEFAULT_GOAL := help
