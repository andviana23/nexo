/**
 * NEXO - Sistema de Gestão para Barbearias
 * Unit Store (Zustand)
 *
 * Gerenciamento de estado de unidade/filial ativa.
 * Persiste no localStorage e sincroniza com headers HTTP.
 */

'use client';

import type { UserUnit } from '@/types/unit';
import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

// =============================================================================
// TIPOS
// =============================================================================

interface UnitState {
  // Estado
  units: UserUnit[];
  activeUnit: UserUnit | null;
  isLoading: boolean;
  isHydrated: boolean;
  error: string | null;

  // Computed
  isMultiUnit: boolean;

  // Ações
  setUnits: (units: UserUnit[]) => void;
  setActiveUnit: (unit: UserUnit) => void;
  clearActiveUnit: () => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setHydrated: () => void;
  reset: () => void;
}

// =============================================================================
// CONSTANTES
// =============================================================================

const UNIT_STORAGE_KEY = 'nexo-unit';

// =============================================================================
// STORE
// =============================================================================

const initialState = {
  units: [],
  activeUnit: null,
  isLoading: false,
  isHydrated: false,
  error: null,
  isMultiUnit: false,
};

export const useUnitStore = create<UnitState>()(
  persist(
    (set) => ({
      ...initialState,

      // Ações
      setUnits: (units) =>
        set({
          units,
          isMultiUnit: units.length > 1,
          error: null,
        }),

      setActiveUnit: (unit) => {
        console.log('📍 [UnitStore] Definindo unidade ativa:', {
          id: unit.unit_id,
          nome: unit.unit_nome,
          isMatriz: unit.unit_matriz,
        });
        set({ activeUnit: unit, error: null });
      },

      clearActiveUnit: () => set({ activeUnit: null }),

      setLoading: (isLoading) => set({ isLoading }),

      setError: (error) => set({ error }),

      setHydrated: () => set({ isHydrated: true }),

      reset: () => {
        console.log('📍 [UnitStore] Reset do estado');
        set(initialState);
      },
    }),
    {
      name: UNIT_STORAGE_KEY,
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        activeUnit: state.activeUnit,
        // Não persiste units pois são carregadas do servidor
      }),
      onRehydrateStorage: () => {
        return (state) => {
          if (state) {
            state.setHydrated();
            console.log('📍 [UnitStore] Hidratado do localStorage:', {
              hasActiveUnit: !!state.activeUnit,
              activeUnitName: state.activeUnit?.unit_nome,
            });
          }
        };
      },
    }
  )
);

// =============================================================================
// HOOKS SELETORES (evita re-renders desnecessários)
// =============================================================================

/**
 * Retorna a lista de unidades do usuário
 */
export function useUnits(): UserUnit[] {
  return useUnitStore((state) => state.units);
}

/**
 * Retorna a unidade ativa atual
 */
export function useActiveUnit(): UserUnit | null {
  return useUnitStore((state) => state.activeUnit);
}

/**
 * Retorna o ID da unidade ativa (útil para headers)
 */
export function useActiveUnitId(): string | null {
  return useUnitStore((state) => state.activeUnit?.unit_id ?? null);
}

/**
 * Retorna se o tenant tem múltiplas unidades
 */
export function useIsMultiUnit(): boolean {
  return useUnitStore((state) => state.isMultiUnit);
}

/**
 * Retorna se o store foi hidratado do localStorage
 */
export function useUnitHydrated(): boolean {
  return useUnitStore((state) => state.isHydrated);
}

/**
 * Retorna se está carregando unidades
 */
export function useUnitLoading(): boolean {
  return useUnitStore((state) => state.isLoading);
}

/**
 * Retorna erro se houver
 */
export function useUnitError(): string | null {
  return useUnitStore((state) => state.error);
}

// =============================================================================
// HELPER - GET ACTIVE UNIT ID (para uso fora de React)
// =============================================================================

/**
 * Obtém o ID da unidade ativa de forma síncrona
 * Útil para interceptors do axios
 */
export function getActiveUnitId(): string | null {
  const state = useUnitStore.getState();
  return state.activeUnit?.unit_id ?? null;
}
