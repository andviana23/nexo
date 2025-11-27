'use client';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * Página de Detalhes do Profissional
 *
 * @page /profissionais/[id]
 * @description Exibe informações detalhadas de um profissional específico
 */

import { format } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import {
    ArrowLeftIcon,
    BanknoteIcon,
    CalendarIcon,
    EditIcon,
    Loader2Icon,
    MailIcon,
    MoreVerticalIcon,
    PhoneIcon,
    ScissorsIcon,
    StarIcon,
    UserCircleIcon,
    UserIcon,
    XCircleIcon,
} from 'lucide-react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { useCallback, useEffect, useState } from 'react';

import {
    ProfessionalModal,
    ProfessionalStatusBadge,
    ProfessionalTypeBadge,
} from '@/components/professionals';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Separator } from '@/components/ui/separator';
import {
    useProfessional,
    useUpdateProfessionalStatus,
} from '@/hooks/use-professionals';
import { maskCPF, maskPhone } from '@/lib/validations/professional';
import { useBreadcrumbs } from '@/store/ui-store';
import type {
    ProfessionalModalState,
    ProfessionalStatus,
} from '@/types/professional';

// =============================================================================
// HELPERS
// =============================================================================

function getInitials(name: string): string {
  return name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase();
}

// =============================================================================
// PAGE COMPONENT
// =============================================================================

