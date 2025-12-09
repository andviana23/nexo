/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Diário Page (Refatorado)
 *
 * Página principal do módulo Caixa Diário.
 * Segue 100% o Design System NEXO v1.0
 *
 * @author NEXO v2.0
 */

'use client';

import {
    AlertCircle,
    ArrowDownCircle,
    ArrowUpCircle,
    Banknote,
    Calendar,
    Clock,
    History,
    Lock,
    LockOpen,
    MinusCircle,
    PlusCircle,
    RefreshCw,
    TrendingDown,
    TrendingUp,
    User,
    Wallet,
} from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from '@/components/ui/tooltip';

import {
    ModalAbrirCaixa,
    ModalFecharCaixa,
    ModalReforco,
    ModalSangria,
} from '@/components/caixa';

import { useCaixaDiario } from '@/hooks/use-caixa';
import { cn } from '@/lib/utils';
import type { OperacaoCaixaResponse } from '@/types/caixa';
import {
    DestinoSangria,
    DestinoSangriaLabels,
    OrigemReforco,
    OrigemReforcoLabels,
    TipoOperacaoCaixa,
} from '@/types/caixa';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: string | number | undefined): string => {
  if (value === undefined || value === null) return 'R$ 0,00';
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return 'R$ 0,00';
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
};

const formatDateTime = (dateStr: string | undefined): string => {
  if (!dateStr) return '-';
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateStr));
};

const formatTime = (dateStr: string | undefined): string => {
  if (!dateStr) return '-';
  return new Intl.DateTimeFormat('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateStr));
};

const formatDate = (dateStr: string | undefined): string => {
  if (!dateStr) return '-';
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  }).format(new Date(dateStr));
};

// =============================================================================
// LOADING SKELETON
// =============================================================================

