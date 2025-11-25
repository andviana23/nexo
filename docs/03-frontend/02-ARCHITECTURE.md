# NEXO — Frontend Architecture

> Esta documentação explica a arquitetura técnica do frontend do NEXO, incluindo stack, gerenciamento de estado, estilização e estrutura de pastas.

---

## 1. Stack Oficial

### 1.1 Core

- **Framework:** Next.js 16.0.4 (App Router)
- **Linguagem:** TypeScript 5.x
- **Runtime:** React 19.2.0 (Server Components + Actions)

### 1.2 UI & Styling

- **CSS Engine:** Tailwind CSS 4.1.17
- **Component Library:** shadcn/ui (Radix UI primitives)
- **Icons:** Lucide React
- **Animations:** Framer Motion 12.x
- **Theme:** CSS Variables (Light/Dark) via `next-themes`

### 1.3 State & Data

- **Global State:** Zustand 5.x (Client-side stores)
- **Server State:** TanStack Query 5.x (Data fetching, caching)
- **HTTP Client:** Axios (com interceptors para Auth)

### 1.4 Forms & Validation

- **Forms:** React Hook Form 7.x
- **Validation:** Zod 4.x
- **Resolver:** `@hookform/resolvers`

---

## 2. Estrutura de Pastas

A estrutura segue o padrão `src/` do Next.js App Router:

```
frontend/
├── src/
│   ├── app/                 # Rotas (App Router)
│   │   ├── (auth)/          # Rotas públicas (login, register)
│   │   ├── (dashboard)/     # Rotas protegidas (layout com sidebar)
│   │   ├── api/             # Route Handlers (se necessário)
│   │   ├── globals.css      # CSS global e variáveis
│   │   ├── layout.tsx       # Root Layout
│   │   └── page.tsx         # Landing page (ou redirect)
│   │
│   ├── components/          # Componentes React
│   │   ├── ui/              # Componentes shadcn/ui (Button, Input, etc.)
│   │   ├── forms/           # Componentes de formulário compostos
│   │   ├── layout/          # Sidebar, Header, Footer
│   │   └── shared/          # Componentes reutilizáveis do projeto
│   │
│   ├── hooks/               # Custom Hooks (useAuth, useToast, etc.)
│   ├── lib/                 # Utilitários e configurações
│   │   ├── axios.ts         # Instância configurada do Axios
│   │   ├── query-client.ts  # Configuração do React Query
│   │   └── utils.ts         # cn() e helpers gerais
│   │
│   ├── store/               # Stores do Zustand
│   │   ├── auth-store.ts    # Estado de autenticação
│   │   └── ui-store.ts      # Estado de UI (sidebar open, theme)
│   │
│   ├── types/               # Definições de tipos TypeScript
│   └── services/            # Camada de serviço (chamadas API)
│
├── public/                  # Assets estáticos
├── tailwind.config.ts       # Configuração do Tailwind (se necessário)
└── next.config.js           # Configuração do Next.js
```

---

## 3. Gerenciamento de Estado

### 3.1 Server State (TanStack Query)

Todo dado que vem da API é considerado **Server State**.

- Usar `useQuery` para buscar dados.
- Usar `useMutation` para alterar dados.
- **Nunca** armazenar dados da API no Zustand, a menos que seja estritamente necessário para persistência entre sessões.

### 3.2 Client State (Zustand)

Usado para estado global da aplicação que não vem do banco de dados.

- Tema (Light/Dark)
- Estado da Sidebar (aberta/fechada)
- Dados de Sessão do Usuário (Token, User Info)
- Filtros complexos que persistem na navegação

---

## 4. Estilização (Tailwind + shadcn/ui)

### 4.1 Tailwind CSS 4

- Usamos a versão 4 do Tailwind.
- Não usamos `@apply` excessivamente. Preferimos classes utilitárias diretamente no JSX.
- Para classes condicionais, usamos a função `cn()` (`clsx` + `tailwind-merge`).

### 4.2 shadcn/ui

- Componentes não são instalados como dependência npm, mas copiados para `src/components/ui`.
- Isso permite total controle e customização.
- **Regra:** Se precisar alterar um componente base (ex: Button), altere diretamente em `src/components/ui/button.tsx`.

### 4.3 Theming

- O tema é controlado via CSS Variables em `src/app/globals.css`.
- Suporte nativo a Dark Mode via classe `.dark`.
- Cores semânticas (`bg-primary`, `text-destructive`) em vez de cores fixas (`bg-blue-500`).

---

## 5. Autenticação

- **Auth Store:** Zustand persiste o token JWT no `localStorage` (ou cookies).
- **Axios Interceptor:** Injeta o token `Authorization: Bearer ...` em toda requisição.
- **Middleware:** Protege rotas `/dashboard` verificando a presença do token.
- **Refresh Token:** (Se implementado) Lógica no interceptor de resposta do Axios.

---

## 6. Performance & Best Practices

- **Server Components:** Por padrão, tudo em `src/app` é Server Component. Use `"use client"` apenas quando precisar de interatividade (hooks, eventos).
- **Imagens:** Use `next/image` para otimização automática.
- **Fontes:** Use `next/font` para carregar Inter e JetBrains Mono sem CLS (Cumulative Layout Shift).
- **Code Splitting:** Automático pelo Next.js. Evite imports gigantes no layout principal.

---

## 7. Integração com Backend

- O frontend consome a API Go via REST.
- URLs da API definidas em variáveis de ambiente (`NEXT_PUBLIC_API_URL`).
- Tipagem: Interfaces TypeScript em `src/types` devem espelhar os DTOs do backend.
