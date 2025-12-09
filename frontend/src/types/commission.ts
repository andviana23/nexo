/**
 * NEXO - Sistema de Gestão para Barbearias
 * Commission Types
 *
 * Tipos TypeScript para o módulo de Comissões.
 * Mapeados a partir do backend Go (commission_dto.go)
 */

// =============================================================================
// ENUMS - STATUS
// =============================================================================

export enum CommissionRuleType {
  PERCENTUAL = 'PERCENTUAL',
  FIXO = 'FIXO',
}

export enum CalculationBase {
  BRUTO = 'BRUTO',
  LIQUIDO = 'LIQUIDO',
}

export enum CommissionPeriodStatus {
  ABERTO = 'ABERTO',
  PROCESSANDO = 'PROCESSANDO',
  FECHADO = 'FECHADO',
  PAGO = 'PAGO',
  CANCELADO = 'CANCELADO',
}

export enum AdvanceStatus {
  PENDING = 'PENDING',
  APPROVED = 'APPROVED',
  REJECTED = 'REJECTED',
  DEDUCTED = 'DEDUCTED',
  CANCELLED = 'CANCELLED',
}

export enum CommissionItemStatus {
  PENDENTE = 'PENDENTE',
  PROCESSADO = 'PROCESSADO',
  PAGO = 'PAGO',
  CANCELADO = 'CANCELADO',
  ESTORNADO = 'ESTORNADO',
}

export enum CommissionSource {
  SERVICO = 'SERVICO',
  PROFISSIONAL = 'PROFISSIONAL',
  REGRA = 'REGRA',
  MANUAL = 'MANUAL',
}

// =============================================================================
// COMMISSION RULE
// =============================================================================

export interface CommissionRule {
  id: string;
  tenant_id: string;
  unit_id?: string;
  name: string;
  description?: string;
  type: CommissionRuleType;
  default_rate: string;
  min_amount?: string;
  max_amount?: string;
  calculation_base?: CalculationBase;
  effective_from: string;
  effective_to?: string;
  priority?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateCommissionRuleRequest {
  unit_id?: string;
  name: string;
  description?: string;
  type: CommissionRuleType;
  default_rate: string;
  min_amount?: string;
  max_amount?: string;
  calculation_base?: CalculationBase;
  effective_from?: string;
  effective_to?: string;
  priority?: number;
}

export interface UpdateCommissionRuleRequest {
  name?: string;
  description?: string;
  type?: CommissionRuleType;
  default_rate?: string;
  min_amount?: string;
  max_amount?: string;
  calculation_base?: CalculationBase;
  effective_from?: string;
  effective_to?: string;
  priority?: number;
  is_active?: boolean;
}

export interface ListCommissionRulesFilters {
  active_only?: boolean;
}

export interface ListCommissionRulesResponse {
  data: CommissionRule[];
  total: number;
}

// =============================================================================
// COMMISSION PERIOD
// =============================================================================

export interface CommissionPeriod {
  id: string;
  tenant_id: string;
  unit_id?: string;
  reference_month: string;
  professional_id?: string;
  professional_name?: string;
  total_gross: string;
  total_commission: string;
  total_advances: string;
  total_adjustments: string;
  total_net: string;
  items_count: number;
  status: CommissionPeriodStatus;
  period_start: string;
  period_end: string;
  closed_at?: string;
  paid_at?: string;
  closed_by?: string;
  paid_by?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCommissionPeriodRequest {
  unit_id?: string;
  reference_month: string;
  professional_id: string;
  period_start: string;
  period_end: string;
  notes?: string;
}

export interface CloseCommissionPeriodRequest {
  notes?: string;
}

export interface ListCommissionPeriodsFilters {
  professional_id?: string;
  status?: CommissionPeriodStatus;
  limit?: number;
  offset?: number;
}

export interface ListCommissionPeriodsResponse {
  data: CommissionPeriod[];
  total: number;
}

export interface CommissionPeriodSummary {
  total_gross: string;
  total_commission: string;
  total_advances: string;
  total_net: string;
  items_count: number;
}

// =============================================================================
// ADVANCE (ADIANTAMENTO)
// =============================================================================

export interface Advance {
  id: string;
  tenant_id: string;
  unit_id?: string;
  professional_id: string;
  professional_name?: string;
  amount: string;
  request_date: string;
  reason?: string;
  status: AdvanceStatus;
  approved_at?: string;
  approved_by?: string;
  rejected_at?: string;
  rejected_by?: string;
  rejection_reason?: string;
  deducted_at?: string;
  deduction_period_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateAdvanceRequest {
  unit_id?: string;
  professional_id: string;
  amount: string;
  reason?: string;
}

export interface RejectAdvanceRequest {
  rejection_reason: string;
}

export interface MarkAdvanceDeductedRequest {
  period_id: string;
}

export interface ListAdvancesFilters {
  professional_id?: string;
  status?: AdvanceStatus;
  limit?: number;
  offset?: number;
}

export interface ListAdvancesResponse {
  data: Advance[];
  total: number;
}

export interface AdvancesTotals {
  advances: Advance[];
  total_pending: string;
  total_approved: string;
}

// =============================================================================
// COMMISSION ITEM
// =============================================================================

export interface CommissionItem {
  id: string;
  tenant_id: string;
  unit_id?: string;
  professional_id: string;
  professional_name?: string;
  command_id?: string;
  command_item_id?: string;
  appointment_id?: string;
  service_id?: string;
  service_name?: string;
  gross_value: string;
  commission_rate: string;
  commission_type: CommissionRuleType;
  commission_value: string;
  commission_source: CommissionSource;
  rule_id?: string;
  reference_date: string;
  description?: string;
  status: CommissionItemStatus;
  period_id?: string;
  processed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCommissionItemRequest {
  unit_id?: string;
  professional_id: string;
  command_id?: string;
  command_item_id?: string;
  appointment_id?: string;
  service_id?: string;
  service_name?: string;
  gross_value: string;
  commission_rate: string;
  commission_type: CommissionRuleType;
  commission_source: CommissionSource;
  rule_id?: string;
  reference_date: string;
  description?: string;
}

export interface CreateCommissionItemBatchRequest {
  items: CreateCommissionItemRequest[];
}

export interface ProcessCommissionItemRequest {
  period_id: string;
}

export interface AssignItemsToPeriodRequest {
  professional_id: string;
  period_id: string;
  start_date: string;
  end_date: string;
}

export interface ListCommissionItemsFilters {
  professional_id?: string;
  period_id?: string;
  status?: CommissionItemStatus;
  limit?: number;
  offset?: number;
}

export interface ListCommissionItemsResponse {
  data: CommissionItem[];
  total: number;
}

// =============================================================================
// SUMMARIES
// =============================================================================

export interface CommissionSummary {
  professional_id: string;
  professional_name: string;
  total_gross: string;
  total_commission: string;
  items_count: number;
}

export interface CommissionByService {
  service_id: string;
  service_name: string;
  total_gross: string;
  total_commission: string;
  items_count: number;
}

export interface CommissionSummaries {
  by_professional?: CommissionSummary[];
  by_service?: CommissionByService[];
  start_date: string;
  end_date: string;
}

// =============================================================================
// HELPER TYPES
// =============================================================================

export interface CommissionFilters {
  professional_id?: string;
  start_date?: string;
  end_date?: string;
}
