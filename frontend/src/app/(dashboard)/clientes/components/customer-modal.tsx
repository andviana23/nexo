'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Modal de Cliente
 *
 * @component CustomerModal
 * @description Modal para criar, editar e visualizar clientes
 * Conforme FLUXO_CADASTROS_CLIENTE.md
 */

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { useCheckCpfExists, useCheckPhoneExists } from '@/hooks/use-customers';
import { cleanCPF, cleanPhone, formatCPF, formatPhone, isValidCPF, isValidPhone } from '@/services/customer-service';
import type {
    CreateCustomerRequest,
    CustomerGender,
    CustomerModalState,
    CustomerResponse,
    UpdateCustomerRequest,
} from '@/types/customer';
import { DEFAULT_TAGS, ESTADOS_BR, GENDER_LABELS } from '@/types/customer';
import { AlertCircleIcon, CheckCircleIcon, LoaderIcon, XIcon } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';

// =============================================================================
// PROPS
// =============================================================================

interface CustomerModalProps {
  state: CustomerModalState;
  onClose: () => void;
  onSave: (data: CreateCustomerRequest | UpdateCustomerRequest) => Promise<void>;
  isLoading?: boolean;
}

// =============================================================================
// HELPERS
// =============================================================================

function getInitialFormData(customer?: CustomerResponse): CreateCustomerRequest {
  if (!customer) {
    return {
      nome: '',
      telefone: '',
      email: '',
      cpf: '',
      data_nascimento: '',
      genero: undefined,
      endereco_logradouro: '',
      endereco_numero: '',
      endereco_complemento: '',
      endereco_bairro: '',
      endereco_cidade: '',
      endereco_estado: '',
      endereco_cep: '',
      observacoes: '',
      tags: [],
    };
  }
  return {
    nome: customer.nome,
    telefone: customer.telefone,
    email: customer.email || '',
    cpf: customer.cpf || '',
    data_nascimento: customer.data_nascimento || '',
    genero: customer.genero,
    endereco_logradouro: customer.endereco_logradouro || '',
    endereco_numero: customer.endereco_numero || '',
    endereco_complemento: customer.endereco_complemento || '',
    endereco_bairro: customer.endereco_bairro || '',
    endereco_cidade: customer.endereco_cidade || '',
    endereco_estado: customer.endereco_estado || '',
    endereco_cep: customer.endereco_cep || '',
    observacoes: customer.observacoes || '',
    tags: customer.tags || [],
  };
}

// =============================================================================
// COMPONENT
// =============================================================================

