'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente de Calendário de Agendamentos
 *
 * @component AppointmentCalendar
 * @description Calendário visual usando FullCalendar com view de recursos (barbeiros)
 * 
 * ⚠️ LICENÇA: FullCalendar Scheduler – Modo Avaliação
 * Esta chave só pode ser usada em desenvolvimento.
 * Antes da produção, adquirir licença comercial.
 */

import type {
    DateSelectArg,
    EventApi,
    EventClickArg,
    EventDropArg,
} from '@fullcalendar/core';
import interactionPlugin, { type EventResizeDoneArg } from '@fullcalendar/interaction';
import listPlugin from '@fullcalendar/list';
import type FullCalendarType from '@fullcalendar/react';
import FullCalendar from '@fullcalendar/react';
import resourceTimeGridPlugin from '@fullcalendar/resource-timegrid';
import { useCallback, useMemo, useRef, useState } from 'react';

import { Skeleton } from '@/components/ui/skeleton';
import {
    useCalendarEvents,
    useCalendarResources,
    useUpdateAppointment,
} from '@/hooks/use-appointments';
import {
    FULLCALENDAR_DEFAULTS,
    FULLCALENDAR_LICENSE_KEY,
    FULLCALENDAR_LOCALE_PT_BR,
} from '@/lib/fullcalendar-config';
import type { AppointmentModalState, CalendarEvent } from '@/types/appointment';

// =============================================================================
// TYPES
// =============================================================================

