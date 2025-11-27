/**
 * NEXO - Sistema de Gestão para Barbearias
 * React Query Hooks para Agendamentos
 *
 * @module hooks/use-appointments
 * @description Hooks com Optimistic Updates para melhor UX
 */

import {
  appointmentService,
  appointmentsToCalendarEvents,
  CustomerNotFoundError,
  ProfessionalInactiveError,
  professionalsToCalendarResources,
  TimeSlotConflictError,
} from '@/services/appointment-service';
import type {
  AppointmentResponse,
  AppointmentStatus,
  CheckAvailabilityParams,
  CreateAppointmentRequest,
  ListAppointmentsFilters,
  ListAppointmentsResponse,
  UpdateAppointmentRequest,
  UpdateAppointmentStatusRequest
} from '@/types/appointment';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const appointmentKeys = {
  all: ['appointments'] as const,
  lists: () => [...appointmentKeys.all, 'list'] as const,
  list: (filters: ListAppointmentsFilters) =>
    [...appointmentKeys.lists(), filters] as const,
  details: () => [...appointmentKeys.all, 'detail'] as const,
  detail: (id: string) => [...appointmentKeys.details(), id] as const,
  availability: (params: CheckAvailabilityParams) =>
    [...appointmentKeys.all, 'availability', params] as const,
  professionals: () => ['professionals'] as const,
};

// =============================================================================
// HOOKS DE CONSULTA (READ)
// =============================================================================

/**
 * Hook para listar agendamentos com filtros
 */
export function useAppointments(filters: ListAppointmentsFilters = {}) {
  return useQuery({
    queryKey: appointmentKeys.list(filters),
    queryFn: () => appointmentService.list(filters),
    staleTime: 30_000, // 30 segundos
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook para buscar agendamentos e converter para eventos do FullCalendar
 */
export function useCalendarEvents(filters: ListAppointmentsFilters = {}) {
  console.log('[useCalendarEvents] Hook chamado com filtros:', filters);
  return useQuery({
    queryKey: [...appointmentKeys.list(filters), 'calendar'],
    queryFn: async () => {
      console.log('[useCalendarEvents] queryFn executando...');
      const response = await appointmentService.list(filters);
      console.log('[useCalendarEvents] response:', response);
      return appointmentsToCalendarEvents(response.data);
    },
    staleTime: 30_000,
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook para buscar um agendamento específico
 */
export function useAppointment(id: string) {
  return useQuery({
    queryKey: appointmentKeys.detail(id),
    queryFn: () => appointmentService.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook para verificar disponibilidade
 */
export function useAvailability(params: CheckAvailabilityParams) {
  return useQuery({
    queryKey: appointmentKeys.availability(params),
    queryFn: () => appointmentService.checkAvailability(params),
    enabled: !!params.professional_id && !!params.date,
    staleTime: 60_000, // 1 minuto
  });
}

/**
 * Hook para listar profissionais (barbeiros)
 */
export function useProfessionals() {
  return useQuery({
    queryKey: appointmentKeys.professionals(),
    queryFn: () => appointmentService.listProfessionals(),
    staleTime: 5 * 60_000, // 5 minutos
  });
}

/**
 * Hook para buscar profissionais como recursos do FullCalendar
 */
export function useCalendarResources() {
  console.log('[useCalendarResources] Hook chamado');
  return useQuery({
    queryKey: [...appointmentKeys.professionals(), 'calendar'],
    queryFn: async () => {
      console.log('[useCalendarResources] queryFn executando...');
      const professionals = await appointmentService.listProfessionals();
      console.log('[useCalendarResources] professionals:', professionals);
      return professionalsToCalendarResources(professionals);
    },
    staleTime: 5 * 60_000,
  });
}

// =============================================================================
// HOOKS DE MUTAÇÃO (CREATE, UPDATE, DELETE) - COM OPTIMISTIC UPDATES
// =============================================================================

/**
 * Hook para criar um novo agendamento
 * Não usa optimistic update pois o ID é gerado pelo servidor
 */
export function useCreateAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateAppointmentRequest) =>
      appointmentService.create(data),
    onSuccess: (newAppointment) => {
      // Adiciona o novo agendamento ao cache otimisticamente
      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: [newAppointment, ...old.data],
            total: old.total + 1,
          };
        }
      );
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
      toast.success('Agendamento criado com sucesso!');
    },
    onError: (error: Error) => {
      handleAppointmentError(error);
    },
  });
}

/**
 * Hook para atualizar um agendamento
 * Usa optimistic update para resposta imediata na UI
 */
export function useUpdateAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateAppointmentRequest }) =>
      appointmentService.update(id, data),

    // Optimistic Update
    onMutate: async ({ id, data }) => {
      // Cancelar queries em andamento para não sobrescrever nosso update otimista
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });
      await queryClient.cancelQueries({ queryKey: appointmentKeys.detail(id) });

      // Snapshot do estado anterior para rollback
      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });
      const previousDetail = queryClient.getQueryData<AppointmentResponse>(
        appointmentKeys.detail(id)
      );

      // Atualização otimista nas listas
      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? {
                    ...apt,
                    start_time: data.start_time ?? apt.start_time,
                    notes: data.notes ?? apt.notes,
                  }
                : apt
            ),
          };
        }
      );

      // Atualização otimista no detalhe
      if (previousDetail) {
        queryClient.setQueryData<AppointmentResponse>(
          appointmentKeys.detail(id),
          {
            ...previousDetail,
            start_time: data.start_time ?? previousDetail.start_time,
            notes: data.notes ?? previousDetail.notes,
          }
        );
      }

      return { previousLists, previousDetail, id };
    },

    // Rollback em caso de erro
    onError: (error, variables, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      if (context?.previousDetail) {
        queryClient.setQueryData(
          appointmentKeys.detail(context.id),
          context.previousDetail
        );
      }
      handleAppointmentError(error);
    },

    // Sempre refetch para garantir consistência
    onSettled: (_, __, variables) => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: appointmentKeys.detail(variables.id),
      });
    },

    onSuccess: () => {
      toast.success('Agendamento atualizado com sucesso!');
    },
  });
}

