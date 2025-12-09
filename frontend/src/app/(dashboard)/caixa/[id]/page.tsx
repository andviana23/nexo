/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Detalhes Page
 *
 * Página de detalhes de um caixa específico.
 * Segue os padrões do Design System v1.0.
 *
 * @author NEXO v2.0
 */

'use client';

import {
    AlertTriangle,
    ArrowLeft,
    BadgeCheck,
    Clock,
    RefreshCw,
    User,
    Wallet
} from 'lucide-react';
import Link from 'next/link';
import { use } from 'react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

import { ExtratoDia } from '@/components/caixa';
import { useCaixaById } from '@/hooks/use-caixa';
import { cn } from '@/lib/utils';
import { StatusCaixa, StatusCaixaColors, StatusCaixaLabels } from '@/types/caixa';

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

const formatDateOnly = (dateStr: string | undefined) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date);
};

const formatTimeOnly = (dateStr: string | undefined) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

// =============================================================================
// TYPES
// =============================================================================

interface CaixaDetalhesPageProps {
  params: Promise<{ id: string }>;
}

// =============================================================================
// PAGE COMPONENT
// =============================================================================

export default function CaixaDetalhesPage({ params }: CaixaDetalhesPageProps) {
  const { id } = use(params);
  const { data: caixa, isLoading, isError, error, refetch } = useCaixaById(id);

  // Loading state
  if (isLoading) {
    return (
      <div className="flex-1 space-y-4 p-4 md:p-6 lg:p-8">
        <div className="flex items-center gap-3">
          <Skeleton className="h-10 w-10" />
          <div className="space-y-2">
            <Skeleton className="h-5 w-40" />
            <Skeleton className="h-4 w-24" />
          </div>
        </div>
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-48" />
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <Skeleton className="h-20" />
              <Skeleton className="h-20" />
            </div>
            <Skeleton className="h-32" />
            <Skeleton className="h-24" />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <Skeleton className="h-5 w-32" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-48" />
          </CardContent>
        </Card>
      </div>
    );
  }

  // Error state
  if (isError) {
    return (
      <div className="flex-1 p-4 md:p-6 lg:p-8">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16 text-center">
            <div className="rounded-full bg-destructive/10 p-3 mb-4">
              <AlertTriangle className="h-6 w-6 text-destructive" />
            </div>
            <h2 className="text-base font-semibold mb-2">Erro ao carregar caixa</h2>
            <p className="text-sm text-muted-foreground mb-6 max-w-md">
              {error?.message || 'Ocorreu um erro inesperado.'}
            </p>
            <div className="flex gap-3">
              <Button variant="outline" size="sm" asChild>
                <Link href="/caixa/historico">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Voltar
                </Link>
              </Button>
              <Button size="sm" onClick={() => refetch()}>
                <RefreshCw className="mr-2 h-4 w-4" />
                Tentar novamente
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Not found state
  if (!caixa) {
    return (
      <div className="flex-1 p-4 md:p-6 lg:p-8">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16 text-center">
            <div className="rounded-full bg-muted p-3 mb-4">
              <Wallet className="h-6 w-6 text-muted-foreground" />
            </div>
            <h2 className="text-base font-semibold mb-2">Caixa não encontrado</h2>
            <p className="text-sm text-muted-foreground mb-6">
              O caixa solicitado não existe ou foi removido.
            </p>
            <Button variant="outline" size="sm" asChild>
              <Link href="/caixa/historico">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Voltar
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const temDivergencia =
    caixa.divergencia && parseFloat(caixa.divergencia) !== 0;

  return (
    <div className="flex-1 space-y-4 p-4 md:p-6 lg:p-8">
      {/* Header Clean */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/caixa/historico">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-xl font-semibold tracking-tight">Detalhes do Caixa</h1>
            <p className="text-sm text-muted-foreground">
              {formatDateOnly(caixa.data_abertura)}
            </p>
          </div>
        </div>
        
        <Badge className={cn(StatusCaixaColors[caixa.status as StatusCaixa])}>
          {StatusCaixaLabels[caixa.status as StatusCaixa]}
        </Badge>
      </div>

      {/* Main Card - Informações Principais */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Informações do Caixa</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Horários e Responsáveis */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Abertura</p>
              <div className="flex items-center gap-2">
                <Clock className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">{formatTimeOnly(caixa.data_abertura)}</span>
              </div>
              <div className="flex items-center gap-2">
                <User className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">{caixa.usuario_abertura_nome}</span>
              </div>
            </div>

            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Fechamento</p>
              {caixa.data_fechamento ? (
                <>
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm">{formatTimeOnly(caixa.data_fechamento)}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm">{caixa.usuario_fechamento_nome || '-'}</span>
                  </div>
                </>
              ) : (
                <div className="flex items-center gap-2">
                  <div className="h-2 w-2 rounded-full bg-chart-2 animate-pulse" />
                  <span className="text-sm text-muted-foreground">Caixa em andamento</span>
                </div>
              )}
            </div>
          </div>

          {/* Valores */}
          <div className="space-y-4 border-t pt-4">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div className="space-y-1">
                <p className="text-xs text-muted-foreground">Saldo Inicial</p>
                <p className="text-lg font-semibold">{formatCurrency(caixa.saldo_inicial)}</p>
              </div>
              
              <div className="space-y-1">
                <p className="text-xs text-muted-foreground">Total Entradas</p>
                <p className="text-lg font-semibold text-chart-2">
                  {formatCurrency(caixa.total_entradas)}
                </p>
              </div>
              
              <div className="space-y-1">
                <p className="text-xs text-muted-foreground">Total Saídas</p>
                <p className="text-lg font-semibold text-destructive">
                  {formatCurrency(caixa.total_saidas)}
                </p>
              </div>
              
              <div className="space-y-1">
                <p className="text-xs text-muted-foreground">Saldo Esperado</p>
                <p className="text-lg font-semibold">{formatCurrency(caixa.saldo_esperado)}</p>
              </div>
            </div>
          </div>

          {/* Conferência */}
          <div className={cn(
            "space-y-2 rounded-lg border p-4",
            temDivergencia ? "border-chart-5/30 bg-chart-5/5" : "border-chart-2/30 bg-chart-2/5"
          )}>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                {temDivergencia ? (
                  <AlertTriangle className="h-4 w-4 text-chart-5" />
                ) : (
                  <BadgeCheck className="h-4 w-4 text-chart-2" />
                )}
                <p className="text-sm font-medium">
                  {temDivergencia ? 'Divergência Detectada' : 'Caixa Conferido'}
                </p>
              </div>
              {temDivergencia && (
                <span className="text-lg font-bold text-chart-5">
                  {formatCurrency(caixa.divergencia)}
                </span>
              )}
            </div>
            
            {caixa.saldo_real && (
              <p className="text-sm text-muted-foreground">
                Saldo real informado: <span className="font-medium text-foreground">{formatCurrency(caixa.saldo_real)}</span>
              </p>
            )}

            {temDivergencia && caixa.justificativa_divergencia && (
              <p className="text-sm text-muted-foreground pt-2 border-t">
                <span className="font-medium">Justificativa:</span> {caixa.justificativa_divergencia}
              </p>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Extrato de Operações */}
      <ExtratoDia operacoes={caixa.operacoes} isLoading={false} />
    </div>
  );
}
