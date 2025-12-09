# ✅ CHECKLIST — SPRINT 4: FRONTEND CONFIG + FECHAMENTO

> **Status:** ❌ Não Iniciado  
> **Dependência:** Sprint 3 (Handlers + API)  
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
| Services | 0/3 | 3 |
| Hooks | 0/3 | 3 |
| Testes E2E | 0/2 | 2 |

---

## 1️⃣ SERVICES (API Integration)

### 1.1 Service: `commissionRulesService` (Esforço: 1h)

- [ ] Criar `frontend/src/services/commissionRulesService.ts`

```typescript
import { api } from '@/lib/api';

export interface CommissionRule {
  id: string;
  unit_id?: string;
  professional_id?: string;
  service_id?: string;
  type: 'PERCENTAGE' | 'FIXED' | 'HYBRID' | 'PROGRESSIVE';
  value: string;
  fixed_value?: string;
  tiers?: ProgressiveTier[];
  priority: number;
  active: boolean;
  created_at: string;
}

export interface ProgressiveTier {
  min: number;
  max?: number;
  pct: number;
}

export interface CreateCommissionRuleRequest {
  unit_id?: string;
  professional_id?: string;
  service_id?: string;
  type: string;
  value: string;
  fixed_value?: string;
  tiers?: ProgressiveTier[];
  priority: number;
  active: boolean;
}

export const commissionRulesService = {
  async list(params?: { unit_id?: string; professional_id?: string; service_id?: string; active?: boolean }) {
    const response = await api.get<CommissionRule[]>('/commission-rules', { params });
    return response.data;
  },

  async getById(id: string) {
    const response = await api.get<CommissionRule>(`/commission-rules/${id}`);
    return response.data;
  },

  async create(data: CreateCommissionRuleRequest) {
    const response = await api.post<CommissionRule>('/commission-rules', data);
    return response.data;
  },

  async update(id: string, data: Partial<CreateCommissionRuleRequest>) {
    const response = await api.put<CommissionRule>(`/commission-rules/${id}`, data);
    return response.data;
  },

  async toggle(id: string) {
    const response = await api.patch<CommissionRule>(`/commission-rules/${id}/toggle`);
    return response.data;
  },

  async delete(id: string) {
    await api.delete(`/commission-rules/${id}`);
  },
};
```

#### Checklist

- [ ] Interface CommissionRule
- [ ] Interface CreateCommissionRuleRequest
- [ ] list (com filtros)
- [ ] getById
- [ ] create
- [ ] update
- [ ] toggle
- [ ] delete

---

### 1.2 Service: `commissionPeriodsService` (Esforço: 1h)

- [ ] Criar `frontend/src/services/commissionPeriodsService.ts`

```typescript
import { api } from '@/lib/api';

export interface CommissionPeriod {
  id: string;
  unit_id?: string;
  professional_id: string;
  professional_name?: string;
  start_date: string;
  end_date: string;
  total_services: string;
  total_products: string;
  total_commission: string;
  total_bonus: string;
  total_deductions: string;
  net_value: string;
  qty_services: number;
  qty_products: number;
  status: 'DRAFT' | 'CLOSED' | 'PAID';
  bill_id?: string;
  notes?: string;
  closed_at?: string;
}

export interface PeriodPreviewRequest {
  unit_id?: string;
  professional_id: string;
  start_date: string;
  end_date: string;
}

export interface PeriodPreview {
  professional_id: string;
  professional_name: string;
  total_commission: string;
  total_bonus: string;
  total_deductions: string;
  net_value: string;
  items: CommissionItem[];
}

export interface CommissionItem {
  id: string;
  date: string;
  description: string;
  base_value: string;
  commission_value: string;
}

export interface ClosePeriodRequest {
  notes?: string;
}

export const commissionPeriodsService = {
  async list(params?: { professional_id?: string; unit_id?: string; status?: string; start_date?: string; end_date?: string }) {
    const response = await api.get<CommissionPeriod[]>('/commission-periods', { params });
    return response.data;
  },

  async getById(id: string) {
    const response = await api.get<CommissionPeriod>(`/commission-periods/${id}`);
    return response.data;
  },

  async generatePreview(data: PeriodPreviewRequest) {
    const response = await api.post<PeriodPreview>('/commission-periods/preview', data);
    return response.data;
  },

  async create(data: PeriodPreviewRequest) {
    const response = await api.post<CommissionPeriod>('/commission-periods', data);
    return response.data;
  },

  async update(id: string, data: Partial<CommissionPeriod>) {
    const response = await api.put<CommissionPeriod>(`/commission-periods/${id}`, data);
    return response.data;
  },

  async close(id: string, data?: ClosePeriodRequest) {
    const response = await api.post<CommissionPeriod>(`/commission-periods/${id}/close`, data);
    return response.data;
  },

  async delete(id: string) {
    await api.delete(`/commission-periods/${id}`);
  },
};
```

