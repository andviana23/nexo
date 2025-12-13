/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Hooks
 *
 * Hooks React Query para gerenciar estado do módulo Caixa Diário.
 *
 * @author NEXO v2.0
 */

import { caixaService } from '@/services/caixa-service';
import { useActiveUnitId, useNeedsSelection, useUnitHydrated } from '@/store/unit-store';
import type {
    AbrirCaixaRequest,
    FecharCaixaRequest,
    ListCaixaHistoricoFilters,
    ReforcoRequest,
    SangriaRequest,
} from '@/types/caixa';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const caixaKeys = {
  all: ['caixa'] as const,

  // Status
  status: () => [...caixaKeys.all, 'status'] as const,

  // Caixa aberto
  aberto: () => [...caixaKeys.all, 'aberto'] as const,

  // Totais
  totais: () => [...caixaKeys.all, 'totais'] as const,

  // Histórico
  historico: () => [...caixaKeys.all, 'historico'] as const,
  historicoList: (filters: ListCaixaHistoricoFilters) =>
    [...caixaKeys.historico(), 'list', filters] as const,

  // Detalhes
  detail: (id: string) => [...caixaKeys.all, 'detail', id] as const,
};

// =============================================================================
// QUERIES
// =============================================================================

/**
 * Hook para verificar status do caixa (aberto ou fechado)
 */
export function useCaixaStatus() {
  const activeUnitId = useActiveUnitId();
  const isUnitHydrated = useUnitHydrated();
  const needsSelection = useNeedsSelection();
  const unitReady = isUnitHydrated && !needsSelection && !!activeUnitId;

  return useQuery({
    queryKey: caixaKeys.status(),
    queryFn: () => caixaService.getStatus(),
    staleTime: 30 * 1000, // 30 segundos
    refetchInterval: 60 * 1000, // Atualiza a cada 1 minuto
    enabled: unitReady,
  });
}

/**
 * Hook para buscar caixa aberto atual
 */
export function useCaixaAberto() {
  const activeUnitId = useActiveUnitId();
  const isUnitHydrated = useUnitHydrated();
  const needsSelection = useNeedsSelection();
  const unitReady = isUnitHydrated && !needsSelection && !!activeUnitId;

  return useQuery({
    queryKey: caixaKeys.aberto(),
    queryFn: () => caixaService.getCaixaAberto(),
    retry: false, // Não retentar se não houver caixa aberto
    enabled: unitReady,
  });
}

/**
 * Hook para buscar totais do caixa
 */
export function useCaixaTotais(enabled: boolean = true) {
  const activeUnitId = useActiveUnitId();
  const isUnitHydrated = useUnitHydrated();
  const needsSelection = useNeedsSelection();
  const unitReady = isUnitHydrated && !needsSelection && !!activeUnitId;

  return useQuery({
    queryKey: caixaKeys.totais(),
    queryFn: () => caixaService.getTotais(),
    staleTime: 10 * 1000, // 10 segundos - atualiza mais rápido para refletir vendas
    enabled: enabled && unitReady,
    retry: false, // Não retentar se falhar (ex: caixa fechado)
  });
}

/**
 * Hook para listar histórico de caixas
 */
export function useCaixaHistorico(filters: ListCaixaHistoricoFilters = {}) {
  const activeUnitId = useActiveUnitId();
  const isUnitHydrated = useUnitHydrated();
  const needsSelection = useNeedsSelection();
  const unitReady = isUnitHydrated && !needsSelection && !!activeUnitId;

  return useQuery({
    queryKey: caixaKeys.historicoList(filters),
    queryFn: () => caixaService.getHistorico(filters),
    enabled: unitReady,
  });
}

/**
 * Hook para buscar caixa por ID
 */
export function useCaixaById(id: string) {
  const activeUnitId = useActiveUnitId();
  const isUnitHydrated = useUnitHydrated();
  const needsSelection = useNeedsSelection();
  const unitReady = isUnitHydrated && !needsSelection && !!activeUnitId;

  return useQuery({
    queryKey: caixaKeys.detail(id),
    queryFn: () => caixaService.getCaixaById(id),
    enabled: unitReady && !!id,
  });
}

// =============================================================================
// MUTATIONS
// =============================================================================

/**
 * Hook para abrir caixa
 */
export function useAbrirCaixa() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AbrirCaixaRequest) => caixaService.abrirCaixa(data),
    onSuccess: () => {
      // Invalida queries relacionadas
      queryClient.invalidateQueries({ queryKey: caixaKeys.status() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.aberto() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.totais() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.historico() });
      toast.success('Caixa aberto com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useAbrirCaixa] Erro:', error);
      toast.error(error.message || 'Erro ao abrir caixa');
    },
  });
}

/**
 * Hook para registrar sangria
 */
export function useSangria() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SangriaRequest) => caixaService.registrarSangria(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: caixaKeys.aberto() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.totais() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.status() });
      toast.success('Sangria registrada com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useSangria] Erro:', error);
      toast.error(error.message || 'Erro ao registrar sangria');
    },
  });
}

/**
 * Hook para registrar reforço
 */
export function useReforco() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ReforcoRequest) => caixaService.registrarReforco(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: caixaKeys.aberto() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.totais() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.status() });
      toast.success('Reforço registrado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useReforco] Erro:', error);
      toast.error(error.message || 'Erro ao registrar reforço');
    },
  });
}

/**
 * Hook para fechar caixa
 */
export function useFecharCaixa() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: FecharCaixaRequest) => caixaService.fecharCaixa(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: caixaKeys.status() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.aberto() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.totais() });
      queryClient.invalidateQueries({ queryKey: caixaKeys.historico() });
      toast.success('Caixa fechado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useFecharCaixa] Erro:', error);
      toast.error(error.message || 'Erro ao fechar caixa');
    },
  });
}

// =============================================================================
// HOOKS COMPOSTOS
// =============================================================================

/**
 * Hook que combina status e caixa aberto para facilitar uso na UI
 */
export function useCaixaDiario() {
  const status = useCaixaStatus();
  const isAberto = status.data?.aberto ?? false;
  
  // Só busca totais se houver caixa aberto
  const totais = useCaixaTotais(isAberto);

  return {
    // Status
    isLoading: status.isLoading,
    isError: status.isError,
    error: status.error,

    // Dados
    isAberto,
    caixaAtual: status.data?.caixa_atual,
    ultimoFechamento: status.data?.ultimo_fechamento,

    // Totais
    totais: totais.data,
    totaisLoading: totais.isLoading,

    // Refetch
    refetch: () => {
      status.refetch();
      if (isAberto) {
        totais.refetch();
      }
    },
  };
}
