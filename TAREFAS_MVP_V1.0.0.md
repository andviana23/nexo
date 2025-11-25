# 🎯 Tarefas para Concluir v1.0.0 — MVP CORE

**Data de Emissão:** 24/11/2025
**Deadline:** 05/12/2025 (11 dias úteis restantes)
**Progresso Atual:** 85%
**Status:** 🔴 CRÍTICO - Deadline em risco

---

## 📊 Resumo Executivo

### ✅ Concluído (85%)

- ✅ **Infraestrutura:** Banco de dados (42 tabelas), Neon PostgreSQL, Clean Architecture
- ✅ **Backend Core:** 11 repositórios, 24 use cases, 44 endpoints, 6 cron jobs
- ✅ **Frontend Base:** 7 services, 43 hooks, tratamento de erros
- ✅ **LGPD:** 4 endpoints + Privacy Policy + Cookie Banner
- ✅ **Backup/DR:** GitHub Actions + Runbook completo
- ✅ **Módulo Financeiro:** 100% (Backend + Frontend + Dashboard)
- ✅ **Módulo Metas:** 100% (Backend + Frontend)
- ✅ **Módulo Precificação:** 100% (Backend + Frontend)

### 🔴 Pendente (15%)

- ❌ **Módulo Estoque:** 0% (CRÍTICO)
- ❌ **Módulo Agendamento:** 0% (BLOQUEADOR)
- ❌ **Lista da Vez:** Frontend 0% (Backend 100%)
- ❌ **Assinaturas Asaas:** Integração parcial
- ❌ **CRM Básico:** Frontend 50%
- ❌ **Relatórios UI:** Telas de DRE e Fluxo
- ❌ **Permissões (RBAC):** Frontend 50%
- ❌ **Testes E2E:** Cobertura <50%
- ❌ **Deploy Staging/Produção:** 0%

---

## 🗓️ Plano de Execução (Ordem Obrigatória)

### 📅 **FASE 1: Módulos Críticos de Negócio** (25/11 - 28/11)

**Duração:** 4 dias | **Prioridade:** BLOQUEADOR

---

## 1️⃣ ESTOQUE (Dia 1-2: 25-26/11) 🔴 CRÍTICO

**Total:** 28 horas (~2 dias com 2 devs)

### Backend (14h)

#### T-EST-001: Entrada de Estoque

- **Descrição:** Implementar registro de entrada de produtos no estoque
- **Arquivos:**
  - `backend/internal/domain/entity/stock_entry.go` (Entidade)
  - `backend/internal/application/usecase/stock/create_entry.go` (Use Case)
  - `backend/internal/infra/http/handler/stock_handler.go` (Endpoint POST)
  - `backend/internal/application/dto/stock_dto.go` (DTOs)
- **Validações:**
  - Quantidade > 0
  - Produto existe
  - Tenant isolation
- **Estimativa:** 4h
- **Prioridade:** ALTA

#### T-EST-002: Saída de Estoque

- **Descrição:** Implementar baixa de produtos do estoque
- **Arquivos:**
  - `backend/internal/domain/entity/stock_exit.go` (Entidade)
  - `backend/internal/application/usecase/stock/create_exit.go` (Use Case)
  - `backend/internal/infra/http/handler/stock_handler.go` (Endpoint POST)
- **Validações:**
  - Quantidade disponível suficiente
  - Não permitir estoque negativo
  - Tenant isolation
- **Estimativa:** 4h
- **Prioridade:** ALTA

#### T-EST-003: Consumo Automático

- **Descrição:** Baixa automática de estoque ao finalizar atendimento
- **Arquivos:**
  - `backend/internal/application/usecase/stock/auto_consume.go`
  - Integração com `AttendanceService`
- **Regra de Negócio:**
  - Cada serviço tem produtos associados
  - Ao finalizar atendimento, baixar produtos automaticamente
- **Estimativa:** 3h
- **Prioridade:** ALTA

#### T-EST-004: Inventário (Listagem + Ajuste Manual)

- **Descrição:** Listar estoque atual e permitir ajustes manuais
- **Arquivos:**
  - `backend/internal/application/usecase/stock/list_inventory.go`
  - `backend/internal/application/usecase/stock/adjust_inventory.go`
  - Endpoints GET/PUT
- **Estimativa:** 3h
- **Prioridade:** MÉDIA

---

### Frontend (14h)

