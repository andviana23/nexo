# 🚀 GitHub Copilot — Quick Reference (TL;DR)

**Barber Analytics Pro v2.0 — Versão Rápida para Chat**

---

## ⚡ Regras de Ouro (NÃO NEGOCIÁVEIS)

### 1. Banco de Dados

- ❌ Nunca SQL direto no código
- ✅ Sempre usar repositories (`internal/infrastructure/repository`)
- ✅ Sempre filtrar por `tenant_id`

### 2. Arquitetura

- ✅ Clean Architecture: Domain → Application → Infrastructure
- ❌ Nunca lógica de negócio em handlers/componentes React
- ✅ Use Cases retornam `(data, error)`

### 3. Frontend (Design System)

- ❌ Nunca cores hardcoded (#3B82F6)
- ✅ Sempre usar tokens de `@/app/theme/tokens`
- ✅ MUI 5 via `sx` prop ou `useTheme()`
- ✅ Contrast mínimo 4.5:1 (WCAG AA)

### 4. Multi-Tenancy

- ✅ Sempre extrair `tenant_id` do contexto
- ✅ Sempre filtrar por `tenant_id` em queries
- ❌ Nunca cruzar dados entre tenants

---

## 📁 Estrutura

### Backend (Go)

```
internal/
├── domain/         → Entidades, Value Objects, interfaces
├── application/    → Use Cases, DTOs, Mappers
└── infrastructure/ → HTTP, repositories, scheduler
```

### Frontend (Next.js 15.5.6)

```
app/
├── (auth)/         → Rotas públicas
├── (private)/      → Dashboards
├── components/     → UI reutilizável
├── lib/            → hooks, api, utils
└── theme/          → tokens.ts (FONTE DA VERDADE)
```

---

## 🎯 Convenções

### Backend

- Pacotes: `package financial` (lowercase)
- Entidades: `type Receita struct` (PascalCase)
- Use Cases: `CreateReceitaUseCase` (PascalCase + UseCase)
- DTOs: `CreateReceitaRequest` (PascalCase + Request/Response)

### Frontend

- Componentes: `function ReceitaForm()` (PascalCase)
- Hooks: `function useReceitas()` (camelCase + use)
- Types: `type Receita = {...}` (PascalCase)

---

## ✅ Exemplo Rápido: Backend

```go
// 1. Entity (domain/entity/receita.go)
type Receita struct {
    ID       string
    TenantID string
    Valor    valueobject.Money
}

// 2. Repository Interface (domain/repository/receita.go)
type ReceitaRepository interface {
    Save(ctx context.Context, tenantID string, r *Receita) error
}

// 3. Use Case (application/usecase/financial/create_receita.go)
func (uc *CreateReceitaUseCase) Execute(
    ctx context.Context,
    tenantID, userID string,
    input dto.CreateReceitaRequest,
) (*dto.ReceitaResponse, error) {
    // Validar, criar entidade, persistir
    return response, nil
}

// 4. Handler (infrastructure/http/handler/receita.go)
func (h *ReceitaHandler) Create(c echo.Context) error {
    tenantID := c.Get("tenant_id").(string)
    response, err := h.createUC.Execute(ctx, tenantID, userID, input)
    return c.JSON(201, response)
}
```

---

## ✅ Exemplo Rápido: Frontend

```tsx
// 1. Hook (lib/hooks/useReceitas.ts)
export function useReceitas(tenantId: string) {
  return useQuery({
    queryKey: ['receitas', tenantId],
    queryFn: () => api.receitas.list(tenantId),
  });
}

// 2. Componente (components/ui/Button.tsx)
import { tokens } from '@/app/theme/tokens';

<Box
  sx={{
    padding: tokens.spacing.md,
    color: tokens.colors.primary[500],
    borderRadius: tokens.borders.radius.md,
  }}
/>;

// 3. Form (components/financial/ReceitaForm.tsx)
const schema = z.object({
  descricao: z.string().min(1).max(255),
  valor: z.string().regex(/^\d+(\.\d{2})?$/),
});

export function ReceitaForm() {
  const { mutateAsync } = useCreateReceita();
  const { handleSubmit } = useForm({ resolver: zodResolver(schema) });

  return <form onSubmit={handleSubmit(mutateAsync)}>...</form>;
}
```

---

## 🗄️ Banco: Padrão de Tabelas

```sql
CREATE TABLE nome_tabela (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_nome_tenant ON nome_tabela(tenant_id);
```

---

## 🚫 NUNCA Fazer

❌ SQL direto fora de repositories
❌ Criar `.md` sem solicitação
❌ Ignorar `tenant_id`
❌ Lógica de negócio em handlers/componentes
❌ Cores/spacing hardcoded (`#fff`, `16px`)
❌ Expor dados sensíveis em logs
❌ Misturar camadas

---

## 📖 Docs Principais

| Doc                                                  | Quando Usar                                   |
| ---------------------------------------------------- | --------------------------------------------- |
| [Designer-System.md](../docs/Designer-System.md)     | **SEMPRE** antes de criar componentes visuais |
| [ARQUITETURA.md](../docs/ARQUITETURA.md)             | Dúvidas sobre camadas                         |
| [GUIA_DEV_BACKEND.md](../docs/GUIA_DEV_BACKEND.md)   | Padrões Go                                    |
| [GUIA_DEV_FRONTEND.md](../docs/GUIA_DEV_FRONTEND.md) | Padrões React/Next.js                         |

---

**Versão:** 5.0 | **Idioma:** Português (pt-BR)
