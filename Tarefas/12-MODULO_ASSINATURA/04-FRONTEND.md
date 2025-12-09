# 🎨 Sprint 4: Frontend — Módulo Assinaturas

**Sprint:** 4 de 5  
**Status:** ✅ CONCLUÍDO  
**Progresso:** 100%  
**Estimativa:** 25-30 horas  
**Prioridade:** 🟠 ALTA  
**Dependência:** ✅ Sprint 2 (Backend Core) deve estar concluída

---

## 📚 Referência Obrigatória

> ⚠️ **ANTES DE INICIAR**, leia completamente:
> 
> - **[FLUXO_ASSINATURA.md](../../docs/11-Fluxos/FLUXO_ASSINATURA.md)** — Fonte da verdade
>   - Seção 1.1: Navegação
>   - Seção 2: Página Planos (UI e regras)
>   - Seção 3: Página Assinantes (UI e regras)
>   - Seção 5: Página Relatórios
> - **[DESIGN_SYSTEM.md](../../docs/03-frontend/DESIGN_SYSTEM.md)** — Componentes shadcn/ui

---

## 📊 Progresso das Tarefas

| ID | Tarefa | Estimativa | Status | Progresso |
|----|--------|------------|--------|-----------|
| **Estrutura** |
| FE-001 | Rotas e Layout | 1h | ✅ Concluído | 100% |
| FE-002 | Types e Services | 2h | ✅ Concluído | 100% |
| **Página: Planos** |
| FE-003 | Lista de Planos (DataGrid) | 2h | ✅ Concluído | 100% |
| FE-004 | Modal: Criar/Editar Plano | 3h | ✅ Concluído | 100% |
| FE-005 | Ação: Desativar Plano | 1h | ✅ Concluído | 100% |
| **Página: Assinantes** |
| FE-006 | Lista de Assinantes (Filtros) | 3h | ✅ Concluído | 100% |
| FE-007 | Wizard: Nova Assinatura | 4h | ✅ Concluído | 100% |
| FE-008 | Modal: Detalhes da Assinatura | 2h | ✅ Concluído | 100% |
| FE-009 | Ação: Renovar (PIX/Dinheiro) | 2h | ✅ Concluído | 100% |
| FE-010 | Ação: Cancelar Assinatura | 1h | ✅ Concluído | 100% |
| **Página: Relatórios** |
| FE-011 | Cards de Métricas | 2h | ✅ Concluído | 100% |
| FE-012 | Gráfico: Receita Mensal | 2h | ✅ Concluído | 100% |
| FE-013 | Breakdown por Plano | 2h | ✅ Concluído | 100% |
| **Integração** |
| FE-014 | Tratamento de Erros | 1h | ✅ Concluído | 100% |
| FE-015 | Validação de Permissões (RBAC) | 1h | ✅ Concluído | 100% |

**📈 PROGRESSO SPRINT: 15/15 (100%)**

---

## ✅ ARQUIVOS CRIADOS

### Types
- `frontend/src/types/subscription.ts` — Tipos completos (Plan, Subscription, Metrics, etc.)

### Services
- `frontend/src/services/plan-service.ts` — CRUD de planos
- `frontend/src/services/subscription-service.ts` — CRUD de assinaturas

### Hooks
- `frontend/src/hooks/use-subscriptions.ts` — React Query hooks para planos e assinaturas

### Componentes UI
- `frontend/src/components/ui/radio-group.tsx` — Componente RadioGroup (shadcn/ui)

### Páginas
- `frontend/src/app/(dashboard)/assinatura/page.tsx` — Dashboard principal
- `frontend/src/app/(dashboard)/assinatura/planos/page.tsx` — Lista de planos
- `frontend/src/app/(dashboard)/assinatura/planos/components/plan-modal.tsx` — Modal CRUD
- `frontend/src/app/(dashboard)/assinatura/assinantes/page.tsx` — Lista de assinantes
- `frontend/src/app/(dashboard)/assinatura/assinantes/components/subscription-modal.tsx` — Modal detalhes/renovar
- `frontend/src/app/(dashboard)/assinatura/nova/page.tsx` — Wizard nova assinatura

---

## 📋 Tarefas Detalhadas

### 🏗️ FASE 1: Estrutura Base

#### FE-001: Rotas e Layout

**Objetivo:** Criar estrutura de pastas e rotas no Next.js

**Arquivos:**
- `frontend/src/app/(dashboard)/assinaturas/planos/page.tsx`
- `frontend/src/app/(dashboard)/assinaturas/assinantes/page.tsx`
- `frontend/src/app/(dashboard)/assinaturas/relatorios/page.tsx`
- `frontend/src/components/layout/Sidebar.tsx` (adicionar links)

**Estimativa:** 1h

---

#### FE-002: Types e Services

**Objetivo:** Criar definições de tipos e serviços de API

**Arquivos:**
- `frontend/src/types/subscription.ts`
- `frontend/src/services/subscription-service.ts`
- `frontend/src/services/plan-service.ts`

**Conteúdo:**
- Interfaces TypeScript alinhadas com DTOs do Backend
- Funções de fetch com tratamento de erro padrão

**Estimativa:** 2h

---

### 📋 FASE 2: Página Planos

