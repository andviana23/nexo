# Plano de Ajuste – Integração Asaas x NEXO

Objetivo: alinhar modelagem, webhooks e financeiro do NEXO ao comportamento oficial do Asaas para assegurar competência (DRE) e caixa corretos.

## Escopo imediato
- Assinaturas e cobranças recorrentes do Asaas (cartão/PIX/boleto).
- Reconhecimento de receita (CONFIRMED) e recebimento de caixa (RECEIVED/RECEIVED_IN_CASH).
- Auditoria e reconciliação automática com o Asaas.

## Ações técnicas
1) Modelagem de dados  
   - `subscriptions`: adicionar `next_due_date`, `cycle`, `asaas_status`, `last_confirmed_at`.  
   - `subscription_payments`: adicionar `status_asaas`, `due_date`, `confirmed_date`, `client_payment_date`, `credit_date`, `estimated_credit_date`, `billing_type`, `net_value`; índice único em `asaas_payment_id`.  
   - Vincular `contas_receber` a `asaas_payment_id` e `subscription_id`.

2) Webhooks Asaas  
   - Tratar eventos de forma distinta:  
     - CREATED → inserir/atualizar payment PENDING + `due_date`.  
     - OVERDUE → marcar payment OVERDUE e assinatura INADIMPLENTE.  
     - CONFIRMED → marcar payment CONFIRMED (`confirmed_date`), gerar/atualizar conta a receber (competência do mês) e assinatura ATIVA com `next_due_date` real.  
     - RECEIVED / RECEIVED_IN_CASH → marcar payment RECEIVED com `credit_date`, quitar conta a receber e lançar no fluxo diário.  
     - REFUNDED → marcar REFUNDED, estornar conta/DRE/fluxo.  
   - Idempotência: upsert por `asaas_payment_id`.

3) Financeiro (DRE)  
   - Implementar `ContaReceberRepository.SumByPeriod` filtrando por status PAGO/CONFIRMADO e subtipo PLANO.  
   - Gerar DRE usando mês de `confirmed_date` (competência).  
   - Reprocessar meses afetados após migração.

4) Fluxo de Caixa  
   - Usar `credit_date/estimated_credit_date` das cobranças RECEIVED para entradas confirmadas e previstas.  
   - Atualizar fluxo diário ao receber webhooks ou conciliação.

5) Conciliação e jobs  
   - Job diário: sincronizar `subscriptions` (status/`next_due_date`) via `/subscriptions/{id}`.  
   - Job diário: sincronizar payments do mês corrente via `/payments?subscription=...` e reconciliar com `subscription_payments`, `contas_receber` e fluxo.  
   - Ajustar cron de inadimplência para usar `next_due_date` do Asaas (não `data_vencimento` interna).

6) Migração e scripts  
   - Migration para novos campos e índice único.  
   - Script de reprocessamento: puxar pagamentos Asaas do mês atual/anterior, recriar contas a receber, recalcular DRE e fluxos diários.

7) Testes e observabilidade  
   - Testes de integração com payloads reais de webhooks (Created→Confirmed→Received, Overdue, Refunded).  
   - Teste de idempotência (reenvio do mesmo webhook).  
   - Métricas/logs: contagem por evento, divergências de conciliação, payments recebidos sem fluxo lançado.

## Sequência recomendada (sprints curtas)
1. Migração de schema + enum/status internos.  
2. Refatorar webhooks com idempotência e datas corretas.  
3. Implementar `SumByPeriod` e ajustar DRE/fluxo diário.  
4. Rodar script de reprocessamento e validar relatórios.  
5. Entregar jobs de conciliação e monitoramento.

## Critérios de aceite
- Toda cobrança CONFIRMED do mês aparece no DRE do mesmo mês.  
- Toda cobrança RECEIVED gera entrada no fluxo de caixa na `credit_date`.  
- Nenhuma assinatura PENDING/OVERDUE consta como ATIVA; nenhuma assinatura com último payment CONFIRMED/RECEIVED consta como INATIVA.  
- Reenvio de webhook não cria lançamentos duplicados.
