/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Detalhes do Cliente + Histórico de Atendimentos
 *
 * @page /clientes/[id]
 * @description Perfil completo do cliente com histórico CRM
 * Conforme FLUXO_CADASTROS_CLIENTE.md e P4.4
 */

'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useAppointments } from '@/hooks/use-appointments';
import { useCustomer, useCustomerWithHistory } from '@/hooks/use-customers';
import { cn } from '@/lib/utils';
import { formatPhone, getInitials } from '@/services/customer-service';
import { useBreadcrumbs } from '@/store/ui-store';
import type { AppointmentResponse, AppointmentStatus } from '@/types/appointment';
import { STATUS_CONFIG } from '@/types/appointment';
import { GENDER_LABELS, TAG_COLORS } from '@/types/customer';
import {
    ArrowLeft,
    CakeIcon,
    Calendar,
    CalendarDays,
    ClipboardList,
    DollarSign,
    Edit,
    History,
    Mail,
    MapPin,
    Phone,
    Scissors,
    TrendingUp,
    User,
} from 'lucide-react';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import { useEffect, useMemo, useState } from 'react';

// =============================================================================
// UTILITÁRIOS
// =============================================================================

function formatCurrency(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}

function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatBirthday(dateStr: string | undefined): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr + 'T00:00:00');
  return date.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: 'long',
  });
}

function getStatusBadge(status: AppointmentStatus) {
  const config = STATUS_CONFIG[status];
  return (
    <Badge className={cn(config.bgColor, config.textColor, 'border-0')}>
      {config.label}
    </Badge>
  );
}

function getRelativeTime(dateStr: string | undefined): string {
  if (!dateStr) return 'Nunca';
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) return 'Hoje';
  if (diffDays === 1) return 'Ontem';
  if (diffDays < 7) return `${diffDays} dias atrás`;
  if (diffDays < 30) return `${Math.floor(diffDays / 7)} semanas atrás`;
  if (diffDays < 365) return `${Math.floor(diffDays / 30)} meses atrás`;
  return `${Math.floor(diffDays / 365)} anos atrás`;
}

// =============================================================================
// COMPONENTE: Card de Estatísticas
// =============================================================================

interface StatsCardProps {
  icon: React.ReactNode;
  label: string;
  value: string | number;
  description?: string;
  className?: string;
}

