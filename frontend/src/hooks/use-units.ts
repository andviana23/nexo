/**
 * NEXO - Sistema de Gestão para Barbearias
 * Hook de Unidades
 *
 * Hook que combina Zustand (estado) + React Query (server state) para unidades.
 * Gerencia seleção de unidade ativa e troca entre unidades.
 */

'use client';

import { analytics } from '@/lib/analytics';
import { getErrorMessage, isAxiosError } from '@/lib/axios';
import { queryKeys } from '@/lib/query-client';
import { unitService } from '@/services/unit-service';
import {
    getActiveUnitId,
    useActiveUnit,
    useIsMultiUnit,
    useUnitError,
    useUnitHydrated,
    useUnitLoading,
    useUnits,
    useUnitStore,
} from '@/store/unit-store';
import type { UserUnit } from '@/types/unit';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect } from 'react';
import { useMultiUnitEnabled } from './use-feature-flags';

// =============================================================================
// TIPOS
// =============================================================================

interface UseUnitsReturn {
  // Estado
  units: UserUnit[];
  activeUnit: UserUnit | null;
  activeUnitId: string | null;
  isLoading: boolean;
  isHydrated: boolean;
  isMultiUnit: boolean;
  isMultiUnitEnabled: boolean; // Feature flag
  error: string | null;

  // Ações
  switchUnit: (unitId: string) => Promise<void>;
  setDefaultUnit: (unitId: string) => Promise<void>;
  refreshUnits: () => Promise<void>;

  // Status das mutations
  isSwitching: boolean;
  switchError: string | null;
}

// =============================================================================
// HOOK PRINCIPAL
// =============================================================================

