# 🎯 ROADMAP MILITAR — NEXO v1.0 → v2.0

**Emissão:** 22/11/2025
**Responsável:** Chief Engineering Officer
**Validade:** Até 20/12/2026
**Classificação:** CONFIDENCIAL - USO INTERNO

---

## ⚠️ AVISO CRÍTICO

Este documento é um **plano de combate técnico**, não um roadmap de apresentação.

**Regras invioláveis:**

- ❌ Escopo congelado após aprovação de cada release
- ❌ Datas NÃO são negociáveis sem autorização CEO
- ❌ Mudanças de escopo exigem ADR formal
- ❌ "Parcialmente pronto" = NÃO PRONTO
- ✅ Gate de aprovação obrigatório em CADA milestone

---

## 📊 ESTADO ATUAL DO SISTEMA (24/11/2025)

### Infraestrutura

- ✅ **Banco de Dados:** 100% (42 tabelas, migrations 001-042)
- ✅ **Neon PostgreSQL:** Configurado e operacional
- ✅ **Backend Go:** Estrutura base (Clean Architecture)
- ✅ **Frontend Next.js:** App Router configurado

### Progresso Técnico

| Componente                   | Status       | % Completo |
| ---------------------------- | ------------ | ---------- |
| **Domínio (Entities + VOs)** | ✅ Concluído | 100%       |
| **Repository Ports**         | ✅ Concluído | 100%       |
| **Repositórios PostgreSQL**  | ✅ Concluído | 100%       |
| **Use Cases**                | ✅ Concluído | 100%       |
| **Handlers HTTP**            | ✅ Concluído | 100%       |
| **Cron Jobs**                | ✅ Concluído | 100%       |
| **Frontend Services**        | ✅ Concluído | 100%       |
| **React Hooks**              | ✅ Concluído | 100%       |
| **LGPD**                     | ✅ Concluído | 100%       |
| **Backup Automático**        | ✅ Concluído | 100%       |
| **Módulo Financeiro**        | ✅ Concluído | 100%       |
| **Módulo Metas**             | ✅ Concluído | 100%       |
| **Módulo Precificação**      | ✅ Concluído | 100%       |

**Progresso Global MVP:** 85%

---

## 🎯 RELEASES OFICIAIS

### v1.0.0 — MVP CORE

- **Entrega:** 05/12/2025 (13 dias úteis)
- **Escopo:** CONGELADO
- **Criticidade:** MÁXIMA

### v1.1.0 — FIDELIDADE + GAMIFICAÇÃO

- **Início:** 12/12/2025
- **Entrega:** 10/02/2026 (42 dias úteis)
- **Escopo:** CONGELADO

### v1.2.0 — RELATÓRIOS AVANÇADOS

- **Início:** 11/02/2026
- **Entrega:** 30/03/2026 (33 dias úteis)
- **Escopo:** CONGELADO

### v2.0 — REDE/FRANQUIA + IA

- **Início Planejamento:** 10/04/2026
- **Estimativa:** Q4 2026

---

## 🔥 v1.0.0 — MVP CORE (05/12/2025)

**Prazo:** 13 dias úteis | **Status:** BLOQUEADOR CRÍTICO

### GATE 0: PRÉ-REQUISITOS (CONCLUÍDO - 24/11/2025)

- ✅ Banco de dados migrado (42 tabelas)
- ✅ Domínio completo (19 entidades)
- ✅ Ports definidas (11 interfaces)
- ✅ 100% repositórios implementados (11/11)
- ✅ 100% handlers HTTP prontos (44 endpoints)
- ✅ Frontend Services completo (7 services)
- ✅ React Hooks completo (43 hooks)
- ✅ Cron Jobs implementados (6/6)
- ✅ LGPD endpoints funcionais (4/4)
- ✅ Backup automático configurado

**Decisão GO/NO-GO:** ✅ APROVADO E CONCLUÍDO

---

### SEMANA 1: 22/11 - 29/11/2025 (CRÍTICA)

#### Milestone 1.1: Completar Base Técnica ✅ CONCLUÍDO

**Entrega:** 24/11/2025 (Domingo) - ANTECIPADO

