/**
 * NEXO - Sistema de Gestão para Barbearias
 * Auth Store (Zustand)
 *
 * Gerenciamento de estado de autenticação do cliente.
 * Persiste token no localStorage e sincroniza com cookies para SSR.
 */

import type { Tenant, User } from '@/types';
import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

// =============================================================================
// TIPOS
// =============================================================================

interface AuthState {
  // Estado
  token: string | null;
  user: User | null;
  tenant: Tenant | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isHydrated: boolean;

  // Ações
  setAuth: (token: string, user: User, tenant: Tenant) => void;
  updateUser: (user: Partial<User>) => void;
  updateTenant: (tenant: Partial<Tenant>) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;
  setHydrated: () => void;
}

// =============================================================================
// CONSTANTES
// =============================================================================

const AUTH_STORAGE_KEY = 'nexo-auth';
const TOKEN_COOKIE_NAME = 'nexo-token';

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Define um cookie com configurações seguras
 */
function setCookie(name: string, value: string, days = 7): void {
  if (typeof window === 'undefined') return;

  // DEBUG: Ver o que está sendo salvo
  console.log('🍪 [setCookie] Salvando cookie:', {
    name,
    valueType: typeof value,
    valueLength: value?.length,
    valuePreview: value?.substring(0, 30),
    isUndefined: value === undefined,
    isNull: value === null,
    isString: typeof value === 'string',
    rawValue: value,
  });

  const expires = new Date();
  expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);

  // Configurações de segurança
  const secure = window.location.protocol === 'https:';
  const sameSite = 'Lax';

  const cookieString = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=${sameSite}${
    secure ? ';Secure' : ''
  }`;

  console.log('🍪 [setCookie] Cookie string:', cookieString.substring(0, 100) + '...');

  document.cookie = cookieString;

  // Verifica se foi salvo
  const saved = getCookie(name);
  console.log('🍪 [setCookie] Verificação após salvar:', {
    wasSaved: !!saved,
    savedValue: saved?.substring(0, 30),
  });
}

/**
 * Remove um cookie
 */
function removeCookie(name: string): void {
  if (typeof window === 'undefined') return;
  document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
}

/**
 * Obtém um cookie pelo nome
 */
function getCookie(name: string): string | null {
  if (typeof window === 'undefined') return null;

  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);

  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() || null;
  }

  return null;
}

// =============================================================================
// STORE
// =============================================================================

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // Estado inicial
      token: null,
      user: null,
      tenant: null,
      isAuthenticated: false,
      isLoading: true,
      isHydrated: false,

      /**
       * Define autenticação após login bem-sucedido
       */
      setAuth: (token, user, tenant) => {
        // DEBUG: Ver o que está sendo passado
        console.log('🔐 [setAuth] Recebeu:', {
          tokenType: typeof token,
          tokenLength: token?.length,
          tokenPreview: token?.substring(0, 30),
          hasUser: !!user,
          hasTenant: !!tenant,
        });

        // Salva token no cookie para SSR/middleware
        setCookie(TOKEN_COOKIE_NAME, token, 7);

        set({
          token,
          user,
          tenant,
          isAuthenticated: true,
          isLoading: false,
        });
      },

      /**
       * Atualiza dados do usuário parcialmente
       */
      updateUser: (userData) => {
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null,
        }));
      },

      /**
       * Atualiza dados do tenant parcialmente
       */
      updateTenant: (tenantData) => {
        set((state) => ({
          tenant: state.tenant ? { ...state.tenant, ...tenantData } : null,
        }));
      },

      /**
       * Efetua logout e limpa estado
       */
      logout: () => {
        // Remove cookie
        removeCookie(TOKEN_COOKIE_NAME);

        // Limpa estado
        set({
          token: null,
          user: null,
          tenant: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      /**
       * Define estado de loading
       */
      setLoading: (loading) => {
        set({ isLoading: loading });
      },

      /**
       * Marca store como hidratado (sincronizado com localStorage)
       */
      setHydrated: () => {
        set({ isHydrated: true, isLoading: false });
      },
    }),
    {
      name: AUTH_STORAGE_KEY,
      storage: createJSONStorage(() => localStorage),

      // Apenas persiste esses campos
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        tenant: state.tenant,
        isAuthenticated: state.isAuthenticated,
      }),

      // Callback quando hidratação completa
      onRehydrateStorage: () => (state) => {
        if (state) {
          // Verifica se token no localStorage bate com cookie
          const cookieToken = getCookie(TOKEN_COOKIE_NAME);
          if (state.token && !cookieToken) {
            // Token no localStorage mas não no cookie - re-sincroniza
            setCookie(TOKEN_COOKIE_NAME, state.token, 7);
          } else if (!state.token && cookieToken) {
            // Cookie existe mas localStorage não - limpa cookie
            removeCookie(TOKEN_COOKIE_NAME);
          }

          state.setHydrated();
        }
      },
    }
  )
);

// =============================================================================
// SELETORES OTIMIZADOS
// =============================================================================

/**
 * Seleciona apenas o token
 */
export const useAuthToken = () => useAuthStore((state) => state.token);

/**
 * Seleciona apenas o usuário
 */
export const useCurrentUser = () => useAuthStore((state) => state.user);

/**
 * Seleciona apenas o tenant
 */
export const useCurrentTenant = () => useAuthStore((state) => state.tenant);

/**
 * Seleciona estado de autenticação
 */
export const useIsAuthenticated = () =>
  useAuthStore((state) => state.isAuthenticated);

/**
 * Seleciona estado de loading
 */
export const useAuthLoading = () => useAuthStore((state) => state.isLoading);

/**
 * Seleciona estado de hidratação
 */
export const useAuthHydrated = () => useAuthStore((state) => state.isHydrated);

/**
 * Seleciona role do usuário
 */
export const useUserRole = () => useAuthStore((state) => state.user?.role);

/**
 * Verifica se usuário tem permissão específica
 */
export const useHasRole = (roles: string[]) =>
  useAuthStore((state) => {
    if (!state.user?.role) return false;
    return roles.includes(state.user.role);
  });

// =============================================================================
// HELPERS PARA ACESSO FORA DE COMPONENTES
// =============================================================================

/**
 * Obtém estado atual (para uso em funções não-React)
 */
export const getAuthState = () => useAuthStore.getState();

/**
 * Obtém token atual (para uso em interceptors)
 */
export const getToken = () => useAuthStore.getState().token;

/**
 * Verifica se está autenticado (para uso em funções)
 */
export const isAuthenticated = () => useAuthStore.getState().isAuthenticated;

/**
 * Efetua logout (para uso em interceptors)
 */
export const logout = () => useAuthStore.getState().logout();
