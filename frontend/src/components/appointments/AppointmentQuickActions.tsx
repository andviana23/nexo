'use client';

/**
 * Menu de Ações Rápidas para Appointments
 * Aparece ao clicar em um evento no calendário
 */

import {
    CheckCircle2,
    Clock,
    CreditCard,
    Edit,
    Scissors,
    UserCheck,
    UserX,
    XCircle,
} from 'lucide-react';
import { useCallback } from 'react';

import { Button } from '@/components/ui/button';
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '@/components/ui/popover';
import {
    useCheckInAppointment,
    useCompleteAppointment,
    useConfirmAppointment,
    useFinishServiceAppointment,
    useNoShowAppointment,
    useStartServiceAppointment,
} from '@/hooks/use-appointments';
import type { AppointmentResponse, AppointmentStatus } from '@/types/appointment';
interface AppointmentQuickActionsProps {
  appointment: AppointmentResponse;
  children: React.ReactNode;
  onEdit?: () => void;
  onCancel?: () => void;
  onCloseCommand?: () => void;
}

export function AppointmentQuickActions({
  appointment,
  children,
  onEdit,
  onCancel,
  onCloseCommand,
}: AppointmentQuickActionsProps) {
  const confirm = useConfirmAppointment();
  const checkIn = useCheckInAppointment();
  const startService = useStartServiceAppointment();
  const finishService = useFinishServiceAppointment();
  const complete = useCompleteAppointment();
  const noShow = useNoShowAppointment();

  const handleConfirm = useCallback(() => {
    confirm.mutate(appointment.id);
  }, [confirm, appointment.id]);

  const handleCheckIn = useCallback(() => {
    checkIn.mutate(appointment.id);
  }, [checkIn, appointment.id]);

  const handleStartService = useCallback(() => {
    startService.mutate(appointment.id);
  }, [startService, appointment.id]);

  const handleFinishService = useCallback(() => {
    finishService.mutate(appointment.id);
  }, [finishService, appointment.id]);

  const handleComplete = useCallback(() => {
    complete.mutate(appointment.id);
  }, [complete, appointment.id]);

  const handleNoShow = useCallback(() => {
    noShow.mutate(appointment.id);
  }, [noShow, appointment.id]);

  // Determinar ações disponíveis baseado no status
  const getAvailableActions = (status: AppointmentStatus) => {
    const actions: React.ReactNode[] = [];

    switch (status) {
      case 'CREATED':
        actions.push(
          <Button
            key="confirm"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleConfirm}
          >
            <CheckCircle2 className="mr-2 h-4 w-4 text-blue-600" />
            Confirmar
          </Button>
        );
        actions.push(
          <Button
            key="checkin"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleCheckIn}
          >
            <UserCheck className="mr-2 h-4 w-4 text-violet-600" />
            Cliente Chegou
          </Button>
        );
        actions.push(
          <Button
            key="noshow"
            variant="ghost"
            size="sm"
            className="w-full justify-start text-destructive"
            onClick={handleNoShow}
          >
            <UserX className="mr-2 h-4 w-4" />
            Não Compareceu
          </Button>
        );
        break;

      case 'CONFIRMED':
        actions.push(
          <Button
            key="checkin"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleCheckIn}
          >
            <UserCheck className="mr-2 h-4 w-4 text-violet-600" />
            Cliente Chegou
          </Button>
        );
        actions.push(
          <Button
            key="start"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleStartService}
          >
            <Scissors className="mr-2 h-4 w-4 text-amber-600" />
            Iniciar Atendimento
          </Button>
        );
        actions.push(
          <Button
            key="noshow"
            variant="ghost"
            size="sm"
            className="w-full justify-start text-destructive"
            onClick={handleNoShow}
          >
            <UserX className="mr-2 h-4 w-4" />
            Não Compareceu
          </Button>
        );
        break;

      case 'CHECKED_IN':
        actions.push(
          <Button
            key="start"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleStartService}
          >
            <Scissors className="mr-2 h-4 w-4 text-amber-600" />
            Iniciar Atendimento
          </Button>
        );
        break;

      case 'IN_SERVICE':
        actions.push(
          <Button
            key="finish"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleFinishService}
          >
            <Clock className="mr-2 h-4 w-4 text-pink-600" />
            Finalizar Atendimento
          </Button>
        );
        break;

      case 'AWAITING_PAYMENT':
        if (onCloseCommand) {
          actions.push(
            <Button
              key="command"
              variant="ghost"
              size="sm"
              className="w-full justify-start"
              onClick={onCloseCommand}
            >
              <CreditCard className="mr-2 h-4 w-4 text-green-600" />
              Fechar Comanda
            </Button>
          );
        }
        actions.push(
          <Button
            key="complete"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={handleComplete}
          >
            <CheckCircle2 className="mr-2 h-4 w-4 text-green-600" />
            Concluir (Pago)
          </Button>
        );
        break;
    }

    // Ações sempre disponíveis (exceto para status finais)
    if (!['DONE', 'NO_SHOW', 'CANCELED'].includes(status)) {
      if (onEdit) {
        actions.push(
          <Button
            key="edit"
            variant="ghost"
            size="sm"
            className="w-full justify-start"
            onClick={onEdit}
          >
            <Edit className="mr-2 h-4 w-4" />
            Editar
          </Button>
        );
      }
      if (onCancel) {
        actions.push(
          <Button
            key="cancel"
            variant="ghost"
            size="sm"
            className="w-full justify-start text-destructive"
            onClick={onCancel}
          >
            <XCircle className="mr-2 h-4 w-4" />
            Cancelar
          </Button>
        );
      }
    }

    return actions;
  };

  const actions = getAvailableActions(appointment.status);

  if (actions.length === 0) {
    return <>{children}</>;
  }

  return (
    <Popover>
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent className="w-64 p-2" align="start">
        <div className="space-y-1">
          <div className="px-2 py-1.5 text-sm font-semibold">
            {appointment.customer_name}
          </div>
          <div className="px-2 pb-2 text-xs text-muted-foreground">
            {appointment.services.map((s) => s.service_name).join(', ')}
          </div>
          <div className="border-t pt-1">{actions}</div>
        </div>
      </PopoverContent>
    </Popover>
  );
}
