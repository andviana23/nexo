/**
 * NEXO - Hooks para Meios de Pagamento
 * React Query hooks para gerenciamento de Tipos de Recebimento
 */

import { meioPagamentoService } from '@/services/meio-pagamento-service';
import {
    CreateMeioPagamentoDTO,
    MeioPagamentoFilters,
    UpdateMeioPagamentoDTO,
} from '@/types/meio-pagamento';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import { toast } from 'sonner';

export const MEIOS_PAGAMENTO_QUERY_KEY = ['meios-pagamento'];

/**
 * Hook para listar meios de pagamento
 */
export function useMeiosPagamento(filters?: MeioPagamentoFilters) {
  return useQuery({
    queryKey: [...MEIOS_PAGAMENTO_QUERY_KEY, filters],
    queryFn: () => meioPagamentoService.list(filters),
  });
}

/**
 * Hook para buscar um meio de pagamento por ID
 */
export function useMeioPagamento(id: string) {
  return useQuery({
    queryKey: [...MEIOS_PAGAMENTO_QUERY_KEY, id],
    queryFn: () => meioPagamentoService.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar meio de pagamento
 */
export function useCreateMeioPagamento() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateMeioPagamentoDTO) =>
      meioPagamentoService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: MEIOS_PAGAMENTO_QUERY_KEY,
        exact: false 
      });
      toast.success('Meio de pagamento criado com sucesso!');
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao criar meio de pagamento'
      );
    },
  });
}

/**
 * Hook para atualizar meio de pagamento
 */
export function useUpdateMeioPagamento() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMeioPagamentoDTO }) =>
      meioPagamentoService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: MEIOS_PAGAMENTO_QUERY_KEY,
        exact: false 
      });
      toast.success('Meio de pagamento atualizado com sucesso!');
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao atualizar meio de pagamento'
      );
    },
  });
}

/**
 * Hook para deletar meio de pagamento
 */
export function useDeleteMeioPagamento() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => meioPagamentoService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: MEIOS_PAGAMENTO_QUERY_KEY,
        exact: false 
      });
      toast.success('Meio de pagamento removido com sucesso!');
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao remover meio de pagamento'
      );
    },
  });
}

/**
 * Hook para alternar status ativo/inativo
 */
export function useToggleMeioPagamento() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => meioPagamentoService.toggle(id),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ 
        queryKey: MEIOS_PAGAMENTO_QUERY_KEY,
        exact: false 
      });
      toast.success(
        data.ativo
          ? 'Meio de pagamento ativado!'
          : 'Meio de pagamento desativado!'
      );
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao alterar status'
      );
    },
  });
}