#### Checklist

- [ ] Interface CommissionPeriod
- [ ] Interface PeriodPreview
- [ ] Interface CommissionItem
- [ ] list
- [ ] getById
- [ ] generatePreview
- [ ] create
- [ ] update
- [ ] close
- [ ] delete

---

### 1.3 Service: `commissionsService` (Esforço: 0.5h)

- [ ] Criar `frontend/src/services/commissionsService.ts`

```typescript
import { api } from '@/lib/api';

export interface Commission {
  id: string;
  professional_id: string;
  professional_name?: string;
  command_item_id?: string;
  description: string;
  base_value: string;
  commission_value: string;
  status: 'PENDENTE' | 'PROCESSADO' | 'PAGO' | 'CANCELADO';
  data_competencia: string;
}

export interface CommissionSummary {
  total_pending: string;
  total_processed: string;
  total_paid: string;
  count_pending: number;
  count_processed: number;
}

export const commissionsService = {
  async list(params?: { professional_id?: string; status?: string; start_date?: string; end_date?: string }) {
    const response = await api.get<Commission[]>('/commissions', { params });
    return response.data;
  },

  async getSummary(params?: { start_date?: string; end_date?: string }) {
    const response = await api.get<CommissionSummary>('/commissions/summary', { params });
    return response.data;
  },

  async listByProfessional(professionalId: string, params?: { start_date?: string; end_date?: string }) {
    const response = await api.get<Commission[]>(`/professionals/${professionalId}/commissions`, { params });
    return response.data;
  },
};
```

#### Checklist

- [ ] Interface Commission
- [ ] Interface CommissionSummary
- [ ] list
- [ ] getSummary
- [ ] listByProfessional

---

## 2️⃣ HOOKS

### 2.1 Hook: `useCommissionRules` (Esforço: 0.5h)

- [ ] Criar `frontend/src/hooks/useCommissionRules.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { commissionRulesService, CreateCommissionRuleRequest } from '@/services/commissionRulesService';
import { toast } from 'sonner';

export function useCommissionRules(filters?: { unit_id?: string; professional_id?: string; service_id?: string; active?: boolean }) {
  return useQuery({
    queryKey: ['commission-rules', filters],
    queryFn: () => commissionRulesService.list(filters),
  });
}

export function useCommissionRule(id: string) {
  return useQuery({
    queryKey: ['commission-rule', id],
    queryFn: () => commissionRulesService.getById(id),
    enabled: !!id,
  });
}

export function useCreateCommissionRule() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateCommissionRuleRequest) => commissionRulesService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commission-rules'] });
      toast.success('Regra de comissão criada com sucesso');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao criar regra');
    },
  });
}

export function useUpdateCommissionRule() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateCommissionRuleRequest> }) => 
      commissionRulesService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commission-rules'] });
      toast.success('Regra atualizada com sucesso');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao atualizar regra');
    },
  });
}

export function useToggleCommissionRule() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => commissionRulesService.toggle(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commission-rules'] });
    },
  });
}

export function useDeleteCommissionRule() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => commissionRulesService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['commission-rules'] });
      toast.success('Regra removida com sucesso');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Erro ao remover regra');
    },
  });
}
```