export default function ProfessionalDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;
  const { setBreadcrumbs } = useBreadcrumbs();

  // Estado do modal de edição
  const [modalState, setModalState] = useState<ProfessionalModalState>({
    isOpen: false,
    mode: 'edit',
  });

  // Queries e Mutations
  const { data: professional, isLoading, isError } = useProfessional(id);
  const updateStatus = useUpdateProfessionalStatus();

  // Breadcrumbs
  useEffect(() => {
    setBreadcrumbs([
      { label: 'Dashboard', href: '/' },
      { label: 'Profissionais', href: '/profissionais' },
      { label: professional?.nome || 'Detalhes' },
    ]);
  }, [setBreadcrumbs, professional?.nome]);

  // Handlers de status
  const handleStatusChange = useCallback(
    (newStatus: ProfessionalStatus) => {
      updateStatus.mutate(
        { id, data: { status: newStatus } },
        {
          onSuccess: () => {
            // Recarrega dados
          },
        }
      );
    },
    [id, updateStatus]
  );

  const handleEdit = useCallback(() => {
    if (professional) {
      setModalState({
        isOpen: true,
        mode: 'edit',
        professional,
      });
    }
  }, [professional]);

  const handleCloseModal = useCallback(() => {
    setModalState((prev) => ({ ...prev, isOpen: false }));
  }, []);

  // Loading state
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2Icon className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  // Error state
  if (isError || !professional) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] gap-4">
        <XCircleIcon className="size-12 text-destructive" />
        <h2 className="text-lg font-semibold">Profissional não encontrado</h2>
        <Button variant="outline" asChild>
          <Link href="/profissionais">
            <ArrowLeftIcon className="size-4 mr-2" />
            Voltar
          </Link>
        </Button>
      </div>
    );
  }

  const showCommissionInfo =
    professional.tipo === 'BARBEIRO' ||
    (professional.comissao !== null && professional.comissao !== undefined);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/profissionais">
              <ArrowLeftIcon className="size-5" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">
              {professional.nome}
            </h1>
            <div className="flex items-center gap-2 mt-1">
              <ProfessionalTypeBadge type={professional.tipo} />
              <ProfessionalStatusBadge status={professional.status} />
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* Botão Editar */}
          <Button variant="outline" onClick={handleEdit}>
            <EditIcon className="size-4 mr-2" />
            Editar
          </Button>

          {/* Menu de Ações */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="icon">
                <MoreVerticalIcon className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {/* Ações de Status */}
              {professional.status === 'ATIVO' && (
                <>
                  <DropdownMenuItem
                    onClick={() => handleStatusChange('AFASTADO')}
                  >
                    Afastar Temporariamente
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => handleStatusChange('INATIVO')}
                    className="text-orange-600"
                  >
                    Desativar
                  </DropdownMenuItem>
                </>
              )}

              {professional.status === 'INATIVO' && (
                <DropdownMenuItem
                  onClick={() => handleStatusChange('ATIVO')}
                  className="text-green-600"
                >
                  Reativar
                </DropdownMenuItem>
              )}

              {professional.status === 'AFASTADO' && (
                <DropdownMenuItem
                  onClick={() => handleStatusChange('ATIVO')}
                  className="text-green-600"
                >
                  Retornar do Afastamento
                </DropdownMenuItem>
              )}

              {professional.status !== 'DEMITIDO' && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => {
                      if (
                        confirm(
                          'Tem certeza que deseja desligar este profissional?'
                        )
                      ) {
                        handleStatusChange('DEMITIDO');
                        router.push('/profissionais');
                      }
                    }}
                    className="text-destructive"
                  >
                    Desligar Profissional
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Grid de Cards */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Card de Perfil */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <UserCircleIcon className="size-5 text-primary" />
              Perfil
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <Avatar className="size-20">
                <AvatarImage
                  src={professional.foto || undefined}
                  alt={professional.nome}
                />
                <AvatarFallback className="text-xl">
                  {getInitials(professional.nome)}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-1">
                <p className="text-lg font-semibold">{professional.nome}</p>
                <p className="text-sm text-muted-foreground">
                  {professional.cpf.length === 14 ? 'CNPJ' : 'CPF'}: {maskCPF(professional.cpf)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Card de Contato */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <UserIcon className="size-5 text-primary" />
              Contato
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-3">
              <MailIcon className="size-4 text-muted-foreground" />
              <a
                href={`mailto:${professional.email}`}
                className="text-sm hover:underline"
              >
                {professional.email}
              </a>
            </div>
            <div className="flex items-center gap-3">
              <PhoneIcon className="size-4 text-muted-foreground" />
              <a
                href={`tel:${professional.telefone}`}
                className="text-sm hover:underline"
              >
                {maskPhone(professional.telefone)}
              </a>
            </div>
          </CardContent>
        </Card>

        {/* Card de Dados Profissionais */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <CalendarIcon className="size-5 text-primary" />
              Dados Profissionais
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Tipo</p>
              <p className="font-medium mt-1">
                <ProfessionalTypeBadge type={professional.tipo} />
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Data de Admissão</p>
              <p className="font-medium">
                {format(
                  new Date(professional.data_admissao),
                  "dd 'de' MMMM 'de' yyyy",
                  { locale: ptBR }
                )}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Status</p>
              <p className="font-medium mt-1">
                <ProfessionalStatusBadge status={professional.status} />
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Card de Comissão */}
        {showCommissionInfo && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-lg">
                <BanknoteIcon className="size-5 text-primary" />
                Comissão
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground">
                    Tipo de Comissão
                  </p>
                  <p className="font-medium">
                    {professional.tipo_comissao === 'PERCENTUAL'
                      ? 'Percentual'
                      : 'Valor Fixo'}
                  </p>
                </div>
                <div className="flex gap-8">
                  <div>
                    <p className="text-sm text-muted-foreground">Serviços</p>
                    <p className="text-2xl font-bold text-primary">
                      {professional.comissao ?? 0}%
                    </p>
                  </div>
                  {professional.comissao_produtos !== null &&
                    professional.comissao_produtos !== undefined && (
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Produtos
                        </p>
                        <p className="text-2xl font-bold text-green-600">
                          {professional.comissao_produtos}%
                        </p>
                      </div>
                    )}
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>

      {/* Especialidades */}
      {professional.especialidades && professional.especialidades.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <StarIcon className="size-5 text-primary" />
              Especialidades
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {professional.especialidades.map((especialidade) => (
                <Badge key={especialidade} variant="secondary">
                  <ScissorsIcon className="size-3 mr-1" />
                  {especialidade}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Histórico / Metadados */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Informações do Sistema</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-2 text-sm text-muted-foreground">
            <div>
              <p>
                Cadastrado em:{' '}
                {format(
                  new Date(professional.criado_em),
                  "dd/MM/yyyy 'às' HH:mm",
                  { locale: ptBR }
                )}
              </p>
            </div>
            <div>
              <p>
                Última atualização:{' '}
                {format(
                  new Date(professional.atualizado_em),
                  "dd/MM/yyyy 'às' HH:mm",
                  { locale: ptBR }
                )}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Seções Futuras */}
      <Separator />
      <div className="grid gap-6 md:grid-cols-2 opacity-50">
        {/* Estatísticas - Placeholder */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg text-muted-foreground">
              📊 Estatísticas (Em breve)
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Agendamentos realizados, ticket médio, avaliações...
            </p>
          </CardContent>
        </Card>

        {/* Agenda - Placeholder */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg text-muted-foreground">
              📅 Agenda (Em breve)
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Próximos agendamentos e horários de trabalho...
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Modal de Edição */}
      <ProfessionalModal state={modalState} onClose={handleCloseModal} />
    </div>
  );
}
