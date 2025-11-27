'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página da Lista da Vez (Barber Turn)
 *
 * @page /lista-da-vez
 * @description Sistema de fila giratória para distribuição justa de clientes
 * Conforme FLUXO_LISTA_DA_VEZ.md
 *
 * Regras:
 * - Menor pontuação = próximo da fila
 * - Cada atendimento incrementa +1 ponto
 * - Reset mensal zera todos os pontos
 * - Barbeiros podem ser pausados sem perder posição
 */

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from '@/components/ui/tooltip';
import {
    useAddBarberToTurn,
    useAvailableBarbers,
    useBarberTurnList,
    useRecordTurn,
    useRemoveBarberFromTurn,
    useResetTurnList,
    useToggleBarberStatus,
    useTurnHistory,
    useTurnHistorySummary,
} from '@/hooks/use-barber-turn';
import { cn } from '@/lib/utils';
import {
    formatLastTurn,
    formatMonthYear,
    formatPoints,
    getAvatarColor,
    getBarberStatus,
    getCurrentMonthYear,
    getInitials,
    getMonthOptions,
} from '@/services/barber-turn-service';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    AvailableBarberResponse,
    BarberTurnResponse,
    ListBarbersTurnResponse,
    ListHistorySummaryResponse,
    ListTurnHistoryResponse,
} from '@/types/barber-turn';
import {
    AlertTriangleIcon,
    BarChart3Icon,
    CalendarIcon,
    CheckIcon,
    CrownIcon,
    HistoryIcon,
    ListIcon,
    PauseIcon,
    PlayIcon,
    PlusIcon,
    RefreshCwIcon,
    Trash2Icon,
    UserIcon,
    UsersIcon,
} from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';

// =============================================================================
// TYPES
// =============================================================================

type ViewMode = 'list' | 'history';

interface AddBarberModalState {
  isOpen: boolean;
}

interface ResetModalState {
  isOpen: boolean;
  saveHistory: boolean;
}

interface ConfirmActionState {
  isOpen: boolean;
  type: 'remove' | 'record' | null;
  barberId: string | null;
  barberName: string | null;
}

// =============================================================================
// COMPONENT
// =============================================================================

