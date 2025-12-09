'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal de Criação/Edição de Agendamento
 *
 * @component AppointmentModal
 * @description Modal para criar ou editar agendamentos
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    CalendarIcon,
    ClockIcon,
    Loader2Icon,
    ScissorsIcon,
    UserIcon,
} from 'lucide-react';
import { useCallback, useEffect, useMemo } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
    useCreateAppointment,
    useUpdateAppointment,
} from '@/hooks/use-appointments';
import { useServices } from '@/hooks/useServices';
import {
    formatCurrency,
    type AppointmentModalState,
    type AppointmentResponse,
} from '@/types/appointment';

import { CustomerSelector } from './CustomerSelector';
import { ProfessionalSelector } from './ProfessionalSelector';
import { ServiceSelector } from './ServiceSelector';

// =============================================================================
// SCHEMA DE VALIDAÇÃO
// =============================================================================

const appointmentFormSchema = z.object({
  professional_id: z.string().min(1, 'Selecione um barbeiro'),
  customer_id: z.string().min(1, 'Selecione um cliente'),
  service_ids: z.array(z.string()).min(1, 'Selecione pelo menos um serviço'),
  start_date: z.string().min(1, 'Selecione a data'),
  start_time: z.string().min(1, 'Selecione o horário'),
  notes: z.string().optional(),
}).refine(() => {
  // Validação adicional: duração será calculada automaticamente
  // baseada nos serviços selecionados
  return true;
}, {
  message: 'Configuração de agendamento inválida',
});

type AppointmentFormValues = z.infer<typeof appointmentFormSchema>;

// =============================================================================
// TYPES
// =============================================================================

interface AppointmentModalProps {
  /** Estado do modal */
  state: AppointmentModalState;
  /** Callback para fechar o modal */
  onClose: () => void;
  /** Callback após salvar com sucesso */
  onSuccess?: (appointment: AppointmentResponse) => void;
}

// =============================================================================
// HELPERS
// =============================================================================

function formatTimeFromISO(isoString: string): string {
  const date = new Date(isoString);
  return format(date, 'HH:mm');
}

function formatDateFromISO(isoString: string): string {
  const date = new Date(isoString);
  return format(date, 'yyyy-MM-dd');
}

function combineDateAndTime(dateStr: string, timeStr: string): string {
  return new Date(`${dateStr}T${timeStr}:00`).toISOString();
}

// =============================================================================
// COMPONENT
// =============================================================================

