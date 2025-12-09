/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Status Card Component
 *
 * Exibe o status atual do caixa (aberto/fechado) com informações resumidas.
 *
 * @author NEXO v2.0
 */

'use client';

import { CircleDollarSign, Clock, Lock, LockOpen } from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

import { cn } from '@/lib/utils';
import type { CaixaStatusResponse } from '@/types/caixa';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: string | undefined) => {
  if (!value) return 'R$ 0,00';
  const num = parseFloat(value);
  if (isNaN(num)) return 'R$ 0,00';
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
};

const formatDateTime = (dateStr: string | undefined) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

// =============================================================================
// TYPES
// =============================================================================

interface CaixaStatusCardProps {
  status: CaixaStatusResponse | undefined;
  isLoading: boolean;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function CaixaStatusCard({ status, isLoading }: CaixaStatusCardProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-6 w-6 rounded-full" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-48 mb-2" />
          <Skeleton className="h-4 w-64" />
        </CardContent>
      </Card>
    );
  }

  const isAberto = status?.aberto ?? false;
  const caixa = status?.caixa_atual;

  return (
    <Card className={cn(
      'border-l-4',
      isAberto ? 'border-l-green-500' : 'border-l-gray-400'
    )}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">Status do Caixa</CardTitle>
        {isAberto ? (
          <LockOpen className="h-5 w-5 text-green-500" />
        ) : (
          <Lock className="h-5 w-5 text-gray-400" />
        )}
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-3 mb-3">
          <Badge
            variant={isAberto ? 'default' : 'secondary'}
            className={cn(
              'text-sm px-3 py-1',
              isAberto 
                ? 'bg-green-100 text-green-800 hover:bg-green-200' 
                : 'bg-gray-100 text-gray-600'
            )}
          >
            {isAberto ? 'ABERTO' : 'FECHADO'}
          </Badge>
        </div>

        {isAberto && caixa && (
          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4" />
              <span>Aberto em: {formatDateTime(caixa.data_abertura)}</span>
            </div>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <span>Por: {caixa.usuario_abertura_nome}</span>
            </div>
            <div className="flex items-center gap-2 text-sm font-medium mt-3">
              <CircleDollarSign className="h-4 w-4 text-green-600" />
              <span>Saldo Inicial: {formatCurrency(caixa.saldo_inicial)}</span>
            </div>
          </div>
        )}

        {!isAberto && status?.ultimo_fechamento && (
          <div className="text-sm text-muted-foreground">
            <p>Último fechamento: {formatDateTime(status.ultimo_fechamento)}</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

export default CaixaStatusCard;
