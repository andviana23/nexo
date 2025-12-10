'use client';

import {
    EditIcon,
    MoreHorizontalIcon,
    PowerIcon,
    PowerOffIcon,
    TrashIcon,
    TruckIcon
} from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { cn } from '@/lib/utils';
import { getFornecedorNome, type Fornecedor } from '@/types/fornecedor';

interface SuppliersTableProps {
    suppliers: Fornecedor[];
    onEdit: (supplier: Fornecedor) => void;
    onToggleActive: (supplier: Fornecedor) => void;
    onDelete: (id: string) => void;
    isLoading?: boolean;
}

export function SuppliersTable({
    suppliers,
    onEdit,
    onToggleActive,
    onDelete,
    isLoading
}: SuppliersTableProps) {

    if (!isLoading && suppliers.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center py-16 text-center bg-muted/5 border rounded-md">
                <div className="bg-muted p-4 rounded-full mb-4">
                    <TruckIcon className="size-8 text-muted-foreground" />
                </div>
                <h3 className="text-lg font-semibold">Nenhum fornecedor encontrado</h3>
                <p className="mt-1 text-sm text-muted-foreground max-w-sm">
                    Cadastre fornecedores para gerenciar suas compras.
                </p>
            </div>
        );
    }

    return (
        <div className="rounded-md border shadow-sm overflow-hidden">
            <Table>
                <TableHeader>
                    <TableRow className="bg-muted/50 hover:bg-muted/50">
                        <TableHead className="w-[300px] pl-6">Fornecedor</TableHead>
                        <TableHead>CNPJ</TableHead>
                        <TableHead>Telefone</TableHead>
                        <TableHead>Localização</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead className="text-right pr-6">Ações</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {suppliers.map((supplier) => (
                        <TableRow key={supplier.id} className="cursor-pointer group hover:bg-muted/30">
                            <TableCell className="pl-6 font-medium">
                                <div className="flex flex-col">
                                    <span className={cn(!supplier.ativo && "text-muted-foreground")}>
                                        {getFornecedorNome(supplier)}
                                    </span>
                                    {supplier.nome_fantasia && supplier.nome_fantasia !== supplier.razao_social && (
                                        <span className="text-xs text-muted-foreground">{supplier.razao_social}</span>
                                    )}
                                </div>
                            </TableCell>
                            <TableCell className="text-sm">
                                {supplier.cnpj ? (
                                    supplier.cnpj.replace(/^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$/, '$1.$2.$3/$4-$5')
                                ) : (
                                    <span className="text-muted-foreground">-</span>
                                )}
                            </TableCell>
                            <TableCell className="text-sm">
                                {supplier.telefone ? (
                                    supplier.telefone.replace(/^(\d{2})(\d{4,5})(\d{4})$/, '($1) $2-$3')
                                ) : (
                                    <span className="text-muted-foreground">-</span>
                                )}
                            </TableCell>
                            <TableCell className="text-sm">
                                {supplier.endereco_cidade
                                    ? `${supplier.endereco_cidade}${supplier.endereco_estado ? `/${supplier.endereco_estado}` : ''}`
                                    : <span className="text-muted-foreground">-</span>
                                }
                            </TableCell>
                            <TableCell>
                                <Badge variant={supplier.ativo ? 'default' : 'secondary'} className={cn(!supplier.ativo && "bg-muted text-muted-foreground font-normal")}>
                                    {supplier.ativo ? 'Ativo' : 'Inativo'}
                                </Badge>
                            </TableCell>
                            <TableCell className="text-right pr-6">
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" size="icon" className="size-8 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity">
                                            <MoreHorizontalIcon className="size-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end">
                                        <DropdownMenuLabel>Ações</DropdownMenuLabel>
                                        <DropdownMenuItem onClick={() => onEdit(supplier)}>
                                            <EditIcon className="mr-2 size-4" />
                                            Editar
                                        </DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => onToggleActive(supplier)}>
                                            {supplier.ativo ? (
                                                <>
                                                    <PowerOffIcon className="mr-2 size-4 text-orange-500" />
                                                    Desativar
                                                </>
                                            ) : (
                                                <>
                                                    <PowerIcon className="mr-2 size-4 text-green-500" />
                                                    Ativar
                                                </>
                                            )}
                                        </DropdownMenuItem>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem onClick={() => onDelete(supplier.id)} className="text-destructive focus:text-destructive">
                                            <TrashIcon className="mr-2 size-4" />
                                            Excluir
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
}
