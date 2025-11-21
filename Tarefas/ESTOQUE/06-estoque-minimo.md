# 6. Estoque Mínimo e Alertas

- **Categoria:** ESTOQUE
- **Objetivo:** Monitorar níveis de estoque e alertar quando produtos atingirem o ponto de reposição (estoque mínimo).
- **Escopo:** Backend (Job/Cron), Frontend (Notificações).

## Plano de Execução (prioridade 5)
- **Banco de Dados:** campo `quantidade_minima` em produtos; log de alertas opcional; índices por tenant.
- **Backend:** job `CheckLowStock` (diário ou trigger após SAIDA/CONSUMO) com regra de anti-spam; endpoint para sugestão de compra (PDF/Excel).
- **Frontend:** widget "Reposição Necessária" e notificações; lista de compras sugerida.
- **Cálculos aplicados:** não há fórmula específica na doc de cálculos; consumo diário vem de movimentações; combina-se com curva ABC para priorizar compras.

## Fluxo Operacional

1. Cadastro de Produto define `quantidade_minima`.
2. A cada movimentação de SAIDA (ou via Cron diário), sistema verifica: `quantidade_atual <= quantidade_minima`.
3. Se verdadeiro, gera alerta/notificação para o gestor.
4. Exibe produtos "Abaixo do Mínimo" em destaque no Dashboard.

## Campos Obrigatórios

- `quantidade_minima` (no Produto).

## Comportamentos Esperados

- Notificação proativa (Email, Push ou Central de Notificações no App).
- Lista de compras sugerida baseada nos itens abaixo do mínimo.

## Regras de Baixa

- Não aplica.

## Logs e Auditoria

- Registro de envio de alerta (opcional).

## Dependências

- Cadastro de Produtos.
- Sistema de Notificações.

## Tarefas

1. Adicionar campo `quantidade_minima` na tabela/DTO de produtos (já previsto no schema).
2. Criar Job `CheckLowStockJob` (diário ou trigger-based).
3. Implementar lógica de notificação (evitar spam: notificar apenas 1x até ser reposto).
4. Criar widget "Reposição Necessária" no Dashboard de Estoque.
5. Endpoint para gerar "Sugestão de Compra" (exportar PDF/Excel).

## Critérios de Aceite

- Produtos abaixo do mínimo aparecem no relatório de reposição.
- Alertas são enviados corretamente.
- Reposição (Entrada) remove o produto da lista de alertas.
