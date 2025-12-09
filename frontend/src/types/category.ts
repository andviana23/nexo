export interface Category {
  id: string;
  tenant_id: string;
  nome: string;
  descricao?: string;
  cor?: string;
  icone?: string;
  ativa: boolean;
  criado_em: string;
  atualizado_em: string;
}

export interface CreateCategoryDTO {
  nome: string;
  descricao?: string;
  cor?: string;
  icone?: string;
}

export interface UpdateCategoryDTO {
  nome?: string;
  descricao?: string;
  cor?: string;
  icone?: string;
  ativa?: boolean;
}

export interface CategoryFilters {
  apenas_ativas?: boolean;
  order_by?: string;
}
