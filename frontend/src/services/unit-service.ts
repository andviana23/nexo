/**
 * NEXO - Sistema de Gestão para Barbearias
 * Unit Service
 *
 * Serviço de unidades/filiais - comunicação com API de units do backend.
 */

import { api } from '@/lib/axios';
import type {
    AddUserToUnitRequest,
    CreateUnitRequest,
    ListUnitsResponse,
    ListUserUnitsResponse,
    SetDefaultUnitRequest,
    SwitchUnitResponse,
    Unit,
    UpdateUnitRequest,
    UserUnit,
} from '@/types/unit';

// =============================================================================
// ENDPOINTS
// =============================================================================

const UNIT_ENDPOINTS = {
  // Unidades do usuário logado
  userUnits: '/units/me',
  switchUnit: '/units/switch',
  setDefault: '/units/default',

  // CRUD de unidades (admin)
  list: '/units',
  detail: (id: string) => `/units/${id}`,
  create: '/units',
  update: (id: string) => `/units/${id}`,
  delete: (id: string) => `/units/${id}`,
  activate: (id: string) => `/units/${id}/activate`,
  deactivate: (id: string) => `/units/${id}/deactivate`,

  // Vincular usuários a unidades (admin)
  addUser: '/units/users',
  removeUser: (unitId: string, userId: string) => `/units/${unitId}/users/${userId}`,
} as const;

// =============================================================================
// SERVIÇO
// =============================================================================

export const unitService = {
  // ===========================================================================
  // OPERAÇÕES DO USUÁRIO LOGADO
  // ===========================================================================

  /**
   * Lista unidades às quais o usuário tem acesso
   */
  async getUserUnits(): Promise<ListUserUnitsResponse> {
    const response = await api.get<ListUserUnitsResponse>(UNIT_ENDPOINTS.userUnits);
    return response.data;
  },

  /**
   * Troca a unidade ativa do usuário
   * Pode retornar um novo token com a unit_id atualizada
   */
  async switchUnit(unitId: string): Promise<SwitchUnitResponse> {
    const response = await api.post<SwitchUnitResponse>(UNIT_ENDPOINTS.switchUnit, {
      unit_id: unitId,
    });
    return response.data;
  },

  /**
   * Define a unidade padrão do usuário
   */
  async setDefaultUnit(unitId: string): Promise<UserUnit> {
    const response = await api.post<UserUnit>(UNIT_ENDPOINTS.setDefault, {
      unit_id: unitId,
    } as SetDefaultUnitRequest);
    return response.data;
  },

  // ===========================================================================
  // CRUD DE UNIDADES (ADMIN)
  // ===========================================================================

  /**
   * Lista todas as unidades do tenant (admin)
   */
  async listUnits(): Promise<ListUnitsResponse> {
    const response = await api.get<ListUnitsResponse>(UNIT_ENDPOINTS.list);
    return response.data;
  },

  /**
   * Obtém detalhes de uma unidade específica
   */
  async getUnit(id: string): Promise<Unit> {
    const response = await api.get<Unit>(UNIT_ENDPOINTS.detail(id));
    return response.data;
  },

  /**
   * Cria uma nova unidade (admin)
   */
  async createUnit(data: CreateUnitRequest): Promise<Unit> {
    const response = await api.post<Unit>(UNIT_ENDPOINTS.create, data);
    return response.data;
  },

  /**
   * Atualiza uma unidade existente (admin)
   */
  async updateUnit(id: string, data: UpdateUnitRequest): Promise<Unit> {
    const response = await api.patch<Unit>(UNIT_ENDPOINTS.update(id), data);
    return response.data;
  },

  /**
   * Remove uma unidade (soft delete) (admin)
   */
  async deleteUnit(id: string): Promise<void> {
    await api.delete(UNIT_ENDPOINTS.delete(id));
  },

  /**
   * Ativa uma unidade desativada (admin)
   */
  async activateUnit(id: string): Promise<Unit> {
    const response = await api.post<Unit>(UNIT_ENDPOINTS.activate(id));
    return response.data;
  },

  /**
   * Desativa uma unidade (admin)
   */
  async deactivateUnit(id: string): Promise<Unit> {
    const response = await api.post<Unit>(UNIT_ENDPOINTS.deactivate(id));
    return response.data;
  },

  // ===========================================================================
  // GESTÃO DE USUÁRIOS EM UNIDADES (ADMIN)
  // ===========================================================================

  /**
   * Adiciona um usuário a uma unidade (admin)
   */
  async addUserToUnit(data: AddUserToUnitRequest): Promise<UserUnit> {
    const response = await api.post<UserUnit>(UNIT_ENDPOINTS.addUser, data);
    return response.data;
  },

  /**
   * Remove um usuário de uma unidade (admin)
   */
  async removeUserFromUnit(unitId: string, userId: string): Promise<void> {
    await api.delete(UNIT_ENDPOINTS.removeUser(unitId, userId));
  },
};

export default unitService;
