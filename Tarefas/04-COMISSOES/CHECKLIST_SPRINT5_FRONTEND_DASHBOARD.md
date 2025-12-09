# ✅ CHECKLIST — SPRINT 5: FRONTEND DASHBOARD BARBEIRO

> **Status:** ❌ Não Iniciado  
> **Dependência:** Sprint 4 (Frontend Config + Fechamento)  
> **Esforço Estimado:** 18 horas  
> **Prioridade:** P1 — MVP

---

## 📊 RESUMO

```
░░░░░░░░░░░░░░░░░░░░ 0% COMPLETO
```

| Categoria | Completo | Pendente |
|-----------|:--------:|:--------:|
| Páginas | 0/2 | 2 |
| Componentes | 0/6 | 6 |
| Services | 0/1 | 1 |
| Hooks | 0/1 | 1 |
| Testes E2E | 0/2 | 2 |

---

## 1️⃣ SERVICES (API Integration)

### 1.1 Service: `advancesService` (Esforço: 1h)

- [ ] Criar `frontend/src/services/advancesService.ts`

```typescript
import { api } from '@/lib/api';

export interface Advance {
  id: string;
  unit_id?: string;
  professional_id: string;
  professional_name?: string;
  amount: string;
  request_date: string;
  reason?: string;
  status: 'PENDING' | 'APPROVED' | 'REJECTED' | 'DEDUCTED';
  approved_at?: string;
  approved_by?: string;
  rejected_at?: string;
  rejection_reason?: string;
  deducted_at?: string;
  created_at: string;
}

export interface CreateAdvanceRequest {
  unit_id?: string;
  professional_id: string;
  amount: string;
  request_date?: string;
  reason?: string;
}

export interface RejectAdvanceRequest {
  reason: string;
}

export const advancesService = {
  async list(params?: { professional_id?: string; status?: string }) {
    const response = await api.get<Advance[]>('/advances', { params });
    return response.data;
  },

  async getById(id: string) {
    const response = await api.get<Advance>(`/advances/${id}`);
    return response.data;
  },

  async create(data: CreateAdvanceRequest) {
    const response = await api.post<Advance>('/advances', data);
    return response.data;
  },

  async approve(id: string) {
    const response = await api.post<Advance>(`/advances/${id}/approve`);
    return response.data;
  },

  async reject(id: string, data: RejectAdvanceRequest) {
    const response = await api.post<Advance>(`/advances/${id}/reject`, data);
    return response.data;
  },

  async delete(id: string) {
    await api.delete(`/advances/${id}`);
  },
};
```

#### Checklist

- [ ] Interface Advance
- [ ] Interface CreateAdvanceRequest
- [ ] Interface RejectAdvanceRequest
- [ ] list
- [ ] getById
- [ ] create
- [ ] approve
- [ ] reject
- [ ] delete

---

## 2️⃣ HOOKS

### 2.1 Hook: `useAdvances` (Esforço: 0.5h)

- [ ] Criar `frontend/src/hooks/useAdvances.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { advancesService, CreateAdvanceRequest, RejectAdvanceRequest } from '@/services/advancesService';
import { toast } from 'sonner';

export function useAdvances(filters?: { professional_id?: string; status?: string }) {
  return useQuery({
    queryKey: ['advances', filters],
    queryFn: () => advancesService.list(filters),
  });
}

export function useAdvance(id: string) {
  return useQuery({
    queryKey: ['advance', id],
    queryFn: () => advancesService.getById(id),
    enabled: !!id,
  });
}

export function useCreateAdvance() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateAdvanceRequest) => advancesService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['advances'] });
      toast.success('Adiantamento solicitado com sucesso');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao solicitar adiantamento');
    },
  });
}

export function useApproveAdvance() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => advancesService.approve(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['advances'] });
      toast.success('Adiantamento aprovado');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao aprovar');
    },
  });
}

export function useRejectAdvance() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: RejectAdvanceRequest }) => 
      advancesService.reject(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['advances'] });
      toast.success('Adiantamento rejeitado');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao rejeitar');
    },
  });
}

export function useDeleteAdvance() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => advancesService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['advances'] });
      toast.success('Adiantamento removido');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao remover');
    },
  });
}
```

