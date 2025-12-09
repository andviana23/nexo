/**
 * NEXO - Sistema de Gestão para Barbearias
 * React Query Hooks para Assinaturas e Planos
 *
 * @module hooks/use-subscriptions
 * @description Hooks com Optimistic Updates para módulo de assinaturas
 * Baseado em FLUXO_ASSINATURA.md
 */

import { planService } from '@/services/plan-service';
import { subscriptionService } from '@/services/subscription-service';
import type {
    CreatePlanRequest,
    CreateSubscriptionRequest,
    ListPlansFilters,
    ListSubscriptionsFilters,
    RenewSubscriptionRequest,
    SubscriptionMetrics,
    UpdatePlanRequest
} from '@/types/subscription';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const planKeys = {
  all: ['plans'] as const,
  lists: () => [...planKeys.all, 'list'] as const,
  list: (filters: ListPlansFilters) => [...planKeys.lists(), filters] as const,
  details: () => [...planKeys.all, 'detail'] as const,
  detail: (id: string) => [...planKeys.details(), id] as const,
  active: () => [...planKeys.all, 'active'] as const,
};

export const subscriptionKeys = {
  all: ['subscriptions'] as const,
  lists: () => [...subscriptionKeys.all, 'list'] as const,
  list: (filters: ListSubscriptionsFilters) =>
    [...subscriptionKeys.lists(), filters] as const,
  details: () => [...subscriptionKeys.all, 'detail'] as const,
  detail: (id: string) => [...subscriptionKeys.details(), id] as const,
  metrics: () => [...subscriptionKeys.all, 'metrics'] as const,
  byStatus: (status: string) => [...subscriptionKeys.all, 'status', status] as const,
};

// =============================================================================
// HOOKS DE PLANOS (READ)
// =============================================================================

/**
 * Hook para listar planos com filtros
 */
export function usePlans(filters: ListPlansFilters = {}) {
  return useQuery({
    queryKey: planKeys.list(filters),
    queryFn: () => planService.list(filters),
    staleTime: 5 * 60_000, // 5 minutos - planos mudam pouco
  });
}

/**
 * Hook para listar apenas planos ativos
 */
export function useActivePlans() {
  return useQuery({
    queryKey: planKeys.active(),
    queryFn: () => planService.listActive(),
    staleTime: 5 * 60_000,
  });
}

/**
 * Hook para buscar um plano específico
 */
export function usePlan(id: string) {
  return useQuery({
    queryKey: planKeys.detail(id),
    queryFn: () => planService.getById(id),
    enabled: !!id,
  });
}

// =============================================================================
// HOOKS DE PLANOS (MUTATIONS)
// =============================================================================

/**
 * Hook para criar novo plano
 * RN-PLN-001: Validar campos obrigatórios
 */
export function useCreatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePlanRequest) => planService.create(data),
    onSuccess: (newPlan) => {
      // Invalidar lista de planos
      queryClient.invalidateQueries({ queryKey: planKeys.lists() });
      queryClient.invalidateQueries({ queryKey: planKeys.active() });

      toast.success('Plano criado com sucesso!', {
        description: `${newPlan.nome} - R$ ${newPlan.valor}`,
      });
    },
    onError: (error: Error) => {
      console.error('Erro ao criar plano:', error);
      toast.error('Erro ao criar plano', {
        description: error.message || 'Tente novamente mais tarde.',
      });
    },
  });
}

/**
 * Hook para atualizar plano
 * RN-PLN-002: Não afeta assinaturas já vinculadas
 */
export function useUpdatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePlanRequest }) =>
      planService.update(id, data),
    onSuccess: (updatedPlan) => {
      // Invalidar queries relacionadas
      queryClient.invalidateQueries({ queryKey: planKeys.lists() });
      queryClient.invalidateQueries({ queryKey: planKeys.active() });
      queryClient.invalidateQueries({
        queryKey: planKeys.detail(updatedPlan.id),
      });

      toast.success('Plano atualizado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('Erro ao atualizar plano:', error);
      toast.error('Erro ao atualizar plano', {
        description: error.message || 'Tente novamente mais tarde.',
      });
    },
  });
}

/**
 * Hook para desativar plano
 * RN-PLN-003: Só admin pode desativar
 */
export function useDeactivatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => planService.deactivate(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: planKeys.lists() });
      queryClient.invalidateQueries({ queryKey: planKeys.active() });
      queryClient.invalidateQueries({ queryKey: planKeys.detail(id) });

      toast.success('Plano desativado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('Erro ao desativar plano:', error);
      toast.error('Erro ao desativar plano', {
        description: error.message || 'Tente novamente mais tarde.',
      });
    },
  });
}

/**
 * Hook para reativar plano
 */
