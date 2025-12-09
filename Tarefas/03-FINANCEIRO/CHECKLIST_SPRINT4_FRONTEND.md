# ✅ CHECKLIST — SPRINT 4: FRONTEND

> **Status:** 🟢 100% — Concluído ✅  
> **Dependência:** Sprint 3 (Painel Mensal + Projeções) ✅  
> **Esforço Estimado:** 40 horas  
> **Prioridade:** P1 — Entrega de valor para o usuário  
> **Data Conclusão:** 29/11/2025

---

## 🔧 CORREÇÕES APLICADAS

### Migração Database (29/11/2025)
- [x] Migração `026_despesas_fixas.up.sql` aplicada com sucesso
- [x] Tabela `despesas_fixas` criada no PostgreSQL/Neon
- [x] Índices para performance e multi-tenant
- [x] RLS habilitado para isolamento de tenant
- [x] Trigger para `atualizado_em` automático

### Endpoints Testados
- [x] `POST /api/v1/financial/fixed-expenses` — Criar despesa ✅
- [x] `GET /api/v1/financial/fixed-expenses` — Listar despesas ✅
- [x] `GET /api/v1/financial/fixed-expenses/summary` — Resumo ✅
- [x] `PUT /api/v1/financial/fixed-expenses/:id` — Atualizar ✅
- [x] `DELETE /api/v1/financial/fixed-expenses/:id` — Excluir ✅
- [x] `PATCH /api/v1/financial/fixed-expenses/:id/toggle` — Ativar/Desativar ✅

---

## 📊 OBJETIVO

Implementar todas as **telas do módulo financeiro** usando:

- Next.js App Router
- React Hook Form + Zod
- TanStack Query (React Query)
- Design System (shadcn/ui + tokens)

---

## ✅ PROGRESSO ATUAL — 100% CONCLUÍDO

### 1. Dashboard Financeiro (`financeiro/page.tsx`) ✅
   - [x] KPI Cards (Receita, Despesas, Lucro, Saldo)
   - [x] Próximos Vencimentos
   - [x] Meta do Mês com barra de progresso
   - [x] Projeções (3 meses) com badges de confiança
   - [x] DRE Resumido
   - [x] Ações Rápidas
   - [x] Integração com useDashboard() e useProjections()

### 2. Tipos e Services (`types/financial.ts`, `services/financial-service.ts`) ✅
   - [x] DespesaFixa interface
   - [x] CreateDespesaFixaRequest, UpdateDespesaFixaRequest
   - [x] DespesasFixasListResponse, DespesasFixasSummaryResponse
   - [x] GerarContasRequest, GerarContasResponse
   - [x] ListDespesasFixasFilters
   - [x] Endpoints CRUD despesas fixas
   - [x] Endpoint toggle ativar/desativar
   - [x] Endpoint gerar contas a partir de despesas fixas
   - [x] Endpoint resumo despesas fixas

### 3. Hooks (`hooks/use-financial.ts`) ✅
   - [x] useFixedExpenses()
   - [x] useFixedExpense()
   - [x] useFixedExpensesSummary()
   - [x] useCreateFixedExpense()
   - [x] useUpdateFixedExpense()
   - [x] useDeleteFixedExpense()
   - [x] useToggleFixedExpense()
   - [x] useGeneratePayablesFromFixed()

### 4. Páginas Despesas Fixas ✅
   - [x] Lista (`financeiro/despesas-fixas/page.tsx`)
     - Cards resumo (total, ativas, valor mensal)
     - Tabela com descrição, fornecedor, dia vencimento, valor, status
     - Filtro por status (todas/ativas/inativas)
     - **Paginação completa com navegação**
     - Toggle ativar/desativar
     - Modal gerar contas (ano/mês)
     - Modal confirmar exclusão
   - [x] Nova (`financeiro/despesas-fixas/nova/page.tsx`)
     - Formulário com Zod validation
     - Campos: descrição, fornecedor, valor, dia vencimento, observações
     - Breadcrumbs
   - [x] Editar (`financeiro/despesas-fixas/[id]/editar/page.tsx`)
     - Carrega dados existentes
     - Formulário com validação
     - Loading skeleton
     - Tratamento de erro

