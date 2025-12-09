/**
 * NEXO - Blocked Time Service
 * Serviço para bloqueio de horários na agenda
 */

import { apiClient } from '@/lib/api-client';
import type {
    BlockedTimeResponse,
    CreateBlockedTimeRequest,
    ListBlockedTimesRequest,
    ListBlockedTimesResponse,
} from '@/types/appointment';

const BASE_URL = '/blocked-times';

/**
 * Cria um novo bloqueio de horário
 */
export async function createBlockedTime(
  data: CreateBlockedTimeRequest
): Promise<BlockedTimeResponse> {
  const response = await apiClient.post<BlockedTimeResponse>(BASE_URL, data);
  return response.data;
}

/**
 * Lista bloqueios de horário com filtros opcionais
 */
export async function listBlockedTimes(
  params?: ListBlockedTimesRequest
): Promise<ListBlockedTimesResponse> {
  const response = await apiClient.get<ListBlockedTimesResponse>(BASE_URL, {
    params,
  });
  return response.data;
}

/**
 * Deleta um bloqueio de horário
 */
export async function deleteBlockedTime(id: string): Promise<void> {
  await apiClient.delete(`${BASE_URL}/${id}`);
}
