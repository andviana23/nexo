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
  BlockedTimeError,
  CustomerNotFoundError,
  professionalsToCalendarResources,
  TimeSlotConflictError,
  InsufficientIntervalError,
  InvalidTransitionError,
  ForbiddenScopeError,
  ProfessionalNotFoundError,
  ServiceNotFoundError,
  AppointmentNotFoundError
} from '@/services/appointment-service';
import type {
  AppointmentResponse,
  AppointmentStatus,
  CreateAppointmentRequest,
  ListAppointmentsFilters,
  ListAppointmentsResponse,
  RescheduleAppointmentRequest,
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
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: [newAppointment, ...old.data],
            total: (old.total || 0) + 1,
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
 * Hook para reagendar um agendamento (alterar data/horário)
 * Usa optimistic update para resposta imediata na UI
 */
export function useRescheduleAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: RescheduleAppointmentRequest }) =>
      appointmentService.reschedule(id, data),

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
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? {
                  ...apt,
                  start_time: data.new_start_time ?? apt.start_time,
                  professional_id: data.professional_id ?? apt.professional_id,
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
            start_time: data.new_start_time ?? previousDetail.start_time,
            professional_id: data.professional_id ?? previousDetail.professional_id,
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
      toast.success('Agendamento reagendado com sucesso!');
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
          if (!old || !Array.isArray(old.data)) return old;
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

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
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
 * 
 * IMPORTANTE: Verifica se o status atual já é o desejado para evitar
 * chamadas duplicadas que causam erro "transição de status inválida"
 */
export function useUpdateAppointmentStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      data,
      currentStatus,
    }: {
      id: string;
      data: UpdateAppointmentStatusRequest;
      currentStatus?: AppointmentStatus;
    }) => {
      // Se o status atual já é o desejado, não fazer a chamada
      if (currentStatus && currentStatus === data.status) {
        console.warn(`[useUpdateAppointmentStatus] Status já é ${data.status}, ignorando chamada duplicada`);
        // Retorna dados "falsos" para satisfazer o tipo
        return { id, status: data.status } as AppointmentResponse;
      }
      return appointmentService.updateStatus(id, data);
    },

    // Optimistic Update
    onMutate: async ({ id, data, currentStatus }) => {
      // Se status já é o desejado, não fazer optimistic update
      if (currentStatus && currentStatus === data.status) {
        return { skipped: true };
      }
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
          // Se não há dados no cache, retornar como está
          if (!old) return old;

          // Verifica se old é um objeto válido
          if (typeof old !== 'object') return old;

          // Suporta tanto 'data' quanto 'appointments' como propriedade
          const rawAppointments = (old as { data?: unknown }).data ??
            (old as { appointments?: unknown }).appointments;

          // Se não há array de appointments, retornar sem modificar
          if (!Array.isArray(rawAppointments)) return old;

          // Filtra e mapeia com segurança total
          const updatedAppointments = rawAppointments
            .filter((apt): apt is AppointmentResponse => {
              return apt != null &&
                typeof apt === 'object' &&
                'id' in apt &&
                'status' in apt;
            })
            .map((apt) =>
              apt.id === id ? { ...apt, status: data.status } : apt
            );

          return {
            ...old,
            data: updatedAppointments,
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

    onError: (error, __, context) => {
      // Se foi skipped, não fazer rollback
      if (context?.skipped) return;

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
      handleAppointmentError(error as Error);
    },

    onSettled: (_, __, variables, context) => {
      // Se foi skipped, não invalidar queries
      if (context?.skipped) return;

      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: appointmentKeys.detail(variables.id),
      });
    },

    onSuccess: (_, __, context) => {
      // Se foi skipped, não mostrar toast (já está no status desejado)
      if (context?.skipped) return;

      toast.success('Status atualizado!');
    },
  });
}

// =============================================================================
// HOOKS DE WORKFLOW - Transições de Status Específicas
// =============================================================================

