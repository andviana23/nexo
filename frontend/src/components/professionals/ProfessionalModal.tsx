'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Professional Modal
 *
 * @component ProfessionalModal
 * @description Modal para criar, editar e visualizar profissionais
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
  CalendarIcon,
  Loader2Icon,
  MailIcon,
  PhoneIcon,
  UserIcon,
} from 'lucide-react';
import { useCallback, useEffect, useMemo } from 'react';
import { useForm, useWatch } from 'react-hook-form';

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
import { Textarea } from '@/components/ui/textarea';
import {
  useCreateProfessional,
  useUpdateProfessional,
} from '@/hooks/use-professionals';
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
  /** Estado do modal */
  state: ProfessionalModalState;
  /** Callback para fechar o modal */
  onClose: () => void;
  /** Callback após salvar com sucesso */
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

  const isLoading = createProfessional.isPending || updateProfessional.isPending;
  const isViewMode = mode === 'view';

  // Form com validação Zod
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
    },
  });

  const tipoWatch = useWatch({ control: form.control, name: 'tipo' });
  const tambemBarbeiroWatch = useWatch({ control: form.control, name: 'tambem_barbeiro' });

  // Mostra campos de comissão se for Barbeiro ou Gerente+Barbeiro
  const showCommissionFields = useMemo(() => {
    return tipoWatch === 'BARBEIRO' || (tipoWatch === 'GERENTE' && tambemBarbeiroWatch);
  }, [tipoWatch, tambemBarbeiroWatch]);

  // Preencher formulário quando abrir para editar
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
        observacoes: '',
        tambem_barbeiro: false, // Será calculado a partir do tipo
        tipo_comissao: professional.tipo_comissao || 'PERCENTUAL',
        comissao: professional.comissao ? Number(professional.comissao) : undefined,
        comissao_produtos: professional.comissao_produtos ? Number(professional.comissao_produtos) : undefined,
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
      });
    }
  }, [mode, professional, form]);

  // Submit handler
  const onSubmit = useCallback(
    async (values: CreateProfessionalFormData) => {
      // Prepara dados removendo campos que não existem no backend
      const { tambem_barbeiro: _tambemBarbeiro, ...rest } = values;
      void _tambemBarbeiro;
      
      // Limpa máscaras, formata data e mantém comissao como number
      const data = {
        ...rest,
        cpf: unmaskCPF(rest.cpf),
        telefone: unmaskPhone(rest.telefone),
        data_admissao: format(rest.data_admissao, 'yyyy-MM-dd'),
        comissao: rest.comissao ?? 0,
      };

      // DEBUG: Log do payload
      console.log('[ProfessionalModal] Payload sendo enviado:', data);

      if (mode === 'create') {
        createProfessional.mutate(data as Parameters<typeof createProfessional.mutate>[0], {
          onSuccess: (result) => {
            onSuccess?.(result);
            onClose();
          },
          onError: (error) => {
            console.error('[ProfessionalModal] Erro ao criar:', error);
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

  // Título do modal
  const title = useMemo(() => {
    switch (mode) {
      case 'create':
        return 'Novo Profissional';
      case 'edit':
        return 'Editar Profissional';
      case 'view':
        return 'Detalhes do Profissional';
      default:
        return 'Profissional';
    }
  }, [mode]);

  // ==========================================================================
  // VIEW MODE RENDER
  // ==========================================================================

  if (isViewMode && professional) {
    return (
      <Dialog open={isOpen} onOpenChange={() => onClose()}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <UserIcon className="size-5 text-primary" />
              {title}
            </DialogTitle>
            <DialogDescription>
              Informações do profissional
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {/* Nome e Status */}
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-semibold">{professional.nome}</h3>
                <div className="flex items-center gap-2 mt-1">
                  <ProfessionalTypeBadge type={professional.tipo} />
                  <ProfessionalStatusBadge status={professional.status} />
                </div>
              </div>
            </div>

            <Separator />

            {/* Contato */}
            <div className="space-y-3">
              <div className="flex items-center gap-3">
                <MailIcon className="size-4 text-muted-foreground" />
                <span className="text-sm">{professional.email}</span>
              </div>
              <div className="flex items-center gap-3">
                <PhoneIcon className="size-4 text-muted-foreground" />
                <span className="text-sm">{maskPhone(professional.telefone)}</span>
              </div>
              <div className="flex items-center gap-3">
                <CalendarIcon className="size-4 text-muted-foreground" />
                <span className="text-sm">
                  Admitido em{' '}
                  {format(new Date(professional.data_admissao), "dd 'de' MMMM 'de' yyyy", {
                    locale: ptBR,
                  })}
                </span>
              </div>
            </div>

            {/* Comissão */}
            {professional.comissao && (
              <>
                <Separator />
                <div>
                  <h4 className="text-sm font-medium mb-2">Comissão</h4>
                  <div className="flex items-center gap-4">
                    <Badge variant="outline">
                      Serviços: {professional.comissao}%
                    </Badge>
                    {professional.comissao_produtos && (
                      <Badge variant="outline">
                        Produtos: {professional.comissao_produtos}%
                      </Badge>
                    )}
                  </div>
                </div>
              </>
            )}

            {/* Especialidades */}
            {professional.especialidades && professional.especialidades.length > 0 && (
              <>
                <Separator />
                <div>
                  <h4 className="text-sm font-medium mb-2">Especialidades</h4>
                  <div className="flex flex-wrap gap-2">
                    {professional.especialidades.map((esp) => (
                      <Badge key={esp} variant="secondary">
                        {esp}
                      </Badge>
                    ))}
                  </div>
                </div>
              </>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={onClose}>
              Fechar
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  // ==========================================================================
  // CREATE/EDIT MODE RENDER
  // ==========================================================================

  return (
    <Dialog open={isOpen} onOpenChange={() => onClose()}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <UserIcon className="size-5 text-primary" />
            {title}
          </DialogTitle>
          <DialogDescription>
            {mode === 'create'
              ? 'Preencha os dados para cadastrar um novo profissional'
              : 'Atualize os dados do profissional'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {/* Dados Pessoais */}
            <div className="space-y-4">
              <h3 className="text-sm font-medium text-muted-foreground">
                Dados Pessoais
              </h3>

              {/* Nome */}
              <FormField
                control={form.control}
                name="nome"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nome completo *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="João da Silva"
                        {...field}
                        disabled={isLoading}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Email e Telefone */}
              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email *</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder="joao@email.com"
                          {...field}
                          disabled={isLoading}
                        />
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
                          placeholder="(11) 99999-9999"
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

              {/* CPF e Data Admissão - só em criação */}
              {mode === 'create' && (
                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="cpf"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>CPF/CNPJ *</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="000.000.000-00 ou 00.000.000/0000-00"
                            {...field}
                            onChange={(e) => {
                              const masked = maskCPF(e.target.value);
                              field.onChange(unmaskCPF(masked));
                            }}
                            value={maskCPF(field.value)}
                            disabled={isLoading}
                            maxLength={18}
                          />
                        </FormControl>
                        <FormDescription className="text-xs">
                          Digite CPF (11 dígitos) ou CNPJ (14 dígitos)
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="data_admissao"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Data de Admissão *</FormLabel>
                        <FormControl>
                          <Input
                            type="date"
                            {...field}
                            value={
                              field.value instanceof Date
                                ? format(field.value, 'yyyy-MM-dd')
                                : ''
                            }
                            onChange={(e) =>
                              field.onChange(new Date(e.target.value))
                            }
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
            </div>

            <Separator />

            {/* Dados Profissionais */}
            <div className="space-y-4">
              <h3 className="text-sm font-medium text-muted-foreground">
                Dados Profissionais
              </h3>

              {/* Tipo - só em criação */}
              {mode === 'create' && (
                <FormField
                  control={form.control}
                  name="tipo"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tipo de Profissional *</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                        disabled={isLoading}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione o tipo" />
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

              {/* Checkbox Também Barbeiro (só para Gerente) */}
              {tipoWatch === 'GERENTE' && mode === 'create' && (
                <FormField
                  control={form.control}
                  name="tambem_barbeiro"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                      <FormControl>
                        <Checkbox
                          checked={field.value}
                          onCheckedChange={field.onChange}
                          disabled={isLoading}
                        />
                      </FormControl>
                      <div className="space-y-1 leading-none">
                        <FormLabel>Também atuo como Barbeiro</FormLabel>
                        <FormDescription>
                          Marque se este gerente também realiza atendimentos
                        </FormDescription>
                      </div>
                    </FormItem>
                  )}
                />
              )}

              {/* Campos de Comissão */}
              {showCommissionFields && (
                <div className="space-y-4 p-4 rounded-md border bg-muted/50">
                  <h4 className="text-sm font-medium">Comissão</h4>

                  <div className="grid grid-cols-3 gap-4">
                    <FormField
                      control={form.control}
                      name="tipo_comissao"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Tipo</FormLabel>
                          <Select
                            onValueChange={field.onChange}
                            defaultValue={field.value}
                            disabled={isLoading}
                          >
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="PERCENTUAL">Percentual</SelectItem>
                              <SelectItem value="FIXO">Valor Fixo</SelectItem>
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
                          <FormLabel>Serviços *</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              placeholder="40"
                              {...field}
                              value={field.value ?? ''}
                              onChange={(e) =>
                                field.onChange(
                                  e.target.value ? Number(e.target.value) : undefined
                                )
                              }
                              disabled={isLoading}
                              min={0}
                              max={100}
                              step={1}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="comissao_produtos"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Produtos</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              placeholder="10"
                              {...field}
                              value={field.value ?? ''}
                              onChange={(e) =>
                                field.onChange(
                                  e.target.value ? Number(e.target.value) : undefined
                                )
                              }
                              disabled={isLoading}
                              min={0}
                              max={100}
                              step={1}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </div>
              )}

              {/* Observações */}
              <FormField
                control={form.control}
                name="observacoes"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Observações</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Observações sobre o profissional..."
                        className="resize-none"
                        {...field}
                        disabled={isLoading}
                      />
                    </FormControl>
                    <FormDescription>Máximo de 500 caracteres</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <DialogFooter className="pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={onClose}
                disabled={isLoading}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2Icon className="mr-2 size-4 animate-spin" />}
                {mode === 'create' ? 'Cadastrar' : 'Salvar'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default ProfessionalModal;
