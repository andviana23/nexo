/**
 * NEXO - Sistema de Gestão para Barbearias
 * Types: Pricing (Precificação)
 *
 * @module types/pricing
 * @description Definições de tipos para precificação no NEXO
 * Baseado em docs/10-calculos/preco-servico.md e markup.md
 */

// ============================================================================
// ENUMS
// ============================================================================

/** Tipo de item para simulação */
export type TipoItem = 'SERVICO' | 'PRODUTO';

// ============================================================================
// ENTITIES
// ============================================================================

/** Configuração de precificação do tenant */
export interface PrecificacaoConfig {
  id: string;
  margem_desejada: string;        // Percentual (ex: "35.00" = 35%)
  markup_alvo: string;            // Fator multiplicador (ex: "2.00")
  imposto_percentual: string;     // Percentual de impostos (ex: "6.00" = 6%)
  comissao_percentual_default: string; // Comissão padrão (ex: "30.00" = 30%)
  criado_em: string;
  atualizado_em: string;
}

/** Simulação de preço salva */
export interface PrecificacaoSimulacao {
  id: string;
  item_id: string;
  tipo_item: TipoItem;
  custo_materiais: string;
  custo_mao_de_obra: string;
  custo_total: string;
  margem_desejada: string;
  imposto_percentual: string;
  comissao_percentual: string;
  preco_atual: string;
  preco_sugerido: string;
  diferenca_percentual: string;
  lucro_estimado: string;
  margem_final: string;
  criado_em: string;
}

// ============================================================================
// DTOs - Request
// ============================================================================

/** DTO para salvar/atualizar configuração */
export interface SaveConfigPrecificacaoRequest {
  margem_desejada: string;
  markup_alvo: string;
  imposto_percentual: string;
  comissao_default: string;
}

/** Parâmetros opcionais para simulação */
export interface ParametrosSimulacao {
  margem_desejada?: string;
  imposto_percentual?: string;
  comissao_percentual?: string;
}

/** DTO para simular preço */
export interface SimularPrecoRequest {
  item_id: string;
  tipo_item: TipoItem;
  custo_materiais: string;
  custo_mao_de_obra: string;
  preco_atual: string;
  parametros?: ParametrosSimulacao;
}

/** Filtros para listagem de simulações */
export interface ListSimulacoesFilters {
  item_id?: string;
  tipo_item?: TipoItem;
  page?: number;
  page_size?: number;
}

// ============================================================================
// DTOs - Response
// ============================================================================

/** Resposta da configuração de precificação */
export type PrecificacaoConfigResponse = PrecificacaoConfig;

/** Resposta de uma simulação */
export type PrecificacaoSimulacaoResponse = PrecificacaoSimulacao;

/** Resposta de listagem de simulações */
export interface ListSimulacoesResponse {
  data: PrecificacaoSimulacaoResponse[];
  page: number;
  page_size: number;
  total: number;
}

// ============================================================================
// UTILITÁRIOS
// ============================================================================

/**
 * Calcula o markup a partir da margem desejada
 * Markup = 1 / (1 - Margem)
 */
export function margemParaMarkup(margem: number): number {
  if (margem >= 1) return Infinity;
  return 1 / (1 - margem);
}

/**
 * Calcula a margem a partir do markup
 * Margem = 1 - (1 / Markup)
 */
export function markupParaMargem(markup: number): number {
  if (markup <= 0) return 0;
  return 1 - (1 / markup);
}

/**
 * Calcula o preço sugerido usando a fórmula de precificação
 * Preço = CustoTotal / (1 - Margem - Impostos - Comissão)
 */
export function calcularPrecoSugerido(
  custoTotal: number,
  margem: number,
  imposto: number,
  comissao: number
): number {
  const divisor = 1 - margem - imposto - comissao;
  if (divisor <= 0) return 0;
  return custoTotal / divisor;
}

/**
 * Calcula a margem final real
 * Margem = (Preço - Custo - Impostos - Comissão) / Preço
 */
export function calcularMargemFinal(
  preco: number,
  custoTotal: number,
  imposto: number,
  comissao: number
): number {
  if (preco <= 0) return 0;
  const impostoValor = preco * imposto;
  const comissaoValor = preco * comissao;
  return (preco - custoTotal - impostoValor - comissaoValor) / preco;
}

/**
 * Formata percentual para exibição
 */
export function formatPercentual(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return `${num.toFixed(2)}%`;
}

/**
 * Formata valor monetário
 */
export function formatCurrency(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
}

/**
 * Converte string percentual para decimal
 * Ex: "35.00" -> 0.35
 */
export function percentualParaDecimal(percentual: string): number {
  return parseFloat(percentual) / 100;
}

/**
 * Converte decimal para string percentual
 * Ex: 0.35 -> "35.00"
 */
export function decimalParaPercentual(decimal: number): string {
  return (decimal * 100).toFixed(2);
}
