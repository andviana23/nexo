'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente de Seleção de Serviços
 *
 * @component ServiceSelector
 * @description Multi-select para serviços com busca, preço e duração
 */

import { CheckIcon, Loader2Icon, SearchIcon, XIcon } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';

import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

// =============================================================================
// TYPES
// =============================================================================

interface Service {
  id: string;
  name: string;
  default_price: number; // em centavos
  default_duration: number; // em minutos
  category?: string;
}

interface ServiceSelectorProps {
  /** IDs dos serviços selecionados */
  value: string[];
  /** Callback quando a seleção muda */
  onChange: (value: string[]) => void;
  /** Se está desabilitado */
  disabled?: boolean;
  /** Classe CSS adicional */
  className?: string;
}

// =============================================================================
// MOCK DATA (temporário até API estar pronta)
// =============================================================================

const MOCK_SERVICES: Service[] = [
  { id: '1', name: 'Corte Masculino', default_price: 4500, default_duration: 30, category: 'Corte' },
  { id: '2', name: 'Barba Completa', default_price: 3000, default_duration: 20, category: 'Barba' },
  { id: '3', name: 'Corte + Barba', default_price: 7000, default_duration: 50, category: 'Combo' },
  { id: '4', name: 'Pigmentação', default_price: 5000, default_duration: 45, category: 'Tratamento' },
  { id: '5', name: 'Relaxamento', default_price: 8000, default_duration: 60, category: 'Tratamento' },
  { id: '6', name: 'Hidratação', default_price: 4000, default_duration: 30, category: 'Tratamento' },
  { id: '7', name: 'Sobrancelha', default_price: 1500, default_duration: 10, category: 'Acabamento' },
  { id: '8', name: 'Pezinho', default_price: 1000, default_duration: 10, category: 'Acabamento' },
];

// =============================================================================
// HELPERS
// =============================================================================

function formatPrice(cents: number): string {
  return `R$ ${(cents / 100).toFixed(2)}`;
}

function formatDuration(minutes: number): string {
  if (minutes < 60) return `${minutes}min`;
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return mins > 0 ? `${hours}h ${mins}min` : `${hours}h`;
}

// =============================================================================
// COMPONENT
// =============================================================================

export function ServiceSelector({
  value,
  onChange,
  disabled = false,
  className,
}: ServiceSelectorProps) {
  const [search, setSearch] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  // TODO: Substituir por useQuery quando a API de services estiver pronta
  const services = MOCK_SERVICES;
  const isLoading = false;

  // Filtrar serviços pela busca
  const filteredServices = useMemo(() => {
    if (!search.trim()) return services;
    const searchLower = search.toLowerCase();
    return services.filter(
      (s) =>
        s.name.toLowerCase().includes(searchLower) ||
        s.category?.toLowerCase().includes(searchLower)
    );
  }, [services, search]);

  // Serviços selecionados com detalhes
  const selectedServices = useMemo(() => {
    return services.filter((s) => value.includes(s.id));
  }, [services, value]);

  // Totais
  const totals = useMemo(() => {
    return selectedServices.reduce(
      (acc, s) => ({
        price: acc.price + s.default_price,
        duration: acc.duration + s.default_duration,
      }),
      { price: 0, duration: 0 }
    );
  }, [selectedServices]);

  // Toggle seleção
  const toggleService = useCallback(
    (serviceId: string) => {
      if (disabled) return;
      
      if (value.includes(serviceId)) {
        onChange(value.filter((id) => id !== serviceId));
      } else {
        onChange([...value, serviceId]);
      }
    },
    [value, onChange, disabled]
  );

  // Remover serviço da seleção
  const removeService = useCallback(
    (serviceId: string) => {
      if (disabled) return;
      onChange(value.filter((id) => id !== serviceId));
    },
    [value, onChange, disabled]
  );

  // Agrupar por categoria
  const groupedServices = useMemo(() => {
    const groups: Record<string, Service[]> = {};
    filteredServices.forEach((service) => {
      const category = service.category || 'Outros';
      if (!groups[category]) groups[category] = [];
      groups[category].push(service);
    });
    return groups;
  }, [filteredServices]);

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <div className={cn('relative', className)}>
      {/* Serviços Selecionados */}
      {selectedServices.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-2">
          {selectedServices.map((service) => (
            <Badge
              key={service.id}
              variant="secondary"
              className="flex items-center gap-1 pr-1"
            >
              {service.name}
              <button
                type="button"
                onClick={() => removeService(service.id)}
                disabled={disabled}
                className="ml-1 hover:bg-muted rounded-full p-0.5"
              >
                <XIcon className="size-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}

      {/* Input de busca */}
      <div className="relative">
        <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
        <Input
          placeholder="Buscar serviços..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onFocus={() => setIsOpen(true)}
          disabled={disabled}
          className="pl-10"
        />
        {isLoading && (
          <Loader2Icon className="absolute right-3 top-1/2 -translate-y-1/2 size-4 animate-spin text-muted-foreground" />
        )}
      </div>

      {/* Dropdown de serviços */}
      {isOpen && (
        <>
          {/* Overlay para fechar */}
          <div
            className="fixed inset-0 z-40"
            onClick={() => setIsOpen(false)}
          />
          
          {/* Lista de serviços */}
          <div className="absolute z-50 mt-1 w-full max-h-60 overflow-auto rounded-md border bg-popover p-1 shadow-md">
            {Object.entries(groupedServices).length === 0 ? (
              <div className="py-6 text-center text-sm text-muted-foreground">
                Nenhum serviço encontrado
              </div>
            ) : (
              Object.entries(groupedServices).map(([category, categoryServices]) => (
                <div key={category}>
                  <div className="px-2 py-1.5 text-xs font-semibold text-muted-foreground">
                    {category}
                  </div>
                  {categoryServices.map((service) => {
                    const isSelected = value.includes(service.id);
                    return (
                      <button
                        key={service.id}
                        type="button"
                        onClick={() => toggleService(service.id)}
                        disabled={disabled}
                        className={cn(
                          'flex w-full items-center justify-between rounded-sm px-2 py-1.5 text-sm outline-none',
                          'hover:bg-accent hover:text-accent-foreground',
                          isSelected && 'bg-accent'
                        )}
                      >
                        <div className="flex items-center gap-2">
                          <div
                            className={cn(
                              'size-4 rounded border flex items-center justify-center',
                              isSelected
                                ? 'bg-primary border-primary text-primary-foreground'
                                : 'border-input'
                            )}
                          >
                            {isSelected && <CheckIcon className="size-3" />}
                          </div>
                          <span>{service.name}</span>
                        </div>
                        <div className="flex items-center gap-3 text-muted-foreground">
                          <span className="text-xs">
                            {formatDuration(service.default_duration)}
                          </span>
                          <span className="font-medium text-foreground">
                            {formatPrice(service.default_price)}
                          </span>
                        </div>
                      </button>
                    );
                  })}
                </div>
              ))
            )}
          </div>
        </>
      )}

      {/* Resumo de totais */}
      {selectedServices.length > 0 && (
        <div className="mt-2 flex items-center justify-between text-sm text-muted-foreground">
          <span>
            {selectedServices.length} serviço{selectedServices.length > 1 ? 's' : ''} •{' '}
            {formatDuration(totals.duration)}
          </span>
          <span className="font-medium text-foreground">
            Total: {formatPrice(totals.price)}
          </span>
        </div>
      )}
    </div>
  );
}

export default ServiceSelector;
