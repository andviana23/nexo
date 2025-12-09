/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Fluxo de Caixa
 *
 * Exibe o fluxo de caixa diário/semanal com gráfico e tabela.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import {
    ArrowDownCircle,
    ArrowUpCircle,
    Calendar,
    ChevronLeft,
    ChevronRight,
    TrendingDown,
    TrendingUp,
    Wallet,
} from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';

import { useCashFlow } from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import { formatCurrency, formatDate } from '@/types/financial';

// Interface local que representa exatamente o que usamos na UI
interface FluxoCaixaItem {
  data: string;
  saldo_inicial: string;
  total_entradas: string;
  total_saidas: string;
  saldo_final: string;
  saldo_acumulado: string;
}

// =============================================================================
// TIPOS
// =============================================================================

type ViewMode = 'week' | 'month';

// =============================================================================
// UTILITÁRIOS
// =============================================================================

function getWeekDates(date: Date): { start: Date; end: Date } {
  const start = new Date(date);
  const day = start.getDay();
  const diff = start.getDate() - day + (day === 0 ? -6 : 1);
  start.setDate(diff);
  start.setHours(0, 0, 0, 0);

  const end = new Date(start);
  end.setDate(end.getDate() + 6);
  end.setHours(23, 59, 59, 999);

  return { start, end };
}

function getMonthDates(date: Date): { start: Date; end: Date } {
  const start = new Date(date.getFullYear(), date.getMonth(), 1);
  const end = new Date(date.getFullYear(), date.getMonth() + 1, 0);
  return { start, end };
}

function formatDateRange(start: Date, end: Date): string {
  const startStr = start.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' });
  const endStr = end.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short', year: 'numeric' });
  return `${startStr} - ${endStr}`;
}

function formatMonthYear(date: Date): string {
  return date.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' });
}

// =============================================================================
// COMPONENTE: Gráfico Simples de Barras
// =============================================================================

interface SimpleBarChartProps {
  data: FluxoCaixaItem[];
}

