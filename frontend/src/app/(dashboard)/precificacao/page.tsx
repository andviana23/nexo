/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Precificação
 *
 * @page /precificacao
 * @description Simulador de preços e configuração de margem
 * Baseado em docs/10-calculos/preco-servico.md e markup.md
 */

'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from '@/components/ui/tooltip';
import {
    usePricingConfig,
    useSaveConfig,
    useSimulatePrice,
    useSimulations,
    useUpdateConfig,
} from '@/hooks/use-pricing';
import { useServices } from '@/hooks/useServices';
import { cn } from '@/lib/utils';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    PrecificacaoSimulacaoResponse,
    SimularPrecoRequest,
    TipoItem,
} from '@/types/pricing';
import {
    calcularMargemFinal,
    calcularPrecoSugerido,
    formatCurrency,
    formatPercentual,
    margemParaMarkup,
    percentualParaDecimal,
} from '@/types/pricing';
import type { Service } from '@/types/service';
import {
    AlertCircle,
    ArrowRight,
    Calculator,
    CheckCircle2,
    HelpCircle,
    History,
    Percent,
    PiggyBank,
    RefreshCcw,
    Save,
    Settings,
    TrendingDown,
    TrendingUp,
} from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';

// =============================================================================
// TIPOS LOCAIS
// =============================================================================

interface SimuladorForm {
  tipoItem: TipoItem;
  itemId: string;
  custoMateriais: string;
  custoMaoDeObra: string;
  precoAtual: string;
  margemDesejada: string;
  impostoPercentual: string;
  comissaoPercentual: string;
}

type ConfigForm = {
  margemDesejada: string;
  markupAlvo: string;
  impostoPercentual: string;
  comissaoDefault: string;
};

const INITIAL_FORM: SimuladorForm = {
  tipoItem: 'SERVICO',
  itemId: '',
  custoMateriais: '',
  custoMaoDeObra: '',
  precoAtual: '',
  margemDesejada: '35',
  impostoPercentual: '6',
  comissaoPercentual: '30',
};

const INITIAL_CONFIG_FORM: ConfigForm = {
  margemDesejada: '',
  markupAlvo: '',
  impostoPercentual: '',
  comissaoDefault: '',
};

// =============================================================================
// COMPONENTE: Card de Resultado
// =============================================================================

interface ResultCardProps {
  label: string;
  value: string;
  icon: React.ReactNode;
  variant?: 'default' | 'success' | 'warning' | 'danger';
  description?: string;
}

function ResultCard({ label, value, icon, variant = 'default', description }: ResultCardProps) {
  const variantClasses = {
    default: 'text-foreground',
    success: 'text-green-600',
    warning: 'text-yellow-600',
    danger: 'text-red-600',
  };

  return (
    <div className="flex items-center gap-3 p-4 rounded-lg border bg-card">
      <div className={cn('p-2 rounded-lg bg-muted', variantClasses[variant])}>
        {icon}
      </div>
      <div className="flex-1">
        <p className="text-sm text-muted-foreground">{label}</p>
        <p className={cn('text-xl font-bold', variantClasses[variant])}>{value}</p>
        {description && (
          <p className="text-xs text-muted-foreground mt-0.5">{description}</p>
        )}
      </div>
    </div>
  );
}

// =============================================================================
// COMPONENTE: Configuração de Precificação
// =============================================================================

