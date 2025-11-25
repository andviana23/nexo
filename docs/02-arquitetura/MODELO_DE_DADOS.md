> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🗄️ Design do Banco de Dados

**Versão:** 2.0  
**Data:** 22/11/2025  
**Status:** Alinhado ao estado atual (módulos futuros destacados)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Tabelas Atuais](#tabelas-atuais)
3. [Tabelas Planejadas](#tabelas-planejadas)
4. [Índices & Performance](#índices--performance)
5. [Migrations](#migrations)
6. [Estado Atual vs Planejado](#estado-atual-vs-planejado)

---

## 🎯 Visão Geral

O schema atual cobre apenas os módulos já implementados no backend (financeiro, metas, precificação e preferências de usuário). Módulos como agendamento, assinaturas/Asaas, estoque, comissões e CRM ainda não possuem tabelas.

---

## 📦 Tabelas Atuais

- **Financeiro**
  - `contas_a_pagar`
  - `contas_a_receber`
  - `compensacoes_bancarias`
  - `fluxo_caixa_diario`
  - `dre_mensal`
- **Metas**
  - `metas_mensais`
  - `metas_barbeiro`
  - `metas_ticket_medio`
- **Precificação**
  - `precificacao_config`
  - `precificacao_simulacoes`
- **LGPD/Preferências**
  - `user_preferences`

> Fonte: `backend/internal/infra/db/schema/*.sql`

---

## 🔜 Tabelas Planejadas (não existentes)

- **Agendamento & Lista da Vez:** `agendamentos`, `agendamento_blocos`, `barber_turns`, `barber_turn_history`.
- **Assinaturas/Asaas:** `planos`, `assinaturas`, `faturas_assinatura`, `webhook_events`.
- **Comissões:** `comissoes`, `comissoes_regras`.
- **Estoque:** `produtos`, `movimentacoes_estoque`, `fornecedores`, `consumos_servico`.
- **CRM/Clientes:** `clientes`, `historico_visitas`, `contatos`.

Essas tabelas devem ser especificadas e migradas conforme os módulos forem iniciados.

---

## 📊 Índices & Performance

- Índices por `tenant_id` e datas em todas as tabelas atuais (ver arquivos `.sql`).
- `UNIQUE(id, tenant_id)` adotado para evitar vazamento cross-tenant.
- Gap: ausência de **RLS** (Row Level Security) — ativar quando auth/JWT estiver pronto.

---

## 🧳 Migrations

- Migrations estão em `backend/internal/infra/db/schema/migrations`.
- Cobrem apenas os módulos atuais; novas migrations serão necessárias para agendamento, assinaturas, estoque, etc.

---

## 🧭 Estado Atual vs Planejado

| Área            | Estado atual (22/11/2025)                      | Planejado                                      |
| --------------- | ---------------------------------------------- | ---------------------------------------------- |
| Financeiro      | Tabelas e migrations criadas                   | Completar agregações/índices específicos       |
| Metas           | Tabelas criadas                                | Ajustes para filtros por barbeiro/período      |
| Precificação    | Tabelas criadas                                | Ligar a custos reais (estoque)                 |
| User Prefs      | Tabela criada                                  | Audit logs e histórico de consentimento        |
| Agendamento     | Não existe                                     | Criar schema completo + índices de conflito    |
| Assinaturas     | Não existe                                     | Criar schema para Asaas + webhooks             |
| Estoque/CRM     | Não existe                                     | Criar schema de estoque/cliente/consumo        |
| Segurança       | Sem RLS, sem auditoria                         | Ativar RLS, auditoria e policies por role      |

> Atualizar esta tabela a cada checkpoint do Roadmap Militar.

