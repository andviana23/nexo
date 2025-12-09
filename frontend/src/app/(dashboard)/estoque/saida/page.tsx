'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Tela de Saída de Estoque
 *
 * Permite registrar baixas de estoque (venda, uso interno ou perdas).
 * Depende do contrato de estoque alinhado em P4.1.
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeft, Loader2, PackageMinus } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useMemo } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { Textarea } from '@/components/ui/textarea';
import { useCreateStockExit, useStockItems } from '@/hooks/use-stock';
import { useBreadcrumbs } from '@/store/ui-store';
import { StockItem } from '@/types/stock';

// =============================================================================
// SCHEMA
// =============================================================================

const exitSchema = z.object({
  stock_item_id: z.string().min(1, 'Selecione um produto'),
  quantity: z.coerce.number().positive('Quantidade deve ser maior que zero'),
  reason: z.enum(['VENDA', 'USO_INTERNO', 'PERDA', 'DEVOLUCAO'], {
    required_error: 'Informe o motivo da saída',
  }),
  reference: z.string().max(100, 'Máximo 100 caracteres').optional(),
  notes: z.string().max(500, 'Máximo 500 caracteres').optional(),
});

type ExitFormData = z.infer<typeof exitSchema>;

// =============================================================================
// COMPONENT
// =============================================================================

