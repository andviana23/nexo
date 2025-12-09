/**
 * NEXO - Sistema de Gestão para Barbearias
 * Stock Hooks
 *
 * Hooks React Query para gerenciar estado de estoque.
 */

import * as stockService from '@/services/stock-service';
import type {
    StockFilters,
    StockItemFormData,
    StockMovementFilters,
} from '@/types/stock';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { AxiosError } from 'axios';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const stockKeys = {
  all: ['stock'] as const,
  items: () => [...stockKeys.all, 'items'] as const,
  item: (id: string) => [...stockKeys.items(), id] as const,
  movements: () => [...stockKeys.all, 'movements'] as const,
  movement: (id: string) => [...stockKeys.movements(), id] as const,
  summary: () => [...stockKeys.all, 'summary'] as const,
  itemHistory: (itemId: string) => [...stockKeys.item(itemId), 'history'] as const,
};

// =============================================================================
// STOCK ITEMS
// =============================================================================

/**
 * Hook para listar itens do estoque
 */
export function useStockItems(filters?: StockFilters) {
  return useQuery({
    queryKey: [...stockKeys.items(), filters],
    queryFn: () => stockService.listStockItems(filters),
  });
}

/**
 * Hook para buscar item por ID
 */
export function useStockItem(id: string) {
  return useQuery({
    queryKey: stockKeys.item(id),
    queryFn: () => stockService.getStockItem(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar item
 */
export function useCreateStockItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: stockService.createStockItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Item criado com sucesso!');
    },
    onError: () => {
      toast.error('Erro ao criar item');
    },
  });
}

/**
 * Hook para atualizar item
 */
export function useUpdateStockItem(id: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Partial<StockItemFormData>) =>
      stockService.updateStockItem(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.item(id) });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Item atualizado com sucesso!');
    },
    onError: () => {
      toast.error('Erro ao atualizar item');
    },
  });
}

/**
 * Hook para deletar item
 */
export function useDeleteStockItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: stockService.deleteStockItem,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Item deletado com sucesso!');
    },
    onError: () => {
      toast.error('Erro ao deletar item');
    },
  });
}

/**
 * Hook para criar produto (alias para createStockItem)
 */
export function useCreateProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: stockService.createProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Produto cadastrado com sucesso!');
    },
    onError: (error: AxiosError<{ message?: string }>) => {
      const message = error.response?.data?.message || 'Erro ao cadastrar produto';
      toast.error(message);
    },
  });
}

// =============================================================================
// STOCK MOVEMENTS
// =============================================================================

/**
 * Hook para listar movimentações
 */
export function useStockMovements(filters?: StockMovementFilters) {
  return useQuery({
    queryKey: [...stockKeys.movements(), filters],
    queryFn: () => stockService.listStockMovements(filters),
  });
}

/**
 * Hook para buscar movimentação por ID
 */
export function useStockMovement(id: string) {
  return useQuery({
    queryKey: stockKeys.movement(id),
    queryFn: () => stockService.getStockMovement(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar entrada de estoque
 */
export function useCreateStockEntry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: stockService.createStockEntry,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.movements() });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Entrada registrada com sucesso!');
    },
    onError: () => {
      toast.error('Erro ao registrar entrada');
    },
  });
}

/**
 * Hook para criar entrada de estoque com múltiplos produtos
 */
export function useCreateStockEntryMultiple() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: stockService.createStockEntryMultiple,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.movements() });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Entrada registrada com sucesso!');
    },
    onError: (error: AxiosError<{ message?: string }>) => {
      const message = error.response?.data?.message || 'Erro ao registrar entrada';
      toast.error(message);
    },
  });
}

/**
 * Hook para criar saída de estoque
 */
export function useCreateStockExit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: stockService.createStockExit,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: stockKeys.items() });
      queryClient.invalidateQueries({ queryKey: stockKeys.movements() });
      queryClient.invalidateQueries({ queryKey: stockKeys.summary() });
      toast.success('Saída registrada com sucesso!');
    },
    onError: () => {
      toast.error('Erro ao registrar saída');
    },
  });
}

// =============================================================================
// ANALYTICS
// =============================================================================

/**
 * Hook para obter resumo do estoque
 */
export function useStockSummary() {
  return useQuery({
    queryKey: stockKeys.summary(),
    queryFn: stockService.getStockSummary,
  });
}

/**
 * Hook para obter histórico de um item
 */
export function useStockItemHistory(
  itemId: string,
  filters?: StockMovementFilters
) {
  return useQuery({
    queryKey: [...stockKeys.itemHistory(itemId), filters],
    queryFn: () => stockService.getStockItemHistory(itemId, filters),
    enabled: !!itemId,
  });
}
