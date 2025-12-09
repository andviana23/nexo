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
    useCreateMetaBarbeiro,
    useDeleteMetaBarbeiro,
    useMetasBarbeiro,
    useUpdateMetaBarbeiro,
} from '@/hooks/use-metas';
import {
    extenderMetaBarbeiro,
    formatCurrency,
    formatMesAno,
    formatNivelBonificacao,
    formatPercentual,
    gerarOpcoesMesAno,
    getMesAnoAtual,
    getNivelBonificacaoClass,
    parseMoneyValue,
    type MetaBarbeiroExtended,
    type MetaBarbeiroResponse,
    type SetMetaBarbeiroRequest
} from '@/types/metas';
import { Award, Loader2, Pencil, Plus, Trash2, Users } from 'lucide-react';
import { useMemo, useState } from 'react';

// TODO: Importar hook de profissionais quando disponível
// import { useProfissionais } from '@/hooks/use-profissionais';

// Mock de profissionais para desenvolvimento
const mockProfissionais = [
  { id: '1', nome: 'João Silva' },
  { id: '2', nome: 'Pedro Santos' },
  { id: '3', nome: 'Carlos Oliveira' },
];

export default function MetasBarbeirosPage() {
  // Estado
  const [mesAnoFiltro, setMesAnoFiltro] = useState<string>(getMesAnoAtual());
  const [barbeiroFiltro, setBarbeiroFiltro] = useState<string>('todos');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedMeta, setSelectedMeta] = useState<MetaBarbeiroResponse | null>(null);
  const [formData, setFormData] = useState<SetMetaBarbeiroRequest>({
    barbeiro_id: '',
    mes_ano: '',
    meta_servicos_gerais: '',
    meta_servicos_extras: '',
    meta_produtos: '',
  });

  // Queries & Mutations
  const barbeiroIdParam = barbeiroFiltro !== 'todos' ? barbeiroFiltro : undefined;
  const { data: metas, isLoading } = useMetasBarbeiro(barbeiroIdParam);
  const createMutation = useCreateMetaBarbeiro();
  const updateMutation = useUpdateMetaBarbeiro();
  const deleteMutation = useDeleteMetaBarbeiro();

  // TODO: Usar hook real
  // const { data: profissionais } = useProfissionais();
  const profissionais = mockProfissionais;

  // Filtrar por mês/ano e estender com cálculos
  const metasEstendidas = useMemo<MetaBarbeiroExtended[]>(() => {
    if (!metas) return [];
    return metas
      .filter((m) => m.mes_ano === mesAnoFiltro)
      .map(extenderMetaBarbeiro)
      .sort((a, b) => b.percentual_total - a.percentual_total);
  }, [metas, mesAnoFiltro]);

  // Opções de mês/ano
  const opcoesMesAno = useMemo(() => gerarOpcoesMesAno(), []);

  // Calcular meta total do formulário
  const metaTotalForm = useMemo(() => {
    return (
      parseMoneyValue(formData.meta_servicos_gerais) +
      parseMoneyValue(formData.meta_servicos_extras) +
      parseMoneyValue(formData.meta_produtos)
    );
  }, [formData]);

  // Handlers
  const handleOpenCreate = () => {
    setSelectedMeta(null);
    setFormData({
      barbeiro_id: '',
      mes_ano: mesAnoFiltro,
      meta_servicos_gerais: '',
      meta_servicos_extras: '',
      meta_produtos: '',
    });
    setIsModalOpen(true);
  };

  const handleOpenEdit = (meta: MetaBarbeiroResponse) => {
    setSelectedMeta(meta);
    setFormData({
      barbeiro_id: meta.barbeiro_id,
      mes_ano: meta.mes_ano,
      meta_servicos_gerais: meta.meta_servicos_gerais,
      meta_servicos_extras: meta.meta_servicos_extras,
      meta_produtos: meta.meta_produtos,
    });
    setIsModalOpen(true);
  };

  const handleOpenDelete = (meta: MetaBarbeiroResponse) => {
    setSelectedMeta(meta);
    setIsDeleteDialogOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.barbeiro_id || !formData.mes_ano) return;

    if (selectedMeta) {
      await updateMutation.mutateAsync({
        id: selectedMeta.id,
        payload: {
          meta_servicos_gerais: formData.meta_servicos_gerais,
          meta_servicos_extras: formData.meta_servicos_extras,
          meta_produtos: formData.meta_produtos,
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

  // Obter nome do barbeiro
  const getNomeBarbeiro = (id: string) => {
    const prof = profissionais.find((p) => p.id === id);
    return prof?.nome || 'Barbeiro';
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-lg font-semibold">Metas por Barbeiro</h2>
          <p className="text-sm text-muted-foreground">
            Defina metas individuais para cada profissional
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
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
          <Select value={barbeiroFiltro} onValueChange={setBarbeiroFiltro}>
            <SelectTrigger className="w-[160px]">
              <SelectValue placeholder="Barbeiro" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="todos">Todos</SelectItem>
              {profissionais.map((p) => (
                <SelectItem key={p.id} value={p.id}>
                  {p.nome}
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

      {/* Resumo do Mês */}
      <Card>
        <CardHeader>
          <CardTitle>{formatMesAno(mesAnoFiltro)}</CardTitle>
          <CardDescription>
            {metasEstendidas.length} barbeiro(s) com meta definida
          </CardDescription>
        </CardHeader>
      </Card>

      {/* Grid de Cards */}
      {isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-6 w-32" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-4 w-full mb-2" />
                <Skeleton className="h-4 w-full mb-2" />
                <Skeleton className="h-4 w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : metasEstendidas.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-10 text-center">
            <Users className="h-12 w-12 text-muted-foreground/50 mb-4" />
            <p className="text-muted-foreground">
              Nenhuma meta cadastrada para {formatMesAno(mesAnoFiltro)}
            </p>
            <Button variant="link" onClick={handleOpenCreate}>
              Criar primeira meta
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {metasEstendidas.map((meta, index) => {
            const medals = ['🥇', '🥈', '🥉'];
            const medal = index < 3 ? medals[index] : null;

            return (
              <Card key={meta.id} className="relative overflow-hidden">
                {/* Indicador de posição */}
                {medal && (
                  <div className="absolute top-2 right-2 text-2xl">
                    {medal}
                  </div>
                )}

                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-lg">
                      {meta.barbeiro_nome || getNomeBarbeiro(meta.barbeiro_id)}
                    </CardTitle>
                    <div className="flex gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => handleOpenEdit(meta)}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => handleOpenDelete(meta)}
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    </div>
                  </div>
                  {meta.nivel_bonificacao !== 'NENHUM' && (
                    <Badge className={getNivelBonificacaoClass(meta.nivel_bonificacao)}>
                      <Award className="h-3 w-3 mr-1" />
                      {formatNivelBonificacao(meta.nivel_bonificacao)}
                    </Badge>
                  )}
                </CardHeader>

                <CardContent className="space-y-4">
                  {/* Serviços Gerais */}
                  <div className="space-y-1">
                    <div className="flex justify-between text-sm">
                      <span>Serviços Gerais</span>
                      <span className="font-medium">
                        {formatPercentual(meta.percentual_servicos_gerais)}
                      </span>
                    </div>
                    <Progress
                      value={Math.min(parseFloat(meta.percentual_servicos_gerais) || 0, 100)}
                      className="h-2"
                    />
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>{formatCurrency(meta.realizado_servicos_gerais)}</span>
                      <span>de {formatCurrency(meta.meta_servicos_gerais)}</span>
                    </div>
                  </div>

                  {/* Serviços Extras */}
                  <div className="space-y-1">
                    <div className="flex justify-between text-sm">
                      <span>Serviços Extras</span>
                      <span className="font-medium">
                        {formatPercentual(meta.percentual_servicos_extras)}
                      </span>
                    </div>
                    <Progress
                      value={Math.min(parseFloat(meta.percentual_servicos_extras) || 0, 100)}
                      className="h-2"
                    />
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>{formatCurrency(meta.realizado_servicos_extras)}</span>
                      <span>de {formatCurrency(meta.meta_servicos_extras)}</span>
                    </div>
                  </div>

                  {/* Produtos */}
                  <div className="space-y-1">
                    <div className="flex justify-between text-sm">
                      <span>Produtos</span>
                      <span className="font-medium">
                        {formatPercentual(meta.percentual_produtos)}
                      </span>
                    </div>
                    <Progress
                      value={Math.min(parseFloat(meta.percentual_produtos) || 0, 100)}
                      className="h-2"
                    />
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>{formatCurrency(meta.realizado_produtos)}</span>
                      <span>de {formatCurrency(meta.meta_produtos)}</span>
                    </div>
                  </div>

                  {/* Total */}
                  <div className="pt-3 border-t">
                    <div className="flex justify-between items-center">
                      <div>
                        <p className="text-sm font-medium">Total</p>
                        <p className="text-xs text-muted-foreground">
                          {formatCurrency(meta.realizado_total)} de {formatCurrency(meta.meta_total)}
                        </p>
                      </div>
                      <div className="text-right">
                        <p
                          className={`text-2xl font-bold ${
                            meta.percentual_total >= 100 ? 'text-green-600' : ''
                          }`}
                        >
                          {formatPercentual(meta.percentual_total)}
                        </p>
                        {meta.bonus_valor && meta.bonus_valor > 0 && (
                          <p className="text-xs text-green-600">
                            Bônus: {formatCurrency(meta.bonus_valor)}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}

      {/* Modal Criar/Editar */}
      <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>
              {selectedMeta ? 'Editar Meta do Barbeiro' : 'Nova Meta de Barbeiro'}
            </DialogTitle>
            <DialogDescription>
              {selectedMeta
                ? 'Atualize os valores das metas'
                : 'Defina metas para um profissional'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="barbeiro_id">Barbeiro</Label>
              <Select
                value={formData.barbeiro_id}
                onValueChange={(value) =>
                  setFormData((prev) => ({ ...prev, barbeiro_id: value }))
                }
                disabled={!!selectedMeta}
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
              <Label htmlFor="meta_servicos_gerais">Meta Serviços Gerais (R$)</Label>
              <Input
                id="meta_servicos_gerais"
                type="number"
                step="0.01"
                min="0"
                placeholder="0,00"
                value={formData.meta_servicos_gerais}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, meta_servicos_gerais: e.target.value }))
                }
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="meta_servicos_extras">Meta Serviços Extras (R$)</Label>
              <Input
                id="meta_servicos_extras"
                type="number"
                step="0.01"
                min="0"
                placeholder="0,00"
                value={formData.meta_servicos_extras}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, meta_servicos_extras: e.target.value }))
                }
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="meta_produtos">Meta Produtos (R$)</Label>
              <Input
                id="meta_produtos"
                type="number"
                step="0.01"
                min="0"
                placeholder="0,00"
                value={formData.meta_produtos}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, meta_produtos: e.target.value }))
                }
              />
            </div>

            {/* Preview do Total */}
            <div className="p-3 rounded-lg bg-muted">
              <div className="flex justify-between items-center">
                <span className="font-medium">Meta Total</span>
                <span className="text-xl font-bold">{formatCurrency(metaTotalForm)}</span>
              </div>
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
              <strong>
                {selectedMeta &&
                  (selectedMeta.barbeiro_nome || getNomeBarbeiro(selectedMeta.barbeiro_id))}
              </strong>{' '}
              para <strong>{selectedMeta && formatMesAno(selectedMeta.mes_ano)}</strong>?
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