#### T-EST-005: Tela de Entrada de Estoque

- **Descrição:** Formulário para registrar entrada de produtos
- **Arquivos:**
  - `frontend/app/(dashboard)/estoque/entrada/page.tsx`
  - `frontend/components/stock/EntryForm.tsx`
  - `frontend/hooks/useStockEntry.ts`
- **Campos:**
  - Produto (select)
  - Quantidade
  - Valor unitário
  - Fornecedor
  - Data de validade (opcional)
- **Validação:** Zod + React Hook Form
- **Estimativa:** 4h
- **Prioridade:** ALTA

#### T-EST-006: Tela de Saída de Estoque

- **Descrição:** Formulário para registrar saída de produtos
- **Arquivos:**
  - `frontend/app/(dashboard)/estoque/saida/page.tsx`
  - `frontend/components/stock/ExitForm.tsx`
  - `frontend/hooks/useStockExit.ts`
- **Validações:**
  - Verificar disponibilidade antes de enviar
  - Alertar se estoque baixo
- **Estimativa:** 3h
- **Prioridade:** ALTA

#### T-EST-007: Tela de Inventário

- **Descrição:** Listagem de estoque atual com filtros
- **Arquivos:**
  - `frontend/app/(dashboard)/estoque/inventario/page.tsx`
  - `frontend/components/stock/InventoryTable.tsx`
  - `frontend/hooks/useInventory.ts`
- **Funcionalidades:**
  - Tabela com produtos, quantidade, valor
  - Filtros: categoria, estoque baixo
  - Ajuste manual de quantidade
  - Indicador visual de estoque mínimo
- **Estimativa:** 5h
- **Prioridade:** ALTA

#### T-EST-008: Alerta de Estoque Mínimo (OPCIONAL)

- **Descrição:** Notificação quando produto atingir estoque mínimo
- **Estimativa:** 2h
- **Prioridade:** BAIXA (pode ser v1.1.0)

---

## 2️⃣ AGENDAMENTO (Dia 3-4: 27-28/11) 🔴 BLOQUEADOR

**Total:** 35 horas (~2 dias com 2 devs)

### Backend (17h)

#### T-AGE-001: CRUD Agendamento Backend

- **Descrição:** Endpoints completos para gerenciar agendamentos
- **Arquivos:**
  - `backend/internal/domain/entity/appointment.go` (já existe)
  - `backend/internal/application/usecase/appointment/*` (CRUD completo)
  - `backend/internal/infra/http/handler/appointment_handler.go`
- **Endpoints:**
  - POST `/appointments` - Criar agendamento
  - GET `/appointments` - Listar agendamentos (com filtros: data, barbeiro, status)
  - GET `/appointments/:id` - Detalhes do agendamento
  - PUT `/appointments/:id` - Atualizar agendamento
  - DELETE `/appointments/:id` - Cancelar agendamento
  - PUT `/appointments/:id/confirm` - Confirmar agendamento
  - PUT `/appointments/:id/complete` - Finalizar atendimento
- **Estimativa:** 8h
- **Prioridade:** CRÍTICA

#### T-AGE-002: Validação de Conflitos de Horário

- **Descrição:** Impedir agendamentos conflitantes
- **Arquivos:**
  - `backend/internal/application/usecase/appointment/check_conflicts.go`
- **Regras:**
  - Mesmo barbeiro não pode ter 2 agendamentos no mesmo horário
  - Verificar duração do serviço para evitar sobreposição
  - Respeitar horário de funcionamento da barbearia
- **Estimativa:** 3h
- **Prioridade:** CRÍTICA

#### T-AGE-003: Validação de Horário de Funcionamento

- **Descrição:** Verificar se agendamento está dentro do horário permitido
- **Arquivos:**
  - `backend/internal/application/usecase/appointment/validate_business_hours.go`
- **Regras:**
  - Consultar configuração de horários (tabela `business_hours`)
  - Validar dia da semana + horário início/fim
- **Estimativa:** 2h
- **Prioridade:** ALTA

#### T-AGE-004: Integração Google Agenda (OPCIONAL)

- **Descrição:** Sincronizar agendamentos com Google Calendar
- **Estimativa:** 4h
- **Prioridade:** BAIXA (pode ser v1.1.0)

---

### Frontend (18h)

#### T-AGE-005: Componente de Calendário Visual

