# NEXO — Componentes Críticos

> Lista de componentes e módulos vitais para o funcionamento do sistema.
> **Atenção:** Alterações nestes arquivos exigem revisão redobrada e testes de regressão.

---

## 1. Core Infrastructure

### 1.1 `src/app/layout.tsx` & `src/app/providers.tsx`

Responsáveis por envolver a aplicação com os contextos globais:

- `QueryClientProvider` (TanStack Query)
- `ThemeProvider` (next-themes)
- `Toaster` (Sonner)
- `AuthProvider` (se houver contexto de auth)

### 1.2 `src/lib/axios.ts`

Configuração do cliente HTTP.

- Interceptors de Request (Injeção de Token).
- Interceptors de Response (Tratamento de erro 401/Refresh Token).
- Base URL dinâmica.

### 1.3 `src/store/auth-store.ts`

Gerenciamento do estado de autenticação (Token, User).

- Persistência segura.
- Métodos de login/logout.

---

## 2. UI Components (Base)

Estes componentes são usados em >80% das telas. Quebrá-los quebra o sistema todo.

- **Button:** `src/components/ui/button.tsx`
- **Input:** `src/components/ui/input.tsx`
- **Card:** `src/components/ui/card.tsx`
- **Dialog:** `src/components/ui/dialog.tsx`
- **Form:** `src/components/ui/form.tsx` (Wrapper do RHF)

---

## 3. Business Components

### 3.1 `DataTable` (`src/components/ui/data-table.tsx`)

A tabela padrão do sistema.

- Deve suportar paginação server-side.
- Deve suportar ordenação e filtros.
- Deve ser responsiva.

### 3.2 `Sidebar` (`src/components/layout/sidebar.tsx`)

Navegação principal.

- Deve gerenciar estado mobile (drawer) e desktop (fixo/colapsável).
- Deve destacar a rota ativa.

### 3.3 `UserNav` (`src/components/layout/user-nav.tsx`)

Menu do usuário no header.

- Logout.
- Acesso ao perfil.
- Troca de tema.

---

## 4. Hooks Críticos

- **`useAuth`:** Verifica permissões e estado de login.
- **`useToast`:** Feedback ao usuário.
- **`useDebounce`:** Para inputs de busca em tempo real.

---

## 5. Testes de Regressão

Sempre que tocar nestes componentes, verificar:

1.  O login continua funcionando?
2.  A navegação entre páginas quebra?
3.  O tema Dark Mode aplica corretamente?
4.  Os formulários validam e enviam dados?
5.  As tabelas carregam dados e paginam?
