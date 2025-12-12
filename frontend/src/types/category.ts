export interface Category {
  id: string;
  tenant_id: string;
  unit_id?: string;
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
  unit_id?: string;
}

export interface UpdateCategoryDTO {
  nome?: string;
  descricao?: string;
  cor?: string;
  icone?: string;
  ativa?: boolean;
  unit_id?: string;
}

export interface CategoryFilters {
  apenas_ativas?: boolean;
  order_by?: string;
  unit_id?: string;
}
