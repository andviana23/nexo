/**
 * NEXO - Sistema de Gestão para Barbearias
 * Customer Service
 *
 * Serviço de clientes - comunicação com API de customers do backend.
 * Seguindo padrões do projeto e FLUXO_CADASTROS_CLIENTE.md
 */

import { api } from '@/lib/axios';
import type {
    CreateCustomerRequest,
    CustomerExportResponse,
    CustomerResponse,
    CustomerStatsResponse,
    CustomerWithHistory,
    ListCustomersFilters,
    ListCustomersResponse,
    UpdateCustomerRequest,
} from '@/types/customer'; // =============================================================================
// ENDPOINTS
// =============================================================================

const CUSTOMER_ENDPOINTS = {
  list: '/customers',
  create: '/customers',
  getById: (id: string) => `/customers/${id}`,
  update: (id: string) => `/customers/${id}`,
  inactivate: (id: string) => `/customers/${id}`,
  getWithHistory: (id: string) => `/customers/${id}/history`,
  export: (id: string) => `/customers/${id}/export`,
  stats: '/customers/stats',
  search: '/customers/search',
  checkPhone: '/customers/check-phone',
  checkCpf: '/customers/check-cpf',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const customerService = {
  /**
   * Lista clientes com filtros e paginação
   * RN-CLI-001: Lista com busca, filtros e ordenação
   */
  async list(filters: ListCustomersFilters = {}): Promise<ListCustomersResponse> {
    const params: Record<string, unknown> = {};

    if (filters.page) params.page = filters.page;
    if (filters.page_size) params.page_size = filters.page_size;
    if (filters.search) params.search = filters.search;
    if (filters.ativo !== undefined) params.ativo = filters.ativo;
    if (filters.genero) params.genero = filters.genero;
    if (filters.tag) params.tag = filters.tag;
    if (filters.tags) params.tags = filters.tags;
    if (filters.data_nascimento_inicio) params.data_nascimento_inicio = filters.data_nascimento_inicio;
    if (filters.data_nascimento_fim) params.data_nascimento_fim = filters.data_nascimento_fim;
    if (filters.order_by) params.order_by = filters.order_by;
    if (filters.order_direction) params.order_direction = filters.order_direction;

    const response = await api.get<ListCustomersResponse>(
      CUSTOMER_ENDPOINTS.list,
      { params }
    );
    return response.data;
  },

  /**
   * Busca um cliente pelo ID
   */
  async getById(id: string): Promise<CustomerResponse> {
    const response = await api.get<CustomerResponse>(
      CUSTOMER_ENDPOINTS.getById(id)
    );
    return response.data;
  },

  /**
   * Busca cliente com histórico completo (agendamentos, comandas)
   * RN-CLI-008: Histórico completo do cliente
   */
  async getWithHistory(id: string): Promise<CustomerWithHistory> {
    const response = await api.get<CustomerWithHistory>(
      CUSTOMER_ENDPOINTS.getWithHistory(id)
    );
    return response.data;
  },

  /**
   * Cria um novo cliente
   * RN-CLI-002: Cadastro com validações
   */
  async create(data: CreateCustomerRequest): Promise<CustomerResponse> {
    const response = await api.post<CustomerResponse>(
      CUSTOMER_ENDPOINTS.create,
      data
    );
    return response.data;
  },

  /**
   * Atualiza um cliente existente
   */
  async update(id: string, data: UpdateCustomerRequest): Promise<CustomerResponse> {
    const response = await api.put<CustomerResponse>(
      CUSTOMER_ENDPOINTS.update(id),
      data
    );
    return response.data;
  },

  /**
   * Inativa um cliente (soft delete)
   * RN-CLI-006: Inativação com confirmação
   */
  async inactivate(id: string): Promise<void> {
    await api.delete(CUSTOMER_ENDPOINTS.inactivate(id));
  },

  /**
   * Busca clientes por termo (nome, telefone, CPF)
   * RN-CLI-001: Busca com debounce de 300ms no frontend
   */
  async search(term: string, limit: number = 10): Promise<CustomerResponse[]> {
    const response = await api.get<CustomerResponse[]>(
      CUSTOMER_ENDPOINTS.search,
      { params: { q: term, limit } }
    );
    return response.data;
  },

  /**
   * Exporta dados do cliente (LGPD)
   * RN-CLI-009: Export completo para conformidade LGPD
   */
  async exportData(id: string): Promise<CustomerExportResponse> {
    const response = await api.get<CustomerExportResponse>(
      CUSTOMER_ENDPOINTS.export(id)
    );
    return response.data;
  },

  /**
   * Obtém estatísticas de clientes do tenant
   * RN-CLI-010: Dashboard com métricas
   */
  async getStats(): Promise<CustomerStatsResponse> {
    const response = await api.get<CustomerStatsResponse>(
      CUSTOMER_ENDPOINTS.stats
    );
    return response.data;
  },

  /**
   * Verifica se telefone já existe no tenant
   * RN-CLI-003: Validação em tempo real
   */
  async checkPhoneExists(phone: string, excludeId?: string): Promise<boolean> {
    try {
      const params: Record<string, string> = { phone };
      if (excludeId) params.exclude_id = excludeId;

      const response = await api.get<{ exists: boolean }>(
        CUSTOMER_ENDPOINTS.checkPhone,
        { params }
      );
      return response.data.exists;
    } catch {
      // Em caso de erro, permite continuar
      return false;
    }
  },

  /**
   * Verifica se CPF já existe no tenant
   * RN-CLI-004: Validação em tempo real (CPF opcional)
   */
  async checkCpfExists(cpf: string, excludeId?: string): Promise<boolean> {
    try {
      const params: Record<string, string> = { cpf };
      if (excludeId) params.exclude_id = excludeId;

      const response = await api.get<{ exists: boolean }>(
        CUSTOMER_ENDPOINTS.checkCpf,
        { params }
      );
      return response.data.exists;
    } catch {
      // Em caso de erro, permite continuar
      return false;
    }
  },

  /**
   * Lista apenas clientes ativos (para selects e autocomplete)
   */
  async listActive(limit: number = 100): Promise<CustomerResponse[]> {
    const response = await this.list({
      ativo: true,
      page_size: limit,
    });
    return response.data;
  },

  /**
   * Lista clientes que fazem aniversário em um período
   * Útil para campanhas de marketing
   */
  async listBirthdays(startDate: string, endDate: string): Promise<CustomerResponse[]> {
    const response = await this.list({
      data_nascimento_inicio: startDate,
      data_nascimento_fim: endDate,
      ativo: true,
      page_size: 100,
    });
    return response.data;
  },

  /**
   * Lista clientes por tag
   * Útil para segmentação
   */
  async listByTag(tag: string): Promise<CustomerResponse[]> {
    const response = await this.list({
      tag,
      ativo: true,
      page_size: 100,
    });
    return response.data;
  },
};

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Formata telefone para exibição (XX) XXXXX-XXXX
 */
export function formatPhone(phone: string): string {
  const cleaned = phone.replace(/\D/g, '');
  if (cleaned.length === 11) {
    return `(${cleaned.slice(0, 2)}) ${cleaned.slice(2, 7)}-${cleaned.slice(7)}`;
  }
  if (cleaned.length === 10) {
    return `(${cleaned.slice(0, 2)}) ${cleaned.slice(2, 6)}-${cleaned.slice(6)}`;
  }
  return phone;
}

/**
 * Formata CPF para exibição XXX.XXX.XXX-XX
 */
export function formatCPF(cpf: string): string {
  const cleaned = cpf.replace(/\D/g, '');
  if (cleaned.length === 11) {
    return `${cleaned.slice(0, 3)}.${cleaned.slice(3, 6)}.${cleaned.slice(6, 9)}-${cleaned.slice(9)}`;
  }
  return cpf;
}

/**
 * Remove formatação do telefone
 */
export function cleanPhone(phone: string): string {
  return phone.replace(/\D/g, '');
}

/**
 * Remove formatação do CPF
 */
export function cleanCPF(cpf: string): string {
  return cpf.replace(/\D/g, '');
}

/**
 * Valida CPF (algoritmo oficial)
 */
export function isValidCPF(cpf: string): boolean {
  const cleaned = cpf.replace(/\D/g, '');

  if (cleaned.length !== 11) return false;

  // Verifica se todos os dígitos são iguais
  if (/^(\d)\1{10}$/.test(cleaned)) return false;

  // Validação dos dígitos verificadores
  let sum = 0;
  for (let i = 0; i < 9; i++) {
    sum += parseInt(cleaned.charAt(i)) * (10 - i);
  }
  let remainder = (sum * 10) % 11;
  if (remainder === 10 || remainder === 11) remainder = 0;
  if (remainder !== parseInt(cleaned.charAt(9))) return false;

  sum = 0;
  for (let i = 0; i < 10; i++) {
    sum += parseInt(cleaned.charAt(i)) * (11 - i);
  }
  remainder = (sum * 10) % 11;
  if (remainder === 10 || remainder === 11) remainder = 0;
  if (remainder !== parseInt(cleaned.charAt(10))) return false;

  return true;
}

/**
 * Valida telefone brasileiro (10 ou 11 dígitos)
 */
export function isValidPhone(phone: string): boolean {
  const cleaned = phone.replace(/\D/g, '');
  return cleaned.length === 10 || cleaned.length === 11;
}

/**
 * Formata nome para exibição (primeira letra maiúscula)
 */
export function formatCustomerName(name: string): string {
  return name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

/**
 * Calcula idade a partir da data de nascimento
 */
export function calculateAge(birthDate: string | null | undefined): number | null {
  if (!birthDate) return null;

  const birth = new Date(birthDate);
  const today = new Date();

  let age = today.getFullYear() - birth.getFullYear();
  const monthDiff = today.getMonth() - birth.getMonth();

  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--;
  }

  return age;
}

/**
 * Verifica se cliente faz aniversário hoje
 */
export function isBirthdayToday(birthDate: string | null | undefined): boolean {
  if (!birthDate) return false;

  const birth = new Date(birthDate);
  const today = new Date();

  return (
    birth.getMonth() === today.getMonth() &&
    birth.getDate() === today.getDate()
  );
}

/**
 * Verifica se cliente faz aniversário nos próximos N dias
 */
export function isBirthdaySoon(
  birthDate: string | null | undefined,
  days: number = 7
): boolean {
  if (!birthDate) return false;

  const birth = new Date(birthDate);
  const today = new Date();

  // Ajusta a data de nascimento para o ano atual
  const birthdayThisYear = new Date(
    today.getFullYear(),
    birth.getMonth(),
    birth.getDate()
  );

  // Se já passou, verifica para o próximo ano
  if (birthdayThisYear < today) {
    birthdayThisYear.setFullYear(today.getFullYear() + 1);
  }

  const diffTime = birthdayThisYear.getTime() - today.getTime();
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

  return diffDays >= 0 && diffDays <= days;
}

/**
 * Obtém iniciais do nome para avatar
 */
export function getInitials(name: string): string {
  const parts = name.trim().split(' ').filter(Boolean);
  if (parts.length === 0) return '??';
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
  return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
}

/**
 * Gera cor de avatar baseada no nome (determinística)
 */
export function getAvatarColor(name: string): string {
  const colors = [
    'bg-blue-500',
    'bg-green-500',
    'bg-yellow-500',
    'bg-red-500',
    'bg-purple-500',
    'bg-pink-500',
    'bg-indigo-500',
    'bg-teal-500',
  ];

  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
}

// =============================================================================
// TIPOS DE ERRO ESPECÍFICOS
// =============================================================================

export class CustomerError extends Error {
  constructor(
    message: string,
    public code: string = 'CUSTOMER_ERROR'
  ) {
    super(message);
    this.name = 'CustomerError';
  }
}

export class DuplicatePhoneError extends CustomerError {
  constructor() {
    super('Este telefone já está cadastrado para outro cliente.', 'DUPLICATE_PHONE');
    this.name = 'DuplicatePhoneError';
  }
}

export class DuplicateCpfError extends CustomerError {
  constructor() {
    super('Este CPF já está cadastrado para outro cliente.', 'DUPLICATE_CPF');
    this.name = 'DuplicateCpfError';
  }
}

export class CustomerNotFoundError extends CustomerError {
  constructor() {
    super('Cliente não encontrado.', 'CUSTOMER_NOT_FOUND');
    this.name = 'CustomerNotFoundError';
  }
}

export class InvalidPhoneError extends CustomerError {
  constructor() {
    super('Telefone inválido. Use o formato (XX) XXXXX-XXXX.', 'INVALID_PHONE');
    this.name = 'InvalidPhoneError';
  }
}

export class InvalidCpfError extends CustomerError {
  constructor() {
    super('CPF inválido.', 'INVALID_CPF');
    this.name = 'InvalidCpfError';
  }
}

export class CustomerHasAppointmentsError extends CustomerError {
  constructor() {
    super(
      'Não é possível excluir cliente com agendamentos. Altere o status para Inativo.',
      'HAS_APPOINTMENTS'
    );
    this.name = 'CustomerHasAppointmentsError';
  }
}

export class RequiredFieldError extends CustomerError {
  constructor(field: string) {
    super(`O campo ${field} é obrigatório.`, 'REQUIRED_FIELD');
    this.name = 'RequiredFieldError';
  }
}
