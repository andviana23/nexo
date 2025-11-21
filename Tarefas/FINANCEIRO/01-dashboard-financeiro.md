# 1.1 Dashboard Financeiro em Tempo Real

- **Categoria:** FINANCEIRO
- **Objetivo:** entregar visão consolidada e em tempo quase real do desempenho financeiro por tenant, com KPIs operacionais (faturamento por período/barbeiro, ticket médio, metas, ponto de equilíbrio, contas a pagar/receber, DRE simplificado e tendências de 6 meses).
- **Escopo:** backend (agregações, cache), frontend (cards/gráficos responsivos), integrações com módulos financeiros existentes.

## Plano de Execução (após payables/receivables/fluxo/comissões/DRE)
- **Banco de Dados:** usar views/materializações/snapshots existentes (receitas, despesas, payables, receivables, compensações, comissões, DRE).
- **Backend:** endpoint `/financial/dashboard` agregando KPIs e cache Redis; invalidar em eventos críticos.
- **Frontend:** cards e gráficos (linha/coluna), status de metas/PE, alertas de atraso/inadimplência.
- **Cálculos aplicados:** Ticket Médio (`docs/10-calculos/ticket-medio.md`), Ponto de Equilíbrio (`ponto-de-equilibrio.md`), Faturamento Mínimo (`faturamento-minimo-mensal.md`), Margem de Lucro (`margem-lucro.md`), Fluxo Compensado (`previsao-fluxo-caixa-compensado.md`), Taxa de Ocupação se exibida (`taxa-ocupacao-barbearia.md`, `taxa-ocupacao-barbeiro.md`), e KPIs LTV/CAC se incluídos (`ltv.md`, `cac.md`).

## Regras de Negócio

- RN-FIN-001/002/003/006/007/008: todas as métricas derivam apenas de lançamentos válidos (receitas confirmadas/recebidas, despesas pagas).
- Indicador de ponto de equilíbrio depende de metas configuradas por tenant (baseline obrigatório).
- Comparativo com metas exibe status `ATINGIDO` ou `FALTANDO R$X` calculado em cima do período corrente.
- Gráficos de tendência usam janelas fixas (últimos 6 meses) e respeitam timezone do tenant.

## Dependências Técnicas

- Repositórios `receitas`, `despesas`, `assinatura_invoices`, `contas_a_pagar`, `contas_a_receber`, `metas_financeiras` (novo) via sqlc.
- Jobs de snapshot financeiro (`docs/05-ops-sre/FLUXO_CRONS.md`) e cache Redis para agregações.
- Design System (Next.js 16 + tokens MUI/shadcn) para UI.
- RBAC (owner/manager/accountant) para acesso.

## Riscos

- Consultas pesadas se não houver materializações/snapshots (mitigar com cache + índices descritos em `MODELO_DE_DADOS.md`).
- Divergência de metas caso tenant não configure baseline (precisa fallback e alerta).
- Dados defasados se cron/sync falhar (monitorar Prometheus + alertas).

## Tarefas

1. Definir modelo `financial_goals` por tenant (metas e ponto de equilíbrio) e migrations correspondentes.
2. Implementar serviço de agregação (Go) que consolida faturamento diário, semanal, mensal e por barbeiro, ticket médio e clientes atendidos usando sqlc queries otimizadas.
3. Expor endpoint `GET /financial/dashboard` com payload completo (metas, comparativos, saldo contas a pagar/receber, DRE simplificado, séries de 6 meses) e cache Redis com invalidação a cada 5 minutos ou em eventos críticos.
4. Criar componentes frontend (cards, gráficos line/column) respeitando App Router e tokens do Design System, incluindo estado de carregamento/refetch.
5. Integrar alertas visuais para status de ponto de equilíbrio e metas, com fallback quando baseline inexistente.
6. Implementar testes (unit + integração) garantindo consistência multi-tenant e latência <500 ms p95; adicionar métricas Prometheus para monitorar tempo de resposta.

## Critérios de Aceite

- Endpoint entrega todos os KPIs solicitados, com cache configurável e resposta <500 ms p95.
- Dashboard exibe status de metas (ATINGIDO/FALTANDO R$X) e tendência de 6 meses com dados reais.
- Saldo de contas a pagar/receber e DRE simplificado batem com os módulos correspondentes.
- Testes automatizados cobrindo agregações, RBAC e UI crítica passam no pipeline.
- Observabilidade: métricas/alertas para falhas de agregação ou dados defasados (<5 min).
