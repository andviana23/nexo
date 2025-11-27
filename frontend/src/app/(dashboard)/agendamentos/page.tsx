'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Agendamentos - Layout AppBarber Style
 *
 * Layout idêntico ao AppBarber:
 * - Header compacto com navegação
 * - Profissionais no TOPO das colunas
 * - Slots de 10 minutos
 * - Sidebar direita com mini-calendário e opções
 */

import {
  Calendar,
  ChevronLeft,
  ChevronRight,
  Clock,
  List,
  Lock
} from 'lucide-react';
import { useCallback, useState } from 'react';

import { AgendaCalendar, AppointmentModal, BlockScheduleModal } from '@/components/appointments';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { useProfessionals } from '@/hooks/use-appointments';
import type { AppointmentModalState } from '@/types/appointment';

// =============================================================================
// TIPOS
// =============================================================================

type ViewType = 'day' | 'week' | 'month';

interface BlockModalState {
  isOpen: boolean;
  initialDate?: Date;
  initialProfessionalId?: string;
  initialStartTime?: string;
  initialEndTime?: string;
}

// =============================================================================
// MINI CALENDAR COMPONENT
// =============================================================================

function MiniCalendar({
  currentDate,
  onDateSelect,
}: {
  currentDate: Date;
  onDateSelect: (date: Date) => void;
}) {
  const [viewDate, setViewDate] = useState(currentDate);

  const monthNames = [
    'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
    'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro',
  ];

  const daysOfWeek = ['dom', 'seg', 'ter', 'qua', 'qui', 'sex', 'sáb'];

  const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDay = firstDay.getDay();

    const days: (number | null)[] = [];

    // Dias vazios antes do primeiro dia
    for (let i = 0; i < startingDay; i++) {
      days.push(null);
    }

    // Dias do mês
    for (let i = 1; i <= daysInMonth; i++) {
      days.push(i);
    }

    return days;
  };

  const days = getDaysInMonth(viewDate);

  const goToPrevMonth = () => {
    setViewDate((prev) => new Date(prev.getFullYear(), prev.getMonth() - 1, 1));
  };

  const goToNextMonth = () => {
    setViewDate((prev) => new Date(prev.getFullYear(), prev.getMonth() + 1, 1));
  };

  const isToday = (day: number) => {
    const today = new Date();
    return (
      day === today.getDate() &&
      viewDate.getMonth() === today.getMonth() &&
      viewDate.getFullYear() === today.getFullYear()
    );
  };

  const isSelected = (day: number) => {
    return (
      day === currentDate.getDate() &&
      viewDate.getMonth() === currentDate.getMonth() &&
      viewDate.getFullYear() === currentDate.getFullYear()
    );
  };

  const handleDayClick = (day: number) => {
    const newDate = new Date(viewDate.getFullYear(), viewDate.getMonth(), day);
    onDateSelect(newDate);
  };

  return (
    <div className="bg-card rounded-lg border border-border p-3">
      {/* Header do calendário */}
      <div className="flex items-center justify-between mb-2">
        <button
          onClick={goToPrevMonth}
          className="p-1 hover:bg-accent rounded text-muted-foreground transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
        </button>
        <span className="text-sm font-medium text-foreground">
          {monthNames[viewDate.getMonth()]} {viewDate.getFullYear()}
        </span>
        <button
          onClick={goToNextMonth}
          className="p-1 hover:bg-accent rounded text-muted-foreground transition-colors"
        >
          <ChevronRight className="h-4 w-4" />
        </button>
      </div>

      {/* Dias da semana */}
      <div className="grid grid-cols-7 gap-1 mb-1">
        {daysOfWeek.map((day) => (
          <div
            key={day}
            className="text-center text-xs text-muted-foreground font-medium py-1"
          >
            {day}
          </div>
        ))}
      </div>

      {/* Dias do mês */}
      <div className="grid grid-cols-7 gap-1">
        {days.map((day, index) => (
          <button
            key={index}
            onClick={() => day && handleDayClick(day)}
            disabled={!day}
            className={`
              text-center text-xs py-1.5 rounded-full transition-colors
              ${!day ? 'invisible' : ''}
              ${isSelected(day!) ? 'bg-primary text-primary-foreground' : ''}
              ${isToday(day!) && !isSelected(day!) ? 'border border-primary text-primary' : ''}
              ${!isSelected(day!) && !isToday(day!) ? 'text-foreground hover:bg-accent' : ''}
            `}
          >
            {day}
          </button>
        ))}
      </div>
    </div>
  );
}

// =============================================================================
// PAGE COMPONENT
// =============================================================================

