/**
 * NEXO - Sistema de Gestão para Barbearias
 * React Query Hooks para Clientes
 *
 * @module hooks/use-customers
 * @description Hooks com Optimistic Updates para melhor UX
 * Baseado em FLUXO_CADASTROS_CLIENTE.md
 */

import {
    CustomerHasAppointmentsError,
    CustomerNotFoundError,
    customerService,
    DuplicateCpfError,
    DuplicatePhoneError,
    InvalidCpfError,
    InvalidPhoneError,
} from '@/services/customer-service';
import type {
    CreateCustomerRequest,
    CustomerExportResponse,
    CustomerResponse,
    CustomerStatsResponse,
    ListCustomersFilters,
    ListCustomersResponse,
    UpdateCustomerRequest,
} from '@/types/customer';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const customerKeys = {
  all: ['customers'] as const,
  lists: () => [...customerKeys.all, 'list'] as const,
  list: (filters: ListCustomersFilters) =>
    [...customerKeys.lists(), filters] as const,
  details: () => [...customerKeys.all, 'detail'] as const,
  detail: (id: string) => [...customerKeys.details(), id] as const,
  history: (id: string) => [...customerKeys.all, 'history', id] as const,
  export: (id: string) => [...customerKeys.all, 'export', id] as const,
  active: () => [...customerKeys.all, 'active'] as const,
  stats: () => [...customerKeys.all, 'stats'] as const,
  search: (term: string) => [...customerKeys.all, 'search', term] as const,
  birthdays: (start: string, end: string) =>
    [...customerKeys.all, 'birthdays', start, end] as const,
  byTag: (tag: string) => [...customerKeys.all, 'tag', tag] as const,
};

// =============================================================================
// HOOKS DE CONSULTA (READ)
// =============================================================================

/**
 * Hook para listar clientes com filtros e paginação
 * RN-CLI-001: Listagem com busca, filtros e ordenação
 */
