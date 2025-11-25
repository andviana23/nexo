> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# CI/CD Pipeline

- GitHub Actions (build/test Go 1.24.0, build frontend Next.js 14.2.4, lint/tests).
- Deploy backend: SSH/systemd, make build, make migrate-up, restart serviço.
- Deploy frontend: build estático Next.js, upload para host (Vercel/VPS).
- Referência: GUIA_DEVOPS.md (seção CI/CD) e workflows em .github/workflows/.
