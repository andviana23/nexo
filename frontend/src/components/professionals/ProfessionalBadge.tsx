/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente: Professional Badge
 *
 * @module components/professionals/ProfessionalBadge
 * @description Badge estilizado para tipo e status de profissionais
 */

import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import {
    PROFESSIONAL_STATUS_COLORS,
    PROFESSIONAL_STATUS_LABELS,
    PROFESSIONAL_TYPE_COLORS,
    PROFESSIONAL_TYPE_LABELS,
    type ProfessionalStatus,
    type ProfessionalType,
} from '@/types/professional';

// =============================================================================
// TYPES
// =============================================================================

interface TypeBadgeProps {
  type: ProfessionalType;
  className?: string;
}

interface StatusBadgeProps {
  status: ProfessionalStatus;
  className?: string;
}

// =============================================================================
// COMPONENTS
// =============================================================================

/**
 * Badge para exibir tipo de profissional
 */
export function ProfessionalTypeBadge({ type, className }: TypeBadgeProps) {
  return (
    <Badge
      variant="secondary"
      className={cn(PROFESSIONAL_TYPE_COLORS[type], className)}
    >
      {PROFESSIONAL_TYPE_LABELS[type]}
    </Badge>
  );
}

/**
 * Badge para exibir status do profissional
 */
export function ProfessionalStatusBadge({ status, className }: StatusBadgeProps) {
  return (
    <Badge
      variant="secondary"
      className={cn(PROFESSIONAL_STATUS_COLORS[status], className)}
    >
      {PROFESSIONAL_STATUS_LABELS[status]}
    </Badge>
  );
}

// =============================================================================
// EXPORTS
// =============================================================================

export { type StatusBadgeProps, type TypeBadgeProps };

