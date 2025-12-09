/**
 * NEXO - Sistema de Gestão para Barbearias
 * Feature Flag Service
 *
 * Serviço para consulta de feature flags do tenant.
 */

import { api } from '@/lib/axios';

// =============================================================================
// TIPOS
// =============================================================================

export interface FeatureFlags {
  multi_unit_enabled: boolean;
  use_v2_financial?: boolean;
  use_new_agenda?: boolean;
  enable_online_booking?: boolean;
  enable_stock_management?: boolean;
  enable_comanda?: boolean;
  [key: string]: boolean | undefined;
}

export interface FeatureFlagResponse {
  flags: FeatureFlags;
  tenant_id: string;
}

// =============================================================================
// ENDPOINTS
// =============================================================================

const FEATURE_FLAG_ENDPOINTS = {
  list: '/feature-flags',
  check: (flag: string) => `/feature-flags/${flag}`,
} as const;

// =============================================================================
// SERVIÇO
// =============================================================================

export const featureFlagService = {
  /**
   * Lista todas as feature flags do tenant
   */
  async getFlags(): Promise<FeatureFlagResponse> {
    const response = await api.get<FeatureFlagResponse>(FEATURE_FLAG_ENDPOINTS.list);
    return response.data;
  },

  /**
   * Verifica se uma feature flag específica está habilitada
   */
  async checkFlag(flagName: string): Promise<boolean> {
    try {
      const response = await api.get<{ enabled: boolean }>(
        FEATURE_FLAG_ENDPOINTS.check(flagName)
      );
      return response.data.enabled;
    } catch {
      // Se o endpoint não existir ou der erro, assume desabilitado
      return false;
    }
  },

  /**
   * Verifica múltiplas flags de uma vez
   */
  async checkFlags(flagNames: string[]): Promise<Record<string, boolean>> {
    try {
      const flags = await this.getFlags();
      return flagNames.reduce(
        (acc, flag) => ({
          ...acc,
          [flag]: flags.flags[flag] ?? false,
        }),
        {}
      );
    } catch {
      // Se der erro, retorna todas como false
      return flagNames.reduce((acc, flag) => ({ ...acc, [flag]: false }), {});
    }
  },
};

export default featureFlagService;
