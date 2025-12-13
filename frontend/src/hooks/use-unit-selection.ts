/**
 * NEXO - Sistema de Gestão para Barbearias
 * Hook de Seleção de Unidade
 *
 * Gerencia a lógica de exibição do modal de seleção de unidade.
 * Busca unidades do usuário e chama switch ao selecionar.
 */

'use client';

import { useAuthHydrated, useIsAuthenticated } from '@/store/auth-store';
import { useUnitStore } from '@/store/unit-store';
import type { UserUnit } from '@/types/unit';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { unitService } from '@/services/unit-service';
import { queryKeys } from '@/lib/query-client';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface UseUnitSelectionReturn {
    // Estado
    isModalOpen: boolean;
    units: UserUnit[];
    isLoading: boolean;
    isSelecting: boolean;
    selectedUnitId: string | null;
    error: string | null;

    // Ações
    selectUnit: (unit: UserUnit) => Promise<void>;
}

// =============================================================================
// HOOK
// =============================================================================

export function useUnitSelection(): UseUnitSelectionReturn {
    const router = useRouter();
    const searchParams = useSearchParams();
    const queryClient = useQueryClient();

    // Stores
    const isAuthenticated = useIsAuthenticated();
    const isAuthHydrated = useAuthHydrated();
    const {
        needsSelection,
        activeUnit,
        setUnits,
        setActiveUnit,
        setNeedsSelection,
        isHydrated: isUnitHydrated,
    } = useUnitStore();

    // ==========================================================================
    // QUERY: Buscar unidades do usuário
    // ==========================================================================

    const {
        data: unitsData,
        isLoading: isLoadingUnits,
        error: unitsError,
    } = useQuery({
        queryKey: queryKeys.units.me(),
        queryFn: () => unitService.getUserUnits(),
        enabled: isAuthenticated && isAuthHydrated && needsSelection,
        staleTime: 5 * 60 * 1000, // 5 minutos
        retry: 2,
    });

    // ==========================================================================
    // MUTATION: Selecionar unidade (switch)
    // ==========================================================================

    const switchMutation = useMutation({
        mutationFn: (unitId: string) => unitService.switchUnit(unitId),
        onSuccess: (data, unitId) => {
            // Atualiza unidade ativa no store
            const selectedUnit = unitsData?.units.find((u) => u.unit_id === unitId);
            if (selectedUnit) {
                setActiveUnit(selectedUnit);
            }

            // Remove flag de necessidade de seleção
            setNeedsSelection(false);

            // Invalida queries para refetch com nova unidade
            queryClient.invalidateQueries({ queryKey: queryKeys.units.all });

            // Redireciona para a página desejada ou home
            const returnUrl = searchParams.get('returnUrl');
            setTimeout(() => {
                router.push(returnUrl || '/');
            }, 100);
        },
        onError: (error) => {
            console.error('[useUnitSelection] Erro ao selecionar unidade:', error);
        },
    });

    const switchUnitAsync = switchMutation.mutateAsync;

    // ==========================================================================
    // HANDLERS
    // ==========================================================================

    async function handleSelectUnit(unit: UserUnit): Promise<void> {
        await switchUnitAsync(unit.unit_id);
    }

    // Sincroniza unidades com store quando carrega
    useEffect(() => {
        if (unitsData?.units) {
            setUnits(unitsData.units);

            // Auto-seleciona se só tem 1 unidade
            if (unitsData.units.length === 1) {
                void switchUnitAsync(unitsData.units[0].unit_id);
            }
        }
    }, [unitsData, setUnits, switchUnitAsync]);

    // ==========================================================================
    // COMPUTED: Modal deve estar aberto?
    // ==========================================================================

    // Modal abre quando:
    // 1. Usuário está autenticado
    // 2. Auth store está hidratado
    // 3. Unit store está hidratado
    // 4. Flag needsSelection é true
    // 5. Não tem unidade ativa
    const isModalOpen =
        isAuthenticated &&
        isAuthHydrated &&
        isUnitHydrated &&
        needsSelection &&
        !activeUnit;

    // ==========================================================================
    // RETURN
    // ==========================================================================

    return {
        isModalOpen,
        units: unitsData?.units || [],
        isLoading: isLoadingUnits,
        isSelecting: switchMutation.isPending,
        selectedUnitId: switchMutation.variables || null,
        error: unitsError?.message || switchMutation.error?.message || null,
        selectUnit: handleSelectUnit,
    };
}

export default useUnitSelection;
