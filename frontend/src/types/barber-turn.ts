/**
 * Types da Lista da Vez (Barber Turn)
 *
 * Sistema de fila giratória para distribuição justa de clientes
 * entre barbeiros. Menor pontuação = próximo da fila.
 */

// =============================================================================
// Request DTOs
// =============================================================================

/**
 * Request para adicionar barbeiro à fila
 */
export interface AddBarberToTurnListRequest {
  professional_id: string;
}

/**
 * Request para registrar atendimento
 */
export interface RecordTurnRequest {
  professional_id: string;
}

/**
 * Request para reset mensal
 */
export interface ResetTurnListRequest {
  save_history?: boolean;
}

/**
 * Request para listar barbeiros com filtro
 */
export interface ListBarbersTurnRequest {
  is_active?: boolean;
}

/**
 * Request para buscar histórico
 */
export interface GetTurnHistoryRequest {
  month_year?: string; // formato: "YYYY-MM"
}

// =============================================================================
// Response DTOs
// =============================================================================

/**
 * Resposta de um barbeiro na fila
 */
export interface BarberTurnResponse {
  id: string;
  tenant_id: string;
  professional_id: string;
  professional_name: string;
  professional_type: string;
  professional_photo?: string;
  current_points: number;
  last_turn_at?: string;
  is_active: boolean;
  position: number;
  created_at: string;
  updated_at: string;
}

/**
 * Próximo barbeiro da fila
 */
export interface NextBarberResponse {
  professional_id: string;
  professional_name: string;
  professional_photo?: string;
  current_points: number;
}

/**
 * Estatísticas da fila
 */
export interface BarberTurnStatsResponse {
  total_ativos: number;
  total_pausados: number;
  total_geral: number;
  total_pontos_mes: number;
}

/**
 * Resposta da listagem de barbeiros na fila
 */
export interface ListBarbersTurnResponse {
  barbers: BarberTurnResponse[];
  total: number;
  next_barber?: NextBarberResponse;
  stats: BarberTurnStatsResponse;
}

/**
 * Resposta do registro de atendimento
 */
export interface RecordTurnResponse {
  professional_id: string;
  professional_name: string;
  previous_points: number;
  new_points: number;
  last_turn_at: string;
  message: string;
}

/**
 * Resposta de alteração de status
 */
export interface ToggleStatusResponse {
  professional_id: string;
  professional_name: string;
  is_active: boolean;
  message: string;
}

/**
 * Resposta de remoção de barbeiro
 */
export interface RemoveBarberResponse {
  professional_id: string;
  message: string;
}

/**
 * Snapshot do reset mensal
 */
export interface TurnResetSnapshot {
  month_year: string;
  total_barbers: number;
  total_points_reset: number;
  history_records_created: number;
}

/**
 * Resposta do reset mensal
 */
export interface ResetTurnListResponse {
  message: string;
  snapshot?: TurnResetSnapshot;
}

/**
 * Histórico de atendimentos de um barbeiro
 */
export interface TurnHistoryResponse {
  id: string;
  tenant_id: string;
  professional_id: string;
  professional_name: string;
  month_year: string;
  total_turns: number;
  final_points: number;
  created_at: string;
}

/**
 * Resposta da listagem de histórico
 */
export interface ListTurnHistoryResponse {
  history: TurnHistoryResponse[];
  total: number;
}

/**
 * Resumo mensal do histórico
 */
export interface TurnHistorySummaryResponse {
  month_year: string;
  total_barbeiros: number;
  total_atendimentos: number;
  media_atendimentos: number;
}

/**
 * Resposta do resumo de histórico
 */
export interface ListHistorySummaryResponse {
  summary: TurnHistorySummaryResponse[];
}

/**
 * Barbeiro disponível para adicionar
 */
export interface AvailableBarberResponse {
  id: string;
  nome: string;
  foto?: string;
  status: string;
}

/**
 * Resposta de barbeiros disponíveis
 */
export interface ListAvailableBarbersResponse {
  barbers: AvailableBarberResponse[];
  total: number;
}

// =============================================================================
// UI/Component Types
// =============================================================================

/**
 * Estado do barbeiro na UI
 */
export type BarberTurnStatus = 'active' | 'paused';

/**
 * Ação disponível para um barbeiro
 */
export type BarberTurnAction =
  | 'record' // Registrar atendimento
  | 'pause' // Pausar
  | 'activate' // Ativar
  | 'remove'; // Remover

/**
 * Props do card de barbeiro
 */
export interface BarberTurnCardProps {
  barber: BarberTurnResponse;
  isNext?: boolean;
  onRecordTurn: (professionalId: string) => void;
  onToggleStatus: (professionalId: string) => void;
  onRemove: (professionalId: string) => void;
}

/**
 * Props do modal de adicionar barbeiro
 */
export interface AddBarberModalProps {
  open: boolean;
  onClose: () => void;
  onAdd: (professionalId: string) => void;
  availableBarbers: AvailableBarberResponse[];
  isLoading?: boolean;
}

/**
 * Props do modal de reset mensal
 */
export interface ResetModalProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (saveHistory: boolean) => void;
  stats: BarberTurnStatsResponse;
  isLoading?: boolean;
}

/**
 * Props do componente de histórico
 */
export interface TurnHistoryProps {
  history: TurnHistoryResponse[];
  summary: TurnHistorySummaryResponse[];
  onMonthChange?: (monthYear: string) => void;
}
