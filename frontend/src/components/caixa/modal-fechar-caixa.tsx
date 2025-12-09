/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal Fechar Caixa Component
 *
 * Modal para fechar o caixa informando saldo real.
 *
 * @author NEXO v2.0
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { AlertCircle, Loader2, Lock } from 'lucide-react';
import { useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
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
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';

import { useFecharCaixa } from '@/hooks/use-caixa';
import type { CaixaDiarioResponse } from '@/types/caixa';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: string | undefined) => {
  if (!value) return 'R$ 0,00';
  const num = parseFloat(value);
  if (isNaN(num)) return 'R$ 0,00';
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
};

// =============================================================================
// VALIDATION SCHEMA
// =============================================================================

const fecharCaixaSchema = z.object({
  saldo_real: z
    .string()
    .min(1, 'Informe o saldo real')
    .refine((val) => {
      const num = parseFloat(val.replace(',', '.'));
      return !isNaN(num) && num >= 0;
    }, 'Valor deve ser um número válido maior ou igual a 0'),
  justificativa: z.string().optional(),
});

type FecharCaixaFormData = z.infer<typeof fecharCaixaSchema>;

// =============================================================================
// TYPES
// =============================================================================

interface ModalFecharCaixaProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  caixa: CaixaDiarioResponse | undefined;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ModalFecharCaixa({
  open,
  onOpenChange,
  caixa,
}: ModalFecharCaixaProps) {
  const fecharCaixa = useFecharCaixa();

  const form = useForm<FecharCaixaFormData>({
    resolver: zodResolver(fecharCaixaSchema),
    defaultValues: {
      saldo_real: '',
      justificativa: '',
    },
  });

  const saldoRealValue = useWatch({ control: form.control, name: 'saldo_real' });

  // Calcula divergência em tempo real
  const calcularDivergencia = () => {
    if (!caixa?.saldo_esperado || !saldoRealValue) return null;

    const saldoEsperado = parseFloat(caixa.saldo_esperado);
    const saldoReal = parseFloat(saldoRealValue.replace(',', '.'));

    if (isNaN(saldoEsperado) || isNaN(saldoReal)) return null;

    const divergencia = saldoReal - saldoEsperado;
    return divergencia;
  };

  const divergencia = calcularDivergencia();
  const temDivergencia = divergencia !== null && divergencia !== 0;

  const onSubmit = async (data: FecharCaixaFormData) => {
    // Verifica se precisa de justificativa
    if (temDivergencia && !data.justificativa?.trim()) {
      form.setError('justificativa', {
        message: 'Justificativa obrigatória quando há divergência',
      });
      return;
    }

    try {
      const valorFormatado = data.saldo_real.replace(',', '.');

      await fecharCaixa.mutateAsync({
        saldo_real: valorFormatado,
        justificativa: data.justificativa || undefined,
      });

      form.reset();
      onOpenChange(false);
    } catch (error) {
      console.error('[ModalFecharCaixa] Erro ao fechar caixa:', error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Lock className="h-5 w-5" />
            Fechar Caixa
          </DialogTitle>
          <DialogDescription>
            Informe o valor real em dinheiro na gaveta para encerrar o dia.
          </DialogDescription>
        </DialogHeader>

        {/* Resumo do Caixa */}
        {caixa && (
          <div className="bg-muted rounded-lg p-4 space-y-2">
            <h4 className="font-medium text-sm">Resumo do Caixa</h4>
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div>Saldo Inicial:</div>
              <div className="text-right font-medium">
                {formatCurrency(caixa.saldo_inicial)}
              </div>
              <div>Total Entradas:</div>
              <div className="text-right font-medium text-green-600">
                +{formatCurrency(caixa.total_entradas)}
              </div>
              <div>Total Saídas:</div>
              <div className="text-right font-medium text-red-600">
                -{formatCurrency(caixa.total_saidas)}
              </div>
              <Separator className="col-span-2" />
              <div className="font-semibold">Saldo Esperado:</div>
              <div className="text-right font-bold text-primary">
                {formatCurrency(caixa.saldo_esperado)}
              </div>
            </div>
          </div>
        )}

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="saldo_real"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Saldo Real (R$)</FormLabel>
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
                    Conte o dinheiro na gaveta e informe o valor exato
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Alerta de Divergência */}
            {temDivergencia && (
              <Alert variant={divergencia! > 0 ? 'default' : 'destructive'}>
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>
                  {divergencia! > 0 ? 'Sobra' : 'Falta'} no Caixa
                </AlertTitle>
                <AlertDescription>
                  Divergência de{' '}
                  <span className="font-bold">
                    {formatCurrency(Math.abs(divergencia!).toString())}
                  </span>
                  . Uma justificativa é obrigatória.
                </AlertDescription>
              </Alert>
            )}

            <FormField
              control={form.control}
              name="justificativa"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    Justificativa {temDivergencia && '*'}
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      {...field}
                      placeholder={
                        temDivergencia
                          ? 'Explique a divergência...'
                          : 'Observações opcionais...'
                      }
                      rows={3}
                    />
                  </FormControl>
                  {temDivergencia && (
                    <FormDescription>
                      Obrigatório quando há divergência
                    </FormDescription>
                  )}
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={fecharCaixa.isPending}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={fecharCaixa.isPending}>
                {fecharCaixa.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Fechar Caixa
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default ModalFecharCaixa;
