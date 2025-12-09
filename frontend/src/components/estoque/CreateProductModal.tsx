/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal de Cadastro de Produto
 *
 * Permite criar um novo produto no estoque.
 * Segue DESIGN_SYSTEM.md e padrões do módulo de Estoque.
 *
 * @version 3.0.0
 * REFATORAÇÃO: Novos campos - codigo_barras, estoque_maximo, valor_venda_profissional,
 * valor_entrada, fornecedor_id. Removido: SKU, lead_time_dias, controla_validade, centro_custo.
 */

'use client';

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
import { useCategoriasProdutos } from '@/hooks/use-categorias-produtos';
import { useFornecedores } from '@/hooks/use-fornecedores';
import { useCreateProduct } from '@/hooks/use-stock';
import { getFornecedorNome } from '@/types/fornecedor';
import { zodResolver } from '@hookform/resolvers/zod';
import { AlertCircle, Loader2, Package } from 'lucide-react';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { z } from 'zod';

// =============================================================================
// TIPOS E CONSTANTES
// =============================================================================

/**
 * Unidades de medida conforme constraint do banco de dados:
 * CHECK (unidade_medida IN ('UNIDADE', 'LITRO', 'MILILITRO', 'GRAMA', 'QUILOGRAMA'))
 */
const UNIDADES = [
  { value: 'UNIDADE', label: 'Unidade (UN)' },
  { value: 'QUILOGRAMA', label: 'Quilograma (KG)' },
  { value: 'GRAMA', label: 'Grama (G)' },
  { value: 'MILILITRO', label: 'Mililitro (ML)' },
  { value: 'LITRO', label: 'Litro (L)' },
] as const;

// =============================================================================
// SCHEMA DE VALIDAÇÃO
// =============================================================================

const productSchema = z.object({
  nome: z
    .string()
    .min(2, 'Nome deve ter no mínimo 2 caracteres')
    .max(100, 'Nome deve ter no máximo 100 caracteres'),
  descricao: z
    .string()
    .max(500, 'Descrição deve ter no máximo 500 caracteres')
    .optional(),
  codigo_barras: z
    .string()
    .max(50, 'Código de barras deve ter no máximo 50 caracteres')
    .optional(),
  categoria_produto_id: z
    .string({ required_error: 'Selecione uma categoria' })
    .uuid('Categoria inválida'),
  unidade_medida: z.enum(['UNIDADE', 'QUILOGRAMA', 'GRAMA', 'MILILITRO', 'LITRO'], {
    required_error: 'Selecione uma unidade',
  }),
  valor_unitario: z
    .string()
    .min(1, 'Valor unitário é obrigatório')
    .refine(
      (val) => {
        const cleaned = val.replace(/[^\d,.-]/g, '').replace(',', '.');
        const num = parseFloat(cleaned);
        return !isNaN(num) && num >= 0;
      },
      { message: 'Valor inválido. Use formato: 10,00 ou 10.00' }
    ),
  quantidade_minima: z
    .string()
    .min(1, 'Estoque mínimo é obrigatório')
    .refine(
      (val) => {
        const num = parseFloat(val);
        return !isNaN(num) && num >= 0;
      },
      { message: 'Deve ser um número não negativo' }
    ),
  quantidade_maxima: z
    .string()
    .optional()
    .refine(
      (val) => {
        if (!val || val === '') return true;
        const num = parseFloat(val);
        return !isNaN(num) && num >= 0;
      },
      { message: 'Deve ser um número não negativo' }
    ),
  valor_venda_profissional: z
    .string()
    .optional()
    .refine(
      (val) => {
        if (!val || val === '') return true;
        const cleaned = val.replace(/[^\d,.-]/g, '').replace(',', '.');
        const num = parseFloat(cleaned);
        return !isNaN(num) && num >= 0;
      },
      { message: 'Valor inválido' }
    ),
  valor_entrada: z
    .string()
    .optional()
    .refine(
      (val) => {
        if (!val || val === '') return true;
        const cleaned = val.replace(/[^\d,.-]/g, '').replace(',', '.');
        const num = parseFloat(cleaned);
        return !isNaN(num) && num >= 0;
      },
      { message: 'Valor inválido' }
    ),
  fornecedor_id: z
    .string()
    .uuid('Fornecedor inválido')
    .optional()
    .or(z.literal('')),
});

