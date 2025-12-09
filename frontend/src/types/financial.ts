/**
 * NEXO - Sistema de Gestão para Barbearias
 * Financial Types
 *
 * Tipos TypeScript para o módulo Financeiro.
 * Mapeados a partir do backend Go (financial_handler.go)
 */

// =============================================================================
// ENUMS - STATUS
// =============================================================================

export enum StatusContaPagar {
  PENDENTE = 'PENDENTE',
  PAGO = 'PAGO',
  ATRASADO = 'ATRASADO',
  CANCELADO = 'CANCELADO',
}

export enum StatusContaReceber {
  PENDENTE = 'PENDENTE',
  RECEBIDO = 'RECEBIDO',
  PARCIAL = 'PARCIAL',
  ATRASADO = 'ATRASADO',
  CANCELADO = 'CANCELADO',
}

export enum StatusCompensacao {
  PREVISTO = 'PREVISTO',
  COMPENSADO = 'COMPENSADO',
  DIVERGENTE = 'DIVERGENTE',
  CANCELADO = 'CANCELADO',
}

// =============================================================================
// ENUMS - TIPOS
// =============================================================================

export enum TipoDespesa {
  FIXA = 'FIXA',
  VARIAVEL = 'VARIAVEL',
  EXTRAORDINARIA = 'EXTRAORDINARIA',
}

export enum OrigemReceita {
  SERVICO = 'SERVICO',
  PRODUTO = 'PRODUTO',
  ASSINATURA = 'ASSINATURA',
  OUTRO = 'OUTRO',
}

export enum Periodicidade {
  MENSAL = 'MENSAL',
  QUINZENAL = 'QUINZENAL',
  SEMANAL = 'SEMANAL',
  ANUAL = 'ANUAL',
}

// =============================================================================
// CONTA A PAGAR
// =============================================================================

export interface ContaPagar {
  id: string;
  tenant_id: string;
  descricao: string;
  categoria_id?: string;
  fornecedor?: string;
  valor: string; // Decimal como string (padrão backend)
  tipo: TipoDespesa;
  data_vencimento: string; // ISO date string
  data_pagamento?: string; // ISO date string
  status: StatusContaPagar;
  recorrente: boolean;
  periodicidade?: Periodicidade;
  pix_code?: string;
  comprovante_url?: string;
  observacoes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateContaPagarRequest {
  descricao: string;
  categoria_id?: string;
  fornecedor?: string;
  valor: string;
  tipo: TipoDespesa;
  data_vencimento: string;
  recorrente?: boolean;
  periodicidade?: Periodicidade;
  pix_code?: string;
  observacoes?: string;
}

export interface UpdateContaPagarRequest {
  descricao?: string;
  categoria_id?: string;
  fornecedor?: string;
  valor?: string;
  tipo?: TipoDespesa;
  data_vencimento?: string;
  recorrente?: boolean;
  periodicidade?: Periodicidade;
  pix_code?: string;
  observacoes?: string;
}

export interface MarcarPagamentoRequest {
  data_pagamento: string;
  comprovante_url?: string;
}

// =============================================================================
// CONTA A RECEBER
// =============================================================================

export interface ContaReceber {
  id: string;
  tenant_id: string;
  origem: OrigemReceita;
  descricao_origem: string;
  assinatura_id?: string;
  cliente_id?: string;
  valor: string; // Decimal como string
  valor_pago?: string;
  data_vencimento: string;
  data_recebimento?: string;
  status: StatusContaReceber;
  metodo_pagamento?: string;
  observacoes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateContaReceberRequest {
  origem: OrigemReceita;
  descricao_origem: string;
  assinatura_id?: string;
  cliente_id?: string;
  valor: string;
  data_vencimento: string;
  metodo_pagamento?: string;
  observacoes?: string;
}

export interface UpdateContaReceberRequest {
  origem?: OrigemReceita;
  descricao_origem?: string;
  valor?: string;
  data_vencimento?: string;
  metodo_pagamento?: string;
  observacoes?: string;
}

export interface MarcarRecebimentoRequest {
  data_recebimento: string;
  valor_pago: string;
}

// =============================================================================
// COMPENSAÇÃO BANCÁRIA
// =============================================================================

export interface CompensacaoBancaria {
  id: string;
  tenant_id: string;
  tipo: 'RECEITA' | 'DESPESA';
  receita_id?: string;
  despesa_id?: string;
  data_prevista: string;
  data_compensada?: string;
  valor_previsto: string;
  valor_compensado?: string;
  status: StatusCompensacao;
  observacoes?: string;
  created_at: string;
  updated_at: string;
}

// =============================================================================
// DRE - DEMONSTRATIVO DE RESULTADOS
// =============================================================================

export interface DREMensal {
  id: string;
  tenant_id: string;
  mes_ano: string; // Formato: "2025-11"
  
