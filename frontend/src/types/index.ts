/**
 * NEXO - Sistema de Gestão para Barbearias
 * Tipos Globais e Interfaces Base
 *
 * Tipos compartilhados em toda a aplicação frontend.
 * Devem refletir os DTOs do backend Go.
 */

// Re-export dos módulos de tipos
export * from './appointment';
// Customer module: usar tipos específicos para evitar conflito com appointment.Customer
export {
    DEFAULT_TAGS, ESTADOS_BR, GENDER_LABELS,
    TAG_COLORS, cleanCEP, cleanCPF, cleanPhone, formatCEP, formatCPF, formatPhone, type CheckExistsParams, type CheckExistsResponse, type CreateCustomerRequest, type CustomerAddressExport,
    type CustomerAppointmentExport,
    type CustomerExportResponse, type Customer as CustomerFull, type CustomerGender, type CustomerModalState, type CustomerResponse, type CustomerStatsResponse, type CustomerSummary,
    type CustomerWithHistory, type ListCustomersFilters, type ListCustomersResponse, type SearchCustomersParams, type UpdateCustomerRequest
} from './customer';
export * from './meio-pagamento';
export * from './professional';
export * from './stock';
export * from './unit';

// =============================================================================
// TIPOS UTILITÁRIOS
// =============================================================================

/**
 * Resposta de erro padrão da API
 */
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string[]>;
  timestamp?: string;
}

/**
 * Resposta paginada padrão
 */
export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

/**
 * Filtros base para listagens
 */
