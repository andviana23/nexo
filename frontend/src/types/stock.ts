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
  PRODUCT = 'product',
  CONSUMABLE = 'consumable',
  EQUIPMENT = 'equipment',
  OTHER = 'other',
}

export enum StockUnit {
  UNIT = 'unit',
  KG = 'kg',
  LITER = 'liter',
  ML = 'ml',
  GRAM = 'gram',
  BOX = 'box',
  PACK = 'pack',
}

// =============================================================================
// STOCK ITEM
// =============================================================================

export interface StockItem {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  category: StockCategory;
  unit: StockUnit;
  current_quantity: number;
  min_quantity: number;
  max_quantity?: number;
  cost_price: string; // Decimal como string
  sale_price?: string; // Decimal como string
  barcode?: string;
  sku?: string;
  supplier?: string;
  is_active: boolean;
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
}

export interface StockExitFormData {
  stock_item_id: string;
  quantity: number;
  notes?: string;
  reference?: string;
}

export interface StockItemFormData {
  name: string;
  description?: string;
  category: StockCategory;
  unit: StockUnit;
  min_quantity: number;
  max_quantity?: number;
  cost_price: string;
  sale_price?: string;
  barcode?: string;
  sku?: string;
  supplier?: string;
  is_active?: boolean;
}

// =============================================================================
// API RESPONSES
// =============================================================================

export interface StockInventoryResponse {
  items: StockItem[];
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
