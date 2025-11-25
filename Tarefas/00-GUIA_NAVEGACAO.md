# 🧭 Guia de Navegação — Barber Analytics Pro v2.0

**Versão:** 2.0
**Última Atualização:** 21/11/2025
**Objetivo:** Mapa completo de todas as tarefas, organizado de forma sequencial e lógica

---

## 📊 Status Geral do Projeto

| Componente             | Status      | Progresso                                    |
| ---------------------- | ----------- | -------------------------------------------- |
| **Banco de Dados**     | ✅ Completo | 100% (42 tabelas, migrations 001-038)        |
| **Backend (Go)**       | 🟢 Avançado | ~90% (Domain, Repos, Use Cases, Handlers OK) |
| **Frontend (Next.js)** | 🟢 Avançado | ~85% (Services, Hooks, Components OK)        |
| **Hardening & OPS**    | ✅ Completo | 100% (LGPD, Backup/DR, Testes QA completos)  |
| **Bloqueios de Base**  | ✅ Completo | 100% - **CONCLUÍDO**                         |
| **Financeiro**         | ✅ Completo | 100% Backend, 87.5% Frontend                 |

---

## 🎯 Sequência de Execução (ORDEM OBRIGATÓRIA)

```
┌─────────────────────────────────────────────────────────────────┐
│  🚨 ATENÇÃO: Seguir esta ordem RIGOROSAMENTE!                  │
│  Não pule etapas ou execute fora de ordem                      │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────────┐
│   ETAPA 0 (PRÉ)     │  Status: ✅ CONCLUÍDO
├──────────────────────┤
│ • Banco de Dados     │  ✅ 100% (42 tabelas)
│ • Migrations         │  ✅ 001-038 aplicadas
│ • Infraestrutura     │  ✅ Neon PostgreSQL configurado
└──────────────────────┘
           ↓
┌──────────────────────┐
│   ETAPA 1 (BASE)    │  Status: ✅ CONCLUÍDO
├──────────────────────┤  Tempo: 23 dias (Completado em 22/11/2025)
│ 01-BLOQUEIOS-BASE    │  📂 Tarefas/01-BLOQUEIOS-BASE/
│                      │  📂 Tarefas/CONCLUIR/
│ Sub-etapas:          │
│ ├─ Domain (19 ent)   │  ✅ Concluído
│ ├─ Ports/Interfaces  │  ✅ Concluído
│ ├─ Repositories      │  ✅ Concluído
│ ├─ Use Cases         │  ✅ Concluído
│ ├─ HTTP Handlers     │  ✅ Concluído
│ ├─ Cron Jobs         │  ✅ Concluído (3/6 ativos)
│ ├─ Frontend Services │  ✅ Concluído (7 services)
│ └─ Frontend Hooks    │  ✅ Concluído (43 hooks)
└──────────────────────┘
           ↓
┌──────────────────────┐
│   ETAPA 2 (OPS)     │  Status: ✅ CONCLUÍDO
├──────────────────────┤  Tempo: 1 semana (Completado em 24/11/2025)
│ 02-HARDENING-OPS     │  📂 Tarefas/02-HARDENING-OPS/
│                      │
│ Tarefas:             │
│ ├─ T-HAR-001 LGPD    │  ✅ 4 endpoints + Privacy Page
│ ├─ T-HAR-002 Backup  │  ✅ Workflow + DR Runbook
│ ├─ T-HAR-003 Valid   │  ✅ Checklist Seg/Obs completo
│ └─ QA (8 testes)     │  ✅ 110+ casos de teste criados
└──────────────────────┘
           ↓
┌──────────────────────────────────────────────────────────────┐
│                   ETAPAS 3-6 (MÓDULOS)                        │
│               Status: ✅ BASE PRONTA - PODE EXECUTAR          │
│                                                               │
│  Podem ser executadas EM PARALELO (Etapas 1-2 concluídas)   │
└──────────────────────────────────────────────────────────────┘
           ↓
┌──────────────────────┐    ┌──────────────────────┐
│  ETAPA 3 (FIN)      │    │  ETAPA 4 (EST)      │
├──────────────────────┤    ├──────────────────────┤
│ 03-FINANCEIRO        │    │ 04-ESTOQUE          │
│ ✅ CONCLUÍDO 100%    │    │                      │
│                      │    │ 6 sub-módulos:       │
│ Status:              │    │ 1. Entrada           │
│ ✅ Backend: 21/21    │    │ 2. Saída             │
│ ✅ Repos: 5/5        │    │ 3. Consumo Auto      │
│ ✅ Use Cases: 24/24  │    │ 4. Inventário        │
│ ✅ Endpoints: 100%   │    │ 5. Estoque Mínimo    │
│ ✅ Hooks: 8/8        │    │ 6. Curva ABC         │
│ ✅ Dashboard: 100%   │    │                      │
│ ⏸️ Comissões (baixa) │    │                      │
└──────────────────────┘    └──────────────────────┘

┌──────────────────────┐    ┌──────────────────────┐
│  ETAPA 5 (METAS)    │    │  ETAPA 6 (PREC)     │
├──────────────────────┤    ├──────────────────────┤
│ 05-METAS             │    │ 06-PRECIFICACAO     │
│                      │    │                      │
│ 4 sub-módulos:       │    │ 1 módulo:            │
│ 1. Meta Mensal       │    │ 1. Simulador         │
│ 2. Meta Barbeiro     │    │                      │
│ 3. Meta Ticket       │    │                      │
│ 4. Metas Automáticas │    │                      │
└──────────────────────┘    └──────────────────────┘
           ↓
┌──────────────────────────────────────────────────────────────┐
│              ETAPAS 7-10 (LANÇAMENTO E EVOLUÇÃO)             │
│                  Executar SEQUENCIALMENTE                     │
└──────────────────────────────────────────────────────────────┘
           ↓
┌──────────────────────┐
│  ETAPA 7 (LAUNCH)   │  Status: ⏳ Após 3-6
├──────────────────────┤
│ 07-LANCAMENTO        │  Go-Live + Comunicação
└──────────────────────┘
           ↓
┌──────────────────────┐
│  ETAPA 8 (MON)      │  Status: ⏳ Após 7
├──────────────────────┤
│ 08-MONITORAMENTO     │  Suporte 24/7 + Hotfixes
└──────────────────────┘
           ↓
┌──────────────────────┐
│  ETAPA 9 (EVO)      │  Status: ⏳ Após 8
├──────────────────────┤
│ 09-EVOLUCAO          │  PMF + Crescimento
└──────────────────────┘
           ↓
┌──────────────────────┐
│  ETAPA 10 (AGE)     │  Status: ⏳ Após 9
├──────────────────────┤
│ 10-AGENDAMENTOS      │  DayPilot + Notificações
└──────────────────────┘
```

