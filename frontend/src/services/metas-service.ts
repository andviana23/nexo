/**
 * NEXO - Sistema de Gestão para Barbearias
 * Metas Service
 *
 * Serviço de comunicação com a API de Metas do backend.
 * Endpoints mapeados de: backend/internal/infra/http/handler/metas_handler.go
 */

import { api } from '@/lib/axios';
import type {
    MetaBarbeiroResponse,
    MetaMensalResponse,
    MetasFilters,
    MetaTicketResponse,
    ResumoMetas,
    SetMetaBarbeiroRequest,
    SetMetaMensalRequest,
    SetMetaTicketRequest,
    UpdateMetaBarbeiroRequest,
    UpdateMetaMensalRequest,
    UpdateMetaTicketRequest,
} from '@/types/metas';
import {
    calcularPercentual,
    extenderMetaBarbeiro,
    parseMoneyValue,
} from '@/types/metas';

// =============================================================================
// ENDPOINTS
// =============================================================================

const METAS_ENDPOINTS = {
  // Metas Mensais
  monthly: '/metas/monthly',
  monthlyById: (id: string) => `/metas/monthly/${id}`,

  // Metas Barbeiro
  barbers: '/metas/barbers',
  barbersById: (id: string) => `/metas/barbers/${id}`,

  // Metas Ticket Médio
  ticket: '/metas/ticket',
  ticketById: (id: string) => `/metas/ticket/${id}`,
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const metasService = {
  // ===========================================================================
  // METAS MENSAIS (Faturamento Geral)
  // ===========================================================================

  /**
   * Lista todas as metas mensais
   */
  async listMetasMensais(filters?: MetasFilters): Promise<MetaMensalResponse[]> {
    console.log('[metas-service] Listando metas mensais:', filters);
    const { data } = await api.get<MetaMensalResponse[]>(METAS_ENDPOINTS.monthly, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca uma meta mensal por ID
   */
  async getMetaMensal(id: string): Promise<MetaMensalResponse> {
    console.log('[metas-service] Buscando meta mensal:', id);
    const { data } = await api.get<MetaMensalResponse>(METAS_ENDPOINTS.monthlyById(id));
    return data;
  },

  /**
   * Cria nova meta mensal
   */
  async createMetaMensal(payload: SetMetaMensalRequest): Promise<MetaMensalResponse> {
    console.log('[metas-service] Criando meta mensal:', payload);
    const { data } = await api.post<MetaMensalResponse>(METAS_ENDPOINTS.monthly, payload);
    return data;
  },

  /**
   * Atualiza meta mensal
   */
  async updateMetaMensal(id: string, payload: UpdateMetaMensalRequest): Promise<MetaMensalResponse> {
    console.log('[metas-service] Atualizando meta mensal:', id, payload);
    const { data } = await api.put<MetaMensalResponse>(METAS_ENDPOINTS.monthlyById(id), payload);
    return data;
  },

  /**
   * Deleta meta mensal
   */
  async deleteMetaMensal(id: string): Promise<void> {
    console.log('[metas-service] Deletando meta mensal:', id);
    await api.delete(METAS_ENDPOINTS.monthlyById(id));
  },

  // ===========================================================================
  // METAS POR BARBEIRO
  // ===========================================================================

  /**
   * Lista metas de barbeiros
   */
  async listMetasBarbeiro(barbeiroId?: string): Promise<MetaBarbeiroResponse[]> {
    console.log('[metas-service] Listando metas barbeiro:', barbeiroId);
    const params = barbeiroId ? { barbeiro_id: barbeiroId } : {};
    const { data } = await api.get<MetaBarbeiroResponse[]>(METAS_ENDPOINTS.barbers, { params });
    return data;
  },

  /**
   * Busca meta de barbeiro por ID
   */
  async getMetaBarbeiro(id: string): Promise<MetaBarbeiroResponse> {
    console.log('[metas-service] Buscando meta barbeiro:', id);
    const { data } = await api.get<MetaBarbeiroResponse>(METAS_ENDPOINTS.barbersById(id));
    return data;
  },

  /**
   * Cria nova meta de barbeiro
   */
  async createMetaBarbeiro(payload: SetMetaBarbeiroRequest): Promise<MetaBarbeiroResponse> {
    console.log('[metas-service] Criando meta barbeiro:', payload);
    const { data } = await api.post<MetaBarbeiroResponse>(METAS_ENDPOINTS.barbers, payload);
    return data;
  },

  /**
   * Atualiza meta de barbeiro
   */
  async updateMetaBarbeiro(id: string, payload: UpdateMetaBarbeiroRequest): Promise<MetaBarbeiroResponse> {
    console.log('[metas-service] Atualizando meta barbeiro:', id, payload);
    const { data } = await api.put<MetaBarbeiroResponse>(METAS_ENDPOINTS.barbersById(id), payload);
    return data;
  },

  /**
   * Deleta meta de barbeiro
   */
  async deleteMetaBarbeiro(id: string): Promise<void> {
    console.log('[metas-service] Deletando meta barbeiro:', id);
    await api.delete(METAS_ENDPOINTS.barbersById(id));
  },

  // ===========================================================================
  // METAS TICKET MÉDIO
  // ===========================================================================

  /**
   * Lista metas de ticket médio
   */
  async listMetasTicket(): Promise<MetaTicketResponse[]> {
    console.log('[metas-service] Listando metas ticket');
    const { data } = await api.get<MetaTicketResponse[]>(METAS_ENDPOINTS.ticket);
    return data;
  },

  /**
   * Busca meta de ticket por ID
   */
  async getMetaTicket(id: string): Promise<MetaTicketResponse> {
    console.log('[metas-service] Buscando meta ticket:', id);
    const { data } = await api.get<MetaTicketResponse>(METAS_ENDPOINTS.ticketById(id));
    return data;
  },

  /**
   * Cria nova meta de ticket
   */
  async createMetaTicket(payload: SetMetaTicketRequest): Promise<MetaTicketResponse> {
    console.log('[metas-service] Criando meta ticket:', payload);
    const { data } = await api.post<MetaTicketResponse>(METAS_ENDPOINTS.ticket, payload);
    return data;
  },

  /**
   * Atualiza meta de ticket
   */
  async updateMetaTicket(id: string, payload: UpdateMetaTicketRequest): Promise<MetaTicketResponse> {
    console.log('[metas-service] Atualizando meta ticket:', id, payload);
    const { data } = await api.put<MetaTicketResponse>(METAS_ENDPOINTS.ticketById(id), payload);
    return data;
  },

  /**
   * Deleta meta de ticket
   */
  async deleteMetaTicket(id: string): Promise<void> {
    console.log('[metas-service] Deletando meta ticket:', id);
    await api.delete(METAS_ENDPOINTS.ticketById(id));
  },

  // ===========================================================================
  // HELPERS / AGREGADORES
  // ===========================================================================

  /**
   * Busca resumo geral de metas para o mês atual
   */
  async getResumoMetas(mesAno?: string): Promise<ResumoMetas> {
    console.log('[metas-service] Calculando resumo metas:', mesAno);

    // Busca paralela
    const [metasMensais, metasBarbeiro, metasTicket] = await Promise.all([
      this.listMetasMensais(),
      this.listMetasBarbeiro(),
      this.listMetasTicket(),
    ]);

    // Filtra pelo mês/ano se informado
    const mesAnoFiltro = mesAno || new Date().toISOString().slice(0, 7);

    const metaMensal = metasMensais.find((m) => m.mes_ano === mesAnoFiltro);
    const metasBarbeiroFiltradas = metasBarbeiro.filter((m) => m.mes_ano === mesAnoFiltro);
    const metaTicketGeral = metasTicket.find(
      (m) => m.mes_ano === mesAnoFiltro && m.tipo === 'GERAL'
    );

    // Calcula barbeiros acima da meta
    const barbeirosEstendidos = metasBarbeiroFiltradas.map(extenderMetaBarbeiro);
    const barbeirosAcimaMeta = barbeirosEstendidos.filter((b) => b.percentual_total >= 100).length;

    return {
      mes_ano: mesAnoFiltro,
      meta_faturamento: metaMensal ? parseMoneyValue(metaMensal.meta_faturamento) : 0,
      realizado_faturamento: metaMensal ? parseMoneyValue(metaMensal.realizado) : 0,
      percentual_faturamento: metaMensal
        ? calcularPercentual(
            parseMoneyValue(metaMensal.realizado),
            parseMoneyValue(metaMensal.meta_faturamento)
          )
        : 0,
      total_barbeiros_com_meta: metasBarbeiroFiltradas.length,
      barbeiros_acima_meta: barbeirosAcimaMeta,
      ticket_medio_meta: metaTicketGeral ? parseMoneyValue(metaTicketGeral.meta_valor) : 0,
      ticket_medio_realizado: metaTicketGeral
        ? parseMoneyValue(metaTicketGeral.ticket_medio_realizado)
        : 0,
      percentual_ticket: metaTicketGeral
        ? calcularPercentual(
            parseMoneyValue(metaTicketGeral.ticket_medio_realizado),
            parseMoneyValue(metaTicketGeral.meta_valor)
          )
        : 0,
    };
  },

  /**
   * Busca ranking de barbeiros por percentual de atingimento
   */
  async getRankingBarbeiros(mesAno?: string) {
    const metasBarbeiro = await this.listMetasBarbeiro();
    const mesAnoFiltro = mesAno || new Date().toISOString().slice(0, 7);

    const metasFiltradas = metasBarbeiro.filter((m) => m.mes_ano === mesAnoFiltro);
    const barbeirosEstendidos = metasFiltradas.map(extenderMetaBarbeiro);

    // Ordena por percentual total decrescente
    return barbeirosEstendidos
      .sort((a, b) => b.percentual_total - a.percentual_total)
      .map((b, index) => ({
        posicao: index + 1,
        barbeiro_id: b.barbeiro_id,
        barbeiro_nome: b.barbeiro_nome || 'Barbeiro',
        percentual_total: b.percentual_total,
        nivel_bonificacao: b.nivel_bonificacao,
        realizado_total: b.realizado_total,
        bonus_valor: b.bonus_valor,
      }));
  },

  /**
   * Busca histórico de metas mensais (últimos 6 meses)
   */
  async getHistoricoMetas(meses: number = 6): Promise<MetaMensalResponse[]> {
    const todasMetas = await this.listMetasMensais();

    // Ordena por mes_ano decrescente e pega os últimos N
    return todasMetas
      .sort((a, b) => b.mes_ano.localeCompare(a.mes_ano))
      .slice(0, meses);
  },
};

export default metasService;