export function useUnit(): UseUnitsReturn {
  const queryClient = useQueryClient();

  // Feature flag - verifica se multi-unit está habilitado para o tenant
  const isMultiUnitEnabled = useMultiUnitEnabled();

  // Estado do Zustand
  const { setUnits, setActiveUnit, setError } = useUnitStore();
  const units = useUnits();
  const activeUnit = useActiveUnit();
  const storeIsMultiUnit = useIsMultiUnit();
  const isLoading = useUnitLoading();
  const isHydrated = useUnitHydrated();
  const storeError = useUnitError();

  // Multi-unit só é true se a feature flag está habilitada E tem múltiplas unidades
  const isMultiUnit = isMultiUnitEnabled && storeIsMultiUnit;

  // ==========================================================================
  // QUERY - LISTAR UNIDADES DO USUÁRIO
  // ==========================================================================

  const unitsQuery = useQuery({
    queryKey: queryKeys.units.userUnits(),
    queryFn: () => unitService.getUserUnits(),
    staleTime: 5 * 60 * 1000, // 5 minutos
    enabled: isHydrated && isMultiUnitEnabled, // Só busca se multi-unit habilitado
  });

  // Sincroniza query data com Zustand
  useEffect(() => {
    if (unitsQuery.data?.units) {
      const fetchedUnits = unitsQuery.data.units;
      setUnits(fetchedUnits);

      // Se não tem unidade ativa, seleciona a default ou primeira
      if (!activeUnit && fetchedUnits.length > 0) {
        const defaultUnit = fetchedUnits.find((u) => u.is_default);
        setActiveUnit(defaultUnit || fetchedUnits[0]);
      }

      // Se a unidade ativa não está mais na lista, seleciona a default
      if (activeUnit && !fetchedUnits.find((u) => u.unit_id === activeUnit.unit_id)) {
        const defaultUnit = fetchedUnits.find((u) => u.is_default);
        setActiveUnit(defaultUnit || fetchedUnits[0]);
      }
    }
  }, [unitsQuery.data, activeUnit, setUnits, setActiveUnit]);

  // ==========================================================================
  // MUTATION - TROCAR UNIDADE
  // ==========================================================================

  const switchMutation = useMutation({
    mutationFn: (unitId: string) => unitService.switchUnit(unitId),
    onSuccess: (data) => {
      console.log('📍 [useUnit] Troca de unidade bem-sucedida:', data.unit.unit_nome);
      setActiveUnit(data.unit);

      // Invalida todas as queries que dependem de unidade
      queryClient.invalidateQueries();
    },
    onError: (error) => {
      console.error('📍 [useUnit] Erro ao trocar unidade:', error);
      const message = isAxiosError(error) ? getErrorMessage(error) : 'Erro ao trocar unidade';
      setError(message);
    },
  });

  // ==========================================================================
  // MUTATION - DEFINIR UNIDADE PADRÃO
  // ==========================================================================

  const setDefaultMutation = useMutation({
    mutationFn: (unitId: string) => unitService.setDefaultUnit(unitId),
    onSuccess: () => {
      console.log('📍 [useUnit] Unidade padrão definida');
      queryClient.invalidateQueries({ queryKey: queryKeys.units.userUnits() });
    },
    onError: (error) => {
      console.error('📍 [useUnit] Erro ao definir unidade padrão:', error);
    },
  });

  // ==========================================================================
  // AÇÕES EXPOSTAS
  // ==========================================================================

  const switchUnit = useCallback(
    async (unitId: string) => {
      // Se já é a unidade ativa, não faz nada
      if (activeUnit?.unit_id === unitId) return;

      // Busca a unidade na lista local
      const unit = units.find((u) => u.unit_id === unitId);
      if (!unit) {
        setError('Unidade não encontrada');
        return;
      }

      // Telemetria - rastrear troca de unidade
      analytics.trackUnitSwitch(
        activeUnit?.unit_id ?? 'none',
        unit.unit_id,
        unit.unit_nome
      );

      // Se o backend não precisa de chamada API para troca simples,
      // apenas atualiza o estado local
      // await switchMutation.mutateAsync(unitId);

      // Para MVP: troca local sem chamada API
      setActiveUnit(unit);
      console.log('📍 [useUnit] Unidade trocada localmente:', unit.unit_nome);

      // Invalida todas as queries para recarregar dados da nova unidade
      queryClient.invalidateQueries();
    },
    [activeUnit, units, setActiveUnit, setError, queryClient]
  );

  const setDefaultUnit = useCallback(
    async (unitId: string) => {
      const unit = units.find((u) => u.unit_id === unitId);
      
      // Telemetria - rastrear definição de unidade padrão
      if (unit) {
        analytics.trackSetDefaultUnit(unitId, unit.unit_nome);
      }

      await setDefaultMutation.mutateAsync(unitId);
    },
    [setDefaultMutation, units]
  );

  const refreshUnits = useCallback(async () => {
    await queryClient.invalidateQueries({ queryKey: queryKeys.units.userUnits() });
  }, [queryClient]);

  // ==========================================================================
  // RETORNO
  // ==========================================================================

  return {
    // Estado
    units,
    activeUnit,
    activeUnitId: activeUnit?.unit_id ?? null,
    isLoading: isLoading || unitsQuery.isLoading,
    isHydrated,
    isMultiUnit,
    isMultiUnitEnabled, // Feature flag do tenant
    error: storeError || (unitsQuery.error ? getErrorMessage(unitsQuery.error) : null),

    // Ações
    switchUnit,
    setDefaultUnit,
    refreshUnits,

    // Status das mutations
    isSwitching: switchMutation.isPending,
    switchError: switchMutation.error ? getErrorMessage(switchMutation.error) : null,
  };
}

// =============================================================================
// HOOK UTILITÁRIO - APENAS ID DA UNIDADE ATIVA
// =============================================================================

/**
 * Hook simplificado para obter apenas o ID da unidade ativa
 * Útil em componentes que só precisam do ID para chamadas de API
 */
export function useActiveUnitId(): string | null {
  const activeUnit = useActiveUnit();
  return activeUnit?.unit_id ?? null;
}

// =============================================================================
// HELPER - GET UNIT ID (para uso fora de React)
// =============================================================================

export { getActiveUnitId };
