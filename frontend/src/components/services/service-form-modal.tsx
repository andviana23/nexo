'use client';

import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { useCreateService, useUpdateService } from '@/hooks/useServices';
import { api } from '@/lib/axios';
import { Category } from '@/types/category';
import { Service } from '@/types/service';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import * as z from 'zod';

// Schema de validação
const formSchema = z.object({
  nome: z.string().min(2, 'Nome deve ter pelo menos 2 caracteres'),
  descricao: z.string().optional(),
  preco: z.string().min(1, 'Preço é obrigatório'),
  duracao: z.coerce.number().min(5, 'Duração mínima de 5 minutos'),
  comissao: z.string().optional(),
  categoria_id: z.string().optional(),
  cor: z.string().optional(),
  observacoes: z.string().optional(),
});

interface ServiceFormModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  serviceToEdit?: Service | null;
}

export function ServiceFormModal({
  open,
  onOpenChange,
  serviceToEdit,
}: ServiceFormModalProps) {
  const createService = useCreateService();
  const updateService = useUpdateService();

  // Buscar categorias para o select
  const { data: categoriesData } = useQuery({
    queryKey: ['categories', 'list'],
    queryFn: async () => {
      const { data } = await api.get('/categorias-servicos');
      return data;
    },
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      nome: '',
      descricao: '',
      preco: '',
      duracao: 30,
      comissao: '',
      categoria_id: '',
      cor: '#000000',
      observacoes: '',
    },
  });

  useEffect(() => {
    if (serviceToEdit) {
      form.reset({
        nome: serviceToEdit.nome,
        descricao: serviceToEdit.descricao || '',
        preco: serviceToEdit.preco,
        duracao: serviceToEdit.duracao,
        comissao: serviceToEdit.comissao,
        categoria_id: serviceToEdit.categoria_id || '',
        cor: serviceToEdit.cor || '#000000',
        observacoes: serviceToEdit.observacoes || '',
      });
    } else {
      form.reset({
        nome: '',
        descricao: '',
        preco: '',
        duracao: 30,
        comissao: '',
        categoria_id: '',
        cor: '#000000',
        observacoes: '',
      });
    }
  }, [serviceToEdit, form, open]);

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    const data = {
      ...values,
      // Limpar campos opcionais vazios
      categoria_id: values.categoria_id || undefined,
      descricao: values.descricao || undefined,
      comissao: values.comissao || undefined,
      cor: values.cor || undefined,
      observacoes: values.observacoes || undefined,
    };

    if (serviceToEdit) {
      updateService.mutate(
        { id: serviceToEdit.id, data },
        {
          onSuccess: () => {
            onOpenChange(false);
            form.reset();
          },
        }
      );
    } else {
      createService.mutate(data, {
        onSuccess: () => {
          onOpenChange(false);
          form.reset();
        },
      });
    }
  };

  const isLoading = createService.isPending || updateService.isPending;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>
            {serviceToEdit ? 'Editar Serviço' : 'Novo Serviço'}
          </DialogTitle>
          <DialogDescription>
            Preencha os dados do serviço abaixo.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="nome"
                render={({ field }) => (
                  <FormItem className="col-span-2">
                    <FormLabel>Nome</FormLabel>
                    <FormControl>
                      <Input placeholder="Ex: Corte de Cabelo" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="categoria_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Categoria</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Selecione..." />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {categoriesData?.categorias?.map((cat: Category) => (
                          <SelectItem key={cat.id} value={cat.id}>
                            {cat.nome}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="cor"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Cor</FormLabel>
                    <div className="flex gap-2">
                      <FormControl>
                        <Input type="color" className="w-12 p-1 h-10" {...field} />
                      </FormControl>
                      <Input {...field} placeholder="#000000" />
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="preco"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Preço (R$)</FormLabel>
                    <FormControl>
                      <Input type="number" step="0.01" placeholder="0.00" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="comissao"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Comissão (%)</FormLabel>
                    <FormControl>
                      <Input type="number" step="0.01" placeholder="0.00" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="duracao"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Duração (minutos)</FormLabel>
                    <FormControl>
                      <Input type="number" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="descricao"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Descrição</FormLabel>
                  <FormControl>
                    <Textarea placeholder="Detalhes do serviço..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Salvando...' : 'Salvar'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