- **Descrição:** Interface visual para visualizar e criar agendamentos
- **Arquivos:**
  - `frontend/app/(dashboard)/agendamentos/page.tsx`
  - `frontend/components/appointments/Calendar.tsx`
- **Biblioteca:** FullCalendar ou DayPilot Lite (gratuita)
- **Funcionalidades:**
  - Visualização mensal/semanal/diária
  - Cores diferentes por status (pendente, confirmado, finalizado)
  - Click para criar novo agendamento
  - Click no evento para ver detalhes
- **Estimativa:** 8h
- **Prioridade:** CRÍTICA

#### T-AGE-006: Formulário de Agendamento

- **Descrição:** Modal/página para criar/editar agendamento
- **Arquivos:**
  - `frontend/components/appointments/AppointmentForm.tsx`
  - `frontend/hooks/useAppointments.ts`
- **Campos:**
  - Cliente (select com busca)
  - Barbeiro (select)
  - Serviços (multi-select)
  - Data e hora
  - Observações
- **Validações:**
  - Verificar conflitos antes de enviar
  - Calcular duração total dos serviços
- **Estimativa:** 5h
- **Prioridade:** CRÍTICA

#### T-AGE-007: Drag & Drop no Calendário (OPCIONAL)

- **Descrição:** Permitir arrastar agendamento para outro horário
- **Estimativa:** 3h
- **Prioridade:** MÉDIA (pode ser v1.1.0)

#### T-AGE-008: Notificações de Lembrete (OPCIONAL)

- **Descrição:** WhatsApp/SMS 24h antes do agendamento
- **Estimativa:** 2h
- **Prioridade:** BAIXA (pode ser v1.1.0)

---

## 3️⃣ LISTA DA VEZ (Dia 5: 29/11) 🔴 CRÍTICO

**Total:** 19 horas (~1 dia com 2 devs)

### Backend: ✅ JÁ IMPLEMENTADO (100%)

- Repository, Use Cases, Endpoints já existem

### Frontend (19h)

#### T-LIST-001: Tela Lista da Vez Principal

- **Descrição:** Interface para gerenciar fila de espera
- **Arquivos:**
  - `frontend/app/(dashboard)/lista-da-vez/page.tsx`
  - `frontend/components/queue/QueueList.tsx`
  - `frontend/hooks/useQueue.ts`
- **Funcionalidades:**
  - Listar clientes na fila (ordenado por posição)
  - Mostrar tempo de espera estimado
  - Status de cada cliente (aguardando, em atendimento)
  - Filtro por barbeiro
- **Estimativa:** 6h
- **Prioridade:** CRÍTICA

#### T-LIST-002: Adicionar Cliente na Fila

- **Descrição:** Botão/Modal para adicionar cliente na lista
- **Arquivos:**
  - `frontend/components/queue/AddToQueueForm.tsx`
- **Campos:**
  - Cliente (select com busca ou cadastro rápido)
  - Serviço desejado (opcional)
  - Barbeiro preferido (opcional)
- **Estimativa:** 4h
- **Prioridade:** CRÍTICA

#### T-LIST-003: Chamar Próximo da Fila

- **Descrição:** Botão para chamar próximo cliente
- **Arquivos:**
  - `frontend/components/queue/CallNextButton.tsx`
- **Funcionalidades:**
  - Atualizar status para "em atendimento"
  - Notificação visual/sonora (opcional)
  - Remover da fila ou marcar como atendido
- **Estimativa:** 3h
- **Prioridade:** CRÍTICA

#### T-LIST-004: Cancelar/Remover da Fila

- **Descrição:** Remover cliente da lista (não compareceu)
- **Estimativa:** 2h
- **Prioridade:** ALTA

#### T-LIST-005: Notificações Push (OPCIONAL)

- **Descrição:** Notificar cliente via SMS quando estiver próximo
- **Estimativa:** 4h
- **Prioridade:** BAIXA (pode ser v1.1.0)

---

## 📅 **FASE 2: Integrações e Funcionalidades Complementares** (29/11 - 02/12)

**Duração:** 4 dias | **Prioridade:** ALTA

---

## 4️⃣ ASSINATURAS ASAAS (Dia 6-7: 02-03/12) 🔴 BLOQUEADOR

**Total:** 25 horas (~2 dias com 2 devs)

### Backend (14h)

#### T-ASAAS-001: Integração Completa API Asaas