export function useActivatePlan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => planService.activate(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: planKeys.lists() });
      queryClient.invalidateQueries({ queryKey: planKeys.active() });
      queryClient.invalidateQueries({ queryKey: planKeys.detail(id) });

      toast.success('Plano reativado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('Erro ao reativar plano:', error);
      toast.error('Erro ao reativar plano', {
        description: error.message || 'Tente novamente mais tarde.',
      });
    },
  });
}

// =============================================================================
// HOOKS DE ASSINATURAS (READ)
// =============================================================================

/**
 * Hook para listar assinaturas com filtros
 * RN-SUB-001: Listagem com filtros por status, plano, busca
 */
export function useSubscriptions(filters: ListSubscriptionsFilters = {}) {
  return useQuery({
    queryKey: subscriptionKeys.list(filters),
    queryFn: () => subscriptionService.list(filters),
    staleTime: 30_000, // 30 segundos
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook para buscar uma assinatura específica
 */
export function useSubscription(id: string) {
  return useQuery({
    queryKey: subscriptionKeys.detail(id),
    queryFn: () => subscriptionService.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook para métricas de assinaturas
 * Usado nos cards do dashboard de assinaturas
 */
export function useSubscriptionMetrics() {
  return useQuery<SubscriptionMetrics>({
    queryKey: subscriptionKeys.metrics(),
    queryFn: () => subscriptionService.getMetrics(),
    staleTime: 60_000, // 1 minuto
    refetchInterval: 5 * 60_000, // Atualizar a cada 5 minutos
  });
}

/**
 * Hook para listar assinaturas ativas
 */
export function useActiveSubscriptions() {
  return useQuery({
    queryKey: subscriptionKeys.byStatus('ATIVO'),
    queryFn: () => subscriptionService.listActive(),
    staleTime: 60_000,
  });
}

/**
 * Hook para listar assinaturas inadimplentes
 */
export function useOverdueSubscriptions() {
  return useQuery({
    queryKey: subscriptionKeys.byStatus('INADIMPLENTE'),
    queryFn: () => subscriptionService.listOverdue(),
    staleTime: 60_000,
  });
}

// =============================================================================
// HOOKS DE ASSINATURAS (MUTATIONS)
// =============================================================================

/**
 * Hook para criar nova assinatura
 * RN-SUB-001: Validar se cliente já possui assinatura ativa
 * RN-SUB-002: Calcular data_fim baseado na periodicidade do plano
 */
export function useCreateSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSubscriptionRequest) => subscriptionService.create(data),
    onSuccess: () => {
      // Invalidar todas as queries de assinatura
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });

      toast.success('Assinatura criada com sucesso!', {
        description: 'O cliente agora é um assinante.',
      });
    },
    onError: (error: Error & { status?: number }) => {
      console.error('Erro ao criar assinatura:', error);

      // Tratar erros específicos
      if (error.status === 409) {
        toast.error('Cliente já possui assinatura ativa', {
          description: 'Cancele a assinatura atual antes de criar uma nova.',
        });
      } else if (error.status === 404) {
        toast.error('Cliente ou plano não encontrado', {
          description: 'Verifique os dados e tente novamente.',
        });
      } else {
        toast.error('Erro ao criar assinatura', {
          description: error.message || 'Tente novamente mais tarde.',
        });
      }
    },
  });
}

/**
 * Hook para cancelar assinatura
 * RN-CANC-001: Apenas admin/gerente pode cancelar
 * RN-CANC-002: Muda status para CANCELADO
 */
export function useCancelSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => subscriptionService.cancel(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.detail(id) });

      toast.success('Assinatura cancelada', {
        description: 'A assinatura foi cancelada com sucesso.',
      });
    },
    onError: (error: Error & { status?: number }) => {
      console.error('Erro ao cancelar assinatura:', error);

      if (error.status === 403) {
        toast.error('Sem permissão', {
          description: 'Apenas administradores podem cancelar assinaturas.',
        });
      } else {
        toast.error('Erro ao cancelar assinatura', {
          description: error.message || 'Tente novamente mais tarde.',
        });
      }
    },
  });
}

/**
 * Hook para renovar assinatura (pagamento manual)
 * RN-REN-001: Registra pagamento manual (PIX/Dinheiro)
 * RN-REN-002: Atualiza data_fim e status
 */
export function useRenewSubscription() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data = {},
    }: {
      id: string;
      data?: RenewSubscriptionRequest;
    }) => subscriptionService.renew(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.all });
      queryClient.invalidateQueries({ queryKey: subscriptionKeys.detail(id) });

      toast.success('Assinatura renovada!', {
        description: 'O pagamento foi registrado e a assinatura foi renovada.',
      });
    },
    onError: (error: Error) => {
      console.error('Erro ao renovar assinatura:', error);
      toast.error('Erro ao renovar assinatura', {
        description: error.message || 'Tente novamente mais tarde.',
      });
    },
  });
}
