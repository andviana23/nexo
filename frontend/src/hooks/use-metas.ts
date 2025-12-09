/**
 * NEXO - Sistema de Gestão para Barbearias
 * Metas Hooks
 *
 * React Query hooks para o módulo de Metas.
 */

import { metasService } from '@/services/metas-service';
import type {
    MetaBarbeiroResponse,
    MetaMensalResponse,
    MetasFilters,
    MetaTicketResponse,
    SetMetaBarbeiroRequest,
    SetMetaMensalRequest,
    SetMetaTicketRequest,
    UpdateMetaBarbeiroRequest,
    UpdateMetaMensalRequest,
    UpdateMetaTicketRequest,
} from '@/types/metas';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

// =============================================================================
// QUERY KEYS
// =============================================================================

export const metasKeys = {
  all: ['metas'] as const,
  // Metas Mensais
  mensais: () => [...metasKeys.all, 'mensais'] as const,
  mensaisList: (filters?: MetasFilters) => [...metasKeys.mensais(), 'list', filters] as const,
  mensalDetail: (id: string) => [...metasKeys.mensais(), 'detail', id] as const,
  // Metas Barbeiro
  barbeiros: () => [...metasKeys.all, 'barbeiros'] as const,
  barbeirosList: (barbeiroId?: string) => [...metasKeys.barbeiros(), 'list', barbeiroId] as const,
  barbeiroDetail: (id: string) => [...metasKeys.barbeiros(), 'detail', id] as const,
  barbeirosRanking: (mesAno?: string) => [...metasKeys.barbeiros(), 'ranking', mesAno] as const,
  // Metas Ticket
  ticket: () => [...metasKeys.all, 'ticket'] as const,
  ticketList: () => [...metasKeys.ticket(), 'list'] as const,
  ticketDetail: (id: string) => [...metasKeys.ticket(), 'detail', id] as const,
  // Resumo
  resumo: (mesAno?: string) => [...metasKeys.all, 'resumo', mesAno] as const,
  historico: (meses?: number) => [...metasKeys.all, 'historico', meses] as const,
};

// =============================================================================
// METAS MENSAIS HOOKS
// =============================================================================

/**
 * Hook para listar metas mensais
 */
export function useMetasMensais(filters?: MetasFilters) {
  return useQuery<MetaMensalResponse[], Error>({
    queryKey: metasKeys.mensaisList(filters),
    queryFn: () => metasService.listMetasMensais(filters),
  });
}

/**
 * Hook para buscar uma meta mensal por ID
 */
export function useMetaMensal(id: string) {
  return useQuery<MetaMensalResponse, Error>({
    queryKey: metasKeys.mensalDetail(id),
    queryFn: () => metasService.getMetaMensal(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar meta mensal
 */
export function useCreateMetaMensal() {
  const queryClient = useQueryClient();

  return useMutation<MetaMensalResponse, Error, SetMetaMensalRequest>({
    mutationFn: (payload) => metasService.createMetaMensal(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: metasKeys.mensais() });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta mensal criada com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao criar meta: ${error.message}`);
    },
  });
}

/**
 * Hook para atualizar meta mensal
 */
export function useUpdateMetaMensal() {
  const queryClient = useQueryClient();

  return useMutation<MetaMensalResponse, Error, { id: string; payload: UpdateMetaMensalRequest }>({
    mutationFn: ({ id, payload }) => metasService.updateMetaMensal(id, payload),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: metasKeys.mensais() });
      queryClient.invalidateQueries({ queryKey: metasKeys.mensalDetail(id) });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta mensal atualizada com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao atualizar meta: ${error.message}`);
    },
  });
}

/**
 * Hook para deletar meta mensal
 */
export function useDeleteMetaMensal() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) => metasService.deleteMetaMensal(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: metasKeys.mensais() });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta mensal excluída com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao excluir meta: ${error.message}`);
    },
  });
}

// =============================================================================
// METAS BARBEIRO HOOKS
// =============================================================================

/**
 * Hook para listar metas de barbeiros
 */
export function useMetasBarbeiro(barbeiroId?: string) {
  return useQuery<MetaBarbeiroResponse[], Error>({
    queryKey: metasKeys.barbeirosList(barbeiroId),
    queryFn: () => metasService.listMetasBarbeiro(barbeiroId),
  });
}

/**
 * Hook para buscar meta de barbeiro por ID
 */
export function useMetaBarbeiro(id: string) {
  return useQuery<MetaBarbeiroResponse, Error>({
    queryKey: metasKeys.barbeiroDetail(id),
    queryFn: () => metasService.getMetaBarbeiro(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar meta de barbeiro
 */
export function useCreateMetaBarbeiro() {
  const queryClient = useQueryClient();

  return useMutation<MetaBarbeiroResponse, Error, SetMetaBarbeiroRequest>({
    mutationFn: (payload) => metasService.createMetaBarbeiro(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: metasKeys.barbeiros() });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta de barbeiro criada com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao criar meta: ${error.message}`);
    },
  });
}

