# NEXO — Guia de Desenvolvimento Frontend

Guia prático para desenvolvedores trabalharem no frontend do NEXO.

## 1. Setup Inicial

Certifique-se de ter o Node.js 18+ e pnpm instalados.

```bash
# Instalar dependências
pnpm install

# Rodar servidor de desenvolvimento
pnpm dev
```

O projeto rodará em `http://localhost:3000`.

---

## 2. Criando uma Nova Página

No Next.js App Router, as páginas são definidas pela estrutura de pastas em `src/app`.

### Exemplo: Página de Clientes

1.  Crie a pasta `src/app/(dashboard)/clientes`.
2.  Crie o arquivo `page.tsx` dentro dela.

```tsx
// src/app/(dashboard)/clientes/page.tsx
import { Metadata } from 'next';
import { ClientsList } from './_components/clients-list';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';

export const metadata: Metadata = {
  title: 'Clientes | NEXO',
};

export default function ClientsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Clientes</h1>
        <Button>
          <Plus className="mr-2 h-4 w-4" /> Novo Cliente
        </Button>
      </div>

      <ClientsList />
    </div>
  );
}
```

> **Dica:** Use `_components` dentro da pasta da rota para componentes que são específicos daquela página e não serão reutilizados globalmente.

---

## 3. Adicionando Componentes shadcn/ui

Os componentes não vêm instalados por padrão. Você deve adicioná-los conforme a necessidade.

```bash
# Exemplo: Adicionar um Select
npx shadcn@latest add select
```

Isso criará o arquivo `src/components/ui/select.tsx`.

---

## 4. Fazendo Chamadas à API

Use os hooks customizados que encapsulam o `TanStack Query`.

1.  Crie o hook em `src/hooks/use-something.ts` (ou `src/services/something/hooks.ts`).
2.  Use o hook no seu componente.

```tsx
'use client';

import { useClients } from '@/hooks/use-clients';
import { DataTable } from '@/components/ui/data-table';
import { columns } from './columns';

export function ClientsList() {
  const { data: clients, isLoading, isError } = useClients();

  if (isLoading) return <div>Carregando...</div>;
  if (isError) return <div>Erro ao carregar clientes.</div>;

  return <DataTable columns={columns} data={clients} />;
}
```

---

## 5. Estilização com Tailwind

Use classes utilitárias para quase tudo.

- **Margem/Padding:** `m-4`, `p-6`, `gap-4`
- **Flexbox:** `flex`, `items-center`, `justify-between`
- **Grid:** `grid`, `grid-cols-1`, `md:grid-cols-3`
- **Cores:** `bg-primary`, `text-muted-foreground`, `border-input`
- **Tipografia:** `text-sm`, `font-bold`, `leading-none`

Para estilos condicionais, use `cn()`:

```tsx
import { cn } from "@/lib/utils"

<div className={cn(
  "p-4 border rounded",
  isActive && "border-primary bg-primary/10"
)}>
```

---

## 6. Checklist de Pull Request

Antes de abrir um PR, verifique:

### Build & Lint

- [ ] O código compila sem erros (`pnpm build`).
- [ ] Não há erros de lint (`pnpm lint`).
- [ ] Não há `console.log` esquecido.

### Responsividade (OBRIGATÓRIO)

> ⚠️ **PRs sem responsividade adequada serão REJEITADOS.**

- [ ] Testou em **375px** (mobile pequeno)?
- [ ] Testou em **768px** (tablet)?
- [ ] Testou em **1024px+** (desktop)?
- [ ] Grids usam padrão mobile-first (`grid-cols-1 md:grid-cols-2`)?
- [ ] Tabelas têm scroll horizontal ou versão mobile?
- [ ] Modais não cortam conteúdo em telas pequenas?
- [ ] Textos e botões têm tamanhos adequados para touch (min 44px)?

### Tema & Acessibilidade

- [ ] O tema Dark Mode funciona corretamente?
- [ ] Contraste de cores está adequado (WCAG AA)?

### Funcionalidade

- [ ] Formulários têm validação (Zod)?
- [ ] Estados de Loading e Error foram tratados?
- [ ] Navegação por teclado funciona?

---

## 7. Debugging

- **React Query Devtools:** Ótimo para debugar cache e queries.
- **Zustand Devtools:** Para ver o estado global.
- **Network Tab:** Verifique se o token JWT está sendo enviado no header `Authorization`.
