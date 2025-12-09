'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Menu de Ações do Agendamento
 *
 * @component AppointmentActionsMenu
 * @description Menu dropdown com ações baseadas no status do agendamento
 */

import {
    BanIcon,
    CalendarIcon,
    CheckCircle2Icon,
    CreditCardIcon,
    EditIcon,
    HistoryIcon,
    MoreVerticalIcon,
    PlayIcon,
    ScissorsIcon,
    SquareIcon,
    UserCheckIcon,
    UserIcon,
    UserXIcon,
} from 'lucide-react';
import type { ReactNode } from 'react';

import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { AppointmentResponse, AppointmentStatus } from '@/types/appointment';
import { canTransitionTo } from '@/types/appointment';

// =============================================================================
// TYPES
// =============================================================================

export interface AppointmentAction {
  id: string;
  label: string;
  icon: ReactNode;
  onClick: () => void;
  disabled?: boolean;
  destructive?: boolean;
  group: 'status' | 'edit' | 'client';
}

export interface AppointmentActionsMenuProps {
  /** Dados do agendamento */
  appointment: AppointmentResponse;
  /** Trigger customizado (default: ícone de três pontinhos) */
  trigger?: ReactNode;
  /** Alinhamento do menu */
  align?: 'start' | 'center' | 'end';
  /** Callbacks de ações */
  onConfirm?: () => void;
  onCheckIn?: () => void;
  onStartService?: () => void;
  onFinishService?: () => void;
  onComplete?: () => void;
  onNoShow?: () => void;
  onCancel?: () => void;
  onEdit?: () => void;
  onReschedule?: () => void;
  onViewDetails?: () => void;
  onViewHistory?: () => void;
  onViewClient?: () => void;
  /** Classes CSS extras */
  className?: string;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function AppointmentActionsMenu({
  appointment,
  trigger,
  align = 'end',
  onConfirm,
  onCheckIn,
  onStartService,
  onFinishService,
  onComplete,
  onNoShow,
  onCancel,
  onEdit,
  onReschedule,
  onViewDetails,
  onViewHistory,
  onViewClient,
  className,
}: AppointmentActionsMenuProps) {
  const status = appointment.status as AppointmentStatus;

  // Build available actions based on current status
  const actions: AppointmentAction[] = [];

  // ========== STATUS ACTIONS ==========
  
  // CREATED -> CONFIRMED
  if (canTransitionTo(status, 'CONFIRMED') && onConfirm) {
    actions.push({
      id: 'confirm',
      label: 'Confirmar Agendamento',
      icon: <CheckCircle2Icon className="size-4" />,
      onClick: onConfirm,
      group: 'status',
    });
  }

  // -> CHECKED_IN
  if (canTransitionTo(status, 'CHECKED_IN') && onCheckIn) {
    actions.push({
      id: 'check-in',
      label: 'Cliente Chegou',
      icon: <UserCheckIcon className="size-4" />,
      onClick: onCheckIn,
      group: 'status',
    });
  }

  // -> IN_SERVICE
  if (canTransitionTo(status, 'IN_SERVICE') && onStartService) {
    actions.push({
      id: 'start-service',
      label: 'Iniciar Atendimento',
      icon: <PlayIcon className="size-4" />,
      onClick: onStartService,
      group: 'status',
    });
  }

  // -> AWAITING_PAYMENT
  if (canTransitionTo(status, 'AWAITING_PAYMENT') && onFinishService) {
    actions.push({
      id: 'finish-service',
      label: 'Finalizar Atendimento',
      icon: <SquareIcon className="size-4" />,
      onClick: onFinishService,
      group: 'status',
    });
  }

  // -> DONE
  if (canTransitionTo(status, 'DONE') && onComplete) {
    actions.push({
      id: 'complete',
      label: 'Concluir (Pagamento Recebido)',
      icon: <CreditCardIcon className="size-4" />,
      onClick: onComplete,
      group: 'status',
    });
  }

  // -> NO_SHOW
  if (canTransitionTo(status, 'NO_SHOW') && onNoShow) {
    actions.push({
      id: 'no-show',
      label: 'Cliente Não Compareceu',
      icon: <UserXIcon className="size-4" />,
      onClick: onNoShow,
      destructive: true,
      group: 'status',
    });
  }

  // -> CANCELED
  if (canTransitionTo(status, 'CANCELED') && onCancel) {
    actions.push({
      id: 'cancel',
      label: 'Cancelar Agendamento',
      icon: <BanIcon className="size-4" />,
      onClick: onCancel,
      destructive: true,
      group: 'status',
    });
  }

  // ========== EDIT ACTIONS ==========

  if (onEdit && !['DONE', 'NO_SHOW', 'CANCELED'].includes(status)) {
    actions.push({
      id: 'edit',
      label: 'Editar Agendamento',
      icon: <EditIcon className="size-4" />,
      onClick: onEdit,
      group: 'edit',
    });
  }

  if (onReschedule && !['DONE', 'NO_SHOW', 'CANCELED', 'IN_SERVICE', 'AWAITING_PAYMENT'].includes(status)) {
    actions.push({
      id: 'reschedule',
      label: 'Reagendar',
      icon: <CalendarIcon className="size-4" />,
      onClick: onReschedule,
      group: 'edit',
    });
  }

  // ========== CLIENT ACTIONS ==========

  if (onViewDetails) {
    actions.push({
      id: 'view-details',
      label: 'Ver Detalhes',
      icon: <ScissorsIcon className="size-4" />,
      onClick: onViewDetails,
      group: 'client',
    });
  }

  if (onViewHistory) {
    actions.push({
      id: 'view-history',
      label: 'Histórico de Alterações',
      icon: <HistoryIcon className="size-4" />,
      onClick: onViewHistory,
      group: 'client',
    });
  }

  if (onViewClient) {
    actions.push({
      id: 'view-client',
      label: 'Ver Cliente',
      icon: <UserIcon className="size-4" />,
      onClick: onViewClient,
      group: 'client',
    });
  }

  // Group actions
  const statusActions = actions.filter((a) => a.group === 'status');
  const editActions = actions.filter((a) => a.group === 'edit');
  const clientActions = actions.filter((a) => a.group === 'client');

  // Don't render if no actions
  if (actions.length === 0) {
    return null;
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        {trigger || (
          <Button
            variant="ghost"
            size="icon"
            className={className}
            onClick={(e) => e.stopPropagation()}
          >
            <MoreVerticalIcon className="size-4" />
            <span className="sr-only">Abrir menu de ações</span>
          </Button>
        )}
      </DropdownMenuTrigger>
      <DropdownMenuContent align={align} className="w-56">
        {/* Status Actions */}
        {statusActions.length > 0 && (
          <>
            <DropdownMenuLabel>Ações de Status</DropdownMenuLabel>
            <DropdownMenuGroup>
              {statusActions.map((action) => (
                <DropdownMenuItem
                  key={action.id}
                  onClick={(e) => {
                    e.stopPropagation();
                    action.onClick();
                  }}
                  disabled={action.disabled}
                  className={action.destructive ? 'text-destructive focus:text-destructive' : undefined}
                >
                  {action.icon}
                  <span className="ml-2">{action.label}</span>
                </DropdownMenuItem>
              ))}
            </DropdownMenuGroup>
          </>
        )}

        {/* Edit Actions */}
        {editActions.length > 0 && (
          <>
            {statusActions.length > 0 && <DropdownMenuSeparator />}
            <DropdownMenuLabel>Edição</DropdownMenuLabel>
            <DropdownMenuGroup>
              {editActions.map((action) => (
                <DropdownMenuItem
                  key={action.id}
                  onClick={(e) => {
                    e.stopPropagation();
                    action.onClick();
                  }}
                  disabled={action.disabled}
                >
                  {action.icon}
                  <span className="ml-2">{action.label}</span>
                </DropdownMenuItem>
              ))}
            </DropdownMenuGroup>
          </>
        )}

        {/* Client Actions */}
        {clientActions.length > 0 && (
          <>
            {(statusActions.length > 0 || editActions.length > 0) && <DropdownMenuSeparator />}
            <DropdownMenuLabel>Visualização</DropdownMenuLabel>
            <DropdownMenuGroup>
              {clientActions.map((action) => (
                <DropdownMenuItem
                  key={action.id}
                  onClick={(e) => {
                    e.stopPropagation();
                    action.onClick();
                  }}
                  disabled={action.disabled}
                >
                  {action.icon}
                  <span className="ml-2">{action.label}</span>
                </DropdownMenuItem>
              ))}
            </DropdownMenuGroup>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export default AppointmentActionsMenu;
