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
      date_from: start.toISOString().split('T')[0], // Apenas YYYY-MM-DD
      date_to: end.toISOString().split('T')[0],     // Apenas YYYY-MM-DD
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
   */
  const handleEventDrop = useCallback(
    (info: EventDropArg) => {
      const event = info.event;
      const appointment = (event.extendedProps as CalendarEvent['extendedProps']).appointment;

      updateAppointment.mutate(
        {
          id: appointment.id,
          data: {
            start_time: event.start?.toISOString(),
            professional_id: event.getResources()[0]?.id,
          },
        },
        {
          onError: () => {
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
