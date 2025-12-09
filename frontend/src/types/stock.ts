/**
 * NEXO - Sistema de Gestão para Barbearias
 * Stock Types
 *
 * Tipos TypeScript para o módulo de Estoque.
 */

// =============================================================================
// ENUMS
// =============================================================================

export enum StockMovementType {
  ENTRY = 'entry',
  EXIT = 'exit',
  ADJUSTMENT = 'adjustment',
}

export enum StockCategory {
  POMADA = 'POMADA',
  SHAMPOO = 'SHAMPOO',
  CREME = 'CREME',
  LAMINA = 'LAMINA',
  TOALHA = 'TOALHA',
  LIMPEZA = 'LIMPEZA',
  ESCRITORIO = 'ESCRITORIO',
  BEBIDA = 'BEBIDA',
  REVENDA = 'REVENDA',
  INSUMO = 'INSUMO',
  USO_INTERNO = 'USO_INTERNO',
  PERMANENTE = 'PERMANENTE',
  PROMOCIONAL = 'PROMOCIONAL',
  KIT = 'KIT',
  SERVICO = 'SERVICO',
}

export enum StockCostCenter {
  CUSTO_SERVICO = 'CUSTO_SERVICO',
  DESPESA_OPERACIONAL = 'DESPESA_OPERACIONAL',
  CMV = 'CMV',
}

export enum StockUnit {
  UN = 'UN',
  KG = 'KG',
  L = 'L',
  ML = 'ML',
  G = 'G',
}

// =============================================================================
// STOCK ITEM (alinhado com ProdutoResponse do backend)
// =============================================================================

export interface StockItem {
  id: string;
  tenant_id: string;
  nome: string;
  descricao?: string;
  codigo_barras?: string;
  categoria_produto_id?: string;
  categoria_produto?: {
    id: string;
    nome: string;
  };
  unidade_medida: string;
  valor_unitario: string;
  quantidade_atual: string;
  quantidade_minima: string;
  quantidade_maxima?: string;
  valor_venda_profissional?: string;
  valor_entrada?: string;
  fornecedor_id?: string;
  fornecedor?: {
    id: string;
    razao_social: string;
    nome_fantasia?: string;
  };
  esta_baixo: boolean;
  ativo: boolean;
  created_at: string;
  updated_at: string;
}

// =============================================================================
// STOCK MOVEMENT
// =============================================================================

export interface StockMovement {
  id: string;
  tenant_id: string;
  stock_item_id: string;
  movement_type: StockMovementType;
  quantity: number;
  unit_cost?: string; // Decimal como string
  total_cost?: string; // Decimal como string
  notes?: string;
  reference?: string; // Ex: número de nota fiscal
  created_by: string;
  created_at: string;
  
  // Relacionamentos (join)
  stock_item?: StockItem;
  created_by_name?: string;
}

// =============================================================================
// FORMS
// =============================================================================

export interface StockEntryFormData {
  stock_item_id: string;
  quantity: number;
  unit_cost?: string;
  notes?: string;
  reference?: string;
  supplier_id: string;
  entry_date: string; // YYYY-MM-DD
  generate_financial?: boolean;
}

/**
 * Novo formato de entrada de estoque (múltiplos produtos)
 * Alinhado com RegistrarEntradaRequest do backend
 */
export interface StockEntryMultipleItemsRequest {
  fornecedor_id: string;
  data_entrada: string; // YYYY-MM-DD
  itens: Array<{
    produto_id: string; // UUID do produto
    quantidade: number; // int (será arredondado no backend)
    valor_unitario: string; // decimal como string
  }>;
  observacoes?: string;
  gerar_financeiro: boolean;
}

export interface StockExitFormData {
  stock_item_id: string;
  quantity: number;
  notes?: string;
  reference?: string;
  reason: 'VENDA' | 'USO_INTERNO' | 'PERDA' | 'DEVOLUCAO';
}

export interface StockEntryResponse {
  message: string;
  data: {
    movimentacoes_ids: string[];
    valor_total: string;
    itens_processados: number;
  };
}

export interface StockExitResponse {
  id: string;
  tenant_id: string;
  produto_id: string;
  produto_nome?: string;
  usuario_id: string;
  fornecedor_id?: string;
  tipo: string;
  quantidade: string;
  valor_unitario: string;
  valor_total: string;
  observacoes: string;
  data: string;
  created_at: string;
}

export interface StockItemFormData {
  name: string;
  description?: string;
  category: StockCategory;
  cost_center: StockCostCenter;
  unit: StockUnit;
  min_quantity: number;
  max_quantity?: number;
  cost_price: string;
  sale_price?: string;
  barcode?: string;
  sku?: string;
  supplier?: string;
  is_active?: boolean;
  control_validity: boolean;
  lead_time_days: number;
}

// =============================================================================
// API RESPONSES
// =============================================================================

export interface StockInventoryResponse {
  data: StockItem[];
  total: number;
  low_stock_count: number;
  out_of_stock_count: number;
}

export interface StockMovementResponse {
  movements: StockMovement[];
  total: number;
}

// =============================================================================
// FILTERS
// =============================================================================

export interface StockFilters {
  category?: StockCategory;
  is_active?: boolean;
  low_stock?: boolean; // Apenas itens abaixo do mínimo
  search?: string;
  page?: number;
  page_size?: number;
}

export interface StockMovementFilters {
  stock_item_id?: string;
  movement_type?: StockMovementType;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}
