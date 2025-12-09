'use client';

import {
    AlertCircle,
    ArrowDownCircle,
    ArrowUpCircle,
    Banknote,
    Calendar,
    CheckCircle2,
    DollarSign,
    PiggyBank,
    Target,
    TrendingDown,
    TrendingUp,
    Wallet
} from 'lucide-react';
import Link from 'next/link';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
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
    useDashboard,
    useDREByMonth,
    useFinancialSummary,
    useProjections,
    useProximosVencimentos,
} from '@/hooks/use-financial';
import { cn } from '@/lib/utils';
import financialService from '@/services/financial-service';
import { ProjecaoMensal } from '@/types/financial';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: number | string) => {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return 'R$ 0,00';
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: '2-digit' }).format(date);
};

type ProximoVencimento = Awaited<
  ReturnType<typeof financialService.getProximosVencimentos>
>[number] & {
  pessoa_nome?: string;
  categoria_nome?: string;
};

// =============================================================================
// KPI CARD COMPONENT
// =============================================================================

interface KpiCardProps {
  title: string;
  value: number;
  description?: string;
  icon: React.ReactNode;
  variant?: 'default' | 'success' | 'danger' | 'warning';
  isLoading?: boolean;
}

function KpiCard({ title, value, description, icon, variant = 'default', isLoading }: KpiCardProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <Skeleton className="h-4 w-[100px]" />
          <Skeleton className="h-4 w-4" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-7 w-[120px]" />
          <Skeleton className="h-3 w-[100px] mt-1" />
        </CardContent>
      </Card>
    );
  }

  const variantClasses = {
    default: 'text-primary',
    success: 'text-chart-1',
    danger: 'text-destructive',
    warning: 'text-chart-4',
  };

  const iconWrapperClasses = {
    default: 'text-muted-foreground',
    success: 'text-chart-1',
    danger: 'text-destructive',
    warning: 'text-chart-4',
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <div className={cn(iconWrapperClasses[variant])}>{icon}</div>
      </CardHeader>
      <CardContent>
        <div className={cn('text-2xl font-bold', variantClasses[variant])}>
          {formatCurrency(value)}
        </div>
        {description && <p className="text-xs text-muted-foreground mt-1">{description}</p>}
      </CardContent>
    </Card>
  );
}

// =============================================================================
// PRÓXIMOS VENCIMENTOS
// =============================================================================

