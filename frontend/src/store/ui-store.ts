/**
 * NEXO - Sistema de Gestão para Barbearias
 * UI Store (Zustand)
 *
 * Gerencia estado global da interface:
 * - Sidebar (open/collapsed/mobile)
 * - Theme (light/dark/system)
 * - Breadcrumbs
 * - Loading states globais
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// =============================================================================
// TIPOS
// =============================================================================

interface UIState {
  // Sidebar
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  isMobile: boolean;

  // Breadcrumbs
  breadcrumbs: Array<{ label: string; href?: string }>;

  // Loading Global
  isLoading: boolean;
  loadingMessage?: string;

  // Actions
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleSidebarCollapse: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setIsMobile: (isMobile: boolean) => void;
  setBreadcrumbs: (breadcrumbs: Array<{ label: string; href?: string }>) => void;
  setLoading: (isLoading: boolean, message?: string) => void;
}

// =============================================================================
// STORE
// =============================================================================

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      // Estado inicial
      sidebarOpen: true,
      sidebarCollapsed: false,
      isMobile: false,
      breadcrumbs: [],
      isLoading: false,
      loadingMessage: undefined,

      // Actions
      toggleSidebar: () =>
        set((state) => ({ sidebarOpen: !state.sidebarOpen })),

      setSidebarOpen: (open) => set({ sidebarOpen: open }),

      toggleSidebarCollapse: () =>
        set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),

      setIsMobile: (isMobile) => set({ isMobile }),

      setBreadcrumbs: (breadcrumbs) => set({ breadcrumbs }),

      setLoading: (isLoading, message) =>
        set({ isLoading, loadingMessage: message }),
    }),
    {
      name: 'nexo-ui-storage',
      // Apenas persiste preferências de UI, não estados temporários
      partialize: (state) => ({
        sidebarCollapsed: state.sidebarCollapsed,
      }),
    }
  )
);

// =============================================================================
// HOOKS AUXILIARES
// =============================================================================

/**
 * Hook para controle da sidebar
 */
export const useSidebar = () => {
  const {
    sidebarOpen,
    sidebarCollapsed,
    isMobile,
    toggleSidebar,
    setSidebarOpen,
    toggleSidebarCollapse,
    setSidebarCollapsed,
  } = useUIStore();

  return {
    isOpen: sidebarOpen,
    isCollapsed: sidebarCollapsed,
    isMobile,
    toggle: toggleSidebar,
    setOpen: setSidebarOpen,
    toggleCollapse: toggleSidebarCollapse,
    setCollapsed: setSidebarCollapsed,
  };
};

/**
 * Hook para breadcrumbs
 */
export const useBreadcrumbs = () => {
  const { breadcrumbs, setBreadcrumbs } = useUIStore();

  return {
    breadcrumbs,
    setBreadcrumbs,
  };
};

/**
 * Hook para loading global
 */
export const useGlobalLoading = () => {
  const { isLoading, loadingMessage, setLoading } = useUIStore();

  return {
    isLoading,
    message: loadingMessage,
    setLoading,
    startLoading: (message?: string) => setLoading(true, message),
    stopLoading: () => setLoading(false),
  };
};
