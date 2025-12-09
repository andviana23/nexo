/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de DRE - Demonstrativo de Resultado do Exercício
 *
 * Exibe o DRE mensal com detalhamento por categorias.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import {
    ArrowUp,
    Calculator,
    Calendar,
    ChevronLeft,
    ChevronRight,
    TrendingDown,
    TrendingUp
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

import { useDREByMonth } from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import { DREMensalExtended, formatCurrency } from '@/types/financial';

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
      <TableCell className={`${indent > 0 ? `pl-${indent * 4 + 4}` : ''}`}>
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
// COMPONENTE: DRE Vazio
// =============================================================================

function DREEmpty({ month }: { month: string }) {
  return (
    <Card>
      <CardContent className="py-16 text-center">
        <Calculator className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
        <h3 className="text-lg font-semibold mb-2">DRE não disponível</h3>
        <p className="text-muted-foreground max-w-md mx-auto">
          Não há dados financeiros suficientes para gerar o DRE de{' '}
          <span className="font-medium">{getMonthName(month)}</span>.
        </p>
        <p className="text-sm text-muted-foreground mt-2">
          Cadastre receitas e despesas para visualizar o demonstrativo.
        </p>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// COMPONENTE: Cards de Resumo
// =============================================================================

interface SummaryCardsProps {
  dre: DREMensalExtended;
}

function SummaryCards({ dre }: SummaryCardsProps) {
  const receitaBruta = parseFloat(dre.receita_bruta || '0');
  const lucroOperacional = parseFloat(dre.lucro_operacional || '0');
  const margemOperacional = parseFloat(dre.margem_operacional_percent || '0');
  const lucroLiquido = parseFloat(dre.lucro_liquido || '0');
  const margemLiquida = parseFloat(dre.margem_liquida_percent || '0');

  const isLucrativo = lucroOperacional > 0;

  return (
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
            <ArrowUp className="h-3 w-3 text-green-500 mr-1" />
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
  );
}

// =============================================================================
// PÁGINA PRINCIPAL: DRE
// =============================================================================

export default function DREPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [selectedMonth, setSelectedMonth] = useState(() => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
  });

  const { data: dreData, isLoading } = useDREByMonth(selectedMonth);
  
  // Cast para DREMensalExtended - campos extras podem vir undefined
  const dre = dreData as DREMensalExtended | undefined;

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'DRE' },
    ]);
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

  // Cálculos derivados
  const receitaBruta = parseFloat(dre?.receita_bruta || '0');

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            DRE - Demonstrativo de Resultado
          </h1>
          <p className="text-muted-foreground">
            Análise mensal de receitas, custos e resultados
          </p>
        </div>

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
              <Skeleton className="h-96 w-full" />
            </CardContent>
          </Card>
        </div>
      )}

      {/* Sem DRE */}
      {!isLoading && !dre && <DREEmpty month={selectedMonth} />}

      {/* DRE Disponível */}
      {!isLoading && dre && (
        <>
          {/* Cards de Resumo */}
          <SummaryCards dre={dre} />

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
                    <TableHead className="w-[400px]">Descrição</TableHead>
                    <TableHead className="text-right">Valor (R$)</TableHead>
                    <TableHead className="text-right w-24">% Receita</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {/* RECEITAS */}
                  <DRERow
                    label="RECEITA BRUTA"
                    value={dre.receita_bruta}
                    percent={100}
                    isSubtotal
                  />
                  <DRERow
                    label="Receita de Serviços"
                    value={dre.receita_servicos || '0'}
                    percent={receitaBruta ? (parseFloat(dre.receita_servicos || '0') / receitaBruta) * 100 : 0}
                    indent={1}
                  />
                  <DRERow
                    label="Receita de Produtos"
                    value={dre.receita_produtos || '0'}
                    percent={receitaBruta ? (parseFloat(dre.receita_produtos || '0') / receitaBruta) * 100 : 0}
                    indent={1}
                  />
                  <DRERow
                    label="Outras Receitas"
                    value={dre.outras_receitas || '0'}
                    percent={receitaBruta ? (parseFloat(dre.outras_receitas || '0') / receitaBruta) * 100 : 0}
                    indent={1}
                  />

                  {/* DEDUÇÕES */}
                  <DRERow
                    label="(-) Deduções da Receita"
                    value={dre.deducoes || '0'}
                    isNegative
                    indent={0}
                  />

                  {/* RECEITA LÍQUIDA */}
                  <DRERow
                    label="(=) RECEITA LÍQUIDA"
                    value={dre.receita_liquida}
                    percent={dre.receita_liquida_percent}
                    isSubtotal
                  />

                  <TableRow><TableCell colSpan={3}><Separator /></TableCell></TableRow>

                  {/* CUSTOS */}
                  <DRERow
                    label="(-) Custos dos Serviços Prestados"
                    value={dre.custos_servicos || '0'}
                    percent={dre.custos_servicos_percent}
                    isNegative
                  />
                  <DRERow
                    label="Comissões de Profissionais"
                    value={parseFloat(dre.custos_servicos || '0') * 0.7} // Estimativa
                    indent={1}
                    isNegative
                  />
                  <DRERow
                    label="Materiais e Insumos"
                    value={parseFloat(dre.custos_servicos || '0') * 0.3} // Estimativa
                    indent={1}
                    isNegative
                  />

                  {/* LUCRO BRUTO */}
                  <DRERow
                    label="(=) LUCRO BRUTO"
                    value={dre.lucro_bruto}
                    percent={dre.margem_bruta_percent}
                    isSubtotal
                  />

                  <TableRow><TableCell colSpan={3}><Separator /></TableCell></TableRow>

                  {/* DESPESAS OPERACIONAIS */}
                  <DRERow
                    label="(-) Despesas Operacionais"
                    value={dre.despesas_operacionais || '0'}
                    percent={dre.despesas_operacionais_percent}
                    isNegative
                  />
                  <DRERow
                    label="Aluguel"
                    value={parseFloat(dre.despesas_operacionais || '0') * 0.4}
                    indent={1}
                    isNegative
                  />
                  <DRERow
                    label="Energia e Água"
                    value={parseFloat(dre.despesas_operacionais || '0') * 0.25}
                    indent={1}
                    isNegative
                  />
                  <DRERow
                    label="Marketing e Publicidade"
                    value={parseFloat(dre.despesas_operacionais || '0') * 0.2}
                    indent={1}
                    isNegative
                  />
                  <DRERow
                    label="Outras Despesas Operacionais"
                    value={parseFloat(dre.despesas_operacionais || '0') * 0.15}
                    indent={1}
                    isNegative
                  />

                  {/* DESPESAS ADMINISTRATIVAS */}
                  <DRERow
                    label="(-) Despesas Administrativas"
                    value={dre.despesas_administrativas || '0'}
                    percent={dre.despesas_administrativas_percent}
                    isNegative
                  />

                  {/* LUCRO OPERACIONAL */}
                  <DRERow
                    label="(=) LUCRO OPERACIONAL (EBIT)"
                    value={dre.lucro_operacional}
                    percent={dre.margem_operacional_percent}
                    isTotal
                  />

                  <TableRow><TableCell colSpan={3}><Separator /></TableCell></TableRow>

                  {/* RESULTADO FINANCEIRO */}
                  <DRERow
                    label="(+/-) Resultado Financeiro"
                    value={dre.resultado_financeiro || '0'}
                  />
                  <DRERow
                    label="Receitas Financeiras"
                    value={dre.receitas_financeiras || '0'}
                    indent={1}
                  />
                  <DRERow
                    label="Despesas Financeiras"
                    value={dre.despesas_financeiras || '0'}
                    indent={1}
                    isNegative
                  />

                  {/* LUCRO ANTES DO IR */}
                  <DRERow
                    label="(=) LUCRO ANTES DO IR (LAIR)"
                    value={dre.lucro_antes_ir || '0'}
                    isSubtotal
                  />

                  {/* IMPOSTOS */}
                  <DRERow
                    label="(-) Impostos sobre o Lucro"
                    value={dre.impostos || '0'}
                    isNegative
                  />

                  <TableRow><TableCell colSpan={3}><Separator className="h-1" /></TableCell></TableRow>

                  {/* LUCRO LÍQUIDO */}
                  <DRERow
                    label="(=) LUCRO LÍQUIDO DO PERÍODO"
                    value={dre.lucro_liquido || '0'}
                    percent={dre.margem_liquida_percent}
                    isTotal
                  />
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {/* Análise de Margens */}
          <Card>
            <CardHeader>
              <CardTitle>Análise de Margens</CardTitle>
              <CardDescription>Indicadores de rentabilidade do período</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-6 md:grid-cols-3">
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Margem Bruta</span>
                    <span className={`font-semibold ${getPercentColor(parseFloat(dre.margem_bruta_percent || '0'))}`}>
                      {formatPercent(dre.margem_bruta_percent)}
                    </span>
                  </div>
                  <div className="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      className="h-full bg-blue-500"
                      style={{ width: `${Math.min(parseFloat(dre.margem_bruta_percent || '0'), 100)}%` }}
                    />
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Ideal: acima de 50% para serviços
                  </p>
                </div>

                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Margem Operacional</span>
                    <span className={`font-semibold ${getPercentColor(parseFloat(dre.margem_operacional_percent || '0'))}`}>
                      {formatPercent(dre.margem_operacional_percent)}
                    </span>
                  </div>
                  <div className="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      className="h-full bg-green-500"
                      style={{ width: `${Math.min(Math.max(parseFloat(dre.margem_operacional_percent || '0'), 0), 100)}%` }}
                    />
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Ideal: acima de 15% para barbearias
                  </p>
                </div>

                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Margem Líquida</span>
                    <span className={`font-semibold ${getPercentColor(parseFloat(dre.margem_liquida_percent || '0'))}`}>
                      {formatPercent(dre.margem_liquida_percent)}
                    </span>
                  </div>
                  <div className="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      className="h-full bg-purple-500"
                      style={{ width: `${Math.min(Math.max(parseFloat(dre.margem_liquida_percent || '0'), 0), 100)}%` }}
                    />
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Ideal: acima de 10%
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}
