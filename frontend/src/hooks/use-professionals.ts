/**
 * NEXO - Sistema de Gestão para Barbearias
 * React Query Hooks para Profissionais
 *
 * @module hooks/use-professionals
 * @description Hooks com Optimistic Updates para melhor UX
 */

import {
    DuplicateCpfError,
    DuplicateEmailError,
    InvalidCommissionError,
    ProfessionalHasAppointmentsError,
    ProfessionalNotFoundError,
    professionalService,
} from '@/services/professional-service';
import type {
    CreateProfessionalRequest,
    ListProfessionalsFilters,
    ListProfessionalsResponse,
    ProfessionalResponse,
    ProfessionalStatus,
    UpdateProfessionalRequest,
    UpdateProfessionalStatusRequest,
} from '@/types/professional';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const professionalKeys = {
  all: ['professionals'] as const,
  lists: () => [...professionalKeys.all, 'list'] as const,
  list: (filters: ListProfessionalsFilters) =>
    [...professionalKeys.lists(), filters] as const,
  details: () => [...professionalKeys.all, 'detail'] as const,
  detail: (id: string) => [...professionalKeys.details(), id] as const,
  barbers: () => [...professionalKeys.all, 'barbers'] as const,
  active: () => [...professionalKeys.all, 'active'] as const,
};

// =============================================================================
// HOOKS DE CONSULTA (READ)
// =============================================================================

/**
 * Hook para listar profissionais com filtros e paginação
 */
export function useProfessionals(filters: ListProfessionalsFilters = {}) {
  return useQuery({
    queryKey: professionalKeys.list(filters),
    queryFn: () => professionalService.list(filters),
    staleTime: 60_000, // 1 minuto
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook para buscar um profissional específico
 */
export function useProfessional(id: string) {
  return useQuery({
    queryKey: professionalKeys.detail(id),
    queryFn: () => professionalService.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook para listar apenas barbeiros ativos (para agendamentos)
 */
export function useBarbers() {
  return useQuery({
    queryKey: professionalKeys.barbers(),
    queryFn: () => professionalService.listBarbers(),
    staleTime: 5 * 60_000, // 5 minutos - cache maior para seleção
  });
}

/**
 * Hook para listar todos profissionais ativos (para selects)
 */
export function useActiveProfessionals() {
  return useQuery({
    queryKey: professionalKeys.active(),
    queryFn: () => professionalService.listActive(),
    staleTime: 5 * 60_000,
  });
}

// =============================================================================
// HOOKS DE MUTAÇÃO (CREATE, UPDATE, DELETE) - COM OPTIMISTIC UPDATES
// =============================================================================

/**
 * Hook para criar um novo profissional
 */
export function useCreateProfessional() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProfessionalRequest) =>
      professionalService.create(data),
    onSuccess: (newProfessional) => {
      // Adiciona o novo profissional ao cache
      queryClient.setQueriesData<ListProfessionalsResponse>(
        { queryKey: professionalKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: [newProfessional, ...old.data],
            meta: { ...old.meta, total: old.meta.total + 1 },
          };
        }
      );
      
      // Invalida caches relacionados
      queryClient.invalidateQueries({ queryKey: professionalKeys.lists() });
      queryClient.invalidateQueries({ queryKey: professionalKeys.barbers() });
      queryClient.invalidateQueries({ queryKey: professionalKeys.active() });
      
      toast.success('Profissional cadastrado com sucesso!');
    },
    onError: (error: Error) => {
      handleProfessionalError(error);
    },
  });
}

/**
 * Hook para atualizar um profissional
 * Usa optimistic update para resposta imediata na UI
 */
export function useUpdateProfessional() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProfessionalRequest }) =>
      professionalService.update(id, data),

    // Optimistic Update
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: professionalKeys.lists() });
      await queryClient.cancelQueries({ queryKey: professionalKeys.detail(id) });

      // Snapshot para rollback
      const previousLists = queryClient.getQueriesData<ListProfessionalsResponse>({
        queryKey: professionalKeys.lists(),
      });
      const previousDetail = queryClient.getQueryData<ProfessionalResponse>(
        professionalKeys.detail(id)
      );

      // Atualização otimista nas listas
      queryClient.setQueriesData<ListProfessionalsResponse>(
        { queryKey: professionalKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((prof) =>
              prof.id === id
                ? {
                    ...prof,
                    nome: data.nome ?? prof.nome,
                    email: data.email ?? prof.email,
                    telefone: data.telefone ?? prof.telefone,
                    foto: data.foto ?? prof.foto,
                    especialidades: data.especialidades ?? prof.especialidades,
                    tipo_comissao: data.tipo_comissao ?? prof.tipo_comissao,
                    comissao: data.comissao ?? prof.comissao,
                    status: data.status ?? prof.status,
                  }
                : prof
            ),
          };
        }
      );

      // Atualização otimista no detalhe
      if (previousDetail) {
        queryClient.setQueryData<ProfessionalResponse>(
          professionalKeys.detail(id),
          {
            ...previousDetail,
            ...data,
          }
        );
      }

      return { previousLists, previousDetail, id };
    },

    // Rollback em caso de erro
    onError: (error, _variables, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      if (context?.previousDetail) {
        queryClient.setQueryData(
          professionalKeys.detail(context.id),
          context.previousDetail
        );
      }
      handleProfessionalError(error);
    },

    // Sempre refetch para garantir consistência
    onSettled: (_, __, variables) => {
      queryClient.invalidateQueries({ queryKey: professionalKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: professionalKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: professionalKeys.barbers() });
      queryClient.invalidateQueries({ queryKey: professionalKeys.active() });
    },

    onSuccess: () => {
      toast.success('Profissional atualizado com sucesso!');
    },
  });
}

