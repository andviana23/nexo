/**
 * NEXO - Serviço de Meios de Pagamento
 * API Client para módulo de Tipos de Recebimento
 */

import { api } from '@/lib/axios';
import {
    CreateMeioPagamentoDTO,
    MeioPagamento,
    MeioPagamentoFilters,
    MeioPagamentoListResponse,
    UpdateMeioPagamentoDTO,
} from '@/types/meio-pagamento';

const BASE_URL = '/meios-pagamento';

export const meioPagamentoService = {
  /**
   * Lista todos os meios de pagamento
   */
  list: async (filters?: MeioPagamentoFilters): Promise<MeioPagamentoListResponse> => {
    const response = await api.get<MeioPagamentoListResponse>(BASE_URL, {
      params: filters,
    });
    return response.data;
  },

  /**
   * Busca um meio de pagamento por ID
   */
  getById: async (id: string): Promise<MeioPagamento> => {
    const response = await api.get<MeioPagamento>(`${BASE_URL}/${id}`);
    return response.data;
  },

  /**
   * Cria um novo meio de pagamento
   */
  create: async (data: CreateMeioPagamentoDTO): Promise<MeioPagamento> => {
    const response = await api.post<MeioPagamento>(BASE_URL, data);
    return response.data;
  },

  /**
   * Atualiza um meio de pagamento existente
   */
  update: async (id: string, data: UpdateMeioPagamentoDTO): Promise<MeioPagamento> => {
    const response = await api.put<MeioPagamento>(`${BASE_URL}/${id}`, data);
    return response.data;
  },

  /**
   * Remove um meio de pagamento
   */
  delete: async (id: string): Promise<void> => {
    await api.delete(`${BASE_URL}/${id}`);
  },

  /**
   * Alterna o status ativo/inativo
   */
  toggle: async (id: string): Promise<MeioPagamento> => {
    const response = await api.patch<MeioPagamento>(`${BASE_URL}/${id}/toggle`);
    return response.data;
  },
};
