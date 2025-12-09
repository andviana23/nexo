/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Despesas Fixas
 *
 * Lista despesas fixas recorrentes com toggle de ativo e geração de contas.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import {
    AlertCircle,
    Calendar,
    ChevronLeft,
    ChevronRight,
    Edit,
    Loader2,
    PauseCircle,
    PlayCircle,
    Plus,
    RefreshCw,
    Trash2
} from 'lucide-react';
import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';

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
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from '@/components/ui/tooltip';

import {
    useDeleteFixedExpense,
    useFixedExpenses,
    useFixedExpensesSummary,
    useGeneratePayablesFromFixed,
    useToggleFixedExpense,
} from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import type { DespesaFixa } from '@/types/financial';
import { formatCurrency } from '@/types/financial';

// =============================================================================
// COMPONENTE: Resumo
// =============================================================================

function SummaryCards() {
  const { data: summary, isLoading } = useFixedExpensesSummary();

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {[1, 2, 3].map((i) => (
          <Card key={i}>
            <CardHeader className="pb-2">
              <Skeleton className="h-4 w-24" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!summary) {
    return null;
  }

  return (
    <div className="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Total de Despesas
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{summary.total}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Despesas Ativas
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">{summary.total_ativas}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Valor Total Mensal
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-blue-600">
            {formatCurrency(parseFloat(summary.valor_total))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// COMPONENTE: Modal de Confirmação de Exclusão
// =============================================================================

interface DeleteConfirmModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  despesa: DespesaFixa | null;
  onConfirm: () => void;
  isDeleting: boolean;
}

function DeleteConfirmModal({
  open,
  onOpenChange,
  despesa,
  onConfirm,
  isDeleting,
}: DeleteConfirmModalProps) {
  if (!despesa) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Confirmar Exclusão</DialogTitle>
          <DialogDescription>
            Tem certeza que deseja excluir a despesa fixa &quot;{despesa.descricao}&quot;?
            Esta ação não pode ser desfeita.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isDeleting}>
            Cancelar
          </Button>
          <Button variant="destructive" onClick={onConfirm} disabled={isDeleting}>
            {isDeleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Excluir
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// COMPONENTE: Modal de Geração de Contas
// =============================================================================

interface GenerateModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function GenerateModal({ open, onOpenChange }: GenerateModalProps) {
  const now = new Date();
  const [ano, setAno] = useState(now.getFullYear());
  const [mes, setMes] = useState(now.getMonth() + 1);
  
  const generateMutation = useGeneratePayablesFromFixed();

  const meses = [
    { value: 1, label: 'Janeiro' },
    { value: 2, label: 'Fevereiro' },
    { value: 3, label: 'Março' },
    { value: 4, label: 'Abril' },
    { value: 5, label: 'Maio' },
    { value: 6, label: 'Junho' },
    { value: 7, label: 'Julho' },
    { value: 8, label: 'Agosto' },
    { value: 9, label: 'Setembro' },
    { value: 10, label: 'Outubro' },
    { value: 11, label: 'Novembro' },
    { value: 12, label: 'Dezembro' },
  ];

  const handleGenerate = () => {
    generateMutation.mutate(
      { ano, mes },
      {
        onSuccess: () => {
          onOpenChange(false);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <RefreshCw className="h-5 w-5" />
            Gerar Contas a Pagar
          </DialogTitle>
          <DialogDescription>
            Gera contas a pagar para todas as despesas fixas ativas no mês selecionado.
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="ano">Ano</Label>
            <Input
              id="ano"
              type="number"
              min={2020}
              max={2099}
              value={ano}
              onChange={(e) => setAno(parseInt(e.target.value, 10))}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="mes">Mês</Label>
            <Select value={mes.toString()} onValueChange={(v) => setMes(parseInt(v, 10))}>
              <SelectTrigger id="mes">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {meses.map((m) => (
                  <SelectItem key={m.value} value={m.value.toString()}>
                    {m.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={generateMutation.isPending}>
            Cancelar
          </Button>
          <Button onClick={handleGenerate} disabled={generateMutation.isPending}>
            {generateMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Gerar Contas
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// COMPONENTE: Tabela de Despesas Fixas
// =============================================================================

interface FixedExpensesTableProps {
  filter: 'all' | 'active' | 'inactive';
  page: number;
  onPageChange: (page: number) => void;
}

function FixedExpensesTable({ filter, page, onPageChange }: FixedExpensesTableProps) {
  const [deleteItem, setDeleteItem] = useState<DespesaFixa | null>(null);
  const pageSize = 10;
  
  const filterParams = useMemo(() => {
    const params: { ativo?: boolean; page: number; page_size: number } = {
      page,
      page_size: pageSize,
    };
    if (filter === 'active') params.ativo = true;
    if (filter === 'inactive') params.ativo = false;
    return params;
  }, [filter, page]);

  const { data: response, isLoading, error } = useFixedExpenses(filterParams);
  const deleteMutation = useDeleteFixedExpense();
  const toggleMutation = useToggleFixedExpense();

  const handleDelete = () => {
    if (!deleteItem) return;
    deleteMutation.mutate(deleteItem.id, {
      onSuccess: () => {
        setDeleteItem(null);
      },
    });
  };

  const handleToggle = (id: string) => {
    toggleMutation.mutate(id);
  };

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="space-y-3">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="flex items-center gap-4">
                <Skeleton className="h-4 w-[30%]" />
                <Skeleton className="h-4 w-[20%]" />
                <Skeleton className="h-4 w-[15%]" />
                <Skeleton className="h-4 w-[10%]" />
                <Skeleton className="h-8 w-24" />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          Erro ao carregar despesas fixas. Tente novamente.
        </AlertDescription>
      </Alert>
    );
  }

  const despesas = response?.data ?? [];

  if (despesas.length === 0 && page === 1) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <Calendar className="h-12 w-12 text-muted-foreground mb-4" />
          <h3 className="text-lg font-medium">Nenhuma despesa fixa cadastrada</h3>
          <p className="text-sm text-muted-foreground mb-4">
            Cadastre suas despesas recorrentes para gerenciá-las de forma automática.
          </p>
          <Link href="/financeiro/despesas-fixas/nova">
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Nova Despesa Fixa
            </Button>
          </Link>
        </CardContent>
      </Card>
    );
  }

  const totalPages = response?.total_pages ?? 1;
  const total = response?.total ?? 0;

  return (
    <>
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Descrição</TableHead>
                <TableHead>Fornecedor</TableHead>
                <TableHead>Dia Venc.</TableHead>
                <TableHead className="text-right">Valor</TableHead>
                <TableHead className="text-center">Status</TableHead>
                <TableHead className="text-right">Ações</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {despesas.map((despesa) => (
                <TableRow key={despesa.id}>
                  <TableCell className="font-medium">
                    {despesa.descricao}
                    {despesa.observacoes && (
                      <p className="text-xs text-muted-foreground truncate max-w-[200px]">
                        {despesa.observacoes}
                      </p>
                    )}
                  </TableCell>
                  <TableCell>{despesa.fornecedor || '-'}</TableCell>
                  <TableCell>
                    <Badge variant="outline">Dia {despesa.dia_vencimento}</Badge>
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    {formatCurrency(parseFloat(despesa.valor))}
                  </TableCell>
                  <TableCell className="text-center">
                    <Badge variant={despesa.ativo ? 'default' : 'secondary'}>
                      {despesa.ativo ? 'Ativa' : 'Inativa'}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <TooltipProvider>
                      <div className="flex items-center justify-end gap-1">
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleToggle(despesa.id)}
                              disabled={toggleMutation.isPending}
                            >
                              {despesa.ativo ? (
                                <PauseCircle className="h-4 w-4 text-yellow-500" />
                              ) : (
                                <PlayCircle className="h-4 w-4 text-green-500" />
                              )}
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>
                            {despesa.ativo ? 'Desativar' : 'Ativar'}
                          </TooltipContent>
                        </Tooltip>

                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Link href={`/financeiro/despesas-fixas/${despesa.id}/editar`}>
                              <Button variant="ghost" size="icon">
                                <Edit className="h-4 w-4" />
                              </Button>
                            </Link>
                          </TooltipTrigger>
                          <TooltipContent>Editar</TooltipContent>
                        </Tooltip>

                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => setDeleteItem(despesa)}
                            >
                              <Trash2 className="h-4 w-4 text-destructive" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>Excluir</TooltipContent>
                        </Tooltip>
                      </div>
                    </TooltipProvider>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>

        {/* Paginação */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-6 py-4 border-t">
            <div className="text-sm text-muted-foreground">
              Mostrando {((page - 1) * pageSize) + 1} a {Math.min(page * pageSize, total)} de {total} itens
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => onPageChange(Math.max(1, page - 1))}
                disabled={page === 1}
              >
                <ChevronLeft className="h-4 w-4" />
                Anterior
              </Button>
              <div className="flex items-center gap-1">
                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  let pageNum: number;
                  if (totalPages <= 5) {
                    pageNum = i + 1;
                  } else if (page <= 3) {
                    pageNum = i + 1;
                  } else if (page >= totalPages - 2) {
                    pageNum = totalPages - 4 + i;
                  } else {
                    pageNum = page - 2 + i;
                  }
                  return (
                    <Button
                      key={pageNum}
                      variant={page === pageNum ? 'default' : 'outline'}
                      size="sm"
                      className="w-9"
                      onClick={() => onPageChange(pageNum)}
                    >
                      {pageNum}
                    </Button>
                  );
                })}
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => onPageChange(Math.min(totalPages, page + 1))}
                disabled={page === totalPages}
              >
                Próximo
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        )}
      </Card>

      <DeleteConfirmModal
        open={!!deleteItem}
        onOpenChange={(open) => !open && setDeleteItem(null)}
        despesa={deleteItem}
        onConfirm={handleDelete}
        isDeleting={deleteMutation.isPending}
      />
    </>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL
// =============================================================================

export default function DespesasFixasPage() {
  const [filter, setFilter] = useState<'all' | 'active' | 'inactive'>('all');
  const [page, setPage] = useState(1);
  const [generateModalOpen, setGenerateModalOpen] = useState(false);

  // Breadcrumbs
  const { setBreadcrumbs } = useBreadcrumbs();
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'Despesas Fixas' },
    ]);
  }, [setBreadcrumbs]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Despesas Fixas</h1>
          <p className="text-muted-foreground">
            Gerencie despesas recorrentes mensais
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => setGenerateModalOpen(true)}>
            <RefreshCw className="mr-2 h-4 w-4" />
            Gerar Contas
          </Button>
          <Link href="/financeiro/despesas-fixas/nova">
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Nova Despesa
            </Button>
          </Link>
        </div>
      </div>

      {/* Cards de Resumo */}
      <SummaryCards />

      {/* Filtros */}
      <div className="flex items-center gap-4">
        <Label>Filtrar por status:</Label>
        <Select
          value={filter}
          onValueChange={(v) => {
            setFilter(v as typeof filter);
            setPage(1);
          }}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todas</SelectItem>
            <SelectItem value="active">Ativas</SelectItem>
            <SelectItem value="inactive">Inativas</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Tabela */}
      <FixedExpensesTable filter={filter} page={page} onPageChange={setPage} />

      {/* Modal de Geração */}
      <GenerateModal open={generateModalOpen} onOpenChange={setGenerateModalOpen} />
    </div>
  );
}