function SimpleBarChart({ data }: SimpleBarChartProps) {
  if (data.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center text-muted-foreground">
        Sem dados para o período selecionado
      </div>
    );
  }

  const maxValue = Math.max(
    ...data.map((d) => Math.max(parseFloat(d.total_entradas), parseFloat(d.total_saidas)))
  );

  return (
    <div className="h-64 flex items-end gap-1 px-2">
      {data.map((item, index) => {
        const entradas = parseFloat(item.total_entradas);
        const saidas = parseFloat(item.total_saidas);
        const saldo = parseFloat(item.saldo_final);
        const entradasHeight = maxValue > 0 ? (entradas / maxValue) * 100 : 0;
        const saidasHeight = maxValue > 0 ? (saidas / maxValue) * 100 : 0;
        const dayLabel = new Date(item.data).toLocaleDateString('pt-BR', { weekday: 'short', day: '2-digit' });

        return (
          <div key={index} className="flex-1 flex flex-col items-center gap-1">
            <div className="w-full flex items-end gap-0.5 h-48">
              {/* Barra de Entradas */}
              <div
                className="flex-1 bg-green-500 rounded-t transition-all hover:bg-green-600"
                style={{ height: `${entradasHeight}%`, minHeight: entradas > 0 ? '4px' : '0' }}
                title={`Entradas: ${formatCurrency(entradas)}`}
              />
              {/* Barra de Saídas */}
              <div
                className="flex-1 bg-red-500 rounded-t transition-all hover:bg-red-600"
                style={{ height: `${saidasHeight}%`, minHeight: saidas > 0 ? '4px' : '0' }}
                title={`Saídas: ${formatCurrency(saidas)}`}
              />
            </div>
            {/* Saldo */}
            <div className={`text-xs font-medium ${saldo >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {saldo >= 0 ? '+' : ''}{formatCurrency(saldo).replace('R$', '').trim()}
            </div>
            {/* Label do Dia */}
            <div className="text-xs text-muted-foreground truncate w-full text-center">
              {dayLabel}
            </div>
          </div>
        );
      })}
    </div>
  );
}

// =============================================================================
// COMPONENTE: Cards de Resumo
// =============================================================================

interface SummaryCardsProps {
  data: FluxoCaixaItem[];
}

function SummaryCards({ data }: SummaryCardsProps) {
  const totals = useMemo(() => {
    return data.reduce(
      (acc, item) => ({
        entradas: acc.entradas + parseFloat(item.total_entradas),
        saidas: acc.saidas + parseFloat(item.total_saidas),
        saldo: acc.saldo + (parseFloat(item.total_entradas) - parseFloat(item.total_saidas)),
      }),
      { entradas: 0, saidas: 0, saldo: 0 }
    );
  }, [data]);

  const saldoAcumulado = data.length > 0 ? parseFloat(data[data.length - 1].saldo_acumulado) : 0;

  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Entradas</CardTitle>
          <ArrowUpCircle className="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            {formatCurrency(totals.entradas)}
          </div>
          <p className="text-xs text-muted-foreground">
            {data.filter((d) => parseFloat(d.total_entradas) > 0).length} dias com entrada
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Saídas</CardTitle>
          <ArrowDownCircle className="h-4 w-4 text-red-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-red-600">
            {formatCurrency(totals.saidas)}
          </div>
          <p className="text-xs text-muted-foreground">
            {data.filter((d) => parseFloat(d.total_saidas) > 0).length} dias com saída
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Saldo do Período</CardTitle>
          {totals.saldo >= 0 ? (
            <TrendingUp className="h-4 w-4 text-green-500" />
          ) : (
            <TrendingDown className="h-4 w-4 text-red-500" />
          )}
        </CardHeader>
        <CardContent>
          <div className={`text-2xl font-bold ${totals.saldo >= 0 ? 'text-green-600' : 'text-red-600'}`}>
            {formatCurrency(totals.saldo)}
          </div>
          <p className="text-xs text-muted-foreground">
            Entradas - Saídas
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Saldo Acumulado</CardTitle>
          <Wallet className="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div className={`text-2xl font-bold ${saldoAcumulado >= 0 ? 'text-blue-600' : 'text-red-600'}`}>
            {formatCurrency(saldoAcumulado)}
          </div>
          <p className="text-xs text-muted-foreground">
            Posição final do período
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL: Fluxo de Caixa
// =============================================================================

export default function FluxoCaixaPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [viewMode, setViewMode] = useState<ViewMode>('week');
  const [currentDate, setCurrentDate] = useState(new Date());

  // Calcula período baseado no modo de visualização
  const period = useMemo(() => {
    if (viewMode === 'week') {
      return getWeekDates(currentDate);
    }
    return getMonthDates(currentDate);
  }, [viewMode, currentDate]);

  // Query
  const { data, isLoading } = useCashFlow({
    data_inicio: period.start.toISOString().slice(0, 10),
    data_fim: period.end.toISOString().slice(0, 10),
  });

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'Fluxo de Caixa' },
    ]);
  }, [setBreadcrumbs]);

  // Navegação
  const handlePrevious = () => {
    const newDate = new Date(currentDate);
    if (viewMode === 'week') {
      newDate.setDate(newDate.getDate() - 7);
    } else {
      newDate.setMonth(newDate.getMonth() - 1);
    }
    setCurrentDate(newDate);
  };

  const handleNext = () => {
    const newDate = new Date(currentDate);
    const now = new Date();

    if (viewMode === 'week') {
      newDate.setDate(newDate.getDate() + 7);
    } else {
      newDate.setMonth(newDate.getMonth() + 1);
    }

    // Não avança além do período atual
    if (newDate <= now) {
      setCurrentDate(newDate);
    }
  };

  const handleToday = () => {
    setCurrentDate(new Date());
  };

  const isCurrentPeriod = useMemo(() => {
    const now = new Date();
    return now >= period.start && now <= period.end;
  }, [period]);

  const fluxoData: FluxoCaixaItem[] = data ?? [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Fluxo de Caixa</h1>
          <p className="text-muted-foreground">
            Acompanhe as entradas e saídas diárias
          </p>
        </div>

        <div className="flex items-center gap-4">
          {/* Seletor de Modo */}
          <Select value={viewMode} onValueChange={(v) => setViewMode(v as ViewMode)}>
            <SelectTrigger className="w-[140px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="week">Semanal</SelectItem>
              <SelectItem value="month">Mensal</SelectItem>
            </SelectContent>
          </Select>

          {/* Navegação de Período */}
          <div className="flex items-center gap-2">
            <Button variant="outline" size="icon" onClick={handlePrevious}>
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <div className="flex items-center gap-2 px-4 py-2 border rounded-md min-w-[220px] justify-center">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="font-medium text-sm">
                {viewMode === 'week'
                  ? formatDateRange(period.start, period.end)
                  : formatMonthYear(currentDate)}
              </span>
            </div>
            <Button
              variant="outline"
              size="icon"
              onClick={handleNext}
              disabled={isCurrentPeriod}
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>

          {!isCurrentPeriod && (
            <Button variant="outline" onClick={handleToday}>
              Hoje
            </Button>
          )}
        </div>
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="space-y-4">
          <div className="grid gap-4 md:grid-cols-4">
            {[1, 2, 3, 4].map((i) => (
              <Card key={i}>
                <CardContent className="p-6">
                  <Skeleton className="h-4 w-24 mb-2" />
                  <Skeleton className="h-8 w-32" />
                </CardContent>
              </Card>
            ))}
          </div>
          <Card>
            <CardContent className="p-6">
              <Skeleton className="h-64 w-full" />
            </CardContent>
          </Card>
        </div>
      )}

      {/* Conteúdo */}
      {!isLoading && (
        <>
          {/* Cards de Resumo */}
          <SummaryCards data={fluxoData} />

          {/* Gráfico */}
          <Card>
            <CardHeader>
              <CardTitle>Movimentação por Dia</CardTitle>
              <CardDescription>
                <div className="flex items-center gap-4 mt-2">
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 bg-green-500 rounded" />
                    <span className="text-sm">Entradas</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 bg-red-500 rounded" />
                    <span className="text-sm">Saídas</span>
                  </div>
                </div>
              </CardDescription>
            </CardHeader>
            <CardContent>
              <SimpleBarChart data={fluxoData} />
            </CardContent>
          </Card>

          {/* Tabela Detalhada */}
          <Card>
            <CardHeader>
              <CardTitle>Detalhamento Diário</CardTitle>
              <CardDescription>
                Todas as movimentações do período
              </CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              {fluxoData.length > 0 ? (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Data</TableHead>
                      <TableHead className="text-right">Entradas</TableHead>
                      <TableHead className="text-right">Saídas</TableHead>
                      <TableHead className="text-right">Saldo do Dia</TableHead>
                      <TableHead className="text-right">Saldo Acumulado</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                      {fluxoData.map((item, index) => {
                      const entradas = parseFloat(item.total_entradas);
                      const saidas = parseFloat(item.total_saidas);
                      const saldoDia = entradas - saidas;
                      const saldoAcumulado = parseFloat(item.saldo_acumulado);
                      const isWeekend = [0, 6].includes(new Date(item.data).getDay());

                      return (
                        <TableRow key={index} className={isWeekend ? 'bg-muted/30' : ''}>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              <span className="font-medium">{formatDate(item.data)}</span>
                              <span className="text-xs text-muted-foreground">
                                {new Date(item.data).toLocaleDateString('pt-BR', { weekday: 'short' })}
                              </span>
                            </div>
                          </TableCell>
                          <TableCell className="text-right">
                            {entradas > 0 ? (
                              <span className="text-green-600 font-medium">
                                +{formatCurrency(entradas)}
                              </span>
                            ) : (
                              <span className="text-muted-foreground">-</span>
                            )}
                          </TableCell>
                          <TableCell className="text-right">
                            {saidas > 0 ? (
                              <span className="text-red-600 font-medium">
                                -{formatCurrency(saidas)}
                              </span>
                            ) : (
                              <span className="text-muted-foreground">-</span>
                            )}
                          </TableCell>
                          <TableCell className="text-right">
                            <Badge
                              variant={saldoDia >= 0 ? 'default' : 'destructive'}
                              className="font-mono"
                            >
                              {saldoDia >= 0 ? '+' : ''}{formatCurrency(saldoDia)}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-right">
                            <span className={`font-semibold ${saldoAcumulado >= 0 ? 'text-blue-600' : 'text-red-600'}`}>
                              {formatCurrency(saldoAcumulado)}
                            </span>
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              ) : (
                <div className="p-8 text-center text-muted-foreground">
                  <Wallet className="h-12 w-12 mx-auto mb-3 opacity-50" />
                  <p>Nenhuma movimentação no período selecionado</p>
                  <p className="text-sm mt-1">
                    Cadastre receitas e despesas para visualizar o fluxo de caixa
                  </p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Legenda/Info */}
          <Card>
            <CardContent className="py-4">
              <div className="flex flex-wrap gap-6 text-sm">
                <div className="flex items-center gap-2">
                  <ArrowUpCircle className="h-4 w-4 text-green-500" />
                  <span className="text-muted-foreground">Entradas: receitas e recebimentos</span>
                </div>
                <div className="flex items-center gap-2">
                  <ArrowDownCircle className="h-4 w-4 text-red-500" />
                  <span className="text-muted-foreground">Saídas: despesas e pagamentos</span>
                </div>
                <div className="flex items-center gap-2">
                  <Wallet className="h-4 w-4 text-blue-500" />
                  <span className="text-muted-foreground">Saldo acumulado: posição de caixa</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}