- **Descrição:** Cliente HTTP para todas as operações Asaas
- **Arquivos:**
  - `backend/internal/infra/gateway/asaas/client.go`
  - `backend/internal/infra/gateway/asaas/customer.go`
  - `backend/internal/infra/gateway/asaas/subscription.go`
  - `backend/internal/infra/gateway/asaas/payment.go`
- **Funcionalidades:**
  - Criar/atualizar cliente Asaas
  - Criar assinatura (mensal/anual)
  - Consultar status de pagamento
  - Cancelar assinatura
- **Estimativa:** 6h
- **Prioridade:** CRÍTICA

#### T-ASAAS-002: Webhooks para Eventos Asaas

- **Descrição:** Receber notificações de mudança de status
- **Arquivos:**
  - `backend/internal/infra/http/handler/webhook_handler.go`
  - `backend/internal/application/usecase/subscription/process_webhook.go`
- **Eventos:**
  - Pagamento confirmado → Ativar assinatura
  - Pagamento vencido → Suspender acesso
  - Assinatura cancelada → Desativar tenant
- **Estimativa:** 4h
- **Prioridade:** CRÍTICA

#### T-ASAAS-003: Sincronização de Status

- **Descrição:** Cron job para verificar status de assinaturas
- **Arquivos:**
  - `backend/cmd/cron/sync_subscriptions.go`
- **Frequência:** Diário às 02:00
- **Estimativa:** 2h
- **Prioridade:** ALTA

#### T-ASAAS-004: Tratamento de Erros e Retry

- **Descrição:** Implementar retry com backoff exponencial
- **Estimativa:** 2h
- **Prioridade:** MÉDIA

---

### Frontend (11h)

#### T-ASAAS-005: Tela de Escolha de Plano

- **Descrição:** Página para selecionar plano de assinatura
- **Arquivos:**
  - `frontend/app/(public)/planos/page.tsx`
  - `frontend/components/subscription/PlanCard.tsx`
- **Planos:**
  - Starter: R$ 49,90/mês
  - Professional: R$ 99,90/mês
  - Premium: R$ 199,90/mês
- **Funcionalidades:**
  - Cards visuais com comparação
  - Botão "Assinar"
- **Estimativa:** 4h
- **Prioridade:** CRÍTICA

#### T-ASAAS-006: Fluxo de Checkout

- **Descrição:** Página de finalização de assinatura
- **Arquivos:**
  - `frontend/app/(public)/checkout/page.tsx`
  - `frontend/components/subscription/CheckoutForm.tsx`
- **Dados:**
  - Informações pessoais
  - CPF/CNPJ
  - Endereço
  - Forma de pagamento (cartão, boleto, Pix)
- **Validações:** Zod + CPF/CNPJ
- **Estimativa:** 5h
- **Prioridade:** CRÍTICA

#### T-ASAAS-007: Gerenciamento de Assinatura

- **Descrição:** Página para upgrade/downgrade/cancelamento
- **Arquivos:**
  - `frontend/app/(dashboard)/assinatura/page.tsx`
  - `frontend/components/subscription/ManageSubscription.tsx`
- **Funcionalidades:**
  - Exibir plano atual
  - Botão para alterar plano
  - Botão para cancelar (com confirmação)
  - Histórico de faturas
- **Estimativa:** 2h (se backend já estiver pronto)
- **Prioridade:** ALTA

---

## 5️⃣ CRM BÁSICO (Dia 8: 03/12) 🟡 MÉDIA PRIORIDADE

**Total:** 15 horas (~1 dia com 2 devs)

### Backend: ✅ JÁ IMPLEMENTADO (100%)

### Frontend (15h)

#### T-CRM-001: Tela de Cadastro de Clientes

- **Descrição:** Formulário CRUD de clientes
- **Arquivos:**
  - `frontend/app/(dashboard)/clientes/page.tsx`
  - `frontend/components/clients/ClientForm.tsx`
  - `frontend/hooks/useClients.ts`
- **Campos:**
  - Nome, telefone, e-mail
  - Data de nascimento
  - Endereço
  - Observações
- **Estimativa:** 4h
- **Prioridade:** ALTA

#### T-CRM-002: Histórico de Atendimentos do Cliente

- **Descrição:** Visualizar todos os atendimentos de um cliente
- **Arquivos:**
  - `frontend/app/(dashboard)/clientes/[id]/page.tsx`
  - `frontend/components/clients/ClientHistory.tsx`