function ProximosVencimentosCard() {
  const { data: vencimentos, isLoading } = useProximosVencimentos();

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Próximos Vencimentos</CardTitle>
          <CardDescription>Contas dos próximos 7 dias</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!vencimentos || vencimentos.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-4 w-4" />
            Próximos Vencimentos
          </CardTitle>
          <CardDescription>Próximos 7 dias</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
            <CheckCircle2 className="h-10 w-10 mb-2 opacity-50" />
            <p className="text-sm">Nenhuma conta pendente</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-4 w-4" />
              Próximos Vencimentos
            </CardTitle>
            <CardDescription>Contas dos próximos 7 dias</CardDescription>
          </div>
          <Button variant="ghost" size="sm" asChild>
            <Link href="/financeiro/contas-pagar">Ver todas</Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Vencimento</TableHead>
              <TableHead>Descrição</TableHead>
              <TableHead>Categoria</TableHead>
              <TableHead className="text-right">Valor</TableHead>
              <TableHead className="text-right">Tipo</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {vencimentos.map((item: ProximoVencimento) => {
              const isDespesa = item.tipo === 'PAGAR';
              const isAtrasado =
                new Date(item.dataVencimento) < new Date();

              return (
                <TableRow key={item.id}>
                  <TableCell className="font-medium">
                    <div className={cn('flex items-center gap-1', isAtrasado && 'text-destructive')}>
                      {isAtrasado && <AlertCircle className="h-3 w-3" />}
                      {formatDate(item.dataVencimento)}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{item.descricao}</p>
                      {item.pessoa_nome && (
                        <p className="text-xs text-muted-foreground">{item.pessoa_nome}</p>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">{item.categoria_nome || 'Geral'}</Badge>
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    <span className={isDespesa ? 'text-destructive' : 'text-chart-1'}>
                      {formatCurrency(item.valor)}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <Badge variant={isDespesa ? 'destructive' : 'default'}>
                      {isDespesa ? 'Pagar' : 'Receber'}
                    </Badge>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// DRE RESUMO
// =============================================================================

function DREResumoCard() {
  const mesAtual = new Date().toISOString().slice(0, 7);
  const { data: dre, isLoading } = useDREByMonth(mesAtual);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>DRE Gerencial</CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-[200px] w-full" />
        </CardContent>
      </Card>
    );
  }

  if (!dre) return null;

  const receitaBruta = parseFloat(dre.receita_bruta);
  const custosServicos = parseFloat(dre.custos_servicos);
  const lucroBruto = parseFloat(dre.lucro_bruto);
  const despesasOperacionais = parseFloat(dre.despesas_operacionais) + parseFloat(dre.despesas_administrativas);
  const lucroOperacional = parseFloat(dre.lucro_operacional);

  const margemBrutaPct = receitaBruta > 0 ? (lucroBruto / receitaBruta) * 100 : 0;
  const margemOperacional = receitaBruta > 0 ? (lucroOperacional / receitaBruta) * 100 : 0;
  const isLucrativo = lucroOperacional >= 0;

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Banknote className="h-4 w-4" />
              DRE Gerencial
            </CardTitle>
            <CardDescription>Demonstrativo do Resultado (Mês Atual)</CardDescription>
          </div>
          <Button variant="outline" size="sm" asChild>
            <Link href="/financeiro/dre">Ver completo</Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex justify-between items-center text-sm">
            <span className="text-muted-foreground">Receita Bruta</span>
            <span className="font-semibold text-chart-1">{formatCurrency(receitaBruta)}</span>
          </div>
          <Progress value={100} className="h-2" />
        </div>

        <div className="flex justify-between items-center text-sm">
          <span className="text-muted-foreground">(-) Custos de Serviços</span>
          <span className="font-medium text-destructive">{formatCurrency(custosServicos)}</span>
        </div>

        <Separator />

        <div>
          <div className="flex justify-between items-center text-sm mb-1">
            <span className="font-medium">(=) Lucro Bruto</span>
            <span className="font-semibold">{formatCurrency(lucroBruto)}</span>
          </div>
          <p className="text-xs text-muted-foreground text-right">
            {margemBrutaPct.toFixed(1)}% da receita
          </p>
        </div>

        <div className="flex justify-between items-center text-sm">
          <span className="text-muted-foreground">(-) Despesas Operacionais</span>
          <span className="font-medium text-destructive">{formatCurrency(despesasOperacionais)}</span>
        </div>

        <Separator />

        <div className="rounded-lg bg-muted/50 p-3">
          <div className="flex justify-between items-center mb-1">
            <span className="font-bold text-sm">(=) Resultado Operacional</span>
            <span className={cn('font-bold text-lg', isLucrativo ? 'text-chart-1' : 'text-destructive')}>
              {formatCurrency(lucroOperacional)}
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-xs text-muted-foreground">Margem Operacional</span>
            <Badge variant={isLucrativo ? 'default' : 'destructive'}>
              {margemOperacional.toFixed(1)}%
            </Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// AÇÕES RÁPIDAS
// =============================================================================

function AcoesRapidasCard() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Ações Rápidas</CardTitle>
        <CardDescription>Operações frequentes</CardDescription>
      </CardHeader>
      <CardContent className="grid grid-cols-2 gap-3">
        <Button variant="outline" asChild>
          <Link href="/financeiro/contas-pagar?modal=new" className="h-auto py-4 flex flex-col gap-2">
            <ArrowDownCircle className="h-5 w-5 text-destructive" />
            <span className="text-xs">Nova Despesa</span>
          </Link>
        </Button>

        <Button variant="outline" asChild>
          <Link href="/financeiro/contas-receber?modal=new" className="h-auto py-4 flex flex-col gap-2">
            <ArrowUpCircle className="h-5 w-5 text-chart-1" />
            <span className="text-xs">Nova Receita</span>
          </Link>
        </Button>

        <Button variant="outline" asChild>
          <Link href="/financeiro/dre" className="h-auto py-4 flex flex-col gap-2">
            <Banknote className="h-5 w-5 text-primary" />
            <span className="text-xs">Ver DRE</span>
          </Link>
        </Button>

        <Button variant="outline" asChild>
          <Link href="/financeiro/fluxo-caixa" className="h-auto py-4 flex flex-col gap-2">
            <TrendingUp className="h-5 w-5 text-chart-3" />
            <span className="text-xs">Fluxo de Caixa</span>
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// META DO MÊS
// =============================================================================

function MetaMensalCard() {
  const now = new Date();
  const { data: dashboard, isLoading } = useDashboard(now.getFullYear(), now.getMonth() + 1);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Target className="h-4 w-4" />
            Meta do Mês
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-4 w-full mb-2" />
          <Skeleton className="h-8 w-24" />
        </CardContent>
      </Card>
    );
  }

  if (!dashboard) return null;

  const percentual = parseFloat(dashboard.percentual_meta || '0');
  const metaValor = parseFloat(dashboard.meta_mensal || '0');
  const realizado = parseFloat(dashboard.receita_realizada || '0');

  const getProgressColor = (status: string) => {
    switch (status) {
      case 'Atingida':
        return 'bg-chart-1';
      case 'Em andamento':
        return 'bg-chart-4';
      case 'Abaixo':
        return 'bg-destructive';
      default:
        return 'bg-muted';
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <Target className="h-4 w-4" />
            Meta do Mês
          </CardTitle>
          <Badge variant="outline">{dashboard.nome_mes}</Badge>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {dashboard.status_meta === 'Sem meta' ? (
          <div className="text-center py-6 text-muted-foreground">
            <Target className="h-10 w-10 mx-auto mb-2 opacity-50" />
            <p className="text-sm">Nenhuma meta definida</p>
            <Button variant="link" size="sm" asChild>
              <Link href="/metas">Definir meta</Link>
            </Button>
          </div>
        ) : (
          <>
            <div className="space-y-2">
              <div className="flex justify-between items-end">
                <span className="text-3xl font-bold">{percentual.toFixed(0)}%</span>
                <Badge variant={dashboard.status_meta === 'Atingida' ? 'default' : 'secondary'}>
                  {dashboard.status_meta}
                </Badge>
              </div>
              <Progress value={Math.min(percentual, 100)} className={cn('h-2', getProgressColor(dashboard.status_meta))} />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <p className="text-xs text-muted-foreground">Realizado</p>
                <p className="font-semibold">{formatCurrency(realizado)}</p>
              </div>
              <div className="space-y-1 text-right">
                <p className="text-xs text-muted-foreground">Meta</p>
                <p className="font-semibold">{formatCurrency(metaValor)}</p>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

// =============================================================================
// PROJEÇÕES
// =============================================================================

function ProjecoesCard() {
  const { data: projecoes, isLoading } = useProjections(3);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4" />
            Projeções
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-16 w-full" />
          ))}
        </CardContent>
      </Card>
    );
  }

  if (!projecoes || projecoes.projecoes.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4" />
            Projeções
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
            <TrendingUp className="h-10 w-10 mb-2 opacity-50" />
            <p className="text-sm">Sem dados para projeção</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const getTendenciaIcon = (tendencia: string) => {
    switch (tendencia) {
      case 'Crescente':
        return <TrendingUp className="h-4 w-4 text-chart-1" />;
      case 'Decrescente':
        return <TrendingDown className="h-4 w-4 text-destructive" />;
      default:
        return null;
    }
  };

  const getConfiancaBadge = (confianca: string) => {
    const variants: Record<string, 'default' | 'secondary' | 'destructive'> = {
      Alta: 'default',
      Média: 'secondary',
      Baixa: 'destructive',
    };
    return <Badge variant={variants[confianca] || 'secondary'}>{confianca}</Badge>;
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-4 w-4" />
              Projeções
            </CardTitle>
            <CardDescription>Próximos 3 meses</CardDescription>
          </div>
          <div className="flex items-center gap-2 text-sm">
            {getTendenciaIcon(projecoes.tendencia_receita)}
            <span className="font-medium text-muted-foreground">{projecoes.tendencia_receita}</span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {projecoes.projecoes.map((proj: ProjecaoMensal) => {
          const lucro = parseFloat(proj.lucro_projetado);
          const isPositivo = lucro >= 0;

          return (
            <div key={`${proj.ano}-${proj.mes}`} className="rounded-lg border p-3 space-y-2">
              <div className="flex items-center justify-between">
                <span className="font-semibold text-sm">
                  {proj.nome_mes}/{proj.ano}
                </span>
                {getConfiancaBadge(proj.confianca)}
              </div>

              <div className="grid grid-cols-3 gap-2 text-xs">
                <div>
                  <p className="text-muted-foreground">Receita</p>
                  <p className="font-medium text-chart-1">
                    {formatCurrency(parseFloat(proj.receita_projetada))}
                  </p>
                </div>
                <div>
                  <p className="text-muted-foreground">Despesas</p>
                  <p className="font-medium text-destructive">
                    {formatCurrency(parseFloat(proj.despesas_projetadas))}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-muted-foreground">Lucro</p>
                  <p className={cn('font-bold', isPositivo ? 'text-chart-1' : 'text-destructive')}>
                    {formatCurrency(lucro)}
                  </p>
                </div>
              </div>
            </div>
          );
        })}

        <Separator />

        <div className="flex justify-between text-xs text-muted-foreground">
          <span>Média (3 meses)</span>
          <span className="font-medium">
            {formatCurrency(parseFloat(projecoes.media_receita_3_meses))}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL
// =============================================================================

export default function FinanceiroDashboardPage() {
  const now = new Date();
  const { data: dashboard, isLoading: isLoadingDashboard } = useDashboard(
    now.getFullYear(),
    now.getMonth() + 1
  );
  const { data: summary, isLoading: isLoadingSummary } = useFinancialSummary();

  const receitaRealizada = dashboard
    ? parseFloat(dashboard.receita_realizada)
    : summary?.totalAReceber ?? 0;
  const despesasPagas = dashboard
    ? parseFloat(dashboard.despesas_pagas)
    : summary?.totalAPagar ?? 0;
  const lucroLiquido = dashboard
    ? parseFloat(dashboard.lucro_liquido)
    : receitaRealizada - despesasPagas;
  const saldoCaixa = dashboard ? parseFloat(dashboard.saldo_caixa_atual) : 0;

  const variacao = dashboard ? parseFloat(dashboard.variacao_mes_anterior) : 0;
  const tendencia = dashboard?.tendencia_variacao;

  return (
    <div className="flex-1 space-y-4 p-8 pt-6">
      {/* Header */}
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Dashboard Financeiro</h2>
          <p className="text-muted-foreground">
            Visão consolidada — {dashboard?.nome_mes || 'Mês atual'}
          </p>
        </div>
        <div className="flex items-center space-x-2">
          {tendencia && (
            <Badge variant="outline" className="gap-1">
              {tendencia === 'up' ? (
                <TrendingUp className="h-3 w-3 text-chart-1" />
              ) : tendencia === 'down' ? (
                <TrendingDown className="h-3 w-3 text-destructive" />
              ) : null}
              {variacao > 0 ? '+' : ''}
              {variacao.toFixed(1)}%
            </Badge>
          )}
          <Button>
            <Wallet className="mr-2 h-4 w-4" />
            Novo Lançamento
          </Button>
        </div>
      </div>

      {/* KPIs */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <KpiCard
          title="Receita Realizada"
          value={receitaRealizada}
          description={
            dashboard ? `${formatCurrency(dashboard.receita_pendente)} a receber` : undefined
          }
          icon={<DollarSign className="h-4 w-4" />}
          variant="success"
          isLoading={isLoadingDashboard || isLoadingSummary}
        />
        <KpiCard
          title="Despesas Pagas"
          value={despesasPagas}
          description={
            dashboard ? `${formatCurrency(dashboard.despesas_pendentes)} a pagar` : undefined
          }
          icon={<ArrowDownCircle className="h-4 w-4" />}
          variant="danger"
          isLoading={isLoadingDashboard || isLoadingSummary}
        />
        <KpiCard
          title="Lucro Líquido"
          value={lucroLiquido}
          description={dashboard ? `Margem: ${dashboard.margem_liquida}%` : undefined}
          icon={<TrendingUp className="h-4 w-4" />}
          variant={lucroLiquido >= 0 ? 'default' : 'warning'}
          isLoading={isLoadingDashboard || isLoadingSummary}
        />
        <KpiCard
          title="Saldo em Caixa"
          value={saldoCaixa}
          description="Disponível agora"
          icon={<PiggyBank className="h-4 w-4" />}
          variant={saldoCaixa >= 0 ? 'success' : 'danger'}
          isLoading={isLoadingDashboard}
        />
      </div>

      {/* Grid Principal */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Coluna Esquerda */}
        <div className="space-y-4 lg:col-span-4">
          <ProximosVencimentosCard />
          <DREResumoCard />
        </div>

        {/* Coluna Direita */}
        <div className="space-y-4 lg:col-span-3">
          <AcoesRapidasCard />
          <MetaMensalCard />
          <ProjecoesCard />
        </div>
      </div>
    </div>
  );
}
