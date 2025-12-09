/**
 * NEXO - Sistema de Gestão para Barbearias
 * Hook de Feature Flags
 *
 * Hook para verificar feature flags do tenant.
 */

'use client';

import { featureFlagService, type FeatureFlags } from '@/services/feature-flag-service';
import { useCurrentTenant } from '@/store/auth-store';
import { useQuery } from '@tanstack/react-query';
import { useMemo } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface UseFeatureFlagsReturn {
  flags: FeatureFlags | null;
  isLoading: boolean;
  error: Error | null;
  // Helpers para flags específicas
  isMultiUnitEnabled: boolean;
  isFeatureEnabled: (flag: keyof FeatureFlags) => boolean;
}

// =============================================================================
// QUERY KEYS
// =============================================================================

// Adicionando query key para feature flags
const featureFlagsKeys = {
  all: ['feature-flags'] as const,
  tenant: (tenantId: string) => [...featureFlagsKeys.all, tenantId] as const,
};

// =============================================================================
// HOOK PRINCIPAL
// =============================================================================

export function useFeatureFlags(): UseFeatureFlagsReturn {
  const tenant = useCurrentTenant();

  const { data, isLoading, error } = useQuery({
    queryKey: featureFlagsKeys.tenant(tenant?.id ?? ''),
    queryFn: () => featureFlagService.getFlags(),
    enabled: !!tenant?.id,
    staleTime: 5 * 60 * 1000, // 5 minutos
    gcTime: 30 * 60 * 1000, // 30 minutos
  });

  const flags = data?.flags ?? null;

  // Memoiza o helper para evitar recriação
  const isFeatureEnabled = useMemo(() => {
    return (flag: keyof FeatureFlags): boolean => {
      if (!flags) return false;
      return flags[flag] ?? false;
    };
  }, [flags]);

  // Multi-unit pode vir da flag OU do tenant diretamente
  const isMultiUnitEnabled = useMemo(() => {
    // Prioridade: feature flag > tenant.multi_unit_enabled
    if (flags?.multi_unit_enabled !== undefined) {
      return flags.multi_unit_enabled;
    }
    return tenant?.multi_unit_enabled ?? false;
  }, [flags, tenant?.multi_unit_enabled]);

  return {
    flags,
    isLoading,
    error: error as Error | null,
    isMultiUnitEnabled,
    isFeatureEnabled,
  };
}

// =============================================================================
// HOOK SIMPLIFICADO - APENAS MULTI-UNIT
// =============================================================================

/**
 * Hook simplificado para verificar apenas se multi-unit está habilitado
 */
export function useMultiUnitEnabled(): boolean {
  const { isMultiUnitEnabled } = useFeatureFlags();
  return isMultiUnitEnabled;
}

// =============================================================================
// HOOK PARA FLAG ESPECÍFICA
// =============================================================================

/**
 * Hook para verificar uma feature flag específica
 */
export function useFeatureFlag(flag: keyof FeatureFlags): boolean {
  const { isFeatureEnabled } = useFeatureFlags();
  return isFeatureEnabled(flag);
}
