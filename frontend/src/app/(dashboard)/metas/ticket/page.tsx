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
    useCreateMetaTicket,
    useDeleteMetaTicket,
    useMetasTicket,
    useUpdateMetaTicket,
} from '@/hooks/use-metas';
import {
    formatCurrency,
    formatMesAno,
    formatPercentual,
    gerarOpcoesMesAno,
    getMesAnoAtual,
    parseMoneyValue,
    TipoTicketMeta,
    type MetaTicketResponse,
    type SetMetaTicketRequest
} from '@/types/metas';
import { Building2, Loader2, Pencil, Plus, Trash2, TrendingUp, User } from 'lucide-react';
import { useMemo, useState } from 'react';

// TODO: Importar hook de profissionais quando disponível
const mockProfissionais = [
  { id: '1', nome: 'João Silva' },
  { id: '2', nome: 'Pedro Santos' },
  { id: '3', nome: 'Carlos Oliveira' },
];

export default function MetasTicketPage() {
  // Estado
  const [mesAnoFiltro, setMesAnoFiltro] = useState<string>(getMesAnoAtual());
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedMeta, setSelectedMeta] = useState<MetaTicketResponse | null>(null);
  const [formData, setFormData] = useState<SetMetaTicketRequest>({
    mes_ano: '',
    tipo: TipoTicketMeta.GERAL,
    barbeiro_id: undefined,
    meta_valor: '',
  });

  // Queries & Mutations
  const { data: metas, isLoading } = useMetasTicket();
  const createMutation = useCreateMetaTicket();
  const updateMutation = useUpdateMetaTicket();
  const deleteMutation = useDeleteMetaTicket();

  const profissionais = mockProfissionais;

  // Filtrar por mês/ano
  const metasFiltradas = useMemo(() => {
    if (!metas) return { geral: null, barbeiros: [] };

    const filtradas = metas.filter((m) => m.mes_ano === mesAnoFiltro);
    const geral = filtradas.find((m) => m.tipo === TipoTicketMeta.GERAL) || null;
    const barbeiros = filtradas
      .filter((m) => m.tipo === TipoTicketMeta.BARBEIRO)
      .sort((a, b) => parseFloat(b.percentual) - parseFloat(a.percentual));

    return { geral, barbeiros };
  }, [metas, mesAnoFiltro]);

  // Opções de mês/ano
  const opcoesMesAno = useMemo(() => gerarOpcoesMesAno(), []);

  // Handlers
  const handleOpenCreate = (tipo: TipoTicketMeta = TipoTicketMeta.GERAL) => {
    setSelectedMeta(null);
    setFormData({
      mes_ano: mesAnoFiltro,
      tipo: tipo,
      barbeiro_id: tipo === TipoTicketMeta.BARBEIRO ? '' : undefined,
      meta_valor: '',
    });
    setIsModalOpen(true);
  };

  const handleOpenEdit = (meta: MetaTicketResponse) => {
    setSelectedMeta(meta);
    setFormData({
      mes_ano: meta.mes_ano,
      tipo: meta.tipo as TipoTicketMeta,
      barbeiro_id: meta.barbeiro_id,
      meta_valor: meta.meta_valor,
    });
    setIsModalOpen(true);
  };

  const handleOpenDelete = (meta: MetaTicketResponse) => {
    setSelectedMeta(meta);
    setIsDeleteDialogOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.mes_ano || !formData.meta_valor) return;
    if (formData.tipo === TipoTicketMeta.BARBEIRO && !formData.barbeiro_id) return;

    const payload: SetMetaTicketRequest = {
      ...formData,
      barbeiro_id: formData.tipo === TipoTicketMeta.GERAL ? undefined : formData.barbeiro_id,
    };

    if (selectedMeta) {
      await updateMutation.mutateAsync({
        id: selectedMeta.id,
        payload: { meta_valor: payload.meta_valor },
      });
    } else {
      await createMutation.mutateAsync(payload);
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

  // Obter nome do barbeiro
  const getNomeBarbeiro = (id?: string) => {
    if (!id) return 'N/A';
    const prof = profissionais.find((p) => p.id === id);
    return prof?.nome || 'Barbeiro';
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-lg font-semibold">Metas de Ticket Médio</h2>
          <p className="text-sm text-muted-foreground">
            Defina metas de ticket médio geral e por barbeiro
          </p>
        </div>
        <div className="flex gap-2">
          <Select value={mesAnoFiltro} onValueChange={setMesAnoFiltro}>
            <SelectTrigger className="w-[160px]">
              <SelectValue placeholder="Mês/Ano" />
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
      </div>

      {/* Meta Geral */}
      <Card className="border-2 border-primary/20">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Building2 className="h-5 w-5 text-primary" />
              <CardTitle>Meta Geral da Barbearia</CardTitle>
            </div>
            {!metasFiltradas.geral && (
              <Button size="sm" onClick={() => handleOpenCreate(TipoTicketMeta.GERAL)}>
                <Plus className="h-4 w-4 mr-2" />
                Definir Meta
              </Button>
            )}
          </div>
          <CardDescription>{formatMesAno(mesAnoFiltro)}</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-4">
              <Skeleton className="h-20 w-full" />
              <Skeleton className="h-4 w-full" />
            </div>
          ) : metasFiltradas.geral ? (
            <div className="space-y-6">
              {/* Display principal */}
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-4xl font-bold">
                    {formatCurrency(metasFiltradas.geral.ticket_medio_realizado)}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Ticket Médio Atual
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-2xl font-medium text-muted-foreground">
                    {formatCurrency(metasFiltradas.geral.meta_valor)}
                  </p>
                  <p className="text-sm text-muted-foreground">Meta</p>
                </div>
              </div>

              {/* Progress */}
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Progresso</span>
                  <span
                    className={`font-bold ${
                      parseFloat(metasFiltradas.geral.percentual) >= 100
                        ? 'text-green-600'
                        : ''
                    }`}
                  >
                    {formatPercentual(metasFiltradas.geral.percentual)}
                  </span>
                </div>
                <div className="relative">
                  <Progress
                    value={Math.min(parseFloat(metasFiltradas.geral.percentual), 100)}
                    className="h-4"
                  />
                  {/* Marcadores de referência */}
                  <div className="absolute top-0 left-1/2 h-full w-px bg-muted-foreground/30" />
                  <div className="absolute top-0 left-3/4 h-full w-px bg-muted-foreground/30" />
                </div>
                <div className="flex justify-between text-xs text-muted-foreground">
                  <span>0%</span>
                  <span>50%</span>
                  <span>75%</span>
                  <span>100%</span>
                </div>
              </div>

              {/* Ações */}
              <div className="flex justify-end gap-2 pt-2 border-t">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleOpenEdit(metasFiltradas.geral!)}
                >
                  <Pencil className="h-4 w-4 mr-2" />
                  Editar
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleOpenDelete(metasFiltradas.geral!)}
                >
                  <Trash2 className="h-4 w-4 mr-2 text-destructive" />
                  Excluir
                </Button>
              </div>
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <TrendingUp className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <p className="text-muted-foreground">
                Nenhuma meta geral definida para {formatMesAno(mesAnoFiltro)}
              </p>
              <Button variant="link" onClick={() => handleOpenCreate(TipoTicketMeta.GERAL)}>
                Definir meta agora
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Metas por Barbeiro */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <User className="h-5 w-5" />
              <CardTitle>Metas por Barbeiro</CardTitle>
            </div>
            <Button size="sm" onClick={() => handleOpenCreate(TipoTicketMeta.BARBEIRO)}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Meta
            </Button>
          </div>
          <CardDescription>
            {metasFiltradas.barbeiros.length} barbeiro(s) com meta definida
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-32" />
              ))}
            </div>
          ) : metasFiltradas.barbeiros.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <User className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <p className="text-muted-foreground">
                Nenhuma meta individual definida
              </p>
              <Button
                variant="link"
                onClick={() => handleOpenCreate(TipoTicketMeta.BARBEIRO)}
              >
                Criar primeira meta
              </Button>
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {metasFiltradas.barbeiros.map((meta) => {
                const percentual = parseFloat(meta.percentual) || 0;

                return (
                  <Card key={meta.id} className="bg-muted/30">
                    <CardHeader className="pb-2">
                      <div className="flex items-center justify-between">
                        <CardTitle className="text-base">
                          {meta.barbeiro_nome || getNomeBarbeiro(meta.barbeiro_id)}
                        </CardTitle>
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-7 w-7"
                            onClick={() => handleOpenEdit(meta)}
                          >
                            <Pencil className="h-3 w-3" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-7 w-7"
                            onClick={() => handleOpenDelete(meta)}
                          >
                            <Trash2 className="h-3 w-3 text-destructive" />
                          </Button>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <div className="flex justify-between items-end">
                        <div>
                          <p className="text-2xl font-bold">
                            {formatCurrency(meta.ticket_medio_realizado)}
                          </p>
                          <p className="text-xs text-muted-foreground">Atual</p>
                        </div>
                        <div className="text-right">
                          <p className="text-lg text-muted-foreground">
                            {formatCurrency(meta.meta_valor)}
                          </p>
                          <p className="text-xs text-muted-foreground">Meta</p>
                        </div>
                      </div>
                      <div className="space-y-1">
                        <Progress value={Math.min(percentual, 100)} className="h-2" />
                        <p
                          className={`text-sm font-medium text-right ${
                            percentual >= 100 ? 'text-green-600' : ''
                          }`}
                        >
                          {formatPercentual(percentual)}
                        </p>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modal Criar/Editar */}
      <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {selectedMeta ? 'Editar Meta de Ticket' : 'Nova Meta de Ticket Médio'}
            </DialogTitle>
            <DialogDescription>
              {selectedMeta
                ? 'Atualize o valor da meta'
                : formData.tipo === TipoTicketMeta.GERAL
                ? 'Defina a meta geral de ticket médio'
                : 'Defina a meta de ticket médio para um barbeiro'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            {!selectedMeta && (
              <>
                <div className="grid gap-2">
                  <Label htmlFor="mes_ano">Mês/Ano</Label>
                  <Select
                    value={formData.mes_ano}
                    onValueChange={(value) =>
                      setFormData((prev) => ({ ...prev, mes_ano: value }))
                    }
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
                  <Label htmlFor="tipo">Tipo</Label>
                  <Select
                    value={formData.tipo}
                    onValueChange={(value) =>
                      setFormData((prev) => ({
                        ...prev,
                        tipo: value as TipoTicketMeta,
                        barbeiro_id: value === TipoTicketMeta.GERAL ? undefined : '',
                      }))
                    }
                  >
                    <SelectTrigger id="tipo">
                      <SelectValue placeholder="Selecione o tipo" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value={TipoTicketMeta.GERAL}>
                        <div className="flex items-center gap-2">
                          <Building2 className="h-4 w-4" />
                          Geral (Barbearia)
                        </div>
                      </SelectItem>
                      <SelectItem value={TipoTicketMeta.BARBEIRO}>
                        <div className="flex items-center gap-2">
                          <User className="h-4 w-4" />
                          Por Barbeiro
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {formData.tipo === TipoTicketMeta.BARBEIRO && (
                  <div className="grid gap-2">
                    <Label htmlFor="barbeiro_id">Barbeiro</Label>
                    <Select
                      value={formData.barbeiro_id || ''}
                      onValueChange={(value) =>
                        setFormData((prev) => ({ ...prev, barbeiro_id: value }))
                      }
                    >
                      <SelectTrigger id="barbeiro_id">
                        <SelectValue placeholder="Selecione o barbeiro" />
                      </SelectTrigger>
                      <SelectContent>
                        {profissionais.map((p) => (
                          <SelectItem key={p.id} value={p.id}>
                            {p.nome}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                )}
              </>
            )}

            <div className="grid gap-2">
              <Label htmlFor="meta_valor">Meta de Ticket Médio (R$)</Label>
              <Input
                id="meta_valor"
                type="number"
                step="0.01"
                min="0"
                placeholder="0,00"
                value={formData.meta_valor}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, meta_valor: e.target.value }))
                }
              />
              {formData.meta_valor && (
                <p className="text-sm text-muted-foreground">
                  {formatCurrency(parseMoneyValue(formData.meta_valor))}
                </p>
              )}
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
              Tem certeza que deseja excluir a meta de ticket médio
              {selectedMeta?.tipo === TipoTicketMeta.GERAL
                ? ' geral'
                : ` de ${selectedMeta?.barbeiro_nome || getNomeBarbeiro(selectedMeta?.barbeiro_id)}`}
              {' '}para <strong>{selectedMeta && formatMesAno(selectedMeta.mes_ano)}</strong>?
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
