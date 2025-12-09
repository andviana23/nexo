/**
 * Cores de Status de Agendamento - NEXO v1.0
 * Mapeamento visual conforme FLUXO_STATUS_AGENDAMENTO.md
 */

export type AppointmentStatus =
  | 'CREATED'
  | 'CONFIRMED'
  | 'CHECKED_IN'
  | 'IN_SERVICE'
  | 'AWAITING_PAYMENT'
  | 'DONE'
  | 'NO_SHOW'
  | 'CANCELED';

/**
 * Classes Tailwind para background, border e texto dos cards
 */
export const APPOINTMENT_STATUS_COLORS = {
  CREATED: 'bg-amber-100 border-amber-400 text-amber-900',
  CONFIRMED: 'bg-green-100 border-green-400 text-green-900',
  CHECKED_IN: 'bg-blue-100 border-blue-400 text-blue-900',
  IN_SERVICE: 'bg-purple-100 border-purple-400 text-purple-900',
  AWAITING_PAYMENT: 'bg-orange-100 border-orange-400 text-orange-900',
  DONE: 'bg-slate-100 border-slate-400 text-slate-700',
  NO_SHOW: 'bg-red-100 border-red-400 text-red-900',
  CANCELED: 'bg-slate-200 border-slate-500 text-slate-600',
} as const;

/**
 * Cores sólidas para eventos do FullCalendar
 */
export const APPOINTMENT_CALENDAR_COLORS = {
  CREATED: { backgroundColor: '#fbbf24', borderColor: '#f59e0b', textColor: '#78350f' },
  CONFIRMED: { backgroundColor: '#22c55e', borderColor: '#16a34a', textColor: '#14532d' },
  CHECKED_IN: { backgroundColor: '#3b82f6', borderColor: '#2563eb', textColor: '#1e3a8a' },
  IN_SERVICE: { backgroundColor: '#a855f7', borderColor: '#9333ea', textColor: '#581c87' },
  AWAITING_PAYMENT: { backgroundColor: '#f97316', borderColor: '#ea580c', textColor: '#7c2d12' },
  DONE: { backgroundColor: '#94a3b8', borderColor: '#64748b', textColor: '#1e293b' },
  NO_SHOW: { backgroundColor: '#ef4444', borderColor: '#dc2626', textColor: '#7f1d1d' },
  CANCELED: { backgroundColor: '#64748b', borderColor: '#475569', textColor: '#f1f5f9' },
} as const;

/**
 * Variantes de Badge para cada status
 */
export const APPOINTMENT_BADGE_VARIANTS = {
  CREATED: 'warning',
  CONFIRMED: 'success',
  CHECKED_IN: 'info',
  IN_SERVICE: 'purple',
  AWAITING_PAYMENT: 'orange',
  DONE: 'secondary',
  NO_SHOW: 'destructive',
  CANCELED: 'outline',
} as const;

/**
 * Labels em português para cada status
 */
export const APPOINTMENT_STATUS_LABELS: Record<AppointmentStatus, string> = {
  CREATED: 'Criado',
  CONFIRMED: 'Confirmado',
  CHECKED_IN: 'Cliente Chegou',
  IN_SERVICE: 'Em Atendimento',
  AWAITING_PAYMENT: 'Aguardando Pagamento',
  DONE: 'Concluído',
  NO_SHOW: 'Não Compareceu',
  CANCELED: 'Cancelado',
};

/**
 * Ícones Lucide para cada status
 */
export const APPOINTMENT_STATUS_ICONS = {
  CREATED: 'Calendar',
  CONFIRMED: 'CheckCircle',
  CHECKED_IN: 'UserCheck',
  IN_SERVICE: 'Scissors',
  AWAITING_PAYMENT: 'CreditCard',
  DONE: 'CheckCheck',
  NO_SHOW: 'UserX',
  CANCELED: 'XCircle',
} as const;

/**
 * Helper para obter classes de cor de um status
 */
export function getStatusColorClasses(status: string): string {
  return APPOINTMENT_STATUS_COLORS[status as AppointmentStatus] || APPOINTMENT_STATUS_COLORS.CREATED;
}

/**
 * Helper para obter cores do calendário
 */
export function getCalendarColors(status: string) {
  return APPOINTMENT_CALENDAR_COLORS[status as AppointmentStatus] || APPOINTMENT_CALENDAR_COLORS.CREATED;
}

/**
 * Helper para obter label do status
 */
export function getStatusLabel(status: string): string {
  return APPOINTMENT_STATUS_LABELS[status as AppointmentStatus] || status;
}
