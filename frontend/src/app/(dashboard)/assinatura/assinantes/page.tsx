'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Assinantes
 *
 * @page /assinatura/assinantes
 * @description Lista de assinaturas com filtros e ações
 * Conforme FLUXO_ASSINATURA.md - FE-006
 */

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
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
    useActivePlans,
    useCancelSubscription,
    useRenewSubscription,
    useSubscriptions,
} from '@/hooks/use-subscriptions';
import { cn } from '@/lib/utils';
import { useHasRole } from '@/store/auth-store';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    ListSubscriptionsFilters,
    Subscription,
    SubscriptionModalState,
    SubscriptionStatus,
} from '@/types/subscription';
import {
    SUBSCRIPTION_STATUS_COLORS,
    SUBSCRIPTION_STATUS_LABELS,
} from '@/types/subscription';
import {
    CalendarIcon,
    EyeIcon,
    MoreHorizontalIcon,
    PhoneIcon,
    PlusIcon,
    RefreshCwIcon,
    SearchIcon,
    UsersIcon,
    XCircleIcon,
} from 'lucide-react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useCallback, useEffect, useState } from 'react';
import { SubscriptionModal } from './components/subscription-modal';

// =============================================================================
// HELPERS
// =============================================================================

function getInitials(name: string): string {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase();
}

function formatDate(dateString?: string): string {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleDateString('pt-BR');
}