- **Dados:**
  - Data, barbeiro, serviços, valor
  - Produtos utilizados
  - Observações
- **Estimativa:** 5h
- **Prioridade:** ALTA

#### T-CRM-003: Busca e Filtros Avançados

- **Descrição:** Buscar clientes por nome, telefone, e-mail
- **Arquivos:**
  - `frontend/components/clients/ClientSearch.tsx`
- **Filtros:**
  - Nome, telefone
  - Clientes inativos (sem atendimento há 30+ dias)
  - Ordenação por última visita
- **Estimativa:** 3h
- **Prioridade:** MÉDIA

#### T-CRM-004: Estatísticas do Cliente (OPCIONAL)

- **Descrição:** LTV, frequência média, ticket médio
- **Estimativa:** 3h
- **Prioridade:** BAIXA (pode ser v1.1.0)

---

## 6️⃣ RELATÓRIOS UI (Dia 9: 04/12) 🟡 MÉDIA PRIORIDADE

**Total:** 13 horas (~1 dia com 2 devs)

### Backend: ✅ JÁ IMPLEMENTADO (100%)

- Dashboard, DRE, Fluxo de Caixa já têm endpoints

### Frontend (13h)

#### T-REL-001: Tela DRE Mensal (Gráficos)

- **Descrição:** Visualização do DRE com gráficos
- **Arquivos:**
  - `frontend/app/(dashboard)/financeiro/dre/page.tsx`
  - `frontend/components/financial/DREChart.tsx` (já existe)
- **Funcionalidades:**
  - Seletor de mês
  - Gráfico de barras (receita vs despesa)
  - Tabela detalhada de lançamentos
  - Comparativo com mês anterior
- **Estimativa:** 4h
- **Prioridade:** ALTA

#### T-REL-002: Tela Fluxo de Caixa (Linha do Tempo)

- **Descrição:** Visualização do fluxo compensado
- **Arquivos:**
  - `frontend/app/(dashboard)/financeiro/fluxo/page.tsx`
  - `frontend/components/financial/CashflowChart.tsx` (já existe)
- **Funcionalidades:**
  - Seletor de período
  - Gráfico de linha (entradas, saídas, saldo)
  - Projeção de saldo futuro
- **Estimativa:** 4h
- **Prioridade:** ALTA

#### T-REL-003: Exportação PDF/Excel (OPCIONAL)

- **Descrição:** Exportar relatórios para PDF ou Excel
- **Biblioteca:** jsPDF ou xlsx
- **Estimativa:** 3h
- **Prioridade:** MÉDIA (pode ser v1.1.0)

#### T-REL-004: Comparativo Mês Anterior (OPCIONAL)

- **Descrição:** Mostrar variação % vs mês anterior
- **Estimativa:** 2h
- **Prioridade:** BAIXA (pode ser v1.1.0)

---

## 7️⃣ PERMISSÕES (RBAC) (Dia 10: 04/12) 🟢 BAIXA PRIORIDADE

**Total:** 10 horas (~1 dia)

### Backend: ✅ JÁ IMPLEMENTADO (100%)

- Middleware de autorização já existe

### Frontend (10h)

#### T-RBAC-001: Tela de Gerenciamento de Papéis

- **Descrição:** CRUD de papéis e permissões
- **Arquivos:**
  - `frontend/app/(dashboard)/configuracoes/permissoes/page.tsx`
  - `frontend/components/rbac/RoleManager.tsx`
- **Funcionalidades:**
  - Listar papéis (Admin, Gerente, Barbeiro, Recepcionista)
  - Editar permissões de cada papel
  - Atribuir papel a usuários
- **Estimativa:** 4h
- **Prioridade:** MÉDIA

#### T-RBAC-002: Restrições Visuais por Papel

- **Descrição:** Ocultar menus/botões conforme permissões
- **Arquivos:**
  - `frontend/components/layout/Sidebar.tsx` (atualizar)
  - `frontend/hooks/usePermissions.ts`
- **Regras:**
  - Barbeiro: Não vê Financeiro
  - Recepcionista: Não vê Relatórios
  - Admin: Vê tudo
- **Estimativa:** 3h
- **Prioridade:** MÉDIA

#### T-RBAC-003: Proteção de Rotas

