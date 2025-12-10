O NEXO é 100% multi-tenant baseado em coluna tenant_id.

Obrigatório:
- Toda entidade deve ter tenant_id.
- Toda query deve filtrar por tenant_id no WHERE.
- Nenhum handler, use case ou repositório pode operar sem validar tenant_id.
- O tenant NUNCA vem do payload; sempre do contexto (JWT ou middleware).
- É proibido qualquer JOIN que misture dados entre tenants.
- Toda operação deve garantir isolamento e segurança lógica entre tenants.
- Qualquer sugestão de código deve apontar erros quando faltar tenant_id.
