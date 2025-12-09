'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Textarea } from '@/components/ui/textarea';
import {
    useCloseCommissionPeriod,
    useCommissionPeriods,
    useCommissionPeriodSummary,
    useCreateCommissionPeriod,
    useDeleteCommissionPeriod,
    useMarkPeriodAsPaid,
} from '@/hooks/use-commissions';
import { useProfessionals } from '@/hooks/use-professionals';
import {
    CommissionPeriod,
    CommissionPeriodStatus,
    CreateCommissionPeriodRequest,
} from '@/types/commission';
import {
    CalendarDays,
    CheckCircle2,
    Clock,
    DollarSign,
    Eye,
    Lock,
    MoreHorizontal,
    Plus,
    Trash2,
    XCircle,
} from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: string | number) => {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return 'R$ 0,00';
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR').format(date);
};

const formatMonth = (monthStr: string) => {
  if (!monthStr) return '-';
  const [year, month] = monthStr.split('-');
  const date = new Date(parseInt(year), parseInt(month) - 1);
  return new Intl.DateTimeFormat('pt-BR', { 
    month: 'long', 
    year: 'numeric' 
  }).format(date);
};

const getStatusBadge = (status: CommissionPeriodStatus) => {
  const variants: Record<CommissionPeriodStatus, { variant: 'default' | 'secondary' | 'destructive' | 'outline'; label: string; icon: React.ReactNode }> = {
    [CommissionPeriodStatus.ABERTO]: { variant: 'outline', label: 'Aberto', icon: <Clock className="h-3 w-3 mr-1" /> },
    [CommissionPeriodStatus.PROCESSANDO]: { variant: 'secondary', label: 'Processando', icon: <Clock className="h-3 w-3 mr-1 animate-spin" /> },
    [CommissionPeriodStatus.FECHADO]: { variant: 'default', label: 'Fechado', icon: <Lock className="h-3 w-3 mr-1" /> },
    [CommissionPeriodStatus.PAGO]: { variant: 'default', label: 'Pago', icon: <CheckCircle2 className="h-3 w-3 mr-1" /> },
    [CommissionPeriodStatus.CANCELADO]: { variant: 'destructive', label: 'Cancelado', icon: <XCircle className="h-3 w-3 mr-1" /> },
  };
  const { variant, label, icon } = variants[status] || { variant: 'outline', label: status, icon: null };
  return (
    <Badge variant={variant} className="flex items-center">
      {icon}
      {label}
    </Badge>
  );
};

// =============================================================================
// FORM DIALOG
// =============================================================================

interface PeriodFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function PeriodFormDialog({ open, onOpenChange }: PeriodFormDialogProps) {
  const createMutation = useCreateCommissionPeriod();
  const { data: professionals } = useProfessionals();

