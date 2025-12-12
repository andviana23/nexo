/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal de Seleção de Unidade
 *
 * Modal full-screen obrigatório exibido após login.
 * O usuário deve escolher uma unidade antes de acessar o sistema.
 */

'use client';

import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import type { UserUnit } from '@/types/unit';
import { AnimatePresence, motion } from 'framer-motion';
import { Building2, Check, Loader2, MapPin } from 'lucide-react';

// =============================================================================
// TIPOS
// =============================================================================

interface UnitSelectionModalProps {
    isOpen: boolean;
    units: UserUnit[];
    isLoading: boolean;
    isSelecting: boolean;
    selectedUnitId: string | null;
    error: string | null;
    onSelect: (unit: UserUnit) => void;
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function UnitSelectionModal({
    isOpen,
    units,
    isLoading,
    isSelecting,
    selectedUnitId,
    error,
    onSelect,
}: UnitSelectionModalProps) {
    if (!isOpen) return null;

    return (
        <AnimatePresence>
            {isOpen && (
                <>
                    {/* Backdrop */}
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.3 }}
                        className="fixed inset-0 z-[100] bg-black/60 backdrop-blur-sm"
                        aria-hidden="true"
                    />

                    {/* Modal Container */}
                    <div className="fixed inset-0 z-[101] flex items-center justify-center p-4">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            transition={{ duration: 0.3, ease: [0.16, 1, 0.3, 1] }}
                            className="w-full max-w-md"
                            role="dialog"
                            aria-modal="true"
                            aria-labelledby="unit-selection-title"
                        >
                            {/* Card */}
                            <div className="relative overflow-hidden rounded-2xl bg-white dark:bg-slate-900 shadow-2xl shadow-black/20">
                                {/* Header Gradient */}
                                <div className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500" />

                                {/* Content */}
                                <div className="p-8">
                                    {/* Icon */}
                                    <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500/10 to-purple-500/10 dark:from-blue-500/20 dark:to-purple-500/20">
                                        <Building2 className="h-8 w-8 text-blue-600 dark:text-blue-400" />
                                    </div>

                                    {/* Title */}
                                    <h1
                                        id="unit-selection-title"
                                        className="text-center text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-100"
                                    >
                                        Selecione a Unidade
                                    </h1>

                                    {/* Subtitle */}
                                    <p className="mt-2 text-center text-sm text-slate-500 dark:text-slate-400">
                                        Escolha qual estabelecimento você deseja acessar
                                    </p>

                                    {/* Error Message */}
                                    {error && (
                                        <motion.div
                                            initial={{ opacity: 0, y: -10 }}
                                            animate={{ opacity: 1, y: 0 }}
                                            className="mt-4 flex items-center gap-2 rounded-xl bg-red-50 dark:bg-red-950/50 p-3 text-sm text-red-600 dark:text-red-400"
                                        >
                                            <div className="h-2 w-2 rounded-full bg-red-500 animate-pulse" />
                                            {error}
                                        </motion.div>
                                    )}

                                    {/* Units List */}
                                    <div className="mt-6 space-y-3">
                                        {isLoading ? (
                                            // Loading State
                                            <div className="flex flex-col items-center justify-center py-8">
                                                <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
                                                <p className="mt-3 text-sm text-slate-500 dark:text-slate-400">
                                                    Carregando unidades...
                                                </p>
                                            </div>
                                        ) : units.length === 0 ? (
                                            // Empty State
                                            <div className="text-center py-8">
                                                <p className="text-slate-500 dark:text-slate-400">
                                                    Nenhuma unidade disponível
                                                </p>
                                            </div>
                                        ) : (
                                            // Units Buttons
                                            units.map((unit, index) => (
                                                <motion.div
                                                    key={unit.unit_id}
                                                    initial={{ opacity: 0, y: 10 }}
                                                    animate={{ opacity: 1, y: 0 }}
                                                    transition={{ delay: index * 0.05 }}
                                                >
                                                    <Button
                                                        variant="outline"
                                                        className={cn(
                                                            'relative w-full h-auto py-4 px-5 justify-start text-left',
                                                            'rounded-xl border-2 transition-all duration-200',
                                                            'hover:border-blue-500/50 hover:bg-blue-50/50 dark:hover:bg-blue-950/30',
                                                            'focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2',
                                                            selectedUnitId === unit.unit_id && isSelecting
                                                                ? 'border-blue-500 bg-blue-50 dark:bg-blue-950/50'
                                                                : 'border-slate-200 dark:border-slate-700',
                                                            'disabled:opacity-50 disabled:cursor-not-allowed'
                                                        )}
                                                        onClick={() => onSelect(unit)}
                                                        disabled={isSelecting}
                                                    >
                                                        <div className="flex items-center gap-4 w-full">
                                                            {/* Unit Icon */}
                                                            <div
                                                                className={cn(
                                                                    'flex h-12 w-12 shrink-0 items-center justify-center rounded-xl',
                                                                    'bg-gradient-to-br transition-colors duration-200',
                                                                    unit.unit_matriz
                                                                        ? 'from-amber-500/20 to-orange-500/20 dark:from-amber-500/30 dark:to-orange-500/30'
                                                                        : 'from-slate-100 to-slate-200 dark:from-slate-800 dark:to-slate-700'
                                                                )}
                                                            >
                                                                <Building2
                                                                    className={cn(
                                                                        'h-6 w-6',
                                                                        unit.unit_matriz
                                                                            ? 'text-amber-600 dark:text-amber-400'
                                                                            : 'text-slate-600 dark:text-slate-400'
                                                                    )}
                                                                />
                                                            </div>

                                                            {/* Unit Info */}
                                                            <div className="flex-1 min-w-0">
                                                                <div className="flex items-center gap-2">
                                                                    <span className="font-semibold text-slate-900 dark:text-slate-100 truncate">
                                                                        {unit.unit_nome}
                                                                    </span>
                                                                    {unit.unit_matriz && (
                                                                        <span className="shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800 dark:bg-amber-900/50 dark:text-amber-300">
                                                                            Matriz
                                                                        </span>
                                                                    )}
                                                                </div>
                                                                {unit.unit_apelido && (
                                                                    <div className="flex items-center gap-1 mt-1 text-sm text-slate-500 dark:text-slate-400">
                                                                        <MapPin className="h-3.5 w-3.5" />
                                                                        <span className="truncate">{unit.unit_apelido}</span>
                                                                    </div>
                                                                )}
                                                            </div>

                                                            {/* Selection Indicator */}
                                                            {selectedUnitId === unit.unit_id && isSelecting && (
                                                                <div className="shrink-0 flex items-center justify-center h-6 w-6 rounded-full bg-blue-500">
                                                                    <Loader2 className="h-4 w-4 animate-spin text-white" />
                                                                </div>
                                                            )}

                                                            {selectedUnitId === unit.unit_id && !isSelecting && (
                                                                <div className="shrink-0 flex items-center justify-center h-6 w-6 rounded-full bg-green-500">
                                                                    <Check className="h-4 w-4 text-white" />
                                                                </div>
                                                            )}
                                                        </div>
                                                    </Button>
                                                </motion.div>
                                            ))
                                        )}
                                    </div>

                                    {/* Footer Note */}
                                    {!isLoading && units.length > 0 && (
                                        <p className="mt-6 text-center text-xs text-slate-400 dark:text-slate-500">
                                            Você pode trocar de unidade a qualquer momento pelo menu
                                        </p>
                                    )}
                                </div>
                            </div>
                        </motion.div>
                    </div>
                </>
            )}
        </AnimatePresence>
    );
}

export default UnitSelectionModal;