### 5. Navegação Sidebar ✅
   - [x] Menu Financeiro expandido com subitems:
     - Dashboard
     - Contas a Pagar
     - Contas a Receber
     - Despesas Fixas
     - DRE (owner/admin)
     - Fluxo de Caixa (owner/admin)

### 6. Páginas existentes (já funcionais) ✅
   - [x] Contas a Pagar (`financeiro/contas-pagar/page.tsx`)
   - [x] Contas a Receber (`financeiro/contas-receber/page.tsx`)
   - [x] DRE (`financeiro/dre/page.tsx`)
   - [x] Fluxo de Caixa (`financeiro/fluxo-caixa/page.tsx`)

6. **Melhorias futuras (P2)**
   - [ ] Paginação nas listas
   - [ ] Filtros avançados por período
   - [ ] Export CSV/PDF
   - [ ] Gráficos interativos

---

## 📁 ESTRUTURA DE PASTAS

```
frontend/src/app/(dashboard)/financeiro/
├── page.tsx                         ← Dashboard principal
├── layout.tsx                       ← Layout do módulo
├── contas-a-pagar/
│   ├── page.tsx                     ← Lista
│   ├── nova/
│   │   └── page.tsx                 ← Criar nova
│   └── [id]/
│       ├── page.tsx                 ← Detalhes
│       └── editar/
│           └── page.tsx             ← Editar
├── contas-a-receber/
│   ├── page.tsx                     ← Lista
│   ├── nova/
│   │   └── page.tsx                 ← Criar nova
│   └── [id]/
│       ├── page.tsx                 ← Detalhes
│       └── editar/
│           └── page.tsx             ← Editar
├── despesas-fixas/
│   ├── page.tsx                     ← Lista
│   ├── nova/
│   │   └── page.tsx                 ← Criar nova
│   └── [id]/
│       └── editar/
│           └── page.tsx             ← Editar
├── caixa/
│   └── page.tsx                     ← Fluxo de caixa diário
└── dre/
    └── page.tsx                     ← DRE mensal
```

---

## 📋 TAREFAS

### 1️⃣ INFRAESTRUTURA (Esforço: 4h)

#### 1.1 API Client

- [ ] Criar `frontend/src/lib/api/financial.ts`

