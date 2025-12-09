/**
 * AppointmentCardWithCommand - Wrapper do AppointmentCard com integração de Comanda
 * 
 * Este componente encapsula AppointmentCard e adiciona funcionalidade de fechamento de comanda.
 * Quando o appointment está em status AWAITING_PAYMENT, permite abrir o CommandModal.
 */

'use client';

import { useCommandByAppointment, useCreateCommandFromAppointment } from '@/hooks/use-commands';
import type { AppointmentResponse } from '@/types/appointment';
import { useState } from 'react';
import { CommandModal } from '../agendamentos/CommandModal';
import { AppointmentCard } from './AppointmentCard';

interface AppointmentCardWithCommandProps {
  /** Dados do agendamento */
  appointment: AppointmentResponse;
  /** Callback quando clicado para ver detalhes */
  onClick?: () => void;
  /** Callback para editar */
  onEdit?: () => void;
  /** Callback para cancelar */
  onCancel?: () => void;
  /** Callback para confirmar (CREATED -> CONFIRMED) */
  onConfirm?: () => void;
  /** Callback para check-in (CONFIRMED -> CHECKED_IN) */
  onCheckIn?: () => void;
  /** Callback para iniciar atendimento (CHECKED_IN/CONFIRMED -> IN_SERVICE) */
  onStartService?: () => void;
  /** Callback para finalizar atendimento (IN_SERVICE -> AWAITING_PAYMENT) */
  onFinishService?: () => void;
  /** Callback para concluir (AWAITING_PAYMENT -> DONE) */
  onComplete?: () => void;
  /** Callback para fechar comanda (AWAITING_PAYMENT) - sobrescreve comportamento padrão */
  onCloseCommand?: () => void;
  /** Callback para marcar no-show */
  onNoShow?: () => void;
  /** Variante de exibição */
  variant?: 'default' | 'compact';
  /** Classe CSS adicional */
  className?: string;
}

export function AppointmentCardWithCommand({
  appointment,
  onCloseCommand: onCloseCommandProp,
  ...props
}: AppointmentCardWithCommandProps) {
  const [commandModalOpen, setCommandModalOpen] = useState(false);
  
  // Buscar comanda associada ao appointment (se existir e status for AWAITING_PAYMENT)
  const appointmentIdForCommand = appointment.status === 'AWAITING_PAYMENT' ? appointment.id : undefined;
  const { data: command } = useCommandByAppointment(appointmentIdForCommand);
  
  // Mutation para criar comanda
  const createCommand = useCreateCommandFromAppointment();

  const handleCloseCommand = async () => {
    // Se foi passado um handler externo, usar ele
    if (onCloseCommandProp) {
      onCloseCommandProp();
      return;
    }

    // Se já existe comanda, abrir modal direto
    if (command) {
      setCommandModalOpen(true);
      return;
    }

    // Se não existe, criar comanda primeiro
    try {
      await createCommand.mutateAsync(appointment.id);
      // Após criar, abrir modal (a query será revalidada automaticamente)
      setCommandModalOpen(true);
    } catch {
      // Erro já tratado pelo hook (toast)
    }
  };

  return (
    <>
      <AppointmentCard
        appointment={appointment}
        onCloseCommand={
          appointment.status === 'AWAITING_PAYMENT' ? handleCloseCommand : undefined
        }
        {...props}
      />

      {/* Modal de Fechamento de Comanda */}
      {command && (
        <CommandModal
          commandId={command.id}
          open={commandModalOpen}
          onOpenChange={setCommandModalOpen}
        />
      )}
    </>
  );
}

export default AppointmentCardWithCommand;
