/**
 * CommandModal - Modal de Fechamento de Comanda
 * Design: Estilo Trinks - Compacto, 2 Colunas, SEM SCROLL
 *
 * @version 5.0.0 - Refatoração completa estilo PDV profissional
 *
 * ARQUITETURA DO LAYOUT:
 * ┌──────────────────────────────────────────────────────────────────┐
 * │ HEADER: Fechamento de Conta DD/MM/YYYY                     ? X  │
 * ├─────────────────────────────────┬────────────────────────────────┤
 * │ COLUNA ESQUERDA                 │ COLUNA DIREITA                 │
 * │                                 │                                │
 * │ ┌─────────────────────────────┐ │ ┌────────────────────────────┐ │
 * │ │ Avatar | Nome               │ │ │ Total              R$ XX  │ │
 * │ │        | Telefone    Editar │ │ ├────────────────────────────┤ │
 * │ ├─────────────────────────────┤ │ │ [Select] Forma Pagamento   │ │
 * │ │ Comanda | Desde | Pontos    │ │ │ [Input]  Valor             │ │
 * │ └─────────────────────────────┘ │ ├────────────────────────────┤ │
 * │                                 │ │ Recebido           R$ XX   │ │
 * │ [Serviços][Produtos][Pacotes]   │ │ Falta              R$ XX   │ │
 * │                                 │ ├────────────────────────────┤ │
 * │ ┌─────────────────────────────┐ │ │ [ ] Deixar como dívida     │ │
 * │ │ Item      Preço Desc A pagar│ │ └────────────────────────────┘ │
 * │ │ Corte     50.00 0.00  50.00 │ │                                │
 * │ └─────────────────────────────┘ │                                │
 * │                                 │                                │
 * │ [Produtos usados nos serviços]  │                                │
 * ├─────────────────────────────────┴────────────────────────────────┤
 * │ FOOTER:                              [Cancelar] [Fechar Conta]   │
 * └──────────────────────────────────────────────────────────────────┘
 */

'use client';

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  useAddCommandPayment,
  useCloseCommand,
  useCommand,
  useUpdateCommandItem,
} from '@/hooks/use-commands';
import { useCustomer } from '@/hooks/use-customers';
import { useMeiosPagamento } from '@/hooks/use-meios-pagamento';
import { cn } from '@/lib/utils';
import type { CommandPaymentResponse, SelectedPayment } from '@/types/command';
import {
  calcularResumoFinanceiro,
  calcularValorLiquido,
  canCloseCommand,
  formatMoney,
} from '@/types/command';
import { TIPO_PAGAMENTO_LABELS } from '@/types/meio-pagamento';
import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
  Gift,
  HelpCircle,
  Package,
  Pencil,
  Scissors,
  Star,
  Tag,
  Ticket,
  X,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

