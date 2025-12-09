/**
 * NEXO - Sistema de Gestão para Barbearias
 * Financial Hooks
 *
 * Hooks React Query para gerenciar estado do módulo financeiro.
 */

import { financialService } from '@/services/financial-service';
import type {
    CreateContaPagarRequest,
    CreateContaReceberRequest,
    CreateDespesaFixaRequest,
    GerarContasRequest,
    ListCompensacoesFilters,
    ListContasPagarFilters,
    ListContasReceberFilters,
    ListDespesasFixasFilters,
    ListDREFilters,
    ListFluxoCaixaFilters,
    MarcarPagamentoRequest,
    MarcarRecebimentoRequest,
    UpdateContaPagarRequest,
    UpdateContaReceberRequest,
    UpdateDespesaFixaRequest,
} from '@/types/financial';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const financialKeys = {
  all: ['financial'] as const,
  
  // Contas a Pagar
  payables: () => [...financialKeys.all, 'payables'] as const,
  payablesList: (filters: ListContasPagarFilters) => [...financialKeys.payables(), 'list', filters] as const,
  payable: (id: string) => [...financialKeys.payables(), 'detail', id] as const,
  
  // Contas a Receber
  receivables: () => [...financialKeys.all, 'receivables'] as const,
  receivablesList: (filters: ListContasReceberFilters) => [...financialKeys.receivables(), 'list', filters] as const,
  receivable: (id: string) => [...financialKeys.receivables(), 'detail', id] as const,
  
  // Compensações
  compensations: () => [...financialKeys.all, 'compensations'] as const,
  compensationsList: (filters: ListCompensacoesFilters) => [...financialKeys.compensations(), 'list', filters] as const,
  
  // Fluxo de Caixa
  cashflow: () => [...financialKeys.all, 'cashflow'] as const,
  cashflowList: (filters: ListFluxoCaixaFilters) => [...financialKeys.cashflow(), 'list', filters] as const,
  
  // DRE
  dre: () => [...financialKeys.all, 'dre'] as const,
  dreList: (filters: ListDREFilters) => [...financialKeys.dre(), 'list', filters] as const,
  dreByMonth: (mesAno: string) => [...financialKeys.dre(), 'month', mesAno] as const,
  
  // Dashboard / Resumos
  summary: () => [...financialKeys.all, 'summary'] as const,
  vencimentos: () => [...financialKeys.all, 'vencimentos'] as const,
  dashboard: (year?: number, month?: number) => [...financialKeys.all, 'dashboard', year, month] as const,
  projections: (months?: number) => [...financialKeys.all, 'projections', months] as const,
  
  // Despesas Fixas
  fixedExpenses: () => [...financialKeys.all, 'fixed-expenses'] as const,
  fixedExpensesList: (filters: ListDespesasFixasFilters) => [...financialKeys.fixedExpenses(), 'list', filters] as const,
  fixedExpense: (id: string) => [...financialKeys.fixedExpenses(), 'detail', id] as const,
  fixedExpensesSummary: () => [...financialKeys.fixedExpenses(), 'summary'] as const,
};

// =============================================================================
// CONTAS A PAGAR - QUERIES
// =============================================================================

/**
 * Hook para listar contas a pagar
 */
export function usePayables(filters: ListContasPagarFilters = {}) {
  return useQuery({
    queryKey: financialKeys.payablesList(filters),
    queryFn: () => financialService.listPayables(filters),
  });
}

/**
 * Hook para buscar conta a pagar por ID
 */
export function usePayable(id: string) {
  return useQuery({
    queryKey: financialKeys.payable(id),
    queryFn: () => financialService.getPayable(id),
    enabled: !!id,
  });
}

// =============================================================================
// CONTAS A PAGAR - MUTATIONS
// =============================================================================

/**
 * Hook para criar conta a pagar
 */
export function useCreatePayable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateContaPagarRequest) => financialService.createPayable(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      toast.success('Conta a pagar criada com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useCreatePayable] Erro:', error);
      toast.error('Erro ao criar conta a pagar');
    },
  });
}

/**
 * Hook para atualizar conta a pagar
 */
