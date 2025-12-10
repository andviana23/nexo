'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2Icon } from 'lucide-react';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';

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
import { Separator } from '@/components/ui/separator';
import {
    maskCEP,
    maskCNPJ,
    maskPhone,
    unmaskCEP,
    unmaskCNPJ,
    unmaskPhone
} from '@/lib/validations/professional';
import { supplierSchema, type SupplierFormValues } from '@/lib/validations/supplier';
import type { Fornecedor } from '@/types/fornecedor';

interface SupplierModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (data: SupplierFormValues) => Promise<void>;
    editingSupplier?: Fornecedor | null;
    isLoading?: boolean;
}

export function SupplierModal({
    isOpen,
    onClose,
    onSave,
    editingSupplier,
    isLoading,
}: SupplierModalProps) {
    const form = useForm<SupplierFormValues>({
        resolver: zodResolver(supplierSchema),
        defaultValues: {
            razao_social: '',
            nome_fantasia: '',
            cnpj: '',
            email: '',
            telefone: '',
            endereco_logradouro: '',
            endereco_cidade: '',
            endereco_estado: '',
            endereco_cep: '',
        },
    });

    useEffect(() => {
        if (editingSupplier) {
            form.reset({
                razao_social: editingSupplier.razao_social,
                nome_fantasia: editingSupplier.nome_fantasia || '',
                cnpj: editingSupplier.cnpj || '',
                email: editingSupplier.email || '',
                telefone: editingSupplier.telefone,
                endereco_logradouro: editingSupplier.endereco_logradouro || '',
                endereco_cidade: editingSupplier.endereco_cidade || '',
                endereco_estado: editingSupplier.endereco_estado || '',
                endereco_cep: editingSupplier.endereco_cep || '',
            });
        } else {
            form.reset({
                razao_social: '',
                nome_fantasia: '',
                cnpj: '',
                email: '',
                telefone: '',
                endereco_logradouro: '',
                endereco_cidade: '',
                endereco_estado: '',
                endereco_cep: '',
            });
        }
    }, [editingSupplier, form, isOpen]); // isOpen to reset if closed? better explicit reset on open in parent or here.

    const onSubmit = async (values: SupplierFormValues) => {
        // Clean masks
        const cleaned = {
            ...values,
            cnpj: values.cnpj ? unmaskCNPJ(values.cnpj) : undefined,
            telefone: unmaskPhone(values.telefone),
            endereco_cep: values.endereco_cep ? unmaskCEP(values.endereco_cep) : undefined,
            email: values.email || undefined,
            nome_fantasia: values.nome_fantasia || undefined,
            endereco_logradouro: values.endereco_logradouro || undefined,
            endereco_cidade: values.endereco_cidade || undefined,
            endereco_estado: values.endereco_estado || undefined,
        }
        await onSave(cleaned as any); // Types might misalign on optional vs null/undefined, safe cast for now or strict map.
    };

    const title = editingSupplier ? 'Editar Fornecedor' : 'Novo Fornecedor';
    const description = editingSupplier ? 'Atualize as informações do fornecedor.' : 'Preencha os dados para cadastrar um novo fornecedor.';

    return (
        <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
            <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    <DialogDescription>{description}</DialogDescription>
                </DialogHeader>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">

                        <div className="space-y-4">
                            <h4 className="text-sm font-medium text-muted-foreground uppercase tracking-wider text-[10px]">Dados da Empresa</h4>
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                <FormField
                                    control={form.control}
                                    name="razao_social"
                                    render={({ field }) => (
                                        <FormItem className="sm:col-span-2">
                                            <FormLabel>Razão Social *</FormLabel>
                                            <FormControl>
                                                <Input placeholder="Razão Social Ltda" {...field} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="nome_fantasia"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Nome Fantasia</FormLabel>
                                            <FormControl>
                                                <Input placeholder="Nome Comercial" {...field} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="cnpj"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>CNPJ</FormLabel>
                                            <FormControl>
                                                <Input
                                                    placeholder="00.000.000/0000-00"
                                                    maxLength={18}
                                                    {...field}
                                                    onChange={(e) => {
                                                        field.onChange(maskCNPJ(e.target.value));
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
                                    name="telefone"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Telefone *</FormLabel>
                                            <FormControl>
                                                <Input
                                                    placeholder="(00) 0000-0000"
                                                    {...field}
                                                    disabled={isLoading}
                                                    value={maskPhone(field.value)}
                                                    onChange={(e) => field.onChange(maskPhone(e.target.value))}
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
                                                <Input placeholder="contato@empresa.com" type="email" {...field} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </div>
                        </div>

                        <Separator />

                        <div className="space-y-4">
                            <h4 className="text-sm font-medium text-muted-foreground uppercase tracking-wider text-[10px]">Endereço</h4>
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                <FormField
                                    control={form.control}
                                    name="endereco_cep"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>CEP</FormLabel>
                                            <FormControl>
                                                <Input
                                                    placeholder="00000-000"
                                                    {...field}
                                                    disabled={isLoading}
                                                    value={maskCEP(field.value || '')}
                                                    onChange={(e) => field.onChange(maskCEP(e.target.value))}
                                                    maxLength={9}
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                                <div className="hidden sm:block"></div> {/* Spacer */}

                                <FormField
                                    control={form.control}
                                    name="endereco_logradouro"
                                    render={({ field }) => (
                                        <FormItem className="sm:col-span-2">
                                            <FormLabel>Logradouro</FormLabel>
                                            <FormControl>
                                                <Input placeholder="Rua, Av, Bairro..." {...field} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <FormField
                                    control={form.control}
                                    name="endereco_cidade"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Cidade</FormLabel>
                                            <FormControl>
                                                <Input placeholder="Cidade" {...field} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="endereco_estado"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>UF</FormLabel>
                                            <FormControl>
                                                <Input placeholder="UF" maxLength={2} {...field} onChange={e => field.onChange(e.target.value.toUpperCase())} disabled={isLoading} />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </div>
                        </div>

                        <DialogFooter>
                            <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
                                Cancelar
                            </Button>
                            <Button type="submit" disabled={isLoading}>
                                {isLoading && <Loader2Icon className="mr-2 size-4 animate-spin" />}
                                {editingSupplier ? 'Salvar Alterações' : 'Cadastrar'}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
}
