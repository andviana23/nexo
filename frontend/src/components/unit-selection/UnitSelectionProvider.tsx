/**
 * NEXO - Sistema de Gestão para Barbearias
 * Provider de Seleção de Unidade
 *
 * Wrapper global que renderiza o modal de seleção quando necessário.
 * Deve envolver toda a aplicação.
 */

'use client';

import { UnitSelectionModal } from '@/components/unit-selection/UnitSelectionModal';
import { useUnitSelection } from '@/hooks/use-unit-selection';
import { Suspense } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface UnitSelectionProviderProps {
    children: React.ReactNode;
}

// =============================================================================
// COMPONENTE INTERNO (com useSearchParams que precisa de Suspense)
// =============================================================================

function UnitSelectionModalWrapper() {
    const {
        isModalOpen,
        units,
        isLoading,
        isSelecting,
        selectedUnitId,
        error,
        selectUnit,
    } = useUnitSelection();

    return (
        <UnitSelectionModal
            isOpen={isModalOpen}
            units={units}
            isLoading={isLoading}
            isSelecting={isSelecting}
            selectedUnitId={selectedUnitId}
            error={error}
            onSelect={selectUnit}
        />
    );
}

// =============================================================================
// PROVIDER
// =============================================================================

export function UnitSelectionProvider({ children }: UnitSelectionProviderProps) {
    return (
        <>
            {children}
            <Suspense fallback={null}>
                <UnitSelectionModalWrapper />
            </Suspense>
        </>
    );
}

export default UnitSelectionProvider;
