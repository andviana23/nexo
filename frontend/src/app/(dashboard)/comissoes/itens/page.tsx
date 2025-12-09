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
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
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
import {
    useCommissionItems,
    useDeleteCommissionItem,
} from '@/hooks/use-commissions';
import { useProfessionals } from '@/hooks/use-professionals';
import {
    CommissionItem,
    CommissionItemStatus,
    CommissionRuleType,
    CommissionSource,
} from '@/types/commission';
import {
    Calculator,
    CheckCircle2,
    Clock,
    MoreHorizontal,
    RefreshCw,
    Trash2,
    XCircle,
} from 'lucide-react';
import { useMemo, useState } from 'react';

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

const formatPercentage = (value: string | number) => {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return '0%';
  return `${num}%`;
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR').format(date);
};

const getStatusBadge = (status: CommissionItemStatus) => {
  const variants: Record<CommissionItemStatus, { variant: 'default' | 'secondary' | 'destructive' | 'outline'; label: string; icon: React.ReactNode }> = {
    [CommissionItemStatus.PENDENTE]: { variant: 'outline', label: 'Pendente', icon: <Clock className="h-3 w-3 mr-1" /> },
    [CommissionItemStatus.PROCESSADO]: { variant: 'secondary', label: 'Processado', icon: <RefreshCw className="h-3 w-3 mr-1" /> },
    [CommissionItemStatus.PAGO]: { variant: 'default', label: 'Pago', icon: <CheckCircle2 className="h-3 w-3 mr-1" /> },
    [CommissionItemStatus.CANCELADO]: { variant: 'destructive', label: 'Cancelado', icon: <XCircle className="h-3 w-3 mr-1" /> },
    [CommissionItemStatus.ESTORNADO]: { variant: 'destructive', label: 'Estornado', icon: <RefreshCw className="h-3 w-3 mr-1" /> },
  };
  const { variant, label, icon } = variants[status] || { variant: 'outline', label: status, icon: null };
  return (
    <Badge variant={variant} className="flex items-center">
      {icon}
      {label}
    </Badge>
  );
};

const getSourceBadge = (source: CommissionSource) => {
  const labels: Record<CommissionSource, string> = {
    [CommissionSource.SERVICO]: 'Serviço',
    [CommissionSource.PROFISSIONAL]: 'Profissional',
    [CommissionSource.REGRA]: 'Regra',
    [CommissionSource.MANUAL]: 'Manual',
  };
  return <Badge variant="outline">{labels[source] || source}</Badge>;
};

// =============================================================================
// DELETE DIALOG
// =============================================================================

