'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
    useCommissionPeriods,
    useCommissionRules,
    useCommissionSummaryByProfessional,
    usePendingAdvances
} from '@/hooks/use-commissions';
import { cn } from '@/lib/utils';
import {
    AdvanceStatus,
    CommissionPeriodStatus,
} from '@/types/commission';
import {
    ArrowRight,
    Calculator,
    CalendarDays,
    CheckCircle2,
    DollarSign,
    HandCoins,
    ScrollText,
    Users
} from 'lucide-react';
import Link from 'next/link';
import { useMemo } from 'react';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: string | number) => {
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
  return new Intl.DateTimeFormat('pt-BR', { 
    day: '2-digit', 
    month: '2-digit',
    year: 'numeric'
  }).format(date);
};

const formatMonth = (monthStr: string) => {
  if (!monthStr) return '-';
  const [year, month] = monthStr.split('-');
  const date = new Date(parseInt(year), parseInt(month) - 1);
  return new Intl.DateTimeFormat('pt-BR', { 
    month: 'long', 
    year: 'numeric' 
  }).format(date);
};

const getStatusBadge = (status: CommissionPeriodStatus) => {
  const variants: Record<CommissionPeriodStatus, { variant: 'default' | 'secondary' | 'destructive' | 'outline'; label: string }> = {
    [CommissionPeriodStatus.ABERTO]: { variant: 'outline', label: 'Aberto' },
    [CommissionPeriodStatus.PROCESSANDO]: { variant: 'secondary', label: 'Processando' },
    [CommissionPeriodStatus.FECHADO]: { variant: 'default', label: 'Fechado' },
    [CommissionPeriodStatus.PAGO]: { variant: 'default', label: 'Pago' },
    [CommissionPeriodStatus.CANCELADO]: { variant: 'destructive', label: 'Cancelado' },
  };
  const { variant, label } = variants[status] || { variant: 'outline', label: status };
  return <Badge variant={variant}>{label}</Badge>;
};

const getAdvanceStatusBadge = (status: AdvanceStatus) => {
  const variants: Record<AdvanceStatus, { variant: 'default' | 'secondary' | 'destructive' | 'outline'; label: string }> = {
    [AdvanceStatus.PENDING]: { variant: 'outline', label: 'Pendente' },
    [AdvanceStatus.APPROVED]: { variant: 'default', label: 'Aprovado' },
    [AdvanceStatus.REJECTED]: { variant: 'destructive', label: 'Rejeitado' },
    [AdvanceStatus.DEDUCTED]: { variant: 'secondary', label: 'Deduzido' },
    [AdvanceStatus.CANCELLED]: { variant: 'destructive', label: 'Cancelado' },
  };
  const { variant, label } = variants[status] || { variant: 'outline', label: status };
  return <Badge variant={variant}>{label}</Badge>;
};

// =============================================================================
// KPI CARD
// =============================================================================

interface KpiCardProps {
  title: string;
  value: string | number;
  description?: string;
  icon: React.ReactNode;
  variant?: 'default' | 'success' | 'warning' | 'danger';
  isLoading?: boolean;
  href?: string;
}

function KpiCard({ title, value, description, icon, variant = 'default', isLoading, href }: KpiCardProps) {
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
    default: '',
    success: 'text-green-600',
    warning: 'text-yellow-600',
    danger: 'text-red-600',
  };

  const content = (
    <Card className={cn(href && 'cursor-pointer hover:bg-muted/50 transition-colors')}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className={cn('text-2xl font-bold', variantClasses[variant])}>{value}</div>
        {description && (
          <p className="text-xs text-muted-foreground mt-1">{description}</p>
        )}
      </CardContent>
    </Card>
  );

  if (href) {
    return <Link href={href}>{content}</Link>;
  }

  return content;
}

// =============================================================================
// MAIN PAGE
// =============================================================================