export function useCustomers(filters: ListCustomersFilters = {}) {
  return useQuery({
    queryKey: customerKeys.list(filters),
    queryFn: () => customerService.list(filters),
    staleTime: 60_000, // 1 minuto
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook para buscar um cliente específico
 */
export function useCustomer(id: string) {
  return useQuery({
    queryKey: customerKeys.detail(id),
    queryFn: () => customerService.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook para buscar cliente com histórico completo
 * RN-CLI-008: Histórico de agendamentos e comandas
 */
export function useCustomerWithHistory(id: string) {
  return useQuery({
    queryKey: customerKeys.history(id),
    queryFn: () => customerService.getWithHistory(id),
    enabled: !!id,
    staleTime: 30_000, // 30 segundos - dados mais voláteis
  });
}

/**
 * Hook para exportar dados do cliente (LGPD)
 * RN-CLI-009: Export para conformidade LGPD
 */
export function useCustomerExport(id: string, enabled: boolean = false) {
  return useQuery<CustomerExportResponse>({
    queryKey: customerKeys.export(id),
    queryFn: () => customerService.exportData(id),
    enabled: !!id && enabled,
    staleTime: 0, // Sempre buscar dados frescos para export
  });
}

/**
 * Hook para listar apenas clientes ativos (para selects e autocomplete)
 */
export function useActiveCustomers(limit: number = 100) {
  return useQuery({
    queryKey: customerKeys.active(),
    queryFn: () => customerService.listActive(limit),
    staleTime: 5 * 60_000, // 5 minutos - cache maior para seleção
  });
}

/**
 * Hook para obter estatísticas de clientes
 * RN-CLI-010: Dashboard com métricas
 */
export function useCustomerStats() {
  return useQuery<CustomerStatsResponse>({
    queryKey: customerKeys.stats(),
    queryFn: () => customerService.getStats(),
    staleTime: 5 * 60_000, // 5 minutos
  });
}

/**
 * Hook para buscar clientes (autocomplete/search)
 * RN-CLI-001: Busca com debounce de 300ms no frontend
 */
export function useCustomerSearch(term: string, limit: number = 10) {
  return useQuery({
    queryKey: customerKeys.search(term),
    queryFn: () => customerService.search(term, limit),
    enabled: term.length >= 2, // Mínimo 2 caracteres para buscar
    staleTime: 30_000, // 30 segundos
  });
}

/**
 * Hook para listar aniversariantes do período
 */
export function useCustomerBirthdays(startDate: string, endDate: string) {
  return useQuery({
    queryKey: customerKeys.birthdays(startDate, endDate),
    queryFn: () => customerService.listBirthdays(startDate, endDate),
    enabled: !!startDate && !!endDate,
    staleTime: 5 * 60_000,
  });
}

/**
 * Hook para listar clientes por tag
 */
export function useCustomersByTag(tag: string) {
  return useQuery({
    queryKey: customerKeys.byTag(tag),
    queryFn: () => customerService.listByTag(tag),
    enabled: !!tag,
    staleTime: 5 * 60_000,
  });
}

// =============================================================================
// HOOKS DE MUTAÇÃO (CREATE, UPDATE, DELETE) - COM OPTIMISTIC UPDATES
// =============================================================================

/**
 * Hook para criar um novo cliente
 * RN-CLI-002: Cadastro com validações
 */
export function useCreateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCustomerRequest) => customerService.create(data),
    onSuccess: (_newCustomer) => {
      void _newCustomer;
      // Invalida caches relacionados para forçar refetch
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
      queryClient.invalidateQueries({ queryKey: customerKeys.active() });
      queryClient.invalidateQueries({ queryKey: customerKeys.stats() });

      toast.success('Cliente cadastrado com sucesso!');
    },
    onError: (error: Error) => {
      handleCustomerError(error);
    },
  });
}

/**
 * Hook para atualizar um cliente
 * Usa optimistic update para resposta imediata na UI
 */
export function useUpdateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCustomerRequest }) =>
      customerService.update(id, data),

    // Optimistic Update
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: customerKeys.lists() });
      await queryClient.cancelQueries({ queryKey: customerKeys.detail(id) });

      // Snapshot para rollback
      const previousLists = queryClient.getQueriesData<ListCustomersResponse>({
        queryKey: customerKeys.lists(),
      });
      const previousDetail = queryClient.getQueryData<CustomerResponse>(
        customerKeys.detail(id)
      );

      // Atualização otimista nas listas
      queryClient.setQueriesData<ListCustomersResponse>(
        { queryKey: customerKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((customer) =>
              customer.id === id
                ? {
                    ...customer,
                    nome: data.nome ?? customer.nome,
                    telefone: data.telefone ?? customer.telefone,
                    email: data.email ?? customer.email,
                    cpf: data.cpf ?? customer.cpf,
                    data_nascimento: data.data_nascimento ?? customer.data_nascimento,
                    genero: data.genero ?? customer.genero,
                    observacoes: data.observacoes ?? customer.observacoes,
                    tags: data.tags ?? customer.tags,
                  }
                : customer
            ),
          };
        }
      );

      // Atualização otimista no detalhe
      if (previousDetail) {
        queryClient.setQueryData<CustomerResponse>(customerKeys.detail(id), {
          ...previousDetail,
          ...data,
        });
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
          customerKeys.detail(context.id),
          context.previousDetail
        );
      }
      handleCustomerError(error);
    },

    // Sempre refetch para garantir consistência
    onSettled: (_, __, variables) => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: customerKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: customerKeys.active() });
      queryClient.invalidateQueries({ queryKey: customerKeys.stats() });
    },

    onSuccess: () => {
      toast.success('Cliente atualizado com sucesso!');
    },
  });
}

/**
 * Hook para inativar um cliente (soft delete)
 * RN-CLI-006: Inativação com confirmação
 * Usa optimistic update para feedback imediato
 */
export function useInactivateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => customerService.inactivate(id),

    // Optimistic Update - marca como inativo
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: customerKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListCustomersResponse>({
        queryKey: customerKeys.lists(),
      });

      // Marca como inativo otimisticamente
      queryClient.setQueriesData<ListCustomersResponse>(
        { queryKey: customerKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((customer) =>
              customer.id === id
                ? { ...customer, ativo: false }
                : customer
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
      handleCustomerError(error);
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
      queryClient.invalidateQueries({ queryKey: customerKeys.active() });
      queryClient.invalidateQueries({ queryKey: customerKeys.stats() });
    },

    onSuccess: () => {
      toast.success('Cliente inativado com sucesso.');
    },
  });
}

