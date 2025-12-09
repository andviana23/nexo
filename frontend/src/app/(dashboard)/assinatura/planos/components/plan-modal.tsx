'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal de Plano de Assinatura
 *
 * @component PlanModal
 * @description Modal para criar, editar e visualizar planos
 * Conforme FLUXO_ASSINATURA.md - FE-004
 */

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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';
import type {
    CreatePlanRequest,
    PlanModalState,
    PlanPeriodicity,
    UpdatePlanRequest,
} from '@/types/subscription';
import { PERIODICITY_LABELS } from '@/types/subscription';
import { LoaderIcon } from 'lucide-react';
import { useCallback, useState } from 'react';

// =============================================================================
// PROPS
// =============================================================================

interface PlanModalProps {
  state: PlanModalState;
  onClose: () => void;
  onSave: (data: CreatePlanRequest | UpdatePlanRequest) => Promise<void>;
  isLoading?: boolean;
}

// =============================================================================
// HELPERS
// =============================================================================

interface FormData {
  nome: string;
  descricao: string;
  valor: string;
  periodicidade: PlanPeriodicity;
}

function getInitialFormData(plan?: { id: string; nome: string; descricao?: string; valor: string; periodicidade: PlanPeriodicity }): FormData {
  if (!plan) {
    return {
      nome: '',
      descricao: '',
      valor: '',
      periodicidade: 'MENSAL',
    };
  }
  return {
    nome: plan.nome,
    descricao: plan.descricao || '',
    valor: plan.valor,
    periodicidade: plan.periodicidade,
  };
}