export function useUpdatePayable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateContaPagarRequest }) =>
      financialService.updatePayable(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.payable(id) });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      toast.success('Conta a pagar atualizada!');
    },
    onError: (error: Error) => {
      console.error('[useUpdatePayable] Erro:', error);
      toast.error('Erro ao atualizar conta a pagar');
    },
  });
}

/**
 * Hook para deletar conta a pagar
 */
export function useDeletePayable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.deletePayable(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      toast.success('Conta a pagar excluída!');
    },
    onError: (error: Error) => {
      console.error('[useDeletePayable] Erro:', error);
      toast.error('Erro ao excluir conta a pagar');
    },
  });
}

/**
 * Hook para marcar conta como paga
 */
export function useMarkPayableAsPaid() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: MarcarPagamentoRequest }) =>
      financialService.markPayableAsPaid(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.payable(id) });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      queryClient.invalidateQueries({ queryKey: financialKeys.cashflow() });
      toast.success('Pagamento registrado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useMarkPayableAsPaid] Erro:', error);
      toast.error('Erro ao registrar pagamento');
    },
  });
}

/**
 * Alias para useMarkPayableAsPaid (compatibilidade)
 */
export function usePayPayable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data_pagamento }: { id: string; data_pagamento: string }) =>
      financialService.markPayableAsPaid(id, { data_pagamento }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.payable(id) });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      queryClient.invalidateQueries({ queryKey: financialKeys.cashflow() });
      toast.success('Pagamento registrado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[usePayPayable] Erro:', error);
      toast.error('Erro ao registrar pagamento');
    },
  });
}

/**
 * Hook para cancelar conta a pagar
 */
export function useCancelPayable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.cancelPayable(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      toast.success('Conta cancelada!');
    },
    onError: (error: Error) => {
      console.error('[useCancelPayable] Erro:', error);
      toast.error('Erro ao cancelar conta');
    },
  });
}

// =============================================================================
// CONTAS A RECEBER - QUERIES
// =============================================================================

/**
 * Hook para listar contas a receber
 */
export function useReceivables(filters: ListContasReceberFilters = {}) {
  return useQuery({
    queryKey: financialKeys.receivablesList(filters),
    queryFn: () => financialService.listReceivables(filters),
  });
}

/**
 * Hook para buscar conta a receber por ID
 */
export function useReceivable(id: string) {
  return useQuery({
    queryKey: financialKeys.receivable(id),
    queryFn: () => financialService.getReceivable(id),
    enabled: !!id,
  });
}

// =============================================================================
// CONTAS A RECEBER - MUTATIONS
// =============================================================================

/**
 * Hook para criar conta a receber
 */
export function useCreateReceivable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateContaReceberRequest) => financialService.createReceivable(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.receivables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      toast.success('Conta a receber criada com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useCreateReceivable] Erro:', error);
      toast.error('Erro ao criar conta a receber');
    },
  });
}

/**
 * Hook para atualizar conta a receber
 */
export function useUpdateReceivable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateContaReceberRequest }) =>
      financialService.updateReceivable(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.receivables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.receivable(id) });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      toast.success('Conta a receber atualizada!');
    },
    onError: (error: Error) => {
      console.error('[useUpdateReceivable] Erro:', error);
      toast.error('Erro ao atualizar conta a receber');
    },
  });
}

/**
 * Hook para deletar conta a receber
 */
export function useDeleteReceivable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.deleteReceivable(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.receivables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      toast.success('Conta a receber excluída!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteReceivable] Erro:', error);
      toast.error('Erro ao excluir conta a receber');
    },
  });
}

/**
 * Hook para marcar conta como recebida
 */
export function useMarkReceivableAsReceived() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: MarcarRecebimentoRequest }) =>
      financialService.markReceivableAsReceived(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.receivables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.receivable(id) });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      queryClient.invalidateQueries({ queryKey: financialKeys.cashflow() });
      toast.success('Recebimento registrado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useMarkReceivableAsReceived] Erro:', error);
      toast.error('Erro ao registrar recebimento');
    },
  });
}