```typescript
import { apiClient } from './client';

// Contas a Pagar
export const contasPagarApi = {
  list: (params?: ListParams) => 
    apiClient.get<ContasPagarListResponse>('/financial/payables', { params }),
  
  get: (id: string) => 
    apiClient.get<ContaPagarResponse>(`/financial/payables/${id}`),
  
  create: (data: ContaPagarCreateRequest) => 
    apiClient.post<ContaPagarResponse>('/financial/payables', data),
  
  update: (id: string, data: ContaPagarUpdateRequest) => 
    apiClient.put<ContaPagarResponse>(`/financial/payables/${id}`, data),
  
  delete: (id: string) => 
    apiClient.delete(`/financial/payables/${id}`),
  
  markAsPaid: (id: string, data: MarcarPagamentoRequest) => 
    apiClient.post<ContaPagarResponse>(`/financial/payables/${id}/payment`, data),
};

// Contas a Receber
export const contasReceberApi = {
  list: (params?: ListParams) => 
    apiClient.get<ContasReceberListResponse>('/financial/receivables', { params }),
  
  get: (id: string) => 
    apiClient.get<ContaReceberResponse>(`/financial/receivables/${id}`),
  
  create: (data: ContaReceberCreateRequest) => 
    apiClient.post<ContaReceberResponse>('/financial/receivables', data),
  
  update: (id: string, data: ContaReceberUpdateRequest) => 
    apiClient.put<ContaReceberResponse>(`/financial/receivables/${id}`, data),
  
  delete: (id: string) => 
    apiClient.delete(`/financial/receivables/${id}`),
  
  markAsReceived: (id: string, data: MarcarRecebimentoRequest) => 
    apiClient.post<ContaReceberResponse>(`/financial/receivables/${id}/receipt`, data),
};

// Despesas Fixas
export const despesasFixasApi = {
  list: (params?: ListParams) => 
    apiClient.get<DespesasFixasListResponse>('/financial/fixed-expenses', { params }),
  
  get: (id: string) => 
    apiClient.get<DespesaFixaResponse>(`/financial/fixed-expenses/${id}`),
  
  create: (data: DespesaFixaCreateRequest) => 
    apiClient.post<DespesaFixaResponse>('/financial/fixed-expenses', data),
  
  update: (id: string, data: DespesaFixaUpdateRequest) => 
    apiClient.put<DespesaFixaResponse>(`/financial/fixed-expenses/${id}`, data),
  
  toggle: (id: string) => 
    apiClient.post<DespesaFixaResponse>(`/financial/fixed-expenses/${id}/toggle`),
  
  delete: (id: string) => 
    apiClient.delete(`/financial/fixed-expenses/${id}`),
};

// Dashboard
export const dashboardApi = {
  get: (year?: number, month?: number) => 
    apiClient.get<PainelMensalResponse>('/financial/dashboard', { 
      params: { year, month } 
    }),
  
  getProjections: (monthsAhead?: number) => 
    apiClient.get<ProjecoesResponse>('/financial/projections', { 
      params: { months_ahead: monthsAhead } 
    }),
};

// Caixa
export const caixaApi = {
  list: (params?: { start_date?: string; end_date?: string }) => 
    apiClient.get<FluxoCaixaListResponse>('/financial/cashflow', { params }),
  
  get: (date: string) => 
    apiClient.get<FluxoCaixaResponse>(`/financial/cashflow/${date}`),
};

// DRE
export const dreApi = {
  list: (year?: number) => 
    apiClient.get<DREListResponse>('/financial/dre', { params: { year } }),
  
  get: (year: number, month: number) => 
    apiClient.get<DREResponse>(`/financial/dre/${year}/${month}`),
};
```

#### 1.2 React Query Hooks

- [ ] Criar `frontend/src/hooks/useFinancial.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { contasPagarApi, contasReceberApi, despesasFixasApi, dashboardApi } from '@/lib/api/financial';

// Query Keys
export const financialKeys = {
  all: ['financial'] as const,
  contasPagar: () => [...financialKeys.all, 'contas-pagar'] as const,
  contasPagarList: (params?: ListParams) => [...financialKeys.contasPagar(), 'list', params] as const,
  contasPagarDetail: (id: string) => [...financialKeys.contasPagar(), 'detail', id] as const,
  // ... repetir para outras entidades
  dashboard: (year?: number, month?: number) => [...financialKeys.all, 'dashboard', year, month] as const,
  projections: (months?: number) => [...financialKeys.all, 'projections', months] as const,
};

// Contas a Pagar
export function useContasPagar(params?: ListParams) {
  return useQuery({
    queryKey: financialKeys.contasPagarList(params),
    queryFn: () => contasPagarApi.list(params),
  });
}

export function useContaPagar(id: string) {
  return useQuery({
    queryKey: financialKeys.contasPagarDetail(id),
    queryFn: () => contasPagarApi.get(id),
    enabled: !!id,
  });
}

export function useCreateContaPagar() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: contasPagarApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.contasPagar() });
    },
  });
}

export function useUpdateContaPagar() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ContaPagarUpdateRequest }) => 
      contasPagarApi.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.contasPagar() });
      queryClient.invalidateQueries({ queryKey: financialKeys.contasPagarDetail(id) });
    },
  });
}

export function useMarcarPagamento() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: MarcarPagamentoRequest }) => 
      contasPagarApi.markAsPaid(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.contasPagar() });
      queryClient.invalidateQueries({ queryKey: financialKeys.dashboard() });
    },
  });
}

// Dashboard
export function useDashboard(year?: number, month?: number) {
  return useQuery({
    queryKey: financialKeys.dashboard(year, month),
    queryFn: () => dashboardApi.get(year, month),
  });
}

export function useProjections(monthsAhead?: number) {
  return useQuery({
    queryKey: financialKeys.projections(monthsAhead),
    queryFn: () => dashboardApi.getProjections(monthsAhead),
  });
}

// ... continuar para outras entidades
```

