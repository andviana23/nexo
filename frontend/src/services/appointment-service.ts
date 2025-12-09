/**
 * NEXO - Sistema de Gestão para Barbearias
 * Appointment Service
 *
 * Serviço de agendamentos - comunicação com API de appointments do backend.
 */

import { api, getErrorMessage, isAxiosError } from '@/lib/axios';
import { APPOINTMENT_STATUS_COLORS } from '@/lib/fullcalendar-config';
import type {
    AppointmentResponse,
    CalendarEvent,
    CalendarResource,
    CreateAppointmentRequest,
    ListAppointmentsFilters,
    ListAppointmentsResponse,
    Professional,
    RescheduleAppointmentRequest,
    UpdateAppointmentStatusRequest
} from '@/types/appointment';

// =============================================================================
// ENDPOINTS
// =============================================================================

const APPOINTMENT_ENDPOINTS = {
  list: '/appointments',
  create: '/appointments',
  getById: (id: string) => `/appointments/${id}`,
  updateStatus: (id: string) => `/appointments/${id}/status`,
  reschedule: (id: string) => `/appointments/${id}/reschedule`,
  cancel: (id: string) => `/appointments/${id}/cancel`,
  // Novos endpoints de workflow do agendamento
  confirm: (id: string) => `/appointments/${id}/confirm`,
  checkIn: (id: string) => `/appointments/${id}/check-in`,
  startService: (id: string) => `/appointments/${id}/start`,
  finishService: (id: string) => `/appointments/${id}/finish`,
  complete: (id: string) => `/appointments/${id}/complete`,
  noShow: (id: string) => `/appointments/${id}/no-show`,
  // Recursos relacionados
  professionals: '/professionals',
  customers: '/customers',
  services: '/servicos',
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
      // Para erros de permissão/escopo ou conflitos conhecidos, propaga como erro tipado
      if (isAxiosError(error) && error.response && [403, 404, 409].includes(error.response.status)) {
        mapAppointmentError(error);
      }
      // Retornar resposta vazia em caso de erro inesperado para não travar o calendário
      return { data: [], page: 1, page_size: 20, total: 0 };
    }
  },

  /**
   * Busca um agendamento pelo ID
   */
  async getById(id: string): Promise<AppointmentResponse> {
    const response = await api.get<AppointmentResponse>(
      APPOINTMENT_ENDPOINTS.getById(id)
    );
    return response.data;
  },

  /**
   * Cria um novo agendamento
   */
  async create(data: CreateAppointmentRequest): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.create,
        data
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Atualiza um agendamento existente (reagendamento)
   * Backend usa PATCH /appointments/:id/reschedule
   */
  async reschedule(id: string, data: RescheduleAppointmentRequest): Promise<AppointmentResponse> {
    try {
      const response = await api.patch<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.reschedule(id),
        data
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Cancela um agendamento
   * Backend usa POST /appointments/:id/cancel
   */
  async cancel(id: string, reason?: string): Promise<void> {
    try {
      await api.post(APPOINTMENT_ENDPOINTS.cancel(id), {
        reason: reason,
      });
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Atualiza o status de um agendamento
   * Backend usa PATCH /appointments/:id/status
   */
  async updateStatus(
    id: string,
    data: UpdateAppointmentStatusRequest
  ): Promise<AppointmentResponse> {
    try {
      const response = await api.patch<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.updateStatus(id),
        data
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  // ===========================================================================
  // WORKFLOW DO AGENDAMENTO - Novos endpoints de transição de status
  // ===========================================================================

  /**
   * Confirma um agendamento
   * Transição: CREATED → CONFIRMED
   */
  async confirm(id: string): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.confirm(id)
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Marca cliente como chegou (check-in)
   * Transição: CONFIRMED → CHECKED_IN
   */
  async checkIn(id: string): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.checkIn(id)
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Inicia o atendimento
   * Transição: CONFIRMED/CHECKED_IN → IN_SERVICE
   */
  async startService(id: string): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.startService(id)
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Finaliza o atendimento (aguardando pagamento)
   * Transição: IN_SERVICE → AWAITING_PAYMENT
   */
  async finishService(id: string): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.finishService(id)
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Conclui o agendamento (pagamento recebido)
   * Transição: IN_SERVICE/AWAITING_PAYMENT → DONE
   */
  async complete(id: string): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.complete(id)
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
  },

  /**
   * Marca cliente como não compareceu
   * Transição: CONFIRMED/CHECKED_IN → NO_SHOW
   */
  async noShow(id: string): Promise<AppointmentResponse> {
    try {
      const response = await api.post<AppointmentResponse>(
        APPOINTMENT_ENDPOINTS.noShow(id)
      );
      return response.data;
    } catch (error) {
      mapAppointmentError(error);
    }
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
  // Guard para array undefined/null
  if (!appointments || !Array.isArray(appointments)) {
    console.warn('[appointmentsToCalendarEvents] appointments is undefined or not an array');
    return [];
  }

  return appointments.map((appointment) => {
    const statusColors = APPOINTMENT_STATUS_COLORS[appointment.status];
    const serviceNames = (appointment.services || []).map((s) => s.service_name).join(', ');
    
    return {
      id: appointment.id,
      resourceId: appointment.professional_id,
      title: `${appointment.customer_name} - ${serviceNames}`,
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
  // Guard para array undefined/null
  if (!professionals || !Array.isArray(professionals)) {
    console.warn('[professionalsToCalendarResources] professionals is undefined or not an array');
    return [];
  }

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

export class BlockedTimeError extends AppointmentError {
  constructor() {
    super('Horário bloqueado para o profissional.', 'BLOCKED_TIME');
    this.name = 'BlockedTimeError';
  }
}

export class ProfessionalNotFoundError extends AppointmentError {
  constructor() {
    super('Profissional não encontrado.', 'PROFESSIONAL_NOT_FOUND');
    this.name = 'ProfessionalNotFoundError';
  }
}

export class CustomerNotFoundError extends AppointmentError {
  constructor() {
    super('Cliente não encontrado. Cadastre o cliente primeiro.', 'CUSTOMER_NOT_FOUND');
    this.name = 'CustomerNotFoundError';
  }
}

export class ServiceNotFoundError extends AppointmentError {
  constructor() {
    super('Serviço não encontrado.', 'SERVICE_NOT_FOUND');
    this.name = 'ServiceNotFoundError';
  }
}

export class AppointmentNotFoundError extends AppointmentError {
  constructor() {
    super('Agendamento não encontrado.', 'APPOINTMENT_NOT_FOUND');
    this.name = 'AppointmentNotFoundError';
  }
}

export class ForbiddenScopeError extends AppointmentError {
  constructor() {
    super('Acesso negado: barbeiro só pode agir nos próprios agendamentos.', 'FORBIDDEN_SCOPE');
    this.name = 'ForbiddenScopeError';
  }
}

export class InsufficientIntervalError extends AppointmentError {
  constructor() {
    super('Intervalo mínimo de 10 minutos entre agendamentos.', 'INSUFFICIENT_INTERVAL');
    this.name = 'InsufficientIntervalError';
  }
}

export class InvalidTransitionError extends AppointmentError {
  constructor() {
    super('Transição de status ou fluxo não permitida.', 'INVALID_TRANSITION');
    this.name = 'InvalidTransitionError';
  }
}

// Mantido para compatibilidade com código legado
export class ProfessionalInactiveError extends AppointmentError {
  constructor() {
    super('Este barbeiro não está disponível no momento.', 'PROFESSIONAL_INACTIVE');
    this.name = 'ProfessionalInactiveError';
  }
}

/**
 * Mapeia respostas de erro da API para erros tipados de agendamento.
 */
function mapAppointmentError(error: unknown): never {
  if (!isAxiosError(error) || !error.response) {
    throw error instanceof Error ? error : new AppointmentError('Erro ao processar agendamento.');
  }

  const status = error.response.status;
  const data = error.response.data as { error?: string; message?: string } | undefined;
  const code = (data?.error || '').toUpperCase();
  const message = (data?.message || '').toLowerCase();

  if (status === 403) {
    throw new ForbiddenScopeError();
  }

  if (status === 404) {
    if (code.includes('PROFESSIONAL') || message.includes('profission')) {
      throw new ProfessionalNotFoundError();
    }
    if (code.includes('SERVICE') || message.includes('servi')) {
      throw new ServiceNotFoundError();
    }
    if (code.includes('CUSTOMER') || message.includes('cliente')) {
      throw new CustomerNotFoundError();
    }
    throw new AppointmentNotFoundError();
  }

  if (status === 409) {
    if (code === 'BLOCKED_TIME' || message.includes('bloque')) {
      throw new BlockedTimeError();
    }
    if (code === 'INSUFFICIENT_INTERVAL' || message.includes('interval')) {
      throw new InsufficientIntervalError();
    }
    if (code === 'INVALID_TRANSITION' || message.includes('transi')) {
      throw new InvalidTransitionError();
    }
    throw new TimeSlotConflictError();
  }

  throw new AppointmentError(getErrorMessage(error), code || 'APPOINTMENT_ERROR');
}
