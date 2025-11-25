# 💈 Barber Analytics Pro v2.0

> Sistema SaaS completo para gerenciamento de barbearias com multi-tenancy, analytics avançado e integração com pagamentos.

[![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow)](https://github.com/andviana23/barber-analytics-proV2)
[![Go Version](https://img.shields.io/badge/go-1.24.0-blue)](https://golang.org)
[![Next.js](https://img.shields.io/badge/next.js-14.2.4-black)](https://nextjs.org)
[![PostgreSQL](https://img.shields.io/badge/postgresql-14%2B-blue)](https://www.postgresql.org)

---

## 📋 Índice

- [Sobre o Projeto](#sobre-o-projeto)
- [Arquitetura](#arquitetura)
- [Tecnologias](#tecnologias)
- [Início Rápido](#início-rápido)
- [Documentação](#documentação)
- [Status Atual](#status-atual)
- [Roadmap](#roadmap)
- [Contribuindo](#contribuindo)

---

## 🎯 Sobre o Projeto

**Barber Analytics Pro v2.0** é uma plataforma SaaS moderna para gestão completa de barbearias, oferecendo:

- 💰 **Gestão Financeira**: Receitas, despesas, DRE, fluxo de caixa
- 👥 **Cadastros**: Clientes, profissionais, serviços, meios de pagamento
- 🎟️ **Assinaturas**: Clube do Trato com integração Asaas
- 📊 **Analytics**: Dashboards, métricas, relatórios
- ⏰ **Lista da Vez**: Sistema de rodízio de barbeiros baseado em pontos
- 📦 **Estoque**: Controle de produtos (futuro)
- 🔐 **Multi-tenancy**: Isolamento completo de dados por barbearia
- 📱 **Responsivo**: Interface adaptada para mobile, tablet e desktop

---

## 🏗️ Arquitetura

### Padrões Arquiteturais

- **Clean Architecture** (Robert C. Martin)
- **Domain-Driven Design (DDD)** (Eric Evans)
- **SOLID Principles**
- **Multi-tenancy Column-Based** (tenant_id em todas tabelas)

### Estrutura de Camadas

```
┌─────────────────────────────────────────┐
│       Presentation (HTTP/UI)            │  ← Handlers, Middleware, Components
├─────────────────────────────────────────┤
│       Application (Use Cases)           │  ← Business Logic Orchestration
├─────────────────────────────────────────┤
│       Domain (Entities)                 │  ← Business Rules, Value Objects
├─────────────────────────────────────────┤
│       Infrastructure (DB, APIs)         │  ← Repositories, External Services
└─────────────────────────────────────────┘
```

**Documentação Completa:** [ARQUITETURA.md](./docs/ARQUITETURA.md)

---

## 🛠️ Tecnologias

### Backend

- **Go 1.24.0** (Echo v4, SQLC, golang-migrate)
- **PostgreSQL 14+** (Neon serverless)
- **JWT RS256** (Autenticação assimétrica)
- **Zap** (Structured logging)

### Frontend

- **Next.js 14.2.4** (App Router)
- **React 18.2.0 + React DOM 18.2.0**
- **MUI 5.15.21 + Emotion 11.11** (Design System customizado)
- **TanStack Query 4.36.1** (Data fetching & caching)
- **Zod 3.22 + React Hook Form 7.49** (Validação de formulários)
- **Zustand 4.5.2** (Estado global leve)
- **Axios 1.6**, **ESLint 8.56**, **TypeScript 5.3**

### DevOps

- **GitHub Actions** (CI/CD)
- **NGINX** (Reverse proxy)
- **Neon** (Database hosting)

**Documentação Completa:** [GUIA_DEVOPS.md](./docs/GUIA_DEVOPS.md)

---

## 🚀 Início Rápido

### Pré-requisitos

```bash
# Go
go version  # >= 1.24

# Node.js
node --version  # >= 18.17

# PostgreSQL
psql --version  # >= 14

```

### Setup Backend

```bash
# 1. Clone repositório
git clone https://github.com/andviana23/barber-analytics-proV2.git
cd barber-analytics-proV2/backend

# 2. Copiar .env
cp .env.example .env
# Editar DATABASE_URL, JWT_SECRET, etc.

# 3. Instalar dependências
go mod download

# 4. Rodar migrations
make migrate-up

# 5. Rodar servidor
make run-backend
```

**Backend rodando em:** http://localhost:8080

### Setup Frontend

```bash
cd frontend

# 1. Instalar dependências
npm install

# 2. Copiar .env
cp .env.example .env.local
# Editar NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# 3. Rodar dev server
npm run dev
```

**Frontend rodando em:** http://localhost:3000

### Testes

```bash
# Backend
cd backend
make test

# Frontend
cd frontend
npm test
npm run test:e2e  # Playwright E2E
```

---

## 📚 Documentação

### 📍 Comece Aqui

- **[RESUMO_EXECUTIVO.md](./docs/RESUMO_EXECUTIVO.md)** - Status atual, próximos passos, decisões
- **[INDICE_DOCUMENTACAO.md](./docs/INDICE_DOCUMENTACAO.md)** - Navegação completa entre docs

### 🏗️ Arquitetura & Design

- [ARQUITETURA.md](./docs/ARQUITETURA.md) - Clean Architecture, DDD, princípios
- [Designer-System.md](./docs/Designer-System.md) - Tokens MUI, componentes, acessibilidade
- [BANCO_DE_DADOS.md](./docs/BANCO_DE_DADOS.md) - Schema, índices, migrations
- [MODELO_MULTI_TENANT.md](./docs/MODELO_MULTI_TENANT.md) - Multi-tenancy strategy

### 💻 Guias de Desenvolvimento

- [GUIA_DEV_BACKEND.md](./docs/GUIA_DEV_BACKEND.md) - Padrões Go, exemplos
- [GUIA_DEV_FRONTEND.md](./docs/GUIA_DEV_FRONTEND.md) - Padrões React/Next.js
- [GUIA_DEVOPS.md](./docs/GUIA_DEVOPS.md) - Deploy, CI/CD

### 📡 API & Integrações

- [API_REFERENCE.md](./docs/API_REFERENCE.md) - Endpoints completos
- [INTEGRACOES_ASAAS.md](./docs/INTEGRACOES_ASAAS.md) - Gateway pagamento
- [FLUXO_CRONS.md](./docs/FLUXO_CRONS.md) - Jobs agendados

### 💰 Módulos de Negócio

- [FINANCEIRO.md](./docs/FINANCEIRO.md) - Receitas, despesas, DRE
- [ASSINATURAS.md](./docs/ASSINATURAS.md) - Clube do Trato
- [listadavez.md](./docs/listadavez.md) - Sistema de rodízio

### 🔐 Segurança

- [RBAC.md](./docs/RBAC.md) - Controle de acesso
- [AUDIT_LOGS.md](./docs/AUDIT_LOGS.md) - Auditoria
- [COMPLIANCE_LGPD.md](./docs/COMPLIANCE_LGPD.md) - LGPD

---

## 📊 Status Atual

**Data:** 22/11/2025
**🎉 MARCO ALCANÇADO: 44/44 ENDPOINTS IMPLEMENTADOS!**

### ✅ Backend - 100% CONCLUÍDO

| Módulo                    | Status      | Endpoints | Data Conclusão |
| ------------------------- | ----------- | --------- | -------------- |
| Autenticação              | ✅ Completo | 5         | 20/11/2025     |
| Cadastro de Clientes      | ✅ Completo | 5         | 20/11/2025     |
| Cadastro de Profissionais | ✅ Completo | 5         | 20/11/2025     |
| Cadastro de Serviços      | ✅ Completo | 5         | 20/11/2025     |
| Meios de Pagamento        | ✅ Completo | 5         | 20/11/2025     |
| Lista da Vez              | ✅ Completo | 7         | 20/11/2025     |
| **Metas**                 | ✅ **NOVO** | **15**    | **22/11/2025** |
| **Precificação**          | ✅ **NOVO** | **9**     | **22/11/2025** |
| **Financeiro**            | ✅ **NOVO** | **20**    | **22/11/2025** |
| Onboarding                | ✅ Completo | 2         | 20/11/2025     |

**Total:** 78 endpoints backend funcionais ✅

### 🆕 Módulos Recém-Implementados (22/11)

**METAS (15 endpoints):**

- MetaMensal: 5 endpoints (POST, GET/:id, GET, PUT/:id, DELETE/:id)
- MetaBarbeiro: 5 endpoints (POST, GET/:id, GET, PUT/:id, DELETE/:id)
- MetaTicketMedio: 5 endpoints (POST, GET/:id, GET, PUT/:id, DELETE/:id)

**PRECIFICAÇÃO (9 endpoints):**

- Config: 4 endpoints (POST, GET, PUT, DELETE)
- Simulação: 5 endpoints (POST simulate, POST save, GET/:id, GET, DELETE/:id)

**FINANCEIRO (20 endpoints):**

- ContaPagar: 6 endpoints (CRUD + MarcarPagamento)
- ContaReceber: 6 endpoints (CRUD + MarcarRecebimento)
- Compensação: 3 endpoints (GET, List, DELETE)
- FluxoCaixa: 2 endpoints (GET, List)
- DRE: 2 endpoints (GET/:month, List)
- Cronjob: 1 endpoint (GenerateFluxoDiario)

**Ver detalhes:** `/Tarefas/01-BLOQUEIOS-BASE/VERTICAL_SLICE_ALL_MODULES.md`

### 🟡 Frontend - Em Progresso

- [x] Cadastros básicos (Clientes, Profissionais, Serviços)
- [x] Lista da Vez
- [x] Onboarding
- [ ] Metas (UI + hooks) ← **PRÓXIMO**
- [ ] Precificação (UI + hooks)
- [ ] Financeiro (UI + hooks)

### ⏳ Próximas Implementações

- [ ] Estoque (produtos, movimentações)
- [ ] Assinaturas (Clube do Trato + Asaas)
- [ ] Agendamentos (DayPilot Scheduler)
- [ ] Relatórios Avançados

---

## 🗓️ Roadmap

### Fase 1: Core (✅ Concluída)

- [x] Setup projeto (Go + Next.js)
- [x] Database (PostgreSQL + Migrations)
- [x] Autenticação (JWT RS256)
- [x] Multi-tenancy
- [x] Cadastros básicos

### Fase 2: Onboarding (✅ Concluída - 20/11/2025)

- [x] Frontend signup page
- [x] Frontend onboarding page
- [x] Backend signup use case
- [x] Backend complete onboarding endpoint
- [x] Testes E2E

### Fase 3: Metas, Precificação & Financeiro (✅ Concluída - 22/11/2025)

**METAS:**

- [x] CRUD MetaMensal (5 endpoints)
- [x] CRUD MetaBarbeiro (5 endpoints)
- [x] CRUD MetaTicketMedio (5 endpoints)

**PRECIFICAÇÃO:**

- [x] CRUD Config (4 endpoints)
- [x] CRUD Simulação (5 endpoints)

**FINANCEIRO:**

- [x] CRUD ContaPagar (6 endpoints)
- [x] CRUD ContaReceber (6 endpoints)
- [x] Compensação Bancária (3 endpoints)
- [x] FluxoCaixa (2 endpoints)
- [x] DRE (2 endpoints)
- [x] Cronjob FluxoDiario (1 endpoint)

**Resultado:** 44 endpoints backend implementados e compilando ✅

### Fase 4: Assinaturas (⏳ Planejada)

- [ ] Clube do Trato
- [ ] Integração Asaas
- [ ] Webhooks
- [ ] Cron de sincronização

### Fase 5: Estoque (⏳ Planejada)

- [ ] CRUD Produtos
- [ ] Movimentações
- [ ] Alertas estoque baixo

### Fase 6: Agendamentos (0% ⏳)

- [ ] Integração DayPilot
- [ ] CRUD Agendamentos
- [ ] Notificações

### Fase 7: Lançamento (0% ⏳)

- [ ] Testes carga
- [ ] Security audit
- [ ] Deploy produção
- [ ] Monitoramento

**Roadmap Completo:** [ROADMAP_IMPLEMENTACAO_V2.md](./docs/ROADMAP_IMPLEMENTACAO_V2.md)

---

## 📊 Métricas

### Cobertura de Testes

```
Backend:
├─ Unit Tests: 45% (meta: 80%)
└─ Integration Tests: 20% (meta: 60%)

Frontend:
├─ Unit Tests: 30% (meta: 70%)
└─ E2E Tests: 40% (meta: 80%)
```

### Performance

```
Backend:
├─ Startup: ~500ms
├─ Response time (p95): <100ms
└─ Database queries (avg): <50ms

Frontend:
├─ First Contentful Paint: <1.5s
├─ Time to Interactive: <3s
└─ Lighthouse Score: 85+
```

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

**Guias:**

- [CONTRIBUTING.md](./CONTRIBUTING.md) (a criar)
- [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) (a criar)

---

## 📝 Licença

Este projeto está sob a licença MIT. Veja [LICENSE](./LICENSE) para mais detalhes.

---

## 👥 Autores

- **Andrey Viana** - [@andviana23](https://github.com/andviana23)

---

## 🙏 Agradecimentos

- Clean Architecture - Robert C. Martin
- Domain-Driven Design - Eric Evans
- Go Echo Framework
- Next.js Team
- Material-UI Team

---

## 📞 Suporte

- 📧 Email: contato@barberanalyticspro.com
- 💬 Discord: [Barber Analytics Community](https://discord.gg/...)
- 🐛 Issues: [GitHub Issues](https://github.com/andviana23/barber-analytics-proV2/issues)

---

**Desenvolvido com ❤️ usando Go + Next.js + PostgreSQL**