export default function ListaDaVezPage() {
  const { setBreadcrumbs } = useBreadcrumbs();

  // View state
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [historyMonth, setHistoryMonth] = useState<string>(getCurrentMonthYear());

  // Modal states
  const [addModal, setAddModal] = useState<AddBarberModalState>({ isOpen: false });
  const [resetModal, setResetModal] = useState<ResetModalState>({
    isOpen: false,
    saveHistory: true,
  });
  const [confirmAction, setConfirmAction] = useState<ConfirmActionState>({
    isOpen: false,
    type: null,
    barberId: null,
    barberName: null,
  });

  // Queries
  const { data: turnData, isLoading, isError, refetch } = useBarberTurnList();
  const { data: availableBarbers, isLoading: isLoadingAvailable } = useAvailableBarbers();
  const { data: historyData, isLoading: isLoadingHistory } = useTurnHistory(historyMonth);
  const { data: summaryData } = useTurnHistorySummary();

  // Debug: Log available barbers
  useEffect(() => {
    if (availableBarbers) {
      console.log('[Lista da Vez] Barbeiros disponíveis:', availableBarbers);
    }
  }, [availableBarbers]);

  // Mutations
  const addBarber = useAddBarberToTurn();
  const recordTurn = useRecordTurn();
  const toggleStatus = useToggleBarberStatus();
  const removeBarber = useRemoveBarberFromTurn();
  const resetTurn = useResetTurnList();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Operações' },
      { label: 'Lista da Vez' },
    ]);
  }, [setBreadcrumbs]);

  // ===========================================================================
  // Handlers
  // ===========================================================================

  const handleAddBarber = useCallback(
    async (professionalId: string) => {
      await addBarber.mutateAsync(professionalId);
      setAddModal({ isOpen: false });
    },
    [addBarber]
  );

  const handleRecordTurn = useCallback(
    (barber: BarberTurnResponse) => {
      setConfirmAction({
        isOpen: true,
        type: 'record',
        barberId: barber.professional_id,
        barberName: barber.professional_name,
      });
    },
    []
  );

  const handleToggleStatus = useCallback(
    async (professionalId: string) => {
      await toggleStatus.mutateAsync(professionalId);
    },
    [toggleStatus]
  );

  const handleRemoveBarber = useCallback(
    (barber: BarberTurnResponse) => {
      setConfirmAction({
        isOpen: true,
        type: 'remove',
        barberId: barber.professional_id,
        barberName: barber.professional_name,
      });
    },
    []
  );

  const handleConfirmAction = useCallback(async () => {
    if (!confirmAction.barberId || !confirmAction.type) return;

    if (confirmAction.type === 'record') {
      await recordTurn.mutateAsync(confirmAction.barberId);
    } else if (confirmAction.type === 'remove') {
      await removeBarber.mutateAsync(confirmAction.barberId);
    }

    setConfirmAction({
      isOpen: false,
      type: null,
      barberId: null,
      barberName: null,
    });
  }, [confirmAction, recordTurn, removeBarber]);

  const handleReset = useCallback(async () => {
    await resetTurn.mutateAsync(resetModal.saveHistory);
    setResetModal({ isOpen: false, saveHistory: true });
  }, [resetTurn, resetModal.saveHistory]);

  // ===========================================================================
  // RENDER
  // ===========================================================================

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <ListIcon className="size-8" />
            Lista da Vez
          </h1>
          <p className="text-muted-foreground">
            Fila giratória para distribuição justa de clientes
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
          {/* Toggle View */}
          <div className="flex rounded-lg border">
            <Button
              variant={viewMode === 'list' ? 'secondary' : 'ghost'}
              size="sm"
              onClick={() => setViewMode('list')}
              className="rounded-r-none"
            >
              <UsersIcon className="mr-2 size-4" />
              Fila Atual
            </Button>
            <Button
              variant={viewMode === 'history' ? 'secondary' : 'ghost'}
              size="sm"
              onClick={() => setViewMode('history')}
              className="rounded-l-none"
            >
              <HistoryIcon className="mr-2 size-4" />
              Histórico
            </Button>
          </div>

          {viewMode === 'list' && (
            <>
              <Button
                variant="outline"
                onClick={() => setResetModal({ isOpen: true, saveHistory: true })}
              >
                <RefreshCwIcon className="mr-2 size-4" />
                Reset Mensal
              </Button>
              <Button onClick={() => setAddModal({ isOpen: true })}>
                <PlusIcon className="mr-2 size-4" />
                Adicionar Barbeiro
              </Button>
            </>
          )}
        </div>
      </div>

      {viewMode === 'list' ? (
        <div className="space-y-6">
          {/* Stats Cards - Minimalista */}
          {turnData && <MinimalStatsCards stats={turnData} />}

          {/* Grid Layout: Lista + Gráfico */}
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Ordem de Atendimento */}
            <Card className="border-none shadow-sm">
              <CardHeader className="pb-3">
                <CardTitle className="text-lg font-semibold">Ordem de Atendimento</CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                {isLoading ? (
                  <MinimalSkeleton />
                ) : isError ? (
                  <ErrorState onRetry={refetch} />
                ) : !turnData || turnData.barbers.length === 0 ? (
                  <MinimalEmptyState onAdd={() => setAddModal({ isOpen: true })} />
                ) : (
                  <MinimalTurnList
                    barbers={turnData.barbers}
                    nextBarberId={turnData.next_barber?.professional_id}
                    onRecordTurn={handleRecordTurn}
                    isRecording={recordTurn.isPending}
                  />
                )}
              </CardContent>
            </Card>

            {/* Distribuição de Atendimentos */}
            <Card className="border-none shadow-sm">
              <CardHeader className="pb-3">
                <CardTitle className="text-lg font-semibold">Distribuição de Atendimentos</CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                {isLoading ? (
                  <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-2 border-primary border-t-transparent"></div>
                  </div>
                ) : !turnData || turnData.barbers.length === 0 ? (
                  <div className="flex items-center justify-center h-64 text-sm text-muted-foreground">
                    Nenhum dado disponível
                  </div>
                ) : (
                  <MinimalChart barbers={turnData.barbers} />
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      ) : (
        /* History View */
        <HistoryView
          historyData={historyData}
          summaryData={summaryData}
          selectedMonth={historyMonth}
          onMonthChange={setHistoryMonth}
          isLoading={isLoadingHistory}
        />
      )}

      {/* Add Barber Modal */}
      <AddBarberModal
        isOpen={addModal.isOpen}
        onClose={() => setAddModal({ isOpen: false })}
        onAdd={handleAddBarber}
        availableBarbers={availableBarbers?.barbers || []}
        isLoading={isLoadingAvailable || addBarber.isPending}
      />

      {/* Reset Modal */}
      <ResetModal
        isOpen={resetModal.isOpen}
        saveHistory={resetModal.saveHistory}
        onSaveHistoryChange={(value) => setResetModal((prev) => ({ ...prev, saveHistory: value }))}
        onClose={() => setResetModal({ isOpen: false, saveHistory: true })}
        onConfirm={handleReset}
        stats={turnData?.stats}
        isLoading={resetTurn.isPending}
      />

      {/* Confirm Action Modal */}
      <ConfirmActionModal
        isOpen={confirmAction.isOpen}
        type={confirmAction.type}
        barberName={confirmAction.barberName}
        onClose={() =>
          setConfirmAction({
            isOpen: false,
            type: null,
            barberId: null,
            barberName: null,
          })
        }
        onConfirm={handleConfirmAction}
        isLoading={recordTurn.isPending || removeBarber.isPending}
      />
    </div>
  );
}

// =============================================================================
// STATS CARDS
// =============================================================================

