/**
 * NEXO - Sistema de Gestão para Barbearias
 * Professional Service
 *
 * Serviço de profissionais - comunicação com API de professionals do backend.
 * Seguindo padrões de appointment-service.ts
 */

import { api } from '@/lib/axios';
import type {
    CreateProfessionalRequest,
    ListProfessionalsFilters,
    ListProfessionalsResponse,
    ProfessionalResponse,
    ProfessionalType,
    UpdateProfessionalRequest,
    UpdateProfessionalStatusRequest,
} from '@/types/professional';
import axios from 'axios';

// =============================================================================
// ENDPOINTS
// =============================================================================

const PROFESSIONAL_ENDPOINTS = {
  list: '/professionals',
  create: '/professionals',
  getById: (id: string) => `/professionals/${id}`,
  update: (id: string) => `/professionals/${id}`,
  delete: (id: string) => `/professionals/${id}`,
  updateStatus: (id: string) => `/professionals/${id}/status`,
  checkEmail: '/professionals/check-email',
  checkCpf: '/professionals/check-cpf',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

// Tipo de resposta da API (formato do backend)
interface ApiListProfessionalsResponse {
  data: ProfessionalResponse[];
  page: number;
  page_size: number;
  total: number;
}

export const professionalService = {
  /**
   * Lista profissionais com filtros e paginação
   */
  async list(filters: ListProfessionalsFilters = {}): Promise<ListProfessionalsResponse> {
    const params: Record<string, unknown> = {};
    
    if (filters.page) params.page = filters.page;
    if (filters.page_size) params.page_size = filters.page_size;
    if (filters.search) params.search = filters.search;
    if (filters.tipo) params.tipo = filters.tipo;
    if (filters.status) params.status = filters.status;
    if (filters.order_by) params.order_by = filters.order_by;
    if (filters.order_direction) params.order_direction = filters.order_direction;
    
    const response = await api.get<ApiListProfessionalsResponse>(
      PROFESSIONAL_ENDPOINTS.list,
      { params }
    );
    
    // Adaptar resposta da API para o formato esperado pelo frontend
    const apiData = response.data;
    const pageSize = apiData.page_size || 20;
    const total = apiData.total || 0;
    const totalPages = pageSize > 0 ? Math.ceil(total / pageSize) : 1;
    
    return {
      data: apiData.data || [],
      meta: {
        page: apiData.page || 1,
        page_size: pageSize,
        total: total,
        total_pages: totalPages,
      },
    };
  },

  /**
   * Busca um profissional pelo ID
   */
  async getById(id: string): Promise<ProfessionalResponse> {
    const response = await api.get<{ data: ProfessionalResponse }>(
      PROFESSIONAL_ENDPOINTS.getById(id)
    );
    return response.data.data;
  },

  /**
   * Cria um novo profissional
   */
  async create(data: CreateProfessionalRequest): Promise<ProfessionalResponse> {
    try {
      const response = await api.post<{ data: ProfessionalResponse }>(
        PROFESSIONAL_ENDPOINTS.create,
        data
      );
      return response.data.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        const status = error.response.status;
        const message = error.response.data?.message || '';

        // 409 Conflict - CPF ou Email duplicado
        if (status === 409) {
          if (message.toLowerCase().includes('cpf') || message.toLowerCase().includes('cnpj')) {
            throw new DuplicateCpfError();
          }
          if (message.toLowerCase().includes('email')) {
            throw new DuplicateEmailError();
          }
          // Fallback genérico para 409
          throw new DuplicateCpfError();
        }

        // 400 Bad Request - Validação
        if (status === 400) {
          if (message.toLowerCase().includes('comissão') || message.toLowerCase().includes('comissao')) {
            throw new InvalidCommissionError();
          }
        }

        // 404 Not Found
        if (status === 404) {
          throw new ProfessionalNotFoundError();
        }
      }
      throw error;
    }
  },

  /**
   * Atualiza um profissional existente
   */
  async update(id: string, data: UpdateProfessionalRequest): Promise<ProfessionalResponse> {
    const response = await api.put<{ data: ProfessionalResponse }>(
      PROFESSIONAL_ENDPOINTS.update(id),
      data
    );
    return response.data.data;
  },

  /**
   * Remove um profissional (soft delete)
   * Muda status para DEMITIDO
   */
  async delete(id: string): Promise<void> {
    await api.delete(PROFESSIONAL_ENDPOINTS.delete(id));
  },

  /**
   * Atualiza o status de um profissional
   */
  async updateStatus(
    id: string,
    data: UpdateProfessionalStatusRequest
  ): Promise<ProfessionalResponse> {
    const response = await api.put<{ data: ProfessionalResponse }>(
      PROFESSIONAL_ENDPOINTS.updateStatus(id),
      data
    );
    return response.data.data;
  },

  /**
   * Verifica se email já existe no tenant
   */
  async checkEmailExists(email: string, excludeId?: string): Promise<boolean> {
    try {
      const response = await api.get<{ exists: boolean }>(
        PROFESSIONAL_ENDPOINTS.checkEmail,
        { params: { email, exclude_id: excludeId } }
      );
      return response.data.exists;
    } catch {
      // Em caso de erro, permite continuar
      return false;
    }
  },

  /**
   * Verifica se CPF já existe no tenant
   */
  async checkCpfExists(cpf: string, excludeId?: string): Promise<boolean> {
    try {
      const response = await api.get<{ exists: boolean }>(
        PROFESSIONAL_ENDPOINTS.checkCpf,
        { params: { cpf, exclude_id: excludeId } }
      );
      return response.data.exists;
    } catch {
      // Em caso de erro, permite continuar
      return false;
    }
  },

  /**
   * Lista apenas profissionais ativos (para selects)
   */
  async listActive(): Promise<ProfessionalResponse[]> {
    const response = await this.list({ 
      status: 'ATIVO',
      page_size: 100 
    });
    return response.data;
  },

  /**
   * Lista apenas barbeiros ativos (para agendamentos)
   */
  async listBarbers(): Promise<ProfessionalResponse[]> {
    const response = await this.list({ 
      tipo: 'BARBEIRO',
      status: 'ATIVO',
      page_size: 100 
    });
    return response.data;
  },
};

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Mapeia tipo de profissional para role do sistema
 * RN-PROF-005
 */