function formatCurrency(value: string | number): string {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(numValue);
}

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export default function AssinantesPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const searchParams = useSearchParams();
  const filterParam = searchParams.get('filter');
  
  // RBAC - Apenas admin/gerente pode gerenciar assinaturas
  const canManageSubscriptions = useHasRole(['admin', 'gerente']);

  // Estado de filtros
  const [filters, setFilters] = useState<ListSubscriptionsFilters>(() => ({
    page: 1,
    page_size: 20,
    search: '',
    status: filterParam === 'inadimplente' ? 'INADIMPLENTE' : undefined,
    plano_id: undefined,
  }));

  // Estado do modal
  const [modalState, setModalState] = useState<SubscriptionModalState>({
    isOpen: false,
    mode: 'view',
    subscription: undefined,
  });

  // Queries
  const { data: subscriptions, isLoading, isError, refetch } = useSubscriptions(filters);
  const { data: plans } = useActivePlans();

  // Mutations
  const cancelSubscription = useCancelSubscription();
  const renewSubscription = useRenewSubscription();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Assinaturas', href: '/assinatura' },
      { label: 'Assinantes' },
    ]);
  }, [setBreadcrumbs]);

  // Handlers
  const handleSearch = useCallback((value: string) => {
    setFilters((prev) => ({ ...prev, search: value, page: 1 }));
  }, []);

  const handleFilterStatus = useCallback((value: string) => {
    setFilters((prev) => ({
      ...prev,
      status: value === 'all' ? undefined : (value as SubscriptionStatus),
      page: 1,
    }));
  }, []);

  const handleFilterPlan = useCallback((value: string) => {
    setFilters((prev) => ({
      ...prev,
      plano_id: value === 'all' ? undefined : value,
      page: 1,
    }));
  }, []);

  const handleView = useCallback((subscription: Subscription) => {
    setModalState({ isOpen: true, mode: 'view', subscription });
  }, []);

  const handleRenew = useCallback((subscription: Subscription) => {
    setModalState({ isOpen: true, mode: 'renew', subscription });
  }, []);

  const handleCancel = useCallback(
    (subscription: Subscription) => {
      if (confirm(`Deseja realmente cancelar a assinatura de "${subscription.cliente_nome}"?`)) {
        cancelSubscription.mutate(subscription.id);
      }
    },
    [cancelSubscription]
  );

  const handleCloseModal = useCallback(() => {
    setModalState({ isOpen: false, mode: 'view', subscription: undefined });
  }, []);

  const handleConfirmRenew = useCallback(
    async (data?: { codigo_transacao?: string; observacao?: string }) => {
      if (!modalState.subscription) return;
      
      await renewSubscription.mutateAsync({
        id: modalState.subscription.id,
        data,
      });
      handleCloseModal();
    },
    [modalState.subscription, renewSubscription, handleCloseModal]
  );

  // Loading state
  if (isLoading) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <div className="flex justify-between items-center">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-10 w-40" />
        </div>
        <Card>
          <CardContent className="p-6">
            <div className="space-y-4">
              {[1, 2, 3, 4, 5].map((i) => (
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
          <h1 className="text-3xl font-bold tracking-tight">Assinantes</h1>
          <p className="text-muted-foreground">
            Gerencie as assinaturas de clientes
          </p>
        </div>
        <Link href="/assinatura/nova">
          <Button>
            <PlusIcon className="mr-2 h-4 w-4" />
            Nova Assinatura
          </Button>
        </Link>
      </div>

      {/* Filtros */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Buscar por cliente..."
                className="pl-9"
                value={filters.search || ''}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>
            <Select
              value={filters.status || 'all'}
              onValueChange={handleFilterStatus}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os status</SelectItem>
                {Object.entries(SUBSCRIPTION_STATUS_LABELS).map(([value, label]) => (
                  <SelectItem key={value} value={value}>
                    {label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select
              value={filters.plano_id || 'all'}
              onValueChange={handleFilterPlan}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Plano" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os planos</SelectItem>
                {plans?.map((plan) => (
                  <SelectItem key={plan.id} value={plan.id}>
                    {plan.nome}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Tabela */}
      <Card>
        <CardHeader>
          <CardTitle>Lista de Assinantes</CardTitle>
          <CardDescription>
            {subscriptions?.length ?? 0} assinatura(s) encontrada(s)
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isError ? (
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <p className="text-destructive mb-4">
                Erro ao carregar assinaturas. Tente novamente.
              </p>
              <Button variant="outline" onClick={() => refetch()}>
                Tentar Novamente
              </Button>
            </div>
          ) : !subscriptions || subscriptions.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <UsersIcon className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-semibold mb-2">
                Nenhuma assinatura encontrada
              </h3>
              <p className="text-muted-foreground mb-4">
                Crie uma nova assinatura para um cliente
              </p>
              <Link href="/assinatura/nova">
                <Button>
                  <PlusIcon className="mr-2 h-4 w-4" />
                  Nova Assinatura
                </Button>
              </Link>
            </div>
          ) : (
            <TooltipProvider>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Cliente</TableHead>
                    <TableHead>Plano</TableHead>
                    <TableHead className="text-right">Valor</TableHead>
                    <TableHead className="text-center">Status</TableHead>
                    <TableHead>Vencimento</TableHead>
                    <TableHead className="text-right">Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {subscriptions.map((subscription) => (
                    <TableRow
                      key={subscription.id}
                      className={cn(
                        subscription.status === 'CANCELADO' && 'opacity-60'
                      )}
                    >
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <Avatar className="h-9 w-9">
                            <AvatarFallback>
                              {getInitials(subscription.cliente_nome || 'CL')}
                            </AvatarFallback>
                          </Avatar>
                          <div className="flex flex-col">
                            <span className="font-medium">
                              {subscription.cliente_nome}
                            </span>
                            {subscription.cliente_telefone && (
                              <span className="text-xs text-muted-foreground flex items-center gap-1">
                                <PhoneIcon className="h-3 w-3" />
                                {subscription.cliente_telefone}
                              </span>
                            )}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {subscription.plano_nome || 'Plano'}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right font-medium">
                        {formatCurrency(subscription.valor)}
                      </TableCell>
                      <TableCell className="text-center">
                        <Badge
                          className={cn(
                            'hover:bg-current/90',
                            SUBSCRIPTION_STATUS_COLORS[subscription.status]
                          )}
                        >
                          {SUBSCRIPTION_STATUS_LABELS[subscription.status]}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1 text-sm">
                          <CalendarIcon className="h-3 w-3 text-muted-foreground" />
                          {formatDate(subscription.data_vencimento)}
                        </div>
                      </TableCell>
                      <TableCell className="text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontalIcon className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() => handleView(subscription)}
                            >
                              <EyeIcon className="mr-2 h-4 w-4" />
                              Ver Detalhes
                            </DropdownMenuItem>
                            {(subscription.status === 'ATIVO' ||
                              subscription.status === 'INADIMPLENTE') && canManageSubscriptions && (
                              <>
                                <DropdownMenuItem
                                  onClick={() => handleRenew(subscription)}
                                >
                                  <RefreshCwIcon className="mr-2 h-4 w-4" />
                                  Registrar Renovação
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem
                                  onClick={() => handleCancel(subscription)}
                                  className="text-destructive focus:text-destructive"
                                >
                                  <XCircleIcon className="mr-2 h-4 w-4" />
                                  Cancelar Assinatura
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
      <SubscriptionModal
        state={modalState}
        onClose={handleCloseModal}
        onConfirmRenew={handleConfirmRenew}
        isLoading={renewSubscription.isPending}
      />
    </div>
  );
}
