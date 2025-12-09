/**
 * NEXO - Sistema de Gestão para Barbearias
 * Metas Types
 *
 * Tipos TypeScript para o módulo de Metas.
 * Mapeados a partir do backend Go (metas_dto.go)
 */

// =============================================================================
// ENUMS
// =============================================================================

export enum OrigemMeta {
  MANUAL = 'MANUAL',
  AUTOMATICA = 'AUTOMATICA',
}

export enum StatusMeta {
  PENDENTE = 'PENDENTE',
  ACEITA = 'ACEITA',
  REJEITADA = 'REJEITADA',
}

export enum TipoTicketMeta {
  GERAL = 'GERAL',
  BARBEIRO = 'BARBEIRO',
}

export enum NivelBonificacao {
  NENHUM = 'NENHUM',
  NIVEL_1 = 'NIVEL_1', // >= 100% → 3%
  NIVEL_2 = 'NIVEL_2', // >= 110% → 5%
  NIVEL_3 = 'NIVEL_3', // >= 120% → 8%
}

// =============================================================================
// METAS MENSAIS (Faturamento Geral)
// =============================================================================

export interface MetaMensalResponse {
  id: string;
  mes_ano: string;
  meta_faturamento: string;
  origem: string;
  status: string;
  realizado: string;
  percentual: string;
  criado_em: string;
  atualizado_em: string;
}

export interface SetMetaMensalRequest {
  mes_ano: string;
  meta_faturamento: string;
  origem: string;
}

export interface UpdateMetaMensalRequest {
  meta_faturamento?: string;
  origem?: string;
}

// =============================================================================
// METAS POR BARBEIRO (Individuais)
// =============================================================================

export interface MetaBarbeiroResponse {
  id: string;
  barbeiro_id: string;
  barbeiro_nome?: string;
  mes_ano: string;
  meta_servicos_gerais: string;
  meta_servicos_extras: string;
  meta_produtos: string;
  realizado_servicos_gerais: string;
  realizado_servicos_extras: string;
  realizado_produtos: string;
  percentual_servicos_gerais: string;
  percentual_servicos_extras: string;
  percentual_produtos: string;
  criado_em: string;
  atualizado_em: string;
}

export interface SetMetaBarbeiroRequest {
  barbeiro_id: string;
  mes_ano: string;
  meta_servicos_gerais: string;
  meta_servicos_extras: string;
  meta_produtos: string;
}

export interface UpdateMetaBarbeiroRequest {
  meta_servicos_gerais?: string;
  meta_servicos_extras?: string;
  meta_produtos?: string;
}

// =============================================================================
// METAS TICKET MÉDIO
// =============================================================================

export interface MetaTicketResponse {
  id: string;
  mes_ano: string;
  tipo: string;
  barbeiro_id?: string;
  barbeiro_nome?: string;
  meta_valor: string;
  ticket_medio_realizado: string;
  percentual: string;
  criado_em: string;
  atualizado_em: string;
}

export interface SetMetaTicketRequest {
  mes_ano: string;
  tipo: string;
  barbeiro_id?: string;
  meta_valor: string;
}

export interface UpdateMetaTicketRequest {
  meta_valor?: string;
}

// =============================================================================
// TIPOS DERIVADOS / AUXILIARES
// =============================================================================

export interface MetaBarbeiroExtended extends MetaBarbeiroResponse {
  meta_total: number;
  realizado_total: number;
  percentual_total: number;
  nivel_bonificacao: NivelBonificacao;
  bonus_percentual: number;
  bonus_valor?: number;
}

export interface ResumoMetas {
  mes_ano: string;
  meta_faturamento: number;
  realizado_faturamento: number;
  percentual_faturamento: number;
  total_barbeiros_com_meta: number;
  barbeiros_acima_meta: number;
  ticket_medio_meta: number;
  ticket_medio_realizado: number;
  percentual_ticket: number;
}

export interface MetasFilters {
  mes_ano?: string;
  ano?: number;
  barbeiro_id?: string;
  status?: StatusMeta;
}