/**
 * Hook para atualizar o status de um profissional
 * Usa optimistic update para transições suaves
 */
export function useUpdateProfessionalStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: UpdateProfessionalStatusRequest;
    }) => professionalService.updateStatus(id, data),

    // Optimistic Update
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: professionalKeys.lists() });
      await queryClient.cancelQueries({ queryKey: professionalKeys.detail(id) });

      const previousLists = queryClient.getQueriesData<ListProfessionalsResponse>({
        queryKey: professionalKeys.lists(),
      });
      const previousDetail = queryClient.getQueryData<ProfessionalResponse>(
        professionalKeys.detail(id)
      );

      // Atualiza status otimisticamente nas listas
      queryClient.setQueriesData<ListProfessionalsResponse>(
        { queryKey: professionalKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((prof) =>
              prof.id === id ? { ...prof, status: data.status } : prof
            ),
          };
        }
      );

      // Atualiza status no detalhe
      if (previousDetail) {
        queryClient.setQueryData<ProfessionalResponse>(
          professionalKeys.detail(id),
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
          professionalKeys.detail(context.id),
          context.previousDetail
        );
      }
      toast.error('Erro ao atualizar status. Tente novamente.');
    },

    onSettled: (_, __, variables) => {
      queryClient.invalidateQueries({ queryKey: professionalKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: professionalKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: professionalKeys.barbers() });
      queryClient.invalidateQueries({ queryKey: professionalKeys.active() });
    },

    onSuccess: () => {
      toast.success('Status atualizado!');
    },
  });
}

/**
 * Hook para remover um profissional (soft delete)
 * Usa optimistic update para feedback imediato
 */
export function useDeleteProfessional() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => professionalService.delete(id),

    // Optimistic Update - marca como DEMITIDO
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: professionalKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListProfessionalsResponse>({
        queryKey: professionalKeys.lists(),
      });

      // Marca como DEMITIDO otimisticamente
      queryClient.setQueriesData<ListProfessionalsResponse>(
        { queryKey: professionalKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((prof) =>
              prof.id === id
                ? { ...prof, status: 'DEMITIDO' as ProfessionalStatus }
                : prof
            ),
          };
        }
      );

      return { previousLists };
    },

    onError: (error, _, context) => {
      if (context?.previousLists) {
        context.previousLists.forEach(([queryKey, data]) => {
          if (data) queryClient.setQueryData(queryKey, data);
        });
      }
      handleProfessionalError(error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: professionalKeys.lists() });
      queryClient.invalidateQueries({ queryKey: professionalKeys.barbers() });
      queryClient.invalidateQueries({ queryKey: professionalKeys.active() });
    },

    onSuccess: () => {
      toast.success('Profissional removido.');
    },
  });
}

// =============================================================================
// HOOKS AUXILIARES
// =============================================================================

/**
 * Hook para verificar email duplicado
 */
export function useCheckEmailExists(email: string, excludeId?: string) {
  return useQuery({
    queryKey: ['professionals', 'check-email', email, excludeId],
    queryFn: () => professionalService.checkEmailExists(email, excludeId),
    enabled: email.length > 3 && email.includes('@'),
    staleTime: 5_000,
  });
}

/**
 * Hook para verificar CPF duplicado
 */
export function useCheckCpfExists(cpf: string, excludeId?: string) {
  return useQuery({
    queryKey: ['professionals', 'check-cpf', cpf, excludeId],
    queryFn: () => professionalService.checkCpfExists(cpf, excludeId),
    enabled: cpf.replace(/\D/g, '').length === 11,
    staleTime: 5_000,
  });
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Handler centralizado de erros de profissional
 */
function handleProfessionalError(error: Error) {
  if (error instanceof DuplicateEmailError) {
    toast.error('Este email já está cadastrado para outro profissional.');
    return;
  }

  if (error instanceof DuplicateCpfError) {
    toast.error('Este CPF já está cadastrado para outro profissional.');
    return;
  }

  if (error instanceof InvalidCommissionError) {
    toast.error('Comissão é obrigatória para barbeiros (0-100%).');
    return;
  }

  if (error instanceof ProfessionalNotFoundError) {
    toast.error('Profissional não encontrado.');
    return;
  }

  if (error instanceof ProfessionalHasAppointmentsError) {
    toast.error('Não é possível excluir profissional com agendamentos. Altere o status para Inativo.');
    return;
  }

  if (error instanceof DuplicateCpfError) {
    toast.error('CPF/CNPJ já cadastrado no sistema.');
    return;
  }

  if (error instanceof DuplicateEmailError) {
    toast.error('Email já cadastrado no sistema.');
    return;
  }

  if (error instanceof InvalidCommissionError) {
    toast.error('Valor de comissão inválido.');
    return;
  }

  // Erro genérico
  console.error('[useProfessionals] Error:', error);
  toast.error('Erro ao processar. Tente novamente.');
}
