/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal Abrir Caixa Component
 *
 * Modal para abrir o caixa com saldo inicial.
 *
 * @author NEXO v2.0
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2, Wallet } from 'lucide-react';
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

import { useAbrirCaixa } from '@/hooks/use-caixa';

// =============================================================================
// VALIDATION SCHEMA
// =============================================================================

const abrirCaixaSchema = z.object({
  saldo_inicial: z
    .string()
    .min(1, 'Informe o saldo inicial')
    .refine((val) => {
      const num = parseFloat(val.replace(',', '.'));
      return !isNaN(num) && num >= 0;
    }, 'Valor deve ser um número válido maior ou igual a 0'),
});

type AbrirCaixaFormData = z.infer<typeof abrirCaixaSchema>;

// =============================================================================
// TYPES
// =============================================================================

interface ModalAbrirCaixaProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ModalAbrirCaixa({ open, onOpenChange }: ModalAbrirCaixaProps) {
  const abrirCaixa = useAbrirCaixa();

  const form = useForm<AbrirCaixaFormData>({
    resolver: zodResolver(abrirCaixaSchema),
    defaultValues: {
      saldo_inicial: '',
    },
  });

  const onSubmit = async (data: AbrirCaixaFormData) => {
    try {
      // Converte para formato decimal (troca vírgula por ponto)
      const saldoFormatado = data.saldo_inicial.replace(',', '.');
      
      await abrirCaixa.mutateAsync({
        saldo_inicial: saldoFormatado,
      });

      form.reset();
      onOpenChange(false);
    } catch (error) {
      console.error('[ModalAbrirCaixa] Erro ao abrir caixa:', error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Wallet className="h-5 w-5" />
            Abrir Caixa
          </DialogTitle>
          <DialogDescription>
            Informe o valor inicial em dinheiro na gaveta para iniciar o dia.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="saldo_inicial"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Saldo Inicial (R$)</FormLabel>
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
                  <FormDescription>
                    Valor em dinheiro disponível na gaveta
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
                disabled={abrirCaixa.isPending}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={abrirCaixa.isPending}>
                {abrirCaixa.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Abrir Caixa
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default ModalAbrirCaixa;
