/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Categorias de Produtos
 *
 * CRUD de categorias customizáveis para produtos do estoque.
 * Módulo de Estoque v1.0
 */

'use client';

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
import { Textarea } from '@/components/ui/textarea';
import {
    useCategoriasProdutos,
    useCreateCategoriaProduto,
    useDeleteCategoriaProduto,
    useToggleCategoriaProduto,
    useUpdateCategoriaProduto,
} from '@/hooks/use-categorias-produtos';
import { useBreadcrumbs } from '@/store/ui-store';
import {
    CATEGORIA_CORES,
    CATEGORIA_ICONES,
    CENTRO_CUSTO_OPTIONS,
    CategoriaProduto,
    CentroCusto,
    CreateCategoriaProdutoRequest,
    UpdateCategoriaProdutoRequest,
} from '@/types/categoria-produto';
import {
    ArrowLeft,
    Edit,
    Palette,
    Plus,
    Power,
    Trash2,
} from 'lucide-react';
import Link from 'next/link';
import { useEffect, useState } from 'react';

// Importar ícones dinâmicos
import * as LucideIcons from 'lucide-react';

// Componente para renderizar ícone dinamicamente
function DynamicIcon({ name, className, style }: { name: string; className?: string; style?: React.CSSProperties }) {
  const iconName = name.charAt(0).toUpperCase() + name.slice(1).replace(/-([a-z])/g, (g) => g[1].toUpperCase());
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const Icon = (LucideIcons as any)[iconName];
  if (!Icon) return <LucideIcons.Package className={className} style={style} />;
  return <Icon className={className} style={style} />;
}

