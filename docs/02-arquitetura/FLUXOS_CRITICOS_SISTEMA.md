> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# Fluxos Críticos do Sistema

**Data:** 22/11/2025  
**Status:** Alinhado ao estado atual (futuros sinalizados)

---

## 📋 Índice

1. [Financeiro](#financeiro)
2. [Metas](#metas)
3. [Precificação](#precificação)
4. [LGPD/Preferências](#lgpdpreferências)
5. [Futuros: Agendamento e Assinaturas](#futuros-agendamento-e-assinaturas)

---

## Financeiro
1. Frontend envia contas (payables/receivables) → handlers aplicam bind/validate (validator global ainda não configurado no server) → use cases criam entidades → repositórios SQLC persistem em PostgreSQL (Neon).
2. Marcação de pagamento/recebimento → atualiza status/valores → grava data de pagamento/recebimento.
3. Cron (scheduler) executa `GenerateFluxoDiario` e `GenerateDRE` usando `contas_a_pagar/receber` → salva snapshots em `fluxo_caixa_diario` e `dre_mensal`. Somatórios ainda são placeholders (precisam ser implementados nos repositórios).

## Metas
1. CRUD de metas mensais/barbeiro/ticket → handlers → use cases → repos sqlc (`metas_mensais`, `metas_barbeiro`, `metas_ticket_medio`).
2. Listagens usam filtros básicos; MetaTicket depende de ajuste de repo para listagem por barbeiro.

## Precificação
1. Configuração de precificação salva em `precificacao_config`.
2. Simulações gravadas em `precificacao_simulacoes`; cálculos usam percentuais/impostos/comissão defaults.

## LGPD/Preferências
1. `user_preferences` armazena consentimentos; repositório implementado.
2. Handlers LGPD (export/delete/preferences) estão incompletos e não expostos; banner frontend depende dessas rotas.

## Futuros: Agendamento e Assinaturas
- **Agendamento / Lista da Vez:** Nenhum fluxo implementado. Depende de modelo de dados, regras de conflito, UI de agenda e, opcionalmente, integração Google Calendar.
- **Assinaturas (Asaas):** Nenhum fluxo implementado. Depende de cliente Asaas, webhooks, sync de faturas e bloqueio por inadimplência.

> Revisar estes fluxos a cada checkpoint; mover itens “futuros” para seções principais quando implementados.

