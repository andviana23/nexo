# 1.4 Controle de Contas a Pagar

- **Categoria:** FINANCEIRO
- **Objetivo:** permitir registro, automação e acompanhamento de contas a pagar (fixas e variáveis) com categorias, repetição automática, notificações (D-5/D-1/D0), anexos de comprovante, status e apoio ao dashboard.

## Plano de Execução (prioridade imediata)
- **Ordem:** Primeiro no financeiro (destrava fluxo compensado, DRE, dashboard).
- **Banco de Dados:** criar `contas_a_pagar` e `contas_a_pagar_parcelas` (tenant, categoria, recorrência, status, pix_code, anexos). Índices por tenant/status/vencimento. Eventos para auditoria.
- **Backend:** CRUD `/financial/payables`, recorrência geradora de parcelas, job de notificações (D-5/D-1/D0), upload seguro de comprovante, PATCH status→PAGO com validação de comprovante/justificativa.
- **Frontend:** formulários (cadastro/edição), listagens com filtros (status, vencimento, categoria), upload de anexos, alertas de vencimento/atraso.
- **Cálculos aplicados:** alimenta DRE e fluxo de caixa compensado (entradas/saídas previstas) — usar projeção em `previsao-fluxo-caixa-compensado.md` e metas mínimas/ponto de equilíbrio no dashboard (`faturamento-minimo-mensal.md`, `ponto-de-equilibrio.md`).

## Regras de Negócio

- Categorias seguem estrutura do módulo financeiro; contas fixas devem gerar parcelas futuras automaticamente.
- Notificações obrigatórias (5 dias antes, 1 dia antes, dia do vencimento) via canal configurado do tenant.
- Comprovantes anexados armazenados em bucket seguro e vinculados ao lançamento.
- Status válidos: `ABERTO`, `PAGO`, `ATRASADO`. Mudança para `PAGO` exige comprovante ou justificativa.
- Pagamentos via PIX disponibilizam código “copiar e colar” gerado automaticamente e válido até a data limite.

## Dependências Técnicas

- Nova tabela `contas_a_pagar` + `contas_a_pagar_parcelas` com RLS por tenant.
- Serviço de notificações (email/push/telegram) e cron scheduler.
- Storage de arquivos (S3/Spaces) + serviço de antivírus.
- Integração com dashboard (saldo contas a pagar) e DRE.

## Riscos

- Volume alto de notificações sem throttling (precisa rate limit por tenant).
- Armazenamento de anexos requer compliance LGPD (criptografia/expiração).
- Repetições mal configuradas podem gerar duplicidade (validar antes de criar parcelas).

## Tarefas

1. Modelar migrations para `contas_a_pagar` (campos: categoria, tipo, valor, recorrência, status, pix_code, anexos).
2. Implementar CRUD via API (`POST/GET/PUT/DELETE /financial/payables`) com validações de recorrência e upload seguro de comprovante.
3. Desenvolver serviço de recorrência que gera contas futuras (mensal/quinzenal/anual) e evita duplicidades.
4. Criar job de notificações (D-5/D-1/D0) com logs em `cron_executions` e integração com canal configurado.
5. Expor operações de mudança de status (`PATCH /financial/payables/{id}/status`) exigindo anexos ao marcar como pago.
6. Atualizar dashboard para exibir saldo de contas a pagar e alertas de atraso.
7. Adicionar testes (unitários/integrados) cobrindo recorrência, notificações e RBAC.

## Critérios de Aceite

- Contas podem ser cadastradas com categorias e recorrência; parcelas futuras geradas automaticamente.
- Notificações são disparadas nos três marcos configurados e registradas em log.
- Upload/download de comprovantes funciona e respeita permissões.
- PIX copiar-e-colar gerado e exibido para contas compatíveis.
- Dashboard reflete saldo em tempo real; status muda corretamente e dispara auditoria.
