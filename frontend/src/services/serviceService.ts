import { api } from '@/lib/axios';
import {
    CreateServiceDTO,
    Service,
    ServiceFilters,
    ServiceListResponse,
    ServiceStats,
    UpdateServiceDTO,
} from '@/types/service';

export const serviceService = {
  getAll: async (filters?: ServiceFilters): Promise<ServiceListResponse> => {
    const params = new URLSearchParams();
    
    if (filters) {
      if (filters.apenas_ativos) params.append('apenas_ativos', 'true');
      if (filters.categoria_id) params.append('categoria_id', filters.categoria_id);
      if (filters.profissional_id) params.append('profissional_id', filters.profissional_id);
      if (filters.search) params.append('search', filters.search);
      if (filters.order_by) params.append('order_by', filters.order_by);
    }

    const { data } = await api.get<ServiceListResponse>(`/servicos?${params.toString()}`);
    return data;
  },

  getById: async (id: string): Promise<Service> => {
    const { data } = await api.get<Service>(`/servicos/${id}`);
    return data;
  },

  getStats: async (): Promise<ServiceStats> => {
    const { data } = await api.get<ServiceStats>('/servicos/stats');
    return data;
  },

  create: async (data: CreateServiceDTO): Promise<Service> => {
    const { data: response } = await api.post<Service>('/servicos', data);
    return response;
  },

  update: async (id: string, data: UpdateServiceDTO): Promise<Service> => {
    const { data: response } = await api.put<Service>(`/servicos/${id}`, data);
    return response;
  },

  delete: async (id: string): Promise<void> => {
    console.log('[serviceService.delete] Iniciando deleção do serviço:', id);
    try {
      const response = await api.delete(`/servicos/${id}`);
      console.log('[serviceService.delete] Sucesso! Response:', response.status, response.data);
    } catch (error: unknown) {
      console.error('[serviceService.delete] ERRO ao deletar serviço:', id);
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { status?: number; data?: unknown; headers?: unknown } };
        console.error('[serviceService.delete] Status:', axiosError.response?.status);
        console.error('[serviceService.delete] Response Data:', JSON.stringify(axiosError.response?.data, null, 2));
        console.error('[serviceService.delete] Response Headers:', axiosError.response?.headers);
      }
      if (error && typeof error === 'object' && 'config' in error) {
        const axiosError = error as { config?: { url?: string; method?: string; headers?: unknown } };
        console.error('[serviceService.delete] Request URL:', axiosError.config?.url);
        console.error('[serviceService.delete] Request Method:', axiosError.config?.method);
        console.error('[serviceService.delete] Request Headers:', axiosError.config?.headers);
      }
      throw error;
    }
  },

  toggleStatus: async (id: string): Promise<Service> => {
    const { data } = await api.patch<Service>(`/servicos/${id}/toggle-status`);
    return data;
  },
};
