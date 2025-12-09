/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Nova Despesa Fixa
 *
 * Formulário para criação de nova despesa fixa recorrente.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeft, Loader2, Save } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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

import { useCreateFixedExpense } from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';

// =============================================================================
// SCHEMA DE VALIDAÇÃO
// =============================================================================

const despesaFixaSchema = z.object({
  descricao: z
    .string()
    .min(3, 'Descrição deve ter no mínimo 3 caracteres')
    .max(255, 'Descrição deve ter no máximo 255 caracteres'),
  fornecedor: z.string().optional(),
  valor: z
    .string()
    .min(1, 'Valor é obrigatório')
    .refine(
      (val) => {
        const num = parseFloat(val.replace(',', '.'));
        return !isNaN(num) && num > 0;
      },
      { message: 'Valor deve ser maior que zero' }
    ),
  dia_vencimento: z
    .number()
    .min(1, 'Dia deve ser entre 1 e 31')
    .max(31, 'Dia deve ser entre 1 e 31'),
  categoria_id: z.string().optional(),
  observacoes: z.string().optional(),
});

type DespesaFixaFormData = z.infer<typeof despesaFixaSchema>;

// =============================================================================
// PÁGINA
// =============================================================================

export default function NovaDespesaFixaPage() {
  const router = useRouter();
  const createMutation = useCreateFixedExpense();

  // Breadcrumbs
  const { setBreadcrumbs } = useBreadcrumbs();
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'Despesas Fixas', href: '/financeiro/despesas-fixas' },
      { label: 'Nova Despesa' },
    ]);
  }, [setBreadcrumbs]);

  const form = useForm<DespesaFixaFormData>({
    resolver: zodResolver(despesaFixaSchema),
    defaultValues: {
      descricao: '',
      fornecedor: '',
      valor: '',
      dia_vencimento: 10,
      categoria_id: '',
      observacoes: '',
    },
  });

  const onSubmit = (data: DespesaFixaFormData) => {
    // Converte valor para formato string com ponto decimal
    const valorFormatado = data.valor.replace(',', '.');

    createMutation.mutate(
      {
        descricao: data.descricao,
        fornecedor: data.fornecedor || undefined,
        valor: valorFormatado,
        dia_vencimento: data.dia_vencimento,
        categoria_id: data.categoria_id || undefined,
        observacoes: data.observacoes || undefined,
      },
      {
        onSuccess: () => {
          router.push('/financeiro/despesas-fixas');
        },
      }
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link href="/financeiro/despesas-fixas">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-5 w-5" />
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Nova Despesa Fixa</h1>
          <p className="text-muted-foreground">
            Cadastre uma despesa recorrente mensal
          </p>
        </div>
      </div>

      {/* Formulário */}
      <Card>
        <CardHeader>
          <CardTitle>Dados da Despesa</CardTitle>
          <CardDescription>
            Preencha os dados da despesa fixa. Ela será gerada automaticamente todo mês.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              {/* Descrição */}
              <FormField
                control={form.control}
                name="descricao"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Descrição *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Ex: Aluguel, Energia, Internet..."
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Nome que identificará esta despesa
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Fornecedor */}
              <FormField
                control={form.control}
                name="fornecedor"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Fornecedor</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Ex: Imobiliária XYZ, ENEL, Vivo..."
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Empresa ou pessoa que receberá o pagamento
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Valor e Dia de Vencimento */}
              <div className="grid gap-4 md:grid-cols-2">
                <FormField
                  control={form.control}
                  name="valor"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Valor Mensal *</FormLabel>
                      <FormControl>
                        <div className="relative">
                          <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                            R$
                          </span>
                          <Input
                            placeholder="0,00"
                            className="pl-10"
                            {...field}
                          />
                        </div>
                      </FormControl>
                      <FormDescription>
                        Valor fixo mensal da despesa
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="dia_vencimento"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Dia de Vencimento *</FormLabel>
                      <Select
                        value={field.value?.toString()}
                        onValueChange={(v) => field.onChange(parseInt(v, 10))}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione o dia" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {Array.from({ length: 31 }, (_, i) => i + 1).map((dia) => (
                            <SelectItem key={dia} value={dia.toString()}>
                              Dia {dia}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormDescription>
                        Dia do mês para vencimento
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              {/* Observações */}
              <FormField
                control={form.control}
                name="observacoes"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Observações</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Informações adicionais sobre esta despesa..."
                        className="resize-none"
                        rows={3}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Botões */}
              <div className="flex items-center gap-4 pt-4">
                <Button
                  type="submit"
                  disabled={createMutation.isPending}
                >
                  {createMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  <Save className="mr-2 h-4 w-4" />
                  Salvar Despesa
                </Button>
                <Link href="/financeiro/despesas-fixas">
                  <Button variant="outline" type="button">
                    Cancelar
                  </Button>
                </Link>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