  const now = new Date();
  const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
  const firstDay = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0];
  const lastDay = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().split('T')[0];

  const [formData, setFormData] = useState<CreateCommissionPeriodRequest>({
    professional_id: '',
    reference_month: currentMonth,
    period_start: firstDay,
    period_end: lastDay,
    notes: '',
  });

  const isLoading = createMutation.isPending;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.professional_id || !formData.reference_month) {
      toast.error('Selecione o profissional e o mês de referência');
      return;
    }

    try {
      await createMutation.mutateAsync({
        ...formData,
        period_start: new Date(formData.period_start).toISOString(),
        period_end: new Date(formData.period_end).toISOString(),
      });
      onOpenChange(false);
      setFormData({
        professional_id: '',
        reference_month: currentMonth,
        period_start: firstDay,
        period_end: lastDay,
        notes: '',
      });
    } catch {
      // Error handled by mutation
    }
  };

  const handleMonthChange = (month: string) => {
    const [year, m] = month.split('-');
    const first = new Date(parseInt(year), parseInt(m) - 1, 1);
    const last = new Date(parseInt(year), parseInt(m), 0);
    
    setFormData({
      ...formData,
      reference_month: month,
      period_start: first.toISOString().split('T')[0],
      period_end: last.toISOString().split('T')[0],
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Novo Período de Comissão</DialogTitle>
            <DialogDescription>
              Crie um período para consolidar as comissões de um profissional.
            </DialogDescription>
          </DialogHeader>
          
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="professional">Profissional *</Label>
              <Select
                value={formData.professional_id}
                onValueChange={(value) => setFormData({ ...formData, professional_id: value })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o profissional" />
                </SelectTrigger>
                <SelectContent>
                  {professionals?.data?.map((prof) => (
                    <SelectItem key={prof.id} value={prof.id}>
                      {prof.nome}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="reference_month">Mês de Referência *</Label>
              <Input
                id="reference_month"
                type="month"
                value={formData.reference_month}
                onChange={(e) => handleMonthChange(e.target.value)}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="period_start">Início do Período</Label>
                <Input
                  id="period_start"
                  type="date"
                  value={formData.period_start}
                  onChange={(e) => setFormData({ ...formData, period_start: e.target.value })}
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="period_end">Fim do Período</Label>
                <Input
                  id="period_end"
                  type="date"
                  value={formData.period_end}
                  onChange={(e) => setFormData({ ...formData, period_end: e.target.value })}
                />
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="notes">Observações</Label>
              <Textarea
                id="notes"
                value={formData.notes || ''}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                placeholder="Anotações sobre o período..."
                rows={2}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? 'Criando...' : 'Criar Período'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// DETAIL DIALOG
// =============================================================================

interface PeriodDetailDialogProps {
  period: CommissionPeriod;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function PeriodDetailDialog({ period, open, onOpenChange }: PeriodDetailDialogProps) {
  const { isLoading } = useCommissionPeriodSummary(period.id);
  const closeMutation = useCloseCommissionPeriod();
  const paidMutation = useMarkPeriodAsPaid();

  const canClose = period.status === CommissionPeriodStatus.ABERTO;
  const canMarkPaid = period.status === CommissionPeriodStatus.FECHADO;

  const handleClose = async () => {
    try {
      await closeMutation.mutateAsync({ id: period.id });
    } catch {
      // Error handled by mutation
    }
  };

  const handleMarkPaid = async () => {
    try {
      await paidMutation.mutateAsync(period.id);
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>Detalhes do Período</DialogTitle>
          <DialogDescription>
            {period.professional_name} — {formatMonth(period.reference_month)}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Status */}
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">Status:</span>
            {getStatusBadge(period.status as CommissionPeriodStatus)}
          </div>

          {/* Summary Cards */}
          <div className="grid grid-cols-2 gap-4">
            <Card>
              <CardContent className="pt-4">
                <div className="text-2xl font-bold text-green-600">
                  {isLoading ? <Skeleton className="h-8 w-24" /> : formatCurrency(period.total_net)}
                </div>
                <p className="text-sm text-muted-foreground">Valor Líquido</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4">
                <div className="text-2xl font-bold">
                  {isLoading ? <Skeleton className="h-8 w-24" /> : formatCurrency(period.total_commission)}
                </div>
                <p className="text-sm text-muted-foreground">Total Comissões</p>
              </CardContent>
            </Card>
          </div>

          {/* Details */}
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-muted-foreground">Faturamento Bruto:</span>
              <p className="font-medium">{formatCurrency(period.total_gross)}</p>
            </div>
            <div>
              <span className="text-muted-foreground">Adiantamentos:</span>
              <p className="font-medium text-red-600">- {formatCurrency(period.total_advances)}</p>
            </div>
            <div>
              <span className="text-muted-foreground">Ajustes:</span>
              <p className="font-medium">{formatCurrency(period.total_adjustments)}</p>
            </div>
            <div>
              <span className="text-muted-foreground">Qtd. Itens:</span>
              <p className="font-medium">{period.items_count}</p>
            </div>
            <div>
              <span className="text-muted-foreground">Período:</span>
              <p className="font-medium">
                {formatDate(period.period_start)} a {formatDate(period.period_end)}
              </p>
            </div>
            {period.notes && (
              <div className="col-span-2">
                <span className="text-muted-foreground">Observações:</span>
                <p className="font-medium">{period.notes}</p>
              </div>
            )}
          </div>
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          {canClose && (
            <Button
              variant="outline"
              onClick={handleClose}
              disabled={closeMutation.isPending}
            >
              <Lock className="mr-2 h-4 w-4" />
              {closeMutation.isPending ? 'Fechando...' : 'Fechar Período'}
            </Button>
          )}
          {canMarkPaid && (
            <Button onClick={handleMarkPaid} disabled={paidMutation.isPending}>
              <DollarSign className="mr-2 h-4 w-4" />
              {paidMutation.isPending ? 'Processando...' : 'Marcar como Pago'}
            </Button>
          )}
          <Button variant="ghost" onClick={() => onOpenChange(false)}>
            Fechar
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// DELETE DIALOG
// =============================================================================

interface DeleteDialogProps {
  period: CommissionPeriod;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function DeleteDialog({ period, open, onOpenChange }: DeleteDialogProps) {
  const deleteMutation = useDeleteCommissionPeriod();

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(period.id);
      onOpenChange(false);
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Confirmar Exclusão</DialogTitle>
          <DialogDescription>
            Tem certeza que deseja excluir o período de {period.professional_name} — {formatMonth(period.reference_month)}? 
            Esta ação não pode ser desfeita.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? 'Excluindo...' : 'Excluir'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// MAIN PAGE
// =============================================================================

export default function PeriodosComissaoPage() {
  const [statusFilter, setStatusFilter] = useState<string>('all');
  
  const { data, isLoading } = useCommissionPeriods({
    status: statusFilter === 'all' ? undefined : (statusFilter as CommissionPeriodStatus),
    limit: 50,
  });

  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [detailDialogOpen, setDetailDialogOpen] = useState(false);
  const [selectedPeriod, setSelectedPeriod] = useState<CommissionPeriod | undefined>();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deletingPeriod, setDeletingPeriod] = useState<CommissionPeriod | undefined>();

  const handleViewDetails = (period: CommissionPeriod) => {
    setSelectedPeriod(period);
    setDetailDialogOpen(true);
  };

  const handleDelete = (period: CommissionPeriod) => {
    setDeletingPeriod(period);
    setDeleteDialogOpen(true);
  };

  // Stats
  const stats = useMemo(() => {
    if (!data?.data) return { open: 0, closed: 0, paid: 0, total: 0 };
    
    return {
      open: data.data.filter(p => p.status === CommissionPeriodStatus.ABERTO).length,
      closed: data.data.filter(p => p.status === CommissionPeriodStatus.FECHADO).length,
      paid: data.data.filter(p => p.status === CommissionPeriodStatus.PAGO).length,
      total: data.data.reduce((acc, p) => acc + parseFloat(p.total_net || '0'), 0),
    };
  }, [data]);

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold">Períodos de Comissão</h2>
          <p className="text-sm text-muted-foreground">
            Gerencie os períodos de fechamento de comissões dos profissionais
          </p>
        </div>
        <Button onClick={() => setFormDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Novo Período
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold">{stats.open}</div>
            <p className="text-sm text-muted-foreground">Abertos</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold">{stats.closed}</div>
            <p className="text-sm text-muted-foreground">Fechados</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold text-green-600">{stats.paid}</div>
            <p className="text-sm text-muted-foreground">Pagos</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold">{formatCurrency(stats.total)}</div>
            <p className="text-sm text-muted-foreground">Total Líquido</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center gap-4">
            <div className="grid gap-2">
              <Label>Status</Label>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Todos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Todos</SelectItem>
                  <SelectItem value={CommissionPeriodStatus.ABERTO}>Abertos</SelectItem>
                  <SelectItem value={CommissionPeriodStatus.FECHADO}>Fechados</SelectItem>
                  <SelectItem value={CommissionPeriodStatus.PAGO}>Pagos</SelectItem>
                  <SelectItem value={CommissionPeriodStatus.CANCELADO}>Cancelados</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="pt-6">
          {isLoading ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Profissional</TableHead>
                  <TableHead>Mês</TableHead>
                  <TableHead>Período</TableHead>
                  <TableHead>Comissão</TableHead>
                  <TableHead>Líquido</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-8" /></TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : data?.data && data.data.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Profissional</TableHead>
                  <TableHead>Mês</TableHead>
                  <TableHead>Período</TableHead>
                  <TableHead>Comissão</TableHead>
                  <TableHead>Líquido</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((period) => (
                  <TableRow key={period.id}>
                    <TableCell className="font-medium">
                      {period.professional_name || '-'}
                    </TableCell>
                    <TableCell>{formatMonth(period.reference_month)}</TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(period.period_start)} - {formatDate(period.period_end)}
                    </TableCell>
                    <TableCell>{formatCurrency(period.total_commission)}</TableCell>
                    <TableCell className="font-semibold text-green-600">
                      {formatCurrency(period.total_net)}
                    </TableCell>
                    <TableCell>
                      {getStatusBadge(period.status as CommissionPeriodStatus)}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => handleViewDetails(period)}>
                            <Eye className="mr-2 h-4 w-4" />
                            Ver Detalhes
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            onClick={() => handleDelete(period)}
                            className="text-destructive"
                            disabled={period.status === CommissionPeriodStatus.PAGO}
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            Excluir
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <CalendarDays className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium">Nenhum período encontrado</h3>
              <p className="text-sm text-muted-foreground mb-4">
                Crie um período de comissão para começar
              </p>
              <Button onClick={() => setFormDialogOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Novo Período
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Dialogs */}
      <PeriodFormDialog
        open={formDialogOpen}
        onOpenChange={setFormDialogOpen}
      />

      {selectedPeriod && (
        <PeriodDetailDialog
          period={selectedPeriod}
          open={detailDialogOpen}
          onOpenChange={setDetailDialogOpen}
        />
      )}

      {deletingPeriod && (
        <DeleteDialog
          period={deletingPeriod}
          open={deleteDialogOpen}
          onOpenChange={(open) => {
            setDeleteDialogOpen(open);
            if (!open) setDeletingPeriod(undefined);
          }}
        />
      )}
    </div>
  );
}