#### 1.3 Types

- [ ] Criar `frontend/src/types/financial.ts`

```typescript
// Contas a Pagar
export interface ContaPagar {
  id: string;
  descricao: string;
  categoria_id?: string;
  categoria_nome?: string;
  fornecedor?: string;
  valor: string;
  tipo: 'DESPESA_FIXA' | 'DESPESA_VARIAVEL' | 'INVESTIMENTO' | 'OUTROS';
  recorrente: boolean;
  periodicidade?: 'MENSAL' | 'SEMANAL' | 'QUINZENAL' | 'ANUAL';
  data_vencimento: string;
  data_pagamento?: string;
  status: 'ABERTO' | 'PAGO' | 'ATRASADO' | 'CANCELADO';
  comprovante_url?: string;
  pix_code?: string;
  observacoes?: string;
  criado_em: string;
  atualizado_em: string;
}

export interface ContaPagarCreateRequest {
  descricao: string;
  categoria_id?: string;
  fornecedor?: string;
  valor: string;
  tipo: string;
  recorrente: boolean;
  periodicidade?: string;
  data_vencimento: string;
  observacoes?: string;
}

// ... continuar para outras interfaces
```

#### 1.4 Zod Schemas

- [ ] Criar `frontend/src/lib/validations/financial.ts`

```typescript
import { z } from 'zod';

export const contaPagarSchema = z.object({
  descricao: z.string().min(1, 'Descrição é obrigatória').max(255),
  categoria_id: z.string().uuid().optional(),
  fornecedor: z.string().max(255).optional(),
  valor: z.string().refine((val) => {
    const num = parseFloat(val.replace(',', '.'));
    return !isNaN(num) && num > 0;
  }, 'Valor deve ser maior que zero'),
  tipo: z.enum(['DESPESA_FIXA', 'DESPESA_VARIAVEL', 'INVESTIMENTO', 'OUTROS']),
  recorrente: z.boolean().default(false),
  periodicidade: z.enum(['MENSAL', 'SEMANAL', 'QUINZENAL', 'ANUAL']).optional(),
  data_vencimento: z.string().refine((val) => !isNaN(Date.parse(val)), 'Data inválida'),
  observacoes: z.string().optional(),
});

export type ContaPagarFormData = z.infer<typeof contaPagarSchema>;

export const despesaFixaSchema = z.object({
  descricao: z.string().min(1, 'Descrição é obrigatória').max(255),
  categoria_id: z.string().uuid().optional(),
  unidade_id: z.string().uuid().optional(),
  fornecedor: z.string().max(255).optional(),
  valor: z.string().refine((val) => {
    const num = parseFloat(val.replace(',', '.'));
    return !isNaN(num) && num > 0;
  }, 'Valor deve ser maior que zero'),
  dia_vencimento: z.number().min(1).max(31),
  observacoes: z.string().optional(),
});

export type DespesaFixaFormData = z.infer<typeof despesaFixaSchema>;
```

#### 1.5 Checklist Infraestrutura

- [ ] API Client com todos os endpoints
- [ ] React Query hooks
- [ ] TypeScript interfaces
- [ ] Zod schemas de validação
- [ ] Query keys organizadas

---

### 2️⃣ COMPONENTES COMPARTILHADOS (Esforço: 6h)

#### 2.1 Componentes de Input