function PageSkeleton() {
  return (
    <div className="flex-1 space-y-6 p-4 md:p-6">
      {/* Header Skeleton */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-2">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-4 w-64" />
        </div>
        <div className="flex gap-2">
          <Skeleton className="h-9 w-24" />
          <Skeleton className="h-9 w-28" />
          <Skeleton className="h-9 w-32" />
        </div>
      </div>

      {/* Status Card Skeleton */}
      <Skeleton className="h-40 w-full" />

      {/* Metrics Grid Skeleton */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
        {[1, 2, 3, 4, 5].map((i) => (
          <Skeleton key={i} className="h-28" />
        ))}
      </div>

      {/* Table Skeleton */}
      <Skeleton className="h-80 w-full" />
    </div>
  );
}

// =============================================================================
// ERROR STATE
// =============================================================================

interface ErrorStateProps {
  message: string;
  onRetry: () => void;
}

function ErrorState({ message, onRetry }: ErrorStateProps) {
  return (
    <div className="flex-1 p-4 md:p-6">
      <Card className="border-destructive/50 bg-destructive/5">
        <CardContent className="flex flex-col items-center justify-center py-12">
          <AlertCircle className="h-12 w-12 text-destructive mb-4" />
          <h2 className="text-lg font-semibold text-destructive mb-2">
            Erro ao carregar caixa
          </h2>
          <p className="text-muted-foreground text-center mb-6 max-w-md">
            {message || 'Ocorreu um erro inesperado ao carregar os dados do caixa.'}
          </p>
          <Button onClick={onRetry} variant="outline">
            <RefreshCw className="mr-2 h-4 w-4" />
            Tentar novamente
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// CLOSED STATE
// =============================================================================

interface ClosedStateProps {
  ultimoFechamento?: string;
  onAbrir: () => void;
}

function ClosedState({ ultimoFechamento, onAbrir }: ClosedStateProps) {
  return (
    <Card className="border-dashed border-2 bg-muted/30">
      <CardContent className="flex flex-col items-center justify-center py-16">
        <div className="rounded-full bg-muted p-4 mb-6">
          <Lock className="h-10 w-10 text-muted-foreground" />
        </div>
        <h3 className="text-xl font-semibold mb-2">Caixa Fechado</h3>
        <p className="text-muted-foreground text-center mb-2 max-w-md">
          Abra o caixa para começar a registrar as operações do dia.
        </p>
        {ultimoFechamento && (
          <p className="text-sm text-muted-foreground mb-6">
            Último fechamento: {formatDateTime(ultimoFechamento)}
          </p>
        )}
        <Button size="lg" onClick={onAbrir} className="gap-2">
          <LockOpen className="h-5 w-5" />
          Abrir Caixa
        </Button>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// STATUS BANNER
// =============================================================================

interface StatusBannerProps {
  isAberto: boolean;
  dataAbertura?: string;
  usuarioNome?: string;
  saldoInicial?: string;
}

function StatusBanner({ isAberto, dataAbertura, usuarioNome, saldoInicial }: StatusBannerProps) {
  if (!isAberto) return null;

  return (
    <Card className="border-l-4 border-l-green-500 bg-green-50/50 dark:bg-green-950/20">
      <CardContent className="py-4">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-4">
            <div className="rounded-full bg-green-100 dark:bg-green-900/50 p-2">
              <LockOpen className="h-5 w-5 text-green-600 dark:text-green-400" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <Badge className="bg-green-100 text-green-800 hover:bg-green-200 dark:bg-green-900 dark:text-green-200">
                  Caixa Aberto
                </Badge>
                <span className="text-sm text-muted-foreground">
                  desde {formatTime(dataAbertura)}
                </span>
              </div>
              <p className="text-sm text-muted-foreground mt-1">
                Operador: <span className="font-medium text-foreground">{usuarioNome || '-'}</span>
              </p>
            </div>
          </div>
          <div className="flex items-center gap-6 text-sm">
            <div className="flex items-center gap-2 text-muted-foreground">
              <Calendar className="h-4 w-4" />
              <span>{formatDate(dataAbertura)}</span>
            </div>
            <div className="flex items-center gap-2">
              <Banknote className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">Saldo inicial:</span>
              <span className="font-semibold">{formatCurrency(saldoInicial)}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// METRIC CARD
// =============================================================================

interface MetricCardProps {
  title: string;
  value: string | number | undefined;
  icon: React.ReactNode;
  trend?: 'up' | 'down' | 'neutral';
  description?: string;
  className?: string;
}

function MetricCard({ title, value, icon, trend, description, className }: MetricCardProps) {
  const trendColor = trend === 'up' 
    ? 'text-green-600 dark:text-green-400' 
    : trend === 'down' 
      ? 'text-red-600 dark:text-red-400' 
      : 'text-muted-foreground';

  return (
    <Card className={cn('transition-all hover:shadow-md', className)}>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <div className={cn('rounded-md p-1.5', trendColor)}>
          {icon}
        </div>
      </CardHeader>
      <CardContent>
        <div className={cn('text-2xl font-bold tracking-tight', trendColor)}>
          {formatCurrency(value)}
        </div>
        {description && (
          <p className="text-xs text-muted-foreground mt-1">{description}</p>
        )}
      </CardContent>
    </Card>
  );
}

// =============================================================================
// OPERATIONS TABLE
// =============================================================================

interface OperationsTableProps {
  operacoes: OperacaoCaixaResponse[] | undefined;
}

function OperationsTable({ operacoes }: OperationsTableProps) {
  if (!operacoes || operacoes.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Clock className="h-5 w-5" />
            Extrato do Dia
          </CardTitle>
          <CardDescription>Operações realizadas no caixa atual</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <div className="rounded-full bg-muted p-3 mb-4">
              <Banknote className="h-6 w-6 text-muted-foreground" />
            </div>
            <p className="text-muted-foreground">
              Nenhuma operação registrada ainda.
            </p>
            <p className="text-sm text-muted-foreground mt-1">
              As operações aparecerão aqui conforme forem registradas.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const operacoesOrdenadas = [...operacoes].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );

  const getOperacaoConfig = (tipo: TipoOperacaoCaixa) => {
    switch (tipo) {
      case TipoOperacaoCaixa.VENDA:
        return {
          icon: <ArrowUpCircle className="h-4 w-4" />,
          color: 'text-green-600 dark:text-green-400',
          bgColor: 'bg-green-100 dark:bg-green-900/30',
          label: 'Venda',
        };
      case TipoOperacaoCaixa.SANGRIA:
        return {
          icon: <TrendingDown className="h-4 w-4" />,
          color: 'text-red-600 dark:text-red-400',
          bgColor: 'bg-red-100 dark:bg-red-900/30',
          label: 'Sangria',
        };
      case TipoOperacaoCaixa.REFORCO:
        return {
          icon: <TrendingUp className="h-4 w-4" />,
          color: 'text-blue-600 dark:text-blue-400',
          bgColor: 'bg-blue-100 dark:bg-blue-900/30',
          label: 'Reforço',
        };
      case TipoOperacaoCaixa.DESPESA:
        return {
          icon: <ArrowDownCircle className="h-4 w-4" />,
          color: 'text-orange-600 dark:text-orange-400',
          bgColor: 'bg-orange-100 dark:bg-orange-900/30',
          label: 'Despesa',
        };
      default:
        return {
          icon: <Banknote className="h-4 w-4" />,
          color: 'text-muted-foreground',
          bgColor: 'bg-muted',
          label: tipo,
        };
    }
  };

  const getDestinoOrigemLabel = (op: OperacaoCaixaResponse): string => {
    if (op.destino) {
      return DestinoSangriaLabels[op.destino as DestinoSangria] || op.destino;
    }
    if (op.origem) {
      return OrigemReforcoLabels[op.origem as OrigemReforco] || op.origem;
    }
    return '';
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5" />
              Extrato do Dia
            </CardTitle>
            <CardDescription>
              {operacoesOrdenadas.length} operação(ões) registrada(s)
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow className="bg-muted/50">
                <TableHead className="w-20">Hora</TableHead>
                <TableHead className="w-32">Tipo</TableHead>
                <TableHead>Descrição</TableHead>
                <TableHead className="w-40">Responsável</TableHead>
                <TableHead className="w-32 text-right">Valor</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {operacoesOrdenadas.map((op) => {
                const config = getOperacaoConfig(op.tipo as TipoOperacaoCaixa);
                const destinoOrigem = getDestinoOrigemLabel(op);
                const isNegative = op.tipo === TipoOperacaoCaixa.SANGRIA || op.tipo === TipoOperacaoCaixa.DESPESA;

                return (
                  <TableRow key={op.id} className="hover:bg-muted/30">
                    <TableCell className="font-mono text-sm">
                      {formatTime(op.created_at)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <div className={cn('rounded-full p-1.5', config.bgColor, config.color)}>
                          {config.icon}
                        </div>
                        <span className="font-medium">{config.label}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-col">
                        <span className="text-sm">{op.descricao || '-'}</span>
                        {destinoOrigem && (
                          <span className="text-xs text-muted-foreground">
                            {op.destino ? 'Destino' : 'Origem'}: {destinoOrigem}
                          </span>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-muted-foreground" />
                        <span className="text-sm">{op.usuario_nome || '-'}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-right">
                      <span className={cn(
                        'font-semibold tabular-nums',
                        isNegative ? 'text-red-600 dark:text-red-400' : 'text-green-600 dark:text-green-400'
                      )}>
                        {isNegative ? '-' : '+'} {formatCurrency(op.valor)}
                      </span>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// MAIN PAGE COMPONENT
// =============================================================================

export default function CaixaPage() {
  // Estados dos modais
  const [modalAbrir, setModalAbrir] = useState(false);
  const [modalSangria, setModalSangria] = useState(false);
  const [modalReforco, setModalReforco] = useState(false);
  const [modalFechar, setModalFechar] = useState(false);

  // Dados do caixa
  const {
    isLoading,
    isError,
    error,
    isAberto,
    caixaAtual,
    ultimoFechamento,
    totais,
    refetch,
  } = useCaixaDiario();

  // Loading state
  if (isLoading) {
    return <PageSkeleton />;
  }

  // Error state
  if (isError) {
    return (
      <ErrorState
        message={error?.message || 'Ocorreu um erro inesperado'}
        onRetry={() => refetch()}
      />
    );
  }

  return (
    <TooltipProvider>
      <div className="flex-1 space-y-6 p-4 md:p-6">
        {/* Header */}
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div>
            <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
              <Wallet className="h-6 w-6" />
              Caixa Diário
            </h1>
            <p className="text-muted-foreground mt-1">
              Gerencie a gaveta de dinheiro da barbearia
            </p>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-wrap items-center gap-2">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => refetch()}
                  className="gap-2"
                >
                  <RefreshCw className="h-4 w-4" />
                  <span className="hidden sm:inline">Atualizar</span>
                </Button>
              </TooltipTrigger>
              <TooltipContent>Atualizar dados do caixa</TooltipContent>
            </Tooltip>

            <Link href="/caixa/historico">
              <Button variant="outline" size="sm" className="gap-2">
                <History className="h-4 w-4" />
                <span className="hidden sm:inline">Histórico</span>
              </Button>
            </Link>

            {isAberto && (
              <>
                <Separator orientation="vertical" className="h-6 hidden sm:block" />

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setModalReforco(true)}
                      className="gap-2 border-blue-200 text-blue-700 hover:bg-blue-50 hover:text-blue-800 dark:border-blue-800 dark:text-blue-400 dark:hover:bg-blue-950"
                    >
                      <PlusCircle className="h-4 w-4" />
                      <span className="hidden sm:inline">Reforço</span>
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Adicionar dinheiro ao caixa</TooltipContent>
                </Tooltip>

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setModalSangria(true)}
                      className="gap-2 border-orange-200 text-orange-700 hover:bg-orange-50 hover:text-orange-800 dark:border-orange-800 dark:text-orange-400 dark:hover:bg-orange-950"
                    >
                      <MinusCircle className="h-4 w-4" />
                      <span className="hidden sm:inline">Sangria</span>
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Retirar dinheiro do caixa</TooltipContent>
                </Tooltip>

                <Button
                  variant="destructive"
                  size="sm"
                  onClick={() => setModalFechar(true)}
                  className="gap-2"
                >
                  <Lock className="h-4 w-4" />
                  <span className="hidden sm:inline">Fechar Caixa</span>
                </Button>
              </>
            )}

            {!isAberto && (
              <Button onClick={() => setModalAbrir(true)} className="gap-2">
                <LockOpen className="h-4 w-4" />
                Abrir Caixa
              </Button>
            )}
          </div>
        </div>

        {/* Status Banner (quando aberto) */}
        {isAberto && (
          <StatusBanner
            isAberto={isAberto}
            dataAbertura={caixaAtual?.data_abertura}
            usuarioNome={caixaAtual?.usuario_abertura_nome}
            saldoInicial={caixaAtual?.saldo_inicial}
          />
        )}

        {/* Closed State */}
        {!isAberto && (
          <ClosedState
            ultimoFechamento={ultimoFechamento}
            onAbrir={() => setModalAbrir(true)}
          />
        )}

        {/* Metrics Grid (quando aberto) */}
        {isAberto && (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
            <MetricCard
              title="Saldo Atual"
              value={totais?.saldo_atual || caixaAtual?.saldo_esperado}
              icon={<Wallet className="h-4 w-4" />}
              trend="neutral"
              description="Saldo esperado em caixa"
            />
            <MetricCard
              title="Vendas"
              value={totais?.total_vendas || '0'}
              icon={<ArrowUpCircle className="h-4 w-4" />}
              trend="up"
              description="Total de vendas do dia"
            />
            <MetricCard
              title="Sangrias"
              value={caixaAtual?.total_sangrias || totais?.total_sangrias || '0'}
              icon={<TrendingDown className="h-4 w-4" />}
              trend="down"
              description="Retiradas do caixa"
            />
            <MetricCard
              title="Reforços"
              value={caixaAtual?.total_reforcos || totais?.total_reforcos || '0'}
              icon={<TrendingUp className="h-4 w-4" />}
              trend="up"
              description="Adições ao caixa"
            />
            <MetricCard
              title="Despesas"
              value={totais?.total_despesas || '0'}
              icon={<ArrowDownCircle className="h-4 w-4" />}
              trend="down"
              description="Despesas pagas"
            />
          </div>
        )}

        {/* Operations Table (quando aberto) */}
        {isAberto && (
          <OperationsTable operacoes={caixaAtual?.operacoes} />
        )}

        {/* Modais */}
        <ModalAbrirCaixa
          open={modalAbrir}
          onOpenChange={setModalAbrir}
        />

        <ModalSangria
          open={modalSangria}
          onOpenChange={setModalSangria}
        />

        <ModalReforco
          open={modalReforco}
          onOpenChange={setModalReforco}
        />

        <ModalFecharCaixa
          open={modalFechar}
          onOpenChange={setModalFechar}
          caixa={caixaAtual}
        />
      </div>
    </TooltipProvider>
  );
}
