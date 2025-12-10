/**
 * NEXO - Sistema de Gestão para Barbearias
 * Sidebar Component
 *
 * Navegação lateral com:
 * - Collapse/Expand
 * - Responsivo (mobile drawer)
 * - Suporte a RBAC (filtro por permissões)
 * - Menus colapsáveis (submenu)
 * - Ícones lucide-react
 */

'use client';

import { Button } from '@/components/ui/button';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { Sheet, SheetContent, SheetDescription, SheetTitle, SheetTrigger } from '@/components/ui/sheet';
import { cn } from '@/lib/utils';
import { useAuthStore, useCurrentUser } from '@/store/auth-store';
import { useSidebar } from '@/store/ui-store';
import {
  Banknote,
  BarChart3,
  Calculator,
  Calendar,
  CalendarDays,
  ChevronDown,
  ChevronLeft,
  ClipboardList,
  CreditCard,
  DollarSign,
  FileText,
  FolderOpen,
  HandCoins,
  LayoutDashboard,
  LogOut,
  Menu,
  Package,
  Receipt,
  RefreshCw,
  Scissors,
  ScrollText,
  Settings,
  Tags,
  TrendingDown,
  TrendingUp,
  Truck,
  UserCog,
  Users,
  Wallet
} from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';

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

interface NavGroup {
  title: string;
  icon: React.ComponentType<{ className?: string }>;
  items: NavItem[];
  roles?: string[];
}

type NavigationItem = NavItem | NavGroup;

// Type guard para verificar se é um grupo
function isNavGroup(item: NavigationItem): item is NavGroup {
  return 'items' in item;
}

// =============================================================================
// NAVEGAÇÃO
// =============================================================================