  // Receitas
  receita_bruta: string;
  deducoes: string;
  receita_liquida: string;
  
  // Custos e Despesas
  custos_servicos: string;
  lucro_bruto: string;
  despesas_operacionais: string;
  despesas_administrativas: string;
  
  // Resultado
  lucro_operacional: string;
  
  // Margens
  margem_bruta_percent: string;
  margem_operacional_percent: string;
  
  created_at: string;
  updated_at: string;
}

// Interface simplificada para exibição
export interface DREDisplay {
  mesAno: string;
  mesAnoFormatado: string; // "Novembro 2025"
  receitaBruta: number;
  deducoes: number;
  receitaLiquida: number;
  custosServicos: number;
  lucroBruto: number;
  despesasOperacionais: number;
  despesasAdministrativas: number;
  lucroOperacional: number;
  margemBrutaPercent: number;
  margemOperacionalPercent: number;
}

// =============================================================================
// FLUXO DE CAIXA DIÁRIO
// =============================================================================

export interface FluxoCaixaDiario {
  id: string;
  tenant_id: string;
  data: string; // ISO date string
  saldo_inicial: string;
  total_entradas: string;
  total_saidas: string;
  saldo_final: string;
  saldo_acumulado: string;
  created_at: string;
  updated_at: string;
}

// Interface simplificada para gráficos
export interface FluxoCaixaDisplay {
  data: string;
  dataFormatada: string; // "27/11"
  saldoInicial: number;
  entradas: number;
  saidas: number;
  saldoFinal: number;
  saldoAcumulado: number;
}

// =============================================================================
// FILTROS
// =============================================================================

export interface ListContasPagarFilters {
  status?: StatusContaPagar;
  tipo?: TipoDespesa;
  data_inicio?: string;
  data_fim?: string;
  fornecedor?: string;
  page?: number;
  page_size?: number;
}

export interface ListContasReceberFilters {
  status?: StatusContaReceber;
  origem?: OrigemReceita;
  data_inicio?: string;
  data_fim?: string;
  cliente_id?: string;
  page?: number;
  page_size?: number;
}

export interface ListCompensacoesFilters {
  status?: StatusCompensacao;
  data_inicio?: string;
  data_fim?: string;
  page?: number;
  page_size?: number;
}

export interface ListFluxoCaixaFilters {
  data_inicio?: string;
  data_fim?: string;
}

export interface ListDREFilters {
  mes_ano_inicio?: string; // "2025-01"
  mes_ano_fim?: string; // "2025-12"
}

// =============================================================================
// DASHBOARD / RESUMOS
// =============================================================================

export interface FinancialSummary {
  totalAPagar: number;
  totalAReceber: number;
  saldoAtual: number;
  lucroOperacionalMes: number;
  
  contasPagarPendentes: number;
  contasReceberPendentes: number;
  contasAtrasadas: number;
  
