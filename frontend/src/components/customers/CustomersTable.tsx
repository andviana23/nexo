'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Customers Table
 */

import {
    CakeIcon,
    EditIcon,
    EyeIcon,
    HistoryIcon,
    MoreHorizontalIcon,
    PhoneIcon,
    TrashIcon,
    UserXIcon
} from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useCallback } from 'react';

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
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
import { getInitials } from '@/services/customer-service';
import { formatPhone, TAG_COLORS, type CustomerResponse } from '@/types/customer';

// =============================================================================
// TYPES
// =============================================================================

interface CustomersTableProps {
    customers: CustomerResponse[];
    onView: (customer: CustomerResponse) => void;
    onEdit: (customer: CustomerResponse) => void;
    onInactivate: (customer: CustomerResponse) => void;
    isLoading?: boolean;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function CustomersTable({
    customers,
    onView,
    onEdit,
    onInactivate,
    isLoading,
}: CustomersTableProps) {
    const router = useRouter();

    const handleRowClick = useCallback(
        (customer: CustomerResponse) => {
            onView(customer);
        },
        [onView]
    );

    const safeCustomers = customers.filter(
        (c): c is CustomerResponse => c != null && typeof c === 'object' && 'id' in c
    );

    if (safeCustomers.length === 0 && !isLoading) {
        return (
            <div className="flex flex-col items-center justify-center py-16 text-center bg-muted/5 border rounded-md">
                <div className="bg-muted p-4 rounded-full mb-4">
                    <UserXIcon className="size-8 text-muted-foreground" />
                </div>
                <h3 className="text-lg font-semibold">Nenhum cliente encontrado</h3>
                <p className="mt-1 text-sm text-muted-foreground max-w-sm">
                    Não encontramos clientes com os filtros aplicados.
                </p>
            </div>
        );
    }

    return (
        <div className="rounded-md border shadow-sm overflow-hidden">
            <Table>
                <TableHeader>
                    <TableRow className="bg-muted/50 hover:bg-muted/50">
                        <TableHead className="w-[300px] pl-6">Cliente</TableHead>
                        <TableHead>Contato</TableHead>
                        <TableHead>Tags</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead className="text-right pr-6">Ações</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {safeCustomers.map((customer) => (
                        <TableRow
                            key={customer.id}
                            className="cursor-pointer group hover:bg-muted/30"
                            onClick={() => handleRowClick(customer)}
                        >
                            <TableCell className="pl-6">
                                <div className="flex items-center gap-3">
                                    <Avatar className={cn("size-9 border", !customer.ativo && "opacity-60")}>
                                        <AvatarFallback className={cn(
                                            "text-xs font-medium text-white",
                                            customer.ativo ? "bg-primary" : "bg-muted-foreground"
                                        )}>
                                            {getInitials(customer.nome)}
                                        </AvatarFallback>
                                    </Avatar>
                                    <div>
                                        <p className={cn("font-medium text-sm", !customer.ativo && "text-muted-foreground")}>{customer.nome}</p>
                                        {customer.data_nascimento && (
                                            <div className="flex items-center gap-1 text-xs text-muted-foreground">
                                                <CakeIcon className="size-3" />
                                                {new Date(customer.data_nascimento).toLocaleDateString('pt-BR')}
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </TableCell>

                            <TableCell>
                                <div className="flex flex-col gap-0.5">
                                    <div className="flex items-center gap-1.5 text-sm">
                                        <PhoneIcon className="size-3 text-muted-foreground" />
                                        <span>{formatPhone(customer.telefone)}</span>
                                    </div>
                                    {customer.email && (
                                        <span className="text-xs text-muted-foreground truncate max-w-[150px]">{customer.email}</span>
                                    )}
                                </div>
                            </TableCell>

                            <TableCell>
                                <div className="flex flex-wrap gap-1">
                                    {customer.tags && customer.tags.length > 0 ? (
                                        customer.tags.slice(0, 2).map((tag) => (
                                            <Badge
                                                key={tag}
                                                variant="secondary"
                                                className={cn('text-[10px] px-1.5 h-5 font-normal', TAG_COLORS[tag])}
                                            >
                                                {tag}
                                            </Badge>
                                        ))
                                    ) : (
                                        <span className="text-muted-foreground/30 text-xs">—</span>
                                    )}
                                    {customer.tags && customer.tags.length > 2 && (
                                        <Badge variant="outline" className="text-[10px] px-1.5 h-5">+{customer.tags.length - 2}</Badge>
                                    )}
                                </div>
                            </TableCell>

                            <TableCell>
                                <Badge variant={customer.ativo ? 'default' : 'secondary'} className={cn(
                                    "text-[10px] font-normal",
                                    !customer.ativo && "bg-muted text-muted-foreground"
                                )}>
                                    {customer.ativo ? 'Ativo' : 'Inativo'}
                                </Badge>
                            </TableCell>

                            <TableCell className="text-right pr-6" onClick={(e) => e.stopPropagation()}>
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" size="icon" className="size-8 text-muted-foreground hover:text-foreground opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity">
                                            <span className="sr-only">Abrir menu</span>
                                            <MoreHorizontalIcon className="size-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end">
                                        <DropdownMenuLabel>Ações</DropdownMenuLabel>
                                        <DropdownMenuItem asChild>
                                            <Link href={`/clientes/${customer.id}`}>
                                                <HistoryIcon className="mr-2 size-4" />
                                                Histórico
                                            </Link>
                                        </DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => onView(customer)}>
                                            <EyeIcon className="mr-2 size-4" />
                                            Detalhes
                                        </DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => onEdit(customer)}>
                                            <EditIcon className="mr-2 size-4" />
                                            Editar
                                        </DropdownMenuItem>

                                        {customer.ativo && (
                                            <>
                                                <DropdownMenuSeparator />
                                                <DropdownMenuItem
                                                    onClick={() => onInactivate(customer)}
                                                    className="text-destructive focus:text-destructive"
                                                >
                                                    <TrashIcon className="mr-2 size-4" />
                                                    Inativar
                                                </DropdownMenuItem>
                                            </>
                                        )}
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

export default CustomersTable;