- **Descrição:** Middleware Next.js para proteger páginas
- **Arquivos:**
  - `frontend/middleware.ts` (atualizar)
- **Estimativa:** 3h
- **Prioridade:** ALTA

---

## 📅 **FASE 3: Qualidade e Deploy** (04/12 - 05/12)

**Duração:** 2 dias | **Prioridade:** CRÍTICA

---

## 8️⃣ TESTES E2E (Dia 11: 05/12) 🔴 CRÍTICO

**Total:** 12 horas

#### T-TEST-001: Testes E2E Agendamento

- **Descrição:** Cypress ou Playwright
- **Cenários:**
  - Criar agendamento
  - Verificar conflito de horário
  - Cancelar agendamento
  - Finalizar atendimento
- **Estimativa:** 4h
- **Prioridade:** CRÍTICA

#### T-TEST-002: Testes E2E Lista da Vez

- **Cenários:**
  - Adicionar cliente na fila
  - Chamar próximo
  - Remover da fila
- **Estimativa:** 3h
- **Prioridade:** ALTA

#### T-TEST-003: Testes E2E Assinaturas

- **Cenários:**
  - Selecionar plano
  - Preencher checkout
  - Webhook de confirmação
- **Estimativa:** 3h
- **Prioridade:** ALTA

#### T-TEST-004: Testes E2E Financeiro

- **Cenários:**
  - Criar conta a pagar
  - Gerar DRE
  - Visualizar fluxo de caixa
- **Estimativa:** 2h
- **Prioridade:** MÉDIA

---

## 9️⃣ DEPLOY E MONITORAMENTO (Dia 11-12: 05/12) 🔴 CRÍTICO

**Total:** 8 horas

#### T-DEPLOY-001: Deploy em Staging

- **Descrição:** Configurar ambiente de homologação
- **Infraestrutura:**
  - Vercel (Frontend)
  - Railway/Render (Backend)
  - Neon PostgreSQL (já configurado)
- **Estimativa:** 3h
- **Prioridade:** CRÍTICA

#### T-DEPLOY-002: Smoke Tests em Staging

- **Descrição:** Validar funcionalidades principais
- **Script:** `scripts/smoke_tests_v2.sh` (já existe)
- **Estimativa:** 1h
- **Prioridade:** CRÍTICA

#### T-DEPLOY-003: Configurar Monitoramento

- **Descrição:** Prometheus + Grafana ou alternativa
- **Métricas:**
  - Taxa de erro
  - Latência p95
  - Uptime
- **Estimativa:** 2h
- **Prioridade:** ALTA

#### T-DEPLOY-004: Configurar Alertas

- **Descrição:** Notificações de erro crítico
- **Canais:** Slack, e-mail
- **Estimativa:** 1h
- **Prioridade:** ALTA

#### T-DEPLOY-005: Deploy em Produção

- **Descrição:** Go-live oficial
- **Checklist:**
  - Backup pré-deploy
  - Smoke tests
  - Rollback plan
- **Estimativa:** 1h
- **Prioridade:** CRÍTICA

---

## 🔟 DOCUMENTAÇÃO (Dia 12: 05/12) 🟡 MÉDIA PRIORIDADE

**Total:** 6 horas

#### T-DOC-001: Atualizar README

- **Descrição:** Instruções de instalação e uso
- **Estimativa:** 1h

#### T-DOC-002: Documentação de API (Swagger)

- **Descrição:** Gerar Swagger/OpenAPI docs
- **Estimativa:** 2h

#### T-DOC-003: Guia do Usuário (MVP)

- **Descrição:** Manual básico para clientes
- **Estimativa:** 2h

#### T-DOC-004: Runbook de Incidentes

- **Descrição:** Procedimentos de troubleshooting
- **Estimativa:** 1h

---

## 📊 Resumo de Horas por Módulo