- [ ] `<CurrencyInput />` — Input monetário formatado

```typescript
// frontend/src/components/ui/currency-input.tsx
'use client';

import { forwardRef } from 'react';
import { Input } from '@/components/ui/input';

interface CurrencyInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  onChange?: (value: string) => void;
}

export const CurrencyInput = forwardRef<HTMLInputElement, CurrencyInputProps>(
  ({ onChange, value, ...props }, ref) => {
    const formatCurrency = (val: string) => {
      // Remove non-numeric
      const numbers = val.replace(/\D/g, '');
      // Convert to decimal
      const decimal = (parseInt(numbers, 10) / 100).toFixed(2);
      // Format with locale
      return new Intl.NumberFormat('pt-BR', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      }).format(parseFloat(decimal));
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const formatted = formatCurrency(e.target.value);
      onChange?.(formatted);
    };

    return (
      <div className="relative">
        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
          R$
        </span>
        <Input
          ref={ref}
          {...props}
          value={value}
          onChange={handleChange}
          className="pl-10"
        />
      </div>
    );
  }
);
```

- [ ] `<DatePicker />` — Seletor de data
- [ ] `<CategorySelect />` — Select de categorias
- [ ] `<StatusSelect />` — Select de status

#### 2.2 Componentes de Display

- [ ] `<StatusBadge />` — Badge colorido por status

```typescript
// frontend/src/components/financial/status-badge.tsx
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface StatusBadgeProps {
  status: 'ABERTO' | 'PAGO' | 'ATRASADO' | 'CANCELADO' | 'PENDENTE' | 'RECEBIDO';
}

const statusConfig = {
  ABERTO: { label: 'Aberto', variant: 'outline' as const, className: 'border-blue-500 text-blue-500' },
  PAGO: { label: 'Pago', variant: 'default' as const, className: 'bg-green-500' },
  ATRASADO: { label: 'Atrasado', variant: 'destructive' as const, className: '' },
  CANCELADO: { label: 'Cancelado', variant: 'secondary' as const, className: '' },
  PENDENTE: { label: 'Pendente', variant: 'outline' as const, className: 'border-yellow-500 text-yellow-500' },
  RECEBIDO: { label: 'Recebido', variant: 'default' as const, className: 'bg-green-500' },
};

export function StatusBadge({ status }: StatusBadgeProps) {
  const config = statusConfig[status];
  
  return (
    <Badge variant={config.variant} className={cn(config.className)}>
      {config.label}
    </Badge>
  );
}
```

- [ ] `<FinancialCard />` — Card de métrica

```typescript
// frontend/src/components/financial/financial-card.tsx
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

interface FinancialCardProps {
  title: string;
  value: string;
  subtitle?: string;
  trend?: 'up' | 'down' | 'stable';
  trendValue?: string;
  icon?: React.ReactNode;
  variant?: 'default' | 'success' | 'warning' | 'danger';
}

export function FinancialCard({
  title,
  value,
  subtitle,
  trend,
  trendValue,
  icon,
  variant = 'default',
}: FinancialCardProps) {
  const TrendIcon = trend === 'up' ? TrendingUp : trend === 'down' ? TrendingDown : Minus;
  
  return (
    <Card className={cn(
      variant === 'success' && 'border-green-200 bg-green-50/50',
      variant === 'warning' && 'border-yellow-200 bg-yellow-50/50',
      variant === 'danger' && 'border-red-200 bg-red-50/50',
    )}>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">R$ {value}</div>
        {(subtitle || trend) && (
          <div className="flex items-center gap-1 text-xs text-muted-foreground">
            {trend && (
              <TrendIcon className={cn(
                'h-3 w-3',
                trend === 'up' && 'text-green-500',
                trend === 'down' && 'text-red-500',
              )} />
            )}
            {trendValue && <span>{trendValue}</span>}
            {subtitle && <span>{subtitle}</span>}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
```