| Tarefa                                          | Owner   | Horas | Status      | Conclusão   |
| ----------------------------------------------- | ------- | ----- | ----------- | ----------- |
| T-003-A: Completar 9 repositórios restantes     | Backend | 16h   | ✅ COMPLETO | 22/11 17:00 |
| T-003-B: Testes integração (tenant isolation)   | Backend | 4h    | ✅ COMPLETO | 22/11 18:00 |
| T-005-A: Corrigir handlers HTTP (Input structs) | Backend | 6h    | ✅ COMPLETO | 23/11 14:00 |
| T-005-B: Implementar endpoints GET/PUT/DELETE   | Backend | 10h   | ✅ COMPLETO | 24/11 12:00 |

**Checkpoint:** ✅ Aprovado em 24/11 12:30

**Critérios de Aprovação:**

- ✅ 11/11 repositórios funcionais
- ✅ Testes tenant isolation passando 100%
- ✅ Handlers HTTP compilando sem erros
- ✅ Endpoints CRUD completos testados (44 endpoints ativos)
- ✅ Frontend Services e Hooks implementados

**Resultado:** ✅ SUCESSO - Entregue 1 dia antes do prazo

---

#### Milestone 1.2: LGPD + Backup ✅ CONCLUÍDO

**Entrega:** 24/11/2025 (Domingo) - ANTECIPADO 3 DIAS

| Tarefa                                    | Owner    | Horas | Status      | Conclusão   |
| ----------------------------------------- | -------- | ----- | ----------- | ----------- |
| T-LGPD-001: Endpoint DELETE /me           | Backend  | 3h    | ✅ COMPLETO | 24/11 14:00 |
| T-LGPD-002: Endpoint GET /me/export       | Backend  | 3h    | ✅ COMPLETO | 24/11 15:00 |
| T-LGPD-003: Banner consentimento frontend | Frontend | 4h    | ✅ COMPLETO | 24/11 18:00 |
| T-LGPD-004: Privacy Policy page           | Frontend | 2h    | ✅ COMPLETO | 24/11 19:00 |
| T-OPS-001: GitHub Actions backup workflow | DevOps   | 4h    | ✅ COMPLETO | 24/11 16:00 |
| T-OPS-002: Disaster Recovery Runbook      | DevOps   | 2h    | ✅ COMPLETO | 24/11 17:00 |

**Checkpoint:** ✅ Aprovado em 24/11 20:00

**Critérios de Aprovação:**

- ✅ LGPD endpoints funcionais (4/4)
- ✅ Backup automático rodando (GitHub Actions)
- ✅ Restore testado com sucesso (8 scripts QA)
- ✅ Banner frontend funcionando
- ✅ Privacy Policy completa (600 linhas)
- ✅ 110+ casos de teste QA implementados

**Resultado:** ✅ SUCESSO ANTECIPADO - Entregue 3 dias antes com escopo expandido

---

#### Milestone 1.3: Frontend Services + Hooks ✅ CONCLUÍDO

**Entrega:** 24/11/2025 (Domingo) - ANTECIPADO 5 DIAS

| Tarefa                                         | Owner    | Horas | Status      | Conclusão   |
| ---------------------------------------------- | -------- | ----- | ----------- | ----------- |
| T-007-A: Services (DRE, Fluxo, Payables, etc.) | Frontend | 8h    | ✅ COMPLETO | 23/11 16:00 |
| T-008-A: React Query hooks (43 hooks)          | Frontend | 12h   | ✅ COMPLETO | 24/11 10:00 |
| T-007-B: Tratamento erros padronizado          | Frontend | 2h    | ✅ COMPLETO | 24/11 11:00 |
| T-009-A: Dashboard Financeiro (9 arquivos)     | Frontend | 6h    | ✅ COMPLETO | 24/11 20:00 |

**Checkpoint:** ✅ Aprovado em 24/11 20:30

**Critérios de Aprovação:**

- ✅ Services consumindo API corretamente (7 services)
- ✅ Hooks gerenciando estado/cache (43 hooks)
- ✅ Erro handling funcionando
- ✅ Dashboard completo com 8 cards + 4 métricas + 4 gráficos
- ✅ Recharts integrado (3.5.0)
- ✅ TypeScript sem erros

