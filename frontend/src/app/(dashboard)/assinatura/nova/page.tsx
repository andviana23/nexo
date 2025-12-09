'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Wizard de Nova Assinatura
 *
 * @page /assinatura/nova
 * @description Fluxo em etapas para criar nova assinatura
 * Conforme FLUXO_ASSINATURA.md - FE-007
 */

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { useCustomers } from '@/hooks/use-customers';
import { useActivePlans, useCreateSubscription } from '@/hooks/use-subscriptions';
import { cn } from '@/lib/utils';
import { useBreadcrumbs } from '@/store/ui-store';
import type { CustomerResponse } from '@/types/customer';
import type { CreateSubscriptionRequest, PaymentMethod, Plan } from '@/types/subscription';
import { PAYMENT_METHOD_LABELS, PERIODICITY_LABELS } from '@/types/subscription';
import {
    ArrowLeftIcon,
    ArrowRightIcon,
    CheckCircleIcon,
    CreditCardIcon,
    LoaderIcon,
    SearchIcon,
    UserIcon,
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useState } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

type Step = 'cliente' | 'plano' | 'pagamento' | 'confirmacao';

interface WizardState {
  step: Step;
  cliente?: CustomerResponse;
  plano?: Plan;
  formaPagamento?: PaymentMethod;
}

// =============================================================================
// HELPERS
// =============================================================================

function getInitials(name: string): string {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase();
}

function formatCurrency(value: string | number): string {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(numValue);
}

// =============================================================================
// COMPONENTES
// =============================================================================

