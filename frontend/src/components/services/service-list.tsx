'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useDeleteService, useServices, useToggleServiceStatus } from '@/hooks/useServices';
import { Service } from '@/types/service';
import { Clock, DollarSign, Edit, MoreHorizontal, Power, Scissors, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { ServiceFormModal } from './service-form-modal';

interface ServiceListProps {
  search: string;
}

export function ServiceList({ search }: ServiceListProps) {
  const { data, isLoading } = useServices({ search });
  const deleteService = useDeleteService();
  const toggleStatus = useToggleServiceStatus();

  const [editingService, setEditingService] = useState<Service | null>(null);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-48 mb-2" />
          <Skeleton className="h-4 w-96" />
        </CardHeader>
        <CardContent className="p-0">
          <div className="space-y-4 p-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const services = data?.servicos || [];

  return (
    <>
      <Card className="border-border shadow-sm">
        <CardHeader className="bg-muted/40 pb-4">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg">Catálogo de Serviços</CardTitle>
              <CardDescription>
                Gerencie os serviços oferecidos, preços e comissões.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50 hover:bg-muted/50">
                  <TableHead className="w-[300px] pl-6">Serviço</TableHead>
                  <TableHead>Categoria</TableHead>
                  <TableHead>Preço</TableHead>
                  <TableHead>Duração</TableHead>
                  <TableHead>Comissão</TableHead>
                  <TableHead className="text-center">Status</TableHead>
                  <TableHead className="w-[80px] text-right pr-6">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {services.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-64 text-center">
                      <div className="flex flex-col items-center justify-center gap-3 text-muted-foreground">
                        <div className="h-12 w-12 rounded-full bg-muted flex items-center justify-center">
                          <Scissors className="h-6 w-6" />
                        </div>
                        <div className="space-y-1">
                          <p className="font-medium text-foreground">Nenhum serviço encontrado</p>
                          <p className="text-sm">Cadastre serviços para começar.</p>
                        </div>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  services.map((service) => (
                    <TableRow key={service.id} className="group hover:bg-muted/30">
                      <TableCell className="pl-6">
                        <div className="flex items-center gap-3">
                          <div
                            className="h-9 w-9 rounded-lg flex items-center justify-center shadow-sm"
                            style={{ backgroundColor: service.cor || '#e5e7eb' }}
                          >
                            <Scissors className="h-4 w-4 text-white opacity-90" />
                          </div>
                          <div>
                            <p className="font-medium">{service.nome}</p>
                            {service.descricao && (
                              <p className="text-xs text-muted-foreground line-clamp-1 max-w-[200px]">{service.descricao}</p>
                            )}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        {service.categoria_nome ? (
                          <Badge variant="outline" className="font-normal">
                            {service.categoria_nome}
                          </Badge>
                        ) : (
                          <span className="text-muted-foreground text-xs">-</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center text-sm font-medium">
                          <DollarSign className="h-3 w-3 mr-1 text-muted-foreground" />
                          {Number(service.preco).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center text-sm text-muted-foreground">
                          <Clock className="h-3 w-3 mr-1" />
                          {service.duracao_formatada}
                        </div>
                      </TableCell>
                      <TableCell>
                        <span className="font-mono text-xs">{service.comissao}%</span>
                      </TableCell>
                      <TableCell className="text-center">
                        <Badge
                          variant={service.ativo ? 'default' : 'secondary'}
                          className={service.ativo ? 'bg-emerald-600 hover:bg-emerald-600/80' : ''}
                        >
                          {service.ativo ? 'Ativo' : 'Inativo'}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right pr-6">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity">
                              <span className="sr-only">Abrir menu</span>
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => setEditingService(service)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Editar
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => toggleStatus.mutate(service.id)}>
                              <Power className="mr-2 h-4 w-4" />
                              {service.ativo ? 'Desativar' : 'Ativar'}
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                              className="text-destructive focus:text-destructive"
                              onClick={() => {
                                if (confirm('Tem certeza que deseja excluir este serviço?')) {
                                  deleteService.mutate(service.id);
                                }
                              }}
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Excluir
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <ServiceFormModal
        open={!!editingService}
        onOpenChange={(open) => !open && setEditingService(null)}
        serviceToEdit={editingService}
      />
    </>
  );
}
