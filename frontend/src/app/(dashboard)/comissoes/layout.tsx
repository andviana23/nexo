'use client';

import { cn } from '@/lib/utils';
import {
    Calculator,
    CalendarDays,
    HandCoins,
    LayoutDashboard,
    ScrollText,
} from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ReactNode } from 'react';

interface ComissoesLayoutProps {
  children: ReactNode;
}

const navItems = [
  {
    title: 'Dashboard',
    href: '/comissoes',
    icon: LayoutDashboard,
  },
  {
    title: 'Regras',
    href: '/comissoes/regras',
    icon: ScrollText,
  },
  {
    title: 'Períodos',
    href: '/comissoes/periodos',
    icon: CalendarDays,
  },
  {
    title: 'Adiantamentos',
    href: '/comissoes/adiantamentos',
    icon: HandCoins,
  },
  {
    title: 'Itens',
    href: '/comissoes/itens',
    icon: Calculator,
  },
];

export default function ComissoesLayout({ children }: ComissoesLayoutProps) {
  const pathname = usePathname();

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Comissões</h1>
        <p className="text-muted-foreground">
          Gerencie regras, períodos e adiantamentos de comissões dos profissionais.
        </p>
      </div>

      {/* Sub-navigation */}
      <nav className="flex gap-2 border-b pb-2 overflow-x-auto">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = item.href === '/comissoes' 
            ? pathname === '/comissoes'
            : pathname.startsWith(item.href);

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-md transition-colors whitespace-nowrap',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted'
              )}
            >
              <Icon className="h-4 w-4" />
              {item.title}
            </Link>
          );
        })}
      </nav>

      {/* Content */}
      <div>{children}</div>
    </div>
  );
}
