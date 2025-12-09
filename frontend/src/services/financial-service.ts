/**
 * NEXO - Sistema de Gestão para Barbearias
 * Financial Service
 *
 * Serviço de comunicação com a API financeira do backend.
 * Endpoints mapeados de: backend/internal/infra/http/handler/financial_handler.go
 */

import { api } from '@/lib/axios';
import type {
    CompensacaoBancaria,
    ContaPagar,
    ContaReceber,
    CreateContaPagarRequest,
    CreateContaReceberRequest,
    CreateDespesaFixaRequest,
    DespesaFixa,
    DespesasFixasListResponse,
    DespesasFixasSummaryResponse,
    DREMensal,
    FluxoCaixaDiario,
    GerarContasRequest,
    GerarContasResponse,
    ListCompensacoesFilters,
    ListContasPagarFilters,
    ListContasReceberFilters,
    ListDespesasFixasFilters,
    ListDREFilters,
    ListFluxoCaixaFilters,
    MarcarPagamentoRequest,
    MarcarRecebimentoRequest,
    PainelMensalResponse,
    ProjecoesResponse,
    UpdateContaPagarRequest,
    UpdateContaReceberRequest,
    UpdateDespesaFixaRequest,
} from '@/types/financial';
import {
    StatusContaPagar,
    StatusContaReceber,
} from '@/types/financial';

// =============================================================================
// ENDPOINTS
// =============================================================================

const FINANCIAL_ENDPOINTS = {
  // Contas a Pagar
  payables: '/financial/payables',
  payableById: (id: string) => `/financial/payables/${id}`,
  payablePayment: (id: string) => `/financial/payables/${id}/payment`,

  // Contas a Receber
  receivables: '/financial/receivables',
  receivableById: (id: string) => `/financial/receivables/${id}`,
  receivableReceipt: (id: string) => `/financial/receivables/${id}/receipt`,

  // Despesas Fixas
  fixedExpenses: '/financial/fixed-expenses',
  fixedExpenseById: (id: string) => `/financial/fixed-expenses/${id}`,
  fixedExpenseToggle: (id: string) => `/financial/fixed-expenses/${id}/toggle`,
  fixedExpensesSummary: '/financial/fixed-expenses/summary',
  fixedExpensesGenerate: '/financial/fixed-expenses/generate',

  // Compensações Bancárias
  compensations: '/financial/compensations',
  compensationById: (id: string) => `/financial/compensations/${id}`,

  // Fluxo de Caixa
  cashflow: '/financial/cashflow',
  cashflowById: (id: string) => `/financial/cashflow/${id}`,

  // DRE
  dre: '/financial/dre',
  dreByMonth: (month: string) => `/financial/dre/${month}`,

  // Dashboard e Projeções
  dashboard: '/financial/dashboard',
  projections: '/financial/projections',
} as const;

// =============================================================================
// SERVIÇO PRINCIPAL
// =============================================================================

