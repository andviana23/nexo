> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 👥 Modelo Multi-Tenant

**Versão:** 1.0  
**Data:** 22/11/2025  
**Status:** Em evolução (estado atual vs planejado)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Modelo Selecionado](#modelo-selecionado)
3. [Estado Atual](#estado-atual)
4. [Plano de Implementação](#plano-de-implementação)
5. [Segurança](#segurança)
6. [Performance](#performance)
7. [Estado Atual vs Planejado](#estado-atual-vs-planejado)

---

## 🎯 Visão Geral

Multi-tenancy é o modelo onde múltiplas barbearias compartilham a mesma infraestrutura com isolamento de dados. O projeto adota **column-based** (tenant por linha).

---

## 🏗️ Modelo Selecionado: Column-Based

Cada tabela contém a coluna `tenant_id` e as queries sempre filtram por este valor. Evita migrações complexas de schema múltiplo e reduz custo operacional.

---

## 📌 Estado Atual

- `tenant_id` presente em todas as tabelas atuais (financeiro, metas, precificação, prefs).
- Middleware de tenant no backend é **mock** (header `X-Tenant-ID`); não há JWT/RBAC.
- Não há **RLS** (Row Level Security) no PostgreSQL.
- Não existe tabela `tenants` no código atual; tenants são passados como string.

---

## 🛠️ Plano de Implementação

1. **Autenticação/RBAC:** habilitar JWT RS256; middleware extrai `tenant_id` e roles.
2. **Validator:** registrar validator global no Echo para garantir inputs.
3. **RLS:** criar policies por tabela `USING (tenant_id = current_setting('app.tenant_id')::uuid)` com `SET LOCAL`.
4. **Tabela `tenants`:** cadastrar metadados do tenant (nome/plano/status) e FK de `tenant_id`.
5. **Auditoria:** logar `tenant_id`, `user_id`, operação e horário.

---

## 🔐 Segurança

- **Agora:** apenas campo `tenant_id`; sem enforcement no banco; sem auth.
- **Meta:** JWT + RBAC, middleware de tenant, RLS e auditoria.

---

## ⚡ Performance

- Índices em `tenant_id` + colunas de data/estado já presentes nos schemas atuais.
- Cautela com `SET LOCAL` para RLS; medir overhead quando ativado.

---

## 🧭 Estado Atual vs Planejado

| Item                | Estado atual (22/11/2025)                     | Planejado                                       |
| ------------------- | --------------------------------------------- | ----------------------------------------------- |
| Auth/RBAC           | Inexistente; header mock                      | JWT RS256 + roles + middleware                  |
| RLS                 | Inexistente                                   | Policies por tabela + `SET LOCAL` tenant        |
| Tabela `tenants`    | Inexistente                                   | Criar tabela e FKs                              |
| Auditoria           | Inexistente                                   | Audit log com `tenant_id`/`user_id`             |
| Validação de input  | `c.Validate` usado, mas validator não registrado | Registrar validator global no Echo             |

> Revisar a cada checkpoint do Roadmap Militar.

