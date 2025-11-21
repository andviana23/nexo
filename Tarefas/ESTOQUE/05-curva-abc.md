# 5. Curva ABC e Relatórios

- **Categoria:** ESTOQUE
- **Objetivo:** Classificar produtos por importância (Valor de Consumo/Venda) e fornecer visão analítica do estoque.
- **Escopo:** Backend (Queries Analíticas), Frontend (Dashboards).

## Plano de Execução (prioridade 6)
- **Banco de Dados:** usar movimentações históricas; considerar view/materialização para cálculo rápido de valores movimentados.
- **Backend:** query/serviço para cálculo ABC (ordenar por valor, acumular % e marcar A/B/C); endpoint `GET /stock/reports/abc`.
- **Frontend:** visual Pareto/barras e badges A/B/C na listagem de produtos; filtros por período e critério (custo de consumo ou valor de venda).
- **Cálculos aplicados:** “Curva ABC” (`docs/10-calculos/curva-abc.md`).

## Fluxo Operacional

1. Sistema calcula periodicamente (ou on-demand) a classificação ABC.
   - Classe A: Itens que representam ~80% do valor movimentado (alta importância).
   - Classe B: Itens que representam ~15% (média importância).
   - Classe C: Itens que representam ~5% (baixa importância).
2. Exibir relatório/gráfico no Dashboard de Estoque.
3. Permitir filtrar por período.

## Campos Obrigatórios

- `periodo_analise` (ex: últimos 3, 6, 12 meses).
- Critério: Valor de Custo (Consumo) ou Valor de Venda (Receita).

## Comportamentos Esperados

- Cálculo rápido (usar views materializadas ou cache se necessário).
- Identificação visual clara (tags A, B, C) na lista de produtos.

## Regras de Baixa

- Não aplica. Apenas leitura.

## Logs e Auditoria

- Acesso a relatórios sensíveis pode ser logado.

## Dependências

- Histórico de Movimentações (`movimentacoes`).
- Cadastro de Produtos (Valores).

## Tarefas

1. Criar query SQL otimizada para cálculo ABC (baseado em movimentações de SAIDA \* Valor).
2. Criar endpoint `GET /stock/reports/abc`.
3. Implementar visualização gráfica (Pareto) no frontend.
4. Adicionar badge A/B/C na listagem de produtos.

## Critérios de Aceite

- Classificação ABC condizente com os valores movimentados no período.
- Relatório carrega em tempo aceitável (< 2s).
