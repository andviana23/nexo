/**
 * NEXO - Sistema de Gestão para Barbearias
 * React Query Hooks para Lista da Vez (Barber Turn)
 *
 * @module hooks/use-barber-turn
 * @description Hooks com Optimistic Updates para melhor UX
 * Baseado em FLUXO_LISTA_DA_VEZ.md
 *
 * Sistema de fila giratória para distribuição justa de clientes.
 * Menor pontuação = próximo da fila (menor número de atendimentos no mês).
 */

import {
    BarberAlreadyInListError,
    BarberNotFoundError,
    BarberPausedError,
    barberTurnService,
    NotBarberError,
    ResetFailedError,
} from '@/services/barber-turn-service';
import type {
    BarberTurnResponse,
    ListAvailableBarbersResponse,
    ListBarbersTurnResponse,
    ListHistorySummaryResponse,
    ListTurnHistoryResponse,
    RecordTurnResponse,
    ResetTurnListResponse,
    ToggleStatusResponse,
} from '@/types/barber-turn';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const barberTurnKeys = {
  all: ['barber-turn'] as const,
  lists: () => [...barberTurnKeys.all, 'list'] as const,
  list: (isActive?: boolean) => [...barberTurnKeys.lists(), { isActive }] as const,
  available: () => [...barberTurnKeys.all, 'available'] as const,
  history: () => [...barberTurnKeys.all, 'history'] as const,
  historyByMonth: (monthYear?: string) =>
    [...barberTurnKeys.history(), { monthYear }] as const,
  historySummary: () => [...barberTurnKeys.all, 'history-summary'] as const,
  stats: () => [...barberTurnKeys.all, 'stats'] as const,
};

// =============================================================================
// HOOKS DE CONSULTA (READ)
// =============================================================================

/**
 * Hook para listar barbeiros na fila da vez
 * Ordenados por pontuação (menor primeiro = próximo da fila)
 *
 * @param isActive - Filtrar por status ativo/inativo
 */
export function useBarberTurnList(isActive?: boolean) {
  return useQuery({
    queryKey: barberTurnKeys.list(isActive),
    queryFn: () => barberTurnService.list(isActive),
    staleTime: 30_000, // 30 segundos - dados mais voláteis
    refetchOnWindowFocus: true,
    refetchInterval: 60_000, // Refetch a cada 1 minuto
  });
}

/**
 * Hook para listar barbeiros disponíveis para adicionar
 * Retorna apenas profissionais do tipo BARBEIRO que não estão na fila
 */
export function useAvailableBarbers() {
  return useQuery<ListAvailableBarbersResponse>({
    queryKey: barberTurnKeys.available(),
    queryFn: () => barberTurnService.getAvailableBarbers(),
    staleTime: 60_000, // 1 minuto
  });
}

/**
 * Hook para obter histórico de atendimentos por mês
 *
 * @param monthYear - Mês/ano no formato YYYY-MM
 */
export function useTurnHistory(monthYear?: string) {
  return useQuery<ListTurnHistoryResponse>({
    queryKey: barberTurnKeys.historyByMonth(monthYear),
    queryFn: () => barberTurnService.getHistory(monthYear),
    staleTime: 5 * 60_000, // 5 minutos - histórico muda pouco
  });
}

/**
 * Hook para obter resumo dos últimos 12 meses
 */
export function useTurnHistorySummary() {
  return useQuery<ListHistorySummaryResponse>({
    queryKey: barberTurnKeys.historySummary(),
    queryFn: () => barberTurnService.getHistorySummary(),
    staleTime: 5 * 60_000, // 5 minutos
  });
}

/**
 * Hook para obter estatísticas da fila
 */
export function useTurnStats() {
  return useQuery({
    queryKey: barberTurnKeys.stats(),
    queryFn: () => barberTurnService.getStats(),
    staleTime: 30_000,
  });
}

// =============================================================================
// HOOKS DE MUTAÇÃO (CREATE, UPDATE, DELETE) - COM OPTIMISTIC UPDATES
// =============================================================================

/**
 * Hook para adicionar barbeiro à Lista da Vez
 */
export function useAddBarberToTurn() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (professionalId: string) =>
      barberTurnService.addBarber(professionalId),

    onSuccess: (newBarber) => {
      // Adiciona o novo barbeiro ao cache
      queryClient.setQueriesData<ListBarbersTurnResponse>(
        { queryKey: barberTurnKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            barbers: [...old.barbers, newBarber],
            total: old.total + 1,
            stats: {
              ...old.stats,
              total_ativos: old.stats.total_ativos + 1,
              total_geral: old.stats.total_geral + 1,
            },
          };
        }
      );

      // Invalida caches relacionados
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.lists() });
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.available() });

      toast.success(`${newBarber.professional_name} adicionado à Lista da Vez!`);
    },

    onError: (error: Error) => {
      handleBarberTurnError(error);
    },
  });
}