export const financialService = {
  // ===========================================================================
  // CONTAS A PAGAR
  // ===========================================================================

  /**
   * Lista contas a pagar com filtros opcionais
   */
  async listPayables(filters: ListContasPagarFilters = {}): Promise<ContaPagar[]> {
    console.log('[financial-service] Listando contas a pagar:', filters);
    const { data } = await api.get<ContaPagar[]>(FINANCIAL_ENDPOINTS.payables, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca conta a pagar por ID
   */
  async getPayable(id: string): Promise<ContaPagar> {
    console.log('[financial-service] Buscando conta a pagar:', id);
    const { data } = await api.get<ContaPagar>(FINANCIAL_ENDPOINTS.payableById(id));
    return data;
  },

  /**
   * Cria nova conta a pagar
   */
  async createPayable(payload: CreateContaPagarRequest): Promise<ContaPagar> {
    console.log('[financial-service] Criando conta a pagar:', payload);
    const { data } = await api.post<ContaPagar>(FINANCIAL_ENDPOINTS.payables, payload);
    return data;
  },

  /**
   * Atualiza conta a pagar
   */
  async updatePayable(id: string, payload: UpdateContaPagarRequest): Promise<ContaPagar> {
    console.log('[financial-service] Atualizando conta a pagar:', id, payload);
    const { data } = await api.put<ContaPagar>(FINANCIAL_ENDPOINTS.payableById(id), payload);
    return data;
  },

  /**
   * Deleta conta a pagar
   */
  async deletePayable(id: string): Promise<void> {
    console.log('[financial-service] Deletando conta a pagar:', id);
    await api.delete(FINANCIAL_ENDPOINTS.payableById(id));
  },

  /**
   * Marca conta como paga
   */
  async markPayableAsPaid(id: string, payload: MarcarPagamentoRequest): Promise<void> {
    console.log('[financial-service] Marcando como pago:', id, payload);
    await api.post(FINANCIAL_ENDPOINTS.payablePayment(id), payload);
  },

  /**
   * Cancela conta a pagar
   */
  async cancelPayable(id: string): Promise<void> {
    console.log('[financial-service] Cancelando conta a pagar:', id);
    await api.put(FINANCIAL_ENDPOINTS.payableById(id), { status: 'CANCELADO' });
  },

  // ===========================================================================
  // CONTAS A RECEBER
  // ===========================================================================

  /**
   * Lista contas a receber com filtros opcionais
   */
  async listReceivables(filters: ListContasReceberFilters = {}): Promise<ContaReceber[]> {
    console.log('[financial-service] Listando contas a receber:', filters);
    const { data } = await api.get<ContaReceber[]>(FINANCIAL_ENDPOINTS.receivables, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca conta a receber por ID
   */
  async getReceivable(id: string): Promise<ContaReceber> {
    console.log('[financial-service] Buscando conta a receber:', id);
    const { data } = await api.get<ContaReceber>(FINANCIAL_ENDPOINTS.receivableById(id));
    return data;
  },

  /**
   * Cria nova conta a receber
   */
  async createReceivable(payload: CreateContaReceberRequest): Promise<ContaReceber> {
    console.log('[financial-service] Criando conta a receber:', payload);
    const { data } = await api.post<ContaReceber>(FINANCIAL_ENDPOINTS.receivables, payload);
    return data;
  },

  /**
   * Atualiza conta a receber
   */
  async updateReceivable(id: string, payload: UpdateContaReceberRequest): Promise<ContaReceber> {
    console.log('[financial-service] Atualizando conta a receber:', id, payload);
    const { data } = await api.put<ContaReceber>(FINANCIAL_ENDPOINTS.receivableById(id), payload);
    return data;
  },

  /**
   * Deleta conta a receber
   */
  async deleteReceivable(id: string): Promise<void> {
    console.log('[financial-service] Deletando conta a receber:', id);
    await api.delete(FINANCIAL_ENDPOINTS.receivableById(id));
  },

  /**
   * Marca conta como recebida
   */
  async markReceivableAsReceived(id: string, payload: MarcarRecebimentoRequest): Promise<void> {
    console.log('[financial-service] Marcando como recebido:', id, payload);
    await api.post(FINANCIAL_ENDPOINTS.receivableReceipt(id), payload);
  },

  /**
   * Cancela conta a receber
   */
  async cancelReceivable(id: string): Promise<void> {
    console.log('[financial-service] Cancelando conta a receber:', id);
    await api.put(FINANCIAL_ENDPOINTS.receivableById(id), { status: 'CANCELADO' });
  },

  // ===========================================================================
  // COMPENSAÇÕES BANCÁRIAS
  // ===========================================================================

  /**
   * Lista compensações bancárias
   */
  async listCompensations(filters: ListCompensacoesFilters = {}): Promise<CompensacaoBancaria[]> {
    console.log('[financial-service] Listando compensações:', filters);
    const { data } = await api.get<CompensacaoBancaria[]>(FINANCIAL_ENDPOINTS.compensations, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca compensação por ID
   */
  async getCompensation(id: string): Promise<CompensacaoBancaria> {
    console.log('[financial-service] Buscando compensação:', id);
    const { data } = await api.get<CompensacaoBancaria>(FINANCIAL_ENDPOINTS.compensationById(id));
    return data;
  },

  /**
   * Deleta compensação
   */
  async deleteCompensation(id: string): Promise<void> {
    console.log('[financial-service] Deletando compensação:', id);
    await api.delete(FINANCIAL_ENDPOINTS.compensationById(id));
  },

  // ===========================================================================
  // FLUXO DE CAIXA
  // ===========================================================================

  /**
   * Lista fluxo de caixa por período
   */
  async listCashFlow(filters: ListFluxoCaixaFilters = {}): Promise<FluxoCaixaDiario[]> {
    console.log('[financial-service] Listando fluxo de caixa:', filters);
    const { data } = await api.get<FluxoCaixaDiario[]>(FINANCIAL_ENDPOINTS.cashflow, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca fluxo de caixa por ID/data
   */
  async getCashFlow(id: string): Promise<FluxoCaixaDiario> {
    console.log('[financial-service] Buscando fluxo de caixa:', id);
    const { data } = await api.get<FluxoCaixaDiario>(FINANCIAL_ENDPOINTS.cashflowById(id));
    return data;
  },

  // ===========================================================================
  // DRE - DEMONSTRATIVO DE RESULTADOS
  // ===========================================================================

  /**
   * Lista DREs por período
   */
  async listDRE(filters: ListDREFilters = {}): Promise<DREMensal[]> {
    console.log('[financial-service] Listando DREs:', filters);
    const { data } = await api.get<DREMensal[]>(FINANCIAL_ENDPOINTS.dre, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca DRE de um mês específico
   */
  async getDRE(mesAno: string): Promise<DREMensal> {
    console.log('[financial-service] Buscando DRE:', mesAno);
    const { data } = await api.get<DREMensal>(FINANCIAL_ENDPOINTS.dreByMonth(mesAno));
    return data;
  },

  // ===========================================================================
  // DESPESAS FIXAS
  // ===========================================================================

  /**
   * Lista despesas fixas com filtros opcionais
   */
  async listFixedExpenses(filters: ListDespesasFixasFilters = {}): Promise<DespesasFixasListResponse> {
    console.log('[financial-service] Listando despesas fixas:', filters);
    const { data } = await api.get<DespesasFixasListResponse>(FINANCIAL_ENDPOINTS.fixedExpenses, {
      params: filters,
    });
    return data;
  },

  /**
   * Busca despesa fixa por ID
   */
  async getFixedExpense(id: string): Promise<DespesaFixa> {
    console.log('[financial-service] Buscando despesa fixa:', id);
    const { data } = await api.get<DespesaFixa>(FINANCIAL_ENDPOINTS.fixedExpenseById(id));
    return data;
  },

  /**
   * Cria nova despesa fixa
   */
  async createFixedExpense(payload: CreateDespesaFixaRequest): Promise<DespesaFixa> {
    console.log('[financial-service] Criando despesa fixa:', payload);
    const { data } = await api.post<DespesaFixa>(FINANCIAL_ENDPOINTS.fixedExpenses, payload);
    return data;
  },

  /**
   * Atualiza despesa fixa
   */
  async updateFixedExpense(id: string, payload: UpdateDespesaFixaRequest): Promise<DespesaFixa> {
    console.log('[financial-service] Atualizando despesa fixa:', id, payload);
    const { data } = await api.put<DespesaFixa>(FINANCIAL_ENDPOINTS.fixedExpenseById(id), payload);
    return data;
  },

  /**
   * Deleta despesa fixa
   */
  async deleteFixedExpense(id: string): Promise<void> {
    console.log('[financial-service] Deletando despesa fixa:', id);
    await api.delete(FINANCIAL_ENDPOINTS.fixedExpenseById(id));
  },

  /**
   * Ativa/desativa despesa fixa
   */
  async toggleFixedExpense(id: string): Promise<DespesaFixa> {
    console.log('[financial-service] Toggling despesa fixa:', id);
    const { data } = await api.patch<DespesaFixa>(FINANCIAL_ENDPOINTS.fixedExpenseToggle(id));
    return data;
  },

  /**
   * Busca resumo das despesas fixas
   */
  async getFixedExpensesSummary(): Promise<DespesasFixasSummaryResponse> {
    console.log('[financial-service] Buscando resumo despesas fixas');
    const { data } = await api.get<DespesasFixasSummaryResponse>(FINANCIAL_ENDPOINTS.fixedExpensesSummary);
    return data;
  },

  /**
   * Gera contas a pagar a partir das despesas fixas
   */
  async generatePayablesFromFixed(payload: GerarContasRequest = {}): Promise<GerarContasResponse> {
    console.log('[financial-service] Gerando contas a partir de despesas fixas:', payload);
    const { data } = await api.post<GerarContasResponse>(FINANCIAL_ENDPOINTS.fixedExpensesGenerate, payload);
    return data;
  },

  // ===========================================================================
  // DASHBOARD - PAINEL MENSAL
  // ===========================================================================

  /**
   * Busca dados do painel mensal consolidado
   */
  async getDashboard(year?: number, month?: number): Promise<PainelMensalResponse> {
    console.log('[financial-service] Buscando dashboard:', { year, month });
    const { data } = await api.get<PainelMensalResponse>(FINANCIAL_ENDPOINTS.dashboard, {
      params: { year, month },
    });
    return data;
  },

  /**
   * Busca projeções financeiras
   */
  async getProjections(monthsAhead: number = 3): Promise<ProjecoesResponse> {
    console.log('[financial-service] Buscando projeções:', { monthsAhead });
    const { data } = await api.get<ProjecoesResponse>(FINANCIAL_ENDPOINTS.projections, {
      params: { months_ahead: monthsAhead },
    });
    return data;
  },

  // ===========================================================================
  // HELPERS / AGREGADORES (Calculados no Frontend)
  // ===========================================================================

  /**
   * Calcula resumo financeiro agregando dados de múltiplos endpoints
   */
  async getFinancialSummary(): Promise<{
    totalAPagar: number;
    totalAReceber: number;
    contasPagarPendentes: number;
    contasReceberPendentes: number;
  }> {
    console.log('[financial-service] Calculando resumo financeiro');

    // Busca paralela para performance
    const [payables, receivables] = await Promise.all([
      this.listPayables({ status: StatusContaPagar.PENDENTE }),
      this.listReceivables({ status: StatusContaReceber.PENDENTE }),
    ]);

    const totalAPagar = payables.reduce((sum, p) => sum + parseFloat(p.valor || '0'), 0);
    const totalAReceber = receivables.reduce((sum, r) => sum + parseFloat(r.valor || '0'), 0);

    return {
      totalAPagar,
      totalAReceber,
      contasPagarPendentes: payables.length,
      contasReceberPendentes: receivables.length,
    };
  },

  /**
   * Busca próximos vencimentos (pagar + receber)
   */
  async getProximosVencimentos(limit: number = 5): Promise<Array<{
    tipo: 'PAGAR' | 'RECEBER';
    id: string;
    descricao: string;
    valor: number;
    dataVencimento: string;
  }>> {
    console.log('[financial-service] Buscando próximos vencimentos');

    const hoje = new Date().toISOString().split('T')[0];
    const futuro = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];

    const [payables, receivables] = await Promise.all([
      this.listPayables({ status: StatusContaPagar.PENDENTE, data_inicio: hoje, data_fim: futuro }),
      this.listReceivables({ status: StatusContaReceber.PENDENTE, data_inicio: hoje, data_fim: futuro }),
    ]);

    const vencimentos = [
      ...payables.map((p) => ({
        tipo: 'PAGAR' as const,
        id: p.id,
        descricao: p.descricao,
        valor: parseFloat(p.valor || '0'),
        dataVencimento: p.data_vencimento,
      })),
      ...receivables.map((r) => ({
        tipo: 'RECEBER' as const,
        id: r.id,
        descricao: r.descricao_origem,
        valor: parseFloat(r.valor || '0'),
        dataVencimento: r.data_vencimento,
      })),
    ];

    // Ordena por data de vencimento
    vencimentos.sort((a, b) => 
      new Date(a.dataVencimento).getTime() - new Date(b.dataVencimento).getTime()
    );

    return vencimentos.slice(0, limit);
  },
};

export default financialService;