---

## 📂 Estrutura de Pastas (Organizada)

```
Tarefas/
│
├── 00-GUIA_NAVEGACAO.md          ← VOCÊ ESTÁ AQUI! 🎯
│
├── INDICE_TAREFAS.md              ← Índice oficial com diagrama Mermaid
├── DATABASE_MIGRATIONS_COMPLETED.md  ← Status do banco (FEITO ✅)
├── INTEGRACAO_ASAAS_PLANO.md      ← Referência técnica Asaas
│
├── CONCLUIR/                      ← 🔴 BLOQUEADOR - Ler PRIMEIRO!
│   ├── README.md                  ← Resumo de bloqueios
│   ├── 00-ANALISE_SISTEMA_ATUAL.md   ✅ Análise completa
│   ├── 01-backend-domain-entities.md ❌ 3-4 dias
│   ├── 02-backend-repository-interfaces.md ❌ 2 dias
│   └── 03-08-resumo-tarefas-restantes.md ❌ 17 dias
│
├── 01-BLOQUEIOS-BASE/             ← 🔴 EXECUTAR PRIMEIRO! (Sprint 11-12)
│   ├── README.md                  ← Overview do bloqueador
│   ├── 01-contexto.md             ← Estado atual e lacunas
│   ├── 02-backlog.md              ← Tarefas técnicas detalhadas
│   ├── 03-sprint-plan.md          ← Ordem de execução
│   ├── 04-checklist-dev.md        ← Critérios de pronto (Dev)
│   ├── 05-checklist-qa.md         ← Critérios de qualidade (QA)
│   └── FASE_5_MIGRACAO.md         ← Contexto legado (migrado para cá)
│
├── 02-HARDENING-OPS/              ← 🟡 Após 01-BLOQUEIOS
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   └── FASE_6_HARDENING.md
│
├── 03-FINANCEIRO/                 ← 🟢 Paralelo com 04-06 (após 01)
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   ├── 01-dashboard-financeiro.md    ← Módulo 6
│   ├── 02-dre.md                     ← Módulo 5 (parte 1)
│   ├── 03-contas-a-pagar.md          ← Módulo 1
│   ├── 04-contas-a-receber.md        ← Módulo 2
│   ├── 05-comissoes-automaticas.md   ← Módulo 4
│   ├── 06-dre-completo.md            ← Módulo 5 (parte 2)
│   └── 07-fluxo-caixa-compensado.md  ← Módulo 3
│
├── 04-ESTOQUE/                    ← 🟢 Paralelo com 03,05,06 (após 01)
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   ├── 01-entrada.md                 ← Módulo 1
│   ├── 02-saida.md                   ← Módulo 2
│   ├── 03-consumo-automatico.md      ← Módulo 3
│   ├── 04-inventario.md              ← Módulo 4
│   ├── 05-curva-abc.md               ← Módulo 6
│   └── 06-estoque-minimo.md          ← Módulo 5
│
├── 05-METAS/                      ← 🟢 Paralelo com 03,04,06 (após 01)
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   ├── 01-meta-geral-mes.md          ← Módulo 1
│   ├── 02-meta-por-barbeiro.md       ← Módulo 2
│   ├── 03-meta-ticket-medio.md       ← Módulo 3
│   └── 04-metas-automaticas.md       ← Módulo 4
│
├── 06-PRECIFICACAO/               ← 🟢 Paralelo com 03-05 (após 01)
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   └── 01-precificacao-simulador.md  ← Módulo único
│
├── 07-LANCAMENTO/                 ← 🔵 Após 02-06 concluídos
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   └── FASE_7_LANCAMENTO.md
│
├── 08-MONITORAMENTO/              ← 🔵 Após 07 (Go-Live)
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   └── FASE_8_MONITORING.md
│
├── 09-EVOLUCAO/                   ← 🔵 Após 08 (Ciclos de 2 semanas)
│   ├── README.md
│   ├── 01-contexto.md
│   ├── 02-backlog.md
│   ├── 03-sprint-plan.md
│   ├── 04-checklist-dev.md
│   ├── 05-checklist-qa.md
│   └── FASE_9_EVOLUCAO.md
│
└── 10-AGENDAMENTOS/               ← 🔵 Após 09 (Feature complexa)
    ├── README.md
    ├── 01-contexto.md
    ├── 02-backlog.md
    ├── 03-sprint-plan.md
    ├── 04-checklist-dev.md
    ├── 05-checklist-qa.md
    └── FASE_10_AGENDAMENTOS.md
```

