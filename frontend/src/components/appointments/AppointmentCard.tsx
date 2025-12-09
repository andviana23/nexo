'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Card de Agendamento Compacto
 *
 * @component AppointmentCard
 * @description Card para exibir resumo de um agendamento
 */

import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    CalendarIcon,
    CheckCircle2Icon,
    ClockIcon,
    CreditCardIcon,
    FileTextIcon,
    MoreVerticalIcon,
    PhoneIcon,
    ScissorsIcon,
    UserCheckIcon,
    UserIcon,
    UserXIcon,
    XCircleIcon,
} from 'lucide-react';
import { useMemo } from 'react';

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';
import type { AppointmentResponse, AppointmentStatus } from '@/types/appointment';
import { formatCurrency } from '@/types/appointment';

// =============================================================================
// TYPES
// =============================================================================

interface AppointmentCardProps {
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
  /** Callback para fechar comanda (AWAITING_PAYMENT) */
  onCloseCommand?: () => void;
  /** Callback para marcar no-show */
  onNoShow?: () => void;
  /** Variante de exibição */
  variant?: 'default' | 'compact';
  /** Classe CSS adicional */
  className?: string;
}

// =============================================================================
// STATUS CONFIG
// =============================================================================

const STATUS_CONFIG: Record<
  AppointmentStatus,
  { label: string; color: string; icon: React.ComponentType<{ className?: string }> }
> = {
  CREATED: {
    label: 'Criado',
    color: 'bg-amber-100 text-amber-900 border-amber-400',
    icon: CalendarIcon,
  },
  CONFIRMED: {
    label: 'Confirmado',
    color: 'bg-green-100 text-green-900 border-green-400',
    icon: CheckCircle2Icon,
  },
  CHECKED_IN: {
    label: 'Cliente Chegou',
    color: 'bg-blue-100 text-blue-900 border-blue-400',
    icon: UserCheckIcon,
  },
  IN_SERVICE: {
    label: 'Em Atendimento',
    color: 'bg-purple-100 text-purple-900 border-purple-400',
    icon: ScissorsIcon,
  },
  AWAITING_PAYMENT: {
    label: 'Aguardando Pagamento',
    color: 'bg-orange-100 text-orange-900 border-orange-400',
    icon: CreditCardIcon,
  },
  DONE: {
    label: 'Concluído',
    color: 'bg-slate-100 text-slate-700 border-slate-400',
    icon: CheckCircle2Icon,
  },
  NO_SHOW: {
    label: 'Não Compareceu',
    color: 'bg-red-100 text-red-900 border-red-400',
    icon: UserXIcon,
  },
  CANCELED: {
    label: 'Cancelado',
    color: 'bg-slate-200 text-slate-600 border-slate-500',
    icon: XCircleIcon,
  },
};

// =============================================================================
// HELPERS
// =============================================================================

function getInitials(name: string): string {
  return name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase();
}

// =============================================================================
// COMPONENT
// =============================================================================

