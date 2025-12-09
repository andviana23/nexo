'use client';

import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from '@/components/ui/alert-dialog';
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
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Progress } from '@/components/ui/progress';
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
    useCreateMetaMensal,
    useDeleteMetaMensal,
    useMetasMensais,
    useUpdateMetaMensal,
} from '@/hooks/use-metas';
import {
    formatCurrency,
    formatMesAno,
    formatPercentual,
    gerarOpcoesMesAno,
    getOrigemMetaLabel,
    getStatusMetaBadgeVariant,
    getStatusMetaLabel,
    OrigemMeta,
    parseMoneyValue,
    type MetaMensalResponse,
    type SetMetaMensalRequest
} from '@/types/metas';
import { Loader2, Pencil, Plus, Target, Trash2 } from 'lucide-react';
import { useMemo, useState } from 'react';

export default function MetasMensaisPage() {
  // Estado
  const [anoFiltro, setAnoFiltro] = useState<string>(new Date().getFullYear().toString());
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedMeta, setSelectedMeta] = useState<MetaMensalResponse | null>(null);
  const [formData, setFormData] = useState<SetMetaMensalRequest>({
    mes_ano: '',
    meta_faturamento: '',
    origem: OrigemMeta.MANUAL,
  });

  // Queries & Mutations
  const { data: metas, isLoading } = useMetasMensais();
  const createMutation = useCreateMetaMensal();
  const updateMutation = useUpdateMetaMensal();
  const deleteMutation = useDeleteMetaMensal();

  // Filtrar por ano
  const metasFiltradas = useMemo(() => {
    if (!metas) return [];
    return metas
      .filter((m) => m.mes_ano.startsWith(anoFiltro))
      .sort((a, b) => b.mes_ano.localeCompare(a.mes_ano));
  }, [metas, anoFiltro]);

  // Anos disponíveis para filtro
  const anosDisponiveis = useMemo(() => {
    const currentYear = new Date().getFullYear();
    return [currentYear - 1, currentYear, currentYear + 1].map(String);
  }, []);

  // Opções de mês/ano
  const opcoesMesAno = useMemo(() => gerarOpcoesMesAno(), []);

  // Handlers
  const handleOpenCreate = () => {
    setSelectedMeta(null);
    setFormData({
      mes_ano: `${anoFiltro}-${String(new Date().getMonth() + 1).padStart(2, '0')}`,
      meta_faturamento: '',
      origem: OrigemMeta.MANUAL,
    });
    setIsModalOpen(true);
  };

  const handleOpenEdit = (meta: MetaMensalResponse) => {
    setSelectedMeta(meta);
    setFormData({
      mes_ano: meta.mes_ano,
      meta_faturamento: meta.meta_faturamento,
      origem: meta.origem as OrigemMeta,
    });
    setIsModalOpen(true);
  };

  const handleOpenDelete = (meta: MetaMensalResponse) => {
    setSelectedMeta(meta);
    setIsDeleteDialogOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.mes_ano || !formData.meta_faturamento) return;

    if (selectedMeta) {
      await updateMutation.mutateAsync({
        id: selectedMeta.id,
        payload: {
          meta_faturamento: formData.meta_faturamento,
          origem: formData.origem,
        },
      });
    } else {
      await createMutation.mutateAsync(formData);
    }

    setIsModalOpen(false);
    setSelectedMeta(null);
  };

  const handleDelete = async () => {
    if (!selectedMeta) return;
    await deleteMutation.mutateAsync(selectedMeta.id);
    setIsDeleteDialogOpen(false);
    setSelectedMeta(null);
  };

  const isPending = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-lg font-semibold">Metas Mensais de Faturamento</h2>
          <p className="text-sm text-muted-foreground">
            Defina metas de faturamento para cada mês
          </p>
        </div>
        <div className="flex gap-2">
          <Select value={anoFiltro} onValueChange={setAnoFiltro}>
            <SelectTrigger className="w-[120px]">
              <SelectValue placeholder="Ano" />
            </SelectTrigger>
            <SelectContent>
              {anosDisponiveis.map((ano) => (
                <SelectItem key={ano} value={ano}>
                  {ano}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button onClick={handleOpenCreate}>
            <Plus className="h-4 w-4 mr-2" />
            Nova Meta
          </Button>
        </div>
      </div>

      {/* Tabela */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Target className="h-5 w-5" />
            Metas de {anoFiltro}
          </CardTitle>
          <CardDescription>
            {metasFiltradas.length} meta(s) encontrada(s)
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : metasFiltradas.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <Target className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <p className="text-muted-foreground">Nenhuma meta cadastrada para {anoFiltro}</p>
              <Button variant="link" onClick={handleOpenCreate}>
                Criar primeira meta
              </Button>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Mês/Ano</TableHead>
                    <TableHead className="text-right">Meta</TableHead>
                    <TableHead className="text-right">Realizado</TableHead>
                    <TableHead className="w-[200px]">Progresso</TableHead>
                    <TableHead>Origem</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {metasFiltradas.map((meta) => {
                    const percentual = parseFloat(meta.percentual) || 0;

                    return (
                      <TableRow key={meta.id}>
                        <TableCell className="font-medium">
                          {formatMesAno(meta.mes_ano)}
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          {formatCurrency(meta.meta_faturamento)}
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          {formatCurrency(meta.realizado)}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Progress
                              value={Math.min(percentual, 100)}
                              className="flex-1"
                            />
                            <span
                              className={`text-sm font-medium min-w-[50px] text-right ${
                                percentual >= 100 ? 'text-green-600' : ''
                              }`}
                            >
                              {formatPercentual(percentual)}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline">
                            {getOrigemMetaLabel(meta.origem)}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge variant={getStatusMetaBadgeVariant(meta.status)}>
                            {getStatusMetaLabel(meta.status)}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleOpenEdit(meta)}
                            >
                              <Pencil className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleOpenDelete(meta)}
                            >
                              <Trash2 className="h-4 w-4 text-destructive" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modal Criar/Editar */}
      <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {selectedMeta ? 'Editar Meta Mensal' : 'Nova Meta Mensal'}
            </DialogTitle>
            <DialogDescription>
              {selectedMeta
                ? 'Atualize os dados da meta de faturamento'
                : 'Defina uma nova meta de faturamento mensal'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="mes_ano">Mês/Ano</Label>
              <Select
                value={formData.mes_ano}
                onValueChange={(value) =>
                  setFormData((prev) => ({ ...prev, mes_ano: value }))
                }
                disabled={!!selectedMeta}
              >
                <SelectTrigger id="mes_ano">
                  <SelectValue placeholder="Selecione o mês" />
                </SelectTrigger>
                <SelectContent>
                  {opcoesMesAno.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="meta_faturamento">Meta de Faturamento (R$)</Label>
              <Input
                id="meta_faturamento"
                type="number"
                step="0.01"
                min="0"
                placeholder="0,00"
                value={formData.meta_faturamento}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, meta_faturamento: e.target.value }))
                }
              />
              {formData.meta_faturamento && (
                <p className="text-sm text-muted-foreground">
                  {formatCurrency(parseMoneyValue(formData.meta_faturamento))}
                </p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="origem">Origem</Label>
              <Select
                value={formData.origem}
                onValueChange={(value) =>
                  setFormData((prev) => ({ ...prev, origem: value as OrigemMeta }))
                }
              >
                <SelectTrigger id="origem">
                  <SelectValue placeholder="Selecione a origem" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value={OrigemMeta.MANUAL}>Manual</SelectItem>
                  <SelectItem value={OrigemMeta.AUTOMATICA}>Automática</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsModalOpen(false)}>
              Cancelar
            </Button>
            <Button onClick={handleSubmit} disabled={isPending}>
              {isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
              {selectedMeta ? 'Salvar' : 'Criar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Dialog Confirmar Delete */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Excluir Meta</AlertDialogTitle>
            <AlertDialogDescription>
              Tem certeza que deseja excluir a meta de{' '}
              <strong>{selectedMeta && formatMesAno(selectedMeta.mes_ano)}</strong>?
              <br />
              Esta ação não pode ser desfeita.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteMutation.isPending && (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              )}
              Excluir
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
