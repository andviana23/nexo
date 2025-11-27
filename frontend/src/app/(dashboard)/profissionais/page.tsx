'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Profissionais
 *
 * @page /profissionais
 * @description Listagem e gerenciamento de profissionais
 */

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
import { useProfessionals, useUpdateProfessionalStatus } from '@/hooks/use-professionals';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    ListProfessionalsFilters,
    ProfessionalModalState,
    ProfessionalResponse,
    ProfessionalStatus,
    ProfessionalType,
} from '@/types/professional';
import { PlusIcon, SearchIcon, UsersIcon } from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';

import {
    ProfessionalModal,
    ProfessionalsTable,
} from '@/components/professionals';

// =============================================================================
// COMPONENT
// =============================================================================

export default function ProfissionaisPage() {
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estado de filtros
  const [filters, setFilters] = useState<ListProfessionalsFilters>({
    page: 1,
    page_size: 20,
    search: '',
    tipo: undefined,
    status: undefined,
  });

  // Estado do modal
  const [modalState, setModalState] = useState<ProfessionalModalState>({
    isOpen: false,
    mode: 'create',
    professional: undefined,
  });

  // Query
  const { data, isLoading, isError, refetch } = useProfessionals(filters);
  const updateStatus = useUpdateProfessionalStatus();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Profissionais' },
    ]);
  }, [setBreadcrumbs]);

  // Handlers
  const handleSearch = useCallback((value: string) => {
    setFilters((prev) => ({ ...prev, search: value, page: 1 }));
  }, []);

  const handleFilterTipo = useCallback((value: string) => {
    setFilters((prev) => ({
      ...prev,
      tipo: value === 'all' ? undefined : (value as ProfessionalType),
      page: 1,
    }));
  }, []);

  const handleFilterStatus = useCallback((value: string) => {
    setFilters((prev) => ({
      ...prev,
      status: value === 'all' ? undefined : (value as ProfessionalStatus),
      page: 1,
    }));
  }, []);

  const handleOpenCreate = useCallback(() => {
    setModalState({ isOpen: true, mode: 'create', professional: undefined });
  }, []);

  const handleView = useCallback((professional: ProfessionalResponse) => {
    setModalState({ isOpen: true, mode: 'view', professional });
  }, []);

  const handleEdit = useCallback((professional: ProfessionalResponse) => {
    setModalState({ isOpen: true, mode: 'edit', professional });
  }, []);

  const handleDeactivate = useCallback(
    (professional: ProfessionalResponse) => {
      if (professional.status === 'ATIVO') {
        updateStatus.mutate({
          id: professional.id,
          data: { status: 'INATIVO' },
        });
      }
    },
    [updateStatus]
  );

  const handleCloseModal = useCallback(() => {
    setModalState({ isOpen: false, mode: 'create', professional: undefined });
  }, []);

  const handleSuccess = useCallback(() => {
    refetch();
  }, [refetch]);

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
            Profissionais
          </h1>
          <p className="text-muted-foreground">
            Gerencie barbeiros, gerentes e outros profissionais
          </p>
        </div>
        <Button onClick={handleOpenCreate}>
          <PlusIcon className="mr-2 size-4" />
          Novo Profissional
        </Button>
      </div>

      {/* Filtros */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Filtros</CardTitle>
          <CardDescription>
            Pesquise por nome, email ou filtre por tipo e status
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4">
            {/* Busca */}
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome ou email..."
                className="pl-9"
                value={filters.search || ''}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>

            {/* Tipo */}
            <Select
              value={filters.tipo || 'all'}
              onValueChange={handleFilterTipo}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os tipos</SelectItem>
                <SelectItem value="BARBEIRO">Barbeiro</SelectItem>
                <SelectItem value="GERENTE">Gerente</SelectItem>
                <SelectItem value="RECEPCIONISTA">Recepcionista</SelectItem>
                <SelectItem value="OUTRO">Outro</SelectItem>
              </SelectContent>
            </Select>

            {/* Status */}
            <Select
              value={filters.status || 'all'}
              onValueChange={handleFilterStatus}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os status</SelectItem>
                <SelectItem value="ATIVO">Ativo</SelectItem>
                <SelectItem value="INATIVO">Inativo</SelectItem>
                <SelectItem value="AFASTADO">Afastado</SelectItem>
                <SelectItem value="DEMITIDO">Demitido</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Tabela */}
      <Card>
        <CardContent className="pt-6">
          {isLoading ? (
            <ProfessionalsTableSkeleton />
          ) : isError ? (
            <div className="text-center py-12">
              <p className="text-destructive">Erro ao carregar profissionais.</p>
              <Button variant="outline" onClick={() => refetch()} className="mt-4">
                Tentar novamente
              </Button>
            </div>
          ) : (
            <ProfessionalsTable
              professionals={data?.data || []}
              onView={handleView}
              onEdit={handleEdit}
              onDeactivate={handleDeactivate}
              isLoading={isLoading}
            />
          )}
        </CardContent>
      </Card>

      {/* Paginação */}
      {data && data.meta && data.meta.total_pages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Mostrando {data.data?.length || 0} de {data.meta.total} profissionais
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
              disabled={filters.page === data.meta.total_pages}
              onClick={() => setFilters((prev) => ({ ...prev, page: (prev.page || 1) + 1 }))}
            >
              Próximo
            </Button>
          </div>
        </div>
      )}

      {/* Modal */}
      <ProfessionalModal
        state={modalState}
        onClose={handleCloseModal}
        onSuccess={handleSuccess}
      />
    </div>
  );
}

// =============================================================================
// SKELETON
// =============================================================================

function ProfessionalsTableSkeleton() {
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
