Regras do Banco:

- Queries só via SQLC.
- Toda tabela deve ter tenant_id.
- Todo SELECT, UPDATE, DELETE deve filtrar por tenant_id.
- Migrations usando golang-migrate.
- Nunca usar SQL bruto em código Go.
- Nunca retornar dados sem validar tenant_id.
- Campos financeiros devem usar NUMERIC(18,2).
