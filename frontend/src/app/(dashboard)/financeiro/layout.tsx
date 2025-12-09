'use client';

import { cn } from '@/lib/utils';
import {
    ArrowDownCircle,
    ArrowUpCircle,
    FileText,
    LayoutDashboard,
    TrendingUp,
} from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ReactNode } from 'react';

interface FinanceiroLayoutProps {
  children: ReactNode;
}

const navItems = [
  {
    title: 'Dashboard',
    href: '/financeiro',
    icon: LayoutDashboard,
  },
  {
    title: 'Contas a Pagar',
    href: '/financeiro/contas-pagar',
    icon: ArrowDownCircle,
  },
  {
    title: 'Contas a Receber',
    href: '/financeiro/contas-receber',
    icon: ArrowUpCircle,
  },
  {
    title: 'DRE',
    href: '/financeiro/dre',
    icon: FileText,
  },
  {
    title: 'Fluxo de Caixa',
    href: '/financeiro/fluxo-caixa',
    icon: TrendingUp,
  },
];

export default function FinanceiroLayout({ children }: FinanceiroLayoutProps) {
  const pathname = usePathname();

  return (
    <div className="flex flex-col gap-6">
      {/* Header com Navegação */}
      <div className="border-b">
        <div className="flex items-center gap-1 overflow-x-auto pb-2">
          {navItems.map((item) => {
            const isActive = pathname === item.href || 
              (item.href !== '/financeiro' && pathname?.startsWith(item.href));
            const Icon = item.icon;

            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap',
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
        </div>
      </div>

      {/* Conteúdo */}
      <div>{children}</div>
    </div>
  );
}
