'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Clientes
 */

import { CustomerModal } from '@/components/customers/CustomerModal';
import { CustomersTable } from '@/components/customers/CustomersTable';
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
  useCreateCustomer,
  useCustomers,
  useInactivateCustomer,
  useUpdateCustomer,
} from '@/hooks/use-customers';
import type {
  CustomerFormValues,
} from '@/lib/validations/customer';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
  CustomerGender,
  CustomerModalState,
  CustomerResponse,
  ListCustomersFilters,
} from '@/types/customer';
import {
  PlusIcon,
  SearchIcon,
  UserCheckIcon,
  UsersIcon
} from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';

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

  // Derived Stats (Approximate if paginated, or accurate if backend supports stats in meta)
  // Currently backend returns simple list. We can show "Total Loaded" or similar.
  // Actually the response `ListCustomersResponse` has `total`.
  // We can also calculate active/inactive from current page, but that's misleading. 
  // Ideally, use a stats endpoint. For now, just show Total count.
  const stats = useMemo(() => {
    // Dummy stats derived or static if we don't have aggregates
    // If we want real stats, we'd need a separate query.
    // For now, I'll show "Total" from `data.total`.
    return {
      total: data?.total || 0,
    }
  }, [data]);


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
    async (values: CustomerFormValues) => {
      // Map form values to Request Types. Structure is similar.
      // CustomerFormValues is inferred from Zod.
      // We need to cast or map to expected API payload.
      // Luckily, most fields match.

      // Note: We need to handle `tags`. Zod `tags` is string[] | undefined.

      const payload = {
        ...values,
        tags: values.tags || [],
      };

      if (modalState.mode === 'create') {
        await createCustomer.mutateAsync(payload as any); // using any for simplicity if types vary slightly
      } else if (modalState.mode === 'edit' && modalState.customer) {
        await updateCustomer.mutateAsync({
          id: modalState.customer.id,
          data: payload as any,
        });
      }
      handleCloseModal();
    },
    [modalState, createCustomer, updateCustomer, handleCloseModal]
  );

  return (
    <div className="space-y-8 container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground flex items-center gap-2">
            Clientes
          </h1>
          <p className="text-muted-foreground text-sm mt-1">
            Gerencie o cadastro de clientes, histórico e fidelização.
          </p>
        </div>
        <Button onClick={handleOpenCreate} className="shadow-sm">
          <PlusIcon className="mr-2 h-4 w-4" />
          Novo Cliente
        </Button>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Clientes</CardTitle>
            <UsersIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data?.total ?? 0}</div>
            <p className="text-xs text-muted-foreground">Cadastros na base</p>
          </CardContent>
        </Card>
        {/* We could add more stats if available, e.g. Active Customers */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Clientes Ativos</CardTitle>
            <UserCheckIcon className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">—</div>
            <p className="text-xs text-muted-foreground">Calculado mensalmente</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Novos este Mês</CardTitle>
            <UsersIcon className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">—</div>
            <p className="text-xs text-muted-foreground">Cadastros recentes</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters & Content */}
      <Card className="border-border shadow-sm">
        <CardHeader className="bg-muted/40 pb-4">
          <div className="flex justify-between items-center">
            <div>
              <CardTitle className="text-lg">Base de Clientes</CardTitle>
              <CardDescription>Pesquise e filtre sua lista de clientes.</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 p-4">
          {/* Toolbar */}
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome, telefone, cpf..."
                className="pl-9 bg-background"
                value={filters.search || ''}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>

            <Select
              value={filters.ativo === undefined ? 'all' : String(filters.ativo)}
              onValueChange={handleFilterStatus}
            >
              <SelectTrigger className="w-[160px] bg-background">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos</SelectItem>
                <SelectItem value="true">Ativos</SelectItem>
                <SelectItem value="false">Inativos</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={filters.genero || 'all'}
              onValueChange={handleFilterGenero}
            >
              <SelectTrigger className="w-[160px] bg-background">
                <SelectValue placeholder="Gênero" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos</SelectItem>
                <SelectItem value="M">Masculino</SelectItem>
                <SelectItem value="F">Feminino</SelectItem>
                <SelectItem value="NB">Não Binário</SelectItem>
                <SelectItem value="PNI">Prefiro ñ Inf.</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Table */}
          {isLoading ? (
            <div className="space-y-4 pt-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="flex items-center gap-4">
                  <Skeleton className="size-10 rounded-full" />
                  <div className="flex-1 space-y-2">
                    <Skeleton className="h-4 w-48" />
                    <Skeleton className="h-3 w-32" />
                  </div>
                </div>
              ))}
            </div>
          ) : isError ? (
            <div className="text-center py-12">
              <p className="text-destructive font-medium">Erro ao carregar lista.</p>
              <Button variant="outline" onClick={() => refetch()} className="mt-4">Recarregar</Button>
            </div>
          ) : (
            <div className="rounded-md border overflow-hidden">
              <CustomersTable
                customers={data?.data || []}
                onView={handleView}
                onEdit={handleEdit}
                onInactivate={handleInactivate}
                isLoading={isLoading}
              />
            </div>
          )}
        </CardContent>
      </Card>

      {/* Pagination */}
      {data && data.total > (filters.page_size || 20) && (
        <div className="flex items-center justify-between border-t pt-4">
          <p className="text-sm text-muted-foreground">
            {data.data.length} de {data.total} registros
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
