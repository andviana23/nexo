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
    CORES_TIPO_PAGAMENTO,
    MeioPagamento,
    TIPO_PAGAMENTO_LABELS,
    TipoPagamento,
} from '@/types/meio-pagamento';
import {
    ArrowLeftRight,
    Banknote,
    CircleDollarSign,
    CreditCard,
    Edit,
    FileText,
    MoreHorizontal,
    Plus,
    QrCode,
    Receipt,
    Trash2,
} from 'lucide-react';
import { useState } from 'react';
import { MeioPagamentoModal } from './meio-pagamento-modal';

// Mapa de ícones por tipo
const ICONES_TIPO: Record<TipoPagamento, React.ReactNode> = {
  DINHEIRO: <Banknote className="h-4 w-4 text-white" />,
  PIX: <QrCode className="h-4 w-4 text-white" />,
  CREDITO: <CreditCard className="h-4 w-4 text-white" />,
  DEBITO: <CreditCard className="h-4 w-4 text-white" />,
  TRANSFERENCIA: <ArrowLeftRight className="h-4 w-4 text-white" />,
  BOLETO: <FileText className="h-4 w-4 text-white" />,
  OUTRO: <CircleDollarSign className="h-4 w-4 text-white" />,
};

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
          <Skeleton className="h-10 w-40" />
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <Card key={i}>
              <CardHeader className="pb-2">
                <Skeleton className="h-4 w-32" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-16" />
              </CardContent>
            </Card>
          ))}
        </div>
        <Card>
          <CardContent className="p-0">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="flex items-center gap-4 p-4 border-b last:border-0">
                <Skeleton className="h-10 w-10 rounded-lg" />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-4 w-32" />
                  <Skeleton className="h-3 w-48" />
                </div>
                <Skeleton className="h-8 w-8" />
              </div>
            ))}
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
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Tipos de Recebimento</h1>
          <p className="text-muted-foreground">
            Gerencie os meios de pagamento aceitos e configure taxas e prazos de recebimento.
          </p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Novo Tipo
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Cadastrados</CardTitle>
            <Receipt className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{total}</div>
            <p className="text-xs text-muted-foreground">
              {total === 1 ? 'tipo de recebimento' : 'tipos de recebimento'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ativos</CardTitle>
            <CreditCard className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{ativos}</div>
            <p className="text-xs text-muted-foreground">disponíveis para uso</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Inativos</CardTitle>
            <CreditCard className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-muted-foreground">{inativos}</div>
            <p className="text-xs text-muted-foreground">desabilitados</p>
          </CardContent>
        </Card>
      </div>

      {/* Table Card */}
      <Card>
        <CardHeader>
          <CardTitle>Lista de Tipos de Recebimento</CardTitle>
          <CardDescription>
            Configure taxas, prazos D+ e bandeiras para cada meio de pagamento.
          </CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[250px]">Nome</TableHead>
                <TableHead>Tipo</TableHead>
                <TableHead>Bandeira</TableHead>
                <TableHead className="text-right">Taxa (%)</TableHead>
                <TableHead className="text-right">Taxa Fixa</TableHead>
                <TableHead className="text-center">D+</TableHead>
                <TableHead className="text-center">Status</TableHead>
                <TableHead className="w-[80px] text-right">Ações</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {meiosPagamento.map((meio) => (
                <TableRow key={meio.id} className="group">
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <div
                        className="h-10 w-10 rounded-lg flex items-center justify-center shadow-sm"
                        style={{
                          backgroundColor:
                            meio.cor || CORES_TIPO_PAGAMENTO[meio.tipo as TipoPagamento] || '#6366f1',
                        }}
                      >
                        {ICONES_TIPO[meio.tipo as TipoPagamento] || (
                          <CircleDollarSign className="h-4 w-4 text-white" />
                        )}
                      </div>
                      <div>
                        <p className="font-medium">{meio.nome}</p>
                        {meio.observacoes && (
                          <p className="text-xs text-muted-foreground line-clamp-1">
                            {meio.observacoes}
                          </p>
                        )}
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">
                      {TIPO_PAGAMENTO_LABELS[meio.tipo as TipoPagamento] || meio.tipo}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {meio.bandeira || '-'}
                  </TableCell>
                  <TableCell className="text-right font-mono">
                    {parseFloat(meio.taxa || '0').toFixed(2)}%
                  </TableCell>
                  <TableCell className="text-right font-mono">
                    {parseFloat(meio.taxa_fixa || '0') > 0
                      ? `R$ ${parseFloat(meio.taxa_fixa).toFixed(2)}`
                      : '-'}
                  </TableCell>
                  <TableCell className="text-center">
                    <Badge
                      variant={meio.d_mais === 0 ? 'default' : 'secondary'}
                      className="font-mono"
                    >
                      D+{meio.d_mais}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-center">
                    <Switch
                      checked={meio.ativo}
                      onCheckedChange={() => handleToggle(meio.id)}
                      disabled={toggleMeioPagamento.isPending}
                    />
                  </TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="opacity-0 group-hover:opacity-100 transition-opacity"
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
                  <TableCell colSpan={8} className="h-32 text-center">
                    <div className="flex flex-col items-center gap-2 text-muted-foreground">
                      <Receipt className="h-8 w-8" />
                      <p>Nenhum tipo de recebimento cadastrado.</p>
                      <Button variant="link" onClick={handleCreate}>
                        Criar o primeiro
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
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
              Tem certeza que deseja excluir o tipo de recebimento{' '}
              <strong>{meioPagamentoToDelete?.nome}</strong>? Esta ação não pode ser desfeita.
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
