/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Inventário de Estoque
 *
 * Lista produtos com filtros, saldo atual e alertas de estoque baixo.
 * Conforme FLUXO_ESTOQUE.md
 */

'use client';

import { CreateProductModal } from '@/components/estoque/CreateProductModal';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Progress } from '@/components/ui/progress';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { useDeleteStockItem, useStockItems } from '@/hooks/use-stock';
import { useBreadcrumbs } from '@/store/ui-store';
import { StockCategory } from '@/types/stock';
import {
    AlertCircle,
    ArrowDownCircle,
    ArrowUpCircle,
    Edit,
    Package,
    PackageX,
    Plus,
    Search,
    Trash2,
    TrendingDown
} from 'lucide-react';
import Link from 'next/link';
import { useEffect, useState } from 'react';

export default function EstoquePage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState<StockCategory | 'all'>('all');
  const [lowStockOnly, setLowStockOnly] = useState(false);
  const [createModalOpen, setCreateModalOpen] = useState(false);

  // Queries
  const { data, isLoading } = useStockItems({
    search: search || undefined,
    category: category !== 'all' ? category : undefined,
    low_stock: lowStockOnly || undefined,
  });

  const deleteItem = useDeleteStockItem();

  useEffect(() => {
    setBreadcrumbs([
      { label: 'Estoque', href: '/estoque' },
      { label: 'Inventário' },
    ]);
  }, [setBreadcrumbs]);

  const handleDelete = async (id: string, name: string) => {
    if (confirm(`Deseja realmente excluir "${name}"?`)) {
      await deleteItem.mutateAsync(id);
    }
  };

  // Formatação de moeda
  const formatCurrency = (value: string) => {
    const num = parseFloat(value);
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'BRL',
    }).format(num);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Inventário</h1>
          <p className="text-muted-foreground">
            Gerencie o estoque de produtos e insumos
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" asChild>
            <Link href="/estoque/categorias">
              <Package className="mr-2 h-4 w-4" />
              Categorias
            </Link>
          </Button>
          <Button variant="outline" asChild>
            <Link href="/estoque/entrada">
              <ArrowDownCircle className="mr-2 h-4 w-4" />
              Entrada
            </Link>
          </Button>
          <Button variant="outline" asChild>
            <Link href="/estoque/saida">
              <ArrowUpCircle className="mr-2 h-4 w-4" />
              Saída
            </Link>
          </Button>
          <Button onClick={() => setCreateModalOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Novo Produto
          </Button>
        </div>
      </div>

      {/* Modal de Cadastro de Produto */}
      <CreateProductModal
        open={createModalOpen}
        onOpenChange={setCreateModalOpen}
      />

      {/* Resumo - RN-EST-005 */}
      {data && (
        <div className="grid gap-4 md:grid-cols-3">
          <Card className="transition-all hover:shadow-md">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Total de Itens
              </CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{data.total}</div>
              <p className="text-xs text-muted-foreground mt-1">
                produtos cadastrados
              </p>
            </CardContent>
          </Card>

          <Card className="transition-all hover:shadow-md">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Estoque Baixo
              </CardTitle>
              <AlertCircle className="h-4 w-4 text-orange-500" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-orange-500">
                {data.low_stock_count}
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                requer atenção
              </p>
            </CardContent>
          </Card>

          <Card className="transition-all hover:shadow-md">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Sem Estoque
              </CardTitle>
              <TrendingDown className="h-4 w-4 text-destructive" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-destructive">
                {data.out_of_stock_count}
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                esgotados
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Filtros */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col gap-4 md:flex-row">
            {/* Busca */}
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome ou SKU..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-10"
              />
            </div>

            {/* Categoria - RN-EST-001 */}
            <Select
              value={category}
              onValueChange={(value) =>
                setCategory(value as StockCategory | 'all')
              }
            >
              <SelectTrigger className="w-full md:w-[200px]">
                <SelectValue placeholder="Categoria" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas</SelectItem>
                <SelectItem value={StockCategory.POMADA}>Pomada</SelectItem>
                <SelectItem value={StockCategory.SHAMPOO}>Shampoo</SelectItem>
                <SelectItem value={StockCategory.CREME}>Creme</SelectItem>
                <SelectItem value={StockCategory.LAMINA}>Lâmina</SelectItem>
                <SelectItem value={StockCategory.TOALHA}>Toalha</SelectItem>
                <SelectItem value={StockCategory.LIMPEZA}>Limpeza</SelectItem>
                <SelectItem value={StockCategory.ESCRITORIO}>Escritório</SelectItem>
                <SelectItem value={StockCategory.BEBIDA}>Bebida</SelectItem>
                <SelectItem value={StockCategory.REVENDA}>Revenda</SelectItem>
                <SelectItem value={StockCategory.INSUMO}>Insumo</SelectItem>
                <SelectItem value={StockCategory.USO_INTERNO}>Uso Interno</SelectItem>
                <SelectItem value={StockCategory.PERMANENTE}>Permanente</SelectItem>
                <SelectItem value={StockCategory.PROMOCIONAL}>Promocional</SelectItem>
                <SelectItem value={StockCategory.KIT}>Kit</SelectItem>
                <SelectItem value={StockCategory.SERVICO}>Serviço</SelectItem>
              </SelectContent>
            </Select>

            {/* Filtro Estoque Baixo */}
            <Button
              variant={lowStockOnly ? 'default' : 'outline'}
              onClick={() => setLowStockOnly(!lowStockOnly)}
              className="w-full md:w-auto"
            >
              <AlertCircle className="mr-2 h-4 w-4" />
              Apenas Baixo
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Lista de Produtos - Design Minimalista */}
      {isLoading ? (
        <div className="space-y-2">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <Skeleton key={i} className="h-16 w-full" />
          ))}
        </div>
      ) : data && data.data && data.data.length > 0 ? (
        <Card>
          <CardContent className="p-0">
            <div className="divide-y">
              {data.data.map((item) => {
                const qtdAtual = parseFloat(item.quantidade_atual);
                const qtdMinima = parseFloat(item.quantidade_minima);
                const isLow = qtdAtual <= qtdMinima;
                const isOut = qtdAtual === 0;
                
                const maxForProgress = qtdMinima * 2;
                const stockPercentage = Math.min(100, Math.max(0, (qtdAtual / maxForProgress) * 100));

                return (
                  <div
                    key={item.id}
                    className={`group flex items-center gap-4 p-4 transition-colors hover:bg-muted/50 ${
                      isOut ? 'bg-destructive/5' : isLow ? 'bg-orange-50/50 dark:bg-orange-950/10' : ''
                    }`}
                  >
                    {/* Nome e Categoria */}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h3 className="font-semibold truncate">{item.nome}</h3>
                        {isOut && (
                          <Badge variant="destructive" className="shrink-0">Sem Estoque</Badge>
                        )}
                        {isLow && !isOut && (
                          <Badge variant="outline" className="border-orange-500 text-orange-600 shrink-0">
                            Baixo
                          </Badge>
                        )}
                      </div>
                      {item.categoria_produto && (
                        <p className="text-xs text-muted-foreground mt-0.5">
                          {item.categoria_produto.nome}
                        </p>
                      )}
                    </div>

                    {/* Estoque com Progress */}
                    <div className="hidden md:flex items-center gap-3 w-48">
                      <div className="flex-1">
                        <div className="flex items-center justify-between text-xs mb-1">
                          <span className="text-muted-foreground">Estoque</span>
                          <span className={`font-semibold ${
                            isOut ? 'text-destructive' : isLow ? 'text-orange-600' : 'text-green-600'
                          }`}>
                            {item.quantidade_atual}
                          </span>
                        </div>
                        <Progress 
                          value={stockPercentage} 
                          className={`h-1.5 ${
                            isOut ? '[&>*]:bg-destructive' : isLow ? '[&>*]:bg-orange-500' : '[&>*]:bg-green-500'
                          }`}
                        />
                      </div>
                    </div>

                    {/* Valor */}
                    <div className="hidden sm:block text-right w-24">
                      <p className="text-sm font-semibold">{formatCurrency(item.valor_unitario)}</p>
                      <p className="text-xs text-muted-foreground">{item.unidade_medida}</p>
                    </div>

                    {/* Ações */}
                    <div className="flex gap-1 shrink-0">
                      <Button 
                        variant="ghost" 
                        size="icon"
                        className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
                        asChild
                      >
                        <Link href={`/estoque/${item.id}/editar`}>
                          <Edit className="h-4 w-4" />
                        </Link>
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={() => handleDelete(item.id, item.nome)}
                        disabled={deleteItem.isPending}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      ) : (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <PackageX className="h-12 w-12 text-muted-foreground mb-3" />
            <h3 className="font-semibold mb-1">Nenhum produto encontrado</h3>
            <p className="text-sm text-muted-foreground mb-4 max-w-sm">
              {search || category !== 'all' || lowStockOnly
                ? 'Ajuste os filtros para ver mais produtos.'
                : 'Comece cadastrando seu primeiro produto.'}
            </p>
            {!search && category === 'all' && !lowStockOnly && (
              <Button onClick={() => setCreateModalOpen(true)} size="sm">
                <Plus className="mr-2 h-4 w-4" />
                Cadastrar Produto
              </Button>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
