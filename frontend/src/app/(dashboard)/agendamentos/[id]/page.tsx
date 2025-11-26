'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Detalhes do Agendamento
 *
 * Exibe informações detalhadas de um agendamento específico
 */

import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    ArrowLeftIcon,
    CalendarIcon,
    CheckCircle2Icon,
    ClockIcon,
    Loader2Icon,
    MoreVerticalIcon,
    PhoneIcon,
    ScissorsIcon,
    UserIcon,
    XCircleIcon,
} from 'lucide-react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useCallback, useState } from 'react';

import { AppointmentModal } from '@/components/appointments';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Separator } from '@/components/ui/separator';
import {
    useAppointment,
    useCancelAppointment,
    useUpdateAppointmentStatus,
} from '@/hooks/use-appointments';
import type { AppointmentModalState, AppointmentStatus } from '@/types/appointment';

// =============================================================================
// STATUS CONFIG
// =============================================================================

const STATUS_CONFIG: Record<
  AppointmentStatus,
  { label: string; color: string; icon: React.ComponentType<{ className?: string }> }
> = {
  CREATED: {
    label: 'Pendente',
    color: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    icon: ClockIcon,
  },
  CONFIRMED: {
    label: 'Confirmado',
    color: 'bg-blue-100 text-blue-800 border-blue-200',
    icon: CheckCircle2Icon,
  },
  IN_SERVICE: {
    label: 'Em Atendimento',
    color: 'bg-purple-100 text-purple-800 border-purple-200',
    icon: ScissorsIcon,
  },
  DONE: {
    label: 'Concluído',
    color: 'bg-green-100 text-green-800 border-green-200',
    icon: CheckCircle2Icon,
  },
  NO_SHOW: {
    label: 'Não Compareceu',
    color: 'bg-orange-100 text-orange-800 border-orange-200',
    icon: XCircleIcon,
  },
  CANCELED: {
    label: 'Cancelado',
    color: 'bg-red-100 text-red-800 border-red-200',
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

function formatPrice(cents: number): string {
  return `R$ ${(cents / 100).toFixed(2)}`;
}

// =============================================================================
// PAGE COMPONENT
// =============================================================================

export default function AppointmentDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  // Estado do modal de edição
  const [modalState, setModalState] = useState<AppointmentModalState>({
    isOpen: false,
    mode: 'edit',
  });

  // Queries e Mutations
  const { data: appointment, isLoading, isError } = useAppointment(id);
  const updateStatus = useUpdateAppointmentStatus();
  const cancelAppointment = useCancelAppointment();

  // Handlers de status
  const handleConfirm = useCallback(() => {
    updateStatus.mutate({ id, data: { status: 'CONFIRMED' } });
  }, [id, updateStatus]);

  const handleStartService = useCallback(() => {
    updateStatus.mutate({ id, data: { status: 'IN_SERVICE' } });
  }, [id, updateStatus]);

  const handleFinish = useCallback(() => {
    updateStatus.mutate({ id, data: { status: 'DONE' } });
  }, [id, updateStatus]);

  const handleNoShow = useCallback(() => {
    updateStatus.mutate({ id, data: { status: 'NO_SHOW' } });
  }, [id, updateStatus]);

  const handleCancel = useCallback(() => {
    if (confirm('Tem certeza que deseja cancelar este agendamento?')) {
      cancelAppointment.mutate({ id }, {
        onSuccess: () => router.push('/agendamentos'),
      });
    }
  }, [id, cancelAppointment, router]);

  const handleEdit = useCallback(() => {
    if (appointment) {
      setModalState({
        isOpen: true,
        mode: 'edit',
        appointment,
      });
    }
  }, [appointment]);

  const handleCloseModal = useCallback(() => {
    setModalState((prev) => ({ ...prev, isOpen: false }));
  }, []);

  // Loading state
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2Icon className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  // Error state
  if (isError || !appointment) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] gap-4">
        <XCircleIcon className="size-12 text-destructive" />
        <h2 className="text-lg font-semibold">Agendamento não encontrado</h2>
        <Button variant="outline" asChild>
          <Link href="/agendamentos">
            <ArrowLeftIcon className="size-4 mr-2" />
            Voltar
          </Link>
        </Button>
      </div>
    );
  }

  const statusConfig = STATUS_CONFIG[appointment.status];
  const StatusIcon = statusConfig.icon;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/agendamentos">
              <ArrowLeftIcon className="size-5" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">
              Detalhes do Agendamento
            </h1>
            <p className="text-muted-foreground">
              {format(new Date(appointment.start_time), "EEEE, d 'de' MMMM 'de' yyyy", {
                locale: ptBR,
              })}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* Status Badge */}
          <Badge variant="outline" className={`gap-1 ${statusConfig.color}`}>
            <StatusIcon className="size-3" />
            {statusConfig.label}
          </Badge>

          {/* Ações */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="icon">
                <MoreVerticalIcon className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={handleEdit}>
                Editar Agendamento
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              
              {appointment.status === 'CREATED' && (
                <DropdownMenuItem onClick={handleConfirm}>
                  Confirmar
                </DropdownMenuItem>
              )}
              {appointment.status === 'CONFIRMED' && (
                <>
                  <DropdownMenuItem onClick={handleStartService}>
                    Iniciar Atendimento
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={handleNoShow} className="text-orange-600">
                    Marcar Não Compareceu
                  </DropdownMenuItem>
                </>
              )}
              {appointment.status === 'IN_SERVICE' && (
                <DropdownMenuItem onClick={handleFinish}>
                  Finalizar Atendimento
                </DropdownMenuItem>
              )}
              
              {['CREATED', 'CONFIRMED'].includes(appointment.status) && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={handleCancel} className="text-destructive">
                    Cancelar Agendamento
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Card do Cliente */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <UserIcon className="size-5 text-primary" />
              Cliente
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <Avatar className="size-16">
                <AvatarFallback className="text-lg">
                  {getInitials(appointment.customer.name)}
                </AvatarFallback>
              </Avatar>
              <div>
                <p className="text-lg font-semibold">{appointment.customer.name}</p>
                {appointment.customer.phone && (
                  <p className="text-muted-foreground flex items-center gap-1">
                    <PhoneIcon className="size-4" />
                    {appointment.customer.phone}
                  </p>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Card do Profissional */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <ScissorsIcon className="size-5 text-primary" />
              Barbeiro
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <Avatar className="size-16">
                <AvatarImage
                  src={appointment.professional.avatar_url}
                  alt={appointment.professional.name}
                />
                <AvatarFallback className="text-lg">
                  {getInitials(appointment.professional.name)}
                </AvatarFallback>
              </Avatar>
              <div>
                <p className="text-lg font-semibold">{appointment.professional.name}</p>
                <p className="text-muted-foreground">Profissional</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Card de Data e Hora */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <CalendarIcon className="size-5 text-primary" />
              Data e Hora
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Data</p>
              <p className="font-semibold">
                {format(new Date(appointment.start_time), "EEEE, d 'de' MMMM 'de' yyyy", {
                  locale: ptBR,
                })}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Horário</p>
              <p className="font-semibold flex items-center gap-1">
                <ClockIcon className="size-4" />
                {format(new Date(appointment.start_time), 'HH:mm')} -{' '}
                {format(new Date(appointment.end_time), 'HH:mm')}
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Card de Serviços */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <ScissorsIcon className="size-5 text-primary" />
              Serviços
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {appointment.services.map((service) => (
                <div
                  key={service.id}
                  className="flex items-center justify-between"
                >
                  <div>
                    <p className="font-medium">{service.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {service.duration} min
                    </p>
                  </div>
                  <p className="font-semibold">
                    {formatPrice(service.price)}
                  </p>
                </div>
              ))}
              <Separator />
              <div className="flex items-center justify-between text-lg">
                <p className="font-semibold">Total</p>
                <p className="font-bold text-primary">
                  {formatPrice(appointment.total_price)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Observações */}
      {appointment.notes && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Observações</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">{appointment.notes}</p>
          </CardContent>
        </Card>
      )}

      {/* Timeline de Status (futuro) */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Histórico</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground">
            <p>Criado em: {format(new Date(appointment.created_at), "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })}</p>
            <p>Última atualização: {format(new Date(appointment.updated_at), "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })}</p>
          </div>
        </CardContent>
      </Card>

      {/* Modal de Edição */}
      <AppointmentModal
        state={modalState}
        onClose={handleCloseModal}
      />
    </div>
  );
}