- [ ] `<TrendIndicator />` — Indicador de tendência
- [ ] `<ProgressMeta />` — Barra de progresso da meta

#### 2.3 Componentes de Tabela

- [ ] `<FinancialTable />` — Tabela genérica
- [ ] `<ContasPagarTable />` — Tabela específica
- [ ] `<ContasReceberTable />` — Tabela específica
- [ ] `<DespesasFixasTable />` — Tabela específica

#### 2.4 Checklist Componentes

- [ ] CurrencyInput
- [ ] DatePicker
- [ ] CategorySelect
- [ ] StatusBadge
- [ ] FinancialCard
- [ ] TrendIndicator
- [ ] ProgressMeta
- [ ] FinancialTable
- [ ] Testes de componentes

---

### 3️⃣ TELA: DASHBOARD FINANCEIRO (Esforço: 8h)

#### 3.1 Página Principal

- [ ] Criar `frontend/src/app/(dashboard)/financeiro/page.tsx`

```typescript
'use client';

import { useDashboard, useProjections } from '@/hooks/useFinancial';
import { FinancialCard } from '@/components/financial/financial-card';
import { MonthSelector } from '@/components/financial/month-selector';
import { ProgressMeta } from '@/components/financial/progress-meta';
import { RecentTransactions } from '@/components/financial/recent-transactions';
import { ProjectionsChart } from '@/components/financial/projections-chart';
import { useState } from 'react';

export default function FinanceiroPage() {
  const [selectedDate, setSelectedDate] = useState({
    year: new Date().getFullYear(),
    month: new Date().getMonth() + 1,
  });
  
  const { data: dashboard, isLoading } = useDashboard(
    selectedDate.year, 
    selectedDate.month
  );
  const { data: projections } = useProjections(3);
  
  if (isLoading) return <DashboardSkeleton />;
  
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Painel Financeiro</h1>
        <MonthSelector 
          value={selectedDate} 
          onChange={setSelectedDate} 
        />
      </div>
      
      {/* Cards de Métricas */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <FinancialCard
          title="Receita Realizada"
          value={dashboard.receita_realizada}
          trend={dashboard.tendencia_variacao as 'up' | 'down' | 'stable'}
          trendValue={`${dashboard.variacao_mes_anterior}% vs mês anterior`}
          variant="success"
        />
        <FinancialCard
          title="Despesas"
          value={dashboard.despesas_total}
          subtitle={`Fixas: R$ ${dashboard.despesas_fixas}`}
          variant="danger"
        />
        <FinancialCard
          title="Lucro Líquido"
          value={dashboard.lucro_liquido}
          subtitle={`Margem: ${dashboard.margem_liquida}%`}
          variant={parseFloat(dashboard.lucro_liquido) >= 0 ? 'success' : 'danger'}
        />
        <FinancialCard
          title="Ticket Médio"
          value={dashboard.ticket_medio}
          subtitle={`${dashboard.total_atendimentos} atendimentos`}
        />
      </div>
      
      {/* Meta do Mês */}
      <ProgressMeta
        meta={dashboard.meta_mensal}
        realizado={dashboard.receita_realizada}
        percentual={dashboard.percentual_meta}
        status={dashboard.status_meta}
      />
      
      {/* Grid de Projeções e Transações */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ProjectionsChart data={projections?.projecoes || []} />
        <RecentTransactions />
      </div>
      
      {/* Contas a Vencer */}
      <UpcomingPayments />
    </div>
  );
}
```

#### 3.2 Componentes do Dashboard

- [ ] `<MonthSelector />` — Seletor de mês/ano
- [ ] `<ProgressMeta />` — Barra de progresso da meta
- [ ] `<ProjectionsChart />` — Gráfico de projeções (Recharts)
- [ ] `<RecentTransactions />` — Lista de transações recentes
- [ ] `<UpcomingPayments />` — Contas a vencer na semana
- [ ] `<DashboardSkeleton />` — Loading state

