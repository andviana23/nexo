Regras obrigatórias de backend:

- Linguagem: Go 1.24.x
- Framework: Echo v4
- Queries: SQLC apenas; proibido SQL manual.
- Nunca usar float para dinheiro; usar valueobject.Money.
- DTOs usam snake_case e nunca incluem tenant_id.
- Use Cases SEMPRE:
  - validam tenant,
  - validam RBAC,
  - executam regra de negócio,
  - previnem conflitos,
  - retornam erros padronizados (400, 403, 404, 409).
- Handlers:
  - fazem bind,
  - validam input,
  - chamam use case,
  - retornam DTO mapper.
- Zap é o logger obrigatório.
- Nenhum acesso direto ao banco fora dos repositórios.
