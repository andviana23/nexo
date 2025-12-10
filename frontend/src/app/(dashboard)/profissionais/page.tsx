'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Profissionais
 */

import {
    ProfessionalModal,
    ProfessionalsTable,
} from '@/components/professionals';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
import { Switch } from '@/components/ui/switch';
import { useDeleteProfessional, useProfessionals, useUpdateProfessionalStatus } from '@/hooks/use-professionals';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    ListProfessionalsFilters,
    ProfessionalModalState,
    ProfessionalResponse,
    ProfessionalStatus,
    ProfessionalType,
} from '@/types/professional';
import { Briefcase, PlusIcon, SearchIcon, Users, UserX } from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';

export default function ProfissionaisPage() {
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estado de filtros
  const [filters, setFilters] = useState<ListProfessionalsFilters>({
    page: 1,
    page_size: 100, // Fetch more to handle client-side filtering effectively
    search: '',
    tipo: undefined,
    status: undefined,
  });

  const [showDismissed, setShowDismissed] = useState(false);

  // Estado do modal
  const [modalState, setModalState] = useState<ProfessionalModalState>({
    isOpen: false,
    mode: 'create',
    professional: undefined,
  });

  // Query
  const { data, isLoading, isError, refetch } = useProfessionals(filters);
  const updateStatus = useUpdateProfessionalStatus();
  const deleteProfessional = useDeleteProfessional();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Profissionais' },
    ]);
  }, [setBreadcrumbs]);

  // Derived Data
  const filteredProfessionals = useMemo(() => {
    if (!data?.data) return [];
    return data.data.filter((p) => {
      if (showDismissed) return true;
      return p.status !== 'DEMITIDO';
    });
  }, [data?.data, showDismissed]);

  const stats = useMemo(() => {
    if (!data?.data) return { total: 0, active: 0, dismissed: 0 };
    return {
      total: data.data.filter(p => p.status !== 'DEMITIDO').length,
      active: data.data.filter(p => p.status === 'ATIVO').length,
      dismissed: data.data.filter(p => p.status === 'DEMITIDO').length
    };
  }, [data?.data]);

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

  const handleDismiss = useCallback(
    (professional: ProfessionalResponse) => {
      if (confirm(`Deseja demitir ${professional.nome}? O profissional será marcado como DEMITIDO mas manterá o histórico no sistema.`)) {
        updateStatus.mutate({
          id: professional.id,
          data: { status: 'DEMITIDO' },
        });
      }
    },
    [updateStatus]
  );

  const handleDelete = useCallback(
    (professional: ProfessionalResponse) => {
      if (confirm(`⚠️ ATENÇÃO: Deseja DELETAR PERMANENTEMENTE ${professional.nome}?\n\nEsta ação NÃO pode ser desfeita e todos os dados do profissional serão removidos do sistema.`)) {
        deleteProfessional.mutate(professional.id);
      }
    },
    [deleteProfessional]
  );

  const handleCloseModal = useCallback(() => {
    setModalState({ isOpen: false, mode: 'create', professional: undefined });
  }, []);

  const handleSuccess = useCallback(() => {
    refetch();
  }, [refetch]);

  return (
    <div className="space-y-8 container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Equipe & Profissionais</h1>
          <p className="text-muted-foreground text-sm mt-1">
            Gerencie sua equipe, comissões e permissões de acesso.
          </p>
        </div>
        <Button onClick={handleOpenCreate} className="shadow-sm">
          <PlusIcon className="mr-2 h-4 w-4" />
          Novo Profissional
        </Button>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Equipe Ativa</CardTitle>
            <Users className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.active}</div>
            <p className="text-xs text-muted-foreground">Profissionais atendendo</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Registrado</CardTitle>
            <Briefcase className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-xs text-muted-foreground">Excluindo demitidos</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Demitidos</CardTitle>
            <UserX className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-muted-foreground">{stats.dismissed}</div>
            <p className="text-xs text-muted-foreground">Histórico de desligamentos</p>
          </CardContent>
        </Card>
      </div>

      {/* Filters & Content */}
      <Card className="border-border shadow-sm">
        <CardHeader className="bg-muted/40 pb-4">
          <div className="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
            <div>
              <CardTitle className="text-lg">Gerenciamento</CardTitle>
              <CardDescription>Lista completa de colaboradores.</CardDescription>
            </div>
            <div className="flex items-center space-x-2 bg-white dark:bg-zinc-950 px-3 py-1.5 rounded-md border text-sm shadow-sm">
              <Switch
                id="show-dismissed"
                checked={showDismissed}
                onCheckedChange={setShowDismissed}
              />
              <Label htmlFor="show-dismissed" className="cursor-pointer font-normal text-muted-foreground">Exibir Demitidos</Label>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 p-4">
          {/* Toolbar */}
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome ou email..."
                className="pl-9 bg-background"
                value={filters.search || ''}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>
            <Select
              value={filters.tipo || 'all'}
              onValueChange={handleFilterTipo}
            >
              <SelectTrigger className="w-[180px] bg-background">
                <SelectValue placeholder="Cargo / Tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Cargos</SelectItem>
                <SelectItem value="BARBEIRO">Barbeiro</SelectItem>
                <SelectItem value="GERENTE">Gerente</SelectItem>
                <SelectItem value="RECEPCIONISTA">Recepcionista</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={filters.status || 'all'}
              onValueChange={handleFilterStatus}
            >
              <SelectTrigger className="w-[180px] bg-background">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Status</SelectItem>
                <SelectItem value="ATIVO">Ativo</SelectItem>
                <SelectItem value="INATIVO">Inativo</SelectItem>
                <SelectItem value="AFASTADO">Afastado</SelectItem>
                {/* Demitido removed from here to rely on toggle preference, or kept for specific filter if toggle is on? 
                             Simplest: Keep it, but if user selects it without toggle, they might see nothing if we filter hard.
                             Better: Remove Demitido from here. Rely purely on Toggle + Status (Active/Inactive)
                          */}
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
                  <Skeleton className="h-6 w-20 rounded-full" />
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
              <ProfessionalsTable
                professionals={filteredProfessionals}
                onView={handleView}
                onEdit={handleEdit}
                onDismiss={handleDismiss}
                onDelete={handleDelete}
                isLoading={isLoading}
              />
            </div>
          )}
        </CardContent>
      </Card>

      <ProfessionalModal
        state={modalState}
        onClose={handleCloseModal}
        onSuccess={handleSuccess}
      />
    </div>
  );
}