**Resultado:** ✅ SUCESSO ANTECIPADO - Entregue 5 dias antes com escopo expandido (Dashboard Financeiro)

---

### SEMANA 2: 02/12 - 05/12/2025 (FINALIZAÇÃO)

#### Milestone 2.1: Módulo Financeiro Completo

**Entrega:** 03/12/2025 (Terça-feira) 18:00

| Tarefa                                   | Owner    | Horas | Deadline    | Bloqueadores |
| ---------------------------------------- | -------- | ----- | ----------- | ------------ |
| T-FIN-001: Telas DRE Mensal              | Frontend | 6h    | 03/12 12:00 | M1.3         |
| T-FIN-002: Telas Fluxo Compensado        | Frontend | 6h    | 03/12 16:00 | M1.3         |
| T-FIN-003: Telas Contas Pagar/Receber    | Frontend | 8h    | 03/12 18:00 | M1.3         |
| T-FIN-004: Cron job DRE (dia 1 do mês)   | Backend  | 3h    | 03/12 18:00 | M1.1         |
| T-FIN-005: Cron job Fluxo Diário (00:05) | Backend  | 2h    | 03/12 18:00 | M1.1         |

**Checkpoint:** Teste end-to-end financeiro 03/12 18:30

**Critérios de Aprovação:**

- ✅ DRE gera corretamente
- ✅ Fluxo calcula saldos corretos
- ✅ Contas pagar/receber funcionais
- ✅ Cron jobs agendados corretamente

---

#### Milestone 2.2: Estoque + Metas + Precificação

**Entrega:** 04/12/2025 (Quarta-feira) 18:00

| Tarefa                                       | Owner      | Horas | Deadline    | Bloqueadores |
| -------------------------------------------- | ---------- | ----- | ----------- | ------------ |
| T-EST-001: CRUD Estoque (backend + frontend) | Full-Stack | 10h   | 04/12 14:00 | M2.1         |
| T-MET-001: CRUD Metas (backend + frontend)   | Full-Stack | 8h    | 04/12 16:00 | M2.1         |
| T-PRC-001: Simulador Precificação            | Full-Stack | 6h    | 04/12 18:00 | M2.1         |

**Checkpoint:** Teste módulos auxiliares 04/12 18:30

**Critérios de Aprovação:**

- ✅ Estoque não permite negativo
- ✅ Metas calculam progresso
- ✅ Precificação calcula margem correta

---

#### Milestone 2.3: Agendamento (CRÍTICO)

**Entrega:** 05/12/2025 (Quinta-feira) 18:00

| Tarefa                                | Owner    | Horas | Deadline    | Bloqueadores |
| ------------------------------------- | -------- | ----- | ----------- | ------------ |
| T-AGE-001: Backend agendamento (CRUD) | Backend  | 8h    | 05/12 12:00 | M2.2         |
| T-AGE-002: Frontend calendário visual | Frontend | 10h   | 05/12 18:00 | T-AGE-001    |
| T-AGE-003: Integração Google Agenda   | Backend  | 4h    | 05/12 18:00 | T-AGE-001    |
| T-AGE-004: Bloqueio conflitos         | Backend  | 2h    | 05/12 18:00 | T-AGE-001    |

**Checkpoint:** Demo agendamento funcional 05/12 18:30

**Critérios de Aprovação:**

- ✅ Agenda visual funcionando
- ✅ Sem conflitos de horário
- ✅ Google Agenda sincronizando
- ✅ Drag & drop operacional

---

### GATE 1: APROVAÇÃO v1.0.0 (05/12/2025 - 19:00)

**Checklist Obrigatório:**

#### Funcionalidades (100% ou FALHA)

- [ ] Agendamento funcional end-to-end
- [ ] Lista da Vez operacional
- [ ] Financeiro (DRE + Fluxo + Contas) funcionando
- [ ] Comissões calculando corretamente
- [ ] Estoque CRUD operacional
- [ ] Assinaturas Asaas integradas
- [ ] CRM básico funcional
- [ ] Relatórios mensais gerados
- [ ] Permissões (RBAC) funcionando