---

## 🚀 Por Onde Começar?

### Se você é NOVO no projeto:

1. ✅ **Leia ESTE arquivo** (`00-GUIA_NAVEGACAO.md`)
2. ✅ **Leia** `INDICE_TAREFAS.md` (visão geral + diagrama)
3. ✅ **Leia** `DATABASE_MIGRATIONS_COMPLETED.md` (banco já está pronto!)
4. 🔴 **Leia** `CONCLUIR/README.md` (entenda os bloqueadores)
5. 🔴 **Leia** `CONCLUIR/00-ANALISE_SISTEMA_ATUAL.md` (estado atual)
6. 🔴 **Execute** as tarefas de `CONCLUIR/` na ordem:
   - `01-backend-domain-entities.md`
   - `02-backend-repository-interfaces.md`
   - `03-08-resumo-tarefas-restantes.md`
7. 🔴 **Execute** `01-BLOQUEIOS-BASE/` (seguir sprint-plan)
8. ✅ **Após concluir 01**, execute módulos 02-10 na ordem

### Se você já conhece o projeto:

1. ✅ **Verifique** o status no `INDICE_TAREFAS.md`
2. 🔴 **Se 01-BLOQUEIOS-BASE não foi concluído**: PARE e execute primeiro
3. ✅ **Se 01 foi concluído**: Escolha próxima etapa no diagrama
4. ✅ **Sempre consulte** `02-backlog.md` de cada pasta antes de começar

