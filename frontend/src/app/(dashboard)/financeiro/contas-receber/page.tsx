/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Contas a Receber
 *
 * Lista receitas com filtros, status e ações de gestão.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import {
    AlertCircle,
    ArrowUpCircle,
    Check,
    Edit,
    Plus,
    Trash2,
    X
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
    useCancelReceivable,
    useCreateReceivable,
    useDeleteReceivable,
    useReceivables,
    useReceiveReceivable,
    useUpdateReceivable,
} from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import {
    ContaReceber,
    CreateContaReceberRequest,
    OrigemReceita,
    StatusContaReceber,
    UpdateContaReceberRequest,
    formatCurrency,
    formatDate,
    getDiasParaVencimento,
    getOrigemLabel,
    getStatusBadgeVariant,
    getStatusReceberLabel,
} from '@/types/financial';

// =============================================================================
// COMPONENTE: Modal de Criação/Edição
// =============================================================================

interface ReceivableFormModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editingItem?: ContaReceber | null;
}

function ReceivableFormModal({ open, onOpenChange, editingItem }: ReceivableFormModalProps) {
  const createMutation = useCreateReceivable();
  const updateMutation = useUpdateReceivable();

  // Calcula formData inicial baseado em editingItem
  const getInitialFormData = (): CreateContaReceberRequest => {
    if (editingItem) {
      return {
        origem: editingItem.origem,
        descricao_origem: editingItem.descricao_origem,
        valor: editingItem.valor,
        data_vencimento: editingItem.data_vencimento.slice(0, 10),
        metodo_pagamento: editingItem.metodo_pagamento || '',
        observacoes: editingItem.observacoes || '',
      };
    }
    return {
      origem: OrigemReceita.SERVICO,
      descricao_origem: '',
      valor: '',
      data_vencimento: '',
      metodo_pagamento: '',
      observacoes: '',
    };
  };

  // Estado é inicializado com base no editingItem - resetado via key no componente pai
  const [formData, setFormData] = useState<CreateContaReceberRequest>(getInitialFormData);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      if (editingItem) {
        await updateMutation.mutateAsync({
          id: editingItem.id,
          data: formData as UpdateContaReceberRequest,
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
            {editingItem ? 'Editar Conta a Receber' : 'Nova Conta a Receber'}
          </DialogTitle>
          <DialogDescription>
            {editingItem
              ? 'Atualize os dados da receita'
              : 'Cadastre uma nova receita a receber'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="descricao_origem">Descrição *</Label>
            <Input
              id="descricao_origem"
              value={formData.descricao_origem}
              onChange={(e) =>
                setFormData({ ...formData, descricao_origem: e.target.value })
              }
              placeholder="Ex: Corte de cabelo - João"
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
            <Label htmlFor="origem">Origem *</Label>
            <Select
              value={formData.origem}
              onValueChange={(value) =>
                setFormData({ ...formData, origem: value as OrigemReceita })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione a origem" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value={OrigemReceita.SERVICO}>Serviço</SelectItem>
                <SelectItem value={OrigemReceita.PRODUTO}>Produto</SelectItem>
                <SelectItem value={OrigemReceita.ASSINATURA}>Assinatura</SelectItem>
                <SelectItem value={OrigemReceita.OUTRO}>Outro</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="metodo_pagamento">Método de Pagamento</Label>
            <Select
              value={formData.metodo_pagamento || ''}
              onValueChange={(value) =>
                setFormData({ ...formData, metodo_pagamento: value })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione (opcional)" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="DINHEIRO">Dinheiro</SelectItem>
                <SelectItem value="PIX">PIX</SelectItem>
                <SelectItem value="CARTAO_CREDITO">Cartão Crédito</SelectItem>
                <SelectItem value="CARTAO_DEBITO">Cartão Débito</SelectItem>
                <SelectItem value="BOLETO">Boleto</SelectItem>
              </SelectContent>
            </Select>
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
// COMPONENTE: Modal de Confirmação de Recebimento
// =============================================================================

interface ReceiveConfirmModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  item: ContaReceber | null;
}

function ReceiveConfirmModal({ open, onOpenChange, item }: ReceiveConfirmModalProps) {
  const receiveMutation = useReceiveReceivable();
  const [dataRecebimento, setDataRecebimento] = useState(
    new Date().toISOString().slice(0, 10)
  );

  const handleConfirm = async () => {
    if (!item) return;

    try {
      await receiveMutation.mutateAsync({
        id: item.id,
        data_recebimento: dataRecebimento,
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
          <DialogTitle>Confirmar Recebimento</DialogTitle>
          <DialogDescription>
            Marcar esta conta como recebida
          </DialogDescription>
        </DialogHeader>

        {item && (
          <div className="space-y-4 py-4">
            <div className="rounded-lg bg-muted p-4">
              <p className="font-medium">{item.descricao_origem}</p>
              <p className="text-2xl font-bold text-green-600 mt-1">
                {formatCurrency(parseFloat(item.valor))}
              </p>
              <p className="text-sm text-muted-foreground mt-1">
                Vencimento: {formatDate(item.data_vencimento)}
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="data_recebimento">Data do Recebimento</Label>
              <Input
                id="data_recebimento"
                type="date"
                value={dataRecebimento}
                onChange={(e) => setDataRecebimento(e.target.value)}
              />
            </div>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancelar
          </Button>
          <Button onClick={handleConfirm} disabled={receiveMutation.isPending}>
            <Check className="mr-2 h-4 w-4" />
            {receiveMutation.isPending ? 'Processando...' : 'Confirmar Recebimento'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL: Contas a Receber
// =============================================================================

export default function ContasReceberPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estados
  const [statusFilter, setStatusFilter] = useState<StatusContaReceber | 'all'>('all');
  const [origemFilter, setOrigemFilter] = useState<OrigemReceita | 'all'>('all');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<ContaReceber | null>(null);
  const [receivingItem, setReceivingItem] = useState<ContaReceber | null>(null);

  // Queries
  const { data: receivables, isLoading } = useReceivables({
    status: statusFilter !== 'all' ? statusFilter : undefined,
    origem: origemFilter !== 'all' ? origemFilter : undefined,
  });

  const deleteMutation = useDeleteReceivable();
  const cancelMutation = useCancelReceivable();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'Contas a Receber' },
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
        router.replace('/financeiro/contas-receber');
      }, 0);
      return () => clearTimeout(timer);
    }
  }, [searchParams, router]);

  // Handlers
  const handleEdit = (item: ContaReceber) => {
    setEditingItem(item);
    setIsFormOpen(true);
  };

  const handleCloseForm = (open: boolean) => {
    setIsFormOpen(open);
    if (!open) {
      setEditingItem(null);
    }
  };

  const handleDelete = async (item: ContaReceber) => {
    if (confirm(`Deseja realmente excluir "${item.descricao_origem}"?`)) {
      await deleteMutation.mutateAsync(item.id);
    }
  };

  const handleCancel = async (item: ContaReceber) => {
    if (confirm(`Deseja cancelar a conta "${item.descricao_origem}"?`)) {
      await cancelMutation.mutateAsync(item.id);
    }
  };

  const handleReceive = (item: ContaReceber) => {
    setReceivingItem(item);
  };

  // Cálculos de resumo
  const { totalPendente, totalRecebido, contasAtraso } = useMemo(() => {
    if (!receivables) return { totalPendente: 0, totalRecebido: 0, contasAtraso: 0 };

    const pendentes = receivables.filter((i) => i.status === StatusContaReceber.PENDENTE);
    const recebidos = receivables.filter((i) => i.status === StatusContaReceber.RECEBIDO);
    const atrasadas = receivables.filter((i) => i.status === StatusContaReceber.ATRASADO);

    return {
      totalPendente: pendentes.reduce((acc, i) => acc + parseFloat(i.valor), 0),
      totalRecebido: recebidos.reduce((acc, i) => acc + parseFloat(i.valor), 0),
      contasAtraso: atrasadas.length,
    };
  }, [receivables]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Contas a Receber</h1>
          <p className="text-muted-foreground">
            Gerencie suas receitas e recebimentos
          </p>
        </div>
        <Button onClick={() => setIsFormOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Receita
        </Button>
      </div>

      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">A Receber</CardTitle>
            <ArrowUpCircle className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {formatCurrency(totalPendente)}
            </div>
            <p className="text-xs text-muted-foreground">
              {receivables?.filter((i) => i.status === StatusContaReceber.PENDENTE).length ?? 0} contas
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Recebido</CardTitle>
            <Check className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {formatCurrency(totalRecebido)}
            </div>
            <p className="text-xs text-muted-foreground">
              {receivables?.filter((i) => i.status === StatusContaReceber.RECEBIDO).length ?? 0} contas
            </p>
          </CardContent>
        </Card>

        <Card className={contasAtraso > 0 ? 'border-destructive' : ''}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Em Atraso</CardTitle>
            <AlertCircle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${contasAtraso > 0 ? 'text-destructive' : 'text-muted-foreground'}`}>
              {contasAtraso}
            </div>
            <p className="text-xs text-muted-foreground">contas vencidas</p>
          </CardContent>
        </Card>
      </div>

      {/* Filtros */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col gap-4 md:flex-row">
            <Select
              value={statusFilter}
              onValueChange={(value) => setStatusFilter(value as StatusContaReceber | 'all')}
            >
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos Status</SelectItem>
                <SelectItem value={StatusContaReceber.PENDENTE}>Pendente</SelectItem>
                <SelectItem value={StatusContaReceber.RECEBIDO}>Recebido</SelectItem>
                <SelectItem value={StatusContaReceber.PARCIAL}>Parcial</SelectItem>
                <SelectItem value={StatusContaReceber.ATRASADO}>Atrasado</SelectItem>
                <SelectItem value={StatusContaReceber.CANCELADO}>Cancelado</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={origemFilter}
              onValueChange={(value) => setOrigemFilter(value as OrigemReceita | 'all')}
            >
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Origem" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas Origens</SelectItem>
                <SelectItem value={OrigemReceita.SERVICO}>Serviço</SelectItem>
                <SelectItem value={OrigemReceita.PRODUTO}>Produto</SelectItem>
                <SelectItem value={OrigemReceita.ASSINATURA}>Assinatura</SelectItem>
                <SelectItem value={OrigemReceita.OUTRO}>Outro</SelectItem>
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
          ) : receivables && receivables.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Descrição</TableHead>
                  <TableHead>Origem</TableHead>
                  <TableHead className="text-right">Valor</TableHead>
                  <TableHead>Vencimento</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {receivables.map((item) => {
                  const diasParaVencer = getDiasParaVencimento(item.data_vencimento);
                  const isAtrasado = item.status === StatusContaReceber.ATRASADO;
                  const isUrgente = diasParaVencer <= 3 && diasParaVencer >= 0 && item.status === StatusContaReceber.PENDENTE;

                  return (
                    <TableRow key={item.id} className={isAtrasado ? 'bg-red-50 dark:bg-red-950/20' : ''}>
                      <TableCell>
                        <p className="font-medium">{item.descricao_origem}</p>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {getOrigemLabel(item.origem)}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right font-semibold text-green-600">
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
                          {getStatusReceberLabel(item.status)}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex justify-end gap-1">
                          <TooltipProvider>
                            {item.status === StatusContaReceber.PENDENTE && (
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => handleReceive(item)}
                                    className="h-8 w-8 text-green-600"
                                  >
                                    <Check className="h-4 w-4" />
                                  </Button>
                                </TooltipTrigger>
                                <TooltipContent>Marcar como Recebido</TooltipContent>
                              </Tooltip>
                            )}

                            {(item.status === StatusContaReceber.PENDENTE || item.status === StatusContaReceber.ATRASADO) && (
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
                  Nenhuma conta a receber encontrada. Clique em &ldquo;Nova Receita&rdquo; para cadastrar.
                </AlertDescription>
              </Alert>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modais */}
      <ReceivableFormModal
        key={editingItem?.id || 'new'}
        open={isFormOpen}
        onOpenChange={handleCloseForm}
        editingItem={editingItem}
      />

      <ReceiveConfirmModal
        open={!!receivingItem}
        onOpenChange={(open) => !open && setReceivingItem(null)}
        item={receivingItem}
      />
    </div>
  );
}