#### Qualidade (Mínimo ou FALHA)

- [ ] Cobertura testes backend ≥70%
- [ ] Cobertura testes frontend ≥60%
- [ ] Testes E2E ≥80% passando
- [ ] Performance p95 <300ms
- [ ] Zero erros críticos

#### Compliance (100% ou FALHA)

- [ ] LGPD endpoints funcionando
- [ ] Backup automático rodando
- [ ] Privacy Policy publicada
- [ ] Multi-tenant 100% isolado

#### Operacional (100% ou FALHA)

- [ ] Deploy staging executado
- [ ] Smoke tests passando
- [ ] Monitoramento configurado
- [ ] Alertas funcionando
- [ ] Documentação completa

**Decisão GO/NO-GO:**

- ✅ APROVADO → Deploy produção 06/12/2025
- ❌ REPROVADO → Atraso de 1 semana + root cause analysis

---

## 🎮 v1.1.0 — FIDELIDADE + GAMIFICAÇÃO (10/02/2026)

**Prazo:** 42 dias úteis (12/12/2025 - 10/02/2026)

### FASE 1: PLANEJAMENTO (12/12 - 20/12/2025)

#### Milestone 3.1: Design & Especificação

**Entrega:** 16/12/2025

| Tarefa                | Owner     | Deadline | Bloqueadores   |
| --------------------- | --------- | -------- | -------------- |
| UX/UI Cashback        | Design    | 13/12    | v1.0.0 GO-LIVE |
| UX/UI Gamificação     | Design    | 16/12    | UX Cashback    |
| Especificação técnica | Tech Lead | 16/12    | UX completo    |

#### Milestone 3.2: Aprovação Escopo

**Entrega:** 20/12/2025

| Tarefa               | Owner     | Deadline | Bloqueadores |
| -------------------- | --------- | -------- | ------------ |
| Review Product Owner | PM        | 18/12    | M3.1         |
| Review Tech Lead     | Tech Lead | 19/12    | M3.1         |
| Congelamento escopo  | CEO       | 20/12    | Reviews      |

**GATE 2:** Escopo congelado ou BLOQUEIO total do projeto

---

### FASE 2: DESENVOLVIMENTO (06/01 - 31/01/2026)

#### Milestone 3.3: Backend Cashback

**Entrega:** 13/01/2026 (5 dias úteis)

| Tarefa                  | Horas | Deadline |
| ----------------------- | ----- | -------- |
| Entidade Cashback + VOs | 8h    | 09/01    |
| Repository + Use Cases  | 12h   | 10/01    |
| Endpoints HTTP          | 8h    | 13/01    |
| Cron job expiração      | 4h    | 13/01    |

#### Milestone 3.4: Backend Gamificação

**Entrega:** 20/01/2026 (5 dias úteis)

| Tarefa                | Horas | Deadline |
| --------------------- | ----- | -------- |
| Entidade BarbeiroXP   | 8h    | 16/01    |
| Cálculo de níveis     | 10h   | 17/01    |
| Use Cases gamificação | 10h   | 20/01    |
| Endpoints HTTP        | 6h    | 20/01    |

#### Milestone 3.5: Frontend Web

**Entrega:** 27/01/2026 (5 dias úteis)

| Tarefa                      | Horas | Deadline |
| --------------------------- | ----- | -------- |
| Telas configuração cashback | 8h    | 23/01    |
| Dashboard gamificação       | 10h   | 24/01    |
| Ranking barbeiros           | 6h    | 27/01    |
| Integração completa         | 8h    | 27/01    |

#### Milestone 3.6: Metas Avançadas

**Entrega:** 31/01/2026 (3 dias úteis)

| Tarefa                    | Horas | Deadline |
| ------------------------- | ----- | -------- |
| Cálculo metas automáticas | 8h    | 29/01    |
| Alertas de desvio         | 6h    | 30/01    |
| Dashboard progresso       | 8h    | 31/01    |

---

### FASE 3: TESTES & DEPLOY (03/02 - 10/02/2026)

#### Milestone 3.7: QA Completo

