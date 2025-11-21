# 1. Meta Geral do Mês

- **Categoria:** METAS
- **Objetivo:** Definir e acompanhar a meta de faturamento total do tenant (barbearia) para o mês corrente, com visualização de progresso em tempo real.
- **Escopo:** Backend (cálculo, atualização), Frontend (painel, alertas).

## Plano de Execução (prioridade 1 em Metas)
- **Banco de Dados:** tabela `metas_mensais` (tenant_id, mes_ano, meta_faturamento, criado_por, origem/status), índices por tenant/mes.
- **Backend:** endpoints `POST/PUT /goals/monthly`, `GET /goals/monthly/current`; serviço de progresso (somar receitas CONFIRMADO/RECEBIDO do mês); cache Redis invalidado em eventos de receita.
- **Frontend:** card de meta mensal com barra de progresso, status e alerta “meta não configurada”.
- **Cálculos aplicados:** progresso usa receita acumulada; referências a Ponto de Equilíbrio e Faturamento Mínimo para alertas no dashboard (`docs/10-calculos/ponto-de-equilibrio.md`, `docs/10-calculos/faturamento-minimo-mensal.md`).

## Cálculo da Meta

- Meta definida manualmente pelo Owner/Manager no início do mês (ou herdada do mês anterior).
- Progresso calculado somando todas as receitas `CONFIRMADO` ou `RECEBIDO` do mês corrente.
- Fórmula: `Percentual Atingido = (Receitas do Mês / Meta) * 100`

## Atualização Automática

- Atualizado em tempo real a cada lançamento de receita.
- Cache Redis invalida a cada nova receita criada/atualizada.
- Dashboard consome endpoint com cache de 5 minutos.

## Painel Visual

- Card destacado no Dashboard Financeiro.
- Exibição:
  - Valor da Meta (R$ XXX.XXX,XX)
  - Valor Realizado (R$ XXX.XXX,XX)
  - Percentual (XX%)
  - Barra de progresso
  - Status visual (🟢🟡🔴)

## Alertas e Status

- 🟢 **Verde**: Atingido >= 100%
- 🟡 **Amarelo**: Entre 70% e 99%
- 🔴 **Vermelho**: Abaixo de 70%

## Regras

- RN-META-001: Meta deve ser um valor positivo > 0.
- RN-META-002: Apenas Owner/Manager podem definir/alterar meta.
- RN-META-003: Progresso considera apenas receitas do tipo `CONFIRMADO` ou `RECEBIDO`.
- RN-META-004: Meta não pode ser alterada retroativamente (apenas para mês corrente ou futuro).
- RN-META-005: Se meta não definida, exibe alerta "Meta não configurada" no dashboard.

## Dependências

- Módulo Financeiro: Tabela `receitas`, agregações por período.
- Tabela nova: `metas_mensais` (tenant_id, mes_ano, meta_faturamento, criado_em).
- Dashboard Financeiro (integração com widget de meta).
- RBAC (owner/manager para edição).

## Tarefas

1. Criar tabela `metas_mensais` com campos: `id`, `tenant_id`, `mes_ano` (YYYY-MM), `meta_faturamento`, `criado_por`, `criado_em`, `atualizado_em`.
2. Implementar endpoint `POST/PUT /goals/monthly` para criação/atualização de meta mensal.
3. Implementar endpoint `GET /goals/monthly/current` retornando meta do mês corrente + progresso calculado.
4. Criar serviço de cálculo que agrega receitas do mês e calcula percentual.
5. Integrar cache Redis com invalidação em eventos de receita.
6. Desenvolver componente UI (card de meta) com barra de progresso e status colorido.
7. Adicionar validação: impedir alteração de metas de meses passados.
8. Testes unitários (cálculo de percentual) e integração (criação/atualização de meta).

## Critérios de Aceite

- Meta pode ser definida/editada apenas por Owner/Manager.
- Progresso é calculado corretamente somando receitas do mês.
- Status visual (🟢🟡🔴) muda conforme percentual atingido.
- Dashboard exibe meta e progresso em tempo real (<5 min de defasagem).
- Tentativa de editar meta de mês passado retorna erro 400.
- Testes automatizados cobrem todos os cenários de status.
