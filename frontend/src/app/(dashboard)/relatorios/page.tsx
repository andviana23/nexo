/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Relatórios
 *
 * Visualização consolidada de DRE, Fluxo de Caixa, Faturamento e Despesas.
 * Conforme FLUXO_RELATORIOS.md
 */

'use client';

import {
    ArrowDownCircle,
    ArrowUpCircle,
    BarChart3,
    Calculator,
    Calendar,
    ChevronLeft,
    ChevronRight,
    TrendingDown,
    TrendingUp,
    Wallet
} from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';

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
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

import { useCashFlow, useDREByMonth, usePayables, useReceivables } from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import {
    ContaPagar,
    ContaReceber,
    DREMensalExtended,
    formatCurrency,
    formatDate,
    ListContasPagarFilters,
    ListContasReceberFilters,
    ListFluxoCaixaFilters,
    StatusContaPagar,
    StatusContaReceber,
} from '@/types/financial';

// =============================================================================
// TIPOS
// =============================================================================

type ReportTab = 'dre' | 'cashflow' | 'faturamento' | 'despesas';

interface FluxoCaixaItem {
  data: string;
  saldo_inicial: string;
  total_entradas: string;
  total_saidas: string;
  saldo_final: string;
  saldo_acumulado: string;
}

// =============================================================================
// UTILITÁRIOS
// =============================================================================

function getMonthName(dateString: string): string {
  const [year, month] = dateString.split('-');
  const date = new Date(parseInt(year), parseInt(month) - 1);
  return date.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' });
}

function formatPercent(value: number | string | undefined): string {
  if (value === undefined || value === null) return '0.0%';
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return `${num.toFixed(1)}%`;
}

function getPercentColor(value: number): string {
  if (value > 20) return 'text-green-600';
  if (value > 10) return 'text-blue-600';
  if (value >= 0) return 'text-yellow-600';
  return 'text-red-600';
}

function getMonthStartEnd(month: string): { start: string; end: string } {
  const [year, m] = month.split('-').map(Number);
  const start = new Date(year, m - 1, 1);
  const end = new Date(year, m, 0);
  return {
    start: start.toISOString().split('T')[0],
    end: end.toISOString().split('T')[0],
  };
}

// =============================================================================
// COMPONENTE: Linha do DRE
// =============================================================================

interface DRERowProps {
  label: string;
  value: string | number;
  percent?: string | number;
  isSubtotal?: boolean;
  isTotal?: boolean;
  isNegative?: boolean;
  indent?: number;
}

function DRERow({ label, value, percent, isSubtotal, isTotal, isNegative, indent = 0 }: DRERowProps) {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  const numPercent = percent ? (typeof percent === 'string' ? parseFloat(percent) : percent) : undefined;

  return (
    <TableRow className={isTotal ? 'bg-muted/50 font-bold' : isSubtotal ? 'font-semibold' : ''}>
      <TableCell>
        <span style={{ paddingLeft: `${indent * 16}px` }}>{label}</span>
      </TableCell>
      <TableCell className={`text-right ${isNegative ? 'text-red-600' : ''} ${isTotal ? 'text-lg' : ''}`}>
        {isNegative && numValue !== 0 ? '(' : ''}{formatCurrency(Math.abs(numValue))}{isNegative && numValue !== 0 ? ')' : ''}
      </TableCell>
      <TableCell className="text-right w-24">
        {numPercent !== undefined && (
          <span className={getPercentColor(numPercent)}>
            {formatPercent(numPercent)}
          </span>
        )}
      </TableCell>
    </TableRow>
  );
}

// =============================================================================
// COMPONENTE: Seção DRE
// =============================================================================

interface DRESectionProps {
  selectedMonth: string;
}

