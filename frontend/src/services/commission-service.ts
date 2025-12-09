/**
 * NEXO - Sistema de Gestão para Barbearias
 * Commission Service
 *
 * Serviço de comunicação com a API de Comissões do backend.
 * Endpoints mapeados de: backend/internal/infra/http/handler/commission_handler.go
 */

import { api } from '@/lib/axios';
import type {
    Advance,
    AdvancesTotals,
    AssignItemsToPeriodRequest,
    CloseCommissionPeriodRequest,
    CommissionByService,
    CommissionItem,
    CommissionPeriod,
    CommissionPeriodSummary,
    CommissionRule,
    CommissionSummary,
    CreateAdvanceRequest,
    CreateCommissionItemBatchRequest,
    CreateCommissionItemRequest,
    CreateCommissionPeriodRequest,
    CreateCommissionRuleRequest,
    ListAdvancesFilters,
    ListAdvancesResponse,
    ListCommissionItemsFilters,
    ListCommissionItemsResponse,
    ListCommissionPeriodsFilters,
    ListCommissionPeriodsResponse,
    ListCommissionRulesFilters,
    ListCommissionRulesResponse,
    MarkAdvanceDeductedRequest,
    ProcessCommissionItemRequest,
    RejectAdvanceRequest,
    UpdateCommissionRuleRequest
} from '@/types/commission';

// =============================================================================
// ENDPOINTS
// =============================================================================