const navigationItems: NavigationItem[] = [
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
  // Menu colapsável: Cadastros
  {
    title: 'Cadastros',
    icon: FolderOpen,
    items: [
      {
        title: 'Clientes',
        href: '/clientes',
        icon: Users,
      },
      {
        title: 'Profissionais',
        href: '/profissionais',
        icon: UserCog,
        roles: ['owner', 'admin', 'manager'],
      },
      {
        title: 'Serviços',
        href: '/cadastros/servicos',
        icon: Scissors,
      },
      {
        title: 'Categorias',
        href: '/cadastros/categorias',
        icon: Tags,
      },
      {
        title: 'Tipos de Recebimento',
        href: '/cadastros/tipos-recebimento',
        icon: Receipt,
        roles: ['owner', 'admin', 'manager'],
      },
      {
        title: 'Fornecedores',
        href: '/cadastros/fornecedores',
        icon: Truck,
        roles: ['owner', 'admin', 'manager'],
      },
    ],
  },
  {
    title: 'Estoque',
    href: '/estoque',
    icon: Package,
  },
  // Menu colapsável: Financeiro
  {
    title: 'Financeiro',
    icon: DollarSign,
    items: [
      {
        title: 'Dashboard',
        href: '/financeiro',
        icon: Wallet,
      },
      {
        title: 'Caixa Diário',
        href: '/caixa',
        icon: Banknote,
      },
      {
        title: 'Contas a Pagar',
        href: '/financeiro/contas-pagar',
        icon: TrendingDown,
      },
      {
        title: 'Contas a Receber',
        href: '/financeiro/contas-receber',
        icon: TrendingUp,
      },
      {
        title: 'Despesas Fixas',
        href: '/financeiro/despesas-fixas',
        icon: RefreshCw,
      },
      {
        title: 'DRE',
        href: '/financeiro/dre',
        icon: FileText,
        roles: ['owner', 'admin'],
      },
      {
        title: 'Fluxo de Caixa',
        href: '/financeiro/fluxo-caixa',
        icon: BarChart3,
        roles: ['owner', 'admin'],
      },
    ],
  },
  // Menu colapsável: Comissões
  {
    title: 'Comissões',
    icon: Calculator,
    roles: ['owner', 'admin', 'manager'],
    items: [
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
    ],
  },
  {
    title: 'Relatórios',
    href: '/relatorios',
    icon: BarChart3,
    roles: ['owner', 'admin'], // Owner e admin
  },
  {
    title: 'Assinatura',
    href: '/assinatura',
    icon: CreditCard,
    roles: ['owner', 'admin'], // Owner e admin
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
  user,
}: {
  isCollapsed: boolean;
  setCollapsed: (value: boolean) => void;
  pathname: string;
  filteredItems: NavigationItem[];
  user: any; // Using any for simplicity here to match existing usage patterns or import User type
}) {
  // Estado para controlar menus abertos
  const [openGroups, setOpenGroups] = useState<string[]>(() => {
    // Abre automaticamente o grupo que contém a rota atual
    const activeGroup = filteredItems.find(
      (item) =>
        isNavGroup(item) &&
        item.items.some((subItem) => pathname.startsWith(subItem.href))
    );
    return activeGroup ? [activeGroup.title] : [];
  });

  const toggleGroup = (title: string) => {
    setOpenGroups((prev) =>
      prev.includes(title)
        ? prev.filter((t) => t !== title)
        : [...prev, title]
    );
  };

  return (
    <div
      className={cn(
        'flex h-full flex-col border-r bg-card transition-all duration-300',
        isCollapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Logo / Brand */}
      <div className="flex h-16 items-center justify-center border-b px-4 transition-all duration-300">
        <Link href="/" className="flex items-center justify-center">
          <div className={cn("relative shrink-0 transition-all duration-300", isCollapsed ? "h-10 w-10" : "h-12 w-40")}>
            <Image
              src="/nexo.png"
              alt="NEXO Logo"
              fill
              className="object-contain"
              priority
            />
          </div>
        </Link>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 overflow-y-auto p-2">
        {filteredItems.map((item) => {
          // Renderiza grupo colapsável
          if (isNavGroup(item)) {
            const Icon = item.icon;
            const isGroupOpen = openGroups.includes(item.title);
            const hasActiveChild = item.items.some(
              (subItem) => pathname === subItem.href
            );

            // Se sidebar está colapsada, mostrar apenas ícone com tooltip
            if (isCollapsed) {
              return (
                <div key={item.title} className="relative group">
                  <button
                    onClick={() => toggleGroup(item.title)}
                    className={cn(
                      'flex w-full items-center justify-center rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                      hasActiveChild
                        ? 'bg-primary/10 text-primary'
                        : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                    )}
                    title={item.title}
                  >
                    <Icon className="h-5 w-5 shrink-0" />
                  </button>
                  {/* Tooltip com subitems */}
                  <div className="absolute left-full top-0 ml-2 hidden w-48 rounded-lg border bg-popover p-2 shadow-lg group-hover:block z-50">
                    <div className="text-xs font-semibold text-muted-foreground mb-2 px-2">
                      {item.title}
                    </div>
                    {item.items.map((subItem) => {
                      const SubIcon = subItem.icon;
                      const isSubActive = pathname === subItem.href;
                      return (
                        <Link
                          key={subItem.href}
                          href={subItem.href}
                          className={cn(
                            'flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors',
                            isSubActive
                              ? 'bg-primary text-primary-foreground'
                              : 'hover:bg-accent'
                          )}
                        >
                          <SubIcon className="h-4 w-4" />
                          {subItem.title}
                        </Link>
                      );
                    })}
                  </div>
                </div>
              );
            }

            // Sidebar expandida: menu colapsável completo
            return (
              <Collapsible
                key={item.title}
                open={isGroupOpen}
                onOpenChange={() => toggleGroup(item.title)}
              >
                <CollapsibleTrigger asChild>
                  <button
                    className={cn(
                      'flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                      hasActiveChild
                        ? 'bg-primary/10 text-primary'
                        : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                    )}
                  >
                    <Icon className="h-5 w-5 shrink-0" />
                    <span className="flex-1 text-left">{item.title}</span>
                    <ChevronDown
                      className={cn(
                        'h-4 w-4 transition-transform',
                        isGroupOpen && 'rotate-180'
                      )}
                    />
                  </button>
                </CollapsibleTrigger>
                <CollapsibleContent className="pl-4 space-y-1 mt-1">
                  {item.items.map((subItem) => {
                    const SubIcon = subItem.icon;
                    const isSubActive = pathname === subItem.href;

                    return (
                      <Link
                        key={subItem.href}
                        href={subItem.href}
                        className={cn(
                          'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                          isSubActive
                            ? 'bg-primary text-primary-foreground'
                            : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                        )}
                      >
                        <SubIcon className="h-4 w-4 shrink-0" />
                        <span>{subItem.title}</span>
                        {subItem.badge && (
                          <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
                            {subItem.badge}
                          </span>
                        )}
                      </Link>
                    );
                  })}
                </CollapsibleContent>
              </Collapsible>
            );
          }

          // Renderiza item simples
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

      {/* Footer / User Profile & Logout */}
      <div className="border-t p-4">
        {!isCollapsed ? (
          <div className="flex items-center justify-between gap-2 rounded-lg bg-muted/50 p-2">
            <div className="flex items-center gap-2 overflow-hidden">
              <div className="flex flex-col truncate">
                <span className="text-sm font-medium truncate">{user?.name || 'Usuário'}</span>
                <span className="text-xs text-muted-foreground truncate capitalize">{user?.role || 'Cargo'}</span>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              className="shrink-0 text-muted-foreground hover:text-destructive transition-colors"
              onClick={() => useAuthStore.getState().logout()}
              title="Sair"
            >
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        ) : (
          <Button
            variant="ghost"
            size="icon"
            className="w-full text-muted-foreground hover:text-destructive transition-colors"
            onClick={() => useAuthStore.getState().logout()}
            title="Sair"
          >
            <LogOut className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* Collapse Toggle (Desktop only) */}
      <div className="hidden border-t p-2 md:block">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setCollapsed(!isCollapsed)}
          className={cn(
            'w-full flex items-center gap-2',
            isCollapsed ? 'justify-center' : 'justify-start'
          )}
          title={isCollapsed ? 'Expandir menu' : 'Recolher menu'}
        >
          <ChevronLeft
            className={cn(
              'h-5 w-5 transition-transform',
              isCollapsed && 'rotate-180'
            )}
          />
          {!isCollapsed && <span>Recolher</span>}
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

  // DEBUG: Log do usuário
  useEffect(() => {
    console.log('[Sidebar] User state:', {
      email: user?.email,
      role: user?.role,
      name: user?.name,
    });
  }, [user]);

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

  // Filtra itens por role (incluindo subitens de grupos)
  const filteredItems = navigationItems
    .map((item) => {
      // Se é um grupo, filtra os subitens
      if (isNavGroup(item)) {
        // Verifica se o grupo em si tem restrição de role
        if (item.roles && (!user?.role || !item.roles.includes(user.role))) {
          return null;
        }
        // Filtra subitens do grupo
        const filteredSubItems = item.items.filter((subItem) => {
          if (!subItem.roles) return true;
          if (!user?.role) return false;
          return subItem.roles.includes(user.role);
        });
        // Se não sobrou nenhum subitem, não mostra o grupo
        if (filteredSubItems.length === 0) return null;
        return { ...item, items: filteredSubItems };
      }
      // Item simples
      if (!item.roles) return item;
      if (!user?.role) return null;
      return item.roles.includes(user.role) ? item : null;
    })
    .filter((item): item is NavigationItem => item !== null);

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
            <SheetDescription className="sr-only">
              Navegue pelas opções do sistema NEXO
            </SheetDescription>
            <SidebarContent
              isCollapsed={false}
              setCollapsed={setCollapsed}
              pathname={pathname}
              filteredItems={filteredItems}
              user={user}
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
          user={user}
        />
      </aside>
    </>
  );
}
