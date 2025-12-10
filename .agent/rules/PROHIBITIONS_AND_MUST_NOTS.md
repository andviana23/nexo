É terminantemente proibido:

- SQL manual.
- Payload contendo tenant_id.
- Cores e espaçamentos hardcoded.
- CSS inline.
- Criar estados sem conexão com TanStack Query.
- Criar componentes fora do Design System.
- Colocar lógica de negócio em handlers.
- Usar float para valores financeiros.
- Ignorar RBAC ou tenant_id.
- Criar pastas ou padrões que não seguem a arquitetura documentada.

O Agent deve apontar erros sempre que detectar qualquer violação.