export function AppointmentModal({
  state,
  onClose,
  onSuccess,
}: AppointmentModalProps) {
  const { isOpen, mode, appointment, initialDate, initialProfessionalId } = state;

  // Mutations
  const createAppointment = useCreateAppointment();
  const updateAppointment = useUpdateAppointment();

  const isLoading = createAppointment.isPending || updateAppointment.isPending;
  const isViewMode = mode === 'view';

  // Buscar serviços para calcular duração total
  const { data: servicesData } = useServices({ apenas_ativos: true });

  // Form
  const form = useForm<AppointmentFormValues>({
    resolver: zodResolver(appointmentFormSchema),
    defaultValues: {
      professional_id: '',
      customer_id: '',
      service_ids: [],
      start_date: '',
      start_time: '',
      notes: '',
    },
  });

  // Watch dos serviços selecionados para calcular duração
  const selectedServiceIds = useWatch({ control: form.control, name: 'service_ids' });

  // Calcular duração total dos serviços selecionados
  const totalDuration = useMemo(() => {
    if (!servicesData?.servicos || selectedServiceIds.length === 0) return 0;

    return servicesData.servicos
      .filter((s) => selectedServiceIds.includes(s.id))
      .reduce((sum, service) => sum + service.duracao, 0);
  }, [selectedServiceIds, servicesData]);

  // Preencher formulário quando abrir para editar
  useEffect(() => {
    if (mode === 'edit' && appointment) {
      form.reset({
        professional_id: appointment.professional_id,
        customer_id: appointment.customer_id,
        service_ids: (appointment.services || []).map((s) => s.service_id),
        start_date: formatDateFromISO(appointment.start_time),
        start_time: formatTimeFromISO(appointment.start_time),
        notes: appointment.notes || '',
      });
    } else if (mode === 'create') {
      form.reset({
        professional_id: initialProfessionalId || '',
        customer_id: '',
        service_ids: [],
        start_date: initialDate ? format(initialDate, 'yyyy-MM-dd') : format(new Date(), 'yyyy-MM-dd'),
        start_time: initialDate ? format(initialDate, 'HH:mm') : '09:00',
        notes: '',
      });
    }
  }, [mode, appointment, initialDate, initialProfessionalId, form]);

  // Submit handler
  const onSubmit = useCallback(
    async (values: AppointmentFormValues) => {
      const startTime = combineDateAndTime(values.start_date, values.start_time);

      if (mode === 'create') {
        createAppointment.mutate(
          {
            professional_id: values.professional_id,
            customer_id: values.customer_id,
            service_ids: values.service_ids,
            start_time: startTime,
            notes: values.notes,
          },
          {
            onSuccess: (data) => {
              onSuccess?.(data);
              onClose();
            },
          }
        );
      } else if (mode === 'edit' && appointment) {
        updateAppointment.mutate(
          {
            id: appointment.id,
            data: {
              new_start_time: startTime,
              professional_id: values.professional_id,
            },
          },
          {
            onSuccess: (data) => {
              onSuccess?.(data);
              onClose();
            },
          }
        );
      }
    },
    [mode, appointment, createAppointment, updateAppointment, onSuccess, onClose]
  );

  // Título do modal
  const title = useMemo(() => {
    switch (mode) {
      case 'create':
        return 'Novo Agendamento';
      case 'edit':
        return 'Editar Agendamento';
      case 'view':
        return 'Detalhes do Agendamento';
      default:
        return 'Agendamento';
    }
  }, [mode]);

  // ==========================================================================
  // VIEW MODE RENDER
  // ==========================================================================

  if (isViewMode && appointment) {
    return (
      <Dialog open={isOpen} onOpenChange={() => onClose()}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <CalendarIcon className="size-5 text-primary" />
              {title}
            </DialogTitle>
            <DialogDescription>
              Informações do agendamento
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {/* Cliente */}
            <div className="flex items-start gap-3">
              <UserIcon className="size-5 text-muted-foreground mt-0.5" />
              <div>
                <p className="font-medium">{appointment.customer_name}</p>
                {appointment.customer_phone && (
                  <p className="text-sm text-muted-foreground">
                    {appointment.customer_phone}
                  </p>
                )}
              </div>
            </div>

            {/* Barbeiro */}
            <div className="flex items-start gap-3">
              <ScissorsIcon className="size-5 text-muted-foreground mt-0.5" />
              <div>
                <p className="font-medium">{appointment.professional_name}</p>
                <p className="text-sm text-muted-foreground">Barbeiro</p>
              </div>
            </div>

            {/* Data e Hora */}
            <div className="flex items-start gap-3">
              <ClockIcon className="size-5 text-muted-foreground mt-0.5" />
              <div>
                <p className="font-medium">
                  {format(new Date(appointment.start_time), "EEEE, d 'de' MMMM", {
                    locale: ptBR,
                  })}
                </p>
                <p className="text-sm text-muted-foreground">
                  {format(new Date(appointment.start_time), 'HH:mm')} -{' '}
                  {format(new Date(appointment.end_time), 'HH:mm')}
                </p>
              </div>
            </div>

            {/* Serviços */}
            <div className="border-t pt-4">
              <p className="text-sm font-medium mb-2">Serviços</p>
              <div className="space-y-2">
                {(appointment.services || []).map((service) => (
                  <div
                    key={service.service_id}
                    className="flex items-center justify-between text-sm"
                  >
                    <span>{service.service_name}</span>
                    <span className="text-muted-foreground">
                      {formatCurrency(service.price)}
                    </span>
                  </div>
                ))}
                <div className="flex items-center justify-between font-medium pt-2 border-t">
                  <span>Total</span>
                  <span>{formatCurrency(appointment.total_price)}</span>
                </div>
              </div>
            </div>

            {/* Observações */}
            {appointment.notes && (
              <div className="border-t pt-4">
                <p className="text-sm font-medium mb-1">Observações</p>
                <p className="text-sm text-muted-foreground">{appointment.notes}</p>
              </div>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={onClose}>
              Fechar
            </Button>
            <Button
              onClick={() => {
                onClose();
                // Re-open em modo edição (precisa ser gerenciado pelo parent)
              }}
            >
              Editar
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  // ==========================================================================
  // CREATE/EDIT MODE RENDER
  // ==========================================================================

  return (
    <Dialog open={isOpen} onOpenChange={() => onClose()}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <CalendarIcon className="size-5 text-primary" />
            {title}
          </DialogTitle>
          <DialogDescription>
            {mode === 'create'
              ? 'Preencha os dados para criar um novo agendamento'
              : 'Atualize os dados do agendamento'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* Barbeiro */}
            <FormField
              control={form.control}
              name="professional_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Barbeiro</FormLabel>
                  <FormControl>
                    <ProfessionalSelector
                      value={field.value}
                      onChange={field.onChange}
                      disabled={isLoading}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Cliente */}
            {mode === 'create' && (
              <FormField
                control={form.control}
                name="customer_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Cliente</FormLabel>
                    <FormControl>
                      <CustomerSelector
                        value={field.value}
                        onChange={field.onChange}
                        disabled={isLoading}
                      />
                    </FormControl>
                    <FormDescription>
                      Busque pelo nome ou telefone do cliente
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {/* Serviços */}
            <FormField
              control={form.control}
              name="service_ids"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Serviços</FormLabel>
                  <FormControl>
                    <ServiceSelector
                      value={field.value}
                      onChange={field.onChange}
                      disabled={isLoading}
                    />
                  </FormControl>
                  <FormDescription>
                    Selecione um ou mais serviços
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Exibir duração total calculada */}
            {totalDuration > 0 && (
              <div className="flex items-center gap-2 rounded-md border border-muted bg-muted/50 p-3">
                <ClockIcon className="h-4 w-4 text-muted-foreground" />
                <div className="flex-1">
                  <p className="text-sm font-medium">
                    Duração total dos serviços
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {totalDuration < 60 
                      ? `${totalDuration} minutos`
                      : `${Math.floor(totalDuration / 60)}h ${totalDuration % 60 > 0 ? `${totalDuration % 60}min` : ''}`
                    }
                  </p>
                </div>
                <Badge variant="secondary" className="font-mono">
                  {totalDuration}min
                </Badge>
              </div>
            )}

            {/* Data e Hora */}
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="start_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Data</FormLabel>
                    <FormControl>
                      <Input
                        type="date"
                        {...field}
                        disabled={isLoading}
                        min={format(new Date(), 'yyyy-MM-dd')}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="start_time"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Horário</FormLabel>
                    <FormControl>
                      <Input
                        type="time"
                        {...field}
                        disabled={isLoading}
                        step="600" // 10 minutos
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Observações */}
            <FormField
              control={form.control}
              name="notes"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Observações</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Observações opcionais..."
                      className="resize-none"
                      {...field}
                      disabled={isLoading}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className="pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={onClose}
                disabled={isLoading}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2Icon className="size-4 animate-spin" />}
                {mode === 'create' ? 'Criar Agendamento' : 'Salvar Alterações'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default AppointmentModal;
