'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Tela de Entrada de Estoque
 *
 * 100% Design System compliant
 * Features:
 * - Seleção de fornecedor cadastrado
 * - Forma de pagamento (1x, 2x, 3x)
 * - Data do pedido
 * - Múltiplos produtos por entrada
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { addDays, format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    ArrowLeft,
    Calendar as CalendarIcon,
    Check,
    Loader2,
    PackagePlus,
    Plus,
    Trash2,
    X,
} from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useMemo, useState } from 'react';
import { useFieldArray, useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from '@/components/ui/card';
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from '@/components/ui/command';
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
import { Label } from '@/components/ui/label';
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '@/components/ui/popover';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Textarea } from '@/components/ui/textarea';
import { useFornecedores } from '@/hooks/use-fornecedores';
import { useCreateStockEntryMultiple, useStockItems } from '@/hooks/use-stock';
import { cn } from '@/lib/utils';
import { useBreadcrumbs } from '@/store/ui-store';

// =============================================================================
// SCHEMAS DE VALIDAÇÃO
// =============================================================================

const itemEntradaSchema = z.object({
  produto_id: z.string().min(1, 'Selecione um produto'),
  quantidade: z.coerce
    .number()
    .positive('Quantidade deve ser maior que zero')
    .min(0.01, 'Quantidade mínima: 0.01'),
  valor_unitario: z
    .string()
    .min(1, 'Valor unitário é obrigatório')
    .refine(
      (value) => {
        const parsed = Number(value.replace(',', '.'));
        return !Number.isNaN(parsed) && parsed >= 0;
      },
      { message: 'Valor inválido' }
    ),
});

const entrySchema = z.object({
  fornecedor_id: z.string().min(1, 'Selecione um fornecedor'),
  data_pedido: z.date({ required_error: 'Selecione a data do pedido' }),
  forma_pagamento: z.enum(['1x', '2x', '3x'], {
    required_error: 'Selecione a forma de pagamento',
  }),
  itens: z
    .array(itemEntradaSchema)
    .min(1, 'Adicione pelo menos um produto'),
  observacoes: z.string().max(500, 'Máximo 500 caracteres').optional(),
  gerar_financeiro: z.boolean().default(true),
});

