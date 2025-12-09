/**
 * NEXO - Sistema de Gestão para Barbearias
 * Header Component
 *
 * Barra superior com:
 * - Menu mobile (toggle sidebar)
 * - Breadcrumbs
 * - Unit selector (multi-unit)
 * - User navigation
 * - Theme toggle
 */

'use client';

import { UnitSelector } from '@/components/multi-unit';
import { cn } from '@/lib/utils';
import { useBreadcrumbs } from '@/store/ui-store';
import Link from 'next/link';
import { UserNav } from './UserNav';

export function Header() {
  const { breadcrumbs } = useBreadcrumbs();

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center gap-4 border-b bg-card px-4 md:px-6">
      {/* Mobile Menu Toggle (apenas visível no mobile) */}
      <div className="md:hidden">
        {/* O botão já está no Sidebar como SheetTrigger, então aqui é redundante */}
        {/* Mas mantemos o espaço para consistência */}
      </div>

      {/* Breadcrumbs */}
      <div className="flex flex-1 items-center gap-2 text-sm">
        {breadcrumbs.length > 0 ? (
          <nav className="flex items-center gap-2 text-muted-foreground">
            {breadcrumbs.map((crumb, index) => {
              const isLast = index === breadcrumbs.length - 1;

              return (
                <div key={index} className="flex items-center gap-2">
                  {crumb.href ? (
                    <Link
                      href={crumb.href}
                      className="hover:text-foreground transition-colors"
                    >
                      {crumb.label}
                    </Link>
                  ) : (
                    <span
                      className={cn(
                        isLast && 'text-foreground font-medium'
                      )}
                    >
                      {crumb.label}
                    </span>
                  )}

                  {!isLast && <span>/</span>}
                </div>
              );
            })}
          </nav>
        ) : (
          <div className="h-5" /> // Espaço reservado
        )}
      </div>

      {/* Unit Selector - aparece apenas se multi-unit */}
      <UnitSelector
        variant="outline"
        size="sm"
        collapsible
      />

      {/* User Navigation */}
      <UserNav />
    </header>
  );
}