function ConfiguracaoSection() {
  const { data: config, isLoading } = usePricingConfig();
  const saveConfig = useSaveConfig();
  const updateConfig = useUpdateConfig();

  const baseForm = useMemo<ConfigForm>(() => ({
    margemDesejada: config?.margem_desejada ?? INITIAL_CONFIG_FORM.margemDesejada,
    markupAlvo: config?.markup_alvo ?? INITIAL_CONFIG_FORM.markupAlvo,
    impostoPercentual: config?.imposto_percentual ?? INITIAL_CONFIG_FORM.impostoPercentual,
    comissaoDefault: config?.comissao_percentual_default ?? INITIAL_CONFIG_FORM.comissaoDefault,
  }), [config]);

  const [overrides, setOverrides] = useState<Partial<ConfigForm>>({});

  const form = useMemo(
    () => ({ ...baseForm, ...overrides }),
    [baseForm, overrides]
  );

  const handleSave = useCallback(async () => {
    const data = {
      margem_desejada: form.margemDesejada,
      markup_alvo: form.markupAlvo,
      imposto_percentual: form.impostoPercentual,
      comissao_default: form.comissaoDefault,
    };

    if (config) {
      await updateConfig.mutateAsync(data);
    } else {
      await saveConfig.mutateAsync(data);
    }
  }, [form, config, saveConfig, updateConfig]);

  // Sincronizar margem ↔ markup
  const handleMargemChange = useCallback((value: string) => {
    setOverrides((prev) => {
      const margem = parseFloat(value) / 100;
      const markup = margemParaMarkup(margem);
      return {
        ...prev,
        margemDesejada: value,
        markupAlvo: isFinite(markup) ? markup.toFixed(2) : '',
      };
    });
  }, []);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-48" />
          <Skeleton className="h-4 w-64" />
        </CardHeader>
        <CardContent className="space-y-4">
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-10 w-full" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings className="h-5 w-5" />
          Configuração Padrão
        </CardTitle>
        <CardDescription>
          Defina os parâmetros padrão para cálculo de preços
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-4 sm:grid-cols-2">
          {/* Margem Desejada */}
          <div className="space-y-2">
            <Label htmlFor="margem" className="flex items-center gap-1">
              Margem de Lucro Desejada
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <HelpCircle className="h-3.5 w-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent className="max-w-xs">
                    <p>Percentual de lucro líquido sobre o preço de venda.</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      Ex: 35% significa que 35% do preço é lucro.
                    </p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </Label>
            <div className="relative">
              <Input
                id="margem"
                type="number"
                step="0.01"
                min="0"
                max="99"
                placeholder="35.00"
                value={form.margemDesejada}
                onChange={(e) => handleMargemChange(e.target.value)}
                className="pr-8"
              />
              <span className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                %
              </span>
            </div>
          </div>

          {/* Markup */}
          <div className="space-y-2">
            <Label htmlFor="markup" className="flex items-center gap-1">
              Markup Alvo
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <HelpCircle className="h-3.5 w-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent className="max-w-xs">
                    <p>Fator multiplicador sobre o custo.</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      Ex: Markup 2.0 = Preço é 2x o custo.
                    </p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </Label>
            <div className="relative">
              <Input
                id="markup"
                type="number"
                step="0.01"
                min="1"
                placeholder="2.00"
                value={form.markupAlvo}
                onChange={(e) => setOverrides((prev) => ({ ...prev, markupAlvo: e.target.value }))}
                className="pr-8"
              />
              <span className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                x
              </span>
            </div>
          </div>

          {/* Imposto */}
          <div className="space-y-2">
            <Label htmlFor="imposto" className="flex items-center gap-1">
              Alíquota de Impostos
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <HelpCircle className="h-3.5 w-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent className="max-w-xs">
                    <p>Percentual de impostos sobre a receita (Simples, ISS, etc.)</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </Label>
            <div className="relative">
              <Input
                id="imposto"
                type="number"
                step="0.01"
                min="0"
                max="100"
                placeholder="6.00"
                value={form.impostoPercentual}
                onChange={(e) => setOverrides((prev) => ({ ...prev, impostoPercentual: e.target.value }))}
                className="pr-8"
              />
              <span className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                %
              </span>
            </div>
          </div>

          {/* Comissão */}
          <div className="space-y-2">
            <Label htmlFor="comissao" className="flex items-center gap-1">
              Comissão Padrão
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <HelpCircle className="h-3.5 w-3.5 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent className="max-w-xs">
                    <p>Percentual de comissão pago ao profissional sobre o preço.</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </Label>
            <div className="relative">
              <Input
                id="comissao"
                type="number"
                step="0.01"
                min="0"
                max="100"
                placeholder="30.00"
                value={form.comissaoDefault}
                onChange={(e) => setOverrides((prev) => ({ ...prev, comissaoDefault: e.target.value }))}
                className="pr-8"
              />
              <span className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                %
              </span>
            </div>
          </div>
        </div>

        <Button
          onClick={handleSave}
          disabled={saveConfig.isPending || updateConfig.isPending}
          className="w-full sm:w-auto"
        >
          <Save className="h-4 w-4 mr-2" />
          {saveConfig.isPending || updateConfig.isPending ? 'Salvando...' : 'Salvar Configuração'}
        </Button>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// COMPONENTE: Simulador de Preços
// =============================================================================

function SimuladorSection() {
  const { data: config } = usePricingConfig();
  const { data: servicesData } = useServices({ apenas_ativos: true });
  const services = useMemo(() => servicesData?.servicos ?? [], [servicesData]);
  const simulatePrice = useSimulatePrice();

  const baseForm = useMemo<SimuladorForm>(() => ({
    ...INITIAL_FORM,
    margemDesejada: config?.margem_desejada ?? INITIAL_FORM.margemDesejada,
    impostoPercentual: config?.imposto_percentual ?? INITIAL_FORM.impostoPercentual,
    comissaoPercentual: config?.comissao_percentual_default ?? INITIAL_FORM.comissaoPercentual,
  }), [config]);

  const [overrides, setOverrides] = useState<Partial<SimuladorForm>>({});
  const form = useMemo(() => ({ ...baseForm, ...overrides }), [baseForm, overrides]);
  const [resultado, setResultado] = useState<PrecificacaoSimulacaoResponse | null>(null);

  // Cálculo local em tempo real
  const calculoLocal = useMemo(() => {
    const custoMateriais = parseFloat(form.custoMateriais) || 0;
    const custoMaoDeObra = parseFloat(form.custoMaoDeObra) || 0;
    const precoAtual = parseFloat(form.precoAtual) || 0;
    const margem = percentualParaDecimal(form.margemDesejada || '0');
    const imposto = percentualParaDecimal(form.impostoPercentual || '0');
    const comissao = percentualParaDecimal(form.comissaoPercentual || '0');

    const custoTotal = custoMateriais + custoMaoDeObra;
    const precoSugerido = calcularPrecoSugerido(custoTotal, margem, imposto, comissao);
    const margemFinal = precoAtual > 0 
      ? calcularMargemFinal(precoAtual, custoTotal, imposto, comissao)
      : 0;
    const diferencaPercentual = precoAtual > 0 
      ? ((precoSugerido - precoAtual) / precoAtual) * 100
      : 0;
    const lucroEstimado = precoSugerido > 0 
      ? precoSugerido * margem
      : 0;

    return {
      custoTotal,
      precoSugerido,
      margemFinal,
      diferencaPercentual,
      lucroEstimado,
    };
  }, [form]);

  const handleSimulate = useCallback(async () => {
    if (!form.itemId) {
      toast.error('Selecione um serviço ou produto');
      return;
    }
    if (!form.custoMateriais && !form.custoMaoDeObra) {
      toast.error('Informe ao menos um custo');
      return;
    }
    if (!form.precoAtual) {
      toast.error('Informe o preço atual');
      return;
    }

    const request: SimularPrecoRequest = {
      item_id: form.itemId,
      tipo_item: form.tipoItem,
      custo_materiais: form.custoMateriais || '0',
      custo_mao_de_obra: form.custoMaoDeObra || '0',
      preco_atual: form.precoAtual,
      parametros: {
        margem_desejada: form.margemDesejada,
        imposto_percentual: form.impostoPercentual,
        comissao_percentual: form.comissaoPercentual,
      },
    };

    try {
      const result = await simulatePrice.mutateAsync(request);
      setResultado(result);
      toast.success('Simulação realizada!');
    } catch {
      // Erro já tratado no hook
    }
  }, [form, simulatePrice]);

  const handleReset = useCallback(() => {
    setOverrides({});
    setResultado(null);
  }, []);

  // Quando selecionar um serviço, preencher preço atual
  const handleServiceSelect = useCallback((serviceId: string) => {
    setOverrides((prev) => ({ ...prev, itemId: serviceId }));
    const service = services.find((s: Service) => s.id === serviceId);
    if (service) {
      setOverrides((prev) => ({
        ...prev,
        precoAtual: service.preco,
      }));
    }
  }, [services]);

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      {/* Formulário */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calculator className="h-5 w-5" />
            Simulador de Preços
          </CardTitle>
          <CardDescription>
            Calcule o preço ideal baseado nos custos e margem desejada
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Tipo e Item */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label>Tipo</Label>
              <Select
                value={form.tipoItem}
                onValueChange={(v: TipoItem) => setOverrides((prev) => ({ ...prev, tipoItem: v, itemId: '' }))}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="SERVICO">Serviço</SelectItem>
                  <SelectItem value="PRODUTO">Produto</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>Serviço / Produto</Label>
              <Select
                value={form.itemId}
                onValueChange={handleServiceSelect}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione..." />
                </SelectTrigger>
                <SelectContent>
                  {form.tipoItem === 'SERVICO' && services.map((service: Service) => (
                    <SelectItem key={service.id} value={service.id}>
                      {service.nome}
                    </SelectItem>
                  ))}
                  {form.tipoItem === 'PRODUTO' && (
                    <SelectItem value="custom">Produto personalizado</SelectItem>
                  )}
                </SelectContent>
              </Select>
            </div>
          </div>

          <Separator />

          {/* Custos */}
          <div className="space-y-4">
            <h4 className="font-medium text-sm">Custos</h4>
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="custoMateriais">Custo de Materiais/Insumos</Label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                    R$
                  </span>
                  <Input
                    id="custoMateriais"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0,00"
                    value={form.custoMateriais}
                      onChange={(e) => setOverrides((prev) => ({ ...prev, custoMateriais: e.target.value }))}
                    className="pl-10"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="custoMaoDeObra">Custo de Mão de Obra</Label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                    R$
                  </span>
                  <Input
                    id="custoMaoDeObra"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="0,00"
                    value={form.custoMaoDeObra}
                      onChange={(e) => setOverrides((prev) => ({ ...prev, custoMaoDeObra: e.target.value }))}
                    className="pl-10"
                  />
                </div>
              </div>
            </div>
          </div>

          <Separator />

          {/* Preço Atual */}
          <div className="space-y-2">
            <Label htmlFor="precoAtual">Preço Atual de Venda</Label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                R$
              </span>
              <Input
                id="precoAtual"
                type="number"
                step="0.01"
                min="0"
                placeholder="0,00"
                value={form.precoAtual}
                  onChange={(e) => setOverrides((prev) => ({ ...prev, precoAtual: e.target.value }))}
                className="pl-10"
              />
            </div>
          </div>

          <Separator />

          {/* Parâmetros */}
          <div className="space-y-4">
            <h4 className="font-medium text-sm">Parâmetros de Cálculo</h4>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-2">
                <Label htmlFor="margemSim">Margem (%)</Label>
                <Input
                  id="margemSim"
                  type="number"
                  step="0.01"
                  min="0"
                  max="99"
                  value={form.margemDesejada}
                    onChange={(e) => setOverrides((prev) => ({ ...prev, margemDesejada: e.target.value }))}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="impostoSim">Imposto (%)</Label>
                <Input
                  id="impostoSim"
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  value={form.impostoPercentual}
                    onChange={(e) => setOverrides((prev) => ({ ...prev, impostoPercentual: e.target.value }))}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="comissaoSim">Comissão (%)</Label>
                <Input
                  id="comissaoSim"
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  value={form.comissaoPercentual}
                    onChange={(e) => setOverrides((prev) => ({ ...prev, comissaoPercentual: e.target.value }))}
                />
              </div>
            </div>
          </div>

          {/* Botões */}
          <div className="flex gap-2">
            <Button
              onClick={handleSimulate}
              disabled={simulatePrice.isPending}
              className="flex-1"
            >
              <Calculator className="h-4 w-4 mr-2" />
              {simulatePrice.isPending ? 'Calculando...' : 'Calcular Preço'}
            </Button>
            <Button variant="outline" onClick={handleReset}>
              <RefreshCcw className="h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Resultado */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <PiggyBank className="h-5 w-5" />
            Resultado da Simulação
          </CardTitle>
          <CardDescription>
            {resultado 
              ? 'Resultado salvo no servidor'
              : 'Cálculo em tempo real (prévia)'}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Cards de Resultado */}
          <div className="grid gap-4 sm:grid-cols-2">
            <ResultCard
              label="Custo Total"
              value={formatCurrency(resultado?.custo_total || calculoLocal.custoTotal)}
              icon={<Percent className="h-4 w-4" />}
            />
            <ResultCard
              label="Preço Sugerido"
              value={formatCurrency(resultado?.preco_sugerido || calculoLocal.precoSugerido)}
              icon={<Calculator className="h-4 w-4" />}
              variant="success"
            />
          </div>

          {/* Comparação */}
          {(parseFloat(form.precoAtual) > 0 || resultado) && (
            <>
              <Separator />
              <div className="p-4 rounded-lg bg-muted/50">
                <div className="flex items-center justify-between mb-3">
                  <span className="text-sm font-medium">Comparação com Preço Atual</span>
                  {Number(resultado?.diferenca_percentual || calculoLocal.diferencaPercentual) > 0 ? (
                    <Badge variant="destructive" className="gap-1">
                      <TrendingUp className="h-3 w-3" />
                      Subir {formatPercentual(Math.abs(resultado?.diferenca_percentual ? parseFloat(String(resultado.diferenca_percentual)) : calculoLocal.diferencaPercentual))}
                    </Badge>
                  ) : Number(resultado?.diferenca_percentual || calculoLocal.diferencaPercentual) < 0 ? (
                    <Badge variant="secondary" className="gap-1">
                      <TrendingDown className="h-3 w-3" />
                      Reduzir {formatPercentual(Math.abs(resultado?.diferenca_percentual ? parseFloat(String(resultado.diferenca_percentual)) : calculoLocal.diferencaPercentual))}
                    </Badge>
                  ) : (
                    <Badge variant="outline" className="gap-1">
                      <CheckCircle2 className="h-3 w-3" />
                      Preço OK
                    </Badge>
                  )}
                </div>

                <div className="flex items-center gap-4 text-center">
                  <div className="flex-1">
                    <p className="text-xs text-muted-foreground">Atual</p>
                    <p className="text-lg font-bold">
                      {formatCurrency(resultado?.preco_atual || form.precoAtual || '0')}
                    </p>
                  </div>
                  <ArrowRight className="h-5 w-5 text-muted-foreground" />
                  <div className="flex-1">
                    <p className="text-xs text-muted-foreground">Sugerido</p>
                    <p className="text-lg font-bold text-green-600">
                      {formatCurrency(resultado?.preco_sugerido || calculoLocal.precoSugerido)}
                    </p>
                  </div>
                </div>
              </div>
            </>
          )}

          {/* Detalhes */}
          <Separator />
          <div className="grid gap-4 sm:grid-cols-2">
            <ResultCard
              label="Lucro Estimado"
              value={formatCurrency(resultado?.lucro_estimado || calculoLocal.lucroEstimado)}
              icon={<PiggyBank className="h-4 w-4" />}
              variant="success"
              description="Por unidade vendida"
            />
            <ResultCard
              label="Margem Final"
              value={formatPercentual((resultado?.margem_final ? parseFloat(resultado.margem_final) : calculoLocal.margemFinal) * 100)}
              icon={<Percent className="h-4 w-4" />}
              variant={
                (resultado?.margem_final ? parseFloat(resultado.margem_final) : calculoLocal.margemFinal) >= percentualParaDecimal(form.margemDesejada)
                  ? 'success'
                  : 'warning'
              }
              description={
                (resultado?.margem_final ? parseFloat(resultado.margem_final) : calculoLocal.margemFinal) >= percentualParaDecimal(form.margemDesejada)
                  ? 'Dentro da meta'
                  : 'Abaixo da meta'
              }
            />
          </div>

          {/* Aviso */}
          {calculoLocal.precoSugerido <= 0 && (
            <div className="p-4 rounded-lg bg-destructive/10 text-destructive flex items-start gap-2">
              <AlertCircle className="h-5 w-5 mt-0.5" />
              <div>
                <p className="font-medium">Margem inviável</p>
                <p className="text-sm">
                  A soma de margem, impostos e comissão não pode ser ≥ 100%.
                </p>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// =============================================================================
// COMPONENTE: Histórico de Simulações
// =============================================================================

function HistoricoSection() {
  const { data: simulations, isLoading } = useSimulations();

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <Skeleton className="h-64 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (!simulations || simulations.length === 0) {
    return (
      <Card>
        <CardContent className="py-16 text-center">
          <History className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-semibold mb-2">Nenhuma simulação salva</h3>
          <p className="text-muted-foreground max-w-md mx-auto">
            As simulações realizadas aparecerão aqui para consulta futura.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <History className="h-5 w-5" />
          Histórico de Simulações
        </CardTitle>
        <CardDescription>
          Simulações salvas anteriormente
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {simulations.slice(0, 10).map((sim) => (
            <div
              key={sim.id}
              className="flex items-center justify-between p-4 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
            >
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <Badge variant="outline">{sim.tipo_item}</Badge>
                  <span className="text-sm text-muted-foreground">
                    {new Date(sim.criado_em).toLocaleDateString('pt-BR')}
                  </span>
                </div>
                <div className="flex items-center gap-4 mt-2">
                  <span className="text-sm">
                    Custo: <strong>{formatCurrency(sim.custo_total)}</strong>
                  </span>
                  <ArrowRight className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm">
                    Sugerido: <strong className="text-green-600">{formatCurrency(sim.preco_sugerido)}</strong>
                  </span>
                </div>
              </div>
              <div className="text-right">
                <p className="text-sm text-muted-foreground">Margem</p>
                <p className="font-bold">
                  {formatPercentual(parseFloat(sim.margem_final) * 100)}
                </p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL
// =============================================================================

export default function PrecificacaoPage() {
  const { setBreadcrumbs } = useBreadcrumbs();
  const [activeTab, setActiveTab] = useState<'simulador' | 'config' | 'historico'>('simulador');

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([{ label: 'Precificação' }]);
  }, [setBreadcrumbs]);

  return (
    <div className="container mx-auto py-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
          <Calculator className="h-8 w-8" />
          Precificação
        </h1>
        <p className="text-muted-foreground">
          Calcule preços ideais baseados em custos, margem e impostos
        </p>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as typeof activeTab)}>
        <TabsList>
          <TabsTrigger value="simulador" className="gap-2">
            <Calculator className="h-4 w-4" />
            Simulador
          </TabsTrigger>
          <TabsTrigger value="config" className="gap-2">
            <Settings className="h-4 w-4" />
            Configuração
          </TabsTrigger>
          <TabsTrigger value="historico" className="gap-2">
            <History className="h-4 w-4" />
            Histórico
          </TabsTrigger>
        </TabsList>

        <TabsContent value="simulador" className="mt-6">
          <SimuladorSection />
        </TabsContent>

        <TabsContent value="config" className="mt-6">
          <ConfiguracaoSection />
        </TabsContent>

        <TabsContent value="historico" className="mt-6">
          <HistoricoSection />
        </TabsContent>
      </Tabs>
    </div>
  );
}
