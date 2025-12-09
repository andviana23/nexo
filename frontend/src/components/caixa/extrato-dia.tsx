/**
 * NEXO - Sistema de Gestão para Barbearias
 * Extrato do Dia Component
 *
 * Exibe a lista de operações do caixa do dia (vendas, sangrias, reforços).
 *
 * @author NEXO v2.0
 */

'use client';

import { ArrowDownCircle, ArrowUpCircle, Clock, MinusCircle, PlusCircle, User } from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';

import { cn } from '@/lib/utils';
import type { OperacaoCaixaResponse } from '@/types/caixa';
import {
    DestinoSangria,
    DestinoSangriaLabels,
    OrigemReforco,
    OrigemReforcoLabels,
    TipoOperacaoCaixa,
    TipoOperacaoColors,
    TipoOperacaoLabels,
} from '@/types/caixa';

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

const formatTime = (dateStr: string | undefined) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('pt-BR', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

const getOperacaoIcon = (tipo: TipoOperacaoCaixa) => {
  switch (tipo) {
    case TipoOperacaoCaixa.VENDA:
      return <ArrowUpCircle className="h-4 w-4" />;
    case TipoOperacaoCaixa.SANGRIA:
      return <MinusCircle className="h-4 w-4" />;
    case TipoOperacaoCaixa.REFORCO:
      return <PlusCircle className="h-4 w-4" />;
    case TipoOperacaoCaixa.DESPESA:
      return <ArrowDownCircle className="h-4 w-4" />;
    default:
      return null;
  }
};

const getDestinoOrigemLabel = (operacao: OperacaoCaixaResponse): string => {
  if (operacao.destino) {
    return DestinoSangriaLabels[operacao.destino as DestinoSangria] || operacao.destino;
  }
  if (operacao.origem) {
    return OrigemReforcoLabels[operacao.origem as OrigemReforco] || operacao.origem;
  }
  return '';
};

// =============================================================================
// TYPES
// =============================================================================

interface ExtratoDiaProps {
  operacoes: OperacaoCaixaResponse[] | undefined;
  isLoading: boolean;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ExtratoDia({ operacoes, isLoading }: ExtratoDiaProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Extrato do Dia</CardTitle>
          <CardDescription>Operações realizadas hoje</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const operacoesOrdenadas = [...(operacoes || [])].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Extrato do Dia</CardTitle>
        <CardDescription>
          {operacoesOrdenadas.length > 0
            ? `${operacoesOrdenadas.length} operação(ões) registrada(s)`
            : 'Nenhuma operação registrada ainda'}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {operacoesOrdenadas.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            <p>Nenhuma operação registrada no caixa atual.</p>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[100px]">Hora</TableHead>
                <TableHead>Tipo</TableHead>
                <TableHead>Descrição</TableHead>
                <TableHead>Responsável</TableHead>
                <TableHead className="text-right">Valor</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {operacoesOrdenadas.map((operacao) => (
                <TableRow key={operacao.id}>
                  <TableCell className="font-mono text-sm">
                    <div className="flex items-center gap-1">
                      <Clock className="h-3 w-3 text-muted-foreground" />
                      {formatTime(operacao.created_at)}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant="outline"
                      className={cn(
                        'flex items-center gap-1 w-fit',
                        TipoOperacaoColors[operacao.tipo as TipoOperacaoCaixa]
                      )}
                    >
                      {getOperacaoIcon(operacao.tipo as TipoOperacaoCaixa)}
                      {TipoOperacaoLabels[operacao.tipo as TipoOperacaoCaixa]}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col">
                      <span className="text-sm">{operacao.descricao}</span>
                      {(operacao.destino || operacao.origem) && (
                        <span className="text-xs text-muted-foreground">
                          {operacao.destino ? 'Destino' : 'Origem'}:{' '}
                          {getDestinoOrigemLabel(operacao)}
                        </span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1 text-sm text-muted-foreground">
                      <User className="h-3 w-3" />
                      {operacao.usuario_nome}
                    </div>
                  </TableCell>
                  <TableCell
                    className={cn(
                      'text-right font-medium',
                      TipoOperacaoColors[operacao.tipo as TipoOperacaoCaixa]
                    )}
                  >
                    {operacao.tipo === TipoOperacaoCaixa.SANGRIA ||
                    operacao.tipo === TipoOperacaoCaixa.DESPESA
                      ? '-'
                      : '+'}
                    {formatCurrency(operacao.valor)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  );
}

export default ExtratoDia;