export default function ComissoesDashboardPage() {
  // Data fetching
  const { data: rulesData, isLoading: isLoadingRules } = useCommissionRules({ active_only: true });
  const { data: periodsData, isLoading: isLoadingPeriods } = useCommissionPeriods({ limit: 5 });
  const { data: pendingAdvances, isLoading: isLoadingPendingAdvances } = usePendingAdvances();
  
  // Datas para resumo (mês atual)
  const now = new Date();
  const startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString();
  const endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString();
  
  const { data: summaryByProfessional, isLoading: isLoadingSummary } = useCommissionSummaryByProfessional(
    startDate,
    endDate
  );

  // KPIs calculados
  const kpis = useMemo(() => {
    const totalCommission = summaryByProfessional?.reduce(
      (acc, item) => acc + parseFloat(item.total_commission || '0'),
      0
    ) || 0;
    
    const totalGross = summaryByProfessional?.reduce(
      (acc, item) => acc + parseFloat(item.total_gross || '0'),
      0
    ) || 0;

    const totalPending = parseFloat(pendingAdvances?.total_pending || '0');
    const totalApproved = parseFloat(pendingAdvances?.total_approved || '0');
    
    const openPeriods = periodsData?.data?.filter(
      p => p.status === CommissionPeriodStatus.ABERTO
    ).length || 0;

    const activeRules = rulesData?.data?.length || 0;

    return {
      totalCommission,
      totalGross,
      totalPending,
      totalApproved,
      openPeriods,
      activeRules,
      professionalsCount: summaryByProfessional?.length || 0,
    };
  }, [summaryByProfessional, pendingAdvances, periodsData, rulesData]);

  const isLoading = isLoadingRules || isLoadingPeriods || isLoadingPendingAdvances || isLoadingSummary;

  return (
    <div className="flex flex-col gap-6">
      {/* KPIs */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <KpiCard
          title="Comissões do Mês"
          value={formatCurrency(kpis.totalCommission)}
          description={`Faturamento bruto: ${formatCurrency(kpis.totalGross)}`}
          icon={<DollarSign className="h-4 w-4 text-muted-foreground" />}
          variant={kpis.totalCommission > 0 ? 'success' : 'default'}
          isLoading={isLoading}
        />
        <KpiCard
          title="Adiantamentos Pendentes"
          value={formatCurrency(kpis.totalPending)}
          description={`Aprovados: ${formatCurrency(kpis.totalApproved)}`}
          icon={<HandCoins className="h-4 w-4 text-muted-foreground" />}
          variant={kpis.totalPending > 0 ? 'warning' : 'default'}
          isLoading={isLoading}
          href="/comissoes/adiantamentos"
        />
        <KpiCard
          title="Períodos Abertos"
          value={kpis.openPeriods}
          description={`${kpis.professionalsCount} profissionais com comissões`}
          icon={<CalendarDays className="h-4 w-4 text-muted-foreground" />}
          isLoading={isLoading}
          href="/comissoes/periodos"
        />
        <KpiCard
          title="Regras Ativas"
          value={kpis.activeRules}
          description="Regras de comissão configuradas"
          icon={<ScrollText className="h-4 w-4 text-muted-foreground" />}
          isLoading={isLoading}
          href="/comissoes/regras"
        />
      </div>

      {/* Content Grid */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Resumo por Profissional */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base">Comissões por Profissional</CardTitle>
              <CardDescription>Resumo do mês atual</CardDescription>
            </div>
            <Link href="/comissoes/itens">
              <Button variant="ghost" size="sm">
                Ver Todos <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {isLoadingSummary ? (
              <div className="space-y-3">
                {Array.from({ length: 4 }).map((_, i) => (
                  <div key={i} className="flex items-center justify-between">
                    <div className="space-y-1">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="h-3 w-24" />
                    </div>
                    <Skeleton className="h-6 w-20" />
                  </div>
                ))}
              </div>
            ) : summaryByProfessional && summaryByProfessional.length > 0 ? (
              <div className="space-y-4">
                {summaryByProfessional.slice(0, 5).map((item) => (
                  <div key={item.professional_id} className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">{item.professional_name}</p>
                      <p className="text-sm text-muted-foreground">
                        {item.items_count} {item.items_count === 1 ? 'item' : 'itens'}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-green-600">
                        {formatCurrency(item.total_commission)}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        de {formatCurrency(item.total_gross)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-center">
                <Users className="h-10 w-10 text-muted-foreground/50 mb-2" />
                <p className="text-sm text-muted-foreground">
                  Nenhuma comissão registrada neste mês
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Períodos Recentes */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base">Períodos Recentes</CardTitle>
              <CardDescription>Últimos períodos de comissão</CardDescription>
            </div>
            <Link href="/comissoes/periodos">
              <Button variant="ghost" size="sm">
                Ver Todos <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {isLoadingPeriods ? (
              <div className="space-y-3">
                {Array.from({ length: 4 }).map((_, i) => (
                  <div key={i} className="flex items-center justify-between">
                    <div className="space-y-1">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="h-3 w-24" />
                    </div>
                    <Skeleton className="h-6 w-20" />
                  </div>
                ))}
              </div>
            ) : periodsData?.data && periodsData.data.length > 0 ? (
              <div className="space-y-4">
                {periodsData.data.slice(0, 5).map((period) => (
                  <div key={period.id} className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">{period.professional_name || 'Profissional'}</p>
                      <p className="text-sm text-muted-foreground">
                        {formatMonth(period.reference_month)}
                      </p>
                    </div>
                    <div className="text-right flex flex-col items-end gap-1">
                      {getStatusBadge(period.status as CommissionPeriodStatus)}
                      <p className="text-sm font-medium">
                        {formatCurrency(period.total_net)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-center">
                <CalendarDays className="h-10 w-10 text-muted-foreground/50 mb-2" />
                <p className="text-sm text-muted-foreground">
                  Nenhum período de comissão encontrado
                </p>
                <Link href="/comissoes/periodos">
                  <Button variant="link" size="sm" className="mt-2">
                    Criar período
                  </Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Adiantamentos Pendentes */}
        <Card className="md:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base">Adiantamentos Pendentes</CardTitle>
              <CardDescription>Aguardando aprovação</CardDescription>
            </div>
            <Link href="/comissoes/adiantamentos">
              <Button variant="ghost" size="sm">
                Ver Todos <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {isLoadingPendingAdvances ? (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Profissional</TableHead>
                    <TableHead>Valor</TableHead>
                    <TableHead>Data</TableHead>
                    <TableHead>Motivo</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {Array.from({ length: 3 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : pendingAdvances?.advances && pendingAdvances.advances.length > 0 ? (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Profissional</TableHead>
                    <TableHead>Valor</TableHead>
                    <TableHead>Data</TableHead>
                    <TableHead>Motivo</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {pendingAdvances.advances.slice(0, 5).map((advance) => (
                    <TableRow key={advance.id}>
                      <TableCell className="font-medium">
                        {advance.professional_name || '-'}
                      </TableCell>
                      <TableCell>{formatCurrency(advance.amount)}</TableCell>
                      <TableCell>{formatDate(advance.request_date)}</TableCell>
                      <TableCell className="max-w-[200px] truncate">
                        {advance.reason || '-'}
                      </TableCell>
                      <TableCell>
                        {getAdvanceStatusBadge(advance.status as AdvanceStatus)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-center">
                <CheckCircle2 className="h-10 w-10 text-green-500/50 mb-2" />
                <p className="text-sm text-muted-foreground">
                  Nenhum adiantamento pendente de aprovação
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Ações Rápidas</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-3">
            <Link href="/comissoes/regras">
              <Button variant="outline" size="sm">
                <ScrollText className="mr-2 h-4 w-4" />
                Gerenciar Regras
              </Button>
            </Link>
            <Link href="/comissoes/periodos">
              <Button variant="outline" size="sm">
                <CalendarDays className="mr-2 h-4 w-4" />
                Criar Período
              </Button>
            </Link>
            <Link href="/comissoes/adiantamentos">
              <Button variant="outline" size="sm">
                <HandCoins className="mr-2 h-4 w-4" />
                Novo Adiantamento
              </Button>
            </Link>
            <Link href="/comissoes/itens">
              <Button variant="outline" size="sm">
                <Calculator className="mr-2 h-4 w-4" />
                Ver Itens
              </Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
