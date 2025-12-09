# 🧪 Sprint 5: Testes & QA — Módulo Assinaturas

**Sprint:** 5 de 5  
**Status:** ⬜ NÃO INICIADO  
**Progresso:** 0%  
**Estimativa:** 6-10 horas  
**Prioridade:** 🟡 MÉDIA  
**Dependência:** ✅ Sprint 3 e 4 devem estar concluídas

---

## 📚 Referência Obrigatória

> ⚠️ **ANTES DE INICIAR**, leia completamente:
> 
> - **[FLUXO_ASSINATURA.md](../../docs/11-Fluxos/FLUXO_ASSINATURA.md)** — Fonte da verdade
> - **[smoke_tests_complete.sh](../../scripts/smoke_tests_complete.sh)** — Exemplo de smoke tests
> - **[frontend/tests/e2e/](../../frontend/tests/e2e/)** — Exemplos de testes E2E

---

## 📊 Progresso das Tarefas

| ID | Tarefa | Estimativa | Status | Progresso |
|----|--------|------------|--------|-----------|
| **Backend QA** |
| QA-001 | Smoke Tests: Planos | 1h | ⬜ Não Iniciado | 0% |
| QA-002 | Smoke Tests: Assinaturas | 1h | ⬜ Não Iniciado | 0% |
| QA-003 | Teste de Carga (Cron Job) | 1h | ⬜ Não Iniciado | 0% |
| **Frontend QA (E2E)** |
| QA-004 | E2E: Criar Plano | 1h | ⬜ Não Iniciado | 0% |
| QA-005 | E2E: Nova Assinatura (Wizard) | 2h | ⬜ Não Iniciado | 0% |
| QA-006 | E2E: Renovar Assinatura | 1h | ⬜ Não Iniciado | 0% |
| QA-007 | E2E: Cancelar Assinatura | 1h | ⬜ Não Iniciado | 0% |
| **Validação Manual** |
| QA-008 | Checklist de Aceitação | 2h | ⬜ Não Iniciado | 0% |

**📈 PROGRESSO SPRINT: 0/8 (0%)**

---

## 📋 Tarefas Detalhadas

### 🔧 FASE 1: Backend QA

#### QA-001: Smoke Tests: Planos

**Objetivo:** Validar endpoints de Planos via script shell

**Arquivo:** `scripts/smoke_tests_subscription.sh`

**Cenários:**
- Criar plano (sucesso)
- Criar plano com valor inválido (erro 400)
- Listar planos
- Atualizar plano
- Desativar plano

**Estimativa:** 1h

---

#### QA-002: Smoke Tests: Assinaturas

**Objetivo:** Validar endpoints de Assinaturas

**Cenários:**
- Criar assinatura PIX (sucesso)
- Criar assinatura Cartão (mock Asaas)
- Listar assinaturas
- Renovar assinatura
- Cancelar assinatura
- **Validar flag is_subscriber do cliente após criação**
- **Validar remoção da flag is_subscriber após cancelamento**
- **Validar constraint UNIQUE asaas_customer_id por tenant**

**Estimativa:** 1h

---

#### QA-003: Teste de Carga (Cron Job)

**Objetivo:** Garantir que o cron job aguenta volume de dados

**Lógica:**
1. Script para inserir 1000 assinaturas no banco
2. Executar cron job manualmente
3. Medir tempo de execução e uso de memória

**Estimativa:** 1h

---

### 🎭 FASE 2: Frontend QA (E2E)

#### QA-004: E2E: Criar Plano

**Arquivo:** `frontend/tests/e2e/subscription-plans.spec.ts`

**Cenário:**
1. Login como Admin
2. Navegar para Assinaturas > Planos
3. Clicar "Novo Plano"
4. Preencher formulário
5. Salvar
6. Verificar se aparece na lista

**Estimativa:** 1h

---

#### QA-005: E2E: Nova Assinatura (Wizard)

**Arquivo:** `frontend/tests/e2e/subscription-flow.spec.ts`

**Cenário:**
1. Navegar para Assinaturas > Assinantes
2. Clicar "Nova Assinatura"
3. Selecionar Cliente e Plano
4. Escolher PIX
5. Confirmar
6. Verificar status "ATIVO" na lista

**Estimativa:** 2h

---

#### QA-006: E2E: Renovar Assinatura

**Cenário:**
1. Buscar assinatura vencida/ativa
2. Clicar em "Renovar"
3. Confirmar pagamento
4. Verificar nova data de vencimento

**Estimativa:** 1h

---

### ✅ FASE 3: Validação Manual

#### QA-008: Checklist de Aceitação

**Referência:** [FLUXO_ASSINATURA.md — Seção 10](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#10-checklist-de-implementação)

**Itens a validar:**
- [ ] Plano inativo não aparece no wizard
- [ ] Não permite excluir plano com assinaturas
- [ ] Assinatura PIX ativa imediatamente
- [ ] Assinatura Cartão gera link Asaas
- [ ] Webhook atualiza status corretamente
- [ ] Cron job marca inadimplentes após 3 dias
- [ ] Relatórios batem com banco de dados
- [ ] Barbeiro não tem acesso ao módulo
- [ ] **Cliente recebe flag is_subscriber quando assinatura ativa**
- [ ] **Cliente perde flag is_subscriber quando não possui assinaturas ativas**
- [ ] **Não permite duplicar asaas_customer_id no mesmo tenant**
- [ ] **FindOrCreateCustomer reutiliza cliente Asaas existente (busca por nome+telefone)**

**Estimativa:** 2h

---

## ✅ Critérios de Conclusão do Módulo

- [ ] Todos os testes automatizados passando
- [ ] Checklist de aceitação 100% preenchido
- [ ] Documentação atualizada
- [ ] Código mergeado na main

---

**FIM DO DOCUMENTO — SPRINT 5: TESTES & QA**