export function mapProfessionalTypeToRole(tipo: ProfessionalType): string {
  const mapping: Record<ProfessionalType, string> = {
    BARBEIRO: 'barbeiro',
    GERENTE: 'manager',
    RECEPCIONISTA: 'recepcionista',
    OUTRO: 'staff',
  };
  return mapping[tipo];
}

/**
 * Gera senha temporária segura (8 caracteres alfanuméricos)
 */
export function generateTemporaryPassword(): string {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  let password = '';
  for (let i = 0; i < 8; i++) {
    password += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return password;
}

/**
 * Formata nome para exibição (primeira letra maiúscula)
 */
export function formatProfessionalName(name: string): string {
  return name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

/**
 * Verifica se profissional pode receber comissão
 */
export function canReceiveCommission(
  tipo: ProfessionalType,
  tambemBarbeiro?: boolean
): boolean {
  return tipo === 'BARBEIRO' || (tipo === 'GERENTE' && tambemBarbeiro === true);
}

// =============================================================================
// TIPOS DE ERRO ESPECÍFICOS
// =============================================================================

export class ProfessionalError extends Error {
  constructor(
    message: string,
    public code: string = 'PROFESSIONAL_ERROR'
  ) {
    super(message);
    this.name = 'ProfessionalError';
  }
}

export class DuplicateEmailError extends ProfessionalError {
  constructor() {
    super('Este email já está cadastrado para outro profissional.', 'DUPLICATE_EMAIL');
    this.name = 'DuplicateEmailError';
  }
}

export class DuplicateCpfError extends ProfessionalError {
  constructor() {
    super('Este CPF já está cadastrado para outro profissional.', 'DUPLICATE_CPF');
    this.name = 'DuplicateCpfError';
  }
}

export class ProfessionalNotFoundError extends ProfessionalError {
  constructor() {
    super('Profissional não encontrado.', 'PROFESSIONAL_NOT_FOUND');
    this.name = 'ProfessionalNotFoundError';
  }
}

export class InvalidCommissionError extends ProfessionalError {
  constructor() {
    super('Comissão é obrigatória para barbeiros (0-100%).', 'INVALID_COMMISSION');
    this.name = 'InvalidCommissionError';
  }
}

export class ProfessionalHasAppointmentsError extends ProfessionalError {
  constructor() {
    super(
      'Não é possível excluir profissional com agendamentos. Altere o status para Inativo.',
      'HAS_APPOINTMENTS'
    );
    this.name = 'ProfessionalHasAppointmentsError';
  }
}

export class InvalidStatusTransitionError extends ProfessionalError {
  constructor() {
    super('Transição de status inválida.', 'INVALID_STATUS_TRANSITION');
    this.name = 'InvalidStatusTransitionError';
  }
}
