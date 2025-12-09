/**
 * NEXO - Sistema de Gestão para Barbearias
 * React Query Hooks para Precificação
 *
 * @module hooks/use-pricing
 * @description Hooks para configuração e simulação de preços
 */

import { pricingService } from '@/services/pricing-service';
import type {
    ListSimulacoesFilters,
    PrecificacaoConfigResponse,
    SaveConfigPrecificacaoRequest,
    SimularPrecoRequest
} from '@/types/pricing';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const pricingKeys = {
  all: ['pricing'] as const,
  config: () => [...pricingKeys.all, 'config'] as const,
  simulations: () => [...pricingKeys.all, 'simulations'] as const,
  simulationsList: (filters: ListSimulacoesFilters) =>
    [...pricingKeys.simulations(), filters] as const,
  simulation: (id: string) => [...pricingKeys.simulations(), id] as const,
};

// =============================================================================
// HOOKS DE CONSULTA (READ)
// =============================================================================

/**
 * Hook para buscar configuração de precificação
 */
export function usePricingConfig() {
  return useQuery({
    queryKey: pricingKeys.config(),
    queryFn: () => pricingService.getConfig(),
    staleTime: 5 * 60 * 1000, // 5 minutos - config muda raramente
  });
}

/**
 * Hook para listar simulações
 */
export function useSimulations(filters: ListSimulacoesFilters = {}) {
  return useQuery({
    queryKey: pricingKeys.simulationsList(filters),
    queryFn: () => pricingService.listSimulations(filters),
    staleTime: 60_000, // 1 minuto
  });
}

/**
 * Hook para buscar uma simulação específica
 */
export function useSimulation(id: string) {
  return useQuery({
    queryKey: pricingKeys.simulation(id),
    queryFn: () => pricingService.getSimulation(id),
    enabled: !!id,
  });
}

// =============================================================================
// HOOKS DE MUTAÇÃO (WRITE)
// =============================================================================

/**
 * Hook para salvar configuração
 */
export function useSaveConfig() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SaveConfigPrecificacaoRequest) =>
      pricingService.saveConfig(data),
    onSuccess: (data) => {
      queryClient.setQueryData<PrecificacaoConfigResponse>(
        pricingKeys.config(),
        data
      );
      toast.success('Configuração salva com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useSaveConfig] Erro:', error);
      toast.error('Erro ao salvar configuração');
    },
  });
}

/**
 * Hook para atualizar configuração
 */
export function useUpdateConfig() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SaveConfigPrecificacaoRequest) =>
      pricingService.updateConfig(data),
    onSuccess: (data) => {
      queryClient.setQueryData<PrecificacaoConfigResponse>(
        pricingKeys.config(),
        data
      );
      toast.success('Configuração atualizada!');
    },
    onError: (error: Error) => {
      console.error('[useUpdateConfig] Erro:', error);
      toast.error('Erro ao atualizar configuração');
    },
  });
}

/**
 * Hook para deletar configuração
 */
export function useDeleteConfig() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => pricingService.deleteConfig(),
    onSuccess: () => {
      queryClient.setQueryData(pricingKeys.config(), null);
      toast.success('Configuração removida!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteConfig] Erro:', error);
      toast.error('Erro ao remover configuração');
    },
  });
}

/**
 * Hook para simular preço
 */
export function useSimulatePrice() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SimularPrecoRequest) => pricingService.simulate(data),
    onSuccess: () => {
      // Invalidar lista de simulações
      queryClient.invalidateQueries({ queryKey: pricingKeys.simulations() });
    },
    onError: (error: Error) => {
      console.error('[useSimulatePrice] Erro:', error);
      toast.error('Erro ao simular preço');
    },
  });
}

/**
 * Hook para deletar simulação
 */
export function useDeleteSimulation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => pricingService.deleteSimulation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: pricingKeys.simulations() });
      toast.success('Simulação removida!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteSimulation] Erro:', error);
      toast.error('Erro ao remover simulação');
    },
  });
}