export interface BaseFilters {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

// =============================================================================
// AUTENTICAÇÃO
// =============================================================================

/**
 * Credenciais de login
 */
export interface LoginCredentials {
  email: string;
  password: string;
}

/**
 * Resposta de login
 * NOTA: Campos mapeados do backend Go (snake_case -> camelCase)
 */
export interface LoginResponse {
  access_token: string;  // Backend retorna "access_token"
  user: User;
  tenant: Tenant;
}

/**
 * Dados para refresh de token
 */
export interface RefreshTokenResponse {
  token: string;
}

// =============================================================================
// USUÁRIO
// =============================================================================

/**
 * Roles disponíveis no sistema
 */
export type UserRole = 'owner' | 'admin' | 'manager' | 'accountant' | 'employee' | 'barbeiro' | 'professional' | 'receptionist';

/**
 * Status do usuário
 */
export type UserStatus = 'active' | 'inactive' | 'pending';

/**
 * Usuário do sistema
 */
export interface User {
  id: string;
  tenant_id: string;
  email: string;
  name: string;
  role: UserRole;
  status: UserStatus;
  avatar_url?: string;
  phone?: string;
  created_at: string;
  updated_at: string;
}

/**
 * Dados para criação de usuário
 */
export interface CreateUserRequest {
  email: string;
  name: string;
  password: string;
  role: UserRole;
  phone?: string;
}

/**
 * Dados para atualização de usuário
 */
export interface UpdateUserRequest {
  name?: string;
  role?: UserRole;
  status?: UserStatus;
  phone?: string;
  avatar_url?: string;
}

// =============================================================================
// TENANT (BARBEARIA)
// =============================================================================

/**
 * Status do tenant
 */
export type TenantStatus = 'active' | 'inactive' | 'suspended' | 'trial';

/**
 * Plano de assinatura
 */
export type SubscriptionPlan =
  | 'free'
  | 'starter'
  | 'professional'
  | 'enterprise';

/**
 * Tenant (Barbearia)
 */
export interface Tenant {
  id: string;
  name: string;
  slug: string;
  status: TenantStatus;
  plan: SubscriptionPlan;
  logo_url?: string;
  phone?: string;
  email?: string;
  address?: TenantAddress;
  settings?: TenantSettings;
  // Feature flags
  multi_unit_enabled?: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * Endereço do tenant
 */
export interface TenantAddress {
  street: string;
  number: string;
  complement?: string;
  neighborhood: string;
  city: string;
  state: string;
  zip_code: string;
}

/**
 * Configurações do tenant
 */
export interface TenantSettings {
  timezone: string;
  currency: string;
  date_format: string;
  time_format: '12h' | '24h';
  week_starts_on: 0 | 1; // 0 = Sunday, 1 = Monday
  appointment_interval: number; // em minutos
  cancellation_policy_hours: number;
  allow_online_booking: boolean;
}

// =============================================================================
// CLIENTE
// =============================================================================

/**
 * Cliente da barbearia
 */
export interface Client {
  id: string;
  tenant_id: string;
  name: string;
  email?: string;
  phone: string;
  cpf?: string;
  birth_date?: string;
  gender?: 'male' | 'female' | 'other';
  notes?: string;
  tags?: string[];
  total_visits: number;
  total_spent: string; // Dinheiro como string
  last_visit?: string;
  created_at: string;
  updated_at: string;
}

/**
 * Dados para criação de cliente
 */
export interface CreateClientRequest {
  name: string;
  phone: string;
  email?: string;
  cpf?: string;
  birth_date?: string;
  gender?: 'male' | 'female' | 'other';
  notes?: string;
  tags?: string[];
}

/**
 * Dados para atualização de cliente
 */
export type UpdateClientRequest = Partial<CreateClientRequest>;

// =============================================================================
// SERVIÇO
// =============================================================================

/**
 * Serviço oferecido
 */
export interface Service {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  price: string; // Dinheiro como string
  duration: number; // em minutos
  category?: string;
  is_active: boolean;
  commission_type: 'percentage' | 'fixed';
  commission_value: string;
  created_at: string;
  updated_at: string;
}

/**
 * Dados para criação de serviço
 */
export interface CreateServiceRequest {
  name: string;
  price: string;
  duration: number;
  description?: string;
  category?: string;
  commission_type?: 'percentage' | 'fixed';
  commission_value?: string;
}

/**
 * Dados para atualização de serviço
 */
export interface UpdateServiceRequest extends Partial<CreateServiceRequest> {
  is_active?: boolean;
}

// =============================================================================
// PROFISSIONAL
// =============================================================================

/**
 * Profissional (barbeiro)
 */
export interface Professional {
  id: string;
  tenant_id: string;
  user_id: string;
  name: string;
  email?: string;
  phone?: string;
  avatar_url?: string;
  bio?: string;
  specialties?: string[];
  is_active: boolean;
  commission_type: 'percentage' | 'fixed';
  commission_value: string;
  working_hours: WorkingHours[];
  created_at: string;
  updated_at: string;
}

/**
 * Horário de trabalho
 */
export interface WorkingHours {
  day_of_week: 0 | 1 | 2 | 3 | 4 | 5 | 6;
  start_time: string; // "HH:mm"
  end_time: string; // "HH:mm"
  is_working: boolean;
}

// =============================================================================
// AGENDAMENTO
// =============================================================================

/**
 * Status do agendamento
 */
export type AppointmentStatus =
  | 'scheduled'
  | 'confirmed'
  | 'in_progress'
  | 'completed'
  | 'cancelled'
  | 'no_show';

/**
 * Agendamento
 */
export interface Appointment {
  id: string;
  tenant_id: string;
  client_id: string;
  client?: Client;
  professional_id: string;
  professional?: Professional;
  service_id: string;
  service?: Service;
  date: string;
  start_time: string;
  end_time: string;
  status: AppointmentStatus;
  price: string;
  notes?: string;
  source: 'manual' | 'online' | 'walk_in';
  created_at: string;
  updated_at: string;
}

/**
 * Dados para criação de agendamento
 */
export interface CreateAppointmentRequest {
  client_id: string;
  professional_id: string;
  service_id: string;
  date: string;
  start_time: string;
  notes?: string;
  source?: 'manual' | 'online' | 'walk_in';
}

/**
 * Dados para atualização de agendamento
 */
export interface UpdateAppointmentRequest {
  professional_id?: string;
  service_id?: string;
  date?: string;
  start_time?: string;
  status?: AppointmentStatus;
  notes?: string;
}

// =============================================================================
// LISTA DA VEZ
// =============================================================================

/**
 * Status na lista da vez
 */
export type QueueStatus = 'waiting' | 'in_service' | 'completed' | 'cancelled';

/**
 * Item da lista da vez
 */
export interface QueueItem {
  id: string;
  tenant_id: string;
  client_id: string;
  client?: Client;
  professional_id?: string;
  professional?: Professional;
  service_id?: string;
  service?: Service;
  position: number;
  status: QueueStatus;
  priority: boolean;
  check_in_time: string;
  start_time?: string;
  end_time?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

/**
 * Dados para adicionar à lista da vez
 */
export interface AddToQueueRequest {
  client_id: string;
  professional_id?: string;
  service_id?: string;
  priority?: boolean;
  notes?: string;
}

// =============================================================================
// FINANCEIRO
// =============================================================================

/**
 * Tipo de transação
 */
export type TransactionType = 'income' | 'expense';

/**
 * Categoria de transação
 */
export type TransactionCategory =
  | 'service'
  | 'product'
  | 'commission'
  | 'salary'
  | 'rent'
  | 'utilities'
  | 'supplies'
  | 'marketing'
  | 'other';

/**
 * Método de pagamento
 */
export type PaymentMethod =
  | 'cash'
  | 'credit_card'
  | 'debit_card'
  | 'pix'
  | 'transfer'
  | 'other';

/**
 * Transação financeira
 */
export interface Transaction {
  id: string;
  tenant_id: string;
  type: TransactionType;
  category: TransactionCategory;
  amount: string;
  description: string;
  payment_method?: PaymentMethod;
  reference_id?: string;
  reference_type?: 'appointment' | 'sale' | 'expense';
  date: string;
  created_at: string;
  updated_at: string;
}

/**
 * Resumo financeiro
 */
export interface FinancialSummary {
  period: string;
  total_income: string;
  total_expenses: string;
  net_profit: string;
  income_by_category: Record<string, string>;
  expenses_by_category: Record<string, string>;
  income_by_payment_method: Record<string, string>;
}

// =============================================================================
// METAS
// =============================================================================

/**
 * Tipo de meta
 */
export type GoalType = 'revenue' | 'appointments' | 'clients' | 'services';

/**
 * Período da meta
 */
export type GoalPeriod = 'daily' | 'weekly' | 'monthly' | 'yearly';

/**
 * Meta
 */
export interface Goal {
  id: string;
  tenant_id: string;
  professional_id?: string;
  professional?: Professional;
  type: GoalType;
  period: GoalPeriod;
  target_value: string;
  current_value: string;
  start_date: string;
  end_date: string;
  is_achieved: boolean;
  created_at: string;
  updated_at: string;
}

// =============================================================================
// DASHBOARD
// =============================================================================

/**
 * Estatísticas do dashboard
 */
export interface DashboardStats {
  today_appointments: number;
  today_revenue: string;
  week_revenue: string;
  month_revenue: string;
  active_clients: number;
  queue_size: number;
  avg_ticket: string;
  occupation_rate: number;
}

/**
 * Dados para gráficos do dashboard
 */
export interface DashboardChartData {
  period: string;
  revenue: ChartDataPoint[];
  appointments: ChartDataPoint[];
  clients: ChartDataPoint[];
}

/**
 * Ponto de dado para gráfico
 */
export interface ChartDataPoint {
  label: string;
  value: number;
}