export default function SaidaEstoquePage() {
  const router = useRouter();
  const { setBreadcrumbs } = useBreadcrumbs();
  const createExit = useCreateStockExit();
  const { data: stockData, isLoading: isLoadingItems } = useStockItems({
    is_active: true,
    page_size: 100,
  });

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Estoque', href: '/estoque' },
      { label: 'Saída de Estoque' },
    ]);
  }, [setBreadcrumbs]);

  const form = useForm<ExitFormData>({
    resolver: zodResolver(exitSchema),
    defaultValues: {
      stock_item_id: '',
      quantity: 1,
      reason: 'VENDA',
      reference: '',
      notes: '',
    },
  });

  const selectedItemId = useWatch({ control: form.control, name: 'stock_item_id' });
  const selectedItem: StockItem | undefined = useMemo(
    () => stockData?.data.find((item) => item.id === selectedItemId),
    [selectedItemId, stockData?.data]
  );
  const quantity = useWatch({ control: form.control, name: 'quantity' });
  const projectedQuantity =
    selectedItem && typeof quantity === 'number'
      ? Math.max(parseFloat(selectedItem.quantidade_atual) - quantity, 0)
      : undefined;

  const handleSubmit = async (data: ExitFormData) => {
    if (selectedItem && data.quantity > parseFloat(selectedItem.quantidade_atual)) {
      form.setError('quantity', {
        type: 'manual',
        message: 'Quantidade maior que o estoque disponível',
      });
      return;
    }

    try {
      await createExit.mutateAsync({
        stock_item_id: data.stock_item_id,
        quantity: data.quantity,
        reason: data.reason,
        reference: data.reference || undefined,
        notes: data.notes || undefined,
      });
      router.push('/estoque');
    } catch {
      // toast tratado no hook useCreateStockExit
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button asChild variant="ghost" size="icon">
          <Link href="/estoque">
            <ArrowLeft className="h-5 w-5" />
          </Link>
        </Button>
        <div>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <PackageMinus className="h-7 w-7" />
            Registrar Saída
          </h1>
          <p className="text-muted-foreground">
            Registre consumos, vendas ou perdas com controle de saldo.
          </p>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-[2fr,1fr]">
        {/* Formulário */}
        <Card>
          <CardHeader>
            <CardTitle>Dados da Saída</CardTitle>
            <CardDescription>
              Informe o produto, a quantidade, o motivo e uma referência opcional.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(handleSubmit)}
                className="space-y-6"
              >
                {/* Produto */}
                <FormField
                  control={form.control}
                  name="stock_item_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Produto *</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        disabled={isLoadingItems || createExit.isPending}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione o produto" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {isLoadingItems && (
                            <SelectItem value="loading" disabled>
                              Carregando produtos...
                            </SelectItem>
                          )}
                          {!isLoadingItems &&
                            stockData?.data.map((item) => (
                              <SelectItem key={item.id} value={item.id}>
                                {item.nome} • {item.quantidade_atual} {item.unidade_medida}
                              </SelectItem>
                            ))}
                          {!isLoadingItems && (stockData?.data.length || 0) === 0 && (
                            <SelectItem value="empty" disabled>
                              Nenhum produto cadastrado
                            </SelectItem>
                          )}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Quantidade */}
                <FormField
                  control={form.control}
                  name="quantity"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Quantidade *</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          min={0.01}
                          step={0.01}
                          placeholder="Ex: 2"
                          value={field.value}
                          onChange={(e) => field.onChange(Number(e.target.value))}
                          disabled={createExit.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Motivo */}
                <FormField
                  control={form.control}
                  name="reason"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Motivo *</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        disabled={createExit.isPending}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione o motivo" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="VENDA">Venda</SelectItem>
                          <SelectItem value="USO_INTERNO">Uso interno</SelectItem>
                          <SelectItem value="PERDA">Perda/avaria</SelectItem>
                          <SelectItem value="DEVOLUCAO">Devolução</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Referência */}
                <FormField
                  control={form.control}
                  name="reference"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Referência</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Ex: Venda #123, pedido ou nota fiscal"
                          {...field}
                          disabled={createExit.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Observações */}
                <FormField
                  control={form.control}
                  name="notes"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Observações</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Detalhes adicionais sobre a saída"
                          rows={3}
                          {...field}
                          disabled={createExit.isPending}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="flex justify-end gap-3">
                  <Button type="button" variant="outline" asChild>
                    <Link href="/estoque">Cancelar</Link>
                  </Button>
                  <Button type="submit" disabled={createExit.isPending}>
                    {createExit.isPending ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Registrando...
                      </>
                    ) : (
                      'Registrar Saída'
                    )}
                  </Button>
                </div>
              </form>
            </Form>
          </CardContent>
        </Card>

        {/* Detalhes do Produto Selecionado */}
        <Card>
          <CardHeader>
            <CardTitle>Saldo e alerta</CardTitle>
            <CardDescription>
              Confirme o estoque disponível antes da baixa.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {isLoadingItems ? (
              <div className="space-y-2">
                <Skeleton className="h-6 w-full" />
                <Skeleton className="h-20 w-full" />
              </div>
            ) : selectedItem ? (
              <div className="space-y-4">
                <div className="flex items-center gap-3">
                  <h3 className="text-lg font-semibold">{selectedItem.nome}</h3>
                  {selectedItem.categoria_produto && (
                    <Badge variant="outline">{selectedItem.categoria_produto.nome}</Badge>
                  )}
                </div>
                {selectedItem.descricao && (
                  <p className="text-sm text-muted-foreground">
                    {selectedItem.descricao}
                  </p>
                )}
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <p className="text-xs text-muted-foreground">Saldo Atual</p>
                    <p className="text-2xl font-bold">
                      {selectedItem.quantidade_atual} {selectedItem.unidade_medida}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Após Saída</p>
                    <p
                      className={`text-2xl font-bold ${
                        projectedQuantity !== undefined && projectedQuantity <= parseFloat(selectedItem.quantidade_minima)
                          ? 'text-orange-500'
                          : ''
                      }`}
                    >
                      {projectedQuantity !== undefined
                        ? `${projectedQuantity} ${selectedItem.unidade_medida}`
                        : '-'}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Estoque Mínimo</p>
                    <p className="text-lg font-semibold">
                      {selectedItem.quantidade_minima} {selectedItem.unidade_medida}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Código de Barras</p>
                    <p className="font-mono text-sm">{selectedItem.codigo_barras || '-'}</p>
                  </div>
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">
                Selecione um produto para ver o saldo disponível.
              </p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
