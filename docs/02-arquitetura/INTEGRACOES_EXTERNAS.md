> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🔗 Integrações Externas

**Versão:** 1.0  
**Data:** 22/11/2025  
**Status:** Planejamento (nenhuma integração implementada no código atual)

---

## 📋 Índice

1. [Estado Atual](#estado-atual)
2. [Asaas (assinaturas/pagamentos)](#asaas-assinaturaspagamentos)
3. [Google Calendar (agendamento)](#google-calendar-agendamento)
4. [Outras Integrações Futuras](#outras-integrações-futuras)
5. [Checklist de Implementação](#checklist-de-implementação)

---

## Estado Atual
- Não há clientes HTTP ou SDKs no backend para provedores externos.
- Nenhum endpoint ou cron de integração foi implementado.
- Variáveis de ambiente para integrações não são usadas no código.

---

## Asaas (assinaturas/pagamentos)
- **Motivação:** Cobrança recorrente, emissão de faturas e bloqueio de benefícios por inadimplência.
- **Planejamento:** Cliente REST resiliente (retry/backoff), endpoints de assinatura/fatura, webhooks para eventos de pagamento, sync diário.
- **Situação:** Não iniciado. Documentar quando o módulo de assinaturas começar.

## Google Calendar (agendamento)
- **Motivação:** Sincronizar agenda de compromissos com calendários externos.
- **Planejamento:** OAuth client, criação/atualização/cancelamento de eventos, idempotência, webhook/push notifications para atualizações.
- **Situação:** Não iniciado; depende do módulo de agendamento.

## Outras Integrações Futuras
- **Email/SMS/Push:** notificações de agendamento/financeiro.
- **Open Banking/Conciliação:** importação de extratos e conciliação automática (roadmap v2).
- **BI/Analytics:** exportação de dados para ferramentas externas.

---

## Checklist de Implementação
- [ ] Definir contrato e DTOs por integração.
- [ ] Configurar clientes com timeouts, retry e circuit breaker.
- [ ] Armazenar credenciais de forma segura (env/secrets).
- [ ] Criar webhooks e validar assinatura dos eventos.
- [ ] Testes de contrato/sandbox antes de produção.
- [ ] Observabilidade: métricas, logs e alertas por integração.

