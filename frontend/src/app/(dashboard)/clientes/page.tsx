'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Clientes
 *
 * @page /clientes
 * @description Listagem e gerenciamento de clientes
 * Conforme FLUXO_CADASTROS_CLIENTE.md
 */

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
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
    useCreateCustomer,
    useCustomers,
    useInactivateCustomer,
    useUpdateCustomer,
} from '@/hooks/use-customers';
import { cn } from '@/lib/utils';
import { formatPhone, getInitials } from '@/services/customer-service';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    CreateCustomerRequest,
    CustomerGender,
    CustomerModalState,
    CustomerResponse,
    ListCustomersFilters,
    UpdateCustomerRequest,
} from '@/types/customer';
import { GENDER_LABELS, TAG_COLORS } from '@/types/customer';
import {
    CakeIcon,
    EditIcon,
    EyeIcon,
    MailIcon,
    PhoneIcon,
    PlusIcon,
    SearchIcon,
    TrashIcon,
    UsersIcon,
} from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { CustomerModal } from './components';

// =============================================================================
// COMPONENT
// =============================================================================

export default function ClientesPage() {
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estado de filtros
  const [filters, setFilters] = useState<ListCustomersFilters>({
    page: 1,
    page_size: 20,
    search: '',
    ativo: undefined,
    genero: undefined,
  });

  // Estado do modal
  const [modalState, setModalState] = useState<CustomerModalState>({
    isOpen: false,
    mode: 'create',
    customer: undefined,
  });

  // Query
  const { data, isLoading, isError, refetch } = useCustomers(filters);
  
  // Mutations
  const createCustomer = useCreateCustomer();
  const updateCustomer = useUpdateCustomer();
  const inactivateCustomer = useInactivateCustomer();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Cadastros' },
      { label: 'Clientes' },
    ]);
  }, [setBreadcrumbs]);

  // Handlers
  const handleSearch = useCallback((value: string) => {
    setFilters((prev) => ({ ...prev, search: value, page: 1 }));
  }, []);

  const handleFilterStatus = useCallback((value: string) => {
    setFilters((prev) => ({
      ...prev,
      ativo: value === 'all' ? undefined : value === 'true',
      page: 1,
    }));
  }, []);

  const handleFilterGenero = useCallback((value: string) => {
    setFilters((prev) => ({
      ...prev,
      genero: value === 'all' ? undefined : (value as CustomerGender),
      page: 1,
    }));
  }, []);

  const handleOpenCreate = useCallback(() => {
    setModalState({ isOpen: true, mode: 'create', customer: undefined });
  }, []);

  const handleView = useCallback((customer: CustomerResponse) => {
    setModalState({ isOpen: true, mode: 'view', customer });
  }, []);

  const handleEdit = useCallback((customer: CustomerResponse) => {
    setModalState({ isOpen: true, mode: 'edit', customer });
  }, []);

  const handleInactivate = useCallback(
    (customer: CustomerResponse) => {
      if (customer.ativo) {
        if (confirm(`Deseja realmente inativar o cliente "${customer.nome}"?`)) {
          inactivateCustomer.mutate(customer.id);
        }
      }
    },
    [inactivateCustomer]
  );

  const handleCloseModal = useCallback(() => {
    setModalState({ isOpen: false, mode: 'create', customer: undefined });
  }, []);

  const handleSave = useCallback(
    async (data: CreateCustomerRequest | UpdateCustomerRequest) => {
      if (modalState.mode === 'create') {
        await createCustomer.mutateAsync(data as CreateCustomerRequest);
      } else if (modalState.mode === 'edit' && modalState.customer) {
        await updateCustomer.mutateAsync({
          id: modalState.customer.id,
          data: data as UpdateCustomerRequest,
        });
      }
      handleCloseModal();
    },
    [modalState, createCustomer, updateCustomer, handleCloseModal]
  );

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <UsersIcon className="size-8" />
            Clientes
          </h1>
          <p className="text-muted-foreground">
            Gerencie o cadastro de clientes da barbearia
          </p>
        </div>
        <Button onClick={handleOpenCreate}>
          <PlusIcon className="mr-2 size-4" />
          Novo Cliente
        </Button>
      </div>

      {/* Filtros */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Filtros</CardTitle>
          <CardDescription>
            Pesquise por nome, telefone ou CPF
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4">
            {/* Busca */}
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome, telefone ou CPF..."
                className="pl-9"
                value={filters.search || ''}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>

            {/* Status */}
            <Select
              value={filters.ativo === undefined ? 'all' : String(filters.ativo)}
              onValueChange={handleFilterStatus}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os status</SelectItem>
                <SelectItem value="true">Ativos</SelectItem>
                <SelectItem value="false">Inativos</SelectItem>
              </SelectContent>
            </Select>

            {/* Gênero */}
            <Select
              value={filters.genero || 'all'}
              onValueChange={handleFilterGenero}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Gênero" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os gêneros</SelectItem>
                <SelectItem value="M">Masculino</SelectItem>
                <SelectItem value="F">Feminino</SelectItem>
                <SelectItem value="NB">Não Binário</SelectItem>
                <SelectItem value="PNI">Prefiro não informar</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Tabela */}
      <Card>
        <CardContent className="pt-6">
          {isLoading ? (
            <CustomersTableSkeleton />
          ) : isError ? (
            <div className="text-center py-12">
              <p className="text-destructive">Erro ao carregar clientes.</p>
              <Button variant="outline" onClick={() => refetch()} className="mt-4">
                Tentar novamente
              </Button>
            </div>
          ) : data?.data.length === 0 ? (
            <div className="text-center py-12">
              <UsersIcon className="mx-auto size-12 text-muted-foreground/50" />
              <h3 className="mt-4 text-lg font-semibold">Nenhum cliente encontrado</h3>
              <p className="text-muted-foreground mt-1">
                {filters.search
                  ? 'Tente ajustar os filtros de busca'
                  : 'Comece cadastrando seu primeiro cliente'}
              </p>
              {!filters.search && (
                <Button onClick={handleOpenCreate} className="mt-4">
                  <PlusIcon className="mr-2 size-4" />
                  Cadastrar Cliente
                </Button>
              )}
            </div>
          ) : (
            <CustomersTable
              customers={data?.data || []}
              onView={handleView}
              onEdit={handleEdit}
              onInactivate={handleInactivate}
            />
          )}
        </CardContent>
      </Card>

      {/* Paginação */}
      {data && data.total > (filters.page_size || 20) && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Mostrando {data.data.length} de {data.total} clientes
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={filters.page === 1}
              onClick={() => setFilters((prev) => ({ ...prev, page: (prev.page || 1) - 1 }))}
            >
              Anterior
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={(filters.page || 1) * (filters.page_size || 20) >= data.total}
              onClick={() => setFilters((prev) => ({ ...prev, page: (prev.page || 1) + 1 }))}
            >
              Próximo
            </Button>
          </div>
        </div>
      )}

      {/* Modal */}
      <CustomerModal
        state={modalState}
        onClose={handleCloseModal}
        onSave={handleSave}
        isLoading={createCustomer.isPending || updateCustomer.isPending}
      />
    </div>
  );
}

