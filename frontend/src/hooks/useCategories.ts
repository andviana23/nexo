import { categoryService } from '@/services/category-service';
import {
    CategoryFilters,
    CreateCategoryDTO,
    UpdateCategoryDTO,
} from '@/types/category';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import { toast } from 'sonner';

export const CATEGORIES_QUERY_KEY = ['categories'];

export function useCategories(filters?: CategoryFilters) {
  return useQuery({
    queryKey: [...CATEGORIES_QUERY_KEY, filters],
    queryFn: () => categoryService.list(filters),
  });
}

export function useCategory(id: string) {
  return useQuery({
    queryKey: [...CATEGORIES_QUERY_KEY, id],
    queryFn: () => categoryService.getById(id),
    enabled: !!id,
  });
}

export function useCreateCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCategoryDTO) => categoryService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CATEGORIES_QUERY_KEY });
      toast.success('Categoria criada com sucesso!');
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao criar categoria'
      );
    },
  });
}

export function useUpdateCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCategoryDTO }) =>
      categoryService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CATEGORIES_QUERY_KEY });
      toast.success('Categoria atualizada com sucesso!');
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao atualizar categoria'
      );
    },
  });
}

export function useDeleteCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => categoryService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CATEGORIES_QUERY_KEY });
      toast.success('Categoria removida com sucesso!');
    },
    onError: (error: AxiosError<{ message: string }>) => {
      toast.error(
        error.response?.data?.message || 'Erro ao remover categoria'
      );
    },
  });
}
