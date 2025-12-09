/**
 * Tipos TypeScript para o módulo de Comandas (estilo Trinks)
 * 
 * Representa o sistema de fechamento de conta com múltiplas formas de pagamento,
 * integrado com meios de pagamento e cálculo automático de taxas.
 */

export type CommandStatus = 'OPEN' | 'CLOSED' | 'CANCELED';

export type CommandItemType = 'SERVICO' | 'PRODUTO' | 'PACOTE';

/**
 * Comanda principal (conta do cliente)
 */
export interface Command {
  id: string;
  tenant_id: string;
  appointment_id?: string;
  customer_id: string;
  numero?: string;
  status: CommandStatus;
  
  // Valores financeiros
  subtotal: number;
  desconto: number;
  total: number;
  total_recebido: number;
  troco: number;
  saldo_devedor: number;
  
  // Opções de fechamento
  observacoes?: string;
  deixar_troco_gorjeta: boolean;
  deixar_saldo_divida: boolean;
  
  // Auditoria
  criado_em: string;
  atualizado_em: string;
  fechado_em?: string;
  fechado_por?: string;
  
  // Relacionamentos (eager loaded quando necessário)
  items?: CommandItem[];
  payments?: CommandPayment[];
}

/**
 * Item da comanda (serviço, produto, pacote)
 */
export interface CommandItem {
  id: string;
  command_id: string;
  tipo: CommandItemType;
  item_id: string; // ID do serviço/produto/pacote
  descricao: string;
  
  // Valores
  preco_unitario: number;
  quantidade: number;
  desconto_valor: number;
  desconto_percentual: number;
  preco_final: number;
  
  observacoes?: string;
  criado_em: string;
}

/**
 * Pagamento da comanda
 */
export interface CommandPayment {
  id: string;
  command_id: string;
  meio_pagamento_id: string;
  
  // Valores
  valor_recebido: number;
  taxa_percentual: number;
  taxa_fixa: number;
  valor_liquido: number;
  
  observacoes?: string;
  criado_em: string;
  criado_por?: string;
}

// ============================================================================
// DTOs para Requests
// ============================================================================

/**
 * Request para criar comanda
 */
export interface CreateCommandRequest {
  appointment_id?: string;
  customer_id: string;
  items: CommandItemInput[];
  observacoes?: string;
}

/**
 * Input de item ao criar comanda
 */
export interface CommandItemInput {
  tipo: CommandItemType;
  item_id: string;
  descricao: string;
  preco_unitario: string; // String conforme esperado pelo backend
  quantidade?: number;
  desconto_valor?: number;
  desconto_percentual?: number;
}

/**
 * Request para adicionar item à comanda existente
 */
export interface AddCommandItemRequest {
  tipo: CommandItemType;
  item_id: string;
  descricao: string;
  preco_unitario: string; // String conforme esperado pelo backend
  quantidade?: number;
  desconto_valor?: number;
  desconto_percentual?: number;
  observacoes?: string;
}

/**
 * Request para atualizar item
 */
export interface UpdateCommandItemRequest {
  preco_unitario?: string; // String conforme esperado pelo backend
  quantidade?: number;
  desconto_valor?: number;
  desconto_percentual?: number;
  observacoes?: string;
}

/**
 * Request para adicionar pagamento
 */
export interface AddCommandPaymentRequest {
  meio_pagamento_id: string;
  valor_recebido: string; // Backend espera string para valores monetários
  observacoes?: string;
}

/**
 * Request para fechar comanda
 */
export interface CloseCommandRequest {
  deixar_troco_gorjeta?: boolean;
  deixar_saldo_divida?: boolean;
  observacoes?: string;
}

// ============================================================================
// DTOs para Responses
// ============================================================================

export interface CommandResponse {
  id: string;
  tenant_id: string;
  appointment_id?: string;
  customer_id: string;
  numero?: string;
  status: CommandStatus;
  subtotal: string; // Money formatado
  desconto: string;
  total: string;
  total_recebido: string;
  troco: string;
  saldo_devedor: string;
  observacoes?: string;
  deixar_troco_gorjeta: boolean;
  deixar_saldo_divida: boolean;
  criado_em: string;
  atualizado_em: string;
  fechado_em?: string;
  fechado_por?: string;
  items?: CommandItemResponse[];
  payments?: CommandPaymentResponse[];
}

