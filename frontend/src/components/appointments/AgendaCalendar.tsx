'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente AgendaCalendar - Layout AppBarber Style
 *
 * Calendário profissional com:
 * - Slots de 10 minutos (como AppBarber)
 * - Colunas por profissional
 * - PROFISSIONAIS NO TOPO DAS COLUNAS
 * - Horários na lateral esquerda
 *
 * ⚠️ LICENÇA: FullCalendar Scheduler – Modo Avaliação
 */

import type {
    DateSelectArg,
    EventClickArg,
    EventDropArg,
} from '@fullcalendar/core';
import interactionPlugin, {
    type EventResizeDoneArg,
} from '@fullcalendar/interaction';
import listPlugin from '@fullcalendar/list';
import type FullCalendarType from '@fullcalendar/react';
import FullCalendar from '@fullcalendar/react';
import resourceTimeGridPlugin from '@fullcalendar/resource-timegrid';
import { useCallback, useEffect, useMemo, useRef } from 'react';

import { Skeleton } from '@/components/ui/skeleton';
import {
    useCalendarEvents,
    useCalendarResources,
    useUpdateAppointment,
} from '@/hooks/use-appointments';
import {
    FULLCALENDAR_LICENSE_KEY,
    FULLCALENDAR_LOCALE_PT_BR,
} from '@/lib/fullcalendar-config';
import type { AppointmentModalState, CalendarEvent } from '@/types/appointment';

import './agenda-calendar.css';

// =============================================================================
// TYPES
// =============================================================================

type ViewType = 'day' | 'week' | 'month';

interface AgendaCalendarProps {
  /** Data atual do calendário */
  currentDate: Date;
  /** Tipo de visualização */
  viewType: ViewType;
  /** IDs dos profissionais para filtrar (vazio = todos) */
  professionalIds?: string[];
  /** Callback quando um evento é clicado */
  onEventClick?: (state: AppointmentModalState) => void;
  /** Callback quando um slot vazio é selecionado */
  onDateSelect?: (state: AppointmentModalState) => void;
  /** Se o usuário pode editar */
  editable?: boolean;
  /** Modo de bloqueio de horário ativo */
  isBlockMode?: boolean;
}

// Map de views
const VIEW_MAP: Record<ViewType, string> = {
  day: 'resourceTimeGridDay',
  week: 'resourceTimeGridWeek',
  month: 'resourceTimeGridDay', // Para mês, usamos day por enquanto
};

// =============================================================================
// COMPONENT
// =============================================================================

