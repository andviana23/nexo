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
  | 'IN_SERVICE'
  | 'DONE'
  | 'NO_SHOW'
  | 'CANCELED';

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

/** DTO para atualizar um agendamento */
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
  date_from?: string; // ISO8601
  date_to?: string;   // ISO8601
  status?: AppointmentStatus;
}

/** Filtros para verificar disponibilidade */
export interface CheckAvailabilityParams {
  professional_id: string;
  date: string; // YYYY-MM-DD
}

// ============================================================================
// DTOs - Response
// ============================================================================

/** Resposta de um agendamento com dados expandidos */
export interface AppointmentResponse {
  id: string;
  tenant_id: string;
  professional: {
    id: string;
    name: string;
    avatar_url?: string;
  };
  customer: {
    id: string;
    name: string;
    phone?: string;
  };
  services: {
    id: string;
    name: string;
    price: number;
    duration: number;
  }[];
  start_time: string;
  end_time: string;
  status: AppointmentStatus;
  total_price: number;
  notes?: string;
  created_at: string;
  updated_at: string;
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