  variacaoMesAnterior: {
    receita: number; // percentual
    despesa: number;
    lucro: number;
  };
}

export interface ProximosVencimentos {
  tipo: 'PAGAR' | 'RECEBER';
  id: string;
  descricao: string;
  valor: number;
  dataVencimento: string;
  diasParaVencer: number;
}

// =============================================================================
// API RESPONSES (Type Aliases)
// =============================================================================

export type ContaPagarResponse = ContaPagar;
export type ContaReceberResponse = ContaReceber;
export type CompensacaoBancariaResponse = CompensacaoBancaria;
export type FluxoCaixaDiarioResponse = FluxoCaixaDiario;
export type DREMensalResponse = DREMensal;

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Converte valor string para número
 */
export function parseMoneyValue(value: string | undefined): number {
  if (!value) return 0;
  return parseFloat(value) || 0;
}

/**
 * Formata valor para exibição em BRL
 */
export function formatCurrency(value: number | string): string {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(numValue || 0);
}

/**
 * Formata data ISO para DD/MM/YYYY
 */
export function formatDate(dateString: string): string {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return date.toLocaleDateString('pt-BR');
}

/**
 * Formata mês/ano para nome do mês
 */
export function formatMesAno(mesAno: string): string {
  if (!mesAno) return '-';
  const [ano, mes] = mesAno.split('-');
  const meses = [
    'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
    'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
  ];
  return `${meses[parseInt(mes) - 1]} ${ano}`;
}

/**
 * Retorna cor do badge baseado no status
 */
export function getStatusColor(status: string): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (status) {
    case 'PAGO':
    case 'RECEBIDO':
    case 'COMPENSADO':
      return 'default'; // verde/sucesso
    case 'PENDENTE':
    case 'PREVISTO':
      return 'secondary'; // amarelo/warning
    case 'ATRASADO':
    case 'DIVERGENTE':
      return 'destructive'; // vermelho/erro
    case 'CANCELADO':
      return 'outline'; // cinza
    default:
      return 'outline';
  }
}

/**
 * Calcula dias para vencimento
 */
export function getDiasParaVencimento(dataVencimento: string): number {
  const hoje = new Date();
  hoje.setHours(0, 0, 0, 0);
  const vencimento = new Date(dataVencimento);
  vencimento.setHours(0, 0, 0, 0);
  const diffTime = vencimento.getTime() - hoje.getTime();
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}

/**
 * Retorna variante de badge para status
 */
export function getStatusBadgeVariant(status: StatusContaPagar | StatusContaReceber): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (status) {
    case StatusContaPagar.PAGO:
    case StatusContaReceber.RECEBIDO:
      return 'default';
    case StatusContaPagar.PENDENTE:
    case StatusContaReceber.PENDENTE:
    case StatusContaReceber.PARCIAL:
      return 'secondary';
    case StatusContaPagar.ATRASADO:
    case StatusContaReceber.ATRASADO:
      return 'destructive';
    case StatusContaPagar.CANCELADO:
    case StatusContaReceber.CANCELADO:
      return 'outline';
    default:
      return 'outline';
  }
}

/**
 * Retorna label para status de conta a pagar
 */
export function getStatusPagarLabel(status: StatusContaPagar): string {
  switch (status) {
    case StatusContaPagar.PENDENTE:
      return 'Pendente';
    case StatusContaPagar.PAGO:
      return 'Pago';
    case StatusContaPagar.ATRASADO:
      return 'Atrasado';
    case StatusContaPagar.CANCELADO:
      return 'Cancelado';
    default:
      return status;
  }
}

/**
 * Retorna label para status de conta a receber
 */
export function getStatusReceberLabel(status: StatusContaReceber): string {
  switch (status) {
    case StatusContaReceber.PENDENTE:
      return 'Pendente';
    case StatusContaReceber.RECEBIDO:
      return 'Recebido';
    case StatusContaReceber.PARCIAL:
      return 'Parcial';
    case StatusContaReceber.ATRASADO:
      return 'Atrasado';
    case StatusContaReceber.CANCELADO:
      return 'Cancelado';
    default:
      return status;
  }
}

/**
 * Retorna label para tipo de despesa
 */
export function getTipoDespesaLabel(tipo: TipoDespesa): string {
  switch (tipo) {
    case TipoDespesa.FIXA:
      return 'Fixa';
    case TipoDespesa.VARIAVEL:
      return 'Variável';
    case TipoDespesa.EXTRAORDINARIA:
      return 'Extraordinária';
    default:
      return tipo;
  }
}

/**
 * Retorna label para origem de receita
 */
export function getOrigemLabel(origem: OrigemReceita): string {
  switch (origem) {
    case OrigemReceita.SERVICO:
      return 'Serviço';
    case OrigemReceita.PRODUTO:
      return 'Produto';
    case OrigemReceita.ASSINATURA:
      return 'Assinatura';
    case OrigemReceita.OUTRO:
      return 'Outro';
    default:
      return origem;
  }
}

