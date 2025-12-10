/**
 * NEXO - Sistema de Gestão para Barbearias
 * Dashboard Home Page
 *
 * Página inicial do dashboard com visão geral, cards estatísticos e atalhos.
 */

'use client';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useAppointments } from '@/hooks/use-appointments';
import { useCustomerBirthdays, useCustomerStats } from '@/hooks/use-customers';
import { useDashboard } from '@/hooks/use-financial';
import { useBreadcrumbs } from '@/store/ui-store';
import { endOfMonth, endOfToday, format, parseISO, startOfMonth, startOfToday } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
  CalendarIcon,
  DollarSignIcon,
  GiftIcon,
  LayoutDashboardIcon,
  PlusIcon,
  ScissorsIcon,
  UsersIcon,
  WalletIcon
} from 'lucide-react';
import Link from 'next/link';
import { useEffect } from 'react';

export default function DashboardPage() {
  const { setBreadcrumbs } = useBreadcrumbs();

  // Data Hooks
  const customersStats = useCustomerStats();
  const dashboardStats = useDashboard();

  // Appointments Today
  const today = new Date();
  const { data: todayAppointments, isLoading: isLoadingAppointments } = useAppointments({
    start_date: format(startOfToday(), 'yyyy-MM-dd'),
    end_date: format(endOfToday(), 'yyyy-MM-dd'),
    page_size: 100 // Get all for today
  });

  // Birthdays this month
  const { data: birthdays, isLoading: isLoadingBirthdays } = useCustomerBirthdays(
    format(startOfMonth(today), 'yyyy-MM-dd'),
    format(endOfMonth(today), 'yyyy-MM-dd')
  );

  useEffect(() => {
    setBreadcrumbs([{ label: 'Dashboard' }]);
  }, [setBreadcrumbs]);

  // Derived Values
  const totalClientes = customersStats.data?.total_geral || 0;
  const totalReceita = dashboardStats.data?.receita_total ? parseFloat(dashboardStats.data.receita_total) : 0;
  const totalAgendamentosHoje = todayAppointments?.total || 0;

  // Filter next appointments (limit 5)
  // Note: API already filters by date, we just take the first 5 assuming they are ordered by start_time
  const nextAppointments = todayAppointments?.data?.slice(0, 5) || [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Visão geral do seu negócio hoje, {format(today, "EEEE, d 'de' MMMM", { locale: ptBR })}.
          </p>
        </div>
        <div className="flex gap-2">
          <Button asChild>
            <Link href="/agendamentos">
              <PlusIcon className="mr-2 size-4" /> Novo Agendamento
            </Link>
          </Button>
        </div>
      </div>

      {/* Metrics Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Clientes</CardTitle>
            <UsersIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {customersStats.isLoading ? '...' : totalClientes}
            </div>
            <p className="text-xs text-muted-foreground">
              Cadastrados no sistema
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Agendamentos Hoje</CardTitle>
            <CalendarIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {isLoadingAppointments ? '...' : totalAgendamentosHoje}
            </div>
            <p className="text-xs text-muted-foreground">
              Para o dia de hoje
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Receita Estimada (Mês)</CardTitle>
            <DollarSignIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {dashboardStats.isLoading
                ? '...'
                : new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(totalReceita)
              }
            </div>
            <p className="text-xs text-muted-foreground">
              +20.1% em relação ao mês anterior (mock)
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Status do Sistema</CardTitle>
            <LayoutDashboardIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">Ativo</div>
            <p className="text-xs text-muted-foreground">
              Todos os serviços operando
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions & Recent Activity */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Main Content Area - 4 cols */}
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Próximos Agendamentos</CardTitle>
            <CardDescription>
              Você tem {totalAgendamentosHoje} agendamentos para hoje.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingAppointments ? (
              <div className="flex justify-center p-4">Carregando...</div>
            ) : nextAppointments.length > 0 ? (
              <div className="space-y-4">
                {nextAppointments.map((apt) => (
                  <div key={apt.id} className="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 transition-colors">
                    <div className="flex items-center gap-4">
                      <div className="flex flex-col items-center justify-center size-10 rounded bg-primary/10 text-primary font-bold text-xs">
                        <span>{format(parseISO(apt.start_time), 'HH:mm')}</span>
                      </div>
                      <div>
                        <p className="font-medium text-sm">{apt.customer_name}</p>
                        <p className="text-xs text-muted-foreground truncate max-w-[200px]">
                          {apt.services.map(s => s.service_name).join(', ')}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={
                        apt.status === 'CONFIRMED' ? 'default' :
                          apt.status === 'DONE' ? 'secondary' :
                            'outline'
                      }>
                        {apt.status_display}
                      </Badge>
                    </div>
                  </div>
                ))}
                <Button variant="ghost" className="w-full text-xs" asChild>
                  <Link href="/agendamentos">Ver todos os agendamentos</Link>
                </Button>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
                <CalendarIcon className="size-10 mb-2 opacity-20" />
                <p className="text-sm">Nenhum agendamento para hoje.</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Sidebar Area - 3 cols */}
        <div className="col-span-3 space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Atalhos Rápidos</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-2">
              <Button variant="outline" className="h-20 flex flex-col gap-1 items-center justify-center border-dashed" asChild>
                <Link href="/clientes?action=new">
                  <UsersIcon className="size-5 mb-1" />
                  <span>Novo Cliente</span>
                </Link>
              </Button>
              <Button variant="outline" className="h-20 flex flex-col gap-1 items-center justify-center border-dashed" asChild>
                <Link href="/financeiro?tab=receivables&action=new">
                  <WalletIcon className="size-5 mb-1" />
                  <span>Receita</span>
                </Link>
              </Button>
              <Button variant="outline" className="h-20 flex flex-col gap-1 items-center justify-center border-dashed" asChild>
                <Link href="/financeiro?tab=payables&action=new">
                  <DollarSignIcon className="size-5 mb-1" />
                  <span>Despesa</span>
                </Link>
              </Button>
              <Button variant="outline" className="h-20 flex flex-col gap-1 items-center justify-center border-dashed" asChild>
                <Link href="/cadastros/servicos">
                  <ScissorsIcon className="size-5 mb-1" />
                  <span>Serviços</span>
                </Link>
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <GiftIcon className="size-4 text-pink-500" />
                Aniversariantes do Mês
              </CardTitle>
            </CardHeader>
            <CardContent>
              {isLoadingBirthdays ? (
                <div className="text-sm text-muted-foreground">Carregando...</div>
              ) : birthdays && birthdays.length > 0 ? (
                <ScrollArea className="h-[200px]">
                  <div className="space-y-3">
                    {birthdays.map((customer: any) => (
                      <div key={customer.id} className="flex items-center gap-3 text-sm">
                        <Avatar className="size-8">
                          <AvatarImage src={customer.avatar_url} />
                          <AvatarFallback>{customer.nome.substring(0, 2).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1 overflow-hidden">
                          <p className="font-medium truncate">{customer.nome}</p>
                          <p className="text-xs text-muted-foreground">
                            Dia {customer.data_nascimento ? format(parseISO(customer.data_nascimento), 'dd/MM') : '--/--'}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              ) : (
                <p className="text-sm text-muted-foreground">Nenhum aniversariante encontrado.</p>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