**Entrega:** 07/02/2026

| Tarefa            | Horas | Deadline |
| ----------------- | ----- | -------- |
| Testes unitários  | 12h   | 05/02    |
| Testes integração | 10h   | 06/02    |
| Testes E2E        | 8h    | 07/02    |
| Correção bugs     | 16h   | 07/02    |

#### Milestone 3.8: Deploy Produção

**Entrega:** 10/02/2026

| Tarefa          | Horas | Deadline |
| --------------- | ----- | -------- |
| Deploy staging  | 2h    | 08/02    |
| Smoke tests     | 4h    | 09/02    |
| Deploy produção | 4h    | 10/02    |
| Monitoramento   | 4h    | 10/02    |

**GATE 3:** v1.1.0 APROVADO ou ROLLBACK

**Métricas de Sucesso (90 dias pós-release):**

- Churn mensal <10%
- LTV >R$ 1.200
- NPS Barbeiros >8

---

## 📊 v1.2.0 — RELATÓRIOS AVANÇADOS (30/03/2026)

**Prazo:** 33 dias úteis (11/02/2026 - 30/03/2026)

### FASE 1: BACKEND BI (11/02 - 28/02/2026)

#### Milestone 4.1: KPIs Complexos

**Entrega:** 21/02/2026 (8 dias úteis)

| Tarefa                  | Horas | Deadline |
| ----------------------- | ----- | -------- |
| Taxa ocupação (backend) | 10h   | 17/02    |
| Taxa retorno (backend)  | 8h    | 18/02    |
| Agregações otimizadas   | 12h   | 20/02    |
| Cache Redis relatórios  | 8h    | 21/02    |

#### Milestone 4.2: Comparativos Avançados

**Entrega:** 28/02/2026 (5 dias úteis)

| Tarefa                 | Horas | Deadline |
| ---------------------- | ----- | -------- |
| Comparativo trimestral | 10h   | 26/02    |
| Detecção sazonalidade  | 8h    | 27/02    |
| Projeções (algoritmo)  | 10h   | 28/02    |
| Endpoints relatórios   | 6h    | 28/02    |

---

### FASE 2: FRONTEND DASHBOARDS (03/03 - 14/03/2026)

#### Milestone 4.3: Dashboards Interativos

**Entrega:** 10/03/2026 (6 dias úteis)

| Tarefa                          | Horas | Deadline |
| ------------------------------- | ----- | -------- |
| Gráficos interativos (Recharts) | 12h   | 06/03    |
| Filtros avançados               | 8h    | 07/03    |
| Exportação PDF/CSV/Excel        | 10h   | 10/03    |

#### Milestone 4.4: Precificação Inteligente

**Entrega:** 14/03/2026 (3 dias úteis)

| Tarefa             | Horas | Deadline |
| ------------------ | ----- | -------- |
| Simulador completo | 12h   | 13/03    |
| UI precificação    | 8h    | 14/03    |

---

### FASE 3: APPS MOBILE (17/03 - 28/03/2026)

#### Milestone 4.5: App Barbeiro

**Entrega:** 24/03/2026 (6 dias úteis)

| Tarefa                | Horas | Deadline | Stack    |
| --------------------- | ----- | -------- | -------- |
| Setup React Native    | 8h    | 19/03    | RN 0.73  |
| Agenda própria        | 10h   | 20/03    | -        |
| Comissões + metas     | 10h   | 21/03    | -        |
| Gamificação + ranking | 8h    | 24/03    | -        |
| Push notifications    | 6h    | 24/03    | Firebase |

#### Milestone 4.6: App Cliente

**Entrega:** 28/03/2026 (3 dias úteis)

| Tarefa               | Horas | Deadline |
| -------------------- | ----- | -------- |
| Agendamento mobile   | 10h   | 26/03    |
| Histórico + cashback | 8h    | 27/03    |
| Avaliações           | 6h    | 28/03    |

---

### FASE 4: DEPLOY (30/03/2026)

#### Milestone 4.7: Go-Live v1.2.0

**Entrega:** 30/03/2026

