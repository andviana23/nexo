'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente de Seleção de Profissional (Barbeiro)
 *
 * @component ProfessionalSelector
 * @description Select para escolher o barbeiro do agendamento
 */

import { Loader2Icon } from 'lucide-react';
import { useMemo } from 'react';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { useProfessionals } from '@/hooks/use-appointments';
import { cn } from '@/lib/utils';

// =============================================================================
// TYPES
// =============================================================================

interface ProfessionalSelectorProps {
  /** ID do profissional selecionado */
  value: string;
  /** Callback quando a seleção muda */
  onChange: (value: string) => void;
  /** Se está desabilitado */
  disabled?: boolean;
  /** Placeholder quando nenhum está selecionado */
  placeholder?: string;
  /** Classe CSS adicional */
  className?: string;
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

// =============================================================================
// COMPONENT
// =============================================================================

export function ProfessionalSelector({
  value,
  onChange,
  disabled = false,
  placeholder = 'Selecione um barbeiro',
  className,
}: ProfessionalSelectorProps) {
  const { data: professionals = [], isLoading, isError } = useProfessionals();

  // Profissional selecionado
  const selectedProfessional = useMemo(() => {
    return professionals.find((p) => p.id === value);
  }, [professionals, value]);

  // Loading state
  if (isLoading) {
    return (
      <div
        className={cn(
          'flex h-9 w-full items-center justify-center rounded-md border border-input bg-background px-3',
          className
        )}
      >
        <Loader2Icon className="size-4 animate-spin text-muted-foreground" />
      </div>
    );
  }

  // Error state
  if (isError) {
    return (
      <div
        className={cn(
          'flex h-9 w-full items-center rounded-md border border-destructive bg-background px-3 text-sm text-destructive',
          className
        )}
      >
        Erro ao carregar barbeiros
      </div>
    );
  }

  // Empty state
  if (professionals.length === 0) {
    return (
      <div
        className={cn(
          'flex h-9 w-full items-center rounded-md border border-input bg-background px-3 text-sm text-muted-foreground',
          className
        )}
      >
        Nenhum barbeiro cadastrado
      </div>
    );
  }

  return (
    <Select value={value} onValueChange={onChange} disabled={disabled}>
      <SelectTrigger className={cn('w-full', className)}>
        <SelectValue placeholder={placeholder}>
          {selectedProfessional && (
            <div className="flex items-center gap-2">
              <Avatar className="size-5">
                <AvatarImage
                  src={selectedProfessional.avatar_url}
                  alt={selectedProfessional.name}
                />
                <AvatarFallback className="text-[10px]">
                  {getInitials(selectedProfessional.name)}
                </AvatarFallback>
              </Avatar>
              <span>{selectedProfessional.name}</span>
            </div>
          )}
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {professionals.map((professional) => (
          <SelectItem
            key={professional.id}
            value={professional.id}
            className="cursor-pointer"
          >
            <div className="flex items-center gap-2">
              <Avatar className="size-6">
                <AvatarImage
                  src={professional.avatar_url}
                  alt={professional.name}
                />
                <AvatarFallback className="text-[10px]">
                  {getInitials(professional.name)}
                </AvatarFallback>
              </Avatar>
              <div className="flex flex-col">
                <span>{professional.name}</span>
                {professional.google_calendar_connected && (
                  <span className="text-xs text-muted-foreground">
                    Google Calendar conectado
                  </span>
                )}
              </div>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

export default ProfessionalSelector;
