'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página Principal de Assinaturas
 *
 * @page /assinatura
 * @description Dashboard de assinaturas com métricas e navegação
 * Conforme FLUXO_ASSINATURA.md
 */

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useSubscriptionMetrics } from '@/hooks/use-subscriptions';
import { cn } from '@/lib/utils';
import { useBreadcrumbs } from '@/store/ui-store';
import {
    AlertTriangleIcon,
    ArrowUpRightIcon,
    CalendarIcon,
    CheckCircleIcon,
    CreditCardIcon,
    LayoutListIcon,
    PlusIcon,
    SettingsIcon,
    TrendingUpIcon,
    UsersIcon,
} from 'lucide-react';
import Link from 'next/link';
import { useEffect } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface MetricCardProps {
  title: string;
  value: string | number;
  description?: string;
  icon: React.ReactNode;
  trend?: {
    value: number;
    label: string;
  };
  variant?: 'default' | 'success' | 'warning' | 'danger';
  isLoading?: boolean;
}

// =============================================================================
// COMPONENTES
// =============================================================================

function MetricCard({
  title,
  value,
  description,
  icon,
  trend,
  variant = 'default',
  isLoading = false,
}: MetricCardProps) {
  const variantStyles = {
    default: 'bg-card',
    success: 'bg-green-50 dark:bg-green-950/20 border-green-200 dark:border-green-900',
    warning: 'bg-yellow-50 dark:bg-yellow-950/20 border-yellow-200 dark:border-yellow-900',
    danger: 'bg-red-50 dark:bg-red-950/20 border-red-200 dark:border-red-900',
  };

  const iconStyles = {
    default: 'text-muted-foreground',
    success: 'text-green-600 dark:text-green-400',
    warning: 'text-yellow-600 dark:text-yellow-400',
    danger: 'text-red-600 dark:text-red-400',
  };

  if (isLoading) {
    return (
      <Card className={cn('transition-all', variantStyles[variant])}>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-8 w-8 rounded-full" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-32 mb-2" />
          <Skeleton className="h-3 w-20" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={cn('transition-all hover:shadow-md', variantStyles[variant])}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <div className={cn('h-8 w-8 flex items-center justify-center rounded-full bg-muted', iconStyles[variant])}>
          {icon}
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {description && (
          <p className="text-xs text-muted-foreground">{description}</p>
        )}
        {trend && (
          <div className="flex items-center gap-1 mt-2">
            <TrendingUpIcon
              className={cn(
                'h-3 w-3',
                trend.value >= 0 ? 'text-green-500' : 'text-red-500 rotate-180'
              )}
            />
            <span
              className={cn(
                'text-xs font-medium',
                trend.value >= 0 ? 'text-green-600' : 'text-red-600'
              )}
            >
              {trend.value > 0 ? '+' : ''}
              {trend.value}%
            </span>
            <span className="text-xs text-muted-foreground">{trend.label}</span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export default function AssinaturaPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const { data: metrics, isLoading } = useSubscriptionMetrics();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Assinaturas' },
    ]);
  }, [setBreadcrumbs]);

  // Formatadores
  const formatCurrency = (value?: number) => {
    if (value === undefined || value === null) return 'R$ 0,00';
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
    }).format(value);
  };

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Assinaturas</h1>
          <p className="text-muted-foreground">
            Gerencie planos e assinaturas de clientes
          </p>
        </div>
        <div className="flex gap-2">
          <Link href="/assinatura/planos">
            <Button variant="outline">
              <SettingsIcon className="mr-2 h-4 w-4" />
              Gerenciar Planos
            </Button>
          </Link>
          <Link href="/assinatura/nova">
            <Button>
              <PlusIcon className="mr-2 h-4 w-4" />
              Nova Assinatura
            </Button>
          </Link>
        </div>
      </div>

      {/* Cards de Métricas */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="Assinantes Ativos"
          value={metrics?.total_assinantes_ativos ?? 0}
          description="Clientes com assinatura ativa"
          icon={<UsersIcon className="h-4 w-4" />}
          variant="success"
          isLoading={isLoading}
        />
        <MetricCard
          title="Receita Mensal"
          value={formatCurrency(metrics?.receita_mensal)}
          description="MRR - Monthly Recurring Revenue"
          icon={<CreditCardIcon className="h-4 w-4" />}
          isLoading={isLoading}
        />
        <MetricCard
          title="Inadimplentes"
          value={metrics?.total_inadimplentes ?? 0}
          description="Assinaturas vencidas"
          icon={<AlertTriangleIcon className="h-4 w-4" />}
          variant={metrics?.total_inadimplentes ? 'danger' : 'default'}
          isLoading={isLoading}
        />
        <MetricCard
          title="Taxa de Renovação"
          value={`${metrics?.taxa_renovacao?.toFixed(1) ?? 0}%`}
          description="Últimos 30 dias"
          icon={<TrendingUpIcon className="h-4 w-4" />}
          isLoading={isLoading}
        />
      </div>

      {/* Navegação por Abas */}
      <Tabs defaultValue="visao-geral" className="space-y-4">
        <TabsList>
          <TabsTrigger value="visao-geral">Visão Geral</TabsTrigger>
          <TabsTrigger value="renovacoes">Próximas Renovações</TabsTrigger>
          <TabsTrigger value="inadimplentes">Inadimplentes</TabsTrigger>
        </TabsList>

        <TabsContent value="visao-geral" className="space-y-4">
          {/* Ações Rápidas */}
          <div className="grid gap-4 md:grid-cols-3">
            <Link href="/assinatura/assinantes" className="block">
              <Card className="h-full transition-all hover:shadow-md hover:border-primary/50 cursor-pointer">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <LayoutListIcon className="h-5 w-5 text-primary" />
                    Lista de Assinantes
                  </CardTitle>
                  <CardDescription>
                    Visualizar e gerenciar todos os assinantes
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <Badge variant="secondary">
                      {metrics?.total_assinantes_ativos ?? 0} ativos
                    </Badge>
                    <ArrowUpRightIcon className="h-4 w-4 text-muted-foreground" />
                  </div>
                </CardContent>
              </Card>
            </Link>

            <Link href="/assinatura/planos" className="block">
              <Card className="h-full transition-all hover:shadow-md hover:border-primary/50 cursor-pointer">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <SettingsIcon className="h-5 w-5 text-primary" />
                    Planos
                  </CardTitle>
                  <CardDescription>
                    Criar e editar planos de assinatura
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <Badge variant="secondary">
                      {metrics?.total_planos_ativos ?? 0} planos ativos
                    </Badge>
                    <ArrowUpRightIcon className="h-4 w-4 text-muted-foreground" />
                  </div>
                </CardContent>
              </Card>
            </Link>

            <Link href="/assinatura/nova" className="block">
              <Card className="h-full transition-all hover:shadow-md hover:border-primary/50 cursor-pointer">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <PlusIcon className="h-5 w-5 text-primary" />
                    Nova Assinatura
                  </CardTitle>
                  <CardDescription>
                    Criar assinatura para um cliente
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">
                      Selecionar cliente e plano
                    </span>
                    <ArrowUpRightIcon className="h-4 w-4 text-muted-foreground" />
                  </div>
                </CardContent>
              </Card>
            </Link>
          </div>

          {/* Resumo Financeiro */}
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Resumo Financeiro</CardTitle>
                <CardDescription>
                  Receitas de assinaturas do mês atual
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">
                    Receita Recorrente (MRR)
                  </span>
                  <span className="font-semibold">
                    {formatCurrency(metrics?.receita_mensal)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">
                    Previsão Anual (ARR)
                  </span>
                  <span className="font-semibold">
                    {formatCurrency((metrics?.receita_mensal ?? 0) * 12)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">
                    Ticket Médio
                  </span>
                  <span className="font-semibold">
                    {formatCurrency(
                      metrics?.total_assinantes_ativos
                        ? (metrics?.receita_mensal ?? 0) /
                            metrics.total_assinantes_ativos
                        : 0
                    )}
                  </span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Saúde das Assinaturas</CardTitle>
                <CardDescription>
                  Status geral do módulo de assinaturas
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-2">
                  <CheckCircleIcon className="h-4 w-4 text-green-500" />
                  <span className="text-sm">
                    {metrics?.total_assinantes_ativos ?? 0} assinaturas ativas
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <CalendarIcon className="h-4 w-4 text-blue-500" />
                  <span className="text-sm">
                    {metrics?.renovacoes_proximos_7_dias ?? 0} renovações nos
                    próximos 7 dias
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <AlertTriangleIcon
                    className={cn(
                      'h-4 w-4',
                      metrics?.total_inadimplentes
                        ? 'text-red-500'
                        : 'text-muted-foreground'
                    )}
                  />
                  <span className="text-sm">
                    {metrics?.total_inadimplentes ?? 0} assinaturas inadimplentes
                  </span>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="renovacoes" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Próximas Renovações</CardTitle>
              <CardDescription>
                Assinaturas que vencem nos próximos 7 dias
              </CardDescription>
            </CardHeader>
            <CardContent>
              {metrics?.renovacoes_proximos_7_dias === 0 ? (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <CalendarIcon className="h-12 w-12 text-muted-foreground/50 mb-4" />
                  <p className="text-muted-foreground">
                    Nenhuma renovação prevista para os próximos 7 dias
                  </p>
                </div>
              ) : (
                <div className="text-center py-8">
                  <Link href="/assinatura/assinantes?filter=vencendo">
                    <Button variant="outline">
                      Ver {metrics?.renovacoes_proximos_7_dias ?? 0} assinaturas
                    </Button>
                  </Link>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="inadimplentes" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Assinaturas Inadimplentes</CardTitle>
              <CardDescription>
                Clientes com pagamento em atraso
              </CardDescription>
            </CardHeader>
            <CardContent>
              {metrics?.total_inadimplentes === 0 ? (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <CheckCircleIcon className="h-12 w-12 text-green-500/50 mb-4" />
                  <p className="text-muted-foreground">
                    Nenhuma assinatura inadimplente! 🎉
                  </p>
                </div>
              ) : (
                <div className="text-center py-8">
                  <Link href="/assinatura/assinantes?filter=inadimplente">
                    <Button variant="destructive">
                      Ver {metrics?.total_inadimplentes ?? 0} inadimplentes
                    </Button>
                  </Link>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