#### Checklist

- [ ] useAdvances
- [ ] useAdvance
- [ ] useCreateAdvance
- [ ] useApproveAdvance
- [ ] useRejectAdvance
- [ ] useDeleteAdvance

---

## 3️⃣ COMPONENTES — DASHBOARD BARBEIRO

### 3.1 Componente: `ComissaoCard` (Esforço: 1h)

- [ ] Criar `frontend/src/components/comissoes/ComissaoCard.tsx`

```tsx
'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { TrendingUp, TrendingDown, DollarSign } from 'lucide-react';
import { formatCurrency } from '@/lib/utils';
import { cn } from '@/lib/utils';

interface ComissaoCardProps {
  title: string;
  value: string;
  subtitle?: string;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  icon?: 'money' | 'up' | 'down';
  variant?: 'default' | 'success' | 'warning' | 'destructive';
}

export function ComissaoCard({
  title,
  value,
  subtitle,
  trend,
  icon = 'money',
  variant = 'default',
}: ComissaoCardProps) {
  const iconMap = {
    money: DollarSign,
    up: TrendingUp,
    down: TrendingDown,
  };
  const Icon = iconMap[icon];

  const variantStyles = {
    default: 'border-border',
    success: 'border-green-200 bg-green-50',
    warning: 'border-yellow-200 bg-yellow-50',
    destructive: 'border-red-200 bg-red-50',
  };

  return (
    <Card className={cn('transition-colors', variantStyles[variant])}>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{formatCurrency(value)}</div>
        {subtitle && (
          <p className="text-xs text-muted-foreground mt-1">{subtitle}</p>
        )}
        {trend && (
          <div className={cn(
            'flex items-center text-xs mt-2',
            trend.isPositive ? 'text-green-600' : 'text-red-600'
          )}>
            {trend.isPositive ? (
              <TrendingUp className="h-3 w-3 mr-1" />
            ) : (
              <TrendingDown className="h-3 w-3 mr-1" />
            )}
            {trend.value}% em relação ao mês anterior
          </div>
        )}
      </CardContent>
    </Card>
  );
}
```

#### Checklist

- [ ] Card com título e valor
- [ ] Ícone configurável
- [ ] Variantes de cor
- [ ] Trend (comparação)
- [ ] Subtitle

---

### 3.2 Componente: `ComissaoChart` (Esforço: 2h)

- [ ] Criar `frontend/src/components/comissoes/ComissaoChart.tsx`

```tsx
'use client';

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { formatCurrency } from '@/lib/utils';

interface ChartData {
  date: string;
  value: number;
}

interface ComissaoChartProps {
  title: string;
  data: ChartData[];
  isLoading?: boolean;
}

export function ComissaoChart({ title, data, isLoading }: ComissaoChartProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
        </CardHeader>
        <CardContent className="h-[300px] flex items-center justify-center">
          <div className="animate-pulse bg-muted h-full w-full rounded" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="h-[300px]">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={data} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
            <defs>
              <linearGradient id="colorValue" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.3} />
                <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
            <XAxis
              dataKey="date"
              tick={{ fontSize: 12 }}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              tick={{ fontSize: 12 }}
              tickLine={false}
              axisLine={false}
              tickFormatter={(value) => `R$ ${value}`}
            />
            <Tooltip
              content={({ active, payload }) => {
                if (active && payload && payload.length) {
                  return (
                    <div className="bg-background border rounded-lg p-2 shadow-lg">
                      <p className="text-sm font-medium">
                        {formatCurrency(String(payload[0].value))}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {payload[0].payload.date}
                      </p>
                    </div>
                  );
                }
                return null;
              }}
            />
            <Area
              type="monotone"
              dataKey="value"
              stroke="hsl(var(--primary))"
              fillOpacity={1}
              fill="url(#colorValue)"
            />
          </AreaChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
```

