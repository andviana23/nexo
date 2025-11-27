/**
 * NEXO - Sistema de Gestão para Barbearias
 * Barber Turn Service (Lista da Vez)
 *
 * Serviço de Lista da Vez - comunicação com API de barber-turn do backend.
 * Seguindo padrões do projeto e FLUXO_LISTA_DA_VEZ.md
 *
 * Sistema de fila giratória para distribuição justa de clientes.
 * Menor pontuação = próximo da fila (menor número de atendimentos no mês).
 */

import { api } from '@/lib/axios';
import type {
    AddBarberToTurnListRequest,
    BarberTurnResponse,
    ListAvailableBarbersResponse,
    ListBarbersTurnResponse,
    ListHistorySummaryResponse,
    ListTurnHistoryResponse,
    RecordTurnRequest,
    RecordTurnResponse,
    RemoveBarberResponse,
    ResetTurnListRequest,
    ResetTurnListResponse,
    ToggleStatusResponse,
} from '@/types/barber-turn';

// =============================================================================
// ENDPOINTS
// =============================================================================

const BARBER_TURN_ENDPOINTS = {
  list: '/barber-turn/list',
  add: '/barber-turn/add',
  record: '/barber-turn/record',
  toggleStatus: (professionalId: string) => `/barber-turn/${professionalId}/toggle-status`,
  remove: (professionalId: string) => `/barber-turn/${professionalId}`,
  reset: '/barber-turn/reset',
  history: '/barber-turn/history',
  historySummary: '/barber-turn/history/summary',
  available: '/barber-turn/available',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const barberTurnService = {
  /**
   * Lista todos os barbeiros na fila da vez
   * Ordenados por pontuação (menor primeiro = próximo da fila)
   *
   * @param isActive - Filtrar por status ativo/inativo
   * @returns Lista de barbeiros com estatísticas e próximo da fila
   */
  async list(isActive?: boolean): Promise<ListBarbersTurnResponse> {
    const params: Record<string, unknown> = {};
    if (isActive !== undefined) {
      params.is_active = isActive;
    }

    console.log('[barberTurnService.list] Chamando API com params:', params);
    const response = await api.get<ListBarbersTurnResponse>(
      BARBER_TURN_ENDPOINTS.list,
      { params }
    );
    console.log('[barberTurnService.list] Resposta completa:', {
      data: response.data,
      barbers: response.data?.barbers,
      barbersLength: response.data?.barbers?.length,
      stats: response.data?.stats
    });
    return response.data;
  },

  /**
   * Adiciona um barbeiro à Lista da Vez
   * Apenas profissionais do tipo BARBEIRO podem ser adicionados
   *
   * @param professionalId - ID do profissional
   * @returns Dados do barbeiro adicionado
   */
  async addBarber(professionalId: string): Promise<BarberTurnResponse> {
    const request: AddBarberToTurnListRequest = {
      professional_id: professionalId,
    };

    const response = await api.post<BarberTurnResponse>(
      BARBER_TURN_ENDPOINTS.add,
      request
    );
    return response.data;
  },

  /**
   * Registra um atendimento para um barbeiro
   * Incrementa os pontos (+1) e reordena a fila
   *
   * @param professionalId - ID do profissional
   * @returns Dados do atendimento registrado
   */
  async recordTurn(professionalId: string): Promise<RecordTurnResponse> {
    const request: RecordTurnRequest = {
      professional_id: professionalId,
    };

    const response = await api.post<RecordTurnResponse>(
      BARBER_TURN_ENDPOINTS.record,
      request
    );
    return response.data;
  },

  /**
   * Alterna o status ativo/inativo de um barbeiro
   * Barbeiros pausados não aparecem na fila
   *
   * @param professionalId - ID do profissional
   * @returns Novo status do barbeiro
   */
  async toggleStatus(professionalId: string): Promise<ToggleStatusResponse> {
    const response = await api.put<ToggleStatusResponse>(
      BARBER_TURN_ENDPOINTS.toggleStatus(professionalId)
    );
    return response.data;
  },

  /**
   * Remove um barbeiro da Lista da Vez
   * ATENÇÃO: Não preserva pontos - barbeiro perde posição na fila
   *
   * @param professionalId - ID do profissional
   */
  async removeBarber(professionalId: string): Promise<RemoveBarberResponse> {
    const response = await api.delete<RemoveBarberResponse>(
      BARBER_TURN_ENDPOINTS.remove(professionalId)
    );
    return response.data;
  },

  /**
   * Executa reset mensal da Lista da Vez
   * Zera todos os pontos e opcionalmente salva histórico
   *
   * @param saveHistory - Se true, salva histórico antes do reset
   * @returns Resumo do reset executado
   */
  async resetMonthly(saveHistory: boolean = true): Promise<ResetTurnListResponse> {
    const request: ResetTurnListRequest = {
      save_history: saveHistory,
    };

    const response = await api.post<ResetTurnListResponse>(
      BARBER_TURN_ENDPOINTS.reset,
      request
    );
    return response.data;
  },

  /**
   * Lista histórico de atendimentos por mês
   *
   * @param monthYear - Mês/ano no formato YYYY-MM (opcional)
   * @returns Lista de histórico
   */
  async getHistory(monthYear?: string): Promise<ListTurnHistoryResponse> {
    const params: Record<string, unknown> = {};
    if (monthYear) {
      params.month_year = monthYear;
    }

    const response = await api.get<ListTurnHistoryResponse>(
      BARBER_TURN_ENDPOINTS.history,
      { params }
    );
    return response.data;
  },

  /**
   * Lista resumo dos últimos 12 meses
   *
   * @returns Resumo mensal com totais e médias
   */
  async getHistorySummary(): Promise<ListHistorySummaryResponse> {
    const response = await api.get<ListHistorySummaryResponse>(
      BARBER_TURN_ENDPOINTS.historySummary
    );
    return response.data;
  },

  /**
   * Lista barbeiros disponíveis para adicionar à fila
   * Retorna apenas profissionais do tipo BARBEIRO que ainda não estão na fila
   *
   * @returns Lista de barbeiros disponíveis
   */
  async getAvailableBarbers(): Promise<ListAvailableBarbersResponse> {
    console.log('[barberTurnService] Fetching available barbers...');
    const response = await api.get<ListAvailableBarbersResponse>(
      BARBER_TURN_ENDPOINTS.available
    );
    console.log('[barberTurnService] Available barbers response:', response.data);
    return response.data;
  },

  // ===========================================================================
  // Métodos auxiliares
  // ===========================================================================

  /**
   * Obtém o próximo barbeiro da fila
   * @returns Próximo barbeiro ou null se fila vazia
   */
  async getNextBarber(): Promise<BarberTurnResponse | null> {
    const response = await this.list(true);
    return response.barbers[0] || null;
  },

  /**
   * Verifica se um profissional está na lista
   * @param professionalId - ID do profissional
   * @returns true se está na lista
   */
  async isInList(professionalId: string): Promise<boolean> {
    const response = await this.list();
    return response.barbers.some(b => b.professional_id === professionalId);
  },

  /**
   * Obtém estatísticas da fila
   * @returns Estatísticas gerais
   */
  async getStats(): Promise<{
    totalAtivos: number;
    totalPausados: number;
    totalGeral: number;
    totalPontosMes: number;
  }> {
    const response = await this.list();
    return {
      totalAtivos: response.stats.total_ativos,
      totalPausados: response.stats.total_pausados,
      totalGeral: response.stats.total_geral,
      totalPontosMes: response.stats.total_pontos_mes,
    };
  },
};

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Calcula posição na fila baseado nos pontos
 * Menor pontuação = mais próximo de atender
 */
export function calculatePosition(barbers: BarberTurnResponse[]): BarberTurnResponse[] {
  return [...barbers]
    .sort((a, b) => a.current_points - b.current_points)
    .map((b, index) => ({ ...b, position: index + 1 }));
}

/**
 * Formata pontos para exibição
 * Ex: "5 atendimentos" ou "1 atendimento"
 */
export function formatPoints(points: number): string {
  if (points === 0) return 'Nenhum atendimento';
  if (points === 1) return '1 atendimento';
  return `${points} atendimentos`;
}

/**
 * Formata data do último atendimento
 */
export function formatLastTurn(lastTurnAt: string | undefined): string {
  if (!lastTurnAt) return 'Sem atendimentos';

  const date = new Date(lastTurnAt);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) return 'Agora mesmo';
  if (diffMins < 60) return `Há ${diffMins} min`;
  if (diffHours < 24) return `Há ${diffHours}h`;
  if (diffDays === 1) return 'Ontem';
  if (diffDays < 7) return `Há ${diffDays} dias`;

  return date.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
  });
}