const COMMISSION_ENDPOINTS = {
  // Commission Rules
  rules: '/commissions/rules',
  ruleById: (id: string) => `/commissions/rules/${id}`,
  effectiveRules: '/commissions/rules/effective',

  // Commission Periods
  periods: '/commissions/periods',
  periodById: (id: string) => `/commissions/periods/${id}`,
  periodOpen: '/commissions/periods/open',
  periodSummary: (id: string) => `/commissions/periods/${id}/summary`,
  periodClose: (id: string) => `/commissions/periods/${id}/close`,
  periodPaid: (id: string) => `/commissions/periods/${id}/paid`,

  // Advances
  advances: '/commissions/advances',
  advanceById: (id: string) => `/commissions/advances/${id}`,
  advancesPending: '/commissions/advances/pending',
  advancesApproved: '/commissions/advances/approved',
  advanceApprove: (id: string) => `/commissions/advances/${id}/approve`,
  advanceReject: (id: string) => `/commissions/advances/${id}/reject`,
  advanceDeduct: (id: string) => `/commissions/advances/${id}/deduct`,
  advanceCancel: (id: string) => `/commissions/advances/${id}/cancel`,

  // Commission Items
  items: '/commissions/items',
  itemById: (id: string) => `/commissions/items/${id}`,
  itemsBatch: '/commissions/items/batch',
  itemsPending: '/commissions/items/pending',
  itemProcess: (id: string) => `/commissions/items/${id}/process`,
  itemsAssign: '/commissions/items/assign',
  summaryByProfessional: '/commissions/items/summary/professional',
  summaryByService: '/commissions/items/summary/service',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const commissionService = {
  // ===========================================================================
  // COMMISSION RULES (REGRAS DE COMISSÃO)
  // ===========================================================================

  /**
   * Lista regras de comissão com filtros opcionais
   */
  async listRules(filters: ListCommissionRulesFilters = {}): Promise<ListCommissionRulesResponse> {
    console.log('[commission-service] Listando regras:', filters);
    const { data } = await api.get<ListCommissionRulesResponse>(COMMISSION_ENDPOINTS.rules, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca regra de comissão por ID
   */
  async getRule(id: string): Promise<CommissionRule> {
    console.log('[commission-service] Buscando regra:', id);
    const { data } = await api.get<CommissionRule>(COMMISSION_ENDPOINTS.ruleById(id));
    return data;
  },

  /**
   * Cria nova regra de comissão
   */
  async createRule(payload: CreateCommissionRuleRequest): Promise<CommissionRule> {
    console.log('[commission-service] Criando regra:', payload);
    const { data } = await api.post<CommissionRule>(COMMISSION_ENDPOINTS.rules, payload);
    return data;
  },

  /**
   * Atualiza regra de comissão
   */
  async updateRule(id: string, payload: UpdateCommissionRuleRequest): Promise<CommissionRule> {
    console.log('[commission-service] Atualizando regra:', id, payload);
    const { data } = await api.put<CommissionRule>(COMMISSION_ENDPOINTS.ruleById(id), payload);
    return data;
  },

  /**
   * Deleta regra de comissão
   */
  async deleteRule(id: string): Promise<void> {
    console.log('[commission-service] Deletando regra:', id);
    await api.delete(COMMISSION_ENDPOINTS.ruleById(id));
  },

  /**
   * Busca regras efetivas (ativas e no período)
   */
  async getEffectiveRules(professionalId?: string): Promise<CommissionRule[]> {
    console.log('[commission-service] Buscando regras efetivas:', professionalId);
    const { data } = await api.get<CommissionRule[]>(COMMISSION_ENDPOINTS.effectiveRules, {
      params: professionalId ? { professional_id: professionalId } : undefined,
    });
    return data;
  },

  // ===========================================================================
  // COMMISSION PERIODS (PERÍODOS DE COMISSÃO)
  // ===========================================================================

  /**
   * Lista períodos de comissão com filtros
   */
  async listPeriods(filters: ListCommissionPeriodsFilters = {}): Promise<ListCommissionPeriodsResponse> {
    console.log('[commission-service] Listando períodos:', filters);
    const { data } = await api.get<ListCommissionPeriodsResponse>(COMMISSION_ENDPOINTS.periods, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca período por ID
   */
  async getPeriod(id: string): Promise<CommissionPeriod> {
    console.log('[commission-service] Buscando período:', id);
    const { data } = await api.get<CommissionPeriod>(COMMISSION_ENDPOINTS.periodById(id));
    return data;
  },

  /**
   * Cria novo período de comissão
   */
  async createPeriod(payload: CreateCommissionPeriodRequest): Promise<CommissionPeriod> {
    console.log('[commission-service] Criando período:', payload);
    const { data } = await api.post<CommissionPeriod>(COMMISSION_ENDPOINTS.periods, payload);
    return data;
  },

  /**
   * Busca período aberto de um profissional
   */
  async getOpenPeriod(professionalId: string): Promise<CommissionPeriod | null> {
    console.log('[commission-service] Buscando período aberto:', professionalId);
    try {
      const { data } = await api.get<CommissionPeriod>(COMMISSION_ENDPOINTS.periodOpen, {
        params: { professional_id: professionalId },
      });
      return data;
    } catch {
      return null;
    }
  },

  /**
   * Busca resumo do período
   */
  async getPeriodSummary(id: string): Promise<CommissionPeriodSummary> {
    console.log('[commission-service] Buscando resumo do período:', id);
    const { data } = await api.get<CommissionPeriodSummary>(COMMISSION_ENDPOINTS.periodSummary(id));
    return data;
  },

  /**
   * Fecha período de comissão
   */
  async closePeriod(id: string, payload?: CloseCommissionPeriodRequest): Promise<CommissionPeriod> {
    console.log('[commission-service] Fechando período:', id);
    const { data } = await api.post<CommissionPeriod>(COMMISSION_ENDPOINTS.periodClose(id), payload || {});
    return data;
  },

  /**
   * Marca período como pago
   */
  async markPeriodAsPaid(id: string): Promise<CommissionPeriod> {
    console.log('[commission-service] Marcando período como pago:', id);
    const { data } = await api.post<CommissionPeriod>(COMMISSION_ENDPOINTS.periodPaid(id));
    return data;
  },

  /**
   * Deleta período de comissão
   */
  async deletePeriod(id: string): Promise<void> {
    console.log('[commission-service] Deletando período:', id);
    await api.delete(COMMISSION_ENDPOINTS.periodById(id));
  },

  // ===========================================================================
  // ADVANCES (ADIANTAMENTOS)
  // ===========================================================================

  /**
   * Lista adiantamentos com filtros
   */
  async listAdvances(filters: ListAdvancesFilters = {}): Promise<ListAdvancesResponse> {
    console.log('[commission-service] Listando adiantamentos:', filters);
    const { data } = await api.get<ListAdvancesResponse>(COMMISSION_ENDPOINTS.advances, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca adiantamento por ID
   */
  async getAdvance(id: string): Promise<Advance> {
    console.log('[commission-service] Buscando adiantamento:', id);
    const { data } = await api.get<Advance>(COMMISSION_ENDPOINTS.advanceById(id));
    return data;
  },

  /**
   * Cria novo adiantamento
   */
  async createAdvance(payload: CreateAdvanceRequest): Promise<Advance> {
    console.log('[commission-service] Criando adiantamento:', payload);
    const { data } = await api.post<Advance>(COMMISSION_ENDPOINTS.advances, payload);
    return data;
  },

  /**
   * Lista adiantamentos pendentes (com totais)
   */
  async getPendingAdvances(professionalId?: string): Promise<AdvancesTotals> {
    console.log('[commission-service] Buscando adiantamentos pendentes:', professionalId);
    const { data } = await api.get<AdvancesTotals>(COMMISSION_ENDPOINTS.advancesPending, {
      params: professionalId ? { professional_id: professionalId } : undefined,
    });
    return data;
  },

  /**
   * Lista adiantamentos aprovados (com totais)
   */
  async getApprovedAdvances(professionalId?: string): Promise<AdvancesTotals> {
    console.log('[commission-service] Buscando adiantamentos aprovados:', professionalId);
    const { data } = await api.get<AdvancesTotals>(COMMISSION_ENDPOINTS.advancesApproved, {
      params: professionalId ? { professional_id: professionalId } : undefined,
    });
    return data;
  },

  /**
   * Aprova adiantamento
   */
  async approveAdvance(id: string): Promise<Advance> {
    console.log('[commission-service] Aprovando adiantamento:', id);
    const { data } = await api.post<Advance>(COMMISSION_ENDPOINTS.advanceApprove(id));
    return data;
  },

  /**
   * Rejeita adiantamento
   */
  async rejectAdvance(id: string, payload: RejectAdvanceRequest): Promise<Advance> {
    console.log('[commission-service] Rejeitando adiantamento:', id);
    const { data } = await api.post<Advance>(COMMISSION_ENDPOINTS.advanceReject(id), payload);
    return data;
  },

  /**
   * Marca adiantamento como deduzido
   */
  async markAdvanceDeducted(id: string, payload: MarkAdvanceDeductedRequest): Promise<Advance> {
    console.log('[commission-service] Marcando adiantamento como deduzido:', id);
    const { data } = await api.post<Advance>(COMMISSION_ENDPOINTS.advanceDeduct(id), payload);
    return data;
  },

  /**
   * Cancela adiantamento
   */
  async cancelAdvance(id: string): Promise<Advance> {
    console.log('[commission-service] Cancelando adiantamento:', id);
    const { data } = await api.post<Advance>(COMMISSION_ENDPOINTS.advanceCancel(id));
    return data;
  },

  /**
   * Deleta adiantamento
   */
  async deleteAdvance(id: string): Promise<void> {
    console.log('[commission-service] Deletando adiantamento:', id);
    await api.delete(COMMISSION_ENDPOINTS.advanceById(id));
  },

  // ===========================================================================
  // COMMISSION ITEMS (ITENS DE COMISSÃO)
  // ===========================================================================

  /**
   * Lista itens de comissão com filtros
   */
  async listItems(filters: ListCommissionItemsFilters = {}): Promise<ListCommissionItemsResponse> {
    console.log('[commission-service] Listando itens:', filters);
    const { data } = await api.get<ListCommissionItemsResponse>(COMMISSION_ENDPOINTS.items, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca item de comissão por ID
   */
  async getItem(id: string): Promise<CommissionItem> {
    console.log('[commission-service] Buscando item:', id);
    const { data } = await api.get<CommissionItem>(COMMISSION_ENDPOINTS.itemById(id));
    return data;
  },

  /**
   * Cria novo item de comissão
   */
  async createItem(payload: CreateCommissionItemRequest): Promise<CommissionItem> {
    console.log('[commission-service] Criando item:', payload);
    const { data } = await api.post<CommissionItem>(COMMISSION_ENDPOINTS.items, payload);
    return data;
  },

  /**
   * Cria múltiplos itens de comissão (batch)
   */
  async createItemsBatch(payload: CreateCommissionItemBatchRequest): Promise<CommissionItem[]> {
    console.log('[commission-service] Criando itens em batch:', payload);
    const { data } = await api.post<CommissionItem[]>(COMMISSION_ENDPOINTS.itemsBatch, payload);
    return data;
  },

  /**
   * Lista itens pendentes de um profissional
   */
  async getPendingItems(professionalId: string): Promise<ListCommissionItemsResponse> {
    console.log('[commission-service] Buscando itens pendentes:', professionalId);
    const { data } = await api.get<ListCommissionItemsResponse>(COMMISSION_ENDPOINTS.itemsPending, {
      params: { professional_id: professionalId },
    });
    return data;
  },

  /**
   * Processa item de comissão (vincula a período)
   */
  async processItem(id: string, payload: ProcessCommissionItemRequest): Promise<CommissionItem> {
    console.log('[commission-service] Processando item:', id);
    const { data } = await api.post<CommissionItem>(COMMISSION_ENDPOINTS.itemProcess(id), payload);
    return data;
  },

  /**
   * Vincula itens a um período em lote
   */
  async assignItemsToPeriod(payload: AssignItemsToPeriodRequest): Promise<{ assigned_count: number }> {
    console.log('[commission-service] Vinculando itens ao período:', payload);
    const { data } = await api.post<{ assigned_count: number }>(COMMISSION_ENDPOINTS.itemsAssign, payload);
    return data;
  },

  /**
   * Deleta item de comissão
   */
  async deleteItem(id: string): Promise<void> {
    console.log('[commission-service] Deletando item:', id);
    await api.delete(COMMISSION_ENDPOINTS.itemById(id));
  },

  // ===========================================================================
  // SUMMARIES (RESUMOS)
  // ===========================================================================

  /**
   * Busca resumo de comissões por profissional
   */
  async getSummaryByProfessional(
    startDate: string,
    endDate: string,
    professionalId?: string
  ): Promise<CommissionSummary[]> {
    console.log('[commission-service] Buscando resumo por profissional:', { startDate, endDate, professionalId });
    const { data } = await api.get<CommissionSummary[]>(COMMISSION_ENDPOINTS.summaryByProfessional, {
      params: {
        start_date: startDate,
        end_date: endDate,
        professional_id: professionalId,
      },
    });
    return data;
  },

  /**
   * Busca resumo de comissões por serviço
   */
  async getSummaryByService(
    startDate: string,
    endDate: string,
    serviceId?: string
  ): Promise<CommissionByService[]> {
    console.log('[commission-service] Buscando resumo por serviço:', { startDate, endDate, serviceId });
    const { data } = await api.get<CommissionByService[]>(COMMISSION_ENDPOINTS.summaryByService, {
      params: {
        start_date: startDate,
        end_date: endDate,
        service_id: serviceId,
      },
    });
    return data;
  },
};

export default commissionService;