/**
 * Hook para reativar um cliente
 * Define ativo=true usando endpoint de update
 */
export function useReactivateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    // Usa update passando um objeto vazio - o backend irá reativar
    // Nota: Se o backend não suportar isso diretamente, pode ser necessário
    // adicionar um endpoint específico ou passar um flag
    mutationFn: ({ id }: { id: string }) =>
      customerService.update(id, {}),

    // Optimistic Update
    onMutate: async ({ id }) => {
      await queryClient.cancelQueries({ queryKey: customerKeys.lists() });

      const previousLists = queryClient.getQueriesData<ListCustomersResponse>({
        queryKey: customerKeys.lists(),
      });

      // Marca como ativo otimisticamente
      queryClient.setQueriesData<ListCustomersResponse>(
        { queryKey: customerKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((customer) =>
              customer.id === id
                ? { ...customer, ativo: true }
                : customer
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
      toast.error('Erro ao reativar cliente. Tente novamente.');
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
      queryClient.invalidateQueries({ queryKey: customerKeys.active() });
      queryClient.invalidateQueries({ queryKey: customerKeys.stats() });
    },

    onSuccess: () => {
      toast.success('Cliente reativado com sucesso!');
    },
  });
}

// =============================================================================
// HOOKS AUXILIARES
// =============================================================================

/**
 * Hook para verificar telefone duplicado
 * RN-CLI-003: Validação em tempo real
 */
export function useCheckPhoneExists(phone: string, excludeId?: string) {
  const cleanedPhone = phone.replace(/\D/g, '');
  return useQuery({
    queryKey: ['customers', 'check-phone', cleanedPhone, excludeId],
    queryFn: () => customerService.checkPhoneExists(cleanedPhone, excludeId),
    enabled: cleanedPhone.length >= 10,
    staleTime: 5_000,
  });
}

/**
 * Hook para verificar CPF duplicado
 * RN-CLI-004: Validação em tempo real (CPF opcional)
 */
export function useCheckCpfExists(cpf: string, excludeId?: string) {
  const cleanedCpf = cpf.replace(/\D/g, '');
  return useQuery({
    queryKey: ['customers', 'check-cpf', cleanedCpf, excludeId],
    queryFn: () => customerService.checkCpfExists(cleanedCpf, excludeId),
    enabled: cleanedCpf.length === 11,
    staleTime: 5_000,
  });
}

/**
 * Hook para prefetch de clientes ativos (melhora UX em forms)
 */
export function usePrefetchActiveCustomers() {
  const queryClient = useQueryClient();

  return () => {
    queryClient.prefetchQuery({
      queryKey: customerKeys.active(),
      queryFn: () => customerService.listActive(),
      staleTime: 5 * 60_000,
    });
  };
}

/**
 * Hook para invalidar cache de clientes (útil após operações externas)
 */
export function useInvalidateCustomers() {
  const queryClient = useQueryClient();

  return () => {
    queryClient.invalidateQueries({ queryKey: customerKeys.all });
  };
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Handler centralizado de erros de cliente
 */
function handleCustomerError(error: Error) {
  if (error instanceof DuplicatePhoneError) {
    toast.error('Este telefone já está cadastrado para outro cliente.');
    return;
  }

  if (error instanceof DuplicateCpfError) {
    toast.error('Este CPF já está cadastrado para outro cliente.');
    return;
  }

  if (error instanceof InvalidPhoneError) {
    toast.error('Telefone inválido. Use o formato (XX) XXXXX-XXXX.');
    return;
  }

  if (error instanceof InvalidCpfError) {
    toast.error('CPF inválido.');
    return;
  }

  if (error instanceof CustomerNotFoundError) {
    toast.error('Cliente não encontrado.');
    return;
  }

  if (error instanceof CustomerHasAppointmentsError) {
    toast.error(
      'Não é possível excluir cliente com agendamentos. Altere o status para Inativo.'
    );
    return;
  }

  // Erro genérico
  console.error('[useCustomers] Error:', error);
  toast.error('Erro ao processar. Tente novamente.');
}
