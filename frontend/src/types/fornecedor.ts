/**
 * NEXO - Sistema de Gestão para Barbearias
 * Fornecedor Types
 *
 * Tipos TypeScript para o módulo de Fornecedores.
 * Alinhado com backend/internal/application/dto/stock_dto.go
 */

// =============================================================================
// FORNECEDOR (alinhado com FornecedorResponse do backend)
// =============================================================================

export interface Fornecedor {
  id: string;
  tenant_id: string;
  razao_social: string;
  nome_fantasia?: string;
  nome?: string; // Alias para compatibilidade
  cnpj?: string;
  email?: string;
  telefone: string;
  celular?: string;
  // Endereço
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  endereco?: string; // Endereço formatado
  // Banco
  banco?: string;
  agencia?: string;
  conta?: string;
  observacoes?: string;
  ativo: boolean;
  created_at: string;
  updated_at: string;
}

// Alias para compatibilidade com páginas que usam cidade/estado
export interface FornecedorDisplay extends Fornecedor {
  cidade?: string;
  estado?: string;
  cep?: string;
}

// Helper: nome de exibição do fornecedor
export function getFornecedorNome(f: Fornecedor): string {
  return f.nome_fantasia || f.razao_social;
}

// =============================================================================
// REQUEST/RESPONSE DTOs
// =============================================================================

export interface CreateFornecedorRequest {
  razao_social: string;
  nome_fantasia?: string;
  cnpj?: string;
  email?: string;
  telefone: string;
  celular?: string;
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  banco?: string;
  agencia?: string;
  conta?: string;
  observacoes?: string;
}

export interface UpdateFornecedorRequest {
  razao_social?: string;
  nome_fantasia?: string;
  cnpj?: string;
  email?: string;
  telefone?: string;
  celular?: string;
  endereco_logradouro?: string;
  endereco_numero?: string;
  endereco_complemento?: string;
  endereco_bairro?: string;
  endereco_cidade?: string;
  endereco_estado?: string;
  endereco_cep?: string;
  banco?: string;
  agencia?: string;
  conta?: string;
  observacoes?: string;
  ativo?: boolean;
}

// Resposta da listagem (alinhado com ListFornecedoresResponse do backend)
export interface ListFornecedoresResponse {
  fornecedores: Fornecedor[];
  total: number;
}
