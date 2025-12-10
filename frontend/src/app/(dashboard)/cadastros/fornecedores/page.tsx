'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Cadastro de Fornecedores
 */

import { PlusIcon, SearchIcon, TruckIcon } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

import { SupplierModal } from '@/components/suppliers/SupplierModal';
import { SuppliersTable } from '@/components/suppliers/SuppliersTable';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import {
  useCreateFornecedor,
  useDeleteFornecedor,
  useFornecedores,
  useUpdateFornecedor,
} from '@/hooks/use-fornecedores';
import type { SupplierFormValues } from '@/lib/validations/supplier';
import type { Fornecedor } from '@/types/fornecedor';

export default function FornecedoresPage() {
  // Estado
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingFornecedor, setEditingFornecedor] = useState<Fornecedor | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  // Hooks de dados
  const { data: fornecedores, isLoading, error } = useFornecedores();
  const createFornecedor = useCreateFornecedor();
  const updateFornecedor = useUpdateFornecedor();
  const deleteFornecedor = useDeleteFornecedor();

  // Derived state
  const filteredFornecedores = fornecedores?.filter(f => {
    const search = searchTerm.toLowerCase();
    return (
      f.razao_social.toLowerCase().includes(search) ||
      f.nome_fantasia?.toLowerCase().includes(search) ||
      f.cnpj?.includes(search) ||
      f.telefone.includes(search)
    );
  }) || [];

  // Handlers
  const handleOpenCreate = () => {
    setEditingFornecedor(null);
    setIsModalOpen(true);
  };

  const handleOpenEdit = (fornecedor: Fornecedor) => {
    setEditingFornecedor(fornecedor);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingFornecedor(null);
  };

  const handleSave = async (data: SupplierFormValues) => {
    try {
      if (editingFornecedor) {
        await updateFornecedor.mutateAsync({
          id: editingFornecedor.id,
          data: data as any // Types matching
        });
        toast.success('Fornecedor atualizado com sucesso!');
      } else {
        await createFornecedor.mutateAsync(data as any);
        toast.success('Fornecedor criado com sucesso!');
      }
      handleCloseModal();
    } catch (err) {
      console.error(err);
      toast.error(editingFornecedor ? 'Erro ao atualizar.' : 'Erro ao criar.');
    }
  };

  const handleDelete = async (id: string) => {
    if (confirm("Tem certeza que deseja excluir este fornecedor?")) {
      try {
        await deleteFornecedor.mutateAsync(id);
        toast.success('Fornecedor excluído.');
      } catch (err) {
        console.error(err);
        toast.error('Erro ao excluir.');
      }
    }
  };

  const handleToggleActive = async (fornecedor: Fornecedor) => {
    try {
      await updateFornecedor.mutateAsync({
        id: fornecedor.id,
        data: { ativo: !fornecedor.ativo }
      });
      toast.success(`Fornecedor ${!fornecedor.ativo ? 'ativado' : 'desativado'} com sucesso.`);
    } catch (err) {
      console.error(err);
      toast.error('Erro ao alterar status.');
    }
  };

  return (
    <div className="space-y-8 container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground flex items-center gap-2">
            Fornecedores
          </h1>
          <p className="text-muted-foreground text-sm mt-1">
            Gerencie os fornecedores de produtos e serviços da barbearia.
          </p>
        </div>
        <Button onClick={handleOpenCreate} className="shadow-sm">
          <PlusIcon className="mr-2 h-4 w-4" />
          Novo Fornecedor
        </Button>
      </div>

      {/* Conteúdo */}
      <Card className="border-border shadow-sm">
        <CardHeader className="bg-muted/40 pb-4">
          <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
            <div>
              <CardTitle className="text-lg flex items-center gap-2">
                <TruckIcon className="size-5 text-muted-foreground" />
                Lista de Fornecedores
              </CardTitle>
              <CardDescription>Visualize e gerencie seus parceiros comerciais.</CardDescription>
            </div>
            <div className="relative w-full sm:w-72">
              <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input
                placeholder="Buscar fornecedor..."
                className="pl-9 bg-background"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent className="p-0 sm:p-4">
          {isLoading ? (
            <div className="space-y-4 p-4">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : error ? (
            <div className="flex flex-col items-center justify-center py-12 text-destructive gap-2">
              <span>Erro ao carregar dados.</span>
              <Button variant="outline" onClick={() => window.location.reload()}>Tentar novamente</Button>
            </div>
          ) : (
            <div className="rounded-md border-0 sm:border overflow-hidden">
              <SuppliersTable
                suppliers={filteredFornecedores}
                onEdit={handleOpenEdit}
                onToggleActive={handleToggleActive}
                onDelete={handleDelete}
              />
            </div>
          )}
        </CardContent>
      </Card>

      <SupplierModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSave={handleSave}
        editingSupplier={editingFornecedor}
        isLoading={createFornecedor.isPending || updateFornecedor.isPending}
      />
    </div>
  );
}
