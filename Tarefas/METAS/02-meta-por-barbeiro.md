# 2. Meta por Barbeiro

- **Categoria:** METAS
- **Objetivo:** Acompanhar metas individuais de cada barbeiro por tipo de faturamento (Serviços Gerais, Serviços Extras, Venda de Produtos).
- **Escopo:** Backend, Frontend (painel individual e comparativo).

## Plano de Execução (prioridade 2 em Metas)
- **Banco de Dados:** tabela `metas_barbeiro` (tenant_id, barbeiro_id, mes_ano, metas por componente); índices por tenant/mes/barbeiro.
- **Backend:** endpoints `POST/PUT /goals/barber/{barbeiro_id}`, `GET /goals/barber/{id}/current`, `GET /goals/barbers/ranking`; serviço de agregação por componente; cache invalidado em eventos de serviço/venda.
- **Frontend:** cards individuais e ranking comparativo (percentual atingido, breakdown Serviços Gerais/Extras/Produtos).
- **Cálculos aplicados:** tickets/receitas por barbeiro influenciam LTV/CAC/PE no dashboard; para metas de valor usar receitas agregadas; para alertas gerais pode cruzar com Ponto de Equilíbrio (`ponto-de-equilibrio.md`) e Ticket Médio (`ticket-medio.md`) se exibidos por barbeiro.

## Cálculo da Meta

- Meta definida por barbeiro com três componentes:
  1. **Faturamento Geral**: Soma de todos os serviços executados pelo barbeiro.
  2. **Serviços Extras**: Serviços adicionais/premium (categoria específica).
  3. **Venda de Produtos**: Produtos vendidos diretamente pelo barbeiro.
- Progresso calculado individualmente para cada componente.
- Meta total = Soma dos três componentes.

## Atualização Automática

- Atualizado ao finalizar atendimento/serviço ou registrar venda de produto.
- Listener no evento `ServiceCompleted` e `ProductSold` atualiza cache.

## Painel Visual

- Card individual por barbeiro no Dashboard.
- Visão consolidada comparativa (ranking de barbeiros).
- Exibição por barbeiro:
  - Nome do Barbeiro
  - Meta Total (R$)
  - Realizado Total (R$)
  - Breakdown: Serviços Gerais / Extras / Produtos
  - Percentual Geral
  - Status (🟢🟡🔴)

## Alertas e Status

- 🟢 **Verde**: >= 100% da meta total
- 🟡 **Amarelo**: 70-99% da meta total
- 🔴 **Vermelho**: < 70% da meta total

## Regras

- RN-META-006: Meta de barbeiro deve ser configurada por tenant (Owner/Manager).
- RN-META-007: Apenas serviços finalizados e produtos vendidos contam para o progresso.
- RN-META-008: Barbeiro inativo não aparece no ranking, mas mantém histórico.
- RN-META-009: Metas podem ser diferentes para cada barbeiro.
- RN-META-010: Categorias de "Serviços Extras" devem ser configuráveis por tenant.

## Dependências

- Módulo Financeiro: Receitas vinculadas a barbeiro/profissional.
- Módulo Cadastro: Profissionais (barbeiros).
- Tabela nova: `metas_barbeiro` (tenant_id, barbeiro_id, mes_ano, meta_servicos_gerais, meta_servicos_extras, meta_produtos).
- Configuração de categorias (quais são "Serviços Extras").

## Tarefas

1. Criar tabela `metas_barbeiro` com campos: `id`, `tenant_id`, `barbeiro_id`, `mes_ano`, `meta_servicos_gerais`, `meta_servicos_extras`, `meta_produtos`, `criado_em`, `atualizado_em`.
2. Criar endpoint `POST/PUT /goals/barber/{barbeiro_id}` para configurar meta individual.
3. Criar endpoint `GET /goals/barber/{barbeiro_id}/current` retornando meta + progresso detalhado.
4. Criar endpoint `GET /goals/barbers/ranking` retornando lista ordenada por percentual atingido.
5. Implementar serviço de agregação que calcula progresso por componente (serviços gerais, extras, produtos).
6. Desenvolver UI: card individual de barbeiro + tela de ranking comparativo.
7. Configuração de categorias "Serviços Extras" (tabela ou config JSON por tenant).
8. Testes unitários (cálculo de componentes) e integração (atualização automática).

## Critérios de Aceite

- Meta pode ser configurada individualmente para cada barbeiro.
- Progresso é separado corretamente entre Serviços Gerais, Extras e Produtos.
- Ranking exibe todos os barbeiros ativos ordenados por performance.
- Status visual (🟢🟡🔴) correto para cada barbeiro.
- Eventos de finalização de serviço/venda atualizam progresso automaticamente.
- Barbeiros inativos não aparecem no ranking mas mantêm dados históricos.
