export interface Service {
  id: string;
  tenant_id: string;
  categoria_id?: string;
  categoria_nome?: string;
  categoria_cor?: string;
  nome: string;
  descricao?: string;
  preco: string;
  preco_centavos: number;
  duracao: number;
  duracao_formatada: string;
  comissao: string;
  cor?: string;
  imagem?: string;
  profissionais_ids?: string[];
  observacoes?: string;
  tags?: string[];
  ativo: boolean;
  criado_em: string;
  atualizado_em: string;
}

export interface CreateServiceDTO {
  categoria_id?: string;
  nome: string;
  descricao?: string;
  preco: string;
  duracao: number;
  comissao?: string;
  cor?: string;
  imagem?: string;
  profissionais_ids?: string[];
  observacoes?: string;
  tags?: string[];
}

export type UpdateServiceDTO = Partial<CreateServiceDTO> & {
  ativo?: boolean;
};

export interface ServiceFilters {
  apenas_ativos?: boolean;
  categoria_id?: string;
  profissional_id?: string;
  search?: string;
  order_by?: string;
}

export interface ServiceListResponse {
  servicos: Service[];
  total: number;
}

export interface ServiceStats {
  total_servicos: number;
  servicos_ativos: number;
  servicos_inativos: number;
  preco_medio: string;
  duracao_media: number;
  comissao_media: string;
}
