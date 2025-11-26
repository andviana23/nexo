/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Inventário de Estoque
 *
 * Lista produtos com filtros, saldo atual e alertas de estoque baixo.
 * Conforme FLUXO_ESTOQUE.md
 */

'use client';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
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
    Edit,
    Package,
    Plus,
    Search,
    Trash2,
    TrendingDown,
} from 'lucide-react';
import Link from 'next/link';
import { useEffect, useState } from 'react';

export default function EstoquePage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState<StockCategory | 'all'>('all');
  const [lowStockOnly, setLowStockOnly] = useState(false);

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
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Inventário</h1>
          <p className="text-muted-foreground">
            Gerencie o estoque de produtos e insumos
          </p>
        </div>
        <Button asChild>
          <Link href="/estoque/novo">
            <Plus className="mr-2 h-4 w-4" />
            Novo Produto
          </Link>
        </Button>
      </div>

      {/* Resumo - RN-EST-005 */}
      {data && (
        <div className="grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Total de Itens
              </CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{data.total}</div>
            </CardContent>
          </Card>

          <Card>
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
            </CardContent>
          </Card>

          <Card>
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
                <SelectItem value={StockCategory.PRODUCT}>Produto</SelectItem>
                <SelectItem value={StockCategory.CONSUMABLE}>
                  Consumível
                </SelectItem>
                <SelectItem value={StockCategory.EQUIPMENT}>
                  Equipamento
                </SelectItem>
                <SelectItem value={StockCategory.OTHER}>Outro</SelectItem>
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

      {/* Lista de Produtos */}
      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-24 w-full" />
          ))}
        </div>
      ) : data && data.items.length > 0 ? (
        <div className="space-y-4">
          {data.items.map((item) => {
            const isLow = item.current_quantity <= item.min_quantity;
            const isOut = item.current_quantity === 0;

            return (
              <Card key={item.id} className={isOut ? 'border-destructive' : ''}>
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="flex-1 space-y-2">
                      <div className="flex items-center gap-3">
                        <h3 className="text-lg font-semibold">{item.name}</h3>
                        <Badge variant="outline">{item.category}</Badge>
                        {isOut && (
                          <Badge variant="destructive">Sem Estoque</Badge>
                        )}
                        {isLow && !isOut && (
                          <Badge variant="outline" className="border-orange-500 text-orange-500">
                            Estoque Baixo
                          </Badge>
                        )}
                      </div>

                      {item.description && (
                        <p className="text-sm text-muted-foreground">
                          {item.description}
                        </p>
                      )}

                      <div className="grid grid-cols-2 gap-4 pt-2 md:grid-cols-4">
                        <div>
                          <p className="text-xs text-muted-foreground">SKU</p>
                          <p className="font-mono text-sm font-medium">
                            {item.sku || '-'}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">
                            Estoque Atual
                          </p>
                          <p
                            className={`text-sm font-bold ${
                              isOut
                                ? 'text-destructive'
                                : isLow
                                ? 'text-orange-500'
                                : 'text-green-600'
                            }`}
                          >
                            {item.current_quantity} {item.unit}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">
                            Estoque Mínimo
                          </p>
                          <p className="text-sm font-medium">
                            {item.min_quantity} {item.unit}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground">
                            Valor Unitário
                          </p>
                          <p className="text-sm font-medium">
                            {formatCurrency(item.cost_price)}
                          </p>
                        </div>
                      </div>
                    </div>

                    {/* Ações */}
                    <div className="flex gap-2">
                      <Button variant="ghost" size="icon" asChild>
                        <Link href={`/estoque/${item.id}/editar`}>
                          <Edit className="h-4 w-4" />
                        </Link>
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDelete(item.id, item.name)}
                        disabled={deleteItem.isPending}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      ) : (
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            Nenhum produto encontrado. Cadastre o primeiro produto para começar.
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
