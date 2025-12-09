/**
 * NEXO - Sistema de Gestão para Barbearias
 * Commission Hooks
 *
 * Hooks React Query para gerenciar estado do módulo de Comissões.
 */

import { commissionService } from '@/services/commission-service';
import type {
    AssignItemsToPeriodRequest,
    CloseCommissionPeriodRequest,
    CreateAdvanceRequest,
    CreateCommissionItemBatchRequest,
    CreateCommissionItemRequest,
    CreateCommissionPeriodRequest,
    CreateCommissionRuleRequest,
    ListAdvancesFilters,
    ListCommissionItemsFilters,
    ListCommissionPeriodsFilters,
    ListCommissionRulesFilters,
    MarkAdvanceDeductedRequest,
    ProcessCommissionItemRequest,
    RejectAdvanceRequest,
    UpdateCommissionRuleRequest,
} from '@/types/commission';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const commissionKeys = {
  all: ['commissions'] as const,
  
  // Commission Rules
  rules: () => [...commissionKeys.all, 'rules'] as const,
  rulesList: (filters: ListCommissionRulesFilters) => [...commissionKeys.rules(), 'list', filters] as const,
  rule: (id: string) => [...commissionKeys.rules(), 'detail', id] as const,
  effectiveRules: (professionalId?: string) => [...commissionKeys.rules(), 'effective', professionalId] as const,
  
  // Commission Periods
  periods: () => [...commissionKeys.all, 'periods'] as const,
  periodsList: (filters: ListCommissionPeriodsFilters) => [...commissionKeys.periods(), 'list', filters] as const,
  period: (id: string) => [...commissionKeys.periods(), 'detail', id] as const,
  periodOpen: (professionalId: string) => [...commissionKeys.periods(), 'open', professionalId] as const,
  periodSummary: (id: string) => [...commissionKeys.periods(), 'summary', id] as const,
  
  // Advances
  advances: () => [...commissionKeys.all, 'advances'] as const,
  advancesList: (filters: ListAdvancesFilters) => [...commissionKeys.advances(), 'list', filters] as const,
  advance: (id: string) => [...commissionKeys.advances(), 'detail', id] as const,
  advancesPending: (professionalId?: string) => [...commissionKeys.advances(), 'pending', professionalId] as const,
  advancesApproved: (professionalId?: string) => [...commissionKeys.advances(), 'approved', professionalId] as const,
  
  // Commission Items
  items: () => [...commissionKeys.all, 'items'] as const,
  itemsList: (filters: ListCommissionItemsFilters) => [...commissionKeys.items(), 'list', filters] as const,
  item: (id: string) => [...commissionKeys.items(), 'detail', id] as const,
  itemsPending: (professionalId: string) => [...commissionKeys.items(), 'pending', professionalId] as const,
  
  // Summaries
  summaryByProfessional: (startDate: string, endDate: string, professionalId?: string) => 
    [...commissionKeys.all, 'summary', 'professional', startDate, endDate, professionalId] as const,
  summaryByService: (startDate: string, endDate: string, serviceId?: string) => 
    [...commissionKeys.all, 'summary', 'service', startDate, endDate, serviceId] as const,
};

// =============================================================================
// COMMISSION RULES - QUERIES
// =============================================================================

/**
 * Hook para listar regras de comissão
 */
export function useCommissionRules(filters: ListCommissionRulesFilters = {}) {
  return useQuery({
    queryKey: commissionKeys.rulesList(filters),
    queryFn: () => commissionService.listRules(filters),
  });
}

/**
 * Hook para buscar regra de comissão por ID
 */
export function useCommissionRule(id: string) {
  return useQuery({
    queryKey: commissionKeys.rule(id),
    queryFn: () => commissionService.getRule(id),
    enabled: !!id,
  });
}

/**
 * Hook para buscar regras efetivas
 */
export function useEffectiveCommissionRules(professionalId?: string) {
  return useQuery({
    queryKey: commissionKeys.effectiveRules(professionalId),
    queryFn: () => commissionService.getEffectiveRules(professionalId),
  });
}

// =============================================================================
// COMMISSION RULES - MUTATIONS
// =============================================================================

/**
 * Hook para criar regra de comissão
 */
export function useCreateCommissionRule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommissionRuleRequest) => commissionService.createRule(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.rules() });
      toast.success('Regra de comissão criada com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useCreateCommissionRule] Erro:', error);
      toast.error('Erro ao criar regra de comissão');
    },
  });
}

/**
 * Hook para atualizar regra de comissão
 */
export function useUpdateCommissionRule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCommissionRuleRequest }) =>
      commissionService.updateRule(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.rules() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.rule(id) });
      toast.success('Regra de comissão atualizada!');
    },
    onError: (error: Error) => {
      console.error('[useUpdateCommissionRule] Erro:', error);
      toast.error('Erro ao atualizar regra de comissão');
    },
  });
}

/**
 * Hook para deletar regra de comissão
 */
