/**
 * NEXO - Sistema de Gestão para Barbearias
 * UnitSelector Component
 *
 * Componente para seleção de unidade/filial ativa.
 * Usa Design System shadcn/ui com tokens do Tailwind.
 */

'use client';

import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useUnit } from '@/hooks/use-units';
import { analytics } from '@/lib/analytics';
import { cn } from '@/lib/utils';
import type { UserUnit } from '@/types/unit';
import {
    Building2,
    Check,
    ChevronDown,
    MapPin,
    Star,
} from 'lucide-react';

// =============================================================================
// TIPOS
// =============================================================================

interface UnitSelectorProps {
  /** Classe CSS adicional */
  className?: string;
  /** Tamanho do componente */
  size?: 'sm' | 'default' | 'lg';
  /** Variante do botão trigger */
  variant?: 'default' | 'outline' | 'ghost';
  /** Exibir apenas ícone em mobile */
  collapsible?: boolean;
  /** Callback quando unidade é trocada */
  onUnitChange?: (unit: UserUnit) => void;
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function UnitSelector({
  className,
  size = 'default',
  variant = 'outline',
  collapsible = true,
  onUnitChange,
}: UnitSelectorProps) {
  const {
    units,
    activeUnit,
    isMultiUnit,
    isLoading,
    isSwitching,
    switchUnit,
  } = useUnit();

  // Se não tem múltiplas unidades, não mostra o seletor
  if (!isMultiUnit || units.length <= 1) {
    return null;
  }

  const handleOpenChange = (open: boolean) => {
    if (open) {
      analytics.trackUnitSelectorOpen();
    }
  };

  const handleSelectUnit = async (unit: UserUnit) => {
    if (unit.unit_id === activeUnit?.unit_id) return;

    await switchUnit(unit.unit_id);
    onUnitChange?.(unit);
  };

  const displayName = activeUnit?.unit_apelido || activeUnit?.unit_nome || 'Selecione';

  return (
    <DropdownMenu onOpenChange={handleOpenChange}>
      <DropdownMenuTrigger asChild>
        <Button
          variant={variant}
          size={size}
          className={cn(
            'gap-2 font-medium',
            collapsible && 'max-w-[200px] truncate',
            className
          )}
          disabled={isLoading || isSwitching}
        >
          <Building2 className="size-4 shrink-0" />
          <span className={cn(collapsible && 'hidden sm:inline', 'truncate')}>
            {isLoading ? 'Carregando...' : displayName}
          </span>
          <ChevronDown className="size-4 shrink-0 opacity-50" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="start" className="w-[240px]">
        <DropdownMenuLabel className="flex items-center gap-2 text-xs text-muted-foreground">
          <MapPin className="size-3" />
          Trocar unidade
        </DropdownMenuLabel>
        <DropdownMenuSeparator />

        {units.map((unit) => {
          const isActive = unit.unit_id === activeUnit?.unit_id;
          const displayUnitName = unit.unit_apelido || unit.unit_nome;

          return (
            <DropdownMenuItem
              key={unit.unit_id}
              onClick={() => handleSelectUnit(unit)}
              className={cn(
                'flex items-center justify-between gap-2 cursor-pointer',
                isActive && 'bg-accent'
              )}
            >
              <div className="flex items-center gap-2 min-w-0">
                <Building2
                  className={cn(
                    'size-4 shrink-0',
                    unit.unit_matriz ? 'text-primary' : 'text-muted-foreground'
                  )}
                />
                <div className="flex flex-col min-w-0">
                  <span className="truncate font-medium">{displayUnitName}</span>
                  {unit.unit_matriz && (
                    <span className="text-xs text-muted-foreground">Matriz</span>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-1 shrink-0">
                {unit.is_default && (
                  <Star className="size-3 text-amber-500 fill-amber-500" />
                )}
                {isActive && <Check className="size-4 text-primary" />}
              </div>
            </DropdownMenuItem>
          );
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

// =============================================================================
// COMPONENTE COMPACTO (para mobile/sidebar)
// =============================================================================

export function UnitSelectorCompact({
  className,
  onUnitChange,
}: Pick<UnitSelectorProps, 'className' | 'onUnitChange'>) {
  return (
    <UnitSelector
      className={className}
      size="sm"
      variant="ghost"
      collapsible={false}
      onUnitChange={onUnitChange}
    />
  );
}

// =============================================================================
// INDICADOR DE UNIDADE ATIVA (somente leitura)
// =============================================================================

interface ActiveUnitBadgeProps {
  className?: string;
}

export function ActiveUnitBadge({ className }: ActiveUnitBadgeProps) {
  const { activeUnit, isMultiUnit } = useUnit();

  // Se não tem múltiplas unidades, não mostra o badge
  if (!isMultiUnit || !activeUnit) {
    return null;
  }

  const displayName = activeUnit.unit_apelido || activeUnit.unit_nome;

  return (
    <div
      className={cn(
        'inline-flex items-center gap-1.5 px-2 py-1 text-xs font-medium',
        'bg-secondary text-secondary-foreground rounded-md',
        className
      )}
    >
      <Building2 className="size-3" />
      <span className="truncate max-w-[120px]">{displayName}</span>
      {activeUnit.unit_matriz && (
        <span className="text-[10px] text-muted-foreground">(Matriz)</span>
      )}
    </div>
  );
}

export default UnitSelector;