export function AppointmentCard({
  appointment,
  onClick,
  onEdit,
  onCancel,
  onConfirm,
  onCheckIn,
  onStartService,
  onFinishService,
  onComplete,
  onCloseCommand,
  onNoShow,
  variant = 'default',
  className,
}: AppointmentCardProps) {
  const statusConfig = STATUS_CONFIG[appointment.status];
  const StatusIcon = statusConfig.icon;

  // Formatar data e hora
  const formattedDate = useMemo(() => {
    return format(new Date(appointment.start_time), "EEE, d 'de' MMM", {
      locale: ptBR,
    });
  }, [appointment.start_time]);

  const formattedTime = useMemo(() => {
    return `${format(new Date(appointment.start_time), 'HH:mm')} - ${format(
      new Date(appointment.end_time),
      'HH:mm'
    )}`;
  }, [appointment.start_time, appointment.end_time]);

  // Lista de serviços
  const serviceNames = useMemo(() => {
    return appointment.services.map((s) => s.service_name).join(', ');
  }, [appointment.services]);

  // Ações disponíveis baseadas no status
  const availableActions = useMemo(() => {
    const actions: { label: string; onClick: () => void; destructive?: boolean }[] = [];

    switch (appointment.status) {
      case 'CREATED':
        if (onConfirm) actions.push({ label: 'Confirmar', onClick: onConfirm });
        if (onCheckIn) actions.push({ label: 'Cliente Chegou', onClick: onCheckIn });
        if (onNoShow) actions.push({ label: 'Não Compareceu', onClick: onNoShow, destructive: true });
        if (onCancel) actions.push({ label: 'Cancelar', onClick: onCancel, destructive: true });
        break;
      case 'CONFIRMED':
        if (onCheckIn) actions.push({ label: 'Cliente Chegou', onClick: onCheckIn });
        if (onStartService) actions.push({ label: 'Iniciar Atendimento', onClick: onStartService });
        if (onNoShow) actions.push({ label: 'Não Compareceu', onClick: onNoShow, destructive: true });
        if (onCancel) actions.push({ label: 'Cancelar', onClick: onCancel, destructive: true });
        break;
      case 'CHECKED_IN':
        if (onStartService) actions.push({ label: 'Iniciar Atendimento', onClick: onStartService });
        if (onNoShow) actions.push({ label: 'Não Compareceu', onClick: onNoShow, destructive: true });
        if (onCancel) actions.push({ label: 'Cancelar', onClick: onCancel, destructive: true });
        break;
      case 'IN_SERVICE':
        if (onFinishService) actions.push({ label: 'Finalizar Atendimento', onClick: onFinishService });
        if (onCancel) actions.push({ label: 'Cancelar', onClick: onCancel, destructive: true });
        break;
      case 'AWAITING_PAYMENT':
        if (onCloseCommand) actions.push({ label: 'Fechar Comanda', onClick: onCloseCommand });
        if (onComplete) actions.push({ label: 'Concluir (Pago)', onClick: onComplete });
        if (onCancel) actions.push({ label: 'Cancelar', onClick: onCancel, destructive: true });
        break;
    }

    return actions;
  }, [appointment.status, onConfirm, onCheckIn, onCancel, onStartService, onFinishService, onComplete, onCloseCommand, onNoShow]);

  // ==========================================================================
  // RENDER COMPACT
  // ==========================================================================

  if (variant === 'compact') {
    return (
      <div
        className={cn(
          'flex items-center gap-3 rounded-lg border p-3 cursor-pointer hover:bg-accent/50 transition-colors',
          className
        )}
        onClick={onClick}
      >
        <Avatar className="size-10">
          <AvatarFallback className="text-sm">
            {getInitials(appointment.customer_name)}
          </AvatarFallback>
        </Avatar>

        <div className="flex-1 min-w-0">
          <p className="font-medium truncate">{appointment.customer_name}</p>
          <p className="text-sm text-muted-foreground truncate">{serviceNames}</p>
        </div>

        <div className="text-right">
          <p className="text-sm font-medium">{formattedTime}</p>
          <Badge
            variant="outline"
            className={cn('text-xs', statusConfig.color)}
          >
            {statusConfig.label}
          </Badge>
        </div>
      </div>
    );
  }

  // ==========================================================================
  // RENDER DEFAULT
  // ==========================================================================

  return (
    <Card
      className={cn(
        'cursor-pointer hover:shadow-md transition-shadow',
        className
      )}
      onClick={onClick}
      onContextMenu={(e) => {
        // Prevenir menu nativo do navegador
        e.preventDefault();
        // Simular clique no botão de menu (se houver ações disponíveis)
        if (availableActions.length > 0) {
          const menuButton = e.currentTarget.querySelector('[data-menu-trigger]');
          if (menuButton instanceof HTMLElement) {
            menuButton.click();
          }
        }
      }}
    >
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <Avatar className="size-12">
              <AvatarFallback>
                {getInitials(appointment.customer_name)}
              </AvatarFallback>
            </Avatar>
            <div>
              <p className="font-semibold">{appointment.customer_name}</p>
              {appointment.customer_phone && (
                <p className="text-sm text-muted-foreground flex items-center gap-1">
                  <PhoneIcon className="size-3" />
                  {appointment.customer_phone}
                </p>
              )}
            </div>
          </div>

          {/* Menu de ações */}
          {availableActions.length > 0 && (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8"
                  data-menu-trigger
                  onClick={(e) => e.stopPropagation()}
                >
                  <MoreVerticalIcon className="size-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {onEdit && (
                  <>
                    <DropdownMenuItem onClick={(e) => { e.stopPropagation(); onEdit(); }}>
                      Editar
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                  </>
                )}
                {availableActions.map((action, index) => (
                  <DropdownMenuItem
                    key={index}
                    onClick={(e) => { e.stopPropagation(); action.onClick(); }}
                    className={action.destructive ? 'text-destructive' : undefined}
                  >
                    {action.label}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {/* Status e Comanda */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Badge
              variant="outline"
              className={cn('gap-1', statusConfig.color)}
            >
              <StatusIcon className="size-3" />
              {statusConfig.label}
            </Badge>
            {/* Indicador de comanda */}
            {appointment.command_id && (
              <Badge variant="secondary" className="gap-1 text-xs">
                <FileTextIcon className="size-3" />
                Comanda
              </Badge>
            )}
          </div>
          <span className="font-semibold text-lg">
            {formatCurrency(appointment.total_price)}
          </span>
        </div>

        {/* Botão de Fechar Comanda inline (AWAITING_PAYMENT) */}
        {appointment.status === 'AWAITING_PAYMENT' && onCloseCommand && (
          <Button
            variant="default"
            size="sm"
            className="w-full gap-2"
            onClick={(e) => {
              e.stopPropagation();
              onCloseCommand();
            }}
          >
            <CreditCardIcon className="size-4" />
            Fechar Comanda
          </Button>
        )}

        {/* Data e Hora */}
        <div className="flex items-center gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <CalendarIcon className="size-4" />
            <span>{formattedDate}</span>
          </div>
          <div className="flex items-center gap-1">
            <ClockIcon className="size-4" />
            <span>{formattedTime}</span>
          </div>
        </div>

        {/* Barbeiro */}
        <div className="flex items-center gap-2 text-sm">
          <UserIcon className="size-4 text-muted-foreground" />
          <span>{appointment.professional_name}</span>
        </div>

        {/* Serviços */}
        <div className="flex items-center gap-2 text-sm">
          <ScissorsIcon className="size-4 text-muted-foreground" />
          <span className="truncate">{serviceNames}</span>
        </div>

        {/* Observações */}
        {appointment.notes && (
          <p className="text-sm text-muted-foreground italic truncate">
            &quot;{appointment.notes}&quot;
          </p>
        )}
      </CardContent>
    </Card>
  );
}

export default AppointmentCard;
