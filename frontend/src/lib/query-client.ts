/**
 * NEXO - Sistema de Gestão para Barbearias
 * TanStack Query Client Configuration
 *
 * Configuração centralizada do React Query para gerenciamento
 * de estado de servidor (cache, refetch, stale time, etc.)
 */

import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';

/**
 * Cache de queries com tratamento de erros global
 */
const queryCache = new QueryCache({
  onError: (error, query) => {
    // Log de erros em queries para debugging
    if (process.env.NODE_ENV === 'development') {
      console.error(`[Query Error] ${query.queryKey}:`, error);
    }
  },
});

/**
 * Cache de mutations com tratamento de erros global
 */
const mutationCache = new MutationCache({
  onError: (error, _variables, _context, mutation) => {
    // Log de erros em mutations para debugging
    if (process.env.NODE_ENV === 'development') {
      console.error(`[Mutation Error] ${mutation.mutationId}:`, error);
    }
  },
});

/**
 * Função factory para criar QueryClient
 * Necessária para SSR - cada request precisa de sua própria instância
 */
export function makeQueryClient(): QueryClient {
  return new QueryClient({
    queryCache,
    mutationCache,
    defaultOptions: {
      queries: {
        /**
         * Tempo que os dados ficam "fresh" antes de serem considerados stale
         * 5 minutos para a maioria dos dados
         */
        staleTime: 5 * 60 * 1000,

        /**
         * Tempo que dados ficam em cache após não serem mais usados
         * 30 minutos de garbage collection time
         */
        gcTime: 30 * 60 * 1000,

        /**
         * Número de tentativas em caso de erro
         */
        retry: (failureCount, error) => {
          // Não retry para erros de autenticação/autorização
          if (error instanceof Error && 'status' in error) {
            const status = (error as { status: number }).status;
            if (status === 401 || status === 403 || status === 404) {
              return false;
            }
          }
          // Máximo 3 retries para outros erros
          return failureCount < 3;
        },

        /**
         * Delay entre retries (exponential backoff)
         */
        retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),

        /**
         * Refetch automático quando a janela volta ao foco
         */
        refetchOnWindowFocus: true,

        /**
         * Refetch automático quando reconecta à internet
         */
        refetchOnReconnect: true,

        /**
         * Não refetch automaticamente ao montar o componente
         * se os dados ainda são fresh
         */
        refetchOnMount: true,
      },
      mutations: {
        /**
         * Número de tentativas para mutations
         * Mais conservador que queries
         */
        retry: 1,
      },
    },
  });
}

/**
 * Singleton para uso no cliente
 * Criado apenas uma vez e reutilizado
 */
let browserQueryClient: QueryClient | undefined;

/**
 * Obtém o QueryClient para uso no navegador
 * Cria uma nova instância se necessário
 */
export function getQueryClient(): QueryClient {
  // No servidor, sempre criar nova instância
  if (typeof window === 'undefined') {
    return makeQueryClient();
  }

  // No cliente, reutilizar instância existente
  if (!browserQueryClient) {
    browserQueryClient = makeQueryClient();
  }

  return browserQueryClient;
}

/**
 * Keys padronizadas para queries do sistema
 * Facilita invalidação e organização do cache
 */
export const queryKeys = {
  // Auth
  auth: {
    all: ['auth'] as const,
    me: () => [...queryKeys.auth.all, 'me'] as const,
  },

  // Usuários
  users: {
    all: ['users'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.users.all, 'list', filters] as const,
    detail: (id: string) => [...queryKeys.users.all, 'detail', id] as const,
  },

  // Clientes
  clients: {
    all: ['clients'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.clients.all, 'list', filters] as const,
    detail: (id: string) => [...queryKeys.clients.all, 'detail', id] as const,
  },

  // Serviços
  services: {
    all: ['services'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.services.all, 'list', filters] as const,
    detail: (id: string) => [...queryKeys.services.all, 'detail', id] as const,
  },

  // Profissionais
  professionals: {
    all: ['professionals'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.professionals.all, 'list', filters] as const,
    detail: (id: string) =>
      [...queryKeys.professionals.all, 'detail', id] as const,
  },

  // Agendamentos
  appointments: {
    all: ['appointments'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.appointments.all, 'list', filters] as const,
    detail: (id: string) =>
      [...queryKeys.appointments.all, 'detail', id] as const,
    today: () => [...queryKeys.appointments.all, 'today'] as const,
    week: () => [...queryKeys.appointments.all, 'week'] as const,
  },

  // Lista da Vez
  queue: {
    all: ['queue'] as const,
    current: () => [...queryKeys.queue.all, 'current'] as const,
    history: (date?: string) =>
      [...queryKeys.queue.all, 'history', date] as const,
  },

  // Financeiro
  financial: {
    all: ['financial'] as const,
    summary: (period?: string) =>
      [...queryKeys.financial.all, 'summary', period] as const,
    transactions: (filters?: Record<string, unknown>) =>
      [...queryKeys.financial.all, 'transactions', filters] as const,
  },

  // Dashboard
  dashboard: {
    all: ['dashboard'] as const,
    stats: () => [...queryKeys.dashboard.all, 'stats'] as const,
    charts: (period?: string) =>
      [...queryKeys.dashboard.all, 'charts', period] as const,
  },

  // Metas
  goals: {
    all: ['goals'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.goals.all, 'list', filters] as const,
    detail: (id: string) => [...queryKeys.goals.all, 'detail', id] as const,
  },

  // Unidades/Filiais
  units: {
    all: ['units'] as const,
    userUnits: () => [...queryKeys.units.all, 'user-units'] as const,
    detail: (id: string) => [...queryKeys.units.all, 'detail', id] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.units.all, 'list', filters] as const,
  },
} as const;

/**
 * Tipos utilitários para query keys
 */
export type QueryKeys = typeof queryKeys;