#### Checklist

- [ ] Gráfico de área com Recharts
- [ ] Tooltip customizado
- [ ] Cores do tema
- [ ] Loading state
- [ ] Responsivo

---

### 3.3 Componente: `ExtratoList` (Esforço: 1.5h)

- [ ] Criar `frontend/src/components/comissoes/ExtratoList.tsx`

```tsx
'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { formatCurrency, formatDate } from '@/lib/utils';
import { Commission } from '@/services/commissionsService';

interface ExtratoListProps {
  title: string;
  commissions: Commission[];
  isLoading?: boolean;
  maxItems?: number;
}

export function ExtratoList({
  title,
  commissions,
  isLoading,
  maxItems = 10,
}: ExtratoListProps) {
  const displayedCommissions = commissions.slice(0, maxItems);

  const getStatusBadge = (status: string) => {
    const config: Record<string, { variant: 'default' | 'secondary' | 'success' | 'destructive'; label: string }> = {
      PENDENTE: { variant: 'secondary', label: 'Pendente' },
      PROCESSADO: { variant: 'default', label: 'Processado' },
      PAGO: { variant: 'success', label: 'Pago' },
      CANCELADO: { variant: 'destructive', label: 'Cancelado' },
    };
    const { variant, label } = config[status] || { variant: 'default', label: status };
    return <Badge variant={variant}>{label}</Badge>;
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="animate-pulse">
                <div className="h-4 bg-muted rounded w-3/4 mb-2" />
                <div className="h-3 bg-muted rounded w-1/2" />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[400px]">
          <div className="space-y-4">
            {displayedCommissions.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-4">
                Nenhuma comissão encontrada
              </p>
            ) : (
              displayedCommissions.map((commission) => (
                <div
                  key={commission.id}
                  className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors"
                >
                  <div className="space-y-1">
                    <p className="text-sm font-medium">{commission.description}</p>
                    <p className="text-xs text-muted-foreground">
                      {formatDate(commission.data_competencia)}
                    </p>
                  </div>
                  <div className="text-right space-y-1">
                    <p className="text-sm font-semibold text-green-600">
                      +{formatCurrency(commission.commission_value)}
                    </p>
                    {getStatusBadge(commission.status)}
                  </div>
                </div>
              ))
            )}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  );
}
```

#### Checklist

- [ ] Lista com ScrollArea
- [ ] Item com descrição, data e valor
- [ ] Badge de status
- [ ] Empty state
- [ ] Loading state
- [ ] Limite de itens

---

### 3.4 Componente: `AdiantamentoForm` (Esforço: 1h)

- [ ] Criar `frontend/src/components/comissoes/AdiantamentoForm.tsx`

```tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { CurrencyInput } from '@/components/ui/currency-input';

const formSchema = z.object({
  amount: z.string().min(1, 'Valor é obrigatório'),
  reason: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;

interface AdiantamentoFormProps {
  onSubmit: (data: FormValues) => void;
  isLoading?: boolean;
}

export function AdiantamentoForm({ onSubmit, isLoading }: AdiantamentoFormProps) {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      amount: '',
      reason: '',
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="amount"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Valor do Adiantamento *</FormLabel>
              <FormControl>
                <CurrencyInput {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="reason"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Motivo (opcional)</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="Descreva o motivo do adiantamento..."
                  {...field}
                />
              </FormControl>
            </FormItem>
          )}
        />

        <Button type="submit" disabled={isLoading} className="w-full">
          {isLoading ? 'Solicitando...' : 'Solicitar Adiantamento'}
        </Button>
      </form>
    </Form>
  );
}
```

#### Checklist

