# 🚀 Release v1.0.0 — MVP Core

**Versão:** 1.0.0
**Nome:** MVP Core Operacional
**Status:** 🟡 Em Desenvolvimento (54% completo)
**Data Prevista:** Janeiro 2026
**Objetivo:** Entregar sistema funcional completo para gestão operacional de barbearias

---

## 🎯 Visão Geral

O **MVP Core** é a primeira versão do NEXO, focada em resolver os problemas mais críticos de gestão de barbearias premium:

- ✅ **Agendamento** visual e intuitivo
- ✅ **Lista da Vez** automática e justa
- ✅ **Financeiro básico** (receitas, despesas, DRE, fluxo de caixa)
- ✅ **Comissões** transparentes e automáticas
- ✅ **Estoque** essencial (produtos e insumos)
- ✅ **Assinaturas** com integração Asaas (PIX/cartão)
- ✅ **CRM básico** (cadastro e histórico de clientes)
- ✅ **Relatórios mensais** simples
- ✅ **Permissões** (owner, manager, barbeiro, recepção)

---

## 📋 Funcionalidades Principais

### 1. Agendamento

**Problema que resolve:**
Agenda manual desorganizada, conflitos de horários, cancelamentos sem controle.

**Solução:**

- Calendário visual (estilo AppBarber/Trinks)
- Bloqueio automático de horários ocupados
- Integração com Google Agenda
- Status: `CREATED`, `CONFIRMED`, `IN_SERVICE`, `DONE`, `NO_SHOW`, `CANCELED`

**Critérios de Aceite:**

- [ ] Recepção consegue agendar em <30 segundos
- [ ] Não permite conflitos de horário
- [ ] Sincroniza com Google Agenda em <1 minuto
- [ ] Drag & drop funciona perfeitamente

**Implementação Técnica:**
Ver `/Tarefas/10-AGENDAMENTOS/`

---

### 2. Lista da Vez

**Problema que resolve:**
Brigas por clientes, distribuição injusta de atendimentos.

**Solução:**

- Ordenação automática justa
- Reset mensal automático
- Histórico preservado
- Pausar/retomar barbeiro

**Critérios de Aceite:**

- [ ] Ordenação correta: `current_points ASC, last_turn_at ASC, name ASC`
- [ ] Reset funciona no dia 1 de cada mês
- [ ] Histórico não é perdido
- [ ] Barbeiro pode ser pausado/retomado

**Implementação Técnica:**
Já implementado (backend completo)

---

### 3. Financeiro Básico

**Problema que resolve:**
Dono não sabe se está lucrando, gastos sem controle, DRE feito no papel.

**Solução:**

- Registro de receitas (serviços, produtos, assinaturas)
- Registro de despesas (fixas, variáveis)
- DRE mensal automático
- Fluxo de caixa diário
- Contas a pagar/receber
- Compensação bancária (D+)

**Critérios de Aceite:**

- [ ] DRE gera automaticamente todo dia 1
- [ ] Fluxo de caixa atualiza diariamente
- [ ] Compensações bancárias calculam D+ corretamente
- [ ] Categorias personalizáveis

**Implementação Técnica:**
Ver `/Tarefas/03-FINANCEIRO/`

---

### 4. Comissões

**Problema que resolve:**
Cálculo manual de comissões gera erro e desconfiança.

**Solução:**

- Comissão percentual configurável por barbeiro
- Cálculo automático sobre serviços pagos
- Relatório detalhado (barbeiro, período, valor)
- Status: `PENDING`, `PAID`, `CANCELED`

**Critérios de Aceite:**

- [ ] Comissão só conta serviços pagos
- [ ] Cálculo nunca ultrapassa valor do serviço
- [ ] Barbeiro vê comissão em tempo real

**Implementação Técnica:**
Ver `/Tarefas/03-FINANCEIRO/modulo-05-comissoes-automaticas.md`

---

### 5. Estoque Essencial

**Problema que resolve:**
Falta de produtos, desperdício, custo não rastreado.

**Solução:**

- Cadastro de produtos/insumos
- Entrada e saída manual
- Consumo interno
- Alerta de estoque mínimo
- Custo por serviço (ficha técnica)

**Critérios de Aceite:**

- [ ] Não permite estoque negativo
- [ ] Alertas funcionam quando abaixo do mínimo
- [ ] Consumo por serviço abate automaticamente

**Implementação Técnica:**
Ver `/Tarefas/04-ESTOQUE/`

---

### 6. Assinaturas (Asaas)

**Problema que resolve:**
Clientes esquecem de pagar mensalidade, cobrança manual é difícil.

**Solução:**

- Criação de planos personalizados
- Integração com Asaas (PIX/cartão)
- Cobrança automática
- Controle de limite de uso
- Suspensão automática por inadimplência

**Critérios de Aceite:**

- [ ] Criar assinatura no Asaas em <5 segundos
- [ ] Webhook atualiza status automaticamente
- [ ] Benefícios bloqueiam se inadimplente
- [ ] Sincronização funciona sem falhas

**Implementação Técnica:**
Ver `/Tarefas/v1.0.0 — MVP Core/INTEGRACAO_ASAAS.md`

---

### 7. CRM Básico

**Problema que resolve:**
Dados de clientes espalhados, sem histórico unificado.

**Solução:**

- Cadastro completo de clientes
- Histórico de agendamentos
- Histórico de compras
- Tags (VIP, Risco, Novo)
- Origem (Instagram, Google, indicação)

**Critérios de Aceite:**