export function CustomerModal({
  state,
  onClose,
  onSave,
  isLoading = false,
}: CustomerModalProps) {
  const { isOpen, mode, customer } = state;

  // Derive initial form data from props
  const initialData = useMemo(
    () => getInitialFormData(mode !== 'create' ? customer : undefined),
    [mode, customer]
  );

  // Form state - reinitializes when modal opens with different customer
  const [formData, setFormData] = useState<CreateCustomerRequest>(initialData);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [phoneFormatted, setPhoneFormatted] = useState(() =>
    initialData.telefone ? formatPhone(initialData.telefone) : ''
  );
  const [cpfFormatted, setCpfFormatted] = useState(() =>
    initialData.cpf ? formatCPF(initialData.cpf) : ''
  );

  // Duplicate checks
  const { data: phoneExists, isLoading: checkingPhone } = useCheckPhoneExists(
    formData.telefone,
    mode === 'edit' ? customer?.id : undefined
  );
  const { data: cpfExists, isLoading: checkingCpf } = useCheckCpfExists(
    formData.cpf || '',
    mode === 'edit' ? customer?.id : undefined
  );

  // Handlers
  const handleChange = useCallback(
    (field: keyof CreateCustomerRequest, value: string | string[] | CustomerGender | undefined) => {
      setFormData((prev) => ({ ...prev, [field]: value }));
      // Clear error when field changes
      if (errors[field]) {
        setErrors((prev) => {
          const next = { ...prev };
          delete next[field];
          return next;
        });
      }
    },
    [errors]
  );

  const handlePhoneChange = useCallback((value: string) => {
    const cleaned = cleanPhone(value);
    setPhoneFormatted(formatPhone(cleaned));
    handleChange('telefone', cleaned);
  }, [handleChange]);

  const handleCpfChange = useCallback((value: string) => {
    const cleaned = cleanCPF(value);
    setCpfFormatted(formatCPF(cleaned));
    handleChange('cpf', cleaned);
  }, [handleChange]);

  const handleTagToggle = useCallback((tag: string) => {
    setFormData((prev) => ({
      ...prev,
      tags: prev.tags?.includes(tag)
        ? prev.tags.filter((t) => t !== tag)
        : [...(prev.tags || []), tag],
    }));
  }, []);

  const validate = useCallback(() => {
    const newErrors: Record<string, string> = {};

    // Required fields
    if (!formData.nome.trim()) {
      newErrors.nome = 'Nome é obrigatório';
    }
    if (!formData.telefone) {
      newErrors.telefone = 'Telefone é obrigatório';
    } else if (!isValidPhone(formData.telefone)) {
      newErrors.telefone = 'Telefone inválido';
    } else if (phoneExists) {
      newErrors.telefone = 'Telefone já cadastrado';
    }

    // Optional validations
    if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Email inválido';
    }
    if (formData.cpf) {
      if (!isValidCPF(formData.cpf)) {
        newErrors.cpf = 'CPF inválido';
      } else if (cpfExists) {
        newErrors.cpf = 'CPF já cadastrado';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [formData, phoneExists, cpfExists]);

  const handleSubmit = useCallback(async () => {
    if (!validate()) return;

    // Clean empty optional fields
    const cleanedData = {
      ...formData,
      email: formData.email || undefined,
      cpf: formData.cpf || undefined,
      data_nascimento: formData.data_nascimento || undefined,
      genero: formData.genero || undefined,
      endereco_logradouro: formData.endereco_logradouro || undefined,
      endereco_numero: formData.endereco_numero || undefined,
      endereco_complemento: formData.endereco_complemento || undefined,
      endereco_bairro: formData.endereco_bairro || undefined,
      endereco_cidade: formData.endereco_cidade || undefined,
      endereco_estado: formData.endereco_estado || undefined,
      endereco_cep: formData.endereco_cep || undefined,
      observacoes: formData.observacoes || undefined,
      tags: formData.tags || [], // Backend espera sempre um array, pode ser vazio
    };

    await onSave(cleanedData);
  }, [formData, validate, onSave]);

  // Modal title
  const title = {
    create: 'Novo Cliente',
    edit: 'Editar Cliente',
    view: 'Detalhes do Cliente',
  }[mode];

  const isViewMode = mode === 'view';

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            {mode === 'create' && 'Preencha os dados para cadastrar um novo cliente'}
            {mode === 'edit' && 'Atualize as informações do cliente'}
            {mode === 'view' && 'Visualize as informações do cliente'}
          </DialogDescription>
        </DialogHeader>

        <Tabs defaultValue="dados" className="w-full">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="dados">Dados Básicos</TabsTrigger>
            <TabsTrigger value="endereco">Endereço</TabsTrigger>
            <TabsTrigger value="observacoes">Observações</TabsTrigger>
          </TabsList>

          {/* TAB: Dados Básicos */}
          <TabsContent value="dados" className="space-y-4 mt-4">
            {/* Nome */}
            <div className="space-y-2">
              <Label htmlFor="nome">
                Nome <span className="text-destructive">*</span>
              </Label>
              <Input
                id="nome"
                value={formData.nome}
                onChange={(e) => handleChange('nome', e.target.value)}
                placeholder="Nome completo do cliente"
                disabled={isViewMode}
                className={errors.nome ? 'border-destructive' : ''}
              />
              {errors.nome && (
                <p className="text-xs text-destructive">{errors.nome}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              {/* Telefone */}
              <div className="space-y-2">
                <Label htmlFor="telefone">
                  Telefone <span className="text-destructive">*</span>
                </Label>
                <div className="relative">
                  <Input
                    id="telefone"
                    value={phoneFormatted}
                    onChange={(e) => handlePhoneChange(e.target.value)}
                    placeholder="(11) 99999-9999"
                    disabled={isViewMode}
                    className={errors.telefone ? 'border-destructive pr-10' : 'pr-10'}
                  />
                  <div className="absolute right-3 top-1/2 -translate-y-1/2">
                    {checkingPhone ? (
                      <LoaderIcon className="size-4 animate-spin text-muted-foreground" />
                    ) : phoneExists ? (
                      <AlertCircleIcon className="size-4 text-destructive" />
                    ) : formData.telefone && isValidPhone(formData.telefone) ? (
                      <CheckCircleIcon className="size-4 text-green-500" />
                    ) : null}
                  </div>
                </div>
                {errors.telefone && (
                  <p className="text-xs text-destructive">{errors.telefone}</p>
                )}
              </div>

              {/* CPF */}
              <div className="space-y-2">
                <Label htmlFor="cpf">CPF</Label>
                <div className="relative">
                  <Input
                    id="cpf"
                    value={cpfFormatted}
                    onChange={(e) => handleCpfChange(e.target.value)}
                    placeholder="000.000.000-00"
                    disabled={isViewMode}
                    className={errors.cpf ? 'border-destructive pr-10' : 'pr-10'}
                  />
                  <div className="absolute right-3 top-1/2 -translate-y-1/2">
                    {checkingCpf ? (
                      <LoaderIcon className="size-4 animate-spin text-muted-foreground" />
                    ) : cpfExists ? (
                      <AlertCircleIcon className="size-4 text-destructive" />
                    ) : formData.cpf && isValidCPF(formData.cpf) ? (
                      <CheckCircleIcon className="size-4 text-green-500" />
                    ) : null}
                  </div>
                </div>
                {errors.cpf && (
                  <p className="text-xs text-destructive">{errors.cpf}</p>
                )}
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              {/* Email */}
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => handleChange('email', e.target.value)}
                  placeholder="email@exemplo.com"
                  disabled={isViewMode}
                  className={errors.email ? 'border-destructive' : ''}
                />
                {errors.email && (
                  <p className="text-xs text-destructive">{errors.email}</p>
                )}
              </div>

              {/* Data de Nascimento */}
              <div className="space-y-2">
                <Label htmlFor="data_nascimento">Data de Nascimento</Label>
                <Input
                  id="data_nascimento"
                  type="date"
                  value={formData.data_nascimento}
                  onChange={(e) => handleChange('data_nascimento', e.target.value)}
                  disabled={isViewMode}
                />
              </div>
            </div>

            {/* Gênero */}
            <div className="space-y-2">
              <Label htmlFor="genero">Gênero</Label>
              <Select
                value={formData.genero || ''}
                onValueChange={(value) =>
                  handleChange('genero', value as CustomerGender || undefined)
                }
                disabled={isViewMode}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o gênero" />
                </SelectTrigger>
                <SelectContent>
                  {Object.entries(GENDER_LABELS).map(([value, label]) => (
                    <SelectItem key={value} value={value}>
                      {label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Tags */}
            <div className="space-y-2">
              <Label>Tags</Label>
              <div className="flex flex-wrap gap-2">
                {DEFAULT_TAGS.map((tag) => (
                  <Badge
                    key={tag}
                    variant={formData.tags?.includes(tag) ? 'default' : 'outline'}
                    className="cursor-pointer"
                    onClick={() => !isViewMode && handleTagToggle(tag)}
                  >
                    {tag}
                    {formData.tags?.includes(tag) && !isViewMode && (
                      <XIcon className="ml-1 size-3" />
                    )}
                  </Badge>
                ))}
              </div>
            </div>
          </TabsContent>

          {/* TAB: Endereço */}
          <TabsContent value="endereco" className="space-y-4 mt-4">
            <div className="grid grid-cols-3 gap-4">
              {/* CEP */}
              <div className="space-y-2">
                <Label htmlFor="endereco_cep">CEP</Label>
                <Input
                  id="endereco_cep"
                  value={formData.endereco_cep}
                  onChange={(e) => handleChange('endereco_cep', e.target.value)}
                  placeholder="00000-000"
                  disabled={isViewMode}
                />
              </div>

              {/* Logradouro */}
              <div className="space-y-2 col-span-2">
                <Label htmlFor="endereco_logradouro">Logradouro</Label>
                <Input
                  id="endereco_logradouro"
                  value={formData.endereco_logradouro}
                  onChange={(e) => handleChange('endereco_logradouro', e.target.value)}
                  placeholder="Rua, Avenida, etc."
                  disabled={isViewMode}
                />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4">
              {/* Número */}
              <div className="space-y-2">
                <Label htmlFor="endereco_numero">Número</Label>
                <Input
                  id="endereco_numero"
                  value={formData.endereco_numero}
                  onChange={(e) => handleChange('endereco_numero', e.target.value)}
                  placeholder="123"
                  disabled={isViewMode}
                />
              </div>

              {/* Complemento */}
              <div className="space-y-2 col-span-2">
                <Label htmlFor="endereco_complemento">Complemento</Label>
                <Input
                  id="endereco_complemento"
                  value={formData.endereco_complemento}
                  onChange={(e) => handleChange('endereco_complemento', e.target.value)}
                  placeholder="Apto, Bloco, etc."
                  disabled={isViewMode}
                />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4">
              {/* Bairro */}
              <div className="space-y-2">
                <Label htmlFor="endereco_bairro">Bairro</Label>
                <Input
                  id="endereco_bairro"
                  value={formData.endereco_bairro}
                  onChange={(e) => handleChange('endereco_bairro', e.target.value)}
                  placeholder="Bairro"
                  disabled={isViewMode}
                />
              </div>

              {/* Cidade */}
              <div className="space-y-2">
                <Label htmlFor="endereco_cidade">Cidade</Label>
                <Input
                  id="endereco_cidade"
                  value={formData.endereco_cidade}
                  onChange={(e) => handleChange('endereco_cidade', e.target.value)}
                  placeholder="Cidade"
                  disabled={isViewMode}
                />
              </div>

              {/* Estado */}
              <div className="space-y-2">
                <Label htmlFor="endereco_estado">Estado</Label>
                <Select
                  value={formData.endereco_estado || ''}
                  onValueChange={(value) => handleChange('endereco_estado', value)}
                  disabled={isViewMode}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="UF" />
                  </SelectTrigger>
                  <SelectContent>
                    {ESTADOS_BR.map((estado) => (
                      <SelectItem key={estado.value} value={estado.value}>
                        {estado.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </TabsContent>

          {/* TAB: Observações */}
          <TabsContent value="observacoes" className="space-y-4 mt-4">
            <div className="space-y-2">
              <Label htmlFor="observacoes">Observações</Label>
              <Textarea
                id="observacoes"
                value={formData.observacoes}
                onChange={(e) => handleChange('observacoes', e.target.value)}
                placeholder="Preferências, alergias, informações importantes..."
                rows={6}
                disabled={isViewMode}
              />
              <p className="text-xs text-muted-foreground">
                Adicione informações relevantes sobre o cliente que possam ajudar no atendimento.
              </p>
            </div>

            {/* View mode: show additional info */}
            {isViewMode && customer && (
              <>
                <Separator />
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-muted-foreground">Cadastrado em</p>
                    <p className="font-medium">
                      {new Date(customer.criado_em).toLocaleString('pt-BR')}
                    </p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Última atualização</p>
                    <p className="font-medium">
                      {new Date(customer.atualizado_em).toLocaleString('pt-BR')}
                    </p>
                  </div>
                </div>
              </>
            )}
          </TabsContent>
        </Tabs>

        <DialogFooter className="mt-6">
          <Button variant="outline" onClick={onClose}>
            {isViewMode ? 'Fechar' : 'Cancelar'}
          </Button>
          {!isViewMode && (
            <Button onClick={handleSubmit} disabled={isLoading}>
              {isLoading && <LoaderIcon className="mr-2 size-4 animate-spin" />}
              {mode === 'create' ? 'Cadastrar' : 'Salvar'}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
