'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Professional Modal
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
  BriefcaseIcon,
  CalendarIcon,
  DollarSignIcon,
  Loader2Icon,
  MailIcon,
  PercentIcon,
  PhoneIcon,
  UserIcon
} from 'lucide-react';
import { useCallback, useEffect, useMemo } from 'react';
import { useFieldArray, useForm, useWatch } from 'react-hook-form';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
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
  FormMessage
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';

import {
  useCreateProfessional,
  useUpdateProfessional,
} from '@/hooks/use-professionals';
import { useCategories } from '@/hooks/useCategories';

import {
  createProfessionalSchema,
  maskCPF,
  maskPhone,
  unmaskCPF,
  unmaskPhone,
  type CreateProfessionalFormData,
} from '@/lib/validations/professional';
import type {
  ProfessionalModalState,
  ProfessionalResponse,
} from '@/types/professional';

import {
  ProfessionalStatusBadge,
  ProfessionalTypeBadge,
} from './ProfessionalBadge';

// =============================================================================
// TYPES
// =============================================================================

interface ProfessionalModalProps {
  state: ProfessionalModalState;
  onClose: () => void;
  onSuccess?: (professional: ProfessionalResponse) => void;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ProfessionalModal({
  state,
  onClose,
  onSuccess,
}: ProfessionalModalProps) {
  const { isOpen, mode, professional } = state;

  // Mutations
  const createProfessional = useCreateProfessional();
  const updateProfessional = useUpdateProfessional();

  // Data for Category Commissions
  // Data for Category Commissions
  const { data: categories = [] } = useCategories({ apenas_ativas: true });

  const isLoading = createProfessional.isPending || updateProfessional.isPending;
  const isViewMode = mode === 'view';

  // Form
  const form = useForm<CreateProfessionalFormData>({
    resolver: zodResolver(createProfessionalSchema),
    defaultValues: {
      nome: '',
      email: '',
      telefone: '',
      cpf: '',
      tipo: 'BARBEIRO',
      data_admissao: new Date(),
      foto: '',
      especialidades: [],
      observacoes: '',
      tambem_barbeiro: false,
      tipo_comissao: 'PERCENTUAL',
      comissao: undefined,
      comissao_produtos: undefined,
      comissoes_por_categoria: [],
    },
  });

  const { fields, append, remove, replace } = useFieldArray({
    control: form.control,
    name: 'comissoes_por_categoria',
  });

  const tipoWatch = useWatch({ control: form.control, name: 'tipo' });
  const tambemBarbeiroWatch = useWatch({ control: form.control, name: 'tambem_barbeiro' });

  // Mostra campos de comissão se for Barbeiro ou Gerente+Barbeiro
  const showCommissionFields = useMemo(() => {
    return tipoWatch === 'BARBEIRO' || (tipoWatch === 'GERENTE' && tambemBarbeiroWatch);
  }, [tipoWatch, tambemBarbeiroWatch]);

  // Sync category fields when categories load or modal opens
  useEffect(() => {
    if (categories.length > 0 && showCommissionFields) {
      // Initialize with existing categories if not present
      const currentValues = form.getValues('comissoes_por_categoria') || [];
      const newValues = categories.map(cat => {
        const existing = currentValues.find(c => c.categoria_id === cat.id);
        return existing || { categoria_id: cat.id, comissao: 0 };
      });
      // We don't want to overwrite if user has typed something, but simplified logic:
      // If empty, fill it. If distinct length, merge.
      if (currentValues.length === 0) {
        replace(newValues);
      }
    }
  }, [categories, showCommissionFields, replace, form]);

  useEffect(() => {
    if (mode === 'edit' && professional) {
      form.reset({
        nome: professional.nome,
        email: professional.email,
        telefone: professional.telefone,
        cpf: professional.cpf,
        tipo: professional.tipo,
        data_admissao: new Date(professional.data_admissao),
        foto: professional.foto || '',
        especialidades: professional.especialidades || [],
        observacoes: professional.observacoes || '',
        tambem_barbeiro: false,
        tipo_comissao: professional.tipo_comissao || 'PERCENTUAL',
        comissao: professional.comissao ? Number(professional.comissao) : undefined,
        comissao_produtos: professional.comissao_produtos ? Number(professional.comissao_produtos) : undefined,
        // TODO: Load saved category commissions if backend supported it
        comissoes_por_categoria: [],
      });
    } else if (mode === 'create') {
      form.reset({
        nome: '',
        email: '',
        telefone: '',
        cpf: '',
        tipo: 'BARBEIRO',
        data_admissao: new Date(),
        foto: '',
        especialidades: [],
        observacoes: '',
        tambem_barbeiro: false,
        tipo_comissao: 'PERCENTUAL',
        comissao: undefined,
        comissao_produtos: undefined,
        comissoes_por_categoria: [],
      });
    }
  }, [mode, professional, form]);

  const onSubmit = useCallback(
    async (values: CreateProfessionalFormData) => {
      const { tambem_barbeiro: _tambemBarbeiro, ...rest } = values;
      void _tambemBarbeiro;

      const data = {
        ...rest,
        cpf: unmaskCPF(rest.cpf || ''),
        telefone: unmaskPhone(rest.telefone),
        data_admissao: format(rest.data_admissao, 'yyyy-MM-dd'),
        comissao: rest.comissao ?? 0,
      };

      if (mode === 'create') {
        createProfessional.mutate(data as Parameters<typeof createProfessional.mutate>[0], {
          onSuccess: (result) => {
            onSuccess?.(result);
            onClose();
          },
        });
      } else if (mode === 'edit' && professional) {
        updateProfessional.mutate(
          {
            id: professional.id,
            data: {
              nome: data.nome,
              email: data.email,
              telefone: data.telefone,
              foto: data.foto || undefined,
              especialidades: data.especialidades,
              observacoes: data.observacoes,
              tipo_comissao: data.tipo_comissao,
              comissao: data.comissao ?? undefined,
              comissao_produtos: data.comissao_produtos ?? undefined,
              comissoes_por_categoria: data.comissoes_por_categoria,
            },
          },
          {
            onSuccess: (result) => {
              onSuccess?.(result);
              onClose();
            },
          }
        );
      }
    },
    [mode, professional, createProfessional, updateProfessional, onSuccess, onClose]
  );

  const title = useMemo(() => {
    switch (mode) {
      case 'create': return 'Novo Profissional';
      case 'edit': return 'Editar Profissional';
      case 'view': return 'Detalhes do Profissional';
      default: return 'Profissional';
    }
  }, [mode]);

  // VIEW MODE
  if (isViewMode && professional) {
    return (
      <Dialog open={isOpen} onOpenChange={() => onClose()}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <UserIcon className="size-5 text-primary" />
              {title}
            </DialogTitle>
            <DialogDescription>Visualização completa do perfil</DialogDescription>
          </DialogHeader>

          <div className="space-y-6 pt-2">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-xl font-bold">{professional.nome}</h3>
                <div className="flex items-center gap-2 mt-2">
                  <ProfessionalTypeBadge type={professional.tipo} />
                  <ProfessionalStatusBadge status={professional.status} />
                </div>
              </div>
            </div>

            <Separator />

            <div className="space-y-4">
              <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Contato & Informações</h4>
              <div className="grid gap-3">
                <div className="flex items-center gap-3 text-sm">
                  <MailIcon className="size-4 text-muted-foreground" />
                  <span>{professional.email}</span>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <PhoneIcon className="size-4 text-muted-foreground" />
                  <span>{maskPhone(professional.telefone)}</span>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <CalendarIcon className="size-4 text-muted-foreground" />
                  <span>Admitido em {format(new Date(professional.data_admissao), "dd 'de' MMM 'de' yyyy", { locale: ptBR })}</span>
                </div>
              </div>
            </div>

            {professional.comissao && (
              <>
                <Separator />
                <div className="space-y-4">
                  <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Comissionamento</h4>
                  <div className="flex gap-4">
                    <div className="flex flex-col">
                      <span className="text-xs text-muted-foreground">Serviços (Base)</span>
                      <span className="font-mono font-medium text-lg">{professional.comissao}%</span>
                    </div>
                    {professional.comissao_produtos && (
                      <div className="flex flex-col">
                        <span className="text-xs text-muted-foreground">Produtos</span>
                        <span className="font-mono font-medium text-lg">{professional.comissao_produtos}%</span>
                      </div>
                    )}
                  </div>
                </div>
              </>
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={onClose}>Fechar</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  // CREATE / EDIT MODE
  return (
    <Dialog open={isOpen} onOpenChange={() => onClose()}>
      <DialogContent className="sm:max-w-[700px] h-[90vh] sm:h-auto overflow-hidden flex flex-col p-0">
        <DialogHeader className="px-6 py-4 border-b">
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            {mode === 'create' ? 'Informe os dados para cadastrar.' : 'Atualize as informações do profissional.'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex-1 overflow-hidden flex flex-col">
            <ScrollArea className="flex-1 p-6">
              <Tabs defaultValue="dados-pessoais" className="w-full">
                <TabsList className="grid w-full grid-cols-3 mb-6">
                  <TabsTrigger value="dados-pessoais">
                    <UserIcon className="w-4 h-4 mr-2" />
                    Pessoais
                  </TabsTrigger>
                  <TabsTrigger value="profissional">
                    <BriefcaseIcon className="w-4 h-4 mr-2" />
                    Profissional
                  </TabsTrigger>
                  <TabsTrigger value="financeiro" disabled={!showCommissionFields}>
                    <DollarSignIcon className="w-4 h-4 mr-2" />
                    Comissões
                  </TabsTrigger>
                </TabsList>

                {/* TAB: DADOS PESSOAIS */}
                <TabsContent value="dados-pessoais" className="space-y-4">
                  <FormField
                    control={form.control}
                    name="nome"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Nome Completo *</FormLabel>
                        <FormControl>
                          <Input placeholder="Ex: João Silva" {...field} disabled={isLoading} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <FormField
                      control={form.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Email *</FormLabel>
                          <FormControl>
                            <Input type="email" placeholder="email@exemplo.com" {...field} disabled={isLoading} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="telefone"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Telefone *</FormLabel>
                          <FormControl>
                            <Input
                              placeholder="(00) 00000-0000"
                              {...field}
                              onChange={(e) => {
                                const masked = maskPhone(e.target.value);
                                field.onChange(unmaskPhone(masked));
                              }}
                              value={maskPhone(field.value)}
                              disabled={isLoading}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  {mode === 'create' && (
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                      <FormField
                        control={form.control}
                        name="cpf"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>CPF</FormLabel>
                            <FormControl>
                              <Input
                                placeholder="000.000.000-00"
                                {...field}
                                onChange={(e) => {
                                  const raw = unmaskCPF(e.target.value);
                                  field.onChange(raw);
                                }}
                                value={maskCPF(field.value || '')}
                                disabled={isLoading}
                                maxLength={14}
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <FormField
                        control={form.control}
                        name="data_admissao"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Data Admissão *</FormLabel>
                            <FormControl>
                              <Input
                                type="date"
                                {...field}
                                value={field.value instanceof Date ? format(field.value, 'yyyy-MM-dd') : ''}
                                onChange={(e) => field.onChange(new Date(e.target.value))}
                                disabled={isLoading}
                                max={format(new Date(), 'yyyy-MM-dd')}
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>
                  )}

                  <FormField
                    control={form.control}
                    name="observacoes"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Observações</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder="Informações adicionais..."
                            className="resize-none"
                            rows={3}
                            {...field}
                            disabled={isLoading}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </TabsContent>

                {/* TAB: PROFISSIONAL */}
                <TabsContent value="profissional" className="space-y-6">
                  {mode === 'create' && (
                    <FormField
                      control={form.control}
                      name="tipo"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Cargo *</FormLabel>
                          <Select onValueChange={field.onChange} defaultValue={field.value} disabled={isLoading}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Selecione..." />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="BARBEIRO">Barbeiro</SelectItem>
                              <SelectItem value="GERENTE">Gerente</SelectItem>
                              <SelectItem value="RECEPCIONISTA">Recepcionista</SelectItem>
                              <SelectItem value="OUTRO">Outro</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  )}

                  {tipoWatch === 'GERENTE' && mode === 'create' && (
                    <FormField
                      control={form.control}
                      name="tambem_barbeiro"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center space-x-3 space-y-0 rounded-md border p-4 bg-muted/20">
                          <FormControl>
                            <Checkbox checked={field.value} onCheckedChange={field.onChange} disabled={isLoading} />
                          </FormControl>
                          <div className="space-y-1 leading-none">
                            <FormLabel className="font-normal cursor-pointer">Atua também como <strong>Barbeiro</strong></FormLabel>
                            <FormDescription>Habilita configurações de comissão e agenda</FormDescription>
                          </div>
                        </FormItem>
                      )}
                    />
                  )}

                  {!showCommissionFields && (
                    <Alert>
                      <BriefcaseIcon className="h-4 w-4" />
                      <AlertTitle>Configuração Simplificada</AlertTitle>
                      <AlertDescription>
                        Profissionais administrativos não necessitam de configuração de comissão neste momento.
                      </AlertDescription>
                    </Alert>
                  )}
                </TabsContent>

                {/* TAB: FINANCEIRO (COMISSÕES) */}
                <TabsContent value="financeiro" className="space-y-6">
                  <div className="flex items-center gap-2 mb-4">
                    <Badge variant="outline" className="text-primary border-primary">
                      Configuração Financeira
                    </Badge>
                  </div>

                  {/* Comissões Padrão */}
                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 p-4 rounded-md border bg-muted/10">
                    <FormField
                      control={form.control}
                      name="tipo_comissao"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Tipo</FormLabel>
                          <Select onValueChange={field.onChange} defaultValue={field.value} disabled={isLoading}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="PERCENTUAL">Porcentagem (%)</SelectItem>
                              <SelectItem value="FIXO">Valor Fixo (R$)</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="comissao"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Serviços (Padrão)</FormLabel>
                          <div className="relative">
                            <FormControl>
                              <Input
                                type="number"
                                placeholder="0"
                                {...field}
                                value={field.value ?? ''}
                                onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : undefined)}
                                disabled={isLoading}
                                className="pr-8"
                              />
                            </FormControl>
                            <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-muted-foreground">
                              {form.watch('tipo_comissao') === 'PERCENTUAL' ? <PercentIcon className="size-4" /> : <span className="text-sm">R$</span>}
                            </div>
                          </div>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="comissao_produtos"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Produtos (Vendas)</FormLabel>
                          <div className="relative">
                            <FormControl>
                              <Input
                                type="number"
                                placeholder="0"
                                {...field}
                                value={field.value ?? ''}
                                onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : undefined)}
                                disabled={isLoading}
                                className="pr-8"
                              />
                            </FormControl>
                            <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-muted-foreground">
                              {form.watch('tipo_comissao') === 'PERCENTUAL' ? <PercentIcon className="size-4" /> : <span className="text-sm">R$</span>}
                            </div>
                          </div>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <Separator />

                  {/* Comissões por Categoria */}
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <h4 className="text-sm font-medium">Comissão por Categoria</h4>
                      <Badge variant="secondary" className="text-xs">
                        Personalizado
                      </Badge>
                    </div>

                    {/* Alert removido: Funcionalidade ativada */}

                    <div className="rounded-md border">
                      <div className="bg-muted/50 p-3 grid grid-cols-2 gap-4 text-xs font-medium uppercase text-muted-foreground">
                        <div>Categoria</div>
                        <div>Comissão Personalizada (%)</div>
                      </div>
                      <div className="divide-y max-h-[200px] overflow-y-auto">
                        {categories.map((category) => {
                          // Find field for this category or use empty
                          const fieldIndex = fields.findIndex(f => f.categoria_id === category.id);
                          const fieldName = `comissoes_por_categoria.${fieldIndex}.comissao` as const;

                          return (
                            <div key={category.id} className="p-3 grid grid-cols-2 gap-4 items-center hover:bg-muted/20 transition-colors">
                              <div className="flex items-center gap-2">
                                <div className="size-2 rounded-full" style={{ backgroundColor: category.cor || '#ccc' }} />
                                <span className="text-sm font-medium">{category.nome}</span>
                              </div>
                              <div>
                                {fieldIndex >= 0 ? (
                                  <FormField
                                    control={form.control}
                                    name={fieldName}
                                    render={({ field }) => (
                                      <FormItem className="space-y-0">
                                        <div className="relative">
                                          <FormControl>
                                            <Input
                                              type="number"
                                              className="h-8 pr-8"
                                              placeholder="Padrão"
                                              {...field}
                                              value={field.value ?? ''}
                                              onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : undefined)}
                                            />
                                          </FormControl>
                                          <div className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none text-muted-foreground">
                                            <PercentIcon className="size-3" />
                                          </div>
                                          <FormMessage />
                                        </div>
                                      </FormItem>
                                    )}
                                  />
                                ) : (
                                  <div className="text-xs text-muted-foreground italic">
                                    Não disponível
                                  </div>
                                )}
                              </div>
                            </div>
                          );
                        })}
                        {categories.length === 0 && (
                          <div className="p-4 text-center text-sm text-muted-foreground">
                            Nenhuma categoria de serviço cadastrada.
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </TabsContent>
              </Tabs>
            </ScrollArea>

            <DialogFooter className="px-6 py-4 border-t bg-background z-10">
              <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2Icon className="mr-2 size-4 animate-spin" />}
                {mode === 'create' ? 'Cadastrar Profissional' : 'Salvar Alterações'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default ProfessionalModal;