#### FE-003: Lista de Planos (DataGrid)

**Referência:** [FLUXO_ASSINATURA.md — Seção 2](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#2-página-planos)

**Componentes:**
- `DataTable` (shadcn/ui)
- Colunas: Nome, Valor, Periodicidade, Qtd Serviços, Status, Ações
- Badge de status (Ativo/Inativo)

**Estimativa:** 2h

---

#### FE-004: Modal: Criar/Editar Plano

**Componentes:**
- `Dialog` (shadcn/ui)
- `Form` (react-hook-form + zod)
- Campos: Nome, Descrição, Valor (InputCurrency), Qtd Serviços (InputNumber), Limite Uso

**Validação:**
- Valor > 0
- Nome obrigatório

**Estimativa:** 3h

---

#### FE-005: Ação: Desativar Plano

**Lógica:**
- Botão "Desativar" no menu de ações
- Confirmação com `AlertDialog`
- Chamar API DELETE /plans/:id
- Atualizar lista localmente

**Estimativa:** 1h

---

### 👥 FASE 3: Página Assinantes

#### FE-006: Lista de Assinantes

**Referência:** [FLUXO_ASSINATURA.md — Seção 3](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#3-página-assinantes)

**Filtros:**
- Status (Ativo, Inadimplente, Cancelado)
- Plano
- Busca por nome do cliente

**Colunas:**
- Cliente, Plano, Status, Vencimento, Forma Pagamento, Ações

**Estimativa:** 3h

---

#### FE-007: Wizard: Nova Assinatura

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.1](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#61-fluxo-nova-assinatura-cartão-de-crédito)

**Passos do Wizard:**
1. **Selecionar Cliente:** Autocomplete com busca
2. **Selecionar Plano:** Cards com resumo dos planos ativos
3. **Forma de Pagamento:** Radio group (Cartão, PIX, Dinheiro)
4. **Confirmação:** Resumo e botão "Criar Assinatura"

**Feedback:**
- Se Cartão: Mostrar link de pagamento gerado
- Se PIX/Dinheiro: Mostrar sucesso e data de vencimento

**Estimativa:** 4h

---

#### FE-008: Modal: Detalhes da Assinatura

**Conteúdo:**
- Dados do cliente e plano
- Status atual e histórico de pagamentos
- Barra de progresso de uso de serviços (se houver limite)
- Botões de ação (Renovar, Cancelar)

**Estimativa:** 2h

---

#### FE-009: Ação: Renovar (PIX/Dinheiro)

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.4](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#64-fluxo-renovar-assinatura-manual-pixdinheiro)

**Modal:**
- Confirmar recebimento do valor
- Campo opcional: Código da transação / Observação
- Botão "Confirmar Renovação"

**Estimativa:** 2h

---

#### FE-010: Ação: Cancelar Assinatura

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.5](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#65-fluxo-cancelar-assinatura)

**Lógica:**
- Apenas Admin/Gerente (verificar role)
- `AlertDialog` pedindo confirmação
- Aviso: "Esta ação é irreversível"

**Estimativa:** 1h

---

### 📊 FASE 4: Página Relatórios

#### FE-011: Cards de Métricas

**Referência:** [FLUXO_ASSINATURA.md — Seção 5](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#5-página-relatórios)

**Cards:**
- Total Assinaturas Ativas
- Receita Mensal Recorrente (MRR)
- Taxa de Inadimplência
- Total Cancelamentos (Mês)

**Estimativa:** 2h

---

#### FE-012: Gráfico: Receita Mensal

**Componente:**
- Recharts (BarChart ou LineChart)
- Dados vindos da API de métricas

**Estimativa:** 2h

---

#### FE-013: Breakdown por Plano

**Componente:**
- Tabela simples ou PieChart
- Mostrar qual plano é mais popular e qual gera mais receita

**Estimativa:** 2h

---

### 🔒 FASE 5: Integração e Segurança

#### FE-014: Tratamento de Erros

**Objetivo:** Feedback visual para erros de API

**Componentes:**
- `Toast` (sonner) para sucesso/erro
- `Alert` para erros de validação no form

**Estimativa:** 1h

---

#### FE-015: Validação de Permissões (RBAC)

**Referência:** [FLUXO_ASSINATURA.md — Seção 1.2](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#12-permissões-por-página)

**Lógica:**
- Usar hook `useAuth` ou `useRBAC`
- Esconder botões de ação para roles não autorizadas (ex: Barbeiro não vê menu Assinaturas)
- Redirecionar se tentar acessar rota proibida

**Estimativa:** 1h

---

## ✅ Critérios de Conclusão da Sprint

- [ ] Todas as páginas implementadas conforme wireframes/fluxo
- [ ] Integração com API funcionando (CRUD completo)
- [ ] Wizard de assinatura funcional
- [ ] Relatórios exibindo dados reais
- [ ] Permissões RBAC aplicadas corretamente
- [ ] Sem erros de console ou types

---

## 🔗 Próxima Sprint

Após conclusão, iniciar **Sprint 5: Testes & QA**
📂 [05-TESTES-QA.md](./05-TESTES-QA.md)

---

**FIM DO DOCUMENTO — SPRINT 4: FRONTEND**
