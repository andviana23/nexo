/**
 * NEXO - Sistema de Gestão para Barbearias
 * Auth Service
 *
 * Serviço de autenticação - comunicação com API de auth do backend.
 */

import { api } from '@/lib/axios';
import type {
  LoginCredentials,
  LoginResponse,
  RefreshTokenResponse,
  User,
} from '@/types';

// =============================================================================
// ENDPOINTS
// =============================================================================

const AUTH_ENDPOINTS = {
  login: '/auth/login',
  logout: '/auth/logout',
  me: '/auth/me',
  refresh: '/auth/refresh',
  forgotPassword: '/auth/forgot-password',
  resetPassword: '/auth/reset-password',
  changePassword: '/auth/change-password',
} as const;

// =============================================================================
// SERVIÇO
// =============================================================================

export const authService = {
  /**
   * Realiza login com email e senha
   */
  async login(credentials: LoginCredentials): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>(
      AUTH_ENDPOINTS.login,
      credentials
    );
    return response.data;
  },

  /**
   * Realiza logout no servidor
   */
  async logout(): Promise<void> {
    try {
      await api.post(AUTH_ENDPOINTS.logout);
    } catch {
      // Ignora erro no logout - limpa local de qualquer forma
    }
  },

  /**
   * Obtém dados do usuário autenticado
   */
  async getMe(): Promise<User> {
    const response = await api.get<User>(AUTH_ENDPOINTS.me);
    return response.data;
  },

  /**
   * Renova o token de autenticação
   */
  async refreshToken(): Promise<RefreshTokenResponse> {
    const response = await api.post<RefreshTokenResponse>(
      AUTH_ENDPOINTS.refresh
    );
    return response.data;
  },

  /**
   * Solicita recuperação de senha
   */
  async forgotPassword(email: string): Promise<void> {
    await api.post(AUTH_ENDPOINTS.forgotPassword, { email });
  },

  /**
   * Redefine a senha com token de recuperação
   */
  async resetPassword(token: string, newPassword: string): Promise<void> {
    await api.post(AUTH_ENDPOINTS.resetPassword, {
      token,
      new_password: newPassword,
    });
  },

  /**
   * Altera a senha do usuário logado
   */
  async changePassword(
    currentPassword: string,
    newPassword: string
  ): Promise<void> {
    await api.post(AUTH_ENDPOINTS.changePassword, {
      current_password: currentPassword,
      new_password: newPassword,
    });
  },
};

// =============================================================================
// TIPOS DE ERRO ESPECÍFICOS
// =============================================================================

export class AuthenticationError extends Error {
  constructor(message: string, public code: string = 'AUTH_ERROR') {
    super(message);
    this.name = 'AuthenticationError';
  }
}

export class InvalidCredentialsError extends AuthenticationError {
  constructor() {
    super('Email ou senha inválidos', 'INVALID_CREDENTIALS');
    this.name = 'InvalidCredentialsError';
  }
}

export class SessionExpiredError extends AuthenticationError {
  constructor() {
    super('Sua sessão expirou. Faça login novamente.', 'SESSION_EXPIRED');
    this.name = 'SessionExpiredError';
  }
}

export class AccountLockedError extends AuthenticationError {
  constructor() {
    super(
      'Conta temporariamente bloqueada. Tente novamente mais tarde.',
      'ACCOUNT_LOCKED'
    );
    this.name = 'AccountLockedError';
  }
}

export class AccountInactiveError extends AuthenticationError {
  constructor() {
    super('Conta inativa. Entre em contato com o suporte.', 'ACCOUNT_INACTIVE');
    this.name = 'AccountInactiveError';
  }
}