/**
 * Hook para registrar atendimento (incrementar pontos)
 * Usa optimistic update para feedback imediato
 */
export function useRecordTurn() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (professionalId: string) =>
      barberTurnService.recordTurn(professionalId),

    // Optimistic Update
    onMutate: async (professionalId) => {
      await queryClient.cancelQueries({ queryKey: barberTurnKeys.lists() });

      // Snapshot para rollback
      const previousData = queryClient.getQueriesData<ListBarbersTurnResponse>({
        queryKey: barberTurnKeys.lists(),
      });

      // Atualização otimista - incrementa pontos
      queryClient.setQueriesData<ListBarbersTurnResponse>(
        { queryKey: barberTurnKeys.lists() },
        (old) => {
          if (!old) return old;

          const updatedBarbers = old.barbers.map((barber) =>
            barber.professional_id === professionalId
              ? {
                  ...barber,
                  current_points: barber.current_points + 1,
                  last_turn_at: new Date().toISOString(),
                }
              : barber
          );

          // Reordena por pontos (menor primeiro)
          updatedBarbers.sort((a, b) => a.current_points - b.current_points);

          // Atualiza posições
          const withPositions = updatedBarbers.map((b, i) => ({
            ...b,
            position: i + 1,
          }));

          // Calcula novo total de pontos
          const totalPontos = withPositions.reduce(
            (sum, b) => sum + b.current_points,
            0
          );

          return {
            ...old,
            barbers: withPositions,
            stats: {
              ...old.stats,
              total_pontos_mes: totalPontos,
            },
            next_barber: withPositions[0]
              ? {
                  professional_id: withPositions[0].professional_id,
                  professional_name: withPositions[0].professional_name,
                  professional_photo: withPositions[0].professional_photo,
                  current_points: withPositions[0].current_points,
                }
              : undefined,
          };
        }
      );

      return { previousData };
    },

    // Rollback em caso de erro
    onError: (error, _professionalId, context) => {
      if (context?.previousData) {
        context.previousData.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleBarberTurnError(error);
    },

    // Sempre refetch para garantir consistência
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.lists() });
    },

    onSuccess: (response: RecordTurnResponse) => {
      toast.success(
        `Atendimento registrado para ${response.professional_name}! (${response.new_points} pts)`
      );
    },
  });
}

/**
 * Hook para alternar status ativo/inativo
 * Usa optimistic update
 */
export function useToggleBarberStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (professionalId: string) =>
      barberTurnService.toggleStatus(professionalId),

    // Optimistic Update
    onMutate: async (professionalId) => {
      await queryClient.cancelQueries({ queryKey: barberTurnKeys.lists() });

      const previousData = queryClient.getQueriesData<ListBarbersTurnResponse>({
        queryKey: barberTurnKeys.lists(),
      });

      // Alterna status otimisticamente
      queryClient.setQueriesData<ListBarbersTurnResponse>(
        { queryKey: barberTurnKeys.lists() },
        (old) => {
          if (!old) return old;

          const updatedBarbers = old.barbers.map((barber) =>
            barber.professional_id === professionalId
              ? { ...barber, is_active: !barber.is_active }
              : barber
          );

          // Calcula novos totais
          const totalAtivos = updatedBarbers.filter((b) => b.is_active).length;
          const totalPausados = updatedBarbers.filter((b) => !b.is_active).length;

          return {
            ...old,
            barbers: updatedBarbers,
            stats: {
              ...old.stats,
              total_ativos: totalAtivos,
              total_pausados: totalPausados,
            },
          };
        }
      );

      return { previousData };
    },

    onError: (error, _professionalId, context) => {
      if (context?.previousData) {
        context.previousData.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleBarberTurnError(error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.lists() });
    },

    onSuccess: (response: ToggleStatusResponse) => {
      toast.success(response.message);
    },
  });
}

/**
 * Hook para remover barbeiro da fila
 * Usa optimistic update
 */