// Campos extras para compatibilidade com DRE extendido
export interface DREMensalExtended extends DREMensal {
  receita_servicos?: string;
  receita_produtos?: string;
  outras_receitas?: string;
  receita_liquida_percent?: string;
  custos_servicos_percent?: string;
  despesas_operacionais_percent?: string;
  despesas_administrativas_percent?: string;
  receitas_financeiras?: string;
  despesas_financeiras?: string;
  resultado_financeiro?: string;
  lucro_antes_ir?: string;
  impostos?: string;
  lucro_liquido?: string;
  margem_liquida_percent?: string;
}

// =============================================================================
// EXTENDED FLUXO CAIXA
// =============================================================================

export interface FluxoCaixaResponse {
  items: FluxoCaixaDiario[];
  total: number;
}

// Para compatibilidade com o campo saldo_dia
export interface FluxoCaixaDiarioExtended extends FluxoCaixaDiario {
  saldo_dia: string;
}

// =============================================================================
// PAINEL MENSAL (DASHBOARD)
// =============================================================================

export interface PainelMensalResponse {
  ano: number;
  mes: number;
  nome_mes: string;

  // Receitas
  receita_realizada: string;
  receita_pendente: string;
  receita_total: string;

  // Despesas
  despesas_fixas: string;
  despesas_variaveis: string;
  despesas_pagas: string;
  despesas_pendentes: string;
  despesas_total: string;

  // Resultados
  lucro_bruto: string;
  lucro_liquido: string;
  margem_liquida: string;

  // Metas
  meta_mensal: string;
  percentual_meta: string;
  diferenca_meta: string;
  status_meta: 'Atingida' | 'Em andamento' | 'Abaixo' | 'Sem meta';

  // Caixa
  saldo_caixa_atual: string;

  // Comparativo
  variacao_mes_anterior: string;
  tendencia_variacao: 'up' | 'down' | 'stable';
}

// =============================================================================
// PROJEÇÕES FINANCEIRAS
// =============================================================================

export interface ProjecaoMensal {
  ano: number;
  mes: number;
  nome_mes: string;
  receita_projetada: string;
  despesas_projetadas: string;
  despesas_fixas: string;
  lucro_projetado: string;
  dias_uteis: number;
  meta_diaria: string;
  confianca: 'Alta' | 'Média' | 'Baixa';
}

export interface ProjecoesResponse {
  projecoes: ProjecaoMensal[];
  media_receita_3_meses: string;
  media_despesas_3_meses: string;
  tendencia_receita: 'Crescente' | 'Estável' | 'Decrescente';
  data_geracao: string;
}

// =============================================================================
// DESPESAS FIXAS
// =============================================================================

export interface DespesaFixa {
  id: string;
  tenant_id?: string;
  unidade_id?: string;
  descricao: string;
  categoria_id?: string;
  fornecedor?: string;
  valor: string; // Dinheiro sempre string
  dia_vencimento: number;
  ativo: boolean;
  observacoes?: string;
  criado_em: string;
  atualizado_em: string;
}

export interface CreateDespesaFixaRequest {
  descricao: string;
  categoria_id?: string;
  fornecedor?: string;
  valor: string;
  dia_vencimento: number;
  unidade_id?: string;
  observacoes?: string;
}

export interface UpdateDespesaFixaRequest {
  descricao: string;
  categoria_id?: string;
  fornecedor?: string;
  valor: string;
  dia_vencimento: number;
  unidade_id?: string;
  observacoes?: string;
}

export interface DespesasFixasListResponse {
  data: DespesaFixa[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface DespesasFixasSummaryResponse {
  total: number;
  total_ativas: number;
  valor_total: string;
}

export interface GerarContasRequest {
  ano?: number;
  mes?: number;
}

export interface GerarContasResponse {
  total_despesas: number;
  contas_criadas: number;
  erros: number;
  detalhes_erros?: string[];
  tempo_execucao_ms: number;
}

export interface ListDespesasFixasFilters {
  ativo?: boolean;
  page?: number;
  page_size?: number;
}
