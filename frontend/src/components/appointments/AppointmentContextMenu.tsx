'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Menu de Contexto (Botão Direito) para Agendamentos
 *
 * @component AppointmentContextMenu
 * @description Menu que aparece ao clicar com botão direito em um agendamento
 */

import {
    CalendarPlusIcon,
    CheckCircle2Icon,
    CreditCardIcon,
    EditIcon,
    EyeIcon,
    ScissorsIcon,
    UserCheckIcon,
    UserXIcon,
    XCircleIcon
} from 'lucide-react';
import { useEffect, useRef } from 'react';

import { cn } from '@/lib/utils';
import type { AppointmentResponse, AppointmentStatus } from '@/types/appointment';

interface AppointmentContextMenuProps {
  /** Se o menu está aberto */
  isOpen: boolean;
  /** Posição X do mouse */
  x: number;
  /** Posição Y do mouse */
  y: number;
  /** Dados do agendamento */
  appointment: AppointmentResponse | null;
  /** Callback para fechar */
  onClose: () => void;
  /** Callback para visualizar */
  onView?: () => void;
  /** Callback para editar */
  onEdit?: () => void;
  /** Callback para confirmar */
  onConfirm?: () => void;
  /** Callback para check-in */
  onCheckIn?: () => void;
  /** Callback para iniciar atendimento */
  onStartService?: () => void;
  /** Callback para finalizar atendimento */
  onFinishService?: () => void;
  /** Callback para abrir comanda */
  onOpenCommand?: () => void;
  /** Callback para fechar comanda */
  onCloseCommand?: () => void;
  /** Callback para concluir */
  onComplete?: () => void;
  /** Callback para marcar no-show */
  onNoShow?: () => void;
  /** Callback para cancelar */
  onCancel?: () => void;
  /** Callback para reagendar (status finais) */
  onReschedule?: () => void;
}

interface MenuAction {
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  onClick: () => void;
  destructive?: boolean;
  variant?: 'default' | 'primary';
}