interface CommandModalProps {
  commandId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CommandModal({
  commandId,
  open,
  onOpenChange,
}: CommandModalProps) {
  const mapPayments = (payments: CommandPaymentResponse[] | undefined): SelectedPayment[] =>
    (payments || []).map((p) => ({
      meio_pagamento_id: p.meio_pagamento_id,
      nome: '',
      tipo: '',
      valor_recebido: parseFloat(p.valor_recebido),
      taxa_percentual: p.taxa_percentual,
      taxa_fixa: parseFloat(p.taxa_fixa),
      valor_liquido: parseFloat(p.valor_liquido),
    }));

  const [selectedPayments, setSelectedPayments] = useState<SelectedPayment[]>([]);
  const [deixarSaldoDivida, setDeixarSaldoDivida] = useState(false);
  const [selectedMeioId, setSelectedMeioId] = useState<string>('');
  const [valorInput, setValorInput] = useState<string>('');

  const resetState = (payments?: CommandPaymentResponse[]) => {
    setSelectedPayments(mapPayments(payments));
    setDeixarSaldoDivida(false);
    setSelectedMeioId('');
    setValorInput('');
  };

  const { data: command } = useCommand(commandId);
  const { data: customer } = useCustomer(command?.customer_id || '');
  const { data: meiosResponse } = useMeiosPagamento({ apenas_ativos: true });
  const addPayment = useAddCommandPayment();
  const closeCommand = useCloseCommand();
  const updateItem = useUpdateCommandItem();

  const handleOpenChange = (value: boolean) => {
    resetState(command?.payments);
    onOpenChange(value);
  };

  if (!command) return null;

  const meios = meiosResponse?.data || [];
  const subtotal = parseFloat(command.subtotal);
  const desconto = parseFloat(command.desconto);
  const resumo = calcularResumoFinanceiro(subtotal, desconto, selectedPayments);
  const validation = canCloseCommand(resumo, deixarSaldoDivida);

  // Helpers
  const getInitials = (nome: string) => {
    const parts = nome.split(' ').filter(Boolean);
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[parts.length - 1][0]}`.toUpperCase();
    }
    return nome.slice(0, 2).toUpperCase();
  };

  const getMemberYear = (createdAt: string) => {
    return new Date(createdAt).getFullYear();
  };

  // Adicionar pagamento via select + input
  const handleAddPayment = () => {
    if (!selectedMeioId || !valorInput) return;

    const meio = meios.find((m) => m.id === selectedMeioId);
    if (!meio) return;

    const valor = parseFloat(valorInput) || 0;
    if (valor <= 0) return;

    const taxa = parseFloat(meio.taxa) || 0;
    const taxaFixa = parseFloat(meio.taxa_fixa) || 0;

    const label = TIPO_PAGAMENTO_LABELS[meio.tipo as import('@/types/meio-pagamento').TipoPagamento] || meio.tipo;
    const nomeFormatado = meio.bandeira ? `${label} - ${meio.bandeira}` : label;

    const payment: SelectedPayment = {
      meio_pagamento_id: meio.id,
      nome: nomeFormatado,
      tipo: meio.tipo,
      valor_recebido: valor,
      taxa_percentual: taxa,
      taxa_fixa: taxaFixa,
      valor_liquido: calcularValorLiquido(valor, taxa, taxaFixa),
    };

    setSelectedPayments((prev) => [...prev, payment]);
    setSelectedMeioId('');
    setValorInput('');
  };

  // Preencher valor restante automaticamente
  const handleFillRemaining = () => {
    if (resumo.falta > 0) {
      setValorInput(resumo.falta.toFixed(2));
    }
  };

  // Remover pagamento
  const handleRemovePayment = (index: number) => {
    setSelectedPayments((prev) => prev.filter((_, i) => i !== index));
  };

  // Fechar comanda
  const handleClose = async () => {
    if (!validation.valid) {
      validation.errors.forEach((err) => toast.error(err));
      return;
    }

    try {
      for (const payment of selectedPayments) {
        if (payment.valor_recebido > 0) {
          await addPayment.mutateAsync({
            commandId: command.id,
            data: {
              meio_pagamento_id: payment.meio_pagamento_id,
              valor_recebido: payment.valor_recebido.toFixed(2), // Enviar como string
            },
          });
        }
      }

      await closeCommand.mutateAsync({
        commandId: command.id,
        data: {
          deixar_troco_gorjeta: false,
          deixar_saldo_divida: deixarSaldoDivida,
          observacoes: '',
        },
      });

      toast.success('Comanda fechada com sucesso!');
      onOpenChange(false);
    } catch {
      // Erros tratados pelos hooks
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent
        className="max-w-4xl w-[95vw] p-0 gap-0 overflow-hidden"
        showCloseButton={false}
      >
        {/* ══════════════ HEADER ══════════════ */}
        <DialogHeader className="bg-primary px-5 py-3 flex-row items-center justify-between">
          <DialogTitle className="text-sm font-semibold text-primary-foreground uppercase tracking-wider">
            Fechamento de Conta do Dia {format(new Date(), 'dd/MM/yyyy', { locale: ptBR })}
          </DialogTitle>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-primary-foreground/80 hover:text-primary-foreground hover:bg-primary-foreground/10"
            >
              <HelpCircle className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-primary-foreground/80 hover:text-primary-foreground hover:bg-primary-foreground/10"
              onClick={() => handleOpenChange(false)}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        </DialogHeader>

        {/* ══════════════ BODY - 2 COLUNAS ══════════════ */}
        <div className="grid grid-cols-1 lg:grid-cols-[1fr,320px] min-h-[400px]">
          {/* ────────── COLUNA ESQUERDA ────────── */}
          <div className="p-5 space-y-4 border-r border-border">
            {/* Cliente Card - Compacto */}
            {customer && (
              <div className="flex items-center justify-between gap-4 pb-3 border-b border-border">
                <div className="flex items-center gap-3">
                  <Avatar className="h-10 w-10 bg-primary text-primary-foreground">
                    <AvatarFallback className="bg-primary text-primary-foreground font-semibold text-sm">
                      {getInitials(customer.nome)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-semibold text-foreground leading-tight">{customer.nome}</p>
                    <p className="text-xs text-muted-foreground">{customer.telefone || '(00) 0000-0000'}</p>
                  </div>
                </div>
                <button className="text-xs text-primary hover:underline flex items-center gap-1">
                  <Pencil className="h-3 w-3" />
                  Editar cliente
                </button>
              </div>
            )}

            {/* Info Grid - Inline */}
            {customer && (
              <div className="grid grid-cols-3 text-center text-xs py-2 bg-muted/30 rounded-md">
                <div>
                  <span className="text-muted-foreground">Comanda</span>
                  <p className="font-semibold">{command.numero || '-'}</p>
                </div>
                <div className="border-x border-border">
                  <span className="text-muted-foreground">Cliente desde</span>
                  <p className="font-semibold">{getMemberYear(customer.criado_em)}</p>
                </div>
                <div className="flex flex-col items-center">
                  <div className="flex items-center gap-1">
                    <Star className="h-3 w-3 text-amber-500 fill-amber-500" />
                    <span className="text-muted-foreground">Pontos</span>
                  </div>
                  <button className="text-primary text-[10px] hover:underline">Carregar</button>
                </div>
              </div>
            )}

            {/* Tabs - Compactas */}
            <div className="flex flex-wrap gap-1.5">
              <Button size="sm" variant="outline" className="h-7 text-xs gap-1 px-2">
                <Scissors className="h-3 w-3" />
                Serviços
              </Button>
              <Button size="sm" className="h-7 text-xs gap-1 px-2 bg-primary">
                <Package className="h-3 w-3" />
                Produtos
              </Button>
              <Button size="sm" className="h-7 text-xs gap-1 px-2 bg-emerald-600 hover:bg-emerald-700">
                <Gift className="h-3 w-3" />
                Pacotes
              </Button>
              <Button size="sm" className="h-7 text-xs gap-1 px-2 bg-cyan-600 hover:bg-cyan-700">
                <Ticket className="h-3 w-3" />
                Vales
              </Button>
              <Button size="sm" className="h-7 text-xs gap-1 px-2 bg-amber-600 hover:bg-amber-700">
                <Tag className="h-3 w-3" />
                Cupom
              </Button>
            </div>

            {/* Tabela de Itens - Compacta */}
            <div className="border border-border rounded-md overflow-hidden">
              <table className="w-full text-xs">
                <thead>
                  <tr className="bg-muted/50 border-b border-border">
                    <th className="text-left p-2 font-medium text-muted-foreground">Item</th>
                    <th className="text-right p-2 font-medium text-muted-foreground w-20">Preço</th>
                    <th className="text-right p-2 font-medium text-muted-foreground w-16">Desc</th>
                    <th className="text-right p-2 font-medium text-muted-foreground w-20">A pagar</th>
                  </tr>
                </thead>
                <tbody>
                  {(command.items || []).length === 0 ? (
                    <tr>
                      <td colSpan={4} className="p-4 text-center text-muted-foreground">
                        <Scissors className="h-5 w-5 mx-auto mb-1 opacity-40" />
                        <span className="text-[10px]">Nenhum item</span>
                      </td>
                    </tr>
                  ) : (
                    (command.items || []).map((item) => (
                      <tr key={item.id} className="border-b border-border last:border-0 hover:bg-muted/20">
                        <td className="p-2">
                          <div className="flex items-center gap-1.5">
                            <div className="w-0.5 h-4 bg-primary rounded-full" />
                            <Scissors className="h-3 w-3 text-primary" />
                            <span className="font-medium truncate max-w-[150px]">{item.descricao}</span>
                          </div>
                        </td>
                        <td className="p-2">
                          <Input
                            type="text"
                            defaultValue={item.preco_unitario}
                            className="h-6 w-16 text-right text-xs px-1 ml-auto"
                            onChange={(e) => {
                              updateItem.mutate({
                                commandId: command.id,
                                itemId: item.id,
                                data: { preco_unitario: e.target.value },
                              });
                            }}
                          />
                        </td>
                        <td className="p-2">
                          <Input
                            type="text"
                            defaultValue={item.desconto_valor}
                            className="h-6 w-12 text-right text-xs px-1 ml-auto"
                            onChange={(e) => {
                              updateItem.mutate({
                                commandId: command.id,
                                itemId: item.id,
                                data: { desconto_valor: parseFloat(e.target.value) || 0 },
                              });
                            }}
                          />
                        </td>
                        <td className="text-right p-2 font-semibold">
                          {formatMoney(parseFloat(item.preco_final))}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* Produtos usados - Compacto */}
            <Button
              variant="outline"
              size="sm"
              className="w-full h-8 text-xs justify-start gap-2 border-dashed text-muted-foreground"
            >
              <Package className="h-3 w-3" />
              Produtos usados nos serviços
            </Button>
          </div>

          {/* ────────── COLUNA DIREITA ────────── */}
          <div className="p-5 bg-muted/20 flex flex-col">
            {/* Total - Destaque */}
            <div className="flex justify-between items-center pb-3 border-b border-border mb-4">
              <span className="text-sm font-medium text-muted-foreground">Total</span>
              <span className="text-2xl font-bold">{formatMoney(resumo.total)}</span>
            </div>

            {/* Seletor de Pagamento - Compacto */}
            <div className="space-y-3 flex-1">
              <div className="space-y-2">
                <Select value={selectedMeioId} onValueChange={setSelectedMeioId}>
                  <SelectTrigger className="h-9 text-sm">
                    <SelectValue placeholder="Selecione forma de pagamento" />
                  </SelectTrigger>
                  <SelectContent>
                    {meios.map((meio) => (
                      <SelectItem key={meio.id} value={meio.id} className="text-sm">
                        {TIPO_PAGAMENTO_LABELS[meio.tipo as import('@/types/meio-pagamento').TipoPagamento] || meio.tipo}
                        {meio.bandeira ? ` - ${meio.bandeira}` : ''}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <div className="flex gap-2">
                  <Input
                    type="number"
                    step="0.01"
                    placeholder="Valor"
                    value={valorInput}
                    onChange={(e) => setValorInput(e.target.value)}
                    className="h-9 text-sm flex-1"
                  />
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-9 text-xs px-2"
                    onClick={handleFillRemaining}
                    disabled={resumo.falta <= 0}
                  >
                    Preencher
                  </Button>
                  <Button
                    size="sm"
                    className="h-9 text-xs px-3"
                    onClick={handleAddPayment}
                    disabled={!selectedMeioId || !valorInput}
                  >
                    Add
                  </Button>
                </div>
              </div>

              {/* Pagamentos Adicionados */}
              {selectedPayments.length > 0 && (
                <div className="space-y-1 max-h-24 overflow-y-auto">
                  {selectedPayments.map((p, idx) => (
                    <div
                      key={idx}
                      className="flex items-center justify-between text-xs bg-background rounded px-2 py-1.5 border border-border"
                    >
                      <span className="truncate">{p.nome || 'Pagamento'}</span>
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{formatMoney(p.valor_recebido)}</span>
                        <button
                          onClick={() => handleRemovePayment(idx)}
                          className="text-destructive hover:text-destructive/80"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Resumo */}
              <div className="space-y-2 pt-3 border-t border-border mt-auto">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Recebido</span>
                  <span className={cn('font-semibold', resumo.total_recebido > 0 && 'text-primary')}>
                    {formatMoney(resumo.total_recebido)}
                  </span>
                </div>
                {resumo.falta > 0 && (
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Falta</span>
                    <span className="font-semibold text-destructive">{formatMoney(resumo.falta)}</span>
                  </div>
                )}
                {resumo.troco > 0 && (
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Troco</span>
                    <span className="font-semibold text-emerald-600">{formatMoney(resumo.troco)}</span>
                  </div>
                )}
              </div>

              {/* Checkbox Dívida */}
              {resumo.falta > 0 && (
                <label className="flex items-center gap-2 p-2 bg-amber-50 dark:bg-amber-950/30 border border-amber-200 dark:border-amber-800 rounded-md cursor-pointer mt-3">
                  <Checkbox
                    id="deixar-divida"
                    checked={deixarSaldoDivida}
                    onCheckedChange={(checked) => setDeixarSaldoDivida(checked as boolean)}
                  />
                  <span className="text-xs font-medium text-amber-800 dark:text-amber-200">
                    Deixar saldo como dívida ({formatMoney(resumo.falta)})
                  </span>
                </label>
              )}
            </div>
          </div>
        </div>

        {/* ══════════════ FOOTER ══════════════ */}
        <DialogFooter className="border-t border-border px-5 py-3 bg-background">
          <div className="flex items-center justify-end w-full gap-3">
            <Button
              variant="outline"
              size="sm"
              onClick={() => handleOpenChange(false)}
              disabled={closeCommand.isPending}
            >
              Cancelar
            </Button>
            <Button
              size="sm"
              onClick={handleClose}
              disabled={!validation.valid || closeCommand.isPending}
              className="min-w-[100px]"
            >
              {closeCommand.isPending ? (
                <span className="flex items-center gap-2">
                  <span className="h-3 w-3 animate-spin rounded-full border-2 border-current border-t-transparent" />
                  Fechando...
                </span>
              ) : (
                'Fechar Conta'
              )}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
