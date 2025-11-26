'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Componente de Seleção de Cliente
 *
 * @component CustomerSelector
 * @description Combobox com busca para selecionar cliente
 */

import { Loader2Icon, PlusIcon, SearchIcon } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

// =============================================================================
// TYPES
// =============================================================================

interface Customer {
  id: string;
  name: string;
  phone?: string;
  email?: string;
  avatar_url?: string;
}

interface CustomerSelectorProps {
  /** ID do cliente selecionado */
  value: string;
  /** Callback quando a seleção muda */
  onChange: (value: string) => void;
  /** Se está desabilitado */
  disabled?: boolean;
  /** Callback para criar novo cliente */
  onCreateNew?: () => void;
  /** Classe CSS adicional */
  className?: string;
}

// =============================================================================
// MOCK DATA (temporário até API estar pronta)
// =============================================================================

const MOCK_CUSTOMERS: Customer[] = [
  { id: '1', name: 'João Silva', phone: '(11) 99999-1234', email: 'joao@email.com' },
  { id: '2', name: 'Pedro Santos', phone: '(11) 98888-5678', email: 'pedro@email.com' },
  { id: '3', name: 'Carlos Oliveira', phone: '(11) 97777-9012' },
  { id: '4', name: 'Ricardo Lima', phone: '(11) 96666-3456', email: 'ricardo@email.com' },
  { id: '5', name: 'Fernando Costa', phone: '(11) 95555-7890' },
  { id: '6', name: 'Bruno Almeida', phone: '(11) 94444-2345', email: 'bruno@email.com' },
  { id: '7', name: 'André Pereira', phone: '(11) 93333-6789' },
  { id: '8', name: 'Marcos Souza', phone: '(11) 92222-0123', email: 'marcos@email.com' },
];

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

export function CustomerSelector({
  value,
  onChange,
  disabled = false,
  onCreateNew,
  className,
}: CustomerSelectorProps) {
  const [search, setSearch] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  // TODO: Substituir por useQuery quando a API de customers estiver pronta
  const customers = MOCK_CUSTOMERS;
  const isLoading = false;

  // Filtrar clientes pela busca
  const filteredCustomers = useMemo(() => {
    if (!search.trim()) return customers;
    const searchLower = search.toLowerCase();
    const searchNumbers = search.replace(/\D/g, '');
    
    return customers.filter(
      (c) =>
        c.name.toLowerCase().includes(searchLower) ||
        (c.phone && c.phone.replace(/\D/g, '').includes(searchNumbers)) ||
        (c.email && c.email.toLowerCase().includes(searchLower))
    );
  }, [customers, search]);

  // Cliente selecionado
  const selectedCustomer = useMemo(() => {
    return customers.find((c) => c.id === value);
  }, [customers, value]);

  // Selecionar cliente
  const selectCustomer = useCallback(
    (customerId: string) => {
      if (disabled) return;
      onChange(customerId);
      setIsOpen(false);
      setSearch('');
    },
    [onChange, disabled]
  );

  // Limpar seleção
  const clearSelection = useCallback(() => {
    if (disabled) return;
    onChange('');
    setSearch('');
  }, [onChange, disabled]);

  // ==========================================================================
  // RENDER - Cliente selecionado
  // ==========================================================================

  if (selectedCustomer) {
    return (
      <div
        className={cn(
          'flex items-center justify-between rounded-md border border-input bg-background px-3 py-2',
          disabled && 'opacity-50',
          className
        )}
      >
        <div className="flex items-center gap-3">
          <Avatar className="size-8">
            <AvatarImage
              src={selectedCustomer.avatar_url}
              alt={selectedCustomer.name}
            />
            <AvatarFallback className="text-xs">
              {getInitials(selectedCustomer.name)}
            </AvatarFallback>
          </Avatar>
          <div className="flex flex-col">
            <span className="text-sm font-medium">{selectedCustomer.name}</span>
            {selectedCustomer.phone && (
              <span className="text-xs text-muted-foreground">
                {selectedCustomer.phone}
              </span>
            )}
          </div>
        </div>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={clearSelection}
          disabled={disabled}
          className="text-muted-foreground hover:text-foreground"
        >
          Alterar
        </Button>
      </div>
    );
  }

  // ==========================================================================
  // RENDER - Busca
  // ==========================================================================

  return (
    <div className={cn('relative', className)}>
      {/* Input de busca */}
      <div className="relative">
        <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
        <Input
          placeholder="Buscar por nome ou telefone..."
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

      {/* Dropdown de clientes */}
      {isOpen && (
        <>
          {/* Overlay para fechar */}
          <div
            className="fixed inset-0 z-40"
            onClick={() => setIsOpen(false)}
          />

          {/* Lista de clientes */}
          <div className="absolute z-50 mt-1 w-full max-h-60 overflow-auto rounded-md border bg-popover p-1 shadow-md">
            {filteredCustomers.length === 0 ? (
              <div className="py-6 text-center">
                <p className="text-sm text-muted-foreground mb-3">
                  Nenhum cliente encontrado
                </p>
                {onCreateNew && (
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => {
                      setIsOpen(false);
                      onCreateNew();
                    }}
                  >
                    <PlusIcon className="size-4 mr-2" />
                    Cadastrar novo cliente
                  </Button>
                )}
              </div>
            ) : (
              <>
                {filteredCustomers.map((customer) => (
                  <button
                    key={customer.id}
                    type="button"
                    onClick={() => selectCustomer(customer.id)}
                    disabled={disabled}
                    className={cn(
                      'flex w-full items-center gap-3 rounded-sm px-2 py-2 text-sm outline-none',
                      'hover:bg-accent hover:text-accent-foreground',
                      value === customer.id && 'bg-accent'
                    )}
                  >
                    <Avatar className="size-8">
                      <AvatarImage
                        src={customer.avatar_url}
                        alt={customer.name}
                      />
                      <AvatarFallback className="text-xs">
                        {getInitials(customer.name)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex flex-col items-start">
                      <span className="font-medium">{customer.name}</span>
                      {customer.phone && (
                        <span className="text-xs text-muted-foreground">
                          {customer.phone}
                        </span>
                      )}
                    </div>
                  </button>
                ))}

                {/* Opção de criar novo */}
                {onCreateNew && (
                  <>
                    <div className="my-1 h-px bg-border" />
                    <button
                      type="button"
                      onClick={() => {
                        setIsOpen(false);
                        onCreateNew();
                      }}
                      className={cn(
                        'flex w-full items-center gap-2 rounded-sm px-2 py-2 text-sm outline-none',
                        'text-primary hover:bg-accent'
                      )}
                    >
                      <PlusIcon className="size-4" />
                      <span>Cadastrar novo cliente</span>
                    </button>
                  </>
                )}
              </>
            )}
          </div>
        </>
      )}
    </div>
  );
}

export default CustomerSelector;