export interface RankingBarbeiro {
  posicao: number;
  barbeiro_id: string;
  barbeiro_nome: string;
  percentual_total: number;
  nivel_bonificacao: NivelBonificacao;
  variacao_mes_anterior?: number;
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Formata um valor percentual para exibição
 */
export function formatPercentual(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return '0%';
  return `${num.toFixed(1)}%`;
}

/**
 * Retorna classe de cor para o percentual
 */
export function getPercentualColor(percentual: number): string {
  if (percentual >= 120) return 'text-purple-600';
  if (percentual >= 100) return 'text-green-600';
  if (percentual >= 80) return 'text-yellow-600';
  if (percentual >= 50) return 'text-orange-500';
  return 'text-red-600';
}

/**
 * Retorna classe de cor de fundo para progress bar
 */
export function getProgressColor(percentual: number): string {
  if (percentual >= 120) return 'bg-purple-500';
  if (percentual >= 100) return 'bg-green-500';
  if (percentual >= 80) return 'bg-yellow-500';
  if (percentual >= 50) return 'bg-orange-500';
  return 'bg-red-500';
}

/**
 * Calcula o nível de bonificação baseado no percentual atingido
 */
export function calcularNivelBonificacao(percentual: number): NivelBonificacao {
  if (percentual >= 120) return NivelBonificacao.NIVEL_3;
  if (percentual >= 110) return NivelBonificacao.NIVEL_2;
  if (percentual >= 100) return NivelBonificacao.NIVEL_1;
  return NivelBonificacao.NENHUM;
}

/**
 * Retorna o percentual de bônus para cada nível
 */
export function getBonusPercentual(nivel: NivelBonificacao): number {
  switch (nivel) {
    case NivelBonificacao.NIVEL_3: return 8;
    case NivelBonificacao.NIVEL_2: return 5;
    case NivelBonificacao.NIVEL_1: return 3;
    default: return 0;
  }
}

/**
 * Formata o nível de bonificação para exibição
 */
export function formatNivelBonificacao(nivel: NivelBonificacao): string {
  switch (nivel) {
    case NivelBonificacao.NIVEL_3: return 'Nível 3 (+8%)';
    case NivelBonificacao.NIVEL_2: return 'Nível 2 (+5%)';
    case NivelBonificacao.NIVEL_1: return 'Nível 1 (+3%)';
    default: return 'Sem bônus';
  }
}

/**
 * Retorna cor do badge de nível de bonificação
 */
export function getNivelBonificacaoColor(nivel: NivelBonificacao): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (nivel) {
    case NivelBonificacao.NIVEL_3: return 'default'; // purple/gold
    case NivelBonificacao.NIVEL_2: return 'default'; // green
    case NivelBonificacao.NIVEL_1: return 'secondary'; // blue
    default: return 'outline'; // gray
  }
}

/**
 * Retorna classe CSS para badge de nível
 */
export function getNivelBonificacaoClass(nivel: NivelBonificacao): string {
  switch (nivel) {
    case NivelBonificacao.NIVEL_3: return 'bg-purple-100 text-purple-800 border-purple-200';
    case NivelBonificacao.NIVEL_2: return 'bg-green-100 text-green-800 border-green-200';
    case NivelBonificacao.NIVEL_1: return 'bg-blue-100 text-blue-800 border-blue-200';
    default: return 'bg-gray-100 text-gray-600 border-gray-200';
  }
}

/**
 * Formata mês/ano para exibição
 */
export function formatMesAno(mesAno: string): string {
  if (!mesAno) return '-';
  const [ano, mes] = mesAno.split('-');
  const meses = [
    'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
    'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
  ];
  const mesIndex = parseInt(mes) - 1;
  if (mesIndex < 0 || mesIndex > 11) return mesAno;
  return `${meses[mesIndex]} ${ano}`;
}

/**
 * Formata mês/ano para exibição curta
 */
export function formatMesAnoShort(mesAno: string): string {
  if (!mesAno) return '-';
  const [ano, mes] = mesAno.split('-');
  const meses = ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'];
  const mesIndex = parseInt(mes) - 1;
  if (mesIndex < 0 || mesIndex > 11) return mesAno;
  return `${meses[mesIndex]}/${ano.slice(2)}`;
}

