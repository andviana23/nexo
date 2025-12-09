/**
 * NEXO - Sistema de Gestão para Barbearias
 * Stock Service
 *
 * Serviço para gerenciar estoque (itens, movimentações).
 */

import { api } from '@/lib/axios';
import type {
    StockEntryFormData,
    StockEntryMultipleItemsRequest,
    StockEntryResponse,
    StockExitFormData,
    StockExitResponse,
    StockFilters,
    StockInventoryResponse,
    StockItem,
    StockItemFormData,
    StockMovement,
    StockMovementFilters,
    StockMovementResponse,
} from '@/types/stock';

// =============================================================================
// STOCK ITEMS (INVENTÁRIO)
// =============================================================================

/**
 * Lista itens do estoque com filtros
 */
export async function listStockItems(
  filters?: StockFilters
): Promise<StockInventoryResponse> {
  const { data } = await api.get<StockInventoryResponse>('/stock/items', {
    params: filters,
  });
  return data;
}

/**
 * Busca item do estoque por ID
 */
export async function getStockItem(id: string): Promise<StockItem> {
  const { data } = await api.get<StockItem>(`/stock/items/${id}`);
  return data;
}

/**
 * Cria novo item no estoque
 */
export async function createStockItem(
  itemData: StockItemFormData
): Promise<StockItem> {
  const { data } = await api.post<StockItem>('/stock/items', itemData);
  return data;
}

/**
 * Atualiza item do estoque
 */
export async function updateStockItem(
  id: string,
  itemData: Partial<StockItemFormData>
): Promise<StockItem> {
  const { data } = await api.patch<StockItem>(
    `/stock/items/${id}`,
    itemData
  );
  return data;
}

/**
 * Deleta item do estoque
 */
export async function deleteStockItem(id: string): Promise<void> {
  await api.delete(`/stock/items/${id}`);
}

// =============================================================================
// PRODUTOS (CRUD conforme backend DTO)
// =============================================================================

/**
 * Tipo para criar produto (alinhado com CreateProdutoRequest do backend)
 */
export interface CreateProductRequest {
  nome: string;
  descricao?: string;
  codigo_barras?: string;
  categoria_produto_id: string; // FK para categoria customizada (obrigatório)
  unidade_medida: string; // StockUnit
  valor_unitario: string;
  quantidade_minima: number;
  quantidade_maxima?: string;
  valor_venda_profissional?: string;
  valor_entrada?: string;
  fornecedor_id?: string;
}

/**
 * Cria novo produto no estoque
 */
export async function createProduct(data: CreateProductRequest): Promise<StockItem> {
  const { data: response } = await api.post<StockItem>('/stock/products', data);
  return response;
}

// =============================================================================
// STOCK MOVEMENTS (ENTRADAS/SAÍDAS)
// =============================================================================

/**
 * Lista movimentações de estoque
 */
export async function listStockMovements(
  filters?: StockMovementFilters
): Promise<StockMovementResponse> {
  const { data } = await api.get<StockMovementResponse>(
    '/stock/movements',
    { params: filters }
  );
  return data;
}

/**
 * Busca movimentação por ID
 */
export async function getStockMovement(id: string): Promise<StockMovement> {
  const { data } = await api.get<StockMovement>(
    `/stock/movements/${id}`
  );
  return data;
}

/**
 * Registra entrada de estoque
 */
export async function createStockEntry(
  entryData: StockEntryFormData
): Promise<StockEntryResponse> {
  const payload = {
    fornecedor_id: entryData.supplier_id,
    data_entrada: entryData.entry_date,
    observacoes: [entryData.reference, entryData.notes].filter(Boolean).join(' - ') || undefined,
    gerar_financeiro: entryData.generate_financial ?? false,
    itens: [
      {
        produto_id: entryData.stock_item_id,
        quantidade: entryData.quantity,
        valor_unitario: entryData.unit_cost
          ? entryData.unit_cost.replace(',', '.')
          : '0',
      },
    ],
  };

  const { data } = await api.post<StockEntryResponse>('/stock/entries', payload);
  return data;
}

/**
 * Registra entrada de estoque com múltiplos produtos
 * Alinhado com RegistrarEntradaRequest do backend
 */
export async function createStockEntryMultiple(
  entryData: StockEntryMultipleItemsRequest
): Promise<StockEntryResponse> {
  const { data } = await api.post<StockEntryResponse>('/stock/entries', entryData);
  return data;
}

/**
 * Registra saída de estoque
 */
export async function createStockExit(
  exitData: StockExitFormData
): Promise<StockExitResponse> {
  const payload = {
    produto_id: exitData.stock_item_id,
    quantidade: String(exitData.quantity),
    motivo: exitData.reason,
    observacoes: [exitData.reference, exitData.notes].filter(Boolean).join(' - ') || undefined,
  };

  const { data } = await api.post<StockExitResponse>('/stock/exit', payload);
  return data;
}

// =============================================================================
// RELATÓRIOS/ANALYTICS
// =============================================================================

/**
 * Obtém resumo do estoque (itens em falta, baixo estoque, etc)
 */
export async function getStockSummary(): Promise<{
  total_items: number;
  low_stock_count: number;
  out_of_stock_count: number;
  total_value: string;
}> {
  const { data } = await api.get('/stock/summary');
  return data;
}

/**
 * Obtém histórico de movimentações de um item
 */
export async function getStockItemHistory(
  itemId: string,
  filters?: StockMovementFilters
): Promise<StockMovementResponse> {
  const { data } = await api.get<StockMovementResponse>(
    `/stock/items/${itemId}/history`,
    { params: filters }
  );
  return data;
}