function StatsCard({ icon, label, value, description, className }: StatsCardProps) {
  return (
    <Card className={className}>
      <CardContent className="p-4">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-muted">{icon}</div>
          <div className="flex-1">
            <p className="text-sm text-muted-foreground">{label}</p>
            <p className="text-xl font-bold">{value}</p>
            {description && (
              <p className="text-xs text-muted-foreground mt-0.5">{description}</p>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// COMPONENTE: Timeline de Atendimentos
// =============================================================================

interface HistoryTimelineProps {
  appointments: AppointmentResponse[];
  isLoading: boolean;
}

function HistoryTimeline({ appointments, isLoading }: HistoryTimelineProps) {
  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="flex gap-4">
            <Skeleton className="h-12 w-12 rounded-full" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-4 w-1/3" />
              <Skeleton className="h-4 w-2/3" />
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (appointments.length === 0) {
    return (
      <div className="text-center py-12">
        <History className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
        <h3 className="text-lg font-semibold mb-2">Sem histórico de atendimentos</h3>
        <p className="text-muted-foreground max-w-md mx-auto">
          Este cliente ainda não possui nenhum atendimento registrado.
        </p>
      </div>
    );
  }

  // Agrupar por mês/ano
  const groupedByMonth = appointments.reduce((acc, apt) => {
    const date = new Date(apt.start_time);
    const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(apt);
    return acc;
  }, {} as Record<string, AppointmentResponse[]>);

  const sortedMonths = Object.keys(groupedByMonth).sort().reverse();

  return (
    <div className="space-y-8">
      {sortedMonths.map((monthKey) => {
        const [year, month] = monthKey.split('-');
        const monthName = new Date(parseInt(year), parseInt(month) - 1).toLocaleDateString('pt-BR', {
          month: 'long',
          year: 'numeric',
        });

        return (
          <div key={monthKey}>
            <h4 className="text-sm font-medium text-muted-foreground mb-4 capitalize">
              {monthName}
            </h4>
            <div className="space-y-4">
              {groupedByMonth[monthKey].map((apt) => (
                <div
                  key={apt.id}
                  className="flex gap-4 p-4 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
                >
                  {/* Ícone/Data */}
                  <div className="flex flex-col items-center">
                    <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                      <Scissors className="h-5 w-5 text-primary" />
                    </div>
                    <div className="text-xs text-muted-foreground mt-1 text-center">
                      {new Date(apt.start_time).toLocaleDateString('pt-BR', {
                        day: '2-digit',
                        month: 'short',
                      })}
                    </div>
                  </div>

                  {/* Conteúdo */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <p className="font-medium">
                          {apt.services.map((s) => s.service_name).join(', ')}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          com {apt.professional_name}
                        </p>
                      </div>
                      {getStatusBadge(apt.status)}
                    </div>
                    <div className="flex items-center gap-4 mt-2 text-sm text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Calendar className="h-3.5 w-3.5" />
                        {formatDateTime(apt.start_time)}
                      </span>
                      <span className="flex items-center gap-1">
                        <DollarSign className="h-3.5 w-3.5" />
                        {formatCurrency(apt.total_price)}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}

// =============================================================================
// COMPONENTE: Tabela de Atendimentos
// =============================================================================

interface HistoryTableProps {
  appointments: AppointmentResponse[];
  isLoading: boolean;
}

function HistoryTable({ appointments, isLoading }: HistoryTableProps) {
  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-6">
          <Skeleton className="h-64 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (appointments.length === 0) {
    return (
      <Card>
        <CardContent className="py-12 text-center">
          <ClipboardList className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-semibold mb-2">Nenhum atendimento encontrado</h3>
          <p className="text-muted-foreground">
            Os atendimentos do cliente aparecerão aqui.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Data/Hora</TableHead>
              <TableHead>Serviços</TableHead>
              <TableHead>Profissional</TableHead>
              <TableHead className="text-right">Valor</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {appointments.map((apt) => (
              <TableRow key={apt.id}>
                <TableCell>
                  <div>
                    <p className="font-medium">{formatDate(apt.start_time)}</p>
                    <p className="text-xs text-muted-foreground">
                      {new Date(apt.start_time).toLocaleTimeString('pt-BR', {
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </p>
                  </div>
                </TableCell>
                <TableCell>
                  <div className="flex flex-wrap gap-1">
                    {apt.services.map((s, idx) => (
                      <Badge key={idx} variant="outline" className="text-xs">
                        {s.service_name}
                      </Badge>
                    ))}
                  </div>
                </TableCell>
                <TableCell>{apt.professional_name}</TableCell>
                <TableCell className="text-right font-medium">
                  {formatCurrency(apt.total_price)}
                </TableCell>
                <TableCell>{getStatusBadge(apt.status)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

// =============================================================================
// PÁGINA PRINCIPAL
// =============================================================================

export default function ClienteDetalhePage() {
  const params = useParams();
  const { setBreadcrumbs } = useBreadcrumbs();
  const clienteId = params.id as string;

  const [activeTab, setActiveTab] = useState<'timeline' | 'table'>('timeline');

  // Queries
  const { data: customer, isLoading: isLoadingCustomer } = useCustomer(clienteId);
  const { data: customerHistory } = useCustomerWithHistory(clienteId);
  
  // Buscar agendamentos do cliente (todos os status)
  const { data: appointmentsData, isLoading: isLoadingAppointments } = useAppointments({
    customer_id: clienteId,
    page_size: 100,
  });

  const appointments = useMemo(() => {
    if (!appointmentsData?.data) return [];
    // Ordenar por data decrescente
    return [...appointmentsData.data].sort(
      (a, b) => new Date(b.start_time).getTime() - new Date(a.start_time).getTime()
    );
  }, [appointmentsData]);

  // Métricas calculadas
  const metrics = useMemo(() => {
    if (!appointments.length) {
      return {
        concluidos: 0,
        cancelados: 0,
        noShow: 0,
        taxaConclusao: '0%',
      };
    }

    const concluidos = appointments.filter((a) => a.status === 'DONE').length;
    const cancelados = appointments.filter((a) => a.status === 'CANCELED').length;
    const noShow = appointments.filter((a) => a.status === 'NO_SHOW').length;
    const total = appointments.length;
    const taxaConclusao = total > 0 ? ((concluidos / total) * 100).toFixed(0) + '%' : '0%';

    return { concluidos, cancelados, noShow, taxaConclusao };
  }, [appointments]);

  // Breadcrumbs
  useEffect(() => {
    if (customer) {
      setBreadcrumbs([
        { label: 'Clientes', href: '/clientes' },
        { label: customer.nome },
      ]);
    }
  }, [customer, setBreadcrumbs]);

  // Loading state
  if (isLoadingCustomer) {
    return (
      <div className="container mx-auto py-6 space-y-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-6 md:grid-cols-3">
          <Skeleton className="h-32" />
          <Skeleton className="h-32 md:col-span-2" />
        </div>
        <Skeleton className="h-96" />
      </div>
    );
  }

  // Not found
  if (!customer) {
    return (
      <div className="container mx-auto py-6">
        <Card>
          <CardContent className="py-16 text-center">
            <User className="h-16 w-16 mx-auto mb-4 text-muted-foreground opacity-50" />
            <h3 className="text-lg font-semibold mb-2">Cliente não encontrado</h3>
            <p className="text-muted-foreground mb-4">
              O cliente solicitado não existe ou foi removido.
            </p>
            <Button asChild>
              <Link href="/clientes">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Voltar para Clientes
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/clientes">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-bold">{customer.nome}</h1>
            <p className="text-muted-foreground">
              Cliente desde {formatDate(customer.criado_em)}
            </p>
          </div>
        </div>
        <Button variant="outline" asChild>
          <Link href={`/clientes?edit=${clienteId}`}>
            <Edit className="h-4 w-4 mr-2" />
            Editar
          </Link>
        </Button>
      </div>

      {/* Cards de Informações e Métricas */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Informações do Cliente */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base">Informações</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Avatar e Nome */}
            <div className="flex items-center gap-3">
              <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-xl font-bold text-primary">
                {getInitials(customer.nome)}
              </div>
              <div>
                <p className="font-medium">{customer.nome}</p>
                {customer.genero && (
                  <Badge variant="outline" className="text-xs mt-1">
                    {GENDER_LABELS[customer.genero]}
                  </Badge>
                )}
              </div>
            </div>

            <Separator />

            {/* Contato */}
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-sm">
                <Phone className="h-4 w-4 text-muted-foreground" />
                <span>{formatPhone(customer.telefone)}</span>
              </div>
              {customer.email && (
                <div className="flex items-center gap-2 text-sm">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="truncate">{customer.email}</span>
                </div>
              )}
              {customer.data_nascimento && (
                <div className="flex items-center gap-2 text-sm">
                  <CakeIcon className="h-4 w-4 text-muted-foreground" />
                  <span>{formatBirthday(customer.data_nascimento)}</span>
                </div>
              )}
            </div>

            {/* Endereço */}
            {customer.endereco_logradouro && (
              <>
                <Separator />
                <div className="flex items-start gap-2 text-sm">
                  <MapPin className="h-4 w-4 text-muted-foreground mt-0.5" />
                  <div>
                    <p>
                      {customer.endereco_logradouro}
                      {customer.endereco_numero && `, ${customer.endereco_numero}`}
                    </p>
                    {customer.endereco_bairro && (
                      <p className="text-muted-foreground">{customer.endereco_bairro}</p>
                    )}
                    {customer.endereco_cidade && (
                      <p className="text-muted-foreground">
                        {customer.endereco_cidade}
                        {customer.endereco_estado && ` - ${customer.endereco_estado}`}
                      </p>
                    )}
                  </div>
                </div>
              </>
            )}

            {/* Tags */}
            {customer.tags && customer.tags.length > 0 && (
              <>
                <Separator />
                <div className="flex flex-wrap gap-1">
                  {customer.tags.map((tag) => (
                    <Badge
                      key={tag}
                      className={cn(
                        TAG_COLORS[tag] || 'bg-gray-100 text-gray-800'
                      )}
                    >
                      {tag}
                    </Badge>
                  ))}
                </div>
              </>
            )}

            {/* Observações */}
            {customer.observacoes && (
              <>
                <Separator />
                <div>
                  <p className="text-xs text-muted-foreground mb-1">Observações</p>
                  <p className="text-sm">{customer.observacoes}</p>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        {/* Métricas */}
        <div className="lg:col-span-2 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard
            icon={<ClipboardList className="h-5 w-5 text-primary" />}
            label="Total Atendimentos"
            value={customerHistory?.total_atendimentos || appointments.length}
            description={`${metrics.concluidos} concluídos`}
          />
          <StatsCard
            icon={<DollarSign className="h-5 w-5 text-green-600" />}
            label="Total Gasto"
            value={customerHistory?.total_gasto ? formatCurrency(customerHistory.total_gasto) : 'R$ 0,00'}
            description={`Ticket médio: ${customerHistory?.ticket_medio ? formatCurrency(customerHistory.ticket_medio) : '-'}`}
          />
          <StatsCard
            icon={<CalendarDays className="h-5 w-5 text-blue-600" />}
            label="Último Atendimento"
            value={getRelativeTime(customerHistory?.ultimo_atendimento || undefined)}
            description={customerHistory?.ultimo_atendimento ? formatDate(customerHistory.ultimo_atendimento) : undefined}
          />
          <StatsCard
            icon={<TrendingUp className="h-5 w-5 text-amber-600" />}
            label="Taxa de Conclusão"
            value={metrics.taxaConclusao}
            description={`${metrics.cancelados} cancelados, ${metrics.noShow} não compareceu`}
          />
        </div>
      </div>

      {/* Histórico de Atendimentos */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <History className="h-5 w-5" />
                Histórico de Atendimentos
              </CardTitle>
              <CardDescription>
                Todos os atendimentos realizados para este cliente
              </CardDescription>
            </div>
            <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'timeline' | 'table')}>
              <TabsList>
                <TabsTrigger value="timeline">Timeline</TabsTrigger>
                <TabsTrigger value="table">Tabela</TabsTrigger>
              </TabsList>
            </Tabs>
          </div>
        </CardHeader>
        <CardContent>
          {activeTab === 'timeline' ? (
            <HistoryTimeline
              appointments={appointments}
              isLoading={isLoadingAppointments}
            />
          ) : (
            <HistoryTable
              appointments={appointments}
              isLoading={isLoadingAppointments}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
