/**
 * API Client - Wrapper sobre axios para uso nos services
 * 
 * Re-exporta a instância configurada do axios como apiClient
 * para manter compatibilidade com services existentes.
 */

import { api } from './axios';

// Re-exporta a instância do axios como apiClient
export const apiClient = api;

// Também exporta como default
export default apiClient;
