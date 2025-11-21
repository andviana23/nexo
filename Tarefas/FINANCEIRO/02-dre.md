# 1.2 Demonstrativo de Resultado (DRE)

- **Categoria:** FINANCEIRO
- **Objetivo:** disponibilizar DRE completo com receitas por categoria (serviços, produtos, planos), custos variáveis (comissão automática, insumos), despesas fixas/variáveis, resultado operacional, margem de lucro, lucro líquido, comparação mensal e exportação em PDF.
- **Escopo:** cálculos backend, APIs, armazenamento de configurações, UI analítica e exportação.

## Plano de Execução (quinta prioridade, depende de payables/receivables/fluxo/comissões)
- **Banco de Dados:** mapear categorias/flags fixo vs variável; opcional snapshot mensal (`dre_mensal`); garantir subtipo de receitas e integração com comissões e insumos.
- **Backend:** agregações mensais, variação m/m, exportação PDF; consumir comissões, insumos, contas a pagar/receber e compensações; APIs `/financial/dre` e `/financial/dre/export`.
- **Frontend:** tela analítica com filtros (período/categoria), gráficos/tabelas e exportação PDF.
- **Cálculos aplicados:** Margem de Lucro (`docs/10-calculos/margem-lucro.md`); Custo de Insumo (`docs/10-calculos/custo-insumo-servico.md`); Markup (`docs/10-calculos/markup.md`) para comparação; Faturamento Mínimo e Ponto de Equilíbrio para alertas (`faturamento-minimo-mensal.md`, `ponto-de-equilibrio.md`); Ticket Médio / LTV / CAC como KPIs complementares (`ticket-medio.md`, `ltv.md`, `cac.md`).

## Regras de Negócio

- Categorias devem seguir tipificação oficial (`FINANCEIRO.md`) e mapear lançamentos corretamente.
- Custos variáveis incluem comissões automáticas e insumos vinculados ao serviço.
- Despesas fixas (aluguel, internet, marketing, contabilidade, luz/água etc.) precisam estar registradas como contas a pagar/despesas com flag apropriada.
- Comparação mês a mês considera períodos consecutivos e exibe variação percentual.
- Exportação PDF deve bloquear acesso a usuários sem permissão (apenas owner/manager/accountant).

## Dependências Técnicas

- Tabelas `receitas`, `despesas`, `contas_a_pagar`, `contas_a_receber`, `comissoes_calculadas`, `insumos_por_servico`.
- Engine de comissões automáticas e módulo de contas.
- Serviço de geração de PDF (Node/Go + template) e storage seguro.
- Redis ou materialized view para somatórios mensais.

## Riscos

- Falta de dados categorizados causa DRE inconsistente (precisa validação e alertas).
- Geração de PDF pesada; necessário job assíncrono ou streaming.
- Multi-tenant sem RLS apropriado pode vazar dados (reforçar filtros).

## Tarefas

1. Mapear categorias em tabela de configuração (`dre_category_mapping`) por tenant (serviços, produtos, planos, custos, despesas fixas/variáveis).
2. Implementar agregações mensais (`dre_service`) calculando receitas, custos variáveis, despesas e resultados, com suporte a comparação mês a mês.
3. Criar endpoint `GET /financial/dre?periodo=YYYY-MM` retornando todos os blocos, variações e margens.
4. Desenvolver exportação PDF (layout oficial) acionada via `POST /financial/dre/export`, armazenando arquivo temporário e enviando link seguro.
5. Construir tela Next.js com filtros (período, categorias) e componentes gráficos/tabelas.
6. Escrever testes (unit, integração, snapshot do PDF) garantindo arredondamento decimal correto e consistência entre API e UI.

## Critérios de Aceite

- DRE apresenta todos os blocos solicitados e fecha matematicamente (receitas - custos - despesas = lucro líquido).
- Comparativo mês a mês mostra variação percentual e destaque visual.
- Exportação PDF disponível e protegida, com geração <30s e link expira em 15 min.
- APIs respondem <700 ms p95 com cache habilitado.
- Testes automatizados e documentação atualizada em `docs/07-produto-e-funcionalidades/FINANCEIRO.md`.
