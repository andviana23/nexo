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
  FormDescription,
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
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { useCreateMeioPagamento, useUpdateMeioPagamento } from '@/hooks/use-meios-pagamento';
import {
  BANDEIRAS_CARTAO,
  CORES_TIPO_PAGAMENTO,
  CreateMeioPagamentoDTO,
  D_MAIS_PADRAO,
  MeioPagamento,
  TIPO_PAGAMENTO_LABELS,
  TipoPagamento,
} from '@/types/meio-pagamento';
import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';

// Schema de validação
const meioPagamentoSchema = z.object({
  nome: z.string().optional(),
  tipo: z.enum(['DINHEIRO', 'PIX', 'CREDITO', 'DEBITO', 'TRANSFERENCIA', 'BOLETO', 'OUTRO']),
  bandeira: z.string().optional(),
  taxa: z.string().optional(),
  taxa_fixa: z.string().optional(),
  d_mais: z.coerce.number().min(0, 'D+ deve ser maior ou igual a 0').max(365, 'D+ máximo é 365'),
  icone: z.string().optional(),
  cor: z.string().optional(),
  ordem_exibicao: z.coerce.number().min(0).optional(),
  observacoes: z.string().optional(),
  ativo: z.boolean().optional(),
});

type MeioPagamentoFormData = z.infer<typeof meioPagamentoSchema>;

interface MeioPagamentoModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  meioPagamento?: MeioPagamento | null;
}

