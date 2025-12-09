/**
 * NEXO - Sistema de Gestão para Barbearias
 * Subscription Service
 *
 * Serviço de assinaturas - comunicação com API de subscriptions do backend.
 * Seguindo padrões do projeto e FLUXO_ASSINATURA.md
 */

import { api } from '@/lib/axios';
import type {
    CreateSubscriptionRequest,
    ListSubscriptionsFilters,
    ListSubscriptionsResponse,
    RenewSubscriptionRequest,
    Subscription,
    SubscriptionMetrics,
    SubscriptionMetricsResponse,
    SubscriptionResponse,
} from '@/types/subscription';

// =============================================================================
// ENDPOINTS
// =============================================================================

const SUBSCRIPTION_ENDPOINTS = {
  list: '/subscriptions',
  create: '/subscriptions',
  getById: (id: string) => `/subscriptions/${id}`,
  cancel: (id: string) => `/subscriptions/${id}`,
  renew: (id: string) => `/subscriptions/${id}/renew`,
  metrics: '/subscriptions/metrics',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const subscriptionService = {
  /**
   * Lista assinaturas com filtros
   */
  async list(filters: ListSubscriptionsFilters = {}): Promise<Subscription[]> {
    const params: Record<string, unknown> = {};

    if (filters.page) params.page = filters.page;
    if (filters.page_size) params.page_size = filters.page_size;
    if (filters.status) params.status = filters.status;
    if (filters.plano_id) params.plano_id = filters.plano_id;
    if (filters.search) params.search = filters.search;

    console.log('[subscriptionService.list] Buscando assinaturas com filtros:', params);
    
    const response = await api.get<ListSubscriptionsResponse>(
      SUBSCRIPTION_ENDPOINTS.list,
      { params }
    );
    
    console.log('[subscriptionService.list] Resposta da API:', {
      status: response.status,
      dataType: typeof response.data,
      isArray: Array.isArray(response.data),
      hasDataProp: 'data' in (response.data || {}),
      length: Array.isArray(response.data) ? response.data.length : (response.data?.data?.length ?? 0),
    });
    
    // API retorna array diretamente, não { data: [...] }
    const subscriptions = Array.isArray(response.data) 
      ? response.data 
      : (response.data?.data || []);
    
    console.log('[subscriptionService.list] Assinaturas encontradas:', subscriptions.length);
    
    return subscriptions;
  },

  /**
   * Busca uma assinatura pelo ID
   */
  async getById(id: string): Promise<Subscription> {
    const response = await api.get<SubscriptionResponse>(
      SUBSCRIPTION_ENDPOINTS.getById(id)
    );
    return response.data.data || response.data as unknown as Subscription;
  },

  /**
   * Cria uma nova assinatura
   * RN-SUB-001: Validar se cliente já possui assinatura ativa
   */
  async create(data: CreateSubscriptionRequest): Promise<Subscription> {
    const response = await api.post<SubscriptionResponse>(
      SUBSCRIPTION_ENDPOINTS.create,
      data
    );
    return response.data.data || response.data as unknown as Subscription;
  },

  /**
   * Cancela uma assinatura
   * RN-CANC-003: Apenas admin/gerente pode cancelar
   */
  async cancel(id: string): Promise<void> {
    await api.delete(SUBSCRIPTION_ENDPOINTS.cancel(id));
  },

  /**
   * Renova uma assinatura (PIX/Dinheiro)
   * RN-REN-001: Registra pagamento manual
   */
  async renew(id: string, data: RenewSubscriptionRequest = {}): Promise<Subscription> {
    const response = await api.post<SubscriptionResponse>(
      SUBSCRIPTION_ENDPOINTS.renew(id),
      data
    );
    return response.data.data || response.data as unknown as Subscription;
  },

  /**
   * Obtém métricas de assinaturas
   */
  async getMetrics(): Promise<SubscriptionMetrics> {
    const response = await api.get<SubscriptionMetricsResponse>(
      SUBSCRIPTION_ENDPOINTS.metrics
    );
    return response.data.data || response.data as unknown as SubscriptionMetrics;
  },

  /**
   * Lista assinaturas por status
   */
  async listByStatus(status: ListSubscriptionsFilters['status']): Promise<Subscription[]> {
    return this.list({ status });
  },

  /**
   * Lista assinaturas ativas
   */
  async listActive(): Promise<Subscription[]> {
    return this.listByStatus('ATIVO');
  },

  /**
   * Lista assinaturas inadimplentes
   */
  async listOverdue(): Promise<Subscription[]> {
    return this.listByStatus('INADIMPLENTE');
  },
};

export default subscriptionService;
