/**
 * NEXO - Sistema de Gestão para Barbearias
 * Hook de Fornecedores
 *
 * Gerencia estado e operações de fornecedores.
 * Alinhado com FornecedorResponse do backend.
 */

'use client';

import { api } from '@/lib/axios';
import {
    CreateFornecedorRequest,
    Fornecedor,
    ListFornecedoresResponse,
    UpdateFornecedorRequest,
} from '@/types/fornecedor';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const fornecedoresKeys = {
  all: ['fornecedores'] as const,
  list: () => [...fornecedoresKeys.all, 'list'] as const,
  detail: (id: string) => [...fornecedoresKeys.all, 'detail', id] as const,
};

// =============================================================================
// QUERIES
// =============================================================================

/**
 * Hook para listar fornecedores
 */
export function useFornecedores() {
  return useQuery<Fornecedor[]>({
    queryKey: fornecedoresKeys.list(),
    queryFn: async () => {
      const response = await api.get<ListFornecedoresResponse>('/fornecedores');
      return response.data.fornecedores ?? [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutos
  });
}

/**
 * Hook para buscar fornecedor por ID
 */
export function useFornecedor(id: string) {
  return useQuery<Fornecedor>({
    queryKey: fornecedoresKeys.detail(id),
    queryFn: async () => {
      const response = await api.get<Fornecedor>(`/fornecedores/${id}`);
      return response.data;
    },
    enabled: !!id,
  });
}

// =============================================================================
// MUTATIONS
// =============================================================================

/**
 * Hook para criar fornecedor
 */
export function useCreateFornecedor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateFornecedorRequest) => {
      const response = await api.post<Fornecedor>('/fornecedores', data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fornecedoresKeys.all });
    },
  });
}

/**
 * Hook para atualizar fornecedor
 */
export function useUpdateFornecedor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string;
      data: UpdateFornecedorRequest;
    }) => {
      const response = await api.put<Fornecedor>(`/fornecedores/${id}`, data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fornecedoresKeys.all });
    },
  });
}

/**
 * Hook para deletar fornecedor
 */
export function useDeleteFornecedor() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/fornecedores/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fornecedoresKeys.all });
    },
  });
}
