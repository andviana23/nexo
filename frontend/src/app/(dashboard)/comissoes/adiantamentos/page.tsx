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
    useAdvances,
    useApproveAdvance,
    useCancelAdvance,
    useCreateAdvance,
    useDeleteAdvance,
    usePendingAdvances,
    useRejectAdvance,
} from '@/hooks/use-commissions';
import { useProfessionals } from '@/hooks/use-professionals';
import { cn } from '@/lib/utils';
import {
    Advance,
    AdvanceStatus,
    CreateAdvanceRequest,
} from '@/types/commission';
import {
    AlertCircle,
    Ban,
    Check,
    CheckCircle2,
    Clock,
    HandCoins,
    MoreHorizontal,
    Plus,
    Trash2,
    X,
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

const getStatusBadge = (status: AdvanceStatus) => {
  const variants: Record<AdvanceStatus, { variant: 'default' | 'secondary' | 'destructive' | 'outline'; label: string; icon: React.ReactNode }> = {
    [AdvanceStatus.PENDING]: { variant: 'outline', label: 'Pendente', icon: <Clock className="h-3 w-3 mr-1" /> },
    [AdvanceStatus.APPROVED]: { variant: 'default', label: 'Aprovado', icon: <CheckCircle2 className="h-3 w-3 mr-1" /> },
    [AdvanceStatus.REJECTED]: { variant: 'destructive', label: 'Rejeitado', icon: <XCircle className="h-3 w-3 mr-1" /> },
    [AdvanceStatus.DEDUCTED]: { variant: 'secondary', label: 'Deduzido', icon: <Check className="h-3 w-3 mr-1" /> },
    [AdvanceStatus.CANCELLED]: { variant: 'destructive', label: 'Cancelado', icon: <Ban className="h-3 w-3 mr-1" /> },
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

interface AdvanceFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function AdvanceFormDialog({ open, onOpenChange }: AdvanceFormDialogProps) {
  const createMutation = useCreateAdvance();
  const { data: professionals } = useProfessionals();

  const [formData, setFormData] = useState<CreateAdvanceRequest>({
    professional_id: '',
    amount: '',
    reason: '',
  });

  const isLoading = createMutation.isPending;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.professional_id || !formData.amount) {
      toast.error('Preencha os campos obrigatórios');
      return;
    }

    try {
      await createMutation.mutateAsync(formData);
      onOpenChange(false);
      setFormData({
        professional_id: '',
        amount: '',
        reason: '',
      });
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[450px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Solicitar Adiantamento</DialogTitle>
            <DialogDescription>
              Registre uma solicitação de adiantamento para um profissional.
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
              <Label htmlFor="amount">Valor (R$) *</Label>
              <Input
                id="amount"
                type="number"
                step="0.01"
                min="0"
                value={formData.amount}
                onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                placeholder="0,00"
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="reason">Motivo</Label>
              <Textarea
                id="reason"
                value={formData.reason || ''}
                onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
                placeholder="Motivo do adiantamento..."
                rows={3}
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
              {isLoading ? 'Solicitando...' : 'Solicitar'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// REJECT DIALOG
// =============================================================================

interface RejectDialogProps {
  advance: Advance;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function RejectDialog({ advance, open, onOpenChange }: RejectDialogProps) {
  const rejectMutation = useRejectAdvance();
  const [reason, setReason] = useState('');

  const handleReject = async () => {
    if (!reason.trim()) {
      toast.error('Informe o motivo da rejeição');
      return;
    }

    try {
      await rejectMutation.mutateAsync({
        id: advance.id,
        data: { rejection_reason: reason },
      });
      onOpenChange(false);
      setReason('');
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Rejeitar Adiantamento</DialogTitle>
          <DialogDescription>
            Informe o motivo da rejeição do adiantamento de {formatCurrency(advance.amount)} para {advance.professional_name}.
          </DialogDescription>
        </DialogHeader>
        
        <div className="grid gap-2 py-4">
          <Label htmlFor="rejection_reason">Motivo da Rejeição *</Label>
          <Textarea
            id="rejection_reason"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Descreva o motivo..."
            rows={3}
          />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={handleReject}
            disabled={rejectMutation.isPending}
          >
            {rejectMutation.isPending ? 'Rejeitando...' : 'Rejeitar'}
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
  advance: Advance;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function DeleteDialog({ advance, open, onOpenChange }: DeleteDialogProps) {
  const deleteMutation = useDeleteAdvance();

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(advance.id);
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
            Tem certeza que deseja excluir o adiantamento de {formatCurrency(advance.amount)} para {advance.professional_name}? 
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

export default function AdiantamentosPage() {
  const [statusFilter, setStatusFilter] = useState<string>('all');
  
  const { data, isLoading } = useAdvances({
    status: statusFilter === 'all' ? undefined : (statusFilter as AdvanceStatus),
    limit: 50,
  });
  
  const { data: pendingData } = usePendingAdvances();

  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [rejectDialogOpen, setRejectDialogOpen] = useState(false);
  const [rejectingAdvance, setRejectingAdvance] = useState<Advance | undefined>();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deletingAdvance, setDeletingAdvance] = useState<Advance | undefined>();

  const approveMutation = useApproveAdvance();
  const cancelMutation = useCancelAdvance();

  const handleApprove = async (advance: Advance) => {
    try {
      await approveMutation.mutateAsync(advance.id);
    } catch {
      // Error handled by mutation
    }
  };

  const handleReject = (advance: Advance) => {
    setRejectingAdvance(advance);
    setRejectDialogOpen(true);
  };

  const handleCancel = async (advance: Advance) => {
    try {
      await cancelMutation.mutateAsync(advance.id);
    } catch {
      // Error handled by mutation
    }
  };

  const handleDelete = (advance: Advance) => {
    setDeletingAdvance(advance);
    setDeleteDialogOpen(true);
  };

  // Stats
  const stats = useMemo(() => {
    return {
      pending: parseFloat(pendingData?.total_pending || '0'),
      approved: parseFloat(pendingData?.total_approved || '0'),
      pendingCount: pendingData?.advances?.filter(a => a.status === AdvanceStatus.PENDING).length || 0,
    };
  }, [pendingData]);

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold">Adiantamentos</h2>
          <p className="text-sm text-muted-foreground">
            Gerencie solicitações e aprovações de adiantamentos
          </p>
        </div>
        <Button onClick={() => setFormDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Novo Adiantamento
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardContent className="pt-4">
            <div className="flex items-center gap-2">
              <Clock className="h-5 w-5 text-yellow-500" />
              <div>
                <div className="text-2xl font-bold">{formatCurrency(stats.pending)}</div>
                <p className="text-sm text-muted-foreground">
                  {stats.pendingCount} pendente{stats.pendingCount !== 1 ? 's' : ''}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-4">
            <div className="flex items-center gap-2">
              <CheckCircle2 className="h-5 w-5 text-green-500" />
              <div>
                <div className="text-2xl font-bold text-green-600">{formatCurrency(stats.approved)}</div>
                <p className="text-sm text-muted-foreground">Aprovados (a deduzir)</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className={cn(stats.pendingCount > 0 && 'border-yellow-500')}>
          <CardContent className="pt-4">
            <div className="flex items-center gap-2">
              <AlertCircle className={cn('h-5 w-5', stats.pendingCount > 0 ? 'text-yellow-500' : 'text-muted-foreground')} />
              <div>
                <p className="text-sm font-medium">
                  {stats.pendingCount > 0 
                    ? `${stats.pendingCount} adiantamento${stats.pendingCount !== 1 ? 's' : ''} aguardando aprovação`
                    : 'Nenhum adiantamento pendente'}
                </p>
              </div>
            </div>
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
                  <SelectItem value={AdvanceStatus.PENDING}>Pendentes</SelectItem>
                  <SelectItem value={AdvanceStatus.APPROVED}>Aprovados</SelectItem>
                  <SelectItem value={AdvanceStatus.REJECTED}>Rejeitados</SelectItem>
                  <SelectItem value={AdvanceStatus.DEDUCTED}>Deduzidos</SelectItem>
                  <SelectItem value={AdvanceStatus.CANCELLED}>Cancelados</SelectItem>
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
                  <TableHead>Valor</TableHead>
                  <TableHead>Data</TableHead>
                  <TableHead>Motivo</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-40" /></TableCell>
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
                  <TableHead>Valor</TableHead>
                  <TableHead>Data</TableHead>
                  <TableHead>Motivo</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((advance) => (
                  <TableRow key={advance.id}>
                    <TableCell className="font-medium">
                      {advance.professional_name || '-'}
                    </TableCell>
                    <TableCell className="font-semibold">
                      {formatCurrency(advance.amount)}
                    </TableCell>
                    <TableCell>{formatDate(advance.request_date)}</TableCell>
                    <TableCell className="max-w-[200px] truncate">
                      {advance.reason || '-'}
                    </TableCell>
                    <TableCell>
                      {getStatusBadge(advance.status as AdvanceStatus)}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          {advance.status === AdvanceStatus.PENDING && (
                            <>
                              <DropdownMenuItem 
                                onClick={() => handleApprove(advance)}
                                disabled={approveMutation.isPending}
                              >
                                <Check className="mr-2 h-4 w-4 text-green-600" />
                                Aprovar
                              </DropdownMenuItem>
                              <DropdownMenuItem onClick={() => handleReject(advance)}>
                                <X className="mr-2 h-4 w-4 text-red-600" />
                                Rejeitar
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                            </>
                          )}
                          {advance.status === AdvanceStatus.APPROVED && (
                            <>
                              <DropdownMenuItem 
                                onClick={() => handleCancel(advance)}
                                disabled={cancelMutation.isPending}
                              >
                                <Ban className="mr-2 h-4 w-4" />
                                Cancelar
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                            </>
                          )}
                          <DropdownMenuItem
                            onClick={() => handleDelete(advance)}
                            className="text-destructive"
                            disabled={advance.status === AdvanceStatus.DEDUCTED}
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
              <HandCoins className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium">Nenhum adiantamento encontrado</h3>
              <p className="text-sm text-muted-foreground mb-4">
                Registre um adiantamento para um profissional
              </p>
              <Button onClick={() => setFormDialogOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Novo Adiantamento
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Dialogs */}
      <AdvanceFormDialog
        open={formDialogOpen}
        onOpenChange={setFormDialogOpen}
      />

      {rejectingAdvance && (
        <RejectDialog
          advance={rejectingAdvance}
          open={rejectDialogOpen}
          onOpenChange={(open) => {
            setRejectDialogOpen(open);
            if (!open) setRejectingAdvance(undefined);
          }}
        />
      )}

      {deletingAdvance && (
        <DeleteDialog
          advance={deletingAdvance}
          open={deleteDialogOpen}
          onOpenChange={(open) => {
            setDeleteDialogOpen(open);
            if (!open) setDeletingAdvance(undefined);
          }}
        />
      )}
    </div>
  );
}
