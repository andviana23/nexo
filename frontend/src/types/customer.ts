/**
 * NEXO - Sistema de Gestão para Barbearias
 * Types: Customer (Cliente)
 * 
 * @module types/customer
 * @description Definições de tipos para clientes no NEXO
 * Conforme FLUXO_CADASTROS_CLIENTE.md
 */

// ============================================================================
// ENUMS
// ============================================================================

/** Gêneros disponíveis */
export type CustomerGender = 'M' | 'F' | 'NB' | 'PNI';

// ============================================================================
// ENTITIES
// ============================================================================

/** Entidade principal de Cliente */
export interface Customer {
  id: string;
  tenant_id: string;
  
  // Dados Básicos (obrigatórios)
  nome: string;
  telefone: string;
  
  // Dados Opcionais
  email?: string;
  cpf?: string;
  data_nascimento?: string; // ISO8601 (YYYY-MM-DD)
  genero?: CustomerGender;
  
  // Endereço
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  
  // CRM
  observacoes?: string;
  tags: string[];
  
  // Status
  ativo: boolean;
  
  // Timestamps
  criado_em: string;
  atualizado_em: string;
}

/** Resumo do cliente para listas e selects */
export interface CustomerSummary {
  id: string;
  nome: string;
  telefone: string;
  email?: string;
  tags: string[];
}

/** Cliente com histórico de atendimentos */
export interface CustomerWithHistory extends Customer {
  total_atendimentos: number;
  total_gasto: string;
  ticket_medio: string;
  ultimo_atendimento?: string;
  frequencia_media_dias?: number;
}

// ============================================================================
// DTOs - Request
// ============================================================================

/** DTO para criar cliente */
export interface CreateCustomerRequest {
  // Campos Obrigatórios
  nome: string;
  telefone: string;
  
  // Campos Opcionais
  email?: string;
  cpf?: string;
  data_nascimento?: string; // YYYY-MM-DD
  genero?: CustomerGender;
  
  // Endereço
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  
  // CRM
  observacoes?: string;
  tags?: string[];
}

/** DTO para atualizar cliente */
export interface UpdateCustomerRequest {
  nome?: string;
  telefone?: string;
  email?: string;
  cpf?: string;
  data_nascimento?: string;
  genero?: CustomerGender;
  
  // Endereço
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  
  // CRM
  observacoes?: string;
  tags?: string[];
}

/** Filtros para listagem */
export interface ListCustomersFilters {
  page?: number;
  page_size?: number;
  search?: string;
  ativo?: boolean;
  tags?: string[];
  tag?: string; // Filtro por tag única
  genero?: CustomerGender;
  data_nascimento_inicio?: string;
  data_nascimento_fim?: string;
  order_by?: 'nome' | 'criado_em' | 'atualizado_em';
  order_direction?: 'asc' | 'desc';
}

/** Parâmetros de busca rápida */
export interface SearchCustomersParams {
  q: string;
}

/** Parâmetros para verificar duplicidade */
export interface CheckExistsParams {
  telefone?: string;
  cpf?: string;
  exclude_id?: string;
}

// ============================================================================
// DTOs - Response
// ============================================================================

/** Resposta de um cliente */
export interface CustomerResponse {
  id: string;
  tenant_id: string;
  nome: string;
  telefone: string;
  email?: string;
  cpf?: string;
  data_nascimento?: string;
  genero?: CustomerGender;
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  observacoes?: string;
  tags: string[];
  ativo: boolean;
  criado_em: string;
  atualizado_em: string;
}

/** Resposta paginada de listagem */
export interface ListCustomersResponse {
  data: CustomerResponse[];
  page: number;
  page_size: number;
  total: number;
}

/** Resposta de estatísticas */
export interface CustomerStatsResponse {
  total_ativos: number;
  total_inativos: number;
  novos_ultimos_30_dias: number;
  total_geral: number;
}

/** Resposta de verificação de existência */
export interface CheckExistsResponse {
  exists: boolean;
}

// ============================================================================
// LGPD Export Types
// ============================================================================

/** Endereço para exportação */
export interface CustomerAddressExport {
  logradouro?: string;
  numero?: string;
  complemento?: string;
  bairro?: string;
  cidade?: string;
  estado?: string;
  cep?: string;
}

