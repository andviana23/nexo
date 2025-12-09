# Copilot/Codex — Instruções Oficiais (NEXO)

## Fonte de Verdade (ordem)

1. PRD e fluxos (`docs/07-produto-e-funcionalidades/PRD-NEXO.md`, `docs/11-Fluxos/FLUXO_*`)
2. Arquitetura (`docs/02-arquitetura/*`)
3. Design System e guias frontend (`docs/03-frontend/*`)
4. Guias backend (`docs/04-backend/*`)
5. Segurança/RBAC (`docs/06-seguranca/*`)
6. Código existente (pode estar desatualizado)

Sempre priorizar segurança, isolamento multi-tenant e aderência ao Design System.

## Leitura mínima antes de alterar

- PRD + fluxo do módulo impactado.
- Arquitetura geral (`docs/02-arquitetura/ARQUITETURA.md`, `FLUXOS_CRITICOS_SISTEMA.md`, `MODELO_DE_DADOS.md`).
- Segurança/RBAC (`docs/06-seguranca/RBAC.md`).
- Se backend: `docs/04-backend/API_PUBLICA.md` e DTOs.
- Se frontend: Design System (`docs/03-frontend/DESIGN_SYSTEM.md`, `01-FOUNDATIONS.md`, `03-COMPONENTS.md`).

## Princípios de Arquitetura

- Handlers apenas orquestram; regras de negócio em use cases/domain.
- Repositórios expostos por interface; SQL somente via sqlc/migrations.
- Toda operação deve receber e filtrar por `tenant_id`; nunca inferir do payload.
- RBAC aplicado em handlers; barbeiro só vê/atua nos próprios dados.
- Erros padronizados (400/403/404/409) e logs estruturados.

## Backend

- Handlers finos (binding, auth/RBAC, chamar use case, mapear erro).
- DTOs em `backend/internal/application/dto`, snake_case no JSON, `omitempty` em opcionais, dinheiro como string, sem `tenant_id` no payload.
- Use cases validam tenant, RBAC, regras de negócio, conflitos e bloqueios.
- Repositórios via sqlc; nenhuma query sem `tenant_id`. Sem SQL direto em código.
- Status lifecycle único com mapeamento de erros 400/403/404/409; logs estruturados.

## Frontend

- Usar somente componentes/tokens do Design System (shadcn/ui + Tailwind tokens). Sem cores/espaçamentos/fontes hardcoded; sem inline styles; sem `any`.
- Arquitetura Next.js App Router; hooks/serviços alinhados ao contrato do backend (payloads, status, datas).
- Responsividade e acessibilidade obrigatórias (foco, aria, contraste, navegação por teclado).
- Conferir componente existente antes de criar novo.

## Dados e Banco

- Migrations via golang-migrate; queries geradas por sqlc.
- Todas as queries com `tenant_id`; sem JOIN entre tenants; sem SQL manual em código.

## Processo e Revisão

- Antes de codar: citar docs lidos relevantes.
- Validar contrato de API, RBAC, filtros de tenant, tokens do DS e lifecycle de status.
- Apontar e corrigir violações de arquitetura, segurança ou DS; adicionar testes quando alterar fluxos críticos.

## Nunca permitir

- SQL manual ou sem `tenant_id`; lógica de domínio em handler.
- DTO com float para dinheiro ou `tenant_id` em payload.
- Cores/espaçamentos hardcoded, `any`, CSS inline, componentes fora do DS.

## Permitido e incentivado

- Refatorações que reforcem clean architecture, multi-tenant, RBAC e DS.
- Melhorias de tipagem, validações, testes e observabilidade.
