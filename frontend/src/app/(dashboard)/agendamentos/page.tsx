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

import { format } from 'date-fns';
import {
    Calendar,
    CalendarDays,
    ChevronLeft,
    ChevronRight,
    Clock,
    List,
    Lock
} from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';

import { CommandModal } from '@/components/agendamentos/CommandModal';
import {
    AgendaCalendar,
    AppointmentCardWithCommand,
    AppointmentContextMenu,
    AppointmentModal,
    BlockScheduleModal
} from '@/components/appointments';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useAppointments, useUpdateAppointmentStatus } from '@/hooks/use-appointments';
import { useCreateCommandFromAppointment } from '@/hooks/use-commands';
import type { AppointmentModalState, AppointmentResponse } from '@/types/appointment';
import { toast } from 'sonner';

// =============================================================================
// TIPOS
// =============================================================================

type ViewType = 'day' | 'week' | 'month';
type DisplayMode = 'calendar' | 'list';

interface BlockModalState {
  isOpen: boolean;
  initialDate?: Date;
  initialProfessionalId?: string;
  initialStartTime?: string;
  initialEndTime?: string;
}

interface ContextMenuState {
  isOpen: boolean;
  x: number;
  y: number;
  appointment: AppointmentResponse | null;
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
  const [displayMode, setDisplayMode] = useState<DisplayMode>('calendar');

  // Estado do modal
  const [modalState, setModalState] = useState<AppointmentModalState>({
    isOpen: false,
    mode: 'create',
  });

  // Estado do modal de comanda
  const [commandModalState, setCommandModalState] = useState({
    isOpen: false,
    commandId: '',
  });

  // Estado do menu de contexto (botão direito)
  const [contextMenuState, setContextMenuState] = useState<ContextMenuState>({
    isOpen: false,
    x: 0,
    y: 0,
    appointment: null,
  });

  // Estados para checkboxes da sidebar
  const [showBlockSchedule, setShowBlockSchedule] = useState(false);
  const [showOnlyAwaitingPayment, setShowOnlyAwaitingPayment] = useState(false);

  // Estado do modal de bloqueio
  const [blockModalState, setBlockModalState] = useState<BlockModalState>({
    isOpen: false,
  });

  // Hook para atualizar status do agendamento
  const updateStatus = useUpdateAppointmentStatus();

  // Hook para criar comanda a partir de agendamento
  const createCommand = useCreateCommandFromAppointment();

  // Buscar appointments para view de lista
  // Formatar datas como YYYY-MM-DD conforme esperado pelo backend
  const startDate = useMemo(() => {
    return format(new Date(currentDate), 'yyyy-MM-dd');
  }, [currentDate]);

  const endDate = useMemo(() => {
    return format(new Date(currentDate), 'yyyy-MM-dd');
  }, [currentDate]);

  const {data: appointmentsData, isLoading: isLoadingAppointments} = useAppointments({
    start_date: displayMode === 'list' ? startDate : undefined,
    end_date: displayMode === 'list' ? endDate : undefined,
    status: showOnlyAwaitingPayment ? ['AWAITING_PAYMENT'] : undefined,
  });

  // Suporta tanto 'data' quanto 'appointments' para compatibilidade
  const appointments = useMemo(
    () => appointmentsData?.data || (appointmentsData as unknown as { appointments?: AppointmentResponse[] })?.appointments || [],
    [appointmentsData]
  );

  // Filtrar appointments para lista
  const filteredAppointments = useMemo(() => {
    if (displayMode !== 'list') return [];
    
    let filtered = appointments;
    
    if (showOnlyAwaitingPayment) {
      filtered = filtered.filter(apt => apt.status === 'AWAITING_PAYMENT');
    }
    
    // Ordenar por horário
    return filtered.sort((a, b) => 
      new Date(a.start_time).getTime() - new Date(b.start_time).getTime()
    );
  }, [appointments, displayMode, showOnlyAwaitingPayment]);

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

  // Handler para abrir comanda a partir de um agendamento
  const handleOpenCommand = useCallback((appointment: AppointmentResponse) => {
    if (appointment.command_id) {
      // Se já tem comanda, abrir direto
      setCommandModalState({
        isOpen: true,
        commandId: appointment.command_id,
      });
    } else {
      // Criar comanda a partir do agendamento
      createCommand.mutate(appointment.id, {
        onSuccess: (data) => {
          setCommandModalState({
            isOpen: true,
            commandId: data.id,
          });
        },
        onError: () => {
          toast.error('Erro ao criar comanda');
        },
      });
    }
  }, [createCommand]);

