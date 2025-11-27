/**
 * NEXO - Sistema de Gestão para Barbearias
 * Appointment Service
 *
 * Serviço de agendamentos - comunicação com API de appointments do backend.
 */

import { api } from '@/lib/axios';
import { APPOINTMENT_STATUS_COLORS } from '@/lib/fullcalendar-config';
import type {
  AppointmentResponse,
  AvailabilityResponse,
  CalendarEvent,
  CalendarResource,
  CheckAvailabilityParams,
  CreateAppointmentRequest,
  ListAppointmentsFilters,
  ListAppointmentsResponse,
  Professional,
  UpdateAppointmentRequest,
  UpdateAppointmentStatusRequest,
} from '@/types/appointment';

// =============================================================================
// ENDPOINTS
// =============================================================================

const APPOINTMENT_ENDPOINTS = {
  list: '/appointments',
  create: '/appointments',
  getById: (id: string) => `/appointments/${id}`,
  update: (id: string) => `/appointments/${id}`,
  delete: (id: string) => `/appointments/${id}`,
  updateStatus: (id: string) => `/appointments/${id}/status`,
  availability: '/appointments/availability',
  professionals: '/professionals',
  customers: '/customers',
  services: '/services',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const appointmentService = {
  /**
   * Lista agendamentos com filtros
   */
  async list(filters: ListAppointmentsFilters = {}): Promise<ListAppointmentsResponse> {
    console.log('[appointment-service] Chamando list com filtros:', filters);
    try {
      const response = await api.get<ListAppointmentsResponse>(
        APPOINTMENT_ENDPOINTS.list,
        { params: filters }
      );
      console.log('[appointment-service] Resposta list:', response.data);
      return response.data;
    } catch (error) {
      console.error('[appointment-service] Erro ao listar agendamentos:', error);
      // Retornar resposta vazia em caso de erro para não travar o calendário
      return { data: [], page: 1, page_size: 20, total: 0 };
    }
  },

  /**
   * Busca um agendamento pelo ID
   */
  async getById(id: string): Promise<AppointmentResponse> {
    const response = await api.get<{ data: AppointmentResponse }>(
      APPOINTMENT_ENDPOINTS.getById(id)
    );
    return response.data.data;
  },

  /**
   * Cria um novo agendamento
   */
  async create(data: CreateAppointmentRequest): Promise<AppointmentResponse> {
    const response = await api.post<{ data: AppointmentResponse }>(
      APPOINTMENT_ENDPOINTS.create,
      data
    );
    return response.data.data;
  },

  /**
   * Atualiza um agendamento existente
   */
  async update(id: string, data: UpdateAppointmentRequest): Promise<AppointmentResponse> {
    const response = await api.put<{ data: AppointmentResponse }>(
      APPOINTMENT_ENDPOINTS.update(id),
      data
    );
    return response.data.data;
  },

  /**
   * Cancela um agendamento
   */
  async cancel(id: string, reason?: string): Promise<void> {
    await api.delete(APPOINTMENT_ENDPOINTS.delete(id), {
      data: reason ? { reason } : undefined,
    });
  },

  /**
   * Atualiza o status de um agendamento
   */
  async updateStatus(
    id: string,
    data: UpdateAppointmentStatusRequest
  ): Promise<AppointmentResponse> {
    const response = await api.put<{ data: AppointmentResponse }>(
      APPOINTMENT_ENDPOINTS.updateStatus(id),
      data
    );
    return response.data.data;
  },

  /**
   * Verifica disponibilidade de um profissional em uma data
   */
  async checkAvailability(params: CheckAvailabilityParams): Promise<AvailabilityResponse> {
    const response = await api.get<AvailabilityResponse>(
      APPOINTMENT_ENDPOINTS.availability,
      { params }
    );
    return response.data;
  },

  /**
   * Lista profissionais (barbeiros) ativos
   */
  async listProfessionals(): Promise<Professional[]> {
    try {
      const response = await api.get<{ data: Array<{
        id: string;
        tenant_id: string;
        nome: string;
        email: string;
        telefone: string;
        foto?: string;
        status: string;
        tipo: string;
      }>; total: number }>(
        APPOINTMENT_ENDPOINTS.professionals,
        { params: { status: 'ATIVO' } }
      );
      
      // Mapear resposta do backend para o tipo Professional usado no frontend
      return response.data.data.map(p => ({
        id: p.id,
        tenant_id: p.tenant_id,
        name: p.nome,
        email: p.email,
        phone: p.telefone,
        avatar_url: p.foto,
        is_active: p.status === 'ATIVO',
        google_calendar_connected: false,
        created_at: new Date().toISOString(),
      }));
    } catch (error) {
      console.error('[appointment-service] Erro ao buscar profissionais:', error);
      return [];
    }
  },
};

// =============================================================================
// HELPERS - Conversão para FullCalendar
// =============================================================================

/**
 * Converte agendamentos da API para eventos do FullCalendar
 */
export function appointmentsToCalendarEvents(
  appointments: AppointmentResponse[]
): CalendarEvent[] {
  return appointments.map((appointment) => {
    const statusColors = APPOINTMENT_STATUS_COLORS[appointment.status];
    const serviceNames = appointment.services.map((s) => s.name).join(', ');
    
    return {
      id: appointment.id,
      resourceId: appointment.professional.id,
      title: `${appointment.customer.name} - ${serviceNames}`,
      start: appointment.start_time,
      end: appointment.end_time,
      backgroundColor: statusColors?.backgroundColor,
      borderColor: statusColors?.borderColor,
      textColor: statusColors?.textColor,
      extendedProps: {
        appointment,
      },
    };
  });
}

/**
 * Converte profissionais da API para recursos do FullCalendar
 */
export function professionalsToCalendarResources(
  professionals: Professional[]
): CalendarResource[] {
  return professionals.map((professional) => ({
    id: professional.id,
    title: professional.name,
    extendedProps: {
      professional,
    },
  }));
}

// =============================================================================
// TIPOS DE ERRO ESPECÍFICOS
// =============================================================================

export class AppointmentError extends Error {
  constructor(
    message: string,
    public code: string = 'APPOINTMENT_ERROR'
  ) {
    super(message);
    this.name = 'AppointmentError';
  }
}

export class TimeSlotConflictError extends AppointmentError {
  constructor() {
    super('Este horário já está ocupado. Escolha outro.', 'TIME_SLOT_CONFLICT');
    this.name = 'TimeSlotConflictError';
  }
}

export class ProfessionalInactiveError extends AppointmentError {
  constructor() {
    super('Este barbeiro não está disponível no momento.', 'PROFESSIONAL_INACTIVE');
    this.name = 'ProfessionalInactiveError';
  }
}

export class CustomerNotFoundError extends AppointmentError {
  constructor() {
    super('Cliente não encontrado. Cadastre o cliente primeiro.', 'CUSTOMER_NOT_FOUND');
    this.name = 'CustomerNotFoundError';
  }
}

export class InsufficientIntervalError extends AppointmentError {
  constructor() {
    super('Intervalo mínimo de 10 minutos entre agendamentos.', 'INSUFFICIENT_INTERVAL');
    this.name = 'InsufficientIntervalError';
  }
}

export class InvalidStatusTransitionError extends AppointmentError {
  constructor() {
    super('Transição de status inválida.', 'INVALID_STATUS_TRANSITION');
    this.name = 'InvalidStatusTransitionError';
  }
}