export function useDeleteCommissionRule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.deleteRule(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.rules() });
      toast.success('Regra de comissão removida!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteCommissionRule] Erro:', error);
      toast.error('Erro ao remover regra de comissão');
    },
  });
}

// =============================================================================
// COMMISSION PERIODS - QUERIES
// =============================================================================

/**
 * Hook para listar períodos de comissão
 */
export function useCommissionPeriods(filters: ListCommissionPeriodsFilters = {}) {
  return useQuery({
    queryKey: commissionKeys.periodsList(filters),
    queryFn: () => commissionService.listPeriods(filters),
  });
}

/**
 * Hook para buscar período de comissão por ID
 */
export function useCommissionPeriod(id: string) {
  return useQuery({
    queryKey: commissionKeys.period(id),
    queryFn: () => commissionService.getPeriod(id),
    enabled: !!id,
  });
}

/**
 * Hook para buscar período aberto de um profissional
 */
export function useOpenCommissionPeriod(professionalId: string) {
  return useQuery({
    queryKey: commissionKeys.periodOpen(professionalId),
    queryFn: () => commissionService.getOpenPeriod(professionalId),
    enabled: !!professionalId,
  });
}

/**
 * Hook para buscar resumo do período
 */
export function useCommissionPeriodSummary(periodId: string) {
  return useQuery({
    queryKey: commissionKeys.periodSummary(periodId),
    queryFn: () => commissionService.getPeriodSummary(periodId),
    enabled: !!periodId,
  });
}

// =============================================================================
// COMMISSION PERIODS - MUTATIONS
// =============================================================================

/**
 * Hook para criar período de comissão
 */
export function useCreateCommissionPeriod() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommissionPeriodRequest) => commissionService.createPeriod(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Período de comissão criado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useCreateCommissionPeriod] Erro:', error);
      toast.error('Erro ao criar período de comissão');
    },
  });
}

/**
 * Hook para fechar período de comissão
 */
export function useCloseCommissionPeriod() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data?: CloseCommissionPeriodRequest }) =>
      commissionService.closePeriod(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.period(id) });
      toast.success('Período fechado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useCloseCommissionPeriod] Erro:', error);
      toast.error('Erro ao fechar período');
    },
  });
}

/**
 * Hook para marcar período como pago
 */
export function useMarkPeriodAsPaid() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.markPeriodAsPaid(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Período marcado como pago!');
    },
    onError: (error: Error) => {
      console.error('[useMarkPeriodAsPaid] Erro:', error);
      toast.error('Erro ao marcar período como pago');
    },
  });
}

/**
 * Hook para deletar período de comissão
 */
export function useDeleteCommissionPeriod() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.deletePeriod(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Período de comissão removido!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteCommissionPeriod] Erro:', error);
      toast.error('Erro ao remover período');
    },
  });
}

// =============================================================================
// ADVANCES - QUERIES
// =============================================================================

/**
 * Hook para listar adiantamentos
 */
export function useAdvances(filters: ListAdvancesFilters = {}) {
  return useQuery({
    queryKey: commissionKeys.advancesList(filters),
    queryFn: () => commissionService.listAdvances(filters),
  });
}

/**
 * Hook para buscar adiantamento por ID
 */
export function useAdvance(id: string) {
  return useQuery({
    queryKey: commissionKeys.advance(id),
    queryFn: () => commissionService.getAdvance(id),
    enabled: !!id,
  });
}

/**
 * Hook para buscar adiantamentos pendentes
 */
export function usePendingAdvances(professionalId?: string) {
  return useQuery({
    queryKey: commissionKeys.advancesPending(professionalId),
    queryFn: () => commissionService.getPendingAdvances(professionalId),
  });
}

/**
 * Hook para buscar adiantamentos aprovados
 */
export function useApprovedAdvances(professionalId?: string) {
  return useQuery({
    queryKey: commissionKeys.advancesApproved(professionalId),
    queryFn: () => commissionService.getApprovedAdvances(professionalId),
  });
}

// =============================================================================
// ADVANCES - MUTATIONS
// =============================================================================

/**
 * Hook para criar adiantamento
 */
export function useCreateAdvance() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateAdvanceRequest) => commissionService.createAdvance(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.advances() });
      toast.success('Adiantamento solicitado com sucesso!');
    },
    onError: (error: Error) => {
      console.error('[useCreateAdvance] Erro:', error);
      toast.error('Erro ao solicitar adiantamento');
    },
  });
}

/**
 * Hook para aprovar adiantamento
 */
export function useApproveAdvance() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.approveAdvance(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.advances() });
      toast.success('Adiantamento aprovado!');
    },
    onError: (error: Error) => {
      console.error('[useApproveAdvance] Erro:', error);
      toast.error('Erro ao aprovar adiantamento');
    },
  });
}

/**
 * Hook para rejeitar adiantamento
 */
export function useRejectAdvance() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: RejectAdvanceRequest }) =>
      commissionService.rejectAdvance(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.advances() });
      toast.success('Adiantamento rejeitado!');
    },
    onError: (error: Error) => {
      console.error('[useRejectAdvance] Erro:', error);
      toast.error('Erro ao rejeitar adiantamento');
    },
  });
}