- [ ] Form com react-hook-form + zod
- [ ] CurrencyInput para valor
- [ ] Textarea para motivo
- [ ] Validação
- [ ] Loading state

---

### 3.5 Componente: `AdiantamentosList` (Esforço: 1.5h)

- [ ] Criar `frontend/src/components/comissoes/AdiantamentosList.tsx`

```tsx
'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { formatCurrency, formatDate } from '@/lib/utils';
import { Advance } from '@/services/advancesService';
import { Check, X, Trash2 } from 'lucide-react';

interface AdiantamentosListProps {
  advances: Advance[];
  isLoading?: boolean;
  isManager?: boolean;
  onApprove?: (id: string) => void;
  onReject?: (id: string) => void;
  onDelete?: (id: string) => void;
}

export function AdiantamentosList({
  advances,
  isLoading,
  isManager = false,
  onApprove,
  onReject,
  onDelete,
}: AdiantamentosListProps) {
  const getStatusBadge = (status: string) => {
    const config: Record<string, { variant: 'default' | 'secondary' | 'success' | 'destructive' | 'outline'; label: string }> = {
      PENDING: { variant: 'secondary', label: 'Pendente' },
      APPROVED: { variant: 'success', label: 'Aprovado' },
      REJECTED: { variant: 'destructive', label: 'Rejeitado' },
      DEDUCTED: { variant: 'outline', label: 'Deduzido' },
    };
    const { variant, label } = config[status] || { variant: 'default', label: status };
    return <Badge variant={variant}>{label}</Badge>;
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Adiantamentos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="animate-pulse space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-16 bg-muted rounded" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Adiantamentos</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {advances.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-4">
              Nenhum adiantamento encontrado
            </p>
          ) : (
            advances.map((advance) => (
              <div
                key={advance.id}
                className="flex items-center justify-between p-4 rounded-lg border"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <p className="font-medium">{formatCurrency(advance.amount)}</p>
                    {getStatusBadge(advance.status)}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    Solicitado em {formatDate(advance.request_date)}
                  </p>
                  {advance.reason && (
                    <p className="text-sm text-muted-foreground">
                      Motivo: {advance.reason}
                    </p>
                  )}
                  {advance.rejection_reason && (
                    <p className="text-sm text-red-600">
                      Motivo da rejeição: {advance.rejection_reason}
                    </p>
                  )}
                </div>

                <div className="flex gap-2">
                  {isManager && advance.status === 'PENDING' && (
                    <>
                      <Button
                        variant="outline"
                        size="icon"
                        onClick={() => onApprove?.(advance.id)}
                      >
                        <Check className="h-4 w-4 text-green-600" />
                      </Button>
                      <Button
                        variant="outline"
                        size="icon"
                        onClick={() => onReject?.(advance.id)}
                      >
                        <X className="h-4 w-4 text-red-600" />
                      </Button>
                    </>
                  )}
                  {advance.status === 'PENDING' && (
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => onDelete?.(advance.id)}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  )}
                </div>
              </div>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
}
```

#### Checklist

- [ ] Lista de adiantamentos
- [ ] Badge de status
- [ ] Botões de ação (aprovar/rejeitar para gestor)
- [ ] Motivo de rejeição
- [ ] Empty state
- [ ] Loading state

---

### 3.6 Componente: `MetaProgressBar` (Esforço: 1h)

- [ ] Criar `frontend/src/components/comissoes/MetaProgressBar.tsx`

