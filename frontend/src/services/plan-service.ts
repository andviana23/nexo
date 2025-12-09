/**
 * NEXO - Sistema de Gestão para Barbearias
 * Plan Service
 *
 * Serviço de planos - comunicação com API de plans do backend.
 * Seguindo padrões do projeto e FLUXO_ASSINATURA.md
 */

import { api } from '@/lib/axios';
import type {
    CreatePlanRequest,
    ListPlansFilters,
    ListPlansResponse,
    Plan,
    PlanResponse,
    UpdatePlanRequest,
} from '@/types/subscription';

// =============================================================================
// ENDPOINTS
// =============================================================================

const PLAN_ENDPOINTS = {
  list: '/plans',
  create: '/plans',
  getById: (id: string) => `/plans/${id}`,
  update: (id: string) => `/plans/${id}`,
  deactivate: (id: string) => `/plans/${id}`,
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const planService = {
  /**
   * Lista todos os planos
   */
  async list(filters: ListPlansFilters = {}): Promise<Plan[]> {
    const params: Record<string, unknown> = {};
    
    if (filters.ativo !== undefined) {
      params.ativo = filters.ativo;
    }

    const response = await api.get<ListPlansResponse>(
      PLAN_ENDPOINTS.list,
      { params }
    );
    return response.data.data || response.data as unknown as Plan[];
  },

  /**
   * Lista planos ativos (para selects)
   */
  async listActive(): Promise<Plan[]> {
    return this.list({ ativo: true });
  },

  /**
   * Busca um plano pelo ID
   */
  async getById(id: string): Promise<Plan> {
    const response = await api.get<PlanResponse>(
      PLAN_ENDPOINTS.getById(id)
    );
    return response.data.data || response.data as unknown as Plan;
  },

  /**
   * Cria um novo plano
   */
  async create(data: CreatePlanRequest): Promise<Plan> {
    const response = await api.post<PlanResponse>(
      PLAN_ENDPOINTS.create,
      data
    );
    return response.data.data || response.data as unknown as Plan;
  },

  /**
   * Atualiza um plano existente
   */
  async update(id: string, data: UpdatePlanRequest): Promise<Plan> {
    const response = await api.put<PlanResponse>(
      PLAN_ENDPOINTS.update(id),
      data
    );
    return response.data.data || response.data as unknown as Plan;
  },

  /**
   * Desativa um plano
   * RN-PLAN-002: Plano não é deletado, apenas desativado
   */
  async deactivate(id: string): Promise<void> {
    await api.delete(PLAN_ENDPOINTS.deactivate(id));
  },

  /**
   * Reativa um plano
   */
  async activate(id: string): Promise<Plan> {
    return this.update(id, { ativo: true });
  },
};

export default planService;
