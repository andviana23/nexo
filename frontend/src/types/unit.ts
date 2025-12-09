/**
 * NEXO - Sistema de Gestão para Barbearias
 * Tipos para Multi-Unidade
 *
 * Tipos e interfaces para suporte a múltiplas unidades/filiais.
 * Alinhado com DTOs do backend Go.
 */

// =============================================================================
// UNIT (Unidade/Filial)
// =============================================================================

/**
 * Unidade/Filial do tenant
 */
export interface Unit {
  id: string;
  tenant_id: string;
  nome: string;
  apelido?: string;
  descricao?: string;
  endereco_resumo?: string;
  cidade?: string;
  estado?: string;
  timezone: string;
  ativa: boolean;
  is_matriz: boolean;
  criado_em: string;
  atualizado_em: string;
}

/**
 * Resposta de criação de unidade
 */
export interface CreateUnitRequest {
  nome: string;
  apelido?: string;
  descricao?: string;
  endereco_resumo?: string;
  cidade?: string;
  estado?: string;
  timezone?: string;
  is_matriz?: boolean;
}

/**
 * Resposta de atualização de unidade
 */
export interface UpdateUnitRequest {
  nome?: string;
  apelido?: string;
  descricao?: string;
  endereco_resumo?: string;
  cidade?: string;
  estado?: string;
  timezone?: string;
}

/**
 * Resposta de listagem de unidades
 */
export interface ListUnitsResponse {
  units: Unit[];
  total: number;
}

// =============================================================================
// USER UNIT (Vínculo Usuário-Unidade)
// =============================================================================

/**
 * Vínculo entre usuário e unidade com detalhes
 */
export interface UserUnit {
  id: string;
  user_id: string;
  unit_id: string;
  unit_nome: string;
  unit_apelido?: string;
  unit_matriz: boolean;
  unit_ativa: boolean;
  is_default: boolean;
  role_override?: string;
  tenant_id: string;
}

/**
 * Resposta de listagem de unidades do usuário
 */
export interface ListUserUnitsResponse {
  units: UserUnit[];
  total: number;
}

/**
 * Request para trocar de unidade
 */
export interface SwitchUnitRequest {
  unit_id: string;
}

/**
 * Resposta de troca de unidade
 */
export interface SwitchUnitResponse {
  unit: UserUnit;
  access_token: string; // Novo token com unit_id (se aplicável)
}

/**
 * Request para definir unidade padrão
 */
export interface SetDefaultUnitRequest {
  unit_id: string;
}

/**
 * Request para adicionar usuário à unidade (admin)
 */
export interface AddUserToUnitRequest {
  user_id: string;
  unit_id: string;
  is_default?: boolean;
  role_override?: string;
}

// =============================================================================
// TIPOS DE CONTEXTO E ESTADO
// =============================================================================

/**
 * Estado da unidade ativa no contexto
 */
export interface ActiveUnitState {
  unit: UserUnit | null;
  isLoading: boolean;
  error: Error | null;
}

/**
 * Contexto de unidades disponível via hook
 */
export interface UnitContextValue {
  // Estado
  units: UserUnit[];
  activeUnit: UserUnit | null;
  isLoading: boolean;
  isMultiUnitEnabled: boolean;
  error: Error | null;

  // Ações
  switchUnit: (unitId: string) => Promise<void>;
  setDefaultUnit: (unitId: string) => Promise<void>;
  refreshUnits: () => Promise<void>;
}

// =============================================================================
// CONSTANTES
// =============================================================================

export const UNIT_STORAGE_KEY = 'nexo-active-unit';

/**
 * Header HTTP para identificar a unidade ativa
 */
export const UNIT_HEADER = 'X-Unit-ID';

/**
 * Estados brasileiros para seleção
 */
export const ESTADOS_BRASIL = [
  'AC', 'AL', 'AP', 'AM', 'BA', 'CE', 'DF', 'ES', 'GO',
  'MA', 'MT', 'MS', 'MG', 'PA', 'PB', 'PR', 'PE', 'PI',
  'RJ', 'RN', 'RS', 'RO', 'RR', 'SC', 'SP', 'SE', 'TO'
] as const;

export type EstadoBrasil = typeof ESTADOS_BRASIL[number];

/**
 * Timezones comuns no Brasil
 */
export const TIMEZONES_BRASIL = [
  'America/Sao_Paulo',
  'America/Manaus',
  'America/Cuiaba',
  'America/Recife',
  'America/Fortaleza',
  'America/Belem',
  'America/Rio_Branco',
  'America/Porto_Velho',
  'America/Boa_Vista',
  'America/Noronha',
] as const;

export type TimezoneBrasil = typeof TIMEZONES_BRASIL[number];