```tsx
'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { formatCurrency } from '@/lib/utils';
import { Target } from 'lucide-react';

interface MetaProgressBarProps {
  current: number;
  goal: number;
  title?: string;
}

export function MetaProgressBar({
  current,
  goal,
  title = 'Meta do Mês',
}: MetaProgressBarProps) {
  const percentage = goal > 0 ? Math.min((current / goal) * 100, 100) : 0;
  const remaining = Math.max(goal - current, 0);

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Target className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex justify-between text-sm">
          <span>Atual: {formatCurrency(String(current))}</span>
          <span className="text-muted-foreground">
            Meta: {formatCurrency(String(goal))}
          </span>
        </div>
        <Progress value={percentage} className="h-3" />
        <div className="flex justify-between text-xs text-muted-foreground">
          <span>{percentage.toFixed(1)}% atingido</span>
          <span>Faltam {formatCurrency(String(remaining))}</span>
        </div>
      </CardContent>
    </Card>
  );
}
```

#### Checklist

- [ ] Progress bar
- [ ] Valor atual vs meta
- [ ] Percentual
- [ ] Valor restante

---

## 4️⃣ PÁGINAS

### 4.1 Página: `/barbeiro/painel` (Esforço: 3h)

- [ ] Criar `frontend/src/app/(authenticated)/barbeiro/painel/page.tsx`

```tsx
'use client';

import { useState } from 'react';
import { useSession } from 'next-auth/react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ComissaoCard } from '@/components/comissoes/ComissaoCard';
import { ComissaoChart } from '@/components/comissoes/ComissaoChart';
import { ExtratoList } from '@/components/comissoes/ExtratoList';
import { MetaProgressBar } from '@/components/comissoes/MetaProgressBar';
import { AdiantamentoForm } from '@/components/comissoes/AdiantamentoForm';
import { AdiantamentosList } from '@/components/comissoes/AdiantamentosList';
import { useCommissions, useCommissionSummary } from '@/hooks/useCommissions';
import { useAdvances, useCreateAdvance } from '@/hooks/useAdvances';
import { useProfessionalMeta } from '@/hooks/useMetas';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Plus, DollarSign, TrendingUp, Calendar } from 'lucide-react';

export default function BarbeiroPainelPage() {
  const { data: session } = useSession();
  const [isAdvanceFormOpen, setIsAdvanceFormOpen] = useState(false);
  
  const professionalId = session?.user?.professionalId;
  
  const { data: commissions, isLoading: commissionsLoading } = useCommissions({
    professional_id: professionalId,
  });
  
  const { data: summary, isLoading: summaryLoading } = useCommissionSummary();
  
  const { data: advances, isLoading: advancesLoading } = useAdvances({
    professional_id: professionalId,
  });
  
  const { data: meta } = useProfessionalMeta(professionalId);
  
  const createAdvance = useCreateAdvance();

  const handleCreateAdvance = async (data: { amount: string; reason?: string }) => {
    if (!professionalId) return;
    await createAdvance.mutateAsync({
      professional_id: professionalId,
      ...data,
    });
    setIsAdvanceFormOpen(false);
  };

  // Dados para o gráfico (últimos 7 dias)
  const chartData = commissions?.slice(0, 7).map((c) => ({
    date: new Date(c.data_competencia).toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: 'short',
    }),
    value: parseFloat(c.commission_value),
  })).reverse() || [];

  return (
    <div className="container py-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Meu Painel</h1>
          <p className="text-muted-foreground">
            Acompanhe suas comissões e metas
          </p>
        </div>
        <Dialog open={isAdvanceFormOpen} onOpenChange={setIsAdvanceFormOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Solicitar Adiantamento
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Solicitar Adiantamento</DialogTitle>
            </DialogHeader>
            <AdiantamentoForm
              onSubmit={handleCreateAdvance}
              isLoading={createAdvance.isPending}
            />
          </DialogContent>
        </Dialog>
      </div>

      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <ComissaoCard
          title="Comissões do Mês"
          value={summary?.total_pending || '0'}
          subtitle="Pendentes de fechamento"
          icon="money"
        />
        <ComissaoCard
          title="Já Processado"
          value={summary?.total_processed || '0'}
          subtitle="Aguardando pagamento"
          icon="up"
          variant="warning"
        />
        <ComissaoCard
          title="Total Pago"
          value={summary?.total_paid || '0'}
          subtitle="Neste mês"
          icon="up"
          variant="success"
        />
        <ComissaoCard
          title="Atendimentos"
          value={String(summary?.count_pending || 0)}
          subtitle="No período"
          icon="money"
        />
      </div>

      {/* Meta */}
      {meta && (
        <MetaProgressBar
          current={parseFloat(summary?.total_pending || '0')}
          goal={parseFloat(meta.meta_valor)}
        />
      )}

      <Tabs defaultValue="extrato" className="space-y-4">
        <TabsList>
          <TabsTrigger value="extrato">Extrato</TabsTrigger>
          <TabsTrigger value="grafico">Evolução</TabsTrigger>
          <TabsTrigger value="adiantamentos">Adiantamentos</TabsTrigger>
        </TabsList>

        <TabsContent value="extrato">
          <ExtratoList
            title="Últimos Atendimentos"
            commissions={commissions || []}
            isLoading={commissionsLoading}
          />
        </TabsContent>

        <TabsContent value="grafico">
          <ComissaoChart
            title="Evolução Semanal"
            data={chartData}
            isLoading={commissionsLoading}
          />
        </TabsContent>

        <TabsContent value="adiantamentos">
          <AdiantamentosList
            advances={advances || []}
            isLoading={advancesLoading}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
```