type EntryFormData = z.infer<typeof entrySchema>;

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export default function EntradaEstoquePage() {
  const router = useRouter();
  const { setBreadcrumbs } = useBreadcrumbs();
  const createEntry = useCreateStockEntryMultiple();

  // Queries
  const { data: stockData, isLoading: isLoadingItems } = useStockItems({
    is_active: true,
    page_size: 100,
  });
  const {
    data: fornecedores,
    isLoading: isLoadingFornecedores,
  } = useFornecedores();

  // Estados locais
  const [openProductSelect, setOpenProductSelect] = useState<number | null>(
    null
  );
  const [openFornecedorSelect, setOpenFornecedorSelect] = useState(false);

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Estoque', href: '/estoque' },
      { label: 'Entrada de Estoque' },
    ]);
  }, [setBreadcrumbs]);

  const form = useForm<EntryFormData>({
    resolver: zodResolver(entrySchema),
    defaultValues: {
      fornecedor_id: '',
      data_pedido: new Date(),
      forma_pagamento: '1x',
      itens: [{ produto_id: '', quantidade: 1, valor_unitario: '' }],
      observacoes: '',
      gerar_financeiro: true,
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: 'itens',
  });

  const itens = useWatch({ control: form.control, name: 'itens' });
  const dataPedido = useWatch({ control: form.control, name: 'data_pedido' });
  const formaPagamento = useWatch({ control: form.control, name: 'forma_pagamento' });
  const fornecedorId = useWatch({ control: form.control, name: 'fornecedor_id' });

  // Cálculo do total
  const totalGeral = useMemo(() => {
    return itens.reduce((acc, item) => {
      const valor = Number(item.valor_unitario?.replace(',', '.') || 0);
      const qtd = Number(item.quantidade || 0);
      return acc + valor * qtd;
    }, 0);
  }, [itens]);

  // Calcular datas de vencimento baseado na forma de pagamento
  const datasVencimento = useMemo(() => {
    if (!dataPedido) return [];

    const parcelas = parseInt(formaPagamento);
    const datas: Date[] = [];

    for (let i = 0; i < parcelas; i++) {
      datas.push(addDays(dataPedido, 30 * (i + 1)));
    }

    return datas;
  }, [dataPedido, formaPagamento]);

  const formatCurrency = (value: number) =>
    new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
    }).format(value);

  const handleSubmit = async (data: EntryFormData) => {
    try {
      await createEntry.mutateAsync({
        fornecedor_id: data.fornecedor_id,
        data_entrada: format(data.data_pedido, 'yyyy-MM-dd'),
        itens: data.itens.map((item) => ({
          produto_id: item.produto_id,
          quantidade: Math.round(item.quantidade), // Backend espera int
          valor_unitario: item.valor_unitario.replace(',', '.'),
        })),
        observacoes: data.observacoes || '',
        gerar_financeiro: data.gerar_financeiro,
      });
      router.push('/estoque');
    } catch {
      // toast tratado no hook
    }
  };

  const fornecedorSelecionado = useMemo(() => {
    return fornecedores?.find((f) => f.id === fornecedorId);
  }, [fornecedorId, fornecedores]);

  return (
    <div className="container mx-auto space-y-6 py-6">
      {/* ========== HEADER ========== */}
      <div className="flex items-start gap-4">
        <Button
          asChild
          variant="outline"
          size="icon"
          className="shrink-0"
        >
          <Link href="/estoque">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <PackagePlus className="h-7 w-7 text-primary" />
            <h1 className="text-3xl font-bold tracking-tight">
              Registrar Entrada de Estoque
            </h1>
          </div>
          <p className="mt-2 text-base text-muted-foreground">
            Registre a compra de produtos, selecione o fornecedor e defina a forma de pagamento
          </p>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-3">
            {/* ========== COLUNA ESQUERDA: DADOS DO PEDIDO ========== */}
            <div className="lg:col-span-2 space-y-6">
              {/* Card: Informações do Pedido */}
              <Card>
                <CardHeader>
                  <CardTitle>Informações do Pedido</CardTitle>
                  <CardDescription>
                    Defina o fornecedor, data e forma de pagamento
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  {/* Fornecedor com Combobox */}
                  <FormField
                    control={form.control}
                    name="fornecedor_id"
                    render={({ field }) => (
                      <FormItem className="flex flex-col">
                        <FormLabel>Fornecedor *</FormLabel>
                        <Popover
                          open={openFornecedorSelect}
                          onOpenChange={setOpenFornecedorSelect}
                        >
                          <PopoverTrigger asChild>
                            <FormControl>
                              <Button
                                variant="outline"
                                role="combobox"
                                disabled={isLoadingFornecedores || createEntry.isPending}
                                className={cn(
                                  'w-full justify-between',
                                  !field.value && 'text-muted-foreground'
                                )}
                              >
                                {field.value
                                  ? fornecedores?.find((f) => f.id === field.value)
                                      ?.nome || 'Selecione...'
                                  : 'Selecione um fornecedor'}
                                <CalendarIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                              </Button>
                            </FormControl>
                          </PopoverTrigger>
                          <PopoverContent className="w-full p-0" align="start">
                            <Command>
                              <CommandInput placeholder="Buscar fornecedor..." />
                              <CommandList>
                                <CommandEmpty>
                                  Nenhum fornecedor encontrado.
                                </CommandEmpty>
                                <CommandGroup>
                                  <ScrollArea className="h-[200px]">
                                    {fornecedores?.map((fornecedor) => (
                                      <CommandItem
                                        key={fornecedor.id}
                                        value={fornecedor.nome}
                                        onSelect={() => {
                                          field.onChange(fornecedor.id);
                                          setOpenFornecedorSelect(false);
                                        }}
                                      >
                                        <Check
                                          className={cn(
                                            'mr-2 h-4 w-4',
                                            fornecedor.id === field.value
                                              ? 'opacity-100'
                                              : 'opacity-0'
                                          )}
                                        />
                                        <div className="flex flex-col">
                                          <span className="font-medium">
                                            {fornecedor.nome}
                                          </span>
                                          {fornecedor.email && (
                                            <span className="text-sm text-muted-foreground">
                                              {fornecedor.email}
                                            </span>
                                          )}
                                        </div>
                                      </CommandItem>
                                    ))}
                                  </ScrollArea>
                                </CommandGroup>
                              </CommandList>
                            </Command>
                          </PopoverContent>
                        </Popover>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="grid gap-4 sm:grid-cols-2">
                    {/* Data do Pedido */}
                    <FormField
                      control={form.control}
                      name="data_pedido"
                      render={({ field }) => (
                        <FormItem className="flex flex-col">
                          <FormLabel>Data do Pedido *</FormLabel>
                          <Popover>
                            <PopoverTrigger asChild>
                              <FormControl>
                                <Button
                                  variant="outline"
                                  disabled={createEntry.isPending}
                                  className={cn(
                                    'w-full pl-3 text-left font-normal',
                                    !field.value && 'text-muted-foreground'
                                  )}
                                >
                                  {field.value ? (
                                    format(field.value, 'PPP', { locale: ptBR })
                                  ) : (
                                    <span>Selecione a data</span>
                                  )}
                                  <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                </Button>
                              </FormControl>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0" align="start">
                              <Calendar
                                mode="single"
                                selected={field.value}
                                onSelect={field.onChange}
                                disabled={(date) =>
                                  date > new Date() || date < new Date('1900-01-01')
                                }
                                initialFocus
                                locale={ptBR}
                              />
                            </PopoverContent>
                          </Popover>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    {/* Forma de Pagamento */}
                    <FormField
                      control={form.control}
                      name="forma_pagamento"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Forma de Pagamento *</FormLabel>
                          <Select
                            onValueChange={field.onChange}
                            value={field.value}
                            disabled={createEntry.isPending}
                          >
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Selecione" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="1x">À vista (1x)</SelectItem>
                              <SelectItem value="2x">Parcelado em 2x</SelectItem>
                              <SelectItem value="3x">Parcelado em 3x</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormDescription>
                            Define o parcelamento da conta a pagar
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Card: Produtos */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle>Produtos</CardTitle>
                      <CardDescription>
                        Adicione os produtos desta entrada
                      </CardDescription>
                    </div>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        append({ produto_id: '', quantidade: 1, valor_unitario: '' })
                      }
                      disabled={createEntry.isPending}
                    >
                      <Plus className="mr-2 h-4 w-4" />
                      Adicionar Produto
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  {fields.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-12 text-center">
                      <PackagePlus className="h-12 w-12 text-muted-foreground/50" />
                      <p className="mt-4 text-sm text-muted-foreground">
                        Nenhum produto adicionado ainda
                      </p>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="mt-4"
                        onClick={() =>
                          append({
                            produto_id: '',
                            quantidade: 1,
                            valor_unitario: '',
                          })
                        }
                      >
                        <Plus className="mr-2 h-4 w-4" />
                        Adicionar Primeiro Produto
                      </Button>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      {/* Desktop: Tabela */}
                      <div className="hidden md:block">
                        <Table>
                          <TableHeader>
                            <TableRow>
                              <TableHead className="w-[40%]">Produto</TableHead>
                              <TableHead className="w-[20%]">Quantidade</TableHead>
                              <TableHead className="w-[25%]">Valor Unit.</TableHead>
                              <TableHead className="w-[15%] text-right">
                                Subtotal
                              </TableHead>
                              <TableHead className="w-[50px]"></TableHead>
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            {fields.map((field, index) => {
                              const item = itens?.[index] || { valor_unitario: '', quantidade: 0, produto_id: '' };
                              const subtotal =
                                Number(
                                  item.valor_unitario?.replace(',', '.') || 0
                                ) * Number(item.quantidade || 0);

                              return (
                                <TableRow key={field.id}>
                                  <TableCell>
                                    <FormField
                                      control={form.control}
                                      name={`itens.${index}.produto_id`}
                                      render={({ field }) => (
                                        <FormItem>
                                          <Popover
                                            open={openProductSelect === index}
                                            onOpenChange={(open) =>
                                              setOpenProductSelect(open ? index : null)
                                            }
                                          >
                                            <PopoverTrigger asChild>
                                              <FormControl>
                                                <Button
                                                  variant="outline"
                                                  role="combobox"
                                                  disabled={
                                                    isLoadingItems ||
                                                    createEntry.isPending
                                                  }
                                                  className={cn(
                                                    'w-full justify-between',
                                                    !field.value &&
                                                      'text-muted-foreground'
                                                  )}
                                                >
                                                  {field.value
                                                    ? stockData?.data.find(
                                                        (p) => p.id === field.value
                                                      )?.nome || 'Selecione...'
                                                    : 'Selecione um produto'}
                                                  <CalendarIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                                                </Button>
                                              </FormControl>
                                            </PopoverTrigger>
                                            <PopoverContent
                                              className="w-[400px] p-0"
                                              align="start"
                                            >
                                              <Command>
                                                <CommandInput placeholder="Buscar produto..." />
                                                <CommandList>
                                                  <CommandEmpty>
                                                    Nenhum produto encontrado.
                                                  </CommandEmpty>
                                                  <CommandGroup>
                                                    <ScrollArea className="h-[200px]">
                                                      {stockData?.data.map(
                                                        (produto) => (
                                                          <CommandItem
                                                            key={produto.id}
                                                            value={produto.nome}
                                                            onSelect={() => {
                                                              field.onChange(
                                                                produto.id
                                                              );
                                                              setOpenProductSelect(
                                                                null
                                                              );
                                                            }}
                                                          >
                                                            <Check
                                                              className={cn(
                                                                'mr-2 h-4 w-4',
                                                                produto.id ===
                                                                  field.value
                                                                  ? 'opacity-100'
                                                                  : 'opacity-0'
                                                              )}
                                                            />
                                                            <div className="flex flex-col">
                                                              <span className="font-medium">
                                                                {produto.nome}
                                                              </span>
                                                              <span className="text-sm text-muted-foreground">
                                                                Estoque:{' '}
                                                                {
                                                                  produto.quantidade_atual
                                                                }{' '}
                                                                {
                                                                  produto.unidade_medida
                                                                }
                                                              </span>
                                                            </div>
                                                          </CommandItem>
                                                        )
                                                      )}
                                                    </ScrollArea>
                                                  </CommandGroup>
                                                </CommandList>
                                              </Command>
                                            </PopoverContent>
                                          </Popover>
                                          <FormMessage />
                                        </FormItem>
                                      )}
                                    />
                                  </TableCell>
                                  <TableCell>
                                    <FormField
                                      control={form.control}
                                      name={`itens.${index}.quantidade`}
                                      render={({ field }) => (
                                        <FormItem>
                                          <FormControl>
                                            <Input
                                              type="number"
                                              min={0.01}
                                              step={0.01}
                                              placeholder="1"
                                              disabled={createEntry.isPending}
                                              {...field}
                                              onChange={(e) =>
                                                field.onChange(
                                                  Number(e.target.value)
                                                )
                                              }
                                            />
                                          </FormControl>
                                          <FormMessage />
                                        </FormItem>
                                      )}
                                    />
                                  </TableCell>
                                  <TableCell>
                                    <FormField
                                      control={form.control}
                                      name={`itens.${index}.valor_unitario`}
                                      render={({ field }) => (
                                        <FormItem>
                                          <FormControl>
                                            <div className="relative">
                                              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted-foreground">
                                                R$
                                              </span>
                                              <Input
                                                placeholder="0,00"
                                                className="pl-10"
                                                disabled={createEntry.isPending}
                                                {...field}
                                              />
                                            </div>
                                          </FormControl>
                                          <FormMessage />
                                        </FormItem>
                                      )}
                                    />
                                  </TableCell>
                                  <TableCell className="text-right font-medium">
                                    {formatCurrency(subtotal)}
                                  </TableCell>
                                  <TableCell>
                                    <Button
                                      type="button"
                                      variant="ghost"
                                      size="icon"
                                      onClick={() => remove(index)}
                                      disabled={
                                        fields.length === 1 || createEntry.isPending
                                      }
                                    >
                                      <Trash2 className="h-4 w-4 text-destructive" />
                                    </Button>
                                  </TableCell>
                                </TableRow>
                              );
                            })}
                          </TableBody>
                        </Table>
                      </div>

                      {/* Mobile: Cards */}
                      <div className="md:hidden space-y-4">
                        {fields.map((field, index) => {
                          const item = itens?.[index] || { valor_unitario: '', quantidade: 0, produto_id: '' };
                          const subtotal =
                            Number(item.valor_unitario?.replace(',', '.') || 0) *
                            Number(item.quantidade || 0);

                          return (
                            <Card key={field.id}>
                              <CardHeader className="pb-3">
                                <div className="flex items-start justify-between">
                                  <CardTitle className="text-base">
                                    Item {index + 1}
                                  </CardTitle>
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => remove(index)}
                                    disabled={
                                      fields.length === 1 || createEntry.isPending
                                    }
                                  >
                                    <X className="h-4 w-4" />
                                  </Button>
                                </div>
                              </CardHeader>
                              <CardContent className="space-y-4">
                                <FormField
                                  control={form.control}
                                  name={`itens.${index}.produto_id`}
                                  render={({ field }) => (
                                    <FormItem>
                                      <FormLabel>Produto</FormLabel>
                                      <Select
                                        onValueChange={field.onChange}
                                        value={field.value}
                                        disabled={
                                          isLoadingItems || createEntry.isPending
                                        }
                                      >
                                        <FormControl>
                                          <SelectTrigger>
                                            <SelectValue placeholder="Selecione" />
                                          </SelectTrigger>
                                        </FormControl>
                                        <SelectContent>
                                          {stockData?.data.map((produto) => (
                                            <SelectItem
                                              key={produto.id}
                                              value={produto.id}
                                            >
                                              {produto.nome}
                                            </SelectItem>
                                          ))}
                                        </SelectContent>
                                      </Select>
                                      <FormMessage />
                                    </FormItem>
                                  )}
                                />
                                <div className="grid grid-cols-2 gap-3">
                                  <FormField
                                    control={form.control}
                                    name={`itens.${index}.quantidade`}
                                    render={({ field }) => (
                                      <FormItem>
                                        <FormLabel>Qtd.</FormLabel>
                                        <FormControl>
                                          <Input
                                            type="number"
                                            min={0.01}
                                            step={0.01}
                                            disabled={createEntry.isPending}
                                            {...field}
                                            onChange={(e) =>
                                              field.onChange(Number(e.target.value))
                                            }
                                          />
                                        </FormControl>
                                        <FormMessage />
                                      </FormItem>
                                    )}
                                  />
                                  <FormField
                                    control={form.control}
                                    name={`itens.${index}.valor_unitario`}
                                    render={({ field }) => (
                                      <FormItem>
                                        <FormLabel>Valor Unit.</FormLabel>
                                        <FormControl>
                                          <div className="relative">
                                            <span className="absolute left-2 top-1/2 -translate-y-1/2 text-xs text-muted-foreground">
                                              R$
                                            </span>
                                            <Input
                                              placeholder="0,00"
                                              className="pl-8"
                                              disabled={createEntry.isPending}
                                              {...field}
                                            />
                                          </div>
                                        </FormControl>
                                        <FormMessage />
                                      </FormItem>
                                    )}
                                  />
                                </div>
                                <div className="pt-2 border-t">
                                  <div className="flex justify-between">
                                    <span className="text-sm text-muted-foreground">
                                      Subtotal:
                                    </span>
                                    <span className="font-semibold">
                                      {formatCurrency(subtotal)}
                                    </span>
                                  </div>
                                </div>
                              </CardContent>
                            </Card>
                          );
                        })}
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Card: Observações */}
              <Card>
                <CardHeader>
                  <CardTitle>Informações Adicionais</CardTitle>
                </CardHeader>
                <CardContent>
                  <FormField
                    control={form.control}
                    name="observacoes"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Observações</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder="Detalhes sobre o pedido, notas fiscais, etc."
                            rows={4}
                            disabled={createEntry.isPending}
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          Informações complementares sobre esta entrada (máx. 500
                          caracteres)
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>
            </div>

            {/* ========== COLUNA DIREITA: RESUMO ========== */}
            <div className="space-y-6">
              {/* Card: Resumo Financeiro */}
              <Card className="sticky top-6">
                <CardHeader>
                  <CardTitle>Resumo Financeiro</CardTitle>
                  <CardDescription>Valores e parcelamento</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  {/* Total Geral */}
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">
                        Total de Produtos
                      </span>
                      <span className="font-medium">{fields.length}</span>
                    </div>
                    <Separator />
                    <div className="flex justify-between items-center pt-2">
                      <span className="text-base font-semibold">Total Geral</span>
                      <span className="text-2xl font-bold text-primary">
                        {formatCurrency(totalGeral)}
                      </span>
                    </div>
                  </div>

                  {/* Parcelamento */}
                  {formaPagamento && datasVencimento.length > 0 && (
                    <>
                      <Separator />
                      <div className="space-y-3">
                        <Label className="text-sm font-semibold">
                          Parcelamento ({formaPagamento})
                        </Label>
                        <div className="space-y-2">
                          {datasVencimento.map((data, index) => (
                            <div
                              key={index}
                              className="flex justify-between items-center rounded-lg bg-muted/50 p-3"
                            >
                              <div className="flex items-center gap-2">
                                <Badge variant="outline" className="font-mono">
                                  {index + 1}/{datasVencimento.length}
                                </Badge>
                                <span className="text-sm text-muted-foreground">
                                  {format(data, 'dd/MM/yyyy')}
                                </span>
                              </div>
                              <span className="font-semibold">
                                {formatCurrency(totalGeral / datasVencimento.length)}
                              </span>
                            </div>
                          ))}
                        </div>
                      </div>
                    </>
                  )}

                  {/* Fornecedor Selecionado */}
                  {fornecedorSelecionado && (
                    <>
                      <Separator />
                      <div className="space-y-3">
                        <Label className="text-sm font-semibold">Fornecedor</Label>
                        <div className="rounded-lg bg-muted/50 p-3 space-y-1">
                          <p className="font-medium">{fornecedorSelecionado.nome}</p>
                          {fornecedorSelecionado.telefone && (
                            <p className="text-sm text-muted-foreground">
                              {fornecedorSelecionado.telefone}
                            </p>
                          )}
                          {fornecedorSelecionado.email && (
                            <p className="text-sm text-muted-foreground">
                              {fornecedorSelecionado.email}
                            </p>
                          )}
                        </div>
                      </div>
                    </>
                  )}

                  {/* Gerar Financeiro */}
                  <Separator />
                  <FormField
                    control={form.control}
                    name="gerar_financeiro"
                    render={({ field }) => (
                      <FormItem className="flex items-center justify-between rounded-lg border p-3">
                        <div className="space-y-0.5">
                          <FormLabel className="text-sm font-semibold">
                            Gerar Financeiro
                          </FormLabel>
                          <FormDescription className="text-xs">
                            Cria contas a pagar automaticamente
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            disabled={createEntry.isPending}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                </CardContent>
              </Card>

              {/* Ações */}
              <Card>
                <CardContent className="pt-6">
                  <div className="flex flex-col gap-3">
                    <Button
                      type="submit"
                      size="lg"
                      disabled={createEntry.isPending || fields.length === 0}
                      className="w-full"
                    >
                      {createEntry.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          Registrando...
                        </>
                      ) : (
                        <>
                          <Check className="mr-2 h-4 w-4" />
                          Registrar Entrada
                        </>
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      size="lg"
                      asChild
                      className="w-full"
                    >
                      <Link href="/estoque">Cancelar</Link>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </form>
      </Form>
    </div>
  );
}
