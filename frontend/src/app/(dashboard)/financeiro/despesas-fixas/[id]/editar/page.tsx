/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Edição de Despesa Fixa
 *
 * Formulário para edição de despesa fixa existente.
 * Conforme FLUXO_FINANCEIRO.md
 */

'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { AlertCircle, ArrowLeft, Loader2, Save } from 'lucide-react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { Alert, AlertDescription } from '@/components/ui/alert';
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
import { Skeleton } from '@/components/ui/skeleton';
import { Textarea } from '@/components/ui/textarea';

import { useFixedExpense, useUpdateFixedExpense } from '@/hooks/use-financial';
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
// COMPONENTE: Loading Skeleton
// =============================================================================

function LoadingSkeleton() {
  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-6 w-48" />
        <Skeleton className="h-4 w-64 mt-2" />
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-10 w-full" />
        </div>
        <div className="space-y-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-10 w-full" />
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-10 w-full" />
          </div>
          <div className="space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-10 w-full" />
          </div>
        </div>
        <div className="space-y-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-24 w-full" />
        </div>
        <div className="flex gap-4 pt-4">
          <Skeleton className="h-10 w-32" />
          <Skeleton className="h-10 w-24" />
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// PÁGINA
// =============================================================================

export default function EditarDespesaFixaPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const { data: despesa, isLoading, error } = useFixedExpense(id);
  const updateMutation = useUpdateFixedExpense();

  // Breadcrumbs
  const { setBreadcrumbs } = useBreadcrumbs();
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Financeiro', href: '/financeiro' },
      { label: 'Despesas Fixas', href: '/financeiro/despesas-fixas' },
      { label: 'Editar Despesa' },
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

  // Preenche o formulário quando os dados chegam
  useEffect(() => {
    if (despesa) {
      form.reset({
        descricao: despesa.descricao,
        fornecedor: despesa.fornecedor || '',
        valor: despesa.valor,
        dia_vencimento: despesa.dia_vencimento,
        categoria_id: despesa.categoria_id || '',
        observacoes: despesa.observacoes || '',
      });
    }
  }, [despesa, form]);

  const onSubmit = (data: DespesaFixaFormData) => {
    // Converte valor para formato string com ponto decimal
    const valorFormatado = data.valor.replace(',', '.');

    updateMutation.mutate(
      {
        id,
        data: {
          descricao: data.descricao,
          fornecedor: data.fornecedor || undefined,
          valor: valorFormatado,
          dia_vencimento: data.dia_vencimento,
          categoria_id: data.categoria_id || undefined,
          observacoes: data.observacoes || undefined,
        },
      },
      {
        onSuccess: () => {
          router.push('/financeiro/despesas-fixas');
        },
      }
    );
  };

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/financeiro/despesas-fixas">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-5 w-5" />
            </Button>
          </Link>
          <h1 className="text-3xl font-bold tracking-tight">Editar Despesa Fixa</h1>
        </div>
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            Erro ao carregar despesa fixa. Verifique se o ID é válido.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

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
          <h1 className="text-3xl font-bold tracking-tight">Editar Despesa Fixa</h1>
          <p className="text-muted-foreground">
            Atualize os dados da despesa recorrente
          </p>
        </div>
      </div>

      {/* Formulário */}
      {isLoading ? (
        <LoadingSkeleton />
      ) : (
        <Card>
          <CardHeader>
            <CardTitle>Dados da Despesa</CardTitle>
            <CardDescription>
              Altere os dados conforme necessário. As contas futuras serão geradas com os novos valores.
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
                    disabled={updateMutation.isPending}
                  >
                    {updateMutation.isPending && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    <Save className="mr-2 h-4 w-4" />
                    Salvar Alterações
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
      )}
    </div>
  );
}
