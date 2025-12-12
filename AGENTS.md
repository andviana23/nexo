# Repository Guidelines

## Project Structure & Module Organization
- `backend/`: Go API (Echo + sqlc) following Clean Architecture/DDD. Key areas: `internal/domain` (entities/value objects/ports), `internal/application/usecase`, `internal/infra` (repos, db, http handlers), `cmd/api` (entrypoint), `migrations/`.
- `frontend/`: Next.js App Router UI with shadcn/ui + Tailwind tokens. Hooks/services/types live under `src/`.
- `docs/`: PRD, architecture, security (RBAC), and module flows. Treat these as source of truth.
- `scripts/`: smoke/E2E helpers. Root `Makefile` orchestrates local dev.

## Build, Test, and Development Commands
- `make dev`: start backend (Air hot‑reload) + frontend in parallel.
- `make backend` / `make frontend`: start each side separately.
- `make stop`, `make logs`: stop services or tail logs.
- `make install`: download Go modules + install frontend deps via `pnpm`.
- `make build`: production builds for both sides.
- `make validate-schema`: validate DB schema (requires `DATABASE_URL`).
- `make smoke-tests`: run API smoke tests against `API_URL`.
- Backend tests: `cd backend && go test ./...`.
- Frontend dev/lint/E2E: `cd frontend && pnpm dev | pnpm lint | pnpm test:e2e`.

## Coding Style & Naming Conventions
- **Backend (Go)**: keep handlers thin; business rules in use cases/domain. Never write manual SQL—use sqlc + migrations only. Always filter by `tenant_id` from auth context; never accept it in payloads. Run `gofmt` (tabs, stdlib style).
- **Frontend (TS/React)**: use Design System components/tokens only; no hardcoded colors/spacing, no `any`, no inline CSS. Types and enums must match backend DTOs. Prefer small, typed hooks/services.
- JSON fields are `snake_case`; money values are strings/decimals.

## Testing Guidelines
- Backend uses Go `testing` with integration tests in `backend/internal/infra/http/handler/*_integration_test.go`.
- Frontend uses Playwright E2E tests in `frontend/tests`.

## Commit & Pull Request Guidelines
- Prefer Conventional Commits: `feat(scope):`, `fix(scope):`, `docs:`, `chore:` (history is mixed, but new commits should follow this).
- PRs must describe affected modules, link the related task/PRD, note RBAC/multi‑tenant impacts, and include screenshots for UI changes.

