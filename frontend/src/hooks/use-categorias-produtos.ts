/**
 * NEXO - Hooks para Categorias de Produtos
 * React Query hooks para gerenciamento de categorias de produtos
 * Módulo de Estoque v1.0
 */

import { categoriaProdutoService } from '@/services/categoria-produto-service';
import {
    CreateCategoriaProdutoRequest,
    UpdateCategoriaProdutoRequest,
} from '@/types/categoria-produto';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import { toast } from 'sonner';

export const CATEGORIAS_PRODUTOS_QUERY_KEY = ['categorias-produtos'];

/**
 * Hook para listar categorias de produtos
 * @param apenasAtivas - Se true, retorna apenas categorias ativas
 */
export function useCategoriasProdutos(apenasAtivas?: boolean) {
  return useQuery({
    queryKey: [...CATEGORIAS_PRODUTOS_QUERY_KEY, { apenasAtivas }],
    queryFn: () => categoriaProdutoService.list(apenasAtivas),
  });
}

/**
 * Hook para buscar uma categoria por ID
 */
export function useCategoriaProduto(id: string) {
  return useQuery({
    queryKey: [...CATEGORIAS_PRODUTOS_QUERY_KEY, id],
    queryFn: () => categoriaProdutoService.getById(id),
    enabled: !!id,
  });
}

/**
 * Hook para criar categoria de produto
 */
export function useCreateCategoriaProduto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCategoriaProdutoRequest) =>
      categoriaProdutoService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: CATEGORIAS_PRODUTOS_QUERY_KEY,
        exact: false,
      });
      toast.success('Categoria criada com sucesso!');
    },
    onError: (error: AxiosError<{ error: string }>) => {
      toast.error(
        error.response?.data?.error || 'Erro ao criar categoria'
      );
    },
  });
}

/**
 * Hook para atualizar categoria de produto
 */
export function useUpdateCategoriaProduto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCategoriaProdutoRequest }) =>
      categoriaProdutoService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: CATEGORIAS_PRODUTOS_QUERY_KEY,
        exact: false,
      });
      toast.success('Categoria atualizada com sucesso!');
    },
    onError: (error: AxiosError<{ error: string }>) => {
      toast.error(
        error.response?.data?.error || 'Erro ao atualizar categoria'
      );
    },
  });
}

/**
 * Hook para deletar categoria de produto
 */
export function useDeleteCategoriaProduto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => categoriaProdutoService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: CATEGORIAS_PRODUTOS_QUERY_KEY,
        exact: false,
      });
      toast.success('Categoria removida com sucesso!');
    },
    onError: (error: AxiosError<{ error: string }>) => {
      toast.error(
        error.response?.data?.error || 'Erro ao remover categoria'
      );
    },
  });
}

/**
 * Hook para alternar status ativa/inativa da categoria
 */
export function useToggleCategoriaProduto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => categoriaProdutoService.toggle(id),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: CATEGORIAS_PRODUTOS_QUERY_KEY,
        exact: false,
      });
      toast.success(
        data.ativa
          ? 'Categoria ativada com sucesso!'
          : 'Categoria desativada com sucesso!'
      );
    },
    onError: (error: AxiosError<{ error: string }>) => {
      toast.error(
        error.response?.data?.error || 'Erro ao alterar status da categoria'
      );
    },
  });
}
