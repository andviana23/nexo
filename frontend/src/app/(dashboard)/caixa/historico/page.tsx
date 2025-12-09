/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Histórico Page
 *
 * Página de histórico do módulo Caixa Diário.
 *
 * @author NEXO v2.0
 */

'use client';

import {
    AlertCircle,
    ArrowLeft,
    Calendar,
    CheckCircle2,
    ChevronLeft,
    ChevronRight,
    Eye,
    RefreshCw,
} from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';

import { useCaixaHistorico } from '@/hooks/use-caixa';
import { cn } from '@/lib/utils';
import type { ListCaixaHistoricoFilters } from '@/types/caixa';
import { StatusCaixa, StatusCaixaColors, StatusCaixaLabels } from '@/types/caixa';

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

const formatDateTime = (dateStr: string | undefined) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

// =============================================================================
// PAGE COMPONENT
// =============================================================================

export default function CaixaHistoricoPage() {
  const [filters, setFilters] = useState<ListCaixaHistoricoFilters>({
    page: 1,
    page_size: 10,
  });

  const { data, isLoading, isError, error, refetch } = useCaixaHistorico(filters);

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const handleFilterChange = (key: keyof ListCaixaHistoricoFilters, value: string) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value || undefined,
      page: 1, // Reset para primeira página ao filtrar
    }));
  };

  // Loading state
  if (isLoading && !data) {
    return (
      <div className="flex-1 space-y-6 p-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-64" />
      </div>
    );
  }

  // Error state
  if (isError) {
    return (
      <div className="flex-1 p-6">
        <div className="rounded-lg border border-destructive bg-destructive/10 p-6 text-center">
          <h2 className="text-lg font-semibold text-destructive mb-2">
            Erro ao carregar histórico
          </h2>
          <p className="text-muted-foreground mb-4">
            {error?.message || 'Ocorreu um erro inesperado'}
          </p>
          <Button onClick={() => refetch()}>
            <RefreshCw className="mr-2 h-4 w-4" />
            Tentar novamente
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 space-y-6 p-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-4">
          <Link href="/caixa">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">
              Histórico de Caixas
            </h1>
            <p className="text-muted-foreground">
              Consulte os caixas anteriores
            </p>
          </div>
        </div>

        <Button variant="outline" size="sm" onClick={() => refetch()}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Atualizar
        </Button>
      </div>

      <Separator />

      {/* Filtros */}
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">Filtros</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-4">
            <div className="flex flex-col gap-1.5">
              <label className="text-sm text-muted-foreground">
                Data Início
              </label>
              <Input
                type="date"
                value={filters.data_inicio || ''}
                onChange={(e) =>
                  handleFilterChange('data_inicio', e.target.value)
                }
                className="w-40"
              />
            </div>
            <div className="flex flex-col gap-1.5">
              <label className="text-sm text-muted-foreground">Data Fim</label>
              <Input
                type="date"
                value={filters.data_fim || ''}
                onChange={(e) => handleFilterChange('data_fim', e.target.value)}
                className="w-40"
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tabela */}
      <Card>
        <CardHeader>
          <CardTitle>Caixas</CardTitle>
          <CardDescription>
            {data?.total
              ? `${data.total} registro(s) encontrado(s)`
              : 'Nenhum registro encontrado'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {!data?.items?.length ? (
            <div className="text-center py-8 text-muted-foreground">
              <Calendar className="mx-auto h-12 w-12 mb-4 opacity-50" />
              <p>Nenhum caixa encontrado para os filtros selecionados.</p>
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Abertura</TableHead>
                    <TableHead>Fechamento</TableHead>
                    <TableHead>Responsável</TableHead>
                    <TableHead>Saldo Inicial</TableHead>
                    <TableHead>Saldo Final</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Divergência</TableHead>
                    <TableHead className="w-20">Ações</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {data.items.map((caixa) => (
                    <TableRow key={caixa.id}>
                      <TableCell className="font-mono text-sm">
                        {formatDateTime(caixa.data_abertura)}
                      </TableCell>
                      <TableCell className="font-mono text-sm">
                        {caixa.data_fechamento
                          ? formatDateTime(caixa.data_fechamento)
                          : '-'}
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-col">
                          <span className="text-sm">
                            {caixa.usuario_abertura_nome}
                          </span>
                          {caixa.usuario_fechamento_nome && (
                            <span className="text-xs text-muted-foreground">
                              Fechado por: {caixa.usuario_fechamento_nome}
                            </span>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>{formatCurrency(caixa.saldo_inicial)}</TableCell>
                      <TableCell>
                        {caixa.saldo_real
                          ? formatCurrency(caixa.saldo_real)
                          : formatCurrency(caixa.saldo_esperado)}
                      </TableCell>
                      <TableCell>
                        <Badge
                          className={cn(
                            StatusCaixaColors[caixa.status as StatusCaixa]
                          )}
                        >
                          {StatusCaixaLabels[caixa.status as StatusCaixa]}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {caixa.tem_divergencia ? (
                          <div className="flex items-center gap-1 text-amber-600">
                            <AlertCircle className="h-4 w-4" />
                            <span className="font-medium">
                              {formatCurrency(caixa.divergencia)}
                            </span>
                          </div>
                        ) : (
                          <div className="flex items-center gap-1 text-green-600">
                            <CheckCircle2 className="h-4 w-4" />
                            <span className="text-sm">OK</span>
                          </div>
                        )}
                      </TableCell>
                      <TableCell>
                        <Link href={`/caixa/${caixa.id}`}>
                          <Button variant="ghost" size="icon">
                            <Eye className="h-4 w-4" />
                          </Button>
                        </Link>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>

              {/* Paginação */}
              {data.total_pages > 1 && (
                <div className="flex items-center justify-between mt-4 pt-4 border-t">
                  <p className="text-sm text-muted-foreground">
                    Página {data.page} de {data.total_pages}
                  </p>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePageChange(data.page - 1)}
                      disabled={data.page <= 1}
                    >
                      <ChevronLeft className="h-4 w-4" />
                      Anterior
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePageChange(data.page + 1)}
                      disabled={data.page >= data.total_pages}
                    >
                      Próxima
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
