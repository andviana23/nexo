# Frontend Multi-Unit - Documentação Técnica

## Visão Geral

Este documento descreve a implementação do suporte a múltiplas unidades/filiais no frontend do NEXO.

## Arquitetura

```
src/
├── types/
│   └── unit.ts                    # Tipos TypeScript para units
├── store/
│   └── unit-store.ts              # Zustand store para estado de unit
├── hooks/
│   └── use-units.ts               # Hook que combina Zustand + React Query
├── services/
│   └── unit-service.ts            # API service para endpoints de units
├── components/
│   └── multi-unit/
│       ├── index.ts               # Exports centralizados
│       ├── UnitSelector.tsx       # Dropdown para troca de unidade
│       ├── UnitGuard.tsx          # HOC para proteção de rotas
│       └── UnitContextBanner.tsx  # Banner de contexto de unidade
└── lib/
    └── axios.ts                   # Interceptor adiciona X-Unit-ID
```

## Componentes

### UnitSelector
Dropdown para seleção de unidade ativa.

```tsx
import { UnitSelector } from '@/components/multi-unit';

// Variantes
<UnitSelector />                              // Default
<UnitSelector size="sm" variant="ghost" />    // Compacto
<UnitSelector collapsible={false} />          // Sempre mostra texto
```

**Props:**
- `size`: 'sm' | 'default' | 'lg'
- `variant`: 'default' | 'outline' | 'ghost'
- `collapsible`: boolean - esconde texto em mobile
- `onUnitChange`: callback quando unidade muda

### UnitGuard
Protege rotas/componentes que requerem unidade selecionada.

```tsx
import { UnitGuard } from '@/components/multi-unit';

<UnitGuard>
  <DashboardContent />
</UnitGuard>

// Com fallbacks customizados
<UnitGuard
  loadingFallback={<MyLoader />}
  noUnitFallback={<SelectUnitPrompt />}
>
  <Content />
</UnitGuard>
```

### UnitContextBanner
Banner que indica a unidade ativa.

```tsx
import { UnitContextBanner } from '@/components/multi-unit';

<UnitContextBanner />                    // Sutil
<UnitContextBanner variant="info" />     // Informativo
<UnitContextBanner dismissible />        // Com botão fechar
```

## Hooks

### useUnit()
Hook principal para interagir com unidades.

```tsx
import { useUnit } from '@/hooks/use-units';

function MyComponent() {
  const {
    units,           // Lista de unidades do usuário
    activeUnit,      // Unidade ativa
    activeUnitId,    // ID da unidade ativa
    isMultiUnit,     // Se tem múltiplas unidades
    isLoading,       // Carregando do servidor
    switchUnit,      // Função para trocar unidade
    setDefaultUnit,  // Função para definir padrão
    refreshUnits,    // Recarrega lista
  } = useUnit();

  return (
    <div>
      {isMultiUnit && <span>Multi-unit ativo</span>}
      <span>Unidade: {activeUnit?.unit_nome}</span>
    </div>
  );
}
```

### useActiveUnitId()
Hook simplificado para obter apenas o ID.

```tsx
import { useActiveUnitId } from '@/hooks/use-units';

const unitId = useActiveUnitId(); // string | null
```

## Store (Zustand)

### unit-store.ts
Estado persistido no localStorage.

```tsx
import {
  useUnits,
  useActiveUnit,
  useIsMultiUnit,
  getActiveUnitId, // Para uso fora de React
} from '@/store/unit-store';
```

**Estado persistido:**
- `activeUnit`: Unidade selecionada

**Estado não-persistido:**
- `units`: Lista de unidades (carregada do servidor)
- `isLoading`, `error`, `isHydrated`

## Service (API)

### unit-service.ts
Comunicação com backend.

```tsx
import { unitService } from '@/services/unit-service';

// Unidades do usuário logado
const { units } = await unitService.getUserUnits();

// Trocar unidade
await unitService.switchUnit(unitId);

// Definir padrão
await unitService.setDefaultUnit(unitId);

// CRUD admin
await unitService.listUnits();
await unitService.createUnit({ nome: 'Nova Filial' });
await unitService.updateUnit(id, { nome: 'Nome Atualizado' });
```

## Interceptor HTTP

O arquivo `src/lib/axios.ts` adiciona automaticamente o header `X-Unit-ID` em todas as requests quando há unidade ativa.

```typescript
// Adicionado automaticamente pelo interceptor
headers: {
  'Authorization': 'Bearer <token>',
  'X-Unit-ID': '<unit-id>'
}
```

## Fluxo de Login

1. Usuário faz login
2. Auth store salva token/user/tenant
3. React Query invalida `queryKeys.units.userUnits()`
4. Hook `useUnit` carrega unidades do servidor
5. Se `activeUnit` existe no localStorage, mantém
6. Senão, seleciona `is_default` ou primeira

## Fluxo de Logout

1. Usuário faz logout
2. Auth store limpa estado
3. Unit store reseta estado
4. Query cache limpo
5. Redirect para /login

## Integração com Componentes Existentes

O `UnitSelector` foi integrado no `Header.tsx` e aparece automaticamente quando o usuário tem acesso a múltiplas unidades.

```tsx
// src/components/layout/Header.tsx
<UnitSelector variant="outline" size="sm" collapsible />
```

## Tipos TypeScript

### Unit
```typescript
interface Unit {
  id: string;
  tenant_id: string;
  nome: string;
  apelido?: string;
  descricao?: string;
  endereco_resumo?: string;
  cidade?: string;
  estado?: string;
  timezone: string;
  ativa: boolean;
  is_matriz: boolean;
  criado_em: string;
  atualizado_em: string;
}
```

### UserUnit
```typescript
interface UserUnit {
  id: string;
  user_id: string;
  unit_id: string;
  unit_nome: string;
  unit_apelido?: string;
  unit_matriz: boolean;
  unit_ativa: boolean;
  is_default: boolean;
  role_override?: string;
  tenant_id: string;
}
```

## Query Keys

```typescript
queryKeys.units.all          // ['units']
queryKeys.units.userUnits()  // ['units', 'user-units']
queryKeys.units.detail(id)   // ['units', 'detail', id]
queryKeys.units.list()       // ['units', 'list']
```

## Considerações de UX

1. **Visibilidade**: O seletor só aparece quando há múltiplas unidades
2. **Persistência**: Unidade ativa persiste entre sessões
3. **Loading**: Skeleton enquanto carrega unidades
4. **Mobile**: Texto colapsa, mostra apenas ícone
5. **Feedback**: Badge mostra matriz, estrela mostra padrão

## Próximos Passos (Backend)

Para completar a funcionalidade, o backend precisa:

1. **Endpoint `/api/v1/units/me`**: Retornar unidades do usuário logado
2. **Leitura do header `X-Unit-ID`**: Filtrar dados por unidade
3. **Endpoint `/api/v1/units/switch`** (opcional): Trocar unidade com novo token
4. **Endpoint `/api/v1/units/default`**: Definir unidade padrão

## Testes

Recomendações para testes:

```tsx
// Mock do hook
jest.mock('@/hooks/use-units', () => ({
  useUnit: () => ({
    units: mockUnits,
    activeUnit: mockUnits[0],
    isMultiUnit: true,
    switchUnit: jest.fn(),
  }),
}));
```
