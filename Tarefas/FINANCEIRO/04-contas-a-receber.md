# 1.5 Controle de Contas a Receber

- **Categoria:** FINANCEIRO
- **Objetivo:** acompanhar valores a receber de planos/mensalidades e serviços com pagamento futuro, com alertas de inadimplência e integração ao dashboard e assinaturas/Asaas.

## Plano de Execução (segunda prioridade)
- **Banco de Dados:** tabela `contas_a_receber` (origem: assinatura/serviço/manual, status, valor_pago, datas); índices por tenant/status/vencimento.
- **Backend:** sync Asaas/assinaturas para gerar/atualizar; CRUD manual; job de inadimplência (N dias); conciliação `PATCH /financial/receivables/{id}/settle`; webhooks/cron de reconciliação.
- **Frontend:** listagem com filtros (status/período/origem), alertas de atraso, ação de recebimento manual, detalhes da fatura.
- **Cálculos aplicados:** insumos para fluxo de caixa compensado (previsões) usando `previsao-fluxo-caixa-compensado.md`; refletem receitas no DRE e dashboard (margem/ponto de equilíbrio/metas).

## Regras de Negócio

- Recebíveis provenientes de assinaturas sincronizam automaticamente via Asaas; serviços futuros podem ser lançados manualmente.
- Estados: `PENDENTE`, `RECEBIDO`, `ATRASADO`, `CANCELADO`.
- Inadimplência sinalizada após N dias configuráveis (default 7) sem recebimento.
- Liquidação de um recebível ocorre apenas quando fatura/status corresponde a `RECEBIDO` ou pagamento manual confirmado.

## Dependências Técnicas

- Tabelas `contas_a_receber`, `assinatura_invoices`, `planos_assinatura`, `servicos_agendados`.
- Jobs de sincronização Asaas e alertas.
- Dashboard financeiro e DRE.

## Riscos

- Desalinhamento entre assinaturas e recebíveis se webhooks/cron falharem (exigir reconciliação manual).
- Alertas de inadimplência excessivos se parametrização por tenant não existir (criar configuração personalizada).
- Necessidade de conciliação manual em caso de pagamentos parciais (prever campo `valor_pago`).

## Tarefas

1. Criar migrations `contas_a_receber` vinculadas a origem (assinatura, serviço, outro) e status.
2. Implementar sincronização com assinaturas/Asaas para gerar/atualizar recebíveis automaticamente.
3. Disponibilizar CRUD/API para lançamentos manuais de serviços com pagamento futuro.
4. Construir job de alertas de inadimplência que dispara notificações e atualiza dashboard.
5. Expor endpoints para conciliação manual (`PATCH /financial/receivables/{id}/settle`).
6. Atualizar dashboard/DRE para usar esses dados em saldos e projeções.
7. Escrever testes cobrindo sincronização, inadimplência e RBAC.

## Critérios de Aceite

- Todos os planos/mensalidades geram recebíveis automaticamente; serviços futuros podem ser registrados manualmente.
- Alertas de inadimplência ocorrem após o período configurado e são exibidos no dashboard.
- Liquidação (manual ou automática) remove o recebível das pendências e atualiza relatórios.
- APIs respondem respeitando filtros de período e status, com isolamento por tenant.
- Documentação atualizada em `docs/07-produto-e-funcionalidades/FINANCEIRO.md`.
