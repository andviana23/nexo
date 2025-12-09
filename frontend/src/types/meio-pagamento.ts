/**
 * NEXO - Tipos para Meios de Pagamento
 * Tipos de Recebimento / Formas de Pagamento
 */

// Tipos de pagamento disponíveis
export type TipoPagamento =
  | 'DINHEIRO'
  | 'PIX'
  | 'CREDITO'
  | 'DEBITO'
  | 'TRANSFERENCIA'
  | 'BOLETO'
  | 'OUTRO';

// Labels amigáveis para os tipos
export const TIPO_PAGAMENTO_LABELS: Record<TipoPagamento, string> = {
  DINHEIRO: 'Dinheiro',
  PIX: 'PIX',
  CREDITO: 'Cartão de Crédito',
  DEBITO: 'Cartão de Débito',
  TRANSFERENCIA: 'Transferência',
  BOLETO: 'Boleto',
  OUTRO: 'Outro',
};

// Bandeiras de cartão
export const BANDEIRAS_CARTAO = [
  'Visa',
  'Mastercard',
  'Elo',
  'Hipercard',
  'American Express',
  'Diners',
  'Alelo',
  'Sodexo',
  'VR',
  'Ticket',
  'Outra',
] as const;

export type BandeiraCartao = (typeof BANDEIRAS_CARTAO)[number];

// Interface principal
export interface MeioPagamento {
  id: string;
  tenant_id: string;
  nome: string;
  tipo: TipoPagamento;
  bandeira?: string;
  taxa: string;
  taxa_fixa: string;
  d_mais: number;
  icone?: string;
  cor?: string;
  ordem_exibicao: number;
  observacoes?: string;
  ativo: boolean;
  criado_em: string;
  atualizado_em: string;
}

// DTO para criação
export interface CreateMeioPagamentoDTO {
  nome: string;
  tipo: TipoPagamento;
  bandeira?: string;
  taxa?: string;
  taxa_fixa?: string;
  d_mais?: number;
  icone?: string;
  cor?: string;
  ordem_exibicao?: number;
  observacoes?: string;
  ativo?: boolean;
}

// DTO para atualização
export type UpdateMeioPagamentoDTO = Partial<CreateMeioPagamentoDTO>;

// Filtros para listagem
export interface MeioPagamentoFilters {
  apenas_ativos?: boolean;
  tipo?: TipoPagamento;
}

// Resposta de listagem
export interface MeioPagamentoListResponse {
  data: MeioPagamento[];
  total: number;
  total_ativo: number;
}

// Valores padrão de D+ por tipo
export const D_MAIS_PADRAO: Record<TipoPagamento, number> = {
  DINHEIRO: 0,
  PIX: 0,
  DEBITO: 1,
  CREDITO: 30,
  TRANSFERENCIA: 1,
  BOLETO: 2,
  OUTRO: 0,
};

// Cores padrão por tipo
export const CORES_TIPO_PAGAMENTO: Record<TipoPagamento, string> = {
  DINHEIRO: '#22c55e', // green-500
  PIX: '#06b6d4', // cyan-500
  CREDITO: '#3b82f6', // blue-500
  DEBITO: '#8b5cf6', // violet-500
  TRANSFERENCIA: '#f59e0b', // amber-500
  BOLETO: '#6b7280', // gray-500
  OUTRO: '#64748b', // slate-500
};

// Ícones por tipo (nomes do lucide-react)
export const ICONES_TIPO_PAGAMENTO: Record<TipoPagamento, string> = {
  DINHEIRO: 'banknote',
  PIX: 'qr-code',
  CREDITO: 'credit-card',
  DEBITO: 'credit-card',
  TRANSFERENCIA: 'arrow-left-right',
  BOLETO: 'file-text',
  OUTRO: 'circle-dollar-sign',
};
