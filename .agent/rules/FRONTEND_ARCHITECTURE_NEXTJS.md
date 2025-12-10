Regras oficiais do frontend:

- Next.js App Router é obrigatório.
- Estrutura:
  app/ ← páginas e layouts
  components/ ← UI reutilizável
  hooks/ ← lógica cliente (TanStack Query)
  lib/services/ ← chamadas de API com Zod
- Formulários: React Hook Form + Zod.
- Estado de dados: TanStack Query.
- Nunca usar useState para dados de API.
- Nunca criar service fora de lib/services.
- Responsividade sempre implementada.