/**
 * Hook para atualizar meta de barbeiro
 */
export function useUpdateMetaBarbeiro() {
  const queryClient = useQueryClient();

  return useMutation<MetaBarbeiroResponse, Error, { id: string; payload: UpdateMetaBarbeiroRequest }>({
    mutationFn: ({ id, payload }) => metasService.updateMetaBarbeiro(id, payload),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: metasKeys.barbeiros() });
      queryClient.invalidateQueries({ queryKey: metasKeys.barbeiroDetail(id) });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta de barbeiro atualizada com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao atualizar meta: ${error.message}`);
    },
  });
}

/**
 * Hook para deletar meta de barbeiro
 */
export function useDeleteMetaBarbeiro() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) => metasService.deleteMetaBarbeiro(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: metasKeys.barbeiros() });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta de barbeiro excluída com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao excluir meta: ${error.message}`);
    },
  });
}

/**
 * Hook para buscar ranking de barbeiros
 */
export function useRankingBarbeiros(mesAno?: string) {
  return useQuery({
    queryKey: metasKeys.barbeirosRanking(mesAno),
    queryFn: () => metasService.getRankingBarbeiros(mesAno),
  });
}

// =============================================================================
// METAS TICKET MÉDIO HOOKS
// =============================================================================

/**
 * Hook para listar metas de ticket médio
 */
export function useMetasTicket() {
  return useQuery<MetaTicketResponse[], Error>({
    queryKey: metasKeys.ticketList(),
    queryFn: () => metasService.listMetasTicket(),
  });
}

/**
 * Hook para buscar meta de ticket por ID
 */
export function useMetaTicket(id: string) {
  return useQuery<MetaTicketResponse, Error>({
    queryKey: metasKeys.ticketDetail(id),
    queryFn: () => metasService.getMetaTicket(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar meta de ticket
 */
export function useCreateMetaTicket() {
  const queryClient = useQueryClient();

  return useMutation<MetaTicketResponse, Error, SetMetaTicketRequest>({
    mutationFn: (payload) => metasService.createMetaTicket(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: metasKeys.ticket() });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta de ticket médio criada com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao criar meta: ${error.message}`);
    },
  });
}

/**
 * Hook para atualizar meta de ticket
 */
export function useUpdateMetaTicket() {
  const queryClient = useQueryClient();

  return useMutation<MetaTicketResponse, Error, { id: string; payload: UpdateMetaTicketRequest }>({
    mutationFn: ({ id, payload }) => metasService.updateMetaTicket(id, payload),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: metasKeys.ticket() });
      queryClient.invalidateQueries({ queryKey: metasKeys.ticketDetail(id) });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta de ticket médio atualizada com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao atualizar meta: ${error.message}`);
    },
  });
}

/**
 * Hook para deletar meta de ticket
 */
export function useDeleteMetaTicket() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) => metasService.deleteMetaTicket(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: metasKeys.ticket() });
      queryClient.invalidateQueries({ queryKey: metasKeys.resumo() });
      toast.success('Meta de ticket médio excluída com sucesso!');
    },
    onError: (error) => {
      toast.error(`Erro ao excluir meta: ${error.message}`);
    },
  });
}

// =============================================================================
// HOOKS AGREGADORES
// =============================================================================

/**
 * Hook para buscar resumo geral de metas
 */
export function useResumoMetas(mesAno?: string) {
  return useQuery({
    queryKey: metasKeys.resumo(mesAno),
    queryFn: () => metasService.getResumoMetas(mesAno),
  });
}

/**
 * Hook para buscar histórico de metas
 */
export function useHistoricoMetas(meses: number = 6) {
  return useQuery({
    queryKey: metasKeys.historico(meses),
    queryFn: () => metasService.getHistoricoMetas(meses),
  });
}