/** Atendimento para exportação */
export interface CustomerAppointmentExport {
  data: string;
  status: string;
  profissional: string;
  valor_total: string;
}

/** Resposta de exportação LGPD */
export interface CustomerExportResponse {
  dados_pessoais: {
    nome: string;
    email?: string;
    telefone: string;
    cpf?: string;
    data_nascimento?: string;
    genero?: string;
    endereco?: CustomerAddressExport;
  };
  historico_atendimentos: CustomerAppointmentExport[];
  metricas: {
    total_gasto: string;
    ticket_medio: string;
    total_visitas: number;
  };
  data_exportacao: string;
}

// ============================================================================
// UI State Types
// ============================================================================

/** Estado do modal de cliente */
export interface CustomerModalState {
  isOpen: boolean;
  mode: 'create' | 'edit' | 'view';
  customer?: CustomerResponse;
}

// ============================================================================
// HELPERS
// ============================================================================

/** Labels para gêneros */
export const GENDER_LABELS: Record<CustomerGender, string> = {
  M: 'Masculino',
  F: 'Feminino',
  NB: 'Não Binário',
  PNI: 'Prefiro não informar',
};

/** Tags padrão do sistema */
export const DEFAULT_TAGS = [
  'VIP',
  'Recorrente',
  'Inadimplente',
  'Novo',
  'Gastador',
  'Inativo',
] as const;

/** Estados brasileiros */
export const ESTADOS_BR = [
  { value: 'AC', label: 'Acre' },
  { value: 'AL', label: 'Alagoas' },
  { value: 'AP', label: 'Amapá' },
  { value: 'AM', label: 'Amazonas' },
  { value: 'BA', label: 'Bahia' },
  { value: 'CE', label: 'Ceará' },
  { value: 'DF', label: 'Distrito Federal' },
  { value: 'ES', label: 'Espírito Santo' },
  { value: 'GO', label: 'Goiás' },
  { value: 'MA', label: 'Maranhão' },
  { value: 'MT', label: 'Mato Grosso' },
  { value: 'MS', label: 'Mato Grosso do Sul' },
  { value: 'MG', label: 'Minas Gerais' },
  { value: 'PA', label: 'Pará' },
  { value: 'PB', label: 'Paraíba' },
  { value: 'PR', label: 'Paraná' },
  { value: 'PE', label: 'Pernambuco' },
  { value: 'PI', label: 'Piauí' },
  { value: 'RJ', label: 'Rio de Janeiro' },
  { value: 'RN', label: 'Rio Grande do Norte' },
  { value: 'RS', label: 'Rio Grande do Sul' },
  { value: 'RO', label: 'Rondônia' },
  { value: 'RR', label: 'Roraima' },
  { value: 'SC', label: 'Santa Catarina' },
  { value: 'SP', label: 'São Paulo' },
  { value: 'SE', label: 'Sergipe' },
  { value: 'TO', label: 'Tocantins' },
] as const;

/** Cores para tags padrão (Tailwind classes) */
export const TAG_COLORS: Record<string, string> = {
  VIP: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
  Recorrente: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  Inadimplente: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  Novo: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  Gastador: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  Inativo: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
};

/** Formatar telefone para exibição */
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

/** Formatar CPF para exibição */
export function formatCPF(cpf: string): string {
  const cleaned = cpf.replace(/\D/g, '');
  if (cleaned.length === 11) {
    return `${cleaned.slice(0, 3)}.${cleaned.slice(3, 6)}.${cleaned.slice(6, 9)}-${cleaned.slice(9)}`;
  }
  return cpf;
}

/** Formatar CEP para exibição */
export function formatCEP(cep: string): string {
  const cleaned = cep.replace(/\D/g, '');
  if (cleaned.length === 8) {
    return `${cleaned.slice(0, 5)}-${cleaned.slice(5)}`;
  }
  return cep;
}

/** Limpar telefone para envio */
export function cleanPhone(phone: string): string {
  return phone.replace(/\D/g, '');
}

/** Limpar CPF para envio */
export function cleanCPF(cpf: string): string {
  return cpf.replace(/\D/g, '');
}

/** Limpar CEP para envio */
export function cleanCEP(cep: string): string {
  return cep.replace(/\D/g, '');
}
