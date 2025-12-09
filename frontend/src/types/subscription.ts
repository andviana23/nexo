/**
 * NEXO - Sistema de Gestão para Barbearias
 * Types: Subscription (Assinaturas)
 * 
 * @module types/subscription
 * @description Definições de tipos para o módulo de assinaturas
 * Conforme FLUXO_ASSINATURA.md
 */

// ============================================================================
// ENUMS
// ============================================================================

/** Status da assinatura */
export type SubscriptionStatus = 
  | 'PENDENTE'      // Aguardando pagamento inicial
  | 'ATIVO'         // Pagamento confirmado, assinatura ativa
  | 'INADIMPLENTE'  // Pagamento vencido
  | 'CANCELADO'     // Cancelado manualmente
  | 'EXPIRADO';     // Venceu e não foi renovado

/** Formas de pagamento */
export type PaymentMethod = 'CARTAO' | 'PIX' | 'DINHEIRO';

/** Status do pagamento */
export type PaymentStatus = 'PENDENTE' | 'CONFIRMADO' | 'VENCIDO' | 'CANCELADO';

/** Periodicidade do plano */
export type PlanPeriodicity = 'MENSAL' | 'TRIMESTRAL' | 'SEMESTRAL' | 'ANUAL';

// ============================================================================
// ENTITIES - PLANOS
// ============================================================================

/** Entidade Plano de Assinatura */
export interface Plan {
  id: string;
  tenant_id: string;
  nome: string;
  descricao?: string;
  valor: string; // Decimal como string
  periodicidade: PlanPeriodicity;
  qtd_servicos?: number;
  limite_uso_mensal?: number;
  ativo: boolean;
  created_at: string;
  updated_at: string;
}

/** Resumo do plano para selects */
export interface PlanSummary {
  id: string;
  nome: string;
  valor: string;
  periodicidade: PlanPeriodicity;
  ativo: boolean;
}

// ============================================================================
// ENTITIES - ASSINATURAS
// ============================================================================

/** Entidade Assinatura */
export interface Subscription {
  id: string;
  tenant_id: string;
  cliente_id: string;
  plano_id: string;
  
  // Dados do Asaas
  asaas_customer_id?: string;
  asaas_subscription_id?: string;
  
  // Pagamento
  forma_pagamento: PaymentMethod;
  status: SubscriptionStatus;
  valor: string; // Decimal como string
  link_pagamento?: string;
  codigo_transacao?: string;
  
  // Datas
  data_ativacao?: string;
  data_vencimento?: string;
  data_cancelamento?: string;
  cancelado_por?: string;
  
  // Uso
  servicos_utilizados: number;
  
  // Timestamps
  created_at: string;
  updated_at: string;
  
  // Dados relacionados (quando JOIN)
  cliente_nome?: string;
  cliente_telefone?: string;
  cliente_email?: string;
  plano_nome?: string;
  plano_qtd_servicos?: number;
  plano_limite_uso_mensal?: number;
}

/** Assinatura com dados completos */
export interface SubscriptionWithDetails extends Subscription {
  cliente: {
    id: string;
    nome: string;
    telefone: string;
    email?: string;
  };
  plano: Plan;
  pagamentos: SubscriptionPayment[];
}

// ============================================================================
// ENTITIES - PAGAMENTOS
// ============================================================================

/** Pagamento de assinatura */
export interface SubscriptionPayment {
  id: string;
  tenant_id: string;
  subscription_id: string;
  asaas_payment_id?: string;
  valor: string;
  forma_pagamento: PaymentMethod;
  status: PaymentStatus;
  data_pagamento?: string;
  codigo_transacao?: string;
  observacao?: string;
  created_at: string;
}

// ============================================================================
// ENTITIES - MÉTRICAS
// ============================================================================

/** Métricas de assinaturas */
export interface SubscriptionMetrics {
  total_assinantes_ativos: number;
  total_inativas: number;
  total_inadimplentes: number;
  total_planos_ativos: number;
  receita_mensal: number;
  taxa_renovacao: number;
  renovacoes_proximos_7_dias: number;
}

// ============================================================================
// DTOs - REQUEST (PLANOS)
// ============================================================================

