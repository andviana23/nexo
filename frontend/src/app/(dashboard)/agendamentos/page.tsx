'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Agendamentos
 *
 * Tela principal do módulo de Agendamentos com calendário visual.
 */

import { CalendarDays, Filter, Plus } from 'lucide-react';
import { useCallback, useState } from 'react';

import { AppointmentCalendar, AppointmentModal } from '@/components/appointments';
import { Button } from '@/components/ui/button';
import type { AppointmentModalState } from '@/types/appointment';

// =============================================================================
// PAGE COMPONENT
// =============================================================================

export default function AgendamentosPage() {
  // Estado do modal
  const [modalState, setModalState] = useState<AppointmentModalState>({
    isOpen: false,
    mode: 'create',
  });

  // Estado dos filtros
  const [selectedProfessionals] = useState<string[]>([]);

  // Handlers
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
      initialDate: new Date(),
    });
  }, []);

  const handleCloseModal = useCallback(() => {
    setModalState((prev) => ({ ...prev, isOpen: false }));
  }, []);

  return (
    <div className="space-y-6">
      {/* Header da Página */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <CalendarDays className="h-6 w-6 text-primary" />
            Agendamentos
          </h1>
          <p className="text-muted-foreground">
            Gerencie os agendamentos da sua barbearia
          </p>
        </div>

        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Filter className="mr-2 h-4 w-4" />
            Filtros
          </Button>
          <Button size="sm" onClick={handleNewAppointment}>
            <Plus className="mr-2 h-4 w-4" />
            Novo Agendamento
          </Button>
        </div>
      </div>

      {/* Calendário */}
      <div className="rounded-lg border bg-card p-4 shadow-sm">
        <AppointmentCalendar
          professionalIds={selectedProfessionals}
          onEventClick={handleEventClick}
          onDateSelect={handleDateSelect}
          editable={true}
        />
      </div>

      {/* Modal de Agendamento */}
      <AppointmentModal
        state={modalState}
        onClose={handleCloseModal}
      />
    </div>
  );
}
