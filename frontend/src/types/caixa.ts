/**
 * NEXO - Sistema de Gestão para Barbearias
 * Caixa Types
 *
 * Tipos TypeScript para o módulo Caixa Diário.
 * Espelhados de: backend/internal/application/dto/caixa_dto.go
 *
 * @author NEXO v2.0
 */

// =============================================================================
// ENUMS
// =============================================================================

/**
 * Status do Caixa Diário
 */
export enum StatusCaixa {
  ABERTO = 'ABERTO',
  FECHADO = 'FECHADO',
}

/**
 * Tipo de operação no caixa
 */
export enum TipoOperacaoCaixa {
  SANGRIA = 'SANGRIA',
  REFORCO = 'REFORCO',
  VENDA = 'VENDA',
  DESPESA = 'DESPESA',
}

/**
 * Destino da sangria
 */
export enum DestinoSangria {
  DEPOSITO = 'DEPOSITO',
  PAGAMENTO = 'PAGAMENTO',
  COFRE = 'COFRE',
  OUTROS = 'OUTROS',
}

/**
 * Origem do reforço
 */
export enum OrigemReforco {
  TROCO = 'TROCO',
  CAPITAL_GIRO = 'CAPITAL_GIRO',
  TRANSFERENCIA = 'TRANSFERENCIA',
  OUTROS = 'OUTROS',
}

// =============================================================================
// REQUESTS
// =============================================================================

/**
 * Requisição para abrir o caixa
 */
export interface AbrirCaixaRequest {
  saldo_inicial: string; // Decimal como string
}

/**
 * Requisição para registrar sangria
 */
export interface SangriaRequest {
  valor: string; // Decimal como string
  destino: DestinoSangria;
  descricao: string;
}

/**
 * Requisição para registrar reforço
 */
export interface ReforcoRequest {
  valor: string; // Decimal como string
  origem: OrigemReforco;
  descricao: string;
}

/**
 * Requisição para fechar o caixa
 */
export interface FecharCaixaRequest {
  saldo_real: string; // Decimal como string
  justificativa?: string;
}

/**
 * Filtros para listagem do histórico
 */
export interface ListCaixaHistoricoFilters {
  data_inicio?: string;
  data_fim?: string;
  usuario_id?: string;
  page?: number;
  page_size?: number;
}

// =============================================================================
// RESPONSES
// =============================================================================

/**
 * Resposta de uma operação do caixa
 */
export interface OperacaoCaixaResponse {
  id: string;
  tipo: TipoOperacaoCaixa;
  valor: string; // Decimal como string
  descricao: string;
  destino?: string;
  origem?: string;
  usuario_id: string;
  usuario_nome: string;
  created_at: string;
}

/**
 * Resposta completa de um Caixa Diário
 */
export interface CaixaDiarioResponse {
  id: string;
  usuario_abertura_id: string;
  usuario_abertura_nome: string;
  usuario_fechamento_id?: string;
  usuario_fechamento_nome?: string;
  data_abertura: string;
  data_fechamento?: string;
  saldo_inicial: string; // Decimal como string
  total_entradas: string;
  total_saidas: string;
  total_sangrias: string;
  total_reforcos: string;
  saldo_esperado: string;
  saldo_real?: string;
  divergencia?: string;
  status: StatusCaixa;
  justificativa_divergencia?: string;
  created_at: string;
  updated_at: string;
  operacoes?: OperacaoCaixaResponse[];
}

/**
 * Resumo do caixa para listagens
 */
export interface CaixaDiarioResumoResponse {
  id: string;
  usuario_abertura_nome: string;
  usuario_fechamento_nome?: string;
  data_abertura: string;
  data_fechamento?: string;
  saldo_inicial: string;
  saldo_esperado: string;
  saldo_real?: string;
  divergencia?: string;
  status: StatusCaixa;
  tem_divergencia: boolean;
}

/**
 * Status atual do caixa
 */
export interface CaixaStatusResponse {
  aberto: boolean;
  caixa_atual?: CaixaDiarioResponse;
  ultimo_fechamento?: string;
}

/**
 * Resposta paginada do histórico
 */
export interface ListCaixaHistoricoResponse {
  items: CaixaDiarioResumoResponse[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

/**
 * Totais do caixa por tipo de operação
 */
export interface TotaisCaixaResponse {
  total_vendas: string;
  total_sangrias: string;
  total_reforcos: string;
  total_despesas: string;
  saldo_atual: string;
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Labels para destino de sangria
 */
export const DestinoSangriaLabels: Record<DestinoSangria, string> = {
  [DestinoSangria.DEPOSITO]: 'Depósito Bancário',
  [DestinoSangria.PAGAMENTO]: 'Pagamento',
  [DestinoSangria.COFRE]: 'Cofre',
  [DestinoSangria.OUTROS]: 'Outros',
};

/**
 * Labels para origem de reforço
 */
export const OrigemReforcoLabels: Record<OrigemReforco, string> = {
  [OrigemReforco.TROCO]: 'Troco',
  [OrigemReforco.CAPITAL_GIRO]: 'Capital de Giro',
  [OrigemReforco.TRANSFERENCIA]: 'Transferência',
  [OrigemReforco.OUTROS]: 'Outros',
};

/**
 * Labels para tipo de operação
 */
export const TipoOperacaoLabels: Record<TipoOperacaoCaixa, string> = {
  [TipoOperacaoCaixa.SANGRIA]: 'Sangria',
  [TipoOperacaoCaixa.REFORCO]: 'Reforço',
  [TipoOperacaoCaixa.VENDA]: 'Venda',
  [TipoOperacaoCaixa.DESPESA]: 'Despesa',
};

/**
 * Labels para status do caixa
 */
export const StatusCaixaLabels: Record<StatusCaixa, string> = {
  [StatusCaixa.ABERTO]: 'Aberto',
  [StatusCaixa.FECHADO]: 'Fechado',
};

/**
 * Cores para status do caixa (Tailwind classes usando tokens semânticos)
 */
export const StatusCaixaColors: Record<StatusCaixa, string> = {
  [StatusCaixa.ABERTO]: 'bg-chart-2/10 text-chart-2 border-chart-2/30',
  [StatusCaixa.FECHADO]: 'bg-muted text-muted-foreground border-muted',
};

/**
 * Cores para tipo de operação (Tailwind classes usando tokens semânticos)
 */
export const TipoOperacaoColors: Record<TipoOperacaoCaixa, string> = {
  [TipoOperacaoCaixa.SANGRIA]: 'text-destructive',
  [TipoOperacaoCaixa.REFORCO]: 'text-chart-1',
  [TipoOperacaoCaixa.VENDA]: 'text-chart-2',
  [TipoOperacaoCaixa.DESPESA]: 'text-chart-5',
};