/**
 * Alias para useMarkReceivableAsReceived (compatibilidade)
 */
export function useReceiveReceivable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data_recebimento }: { id: string; data_recebimento: string }) =>
      financialService.markReceivableAsReceived(id, { 
        data_recebimento,
        valor_pago: '0', // Será calculado no backend
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.receivables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.receivable(id) });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      queryClient.invalidateQueries({ queryKey: financialKeys.cashflow() });
      toast.success('Recebimento registrado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useReceiveReceivable] Erro:', error);
      toast.error('Erro ao registrar recebimento');
    },
  });
}

/**
 * Hook para cancelar conta a receber
 */
export function useCancelReceivable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.cancelReceivable(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.receivables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      queryClient.invalidateQueries({ queryKey: financialKeys.vencimentos() });
      toast.success('Conta cancelada!');
    },
    onError: (error: Error) => {
      console.error('[useCancelReceivable] Erro:', error);
      toast.error('Erro ao cancelar conta');
    },
  });
}

// =============================================================================
// COMPENSAÇÕES BANCÁRIAS
// =============================================================================

/**
 * Hook para listar compensações
 */
export function useCompensations(filters: ListCompensacoesFilters = {}) {
  return useQuery({
    queryKey: financialKeys.compensationsList(filters),
    queryFn: () => financialService.listCompensations(filters),
  });
}

/**
 * Hook para deletar compensação
 */
export function useDeleteCompensation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.deleteCompensation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.compensations() });
      toast.success('Compensação excluída!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteCompensation] Erro:', error);
      toast.error('Erro ao excluir compensação');
    },
  });
}

// =============================================================================
// FLUXO DE CAIXA
// =============================================================================

/**
 * Hook para listar fluxo de caixa
 */
export function useCashFlow(filters: ListFluxoCaixaFilters = {}) {
  return useQuery({
    queryKey: financialKeys.cashflowList(filters),
    queryFn: () => financialService.listCashFlow(filters),
  });
}

/**
 * Hook para buscar fluxo de caixa por ID
 */
export function useCashFlowById(id: string) {
  return useQuery({
    queryKey: [...financialKeys.cashflow(), id],
    queryFn: () => financialService.getCashFlow(id),
    enabled: !!id,
  });
}

// =============================================================================
// DRE - DEMONSTRATIVO DE RESULTADOS
// =============================================================================

/**
 * Hook para listar DREs por período
 */
export function useDRE(filters: ListDREFilters = {}) {
  return useQuery({
    queryKey: financialKeys.dreList(filters),
    queryFn: () => financialService.listDRE(filters),
    staleTime: 5 * 60 * 1000, // 5 minutos - DRE muda pouco
  });
}

/**
 * Hook para buscar DRE de um mês específico
 */
export function useDREByMonth(mesAno: string) {
  return useQuery({
    queryKey: financialKeys.dreByMonth(mesAno),
    queryFn: () => financialService.getDRE(mesAno),
    enabled: !!mesAno,
    staleTime: 5 * 60 * 1000,
  });
}

// =============================================================================
// DASHBOARD / RESUMOS
// =============================================================================

/**
 * Hook para resumo financeiro (Dashboard)
 */
export function useFinancialSummary() {
  return useQuery({
    queryKey: financialKeys.summary(),
    queryFn: () => financialService.getFinancialSummary(),
    staleTime: 60 * 1000, // 1 minuto
    refetchInterval: 5 * 60 * 1000, // Atualiza a cada 5 minutos
  });
}

/**
 * Hook para próximos vencimentos
 */
export function useProximosVencimentos(limit: number = 5) {
  return useQuery({
    queryKey: [...financialKeys.vencimentos(), limit],
    queryFn: () => financialService.getProximosVencimentos(limit),
    staleTime: 60 * 1000,
  });
}

/**
 * Hook para painel mensal consolidado (Dashboard real do backend)
 */
export function useDashboard(year?: number, month?: number) {
  return useQuery({
    queryKey: financialKeys.dashboard(year, month),
    queryFn: () => financialService.getDashboard(year, month),
    staleTime: 60 * 1000, // 1 minuto
    refetchInterval: 5 * 60 * 1000, // Atualiza a cada 5 minutos
  });
}