| Tarefa               | Horas | Deadline |
| -------------------- | ----- | -------- |
| Testes E2E completos | 12h   | 29/03    |
| Deploy staging       | 4h    | 29/03    |
| Deploy produção      | 6h    | 30/03    |
| Publish app stores   | 4h    | 30/03    |

**GATE 4:** v1.2.0 APROVADO

**Métricas de Sucesso (60 dias pós-release):**

- 80% clientes usando relatórios
- 60% usando apps mobile
- Apps >4.5 estrelas

---

## 🏢 v2.0 — REDE/FRANQUIA + IA (20/12/2026)

**Prazo:** 34 semanas (10/04/2026 - 20/12/2026)

### FASE 1: PLANEJAMENTO (10/04 - 30/04/2026)

#### Milestone 5.1: Research & Design

**Entrega:** 30/04/2026

| Tarefa                       | Owner        | Duração   |
| ---------------------------- | ------------ | --------- |
| Research IA (time series)    | Data Science | 2 semanas |
| Design multi-tenant avançado | Arquiteto    | 2 semanas |
| Integrações (Open Banking)   | Tech Lead    | 2 semanas |

---

### FASE 2: NOTAS FISCAIS (05/05 - 30/06/2026)

#### Milestone 5.2: Gateway NFSe/NFe

**Entrega:** 30/06/2026 (8 semanas)

| Tarefa                     | Duração   |
| -------------------------- | --------- |
| Integração eNotas.io       | 3 semanas |
| Backend emissão automática | 2 semanas |
| Frontend configuração      | 2 semanas |
| Testes certificação        | 1 semana  |

---

### FASE 3: INTEGRAÇÕES BANCÁRIAS (01/07 - 31/08/2026)

#### Milestone 5.3: Open Banking

**Entrega:** 31/08/2026 (8 semanas)

| Tarefa                       | Duração   |
| ---------------------------- | --------- |
| Integração bancos (6 bancos) | 4 semanas |
| Conciliação automática       | 2 semanas |
| UI conciliação               | 1 semana  |
| Testes + homologação         | 1 semana  |

---

### FASE 4: FRANQUIAS (01/09 - 31/10/2026)

#### Milestone 5.4: Multi-Unidade Avançado

**Entrega:** 31/10/2026 (8 semanas)

| Tarefa                | Duração   |
| --------------------- | --------- |
| Painel consolidado    | 3 semanas |
| Repasse royalties     | 2 semanas |
| Comparativos unidades | 2 semanas |
| Testes multi-tenant   | 1 semana  |

---

### FASE 5: IA PREDITIVA (01/11 - 10/12/2026)

#### Milestone 5.5: Microserviço IA

**Entrega:** 10/12/2026 (6 semanas)

| Tarefa                  | Duração   | Stack                  |
| ----------------------- | --------- | ---------------------- |
| Setup Python + ML       | 1 semana  | FastAPI + Scikit-learn |
| Modelo previsão demanda | 2 semanas | Prophet                |
| Modelo churn            | 1 semana  | RandomForest           |
| API REST integração     | 1 semana  | -                      |
| Testes acurácia         | 1 semana  | ≥70% accuracy          |

---

### FASE 6: API PÚBLICA (11/12 - 18/12/2026)

#### Milestone 5.6: Developer Platform

**Entrega:** 18/12/2026 (1 semana)

| Tarefa           | Duração |
| ---------------- | ------- |
| OAuth2 completo  | 3 dias  |
| Swagger/OpenAPI  | 2 dias  |
| SDKs (JS/Python) | 2 dias  |

---

### FASE 7: GO-LIVE v2.0 (20/12/2026)

**GATE 5:** v2.0 FINAL

**Métricas de Sucesso (6 meses pós-release):**

- 40% clientes multi-unidade
- MRR >R$ 200k
- IA acurácia >70%

---

## 🚨 BLOQUEADORES CRÍTICOS & CONTINGÊNCIAS

### v1.0.0

| Bloqueador               | Probabilidade | Impacto | Contingência                      |
| ------------------------ | ------------- | ------- | --------------------------------- |
| **Atraso repositórios**  | ALTA          | CRÍTICO | Adicionar 1 dev backend full-time |
| **Bug tenant isolation** | MÉDIA         | CRÍTICO | Code freeze total até fix         |
| **Google Agenda falha**  | BAIXA         | ALTO    | Remover integração do MVP         |
| **Asaas indisponível**   | BAIXA         | ALTO    | Modo manual temporário            |

