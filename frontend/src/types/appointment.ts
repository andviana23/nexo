/**
 * Tipos TypeScript para o Módulo de Agendamento
 * 
 * @module types/appointment
 * @description Definições de tipos para agendamentos no NEXO
 */

// ============================================================================
// ENUMS
// ============================================================================

/** Status possíveis de um agendamento */
export type AppointmentStatus = 
  | 'CREATED'
  | 'CONFIRMED'
  | 'CHECKED_IN'       // Cliente chegou
  | 'IN_SERVICE'
  | 'AWAITING_PAYMENT' // Aguardando pagamento
  | 'DONE'
  | 'NO_SHOW'
  | 'CANCELED';

/** Configuração de cada status */
export interface StatusConfig {
  label: string;
  color: string;
  bgColor: string;
  textColor: string;
  icon: string;
  allowedTransitions: AppointmentStatus[];
}

/** Mapa de configurações de status */
export const STATUS_CONFIG: Record<AppointmentStatus, StatusConfig> = {
  CREATED: {
    label: 'Criado',
    color: '#3B82F6',
    bgColor: 'bg-blue-100',
    textColor: 'text-blue-700',
    icon: 'calendar-plus',
    allowedTransitions: ['CONFIRMED', 'CHECKED_IN', 'CANCELED', 'NO_SHOW'],
  },
  CONFIRMED: {
    label: 'Confirmado',
    color: '#10B981',
    bgColor: 'bg-emerald-100',
    textColor: 'text-emerald-700',
    icon: 'check-circle',
    allowedTransitions: ['CHECKED_IN', 'IN_SERVICE', 'CANCELED', 'NO_SHOW'],
  },
  CHECKED_IN: {
    label: 'Cliente Chegou',
    color: '#8B5CF6',
    bgColor: 'bg-violet-100',
    textColor: 'text-violet-700',
    icon: 'user-check',
    allowedTransitions: ['IN_SERVICE', 'CANCELED', 'NO_SHOW'],
  },
  IN_SERVICE: {
    label: 'Em Atendimento',
    color: '#F59E0B',
    bgColor: 'bg-amber-100',
    textColor: 'text-amber-700',
    icon: 'scissors',
    allowedTransitions: ['AWAITING_PAYMENT', 'DONE', 'CANCELED'],
  },
  AWAITING_PAYMENT: {
    label: 'Aguardando Pagamento',
    color: '#EC4899',
    bgColor: 'bg-pink-100',
    textColor: 'text-pink-700',
    icon: 'credit-card',
    allowedTransitions: ['DONE', 'CANCELED'],
  },
  DONE: {
    label: 'Concluído',
    color: '#22C55E',
    bgColor: 'bg-green-100',
    textColor: 'text-green-700',
    icon: 'check',
    allowedTransitions: [],
  },
  NO_SHOW: {
    label: 'Não Compareceu',
    color: '#EF4444',
    bgColor: 'bg-red-100',
    textColor: 'text-red-700',
    icon: 'user-x',
    allowedTransitions: [],
  },
  CANCELED: {
    label: 'Cancelado',
    color: '#6B7280',
    bgColor: 'bg-gray-100',
    textColor: 'text-gray-700',
    icon: 'x-circle',
    allowedTransitions: [],
  },
};

/** Verifica se uma transição de status é permitida */
export function canTransitionTo(currentStatus: AppointmentStatus, newStatus: AppointmentStatus): boolean {
  return STATUS_CONFIG[currentStatus].allowedTransitions.includes(newStatus);
}

/** Retorna a configuração de um status */
export function getStatusConfig(status: AppointmentStatus): StatusConfig {
  return STATUS_CONFIG[status] || STATUS_CONFIG.CREATED;
}

/** Verifica se é um status final */
export function isFinalStatus(status: AppointmentStatus): boolean {
  return ['DONE', 'NO_SHOW', 'CANCELED'].includes(status);
}

/** Verifica se é um status ativo */
export function isActiveStatus(status: AppointmentStatus): boolean {
  return ['CREATED', 'CONFIRMED', 'CHECKED_IN', 'IN_SERVICE', 'AWAITING_PAYMENT'].includes(status);
}

// ============================================================================
// ENTITIES
// ============================================================================