export function MeioPagamentoModal({
  open,
  onOpenChange,
  meioPagamento,
}: MeioPagamentoModalProps) {
  const createMutation = useCreateMeioPagamento();
  const updateMutation = useUpdateMeioPagamento();

  const isEditing = !!meioPagamento;
  const isPending = createMutation.isPending || updateMutation.isPending;

  const form = useForm<MeioPagamentoFormData>({
    resolver: zodResolver(meioPagamentoSchema),
    defaultValues: {
      nome: '',
      tipo: 'DINHEIRO',
      bandeira: '',
      taxa: '0',
      taxa_fixa: '0',
      d_mais: 0,
      icone: '',
      cor: '',
      ordem_exibicao: 0,
      observacoes: '',
      ativo: true,
    },
  });

  const tipoSelecionado = useWatch({ control: form.control, name: 'tipo' }) as TipoPagamento;
  const mostrarBandeira = tipoSelecionado === 'CREDITO' || tipoSelecionado === 'DEBITO';

  // Atualiza D+ padrão quando muda o tipo (apenas em criação)
  useEffect(() => {
    if (!isEditing && tipoSelecionado) {
      form.setValue('d_mais', D_MAIS_PADRAO[tipoSelecionado] || 0);
      form.setValue('cor', CORES_TIPO_PAGAMENTO[tipoSelecionado] || '');

      // Limpa bandeira se não for cartão
      if (tipoSelecionado !== 'CREDITO' && tipoSelecionado !== 'DEBITO') {
        form.setValue('bandeira', '');
      }
    }
  }, [tipoSelecionado, isEditing, form]);

  // Carrega dados para edição
  useEffect(() => {
    if (meioPagamento) {
      form.reset({
        nome: meioPagamento.nome,
        tipo: meioPagamento.tipo as TipoPagamento,
        bandeira: meioPagamento.bandeira || '',
        taxa: meioPagamento.taxa || '0',
        taxa_fixa: meioPagamento.taxa_fixa || '0',
        d_mais: meioPagamento.d_mais,
        icone: meioPagamento.icone || '',
        cor: meioPagamento.cor || '',
        ordem_exibicao: meioPagamento.ordem_exibicao || 0,
        observacoes: meioPagamento.observacoes || '',
        ativo: meioPagamento.ativo,
      });
    } else {
      form.reset({
        nome: '',
        tipo: 'DINHEIRO',
        bandeira: '',
        taxa: '0',
        taxa_fixa: '0',
        d_mais: 0,
        icone: '',
        cor: CORES_TIPO_PAGAMENTO.DINHEIRO,
        ordem_exibicao: 0,
        observacoes: '',
        ativo: true,
      });
    }
  }, [meioPagamento, form]);

  const onSubmit = async (data: MeioPagamentoFormData) => {
    const payload: CreateMeioPagamentoDTO = {
      nome: data.nome || '',
      tipo: data.tipo,
      taxa: data.taxa || '0',
      taxa_fixa: data.taxa_fixa || '0',
      d_mais: data.d_mais || 0,
      ordem_exibicao: data.ordem_exibicao || 0,
      ativo: data.ativo ?? true,
    };

    // Adiciona campos opcionais apenas se tiverem valor
    if (data.bandeira && data.bandeira.trim() !== '') {
      payload.bandeira = data.bandeira.trim();
    }
    if (data.icone && data.icone.trim() !== '') {
      payload.icone = data.icone.trim();
    }
    if (data.cor && data.cor.trim() !== '') {
      payload.cor = data.cor.trim();
    }
    if (data.observacoes && data.observacoes.trim() !== '') {
      payload.observacoes = data.observacoes.trim();
    }

    try {
      if (isEditing && meioPagamento) {
        await updateMutation.mutateAsync({ id: meioPagamento.id, data: payload });
      } else {
        await createMutation.mutateAsync(payload);
      }
      onOpenChange(false);
    } catch {
      // Erro já tratado no hook
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Editar Tipo de Recebimento' : 'Novo Tipo de Recebimento'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Atualize as informações do meio de pagamento.'
              : 'Configure um novo meio de pagamento aceito pela sua barbearia.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 pt-4">

            {/* GRUPO 1: Informações Básicas */}
            <div className="space-y-4">
              <h4 className="text-sm font-medium text-muted-foreground uppercase tracking-wider text-[10px]">Informações Básicas</h4>

              <div className="grid gap-4 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="tipo"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tipo *</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione o tipo" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {Object.entries(TIPO_PAGAMENTO_LABELS).map(([value, label]) => (
                            <SelectItem key={value} value={value}>
                              {label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {mostrarBandeira && (
                  <FormField
                    control={form.control}
                    name="bandeira"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Bandeira</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Selecione a bandeira" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {BANDEIRAS_CARTAO.map((bandeira) => (
                              <SelectItem key={bandeira} value={bandeira}>
                                {bandeira}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                )}
              </div>
            </div>

            <Separator />

            {/* GRUPO 2: Financeiro */}
            <div className="space-y-4">
              <h4 className="text-sm font-medium text-muted-foreground uppercase tracking-wider text-[10px]">Configuração Financeira</h4>

              <div className="grid gap-4 sm:grid-cols-3">
                <FormField
                  control={form.control}
                  name="taxa"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Taxa (%)</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          step="0.01"
                          min="0"
                          max="100"
                          placeholder="0.00"
                          className="text-right"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="taxa_fixa"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Taxa Fixa (R$)</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          step="0.01"
                          min="0"
                          placeholder="0.00"
                          className="text-right"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="d_mais"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Recebimento (D+)</FormLabel>
                      <FormControl>
                        <Input type="number" min="0" max="365" placeholder="0" className="text-center" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>

            <Separator />

            {/* GRUPO 3: Visual & Outros */}
            <div className="space-y-4">
              <h4 className="text-sm font-medium text-muted-foreground uppercase tracking-wider text-[10px]">Aparência & Outros</h4>

              <div className="grid gap-4 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="cor"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Cor de Identificação</FormLabel>
                      <FormControl>
                        <div className="flex gap-2">
                          <Input
                            type="color"
                            className="w-10 h-9 p-0.5 cursor-pointer border-2"
                            {...field}
                          />
                          <Input
                            type="text"
                            placeholder="#000000"
                            className="flex-1 font-mono uppercase"
                            value={field.value}
                            onChange={field.onChange}
                          />
                        </div>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="ordem_exibicao"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Ordem</FormLabel>
                      <FormControl>
                        <Input type="number" min="0" placeholder="0" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name="observacoes"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Observações</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Anotações internas..."
                        className="resize-none"
                        rows={2}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="ativo"
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm bg-muted/20">
                    <div className="space-y-0.5">
                      <FormLabel className="text-sm font-semibold">Status Ativo</FormLabel>
                      <FormDescription className="text-xs">
                        Habilitar este meio de pagamento no sistema.
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>

            <DialogFooter className="pt-4">
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? 'Salvando...' : isEditing ? 'Salvar Alterações' : 'Criar Tipo'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