export function AppointmentContextMenu({
  isOpen,
  x,
  y,
  appointment,
  onClose,
  onView,
  onEdit,
  onConfirm,
  onCheckIn,
  onStartService,
  onFinishService,
  onOpenCommand,
  onCloseCommand,
  onComplete,
  onNoShow,
  onCancel,
  onReschedule,
}: AppointmentContextMenuProps) {
  const menuRef = useRef<HTMLDivElement>(null);

  // Fechar ao clicar fora
  useEffect(() => {
    if (!isOpen) return;

    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen, onClose]);

  // Fechar ao pressionar ESC
  useEffect(() => {
    if (!isOpen) return;

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  if (!isOpen || !appointment) return null;

  // Construir ações baseadas no status
  const actions: MenuAction[] = [];

  // Ações baseadas no status
  switch (appointment.status as AppointmentStatus) {
    case 'CREATED':
      // CREATED: Confirmar, Editar, Abrir Comanda, Cancelar
      if (onConfirm) {
        actions.push({
          label: 'Confirmar Agendamento',
          icon: CheckCircle2Icon,
          onClick: () => {
            onConfirm();
            onClose();
          },
          variant: 'primary',
        });
      }
      if (onEdit) {
        actions.push({
          label: 'Editar Agendamento',
          icon: EditIcon,
          onClick: () => {
            onEdit();
            onClose();
          },
        });
      }
      if (onOpenCommand) {
        actions.push({
          label: 'Abrir Comanda',
          icon: CreditCardIcon,
          onClick: () => {
            onOpenCommand();
            onClose();
          },
        });
      }
      if (onCancel) {
        actions.push({
          label: 'Cancelar Agendamento',
          icon: XCircleIcon,
          onClick: () => {
            onCancel();
            onClose();
          },
          destructive: true,
        });
      }
      break;

    case 'CONFIRMED':
      // CONFIRMED: Check-In, Editar, Abrir Comanda, No-Show, Cancelar
      if (onCheckIn) {
        actions.push({
          label: 'Fazer Check-In',
          icon: UserCheckIcon,
          onClick: () => {
            onCheckIn();
            onClose();
          },
          variant: 'primary',
        });
      }
      if (onEdit) {
        actions.push({
          label: 'Editar Agendamento',
          icon: EditIcon,
          onClick: () => {
            onEdit();
            onClose();
          },
        });
      }
      if (onOpenCommand) {
        actions.push({
          label: 'Abrir Comanda',
          icon: CreditCardIcon,
          onClick: () => {
            onOpenCommand();
            onClose();
          },
        });
      }
      if (onNoShow) {
        actions.push({
          label: 'Não Compareceu',
          icon: UserXIcon,
          onClick: () => {
            onNoShow();
            onClose();
          },
          destructive: true,
        });
      }
      if (onCancel) {
        actions.push({
          label: 'Cancelar',
          icon: XCircleIcon,
          onClick: () => {
            onCancel();
            onClose();
          },
          destructive: true,
        });
      }
      break;

    case 'CHECKED_IN':
      // CHECKED_IN: Iniciar, Editar, Abrir Comanda, No-Show, Cancelar
      if (onStartService) {
        actions.push({
          label: 'Iniciar Atendimento',
          icon: ScissorsIcon,
          onClick: () => {
            onStartService();
            onClose();
          },
          variant: 'primary',
        });
      }
      if (onEdit) {
        actions.push({
          label: 'Editar Agendamento',
          icon: EditIcon,
          onClick: () => {
            onEdit();
            onClose();
          },
        });
      }
      if (onOpenCommand) {
        actions.push({
          label: 'Abrir Comanda',
          icon: CreditCardIcon,
          onClick: () => {
            onOpenCommand();
            onClose();
          },
        });
      }
      if (onNoShow) {
        actions.push({
          label: 'Não Compareceu',
          icon: UserXIcon,
          onClick: () => {
            onNoShow();
            onClose();
          },
          destructive: true,
        });
      }
      if (onCancel) {
        actions.push({
          label: 'Cancelar',
          icon: XCircleIcon,
          onClick: () => {
            onCancel();
            onClose();
          },
          destructive: true,
        });
      }
      break;

    case 'IN_SERVICE':
      // IN_SERVICE: Finalizar, Abrir Comanda, Cancelar
      if (onFinishService) {
        actions.push({
          label: 'Finalizar Atendimento',
          icon: CheckCircle2Icon,
          onClick: () => {
            onFinishService();
            onClose();
          },
          variant: 'primary',
        });
      }
      if (onOpenCommand) {
        actions.push({
          label: 'Abrir Comanda',
          icon: CreditCardIcon,
          onClick: () => {
            onOpenCommand();
            onClose();
          },
        });
      }
      if (onCancel) {
        actions.push({
          label: 'Cancelar',
          icon: XCircleIcon,
          onClick: () => {
            onCancel();
            onClose();
          },
          destructive: true,
        });
      }
      break;

    case 'AWAITING_PAYMENT':
      // AWAITING_PAYMENT: Fechar Comanda, Concluir, Cancelar
      if (onCloseCommand) {
        actions.push({
          label: 'Fechar Comanda',
          icon: CreditCardIcon,
          onClick: () => {
            onCloseCommand();
            onClose();
          },
          variant: 'primary',
        });
      }
      if (onComplete) {
        actions.push({
          label: 'Concluir (Já Pago)',
          icon: CheckCircle2Icon,
          onClick: () => {
            onComplete();
            onClose();
          },
        });
      }
      if (onCancel) {
        actions.push({
          label: 'Cancelar',
          icon: XCircleIcon,
          onClick: () => {
            onCancel();
            onClose();
          },
          destructive: true,
        });
      }
      break;

    case 'DONE':
    case 'NO_SHOW':
    case 'CANCELED':
      // Estados finais: Visualizar, Reagendar
      if (onView) {
        actions.push({
          label: 'Visualizar Detalhes',
          icon: EyeIcon,
          onClick: () => {
            onView();
            onClose();
          },
        });
      }
      if (onReschedule) {
        actions.push({
          label: 'Reagendar',
          icon: CalendarPlusIcon,
          onClick: () => {
            onReschedule();
            onClose();
          },
          variant: 'primary',
        });
      }
      break;
  }

  return (
    <div
      ref={menuRef}
      className="fixed z-50 min-w-[200px] rounded-md border bg-popover p-1 shadow-md animate-in fade-in-0 zoom-in-95"
      style={{
        left: `${x}px`,
        top: `${y}px`,
      }}
    >
      {/* Header com nome do cliente */}
      <div className="px-2 py-1.5 text-sm font-semibold border-b mb-1">
        {appointment.customer_name}
      </div>

      {/* Ações */}
      <div className="space-y-0.5">
        {actions.map((action, index) => {
          const Icon = action.icon;
          return (
            <button
              key={index}
              onClick={action.onClick}
              className={cn(
                'relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors',
                'hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground',
                action.destructive && 'text-destructive focus:text-destructive',
                action.variant === 'primary' && 'font-medium'
              )}
            >
              <Icon className="mr-2 h-4 w-4" />
              <span>{action.label}</span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
