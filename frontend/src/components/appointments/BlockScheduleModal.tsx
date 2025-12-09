'use client';

/**
 * NEXO - Modal de Bloqueio de Horário
 * Permite bloquear horários na agenda para profissionais específicos
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { Clock, Lock } from 'lucide-react';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { useProfessionals } from '@/hooks/use-appointments';
import { useCreateBlockedTime } from '@/hooks/use-blocked-times';

// =============================================================================
// SCHEMA
// =============================================================================

const blockScheduleSchema = z.object({
  professional_id: z.string().min(1, 'Selecione um profissional'),
  date: z.string().min(1, 'Selecione uma data'),
  start_time: z.string().min(1, 'Informe o horário de início'),
  end_time: z.string().min(1, 'Informe o horário de término'),
  reason: z.string().min(3, 'Motivo deve ter pelo menos 3 caracteres'),
  is_recurring: z.boolean().default(false),
}).refine((data) => data.end_time > data.start_time, {
  message: 'Horário de término deve ser maior que o de início',
  path: ['end_time'],
});

type BlockScheduleFormData = z.infer<typeof blockScheduleSchema>;

// =============================================================================
// TYPES
// =============================================================================

interface BlockScheduleModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialDate?: Date;
  initialProfessionalId?: string;
  initialStartTime?: string;
  initialEndTime?: string;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function BlockScheduleModal({
  isOpen,
  onClose,
  initialDate,
  initialProfessionalId,
  initialStartTime,
  initialEndTime,
}: BlockScheduleModalProps) {
  const { data: professionals = [] } = useProfessionals();
  const createBlockedTime = useCreateBlockedTime();

  const form = useForm<BlockScheduleFormData>({
    resolver: zodResolver(blockScheduleSchema),
    defaultValues: {
      professional_id: initialProfessionalId || '',
      date: initialDate?.toISOString().split('T')[0] || '',
      start_time: initialStartTime || '08:00',
      end_time: initialEndTime || '09:00',
      reason: '',
      is_recurring: false,
    },
  });

  // Reset form when modal opens with new initial values
  useEffect(() => {
    if (isOpen) {
      form.reset({
        professional_id: initialProfessionalId || '',
        date: initialDate?.toISOString().split('T')[0] || '',
        start_time: initialStartTime || '08:00',
        end_time: initialEndTime || '09:00',
        reason: '',
        is_recurring: false,
      });
    }
  }, [isOpen, initialDate, initialProfessionalId, initialStartTime, initialEndTime, form]);

  const onSubmit = async (data: BlockScheduleFormData) => {
    // Combinar data + horário em ISO 8601
    const startDateTime = new Date(`${data.date}T${data.start_time}:00`);
    const endDateTime = new Date(`${data.date}T${data.end_time}:00`);

    try {
      await createBlockedTime.mutateAsync({
        professional_id: data.professional_id,
        start_time: startDateTime.toISOString(),
        end_time: endDateTime.toISOString(),
        reason: data.reason,
        is_recurring: data.is_recurring,
      });
      
      form.reset();
      onClose();
    } catch {
      // error já tratado no hook
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-foreground">
            <Lock className="h-5 w-5 text-destructive" />
            Bloquear Horário
          </DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* Profissional */}
            <FormField
              control={form.control}
              name="professional_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Profissional</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Selecione o profissional" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="all">Todos os profissionais</SelectItem>
                      {professionals.map((prof) => (
                        <SelectItem key={prof.id} value={prof.id}>
                          {prof.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Data */}
            <FormField
              control={form.control}
              name="date"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Data</FormLabel>
                  <FormControl>
                    <Input type="date" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Horários */}
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="start_time"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Início</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Clock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <Input type="time" className="pl-10" {...field} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="end_time"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Término</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Clock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <Input type="time" className="pl-10" {...field} />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Motivo */}
            <FormField
              control={form.control}
              name="reason"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Motivo</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Ex: Almoço, Reunião, Folga..."
                      className="resize-none"
                      rows={2}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Ações */}
            <div className="flex justify-end gap-2 pt-4 border-t">
              <Button 
                type="button" 
                variant="outline" 
                onClick={onClose}
                disabled={createBlockedTime.isPending}
              >
                Cancelar
              </Button>
              <Button 
                type="submit" 
                variant="destructive"
                disabled={createBlockedTime.isPending}
              >
                <Lock className="h-4 w-4 mr-2" />
                {createBlockedTime.isPending ? 'Bloqueando...' : 'Bloquear Horário'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default BlockScheduleModal;
