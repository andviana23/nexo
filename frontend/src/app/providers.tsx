/**
 * NEXO - Sistema de Gestão para Barbearias
 * Providers
 *
 * Componente que agrupa todos os providers da aplicação.
 * QueryClientProvider, ThemeProvider, Toaster, etc.
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
  // Cria QueryClient uma vez por montagem do componente
  // Isso garante que cada sessão no navegador tem seu próprio cache
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

        {/* Toast Notifications */}
        <Toaster
          position="top-right"
          expand={false}
          richColors
          closeButton
          duration={4000}
        />
      </ThemeProvider>

      {/* DevTools apenas em desenvolvimento */}
      {process.env.NODE_ENV === 'development' && (
        <ReactQueryDevtools
          initialIsOpen={false}
          buttonPosition="bottom-left"
        />
      )}
    </QueryClientProvider>
  );
}
