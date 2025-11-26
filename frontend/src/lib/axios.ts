import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

// Instância do Axios configurada para a API
export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor de Request - Adiciona token de autenticação
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Só executa no client-side
    if (typeof window !== 'undefined') {
      // Busca token do Zustand persist storage
      const authData = localStorage.getItem('nexo-auth');
      if (authData) {
        try {
          const { state } = JSON.parse(authData);
          const token = state?.token;
          if (token && config.headers) {
            config.headers.Authorization = `Bearer ${token}`;
          }
        } catch (error) {
          console.error('Erro ao parsear token:', error);
        }
      }
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Interceptor de Response - Tratamento de erros globais
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // Erro 401 - Não autorizado
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      // Limpa token e redireciona para login
      if (typeof window !== 'undefined') {
        localStorage.removeItem('nexo-auth');

        // Só redireciona se não estiver na página de login
        if (!window.location.pathname.includes('/login')) {
          window.location.href = '/login';
        }
      }
    }

    // Erro 403 - Forbidden (sem permissão)
    if (error.response?.status === 403) {
      console.error('Acesso negado:', error.response.data);
    }

    // Erro 500 - Erro interno do servidor
    if (error.response?.status === 500) {
      console.error('Erro interno do servidor:', error.response.data);
    }

    return Promise.reject(error);
  }
);

// Helper para verificar se é erro do Axios
export function isAxiosError(error: unknown): error is AxiosError {
  return axios.isAxiosError(error);
}

// Helper para extrair mensagem de erro
export function getErrorMessage(error: unknown): string {
  if (isAxiosError(error)) {
    const data = error.response?.data as
      | { message?: string; error?: string }
      | undefined;
    return data?.message || data?.error || error.message || 'Erro desconhecido';
  }
  if (error instanceof Error) {
    return error.message;
  }
  return 'Erro desconhecido';
}

export default api;