export function useRemoveBarberFromTurn() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (professionalId: string) =>
      barberTurnService.removeBarber(professionalId),

    // Optimistic Update
    onMutate: async (professionalId) => {
      await queryClient.cancelQueries({ queryKey: barberTurnKeys.lists() });

      const previousData = queryClient.getQueriesData<ListBarbersTurnResponse>({
        queryKey: barberTurnKeys.lists(),
      });

      // Remove otimisticamente
      queryClient.setQueriesData<ListBarbersTurnResponse>(
        { queryKey: barberTurnKeys.lists() },
        (old) => {
          if (!old) return old;

          const removedBarber = old.barbers.find(
            (b) => b.professional_id === professionalId
          );
          const updatedBarbers = old.barbers.filter(
            (b) => b.professional_id !== professionalId
          );

          // Recalcula posições
          const withPositions = updatedBarbers.map((b, i) => ({
            ...b,
            position: i + 1,
          }));

          // Recalcula estatísticas
          const totalAtivos = withPositions.filter((b) => b.is_active).length;
          const totalPausados = withPositions.filter((b) => !b.is_active).length;
          const totalPontos = withPositions.reduce(
            (sum, b) => sum + b.current_points,
            0
          );

          return {
            ...old,
            barbers: withPositions,
            total: withPositions.length,
            stats: {
              total_ativos: totalAtivos,
              total_pausados: totalPausados,
              total_geral: withPositions.length,
              total_pontos_mes: totalPontos,
            },
            next_barber: withPositions[0]
              ? {
                  professional_id: withPositions[0].professional_id,
                  professional_name: withPositions[0].professional_name,
                  professional_photo: withPositions[0].professional_photo,
                  current_points: withPositions[0].current_points,
                }
              : undefined,
          };
        }
      );

      return { previousData };
    },

    onError: (error, _professionalId, context) => {
      if (context?.previousData) {
        context.previousData.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleBarberTurnError(error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.lists() });
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.available() });
    },

    onSuccess: () => {
      toast.success('Barbeiro removido da Lista da Vez.');
    },
  });
}

/**
 * Hook para executar reset mensal
 * NÃO usa optimistic update por ser operação crítica
 */
export function useResetTurnList() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (saveHistory: boolean = true) =>
      barberTurnService.resetMonthly(saveHistory),

    onSuccess: (response: ResetTurnListResponse) => {
      // Invalida todos os caches
      queryClient.invalidateQueries({ queryKey: barberTurnKeys.all });

      if (response.snapshot) {
        toast.success(
          `Reset mensal executado! ${response.snapshot.total_points_reset} pontos zerados para ${response.snapshot.total_barbers} barbeiros.`
        );
      } else {
        toast.success(response.message);
      }
    },

    onError: (error: Error) => {
      handleBarberTurnError(error);
    },
  });
}

// =============================================================================
// HOOKS AUXILIARES
// =============================================================================

/**
 * Hook para prefetch da lista de barbeiros na fila
 */
export function usePrefetchBarberTurnList() {
  const queryClient = useQueryClient();

  return () => {
    queryClient.prefetchQuery({
      queryKey: barberTurnKeys.list(),
      queryFn: () => barberTurnService.list(),
      staleTime: 30_000,
    });
  };
}

/**
 * Hook para prefetch de barbeiros disponíveis
 */
export function usePrefetchAvailableBarbers() {
  const queryClient = useQueryClient();

  return () => {
    queryClient.prefetchQuery({
      queryKey: barberTurnKeys.available(),
      queryFn: () => barberTurnService.getAvailableBarbers(),
      staleTime: 60_000,
    });
  };
}

/**
 * Hook para invalidar cache da Lista da Vez
 */
export function useInvalidateBarberTurn() {
  const queryClient = useQueryClient();

  return () => {
    queryClient.invalidateQueries({ queryKey: barberTurnKeys.all });
  };
}

/**
 * Hook para obter o próximo barbeiro da fila
 */
export function useNextBarber() {
  const { data } = useBarberTurnList(true);
  return data?.next_barber || null;
}

/**
 * Hook para verificar se um barbeiro está na lista
 */
export function useIsBarberInList(professionalId: string) {
  const { data } = useBarberTurnList();
  return data?.barbers.some((b) => b.professional_id === professionalId) || false;
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Handler centralizado de erros da Lista da Vez
 */
function handleBarberTurnError(error: Error) {
  if (error instanceof BarberNotFoundError) {
    toast.error('Barbeiro não encontrado na lista.');
    return;
  }

  if (error instanceof BarberAlreadyInListError) {
    toast.error('Este barbeiro já está na Lista da Vez.');
    return;
  }

  if (error instanceof NotBarberError) {
    toast.error('Este profissional não é do tipo Barbeiro.');
    return;
  }

  if (error instanceof BarberPausedError) {
    toast.error('Não é possível registrar atendimento para barbeiro pausado.');
    return;
  }

  if (error instanceof ResetFailedError) {
    toast.error('Erro ao executar reset mensal. Tente novamente.');
    return;
  }

  // Erro genérico
  console.error('[useBarberTurn] Error:', error);
  toast.error('Erro ao processar. Tente novamente.');
}

// =============================================================================
// TIPOS EXPORTADOS
// =============================================================================

export type { BarberTurnResponse, ListBarbersTurnResponse };