function StepIndicator({
  steps,
  currentStep,
}: {
  steps: { key: Step; label: string }[];
  currentStep: Step;
}) {
  const currentIndex = steps.findIndex((s) => s.key === currentStep);

  return (
    <div className="flex items-center justify-center gap-2">
      {steps.map((step, index) => (
        <div key={step.key} className="flex items-center">
          <div
            className={cn(
              'flex items-center justify-center w-8 h-8 rounded-full text-sm font-medium transition-colors',
              index <= currentIndex
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted text-muted-foreground'
            )}
          >
            {index < currentIndex ? (
              <CheckCircleIcon className="h-4 w-4" />
            ) : (
              index + 1
            )}
          </div>
          {index < steps.length - 1 && (
            <div
              className={cn(
                'w-12 h-0.5 mx-2 transition-colors',
                index < currentIndex ? 'bg-primary' : 'bg-muted'
              )}
            />
          )}
        </div>
      ))}
    </div>
  );
}

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export default function NovaAssinaturaPage() {
  const router = useRouter();
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estado do wizard
  const [state, setState] = useState<WizardState>({ step: 'cliente' });
  const [searchTerm, setSearchTerm] = useState('');

  // Queries
  const { data: customers, isLoading: loadingCustomers } = useCustomers({
    page: 1,
    page_size: 50,
    search: searchTerm,
    ativo: true,
  });
  const { data: plans, isLoading: loadingPlans } = useActivePlans();

  // Mutation
  const createSubscription = useCreateSubscription();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Assinaturas', href: '/assinatura' },
      { label: 'Nova Assinatura' },
    ]);
  }, [setBreadcrumbs]);

  // Definição das etapas
  const steps: { key: Step; label: string }[] = [
    { key: 'cliente', label: 'Cliente' },
    { key: 'plano', label: 'Plano' },
    { key: 'pagamento', label: 'Pagamento' },
    { key: 'confirmacao', label: 'Confirmação' },
  ];

  // Handlers
  const handleSelectCustomer = useCallback((customer: CustomerResponse) => {
    setState((prev) => ({ ...prev, cliente: customer }));
  }, []);

  const handleSelectPlan = useCallback((plan: Plan) => {
    setState((prev) => ({ ...prev, plano: plan }));
  }, []);

  const handleSelectPayment = useCallback((method: PaymentMethod) => {
    setState((prev) => ({ ...prev, formaPagamento: method }));
  }, []);

  const handleNext = useCallback(() => {
    const stepOrder: Step[] = ['cliente', 'plano', 'pagamento', 'confirmacao'];
    const currentIndex = stepOrder.indexOf(state.step);
    if (currentIndex < stepOrder.length - 1) {
      setState((prev) => ({ ...prev, step: stepOrder[currentIndex + 1] }));
    }
  }, [state.step]);

  const handleBack = useCallback(() => {
    const stepOrder: Step[] = ['cliente', 'plano', 'pagamento', 'confirmacao'];
    const currentIndex = stepOrder.indexOf(state.step);
    if (currentIndex > 0) {
      setState((prev) => ({ ...prev, step: stepOrder[currentIndex - 1] }));
    }
  }, [state.step]);

  const handleSubmit = useCallback(async () => {
    if (!state.cliente || !state.plano || !state.formaPagamento) return;

    const data: CreateSubscriptionRequest = {
      cliente_id: state.cliente.id,
      plano_id: state.plano.id,
      forma_pagamento: state.formaPagamento,
    };

    await createSubscription.mutateAsync(data);
    router.push('/assinatura/assinantes');
  }, [state, createSubscription, router]);

  const canProceed = useCallback(() => {
    switch (state.step) {
      case 'cliente':
        return !!state.cliente;
      case 'plano':
        return !!state.plano;
      case 'pagamento':
        return !!state.formaPagamento;
      case 'confirmacao':
        return true;
      default:
        return false;
    }
  }, [state]);

  // Renderização das etapas
  const renderStep = () => {
    switch (state.step) {
      case 'cliente':
        return (
          <div className="space-y-4">
            <div className="relative">
              <SearchIcon className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Buscar cliente por nome ou telefone..."
                className="pl-9"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>

            {loadingCustomers ? (
              <div className="space-y-2">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-16 w-full" />
                ))}
              </div>
            ) : !customers?.data || customers.data.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                Nenhum cliente encontrado
              </div>
            ) : (
              <div className="space-y-2 max-h-[400px] overflow-y-auto">
                {customers.data.map((customer) => (
                  <div
                    key={customer.id}
                    onClick={() => handleSelectCustomer(customer)}
                    className={cn(
                      'flex items-center gap-3 p-4 rounded-lg border cursor-pointer transition-colors',
                      state.cliente?.id === customer.id
                        ? 'border-primary bg-primary/5'
                        : 'hover:bg-muted/50'
                    )}
                  >
                    <Avatar className="h-10 w-10">
                      <AvatarFallback>
                        {getInitials(customer.nome)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <p className="font-medium">{customer.nome}</p>
                      <p className="text-sm text-muted-foreground">
                        {customer.telefone}
                      </p>
                    </div>
                    {state.cliente?.id === customer.id && (
                      <CheckCircleIcon className="h-5 w-5 text-primary" />
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        );

      case 'plano':
        return (
          <div className="space-y-4">
            {loadingPlans ? (
              <div className="space-y-2">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-24 w-full" />
                ))}
              </div>
            ) : !plans || plans.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                Nenhum plano ativo disponível
              </div>
            ) : (
              <div className="grid gap-4 md:grid-cols-2">
                {plans.map((plan) => (
                  <Card
                    key={plan.id}
                    onClick={() => handleSelectPlan(plan)}
                    className={cn(
                      'cursor-pointer transition-all',
                      state.plano?.id === plan.id
                        ? 'border-primary ring-2 ring-primary/20'
                        : 'hover:border-muted-foreground/50'
                    )}
                  >
                    <CardHeader className="pb-2">
                      <div className="flex items-center justify-between">
                        <CardTitle className="text-lg">{plan.nome}</CardTitle>
                        {state.plano?.id === plan.id && (
                          <CheckCircleIcon className="h-5 w-5 text-primary" />
                        )}
                      </div>
                      <CardDescription>{plan.descricao}</CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="flex items-baseline gap-1">
                        <span className="text-2xl font-bold">
                          {formatCurrency(plan.valor)}
                        </span>
                        <span className="text-muted-foreground">
                          /{PERIODICITY_LABELS[plan.periodicidade].toLowerCase()}
                        </span>
                      </div>
                      {plan.limite_uso_mensal && (
                        <p className="text-sm text-muted-foreground mt-2">
                          Até {plan.limite_uso_mensal} serviços/mês
                        </p>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>
        );

      case 'pagamento':
        return (
          <div className="space-y-4">
            <RadioGroup
              value={state.formaPagamento}
              onValueChange={(value) =>
                handleSelectPayment(value as PaymentMethod)
              }
              className="space-y-3"
            >
              {(Object.entries(PAYMENT_METHOD_LABELS) as [PaymentMethod, string][]).map(
                ([value, label]) => (
                  <Label
                    key={value}
                    htmlFor={value}
                    className={cn(
                      'flex items-center gap-4 p-4 rounded-lg border cursor-pointer transition-colors',
                      state.formaPagamento === value
                        ? 'border-primary bg-primary/5'
                        : 'hover:bg-muted/50'
                    )}
                  >
                    <RadioGroupItem value={value} id={value} />
                    <CreditCardIcon className="h-5 w-5 text-muted-foreground" />
                    <span className="font-medium">{label}</span>
                  </Label>
                )
              )}
            </RadioGroup>

            {state.formaPagamento === 'CARTAO' && (
              <div className="p-4 rounded-lg bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-900">
                <p className="text-sm text-blue-700 dark:text-blue-400">
                  <strong>Pagamento com Cartão:</strong> O cliente receberá um
                  link de pagamento do Asaas por e-mail/SMS para configurar a
                  cobrança recorrente.
                </p>
              </div>
            )}

            {(state.formaPagamento === 'PIX' ||
              state.formaPagamento === 'DINHEIRO') && (
              <div className="p-4 rounded-lg bg-yellow-50 dark:bg-yellow-950/20 border border-yellow-200 dark:border-yellow-900">
                <p className="text-sm text-yellow-700 dark:text-yellow-400">
                  <strong>Pagamento Manual:</strong> A renovação deverá ser
                  registrada manualmente quando o cliente efetuar o pagamento.
                </p>
              </div>
            )}
          </div>
        );

      case 'confirmacao':
        return (
          <div className="space-y-6">
            {/* Resumo do Cliente */}
            <div className="flex items-center gap-4 p-4 rounded-lg bg-muted/50">
              <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center">
                <UserIcon className="h-6 w-6 text-primary" />
              </div>
              <div>
                <p className="font-medium">{state.cliente?.nome}</p>
                <p className="text-sm text-muted-foreground">
                  {state.cliente?.telefone}
                </p>
              </div>
            </div>

            {/* Resumo do Plano */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label className="text-muted-foreground">Plano</Label>
                <p className="font-medium">{state.plano?.nome}</p>
              </div>
              <div>
                <Label className="text-muted-foreground">Valor</Label>
                <p className="font-medium">
                  {state.plano && formatCurrency(state.plano.valor)}
                </p>
              </div>
              <div>
                <Label className="text-muted-foreground">Periodicidade</Label>
                <p className="font-medium">
                  {state.plano &&
                    PERIODICITY_LABELS[state.plano.periodicidade]}
                </p>
              </div>
              <div>
                <Label className="text-muted-foreground">
                  Forma de Pagamento
                </Label>
                <p className="font-medium">
                  {state.formaPagamento &&
                    PAYMENT_METHOD_LABELS[state.formaPagamento]}
                </p>
              </div>
            </div>

            <Separator />

            <div className="p-4 rounded-lg bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-900">
              <p className="text-sm text-green-700 dark:text-green-400">
                Ao confirmar, a assinatura será criada e o cliente será
                notificado. Para cartão de crédito, o link de pagamento será
                enviado automaticamente.
              </p>
            </div>
          </div>
        );
    }
  };

  return (
    <div className="flex flex-col gap-6 p-6 max-w-3xl mx-auto">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold tracking-tight">Nova Assinatura</h1>
        <p className="text-muted-foreground mt-2">
          Siga os passos para criar uma nova assinatura
        </p>
      </div>

      {/* Step Indicator */}
      <StepIndicator steps={steps} currentStep={state.step} />

      {/* Card Principal */}
      <Card>
        <CardHeader>
          <CardTitle>
            {steps.find((s) => s.key === state.step)?.label}
          </CardTitle>
          <CardDescription>
            {state.step === 'cliente' &&
              'Selecione o cliente que receberá a assinatura'}
            {state.step === 'plano' && 'Escolha o plano de assinatura'}
            {state.step === 'pagamento' &&
              'Selecione a forma de pagamento preferida'}
            {state.step === 'confirmacao' &&
              'Revise os dados e confirme a assinatura'}
          </CardDescription>
        </CardHeader>
        <Separator />
        <CardContent className="pt-6">{renderStep()}</CardContent>
        <Separator />
        <div className="flex justify-between p-6">
          <Button
            variant="outline"
            onClick={handleBack}
            disabled={state.step === 'cliente'}
          >
            <ArrowLeftIcon className="mr-2 h-4 w-4" />
            Voltar
          </Button>

          {state.step === 'confirmacao' ? (
            <Button
              onClick={handleSubmit}
              disabled={!canProceed() || createSubscription.isPending}
            >
              {createSubscription.isPending && (
                <LoaderIcon className="mr-2 h-4 w-4 animate-spin" />
              )}
              Criar Assinatura
            </Button>
          ) : (
            <Button onClick={handleNext} disabled={!canProceed()}>
              Próximo
              <ArrowRightIcon className="ml-2 h-4 w-4" />
            </Button>
          )}
        </div>
      </Card>
    </div>
  );
}
