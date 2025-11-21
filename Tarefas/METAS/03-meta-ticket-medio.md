# 3. Meta por Ticket Médio

- **Categoria:** METAS
- **Objetivo:** Definir e acompanhar a meta de ticket médio (valor médio por atendimento/transação) para incentivar vendas de maior valor.
- **Escopo:** Backend, Frontend.

## Plano de Execução (prioridade 3 em Metas)
- **Banco de Dados:** tabela `metas_ticket_medio` (tenant_id, mes_ano, meta_valor, tipo=geral|barbeiro, barbeiro_id).
- **Backend:** endpoints `POST/PUT /goals/average-ticket`, `GET /goals/average-ticket/current`; serviço que calcula ticket médio = receitas/atendimentos; tratar divisão por zero.
- **Frontend:** widget de ticket médio (meta, valor atual, percentual, atendimentos).
- **Cálculos aplicados:** Ticket Médio (`docs/10-calculos/ticket-medio.md`); se exibir metas cruzadas no dashboard, pode alimentar LTV (`ltv.md`) e comparação com CAC (`cac.md`).

## Cálculo da Meta

- Meta definida como valor mínimo esperado de ticket médio (ex: R$ 80,00).
- Ticket Médio Realizado = `Total de Receitas do Período / Quantidade de Atendimentos`.
- Fórmula de Progresso: `Percentual = (Ticket Médio Realizado / Meta Ticket Médio) * 100`

## Atualização Automática

- Recalculado a cada nova receita/atendimento finalizado.
- Cache atualizado em tempo real ou a cada 5 minutos.

## Painel Visual

- Card no Dashboard mostrando:
  - Meta de Ticket Médio (R$)
  - Ticket Médio Atual (R$)
  - Percentual (%)
  - Quantidade de Atendimentos no Período
  - Status (🟢🟡🔴)

## Alertas e Status

- 🟢 **Verde**: Ticket Médio >= Meta
- 🟡 **Amarelo**: Entre 80% e 99% da meta
- 🔴 **Vermelho**: Abaixo de 80% da meta

## Regras

- RN-META-011: Meta de ticket médio deve ser valor positivo > 0.
- RN-META-012: Cálculo considera apenas receitas vinculadas a atendimentos (não despesas, não ajustes).
- RN-META-013: Se quantidade de atendimentos = 0, exibir "Sem dados" ao invés de divisão por zero.
- RN-META-014: Owner/Manager podem definir/alterar meta.
- RN-META-015: Pode ser definida meta geral (tenant) ou meta por barbeiro.

## Dependências

- Módulo Financeiro: Receitas.
- Módulo Agendamento/Atendimento: Contagem de atendimentos finalizados.
- Tabela: `metas_ticket_medio` (tenant_id, mes_ano, meta_valor, tipo [geral|barbeiro], barbeiro_id).

## Tarefas

1. Criar tabela `metas_ticket_medio`.
2. Implementar endpoint `POST/PUT /goals/average-ticket` para definir meta.
3. Implementar endpoint `GET /goals/average-ticket/current` retornando meta + ticket médio calculado.
4. Criar serviço de cálculo que divide total de receitas por quantidade de atendimentos.
5. Integrar com eventos de atendimento para atualizar contadores.
6. Desenvolver widget UI de Ticket Médio com status visual.
7. Permitir configuração opcional por barbeiro (meta individual de ticket médio).
8. Testes unitários (cálculo com diferentes cenários) e integração.

## Critérios de Aceite

- Meta de ticket médio pode ser definida (geral ou por barbeiro).
- Cálculo correto: Total Receitas / Quantidade de Atendimentos.
- Status visual (🟢🟡🔴) reflete percentual atingido.
- Dashboard exibe meta e ticket médio atual com precisão.
- Divisão por zero tratada corretamente (exibe "Sem dados").
- Testes cobrem cenários: sem atendimentos, meta não configurada, múltiplos barbeiros.
