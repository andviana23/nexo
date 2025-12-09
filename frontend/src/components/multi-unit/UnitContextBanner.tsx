/**
 * NEXO - Sistema de Gestão para Barbearias
 * UnitContextBanner Component
 *
 * Banner que exibe o contexto da unidade ativa.
 * Aparece no topo da página quando há múltiplas unidades.
 */

'use client';

import { Button } from '@/components/ui/button';
import { useUnit } from '@/hooks/use-units';
import { cn } from '@/lib/utils';
import { Building2, Info, X } from 'lucide-react';
import { useCallback, useState } from 'react';
import { UnitSelector } from './UnitSelector';

// =============================================================================
// TIPOS
// =============================================================================

interface UnitContextBannerProps {
  /** Classe CSS adicional */
  className?: string;
  /** Permite fechar o banner */
  dismissible?: boolean;
  /** Variante visual */
  variant?: 'info' | 'subtle' | 'prominent';
  /** Exibe o seletor de unidade inline */
  showSelector?: boolean;
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function UnitContextBanner({
  className,
  dismissible = false,
  variant = 'subtle',
  showSelector = true,
}: UnitContextBannerProps) {
  const { activeUnit, isMultiUnit, isLoading } = useUnit();
  const [isDismissed, setIsDismissed] = useState(false);

  const handleDismiss = useCallback(() => {
    setIsDismissed(true);
  }, []);

  // Não mostra se não é multi-unidade ou foi fechado
  if (!isMultiUnit || isDismissed || isLoading) {
    return null;
  }

  // Não mostra se não tem unidade ativa
  if (!activeUnit) {
    return null;
  }

  const displayName = activeUnit.unit_apelido || activeUnit.unit_nome;

  const variantStyles = {
    info: 'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-950/50 dark:border-blue-800 dark:text-blue-200',
    subtle: 'bg-muted/50 border-border text-muted-foreground',
    prominent: 'bg-primary/10 border-primary/20 text-primary',
  };

  return (
    <div
      className={cn(
        'flex items-center justify-between gap-3 px-4 py-2 border-b text-sm',
        variantStyles[variant],
        className
      )}
    >
      <div className="flex items-center gap-2 min-w-0">
        <Building2 className="size-4 shrink-0" />
        <span className="font-medium">Trabalhando em:</span>
        <span className="truncate">{displayName}</span>
        {activeUnit.unit_matriz && (
          <span className="text-xs opacity-75">(Matriz)</span>
        )}
      </div>

      <div className="flex items-center gap-2 shrink-0">
        {showSelector && (
          <UnitSelector
            size="sm"
            variant="ghost"
            collapsible={false}
          />
        )}

        {dismissible && (
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={handleDismiss}
            className="size-6"
          >
            <X className="size-3" />
            <span className="sr-only">Fechar</span>
          </Button>
        )}
      </div>
    </div>
  );
}

// =============================================================================
// VARIANTE INLINE (para uso em cards/sections)
// =============================================================================

interface UnitContextInlineProps {
  className?: string;
}

export function UnitContextInline({ className }: UnitContextInlineProps) {
  const { activeUnit, isMultiUnit } = useUnit();

  if (!isMultiUnit || !activeUnit) {
    return null;
  }

  const displayName = activeUnit.unit_apelido || activeUnit.unit_nome;

  return (
    <div
      className={cn(
        'inline-flex items-center gap-1.5 text-xs text-muted-foreground',
        className
      )}
    >
      <Info className="size-3" />
      <span>
        Dados de: <strong>{displayName}</strong>
      </span>
    </div>
  );
}

// =============================================================================
// VARIANTE ALERTA (para confirmar operações cross-unit)
// =============================================================================

interface UnitConfirmAlertProps {
  className?: string;
  message?: string;
}

export function UnitConfirmAlert({
  className,
  message = 'Esta operação será realizada na unidade atual.',
}: UnitConfirmAlertProps) {
  const { activeUnit, isMultiUnit } = useUnit();

  if (!isMultiUnit || !activeUnit) {
    return null;
  }

  const displayName = activeUnit.unit_apelido || activeUnit.unit_nome;

  return (
    <div
      className={cn(
        'flex items-start gap-2 p-3 rounded-md',
        'bg-amber-50 border border-amber-200 text-amber-800',
        'dark:bg-amber-950/50 dark:border-amber-800 dark:text-amber-200',
        className
      )}
    >
      <Info className="size-4 shrink-0 mt-0.5" />
      <div className="text-sm">
        <p className="font-medium">Unidade: {displayName}</p>
        <p className="text-xs opacity-80 mt-0.5">{message}</p>
      </div>
    </div>
  );
}

export default UnitContextBanner;