/**
 * Hook para confirmar agendamento
 * Transição: CREATED → CONFIRMED
 */
export function useConfirmAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => appointmentService.confirm(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'CONFIRMED' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
    },

    onSuccess: () => {
      toast.success('Agendamento confirmado!');
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },
  });
}

/**
 * Hook para marcar cliente como chegou (check-in)
 * Transição: CONFIRMED → CHECKED_IN
 */
export function useCheckInAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => appointmentService.checkIn(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'CHECKED_IN' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },

    onSuccess: () => {
      toast.success('Cliente chegou! ✅');
    },
  });
}

/**
 * Hook para iniciar atendimento
 * Transição: CONFIRMED/CHECKED_IN → IN_SERVICE
 */
export function useStartServiceAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => appointmentService.startService(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'IN_SERVICE' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },

    onSuccess: () => {
      toast.success('Atendimento iniciado! ✂️');
    },
  });
}

/**
 * Hook para finalizar atendimento (aguardando pagamento)
 * Transição: IN_SERVICE → AWAITING_PAYMENT
 */
export function useFinishServiceAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => appointmentService.finishService(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'AWAITING_PAYMENT' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },

    onSuccess: () => {
      toast.success('Atendimento finalizado! Aguardando pagamento 💰');
    },
  });
}

/**
 * Hook para concluir agendamento (pagamento recebido)
 * Transição: IN_SERVICE/AWAITING_PAYMENT → DONE
 */
export function useCompleteAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => appointmentService.complete(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'DONE' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },

    onSuccess: () => {
      toast.success('Agendamento concluído! ✅💰');
    },
  });
}

/**
 * Hook para marcar cliente como não compareceu
 * Transição: CONFIRMED/CHECKED_IN → NO_SHOW
 */
export function useNoShowAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => appointmentService.noShow(id),

    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: appointmentKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListAppointmentsResponse>({
        queryKey: appointmentKeys.lists(),
      });

      queryClient.setQueriesData<ListAppointmentsResponse>(
        { queryKey: appointmentKeys.lists() },
        (old) => {
          if (!old || !Array.isArray(old.data)) return old;
          return {
            ...old,
            data: old.data.map((apt) =>
              apt.id === id
                ? { ...apt, status: 'NO_SHOW' as AppointmentStatus }
                : apt
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, __, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleAppointmentError(error as Error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },

    onSuccess: () => {
      toast.warning('Cliente marcado como não compareceu.');
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

  if (error instanceof BlockedTimeError) {
    toast.error('Horário bloqueado para o profissional. Escolha outro horário.');
    return;
  }

  if (error instanceof InsufficientIntervalError) {
    toast.error('Intervalo mínimo de 10 minutos não respeitado.');
    return;
  }

  if (error instanceof InvalidTransitionError) {
    toast.error('Transição de status não permitida para este agendamento.');
    return;
  }

  if (error instanceof ForbiddenScopeError) {
    toast.error('Acesso negado: barbeiro só pode agir nos próprios agendamentos.');
    return;
  }

  if (error instanceof ProfessionalNotFoundError) {
    toast.error('Profissional não encontrado. Selecione outro.');
    return;
  }

  if (error instanceof ServiceNotFoundError) {
    toast.error('Serviço não encontrado. Verifique a lista de serviços.');
    return;
  }

  if (error instanceof CustomerNotFoundError) {
    toast.error('Cliente não encontrado. Cadastre o cliente primeiro.');
    return;
  }

  if (error instanceof AppointmentNotFoundError) {
    toast.error('Agendamento não encontrado. Atualize a página e tente novamente.');
    return;
  }

  // Erro genérico
  toast.error('Erro ao processar agendamento. Tente novamente.');
}

// =============================================================================
// ALIASES (para manter compatibilidade com código existente)
// =============================================================================

/**
 * @deprecated Use useRescheduleAppointment em vez disso
 * Mantido para compatibilidade com código legado
 */
export const useUpdateAppointment = useRescheduleAppointment;