### v1.1.0

| Bloqueador                 | Probabilidade | Impacto | Contingência                |
| -------------------------- | ------------- | ------- | --------------------------- |
| **Cálculo XP incorreto**   | MÉDIA         | MÉDIO   | Audit completo antes deploy |
| **Cashback gera negativo** | MÉDIA         | ALTO    | Validação extra + testes    |

### v1.2.0

| Bloqueador                 | Probabilidade | Impacto | Contingência          |
| -------------------------- | ------------- | ------- | --------------------- |
| **Apps mobile rejeitados** | MÉDIA         | ALTO    | PWA como fallback     |
| **Relatórios lentos**      | ALTA          | MÉDIO   | Cache Redis agressivo |

### v2.0

| Bloqueador                | Probabilidade | Impacto | Contingência              |
| ------------------------- | ------------- | ------- | ------------------------- |
| **IA baixa acurácia**     | ALTA          | MÉDIO   | Remover IA ou simplificar |
| **Open Banking bloqueio** | MÉDIA         | ALTO    | Importação manual CSV     |

---

## 📅 LINHA DO TEMPO GANTT (TEXTO)

```
2025
NOV: ████████░░░░░░░░░░░░ v1.0 (Semana 1)
DEC: ██████████████████░░ v1.0 DONE + v1.1 Planejamento

2026
JAN: ████████████████████ v1.1 Desenvolvimento
FEV: ██████████████████░░ v1.1 DONE + v1.2 Backend
MAR: ████████████████████ v1.2 Frontend + Apps
ABR: ██████░░░░░░░░░░░░░░ v1.2 DONE + v2.0 Research
MAI: ████████████░░░░░░░░ v2.0 Notas Fiscais (início)
JUN: ████████████████░░░░ v2.0 Notas Fiscais (fim)
JUL: ████████████████░░░░ v2.0 Open Banking (início)
AGO: ████████████████████ v2.0 Open Banking (fim)
SET: ████████████████░░░░ v2.0 Franquias (início)
OUT: ████████████████████ v2.0 Franquias (fim)
NOV: ████████████████░░░░ v2.0 IA Preditiva
DEC: ██████████████████░░ v2.0 DONE (20/12)
```

---

## 🎯 CRITÉRIOS MILITARES DE "PRONTO"

### Definição ABSOLUTA

- ✅ **PRONTO** = 100% funcional, testado, documentado, deployado
- ❌ **PARCIAL** = NÃO PRONTO (é falha)
- ❌ **90%** = NÃO PRONTO (é falha)

### Checklist por Feature

| Item                    | Obrigatório |
| ----------------------- | ----------- |
| Código no repositório   | SIM         |
| Testes unitários ≥70%   | SIM         |
| Testes integração       | SIM         |
| Code review aprovado    | SIM         |
| Deploy em staging       | SIM         |
| Smoke tests passando    | SIM         |
| Documentação atualizada | SIM         |
| Deploy em produção      | SIM         |
| Monitoramento ativo     | SIM         |

**Qualquer NÃO = Feature NÃO PRONTA**

---

## 🔒 CONGELAMENTO DE ESCOPO

### v1.0.0

**Status:** ✅ CONGELADO (22/11/2025)

**Funcionalidades aprovadas:**

1. Agendamento
2. Lista da Vez
3. Financeiro básico
4. Comissões
5. Estoque essencial
6. Assinaturas Asaas
7. CRM básico
8. Relatórios mensais
9. Permissões

**Mudanças:** PROIBIDAS sem ADR + aprovação CEO

### v1.1.0

**Status:** ⚪ AGUARDANDO CONGELAMENTO (20/12/2025)

**Funcionalidades planejadas:**

1. Cashback
2. Gamificação
3. Metas avançadas

**Mudanças:** Permitidas até 20/12/2025

### v1.2.0

**Status:** ⚪ AGUARDANDO CONGELAMENTO (28/02/2026)