type ProductFormData = z.infer<typeof productSchema>;

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Formata valor monetário para exibição (BR)
 */
function formatCurrency(value: string): string {
  const cleaned = value.replace(/\D/g, '');
  if (!cleaned) return '';
  const num = parseFloat(cleaned) / 100;
  return num.toLocaleString('pt-BR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

/**
 * Converte valor formatado para backend (string decimal)
 */
function currencyToDecimal(value: string | undefined): string | undefined {
  if (!value || value === '') return undefined;
  const cleaned = value.replace(/[^\d,.-]/g, '').replace(',', '.');
  const num = parseFloat(cleaned);
  if (isNaN(num)) return undefined;
  return num.toFixed(2);
}

// =============================================================================
// PROPS
// =============================================================================

interface CreateProductModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function CreateProductModal({
  open,
  onOpenChange,
  onSuccess,
}: CreateProductModalProps) {
  const createProduct = useCreateProduct();
  const {
    data: categorias,
    isLoading: categoriasLoading,
    error: categoriasError,
  } = useCategoriasProdutos();
  const {
    data: fornecedores,
    isLoading: fornecedoresLoading,
    error: fornecedoresError,
  } = useFornecedores();

  // Filtrar apenas categorias e fornecedores ativos
  const categoriasAtivas = categorias?.filter((c) => c.ativa) ?? [];
  const fornecedoresAtivos = fornecedores?.filter((f) => f.ativo) ?? [];

  const form = useForm<ProductFormData>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      nome: '',
      descricao: '',
      codigo_barras: '',
      categoria_produto_id: '',
      unidade_medida: 'UNIDADE',
      valor_unitario: '',
      quantidade_minima: '0',
      quantidade_maxima: '',
      valor_venda_profissional: '',
      valor_entrada: '',
      fornecedor_id: '',
    },
    mode: 'onBlur',
  });

  // Reset form ao abrir/fechar modal
  useEffect(() => {
    if (!open) {
      form.reset();
    }
  }, [open, form]);

  const onSubmit = async (data: ProductFormData) => {
    try {
      await createProduct.mutateAsync({
        nome: data.nome.trim(),
        descricao: data.descricao?.trim() || '',
        codigo_barras: data.codigo_barras?.trim() || undefined,
        categoria_produto_id: data.categoria_produto_id,
        unidade_medida: data.unidade_medida as string,
        valor_unitario: currencyToDecimal(data.valor_unitario) ?? '0.00',
        quantidade_minima: parseFloat(data.quantidade_minima) || 0,
        quantidade_maxima: data.quantidade_maxima
          ? data.quantidade_maxima
          : undefined,
        valor_venda_profissional: currencyToDecimal(
          data.valor_venda_profissional
        ),
        valor_entrada: currencyToDecimal(data.valor_entrada),
        fornecedor_id:
          data.fornecedor_id && data.fornecedor_id !== 'none'
            ? data.fornecedor_id
            : undefined,
      });

      toast.success('Produto cadastrado', {
        description: `${data.nome} foi adicionado ao estoque com sucesso.`,
      });

      onOpenChange(false);
      onSuccess?.();
    } catch (error: unknown) {
      const err = error as {
        response?: { data?: { message?: string } };
        message?: string;
      };
      const errorMessage =
        err?.response?.data?.message ||
        err?.message ||
        'Erro ao cadastrar produto';

      toast.error('Erro ao cadastrar produto', {
        description: errorMessage,
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Package className="h-5 w-5 text-primary" />
            Novo Produto
          </DialogTitle>
          <DialogDescription>
            Cadastre um novo produto no estoque. Preencha os campos obrigatórios
            (*).
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            {/* Linha 1: Nome e Código de Barras */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <FormField
                control={form.control}
                name="nome"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormLabel>Nome *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Pomada Modeladora"
                        {...field}
                        disabled={createProduct.isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="codigo_barras"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Código de Barras</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="7891234567890"
                        {...field}
                        disabled={createProduct.isPending}
                        className="font-mono"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Linha 2: Categoria, Unidade e Fornecedor */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <FormField
                control={form.control}
                name="categoria_produto_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Categoria *</FormLabel>
                    {categoriasError ? (
                      <div className="flex items-center gap-2 text-destructive text-sm">
                        <AlertCircle className="h-4 w-4" />
                        <span>Erro ao carregar</span>
                      </div>
                    ) : (
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        disabled={createProduct.isPending || categoriasLoading}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue
                              placeholder={
                                categoriasLoading
                                  ? 'Carregando...'
                                  : 'Selecione'
                              }
                            />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {categoriasAtivas.map((cat) => (
                            <SelectItem key={cat.id} value={cat.id}>
                              <span className="flex items-center gap-2">
                                <span
                                  className="w-3 h-3 rounded-full"
                                  style={{ backgroundColor: cat.cor }}
                                />
                                {cat.nome}
                              </span>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="unidade_medida"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Unidade *</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      value={field.value}
                      disabled={createProduct.isPending}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Selecione" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {UNIDADES.map((un) => (
                          <SelectItem key={un.value} value={un.value}>
                            {un.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="fornecedor_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Fornecedor</FormLabel>
                    {fornecedoresError ? (
                      <div className="flex items-center gap-2 text-destructive text-sm">
                        <AlertCircle className="h-4 w-4" />
                        <span>Erro ao carregar</span>
                      </div>
                    ) : (
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        disabled={
                          createProduct.isPending || fornecedoresLoading
                        }
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue
                              placeholder={
                                fornecedoresLoading
                                  ? 'Carregando...'
                                  : 'Selecione'
                              }
                            />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="none">Nenhum</SelectItem>
                          {fornecedoresAtivos.map((forn) => (
                            <SelectItem key={forn.id} value={forn.id}>
                              {getFornecedorNome(forn)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Linha 3: Valores */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <FormField
                control={form.control}
                name="valor_unitario"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Valor Unitário (R$) *</FormLabel>
                    <FormControl>
                      <Input
                        type="text"
                        inputMode="decimal"
                        placeholder="0,00"
                        {...field}
                        onChange={(e) => {
                          const formatted = formatCurrency(e.target.value);
                          field.onChange(formatted);
                        }}
                        disabled={createProduct.isPending}
                        className="text-right"
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Preço de venda ao cliente
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="valor_venda_profissional"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Valor Profissional (R$)</FormLabel>
                    <FormControl>
                      <Input
                        type="text"
                        inputMode="decimal"
                        placeholder="0,00"
                        {...field}
                        onChange={(e) => {
                          const formatted = formatCurrency(e.target.value);
                          field.onChange(formatted);
                        }}
                        disabled={createProduct.isPending}
                        className="text-right"
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Preço para profissionais
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="valor_entrada"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Custo de Aquisição (R$)</FormLabel>
                    <FormControl>
                      <Input
                        type="text"
                        inputMode="decimal"
                        placeholder="0,00"
                        {...field}
                        onChange={(e) => {
                          const formatted = formatCurrency(e.target.value);
                          field.onChange(formatted);
                        }}
                        disabled={createProduct.isPending}
                        className="text-right"
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Custo unitário de compra
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Linha 4: Estoque */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="quantidade_minima"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Estoque Mínimo</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min="0"
                        step="1"
                        placeholder="0"
                        {...field}
                        disabled={createProduct.isPending}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Alerta quando abaixo desse valor
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="quantidade_maxima"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Estoque Máximo</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min="0"
                        step="1"
                        placeholder="0"
                        {...field}
                        disabled={createProduct.isPending}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Quantidade máxima em estoque
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Descrição */}
            <FormField
              control={form.control}
              name="descricao"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Descrição</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Descrição detalhada do produto..."
                      className="resize-none h-20"
                      {...field}
                      disabled={createProduct.isPending}
                    />
                  </FormControl>
                  <FormDescription className="text-xs">
                    {field.value?.length || 0}/500 caracteres
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={createProduct.isPending}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={createProduct.isPending}>
                {createProduct.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Salvando...
                  </>
                ) : (
                  'Cadastrar Produto'
                )}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
