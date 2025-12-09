/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Cadastro de Fornecedores
 *
 * Lista, cria, edita e exclui fornecedores.
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
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
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
    useCreateFornecedor,
    useDeleteFornecedor,
    useFornecedores,
    useUpdateFornecedor,
} from '@/hooks/use-fornecedores';
import { Fornecedor, getFornecedorNome } from '@/types/fornecedor';
import { zodResolver } from '@hookform/resolvers/zod';
import {
    AlertCircle,
    ChevronRight,
    Edit,
    MoreHorizontal,
    Plus,
    Power,
    PowerOff,
    Trash2,
    Truck,
} from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

// =============================================================================
// SCHEMA DE VALIDAÇÃO
// =============================================================================

const fornecedorSchema = z.object({
  razao_social: z
    .string()
    .min(2, 'Razão social deve ter pelo menos 2 caracteres')
    .max(200, 'Razão social deve ter no máximo 200 caracteres'),
  nome_fantasia: z.string().max(200).optional(),
  cnpj: z
    .string()
    .regex(/^\d{14}$/, 'CNPJ deve ter 14 dígitos')
    .optional()
    .or(z.literal('')),
  email: z.string().email('Email inválido').optional().or(z.literal('')),
  telefone: z
    .string()
    .min(10, 'Telefone deve ter pelo menos 10 dígitos')
    .max(15, 'Telefone deve ter no máximo 15 dígitos'),
  endereco_logradouro: z.string().max(300).optional(),
  endereco_cidade: z.string().max(100).optional(),
  endereco_estado: z
    .string()
    .length(2, 'Estado deve ter 2 caracteres (UF)')
    .optional()
    .or(z.literal('')),
  endereco_cep: z
    .string()
    .regex(/^\d{8}$/, 'CEP deve ter 8 dígitos')
    .optional()
    .or(z.literal('')),
});