/** Entidade principal de Agendamento */
export interface Appointment {
  id: string;
  tenant_id: string;
  professional_id: string;
  customer_id: string;
  start_time: string; // ISO8601
  end_time: string;   // ISO8601
  status: AppointmentStatus;
  notes?: string;
  google_calendar_event_id?: string;
  canceled_reason?: string;
  total_price: number;
  created_at: string;
  updated_at: string;
  
  // Relacionamentos (quando expandidos)
  professional?: Professional;
  customer?: Customer;
  services?: AppointmentService[];
}

/** Serviço vinculado a um agendamento */
export interface AppointmentService {
  appointment_id: string;
  service_id: string;
  price_at_booking: number;
  duration_at_booking: number;
  
  // Relacionamento expandido
  service?: Service;
}

/** Profissional (Barbeiro) */
export interface Professional {
  id: string;
  tenant_id: string;
  name: string;
  email?: string;
  phone?: string;
  avatar_url?: string;
  is_active: boolean;
  google_calendar_connected?: boolean;
  created_at: string;
}

/** Cliente */
export interface Customer {
  id: string;
  tenant_id: string;
  name: string;
  email?: string;
  phone?: string;
  avatar_url?: string;
  is_active: boolean;
  created_at: string;
}

/** Serviço */
export interface Service {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  default_price: number;
  default_duration: number; // em minutos
  is_active: boolean;
  category?: string;
  created_at: string;
}

// ============================================================================
// DTOs - Request
// ============================================================================

/** DTO para criar um novo agendamento */
export interface CreateAppointmentRequest {
  professional_id: string;
  customer_id: string;
  service_ids: string[];
  start_time: string; // ISO8601
  notes?: string;
}

/** DTO para reagendar um agendamento */
export interface RescheduleAppointmentRequest {
  new_start_time: string; // ISO8601
  professional_id?: string;
}

/** DTO para atualizar um agendamento (não usado atualmente) */
export interface UpdateAppointmentRequest {
  professional_id?: string;
  service_ids?: string[];
  start_time?: string;
  notes?: string;
}

/** DTO para alterar status de um agendamento */
export interface UpdateAppointmentStatusRequest {
  status: AppointmentStatus;
  reason?: string;
}

/** Filtros para listagem de agendamentos */
export interface ListAppointmentsFilters {
  page?: number;
  page_size?: number;
  professional_id?: string;
  customer_id?: string;
  start_date?: string; // YYYY-MM-DD ou ISO8601
  end_date?: string;   // YYYY-MM-DD ou ISO8601
  status?: AppointmentStatus | AppointmentStatus[]; // String única ou array de status
}

/** Filtros para verificar disponibilidade */
export interface CheckAvailabilityParams {
  professional_id: string;
  date: string; // YYYY-MM-DD
}

// ============================================================================
// DTOs - Response
// ============================================================================

/** Serviço dentro de um agendamento (resposta da API) */
export interface AppointmentServiceResponse {
  service_id: string;
  service_name: string;
  price: string;
  duration: number;
}

/** Resposta de um agendamento com dados expandidos (formato do backend) */
export interface AppointmentResponse {
  id: string;
  tenant_id: string;
  professional_id: string;
  professional_name: string;
  customer_id: string;
  customer_name: string;
  customer_phone: string;
  services: AppointmentServiceResponse[];
  start_time: string;
  end_time: string;
  duration: number;
  status: AppointmentStatus;
  status_display: string;
  status_color: string;
  total_price: string;
  notes?: string;
  canceled_reason?: string;
  google_calendar_event_id?: string;
  command_id?: string; // ID da comanda vinculada (quando status = AWAITING_PAYMENT)
  created_at: string;
  updated_at: string;
}

/** Helper para acessar dados de professional/customer como objetos */
export function getAppointmentProfessional(apt: AppointmentResponse) {
  return {
    id: apt.professional_id,
    name: apt.professional_name,
  };
}

export function getAppointmentCustomer(apt: AppointmentResponse) {
  return {
    id: apt.customer_id,
    name: apt.customer_name,
    phone: apt.customer_phone,
  };
}

export function getAppointmentServices(apt: AppointmentResponse) {
  return apt.services.map(s => ({
    id: s.service_id,
    name: s.service_name,
    price: parseFloat(s.price),
    duration: s.duration,
  }));
}

/** Resposta paginada de listagem */
export interface ListAppointmentsResponse {
  data: AppointmentResponse[];
  page: number;
  page_size: number;
  total: number;
}