---

## 📋 Padrão de Arquivos (Todas as Pastas)

Cada pasta `XX-NOME/` segue esta estrutura:

| Arquivo                   | Descrição                                     | Quando Ler           |
| ------------------------- | --------------------------------------------- | -------------------- |
| `README.md`               | Overview da etapa, objetivos, status          | Sempre PRIMEIRO      |
| `01-contexto.md`          | Estado atual, lacunas, pré-requisitos         | Antes de planejar    |
| `02-backlog.md`           | Lista detalhada de tarefas técnicas           | Antes de executar    |
| `03-sprint-plan.md`       | Ordem de execução, dependências               | Ao iniciar sprint    |
| `04-checklist-dev.md`     | Critérios de "pronto" (desenvolvedor)         | Durante dev          |
| `05-checklist-qa.md`      | Critérios de qualidade (QA/testes)            | Antes de deploy      |
| `0X-modulo-especifico.md` | Detalhes de sub-módulos (se aplicável)        | Conforme necessidade |
| `FASE_X_NOME.md`          | Documento legado de planejamento (referência) | Consulta opcional    |

---

## 🎯 Estimativas de Tempo

| Etapa                 | Tempo Estimado | Pode Paralelizar?              |
| --------------------- | -------------- | ------------------------------ |
| **01-BLOQUEIOS-BASE** | 23 dias        | ❌ NÃO (é pré-requisito)       |
| **02-HARDENING-OPS**  | 5-7 dias       | ❌ NÃO (após 01)               |
| **03-FINANCEIRO**     | 10-12 dias     | ✅ SIM (com 04-06, após 01)    |
| **04-ESTOQUE**        | 8-10 dias      | ✅ SIM (com 03,05,06, após 01) |
| **05-METAS**          | 6-8 dias       | ✅ SIM (com 03,04,06, após 01) |
| **06-PRECIFICACAO**   | 4-5 dias       | ✅ SIM (com 03-05, após 01)    |
| **07-LANCAMENTO**     | 3-5 dias       | ❌ NÃO (após 02-06)            |
| **08-MONITORAMENTO**  | Contínuo       | ❌ NÃO (após 07)               |
| **09-EVOLUCAO**       | Ciclos 2 sem   | ❌ NÃO (após 08)               |
| **10-AGENDAMENTOS**   | 10-12 dias     | ❌ NÃO (após 09)               |

**Total (sem paralelização):** ~90 dias úteis
**Total (com paralelização de 03-06):** ~60 dias úteis

---

## 🔗 Referências Importantes

### Documentação Técnica:

- `docs/02-arquitetura/ARQUITETURA.md` - Arquitetura Clean + DDD
- `docs/02-arquitetura/MODELO_DE_DADOS.md` - Schema do banco completo
- `docs/02-arquitetura/FLUXOS_CRITICOS_SISTEMA.md` - Fluxos principais
- `docs/04-backend/GUIA_DEV_BACKEND.md` - Padrões Go
- `docs/03-frontend/GUIA_FRONTEND.md` - Padrões Next.js + React
- `docs/03-frontend/DESIGN_SYSTEM.md` - Componentes UI
- `.github/copilot-instructions.md` - Regras para IA

### Diagramas:

- `docs/DIAGRAMA_DEPENDENCIAS_COMPLETO.md` - Visualização completa
- `docs/diagrams/` - Diagramas MermaidJS

### Produto:

- `docs/07-produto-e-funcionalidades/CATALOGO_FUNCIONALIDADES.md`
- `docs/07-produto-e-funcionalidades/ROADMAP_PRODUTO.md`

---

## ⚠️ Regras CRÍTICAS

### ❌ NÃO FAÇA:

1. ❌ Pular a etapa `01-BLOQUEIOS-BASE`
2. ❌ Executar módulos 03-10 antes de concluir 01
3. ❌ Criar código que viole o modelo multi-tenant
4. ❌ Acessar repositório diretamente de cron jobs (sempre via use case)
5. ❌ Ignorar validações de tenant_id em queries SQL
6. ❌ Usar npm ao invés de pnpm no frontend
7. ❌ Criar migrations novas sem seguir padrão existente
8. ❌ Hardcodar valores que devem ser configuráveis