export interface CommandItemResponse {
  id: string;
  command_id: string;
  tipo: CommandItemType;
  item_id: string;
  descricao: string;
  preco_unitario: string;
  quantidade: number;
  desconto_valor: string;
  desconto_percentual: number;
  preco_final: string;
  observacoes?: string;
  criado_em: string;
}

export interface CommandPaymentResponse {
  id: string;
  command_id: string;
  meio_pagamento_id: string;
  valor_recebido: string;
  taxa_percentual: number;
  taxa_fixa: string;
  valor_liquido: string;
  observacoes?: string;
  criado_em: string;
  criado_por?: string;
}

// ============================================================================
// State Management Types
// ============================================================================

/**
 * Estado de um pagamento selecionado no form
 */
export interface SelectedPayment {
  meio_pagamento_id: string;
  nome: string;
  tipo: string;
  valor_recebido: number;
  taxa_percentual: number;
  taxa_fixa: number;
  valor_liquido: number;
  icone?: string;
  cor?: string;
}

/**
 * Estado do formulário de pagamento
 */
export interface PaymentFormState {
  selected_payments: SelectedPayment[];
  deixar_troco_gorjeta: boolean;
  deixar_saldo_divida: boolean;
  observacoes: string;
}

/**
 * Resumo financeiro calculado
 */
export interface ResumoFinanceiro {
  subtotal: number;
  desconto: number;
  total: number;
  total_recebido: number;
  total_taxas: number;
  total_liquido: number;
  falta: number;
  troco: number;
}

// ============================================================================
// Helpers
// ============================================================================

/**
 * Calcula valor líquido após aplicar taxas
 */
export function calcularValorLiquido(
  valorRecebido: number,
  taxaPercentual: number,
  taxaFixa: number
): number {
  const valorTaxaPercentual = (valorRecebido * taxaPercentual) / 100;
  return valorRecebido - valorTaxaPercentual - taxaFixa;
}

/**
 * Calcula resumo financeiro da comanda
 */
export function calcularResumoFinanceiro(
  subtotal: number,
  desconto: number,
  selectedPayments: SelectedPayment[]
): ResumoFinanceiro {
  const total = subtotal - desconto;
  const totalRecebido = selectedPayments.reduce((sum, p) => sum + p.valor_recebido, 0);
  const totalTaxas = selectedPayments.reduce(
    (sum, p) => sum + (p.valor_recebido - p.valor_liquido),
    0
  );
  const totalLiquido = selectedPayments.reduce((sum, p) => sum + p.valor_liquido, 0);
  const diferenca = totalRecebido - total;
  
  return {
    subtotal,
    desconto,
    total,
    total_recebido: totalRecebido,
    total_taxas: totalTaxas,
    total_liquido: totalLiquido,
    falta: diferenca < 0 ? Math.abs(diferenca) : 0,
    troco: diferenca > 0 ? diferenca : 0,
  };
}

/**
 * Valida se pode fechar a comanda
 */
export function canCloseCommand(
  resumo: ResumoFinanceiro,
  deixarSaldoDivida: boolean
): { valid: boolean; errors: string[] } {
  const errors: string[] = [];
  
  if (resumo.total_recebido === 0) {
    errors.push('Nenhum pagamento foi registrado');
  }
  
  if (resumo.falta > 0 && !deixarSaldoDivida) {
    errors.push(`Falta receber R$ ${resumo.falta.toFixed(2)}`);
  }
  
  return {
    valid: errors.length === 0,
    errors,
  };
}

/**
 * Formata valor monetário
 */
export function formatMoney(value: number | string): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
}

/**
 * Parse de string monetária para number
 */
export function parseMoney(value: string): number {
  return parseFloat(value.replace(/[^0-9,-]/g, '').replace(',', '.')) || 0;
}
