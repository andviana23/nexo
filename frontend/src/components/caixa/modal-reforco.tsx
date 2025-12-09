/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal Reforço Component
 *
 * Modal para registrar reforço (entrada) no caixa.
 *
 * @author NEXO v2.0
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2, PlusCircle } from 'lucide-react';
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

import { useReforco } from '@/hooks/use-caixa';
import { OrigemReforco, OrigemReforcoLabels } from '@/types/caixa';

// =============================================================================
// VALIDATION SCHEMA
// =============================================================================

const reforcoSchema = z.object({
  valor: z
    .string()
    .min(1, 'Informe o valor')
    .refine((val) => {
      const num = parseFloat(val.replace(',', '.'));
      return !isNaN(num) && num > 0;
    }, 'Valor deve ser maior que zero'),
  origem: z.nativeEnum(OrigemReforco, {
    required_error: 'Selecione a origem',
  }),
  descricao: z
    .string()
    .min(5, 'Descrição deve ter no mínimo 5 caracteres')
    .max(255, 'Descrição muito longa'),
});

type ReforcoFormData = z.infer<typeof reforcoSchema>;

// =============================================================================
// TYPES
// =============================================================================

interface ModalReforcoProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ModalReforco({ open, onOpenChange }: ModalReforcoProps) {
  const reforco = useReforco();

  const form = useForm<ReforcoFormData>({
    resolver: zodResolver(reforcoSchema),
    defaultValues: {
      valor: '',
      origem: undefined,
      descricao: '',
    },
  });

  const onSubmit = async (data: ReforcoFormData) => {
    try {
      const valorFormatado = data.valor.replace(',', '.');

      await reforco.mutateAsync({
        valor: valorFormatado,
        origem: data.origem,
        descricao: data.descricao,
      });

      form.reset();
      onOpenChange(false);
    } catch (error) {
      console.error('[ModalReforco] Erro ao registrar reforço:', error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-blue-600">
            <PlusCircle className="h-5 w-5" />
            Registrar Reforço
          </DialogTitle>
          <DialogDescription>
            Registre uma entrada de dinheiro no caixa.
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
              name="origem"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Origem</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Selecione a origem" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {Object.entries(OrigemReforcoLabels).map(([key, label]) => (
                        <SelectItem key={key} value={key}>
                          {label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    De onde o dinheiro está vindo
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
                      placeholder="Descreva o motivo do reforço..."
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
                disabled={reforco.isPending}
              >
                Cancelar
              </Button>
              <Button
                type="submit"
                disabled={reforco.isPending}
              >
                {reforco.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Registrar Reforço
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default ModalReforco;