#### 3.3 Checklist Dashboard

- [ ] Cards de métricas principais
- [ ] Seletor de mês/ano
- [ ] Barra de progresso da meta
- [ ] Gráfico de projeções
- [ ] Lista de transações recentes
- [ ] Contas a vencer
- [ ] Loading states
- [ ] Error handling
- [ ] Responsividade

---

### 4️⃣ TELA: CONTAS A PAGAR (Esforço: 8h)

#### 4.1 Lista

- [ ] Criar `frontend/src/app/(dashboard)/financeiro/contas-a-pagar/page.tsx`

```typescript
'use client';

import { useContasPagar, useMarcarPagamento, useDeleteContaPagar } from '@/hooks/useFinancial';
import { DataTable } from '@/components/ui/data-table';
import { Button } from '@/components/ui/button';
import { Plus, Filter } from 'lucide-react';
import Link from 'next/link';
import { columns } from './columns';
import { ContasPagarFilters } from './filters';

export default function ContasPagarPage() {
  const [filters, setFilters] = useState<ContasPagarFilters>({});
  const { data, isLoading } = useContasPagar(filters);
  const marcarPagamento = useMarcarPagamento();
  const deleteConta = useDeleteContaPagar();
  
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Contas a Pagar</h1>
          <p className="text-muted-foreground">
            Gerencie suas contas e despesas
          </p>
        </div>
        <Button asChild>
          <Link href="/financeiro/contas-a-pagar/nova">
            <Plus className="mr-2 h-4 w-4" />
            Nova Conta
          </Link>
        </Button>
      </div>
      
      {/* Resumo */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <SummaryCard title="Total a Pagar" value={data?.summary.total || '0,00'} />
        <SummaryCard title="Vencidas" value={data?.summary.vencidas || '0,00'} variant="danger" />
        <SummaryCard title="A Vencer" value={data?.summary.a_vencer || '0,00'} variant="warning" />
        <SummaryCard title="Pagas no Mês" value={data?.summary.pagas || '0,00'} variant="success" />
      </div>
      
      {/* Filtros */}
      <ContasPagarFilters filters={filters} onFilterChange={setFilters} />
      
      {/* Tabela */}
      <DataTable
        columns={columns}
        data={data?.items || []}
        isLoading={isLoading}
        onMarkAsPaid={(id) => marcarPagamento.mutate({ id, data: { data_pagamento: new Date().toISOString() } })}
        onDelete={(id) => deleteConta.mutate(id)}
      />
    </div>
  );
}
```

#### 4.2 Formulário de Criação

- [ ] Criar `frontend/src/app/(dashboard)/financeiro/contas-a-pagar/nova/page.tsx`

```typescript
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { contaPagarSchema, ContaPagarFormData } from '@/lib/validations/financial';
import { useCreateContaPagar } from '@/hooks/useFinancial';
import { useRouter } from 'next/navigation';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { CurrencyInput } from '@/components/ui/currency-input';
import { DatePicker } from '@/components/ui/date-picker';
import { CategorySelect } from '@/components/ui/category-select';
import { toast } from 'sonner';

export default function NovaContaPagarPage() {
  const router = useRouter();
  const createConta = useCreateContaPagar();
  
  const form = useForm<ContaPagarFormData>({
    resolver: zodResolver(contaPagarSchema),
    defaultValues: {
      recorrente: false,
      tipo: 'DESPESA_VARIAVEL',
    },
  });
  
  const onSubmit = async (data: ContaPagarFormData) => {
    try {
      await createConta.mutateAsync(data);
      toast.success('Conta criada com sucesso!');
      router.push('/financeiro/contas-a-pagar');
    } catch (error) {
      toast.error('Erro ao criar conta');
    }
  };
  
  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Nova Conta a Pagar</h1>
      
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <FormField
            control={form.control}
            name="descricao"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Descrição</FormLabel>
                <FormControl>
                  <Input placeholder="Ex: Aluguel" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          
          <div className="grid grid-cols-2 gap-4">
            <FormField
              control={form.control}
              name="valor"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Valor</FormLabel>
                  <FormControl>
                    <CurrencyInput {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            
            <FormField
              control={form.control}
              name="data_vencimento"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Data de Vencimento</FormLabel>
                  <FormControl>
                    <DatePicker {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          
          {/* ... outros campos */}
          
          <div className="flex gap-4">
            <Button type="button" variant="outline" onClick={() => router.back()}>
              Cancelar
            </Button>
            <Button type="submit" disabled={createConta.isPending}>
              {createConta.isPending ? 'Salvando...' : 'Salvar'}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
```

