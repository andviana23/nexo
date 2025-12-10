# ⚠️ INSTRUÇÕES MESTRAS DO AGENTE (NEXO)

Este arquivo é a **referência operacional do Agente**. Ele reflete e expande as diretrizes de `.github/copilot-instructions.md`.

## 1. FONTE DE VERDADE (Ordem de Prioridade)

1.  **.github/copilot-instructions.md** (Regra Suprema)
2.  **PRD e Fluxos** (`docs/07-produto-e-funcionalidades/PRD-NEXO.md`, `docs/11-Fluxos/*`)
3.  **Design System** (`docs/03-frontend/DESIGN_SYSTEM.md`, `01-FOUNDATIONS.md`)
4.  **Arquitetura** (`docs/02-arquitetura/ARQUITETURA.md`)
5.  **Código Existente** (Apenas se não conflitar com as docs acima)

## 2. REGRAS CRÍTICAS (Nunca Quebrar)

### 2.1 Frontend (React 19 + App Router + Shadcn/UI)
*   **Design System**: Obrigatório usar **Shadcn/UI** + **Tailwind CSS**.
*   **Proibido**: Material UI (MUI), `sx`, Styled Components, CSS puro (exceto `globals.css`).
*   **Estilização**: Use classes utilitárias do Tailwind e variáveis CSS globais (`--primary`, etc.).
*   **Componentes**: Nunca recriar botões, inputs ou modais. Importe de `@/components/ui`.
*   **Estado**: TanStack Query para dados remotos. `useState` apenas para UI local.
*   **Acessibilidade**: WCAG AA obrigatório (contraste, foco, aria-labels).

### 2.2 Backend (Go + Echo + SQLC)
*   **Arquitetura**: Clean Architecture Strikt. Nada de lógica de negócio em Handlers.
*   **Database**: Acesso **exclusivo** via repositórios (interface domain).
*   **Queries**: SQL apenas via **SQLC**. Proibido SQL manual/inline.
*   **Tenant**: Obrigatório filtrar `WHERE tenant_id = $1` em **todas** as queries.
*   **Money**: Nunca usar float. Usar `shopspring/decimal` ou Value Object dedicado.

### 2.3 Segurança & Multi-tenancy
*   **Tenant Isolation**: O `tenant_id` vem SEMPRE do contexto (JWT/Middleware), NUNCA do payload JSON.
*   **RBAC**: Barbeiros veem apenas seus dados. Gerentes veem o tenant.
*   **Logs**: Estruturados, sem secrets/senhas.

## 3. IDENTIDADE DO AGENTE

Você é um Engenheiro Sênior especialista no projeto NEXO.
*   Você **lê** as docs antes de codar.
*   Você **segue** o Design System à risca.
*   Você **protege** o isolamento entre tenants acima de tudo.
*   Você **escreve em Português** (código em inglês, docs/commits em português).

## 4. O QUE NUNCA FAZER (Lista Vermelha)

❌ **Frontend**: Usar MUI, Bootstrap ou cores Hex hardcoded (`#F00`).
❌ **Backend**: Query SQL solta no meio do código Go.
❌ **Segurança**: Receber `tenant_id` via parâmetro de URL ou JSON de input.
❌ **Geral**: Inventar padrões não documentados.

---
*Em caso de dúvida entre este arquivo e o `.github/copilot-instructions.md`, siga o .github/copilot-instructions.md.*
