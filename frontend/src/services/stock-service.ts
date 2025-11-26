/**
 * NEXO - Sistema de Gestão para Barbearias
 * Stock Service
 *
 * Serviço para gerenciar estoque (itens, movimentações).
 */

import { api } from '@/lib/axios';
import type {
  StockEntryFormData,
  StockExitFormData,
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
): Promise<StockMovement> {
  const { data } = await api.post<StockMovement>(
    '/stock/entries',
    entryData
  );
  return data;
}

/**
 * Registra saída de estoque
 */
export async function createStockExit(
  exitData: StockExitFormData
): Promise<StockMovement> {
  const { data } = await api.post<StockMovement>(
    '/stock/exits',
    exitData
  );
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