type FornecedorFormData = z.infer<typeof fornecedorSchema>;

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export default function FornecedoresPage() {
  // Estado
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingFornecedor, setEditingFornecedor] = useState<Fornecedor | null>(
    null
  );
  const [deletingId, setDeletingId] = useState<string | null>(null);

  // Hooks de dados
  const { data: fornecedores, isLoading, error } = useFornecedores();
  const createFornecedor = useCreateFornecedor();
  const updateFornecedor = useUpdateFornecedor();
  const deleteFornecedor = useDeleteFornecedor();

  // Form
  const form = useForm<FornecedorFormData>({
    resolver: zodResolver(fornecedorSchema),
    defaultValues: {
      razao_social: '',
      nome_fantasia: '',
      cnpj: '',
      email: '',
      telefone: '',
      endereco_logradouro: '',
      endereco_cidade: '',
      endereco_estado: '',
      endereco_cep: '',
    },
  });

  // ==========================================================================
  // HANDLERS
  // ==========================================================================

  const openCreate = () => {
    form.reset({
      razao_social: '',
      nome_fantasia: '',
      cnpj: '',
      email: '',
      telefone: '',
      endereco_logradouro: '',
      endereco_cidade: '',
      endereco_estado: '',
      endereco_cep: '',
    });
    setIsCreateOpen(true);
  };

  const openEdit = (fornecedor: Fornecedor) => {
    form.reset({
      razao_social: fornecedor.razao_social,
      nome_fantasia: fornecedor.nome_fantasia || '',
      cnpj: fornecedor.cnpj || '',
      email: fornecedor.email || '',
      telefone: fornecedor.telefone,
      endereco_logradouro: fornecedor.endereco_logradouro || '',
      endereco_cidade: fornecedor.endereco_cidade || '',
      endereco_estado: fornecedor.endereco_estado || '',
      endereco_cep: fornecedor.endereco_cep || '',
    });
    setEditingFornecedor(fornecedor);
  };

  const closeModals = () => {
    setIsCreateOpen(false);
    setEditingFornecedor(null);
    form.reset();
  };

  const onSubmitCreate = async (data: FornecedorFormData) => {
    try {
      await createFornecedor.mutateAsync({
        razao_social: data.razao_social,
        nome_fantasia: data.nome_fantasia || undefined,
        cnpj: data.cnpj || undefined,
        email: data.email || undefined,
        telefone: data.telefone,
        endereco_logradouro: data.endereco_logradouro || undefined,
        endereco_cidade: data.endereco_cidade || undefined,
        endereco_estado: data.endereco_estado || undefined,
        endereco_cep: data.endereco_cep || undefined,
      });
      toast.success('Fornecedor criado com sucesso!');
      closeModals();
    } catch (err) {
      toast.error('Erro ao criar fornecedor');
      console.error(err);
    }
  };

  const onSubmitEdit = async (data: FornecedorFormData) => {
    if (!editingFornecedor) return;

    try {
      await updateFornecedor.mutateAsync({
        id: editingFornecedor.id,
        data: {
          razao_social: data.razao_social,
          nome_fantasia: data.nome_fantasia || undefined,
          cnpj: data.cnpj || undefined,
          email: data.email || undefined,
          telefone: data.telefone,
          endereco_logradouro: data.endereco_logradouro || undefined,
          endereco_cidade: data.endereco_cidade || undefined,
          endereco_estado: data.endereco_estado || undefined,
          endereco_cep: data.endereco_cep || undefined,
        },
      });
      toast.success('Fornecedor atualizado com sucesso!');
      closeModals();
    } catch (err) {
      toast.error('Erro ao atualizar fornecedor');
      console.error(err);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteFornecedor.mutateAsync(id);
      toast.success('Fornecedor excluído com sucesso!');
      setDeletingId(null);
    } catch (err) {
      toast.error('Erro ao excluir fornecedor');
      console.error(err);
    }
  };

  const handleToggleAtivo = async (fornecedor: Fornecedor) => {
    try {
      await updateFornecedor.mutateAsync({
        id: fornecedor.id,
        data: { ativo: !fornecedor.ativo },
      });
      toast.success(
        fornecedor.ativo
          ? 'Fornecedor desativado com sucesso!'
          : 'Fornecedor ativado com sucesso!'
      );
    } catch (err) {
      toast.error('Erro ao alterar status do fornecedor');
      console.error(err);
    }
  };

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <div className="flex flex-col gap-6 p-4 md:p-6">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1 text-sm text-muted-foreground">
        <Link href="/dashboard" className="hover:text-foreground transition-colors">
          Dashboard
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span>Cadastros</span>
        <ChevronRight className="h-4 w-4" />
        <span className="text-foreground font-medium">Fornecedores</span>
      </nav>

      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Fornecedores</h1>
          <p className="text-muted-foreground">
            Gerencie os fornecedores de produtos da barbearia
          </p>
        </div>
        <Button onClick={openCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Novo Fornecedor
        </Button>
      </div>

      {/* Conteúdo */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Truck className="h-5 w-5" />
            Lista de Fornecedores
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : error ? (
            <div className="flex items-center justify-center gap-2 py-8 text-destructive">
              <AlertCircle className="h-5 w-5" />
              <span>Erro ao carregar fornecedores</span>
            </div>
          ) : !fornecedores?.length ? (
            <div className="flex flex-col items-center justify-center gap-2 py-12 text-muted-foreground">
              <Truck className="h-12 w-12 opacity-50" />
              <p>Nenhum fornecedor cadastrado</p>
              <Button variant="outline" size="sm" onClick={openCreate}>
                <Plus className="mr-2 h-4 w-4" />
                Cadastrar primeiro fornecedor
              </Button>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Nome</TableHead>
                    <TableHead>CNPJ</TableHead>
                    <TableHead>Telefone</TableHead>
                    <TableHead>Cidade/UF</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="w-[70px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {fornecedores.map((fornecedor) => (
                    <TableRow key={fornecedor.id}>
                      <TableCell className="font-medium">
                        {getFornecedorNome(fornecedor)}
                        {fornecedor.nome_fantasia &&
                          fornecedor.nome_fantasia !== fornecedor.razao_social && (
                            <span className="block text-xs text-muted-foreground">
                              {fornecedor.razao_social}
                            </span>
                          )}
                      </TableCell>
                      <TableCell>
                        {fornecedor.cnpj
                          ? fornecedor.cnpj.replace(
                              /^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$/,
                              '$1.$2.$3/$4-$5'
                            )
                          : '-'}
                      </TableCell>
                      <TableCell>
                        {fornecedor.telefone
                          ? fornecedor.telefone.replace(
                              /^(\d{2})(\d{4,5})(\d{4})$/,
                              '($1) $2-$3'
                            )
                          : '-'}
                      </TableCell>
                      <TableCell>
                        {fornecedor.endereco_cidade && fornecedor.endereco_estado
                          ? `${fornecedor.endereco_cidade}/${fornecedor.endereco_estado}`
                          : fornecedor.endereco_cidade || fornecedor.endereco_estado || '-'}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={fornecedor.ativo ? 'default' : 'secondary'}
                        >
                          {fornecedor.ativo ? 'Ativo' : 'Inativo'}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() => openEdit(fornecedor)}
                            >
                              <Edit className="mr-2 h-4 w-4" />
                              Editar
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => handleToggleAtivo(fornecedor)}
                            >
                              {fornecedor.ativo ? (
                                <>
                                  <PowerOff className="mr-2 h-4 w-4" />
                                  Desativar
                                </>
                              ) : (
                                <>
                                  <Power className="mr-2 h-4 w-4" />
                                  Ativar
                                </>
                              )}
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                              className="text-destructive"
                              onClick={() => setDeletingId(fornecedor.id)}
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Excluir
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modal de Criar */}
      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Novo Fornecedor</DialogTitle>
            <DialogDescription>
              Preencha os dados para cadastrar um novo fornecedor
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmitCreate)}
              className="space-y-4"
            >
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="razao_social"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Razão Social *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Nome da empresa"
                          {...field}
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="nome_fantasia"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Nome Fantasia</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Nome comercial"
                          {...field}
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="cnpj"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>CNPJ</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="00000000000000"
                          maxLength={14}
                          {...field}
                          onChange={(e) =>
                            field.onChange(e.target.value.replace(/\D/g, ''))
                          }
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="telefone"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Telefone *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="11999999999"
                          maxLength={15}
                          {...field}
                          onChange={(e) =>
                            field.onChange(e.target.value.replace(/\D/g, ''))
                          }
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder="contato@empresa.com"
                          {...field}
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="endereco_logradouro"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Endereço</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Rua, número, complemento, bairro"
                          rows={2}
                          {...field}
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="endereco_cidade"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Cidade</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="São Paulo"
                          {...field}
                          disabled={createFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="endereco_estado"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>UF</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="SP"
                            maxLength={2}
                            {...field}
                            onChange={(e) =>
                              field.onChange(e.target.value.toUpperCase())
                            }
                            disabled={createFornecedor.isPending}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="endereco_cep"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>CEP</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="00000000"
                            maxLength={8}
                            {...field}
                            onChange={(e) =>
                              field.onChange(e.target.value.replace(/\D/g, ''))
                            }
                            disabled={createFornecedor.isPending}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>

              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={closeModals}
                  disabled={createFornecedor.isPending}
                >
                  Cancelar
                </Button>
                <Button type="submit" disabled={createFornecedor.isPending}>
                  {createFornecedor.isPending ? 'Salvando...' : 'Salvar'}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Modal de Editar */}
      <Dialog
        open={!!editingFornecedor}
        onOpenChange={() => setEditingFornecedor(null)}
      >
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Editar Fornecedor</DialogTitle>
            <DialogDescription>
              Atualize os dados do fornecedor
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmitEdit)}
              className="space-y-4"
            >
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="razao_social"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Razão Social *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Nome da empresa"
                          {...field}
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="nome_fantasia"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Nome Fantasia</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Nome comercial"
                          {...field}
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="cnpj"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>CNPJ</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="00000000000000"
                          maxLength={14}
                          {...field}
                          onChange={(e) =>
                            field.onChange(e.target.value.replace(/\D/g, ''))
                          }
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="telefone"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Telefone *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="11999999999"
                          maxLength={15}
                          {...field}
                          onChange={(e) =>
                            field.onChange(e.target.value.replace(/\D/g, ''))
                          }
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder="contato@empresa.com"
                          {...field}
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="endereco_logradouro"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormLabel>Endereço</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Rua, número, complemento, bairro"
                          rows={2}
                          {...field}
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="endereco_cidade"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Cidade</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="São Paulo"
                          {...field}
                          disabled={updateFornecedor.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="endereco_estado"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>UF</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="SP"
                            maxLength={2}
                            {...field}
                            onChange={(e) =>
                              field.onChange(e.target.value.toUpperCase())
                            }
                            disabled={updateFornecedor.isPending}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="endereco_cep"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>CEP</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="00000000"
                            maxLength={8}
                            {...field}
                            onChange={(e) =>
                              field.onChange(e.target.value.replace(/\D/g, ''))
                            }
                            disabled={updateFornecedor.isPending}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>

              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={closeModals}
                  disabled={updateFornecedor.isPending}
                >
                  Cancelar
                </Button>
                <Button type="submit" disabled={updateFornecedor.isPending}>
                  {updateFornecedor.isPending ? 'Salvando...' : 'Salvar'}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Modal de Confirmação de Exclusão */}
      <Dialog open={!!deletingId} onOpenChange={() => setDeletingId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirmar Exclusão</DialogTitle>
            <DialogDescription>
              Tem certeza que deseja excluir este fornecedor? Esta ação não pode
              ser desfeita.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeletingId(null)}
              disabled={deleteFornecedor.isPending}
            >
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={() => deletingId && handleDelete(deletingId)}
              disabled={deleteFornecedor.isPending}
            >
              {deleteFornecedor.isPending ? 'Excluindo...' : 'Excluir'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
