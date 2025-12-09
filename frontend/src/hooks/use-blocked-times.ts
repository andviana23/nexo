/**
 * NEXO - Blocked Time Hooks
 * React Query hooks para gerenciamento de bloqueios de horário
 */

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { getErrorMessage } from '@/lib/axios';

import {
    createBlockedTime,
    deleteBlockedTime,
    listBlockedTimes,
} from '@/services/blocked-time-service';
import type {
    CreateBlockedTimeRequest,
    ListBlockedTimesRequest,
} from '@/types/appointment';

/**
 * Query key factory para bloqueios
 */
const blockedTimesKeys = {
  all: ['blocked-times'] as const,
  lists: () => [...blockedTimesKeys.all, 'list'] as const,
  list: (filters?: ListBlockedTimesRequest) =>
    [...blockedTimesKeys.lists(), filters] as const,
};

/**
 * Hook para listar bloqueios de horário
 */
export function useListBlockedTimes(filters?: ListBlockedTimesRequest) {
  return useQuery({
    queryKey: blockedTimesKeys.list(filters),
    queryFn: () => listBlockedTimes(filters),
    staleTime: 1000 * 60 * 5, // 5 minutos
  });
}

/**
 * Hook para criar bloqueio de horário
 */
export function useCreateBlockedTime() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateBlockedTimeRequest) => createBlockedTime(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: blockedTimesKeys.lists() });
      toast.success('Horário bloqueado com sucesso!');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao bloquear horário');
    },
  });
}

/**
 * Hook para deletar bloqueio de horário
 */
export function useDeleteBlockedTime() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteBlockedTime(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: blockedTimesKeys.lists() });
      toast.success('Bloqueio removido com sucesso!');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao remover bloqueio');
    },
  });
}
