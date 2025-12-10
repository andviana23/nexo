'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Customer Modal
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { format } from 'date-fns';
import {
    CalendarIcon,
    Loader2Icon,
    MailIcon,
    MapPinIcon,
    PhoneIcon,
    UserIcon,
    XIcon
} from 'lucide-react';
import { useCallback, useEffect, useMemo } from 'react';
import { useForm } from 'react-hook-form';

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
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { customerSchema, type CustomerFormValues } from '@/lib/validations/customer';
import {
    maskCPF,
    maskPhone,
    unmaskCPF,
    unmaskPhone
} from '@/lib/validations/professional';
import { ESTADOS_BR, type CustomerModalState } from '@/types/customer';

// Utility for Tag Colors (reuse or define)
const DEFAULT_TAGS = ['Vip', 'Frequente', 'Novo', 'Atrasado', 'Problema'];

interface CustomerModalProps {
    state: CustomerModalState;
    onClose: () => void;
    onSave: (data: CustomerFormValues) => Promise<void>;
    isLoading?: boolean;
}

export function CustomerModal({
    state,
    onClose,
    onSave,
    isLoading = false,
}: CustomerModalProps) {
    const { isOpen, mode, customer } = state;
    const isViewMode = mode === 'view';

    const form = useForm<CustomerFormValues>({
        resolver: zodResolver(customerSchema),
        defaultValues: {
            nome: '',
            telefone: '',
            email: '',
            cpf: '',
            data_nascimento: '',
            genero: undefined,
            endereco_cep: '',
            endereco_logradouro: '',
            endereco_numero: '',
            endereco_complemento: '',
            endereco_bairro: '',
            endereco_cidade: '',
            endereco_estado: '',
            observacoes: '',
            tags: [],
            ativo: true,
        },
    });

    useEffect(() => {
        if ((mode === 'edit' || mode === 'view') && customer) {
            form.reset({
                nome: customer.nome,
                telefone: customer.telefone,
                email: customer.email || '',
                cpf: customer.cpf || '',
                data_nascimento: customer.data_nascimento ? format(new Date(customer.data_nascimento), 'yyyy-MM-dd') : '',
                genero: customer.genero,
                endereco_cep: customer.endereco_cep || '',
                endereco_logradouro: customer.endereco_logradouro || '',
                endereco_numero: customer.endereco_numero || '',
                endereco_complemento: customer.endereco_complemento || '',
                endereco_bairro: customer.endereco_bairro || '',
                endereco_cidade: customer.endereco_cidade || '',
                endereco_estado: customer.endereco_estado || '',
                observacoes: customer.observacoes || '',
                tags: customer.tags || [],
                ativo: customer.ativo,
            });
        } else if (mode === 'create') {
            form.reset({
                nome: '',
                telefone: '',
                email: '',
                cpf: '',
                data_nascimento: '',
                genero: undefined,
                endereco_cep: '',
                endereco_logradouro: '',
                endereco_numero: '',
                endereco_complemento: '',
                endereco_bairro: '',
                endereco_cidade: '',
                endereco_estado: '',
                observacoes: '',
                tags: [],
                ativo: true,
            });
        }
    }, [mode, customer, form]);

    const onSubmit = useCallback(async (values: CustomerFormValues) => {
        // Unmask handled by caller? Or here?
        // Zod schema expects strings. I should clean phone/cpf before sending if backend expects clean.
        // Usually hooks handle it, but I'll do it here if needed.
        // For consistency, I will clean here.
        const cleaned = {
            ...values,
            telefone: unmaskPhone(values.telefone),
            cpf: values.cpf ? unmaskCPF(values.cpf) : undefined,
            // Ensure optional fields are undefined if empty string
            email: values.email || undefined,
            data_nascimento: values.data_nascimento || undefined,
            // ...
        }
        await onSave(cleaned);
    }, [onSave]);

    const title = useMemo(() => {
        switch (mode) {
            case 'create': return 'Novo Cliente';
            case 'edit': return 'Editar Cliente';
            case 'view': return 'Detalhes do Cliente';
            default: return 'Cliente';
        }
    }, [mode]);

    const handleTagToggle = (tag: string) => {
        const currentTags = form.getValues('tags') || [];
        if (currentTags.includes(tag)) {
            form.setValue('tags', currentTags.filter(t => t !== tag));
        } else {
            form.setValue('tags', [...currentTags, tag]);
        }
    };

    if (isViewMode && customer) {
        return (
            <Dialog open={isOpen} onOpenChange={() => onClose()}>
                <DialogContent className="sm:max-w-[500px]">
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <UserIcon className="size-5 text-primary" />
                            {title}
                        </DialogTitle>
                        <DialogDescription>Informações completas do cliente</DialogDescription>
                    </DialogHeader>
                    <div className="space-y-6 pt-2">
                        {/* Header Info */}
                        <div>
                            <h3 className="text-xl font-bold">{customer.nome}</h3>
                            <div className="flex items-center gap-2 mt-2">
                                <Badge variant={customer.ativo ? 'default' : 'secondary'}>{customer.ativo ? 'Ativo' : 'Inativo'}</Badge>
                                {customer.genero && <Badge variant="outline">{customer.genero}</Badge>}
                            </div>
                        </div>

                        <Separator />

                        {/* Contact */}
                        <div className="space-y-3">
                            <div className="flex items-center gap-3 text-sm">
                                <PhoneIcon className="size-4 text-muted-foreground" />
                                <span>{maskPhone(customer.telefone)}</span>
                            </div>
                            {customer.email && (
                                <div className="flex items-center gap-3 text-sm">
                                    <MailIcon className="size-4 text-muted-foreground" />
                                    <span>{customer.email}</span>
                                </div>
                            )}
                            {customer.data_nascimento && (
                                <div className="flex items-center gap-3 text-sm">
                                    <CalendarIcon className="size-4 text-muted-foreground" />
                                    <span>Nascimento: {format(new Date(customer.data_nascimento), 'dd/MM/yyyy')}</span>
                                </div>
                            )}
                        </div>

                        {/* Address */}
                        {(customer.endereco_logradouro || customer.endereco_cidade) && (
                            <>
                                <Separator />
                                <div className="space-y-2">
                                    <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                                        <MapPinIcon className="size-4" />
                                        Endereço
                                    </div>
                                    <p className="text-sm pl-6">
                                        {customer.endereco_logradouro}, {customer.endereco_numero} {customer.endereco_complemento}
                                        <br />
                                        {customer.endereco_bairro} - {customer.endereco_cidade}/{customer.endereco_estado}
                                        <br />
                                        CEP: {customer.endereco_cep}
                                    </p>
                                </div>
                            </>
                        )}

                        {/* Tags */}
                        {customer.tags && customer.tags.length > 0 && (
                            <>
                                <Separator />
                                <div className="flex flex-wrap gap-2">
                                    {customer.tags.map(tag => (
                                        <Badge key={tag} variant="secondary">{tag}</Badge>
                                    ))}
                                </div>
                            </>
                        )}

                        {/* Obs */}
                        {customer.observacoes && (
                            <div className="bg-muted/30 p-3 rounded-md text-sm italic">
                                "{customer.observacoes}"
                            </div>
                        )}

                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={onClose}>Fechar</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        );
    }

    return (
        <Dialog open={isOpen} onOpenChange={() => onClose()}>
            <DialogContent className="sm:max-w-[650px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    <DialogDescription>Preencha os dados do cliente.</DialogDescription>
                </DialogHeader>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <Tabs defaultValue="basic" className="w-full">
                            <TabsList className="grid w-full grid-cols-3">
                                <TabsTrigger value="basic">Dados Básicos</TabsTrigger>
                                <TabsTrigger value="address">Endereço</TabsTrigger>
                                <TabsTrigger value="notes">Observações</TabsTrigger>
                            </TabsList>

                            <TabsContent value="basic" className="space-y-4 mt-4">
                                <FormField
                                    control={form.control}
                                    name="nome"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Nome Completo *</FormLabel>
                                            <FormControl>
                                                <Input placeholder="Ex: Maria Oliveira" {...field} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
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
                                                            const masked = maskPhone(e.target.value); // Assume maskPhone handles typing
                                                            field.onChange(masked);
                                                        }}
                                                        disabled={isLoading}
                                                    />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="email"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Email</FormLabel>
                                                <FormControl>
                                                    <Input type="email" placeholder="email@exemplo.com" {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </div>

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
                                                        onChange={(e) => field.onChange(maskCPF(e.target.value))}
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
                                        name="data_nascimento"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Data de Nascimento</FormLabel>
                                                <FormControl>
                                                    <Input type="date" {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </div>

                                <FormField
                                    control={form.control}
                                    name="genero"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Gênero</FormLabel>
                                            <Select onValueChange={field.onChange} defaultValue={field.value} disabled={isLoading}>
                                                <FormControl>
                                                    <SelectTrigger>
                                                        <SelectValue placeholder="Selecione..." />
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    <SelectItem value="M">Masculino</SelectItem>
                                                    <SelectItem value="F">Feminino</SelectItem>
                                                    <SelectItem value="NB">Não Binário</SelectItem>
                                                    <SelectItem value="PNI">Prefiro não Informar</SelectItem>
                                                </SelectContent>
                                            </Select>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </TabsContent>

                            <TabsContent value="address" className="space-y-4 mt-4">
                                <div className="grid grid-cols-3 gap-4">
                                    <FormField
                                        control={form.control}
                                        name="endereco_cep"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>CEP</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="00000-000" {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="endereco_cidade"
                                        render={({ field }) => (
                                            <FormItem className="col-span-2">
                                                <FormLabel>Cidade</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="Cidade" {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </div>

                                <div className="grid grid-cols-4 gap-4">
                                    <FormField
                                        control={form.control}
                                        name="endereco_logradouro"
                                        render={({ field }) => (
                                            <FormItem className="col-span-3">
                                                <FormLabel>Logradouro</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="Rua, Av..." {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="endereco_numero"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Nº</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="123" {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </div>

                                <div className="grid grid-cols-2 gap-4">
                                    <FormField
                                        control={form.control}
                                        name="endereco_bairro"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Bairro</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="Bairro" {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="endereco_complemento"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Complemento</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="Apto, Bloco..." {...field} disabled={isLoading} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                </div>

                                <FormField
                                    control={form.control}
                                    name="endereco_estado"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Estado</FormLabel>
                                            <Select onValueChange={field.onChange} defaultValue={field.value} disabled={isLoading}>
                                                <FormControl>
                                                    <SelectTrigger>
                                                        <SelectValue placeholder="UF" />
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    {ESTADOS_BR.map(uf => (
                                                        <SelectItem key={uf.value} value={uf.value}>{uf.label}</SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </TabsContent>

                            <TabsContent value="notes" className="space-y-4 mt-4">
                                <FormField
                                    control={form.control}
                                    name="tags"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Tags</FormLabel>
                                            <div className="flex flex-wrap gap-2">
                                                {DEFAULT_TAGS.map(tag => {
                                                    const isSelected = field.value?.includes(tag);
                                                    return (
                                                        <Badge
                                                            key={tag}
                                                            variant={isSelected ? 'default' : 'outline'}
                                                            className="cursor-pointer select-none"
                                                            onClick={() => handleTagToggle(tag)}
                                                        >
                                                            {tag}
                                                            {isSelected && <XIcon className="ml-1 size-3" />}
                                                        </Badge>
                                                    );
                                                })}
                                            </div>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <FormField
                                    control={form.control}
                                    name="observacoes"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Observações Gerais</FormLabel>
                                            <FormControl>
                                                <Textarea
                                                    placeholder="Preferências, histórico, etc."
                                                    rows={5}
                                                    {...field}
                                                    disabled={isLoading}
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </TabsContent>
                        </Tabs>

                        <DialogFooter>
                            <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
                                Cancelar
                            </Button>
                            <Button type="submit" disabled={isLoading}>
                                {isLoading && <Loader2Icon className="mr-2 size-4 animate-spin" />}
                                {mode === 'create' ? 'Cadastrar' : 'Salvar Alterações'}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
}