### ✅ SEMPRE FAÇA:

1. ✅ Consulte `02-backlog.md` antes de iniciar qualquer tarefa
2. ✅ Valide com `04-checklist-dev.md` antes de considerar "pronto"
3. ✅ Execute testes com `05-checklist-qa.md`
4. ✅ Mantenha aderência ao Design System (MUI + tokens)
5. ✅ Use sqlc para queries SQL no backend Go
6. ✅ Use Zod + React Hook Form para formulários frontend
7. ✅ Documente decisões arquiteturais importantes (ADR)
8. ✅ Mantenha cobertura de testes > 70%

---

## 🆘 Suporte

### Dúvidas sobre:

- **Banco de Dados**: `DATABASE_MIGRATIONS_COMPLETED.md`
- **Backend Go**: `docs/04-backend/GUIA_DEV_BACKEND.md`
- **Frontend Next.js**: `docs/03-frontend/GUIA_FRONTEND.md`
- **Arquitetura**: `docs/02-arquitetura/ARQUITETURA.md`
- **Design System**: `docs/03-frontend/DESIGN_SYSTEM.md`
- **Integrações**: `docs/02-arquitetura/INTEGRACOES_EXTERNAS.md`
- **IA (Copilot)**: `.github/copilot-instructions.md`

### Dúvidas gerais:

- Abra uma issue no repositório
- Consulte o Tech Lead / Arquiteto-Chefe

---

## 📊 Dashboard de Progresso

Atualize esta seção conforme avançar:

```
┌─────────────────────────────────────────────────────────┐
│  PROGRESSO GERAL DO PROJETO                            │
├─────────────────────────────────────────────────────────┤
│  █████████████████████░░░░░░░░░░░░░░░░░  50% (Atualizado)│
│                                                         │
│  Etapa 0 (Pré):         ████████████████████████ 100%  │
│  Etapa 1 (Base):        ████████████████████████ 100%  │
│  Etapa 2 (OPS):         ████████████████████████ 100%  │
│  Etapa 3 (Financeiro):  ███████████████████████░  95%  │
│  Etapa 4 (Estoque):     ░░░░░░░░░░░░░░░░░░░░░░░░   0%  │
│  Etapa 5 (Metas):       ████████████████████████ 100%  │
│  Etapa 6 (Precific):    ████████████████████████ 100%  │
│  Etapa 7 (Lançamento):  ░░░░░░░░░░░░░░░░░░░░░░░░   0%  │
│  Etapa 8 (Monitor):     ░░░░░░░░░░░░░░░░░░░░░░░░   0%  │
│  Etapa 9 (Evolução):    ░░░░░░░░░░░░░░░░░░░░░░░░   0%  │
│  Etapa 10 (Agenda):     ░░░░░░░░░░░░░░░░░░░░░░░░   0%  │
└─────────────────────────────────────────────────────────┘
```

---

## ✅ Checklist Rápido de Início

Marque conforme avançar:

- [ ] Li este arquivo (`00-GUIA_NAVEGACAO.md`)
- [ ] Li `INDICE_TAREFAS.md`
- [ ] Li `DATABASE_MIGRATIONS_COMPLETED.md`
- [ ] Li `CONCLUIR/README.md`
- [ ] Li `CONCLUIR/00-ANALISE_SISTEMA_ATUAL.md`
- [ ] Entendi que NÃO posso pular `01-BLOQUEIOS-BASE`
- [ ] Entendi a estrutura de pastas e arquivos
- [ ] Entendi as regras críticas (SEMPRE/NUNCA)
- [ ] Configurei ambiente de desenvolvimento
- [ ] Testei `make dev` (backend + frontend)
- [ ] Li `.github/copilot-instructions.md`
- [ ] Pronto para começar `CONCLUIR/01-backend-domain-entities.md`

---

**Última Atualização:** 21/11/2025
**Próxima Revisão:** Após conclusão de 01-BLOQUEIOS-BASE

---

**BOA SORTE! 🚀**