### v2.0

**Status:** ⚪ PLANEJAMENTO (Research fase)

---

## 📊 MÉTRICAS DE ACOMPANHAMENTO

### Dashboard Semanal Obrigatório

**Toda Segunda-feira 09:00:**

| Métrica                 | Target         | Atual | Status |
| ----------------------- | -------------- | ----- | ------ |
| Tasks completas semana  | 100% planejado | -     | -      |
| Bugs críticos abertos   | 0              | -     | -      |
| Code coverage backend   | ≥70%           | -     | -      |
| Code coverage frontend  | ≥60%           | -     | -      |
| Deploy frequency        | ≥2x/semana     | -     | -      |
| MTTR (Mean Time Repair) | <4h            | -     | -      |

### KPIs por Release

**v1.0.0:**

- Velocidade deploy: <30min
- Uptime primeira semana: >99%
- Bugs críticos primeira semana: 0

**v1.1.0:**

- Churn mensal: <10%
- LTV: >R$ 1.200
- NPS Barbeiros: >8

**v1.2.0:**

- Uso relatórios: >80%
- Uso apps: >60%
- App rating: >4.5★

**v2.0:**

- Clientes multi-unidade: >40%
- MRR: >R$ 200k
- IA accuracy: >70%

---

## ⚠️ GATES DE APROVAÇÃO

### Gate 0: Pré-Requisitos (22/11/2025)

- ✅ APROVADO

### Gate 1: v1.0.0 Final (05/12/2025 19:00)

- ⚪ PENDENTE
- **Critério:** 100% funcionalidades ou FALHA

### Gate 2: v1.1.0 Escopo (20/12/2025)

- ⚪ PENDENTE
- **Critério:** Congelamento ou BLOQUEIO

### Gate 3: v1.1.0 Final (10/02/2026)

- ⚪ PENDENTE
- **Critério:** Métricas sucesso ou ROLLBACK

### Gate 4: v1.2.0 Final (30/03/2026)

- ⚪ PENDENTE
- **Critério:** Apps publicados ou FALHA

### Gate 5: v2.0 Final (20/12/2026)

- ⚪ PENDENTE
- **Critério:** IA acurácia >70% ou DEGRADAÇÃO

---

## 🔥 PLANO DE AÇÃO EMERGENCIAL

### Situação: Atraso v1.0.0

**Trigger:** 03/12/2025 e progresso <90%

**Ações Imediatas:**

1. Code freeze total
2. War room diária (2x/dia)
3. Paralizar features não-críticas
4. Adicionar recursos (devs extras)
5. Comunicação CEO imediata

### Situação: Bug Crítico Produção

**Trigger:** Severity 1 (sistema inoperante)

**Ações Imediatas (SLA: 1h):**

1. Rollback automático
2. Incident commander ativado
3. Time on-call mobilizado
4. Comunicação clientes
5. Post-mortem obrigatório em 24h

### Situação: Escopo Creep

**Trigger:** Solicitação mudança escopo

**Ações:**

1. ADR formal obrigatório
2. Análise impacto (datas)
3. Aprovação CEO necessária
4. Comunicação stakeholders
5. Ajuste roadmap formal

---

## 📞 RESPONSABILIDADES

| Papel               | Responsável  | Contato |
| ------------------- | ------------ | ------- |
| CEO (Aprovar gates) | Andrey Viana | -       |
| Tech Lead           | TBD          | -       |
| Backend Lead        | TBD          | -       |
| Frontend Lead       | TBD          | -       |
| DevOps              | TBD          | -       |
| QA Lead             | TBD          | -       |

---

## 📄 APROVAÇÕES

| Versão | Data       | Aprovador          | Status      |
| ------ | ---------- | ------------------ | ----------- |
| 1.0    | 22/11/2025 | Andrey Viana (CEO) | ⚪ PENDENTE |

---

**FIM DO ROADMAP MILITAR**

**Próxima Revisão:** 25/11/2025 (Checkpoint Milestone 1.1)

---

**ASSINATURA DIGITAL:**

```
___________________________
Andrey Viana
CEO - NEXO
22/11/2025
```
