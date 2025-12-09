/**
 * Serviço de API para Comandas (Fechamento de Conta)
 * 
 * Integra com backend Go para gerenciar comandas, itens e pagamentos.
 */

import { apiClient } from '@/lib/api-client';
import { getErrorMessage, isAxiosError } from '@/lib/axios';
import type { AppointmentResponse, AppointmentServiceResponse } from '@/types/appointment';
import type {
    AddCommandItemRequest,
    AddCommandPaymentRequest,
    CloseCommandRequest,
    CommandItemResponse,
    CommandPaymentResponse,
    CommandResponse,
    CreateCommandRequest,
    UpdateCommandItemRequest
} from '@/types/command';

const BASE_PATH = '/commands';

// ============================================================================
// Comandas
// ============================================================================

/**
 * Criar nova comanda
 */
export async function createCommand(data: CreateCommandRequest): Promise<CommandResponse> {
  const response = await apiClient.post<CommandResponse>(BASE_PATH, data);
  return response.data;
}

/**
 * Buscar comanda por ID (com items e payments)
 */
export async function getCommand(commandId: string): Promise<CommandResponse> {
  const response = await apiClient.get<CommandResponse>(`${BASE_PATH}/${commandId}`);
  return response.data;
}

/**
 * Buscar comanda por appointment_id
 */
export async function getCommandByAppointment(appointmentId: string): Promise<CommandResponse | null> {
  try {
    const response = await apiClient.get<CommandResponse>(`${BASE_PATH}/by-appointment/${appointmentId}`);
    return response.data;
  } catch (error) {
    if (isAxiosError(error) && error.response?.status === 404) {
      return null;
    }
    throw new Error(getErrorMessage(error));
  }
}

/**
 * Listar comandas com filtros
 */
export interface ListCommandsFilters {
  customer_id?: string;
  status?: 'OPEN' | 'CLOSED' | 'CANCELED';
  data_inicio?: string;
  data_fim?: string;
  page?: number;
  limit?: number;
}

export async function listCommands(filters?: ListCommandsFilters): Promise<{
  commands: CommandResponse[];
  total: number;
  page: number;
  limit: number;
}> {
  const response = await apiClient.get(BASE_PATH, { params: filters });
  return response.data;
}

// ============================================================================
// Itens da Comanda
// ============================================================================

/**
 * Adicionar item à comanda
 */
export async function addCommandItem(
  commandId: string,
  data: AddCommandItemRequest
): Promise<CommandItemResponse> {
  const response = await apiClient.post<CommandItemResponse>(
    `${BASE_PATH}/${commandId}/items`,
    data
  );
  return response.data;
}

/**
 * Atualizar item da comanda
 */
export async function updateCommandItem(
  commandId: string,
  itemId: string,
  data: UpdateCommandItemRequest
): Promise<CommandItemResponse> {
  const response = await apiClient.patch<CommandItemResponse>(
    `${BASE_PATH}/${commandId}/items/${itemId}`,
    data
  );
  return response.data;
}

/**
 * Remover item da comanda
 */
export async function removeCommandItem(
  commandId: string,
  itemId: string
): Promise<void> {
  await apiClient.delete(`${BASE_PATH}/${commandId}/items/${itemId}`);
}

// ============================================================================
// Pagamentos da Comanda
// ============================================================================

/**
 * Adicionar pagamento à comanda
 */
export async function addCommandPayment(
  commandId: string,
  data: AddCommandPaymentRequest
): Promise<CommandPaymentResponse> {
  const response = await apiClient.post<CommandPaymentResponse>(
    `${BASE_PATH}/${commandId}/payments`,
    data
  );
  return response.data;
}

/**
 * Remover pagamento da comanda
 */
export async function removeCommandPayment(
  commandId: string,
  paymentId: string
): Promise<void> {
  await apiClient.delete(`${BASE_PATH}/${commandId}/payments/${paymentId}`);
}

// ============================================================================
// Fechamento
// ============================================================================

/**
 * Fechar comanda (finalizar conta)
 */
export async function closeCommand(
  commandId: string,
  data: CloseCommandRequest
): Promise<CommandResponse> {
  const response = await apiClient.post<CommandResponse>(
    `${BASE_PATH}/${commandId}/close`,
    data
  );
  return response.data;
}

/**
 * Cancelar comanda
 */
export async function cancelCommand(commandId: string, motivo?: string): Promise<CommandResponse> {
  const response = await apiClient.post<CommandResponse>(
    `${BASE_PATH}/${commandId}/cancel`,
    { motivo }
  );
  return response.data;
}

// ============================================================================
// Helpers
// ============================================================================

/**
 * Criar comanda a partir de appointment
 * 
 * Esta função busca o appointment pelo ID e cria uma comanda com os serviços associados.
 * Se o appointment já tiver uma comanda vinculada (command_id), busca e retorna a existente.
 */
export async function createCommandFromAppointment(appointmentId: string): Promise<CommandResponse> {
  console.log('[createCommandFromAppointment] Iniciando para appointment:', appointmentId);
  
  // Primeiro, verificar se já existe uma comanda para este appointment
  const existingCommand = await getCommandByAppointment(appointmentId);
  if (existingCommand) {
    console.log('[createCommandFromAppointment] Comanda já existe:', existingCommand.id);
    return existingCommand;
  }
  
  // Buscar appointment para obter dados
  // Nota: apiClient já tem /api/v1 como baseURL, então usamos apenas /appointments
  console.log('[createCommandFromAppointment] Buscando dados do appointment...');
  const appointmentResponse = await apiClient.get<AppointmentResponse>(`/appointments/${appointmentId}`);
  const appointment = appointmentResponse.data;
  
  console.log('[createCommandFromAppointment] Appointment encontrado:', {
    id: appointment.id,
    customer_id: appointment.customer_id,
    services_count: appointment.services?.length || 0,
  });
  
  // Criar comanda com serviços do appointment
  // Nota: preco_unitario deve ser STRING conforme esperado pelo backend
  const items = (appointment.services || []).map((service: AppointmentServiceResponse) => ({
    tipo: 'SERVICO' as const,
    item_id: service.service_id,
    descricao: service.service_name || 'Serviço',
    preco_unitario: String(service.price || '0'),
    quantidade: 1,
  }));
  
  console.log('[createCommandFromAppointment] Criando comanda com', items.length, 'itens');
  
  const command = await createCommand({
    appointment_id: appointmentId,
    customer_id: appointment.customer_id,
    items,
  });
  
  console.log('[createCommandFromAppointment] Comanda criada:', command.id);
  return command;
}
