/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Service
 *
 * Serviço de comunicação com a API de Caixa Diário do backend.
 * Endpoints mapeados de: backend/internal/infra/http/handler/caixa_handler.go
 *
 * @author NEXO v2.0
 */

import { api } from '@/lib/axios';
import type {
    AbrirCaixaRequest,
    CaixaDiarioResponse,
    CaixaStatusResponse,
    FecharCaixaRequest,
    ListCaixaHistoricoFilters,
    ListCaixaHistoricoResponse,
    ReforcoRequest,
    SangriaRequest,
    TotaisCaixaResponse,
} from '@/types/caixa';

// =============================================================================
// ENDPOINTS
// =============================================================================

const CAIXA_ENDPOINTS = {
  // Status
  status: '/caixa/status',

  // Operações principais
  abrir: '/caixa/abrir',
  sangria: '/caixa/sangria',
  reforco: '/caixa/reforco',
  fechar: '/caixa/fechar',

  // Consultas
  aberto: '/caixa/aberto',
  historico: '/caixa/historico',
  totais: '/caixa/totais',
  byId: (id: string) => `/caixa/${id}`,
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const caixaService = {
  // ===========================================================================
  // STATUS
  // ===========================================================================

  /**
   * Verifica o status atual do caixa (aberto ou fechado)
   */
  async getStatus(): Promise<CaixaStatusResponse> {
    console.log('[caixa-service] Verificando status do caixa');
    const { data } = await api.get<CaixaStatusResponse>(CAIXA_ENDPOINTS.status);
    return data;
  },

  // ===========================================================================
  // OPERAÇÕES
  // ===========================================================================

  /**
   * Abre o caixa com saldo inicial
   */
  async abrirCaixa(payload: AbrirCaixaRequest): Promise<CaixaDiarioResponse> {
    console.log('[caixa-service] Abrindo caixa:', payload);
    const { data } = await api.post<CaixaDiarioResponse>(
      CAIXA_ENDPOINTS.abrir,
      payload
    );
    return data;
  },

  /**
   * Registra sangria (retirada) do caixa
   */
  async registrarSangria(payload: SangriaRequest): Promise<CaixaDiarioResponse> {
    console.log('[caixa-service] Registrando sangria:', payload);
    const { data } = await api.post<CaixaDiarioResponse>(
      CAIXA_ENDPOINTS.sangria,
      payload
    );
    return data;
  },

  /**
   * Registra reforço (entrada) no caixa
   */
  async registrarReforco(payload: ReforcoRequest): Promise<CaixaDiarioResponse> {
    console.log('[caixa-service] Registrando reforço:', payload);
    const { data } = await api.post<CaixaDiarioResponse>(
      CAIXA_ENDPOINTS.reforco,
      payload
    );
    return data;
  },

  /**
   * Fecha o caixa com saldo real informado
   */
  async fecharCaixa(payload: FecharCaixaRequest): Promise<CaixaDiarioResponse> {
    console.log('[caixa-service] Fechando caixa:', payload);
    const { data } = await api.post<CaixaDiarioResponse>(
      CAIXA_ENDPOINTS.fechar,
      payload
    );
    return data;
  },

  // ===========================================================================
  // CONSULTAS
  // ===========================================================================

  /**
   * Busca o caixa aberto atual
   */
  async getCaixaAberto(): Promise<CaixaDiarioResponse> {
    console.log('[caixa-service] Buscando caixa aberto');
    const { data } = await api.get<CaixaDiarioResponse>(CAIXA_ENDPOINTS.aberto);
    return data;
  },

  /**
   * Lista histórico de caixas com filtros
   */
  async getHistorico(
    filters: ListCaixaHistoricoFilters = {}
  ): Promise<ListCaixaHistoricoResponse> {
    console.log('[caixa-service] Listando histórico:', filters);
    const params = {
      page: filters.page || 1,
      page_size: filters.page_size || 20,
      ...(filters.data_inicio && { data_inicio: filters.data_inicio }),
      ...(filters.data_fim && { data_fim: filters.data_fim }),
      ...(filters.usuario_id && { usuario_id: filters.usuario_id }),
    };
    const { data } = await api.get<ListCaixaHistoricoResponse>(
      CAIXA_ENDPOINTS.historico,
      { params }
    );
    return data;
  },

  /**
   * Busca totais do caixa atual
   */
  async getTotais(): Promise<TotaisCaixaResponse> {
    console.log('[caixa-service] Buscando totais do caixa');
    const { data } = await api.get<TotaisCaixaResponse>(CAIXA_ENDPOINTS.totais);
    return data;
  },

  /**
   * Busca caixa por ID
   */
  async getCaixaById(id: string): Promise<CaixaDiarioResponse> {
    console.log('[caixa-service] Buscando caixa por ID:', id);
    const { data } = await api.get<CaixaDiarioResponse>(
      CAIXA_ENDPOINTS.byId(id)
    );
    return data;
  },
};

// Export default para import simplificado
export default caixaService;