function DRESection({ selectedMonth }: DRESectionProps) {
  const { data: dreData, isLoading } = useDREByMonth(selectedMonth);
  const dre = dreData as DREMensalExtended | undefined;

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <Skeleton className="h-96 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (!dre) {
    return (
      <Card>
        <CardContent className="py-16 text-center">
          <Calculator className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-semibold mb-2">DRE não disponível</h3>
          <p className="text-muted-foreground max-w-md mx-auto">
            Não há dados financeiros suficientes para gerar o DRE de{' '}
            <span className="font-medium">{getMonthName(selectedMonth)}</span>.
          </p>
          <p className="text-sm text-muted-foreground mt-2">
            Cadastre receitas e despesas para visualizar o demonstrativo.
          </p>
        </CardContent>
      </Card>
    );
  }

  const receitaBruta = parseFloat(dre.receita_bruta || '0');
  const lucroOperacional = parseFloat(dre.lucro_operacional || '0');
  const margemOperacional = parseFloat(dre.margem_operacional_percent || '0');
  const lucroLiquido = parseFloat(dre.lucro_liquido || '0');
  const margemLiquida = parseFloat(dre.margem_liquida_percent || '0');
  const isLucrativo = lucroOperacional > 0;

  return (
    <div className="space-y-6">
      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Receita Bruta
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(receitaBruta)}</div>
            <div className="flex items-center text-xs text-muted-foreground mt-1">
              <TrendingUp className="h-3 w-3 text-green-500 mr-1" />
              Total de entradas
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Lucro Operacional
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${isLucrativo ? 'text-green-600' : 'text-red-600'}`}>
              {formatCurrency(lucroOperacional)}
            </div>
            <div className="flex items-center text-xs mt-1">
              <Badge variant={isLucrativo ? 'default' : 'destructive'} className="text-xs">
                Margem: {formatPercent(margemOperacional)}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Lucro Líquido
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${lucroLiquido >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {formatCurrency(lucroLiquido)}
            </div>
            <div className="flex items-center text-xs mt-1">
              <Badge variant={lucroLiquido >= 0 ? 'default' : 'destructive'} className="text-xs">
                Margem: {formatPercent(margemLiquida)}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Indicador
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              {isLucrativo ? (
                <>
                  <TrendingUp className="h-8 w-8 text-green-500" />
                  <div>
                    <p className="font-semibold text-green-600">Positivo</p>
                    <p className="text-xs text-muted-foreground">Mês lucrativo</p>
                  </div>
                </>
              ) : (
                <>
                  <TrendingDown className="h-8 w-8 text-red-500" />
                  <div>
                    <p className="font-semibold text-red-600">Negativo</p>
                    <p className="text-xs text-muted-foreground">Prejuízo no mês</p>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabela DRE */}
      <Card>
        <CardHeader>
          <CardTitle>Demonstrativo Detalhado</CardTitle>
          <CardDescription>
            Estrutura completa de receitas, custos e despesas
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Descrição</TableHead>
                <TableHead className="text-right">Valor</TableHead>
                <TableHead className="text-right w-24">%</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {/* RECEITAS */}
              <DRERow label="RECEITA BRUTA" value={dre.receita_bruta} isSubtotal />
              <DRERow label="Receita de Serviços" value={dre.receita_servicos || '0'} percent={receitaBruta > 0 ? (parseFloat(dre.receita_servicos || '0') / receitaBruta) * 100 : 0} indent={1} />
              <DRERow label="Receita de Produtos" value={dre.receita_produtos || '0'} percent={receitaBruta > 0 ? (parseFloat(dre.receita_produtos || '0') / receitaBruta) * 100 : 0} indent={1} />
              
              <DRERow label="(-) Deduções" value={dre.deducoes} isNegative indent={1} />
              
              <TableRow className="bg-muted/30">
                <TableCell colSpan={3}><Separator /></TableCell>
              </TableRow>
              
              <DRERow label="RECEITA LÍQUIDA" value={dre.receita_liquida} isSubtotal />
              
              {/* CUSTOS */}
              <DRERow label="(-) Custos dos Serviços" value={dre.custos_servicos} isNegative indent={1} />
              
              <DRERow label="LUCRO BRUTO" value={dre.lucro_bruto} percent={dre.margem_bruta_percent} isSubtotal />
              
              {/* DESPESAS */}
              <DRERow label="(-) Despesas Operacionais" value={dre.despesas_operacionais} isNegative indent={1} />
              <DRERow label="(-) Despesas Administrativas" value={dre.despesas_administrativas} isNegative indent={1} />
              
              <TableRow className="bg-muted/30">
                <TableCell colSpan={3}><Separator /></TableCell>
              </TableRow>
              
              {/* RESULTADO */}
              <DRERow label="LUCRO OPERACIONAL" value={dre.lucro_operacional} percent={dre.margem_operacional_percent} isTotal />
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// COMPONENTE: Seção Fluxo de Caixa
// =============================================================================

interface CashFlowSectionProps {
  selectedMonth: string;
}

function CashFlowSection({ selectedMonth }: CashFlowSectionProps) {
  const { start, end } = getMonthStartEnd(selectedMonth);
  const filters: ListFluxoCaixaFilters = { data_inicio: start, data_fim: end };
  
  const { data, isLoading } = useCashFlow(filters);
  const cashFlowData = (data || []) as FluxoCaixaItem[];

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <Skeleton className="h-96 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (cashFlowData.length === 0) {
    return (
      <Card>
        <CardContent className="py-16 text-center">
          <Wallet className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-semibold mb-2">Fluxo de Caixa não disponível</h3>
          <p className="text-muted-foreground max-w-md mx-auto">
            Não há movimentações registradas em{' '}
            <span className="font-medium">{getMonthName(selectedMonth)}</span>.
          </p>
        </CardContent>
      </Card>
    );
  }

  // Calcular totais
  const totais = cashFlowData.reduce(
    (acc, item) => ({
      entradas: acc.entradas + parseFloat(item.total_entradas),
      saidas: acc.saidas + parseFloat(item.total_saidas),
    }),
    { entradas: 0, saidas: 0 }
  );
  const saldoFinal = parseFloat(cashFlowData[cashFlowData.length - 1]?.saldo_acumulado || '0');

  return (
    <div className="space-y-6">
      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Entradas
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(totais.entradas)}</div>
            <div className="flex items-center text-xs text-muted-foreground mt-1">
              <ArrowUpCircle className="h-3 w-3 text-green-500 mr-1" />
              {cashFlowData.length} dias com movimentação
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Saídas
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{formatCurrency(totais.saidas)}</div>
            <div className="flex items-center text-xs text-muted-foreground mt-1">
              <ArrowDownCircle className="h-3 w-3 text-red-500 mr-1" />
              Pagamentos e despesas
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Saldo Líquido
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${totais.entradas - totais.saidas >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {formatCurrency(totais.entradas - totais.saidas)}
            </div>
            <div className="flex items-center text-xs text-muted-foreground mt-1">
              Entradas - Saídas
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Saldo Acumulado
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${saldoFinal >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {formatCurrency(saldoFinal)}
            </div>
            <div className="flex items-center text-xs text-muted-foreground mt-1">
              Posição final do período
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabela de Movimentações */}
      <Card>
        <CardHeader>
          <CardTitle>Movimentações Diárias</CardTitle>
          <CardDescription>
            Detalhamento do fluxo de caixa por dia
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Data</TableHead>
                <TableHead className="text-right">Saldo Inicial</TableHead>
                <TableHead className="text-right">Entradas</TableHead>
                <TableHead className="text-right">Saídas</TableHead>
                <TableHead className="text-right">Saldo Final</TableHead>
                <TableHead className="text-right">Acumulado</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {cashFlowData.map((item, index) => {
                const entradas = parseFloat(item.total_entradas);
                const saidas = parseFloat(item.total_saidas);
                const saldoFinalItem = parseFloat(item.saldo_final);
                
                return (
                  <TableRow key={index}>
                    <TableCell className="font-medium">{formatDate(item.data)}</TableCell>
                    <TableCell className="text-right">{formatCurrency(item.saldo_inicial)}</TableCell>
                    <TableCell className="text-right text-green-600">
                      {entradas > 0 ? `+${formatCurrency(entradas)}` : '-'}
                    </TableCell>
                    <TableCell className="text-right text-red-600">
                      {saidas > 0 ? `-${formatCurrency(saidas)}` : '-'}
                    </TableCell>
                    <TableCell className={`text-right font-medium ${saldoFinalItem >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {formatCurrency(saldoFinalItem)}
                    </TableCell>
                    <TableCell className="text-right">{formatCurrency(item.saldo_acumulado)}</TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// COMPONENTE: Seção Faturamento
// =============================================================================

interface FaturamentoSectionProps {
  selectedMonth: string;
}

function FaturamentoSection({ selectedMonth }: FaturamentoSectionProps) {
  const { start, end } = getMonthStartEnd(selectedMonth);
  const filters: ListContasReceberFilters = {
    data_inicio: start,
    data_fim: end,
    page_size: 100,
  };
  
  const { data: receivables = [], isLoading } = useReceivables(filters);

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <Skeleton className="h-96 w-full" />
        </CardContent>
      </Card>
    );
  }

  // Calcular totais
  const totalRecebido = receivables
    .filter((r: ContaReceber) => r.status === StatusContaReceber.RECEBIDO)
    .reduce((acc: number, r: ContaReceber) => acc + parseFloat(r.valor), 0);
  
  const totalPendente = receivables
    .filter((r: ContaReceber) => r.status === StatusContaReceber.PENDENTE)
    .reduce((acc: number, r: ContaReceber) => acc + parseFloat(r.valor), 0);

  const totalAtrasado = receivables
    .filter((r: ContaReceber) => r.status === StatusContaReceber.ATRASADO)
    .reduce((acc: number, r: ContaReceber) => acc + parseFloat(r.valor), 0);

  const totalGeral = receivables.reduce((acc: number, r: ContaReceber) => acc + parseFloat(r.valor), 0);

  if (receivables.length === 0) {
    return (
      <Card>
        <CardContent className="py-16 text-center">
          <ArrowUpCircle className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-semibold mb-2">Sem faturamento registrado</h3>
          <p className="text-muted-foreground max-w-md mx-auto">
            Não há receitas registradas em{' '}
            <span className="font-medium">{getMonthName(selectedMonth)}</span>.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Faturado
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalGeral)}</div>
            <div className="text-xs text-muted-foreground mt-1">
              {receivables.length} lançamentos
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Recebido
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(totalRecebido)}</div>
            <Badge variant="default" className="mt-1">Confirmado</Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Pendente
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{formatCurrency(totalPendente)}</div>
            <Badge variant="secondary" className="mt-1">A receber</Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Atrasado
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{formatCurrency(totalAtrasado)}</div>
            <Badge variant="destructive" className="mt-1">Em atraso</Badge>
          </CardContent>
        </Card>
      </div>

      {/* Tabela de Receitas */}
      <Card>
        <CardHeader>
          <CardTitle>Detalhamento do Faturamento</CardTitle>
          <CardDescription>
            Todas as receitas do período
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Descrição</TableHead>
                <TableHead>Origem</TableHead>
                <TableHead className="text-right">Valor</TableHead>
                <TableHead>Vencimento</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {receivables.slice(0, 20).map((item) => (
                <TableRow key={item.id}>
                  <TableCell className="font-medium">{item.descricao_origem}</TableCell>
                  <TableCell>
                    <Badge variant="outline">{item.origem}</Badge>
                  </TableCell>
                  <TableCell className="text-right">{formatCurrency(item.valor)}</TableCell>
                  <TableCell>{formatDate(item.data_vencimento)}</TableCell>
                  <TableCell>
                    <Badge
                      variant={
                        item.status === StatusContaReceber.RECEBIDO
                          ? 'default'
                          : item.status === StatusContaReceber.ATRASADO
                          ? 'destructive'
                          : 'secondary'
                      }
                    >
                      {item.status}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          {receivables.length > 20 && (
            <p className="text-sm text-muted-foreground text-center mt-4">
              Mostrando 20 de {receivables.length} registros
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// COMPONENTE: Seção Despesas
// =============================================================================

interface DespesasSectionProps {
  selectedMonth: string;
}

function DespesasSection({ selectedMonth }: DespesasSectionProps) {
  const { start, end } = getMonthStartEnd(selectedMonth);
  const filters: ListContasPagarFilters = {
    data_inicio: start,
    data_fim: end,
    page_size: 100,
  };
  
  const { data: payables = [], isLoading } = usePayables(filters);

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <Skeleton className="h-96 w-full" />
        </CardContent>
      </Card>
    );
  }

  // Calcular totais
  const totalPago = payables
    .filter((p: ContaPagar) => p.status === StatusContaPagar.PAGO)
    .reduce((acc: number, p: ContaPagar) => acc + parseFloat(p.valor), 0);
  
  const totalPendente = payables
    .filter((p: ContaPagar) => p.status === StatusContaPagar.PENDENTE)
    .reduce((acc: number, p: ContaPagar) => acc + parseFloat(p.valor), 0);

  const totalAtrasado = payables
    .filter((p: ContaPagar) => p.status === StatusContaPagar.ATRASADO)
    .reduce((acc: number, p: ContaPagar) => acc + parseFloat(p.valor), 0);

  const totalGeral = payables.reduce((acc: number, p: ContaPagar) => acc + parseFloat(p.valor), 0);

  if (payables.length === 0) {
    return (
      <Card>
        <CardContent className="py-16 text-center">
          <ArrowDownCircle className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-semibold mb-2">Sem despesas registradas</h3>
          <p className="text-muted-foreground max-w-md mx-auto">
            Não há despesas registradas em{' '}
            <span className="font-medium">{getMonthName(selectedMonth)}</span>.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Cards de Resumo */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Despesas
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(totalGeral)}</div>
            <div className="text-xs text-muted-foreground mt-1">
              {payables.length} lançamentos
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Pago
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(totalPago)}</div>
            <Badge variant="default" className="mt-1">Quitado</Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Pendente
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{formatCurrency(totalPendente)}</div>
            <Badge variant="secondary" className="mt-1">A pagar</Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Atrasado
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{formatCurrency(totalAtrasado)}</div>
            <Badge variant="destructive" className="mt-1">Em atraso</Badge>
          </CardContent>
        </Card>
      </div>

      {/* Tabela de Despesas */}
      <Card>
        <CardHeader>
          <CardTitle>Detalhamento das Despesas</CardTitle>
          <CardDescription>
            Todas as despesas do período
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Descrição</TableHead>
                <TableHead>Tipo</TableHead>
                <TableHead>Fornecedor</TableHead>
                <TableHead className="text-right">Valor</TableHead>
                <TableHead>Vencimento</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {payables.slice(0, 20).map((item) => (
                <TableRow key={item.id}>
                  <TableCell className="font-medium">{item.descricao}</TableCell>
                  <TableCell>
                    <Badge variant="outline">{item.tipo}</Badge>
                  </TableCell>
                  <TableCell>{item.fornecedor || '-'}</TableCell>
                  <TableCell className="text-right">{formatCurrency(item.valor)}</TableCell>
                  <TableCell>{formatDate(item.data_vencimento)}</TableCell>
                  <TableCell>
                    <Badge
                      variant={
                        item.status === StatusContaPagar.PAGO
                          ? 'default'
                          : item.status === StatusContaPagar.ATRASADO
                          ? 'destructive'
                          : 'secondary'
                      }
                    >
                      {item.status}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          {payables.length > 20 && (
            <p className="text-sm text-muted-foreground text-center mt-4">
              Mostrando 20 de {payables.length} registros
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL: Relatórios
// =============================================================================

export default function RelatoriosPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [activeTab, setActiveTab] = useState<ReportTab>('dre');
  const [selectedMonth, setSelectedMonth] = useState(() => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
  });

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([{ label: 'Relatórios' }]);
  }, [setBreadcrumbs]);

  // Navegação de meses
  const handlePreviousMonth = () => {
    const [year, month] = selectedMonth.split('-').map(Number);
    const date = new Date(year, month - 2);
    setSelectedMonth(
      `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
    );
  };

  const handleNextMonth = () => {
    const [year, month] = selectedMonth.split('-').map(Number);
    const date = new Date(year, month);
    const now = new Date();
    // Não permite avançar além do mês atual
    if (date <= now) {
      setSelectedMonth(
        `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
      );
    }
  };

  const isCurrentMonth = useMemo(() => {
    const now = new Date();
    const current = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
    return selectedMonth === current;
  }, [selectedMonth]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <BarChart3 className="h-8 w-8" />
            Relatórios
          </h1>
          <p className="text-muted-foreground">
            Visualize DRE, Fluxo de Caixa, Faturamento e Despesas consolidados
          </p>
        </div>

        {/* Navegação de Período */}
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={handlePreviousMonth}>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <div className="flex items-center gap-2 px-4 py-2 border rounded-md min-w-[200px] justify-center">
            <Calendar className="h-4 w-4 text-muted-foreground" />
            <span className="font-medium capitalize">{getMonthName(selectedMonth)}</span>
          </div>
          <Button
            variant="outline"
            size="icon"
            onClick={handleNextMonth}
            disabled={isCurrentMonth}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Tabs de Relatórios */}
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as ReportTab)}>
        <TabsList className="grid w-full grid-cols-4 lg:w-[500px]">
          <TabsTrigger value="dre" className="flex items-center gap-2">
            <Calculator className="h-4 w-4" />
            <span className="hidden sm:inline">DRE</span>
          </TabsTrigger>
          <TabsTrigger value="cashflow" className="flex items-center gap-2">
            <Wallet className="h-4 w-4" />
            <span className="hidden sm:inline">Fluxo</span>
          </TabsTrigger>
          <TabsTrigger value="faturamento" className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4" />
            <span className="hidden sm:inline">Receitas</span>
          </TabsTrigger>
          <TabsTrigger value="despesas" className="flex items-center gap-2">
            <TrendingDown className="h-4 w-4" />
            <span className="hidden sm:inline">Despesas</span>
          </TabsTrigger>
        </TabsList>

        <TabsContent value="dre" className="mt-6">
          <DRESection selectedMonth={selectedMonth} />
        </TabsContent>

        <TabsContent value="cashflow" className="mt-6">
          <CashFlowSection selectedMonth={selectedMonth} />
        </TabsContent>

        <TabsContent value="faturamento" className="mt-6">
          <FaturamentoSection selectedMonth={selectedMonth} />
        </TabsContent>

        <TabsContent value="despesas" className="mt-6">
          <DespesasSection selectedMonth={selectedMonth} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