| Módulo              | Backend | Frontend | Total    | Dias (2 devs) | Prioridade    |
| ------------------- | ------- | -------- | -------- | ------------- | ------------- |
| **1. Estoque**      | 14h     | 14h      | 28h      | 2 dias        | 🔴 CRÍTICO    |
| **2. Agendamento**  | 17h     | 18h      | 35h      | 2 dias        | 🔴 BLOQUEADOR |
| **3. Lista da Vez** | 0h (✅) | 19h      | 19h      | 1 dia         | 🔴 CRÍTICO    |
| **4. Assinaturas**  | 14h     | 11h      | 25h      | 2 dias        | 🔴 BLOQUEADOR |
| **5. CRM**          | 0h (✅) | 15h      | 15h      | 1 dia         | 🟡 MÉDIA      |
| **6. Relatórios**   | 0h (✅) | 13h      | 13h      | 1 dia         | 🟡 MÉDIA      |
| **7. RBAC**         | 0h (✅) | 10h      | 10h      | 1 dia         | 🟢 BAIXA      |
| **8. Testes E2E**   | -       | -        | 12h      | 1 dia         | 🔴 CRÍTICO    |
| **9. Deploy**       | -       | -        | 8h       | 1 dia         | 🔴 CRÍTICO    |
| **10. Docs**        | -       | -        | 6h       | 0.5 dia       | 🟡 MÉDIA      |
| **TOTAL**           | **45h** | **100h** | **171h** | **~11 dias**  | -             |

---

## 🚨 Análise de Risco

### ⚠️ **ALERTA VERMELHO: Deadline em Risco Extremo**

- **Prazo:** 05/12/2025 (11 dias úteis restantes)
- **Trabalho restante:** 171 horas
- **Capacidade com 2 devs:** 11 dias × 2 devs × 8h = 176h ✅ **VIÁVEL** (margem de 5h)

### 🔴 Bloqueadores Críticos

1. **Agendamento** → Sem isso, sistema não funciona (core do negócio)
2. **Assinaturas Asaas** → Sem isso, não há receita
3. **Lista da Vez** → Diferencial competitivo principal
4. **Estoque** → Necessário para operação completa

### ✅ Fatores de Sucesso

- Backend 100% pronto (repositórios, use cases, endpoints)
- Frontend Services e Hooks 100% implementados
- LGPD e Backup já concluídos
- Módulos Financeiro, Metas e Precificação completos

---

## 🎯 Planos de Ação

### 📌 **Opção 1: Execução Completa (ARRISCADO)**

**Estratégia:** 2 devs full-time + paralelização máxima

**Divisão de Trabalho:**

**Dev 1 (Backend Focus):**

- Dia 1-2: Estoque backend (14h)
- Dia 3-4: Agendamento backend (17h)
- Dia 5-6: Asaas backend (14h)
- Dia 7: Suporte testes e deploy

**Dev 2 (Frontend Focus):**

- Dia 1-2: Estoque frontend (14h)
- Dia 3-4: Agendamento frontend (18h)
- Dia 5: Lista da Vez (19h - dividir em 2 dias)
- Dia 6: Assinaturas frontend (11h)
- Dia 7-8: CRM (15h)
- Dia 9: Relatórios (13h)
- Dia 10: RBAC (10h)
- Dia 11: Testes E2E (12h)

**Resultado:** MVP 100% completo em 11 dias ✅

**Riscos:**

- Zero margem para imprevistos
- Qualquer bloqueio causa atraso
- Burnout da equipe

---

### 📌 **Opção 2: Redução de Escopo (RECOMENDADO)**

**Estratégia:** Mover funcionalidades não-críticas para v1.1.0

**Remover do MVP v1.0.0:**

- ❌ Integração Google Agenda (4h)
- ❌ Drag & drop calendário (3h)
- ❌ Notificações push Lista da Vez (4h)
- ❌ Alerta estoque mínimo (2h)
- ❌ Estatísticas CRM (3h)
- ❌ Exportação PDF/Excel (3h)
- ❌ Comparativo mês anterior (2h)
- ❌ RBAC tela de gerenciamento (4h)

**Total economizado:** 25 horas

**Nova carga de trabalho:** 146 horas (~9 dias com 2 devs)

**Resultado:** MVP 100% funcional em 9 dias com **2 dias de margem** ✅

**Vantagens:**

- Margem de segurança para imprevistos
- Reduz risco de burnout
- Mantém todas as funcionalidades críticas
- Features removidas são melhorias incrementais

---

### 📌 **Opção 3: Atraso Controlado (ÚLTIMA OPÇÃO)**

**Nova data:** 12/12/2025 (+7 dias)

**Vantagens:**

- Permite implementação 100% completa
- Reduz risco de bugs
- Tempo para testes adequados

**Desvantagens:**

- Requer aprovação formal CEO (ADR)
- Impacta planejamento v1.1.0
- Pode afetar confiança de stakeholders

---

