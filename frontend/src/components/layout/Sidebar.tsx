/**
 * NEXO - Sistema de Gestão para Barbearias
 * Sidebar Component
 *
 * Navegação lateral com:
 * - Collapse/Expand
 * - Responsivo (mobile drawer)
 * - Suporte a RBAC (filtro por permissões)
 * - Ícones lucide-react
 */

'use client';

import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetTitle, SheetTrigger } from '@/components/ui/sheet';
import { cn } from '@/lib/utils';
import { useCurrentUser } from '@/store/auth-store';
import { useSidebar } from '@/store/ui-store';
import {
    BarChart3,
    Calendar,
    ChevronLeft,
    ClipboardList,
    CreditCard,
    DollarSign,
    LayoutDashboard,
    Menu,
    Package,
    Scissors,
    Settings,
    Users,
} from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect } from 'react';

// =============================================================================
// TIPOS
// =============================================================================

interface NavItem {
  title: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  roles?: string[]; // Se definido, apenas essas roles podem ver
  badge?: string; // Badge opcional (ex: "Beta", "Novo")
}

// =============================================================================
// NAVEGAÇÃO
// =============================================================================

const navItems: NavItem[] = [
  {
    title: 'Dashboard',
    href: '/',
    icon: LayoutDashboard,
  },
  {
    title: 'Agendamentos',
    href: '/agendamentos',
    icon: Calendar,
  },
  {
    title: 'Lista da Vez',
    href: '/lista-da-vez',
    icon: ClipboardList,
  },
  {
    title: 'Clientes',
    href: '/clientes',
    icon: Users,
  },
  {
    title: 'Serviços',
    href: '/servicos',
    icon: Scissors,
  },
  {
    title: 'Estoque',
    href: '/estoque',
    icon: Package,
  },
  {
    title: 'Financeiro',
    href: '/financeiro',
    icon: DollarSign,
  },
  {
    title: 'Relatórios',
    href: '/relatorios',
    icon: BarChart3,
    roles: ['owner', 'admin'], // Apenas owner e admin
  },
  {
    title: 'Assinatura',
    href: '/assinatura',
    icon: CreditCard,
    roles: ['owner'], // Apenas owner
  },
  {
    title: 'Configurações',
    href: '/configuracoes',
    icon: Settings,
  },
];

// =============================================================================
// COMPONENTE
// =============================================================================

// =============================================================================
// COMPONENTE SIDEBAR CONTENT (EXTRAÍDO PARA EVITAR RECREAÇÃO)
// =============================================================================

function SidebarContent({
  isCollapsed,
  setCollapsed,
  pathname,
  filteredItems,
}: {
  isCollapsed: boolean;
  setCollapsed: (value: boolean) => void;
  pathname: string;
  filteredItems: NavItem[];
}) {
  return (
    <div
      className={cn(
        'flex h-full flex-col border-r bg-card transition-all duration-300',
        isCollapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Logo / Brand */}
      <div className="flex h-16 items-center border-b px-4">
        {!isCollapsed ? (
          <Link href="/" className="flex items-center gap-2 font-semibold">
            <Scissors className="h-6 w-6 text-primary" />
            <span className="text-xl">NEXO</span>
          </Link>
        ) : (
          <Link href="/" className="flex items-center justify-center">
            <Scissors className="h-6 w-6 text-primary" />
          </Link>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 overflow-y-auto p-2">
        {filteredItems.map((item) => {
          const isActive = pathname === item.href;
          const Icon = item.icon;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                isCollapsed && 'justify-center'
              )}
              title={isCollapsed ? item.title : undefined}
            >
              <Icon className="h-5 w-5 shrink-0" />
              {!isCollapsed && <span className="flex-1">{item.title}</span>}
              {!isCollapsed && item.badge && (
                <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
                  {item.badge}
                </span>
              )}
            </Link>
          );
        })}
      </nav>

      {/* Collapse Toggle (Desktop only) */}
      <div className="hidden border-t p-2 md:block">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setCollapsed(!isCollapsed)}
          className={cn('w-full', isCollapsed && 'justify-center')}
        >
          <ChevronLeft
            className={cn(
              'h-5 w-5 transition-transform',
              isCollapsed && 'rotate-180'
            )}
          />
          {!isCollapsed && <span className="ml-2">Recolher</span>}
        </Button>
      </div>
    </div>
  );
}

// =============================================================================
// COMPONENTE PRINCIPAL
// =============================================================================

export function Sidebar() {
  const { isCollapsed, setCollapsed, isOpen, setOpen } = useSidebar();
  const pathname = usePathname();
  const user = useCurrentUser();

  // Detecta mobile
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth < 768) {
        setCollapsed(false); // Mobile nunca está colapsado
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [setCollapsed]);

  // Filtra itens por role
  const filteredItems = navItems.filter((item) => {
    if (!item.roles) return true;
    if (!user?.role) return false;
    return item.roles.includes(user.role);
  });

  // Mobile: Sheet (Drawer)
  // Desktop: Sempre visível
  return (
    <>
      {/* Mobile */}
      <div className="md:hidden">
        <Sheet open={isOpen} onOpenChange={setOpen}>
          <SheetTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="fixed left-4 top-4 z-40 md:hidden"
            >
              <Menu className="h-6 w-6" />
              <span className="sr-only">Toggle Menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-64 p-0">
            <SheetTitle className="sr-only">Menu de Navegação</SheetTitle>
            <SidebarContent
              isCollapsed={false}
              setCollapsed={setCollapsed}
              pathname={pathname}
              filteredItems={filteredItems}
            />
          </SheetContent>
        </Sheet>
      </div>

      {/* Desktop */}
      <aside className="hidden md:block">
        <SidebarContent
          isCollapsed={isCollapsed}
          setCollapsed={setCollapsed}
          pathname={pathname}
          filteredItems={filteredItems}
        />
      </aside>
    </>
  );
}