/**
 * Hook para marcar adiantamento como deduzido
 */
export function useMarkAdvanceDeducted() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: MarkAdvanceDeductedRequest }) =>
      commissionService.markAdvanceDeducted(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.advances() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Adiantamento deduzido do período!');
    },
    onError: (error: Error) => {
      console.error('[useMarkAdvanceDeducted] Erro:', error);
      toast.error('Erro ao deduzir adiantamento');
    },
  });
}

/**
 * Hook para cancelar adiantamento
 */
export function useCancelAdvance() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.cancelAdvance(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.advances() });
      toast.success('Adiantamento cancelado!');
    },
    onError: (error: Error) => {
      console.error('[useCancelAdvance] Erro:', error);
      toast.error('Erro ao cancelar adiantamento');
    },
  });
}

/**
 * Hook para deletar adiantamento
 */
export function useDeleteAdvance() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.deleteAdvance(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.advances() });
      toast.success('Adiantamento removido!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteAdvance] Erro:', error);
      toast.error('Erro ao remover adiantamento');
    },
  });
}

// =============================================================================
// COMMISSION ITEMS - QUERIES
// =============================================================================

/**
 * Hook para listar itens de comissão
 */
export function useCommissionItems(filters: ListCommissionItemsFilters = {}) {
  return useQuery({
    queryKey: commissionKeys.itemsList(filters),
    queryFn: () => commissionService.listItems(filters),
  });
}

/**
 * Hook para buscar item de comissão por ID
 */
export function useCommissionItem(id: string) {
  return useQuery({
    queryKey: commissionKeys.item(id),
    queryFn: () => commissionService.getItem(id),
    enabled: !!id,
  });
}

/**
 * Hook para buscar itens pendentes de um profissional
 */
export function usePendingCommissionItems(professionalId: string) {
  return useQuery({
    queryKey: commissionKeys.itemsPending(professionalId),
    queryFn: () => commissionService.getPendingItems(professionalId),
    enabled: !!professionalId,
  });
}

// =============================================================================
// COMMISSION ITEMS - MUTATIONS
// =============================================================================

/**
 * Hook para criar item de comissão
 */
export function useCreateCommissionItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommissionItemRequest) => commissionService.createItem(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.items() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Item de comissão criado!');
    },
    onError: (error: Error) => {
      console.error('[useCreateCommissionItem] Erro:', error);
      toast.error('Erro ao criar item de comissão');
    },
  });
}

/**
 * Hook para criar múltiplos itens de comissão
 */
export function useCreateCommissionItemsBatch() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommissionItemBatchRequest) => commissionService.createItemsBatch(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.items() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Itens de comissão criados!');
    },
    onError: (error: Error) => {
      console.error('[useCreateCommissionItemsBatch] Erro:', error);
      toast.error('Erro ao criar itens de comissão');
    },
  });
}

/**
 * Hook para processar item de comissão
 */
export function useProcessCommissionItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ProcessCommissionItemRequest }) =>
      commissionService.processItem(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.items() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Item processado!');
    },
    onError: (error: Error) => {
      console.error('[useProcessCommissionItem] Erro:', error);
      toast.error('Erro ao processar item');
    },
  });
}

/**
 * Hook para vincular itens a um período
 */
export function useAssignItemsToPeriod() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AssignItemsToPeriodRequest) => commissionService.assignItemsToPeriod(data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.items() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success(`${result.assigned_count} itens vinculados ao período!`);
    },
    onError: (error: Error) => {
      console.error('[useAssignItemsToPeriod] Erro:', error);
      toast.error('Erro ao vincular itens ao período');
    },
  });
}

/**
 * Hook para deletar item de comissão
 */
export function useDeleteCommissionItem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => commissionService.deleteItem(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commissionKeys.items() });
      queryClient.invalidateQueries({ queryKey: commissionKeys.periods() });
      toast.success('Item de comissão removido!');
    },
    onError: (error: Error) => {
      console.error('[useDeleteCommissionItem] Erro:', error);
      toast.error('Erro ao remover item');
    },
  });
}

// =============================================================================
// SUMMARIES - QUERIES
// =============================================================================

/**
 * Hook para buscar resumo de comissões por profissional
 */
export function useCommissionSummaryByProfessional(
  startDate: string,
  endDate: string,
  professionalId?: string
) {
  return useQuery({
    queryKey: commissionKeys.summaryByProfessional(startDate, endDate, professionalId),
    queryFn: () => commissionService.getSummaryByProfessional(startDate, endDate, professionalId),
    enabled: !!startDate && !!endDate,
  });
}

/**
 * Hook para buscar resumo de comissões por serviço
 */
export function useCommissionSummaryByService(
  startDate: string,
  endDate: string,
  serviceId?: string
) {
  return useQuery({
    queryKey: commissionKeys.summaryByService(startDate, endDate, serviceId),
    queryFn: () => commissionService.getSummaryByService(startDate, endDate, serviceId),
    enabled: !!startDate && !!endDate,
  });
}