/** Slot de disponibilidade */
export interface AvailabilitySlot {
  time: string; // HH:mm
  available: boolean;
  reason?: 'BOOKED' | 'BLOCKED' | 'OUTSIDE_HOURS';
}

/** Resposta de verificação de disponibilidade */
export interface AvailabilityResponse {
  data: AvailabilitySlot[];
}

// ============================================================================
// FullCalendar Types
// ============================================================================

/** Evento do FullCalendar (adapter para Appointment) */
export interface CalendarEvent {
  id: string;
  resourceId: string; // professional_id
  title: string;
  start: string | Date;
  end: string | Date;
  backgroundColor?: string;
  borderColor?: string;
  textColor?: string;
  extendedProps: {
    appointment: AppointmentResponse;
  };
}

/** Recurso do FullCalendar (Barbeiro) */
export interface CalendarResource {
  id: string;
  title: string;
  eventColor?: string;
  extendedProps?: {
    professional: Professional;
  };
}

// ============================================================================
// Form Types (React Hook Form)
// ============================================================================

/** Schema do formulário de agendamento */
export interface AppointmentFormData {
  professional_id: string;
  customer_id: string;
  service_ids: string[];
  start_date: Date;
  start_time: string; // HH:mm
  notes?: string;
}

// ============================================================================
// UI State Types
// ============================================================================

/** Estado do modal de agendamento */
export interface AppointmentModalState {
  isOpen: boolean;
  mode: 'create' | 'edit' | 'view';
  appointment?: AppointmentResponse;
  initialDate?: Date;
  initialProfessionalId?: string;
}

/** Filtros ativos na visualização */
export interface CalendarFilters {
  professionalIds: string[];
  statuses: AppointmentStatus[];
  view: 'day' | 'week' | 'month' | 'list';
  date: Date;
}

// ============================================================================
// Blocked Times Types
// ============================================================================

/** Horário bloqueado */
export interface BlockedTime {
  id: string;
  tenant_id: string;
  professional_id: string;
  start_time: string; // ISO 8601
  end_time: string;   // ISO 8601
  reason: string;
  is_recurring: boolean;
  recurrence_rule?: string;
  created_at: string;
  updated_at: string;
  created_by?: string;
}

/** Request para criar bloqueio */
export interface CreateBlockedTimeRequest {
  professional_id: string;
  start_time: string; // ISO 8601
  end_time: string;   // ISO 8601
  reason: string;
  is_recurring?: boolean;
  recurrence_rule?: string;
}

/** Response de bloqueio criado */
export interface BlockedTimeResponse {
  id: string;
  tenant_id: string;
  professional_id: string;
  start_time: string;
  end_time: string;
  reason: string;
  is_recurring: boolean;
  recurrence_rule?: string;
  created_at: string;
  updated_at: string;
  created_by?: string;
}

/** Request para listar bloqueios */
export interface ListBlockedTimesRequest {
  professional_id?: string;
  start_date?: string; // ISO 8601
  end_date?: string;   // ISO 8601
}

/** Response de listagem de bloqueios */
export interface ListBlockedTimesResponse {
  blocked_times: BlockedTimeResponse[];
  total: number;
}

// ============================================================================
// HELPERS DE FORMATAÇÃO
// ============================================================================

/**
 * Formata um valor monetário para exibição em BRL.
 * Aceita string numérica do backend (ex: "50.00") ou number.
 * 
 * @param value - Valor numérico como string ou number
 * @returns String formatada em BRL (ex: "R$ 50,00")
 * 
 * @example
 * formatCurrency("50.00") // "R$ 50,00"
 * formatCurrency(50)      // "R$ 50,00"
 * formatCurrency("R$ 50") // "R$ 50,00" (tenta parsear mesmo com prefixo)
 */
export function formatCurrency(value: string | number | undefined | null): string {
  if (value === undefined || value === null) {
    return 'R$ 0,00';
  }
  
  // Se já é number, usa direto
  let numValue: number;
  if (typeof value === 'number') {
    numValue = value;
  } else {
    // Remove "R$", espaços e substitui vírgula por ponto para parse
    const cleaned = value.replace(/[R$\s]/g, '').replace(',', '.');
    numValue = parseFloat(cleaned);
  }
  
  // Trata NaN
  if (isNaN(numValue)) {
    return 'R$ 0,00';
  }
  
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(numValue);
}