export function AgendaCalendar({
  currentDate,
  viewType,
  professionalIds = [],
  onEventClick,
  onDateSelect,
  editable = true,
  isBlockMode = false,
}: AgendaCalendarProps) {
  const calendarRef = useRef<FullCalendarType>(null);

  // Datas para filtro - estabilizadas
  const dateRange = useMemo(() => {
    const start = new Date(currentDate);
    start.setHours(0, 0, 0, 0);
    start.setDate(start.getDate() - start.getDay() + 1); // Segunda
    const end = new Date(start);
    end.setDate(end.getDate() + 6); // Domingo
    return {
      date_from: start.toISOString().split('T')[0],
      date_to: end.toISOString().split('T')[0],
    };
  }, [currentDate]);

  // ==========================================================================
  // QUERIES
  // ==========================================================================

  const { data: events = [], isLoading: isLoadingEvents } = useCalendarEvents({
    ...dateRange,
    professional_id:
      professionalIds.length === 1 ? professionalIds[0] : undefined,
  });

  const { data: resources = [], isLoading: isLoadingResources } =
    useCalendarResources();

  const updateAppointment = useUpdateAppointment();

  // Filtrar recursos
  const filteredResources = useMemo(() => {
    if (professionalIds.length === 0) return resources;
    return resources.filter((r) => professionalIds.includes(r.id));
  }, [resources, professionalIds]);

  // ==========================================================================
  // SYNC COM PROPS EXTERNAS
  // ==========================================================================

  // Atualizar data quando currentDate mudar
  useEffect(() => {
    const calendarApi = calendarRef.current?.getApi();
    if (calendarApi) {
      calendarApi.gotoDate(currentDate);
    }
  }, [currentDate]);

  // Atualizar view quando viewType mudar
  useEffect(() => {
    const calendarApi = calendarRef.current?.getApi();
    if (calendarApi) {
      calendarApi.changeView(VIEW_MAP[viewType]);
    }
  }, [viewType]);

  // ==========================================================================
  // HANDLERS
  // ==========================================================================

  const handleEventClick = useCallback(
    (info: EventClickArg) => {
      const calendarEvent = info.event
        .extendedProps as CalendarEvent['extendedProps'];
      onEventClick?.({
        isOpen: true,
        mode: 'view',
        appointment: calendarEvent.appointment,
      });
    },
    [onEventClick]
  );

  const handleDateSelect = useCallback(
    (info: DateSelectArg) => {
      onDateSelect?.({
        isOpen: true,
        mode: 'create',
        initialDate: info.start,
        initialProfessionalId: info.resource?.id,
      });
      calendarRef.current?.getApi()?.unselect();
    },
    [onDateSelect]
  );

  const handleEventDrop = useCallback(
    (info: EventDropArg) => {
      const event = info.event;
      const appointment = (
        event.extendedProps as CalendarEvent['extendedProps']
      ).appointment;

      updateAppointment.mutate(
        {
          id: appointment.id,
          data: {
            start_time: event.start?.toISOString(),
            professional_id: event.getResources()[0]?.id,
          },
        },
        { onError: () => info.revert() }
      );
    },
    [updateAppointment]
  );

  const handleEventResize = useCallback(
    (info: EventResizeDoneArg) => {
      const appointment = (
        info.event.extendedProps as CalendarEvent['extendedProps']
      ).appointment;
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
  // LOADING
  // ==========================================================================

  if (isLoadingEvents || isLoadingResources) {
    return (
      <div className="h-full flex flex-col bg-white">
        {/* Body skeleton */}
        <div className="flex-1 p-4">
          <Skeleton className="h-full w-full" />
        </div>
      </div>
    );
  }

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <div className="agenda-calendar-wrapper h-full flex flex-col bg-white">
      <FullCalendar
        ref={calendarRef}
        // Plugins
        plugins={[resourceTimeGridPlugin, interactionPlugin, listPlugin]}
        // Licença
        schedulerLicenseKey={FULLCALENDAR_LICENSE_KEY}
        // Localização
        locale={FULLCALENDAR_LOCALE_PT_BR}
        // View
        initialView={VIEW_MAP[viewType]}
        initialDate={currentDate}
        // Recursos (profissionais)
        resources={filteredResources}
        resourceOrder="title"
        // Eventos
        events={events}
        // ========================================
        // HORÁRIOS - Slots de 10 minutos (AppBarber)
        // ========================================
        slotMinTime="07:00:00"
        slotMaxTime="22:00:00"
        slotDuration="00:10:00"
        slotLabelInterval="00:10:00"
        // Layout
        height="100%"
        expandRows
        nowIndicator
        // Header - OCULTO (usamos header customizado no page)
        headerToolbar={false}
        // Dias úteis
        businessHours={{
          daysOfWeek: [1, 2, 3, 4, 5, 6],
          startTime: '08:00',
          endTime: '21:00',
        }}
        // Interatividade
        editable={editable}
        selectable={editable}
        selectMirror
        navLinks={false}
        // Formato de hora
        eventTimeFormat={{
          hour: '2-digit',
          minute: '2-digit',
          meridiem: false,
          hour12: false,
        }}
        slotLabelFormat={{
          hour: '2-digit',
          minute: '2-digit',
          hour12: false,
        }}
        // Handlers
        eventClick={handleEventClick}
        select={handleDateSelect}
        eventDrop={handleEventDrop}
        eventResize={handleEventResize}
        // Texto
        allDayText=""
        allDaySlot={false}
        noEventsText="Nenhum agendamento"
        // Customização de eventos - Design System
        eventContent={(eventInfo) => {
          const appointment = (
            eventInfo.event.extendedProps as CalendarEvent['extendedProps']
          ).appointment;
          const status = appointment?.status || 'CREATED';
          const serviceName = appointment?.services?.[0]?.name;

          return (
            <div className="nexo-event-content" data-status={status}>
              <div className="nexo-event-time">{eventInfo.timeText}</div>
              <div className="nexo-event-title">
                {eventInfo.event.title}
              </div>
              {serviceName && (
                <div className="nexo-event-info">{serviceName}</div>
              )}
            </div>
          );
        }}
      />
    </div>
  );
}

export default AgendaCalendar;