#### Checklist

- [ ] Cards de resumo (4 métricas)
- [ ] Progress bar de meta
- [ ] Tabs: Extrato, Gráfico, Adiantamentos
- [ ] Lista de atendimentos com comissão
- [ ] Gráfico de evolução
- [ ] Lista de adiantamentos
- [ ] Modal para solicitar adiantamento
- [ ] RBAC: só vê seus próprios dados
- [ ] Loading states

---

### 4.2 Página: `/financeiro/adiantamentos` (Esforço: 2h)

- [ ] Criar `frontend/src/app/(authenticated)/financeiro/adiantamentos/page.tsx`

```tsx
'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { AdiantamentosList } from '@/components/comissoes/AdiantamentosList';
import { useAdvances, useApproveAdvance, useRejectAdvance } from '@/hooks/useAdvances';
import { useProfessionals } from '@/hooks/useProfessionals';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';

export default function AdiantamentosPage() {
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [professionalFilter, setProfessionalFilter] = useState<string>('');
  const [rejectModalOpen, setRejectModalOpen] = useState(false);
  const [rejectingId, setRejectingId] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState('');

  const { data: advances, isLoading } = useAdvances({
    status: statusFilter || undefined,
    professional_id: professionalFilter || undefined,
  });
  
  const { data: professionals } = useProfessionals();
  const approveAdvance = useApproveAdvance();
  const rejectAdvance = useRejectAdvance();

  const handleApprove = async (id: string) => {
    await approveAdvance.mutateAsync(id);
  };

  const handleRejectClick = (id: string) => {
    setRejectingId(id);
    setRejectReason('');
    setRejectModalOpen(true);
  };

  const handleRejectConfirm = async () => {
    if (!rejectingId || !rejectReason) return;
    await rejectAdvance.mutateAsync({
      id: rejectingId,
      data: { reason: rejectReason },
    });
    setRejectModalOpen(false);
    setRejectingId(null);
  };

  return (
    <div className="container py-6 space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Gestão de Adiantamentos</h1>
        <p className="text-muted-foreground">
          Aprove ou rejeite solicitações de adiantamento
        </p>
      </div>

      {/* Filtros */}
      <div className="flex gap-4">
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">Todos</SelectItem>
            <SelectItem value="PENDING">Pendentes</SelectItem>
            <SelectItem value="APPROVED">Aprovados</SelectItem>
            <SelectItem value="REJECTED">Rejeitados</SelectItem>
            <SelectItem value="DEDUCTED">Deduzidos</SelectItem>
          </SelectContent>
        </Select>

        <Select value={professionalFilter} onValueChange={setProfessionalFilter}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Profissional" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">Todos</SelectItem>
            {professionals?.map((p) => (
              <SelectItem key={p.id} value={p.id}>
                {p.nome}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Lista */}
      <AdiantamentosList
        advances={advances || []}
        isLoading={isLoading}
        isManager={true}
        onApprove={handleApprove}
        onReject={handleRejectClick}
      />

      {/* Modal de Rejeição */}
      <Dialog open={rejectModalOpen} onOpenChange={setRejectModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Rejeitar Adiantamento</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label>Motivo da Rejeição *</Label>
              <Textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                placeholder="Informe o motivo..."
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => setRejectModalOpen(false)}>
                Cancelar
              </Button>
              <Button
                variant="destructive"
                onClick={handleRejectConfirm}
                disabled={!rejectReason || rejectAdvance.isPending}
              >
                Confirmar Rejeição
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
```

