# 🚀 Deploy em Produção — Barber Analytics Pro v2

Guia rápido para publicar backend (Go) e frontend (Next.js standalone) em produção com scripts + GitHub Actions.

---

## 🔒 Pré-requisitos
- Acesso SSH ao servidor (`VPS_HOST`, `VPS_USER`, chave privada).
- Systemd services configurados: `barber-api` (backend) e `barber-frontend` (frontend).
- Diretórios padrão: `/opt/barber-api` e `/opt/barber-frontend`.
- Variáveis sensíveis configuradas no servidor (DATABASE_URL, JWT keys, NEXT_PUBLIC_API_URL, etc).

---

## 📦 Variáveis usadas nos scripts

| Variável                 | Default                     | Uso                                   |
| ------------------------ | --------------------------- | ------------------------------------- |
| `SSH_HOST` / `SSH_USER`  | —                           | Destino do deploy (obrigatório)       |
| `SSH_KEY_PATH`           | `~/.ssh/id_rsa`             | Caminho da chave SSH                  |
| `BACKEND_ARTIFACT`       | `backend/bin/main`          | Binário a enviar                      |
| `BACKEND_REMOTE_DIR`     | `/opt/barber-api`           | Pasta do backend no servidor          |
| `BACKEND_SERVICE`        | `barber-api`                | Nome do service systemd               |
| `FRONT_BUILD_DIR`        | `frontend/.next/standalone` | Build standalone (Next.js)            |
| `FRONT_STATIC_DIR`       | `frontend/.next/static`     | Assets estáticos do Next              |
| `FRONT_PUBLIC_DIR`       | `frontend/public`           | Assets públicos                       |
| `FRONT_REMOTE_DIR`       | `/opt/barber-frontend`      | Pasta do frontend no servidor         |
| `FRONT_SERVICE`          | `barber-frontend`           | Nome do service systemd               |
| `SKIP_BUILD`             | `0`                         | Se `1`, scripts não executam build    |

---

## ✅ Checklist Pré-Deploy
1. Tests verdes (backend `go test ./...`, frontend `pnpm test:unit` + `pnpm test:e2e` se aplicável).
2. Migrations aplicadas no banco (Neon) e compatíveis com o binário.
3. Secrets no GitHub: `VPS_HOST`, `VPS_USER`, `SSH_PRIVATE_KEY`, `NEXT_PUBLIC_API_URL_PROD`.
4. Chaves JWT presentes no servidor (`/opt/barber-api/keys` ou vars `JWT_*_PATH`).

---

## ▶️ Deploy via CLI (scripts)

```bash
# Backend
export SSH_HOST=api.seudominio.com SSH_USER=barber SSH_KEY_PATH=~/.ssh/id_rsa
go build -o backend/bin/main ./backend/cmd/api
./scripts/deploy-backend.sh

# Frontend (usa build standalone gerado pelo Next)
export SSH_HOST=app.seudominio.com SSH_USER=barber
cd frontend && pnpm install --frozen-lockfile && pnpm build && cd ..
./scripts/deploy-frontend.sh
```

Scripts fazem backup automático do binário atual (`/opt/barber-api/backups/main.<timestamp>`) e reiniciam os serviços.

---

## 🚦 Deploy via GitHub Actions (com aprovação)

Workflow: **Deploy Production (Backend + Frontend)** (`.github/workflows/deploy-production.yml`)

Inputs:
- `target`: `both` (padrão) | `backend` | `frontend`
- `ref`: branch/tag/SHA a deployar (padrão `main`)

Características:
- Ambiente `production` exige aprovação prévia antes de executar.
- Constrói backend (Go 1.22) e frontend (Node 20 + pnpm) e chama os scripts de deploy.
- Usa secrets: `VPS_HOST`, `VPS_USER`, `SSH_PRIVATE_KEY`, `NEXT_PUBLIC_API_URL_PROD`.

---

## 🔍 Pós-Deploy / Verificações
- Backend: `curl -f https://api.seudominio.com/health` (ou via SSH `curl -f http://localhost:8080/health`).
- Frontend: acessar `/dashboard` autenticado e `/signup` público.
- Logs: `journalctl -u barber-api -n 100 --no-pager` e `journalctl -u barber-frontend -n 100 --no-pager`.
- Recursos: `systemctl status barber-api` / `barber-frontend` (uptime, últimas falhas).

---

## 🔄 Rollback Rápido
1. Restaurar binário anterior (backend):
   ```bash
   sudo ls /opt/barber-api/backups
   sudo cp /opt/barber-api/backups/main.<timestamp> /opt/barber-api/main
   sudo systemctl restart barber-api
   ```
2. Frontend: reimplantar build anterior (manter último tar em `/tmp` ou gerar build anterior e rodar `deploy-frontend.sh`).

---

## 📈 Monitoramento Inicial
- Healthcheck `/health` exposto (latência + status DB/migrations).
- Logs no journal (systemd) e Prometheus/Grafana já configurados (vide `docs/05-ops-sre/MONITORING_E_ALERTAS.md`).
- Alarmes de erro 5xx e indisponibilidade devem estar habilitados antes do corte.
