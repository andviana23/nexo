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

// Mapeamento de status para cores
export const APPOINTMENT_STATUS_COLORS: Record<string, { backgroundColor: string; borderColor: string; textColor: string }> = {
  CREATED: {
    backgroundColor: '#FEF3C7', // Yellow-100
    borderColor: '#F59E0B',     // Amber-500
    textColor: '#92400E',       // Amber-800
  },
  CONFIRMED: {
    backgroundColor: '#DBEAFE', // Blue-100
    borderColor: '#3B82F6',     // Blue-500
    textColor: '#1E40AF',       // Blue-800
  },
  IN_SERVICE: {
    backgroundColor: '#D1FAE5', // Green-100
    borderColor: '#10B981',     // Emerald-500
    textColor: '#065F46',       // Emerald-800
  },
  DONE: {
    backgroundColor: '#E0E7FF', // Indigo-100
    borderColor: '#6366F1',     // Indigo-500
    textColor: '#3730A3',       // Indigo-800
  },
  NO_SHOW: {
    backgroundColor: '#FEE2E2', // Red-100
    borderColor: '#EF4444',     // Red-500
    textColor: '#991B1B',       // Red-800
  },
  CANCELED: {
    backgroundColor: '#F3F4F6', // Gray-100
    borderColor: '#9CA3AF',     // Gray-400
    textColor: '#4B5563',       // Gray-600
  },
};

// Labels de status em português
export const APPOINTMENT_STATUS_LABELS: Record<string, string> = {
  CREATED: 'Criado',
  CONFIRMED: 'Confirmado',
  IN_SERVICE: 'Em Atendimento',
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
