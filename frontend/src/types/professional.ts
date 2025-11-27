/**
 * NEXO - Sistema de Gestão para Barbearias
 * Types: Professional (Profissional)
 * 
 * @module types/professional
 * @description Definições de tipos para profissionais no NEXO
 */

// ============================================================================
// ENUMS
// ============================================================================

/** Tipos de profissional */
export type ProfessionalType = 'BARBEIRO' | 'GERENTE' | 'RECEPCIONISTA' | 'OUTRO';

/** Status do profissional */
export type ProfessionalStatus = 'ATIVO' | 'INATIVO' | 'AFASTADO' | 'DEMITIDO';

/** Tipo de comissão */
export type CommissionType = 'PERCENTUAL' | 'FIXO';

// ============================================================================
// HORÁRIO DE TRABALHO
// ============================================================================

/** Turno de trabalho */
export interface WorkShift {
  inicio: string; // HH:MM
  fim: string;    // HH:MM
}

/** Dia da semana */
export interface WeekDay {
  ativo: boolean;
  turnos: WorkShift[];
}

/** Horário de trabalho semanal */
export interface WorkSchedule {
  segunda: WeekDay;
  terca: WeekDay;
  quarta: WeekDay;
  quinta: WeekDay;
  sexta: WeekDay;
  sabado: WeekDay;
  domingo: WeekDay;
}

// ============================================================================
// ENTITIES
// ============================================================================

/** Entidade principal de Profissional */
export interface Professional {
  id: string;
  tenant_id: string;
  user_id?: string;
  
  // Dados Pessoais
  nome: string;
  email: string;
  telefone: string;
  cpf: string;
  foto?: string;
  
  // Dados Profissionais
  tipo: ProfessionalType;
  status: ProfessionalStatus;
  data_admissao: string; // ISO8601
  data_demissao?: string;
  especialidades?: string[];
  observacoes?: string;
  
  // Comissão
  tipo_comissao?: CommissionType;
  comissao?: number;
  comissao_produtos?: number;
  
  // Horário de Trabalho
  horario_trabalho?: WorkSchedule;
  
  // Timestamps
  criado_em: string;
  atualizado_em: string;
}

// ============================================================================
// DTOs - Request
// ============================================================================

/** DTO para criar profissional */
export interface CreateProfessionalRequest {
  nome: string;
  email: string;
  telefone: string;
  cpf: string;
  tipo: ProfessionalType;
  data_admissao?: string;
  foto?: string;
  especialidades?: string[];
  observacoes?: string;
  
  // Campos condicionais (Barbeiro/Gerente)
  tambem_barbeiro?: boolean;
  tipo_comissao?: CommissionType;
  comissao?: number;
  comissao_produtos?: number;
  horario_trabalho?: WorkSchedule;
}

/** DTO para atualizar profissional */
export interface UpdateProfessionalRequest {
  nome?: string;
  email?: string;
  telefone?: string;
  foto?: string;
  especialidades?: string[];
  observacoes?: string;
  tipo_comissao?: CommissionType;
  comissao?: number;
  comissao_produtos?: number;
  horario_trabalho?: WorkSchedule;
  status?: ProfessionalStatus;
}

/** DTO para atualizar status do profissional */
export interface UpdateProfessionalStatusRequest {
  status: ProfessionalStatus;
}

/** Filtros para listagem */
export interface ListProfessionalsFilters {
  page?: number;
  page_size?: number;
  tipo?: ProfessionalType;
  status?: ProfessionalStatus;
  search?: string;
  order_by?: 'nome' | 'criado_em' | 'data_admissao';
  order_direction?: 'asc' | 'desc';
}

// ============================================================================
// DTOs - Response
// ============================================================================

/** Resposta de um profissional */
export interface ProfessionalResponse {
  id: string;
  nome: string;
  email: string;
  telefone: string;
  cpf: string;
  tipo: ProfessionalType;
  status: ProfessionalStatus;
  foto?: string;
  especialidades?: string[];
  tipo_comissao?: CommissionType;
  comissao?: number | string; // Backend pode retornar como string (numeric do PostgreSQL)
  comissao_produtos?: number | string;
  horario_trabalho?: WorkSchedule;
  data_admissao: string;
  criado_em: string;
  atualizado_em: string;
}

/** Resposta paginada de listagem */
export interface ListProfessionalsResponse {
  data: ProfessionalResponse[];
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

// ============================================================================
// UI State Types
// ============================================================================

/** Estado do modal de profissional */
export interface ProfessionalModalState {
  isOpen: boolean;
  mode: 'create' | 'edit' | 'view';
  professional?: ProfessionalResponse;
}

// ============================================================================
// HELPERS
// ============================================================================

/** Labels para tipos de profissional */
export const PROFESSIONAL_TYPE_LABELS: Record<ProfessionalType, string> = {
  BARBEIRO: 'Barbeiro',
  GERENTE: 'Gerente',
  RECEPCIONISTA: 'Recepcionista',
  OUTRO: 'Outro',
};

/** Labels para status */
export const PROFESSIONAL_STATUS_LABELS: Record<ProfessionalStatus, string> = {
  ATIVO: 'Ativo',
  INATIVO: 'Inativo',
  AFASTADO: 'Afastado',
  DEMITIDO: 'Demitido',
};

/** Cores para tipos (Tailwind classes) */
export const PROFESSIONAL_TYPE_COLORS: Record<ProfessionalType, string> = {
  BARBEIRO: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  GERENTE: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  RECEPCIONISTA: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
  OUTRO: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
};

/** Cores para status (Tailwind classes) */
export const PROFESSIONAL_STATUS_COLORS: Record<ProfessionalStatus, string> = {
  ATIVO: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  INATIVO: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
  AFASTADO: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300',
  DEMITIDO: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
};

/** Dias da semana para iteração */
export const WEEK_DAYS = [
  { key: 'segunda', label: 'Segunda-feira' },
  { key: 'terca', label: 'Terça-feira' },
  { key: 'quarta', label: 'Quarta-feira' },
  { key: 'quinta', label: 'Quinta-feira' },
  { key: 'sexta', label: 'Sexta-feira' },
  { key: 'sabado', label: 'Sábado' },
  { key: 'domingo', label: 'Domingo' },
] as const;

/** Horário de trabalho padrão */
export const DEFAULT_WORK_SCHEDULE: WorkSchedule = {
  segunda: { ativo: true, turnos: [{ inicio: '08:00', fim: '18:00' }] },
  terca: { ativo: true, turnos: [{ inicio: '08:00', fim: '18:00' }] },
  quarta: { ativo: true, turnos: [{ inicio: '08:00', fim: '18:00' }] },
  quinta: { ativo: true, turnos: [{ inicio: '08:00', fim: '18:00' }] },
  sexta: { ativo: true, turnos: [{ inicio: '08:00', fim: '18:00' }] },
  sabado: { ativo: true, turnos: [{ inicio: '08:00', fim: '14:00' }] },
  domingo: { ativo: false, turnos: [] },
};