#### 4.3 Checklist Contas a Pagar

- [ ] Página de listagem
- [ ] Filtros (status, período, categoria)
- [ ] Tabela com ações
- [ ] Formulário de criação
- [ ] Formulário de edição
- [ ] Modal de confirmação de exclusão
- [ ] Modal de marcar como pago
- [ ] Indicadores visuais de vencimento
- [ ] Paginação

---

### 5️⃣ TELA: CONTAS A RECEBER (Esforço: 6h)

- [ ] Página de listagem
- [ ] Filtros (status, origem, período)
- [ ] Tabela com ações
- [ ] Formulário de criação
- [ ] Formulário de edição
- [ ] Modal de marcar como recebido
- [ ] Indicadores de atraso
- [ ] Paginação

---

### 6️⃣ TELA: DESPESAS FIXAS (Esforço: 4h)

- [ ] Página de listagem
- [ ] Toggle ativo/inativo inline
- [ ] Formulário de criação
- [ ] Formulário de edição
- [ ] Preview do próximo vencimento
- [ ] Total mensal calculado

---

### 7️⃣ TELA: CAIXA DIÁRIO (Esforço: 4h)

- [ ] Calendário de dias
- [ ] Detalhes do dia selecionado
- [ ] Resumo do período
- [ ] Exportar relatório

---

### 8️⃣ TELA: DRE (Esforço: 4h)

- [ ] Visualização mensal
- [ ] Comparativo anual
- [ ] Gráficos de evolução
- [ ] Exportar para PDF

---

### 9️⃣ TESTES (Esforço: 6h)

#### 9.1 Testes de Componentes

- [ ] CurrencyInput
- [ ] StatusBadge
- [ ] FinancialCard
- [ ] ProgressMeta

#### 9.2 Testes de Hooks

- [ ] useContasPagar
- [ ] useDashboard
- [ ] useProjections

#### 9.3 Testes E2E (Playwright)

- [ ] Fluxo: criar conta a pagar
- [ ] Fluxo: marcar como pago
- [ ] Fluxo: dashboard

---

## 📊 ESTIMATIVA DE TEMPO

| Tarefa | Horas |
|--------|:-----:|
| Infraestrutura (API, hooks, types) | 4h |
| Componentes compartilhados | 6h |
| Dashboard Financeiro | 8h |
| Contas a Pagar | 8h |
| Contas a Receber | 6h |
| Despesas Fixas | 4h |
| Caixa Diário | 4h |
| DRE | 4h |
| Testes | 6h |
| **TOTAL** | **50h** |

---

## ✅ CRITÉRIOS DE CONCLUSÃO

- [ ] Todas as telas funcionando
- [ ] Formulários com validação Zod
- [ ] React Query para todos os endpoints
- [ ] Componentes seguindo Design System
- [ ] Responsividade em todos os breakpoints
- [ ] Loading states implementados
- [ ] Error handling implementado
- [ ] Testes passando

---

## 🔗 DEPENDÊNCIAS

| Dependência | Status | Bloqueio |
|-------------|--------|----------|
| Sprint 3 (Painel/Projeções) | ❌ | **BLOQUEIA** |
| Design System | ✅ | — |
| API Client base | ✅ | — |

---

*Próximo Sprint: Sprint 5 — Testes + QA*
