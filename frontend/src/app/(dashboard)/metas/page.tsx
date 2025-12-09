'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Skeleton } from '@/components/ui/skeleton';
import { useHistoricoMetas, useRankingBarbeiros, useResumoMetas } from '@/hooks/use-metas';
import {
    formatCurrency,
    formatMesAno,
    formatMesAnoShort,
    formatNivelBonificacao,
    formatPercentual,
    getMesAnoAtual,
    getNivelBonificacaoClass
} from '@/types/metas';
import { ArrowDownRight, ArrowUpRight, Award, Target, TrendingUp, Users } from 'lucide-react';
import { useMemo } from 'react';

export default function MetasDashboardPage() {
  const mesAnoAtual = getMesAnoAtual();
  
  const { data: resumo, isLoading: isLoadingResumo } = useResumoMetas(mesAnoAtual);
  const { data: historico, isLoading: isLoadingHistorico } = useHistoricoMetas(6);
  const { data: ranking, isLoading: isLoadingRanking } = useRankingBarbeiros(mesAnoAtual);

  // KPIs
  const kpis = useMemo(() => {
    if (!resumo) return null;

    return [
      {
        title: 'Meta de Faturamento',
        value: formatCurrency(resumo.meta_faturamento),
        realizado: formatCurrency(resumo.realizado_faturamento),
        percentual: resumo.percentual_faturamento,
        icon: Target,
        trend: resumo.percentual_faturamento >= 100 ? 'up' : 'down',
      },
      {
        title: 'Atingimento Geral',
        value: formatPercentual(resumo.percentual_faturamento),
        description: `${formatCurrency(resumo.realizado_faturamento)} de ${formatCurrency(resumo.meta_faturamento)}`,
        percentual: resumo.percentual_faturamento,
        icon: TrendingUp,
        trend: resumo.percentual_faturamento >= 100 ? 'up' : 'down',
      },
      {
        title: 'Barbeiros Acima da Meta',
        value: `${resumo.barbeiros_acima_meta}/${resumo.total_barbeiros_com_meta}`,
        description: resumo.total_barbeiros_com_meta > 0
          ? `${((resumo.barbeiros_acima_meta / resumo.total_barbeiros_com_meta) * 100).toFixed(0)}% do time`
          : 'Nenhuma meta definida',
        percentual: resumo.total_barbeiros_com_meta > 0
          ? (resumo.barbeiros_acima_meta / resumo.total_barbeiros_com_meta) * 100
          : 0,
        icon: Users,
        trend: resumo.barbeiros_acima_meta > 0 ? 'up' : 'neutral',
      },
      {
        title: 'Ticket Médio',
        value: formatCurrency(resumo.ticket_medio_realizado),
        description: `Meta: ${formatCurrency(resumo.ticket_medio_meta)}`,
        percentual: resumo.percentual_ticket,
        icon: Award,
        trend: resumo.percentual_ticket >= 100 ? 'up' : 'down',
      },
    ];
  }, [resumo]);

  return (
    <div className="flex flex-col gap-6">
      {/* Mês Atual */}
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">{formatMesAno(mesAnoAtual)}</h2>
        <Badge variant="outline">Mês Atual</Badge>
      </div>

      {/* KPIs Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {isLoadingResumo ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-4 w-4" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-32 mb-2" />
                <Skeleton className="h-3 w-full" />
              </CardContent>
            </Card>
          ))
        ) : kpis ? (
          kpis.map((kpi) => {
            const Icon = kpi.icon;
            const TrendIcon = kpi.trend === 'up' ? ArrowUpRight : ArrowDownRight;
            const trendColor = kpi.trend === 'up' ? 'text-green-600' : 'text-red-600';

            return (
              <Card key={kpi.title}>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">{kpi.title}</CardTitle>
                  <Icon className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="flex items-baseline gap-2">
                    <div className="text-2xl font-bold">{kpi.value}</div>
                    {kpi.trend !== 'neutral' && (
                      <TrendIcon className={`h-4 w-4 ${trendColor}`} />
                    )}
                  </div>
                  {kpi.description && (
                    <p className="text-xs text-muted-foreground mt-1">
                      {kpi.description}
                    </p>
                  )}
                  {kpi.realizado && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Realizado: {kpi.realizado}
                    </p>
                  )}
                  <div className="mt-3">
                    <div className="flex justify-between text-xs mb-1">
                      <span>Progresso</span>
                      <span className={kpi.percentual >= 100 ? 'text-green-600 font-medium' : ''}>
                        {formatPercentual(kpi.percentual)}
                      </span>
                    </div>
                    <Progress
                      value={Math.min(kpi.percentual, 100)}
                      className="h-2"
                    />
                  </div>
                </CardContent>
              </Card>
            );
          })
        ) : (
          <Card className="col-span-full">
            <CardContent className="flex items-center justify-center py-10">
              <p className="text-muted-foreground">Nenhuma meta definida para este mês</p>
            </CardContent>
          </Card>
        )}
      </div>

      {/* Segunda Linha: Ranking + Histórico */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Ranking de Barbeiros */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              Ranking de Barbeiros
            </CardTitle>
            <CardDescription>
              Ordenado por percentual de atingimento da meta
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingRanking ? (
              <div className="space-y-4">
                {Array.from({ length: 5 }).map((_, i) => (
                  <div key={i} className="flex items-center gap-3">
                    <Skeleton className="h-8 w-8 rounded-full" />
                    <Skeleton className="h-4 flex-1" />
                    <Skeleton className="h-4 w-16" />
                  </div>
                ))}
              </div>
            ) : ranking && ranking.length > 0 ? (
              <div className="space-y-4">
                {ranking.slice(0, 5).map((barbeiro, index) => {
                  const medals = ['🥇', '🥈', '🥉'];
                  const medal = index < 3 ? medals[index] : null;

                  return (
                    <div
                      key={barbeiro.barbeiro_id}
                      className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted/50 transition-colors"
                    >
                      <div className="flex items-center justify-center w-8 h-8 rounded-full bg-muted text-sm font-medium">
                        {medal || barbeiro.posicao}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium truncate">{barbeiro.barbeiro_nome}</p>
                        <p className="text-xs text-muted-foreground">
                          {formatCurrency(barbeiro.realizado_total || 0)}
                        </p>
                      </div>
                      <div className="text-right">
                        <p className={`font-bold ${barbeiro.percentual_total >= 100 ? 'text-green-600' : 'text-muted-foreground'}`}>
                          {formatPercentual(barbeiro.percentual_total)}
                        </p>
                        {barbeiro.nivel_bonificacao !== 'NENHUM' && (
                          <Badge className={getNivelBonificacaoClass(barbeiro.nivel_bonificacao)}>
                            {formatNivelBonificacao(barbeiro.nivel_bonificacao)}
                          </Badge>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="flex items-center justify-center py-10">
                <p className="text-muted-foreground">Nenhuma meta de barbeiro cadastrada</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Histórico de Metas */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Histórico de Metas
            </CardTitle>
            <CardDescription>
              Últimos 6 meses
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingHistorico ? (
              <div className="space-y-4">
                {Array.from({ length: 6 }).map((_, i) => (
                  <div key={i} className="flex items-center gap-3">
                    <Skeleton className="h-4 w-16" />
                    <Skeleton className="h-2 flex-1" />
                    <Skeleton className="h-4 w-12" />
                  </div>
                ))}
              </div>
            ) : historico && historico.length > 0 ? (
              <div className="space-y-4">
                {historico.map((meta) => {
                  const percentual = parseFloat(meta.percentual) || 0;

                  return (
                    <div key={meta.id} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-sm font-medium">
                          {formatMesAnoShort(meta.mes_ano)}
                        </span>
                        <span className={`text-sm font-bold ${percentual >= 100 ? 'text-green-600' : 'text-muted-foreground'}`}>
                          {formatPercentual(percentual)}
                        </span>
                      </div>
                      <div className="relative">
                        <Progress
                          value={Math.min(percentual, 100)}
                          className="h-2"
                        />
                        {percentual >= 100 && (
                          <div
                            className="absolute top-0 left-0 h-full bg-green-500/20 rounded-full"
                            style={{ width: `${Math.min(percentual, 150) - 100}%`, marginLeft: '100%' }}
                          />
                        )}
                      </div>
                      <div className="flex justify-between text-xs text-muted-foreground">
                        <span>Meta: {formatCurrency(meta.meta_faturamento)}</span>
                        <span>Real: {formatCurrency(meta.realizado)}</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="flex items-center justify-center py-10">
                <p className="text-muted-foreground">Nenhum histórico disponível</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Card de Bônus Projetado */}
      {ranking && ranking.some((b) => b.nivel_bonificacao !== 'NENHUM') && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="h-5 w-5 text-yellow-500" />
              Bônus Projetado
            </CardTitle>
            <CardDescription>
              Baseado no atingimento atual das metas
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {ranking
                .filter((b) => b.nivel_bonificacao !== 'NENHUM')
                .map((barbeiro) => (
                  <div
                    key={barbeiro.barbeiro_id}
                    className="flex items-center justify-between p-3 rounded-lg border bg-muted/30"
                  >
                    <div>
                      <p className="font-medium">{barbeiro.barbeiro_nome}</p>
                      <Badge className={getNivelBonificacaoClass(barbeiro.nivel_bonificacao)}>
                        {formatNivelBonificacao(barbeiro.nivel_bonificacao)}
                      </Badge>
                    </div>
                    <div className="text-right">
                      <p className="text-lg font-bold text-green-600">
                        {formatCurrency(barbeiro.bonus_valor || 0)}
                      </p>
                      <p className="text-xs text-muted-foreground">bônus</p>
                    </div>
                  </div>
                ))}
            </div>
            <div className="mt-4 pt-4 border-t flex justify-between items-center">
              <span className="font-medium">Total de Bônus Projetado</span>
              <span className="text-xl font-bold text-green-600">
                {formatCurrency(
                  ranking
                    .filter((b) => b.nivel_bonificacao !== 'NENHUM')
                    .reduce((acc, b) => acc + (b.bonus_valor || 0), 0)
                )}
              </span>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
