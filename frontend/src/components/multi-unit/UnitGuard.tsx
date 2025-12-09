/**
 * NEXO - Sistema de Gestão para Barbearias
 * UnitGuard Component
 *
 * Componente de guarda que garante que uma unidade esteja selecionada
 * antes de renderizar conteúdo que depende de contexto de unidade.
 */

'use client';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { useUnit } from '@/hooks/use-units';
import { cn } from '@/lib/utils';
import type { UserUnit } from '@/types/unit';
import { AlertCircle, Building2, RefreshCw } from 'lucide-react';
import { type ReactNode } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface UnitGuardProps {
  /** Conteúdo a ser renderizado quando há unidade selecionada */
  children: ReactNode;
  /** Conteúdo a ser renderizado durante carregamento */
  loadingFallback?: ReactNode;
  /** Conteúdo a ser renderizado quando não há unidade */
  noUnitFallback?: ReactNode;
  /** Se deve mostrar seletor de unidades quando não há nenhuma selecionada */
  showUnitSelector?: boolean;
  /** Classe CSS adicional */
  className?: string;
}

// =============================================================================
// COMPONENTES DE FALLBACK
// =============================================================================

function LoadingFallback() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[200px] gap-4">
      <Skeleton className="h-8 w-8 rounded-full" />
      <Skeleton className="h-4 w-32" />
      <Skeleton className="h-3 w-24" />
    </div>
  );
}

interface NoUnitFallbackProps {
  units: UserUnit[];
  onSelectUnit: (unitId: string) => void;
  onRefresh: () => void;
  isLoading: boolean;
}

function NoUnitFallback({
  units,
  onSelectUnit,
  onRefresh,
  isLoading,
}: NoUnitFallbackProps) {
  if (units.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[300px] gap-4 p-6">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Sem acesso a unidades</AlertTitle>
          <AlertDescription>
            Você não tem acesso a nenhuma unidade. Entre em contato com o
            administrador para solicitar acesso.
          </AlertDescription>
        </Alert>
        <Button variant="outline" onClick={onRefresh} disabled={isLoading}>
          <RefreshCw className={cn('size-4 mr-2', isLoading && 'animate-spin')} />
          Tentar novamente
        </Button>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-[300px] gap-6 p-6">
      <div className="text-center">
        <Building2 className="size-12 text-muted-foreground mx-auto mb-4" />
        <h2 className="text-xl font-semibold mb-2">Selecione uma unidade</h2>
        <p className="text-sm text-muted-foreground max-w-md">
          Para continuar, selecione a unidade onde deseja trabalhar.
        </p>
      </div>

      <div className="grid gap-3 w-full max-w-md">
        {units.map((unit) => (
          <Card
            key={unit.unit_id}
            className="cursor-pointer hover:border-primary transition-colors"
            onClick={() => onSelectUnit(unit.unit_id)}
          >
            <CardHeader className="py-3">
              <CardTitle className="text-base flex items-center gap-2">
                <Building2 className="size-4" />
                {unit.unit_apelido || unit.unit_nome}
                {unit.unit_matriz && (
                  <span className="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded">
                    Matriz
                  </span>
                )}
              </CardTitle>
            </CardHeader>
          </Card>
        ))}
      </div>
    </div>
  );
}

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export function UnitGuard({
  children,
  loadingFallback,
  noUnitFallback,
  showUnitSelector = true,
  className,
}: UnitGuardProps) {
  const {
    units,
    activeUnit,
    isLoading,
    isHydrated,
    switchUnit,
    refreshUnits,
  } = useUnit();

  // Ainda não hidratou do localStorage
  if (!isHydrated) {
    return (
      <div className={className}>
        {loadingFallback || <LoadingFallback />}
      </div>
    );
  }

  // Carregando unidades do servidor
  if (isLoading && units.length === 0) {
    return (
      <div className={className}>
        {loadingFallback || <LoadingFallback />}
      </div>
    );
  }

  // Não tem unidade selecionada
  if (!activeUnit) {
    if (noUnitFallback) {
      return <div className={className}>{noUnitFallback}</div>;
    }

    if (showUnitSelector) {
      return (
        <div className={className}>
          <NoUnitFallback
            units={units}
            onSelectUnit={switchUnit}
            onRefresh={refreshUnits}
            isLoading={isLoading}
          />
        </div>
      );
    }

    return null;
  }

  // Tudo OK, renderiza children
  return <div className={className}>{children}</div>;
}

// =============================================================================
// HOOK UTILITÁRIO
// =============================================================================

/**
 * Hook para verificar se há unidade ativa
 * Útil para renderização condicional simples
 */
export function useRequireUnit(): {
  hasUnit: boolean;
  unitId: string | null;
  isLoading: boolean;
} {
  const { activeUnit, isLoading, isHydrated } = useUnit();

  return {
    hasUnit: !!activeUnit && isHydrated,
    unitId: activeUnit?.unit_id ?? null,
    isLoading: isLoading || !isHydrated,
  };
}

export default UnitGuard;
