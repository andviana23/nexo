> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🧬 Domain Models

**Versão:** 2.0  
**Data:** 22/11/2025  
**Status:** Alinhado ao estado atual (módulos futuros destacados)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Bounded Contexts Implementados](#bounded-contexts-implementados)
3. [Bounded Contexts Planejados](#bounded-contexts-planejados)
4. [Value Objects](#value-objects)
5. [Enums](#enums)
6. [Estado Atual vs Planejado](#estado-atual-vs-planejado)

---

## 🎯 Visão Geral

Modelos de domínio separados por **Bounded Context** (DDD). Este documento reflete o código existente no repositório e aponta lacunas.

---

## ✅ Bounded Contexts Implementados

### Financeiro
- **Agregados:** `ContaPagar`, `ContaReceber`, `CompensacaoBancaria`, `FluxoCaixaDiario`, `DREMensal`.
- **Serviços/Use Cases:** criação/atualização, marcação de pagamento/recebimento, geração de fluxo diário e DRE.
- **Observação:** métodos agregados de soma (`SumByPeriod`, filtros avançados) nos repositórios ainda são placeholders.

### Metas
- **Agregados:** `MetaMensal`, `MetaBarbeiro`, `MetaTicketMedio`.
- **Use cases/handlers:** CRUD completo; MetaTicket depende de ajuste de repositório para listagem por barbeiro.

### Precificação
- **Agregados:** `PrecificacaoConfig`, `PrecificacaoSimulacao`.
- **Use cases:** salvar/atualizar configuração, simular preço, salvar/listar simulações.

### Preferências do Usuário (LGPD)
- **Agregado:** `UserPreferences`.
- **Estado:** repositório implementado; handlers/UC de LGPD ainda incompletos e não expostos.

---

## 🔜 Bounded Contexts Planejados (não implementados)

- **Agendamento & Lista da Vez:** Agenda, bloqueios, conflitos, histórico, ranking.
- **Assinaturas/Asaas:** Plano, Assinatura, Fatura, eventos de webhook.
- **Comissões:** Regra de cálculo por serviço pago, lançamentos de comissão.
- **Estoque:** Produto, Movimentacao, Fornecedor, Consumo por serviço.
- **CRM/Clientes:** Cliente, histórico de visitas, contatos.

---

## 🧱 Value Objects

- `Money` (decimal, BRL implícito)
- `Percentual` (decimal)
- `MesAno` (YYYY-MM)
- `TipoCusto` (FIXO/VARIAVEL)
- `StatusConta` (PENDENTE/PAGO)
- `TipoMetaTicket` (GERAL/BARBEIRO)
- `OrigemMeta` (PLANEJADA/REAL)
- `TipoItemPrecificacao` (servico, produto)

---

## 🔠 Enums

- `StatusConta` (pendente, pago)
- `TipoCusto` (fixo, variavel)
- `TipoMetaTicket` (geral, barbeiro)
- `TipoItemPrecificacao` (servico, produto)

---

## 🧭 Estado Atual vs Planejado

| Contexto         | Estado atual (22/11/2025)                              | Planejado                                   |
| ---------------- | ------------------------------------------------------ | ------------------------------------------- |
| Financeiro       | Agregados e use cases prontos; somatórios agregados placeholders | Somatórios completos, filtros avançados     |
| Metas            | CRUD completo; MetaTicket depende de ajuste de repo    | KPIs derivados e filtros por barbeiro       |
| Precificação     | Config/Simulação funcionando                          | Integrar com custos reais (estoque/serviço) |
| User Preferences | Repo pronto; handlers LGPD incompletos                 | Rotas `/me/preferences|export|delete`       |
| Agenda/Lista     | Não implementado                                      | Agenda completa + conflitos + ranking       |
| Assinaturas      | Não implementado                                      | Integração Asaas + webhooks                 |
| Estoque/CRM      | Não implementado                                      | CRUD estoque, clientes, consumo por serviço |

> Revisar este quadro a cada checkpoint do Roadmap Militar.