function StatsCards({ stats }: { stats: ListBarbersTurnResponse }) {
  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Barbeiros Ativos</CardTitle>
          <UsersIcon className="size-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">{stats.stats.total_ativos}</div>
          <p className="text-xs text-muted-foreground">
            na fila de atendimento
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Barbeiros Pausados</CardTitle>
          <PauseIcon className="size-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-amber-600">{stats.stats.total_pausados}</div>
          <p className="text-xs text-muted-foreground">
            temporariamente fora
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Atendimentos no Mês</CardTitle>
          <BarChart3Icon className="size-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.stats.total_pontos_mes}</div>
          <p className="text-xs text-muted-foreground">
            total acumulado
          </p>
        </CardContent>
      </Card>

      <Card className="bg-linear-to-br from-amber-50 to-amber-100 dark:from-amber-950 dark:to-amber-900 border-amber-200 dark:border-amber-800">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Próximo da Fila</CardTitle>
          <CrownIcon className="size-4 text-amber-600" />
        </CardHeader>
        <CardContent>
          {stats.next_barber ? (
            <>
              <div className="text-2xl font-bold text-amber-700 dark:text-amber-400 truncate">
                {stats.next_barber.professional_name.split(' ')[0]}
              </div>
              <p className="text-xs text-amber-600 dark:text-amber-500">
                {formatPoints(stats.next_barber.current_points)}
              </p>
            </>
          ) : (
            <>
              <div className="text-2xl font-bold text-muted-foreground">—</div>
              <p className="text-xs text-muted-foreground">
                nenhum barbeiro ativo
              </p>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// NEW STATS CARDS (Seguindo design da imagem)
// =============================================================================

function NewStatsCards({ stats }: { stats: ListBarbersTurnResponse }) {
  const totalBarbers = stats.stats.total_ativos + stats.stats.total_pausados;
  const avgAtendimentos = totalBarbers > 0 
    ? (stats.stats.total_pontos_mes / totalBarbers).toFixed(1) 
    : '0.0';

  return (
    <div className="grid gap-4 md:grid-cols-4">
      {/* Total de Barbeiros */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-blue-100 dark:bg-blue-900/20 p-2">
              <UsersIcon className="size-5 text-blue-600 dark:text-blue-400" />
            </div>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total de Barbeiros
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{totalBarbers}</div>
        </CardContent>
      </Card>

      {/* Total de Atendimentos */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-green-100 dark:bg-green-900/20 p-2">
              <BarChart3Icon className="size-5 text-green-600 dark:text-green-400" />
            </div>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total de Atendimentos
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{stats.stats.total_pontos_mes}</div>
        </CardContent>
      </Card>

      {/* Média de Atendimentos */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-purple-100 dark:bg-purple-900/20 p-2">
              <CrownIcon className="size-5 text-purple-600 dark:text-purple-400" />
            </div>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Média de Atendimentos
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{avgAtendimentos}</div>
        </CardContent>
      </Card>

      {/* Última Atualização */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="flex items-center gap-2">
            <div className="rounded-lg bg-orange-100 dark:bg-orange-900/20 p-2">
              <CalendarIcon className="size-5 text-orange-600 dark:text-orange-400" />
            </div>
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Última Atualização
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-xl font-bold">
            {new Date().toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// NEW TURN LIST (Seguindo design da imagem)
// =============================================================================

interface NewTurnListProps {
  barbers: BarberTurnResponse[];
  nextBarberId?: string;
  onRecordTurn: (barber: BarberTurnResponse) => void;
  isRecording?: boolean;
}

function NewTurnList({ barbers, nextBarberId, onRecordTurn, isRecording }: NewTurnListProps) {
  const totalPoints = barbers.reduce((sum, b) => sum + b.current_points, 0);
  const COLORS = [
    'bg-blue-500',
    'bg-red-500', 
    'bg-orange-500',
    'bg-green-500',
    'bg-purple-500'
  ];

  return (
    <div className="space-y-3">
      {barbers.map((barber, index) => {
        const isNext = barber.professional_id === nextBarberId;
        const percentage = totalPoints > 0 ? ((barber.current_points / totalPoints) * 100).toFixed(1) : '0.0';
        const colorClass = COLORS[index % COLORS.length];

        return (
          <div
            key={barber.professional_id}
            className={cn(
              'flex items-center gap-4 p-4 rounded-lg border transition-all',
              isNext && 'bg-green-50 dark:bg-green-950/20 border-green-200 dark:border-green-800'
            )}
          >
            {/* Posição */}
            <div className="flex flex-col items-center min-w-[60px]">
              <div className={cn(
                'flex items-center justify-center size-8 rounded-full text-white font-bold text-sm',
                isNext ? 'bg-green-600' : 'bg-muted-foreground'
              )}>
                {index + 1}°
              </div>
              {isNext && (
                <span className="text-[10px] text-green-600 dark:text-green-400 font-medium mt-1 uppercase">
                  Próximo na vez
                </span>
              )}
            </div>

            {/* Barra de cor (indicador visual) */}
            <div className={cn('w-1 h-12 rounded-full', colorClass)} />

            {/* Nome */}
            <div className="flex-1 min-w-0">
              <p className="font-semibold truncate text-base">{barber.professional_name}</p>
            </div>

            {/* Atendimentos */}
            <div className="text-center min-w-20">
              <div className="text-2xl font-bold">{barber.current_points}</div>
              <div className="text-xs text-muted-foreground">atendimentos</div>
            </div>

            {/* Participação */}
            <div className="text-center min-w-20">
              <div className="text-lg font-semibold">{percentage}%</div>
              <div className="text-xs text-muted-foreground">participação</div>
            </div>

            {/* Botão +1 */}
            <Button
              size="sm"
              onClick={() => onRecordTurn(barber)}
              disabled={isRecording || !barber.is_active}
              className={cn(
                'min-w-[60px]',
                isNext && 'bg-green-600 hover:bg-green-700'
              )}
            >
              +1
            </Button>
          </div>
        );
      })}
    </div>
  );
}

// =============================================================================
// DISTRIBUTION CHART (Gráfico de pizza)
// =============================================================================

interface DistributionChartProps {
  barbers: BarberTurnResponse[];
}

function DistributionChart({ barbers }: DistributionChartProps) {
  const totalPoints = barbers.reduce((sum, b) => sum + b.current_points, 0);
  const COLORS = [
    { bg: '#3B82F6', name: 'blue' },    // blue-500
    { bg: '#EF4444', name: 'red' },     // red-500
    { bg: '#F97316', name: 'orange' },  // orange-500
    { bg: '#10B981', name: 'green' },   // green-500
    { bg: '#A855F7', name: 'purple' },  // purple-500
  ];

  // Calcular ângulos para cada fatia do gráfico
  let currentAngle = 0;
  const slices = barbers.map((barber, index) => {
    const percentage = totalPoints > 0 ? (barber.current_points / totalPoints) * 100 : 0;
    const angle = (percentage / 100) * 360;
    const slice = {
      barber,
      percentage,
      startAngle: currentAngle,
      endAngle: currentAngle + angle,
      color: COLORS[index % COLORS.length],
    };
    currentAngle += angle;
    return slice;
  });

  return (
    <div className="space-y-6">
      {/* Gráfico SVG de Pizza */}
      <div className="flex items-center justify-center">
        <svg viewBox="0 0 200 200" className="w-64 h-64">
          <circle cx="100" cy="100" r="100" fill="#f3f4f6" className="dark:fill-slate-800" />
          {slices.map((slice, index) => {
            const startAngleRad = (slice.startAngle - 90) * (Math.PI / 180);
            const endAngleRad = (slice.endAngle - 90) * (Math.PI / 180);
            
            const x1 = 100 + 100 * Math.cos(startAngleRad);
            const y1 = 100 + 100 * Math.sin(startAngleRad);
            const x2 = 100 + 100 * Math.cos(endAngleRad);
            const y2 = 100 + 100 * Math.sin(endAngleRad);
            
            const largeArc = slice.percentage > 50 ? 1 : 0;
            
            return (
              <path
                key={index}
                d={`M 100 100 L ${x1} ${y1} A 100 100 0 ${largeArc} 1 ${x2} ${y2} Z`}
                fill={slice.color.bg}
                className="transition-opacity hover:opacity-80"
              />
            );
          })}
          {/* Círculo branco central (efeito donut) */}
          <circle cx="100" cy="100" r="60" fill="white" className="dark:fill-slate-950" />
        </svg>
      </div>

      {/* Legenda */}
      <div className="space-y-2">
        {slices.map((slice, index) => (
          <div key={index} className="flex items-center justify-between py-2 border-b last:border-0">
            <div className="flex items-center gap-3">
              <div
                className="w-4 h-4 rounded-full"
                style={{ backgroundColor: slice.color.bg }}
              />
              <span className="font-medium">{slice.barber.professional_name}</span>
            </div>
            <span className="text-sm text-muted-foreground">
              {slice.percentage.toFixed(1)}%
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

// =============================================================================
// MINIMAL UI COMPONENTS (Design Minimalista)
// =============================================================================

function MinimalStatsCards({ stats }: { stats: ListBarbersTurnResponse }) {
  const totalBarbers = stats.stats.total_ativos + stats.stats.total_pausados;
  const barbeirosNaLista = stats.barbers.length;
  const avgAtendimentos = barbeirosNaLista > 0 
    ? (stats.stats.total_pontos_mes / barbeirosNaLista).toFixed(1) 
    : '0.0';

  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card className="border-none shadow-sm hover:shadow-md transition-shadow">
        <CardContent className="p-6">
          <div className="flex items-center gap-3">
            <div className="rounded-xl bg-blue-50 dark:bg-blue-950/20 p-3">
              <UsersIcon className="size-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Barbeiros na Lista</p>
              <p className="text-2xl font-bold">{barbeirosNaLista}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="border-none shadow-sm hover:shadow-md transition-shadow">
        <CardContent className="p-6">
          <div className="flex items-center gap-3">
            <div className="rounded-xl bg-green-50 dark:bg-green-950/20 p-3">
              <BarChart3Icon className="size-5 text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Total de Atendimentos</p>
              <p className="text-2xl font-bold">{stats.stats.total_pontos_mes}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="border-none shadow-sm hover:shadow-md transition-shadow">
        <CardContent className="p-6">
          <div className="flex items-center gap-3">
            <div className="rounded-xl bg-purple-50 dark:bg-purple-950/20 p-3">
              <CrownIcon className="size-5 text-purple-600 dark:text-purple-400" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Média de Atendimentos</p>
              <p className="text-2xl font-bold">{avgAtendimentos}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="border-none shadow-sm hover:shadow-md transition-shadow">
        <CardContent className="p-6">
          <div className="flex items-center gap-3">
            <div className="rounded-xl bg-orange-50 dark:bg-orange-950/20 p-3">
              <CalendarIcon className="size-5 text-orange-600 dark:text-orange-400" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Última Atualização</p>
              <p className="text-xl font-bold">
                {new Date().toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function MinimalTurnList({ barbers, nextBarberId, onRecordTurn, isRecording }: NewTurnListProps) {
  const totalPoints = barbers.reduce((sum, b) => sum + b.current_points, 0);
  const COLORS = [
    'bg-blue-500',
    'bg-red-500', 
    'bg-orange-500',
    'bg-green-500',
    'bg-purple-500'
  ];

  return (
    <div className="space-y-2">
      {barbers.map((barber, index) => {
        const isNext = barber.professional_id === nextBarberId;
        const percentage = totalPoints > 0 ? ((barber.current_points / totalPoints) * 100).toFixed(1) : '0.0';
        const colorClass = COLORS[index % COLORS.length];

        return (
          <div
            key={barber.professional_id}
            className={cn(
              'flex items-center gap-3 p-3 rounded-xl transition-all hover:shadow-sm',
              isNext 
                ? 'bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-800' 
                : 'bg-muted/30 hover:bg-muted/50'
            )}
          >
            {/* Posição + Indicador */}
            <div className="flex items-center gap-2">
              <div className={cn(
                'flex items-center justify-center size-7 rounded-lg text-white font-bold text-xs',
                isNext ? 'bg-green-600' : 'bg-muted-foreground/60'
              )}>
                {index + 1}°
              </div>
              <div className={cn('w-1 h-8 rounded-full', colorClass)} />
            </div>

            {/* Nome */}
            <div className="flex-1 min-w-0">
              <p className="font-medium truncate">{barber.professional_name}</p>
              {isNext && (
                <p className="text-[10px] text-green-600 dark:text-green-400 font-medium uppercase">
                  Próximo na vez
                </p>
              )}
            </div>

            {/* Estatísticas compactas */}
            <div className="flex items-center gap-4 text-sm">
              <div className="text-center">
                <div className="font-bold">{barber.current_points}</div>
                <div className="text-[10px] text-muted-foreground">atend.</div>
              </div>
              <div className="text-center">
                <div className="font-semibold">{percentage}%</div>
                <div className="text-[10px] text-muted-foreground">part.</div>
              </div>
            </div>

            {/* Botão +1 minimalista */}
            <Button
              size="sm"
              onClick={() => onRecordTurn(barber)}
              disabled={isRecording || !barber.is_active}
              className={cn(
                'h-8 w-12 rounded-lg font-bold',
                isNext && 'bg-green-600 hover:bg-green-700'
              )}
            >
              +1
            </Button>
          </div>
        );
      })}
    </div>
  );
}

function MinimalChart({ barbers }: DistributionChartProps) {
  const totalPoints = barbers.reduce((sum, b) => sum + b.current_points, 0);
  const COLORS = [
    { bg: '#3B82F6', name: 'blue' },
    { bg: '#EF4444', name: 'red' },
    { bg: '#F97316', name: 'orange' },
    { bg: '#10B981', name: 'green' },
    { bg: '#A855F7', name: 'purple' },
  ];

  let currentAngle = 0;
  const slices = barbers.map((barber, index) => {
    const percentage = totalPoints > 0 ? (barber.current_points / totalPoints) * 100 : 0;
    const angle = (percentage / 100) * 360;
    const slice = {
      barber,
      percentage,
      startAngle: currentAngle,
      endAngle: currentAngle + angle,
      color: COLORS[index % COLORS.length],
    };
    currentAngle += angle;
    return slice;
  });

  return (
    <div className="space-y-6">
      {/* Gráfico SVG */}
      <div className="flex items-center justify-center">
        <svg viewBox="0 0 200 200" className="w-60 h-60">
          <circle cx="100" cy="100" r="100" fill="#f3f4f6" className="dark:fill-slate-800" />
          {slices.map((slice, index) => {
            const startAngleRad = (slice.startAngle - 90) * (Math.PI / 180);
            const endAngleRad = (slice.endAngle - 90) * (Math.PI / 180);
            
            const x1 = 100 + 100 * Math.cos(startAngleRad);
            const y1 = 100 + 100 * Math.sin(startAngleRad);
            const x2 = 100 + 100 * Math.cos(endAngleRad);
            const y2 = 100 + 100 * Math.sin(endAngleRad);
            
            const largeArc = slice.percentage > 50 ? 1 : 0;
            
            return (
              <path
                key={index}
                d={`M 100 100 L ${x1} ${y1} A 100 100 0 ${largeArc} 1 ${x2} ${y2} Z`}
                fill={slice.color.bg}
                className="transition-all hover:opacity-80 cursor-pointer"
              />
            );
          })}
          <circle cx="100" cy="100" r="60" fill="white" className="dark:fill-slate-950" />
        </svg>
      </div>

      {/* Legenda minimalista */}
      <div className="space-y-1">
        {slices.map((slice, index) => (
          <div key={index} className="flex items-center justify-between px-2 py-1.5 rounded-lg hover:bg-muted/50 transition-colors">
            <div className="flex items-center gap-2">
              <div
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: slice.color.bg }}
              />
              <span className="text-sm font-medium">{slice.barber.professional_name}</span>
            </div>
            <span className="text-sm text-muted-foreground font-semibold">
              {slice.percentage.toFixed(1)}%
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

function MinimalEmptyState({ onAdd }: { onAdd: () => void }) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="rounded-2xl bg-linear-to-br from-muted/50 to-muted/30 p-6 mb-6">
        <UsersIcon className="size-12 text-muted-foreground/40" />
      </div>
      <h3 className="text-lg font-semibold mb-2">Lista vazia</h3>
      <p className="text-sm text-muted-foreground max-w-sm mb-6">
        Adicione barbeiros para iniciar a Lista da Vez e começar a distribuir atendimentos de forma justa.
      </p>
      <Button onClick={onAdd} size="lg" className="rounded-xl shadow-sm">
        <PlusIcon className="mr-2 size-5" />
        Adicionar Barbeiro
      </Button>
    </div>
  );
}

function MinimalSkeleton() {
  return (
    <div className="space-y-2">
      {[1, 2, 3].map((i) => (
        <div key={i} className="flex items-center gap-3 p-3 rounded-xl bg-muted/30">
          <Skeleton className="size-7 rounded-lg" />
          <Skeleton className="w-1 h-8 rounded-full" />
          <Skeleton className="h-4 flex-1" />
          <Skeleton className="h-8 w-16" />
          <Skeleton className="h-8 w-12 rounded-lg" />
        </div>
      ))}
    </div>
  );
}

// =============================================================================
// TURN LIST (ORIGINAL - mantido para compatibilidade)
// =============================================================================

interface TurnListProps {
  barbers: BarberTurnResponse[];
  nextBarberId?: string;
  onRecordTurn: (barber: BarberTurnResponse) => void;
  onToggleStatus: (professionalId: string) => void;
  onRemove: (barber: BarberTurnResponse) => void;
  isRecording?: boolean;
  isToggling?: boolean;
}

function TurnList({
  barbers,
  nextBarberId,
  onRecordTurn,
  onToggleStatus,
  onRemove,
  isRecording,
  isToggling,
}: TurnListProps) {
  return (
    <TooltipProvider>
      <div className="space-y-3">
        {barbers.map((barber, index) => {
          const isNext = barber.professional_id === nextBarberId && barber.is_active;
          const status = getBarberStatus(barber);

          return (
            <div
              key={barber.id}
              className={cn(
                'flex items-center gap-4 p-4 rounded-lg border transition-all',
                isNext && 'bg-amber-50 dark:bg-amber-950 border-amber-300 dark:border-amber-700 shadow-sm',
                !barber.is_active && 'opacity-60 bg-muted/50'
              )}
            >
              {/* Position */}
              <div
                className={cn(
                  'flex size-10 items-center justify-center rounded-full font-bold text-lg',
                  isNext
                    ? 'bg-amber-500 text-white'
                    : 'bg-muted text-muted-foreground'
                )}
              >
                {index + 1}
              </div>

              {/* Avatar */}
              <div
                className={cn(
                  'flex size-12 items-center justify-center rounded-full text-white font-medium',
                  !barber.is_active ? 'bg-muted-foreground' : getAvatarColor(barber.professional_name)
                )}
              >
                {getInitials(barber.professional_name)}
              </div>

              {/* Info */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <p className="font-semibold truncate">{barber.professional_name}</p>
                  {isNext && (
                    <Badge className="bg-amber-500 hover:bg-amber-600">
                      <CrownIcon className="mr-1 size-3" />
                      Próximo
                    </Badge>
                  )}
                  <Badge variant={status.color === 'success' ? 'default' : 'secondary'}>
                    {status.label}
                  </Badge>
                </div>
                <div className="flex items-center gap-4 text-sm text-muted-foreground mt-1">
                  <span className="font-medium">{barber.current_points} pts</span>
                  <span>•</span>
                  <span>{formatLastTurn(barber.last_turn_at)}</span>
                </div>
              </div>

              {/* Actions */}
              <div className="flex items-center gap-2">
                {barber.is_active && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        size="sm"
                        onClick={() => onRecordTurn(barber)}
                        disabled={isRecording}
                        className={cn(
                          isNext && 'bg-amber-500 hover:bg-amber-600'
                        )}
                      >
                        <CheckIcon className="mr-2 size-4" />
                        Atendeu
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Registrar atendimento (+1 ponto)</TooltipContent>
                  </Tooltip>
                )}

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => onToggleStatus(barber.professional_id)}
                      disabled={isToggling}
                    >
                      {barber.is_active ? (
                        <PauseIcon className="size-4" />
                      ) : (
                        <PlayIcon className="size-4" />
                      )}
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {barber.is_active ? 'Pausar barbeiro' : 'Reativar barbeiro'}
                  </TooltipContent>
                </Tooltip>

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => onRemove(barber)}
                      className="text-destructive hover:text-destructive"
                    >
                      <Trash2Icon className="size-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Remover da lista</TooltipContent>
                </Tooltip>
              </div>
            </div>
          );
        })}
      </div>
    </TooltipProvider>
  );
}

// =============================================================================
// HISTORY VIEW
// =============================================================================

interface HistoryViewProps {
  historyData?: ListTurnHistoryResponse;
  summaryData?: ListHistorySummaryResponse;
  selectedMonth: string;
  onMonthChange: (month: string) => void;
  isLoading: boolean;
}

function HistoryView({
  historyData,
  summaryData,
  selectedMonth,
  onMonthChange,
  isLoading,
}: HistoryViewProps) {
  const monthOptions = getMonthOptions();

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      {summaryData && summaryData.summary.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3Icon className="size-5" />
              Resumo dos Últimos Meses
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
              {summaryData.summary.slice(0, 6).map((item) => (
                <div
                  key={item.month_year}
                  className={cn(
                    'p-4 rounded-lg border text-center cursor-pointer transition-all hover:border-primary',
                    item.month_year === selectedMonth && 'border-primary bg-primary/5'
                  )}
                  onClick={() => onMonthChange(item.month_year)}
                >
                  <p className="text-sm font-medium text-muted-foreground">
                    {formatMonthYear(item.month_year)}
                  </p>
                  <p className="text-2xl font-bold mt-1">{item.total_atendimentos}</p>
                  <p className="text-xs text-muted-foreground">atendimentos</p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* History Detail */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <CalendarIcon className="size-5" />
                Histórico de Atendimentos
              </CardTitle>
              <CardDescription>
                Detalhamento por barbeiro no mês selecionado
              </CardDescription>
            </div>
            <Select value={selectedMonth} onValueChange={onMonthChange}>
              <SelectTrigger className="w-[180px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {monthOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <TurnListSkeleton />
          ) : historyData?.history.length === 0 ? (
            <div className="text-center py-12">
              <HistoryIcon className="mx-auto size-12 text-muted-foreground/50" />
              <h3 className="mt-4 text-lg font-semibold">Sem histórico</h3>
              <p className="text-muted-foreground mt-1">
                Nenhum registro encontrado para {formatMonthYear(selectedMonth)}
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {historyData?.history.map((item, index) => (
                <div
                  key={item.id}
                  className="flex items-center gap-4 p-4 rounded-lg border"
                >
                  <div className="flex size-10 items-center justify-center rounded-full bg-muted font-bold">
                    {index + 1}
                  </div>
                  <div className="flex-1">
                    <p className="font-semibold">{item.professional_name}</p>
                    <p className="text-sm text-muted-foreground">
                      {item.total_turns} atendimentos • {item.final_points} pontos finais
                    </p>
                  </div>
                  <Badge variant="outline">
                    {formatMonthYear(item.month_year)}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// MODALS
// =============================================================================

interface AddBarberModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAdd: (professionalId: string) => void;
  availableBarbers: AvailableBarberResponse[];
  isLoading?: boolean;
}

function AddBarberModal({
  isOpen,
  onClose,
  onAdd,
  availableBarbers,
  isLoading,
}: AddBarberModalProps) {
  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="text-xl">Adicionar Barbeiro</DialogTitle>
          <DialogDescription>
            Selecione um barbeiro para adicionar à Lista da Vez
          </DialogDescription>
        </DialogHeader>

        <div className="max-h-[400px] overflow-y-auto pr-1">
          {isLoading ? (
            <div className="space-y-2 p-1">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-14 w-full rounded-xl" />
              ))}
            </div>
          ) : availableBarbers.length === 0 ? (
            <div className="text-center py-12">
              <div className="rounded-2xl bg-muted/50 p-4 w-fit mx-auto mb-4">
                <UserIcon className="size-10 text-muted-foreground/50" />
              </div>
              <p className="text-sm font-medium text-muted-foreground mb-1">
                Nenhum barbeiro disponível
              </p>
              <p className="text-xs text-muted-foreground/70">
                Todos os barbeiros já estão na lista
              </p>
            </div>
          ) : (
            <div className="space-y-2 p-1">
              {availableBarbers.map((barber) => (
                <button
                  key={barber.id}
                  onClick={() => {
                    onAdd(barber.id);
                    onClose();
                  }}
                  className="w-full flex items-center gap-3 p-3 rounded-xl border-none bg-muted/30 hover:bg-muted/50 transition-all hover:shadow-sm text-left"
                >
                  <div
                    className={cn(
                      'flex size-10 items-center justify-center rounded-xl text-white font-medium text-sm',
                      getAvatarColor(barber.nome)
                    )}
                  >
                    {getInitials(barber.nome)}
                  </div>
                  <div className="flex-1">
                    <p className="font-semibold text-sm">{barber.nome}</p>
                    <p className="text-xs text-muted-foreground capitalize">
                      {barber.status}
                    </p>
                  </div>
                  <div className="rounded-lg bg-primary/10 p-2">
                    <PlusIcon className="size-4 text-primary" />
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

interface ResetModalProps {
  isOpen: boolean;
  saveHistory: boolean;
  onSaveHistoryChange: (value: boolean) => void;
  onClose: () => void;
  onConfirm: () => void;
  stats?: { total_ativos: number; total_pontos_mes: number };
  isLoading?: boolean;
}

function ResetModal({
  isOpen,
  saveHistory,
  onSaveHistoryChange,
  onClose,
  onConfirm,
  stats,
  isLoading,
}: ResetModalProps) {
  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-amber-600">
            <AlertTriangleIcon className="size-5" />
            Reset Mensal
          </DialogTitle>
          <DialogDescription>
            Esta ação irá zerar os pontos de todos os barbeiros na Lista da Vez.
          </DialogDescription>
        </DialogHeader>

        {stats && (
          <div className="rounded-lg bg-muted p-4 space-y-2">
            <p className="text-sm">
              <strong>{stats.total_ativos}</strong> barbeiros ativos serão resetados
            </p>
            <p className="text-sm">
              <strong>{stats.total_pontos_mes}</strong> pontos serão zerados
            </p>
          </div>
        )}

        <div className="flex items-center justify-between py-4 border-y">
          <div>
            <p className="font-medium">Salvar histórico</p>
            <p className="text-sm text-muted-foreground">
              Preserva registro dos atendimentos do mês
            </p>
          </div>
          <label className="relative inline-flex items-center cursor-pointer">
            <input
              type="checkbox"
              className="sr-only peer"
              checked={saveHistory}
              onChange={(e) => onSaveHistoryChange(e.target.checked)}
            />
            <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/30 dark:peer-focus:ring-primary/50 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:start-0.5 after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
          </label>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={onConfirm}
            disabled={isLoading}
          >
            {isLoading ? 'Processando...' : 'Confirmar Reset'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

interface ConfirmActionModalProps {
  isOpen: boolean;
  type: 'remove' | 'record' | null;
  barberName: string | null;
  onClose: () => void;
  onConfirm: () => void;
  isLoading?: boolean;
}

function ConfirmActionModal({
  isOpen,
  type,
  barberName,
  onClose,
  onConfirm,
  isLoading,
}: ConfirmActionModalProps) {
  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {type === 'record'
              ? 'Confirmar Atendimento'
              : 'Remover da Lista'}
          </DialogTitle>
          <DialogDescription>
            {type === 'record'
              ? `Registrar atendimento para ${barberName}? Os pontos serão incrementados em +1.`
              : `Tem certeza que deseja remover ${barberName} da Lista da Vez? Os pontos acumulados serão perdidos.`}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancelar
          </Button>
          <Button
            variant={type === 'remove' ? 'destructive' : 'default'}
            onClick={onConfirm}
            disabled={isLoading}
          >
            {isLoading
              ? 'Processando...'
              : type === 'record'
              ? 'Confirmar'
              : 'Remover'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// STATES
// =============================================================================

function TurnListSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="flex items-center gap-4 p-4 rounded-lg border">
          <Skeleton className="size-10 rounded-full" />
          <Skeleton className="size-12 rounded-full" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-3 w-24" />
          </div>
          <Skeleton className="h-9 w-24" />
        </div>
      ))}
    </div>
  );
}

function EmptyState({ onAdd }: { onAdd: () => void }) {
  return (
    <div className="text-center py-12">
      <ListIcon className="mx-auto size-12 text-muted-foreground/50" />
      <h3 className="mt-4 text-lg font-semibold">Lista vazia</h3>
      <p className="text-muted-foreground mt-1">
        Adicione barbeiros para iniciar a Lista da Vez
      </p>
      <Button onClick={onAdd} className="mt-4">
        <PlusIcon className="mr-2 size-4" />
        Adicionar Barbeiro
      </Button>
    </div>
  );
}

function ErrorState({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="text-center py-12">
      <AlertTriangleIcon className="mx-auto size-12 text-destructive/50" />
      <h3 className="mt-4 text-lg font-semibold text-destructive">Erro ao carregar</h3>
      <p className="text-muted-foreground mt-1">
        Não foi possível carregar a Lista da Vez
      </p>
      <Button variant="outline" onClick={onRetry} className="mt-4">
        Tentar novamente
      </Button>
    </div>
  );
}
