/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal Sangria Component
 *
 * Modal para registrar sangria (retirada) do caixa.
 *
 * @author NEXO v2.0
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2, MinusCircle } from 'lucide-react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';

import { useSangria } from '@/hooks/use-caixa';
import { DestinoSangria, DestinoSangriaLabels } from '@/types/caixa';

// =============================================================================
// VALIDATION SCHEMA
// =============================================================================

const sangriaSchema = z.object({
  valor: z
    .string()
    .min(1, 'Informe o valor')
    .refine((val) => {
      const num = parseFloat(val.replace(',', '.'));
      return !isNaN(num) && num > 0;
    }, 'Valor deve ser maior que zero'),
  destino: z.nativeEnum(DestinoSangria, {
    required_error: 'Selecione o destino',
  }),
  descricao: z
    .string()
    .min(5, 'Descrição deve ter no mínimo 5 caracteres')
    .max(255, 'Descrição muito longa'),
});

type SangriaFormData = z.infer<typeof sangriaSchema>;

// =============================================================================
// TYPES
// =============================================================================

interface ModalSangriaProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ModalSangria({ open, onOpenChange }: ModalSangriaProps) {
  const sangria = useSangria();

  const form = useForm<SangriaFormData>({
    resolver: zodResolver(sangriaSchema),
    defaultValues: {
      valor: '',
      destino: undefined,
      descricao: '',
    },
  });

  const onSubmit = async (data: SangriaFormData) => {
    try {
      const valorFormatado = data.valor.replace(',', '.');

      await sangria.mutateAsync({
        valor: valorFormatado,
        destino: data.destino,
        descricao: data.descricao,
      });

      form.reset();
      onOpenChange(false);
    } catch (error) {
      console.error('[ModalSangria] Erro ao registrar sangria:', error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-red-600">
            <MinusCircle className="h-5 w-5" />
            Registrar Sangria
          </DialogTitle>
          <DialogDescription>
            Registre uma retirada de dinheiro do caixa.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="valor"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Valor (R$)</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                        R$
                      </span>
                      <Input
                        {...field}
                        type="text"
                        inputMode="decimal"
                        placeholder="0,00"
                        className="pl-10"
                        autoFocus
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="destino"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Destino</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Selecione o destino" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {Object.entries(DestinoSangriaLabels).map(([key, label]) => (
                        <SelectItem key={key} value={key}>
                          {label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    Para onde o dinheiro será destinado
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="descricao"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Descrição</FormLabel>
                  <FormControl>
                    <Textarea
                      {...field}
                      placeholder="Descreva o motivo da sangria..."
                      rows={3}
                    />
                  </FormControl>
                  <FormDescription>
                    Mínimo 5 caracteres
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={sangria.isPending}
              >
                Cancelar
              </Button>
              <Button
                type="submit"
                variant="destructive"
                disabled={sangria.isPending}
              >
                {sangria.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Registrar Sangria
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default ModalSangria;