#### Checklist

- [ ] useCommissionRules
- [ ] useCommissionRule
- [ ] useCreateCommissionRule
- [ ] useUpdateCommissionRule
- [ ] useToggleCommissionRule
- [ ] useDeleteCommissionRule

---

### 2.2 Hook: `useCommissionPeriods` (Esforço: 0.5h)

- [ ] Criar `frontend/src/hooks/useCommissionPeriods.ts`

#### Checklist

- [ ] useCommissionPeriods
- [ ] useCommissionPeriod
- [ ] useGeneratePreview
- [ ] useCreatePeriod
- [ ] useClosePeriod
- [ ] useDeletePeriod

---

### 2.3 Hook: `useCommissions` (Esforço: 0.5h)

- [ ] Criar `frontend/src/hooks/useCommissions.ts`

#### Checklist

- [ ] useCommissions
- [ ] useCommissionSummary
- [ ] useCommissionsByProfessional

---

## 3️⃣ COMPONENTES

### 3.1 Componente: `RegrasComissaoForm` (Esforço: 2h)

- [ ] Criar `frontend/src/components/comissoes/RegrasComissaoForm.tsx`

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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { CurrencyInput } from '@/components/ui/currency-input';

const formSchema = z.object({
  unit_id: z.string().optional(),
  professional_id: z.string().optional(),
  service_id: z.string().optional(),
  type: z.enum(['PERCENTAGE', 'FIXED', 'HYBRID', 'PROGRESSIVE']),
  value: z.string().min(1, 'Valor é obrigatório'),
  fixed_value: z.string().optional(),
  priority: z.number().min(0),
  active: z.boolean(),
});

type FormValues = z.infer<typeof formSchema>;

interface RegrasComissaoFormProps {
  defaultValues?: Partial<FormValues>;
  onSubmit: (data: FormValues) => void;
  isLoading?: boolean;
  professionals?: { id: string; name: string }[];
  services?: { id: string; name: string }[];
  units?: { id: string; name: string }[];
}