/**
 * Formata mês/ano para exibição
 * Ex: "2024-01" -> "Janeiro 2024"
 */
export function formatMonthYear(monthYear: string): string {
  const [year, month] = monthYear.split('-');
  const months = [
    'Janeiro', 'Fevereiro', 'Março', 'Abril',
    'Maio', 'Junho', 'Julho', 'Agosto',
    'Setembro', 'Outubro', 'Novembro', 'Dezembro',
  ];
  return `${months[parseInt(month) - 1]} ${year}`;
}

/**
 * Obtém mês/ano atual no formato YYYY-MM
 */
export function getCurrentMonthYear(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  return `${year}-${month}`;
}

/**
 * Obtém lista de meses para seleção (últimos 12 meses)
 */
export function getMonthOptions(): { value: string; label: string }[] {
  const options: { value: string; label: string }[] = [];
  const now = new Date();

  for (let i = 0; i < 12; i++) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const value = `${year}-${month}`;
    options.push({
      value,
      label: formatMonthYear(value),
    });
  }

  return options;
}

/**
 * Calcula média de atendimentos
 */
export function calculateAveragePoints(barbers: BarberTurnResponse[]): number {
  if (barbers.length === 0) return 0;
  const total = barbers.reduce((sum, b) => sum + b.current_points, 0);
  return Math.round((total / barbers.length) * 10) / 10;
}

