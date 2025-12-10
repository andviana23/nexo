'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  useDeleteMeioPagamento,
  useMeiosPagamento,
  useToggleMeioPagamento,
} from '@/hooks/use-meios-pagamento';
import {
  MeioPagamento,
  TIPO_PAGAMENTO_LABELS,
  TipoPagamento
} from '@/types/meio-pagamento';
import {
  Banknote,
  CreditCard,
  Edit,
  MoreHorizontal,
  Plus,
  Receipt,
  Trash2,
  Wallet
} from 'lucide-react';
import { useState } from 'react';
import { MeioPagamentoModal } from './meio-pagamento-modal';

export function MeiosPagamentoList() {
  const { data, isLoading } = useMeiosPagamento();
  const deleteMeioPagamento = useDeleteMeioPagamento();
  const toggleMeioPagamento = useToggleMeioPagamento();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [meioPagamentoToEdit, setMeioPagamentoToEdit] = useState<MeioPagamento | null>(null);
  const [meioPagamentoToDelete, setMeioPagamentoToDelete] = useState<MeioPagamento | null>(null);

  const handleCreate = () => {
    setMeioPagamentoToEdit(null);
    setIsModalOpen(true);
  };

  const handleEdit = (meioPagamento: MeioPagamento) => {
    setMeioPagamentoToEdit(meioPagamento);
    setIsModalOpen(true);
  };

  const handleDeleteClick = (meioPagamento: MeioPagamento) => {
    setMeioPagamentoToDelete(meioPagamento);
  };

  const handleConfirmDelete = async () => {
    if (meioPagamentoToDelete) {
      await deleteMeioPagamento.mutateAsync(meioPagamentoToDelete.id);
      setMeioPagamentoToDelete(null);
    }
  };

  const handleToggle = async (id: string) => {
    await toggleMeioPagamento.mutateAsync(id);
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <Skeleton className="h-8 w-64" />
            <Skeleton className="h-4 w-96" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <Card key={i}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-4 w-4 rounded-full" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-16" />
                <Skeleton className="h-3 w-32 mt-1" />
              </CardContent>
            </Card>
          ))}
        </div>
        <Card>
          <CardHeader>
            <Skeleton className="h-6 w-48 mb-2" />
            <Skeleton className="h-4 w-96" />
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {[1, 2, 3, 4].map((i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const meiosPagamento = data?.data || [];
  const total = data?.total || 0;
  const ativos = data?.total_ativo || 0;
  const inativos = total - ativos;

  return (
    <div className="space-y-8 container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header Info */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Tipos de Recebimento</h1>
          <p className="text-muted-foreground text-sm mt-1">
            Configure as formas de pagamento aceitas, taxas e prazos.
          </p>
        </div>
        <Button onClick={handleCreate} className="shadow-sm">
          <Plus className="mr-2 h-4 w-4" />
          Novo Tipo
        </Button>
      </div>

      {/* KPI Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Cadastrados</CardTitle>
            <Wallet className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{total}</div>
            <p className="text-xs text-muted-foreground">
              Opções de pagamento configuradas
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ativos</CardTitle>
            <CreditCard className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">{ativos}</div>
            <p className="text-xs text-muted-foreground">Disponíveis no caixa</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Inativos</CardTitle>
            <Banknote className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-muted-foreground">{inativos}</div>
            <p className="text-xs text-muted-foreground">Desabilitados temporariamente</p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Card className="border-border shadow-sm">
        <CardHeader className="bg-muted/40 pb-4">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg">Gerenciamento</CardTitle>
              <CardDescription>
                Lista de todos os métodos de pagamento cadastrados no sistema.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50 hover:bg-muted/50">
                  <TableHead className="w-[200px] pl-6">Tipo</TableHead>
                  <TableHead>Bandeira</TableHead>
                  <TableHead className="text-right">Taxa (%)</TableHead>
                  <TableHead className="text-right">Taxa Fixa</TableHead>
                  <TableHead className="text-center">Recebimento</TableHead>
                  <TableHead className="text-center">Status</TableHead>
                  <TableHead className="w-[80px] text-right pr-6">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {meiosPagamento.map((meio) => (
                  <TableRow key={meio.id} className="group hover:bg-muted/30">
                    <TableCell className="font-medium pl-6">
                      <div className="flex items-center gap-2">
                        <div
                          className="w-1.5 h-8 rounded-full"
                          style={{ backgroundColor: meio.cor || '#e5e7eb' }}
                        />
                        <span className="font-medium">
                          {TIPO_PAGAMENTO_LABELS[meio.tipo as TipoPagamento] || meio.tipo}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {meio.bandeira || '-'}
                    </TableCell>
                    <TableCell className="text-right font-mono text-xs">
                      {parseFloat(meio.taxa || '0').toFixed(2)}%
                    </TableCell>
                    <TableCell className="text-right font-mono text-xs">
                      {parseFloat(meio.taxa_fixa || '0') > 0
                        ? `R$ ${parseFloat(meio.taxa_fixa).toFixed(2)}`
                        : '-'}
                    </TableCell>
                    <TableCell className="text-center">
                      <Badge
                        variant={meio.d_mais === 0 ? 'secondary' : 'outline'}
                        className="font-mono text-xs"
                      >
                        D+{meio.d_mais}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-center">
                      <Switch
                        checked={meio.ativo}
                        onCheckedChange={() => handleToggle(meio.id)}
                        disabled={toggleMeioPagamento.isPending}
                        className="scale-90"
                      />
                    </TableCell>
                    <TableCell className="text-right pr-6">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-muted-foreground hover:text-foreground opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity"
                          >
                            <MoreHorizontal className="h-4 w-4" />
                            <span className="sr-only">Abrir menu</span>
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => handleEdit(meio)}>
                            <Edit className="mr-2 h-4 w-4" />
                            Editar
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            className="text-destructive focus:text-destructive"
                            onClick={() => handleDeleteClick(meio)}
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            Excluir
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
                {meiosPagamento.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={7} className="h-64 text-center">
                      <div className="flex flex-col items-center justify-center gap-3 text-muted-foreground">
                        <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
                          <Receipt className="h-6 w-6" />
                        </div>
                        <div className="space-y-1">
                          <p className="font-medium text-foreground">Nenhum meio de pagamento</p>
                          <p className="text-sm">Cadastre formas de pagamento para usar no caixa.</p>
                        </div>
                        <Button variant="outline" size="sm" onClick={handleCreate} className="mt-2">
                          Cadastrar Agora
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Modal de Criação/Edição */}
      <MeioPagamentoModal
        open={isModalOpen}
        onOpenChange={setIsModalOpen}
        meioPagamento={meioPagamentoToEdit}
      />

      {/* Dialog de Confirmação de Exclusão */}
      <Dialog
        open={!!meioPagamentoToDelete}
        onOpenChange={(open) => !open && setMeioPagamentoToDelete(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirmar Exclusão</DialogTitle>
            <DialogDescription>
              Tem certeza que deseja excluir este tipo de recebimento?
              <br />
              <span className="font-semibold text-foreground mt-2 block">
                {meioPagamentoToDelete && (TIPO_PAGAMENTO_LABELS[meioPagamentoToDelete.tipo as TipoPagamento] || meioPagamentoToDelete.tipo)}
                {meioPagamentoToDelete?.bandeira && ` - ${meioPagamentoToDelete.bandeira}`}
              </span>
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setMeioPagamentoToDelete(null)}>
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={deleteMeioPagamento.isPending}
            >
              {deleteMeioPagamento.isPending ? 'Excluindo...' : 'Excluir'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