/**
 * Converte string de valor para número
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
 * Calcula percentual de atingimento
 */
export function calcularPercentual(realizado: number, meta: number): number {
  if (meta <= 0) return 0;
  return (realizado / meta) * 100;
}

/**
 * Retorna status label para status de meta
 */
export function getStatusMetaLabel(status: string): string {
  switch (status) {
    case StatusMeta.PENDENTE: return 'Pendente';
    case StatusMeta.ACEITA: return 'Aceita';
    case StatusMeta.REJEITADA: return 'Rejeitada';
    default: return status;
  }
}

/**
 * Retorna variante do badge para status
 */
export function getStatusMetaBadgeVariant(status: string): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (status) {
    case StatusMeta.ACEITA: return 'default';
    case StatusMeta.PENDENTE: return 'secondary';
    case StatusMeta.REJEITADA: return 'destructive';
    default: return 'outline';
  }
}

/**
 * Retorna label para origem da meta
 */
export function getOrigemMetaLabel(origem: string): string {
  switch (origem) {
    case OrigemMeta.MANUAL: return 'Manual';
    case OrigemMeta.AUTOMATICA: return 'Automática';
    default: return origem;
  }
}

/**
 * Retorna label para tipo de ticket
 */
export function getTipoTicketLabel(tipo: string): string {
  switch (tipo) {
    case TipoTicketMeta.GERAL: return 'Geral';
    case TipoTicketMeta.BARBEIRO: return 'Por Barbeiro';
    default: return tipo;
  }
}

/**
 * Gera o mês/ano atual no formato YYYY-MM
 */
export function getMesAnoAtual(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  return `${year}-${month}`;
}

/**
 * Gera lista de meses/anos para seleção (últimos 12 meses + próximos 6)
 */
export function gerarOpcoesMesAno(): { value: string; label: string }[] {
  const opcoes: { value: string; label: string }[] = [];
  const now = new Date();
  
  // Últimos 12 meses
  for (let i = 11; i >= 0; i--) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1);
    const value = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
    opcoes.push({ value, label: formatMesAno(value) });
  }
  
  // Próximos 6 meses
  for (let i = 1; i <= 6; i++) {
    const date = new Date(now.getFullYear(), now.getMonth() + i, 1);
    const value = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
    opcoes.push({ value, label: formatMesAno(value) });
  }
  
  return opcoes;
}

/**
 * Calcula meta total do barbeiro
 */
export function calcularMetaTotalBarbeiro(meta: MetaBarbeiroResponse): number {
  return (
    parseMoneyValue(meta.meta_servicos_gerais) +
    parseMoneyValue(meta.meta_servicos_extras) +
    parseMoneyValue(meta.meta_produtos)
  );
}

/**
 * Calcula realizado total do barbeiro
 */
export function calcularRealizadoTotalBarbeiro(meta: MetaBarbeiroResponse): number {
  return (
    parseMoneyValue(meta.realizado_servicos_gerais) +
    parseMoneyValue(meta.realizado_servicos_extras) +
    parseMoneyValue(meta.realizado_produtos)
  );
}

/**
 * Estende a resposta de meta de barbeiro com cálculos
 */
export function extenderMetaBarbeiro(meta: MetaBarbeiroResponse): MetaBarbeiroExtended {
  const metaTotal = calcularMetaTotalBarbeiro(meta);
  const realizadoTotal = calcularRealizadoTotalBarbeiro(meta);
  const percentualTotal = calcularPercentual(realizadoTotal, metaTotal);
  const nivelBonificacao = calcularNivelBonificacao(percentualTotal);
  const bonusPercentual = getBonusPercentual(nivelBonificacao);
  
  return {
    ...meta,
    meta_total: metaTotal,
    realizado_total: realizadoTotal,
    percentual_total: percentualTotal,
    nivel_bonificacao: nivelBonificacao,
    bonus_percentual: bonusPercentual,
    bonus_valor: realizadoTotal * (bonusPercentual / 100),
  };
}