function formatCurrency(value: string): string {
  // Remove tudo que não é número
  const numbers = value.replace(/\D/g, '');
  if (!numbers) return '';

  // Converte para centavos e formata
  const cents = parseInt(numbers, 10);
  const reais = cents / 100;

  return reais.toLocaleString('pt-BR', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

function parseCurrency(value: string): string {
  // Remove formatação e converte para decimal
  const numbers = value.replace(/\D/g, '');
  if (!numbers) return '0';

  const cents = parseInt(numbers, 10);
  return (cents / 100).toFixed(2);
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function PlanModal({
  state,
  onClose,
  onSave,
  isLoading = false,
}: PlanModalProps) {
  const { isOpen, mode, plan } = state;
  const isView = mode === 'view';

  // Form state
  const [formData, setFormData] = useState<FormData>(getInitialFormData(plan));
  const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>({});

  const handleOpenChange = useCallback((open: boolean) => {
    if (open) {
      setFormData(getInitialFormData(plan));
      setErrors({});
      return;
    }
    onClose();
  }, [onClose, plan]);

  // Handlers
  const handleChange = useCallback(
    (field: keyof FormData, value: string) => {
      setFormData((prev) => ({ ...prev, [field]: value }));
      // Limpar erro ao editar
      if (errors[field]) {
        setErrors((prev) => ({ ...prev, [field]: undefined }));
      }
    },
    [errors]
  );

  const handlePriceChange = useCallback((value: string) => {
    const formatted = formatCurrency(value);
    setFormData((prev) => ({ ...prev, valor: formatted }));
    setErrors((prev) => ({ ...prev, valor: undefined }));
  }, []);

  const validate = useCallback((): boolean => {
    const newErrors: Partial<Record<keyof FormData, string>> = {};

    if (!formData.nome.trim()) {
      newErrors.nome = 'Nome é obrigatório';
    } else if (formData.nome.trim().length < 3) {
      newErrors.nome = 'Nome deve ter pelo menos 3 caracteres';
    }

    if (!formData.valor) {
      newErrors.valor = 'Preço é obrigatório';
    } else {
      const valor = parseFloat(parseCurrency(formData.valor));
      if (valor <= 0) {
        newErrors.valor = 'Preço deve ser maior que zero';
      }
    }

    if (!formData.periodicidade) {
      newErrors.periodicidade = 'Periodicidade é obrigatória';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [formData]);

  const handleSubmit = useCallback(async () => {
    if (!validate()) return;

    const payload: CreatePlanRequest | UpdatePlanRequest = {
      nome: formData.nome.trim(),
      descricao: formData.descricao.trim() || undefined,
      valor: parseCurrency(formData.valor),
      periodicidade: formData.periodicidade,
    };

    await onSave(payload);
  }, [formData, validate, onSave]);

  // Título do modal
  const getTitle = () => {
    switch (mode) {
      case 'create':
        return 'Novo Plano';
      case 'edit':
        return 'Editar Plano';
      case 'view':
        return 'Detalhes do Plano';
      default:
        return 'Plano';
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{getTitle()}</DialogTitle>
          <DialogDescription>
            {mode === 'create' && 'Preencha os dados para criar um novo plano de assinatura.'}
            {mode === 'edit' && 'Edite os dados do plano. Alterações não afetam assinaturas existentes.'}
            {mode === 'view' && 'Visualize os detalhes do plano de assinatura.'}
          </DialogDescription>
        </DialogHeader>

        <Separator />

        <div className="grid gap-4 py-4">
          {/* Nome */}
          <div className="grid gap-2">
            <Label htmlFor="nome" className={errors.nome ? 'text-destructive' : ''}>
              Nome do Plano *
            </Label>
            <Input
              id="nome"
              placeholder="Ex: Plano Mensal, Plano Trimestral..."
              value={formData.nome}
              onChange={(e) => handleChange('nome', e.target.value)}
              disabled={isView}
              className={errors.nome ? 'border-destructive' : ''}
            />
            {errors.nome && (
              <span className="text-xs text-destructive">{errors.nome}</span>
            )}
          </div>

          {/* Descrição */}
          <div className="grid gap-2">
            <Label htmlFor="descricao">Descrição</Label>
            <Textarea
              id="descricao"
              placeholder="Descreva os benefícios do plano..."
              value={formData.descricao}
              onChange={(e) => handleChange('descricao', e.target.value)}
              disabled={isView}
              rows={3}
            />
          </div>

          {/* Preço e Periodicidade */}
          <div className="grid grid-cols-2 gap-4">
            <div className="grid gap-2">
              <Label htmlFor="valor" className={errors.valor ? 'text-destructive' : ''}>
                Preço (R$) *
              </Label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                  R$
                </span>
                <Input
                  id="valor"
                  placeholder="0,00"
                  value={formData.valor}
                  onChange={(e) => handlePriceChange(e.target.value)}
                  disabled={isView}
                  className={`pl-10 ${errors.valor ? 'border-destructive' : ''}`}
                />
              </div>
              {errors.valor && (
                <span className="text-xs text-destructive">{errors.valor}</span>
              )}
            </div>

            <div className="grid gap-2">
              <Label
                htmlFor="periodicidade"
                className={errors.periodicidade ? 'text-destructive' : ''}
              >
                Periodicidade *
              </Label>
              <Select
                value={formData.periodicidade}
                onValueChange={(value) =>
                  handleChange('periodicidade', value as PlanPeriodicity)
                }
                disabled={isView}
              >
                <SelectTrigger
                  id="periodicidade"
                  className={errors.periodicidade ? 'border-destructive' : ''}
                >
                  <SelectValue placeholder="Selecione..." />
                </SelectTrigger>
                <SelectContent>
                  {Object.entries(PERIODICITY_LABELS).map(([value, label]) => (
                    <SelectItem key={value} value={value}>
                      {label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.periodicidade && (
                <span className="text-xs text-destructive">
                  {errors.periodicidade}
                </span>
              )}
            </div>
          </div>
        </div>

        <Separator />

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            {isView ? 'Fechar' : 'Cancelar'}
          </Button>
          {!isView && (
            <Button onClick={handleSubmit} disabled={isLoading}>
              {isLoading && <LoaderIcon className="mr-2 h-4 w-4 animate-spin" />}
              {mode === 'create' ? 'Criar Plano' : 'Salvar Alterações'}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
