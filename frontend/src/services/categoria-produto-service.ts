/**
 * Service para Categorias de Produtos
 * Módulo de Estoque - NEXO v1.0
 */

import { api } from '@/lib/axios';
import {
    CategoriaProduto,
    CreateCategoriaProdutoRequest,
    ListCategoriaProdutoResponse,
    UpdateCategoriaProdutoRequest,
} from '@/types/categoria-produto';

const BASE_URL = '/categorias-produtos';

export const categoriaProdutoService = {
  /**
   * Lista todas as categorias de produtos
   * @param apenasAtivas - Se true, retorna apenas categorias ativas
   */
  list: async (apenasAtivas?: boolean): Promise<CategoriaProduto[]> => {
    const response = await api.get<ListCategoriaProdutoResponse>(BASE_URL, {
      params: apenasAtivas ? { ativas: 'true' } : undefined,
    });
    return response.data.categorias || [];
  },

  /**
   * Busca uma categoria por ID
   */
  getById: async (id: string): Promise<CategoriaProduto> => {
    const response = await api.get<CategoriaProduto>(`${BASE_URL}/${id}`);
    return response.data;
  },

  /**
   * Cria uma nova categoria de produto
   */
  create: async (data: CreateCategoriaProdutoRequest): Promise<CategoriaProduto> => {
    const response = await api.post<CategoriaProduto>(BASE_URL, data);
    return response.data;
  },

  /**
   * Atualiza uma categoria existente
   */
  update: async (id: string, data: UpdateCategoriaProdutoRequest): Promise<CategoriaProduto> => {
    const response = await api.put<CategoriaProduto>(`${BASE_URL}/${id}`, data);
    return response.data;
  },

  /**
   * Remove uma categoria (apenas se não houver produtos vinculados)
   */
  delete: async (id: string): Promise<void> => {
    await api.delete(`${BASE_URL}/${id}`);
  },

  /**
   * Ativa/desativa uma categoria (toggle)
   */
  toggle: async (id: string): Promise<CategoriaProduto> => {
    const response = await api.patch<CategoriaProduto>(`${BASE_URL}/${id}/toggle`);
    return response.data;
  },
};