/**
 * Hook para cancelar um agendamento
 * Usa optimistic update para feedback imediato
 */
export function useCancelAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, reason }: { id: string; reason?: string }) =>
      appointmentService.cancel(id, reason),

    // Optimistic Update
    onMutate: async ({ id }) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      // Marca como cancelado otimisticamente
      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'CANCELED' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (_, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      toast.error('Erro ao cancelar agendamento. Tente novamente.');
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },

    onSuccess: () => {
      toast.success('Agendamento cancelado.');
    },
  });
}

/**
 * Hook para atualizar o status de um agendamento
 * Usa optimistic update para transições de status suaves
 */
export function useUpdateAppointmentStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: UpdateAppointmentStatusRequest;
    }) => appointmentService.updateStatus(id, data),

    // Optimistic Update
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });
      await queryClient.cancelQueries({ queryKey: appointmentKeys.detail(id) });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });
      const previousDetail = queryClient.getQueryData<AppointmentResponse>(
        appointmentKeys.detail(id)
      );

      // Atualiza status otimisticamente nas listas
      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id ? { ...apt, status: data.status } : apt
            ),
          };
        }
      );

      // Atualiza status no detalhe
      if (previousDetail) {
        queryClient.setQueryData<AppointmentResponse>(
          appointmentKeys.detail(id),
          { ...previousDetail, status: data.status }
        );
      }

      return { previousLists, previousDetail, id };
    },

    onError: (_, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      if (context?.previousDetail) {
        queryClient.setQueryData(
          appointmentKeys.detail(context.id),
          context.previousDetail
        );
      }
      toast.error('Erro ao atualizar status. Tente novamente.');
    },

    onSettled: (_, __, variables) => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: appointmentKeys.detail(variables.id),
      });
    },

    onSuccess: () => {
      toast.success('Status atualizado!');
    },
  });
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Handler centralizado de erros de agendamento
 */
function handleAppointmentError(error: Error) {
  if (error instanceof TimeSlotConflictError) {
    toast.error('Horário já está ocupado. Escolha outro horário.');
    return;
  }

  if (error instanceof ProfessionalInactiveError) {
    toast.error('Barbeiro não está disponível no momento.');
    return;
  }

  if (error instanceof CustomerNotFoundError) {
    toast.error('Cliente não encontrado. Cadastre o cliente primeiro.');
    return;
  }

  // Erro genérico
  toast.error('Erro ao processar agendamento. Tente novamente.');
}
