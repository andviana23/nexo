/**
 * NEXO - Sistema de Gestão para Barbearias
 * Hook de Autenticação
 *
 * Hook que combina Zustand (estado) + React Query (server state) para auth.
 */

'use client';

import { getErrorMessage, isAxiosError } from '@/lib/axios';
import { queryKeys } from '@/lib/query-client';
import { authService, InvalidCredentialsError } from '@/services/auth-service';
import {
  useAuthHydrated,
  useAuthLoading,
  useAuthStore,
  useCurrentTenant,
  useCurrentUser,
  useIsAuthenticated,
} from '@/store/auth-store';
import type { LoginCredentials, User } from '@/types';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { useCallback } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface UseAuthReturn {
  // Estado
  user: User | null;
  tenant: ReturnType<typeof useCurrentTenant>;
  isAuthenticated: boolean;
  isLoading: boolean;
  isHydrated: boolean;

  // Ações
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;

  // Status das mutations
  isLoggingIn: boolean;
  isLoggingOut: boolean;
  loginError: string | null;
}

// =============================================================================
// HOOK PRINCIPAL
// =============================================================================

export function useAuth(): UseAuthReturn {
  const router = useRouter();
  const queryClient = useQueryClient();

  // Estado do Zustand
  const { setAuth, logout: storeLogout, setLoading } = useAuthStore();
  const user = useCurrentUser();
  const tenant = useCurrentTenant();
  const isAuthenticated = useIsAuthenticated();
  const isLoading = useAuthLoading();
  const isHydrated = useAuthHydrated();

  // ==========================================================================
  // LOGIN MUTATION
  // ==========================================================================

  const loginMutation = useMutation({
    mutationFn: (credentials: LoginCredentials) =>
      authService.login(credentials),

    onSuccess: (data) => {
      // Salva no Zustand
      setAuth(data.token, data.user, data.tenant);

      // Invalida queries de auth para refetch
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.all });

      // Redireciona para dashboard
      router.push('/dashboard');
    },

    onError: (error) => {
      setLoading(false);

      // Trata erros específicos
      if (isAxiosError(error)) {
        const status = error.response?.status;

        if (status === 401) {
          throw new InvalidCredentialsError();
        }
      }
    },
  });

  // ==========================================================================
  // LOGOUT MUTATION
  // ==========================================================================

  const logoutMutation = useMutation({
    mutationFn: () => authService.logout(),

    onSettled: () => {
      // Limpa estado local
      storeLogout();

      // Limpa todas as queries em cache
      queryClient.clear();

      // Redireciona para login
      router.push('/login');
    },
  });

  // ==========================================================================
  // QUERY PARA USUÁRIO ATUAL (ME)
  // ==========================================================================

  const { refetch: refreshUser } = useQuery({
    queryKey: queryKeys.auth.me(),
    queryFn: () => authService.getMe(),
    enabled: isAuthenticated && isHydrated,
    staleTime: 5 * 60 * 1000, // 5 minutos
    retry: false,
  });

  // ==========================================================================
  // AÇÕES
  // ==========================================================================

  const login = useCallback(
    async (credentials: LoginCredentials) => {
      setLoading(true);
      await loginMutation.mutateAsync(credentials);
    },
    [loginMutation, setLoading]
  );

  const logout = useCallback(async () => {
    await logoutMutation.mutateAsync();
  }, [logoutMutation]);

  const handleRefreshUser = useCallback(async () => {
    await refreshUser();
  }, [refreshUser]);

  // ==========================================================================
  // RETURN
  // ==========================================================================

  return {
    // Estado
    user,
    tenant,
    isAuthenticated,
    isLoading,
    isHydrated,

    // Ações
    login,
    logout,
    refreshUser: handleRefreshUser,

    // Status das mutations
    isLoggingIn: loginMutation.isPending,
    isLoggingOut: logoutMutation.isPending,
    loginError: loginMutation.error
      ? getErrorMessage(loginMutation.error)
      : null,
  };
}

// =============================================================================
// HOOKS AUXILIARES
// =============================================================================

/**
 * Hook para verificar se usuário tem determinada role
 */
export function useRequireRole(allowedRoles: string[]): {
  hasAccess: boolean;
  isLoading: boolean;
} {
  const user = useCurrentUser();
  const isLoading = useAuthLoading();
  const isHydrated = useAuthHydrated();

  if (!isHydrated || isLoading) {
    return { hasAccess: false, isLoading: true };
  }

  const hasAccess = user?.role ? allowedRoles.includes(user.role) : false;

  return { hasAccess, isLoading: false };
}

/**
 * Hook para verificar se é admin
 */
export function useIsAdmin(): boolean {
  const user = useCurrentUser();
  return user?.role === 'admin';
}

/**
 * Hook para verificar se é manager ou admin
 */
export function useIsManagerOrAbove(): boolean {
  const user = useCurrentUser();
  return user?.role === 'admin' || user?.role === 'manager';
}

/**
 * Hook para obter nome de exibição do usuário
 */
export function useDisplayName(): string {
  const user = useCurrentUser();
  if (!user) return '';
  return user.name.split(' ')[0]; // Primeiro nome
}

/**
 * Hook para obter iniciais do usuário (para avatar)
 */
export function useUserInitials(): string {
  const user = useCurrentUser();
  if (!user) return '';

  const names = user.name.split(' ');
  if (names.length === 1) {
    return names[0].substring(0, 2).toUpperCase();
  }

  return (names[0][0] + names[names.length - 1][0]).toUpperCase();
}