export default function AgendamentosPage() {
  // Estado da data atual
  const [currentDate, setCurrentDate] = useState(() => new Date());

  // Estado da visualização
  const [viewType, setViewType] = useState<ViewType>('day');

  // Estado do modal
  const [modalState, setModalState] = useState<AppointmentModalState>({
    isOpen: false,
    mode: 'create',
  });

  // Estados para checkboxes da sidebar
  const [showBlockSchedule, setShowBlockSchedule] = useState(false);

  // Estado do modal de bloqueio
  const [blockModalState, setBlockModalState] = useState<BlockModalState>({
    isOpen: false,
  });

  // Buscar profissionais
  const { data: professionals = [] } = useProfessionals();

  // ==========================================================================
  // NAVEGAÇÃO DE DATAS
  // ==========================================================================

  const goToToday = useCallback(() => {
    setCurrentDate(new Date());
  }, []);

  const goToPrevious = useCallback(() => {
    setCurrentDate((prev) => {
      const newDate = new Date(prev);
      if (viewType === 'day') {
        newDate.setDate(newDate.getDate() - 1);
      } else if (viewType === 'week') {
        newDate.setDate(newDate.getDate() - 7);
      } else {
        newDate.setMonth(newDate.getMonth() - 1);
      }
      return newDate;
    });
  }, [viewType]);

  const goToNext = useCallback(() => {
    setCurrentDate((prev) => {
      const newDate = new Date(prev);
      if (viewType === 'day') {
        newDate.setDate(newDate.getDate() + 1);
      } else if (viewType === 'week') {
        newDate.setDate(newDate.getDate() + 7);
      } else {
        newDate.setMonth(newDate.getMonth() + 1);
      }
      return newDate;
    });
  }, [viewType]);

  // ==========================================================================
  // FORMATAÇÃO DE DATA
  // ==========================================================================

  const formatDateHeader = useCallback(() => {
    const dayNames = [
      'Domingo', 'Segunda', 'Terça', 'Quarta', 'Quinta', 'Sexta', 'Sábado',
    ];
    const monthNames = [
      'Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun',
      'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez',
    ];

    const dayName = dayNames[currentDate.getDay()];
    const day = currentDate.getDate();
    const month = monthNames[currentDate.getMonth()];
    const year = currentDate.getFullYear();

    return `${dayName}, ${day}/${month}/${year}`;
  }, [currentDate]);

  // ==========================================================================
  // HANDLERS DO CALENDÁRIO
  // ==========================================================================

  const handleEventClick = useCallback((state: AppointmentModalState) => {
    setModalState(state);
  }, []);

  const handleDateSelect = useCallback((state: AppointmentModalState) => {
    setModalState(state);
  }, []);

  const handleNewAppointment = useCallback(() => {
    setModalState({
      isOpen: true,
      mode: 'create',
      initialDate: currentDate,
    });
  }, [currentDate]);

  const handleCloseModal = useCallback(() => {
    setModalState((prev) => ({ ...prev, isOpen: false }));
  }, []);

  const handleMiniCalendarDateSelect = useCallback((date: Date) => {
    setCurrentDate(date);
  }, []);

  // Handler para abrir modal de bloqueio
  const handleOpenBlockModal = useCallback((options?: Partial<BlockModalState>) => {
    setBlockModalState({
      isOpen: true,
      initialDate: currentDate,
      ...options,
    });
  }, [currentDate]);

  // Handler para fechar modal de bloqueio
  const handleCloseBlockModal = useCallback(() => {
    setBlockModalState({ isOpen: false });
  }, []);

  // Handler para seleção de slot no calendário (modo bloqueio)
  const handleSlotSelect = useCallback((state: AppointmentModalState) => {
    if (showBlockSchedule) {
      // Modo bloqueio ativo - abrir modal de bloqueio
      handleOpenBlockModal({
        initialDate: state.initialDate,
        initialProfessionalId: state.initialProfessionalId,
        initialStartTime: state.initialDate?.toTimeString().slice(0, 5),
      });
    } else {
      // Modo normal - abrir modal de agendamento
      setModalState(state);
    }
  }, [showBlockSchedule, handleOpenBlockModal]);

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <div className="flex h-[calc(100vh-4rem)] flex-col bg-background">
      {/* ================================================================== */}
      {/* HEADER - Standard Shadcn/UI */}
      {/* ================================================================== */}
      <header className="flex h-16 shrink-0 items-center justify-between border-b px-6">
        {/* Lado Esquerdo - Navegação */}
        <div className="flex items-center gap-2">
          <div className="flex items-center rounded-md border bg-background shadow-sm">
            <Button
              variant="ghost"
              size="icon"
              onClick={goToPrevious}
              className="h-8 w-8 rounded-none rounded-l-md border-r"
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={goToToday}
              className="h-8 rounded-none px-3 text-xs font-medium"
            >
              Hoje
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={goToNext}
              className="h-8 w-8 rounded-none rounded-r-md border-l"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
          
          <span className="ml-4 text-lg font-semibold text-foreground">
            {formatDateHeader()}
          </span>
        </div>

        {/* Lado Direito - Ações e Views */}
        <div className="flex items-center gap-4">
          <div className="flex items-center rounded-md border bg-muted p-1">
            <Button
              variant={viewType === 'day' ? 'default' : 'ghost'}
              size="sm"
              onClick={() => setViewType('day')}
              className="h-7 px-3 text-xs"
            >
              Dia
            </Button>
            <Button
              variant={viewType === 'week' ? 'default' : 'ghost'}
              size="sm"
              onClick={() => setViewType('week')}
              className="h-7 px-3 text-xs"
            >
              Semana
            </Button>
            <Button
              variant={viewType === 'month' ? 'default' : 'ghost'}
              size="sm"
              onClick={() => setViewType('month')}
              className="h-7 px-3 text-xs"
            >
              Mês
            </Button>
          </div>

          <Button onClick={handleNewAppointment} size="sm">
            <Calendar className="mr-2 h-4 w-4" />
            Novo Agendamento
          </Button>
        </div>
      </header>

      {/* ================================================================== */}
      {/* ÁREA PRINCIPAL - Grid Layout */}
      {/* ================================================================== */}
      <div className="flex flex-1 overflow-hidden">
        {/* Calendário Principal */}
        <main className="flex-1 overflow-hidden bg-background p-4">
          <div className="h-full w-full overflow-hidden rounded-lg border bg-card shadow-sm">
            <AgendaCalendar
              currentDate={currentDate}
              viewType={viewType}
              professionalIds={[]}
              onEventClick={handleEventClick}
              onDateSelect={handleSlotSelect}
              editable={true}
              isBlockMode={showBlockSchedule}
            />
          </div>
        </main>

        {/* Sidebar Direita */}
        <aside className="w-80 shrink-0 border-l bg-muted/30 p-4 overflow-y-auto">
          <div className="space-y-6">
            {/* Mini Calendário */}
            <div className="rounded-lg border bg-card p-4 shadow-sm">
              <MiniCalendar
                currentDate={currentDate}
                onDateSelect={handleMiniCalendarDateSelect}
              />
            </div>

            {/* Ações Rápidas */}
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
                Ações Rápidas
              </h3>
              
              <div className="grid gap-2">
                <Button 
                  variant="outline" 
                  className="w-full justify-start"
                  onClick={() => handleOpenBlockModal()}
                >
                  <Lock className="mr-2 h-4 w-4 text-destructive" />
                  Bloquear Horário
                </Button>
                
                <Button variant="outline" className="w-full justify-start">
                  <Clock className="mr-2 h-4 w-4" />
                  Horários Disponíveis
                </Button>
                
                <Button variant="outline" className="w-full justify-start">
                  <List className="mr-2 h-4 w-4" />
                  Lista de Espera
                </Button>
              </div>
            </div>

            {/* Filtros / Opções */}
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
                Visualização
              </h3>
              
              <div className="flex items-center space-x-2 rounded-lg border bg-card p-3 shadow-sm">
                <Checkbox
                  id="block-schedule"
                  checked={showBlockSchedule}
                  onCheckedChange={(checked) => setShowBlockSchedule(!!checked)}
                />
                <label
                  htmlFor="block-schedule"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Modo de Bloqueio
                </label>
              </div>
              
              {showBlockSchedule && (
                <div className="rounded-md bg-destructive/10 p-3 text-xs text-destructive">
                  <p className="font-semibold">Modo Bloqueio Ativo</p>
                  <p>Clique na agenda para bloquear um horário.</p>
                </div>
              )}
            </div>
          </div>
        </aside>
      </div>

      {/* Modais */}
      <AppointmentModal state={modalState} onClose={handleCloseModal} />
      <BlockScheduleModal
        isOpen={blockModalState.isOpen}
        onClose={handleCloseBlockModal}
        initialDate={blockModalState.initialDate}
        initialProfessionalId={blockModalState.initialProfessionalId}
        initialStartTime={blockModalState.initialStartTime}
        initialEndTime={blockModalState.initialEndTime}
      />
    </div>
  );
}