export default function CategoriasProdutosPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [modalOpen, setModalOpen] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<CategoriaProduto | null>(null);
  const [categoryToDelete, setCategoryToDelete] = useState<CategoriaProduto | null>(null);

  // Form state
  const [formData, setFormData] = useState<CreateCategoriaProdutoRequest>({
    nome: '',
    descricao: '',
    cor: '#6B7280',
    icone: 'package',
    centro_custo: 'CMV',
  });

  // Queries e Mutations
  const { data: categorias, isLoading } = useCategoriasProdutos();
  const createMutation = useCreateCategoriaProduto();
  const updateMutation = useUpdateCategoriaProduto();
  const deleteMutation = useDeleteCategoriaProduto();
  const toggleMutation = useToggleCategoriaProduto();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Estoque', href: '/estoque' },
      { label: 'Categorias de Produtos' },
    ]);
  }, [setBreadcrumbs]);

  const handleOpenCreate = () => {
    setEditingCategory(null);
    setFormData({
      nome: '',
      descricao: '',
      cor: '#6B7280',
      icone: 'package',
      centro_custo: 'CMV',
    });
    setModalOpen(true);
  };

  const handleOpenEdit = (categoria: CategoriaProduto) => {
    setEditingCategory(categoria);
    setFormData({
      nome: categoria.nome,
      descricao: categoria.descricao || '',
      cor: categoria.cor || '#6B7280',
      icone: categoria.icone || 'package',
      centro_custo: categoria.centro_custo || 'CMV',
    });
    setModalOpen(true);
  };

  const handleOpenDelete = (categoria: CategoriaProduto) => {
    setCategoryToDelete(categoria);
    setDeleteModalOpen(true);
  };

  const handleSubmit = async () => {
    if (editingCategory) {
      const updateData: UpdateCategoriaProdutoRequest = {
        ...formData,
        ativa: editingCategory.ativa,
      };
      await updateMutation.mutateAsync({ id: editingCategory.id, data: updateData });
    } else {
      await createMutation.mutateAsync(formData);
    }
    setModalOpen(false);
  };

  const handleDelete = async () => {
    if (categoryToDelete) {
      await deleteMutation.mutateAsync(categoryToDelete.id);
      setDeleteModalOpen(false);
      setCategoryToDelete(null);
    }
  };

  const handleToggle = async (categoria: CategoriaProduto) => {
    await toggleMutation.mutateAsync(categoria.id);
  };

  const isSubmitting = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/estoque">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Categorias de Produtos</h1>
            <p className="text-muted-foreground">
              Gerencie as categorias para organização do estoque
            </p>
          </div>
        </div>
        <Button onClick={handleOpenCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Categoria
        </Button>
      </div>

      {/* Lista de Categorias */}
      {isLoading ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <Skeleton key={i} className="h-40" />
          ))}
        </div>
      ) : categorias && categorias.length > 0 ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {categorias.map((categoria) => (
            <Card key={categoria.id} className={!categoria.ativa ? 'opacity-60' : ''}>
              <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div
                      className="flex h-10 w-10 items-center justify-center rounded-lg"
                      style={{ backgroundColor: categoria.cor + '20' }}
                    >
                      <DynamicIcon
                        name={categoria.icone}
                        className="h-5 w-5"
                        style={{ color: categoria.cor }}
                      />
                    </div>
                    <div>
                      <CardTitle className="text-lg">{categoria.nome}</CardTitle>
                      <Badge variant={categoria.ativa ? 'default' : 'secondary'} className="mt-1">
                        {categoria.ativa ? 'Ativa' : 'Inativa'}
                      </Badge>
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
                  {categoria.descricao || 'Sem descrição'}
                </p>
                <div className="flex items-center justify-between">
                  <Badge variant="outline">
                    {CENTRO_CUSTO_OPTIONS.find((c) => c.value === categoria.centro_custo)?.label ||
                      categoria.centro_custo}
                  </Badge>
                  <div className="flex items-center gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleToggle(categoria)}
                      title={categoria.ativa ? 'Desativar' : 'Ativar'}
                    >
                      <Power className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleOpenEdit(categoria)}
                      title="Editar"
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleOpenDelete(categoria)}
                      title="Excluir"
                      className="text-destructive hover:text-destructive"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Palette className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">Nenhuma categoria cadastrada</h3>
            <p className="text-muted-foreground text-center mb-4">
              Crie categorias para organizar os produtos do seu estoque
            </p>
            <Button onClick={handleOpenCreate}>
              <Plus className="mr-2 h-4 w-4" />
              Criar Primeira Categoria
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Modal de Criar/Editar */}
      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              {editingCategory ? 'Editar Categoria' : 'Nova Categoria'}
            </DialogTitle>
            <DialogDescription>
              {editingCategory
                ? 'Atualize as informações da categoria'
                : 'Preencha os dados para criar uma nova categoria'}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="nome">Nome *</Label>
              <Input
                id="nome"
                placeholder="Ex: Pomadas"
                value={formData.nome}
                onChange={(e) => setFormData({ ...formData, nome: e.target.value })}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="descricao">Descrição</Label>
              <Textarea
                id="descricao"
                placeholder="Descrição da categoria..."
                value={formData.descricao}
                onChange={(e) => setFormData({ ...formData, descricao: e.target.value })}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Cor</Label>
                <Select
                  value={formData.cor}
                  onValueChange={(value) => setFormData({ ...formData, cor: value })}
                >
                  <SelectTrigger>
                    <SelectValue>
                      <div className="flex items-center gap-2">
                        <div
                          className="h-4 w-4 rounded-full"
                          style={{ backgroundColor: formData.cor }}
                        />
                        {CATEGORIA_CORES.find((c) => c.value === formData.cor)?.label}
                      </div>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    {CATEGORIA_CORES.map((cor) => (
                      <SelectItem key={cor.value} value={cor.value}>
                        <div className="flex items-center gap-2">
                          <div
                            className="h-4 w-4 rounded-full"
                            style={{ backgroundColor: cor.value }}
                          />
                          {cor.label}
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label>Ícone</Label>
                <Select
                  value={formData.icone}
                  onValueChange={(value) => setFormData({ ...formData, icone: value })}
                >
                  <SelectTrigger>
                    <SelectValue>
                      <div className="flex items-center gap-2">
                        <DynamicIcon name={formData.icone || 'package'} className="h-4 w-4" />
                        {CATEGORIA_ICONES.find((i) => i.value === formData.icone)?.label}
                      </div>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    {CATEGORIA_ICONES.map((icone) => (
                      <SelectItem key={icone.value} value={icone.value}>
                        <div className="flex items-center gap-2">
                          <DynamicIcon name={icone.value} className="h-4 w-4" />
                          {icone.label}
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-2">
              <Label>Centro de Custo (DRE)</Label>
              <Select
                value={formData.centro_custo}
                onValueChange={(value: CentroCusto) =>
                  setFormData({ ...formData, centro_custo: value })
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CENTRO_CUSTO_OPTIONS.map((cc) => (
                    <SelectItem key={cc.value} value={cc.value}>
                      <div>
                        <span className="font-medium">{cc.label}</span>
                        <span className="text-xs text-muted-foreground ml-2">
                          {cc.description}
                        </span>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setModalOpen(false)}>
              Cancelar
            </Button>
            <Button
              onClick={handleSubmit}
              disabled={!formData.nome || isSubmitting}
            >
              {isSubmitting ? 'Salvando...' : editingCategory ? 'Atualizar' : 'Criar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Modal de Confirmação de Exclusão */}
      <Dialog open={deleteModalOpen} onOpenChange={setDeleteModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Excluir Categoria</DialogTitle>
            <DialogDescription>
              Tem certeza que deseja excluir a categoria &quot;{categoryToDelete?.nome}&quot;?
              Esta ação não pode ser desfeita.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteModalOpen(false)}>
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
    </div>
  );
}
