/**
 * NEXO - Sistema de Gestão para Barbearias
 * Providers Globais
 *
 * Componente que agrupa todos os providers da aplicação.
 * - QueryClientProvider (React Query)
 * - ThemeProvider (next-themes)
 * - Toaster (sonner)
 */

'use client';

import { Toaster } from '@/components/ui/sonner';
import { makeQueryClient } from '@/lib/query-client';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { ThemeProvider } from 'next-themes';
import { useState } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface ProvidersProps {
  children: React.ReactNode;
}

// =============================================================================
// COMPONENTE
// =============================================================================

export function Providers({ children }: ProvidersProps) {
  // Cria QueryClient uma única vez (useState garante isso)
  const [queryClient] = useState(() => makeQueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
      >
        {children}

        {/* Toast notifications */}
        <Toaster
          position="top-right"
          expand={false}
          richColors
          closeButton
          toastOptions={{
            duration: 4000,
            classNames: {
              toast: 'group',
              title: 'font-semibold',
              description: 'text-sm text-muted-foreground',
              actionButton: 'bg-primary text-primary-foreground',
              cancelButton: 'bg-muted text-muted-foreground',
            },
          }}
        />
      </ThemeProvider>

      {/* React Query Devtools (apenas em desenvolvimento) */}
      {process.env.NODE_ENV === 'development' && (
        <ReactQueryDevtools
          initialIsOpen={false}
          buttonPosition="bottom-left"
        />
      )}
    </QueryClientProvider>
  );
}