export function RegrasComissaoForm({
  defaultValues,
  onSubmit,
  isLoading,
  professionals = [],
  services = [],
  units = [],
}: RegrasComissaoFormProps) {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: 'PERCENTAGE',
      value: '',
      priority: 0,
      active: true,
      ...defaultValues,
    },
  });

  const selectedType = form.watch('type');

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {/* Escopo */}
        <div className="grid grid-cols-3 gap-4">
          <FormField
            control={form.control}
            name="unit_id"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Unidade (opcional)</FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Todas" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="">Todas</SelectItem>
                    {units.map((unit) => (
                      <SelectItem key={unit.id} value={unit.id}>
                        {unit.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
          {/* Repetir para professional_id e service_id */}
        </div>

        {/* Tipo */}
        <FormField
          control={form.control}
          name="type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Tipo de Comissão</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="PERCENTAGE">Percentual (%)</SelectItem>
                  <SelectItem value="FIXED">Valor Fixo (R$)</SelectItem>
                  <SelectItem value="HYBRID">Híbrido (Fixo + %)</SelectItem>
                  <SelectItem value="PROGRESSIVE">Progressivo (Faixas)</SelectItem>
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />

        {/* Valor */}
        <div className="grid grid-cols-2 gap-4">
          <FormField
            control={form.control}
            name="value"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  {selectedType === 'PERCENTAGE' ? 'Percentual (%)' : 'Valor (R$)'}
                </FormLabel>
                <FormControl>
                  {selectedType === 'PERCENTAGE' ? (
                    <Input type="number" min="0" max="100" step="0.01" {...field} />
                  ) : (
                    <CurrencyInput {...field} />
                  )}
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          {selectedType === 'HYBRID' && (
            <FormField
              control={form.control}
              name="fixed_value"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Valor Fixo (R$)</FormLabel>
                  <FormControl>
                    <CurrencyInput {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
          )}
        </div>

        {/* Prioridade e Status */}
        <div className="flex items-center justify-between">
          <FormField
            control={form.control}
            name="priority"
            render={({ field }) => (
              <FormItem className="flex items-center gap-2">
                <FormLabel>Prioridade</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    min="0"
                    className="w-20"
                    {...field}
                    onChange={(e) => field.onChange(parseInt(e.target.value))}
                  />
                </FormControl>
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="active"
            render={({ field }) => (
              <FormItem className="flex items-center gap-2">
                <FormLabel>Ativa</FormLabel>
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
              </FormItem>
            )}
          />
        </div>

        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Salvando...' : 'Salvar Regra'}
        </Button>
      </form>
    </Form>
  );
}
```

#### Checklist

- [ ] Form com react-hook-form + zod
- [ ] Select de tipo (PERCENTAGE, FIXED, HYBRID, PROGRESSIVE)
- [ ] Campos condicionais (fixed_value para HYBRID)
- [ ] Select de profissional/serviço/unidade (opcional)
- [ ] Prioridade
- [ ] Switch de ativo
- [ ] Validação

---

### 3.2 Componente: `RegrasComissaoTable` (Esforço: 1.5h)

- [ ] Criar `frontend/src/components/comissoes/RegrasComissaoTable.tsx`

#### Checklist

- [ ] Tabela com shadcn/ui DataTable
- [ ] Colunas: Escopo, Tipo, Valor, Prioridade, Status, Ações
- [ ] Formatação de valores
- [ ] Toggle de status
- [ ] Botões Editar/Excluir
- [ ] Empty state

---

### 3.3 Componente: `FechamentoTable` (Esforço: 2h)

- [ ] Criar `frontend/src/components/comissoes/FechamentoTable.tsx`

```tsx
'use client';

import { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatCurrency, formatDate } from '@/lib/utils';
import { ChevronDown, ChevronRight, Eye, Lock, Trash2 } from 'lucide-react';
import { CommissionPeriod } from '@/services/commissionPeriodsService';

interface FechamentoTableProps {
  periods: CommissionPeriod[];
  onView: (period: CommissionPeriod) => void;
  onClose: (period: CommissionPeriod) => void;
  onDelete: (period: CommissionPeriod) => void;
  isLoading?: boolean;
}

export function FechamentoTable({
  periods,
  onView,
  onClose,
  onDelete,
  isLoading,
}: FechamentoTableProps) {
  const [expanded, setExpanded] = useState<string | null>(null);

  const getStatusBadge = (status: string) => {
    const variants: Record<string, 'default' | 'secondary' | 'success' | 'destructive'> = {
      DRAFT: 'secondary',
      CLOSED: 'default',
      PAID: 'success',
    };
    const labels: Record<string, string> = {
      DRAFT: 'Rascunho',
      CLOSED: 'Fechado',
      PAID: 'Pago',
    };
    return <Badge variant={variants[status]}>{labels[status]}</Badge>;
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-8"></TableHead>
          <TableHead>Profissional</TableHead>
          <TableHead>Período</TableHead>
          <TableHead className="text-right">Comissão</TableHead>
          <TableHead className="text-right">Bônus</TableHead>
          <TableHead className="text-right">Deduções</TableHead>
          <TableHead className="text-right">Líquido</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Ações</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {periods.map((period) => (
          <>
            <TableRow key={period.id}>
              <TableCell>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setExpanded(expanded === period.id ? null : period.id)}
                >
                  {expanded === period.id ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </Button>
              </TableCell>
              <TableCell className="font-medium">{period.professional_name}</TableCell>
              <TableCell>
                {formatDate(period.start_date)} - {formatDate(period.end_date)}
              </TableCell>
              <TableCell className="text-right">{formatCurrency(period.total_commission)}</TableCell>
              <TableCell className="text-right text-green-600">
                +{formatCurrency(period.total_bonus)}
              </TableCell>
              <TableCell className="text-right text-red-600">
                -{formatCurrency(period.total_deductions)}
              </TableCell>
              <TableCell className="text-right font-semibold">
                {formatCurrency(period.net_value)}
              </TableCell>
              <TableCell>{getStatusBadge(period.status)}</TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-2">
                  <Button variant="ghost" size="icon" onClick={() => onView(period)}>
                    <Eye className="h-4 w-4" />
                  </Button>
                  {period.status === 'DRAFT' && (
                    <>
                      <Button variant="ghost" size="icon" onClick={() => onClose(period)}>
                        <Lock className="h-4 w-4" />
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => onDelete(period)}>
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </>
                  )}
                </div>
              </TableCell>
            </TableRow>
            {expanded === period.id && (
              <TableRow>
                <TableCell colSpan={9} className="bg-muted/50 p-4">
                  {/* Detalhes expandidos: lista de itens */}
                  <p className="text-sm text-muted-foreground">
                    {period.qty_services} serviços • {period.qty_products} produtos
                  </p>
                </TableCell>
              </TableRow>
            )}
          </>
        ))}
      </TableBody>
    </Table>
  );
}
```

#### Checklist

- [ ] Tabela com dados principais
- [ ] Expansão para ver detalhes
- [ ] Status com Badge
- [ ] Formatação de valores (currency)
- [ ] Ações condicionais (só DRAFT pode fechar/deletar)
- [ ] Loading state

---

### 3.4 Componente: `PreviewModal` (Esforço: 1.5h)

- [ ] Criar `frontend/src/components/comissoes/PreviewModal.tsx`

#### Checklist

- [ ] Modal com shadcn/ui Dialog
- [ ] Resumo: Comissão, Bônus, Deduções, Líquido
- [ ] Lista de itens que compõem o período
- [ ] Campo de observações
- [ ] Botões: Cancelar, Fechar Período
- [ ] Confirmação antes de fechar

---

### 3.5 Componente: `PeriodoFilterForm` (Esforço: 1h)

- [ ] Criar `frontend/src/components/comissoes/PeriodoFilterForm.tsx`

```tsx
'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { DatePicker } from '@/components/ui/date-picker';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';

interface PeriodoFilterFormProps {
  professionals: { id: string; name: string }[];
  units: { id: string; name: string }[];
  onGeneratePreview: (filters: {
    professional_id: string;
    unit_id?: string;
    start_date: string;
    end_date: string;
  }) => void;
  isLoading?: boolean;
}

export function PeriodoFilterForm({
  professionals,
  units,
  onGeneratePreview,
  isLoading,
}: PeriodoFilterFormProps) {
  const [professionalId, setProfessionalId] = useState<string>('');
  const [unitId, setUnitId] = useState<string>('');
  const [startDate, setStartDate] = useState<Date>();
  const [endDate, setEndDate] = useState<Date>();

  const handleSubmit = () => {
    if (!professionalId || !startDate || !endDate) return;

    onGeneratePreview({
      professional_id: professionalId,
      unit_id: unitId || undefined,
      start_date: startDate.toISOString().split('T')[0],
      end_date: endDate.toISOString().split('T')[0],
    });
  };

  return (
    <div className="flex flex-wrap items-end gap-4 p-4 border rounded-lg bg-muted/30">
      <div className="space-y-2">
        <Label>Profissional *</Label>
        <Select value={professionalId} onValueChange={setProfessionalId}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Selecione" />
          </SelectTrigger>
          <SelectContent>
            {professionals.map((p) => (
              <SelectItem key={p.id} value={p.id}>
                {p.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label>Unidade</Label>
        <Select value={unitId} onValueChange={setUnitId}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Todas" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">Todas</SelectItem>
            {units.map((u) => (
              <SelectItem key={u.id} value={u.id}>
                {u.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label>Data Início *</Label>
        <DatePicker date={startDate} onSelect={setStartDate} />
      </div>

      <div className="space-y-2">
        <Label>Data Fim *</Label>
        <DatePicker date={endDate} onSelect={setEndDate} />
      </div>

      <Button
        onClick={handleSubmit}
        disabled={!professionalId || !startDate || !endDate || isLoading}
      >
        {isLoading ? 'Gerando...' : 'Gerar Prévia'}
      </Button>
    </div>
  );
}
```

#### Checklist

- [ ] Select de Profissional (obrigatório)
- [ ] Select de Unidade (opcional)
- [ ] DatePicker início/fim
- [ ] Validação de campos obrigatórios
- [ ] Botão Gerar Prévia

---

### 3.6 Componente: `ConfigComissaoTenantForm` (Esforço: 1h)

- [ ] Criar `frontend/src/components/comissoes/ConfigComissaoTenantForm.tsx`

Configurações globais do tenant:
- Base de cálculo (GROSS_TOTAL, TABLE_PRICE, NET_VALUE)
- % padrão se não houver regra
- Lei Salão Parceiro (flag)

#### Checklist

- [ ] Select de base de cálculo
- [ ] Input de % padrão
- [ ] Switch Lei Salão Parceiro
- [ ] Botão Salvar

---

## 4️⃣ PÁGINAS

### 4.1 Página: `/admin/comissoes/config` (Esforço: 2h)

- [ ] Criar `frontend/src/app/(authenticated)/admin/comissoes/config/page.tsx`

```tsx
'use client';

import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';
import { RegrasComissaoTable } from '@/components/comissoes/RegrasComissaoTable';
import { RegrasComissaoForm } from '@/components/comissoes/RegrasComissaoForm';
import { ConfigComissaoTenantForm } from '@/components/comissoes/ConfigComissaoTenantForm';
import { useCommissionRules, useCreateCommissionRule } from '@/hooks/useCommissionRules';
import { useProfessionals } from '@/hooks/useProfessionals';
import { useServices } from '@/hooks/useServices';
import { useUnits } from '@/hooks/useUnits';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

export default function ConfigComissoesPage() {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingRule, setEditingRule] = useState<string | null>(null);

  const { data: rules, isLoading } = useCommissionRules();
  const { data: professionals } = useProfessionals();
  const { data: services } = useServices();
  const { data: units } = useUnits();
  const createRule = useCreateCommissionRule();

  const handleCreateRule = async (data: any) => {
    await createRule.mutateAsync(data);
    setIsFormOpen(false);
  };

  return (
    <div className="container py-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Configuração de Comissões</h1>
          <p className="text-muted-foreground">
            Defina as regras de comissão para profissionais e serviços
          </p>
        </div>
      </div>

      <Tabs defaultValue="regras" className="space-y-4">
        <TabsList>
          <TabsTrigger value="regras">Regras de Comissão</TabsTrigger>
          <TabsTrigger value="config">Configurações Gerais</TabsTrigger>
        </TabsList>

        <TabsContent value="regras" className="space-y-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Regras de Comissão</CardTitle>
                <CardDescription>
                  Regras específicas por profissional, serviço ou unidade
                </CardDescription>
              </div>
              <Button onClick={() => setIsFormOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Nova Regra
              </Button>
            </CardHeader>
            <CardContent>
              <RegrasComissaoTable
                rules={rules || []}
                isLoading={isLoading}
                onEdit={(id) => setEditingRule(id)}
                onToggle={() => {}}
                onDelete={() => {}}
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="config">
          <Card>
            <CardHeader>
              <CardTitle>Configurações Gerais</CardTitle>
              <CardDescription>
                Configurações padrão aplicadas quando não há regra específica
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ConfigComissaoTenantForm />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Nova Regra de Comissão</DialogTitle>
          </DialogHeader>
          <RegrasComissaoForm
            onSubmit={handleCreateRule}
            isLoading={createRule.isPending}
            professionals={professionals?.map((p) => ({ id: p.id, name: p.nome })) || []}
            services={services?.map((s) => ({ id: s.id, name: s.nome })) || []}
            units={units?.map((u) => ({ id: u.id, name: u.nome })) || []}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
}
```

#### Checklist

- [ ] Tabs: Regras / Configurações Gerais
- [ ] Lista de regras com tabela
- [ ] Modal para criar/editar regra
- [ ] Configurações globais do tenant
- [ ] Loading states
- [ ] Breadcrumb

---

### 4.2 Página: `/financeiro/comissoes` (Esforço: 3h)

- [ ] Criar `frontend/src/app/(authenticated)/financeiro/comissoes/page.tsx`

#### Checklist

- [ ] Filtros (Profissional, Unidade, Período, Status)
- [ ] Botão "Gerar Prévia"
- [ ] Tabela de períodos
- [ ] Modal de prévia/fechamento
- [ ] Toast de sucesso/erro
- [ ] Empty state
- [ ] Loading states

---

## 5️⃣ TESTES E2E

### 5.1 Teste: Configuração de Regras (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/config-regras.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Configuração de Regras de Comissão', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/admin/comissoes/config');
  });

  test('deve criar nova regra percentual', async ({ page }) => {
    await page.click('text=Nova Regra');
    await page.selectOption('[name=type]', 'PERCENTAGE');
    await page.fill('[name=value]', '50');
    await page.click('text=Salvar Regra');
    await expect(page.locator('text=Regra de comissão criada')).toBeVisible();
  });

  test('deve listar regras existentes', async ({ page }) => {
    await expect(page.locator('table')).toBeVisible();
  });

  test('deve ativar/desativar regra', async ({ page }) => {
    // ...
  });
});
```

#### Checklist

- [ ] Criar regra
- [ ] Listar regras
- [ ] Editar regra
- [ ] Toggle ativo
- [ ] Deletar regra

---

### 5.2 Teste: Fechamento de Período (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/fechamento.spec.ts`

#### Checklist

- [ ] Gerar prévia
- [ ] Visualizar detalhes
- [ ] Fechar período
- [ ] Verificar conta a pagar gerada

---

## 📝 NOTAS

### Próximos Passos

Após completar esta sprint:
1. Iniciar Sprint 5 (Frontend Dashboard Barbeiro)
2. Checklist: `CHECKLIST_SPRINT5_FRONTEND_DASHBOARD.md`

### Arquivos Criados

| Arquivo | Status |
|---------|--------|
| `services/commissionRulesService.ts` | ❌ |
| `services/commissionPeriodsService.ts` | ❌ |
| `services/commissionsService.ts` | ❌ |
| `hooks/useCommissionRules.ts` | ❌ |
| `hooks/useCommissionPeriods.ts` | ❌ |
| `hooks/useCommissions.ts` | ❌ |
| `components/comissoes/RegrasComissaoForm.tsx` | ❌ |
| `components/comissoes/RegrasComissaoTable.tsx` | ❌ |
| `components/comissoes/FechamentoTable.tsx` | ❌ |
| `components/comissoes/PreviewModal.tsx` | ❌ |
| `components/comissoes/PeriodoFilterForm.tsx` | ❌ |
| `components/comissoes/ConfigComissaoTenantForm.tsx` | ❌ |
| `app/.../admin/comissoes/config/page.tsx` | ❌ |
| `app/.../financeiro/comissoes/page.tsx` | ❌ |
| `tests/comissoes/*.spec.ts` | ❌ |

---

*Checklist criado em: 05/12/2025*