- [ ] Busca por nome/telefone em <1 segundo
- [ ] Histórico completo visível
- [ ] Barbeiro não vê dados sensíveis (telefone, email)

**Implementação Técnica:**
Já parcialmente implementado

---

### 8. Relatórios Mensais

**Problema que resolve:**
Dono não tem visão clara de resultados.

**Solução:**

- DRE mensal completo
- Fluxo de caixa compensado
- Ticket médio (por barbeiro/unidade)
- Ranking de barbeiros
- Exportação CSV/Excel

**Critérios de Aceite:**

- [ ] Relatórios geram em <3 segundos
- [ ] Exportação funciona sem erros
- [ ] Dados precisos (validados contra banco)

**Implementação Técnica:**
Ver `/Tarefas/03-FINANCEIRO/`

---

### 9. Permissões (RBAC)

**Problema que resolve:**
Barbeiro vê dados que não deveria, risco de vazamento.

**Solução:**

- **Owner**: Acesso total
- **Manager**: Acesso total à unidade
- **Barbeiro**: Apenas dados próprios
- **Recepcionista**: Agenda + cadastros
- **Contador**: Read-only financeiro

**Critérios de Aceite:**

- [ ] Barbeiro não acessa dados de outros
- [ ] Recepcionista não vê financeiro
- [ ] Manager não cruza tenants

**Implementação Técnica:**
Ver `docs/06-seguranca/RBAC.md`

---

## 📊 Status de Implementação

| Módulo       | Backend | Frontend | Testes  | Status       |
| ------------ | ------- | -------- | ------- | ------------ |
| Agendamento  | ⚪ 0%   | ⚪ 0%    | ⚪ 0%   | Planejado    |
| Lista da Vez | ✅ 100% | ✅ 100%  | ✅ 100% | Concluído    |
| Financeiro   | 🟡 70%  | 🟡 60%   | 🟡 40%  | Em Curso     |
| Comissões    | 🟡 80%  | ⚪ 0%    | ⚪ 0%   | Em Curso     |
| Estoque      | ⚪ 0%   | ⚪ 0%    | ⚪ 0%   | Bloqueado    |
| Assinaturas  | 🟡 60%  | 🟡 50%   | 🟡 30%  | Em Curso     |
| CRM          | ✅ 90%  | ✅ 85%   | ✅ 70%  | Quase Pronto |
| Relatórios   | 🟡 50%  | 🟡 40%   | ⚪ 0%   | Em Curso     |
| Permissões   | ✅ 95%  | ✅ 90%   | ✅ 80%  | Quase Pronto |

**Progresso Geral:** 54%

---

## 🔗 Implementação Técnica

Este release é implementado através das seguintes etapas técnicas:

### Obrigatórias (Sequencial)

1. ✅ **CONCLUIR/** - Backlog imediato (domínio, repos, use cases)
2. 🟡 **01-BLOQUEIOS-BASE/** - Base técnica (70% completo)
3. ⚪ **02-HARDENING-OPS/** - LGPD + Backup

### Módulos (Paralelo após #2)

4. 🟡 **03-FINANCEIRO/** - Módulo Financeiro (60%)
5. ⚪ **04-ESTOQUE/** - Módulo Estoque (0%)
6. ⚪ **05-METAS/** - Módulo Metas (0%)
7. ⚪ **06-PRECIFICACAO/** - Módulo Precificação (0%)

### Finalização (Sequencial)

8. ⚪ **07-LANCAMENTO/** - Go-Live
9. ⚪ **10-AGENDAMENTOS/** - Módulo Agendamentos

**Ver `/Tarefas/XX-NOME/` para detalhes técnicos de cada etapa.**

---

## ✅ Critérios de Conclusão

O MVP v1.0.0 estará **PRONTO** quando:

### Funcionalidades

- [ ] Todos os 9 módulos implementados e testados
- [ ] Telas responsivas (desktop + mobile)
- [ ] Exportação CSV/Excel funcionando
- [ ] Integração Asaas estável

### Qualidade

- [ ] Cobertura de testes >70% (backend + frontend)
- [ ] Testes E2E >80% passando
- [ ] Erros tratados amigavelmente
- [ ] Performance: p95 <300ms

### Compliance

- [ ] LGPD completo (export + delete + consent)
- [ ] Backup automático funcionando
- [ ] Privacy Policy publicada
- [ ] Multi-tenant 100% isolado

### Operação

- [ ] Deploy em produção estável
- [ ] Monitoramento configurado
- [ ] Alertas funcionando
- [ ] Documentação completa

---

## 📅 Cronograma

| Milestone            | Data Prevista   | Status       |
| -------------------- | --------------- | ------------ |
| Base Técnica (01-02) | Nov 2025        | 🟡 Em Curso  |
| Módulos (03-06)      | Dez 2025        | ⚪ Pendente  |
| Finalização (07)     | Jan 2026        | ⚪ Pendente  |
| Agendamentos (10)    | Jan 2026        | ⚪ Pendente  |
| **Go-Live v1.0.0**   | **26 Jan 2026** | ⚪ Planejado |

---

## 📚 Referências

- [PRD Completo](../../PRD-NEXO.md)
- [Integração Asaas](./INTEGRACAO_ASAAS.md)
- [Roadmap Produto](../../docs/07-produto-e-funcionalidades/ROADMAP_PRODUTO.md)
- [Catálogo Funcionalidades](../../docs/07-produto-e-funcionalidades/CATALOGO_FUNCIONALIDADES.md)

---

**Última Atualização:** 22/11/2025
**Próxima Revisão:** Conclusão de 01-BLOQUEIOS-BASE