interface AppointmentCalendarProps {
  /** Data inicial do calendário */
  initialDate?: Date;
  /** IDs dos profissionais para filtrar (vazio = todos) */
  professionalIds?: string[];
  /** Callback quando um evento é clicado */
  onEventClick?: (state: AppointmentModalState) => void;
  /** Callback quando um slot vazio é selecionado (criar agendamento) */
  onDateSelect?: (state: AppointmentModalState) => void;
  /** Se o usuário pode editar (arrastar/redimensionar) */
  editable?: boolean;
  /** Classe CSS adicional */
  className?: string;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function AppointmentCalendar({
  initialDate,
  professionalIds = [],
  onEventClick,
  onDateSelect,
  editable = true,
  className,
}: AppointmentCalendarProps) {
  const calendarRef = useRef<FullCalendarType>(null);
  
  // Estabilizar a data inicial para evitar re-renders infinitos
  const [stableDate] = useState(() => initialDate || new Date());

  // Datas para filtro (semana atual) - Usar apenas data sem hora para evitar mudanças
  const dateRange = useMemo(() => {
    const start = new Date(stableDate);
    start.setHours(0, 0, 0, 0); // Zerar hora
    start.setDate(start.getDate() - start.getDay() + 1); // Segunda
    const end = new Date(start);
    end.setDate(end.getDate() + 6); // Domingo
    end.setHours(23, 59, 59, 999); // Final do dia
    return {
      start_date: start.toISOString().split('T')[0], // Apenas YYYY-MM-DD
      end_date: end.toISOString().split('T')[0],     // Apenas YYYY-MM-DD
    };
  }, [stableDate]);

  // Queries
  const { 
    data: events = [], 
    isLoading: isLoadingEvents,
    error: eventsError,
    isError: isEventsError,
  } = useCalendarEvents({
    ...dateRange,
    professional_id: professionalIds.length === 1 ? professionalIds[0] : undefined,
  });

  const { 
    data: resources = [], 
    isLoading: isLoadingResources,
    error: resourcesError,
    isError: isResourcesError,
  } = useCalendarResources();

  // DEBUG: Log detalhado
  console.log('[AppointmentCalendar] Estado:', {
    isLoadingEvents,
    isLoadingResources,
    isEventsError,
    isResourcesError,
    eventsError: eventsError?.message,
    resourcesError: resourcesError?.message,
    eventsCount: events.length,
    resourcesCount: resources.length,
  });

  // Mutation para drag & drop
  const updateAppointment = useUpdateAppointment();

  // Filtrar recursos se necessário
  const filteredResources = useMemo(() => {
    if (professionalIds.length === 0) return resources;
    return resources.filter((r) => professionalIds.includes(r.id));
  }, [resources, professionalIds]);

  // ==========================================================================
  // HANDLERS
  // ==========================================================================

  /**
   * Handler quando um evento é clicado
   */
  const handleEventClick = useCallback(
    (info: EventClickArg) => {
      const calendarEvent = info.event.extendedProps as CalendarEvent['extendedProps'];
      
      onEventClick?.({
        isOpen: true,
        mode: 'view',
        appointment: calendarEvent.appointment,
      });
    },
    [onEventClick]
  );

  /**
   * Handler quando um slot vazio é selecionado (criar novo agendamento)
   */
  const handleDateSelect = useCallback(
    (info: DateSelectArg) => {
      onDateSelect?.({
        isOpen: true,
        mode: 'create',
        initialDate: info.start,
        initialProfessionalId: info.resource?.id,
      });

      // Limpar seleção
      const calendarApi = calendarRef.current?.getApi();
      calendarApi?.unselect();
    },
    [onDateSelect]
  );

  /**
   * Handler quando um evento é arrastado (reagendar)
   * Suporta drag horizontal (horário) e vertical (profissional)
   */
  const handleEventDrop = useCallback(
    (info: EventDropArg) => {
      const event = info.event;
      const appointment = (event.extendedProps as CalendarEvent['extendedProps']).appointment;
      
      // Pegar novo profissional (pode ter mudado com drag vertical)
      const newResource = event.getResources()[0];
      const newProfessionalId = newResource?.id;
      
      // Validações
      if (!newProfessionalId) {
        console.error('[AppointmentCalendar] Profissional não encontrado após drop');
        info.revert();
        return;
      }

      // Detectar se mudou de profissional
      const changedProfessional = newProfessionalId !== appointment.professional_id;
      
      updateAppointment.mutate(
        {
          id: appointment.id,
          data: {
            new_start_time: event.start?.toISOString() || '',
            professional_id: newProfessionalId,
          },
        },
        {
          onSuccess: () => {
            if (changedProfessional) {
              console.log('[AppointmentCalendar] Profissional alterado:', {
                from: appointment.professional_id,
                to: newProfessionalId,
              });
            }
          },
          onError: (error) => {
            console.error('[AppointmentCalendar] Erro ao mover agendamento:', error);
            // Reverter se houver erro
            info.revert();
          },
        }
      );
    },
    [updateAppointment]
  );

  /**
   * Handler quando um evento é redimensionado
   */
  const handleEventResize = useCallback(
    (info: EventResizeDoneArg) => {
      const event = info.event;
      const appointment = (event.extendedProps as CalendarEvent['extendedProps']).appointment;

      // Verificar se a duração mudou significativamente (serviços diferentes)
      // Por enquanto, apenas reverter - reagendar requer seleção de serviços
      info.revert();
      
      onEventClick?.({
        isOpen: true,
        mode: 'edit',
        appointment,
      });
    },
    [onEventClick]
  );

  /**
   * Validação visual durante o arraste
   * Permite mover apenas se não houver conflito de horário
   */
  const handleEventAllow = useCallback(
    (dropInfo: { start: Date; end: Date; resourceId?: string }, draggedEvent: EventApi | null) => {
      const calendarApi = calendarRef.current?.getApi();
      if (!calendarApi || !draggedEvent) return true;

      const draggedAppointment = draggedEvent.extendedProps?.appointment;
      if (!draggedAppointment) return true;

      // Pegar todos os eventos do profissional de destino
      const targetResourceId = dropInfo.resourceId || draggedEvent.getResources()[0]?.id;
      if (!targetResourceId) return false;

      const allEvents = calendarApi.getEvents();
      
      // Verificar se há conflito com outros agendamentos
      const hasConflict = allEvents.some((event) => {
        // Ignorar o próprio evento sendo arrastado
        if (event.id === draggedEvent.id) return false;

        // Verificar se é do mesmo profissional
        const eventResource = event.getResources()[0];
        if (eventResource?.id !== targetResourceId) return false;

        // Verificar sobreposição de horários
        const eventStart = event.start;
        const eventEnd = event.end;
        if (!eventStart || !eventEnd) return false;

        const overlapStart = dropInfo.start < eventEnd;
        const overlapEnd = dropInfo.end > eventStart;

        return overlapStart && overlapEnd;
      });

      // Feedback visual (CSS classe será adicionada automaticamente pelo FullCalendar)
      return !hasConflict;
    },
    []
  );

  // ==========================================================================
  // LOADING STATE
  // ==========================================================================

  if (isLoadingEvents || isLoadingResources) {
    return (
      <div className={`space-y-4 ${className}`}>
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-[600px] w-full" />
      </div>
    );
  }

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <div className={`appointment-calendar ${className}`}>
      <FullCalendar
        ref={calendarRef}
        // Plugins
        plugins={[resourceTimeGridPlugin, interactionPlugin, listPlugin]}
        
        // Licença (Modo Avaliação - NÃO USAR EM PRODUÇÃO)
        schedulerLicenseKey={FULLCALENDAR_LICENSE_KEY}
        
        // Configurações padrão
        {...FULLCALENDAR_DEFAULTS}
        locale={FULLCALENDAR_LOCALE_PT_BR}
        
        // View inicial
        initialView="resourceTimeGridDay"
        initialDate={stableDate}
        
        // Recursos (Barbeiros)
        resources={filteredResources}
        resourceOrder="title"
        
        // Eventos (Agendamentos)
        events={events}
        
        // Interatividade
        editable={editable}
        selectable={editable}
        selectMirror
        
        // Handlers
        eventClick={handleEventClick}
        select={handleDateSelect}
        eventDrop={handleEventDrop}
        eventResize={handleEventResize}
        eventAllow={handleEventAllow}
        
        // Drag & Drop Configuration
        dragRevertDuration={300}
        dragScroll
        snapDuration="00:15:00"
        
        // Visual feedback durante arraste
        eventClassNames={(arg) => {
          const status = arg.event.extendedProps?.appointment?.status;
          return [`event-status-${status?.toLowerCase()}`];
        }}
        
        // Layout
        height="auto"
        stickyHeaderDates
        
        // Customização de eventos
        eventContent={(eventInfo) => (
          <div className="fc-event-content p-1 overflow-hidden">
            <div className="font-medium text-xs truncate">
              {eventInfo.timeText}
            </div>
            <div className="text-xs truncate">
              {eventInfo.event.title}
            </div>
          </div>
        )}
      />
    </div>
  );
}

export default AppointmentCalendar;