/**
 * Obtém status do barbeiro para exibição
 */
export function getBarberStatus(barber: BarberTurnResponse): {
  label: string;
  color: 'success' | 'warning' | 'error';
} {
  if (!barber.is_active) {
    return { label: 'Pausado', color: 'warning' };
  }
  if (barber.position === 1) {
    return { label: 'Próximo', color: 'success' };
  }
  return { label: 'Na fila', color: 'success' };
}

/**
 * Obtém iniciais do nome para avatar
 */
export function getInitials(name: string): string {
  const parts = name.trim().split(' ').filter(Boolean);
  if (parts.length === 0) return '??';
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
  return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
}

/**
 * Gera cor de avatar baseada no nome (determinística)
 */
export function getAvatarColor(name: string): string {
  const colors = [
    'bg-blue-500',
    'bg-green-500',
    'bg-yellow-500',
    'bg-red-500',
    'bg-purple-500',
    'bg-pink-500',
    'bg-indigo-500',
    'bg-teal-500',
  ];

  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
}

// =============================================================================
// TIPOS DE ERRO ESPECÍFICOS
// =============================================================================

export class BarberTurnError extends Error {
  constructor(
    message: string,
    public code: string = 'BARBER_TURN_ERROR'
  ) {
    super(message);
    this.name = 'BarberTurnError';
  }
}

export class BarberNotFoundError extends BarberTurnError {
  constructor() {
    super('Barbeiro não encontrado na lista.', 'BARBER_NOT_FOUND');
    this.name = 'BarberNotFoundError';
  }
}

export class BarberAlreadyInListError extends BarberTurnError {
  constructor() {
    super('Este barbeiro já está na Lista da Vez.', 'ALREADY_IN_LIST');
    this.name = 'BarberAlreadyInListError';
  }
}

export class NotBarberError extends BarberTurnError {
  constructor() {
    super('Este profissional não é do tipo Barbeiro.', 'NOT_BARBER');
    this.name = 'NotBarberError';
  }
}

export class BarberPausedError extends BarberTurnError {
  constructor() {
    super('Não é possível registrar atendimento para barbeiro pausado.', 'BARBER_PAUSED');
    this.name = 'BarberPausedError';
  }
}

export class ResetFailedError extends BarberTurnError {
  constructor() {
    super('Erro ao executar reset mensal.', 'RESET_FAILED');
    this.name = 'ResetFailedError';
  }
}
