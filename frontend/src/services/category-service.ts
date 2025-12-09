import { api } from '@/lib/axios';
import {
    Category,
    CategoryFilters,
    CreateCategoryDTO,
    UpdateCategoryDTO,
} from '@/types/category';

const BASE_URL = '/categorias-servicos';

// Interface para resposta da API de listagem
interface ListCategoriesResponse {
  categorias: Category[];
  total: number;
}

export const categoryService = {
  /**
   * Lista todas as categorias
   */
  list: async (filters?: CategoryFilters): Promise<Category[]> => {
    const response = await api.get<ListCategoriesResponse>(BASE_URL, {
      params: filters,
    });
    // API retorna { categorias: [...], total: N }, extraímos o array
    return response.data.categorias || [];
  },

  /**
   * Busca uma categoria por ID
   */
  getById: async (id: string): Promise<Category> => {
    const response = await api.get<Category>(`${BASE_URL}/${id}`);
    return response.data;
  },

  /**
   * Cria uma nova categoria
   */
  create: async (data: CreateCategoryDTO): Promise<Category> => {
    const response = await api.post<Category>(BASE_URL, data);
    return response.data;
  },

  /**
   * Atualiza uma categoria existente
   */
  update: async (id: string, data: UpdateCategoryDTO): Promise<Category> => {
    const response = await api.put<Category>(`${BASE_URL}/${id}`, data);
    return response.data;
  },

  /**
   * Remove uma categoria
   */
  delete: async (id: string): Promise<void> => {
    await api.delete(`${BASE_URL}/${id}`);
  },
};
