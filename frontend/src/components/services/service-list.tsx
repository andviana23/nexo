'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
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
import { MoreHorizontal, Pencil, Power, Trash } from 'lucide-react';
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
      <div className="space-y-2">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-10 w-full" />
      </div>
    );
  }

  const services = data?.servicos || [];

  return (
    <>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Nome</TableHead>
              <TableHead>Categoria</TableHead>
              <TableHead>Preço</TableHead>
              <TableHead>Duração</TableHead>
              <TableHead>Comissão</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-[70px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {services.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="h-24 text-center">
                  Nenhum serviço encontrado.
                </TableCell>
              </TableRow>
            ) : (
              services.map((service) => (
                <TableRow key={service.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center gap-2">
                      {service.cor && (
                        <div
                          className="h-3 w-3 rounded-full"
                          style={{ backgroundColor: service.cor }}
                        />
                      )}
                      {service.nome}
                    </div>
                  </TableCell>
                  <TableCell>
                    {service.categoria_nome ? (
                      <Badge variant="outline" style={{ 
                        borderColor: service.categoria_cor, 
                        color: service.categoria_cor 
                      }}>
                        {service.categoria_nome}
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground">-</span>
                    )}
                  </TableCell>
                  <TableCell>R$ {service.preco}</TableCell>
                  <TableCell>{service.duracao_formatada}</TableCell>
                  <TableCell>{service.comissao}%</TableCell>
                  <TableCell>
                    <Badge variant={service.ativo ? 'default' : 'secondary'}>
                      {service.ativo ? 'Ativo' : 'Inativo'}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                          <span className="sr-only">Abrir menu</span>
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Ações</DropdownMenuLabel>
                        <DropdownMenuItem onClick={() => setEditingService(service)}>
                          <Pencil className="mr-2 h-4 w-4" />
                          Editar
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => toggleStatus.mutate(service.id)}>
                          <Power className="mr-2 h-4 w-4" />
                          {service.ativo ? 'Desativar' : 'Ativar'}
                        </DropdownMenuItem>
                        <DropdownMenuItem 
                          className="text-destructive focus:text-destructive"
                          onClick={() => {
                            if (confirm('Tem certeza que deseja excluir este serviço?')) {
                              deleteService.mutate(service.id);
                            }
                          }}
                        >
                          <Trash className="mr-2 h-4 w-4" />
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

      <ServiceFormModal
        open={!!editingService}
        onOpenChange={(open) => !open && setEditingService(null)}
        serviceToEdit={editingService}
      />
    </>
  );
}
