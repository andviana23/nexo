/**
 * Tipos para Categorias de Produtos
 * Módulo de Estoque - NEXO v1.0
 */

// Centro de Custo para classificação DRE
export type CentroCusto = 'CMV' | 'CUSTO_SERVICO' | 'DESPESA_OPERACIONAL';

// Constantes de Centro de Custo
export const CENTRO_CUSTO_OPTIONS = [
  { value: 'CMV', label: 'CMV (Revenda)', description: 'Custo Mercadoria Vendida' },
  { value: 'CUSTO_SERVICO', label: 'Custo Serviço', description: 'Insumos usados em serviços' },
  { value: 'DESPESA_OPERACIONAL', label: 'Despesa Operacional', description: 'Material escritório, limpeza' },
] as const;

// Categoria de Produto (response do backend)
export interface CategoriaProduto {
  id: string;
  nome: string;
  descricao: string;
  cor: string;
  icone: string;
  centro_custo: CentroCusto;
  ativa: boolean;
  criado_em: string;
  atualizado_em: string;
}

// Request para criar categoria
export interface CreateCategoriaProdutoRequest {
  nome: string;
  descricao?: string;
  cor?: string;
  icone?: string;
  centro_custo?: CentroCusto;
}

// Request para atualizar categoria
export interface UpdateCategoriaProdutoRequest {
  nome: string;
  descricao?: string;
  cor?: string;
  icone?: string;
  centro_custo?: CentroCusto;
  ativa?: boolean;
}

// Response lista de categorias
export interface ListCategoriaProdutoResponse {
  categorias: CategoriaProduto[];
  total: number;
}

// Cores padrão para categorias
export const CATEGORIA_CORES = [
  { value: '#6B7280', label: 'Cinza' },
  { value: '#EF4444', label: 'Vermelho' },
  { value: '#F97316', label: 'Laranja' },
  { value: '#EAB308', label: 'Amarelo' },
  { value: '#22C55E', label: 'Verde' },
  { value: '#14B8A6', label: 'Turquesa' },
  { value: '#3B82F6', label: 'Azul' },
  { value: '#8B5CF6', label: 'Roxo' },
  { value: '#EC4899', label: 'Rosa' },
] as const;

// Ícones disponíveis (Lucide)
export const CATEGORIA_ICONES = [
  { value: 'package', label: 'Pacote' },
  { value: 'droplet', label: 'Gota (Shampoo)' },
  { value: 'scissors', label: 'Tesoura' },
  { value: 'beer', label: 'Bebida' },
  { value: 'spray-can', label: 'Spray' },
  { value: 'shirt', label: 'Vestuário' },
  { value: 'sparkles', label: 'Brilho' },
  { value: 'brush', label: 'Escova' },
  { value: 'box', label: 'Caixa' },
  { value: 'archive', label: 'Arquivo' },
] as const;
