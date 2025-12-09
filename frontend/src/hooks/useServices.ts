import { getErrorMessage } from '@/lib/axios';
import { serviceService } from '@/services/serviceService';
import { CreateServiceDTO, ServiceFilters, UpdateServiceDTO } from '@/types/service';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';

export const SERVICE_KEYS = {
  all: ['services'] as const,
  lists: () => [...SERVICE_KEYS.all, 'list'] as const,
  list: (filters: ServiceFilters) => [...SERVICE_KEYS.lists(), filters] as const,
  details: () => [...SERVICE_KEYS.all, 'detail'] as const,
  detail: (id: string) => [...SERVICE_KEYS.details(), id] as const,
  stats: () => [...SERVICE_KEYS.all, 'stats'] as const,
};

export function useServices(filters: ServiceFilters = {}) {
  return useQuery({
    queryKey: SERVICE_KEYS.list(filters),
    queryFn: () => serviceService.getAll(filters),
  });
}

export function useService(id: string) {
  return useQuery({
    queryKey: SERVICE_KEYS.detail(id),
    queryFn: () => serviceService.getById(id),
    enabled: !!id,
  });
}

export function useServiceStats() {
  return useQuery({
    queryKey: SERVICE_KEYS.stats(),
    queryFn: () => serviceService.getStats(),
  });
}

export function useCreateService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateServiceDTO) => serviceService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.stats() });
      toast.success('Serviço criado com sucesso!');
    },
    onError: (error) => {
      toast.error(getErrorMessage(error));
    },
  });
}

export function useUpdateService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateServiceDTO }) =>
      serviceService.update(id, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.stats() });
      toast.success('Serviço atualizado com sucesso!');
    },
    onError: (error) => {
      toast.error(getErrorMessage(error));
    },
  });
}

export function useDeleteService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => serviceService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.stats() });
      toast.success('Serviço removido com sucesso!');
    },
    onError: (error) => {
      toast.error(getErrorMessage(error));
    },
  });
}

export function useToggleServiceStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => serviceService.toggleStatus(id),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.lists() });
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: SERVICE_KEYS.stats() });
      toast.success(`Serviço ${data.ativo ? 'ativado' : 'desativado'} com sucesso!`);
    },
    onError: (error) => {
      toast.error(getErrorMessage(error));
    },
  });
}