  const handleEventClick = useCallback((state: AppointmentModalState) => {
    // Se recebeu appointment completo (do FullCalendar)
    if (state.appointment) {
      // SEMPRE abrir a comanda ao clicar em um agendamento
      handleOpenCommand(state.appointment);
    } else {
      // Sem appointment, apenas abrir modal normal
      setModalState(state);
    }
  }, [handleOpenCommand]);

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

  // Handler para menu de contexto (botão direito)
  const handleEventContextMenu = useCallback((state: AppointmentModalState, event: React.MouseEvent) => {
    if (state.appointment) {
      setContextMenuState({
        isOpen: true,
        x: event.clientX,
        y: event.clientY,
        appointment: state.appointment,
      });
    }
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
          {/* Toggle Calendário/Lista */}
          <Tabs value={displayMode} onValueChange={(value) => setDisplayMode(value as DisplayMode)}>
            <TabsList className="grid w-[200px] grid-cols-2">
              <TabsTrigger value="calendar" className="gap-2">
                <CalendarDays className="h-4 w-4" />
                Calendário
              </TabsTrigger>
              <TabsTrigger value="list" className="gap-2">
                <List className="h-4 w-4" />
                Lista
              </TabsTrigger>
            </TabsList>
          </Tabs>

          {/* Seleção de Período (apenas para calendário) */}
          {displayMode === 'calendar' && (
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
          )}

          <Button onClick={handleNewAppointment} size="sm" data-testid="btn-new-appointment">
            <Calendar className="mr-2 h-4 w-4" />
            Novo Agendamento
          </Button>
        </div>
      </header>

      {/* ================================================================== */}
      {/* ÁREA PRINCIPAL - Grid Layout */}
      {/* ================================================================== */}
      <div className="flex flex-1 overflow-hidden">
        {/* Conteúdo Principal - Calendário ou Lista */}
        <main className="flex-1 overflow-hidden bg-background p-4">
          {displayMode === 'calendar' ? (
            <div className="h-full w-full overflow-hidden rounded-lg border bg-card shadow-sm">
              <AgendaCalendar
                currentDate={currentDate}
                viewType={viewType}
                professionalIds={[]}
                onEventClick={handleEventClick}
                onDateSelect={handleSlotSelect}
                onEventContextMenu={handleEventContextMenu}
                editable={true}
                isBlockMode={showBlockSchedule}
              />
            </div>
          ) : (
            <div className="h-full w-full overflow-y-auto rounded-lg border bg-card shadow-sm">
              <div className="p-6">
                <div className="mb-6 flex items-center justify-between">
                  <h2 className="text-xl font-semibold">
                    {showOnlyAwaitingPayment ? 'Aguardando Pagamento' : 'Agendamentos do Dia'}
                  </h2>
                  <span className="text-sm text-muted-foreground">
                    {filteredAppointments.length} {filteredAppointments.length === 1 ? 'agendamento' : 'agendamentos'}
                  </span>
                </div>

                {isLoadingAppointments ? (
                  <div className="space-y-4">
                    {[...Array(3)].map((_, i) => (
                      <Skeleton key={i} className="h-32 w-full" />
                    ))}
                  </div>
                ) : filteredAppointments.length === 0 ? (
                  <div className="flex flex-col items-center justify-center py-12 text-center">
                    <Calendar className="h-12 w-12 text-muted-foreground mb-4" />
                    <p className="text-lg font-medium text-muted-foreground">
                      {showOnlyAwaitingPayment 
                        ? 'Nenhum agendamento aguardando pagamento'
                        : 'Nenhum agendamento para este dia'
                      }
                    </p>
                    <p className="text-sm text-muted-foreground mt-2">
                      {showOnlyAwaitingPayment
                        ? 'Desmarque o filtro para ver todos os agendamentos'
                        : 'Selecione outra data ou crie um novo agendamento'
                      }
                    </p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {filteredAppointments.map((appointment) => (
                      <AppointmentCardWithCommand
                        key={appointment.id}
                        appointment={appointment}
                        onClick={() => {
                          // Sempre abrir a comanda ao clicar em um agendamento
                          handleOpenCommand(appointment);
                        }}
                        onCloseCommand={() => {
                          handleOpenCommand(appointment);
                        }}
                        variant="default"
                      />
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}
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
              
              {/* Modo de Bloqueio (apenas para calendário) */}
              {displayMode === 'calendar' && (
                <>
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
                </>
              )}

              {/* Filtro Aguardando Pagamento (apenas para lista) */}
              {displayMode === 'list' && (
                <div className="flex items-center space-x-2 rounded-lg border bg-card p-3 shadow-sm">
                  <Checkbox
                    id="awaiting-payment"
                    checked={showOnlyAwaitingPayment}
                    onCheckedChange={(checked) => setShowOnlyAwaitingPayment(!!checked)}
                  />
                  <label
                    htmlFor="awaiting-payment"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  >
                    Apenas Aguardando Pagamento
                  </label>
                </div>
              )}
            </div>
          </div>
        </aside>
      </div>

      {/* Modais */}
      <AppointmentModal state={modalState} onClose={handleCloseModal} />
      <CommandModal
        commandId={commandModalState.commandId}
        open={commandModalState.isOpen}
        onOpenChange={(open) => setCommandModalState(prev => ({ ...prev, isOpen: open }))}
      />
      <BlockScheduleModal
        isOpen={blockModalState.isOpen}
        onClose={handleCloseBlockModal}
        initialDate={blockModalState.initialDate}
        initialProfessionalId={blockModalState.initialProfessionalId}
        initialStartTime={blockModalState.initialStartTime}
        initialEndTime={blockModalState.initialEndTime}
      />

      {/* Menu de Contexto (Botão Direito) */}
      <AppointmentContextMenu
        isOpen={contextMenuState.isOpen}
        x={contextMenuState.x}
        y={contextMenuState.y}
        appointment={contextMenuState.appointment}
        onClose={() => setContextMenuState({ isOpen: false, x: 0, y: 0, appointment: null })}
        onView={() => {
          if (contextMenuState.appointment) {
            setModalState({
              isOpen: true,
              mode: 'view',
              appointment: contextMenuState.appointment,
            });
          }
        }}
        onEdit={() => {
          if (contextMenuState.appointment) {
            setModalState({
              isOpen: true,
              mode: 'edit',
              appointment: contextMenuState.appointment,
            });
          }
        }}
        onOpenCommand={() => {
          if (contextMenuState.appointment) {
            const apt = contextMenuState.appointment;
            // Se já tem comanda, abrir direto
            if (apt.command_id) {
              setCommandModalState({
                isOpen: true,
                commandId: apt.command_id,
              });
            } else {
              // Criar comanda a partir do agendamento
              createCommand.mutate(apt.id, {
                onSuccess: (data) => {
                  setCommandModalState({
                    isOpen: true,
                    commandId: data.id,
                  });
                },
                onError: () => {
                  toast.error('Erro ao criar comanda');
                },
              });
            }
          }
        }}
        onConfirm={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'CONFIRMED' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Agendamento confirmado!'),
                onError: () => toast.error('Erro ao confirmar'),
              }
            );
          }
        }}
        onCheckIn={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'CHECKED_IN' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Check-in realizado!'),
                onError: () => toast.error('Erro ao fazer check-in'),
              }
            );
          }
        }}
        onStartService={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'IN_SERVICE' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Atendimento iniciado!'),
                onError: () => toast.error('Erro ao iniciar atendimento'),
              }
            );
          }
        }}
        onFinishService={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'AWAITING_PAYMENT' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Atendimento finalizado! Aguardando pagamento.'),
                onError: () => toast.error('Erro ao finalizar atendimento'),
              }
            );
          }
        }}
        onCloseCommand={() => {
          if (contextMenuState.appointment?.command_id) {
            setCommandModalState({
              isOpen: true,
              commandId: contextMenuState.appointment.command_id,
            });
          } else {
            toast.error('Nenhuma comanda vinculada a este agendamento');
          }
        }}
        onComplete={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'DONE' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Agendamento concluído!'),
                onError: () => toast.error('Erro ao concluir'),
              }
            );
          }
        }}
        onNoShow={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'NO_SHOW' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Marcado como não compareceu'),
                onError: () => toast.error('Erro ao marcar no-show'),
              }
            );
          }
        }}
        onCancel={() => {
          if (contextMenuState.appointment) {
            updateStatus.mutate(
              { id: contextMenuState.appointment.id, data: { status: 'CANCELED' }, currentStatus: contextMenuState.appointment.status },
              {
                onSuccess: () => toast.success('Agendamento cancelado'),
                onError: () => toast.error('Erro ao cancelar'),
              }
            );
          }
        }}
      />
    </div>
  );
}