/** DTO para criar plano */
export interface CreatePlanRequest {
  nome: string;
  descricao?: string;
  valor: string;
  periodicidade: PlanPeriodicity;
  qtd_servicos?: number;
  limite_uso_mensal?: number;
}

/** DTO para atualizar plano */
export interface UpdatePlanRequest {
  nome?: string;
  descricao?: string;
  valor?: string;
  periodicidade?: PlanPeriodicity;
  qtd_servicos?: number;
  limite_uso_mensal?: number;
  ativo?: boolean;
}

// ============================================================================
// DTOs - REQUEST (ASSINATURAS)
// ============================================================================

/** DTO para criar assinatura */
export interface CreateSubscriptionRequest {
  cliente_id: string;
  plano_id: string;
  forma_pagamento: PaymentMethod;
}

/** DTO para renovar assinatura (PIX/Dinheiro) */
export interface RenewSubscriptionRequest {
  codigo_transacao?: string;
  observacao?: string;
}

// ============================================================================
// DTOs - RESPONSE
// ============================================================================

/** Resposta de plano */
export interface PlanResponse {
  data: Plan;
}

/** Resposta de lista de planos */
export interface ListPlansResponse {
  data: Plan[];
  total: number;
}

/** Resposta de assinatura */
export interface SubscriptionResponse {
  data: Subscription;
}

/** Resposta de lista de assinaturas */
export interface ListSubscriptionsResponse {
  data: Subscription[];
  total: number;
}

/** Resposta de métricas */
export interface SubscriptionMetricsResponse {
  data: SubscriptionMetrics;
}

// ============================================================================
// FILTROS
// ============================================================================

/** Filtros para listagem de assinaturas */
export interface ListSubscriptionsFilters {
  page?: number;
  page_size?: number;
  status?: SubscriptionStatus;
  plano_id?: string;
  search?: string;
}

/** Filtros para listagem de planos */
export interface ListPlansFilters {
  page?: number;
  page_size?: number;
  search?: string;
  ativo?: boolean;
}

// ============================================================================
// MODAL STATES
// ============================================================================

/** Estado do modal de plano */
export interface PlanModalState {
  isOpen: boolean;
  mode: 'create' | 'edit' | 'view';
  plan?: Plan;
}

/** Estado do modal de assinatura */
export interface SubscriptionModalState {
  isOpen: boolean;
  mode: 'view' | 'renew' | 'cancel';
  subscription?: Subscription;
}

// ============================================================================
// LABELS E CONSTANTES
// ============================================================================

/** Labels de periodicidade */
export const PERIODICITY_LABELS: Record<PlanPeriodicity, string> = {
  MENSAL: 'Mensal',
  TRIMESTRAL: 'Trimestral',
  SEMESTRAL: 'Semestral',
  ANUAL: 'Anual',
};

/** Labels de status de assinatura */
export const SUBSCRIPTION_STATUS_LABELS: Record<SubscriptionStatus, string> = {
  PENDENTE: 'Pendente',
  ATIVO: 'Ativo',
  INADIMPLENTE: 'Inadimplente',
  CANCELADO: 'Cancelado',
  EXPIRADO: 'Expirado',
};

/** Labels de forma de pagamento */
export const PAYMENT_METHOD_LABELS: Record<PaymentMethod, string> = {
  CARTAO: 'Cartão de Crédito',
  PIX: 'PIX',
  DINHEIRO: 'Dinheiro',
};

/** Labels de status de pagamento */
export const PAYMENT_STATUS_LABELS: Record<PaymentStatus, string> = {
  PENDENTE: 'Pendente',
  CONFIRMADO: 'Confirmado',
  VENCIDO: 'Vencido',
  CANCELADO: 'Cancelado',
};

/** Cores de status de assinatura */
export const SUBSCRIPTION_STATUS_COLORS: Record<SubscriptionStatus, string> = {
  PENDENTE: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400',
  ATIVO: 'bg-green-100 text-green-700 dark:bg-green-900/20 dark:text-green-400',
  INADIMPLENTE: 'bg-red-100 text-red-700 dark:bg-red-900/20 dark:text-red-400',
  CANCELADO: 'bg-gray-100 text-gray-700 dark:bg-gray-900/20 dark:text-gray-400',
  EXPIRADO: 'bg-orange-100 text-orange-700 dark:bg-orange-900/20 dark:text-orange-400',
};
