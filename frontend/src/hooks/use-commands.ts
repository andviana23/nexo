/**
 * React Query hooks para Comandas
 * 
 * Gerencia estado, cache e sincronização com backend para o módulo de comandas.
 */

import { getErrorMessage } from '@/lib/axios';
import {
    addCommandItem,
    addCommandPayment,
    cancelCommand,
    closeCommand,
    createCommand,
    createCommandFromAppointment,
    getCommand,
    getCommandByAppointment,
    listCommands,
    removeCommandItem,
    removeCommandPayment,
    updateCommandItem,
    type ListCommandsFilters,
} from '@/services/command-service';
import type {
    AddCommandItemRequest,
    AddCommandPaymentRequest,
    CloseCommandRequest,
    CreateCommandRequest,
    UpdateCommandItemRequest
} from '@/types/command';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// ============================================================================
// Query Keys
// ============================================================================

export const commandKeys = {
  all: ['commands'] as const,
  lists: () => [...commandKeys.all, 'list'] as const,
  list: (filters: ListCommandsFilters) => [...commandKeys.lists(), filters] as const,
  details: () => [...commandKeys.all, 'detail'] as const,
  detail: (id: string) => [...commandKeys.details(), id] as const,
  byAppointment: (appointmentId: string) => [...commandKeys.all, 'by-appointment', appointmentId] as const,
};

// ============================================================================
// Queries
// ============================================================================

/**
 * Hook para buscar comanda por ID
 */
/**
 * Hook para buscar comanda por ID
 */
export function useCommand(commandId: string | undefined) {
  return useQuery({
    queryKey: commandKeys.detail(commandId!),
    queryFn: () => getCommand(commandId!),
    enabled: !!commandId,
    staleTime: 1000 * 30, // 30 segundos
  });
}

/**
 * Hook para buscar comanda por appointment_id
 */
export function useCommandByAppointment(appointmentId: string | undefined) {
  return useQuery({
    queryKey: commandKeys.byAppointment(appointmentId!),
    queryFn: () => getCommandByAppointment(appointmentId!),
    enabled: !!appointmentId,
    staleTime: 1000 * 30,
  });
}

/**
 * Hook para listar comandas
 */
export function useCommands(filters?: ListCommandsFilters) {
  return useQuery({
    queryKey: commandKeys.list(filters || {}),
    queryFn: () => listCommands(filters),
    staleTime: 1000 * 60, // 1 minuto
  });
}

// ============================================================================
// Mutations - Comandas
// ============================================================================

/**
 * Hook para criar comanda
 */
export function useCreateCommand() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateCommandRequest) => createCommand(data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.lists() });
      if (data.appointment_id) {
        queryClient.invalidateQueries({ 
          queryKey: commandKeys.byAppointment(data.appointment_id) 
        });
      }
      toast.success('Comanda criada com sucesso');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao criar comanda');
    },
  });
}

/**
 * Hook para criar comanda a partir de appointment
 */
export function useCreateCommandFromAppointment() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (appointmentId: string) => createCommandFromAppointment(appointmentId),
    onSuccess: (data, appointmentId) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.lists() });
      queryClient.invalidateQueries({ queryKey: commandKeys.byAppointment(appointmentId) });
      // Invalida appointment para atualizar status
      queryClient.invalidateQueries({ queryKey: ['appointments', 'detail', appointmentId] });
      toast.success('Comanda criada com sucesso');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao criar comanda');
    },
  });
}

/**
 * Hook para fechar comanda
 * 
 * Ao fechar a comanda com close-integrated, o backend registra operações no caixa
 * para TODOS os meios de pagamento. Por isso, invalidamos as queries do caixa.
 */
export function useCloseCommand() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ commandId, data }: { commandId: string; data: CloseCommandRequest }) =>
      closeCommand(commandId, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: commandKeys.lists() });
      if (data.appointment_id) {
        queryClient.invalidateQueries({ queryKey: commandKeys.byAppointment(data.appointment_id) });
        // Invalida appointment para atualizar status para DONE
        queryClient.invalidateQueries({ queryKey: ['appointments', 'detail', data.appointment_id] });
      }
      // Invalida queries do caixa - fechamento de comanda registra operações
      // no caixa diário para TODOS os meios de pagamento via close-integrated
      queryClient.invalidateQueries({ queryKey: ['caixa'] });
      toast.success('Comanda fechada com sucesso');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao fechar comanda');
    },
  });
}

/**
 * Hook para cancelar comanda
 */
export function useCancelCommand() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ commandId, motivo }: { commandId: string; motivo?: string }) =>
      cancelCommand(commandId, motivo),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: commandKeys.lists() });
      toast.success('Comanda cancelada');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao cancelar comanda');
    },
  });
}

// ============================================================================
// Mutations - Itens
// ============================================================================

/**
 * Hook para adicionar item à comanda
 */
export function useAddCommandItem() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ commandId, data }: { commandId: string; data: AddCommandItemRequest }) =>
      addCommandItem(commandId, data),
    onSuccess: (_, { commandId }) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(commandId) });
      toast.success('Item adicionado');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao adicionar item');
    },
  });
}

/**
 * Hook para atualizar item da comanda
 */
export function useUpdateCommandItem() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ 
      commandId, 
      itemId, 
      data 
    }: { 
      commandId: string; 
      itemId: string; 
      data: UpdateCommandItemRequest;
    }) => updateCommandItem(commandId, itemId, data),
    onSuccess: (_, { commandId }) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(commandId) });
      toast.success('Item atualizado');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao atualizar item');
    },
  });
}

/**
 * Hook para remover item da comanda
 */
export function useRemoveCommandItem() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ commandId, itemId }: { commandId: string; itemId: string }) =>
      removeCommandItem(commandId, itemId),
    onSuccess: (_, { commandId }) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(commandId) });
      toast.success('Item removido');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao remover item');
    },
  });
}

// ============================================================================
// Mutations - Pagamentos
// ============================================================================

/**
 * Hook para adicionar pagamento à comanda
 */
export function useAddCommandPayment() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ commandId, data }: { commandId: string; data: AddCommandPaymentRequest }) =>
      addCommandPayment(commandId, data),
    onSuccess: (_, { commandId }) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(commandId) });
      toast.success('Pagamento registrado');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao registrar pagamento');
    },
  });
}

/**
 * Hook para remover pagamento da comanda
 */
export function useRemoveCommandPayment() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ commandId, paymentId }: { commandId: string; paymentId: string }) =>
      removeCommandPayment(commandId, paymentId),
    onSuccess: (_, { commandId }) => {
      queryClient.invalidateQueries({ queryKey: commandKeys.detail(commandId) });
      toast.success('Pagamento removido');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error) || 'Erro ao remover pagamento');
    },
  });
}
