/**
 * Configuração do FullCalendar para o NEXO
 * 
 * ⚠️ ATENÇÃO LEGAL: Licença FullCalendar Scheduler – Modo Avaliação
 * 
 * Durante o período de avaliação gratuita do FullCalendar Premium (Scheduler),
 * o NEXO está autorizado a utilizar a licença não-comercial fornecida pelo
 * próprio FullCalendar, exclusivamente para fins de desenvolvimento interno,
 * sem uso comercial e sem disponibilização aos clientes finais.
 * 
 * RESTRIÇÕES:
 * - ❌ Proibido uso comercial desta licença.
 * - ✅ Permitido apenas para testes internos, homologação e desenvolvimento.
 * - ⚠️ A versão final exigirá a compra da licença oficial.
 * - 🔄 Substituir a chave de desenvolvimento pela licença comercial antes da produção.
 * 
 * @see https://fullcalendar.io/docs/schedulerLicenseKey
 */

// Chave de licença para ambiente de desenvolvimento
// TODO: Substituir pela chave comercial em produção
export const FULLCALENDAR_LICENSE_KEY = 'CC-Attribution-NonCommercial-NoDerivatives';

// Configurações padrão do calendário
export const FULLCALENDAR_DEFAULTS = {
  // Localização
  locale: 'pt-br',
  
  // Horário de funcionamento padrão (barbearias)
  slotMinTime: '08:00:00',
  slotMaxTime: '20:00:00',
  
  // Intervalo de slots (10 minutos)
  slotDuration: '00:10:00',
  
  // Primeiro dia da semana (0 = domingo, 1 = segunda)
  firstDay: 1,
  
  // Altura do slot
  slotLabelInterval: '01:00:00',
  
  // Dias úteis (segunda a sábado)
  businessHours: {
    daysOfWeek: [1, 2, 3, 4, 5, 6], // Segunda a Sábado
    startTime: '08:00',
    endTime: '20:00',
  },
  
  // Formato de hora
  eventTimeFormat: {
    hour: '2-digit',
    minute: '2-digit',
    meridiem: false,
    hour12: false,
  },
  
  // Formato do slot
  slotLabelFormat: {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  },
  
  // Header toolbar padrão
  headerToolbar: {
    left: 'prev,next today',
    center: 'title',
    right: 'resourceTimeGridDay,resourceTimeGridWeek,listWeek',
  },
  
  // Navegação por scroll do mouse
  navLinks: true,
  
  // Permitir seleção de slots
  selectable: true,
  
  // Permitir arrastar eventos
  editable: true,
  
  // Permitir redimensionar eventos
  eventResizableFromStart: true,
  
  // Mostrar indicador de "agora"
  nowIndicator: true,
  
  // Limitar eventos por célula
  dayMaxEvents: true,
  
  // Expandir linhas de recursos
  expandRows: true,
} as const;

// Mapeamento de status para cores - Conforme FLUXO_STATUS_AGENDAMENTO.md
export const APPOINTMENT_STATUS_COLORS: Record<string, { backgroundColor: string; borderColor: string; textColor: string }> = {
  CREATED: {
    backgroundColor: '#fbbf24', // amber-400
    borderColor: '#f59e0b',     // amber-500
    textColor: '#78350f',       // amber-900
  },
  CONFIRMED: {
    backgroundColor: '#22c55e', // green-500
    borderColor: '#16a34a',     // green-600
    textColor: '#14532d',       // green-900
  },
  CHECKED_IN: {
    backgroundColor: '#3b82f6', // blue-500
    borderColor: '#2563eb',     // blue-600
    textColor: '#1e3a8a',       // blue-900
  },
  IN_SERVICE: {
    backgroundColor: '#a855f7', // purple-500
    borderColor: '#9333ea',     // purple-600
    textColor: '#581c87',       // purple-900
  },
  AWAITING_PAYMENT: {
    backgroundColor: '#f97316', // orange-500
    borderColor: '#ea580c',     // orange-600
    textColor: '#7c2d12',       // orange-900
  },
  DONE: {
    backgroundColor: '#94a3b8', // slate-400
    borderColor: '#64748b',     // slate-500
    textColor: '#1e293b',       // slate-800
  },
  NO_SHOW: {
    backgroundColor: '#ef4444', // red-500
    borderColor: '#dc2626',     // red-600
    textColor: '#fef2f2',       // red-50
  },
  CANCELED: {
    backgroundColor: '#64748b', // slate-500
    borderColor: '#475569',     // slate-600
    textColor: '#f1f5f9',       // slate-100
  },
};

// Labels de status em português
export const APPOINTMENT_STATUS_LABELS: Record<string, string> = {
  CREATED: 'Criado',
  CONFIRMED: 'Confirmado',
  CHECKED_IN: 'Cliente Chegou',
  IN_SERVICE: 'Em Atendimento',
  AWAITING_PAYMENT: 'Aguardando Pagamento',
  DONE: 'Concluído',
  NO_SHOW: 'Não Compareceu',
  CANCELED: 'Cancelado',
};

// Textos em português para o FullCalendar
export const FULLCALENDAR_LOCALE_PT_BR = {
  code: 'pt-br',
  week: {
    dow: 1, // Segunda-feira como primeiro dia
    doy: 4,
  },
  buttonText: {
    prev: 'Anterior',
    next: 'Próximo',
    today: 'Hoje',
    month: 'Mês',
    week: 'Semana',
    day: 'Dia',
    list: 'Lista',
  },
  weekText: 'Sm',
  allDayText: 'Dia inteiro',
  moreLinkText: (n: number) => `+${n} mais`,
  noEventsText: 'Nenhum agendamento',
};
