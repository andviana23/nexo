'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Planos de Assinatura
 *
 * @page /assinatura/planos
 * @description CRUD de planos de assinatura
 * Conforme FLUXO_ASSINATURA.md - FE-003, FE-004, FE-005
 */

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
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
    TooltipProvider
} from '@/components/ui/tooltip';
import {
    useActivatePlan,
    useCreatePlan,
    useDeactivatePlan,
    usePlans,
    useUpdatePlan,
} from '@/hooks/use-subscriptions';
import { cn } from '@/lib/utils';
import { useHasRole } from '@/store/auth-store';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    CreatePlanRequest,
    ListPlansFilters,
    Plan,
    PlanModalState,
    UpdatePlanRequest,
} from '@/types/subscription';
import { PERIODICITY_LABELS } from '@/types/subscription';
import {
    EditIcon,
    MoreHorizontalIcon,
    PlusIcon,
    PowerIcon,
    PowerOffIcon,
    SearchIcon,
    TagIcon,
} from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { PlanModal } from './components/plan-modal';

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export default function PlanosPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  
  // RBAC - Apenas admin/gerente pode gerenciar planos
  const canManagePlans = useHasRole(['admin', 'gerente']);

  // Estado de filtros
  const [filters, setFilters] = useState<ListPlansFilters>({
    page: 1,
    page_size: 20,
    search: '',
    ativo: undefined,
  });

  // Estado do modal
  const [modalState, setModalState] = useState<PlanModalState>({
    isOpen: false,
    mode: 'create',
    plan: undefined,
  });

  // Query
  const { data: plans, isLoading, isError, refetch } = usePlans(filters);

  // Mutations
  const createPlan = useCreatePlan();
  const updatePlan = useUpdatePlan();
  const deactivatePlan = useDeactivatePlan();
  const activatePlan = useActivatePlan();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Assinaturas', href: '/assinatura' },
      { label: 'Planos' },
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

  const handleOpenCreate = useCallback(() => {
    setModalState({ isOpen: true, mode: 'create', plan: undefined });
  }, []);

  const handleEdit = useCallback((plan: Plan) => {
    setModalState({ isOpen: true, mode: 'edit', plan });
  }, []);

  const handleView = useCallback((plan: Plan) => {
    setModalState({ isOpen: true, mode: 'view', plan });
  }, []);

  const handleToggleStatus = useCallback(
    (plan: Plan) => {
      if (plan.ativo) {
        if (confirm(`Deseja realmente desativar o plano "${plan.nome}"?`)) {
          deactivatePlan.mutate(plan.id);
        }
      } else {
        activatePlan.mutate(plan.id);
      }
    },
    [deactivatePlan, activatePlan]
  );

  const handleCloseModal = useCallback(() => {
    setModalState({ isOpen: false, mode: 'create', plan: undefined });
  }, []);

  const handleSave = useCallback(
    async (data: CreatePlanRequest | UpdatePlanRequest) => {
      if (modalState.mode === 'create') {
        await createPlan.mutateAsync(data as CreatePlanRequest);
      } else if (modalState.mode === 'edit' && modalState.plan) {
        await updatePlan.mutateAsync({
          id: modalState.plan.id,
          data: data as UpdatePlanRequest,
        });
      }
      handleCloseModal();
    },
    [modalState, createPlan, updatePlan, handleCloseModal]
  );

  // Formatadores
  const formatCurrency = (value: string | number) => {
    const numValue = typeof value === 'string' ? parseFloat(value) : value;
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
    }).format(numValue);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <div className="flex justify-between items-center">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-10 w-32" />
        </div>
        <Card>
          <CardContent className="p-6">
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-16 w-full" />
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Planos</h1>
          <p className="text-muted-foreground">
            Configure os planos de assinatura disponíveis
          </p>
        </div>
        {canManagePlans && (
          <Button onClick={handleOpenCreate}>
            <PlusIcon className="mr-2 h-4 w-4" />
            Novo Plano
          </Button>
        )}
      </div>

      {/* Filtros */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome..."
                className="pl-9"
                value={filters.search || ''}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>
            <Select
              value={
                filters.ativo === undefined
                  ? 'all'
                  : filters.ativo
                  ? 'true'
                  : 'false'
              }
              onValueChange={handleFilterStatus}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos</SelectItem>
                <SelectItem value="true">Ativos</SelectItem>
                <SelectItem value="false">Inativos</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Tabela */}
      <Card>
        <CardHeader>
          <CardTitle>Lista de Planos</CardTitle>
          <CardDescription>
            {plans?.length ?? 0} plano(s) encontrado(s)
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isError ? (
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <p className="text-destructive mb-4">
                Erro ao carregar planos. Tente novamente.
              </p>
              <Button variant="outline" onClick={() => refetch()}>
                Tentar Novamente
              </Button>
            </div>
          ) : !plans || plans.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <TagIcon className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-semibold mb-2">Nenhum plano encontrado</h3>
              <p className="text-muted-foreground mb-4">
                Crie seu primeiro plano de assinatura
              </p>
              <Button onClick={handleOpenCreate}>
                <PlusIcon className="mr-2 h-4 w-4" />
                Criar Plano
              </Button>
            </div>
          ) : (
            <TooltipProvider>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Nome</TableHead>
                    <TableHead>Periodicidade</TableHead>
                    <TableHead className="text-right">Preço</TableHead>
                    <TableHead className="text-center">Status</TableHead>
                    <TableHead className="text-right">Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {plans.map((plan) => (
                    <TableRow
                      key={plan.id}
                      className={cn(!plan.ativo && 'opacity-60')}
                    >
                      <TableCell>
                        <div className="flex flex-col">
                          <span className="font-medium">{plan.nome}</span>
                          {plan.descricao && (
                            <span className="text-xs text-muted-foreground line-clamp-1">
                              {plan.descricao}
                            </span>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {PERIODICITY_LABELS[plan.periodicidade]}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right font-medium">
                        {formatCurrency(plan.valor)}
                      </TableCell>
                      <TableCell className="text-center">
                        <Badge
                          variant={plan.ativo ? 'default' : 'secondary'}
                          className={cn(
                            plan.ativo &&
                              'bg-green-100 text-green-700 hover:bg-green-100 dark:bg-green-900/20 dark:text-green-400'
                          )}
                        >
                          {plan.ativo ? 'Ativo' : 'Inativo'}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontalIcon className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => handleView(plan)}>
                              <TagIcon className="mr-2 h-4 w-4" />
                              Visualizar
                            </DropdownMenuItem>
                            {canManagePlans && (
                              <>
                                <DropdownMenuItem onClick={() => handleEdit(plan)}>
                                  <EditIcon className="mr-2 h-4 w-4" />
                                  Editar
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem
                                  onClick={() => handleToggleStatus(plan)}
                                  className={
                                    plan.ativo
                                      ? 'text-destructive focus:text-destructive'
                                      : 'text-green-600 focus:text-green-600'
                                  }
                                >
                                  {plan.ativo ? (
                                    <>
                                      <PowerOffIcon className="mr-2 h-4 w-4" />
                                      Desativar
                                    </>
                                  ) : (
                                    <>
                                      <PowerIcon className="mr-2 h-4 w-4" />
                                      Ativar
                                    </>
                                  )}
                                </DropdownMenuItem>
                              </>
                            )}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TooltipProvider>
          )}
        </CardContent>
      </Card>

      {/* Modal */}
      <PlanModal
        state={modalState}
        onClose={handleCloseModal}
        onSave={handleSave}
        isLoading={createPlan.isPending || updatePlan.isPending}
      />
    </div>
  );
}
