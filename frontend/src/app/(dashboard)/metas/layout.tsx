'use client';

import { cn } from '@/lib/utils';
import { BarChart3, Target, TrendingUp, Users } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ReactNode } from 'react';

interface MetasLayoutProps {
  children: ReactNode;
}

const tabs = [
  {
    href: '/metas',
    label: 'Dashboard',
    icon: BarChart3,
    exact: true,
  },
  {
    href: '/metas/mensais',
    label: 'Metas Mensais',
    icon: Target,
  },
  {
    href: '/metas/barbeiros',
    label: 'Por Barbeiro',
    icon: Users,
  },
  {
    href: '/metas/ticket',
    label: 'Ticket Médio',
    icon: TrendingUp,
  },
];

export default function MetasLayout({ children }: MetasLayoutProps) {
  const pathname = usePathname();

  const isActive = (href: string, exact?: boolean) => {
    if (exact) {
      return pathname === href;
    }
    return pathname.startsWith(href);
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Metas</h1>
          <p className="text-muted-foreground">
            Gerencie metas de faturamento, barbeiros e ticket médio
          </p>
        </div>
      </div>

      {/* Navigation Tabs */}
      <nav className="flex gap-1 border-b pb-px overflow-x-auto">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const active = isActive(tab.href, tab.exact);

          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={cn(
                'flex items-center gap-2 px-4 py-2 text-sm font-medium transition-colors rounded-t-md whitespace-nowrap',
                'hover:bg-muted/50',
                active
                  ? 'border-b-2 border-primary text-primary bg-muted/30'
                  : 'text-muted-foreground'
              )}
            >
              <Icon className="h-4 w-4" />
              {tab.label}
            </Link>
          );
        })}
      </nav>

      {/* Page Content */}
      <div className="flex-1">{children}</div>
    </div>
  );
}