/**
 * Hook para projeções financeiras
 */
export function useProjections(monthsAhead: number = 3) {
  return useQuery({
    queryKey: financialKeys.projections(monthsAhead),
    queryFn: () => financialService.getProjections(monthsAhead),
    staleTime: 5 * 60 * 1000, // 5 minutos - projeções mudam menos
  });
}

// =============================================================================
// DESPESAS FIXAS - QUERIES
// =============================================================================

/**
 * Hook para listar despesas fixas
 */
export function useFixedExpenses(filters: ListDespesasFixasFilters = {}) {
  return useQuery({
    queryKey: financialKeys.fixedExpensesList(filters),
    queryFn: () => financialService.listFixedExpenses(filters),
  });
}

/**
 * Hook para buscar despesa fixa por ID
 */
export function useFixedExpense(id: string) {
  return useQuery({
    queryKey: financialKeys.fixedExpense(id),
    queryFn: () => financialService.getFixedExpense(id),
    enabled: !!id,
  });
}

/**
 * Hook para resumo de despesas fixas
 */
export function useFixedExpensesSummary() {
  return useQuery({
    queryKey: financialKeys.fixedExpensesSummary(),
    queryFn: () => financialService.getFixedExpensesSummary(),
    staleTime: 5 * 60 * 1000,
  });
}

// =============================================================================
// DESPESAS FIXAS - MUTATIONS
// =============================================================================

/**
 * Hook para criar despesa fixa
 */
export function useCreateFixedExpense() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDespesaFixaRequest) => financialService.createFixedExpense(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.fixedExpenses() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      toast.success('Despesa fixa criada com sucesso');
    },
    onError: (error: Error) => {
      toast.error('Erro ao criar despesa fixa', { description: error.message });
    },
  });
}

/**
 * Hook para atualizar despesa fixa
 */
export function useUpdateFixedExpense() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDespesaFixaRequest }) =>
      financialService.updateFixedExpense(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.fixedExpenses() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      toast.success('Despesa fixa atualizada com sucesso');
    },
    onError: (error: Error) => {
      toast.error('Erro ao atualizar despesa fixa', { description: error.message });
    },
  });
}

/**
 * Hook para deletar despesa fixa
 */
export function useDeleteFixedExpense() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.deleteFixedExpense(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: financialKeys.fixedExpenses() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      toast.success('Despesa fixa excluída com sucesso');
    },
    onError: (error: Error) => {
      toast.error('Erro ao excluir despesa fixa', { description: error.message });
    },
  });
}

/**
 * Hook para ativar/desativar despesa fixa
 */
export function useToggleFixedExpense() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => financialService.toggleFixedExpense(id),
    onSuccess: (despesa) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.fixedExpenses() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      const status = despesa.ativo ? 'ativada' : 'desativada';
      toast.success(`Despesa fixa ${status} com sucesso`);
    },
    onError: (error: Error) => {
      toast.error('Erro ao alterar status da despesa fixa', { description: error.message });
    },
  });
}

/**
 * Hook para gerar contas a pagar a partir de despesas fixas
 */
export function useGeneratePayablesFromFixed() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: GerarContasRequest) => financialService.generatePayablesFromFixed(data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: financialKeys.payables() });
      queryClient.invalidateQueries({ queryKey: financialKeys.summary() });
      toast.success(`Contas geradas: ${result.contas_criadas} de ${result.total_despesas} despesas fixas`);
      if (result.erros > 0) {
        toast.warning(`${result.erros} erro(s) durante a geração`);
      }
    },
    onError: (error: Error) => {
      toast.error('Erro ao gerar contas a pagar', { description: error.message });
    },
  });
}

// =============================================================================
// HELPERS
// =============================================================================

/**
 * Hook para invalidar todo cache financeiro
 * Útil após operações que afetam múltiplos módulos
 */
export function useInvalidateFinancialCache() {
  const queryClient = useQueryClient();

  return () => {
    queryClient.invalidateQueries({ queryKey: financialKeys.all });
  };
}
