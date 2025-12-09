'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal de Detalhes da Assinatura
 *
 * @component SubscriptionModal
 * @description Modal para visualizar detalhes e renovar assinaturas
 * Conforme FLUXO_ASSINATURA.md - FE-008, FE-009
 */

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import type { RenewSubscriptionRequest, SubscriptionModalState } from '@/types/subscription';
import {
    PAYMENT_METHOD_LABELS,
    SUBSCRIPTION_STATUS_COLORS,
    SUBSCRIPTION_STATUS_LABELS,
} from '@/types/subscription';
import {
    CalendarIcon,
    CreditCardIcon,
    LoaderIcon,
    UserIcon,
} from 'lucide-react';
import { useCallback, useState } from 'react';

// =============================================================================
// PROPS
// =============================================================================

interface SubscriptionModalProps {
  state: SubscriptionModalState;
  onClose: () => void;
  onConfirmRenew: (data?: RenewSubscriptionRequest) => Promise<void>;
  isLoading?: boolean;
}

// =============================================================================
// HELPERS
// =============================================================================

function formatDate(dateString?: string): string {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}

function formatCurrency(value: string | number): string {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(numValue);
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function SubscriptionModal({
  state,
  onClose,
  onConfirmRenew,
  isLoading = false,
}: SubscriptionModalProps) {
  const { isOpen, mode, subscription } = state;
  const isRenewMode = mode === 'renew';

  // Form state para renovação
  const [codigoTransacao, setCodigoTransacao] = useState('');
  const [observacao, setObservacao] = useState('');

  const handleOpenChange = useCallback((open: boolean) => {
    if (open) {
      setCodigoTransacao('');
      setObservacao('');
      return;
    }
    onClose();
  }, [onClose]);

  // Handlers
  const handleConfirmRenew = useCallback(async () => {
    const data: RenewSubscriptionRequest = {};
    if (codigoTransacao.trim()) {
      data.codigo_transacao = codigoTransacao.trim();
    }
    if (observacao.trim()) {
      data.observacao = observacao.trim();
    }
    await onConfirmRenew(Object.keys(data).length > 0 ? data : undefined);
  }, [codigoTransacao, observacao, onConfirmRenew]);

  // Título do modal
  const getTitle = () => {
    if (isRenewMode) return 'Renovar Assinatura';
    return 'Detalhes da Assinatura';
  };

  if (!subscription) return null;

  return (
    <Dialog open={isOpen} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[550px]">
        <DialogHeader>
          <DialogTitle>{getTitle()}</DialogTitle>
          <DialogDescription>
            {isRenewMode
              ? 'Registre o pagamento manual para renovar a assinatura.'
              : 'Informações completas da assinatura.'}
          </DialogDescription>
        </DialogHeader>

        <Separator />

        {/* Informações do Cliente */}
        <div className="space-y-4 py-4">
          <div className="flex items-start gap-4 p-4 rounded-lg bg-muted/50">
            <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
              <UserIcon className="h-5 w-5 text-primary" />
            </div>
            <div className="flex-1">
              <h4 className="font-medium">{subscription.cliente_nome}</h4>
              {subscription.cliente_telefone && (
                <p className="text-sm text-muted-foreground">
                  {subscription.cliente_telefone}
                </p>
              )}
              {subscription.cliente_email && (
                <p className="text-sm text-muted-foreground">
                  {subscription.cliente_email}
                </p>
              )}
            </div>
            <Badge
              className={cn(
                'hover:bg-current/90',
                SUBSCRIPTION_STATUS_COLORS[subscription.status]
              )}
            >
              {SUBSCRIPTION_STATUS_LABELS[subscription.status]}
            </Badge>
          </div>

          {/* Informações do Plano */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <Label className="text-muted-foreground">Plano</Label>
              <p className="font-medium">{subscription.plano_nome || 'Plano'}</p>
            </div>
            <div className="space-y-1">
              <Label className="text-muted-foreground">Valor</Label>
              <p className="font-medium">{formatCurrency(subscription.valor)}</p>
            </div>
          </div>

          {/* Forma de Pagamento */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <Label className="text-muted-foreground">Forma de Pagamento</Label>
              <div className="flex items-center gap-2">
                <CreditCardIcon className="h-4 w-4 text-muted-foreground" />
                <span className="font-medium">
                  {PAYMENT_METHOD_LABELS[subscription.forma_pagamento]}
                </span>
              </div>
            </div>
            <div className="space-y-1">
              <Label className="text-muted-foreground">Serviços Utilizados</Label>
              <p className="font-medium">
                {subscription.servicos_utilizados}
                {subscription.plano_limite_uso_mensal && (
                  <span className="text-muted-foreground">
                    {' '}/ {subscription.plano_limite_uso_mensal}
                  </span>
                )}
              </p>
            </div>
          </div>

          {/* Datas */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1">
              <Label className="text-muted-foreground">Data de Ativação</Label>
              <div className="flex items-center gap-2">
                <CalendarIcon className="h-4 w-4 text-muted-foreground" />
                <span className="font-medium">
                  {formatDate(subscription.data_ativacao)}
                </span>
              </div>
            </div>
            <div className="space-y-1">
              <Label className="text-muted-foreground">Data de Vencimento</Label>
              <div className="flex items-center gap-2">
                <CalendarIcon className="h-4 w-4 text-muted-foreground" />
                <span className="font-medium">
                  {formatDate(subscription.data_vencimento)}
                </span>
              </div>
            </div>
          </div>

          {/* Código de Transação */}
          {subscription.codigo_transacao && (
            <div className="space-y-1">
              <Label className="text-muted-foreground">Código de Transação</Label>
              <p className="font-medium font-mono text-sm">
                {subscription.codigo_transacao}
              </p>
            </div>
          )}

          {/* Formulário de Renovação */}
          {isRenewMode && (
            <>
              <Separator />
              <div className="space-y-4">
                <h4 className="font-medium">Dados da Renovação</h4>
                
                <div className="space-y-2">
                  <Label htmlFor="codigo_transacao">
                    Código de Transação (opcional)
                  </Label>
                  <Input
                    id="codigo_transacao"
                    placeholder="Ex: PIX123456789"
                    value={codigoTransacao}
                    onChange={(e) => setCodigoTransacao(e.target.value)}
                  />
                  <p className="text-xs text-muted-foreground">
                    Identificador do pagamento (PIX, comprovante, etc.)
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="observacao">Observação (opcional)</Label>
                  <Textarea
                    id="observacao"
                    placeholder="Informações adicionais sobre o pagamento..."
                    value={observacao}
                    onChange={(e) => setObservacao(e.target.value)}
                    rows={3}
                  />
                </div>
              </div>
            </>
          )}
        </div>

        <Separator />

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            {isRenewMode ? 'Cancelar' : 'Fechar'}
          </Button>
          {isRenewMode && (
            <Button onClick={handleConfirmRenew} disabled={isLoading}>
              {isLoading && <LoaderIcon className="mr-2 h-4 w-4 animate-spin" />}
              Confirmar Renovação
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
