/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Contas a Pagar
 *
 * Lista despesas com filtros, status e ações de gestão.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import {
    AlertCircle,
    ArrowDownCircle,
    Calendar,
    Check,
    Edit,
    Plus,
    Trash2,
    X,
} from 'lucide-react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useMemo, useRef, useState } from 'react';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
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
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from '@/components/ui/tooltip';

import {
    useCancelPayable,
    useCreatePayable,
    useDeletePayable,
    usePayPayable,
    usePayables,
    useUpdatePayable,
} from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import {
    ContaPagar,
    CreateContaPagarRequest,
    StatusContaPagar,
    TipoDespesa,
    UpdateContaPagarRequest,
    formatCurrency,
    formatDate,
    getDiasParaVencimento,
    getStatusBadgeVariant,
    getStatusPagarLabel,
    getTipoDespesaLabel,
} from '@/types/financial';

// =============================================================================
// COMPONENTE: Modal de Criação/Edição
// =============================================================================

interface PayableFormModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editingItem?: ContaPagar | null;
}

function PayableFormModal({ open, onOpenChange, editingItem }: PayableFormModalProps) {
  const createMutation = useCreatePayable();
  const updateMutation = useUpdatePayable();

  // Calcula formData inicial baseado em editingItem
  const getInitialFormData = (): CreateContaPagarRequest => {
    if (editingItem) {
      return {
        descricao: editingItem.descricao,
        valor: editingItem.valor,
        data_vencimento: editingItem.data_vencimento.slice(0, 10),
        tipo: editingItem.tipo,
        fornecedor: editingItem.fornecedor || '',
        observacoes: editingItem.observacoes || '',
        recorrente: editingItem.recorrente || false,
      };
    }
    return {
      descricao: '',
      valor: '',
      data_vencimento: '',
      tipo: TipoDespesa.VARIAVEL,
      fornecedor: '',
      observacoes: '',
      recorrente: false,
    };
  };

  // Estado é inicializado com base no editingItem - resetado via key no componente pai
  const [formData, setFormData] = useState<CreateContaPagarRequest>(getInitialFormData);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      if (editingItem) {
        await updateMutation.mutateAsync({
          id: editingItem.id,
          data: formData as UpdateContaPagarRequest,
        });
      } else {
        await createMutation.mutateAsync(formData);
      }
      onOpenChange(false);
    } catch {
      // Erro já tratado pelo hook
    }
  };

  const isPending = createMutation.isPending || updateMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>
            {editingItem ? 'Editar Conta a Pagar' : 'Nova Conta a Pagar'}
          </DialogTitle>
          <DialogDescription>
            {editingItem
              ? 'Atualize os dados da despesa'
              : 'Cadastre uma nova despesa a pagar'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="descricao">Descrição *</Label>
            <Input
              id="descricao"
              value={formData.descricao}
              onChange={(e) =>
                setFormData({ ...formData, descricao: e.target.value })
              }
              placeholder="Ex: Aluguel do espaço"
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="valor">Valor *</Label>
              <Input
                id="valor"
                type="number"
                step="0.01"
                min="0.01"
                value={formData.valor}
                onChange={(e) =>
                  setFormData({ ...formData, valor: e.target.value })
                }
                placeholder="0,00"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="data_vencimento">Vencimento *</Label>
              <Input
                id="data_vencimento"
                type="date"
                value={formData.data_vencimento}
                onChange={(e) =>
                  setFormData({ ...formData, data_vencimento: e.target.value })
                }
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="tipo">Tipo de Despesa *</Label>
            <Select
              value={formData.tipo}
              onValueChange={(value) =>
                setFormData({ ...formData, tipo: value as TipoDespesa })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione o tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={TipoDespesa.FIXA}>Fixa (recorrente)</SelectItem>
                <SelectItem value={TipoDespesa.VARIAVEL}>Variável</SelectItem>
                <SelectItem value={TipoDespesa.EXTRAORDINARIA}>Extraordinária</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="fornecedor">Fornecedor</Label>
            <Input
              id="fornecedor"
              value={formData.fornecedor}
              onChange={(e) =>
                setFormData({ ...formData, fornecedor: e.target.value })
              }
              placeholder="Nome do fornecedor (opcional)"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="observacoes">Observações</Label>
            <Textarea
              id="observacoes"
              value={formData.observacoes}
              onChange={(e) =>
                setFormData({ ...formData, observacoes: e.target.value })
              }
              placeholder="Observações adicionais..."
              rows={2}
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="recorrente"
              checked={formData.recorrente}
              onChange={(e) =>
                setFormData({ ...formData, recorrente: e.target.checked })
              }
              className="h-4 w-4"
            />
            <Label htmlFor="recorrente" className="cursor-pointer">
              Despesa recorrente (mensal)
            </Label>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Salvando...' : editingItem ? 'Atualizar' : 'Criar'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// COMPONENTE: Modal de Confirmação de Pagamento
// =============================================================================

interface PayConfirmModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  item: ContaPagar | null;
}

function PayConfirmModal({ open, onOpenChange, item }: PayConfirmModalProps) {
  const payMutation = usePayPayable();
  const [dataPagamento, setDataPagamento] = useState(
    new Date().toISOString().slice(0, 10)
  );

  const handleConfirm = async () => {
    if (!item) return;

    try {
      await payMutation.mutateAsync({
        id: item.id,
        data_pagamento: dataPagamento,
      });
      onOpenChange(false);
    } catch {
      // Erro já tratado pelo hook
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>Confirmar Pagamento</DialogTitle>
          <DialogDescription>
            Marcar esta conta como paga
          </DialogDescription>
        </DialogHeader>

        {item && (
          <div className="space-y-4 py-4">
            <div className="rounded-lg bg-muted p-4">
              <p className="font-medium">{item.descricao}</p>
              <p className="text-2xl font-bold text-green-600 mt-1">
                {formatCurrency(parseFloat(item.valor))}
              </p>
              <p className="text-sm text-muted-foreground mt-1">
                Vencimento: {formatDate(item.data_vencimento)}
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="data_pagamento">Data do Pagamento</Label>
              <Input
                id="data_pagamento"
                type="date"
                value={dataPagamento}
                onChange={(e) => setDataPagamento(e.target.value)}
              />
            </div>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancelar
          </Button>
          <Button onClick={handleConfirm} disabled={payMutation.isPending}>
            <Check className="mr-2 h-4 w-4" />
            {payMutation.isPending ? 'Processando...' : 'Confirmar Pagamento'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL: Contas a Pagar
// =============================================================================

export default function ContasPagarPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estados
  const [statusFilter, setStatusFilter] = useState<StatusContaPagar | 'all'>('all');
  const [tipoFilter, setTipoFilter] = useState<TipoDespesa | 'all'>('all');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<ContaPagar | null>(null);
  const [payingItem, setPayingItem] = useState<ContaPagar | null>(null);

  // Queries
  const { data: payables, isLoading } = usePayables({
    status: statusFilter !== 'all' ? statusFilter : undefined,
    tipo: tipoFilter !== 'all' ? tipoFilter : undefined,
  });

  const deleteMutation = useDeletePayable();
  const cancelMutation = useCancelPayable();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'Contas a Pagar' },
    ]);
  }, [setBreadcrumbs]);

  // Abre modal se ?modal=new (usando useRef para evitar re-renders)
  const hasOpenedModal = useRef(false);
  useEffect(() => {
    if (searchParams.get('modal') === 'new' && !hasOpenedModal.current) {
      hasOpenedModal.current = true;
      // Delay para evitar setState during render
      const timer = setTimeout(() => {
        setIsFormOpen(true);
        router.replace('/financeiro/contas-pagar');
      }, 0);
      return () => clearTimeout(timer);
    }
  }, [searchParams, router]);

  // Handlers
  const handleEdit = (item: ContaPagar) => {
    setEditingItem(item);
    setIsFormOpen(true);
  };

  const handleCloseForm = (open: boolean) => {
    setIsFormOpen(open);
    if (!open) {
      setEditingItem(null);
    }
  };

  const handleDelete = async (item: ContaPagar) => {
    if (confirm(`Deseja realmente excluir "${item.descricao}"?`)) {
      await deleteMutation.mutateAsync(item.id);
    }
  };

  const handleCancel = async (item: ContaPagar) => {
    if (confirm(`Deseja cancelar a conta "${item.descricao}"?`)) {
      await cancelMutation.mutateAsync(item.id);
    }
  };

  const handlePay = (item: ContaPagar) => {
    setPayingItem(item);
  };

  // Cálculos de resumo
  const { totalPendente, totalAtrasado, contasVencendo } = useMemo(() => {
    if (!payables) return { totalPendente: 0, totalAtrasado: 0, contasVencendo: 0 };

    const pendentes = payables.filter((i) => i.status === StatusContaPagar.PENDENTE);
    const atrasadas = payables.filter((i) => i.status === StatusContaPagar.ATRASADO);
    const vencendoEmBreve = payables.filter((i) => {
      if (i.status !== StatusContaPagar.PENDENTE) return false;
      const dias = getDiasParaVencimento(i.data_vencimento);
      return dias >= 0 && dias <= 7;
    });

    return {
      totalPendente: pendentes.reduce((acc, i) => acc + parseFloat(i.valor), 0),
      totalAtrasado: atrasadas.reduce((acc, i) => acc + parseFloat(i.valor), 0),
      contasVencendo: vencendoEmBreve.length,
    };
  }, [payables]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Contas a Pagar</h1>
          <p className="text-muted-foreground">
            Gerencie suas despesas e pagamentos
          </p>
        </div>
        <Button onClick={() => setIsFormOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Despesa
        </Button>
      </div>

      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Pendente</CardTitle>
            <ArrowDownCircle className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-500">
              {formatCurrency(totalPendente)}
            </div>
            <p className="text-xs text-muted-foreground">
              {payables?.filter((i) => i.status === StatusContaPagar.PENDENTE).length ?? 0} contas
            </p>
          </CardContent>
        </Card>

        <Card className="border-destructive">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Atrasado</CardTitle>
            <AlertCircle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-destructive">
              {formatCurrency(totalAtrasado)}
            </div>
            <p className="text-xs text-muted-foreground">
              {payables?.filter((i) => i.status === StatusContaPagar.ATRASADO).length ?? 0} contas
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Vencendo em 7 dias</CardTitle>
            <Calendar className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">
              {contasVencendo}
            </div>
            <p className="text-xs text-muted-foreground">contas a vencer</p>
          </CardContent>
        </Card>
      </div>

      {/* Filtros */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col gap-4 md:flex-row">
            <Select
              value={statusFilter}
              onValueChange={(value) => setStatusFilter(value as StatusContaPagar | 'all')}
            >
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos Status</SelectItem>
                <SelectItem value={StatusContaPagar.PENDENTE}>Pendente</SelectItem>
                <SelectItem value={StatusContaPagar.PAGO}>Pago</SelectItem>
                <SelectItem value={StatusContaPagar.ATRASADO}>Atrasado</SelectItem>
                <SelectItem value={StatusContaPagar.CANCELADO}>Cancelado</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={tipoFilter}
              onValueChange={(value) => setTipoFilter(value as TipoDespesa | 'all')}
            >
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos Tipos</SelectItem>
                <SelectItem value={TipoDespesa.FIXA}>Fixa</SelectItem>
                <SelectItem value={TipoDespesa.VARIAVEL}>Variável</SelectItem>
                <SelectItem value={TipoDespesa.EXTRAORDINARIA}>Extraordinária</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Tabela de Contas */}
      <Card>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-6 space-y-4">
              {[1, 2, 3, 4].map((i) => (
                <Skeleton key={i} className="h-16 w-full" />
              ))}
            </div>
          ) : payables && payables.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Descrição</TableHead>
                  <TableHead>Fornecedor</TableHead>
                  <TableHead>Tipo</TableHead>
                  <TableHead className="text-right">Valor</TableHead>
                  <TableHead>Vencimento</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {payables.map((item) => {
                  const diasParaVencer = getDiasParaVencimento(item.data_vencimento);
                  const isAtrasado = item.status === StatusContaPagar.ATRASADO;
                  const isUrgente = diasParaVencer <= 3 && diasParaVencer >= 0 && item.status === StatusContaPagar.PENDENTE;

                  return (
                    <TableRow key={item.id} className={isAtrasado ? 'bg-red-50 dark:bg-red-950/20' : ''}>
                      <TableCell>
                        <div>
                          <p className="font-medium">{item.descricao}</p>
                          {item.recorrente && (
                            <Badge variant="outline" className="text-xs mt-1">
                              Recorrente
                            </Badge>
                          )}
                        </div>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {item.fornecedor || '-'}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {getTipoDespesaLabel(item.tipo)}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right font-semibold">
                        {formatCurrency(parseFloat(item.valor))}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {formatDate(item.data_vencimento)}
                          {isUrgente && (
                            <Badge variant="outline" className="text-xs border-yellow-500 text-yellow-600">
                              {diasParaVencer === 0 ? 'Hoje' : `${diasParaVencer}d`}
                            </Badge>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={getStatusBadgeVariant(item.status)}>
                          {getStatusPagarLabel(item.status)}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex justify-end gap-1">
                          <TooltipProvider>
                            {item.status === StatusContaPagar.PENDENTE && (
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => handlePay(item)}
                                    className="h-8 w-8 text-green-600"
                                  >
                                    <Check className="h-4 w-4" />
                                  </Button>
                                </TooltipTrigger>
                                <TooltipContent>Marcar como Pago</TooltipContent>
                              </Tooltip>
                            )}

                            {(item.status === StatusContaPagar.PENDENTE || item.status === StatusContaPagar.ATRASADO) && (
                              <>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      onClick={() => handleEdit(item)}
                                      className="h-8 w-8"
                                    >
                                      <Edit className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>Editar</TooltipContent>
                                </Tooltip>

                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      onClick={() => handleCancel(item)}
                                      className="h-8 w-8 text-orange-600"
                                    >
                                      <X className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>Cancelar</TooltipContent>
                                </Tooltip>
                              </>
                            )}

                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleDelete(item)}
                                  className="h-8 w-8 text-destructive"
                                  disabled={deleteMutation.isPending}
                                >
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>Excluir</TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          ) : (
            <div className="p-6">
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  Nenhuma conta a pagar encontrada. Clique em &ldquo;Nova Despesa&rdquo; para cadastrar.
                </AlertDescription>
              </Alert>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modais */}
      <PayableFormModal
        key={editingItem?.id || 'new'}
        open={isFormOpen}
        onOpenChange={handleCloseForm}
        editingItem={editingItem}
      />

      <PayConfirmModal
        open={!!payingItem}
        onOpenChange={(open) => !open && setPayingItem(null)}
        item={payingItem}
      />
    </div>
  );
}