// =============================================================================
// TABELA DE CLIENTES
// =============================================================================

interface CustomersTableProps {
  customers: CustomerResponse[];
  onView: (customer: CustomerResponse) => void;
  onEdit: (customer: CustomerResponse) => void;
  onInactivate: (customer: CustomerResponse) => void;
}

function CustomersTable({
  customers,
  onView,
  onEdit,
  onInactivate,
}: CustomersTableProps) {
  // Filtrar possíveis valores undefined/null do optimistic update
  const safeCustomers = customers.filter((c): c is CustomerResponse => 
    c != null && typeof c === 'object' && 'id' in c && 'ativo' in c
  );

  return (
    <TooltipProvider>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Cliente</TableHead>
            <TableHead>Contato</TableHead>
            <TableHead>Gênero</TableHead>
            <TableHead>Tags</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="text-right">Ações</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {safeCustomers.map((customer) => (
            <TableRow key={customer.id}>
              {/* Avatar + Nome */}
              <TableCell>
                <div className="flex items-center gap-3">
                  <div
                    className={cn(
                      'flex size-10 items-center justify-center rounded-full text-white font-medium',
                      customer.ativo ? 'bg-primary' : 'bg-muted-foreground'
                    )}
                  >
                    {getInitials(customer.nome)}
                  </div>
                  <div>
                    <p className="font-medium">{customer.nome}</p>
                    {customer.data_nascimento && (
                      <p className="text-xs text-muted-foreground flex items-center gap-1">
                        <CakeIcon className="size-3" />
                        {new Date(customer.data_nascimento).toLocaleDateString('pt-BR')}
                      </p>
                    )}
                  </div>
                </div>
              </TableCell>

              {/* Contato */}
              <TableCell>
                <div className="space-y-1">
                  <p className="flex items-center gap-1 text-sm">
                    <PhoneIcon className="size-3 text-muted-foreground" />
                    {formatPhone(customer.telefone)}
                  </p>
                  {customer.email && (
                    <p className="flex items-center gap-1 text-xs text-muted-foreground">
                      <MailIcon className="size-3" />
                      {customer.email}
                    </p>
                  )}
                </div>
              </TableCell>

              {/* Gênero */}
              <TableCell>
                {customer.genero ? (
                  <span className="text-sm text-muted-foreground">
                    {GENDER_LABELS[customer.genero]}
                  </span>
                ) : (
                  <span className="text-sm text-muted-foreground/50">—</span>
                )}
              </TableCell>

              {/* Tags */}
              <TableCell>
                <div className="flex flex-wrap gap-1">
                  {customer.tags.length > 0 ? (
                    customer.tags.slice(0, 2).map((tag) => (
                      <Badge
                        key={tag}
                        variant="secondary"
                        className={cn('text-xs', TAG_COLORS[tag] || '')}
                      >
                        {tag}
                      </Badge>
                    ))
                  ) : (
                    <span className="text-sm text-muted-foreground/50">—</span>
                  )}
                  {customer.tags.length > 2 && (
                    <Tooltip>
                      <TooltipTrigger>
                        <Badge variant="outline" className="text-xs">
                          +{customer.tags.length - 2}
                        </Badge>
                      </TooltipTrigger>
                      <TooltipContent>
                        {customer.tags.slice(2).join(', ')}
                      </TooltipContent>
                    </Tooltip>
                  )}
                </div>
              </TableCell>

              {/* Status */}
              <TableCell>
                <Badge variant={customer.ativo ? 'default' : 'secondary'}>
                  {customer.ativo ? 'Ativo' : 'Inativo'}
                </Badge>
              </TableCell>

              {/* Ações */}
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onView(customer)}
                      >
                        <EyeIcon className="size-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Ver detalhes</TooltipContent>
                  </Tooltip>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onEdit(customer)}
                      >
                        <EditIcon className="size-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Editar</TooltipContent>
                  </Tooltip>

                  {customer.ativo && (
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => onInactivate(customer)}
                          className="text-destructive hover:text-destructive"
                        >
                          <TrashIcon className="size-4" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Inativar</TooltipContent>
                    </Tooltip>
                  )}
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TooltipProvider>
  );
}

// =============================================================================
// SKELETON
// =============================================================================

function CustomersTableSkeleton() {
  return (
    <div className="space-y-4">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="flex items-center gap-4">
          <Skeleton className="size-10 rounded-full" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-48" />
            <Skeleton className="h-3 w-32" />
          </div>
          <Skeleton className="h-6 w-20 rounded-full" />
          <Skeleton className="h-6 w-16 rounded-full" />
        </div>
      ))}
    </div>
  );
}