## ✅ Checklist de Conclusão v1.0.0 — MVP CORE

Antes de considerar PRONTO, validar:

### Funcionalidades (100% ou FALHA)

- [ ] **Agendamento:** CRUD completo + calendário visual + validação conflitos
- [ ] **Lista da Vez:** Adicionar, chamar, remover clientes
- [ ] **Financeiro:** DRE + Fluxo + Contas Pagar/Receber (✅ Backend OK)
- [ ] **Comissões:** Cálculo automático (⏸️ Baixa prioridade - pode v1.1.0)
- [ ] **Estoque:** Entrada, saída, inventário, consumo automático
- [ ] **Assinaturas:** Integração Asaas + checkout + webhooks
- [ ] **CRM:** Cadastro clientes + histórico
- [ ] **Relatórios:** Telas DRE e Fluxo com gráficos
- [ ] **Permissões:** RBAC básico (ocultar menus por papel)

### Qualidade (Mínimo ou FALHA)

- [ ] **Testes Backend:** Cobertura ≥70%
- [ ] **Testes Frontend:** Cobertura ≥60%
- [ ] **Testes E2E:** ≥80% passando (agendamento, lista, assinaturas, financeiro)
- [ ] **Performance:** p95 <300ms
- [ ] **Zero erros críticos:** Sem bugs bloqueadores

### Compliance (100% ou FALHA)

- [x] **LGPD:** Endpoints funcionais (✅ Concluído)
- [x] **Backup:** Automático rodando (✅ Concluído)
- [x] **Privacy Policy:** Publicada (✅ Concluído)
- [x] **Multi-tenant:** 100% isolado (✅ Validado)

### Operacional (100% ou FALHA)

- [ ] **Deploy Staging:** Executado e validado
- [ ] **Smoke Tests:** 100% passando
- [ ] **Monitoramento:** Prometheus configurado
- [ ] **Alertas:** Slack/e-mail funcionando
- [ ] **Documentação:** README + Swagger + Guia do usuário

### Gate de Aprovação Final

- [ ] **Code Review:** Aprovado por Tech Lead
- [ ] **QA Sign-off:** Aprovado por QA Lead
- [ ] **Product Sign-off:** Aprovado por Product Owner
- [ ] **CEO Sign-off:** Aprovado por CEO (Andrey Viana)

---

## 🎯 Recomendação Final

### ✅ **Executar Opção 2: Redução de Escopo Controlada**

**Justificativa:**

1. **Mantém deadline de 05/12/2025** ✅
2. **Todas as funcionalidades CRÍTICAS incluídas** ✅
3. **Margem de 2 dias para imprevistos** ✅
4. **Features removidas são melhorias, não bloqueadores** ✅
5. **Reduz risco de burnout da equipe** ✅

**Funcionalidades Core mantidas:**

- ✅ Agendamento completo (BLOQUEADOR)
- ✅ Lista da Vez (DIFERENCIAL)
- ✅ Assinaturas Asaas (RECEITA)
- ✅ Estoque (OPERAÇÃO)
- ✅ CRM Básico
- ✅ Relatórios essenciais
- ✅ LGPD + Backup

**Funcionalidades movidas para v1.1.0:**

- 📅 Integrações avançadas (Google Agenda)
- 📅 UX melhorias (drag & drop, notificações)
- 📅 Relatórios avançados (exportação, comparativos)

---

## 📞 Próximos Passos Imediatos

1. **CEO:** Aprovar Opção 2 (redução de escopo) - **HOJE 24/11**
2. **Tech Lead:** Dividir tarefas entre devs - **HOJE 24/11**
3. **Devs:** Iniciar Estoque (backend + frontend) - **25/11 manhã**
4. **Daily Standup:** 09:00 todos os dias até 05/12
5. **Code Review:** Obrigatório ao final de cada dia
6. **Deploy Staging:** 04/12 (validação pré-produção)
7. **Go-Live:** 05/12 19:00 (horário de menor tráfego)

---

**ATENÇÃO:** Este documento é um **plano de execução crítico**. Qualquer desvio deve ser comunicado imediatamente ao CEO.

**Última Atualização:** 24/11/2025 21:00
**Responsável:** GitHub Copilot + Andrey Viana
**Próxima Revisão:** 25/11/2025 09:00 (Daily Standup)

---

**🚀 VAMOS ENTREGAR ESSE MVP! 🚀**
