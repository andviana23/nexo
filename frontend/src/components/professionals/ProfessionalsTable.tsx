'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Professionals Table
 */

import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    EditIcon,
    EyeIcon,
    MoreHorizontalIcon,
    Trash2Icon,
    UserXIcon,
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useCallback } from 'react';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
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
import { maskPhone } from '@/lib/validations/professional';
import type { ProfessionalResponse } from '@/types/professional';

import {
    ProfessionalStatusBadge,
    ProfessionalTypeBadge,
} from './ProfessionalBadge';

// =============================================================================
// TYPES
// =============================================================================

interface ProfessionalsTableProps {
  professionals: ProfessionalResponse[];
  onView: (professional: ProfessionalResponse) => void;
  onEdit: (professional: ProfessionalResponse) => void;
  onDismiss: (professional: ProfessionalResponse) => void;
  onDelete: (professional: ProfessionalResponse) => void;
  isLoading?: boolean;
}

// =============================================================================
// HELPERS
// =============================================================================

function getInitials(name: string): string {
  return name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase();
}

function formatDate(dateStr: string): string {
  try {
    return format(new Date(dateStr), "dd/MM/yyyy", { locale: ptBR });
  } catch {
    return dateStr;
  }
}

function formatCommission(
  tipoComissao?: string,
  comissao?: number | string
): string {
  if (!comissao) return '-';
  const comissaoNum = typeof comissao === 'string' ? parseFloat(comissao) : comissao;
  if (isNaN(comissaoNum)) return '-';
  if (tipoComissao === 'PERCENTUAL') {
    // Se o valor for menor que 1, assume que está em formato decimal (0.4 = 40%)
    // Se for maior ou igual a 1, assume que já está em formato percentual (40 = 40%)
    const percentual = comissaoNum < 1 ? comissaoNum * 100 : comissaoNum;
    return `${percentual.toFixed(0)}%`;
  }
  return `R$ ${comissaoNum.toFixed(2)}`;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ProfessionalsTable({
  professionals,
  onView,
  onEdit,
  onDismiss,
  onDelete,
  isLoading,
}: ProfessionalsTableProps) {
  const router = useRouter();

  const handleRowClick = useCallback(
    (professional: ProfessionalResponse) => {
      onView(professional);
    },
    [onView]
  );

  const handleEdit = useCallback(
    (professional: ProfessionalResponse) => (e: React.MouseEvent) => {
      e.stopPropagation();
      onEdit(professional);
    },
    [onEdit]
  );

  const handleDismiss = useCallback(
    (professional: ProfessionalResponse) => (e: React.MouseEvent) => {
      e.stopPropagation();
      onDismiss(professional);
    },
    [onDismiss]
  );

  const handleDelete = useCallback(
    (professional: ProfessionalResponse) => (e: React.MouseEvent) => {
      e.stopPropagation();
      onDelete(professional);
    },
    [onDelete]
  );

  if (professionals.length === 0 && !isLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-center border rounded-md bg-muted/5">
        <div className="bg-muted p-4 rounded-full mb-4">
          <UserXIcon className="size-8 text-muted-foreground" />
        </div>
        <h3 className="text-lg font-semibold">Nenhum profissional encontrado</h3>
        <p className="mt-1 text-sm text-muted-foreground max-w-sm">
          Ajuste os filtros ou cadastre um novo profissional para começar a gerenciar sua equipe.
        </p>
      </div>
    );
  }

  return (
    <div className="rounded-md border shadow-sm">
      <Table>
        <TableHeader>
          <TableRow className="bg-muted/50 hover:bg-muted/50">
            <TableHead className="w-[300px] pl-6">Profissional</TableHead>
            <TableHead>Cargo</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Telefone</TableHead>
            <TableHead>Comissão</TableHead>
            <TableHead>Admissão</TableHead>
            <TableHead className="w-[70px] text-right pr-6">Ações</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {professionals.map((professional) => (
            <TableRow
              key={professional.id}
              className="cursor-pointer group hover:bg-muted/30"
              onClick={() => handleRowClick(professional)}
            >
              <TableCell className="pl-6">
                <div className="flex items-center gap-3">
                  <Avatar className="size-9 border">
                    <AvatarImage src={professional.foto} alt={professional.nome} />
                    <AvatarFallback className="bg-primary/5 text-primary text-xs">
                      {getInitials(professional.nome)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium text-sm">{professional.nome}</p>
                    <p className="text-xs text-muted-foreground truncate max-w-[180px]">
                      {professional.email}
                    </p>
                  </div>
                </div>
              </TableCell>

              <TableCell>
                <ProfessionalTypeBadge type={professional.tipo} />
              </TableCell>

              <TableCell>
                <ProfessionalStatusBadge status={professional.status} />
              </TableCell>

              <TableCell className="text-sm text-muted-foreground">
                {maskPhone(professional.telefone)}
              </TableCell>

              <TableCell className="text-sm font-mono text-muted-foreground">
                {formatCommission(professional.tipo_comissao, professional.comissao)}
              </TableCell>

              <TableCell className="text-sm text-muted-foreground">
                {formatDate(professional.data_admissao)}
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
                    <DropdownMenuItem onClick={() => onView(professional)}>
                      <EyeIcon className="mr-2 size-4" />
                      Visualizar
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={handleEdit(professional)}>
                      <EditIcon className="mr-2 size-4" />
                      Editar
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      onClick={handleDismiss(professional)}
                      className="text-orange-600 focus:text-orange-600"
                      disabled={professional.status === 'DEMITIDO'}
                    >
                      <UserXIcon className="mr-2 size-4" />
                      Demitir
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={handleDelete(professional)}
                      className="text-destructive focus:text-destructive"
                    >
                      <Trash2Icon className="mr-2 size-4" />
                      Deletar Permanentemente
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

export default ProfessionalsTable;
