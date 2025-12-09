/**
 * NEXO - Sistema de Gestão para Barbearias
 * Saldo Cards Component
 *
 * Exibe os cards com totais do caixa (vendas, sangrias, reforços, despesas).
 *
 * @author NEXO v2.0
 */

'use client';

import {
    ArrowDownCircle,
    ArrowUpCircle,
    MinusCircle,
    PlusCircle,
    Wallet,
} from 'lucide-react';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

import { cn } from '@/lib/utils';
import type { CaixaDiarioResponse, TotaisCaixaResponse } from '@/types/caixa';

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

// =============================================================================
// TYPES
// =============================================================================

interface SaldoCardsProps {
  caixa: CaixaDiarioResponse | undefined;
  totais: TotaisCaixaResponse | undefined;
  isLoading: boolean;
}

// =============================================================================
// INDIVIDUAL CARD COMPONENT
// =============================================================================

interface SaldoCardProps {
  title: string;
  value: string | undefined;
  icon: React.ReactNode;
  colorClass: string;
  isLoading: boolean;
}

function SaldoCard({ title, value, icon, colorClass, isLoading }: SaldoCardProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-5 w-5 rounded-full" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-7 w-28" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <div className={colorClass}>{icon}</div>
      </CardHeader>
      <CardContent>
        <div className={cn('text-2xl font-bold', colorClass)}>
          {formatCurrency(value)}
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// MAIN COMPONENT
// =============================================================================

export function SaldoCards({ caixa, totais, isLoading }: SaldoCardsProps) {
  // Calcula o saldo atual do caixa
  const saldoAtual = totais?.saldo_atual || caixa?.saldo_esperado;

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
      {/* Saldo Atual */}
      <SaldoCard
        title="Saldo Atual"
        value={saldoAtual}
        icon={<Wallet className="h-5 w-5" />}
        colorClass="text-primary"
        isLoading={isLoading}
      />

      {/* Total Vendas/Entradas */}
      <SaldoCard
        title="Total Vendas"
        value={totais?.total_vendas || caixa?.total_entradas}
        icon={<ArrowUpCircle className="h-5 w-5" />}
        colorClass="text-green-600"
        isLoading={isLoading}
      />

      {/* Total Sangrias */}
      <SaldoCard
        title="Total Sangrias"
        value={totais?.total_sangrias || caixa?.total_sangrias}
        icon={<MinusCircle className="h-5 w-5" />}
        colorClass="text-red-600"
        isLoading={isLoading}
      />

      {/* Total Reforços */}
      <SaldoCard
        title="Total Reforços"
        value={totais?.total_reforcos || caixa?.total_reforcos}
        icon={<PlusCircle className="h-5 w-5" />}
        colorClass="text-blue-600"
        isLoading={isLoading}
      />

      {/* Total Despesas */}
      <SaldoCard
        title="Total Despesas"
        value={totais?.total_despesas || caixa?.total_saidas}
        icon={<ArrowDownCircle className="h-5 w-5" />}
        colorClass="text-orange-600"
        isLoading={isLoading}
      />
    </div>
  );
}

export default SaldoCards;