interface DeleteDialogProps {
  item: CommissionItem;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function DeleteDialog({ item, open, onOpenChange }: DeleteDialogProps) {
  const deleteMutation = useDeleteCommissionItem();

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(item.id);
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
            Tem certeza que deseja excluir este item de comissão de {formatCurrency(item.commission_value)}? 
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

export default function ItensComissaoPage() {
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [professionalFilter, setProfessionalFilter] = useState<string>('all');
  
  const { data: professionals } = useProfessionals();
  const { data, isLoading } = useCommissionItems({
    status: statusFilter === 'all' ? undefined : (statusFilter as CommissionItemStatus),
    professional_id: professionalFilter === 'all' ? undefined : professionalFilter,
    limit: 100,
  });

  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deletingItem, setDeletingItem] = useState<CommissionItem | undefined>();

  const handleDelete = (item: CommissionItem) => {
    setDeletingItem(item);
    setDeleteDialogOpen(true);
  };

  // Stats
  const stats = useMemo(() => {
    if (!data?.data) return { total: 0, pending: 0, processed: 0, totalValue: 0 };
    
    const items = data.data;
    return {
      total: items.length,
      pending: items.filter(i => i.status === CommissionItemStatus.PENDENTE).length,
      processed: items.filter(i => i.status === CommissionItemStatus.PROCESSADO || i.status === CommissionItemStatus.PAGO).length,
      totalValue: items.reduce((acc, i) => acc + parseFloat(i.commission_value || '0'), 0),
    };
  }, [data]);

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold">Itens de Comissão</h2>
          <p className="text-sm text-muted-foreground">
            Visualize todos os itens de comissão gerados
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-sm text-muted-foreground">Total de Itens</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold text-yellow-600">{stats.pending}</div>
            <p className="text-sm text-muted-foreground">Pendentes</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold text-green-600">{stats.processed}</div>
            <p className="text-sm text-muted-foreground">Processados/Pagos</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="text-2xl font-bold">{formatCurrency(stats.totalValue)}</div>
            <p className="text-sm text-muted-foreground">Valor Total</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-wrap items-center gap-4">
            <div className="grid gap-2">
              <Label>Profissional</Label>
              <Select value={professionalFilter} onValueChange={setProfessionalFilter}>
                <SelectTrigger className="w-[200px]">
                  <SelectValue placeholder="Todos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Todos</SelectItem>
                  {professionals?.data?.map((prof) => (
                    <SelectItem key={prof.id} value={prof.id}>
                      {prof.nome}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label>Status</Label>
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="Todos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Todos</SelectItem>
                  <SelectItem value={CommissionItemStatus.PENDENTE}>Pendentes</SelectItem>
                  <SelectItem value={CommissionItemStatus.PROCESSADO}>Processados</SelectItem>
                  <SelectItem value={CommissionItemStatus.PAGO}>Pagos</SelectItem>
                  <SelectItem value={CommissionItemStatus.CANCELADO}>Cancelados</SelectItem>
                  <SelectItem value={CommissionItemStatus.ESTORNADO}>Estornados</SelectItem>
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
                  <TableHead>Serviço</TableHead>
                  <TableHead>Data</TableHead>
                  <TableHead>Bruto</TableHead>
                  <TableHead>Taxa</TableHead>
                  <TableHead>Comissão</TableHead>
                  <TableHead>Origem</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[60px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {Array.from({ length: 8 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-28" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-8" /></TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : data?.data && data.data.length > 0 ? (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Profissional</TableHead>
                    <TableHead>Serviço</TableHead>
                    <TableHead>Data</TableHead>
                    <TableHead>Bruto</TableHead>
                    <TableHead>Taxa</TableHead>
                    <TableHead>Comissão</TableHead>
                    <TableHead>Origem</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="w-[60px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {data.data.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-medium">
                        {item.professional_name || '-'}
                      </TableCell>
                      <TableCell className="max-w-[150px] truncate">
                        {item.service_name || item.description || '-'}
                      </TableCell>
                      <TableCell>{formatDate(item.reference_date)}</TableCell>
                      <TableCell>{formatCurrency(item.gross_value)}</TableCell>
                      <TableCell>
                        {item.commission_type === CommissionRuleType.PERCENTUAL
                          ? formatPercentage(item.commission_rate)
                          : formatCurrency(item.commission_rate)}
                      </TableCell>
                      <TableCell className="font-semibold text-green-600">
                        {formatCurrency(item.commission_value)}
                      </TableCell>
                      <TableCell>
                        {getSourceBadge(item.commission_source as CommissionSource)}
                      </TableCell>
                      <TableCell>
                        {getStatusBadge(item.status as CommissionItemStatus)}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() => handleDelete(item)}
                              className="text-destructive"
                              disabled={item.status === CommissionItemStatus.PAGO}
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
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <Calculator className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium">Nenhum item de comissão</h3>
              <p className="text-sm text-muted-foreground">
                Os itens de comissão são gerados automaticamente a partir dos atendimentos
              </p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete Dialog */}
      {deletingItem && (
        <DeleteDialog
          item={deletingItem}
          open={deleteDialogOpen}
          onOpenChange={(open) => {
            setDeleteDialogOpen(open);
            if (!open) setDeletingItem(undefined);
          }}
        />
      )}
    </div>
  );
}