#### Checklist

- [ ] Filtros (status, profissional)
- [ ] Lista de adiantamentos
- [ ] Botão aprovar
- [ ] Modal de rejeição com motivo
- [ ] RBAC: apenas gestor

---

## 5️⃣ TESTES E2E

### 5.1 Teste: Dashboard do Barbeiro (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/dashboard-barbeiro.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Dashboard do Barbeiro', () => {
  test.beforeEach(async ({ page }) => {
    // Login como barbeiro
    await page.goto('/barbeiro/painel');
  });

  test('deve exibir cards de resumo', async ({ page }) => {
    await expect(page.locator('text=Comissões do Mês')).toBeVisible();
    await expect(page.locator('text=Total Pago')).toBeVisible();
  });

  test('deve exibir extrato de atendimentos', async ({ page }) => {
    await page.click('text=Extrato');
    await expect(page.locator('[data-testid="extrato-list"]')).toBeVisible();
  });

  test('deve solicitar adiantamento', async ({ page }) => {
    await page.click('text=Solicitar Adiantamento');
    await page.fill('[name=amount]', '100');
    await page.fill('[name=reason]', 'Necessidade pessoal');
    await page.click('text=Solicitar Adiantamento');
    await expect(page.locator('text=Adiantamento solicitado')).toBeVisible();
  });

  test('barbeiro não deve ver dados de outros', async ({ page }) => {
    // Verificar que só vê suas próprias comissões
  });
});
```

#### Checklist

- [ ] Exibir cards de resumo
- [ ] Exibir extrato
- [ ] Exibir gráfico
- [ ] Solicitar adiantamento
- [ ] RBAC: só vê seus dados

---

### 5.2 Teste: Gestão de Adiantamentos (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/adiantamentos.spec.ts`

#### Checklist

- [ ] Listar adiantamentos
- [ ] Aprovar adiantamento
- [ ] Rejeitar com motivo
- [ ] Filtros funcionando

---

## 📝 NOTAS

### Próximos Passos

Após completar esta sprint:
1. Iniciar Sprint 6 (Testes E2E + QA)
2. Checklist: `CHECKLIST_SPRINT6_TESTES.md`

### Arquivos Criados

| Arquivo | Status |
|---------|--------|
| `services/advancesService.ts` | ❌ |
| `hooks/useAdvances.ts` | ❌ |
| `components/comissoes/ComissaoCard.tsx` | ❌ |
| `components/comissoes/ComissaoChart.tsx` | ❌ |
| `components/comissoes/ExtratoList.tsx` | ❌ |
| `components/comissoes/MetaProgressBar.tsx` | ❌ |
| `components/comissoes/AdiantamentoForm.tsx` | ❌ |
| `components/comissoes/AdiantamentosList.tsx` | ❌ |
| `app/.../barbeiro/painel/page.tsx` | ❌ |
| `app/.../financeiro/adiantamentos/page.tsx` | ❌ |
| `tests/comissoes/*.spec.ts` | ❌ |

---

*Checklist criado em: 05/12/2025*
