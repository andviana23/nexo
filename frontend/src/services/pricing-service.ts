/**
 * NEXO - Sistema de Gestão para Barbearias
 * Pricing Service
 *
 * Serviço de precificação - comunicação com API de pricing do backend.
 */

import { api } from '@/lib/axios';
import type {
    ListSimulacoesFilters,
    PrecificacaoConfigResponse,
    PrecificacaoSimulacaoResponse,
    SaveConfigPrecificacaoRequest,
    SimularPrecoRequest,
} from '@/types/pricing';

// =============================================================================
// ENDPOINTS
// =============================================================================

const PRICING_ENDPOINTS = {
  config: '/pricing/config',
  simulate: '/pricing/simulate',
  simulations: '/pricing/simulations',
  simulationById: (id: string) => `/pricing/simulations/${id}`,
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const pricingService = {
  // ===========================================================================
  // CONFIGURAÇÃO
  // ===========================================================================

  /**
   * Buscar configuração de precificação do tenant
   */
  async getConfig(): Promise<PrecificacaoConfigResponse | null> {
    try {
      const response = await api.get<PrecificacaoConfigResponse>(
        PRICING_ENDPOINTS.config
      );
      return response.data;
    } catch (error: unknown) {
      // 404 = não tem config ainda
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { status?: number } };
        if (axiosError.response?.status === 404) {
          return null;
        }
      }
      throw error;
    }
  },

  /**
   * Salvar configuração de precificação (cria se não existe)
   */
  async saveConfig(data: SaveConfigPrecificacaoRequest): Promise<PrecificacaoConfigResponse> {
    const response = await api.post<PrecificacaoConfigResponse>(
      PRICING_ENDPOINTS.config,
      data
    );
    return response.data;
  },

  /**
   * Atualizar configuração de precificação
   */
  async updateConfig(data: SaveConfigPrecificacaoRequest): Promise<PrecificacaoConfigResponse> {
    const response = await api.put<PrecificacaoConfigResponse>(
      PRICING_ENDPOINTS.config,
      data
    );
    return response.data;
  },

  /**
   * Deletar configuração de precificação
   */
  async deleteConfig(): Promise<void> {
    await api.delete(PRICING_ENDPOINTS.config);
  },

  // ===========================================================================
  // SIMULAÇÃO
  // ===========================================================================

  /**
   * Simular preço de um serviço ou produto
   */
  async simulate(data: SimularPrecoRequest): Promise<PrecificacaoSimulacaoResponse> {
    const response = await api.post<PrecificacaoSimulacaoResponse>(
      PRICING_ENDPOINTS.simulate,
      data
    );
    return response.data;
  },

  /**
   * Listar simulações salvas
   */
  async listSimulations(filters: ListSimulacoesFilters = {}): Promise<PrecificacaoSimulacaoResponse[]> {
    const response = await api.get<PrecificacaoSimulacaoResponse[]>(
      PRICING_ENDPOINTS.simulations,
      { params: filters }
    );
    return response.data;
  },

  /**
   * Buscar uma simulação específica
   */
  async getSimulation(id: string): Promise<PrecificacaoSimulacaoResponse> {
    const response = await api.get<PrecificacaoSimulacaoResponse>(
      PRICING_ENDPOINTS.simulationById(id)
    );
    return response.data;
  },

  /**
   * Deletar uma simulação
   */
  async deleteSimulation(id: string): Promise<void> {
    await api.delete(PRICING_ENDPOINTS.simulationById(id));
  },
};

export default pricingService;
