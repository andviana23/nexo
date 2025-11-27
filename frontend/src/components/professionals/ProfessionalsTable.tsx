'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Professionals Table
 *
 * @component ProfessionalsTable
 * @description Tabela de listagem de profissionais com ações
 */

import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    EditIcon,
    ExternalLinkIcon,
    EyeIcon,
    MoreHorizontalIcon,
    UserXIcon,
} from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useCallback } from 'react';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
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
  /** Lista de profissionais */
  professionals: ProfessionalResponse[];
  /** Callback ao clicar em visualizar */
  onView: (professional: ProfessionalResponse) => void;
  /** Callback ao clicar em editar */
  onEdit: (professional: ProfessionalResponse) => void;
  /** Callback ao clicar em desativar */
  onDeactivate: (professional: ProfessionalResponse) => void;
  /** Loading state */
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
  return format(new Date(dateStr), "dd 'de' MMM, yyyy", { locale: ptBR });
}

function formatCommission(
  tipoComissao?: string,
  comissao?: number | string
): string {
  if (!comissao) return '-';
  const comissaoNum = typeof comissao === 'string' ? parseFloat(comissao) : comissao;
  if (isNaN(comissaoNum)) return '-';
  if (tipoComissao === 'PERCENTUAL') return `${comissaoNum.toFixed(2)}%`;
  return `R$ ${comissaoNum.toFixed(2)}`;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ProfessionalsTable({
  professionals,
  onView,
  onEdit,
  onDeactivate,
  isLoading,
}: ProfessionalsTableProps) {
  const router = useRouter();

  const handleRowClick = useCallback(
    (professional: ProfessionalResponse) => {
      router.push(`/profissionais/${professional.id}`);
    },
    [router]
  );

  const handleView = useCallback(
    (professional: ProfessionalResponse) => () => onView(professional),
    [onView]
  );

  const handleEdit = useCallback(
    (professional: ProfessionalResponse) => () => onEdit(professional),
    [onEdit]
  );

  const handleDeactivate = useCallback(
    (professional: ProfessionalResponse) => () => onDeactivate(professional),
    [onDeactivate]
  );

  if (professionals.length === 0 && !isLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <UserXIcon className="size-12 text-muted-foreground/50" />
        <h3 className="mt-4 text-lg font-semibold">Nenhum profissional encontrado</h3>
        <p className="mt-1 text-sm text-muted-foreground">
          Cadastre seu primeiro profissional para começar.
        </p>
      </div>
    );
  }

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[300px]">Profissional</TableHead>
            <TableHead>Tipo</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Telefone</TableHead>
            <TableHead>Comissão</TableHead>
            <TableHead>Admissão</TableHead>
            <TableHead className="w-[70px]" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {professionals.map((professional) => (
            <TableRow 
              key={professional.id}
              className="cursor-pointer hover:bg-muted/50"
              onClick={() => handleRowClick(professional)}
            >
              {/* Nome e Email */}
              <TableCell>
                <div className="flex items-center gap-3">
                  <Avatar className="size-10">
                    <AvatarImage src={professional.foto} alt={professional.nome} />
                    <AvatarFallback className="bg-primary/10 text-primary">
                      {getInitials(professional.nome)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-medium">{professional.nome}</p>
                    <p className="text-sm text-muted-foreground">
                      {professional.email}
                    </p>
                  </div>
                </div>
              </TableCell>

              {/* Tipo */}
              <TableCell>
                <ProfessionalTypeBadge type={professional.tipo} />
              </TableCell>

              {/* Status */}
              <TableCell>
                <ProfessionalStatusBadge status={professional.status} />
              </TableCell>

              {/* Telefone */}
              <TableCell className="text-muted-foreground">
                {maskPhone(professional.telefone)}
              </TableCell>

              {/* Comissão */}
              <TableCell className="text-muted-foreground">
                {formatCommission(professional.tipo_comissao, professional.comissao)}
              </TableCell>

              {/* Data Admissão */}
              <TableCell className="text-muted-foreground">
                {formatDate(professional.data_admissao)}
              </TableCell>

              {/* Ações */}
              <TableCell onClick={(e) => e.stopPropagation()}>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="size-8 p-0">
                      <span className="sr-only">Abrir menu</span>
                      <MoreHorizontalIcon className="size-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem asChild>
                      <Link href={`/profissionais/${professional.id}`}>
                        <ExternalLinkIcon className="mr-2 size-4" />
                        Abrir Detalhes
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={handleView(professional)}>
                      <EyeIcon className="mr-2 size-4" />
                      Visualizar
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={handleEdit(professional)}>
                      <EditIcon className="mr-2 size-4" />
                      Editar
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      onClick={handleDeactivate(professional)}
                      className="text-destructive focus:text-destructive"
                      disabled={professional.status === 'DEMITIDO'}
                    >
                      <UserXIcon className="mr-2 size-4" />
                      Desativar
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
