# NEXO — Patterns

> Padrões de desenvolvimento obrigatórios para garantir consistência, manutenibilidade e qualidade de código.

---

## 1. Formulários (Zod + React Hook Form)

Todos os formulários devem seguir a arquitetura **Schema-First**.

### 1.1 Definição do Schema (Zod)

Defina o schema de validação fora do componente ou em um arquivo separado.

```tsx
import { z } from 'zod';

const loginSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string().min(6, 'A senha deve ter no mínimo 6 caracteres'),
});

type LoginFormValues = z.infer<typeof loginSchema>;
```

### 1.2 Implementação (React Hook Form)

Use o hook `useForm` com o `zodResolver`.

```tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

export function LoginForm() {
  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  function onSubmit(data: LoginFormValues) {
    // Call API
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>{/* Fields */}</form>
    </Form>
  );
}
```

---

## 2. Data Fetching (TanStack Query)

### 2.1 Queries (GET)

Use `useQuery` para buscar dados. Sempre tipar o retorno.

```tsx
// hooks/use-clients.ts
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/axios';
import { Client } from '@/types';

export function useClients() {
  return useQuery({
    queryKey: ['clients'],
    queryFn: async () => {
      const { data } = await api.get<Client[]>('/clients');
      return data;
    },
  });
}
```

### 2.2 Mutations (POST, PUT, DELETE)

Use `useMutation` para alterar dados. Lembre-se de invalidar as queries relacionadas após o sucesso.

```tsx
import { useMutation, useQueryClient } from '@tanstack/react-query';

export function useCreateClient() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (newClient: CreateClientDTO) => {
      return api.post('/clients', newClient);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clients'] });
      toast.success('Cliente criado com sucesso!');
    },
  });
}
```

---

## 3. Tratamento de Erros

### 3.1 API Errors

O `axios` interceptor deve tratar erros globais (401, 500).
Para erros de validação (400), o `useMutation` pode receber um callback `onError`.

```tsx
onError: (error) => {
  if (isAxiosError(error) && error.response?.status === 400) {
    toast.error('Dados inválidos. Verifique os campos.');
  } else {
    toast.error('Erro inesperado. Tente novamente.');
  }
};
```

### 3.2 Error Boundaries

Use Error Boundaries (arquivo `error.tsx` no Next.js) para capturar erros de renderização em páginas inteiras.

---

## 4. Layout & Responsividade (OBRIGATÓRIO)

> ⚠️ **Responsividade NÃO é opcional.** Toda página deve funcionar perfeitamente de 375px a 1920px+.

### 4.1 Mobile First (Regra de Ouro)

Escreva classes Tailwind pensando **primeiro no mobile**, depois use breakpoints (`md:`, `lg:`) para telas maiores.

```tsx
// ✅ CORRETO: Base mobile, expande para desktop
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* 1 coluna no mobile, 2 no tablet, 3 no desktop */}
</div>

// ❌ ERRADO: Nunca faça desktop-first
<div className="grid grid-cols-3 sm:grid-cols-1"> {/* NÃO */}
```

### 4.2 Container Padrão

Use um container padrão para limitar a largura do conteúdo:

```tsx
<div className="container mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
  {/* Conteúdo da página */}
</div>
```

### 4.3 Padrões para Grids

```tsx
// Cards de Dashboard/KPIs
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">

// Formulários com campos lado a lado
<div className="grid grid-cols-1 md:grid-cols-2 gap-4">

// Layout com Sidebar
<div className="flex flex-col lg:flex-row">
  <aside className="w-full lg:w-64 lg:shrink-0">
  <main className="flex-1 min-w-0">
</div>
```

### 4.4 Tabelas Responsivas

Tabelas devem ter alternativa mobile:

```tsx
// Opção 1: Scroll horizontal
<div className="overflow-x-auto -mx-4 sm:mx-0">
  <Table className="min-w-[600px]">
    {/* ... */}
  </Table>
</div>

// Opção 2: Cards no mobile (para listas simples)
<div className="hidden md:block">
  <Table>{/* Tabela completa */}</Table>
</div>
<div className="md:hidden space-y-4">
  {items.map(item => <MobileCard key={item.id} {...item} />)}
</div>
```

### 4.5 Modais Responsivos

```tsx
<DialogContent className="sm:max-w-[425px] md:max-w-[600px] max-h-[90vh] overflow-y-auto">
  {/* Em mobile, modal ocupa quase toda a tela */}
  {/* Em desktop, tem largura máxima definida */}
</DialogContent>
```

### 4.6 Tipografia Responsiva

```tsx
// Títulos de página
<h1 className="text-2xl sm:text-3xl lg:text-4xl font-bold">

// Subtítulos
<h2 className="text-xl sm:text-2xl font-semibold">

// Texto de métricas/KPIs
<span className="text-3xl sm:text-4xl lg:text-5xl font-bold">
```

---

## 5. Convenções de Naming

- **Componentes:** PascalCase (`ClientCard.tsx`)
- **Hooks:** camelCase com prefixo use (`useAuth.ts`)
- **Utilitários:** camelCase (`formatDate.ts`)
- **Tipos/Interfaces:** PascalCase (`Client`, `UserResponse`)
- **Pastas:** kebab-case (`client-list`, `auth-provider`)

---

## 6. Imports

Use o alias `@/` para imports absolutos a partir de `src/`.

```tsx
// Bom
import { Button } from '@/components/ui/button';
import { api } from '@/lib/axios';

// Ruim
import { Button } from '../../../components/ui/button';
```
