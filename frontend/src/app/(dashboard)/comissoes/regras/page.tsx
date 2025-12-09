'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from '@/components/ui/dialog';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Textarea } from '@/components/ui/textarea';
import {
    useCommissionRules,
    useCreateCommissionRule,
    useDeleteCommissionRule,
    useUpdateCommissionRule,
} from '@/hooks/use-commissions';
import { cn } from '@/lib/utils';
import {
    CalculationBase,
    CommissionRule,
    CommissionRuleType,
    CreateCommissionRuleRequest,
} from '@/types/commission';
import {
    MoreHorizontal,
    Pencil,
    Plus,
    ScrollText,
    Trash2,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

// =============================================================================
// HELPERS
// =============================================================================

const formatCurrency = (value: string | number) => {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return 'R$ 0,00';
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(num);
};

const formatPercentage = (value: string | number) => {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(num)) return '0%';
  return `${num}%`;
};

// =============================================================================
// FORM DIALOG
// =============================================================================

interface RuleFormDialogProps {
  rule?: CommissionRule;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function RuleFormDialog({ rule, open, onOpenChange }: RuleFormDialogProps) {
  const createMutation = useCreateCommissionRule();
  const updateMutation = useUpdateCommissionRule();
  
  const [formData, setFormData] = useState<CreateCommissionRuleRequest>({
    name: rule?.name || '',
    description: rule?.description || '',
    type: (rule?.type as CommissionRuleType) || CommissionRuleType.PERCENTUAL,
    default_rate: rule?.default_rate || '',
    min_amount: rule?.min_amount || '',
    max_amount: rule?.max_amount || '',
    calculation_base: (rule?.calculation_base as CalculationBase) || CalculationBase.BRUTO,
    priority: rule?.priority || 0,
  });

  const isEditing = !!rule;
  const isLoading = createMutation.isPending || updateMutation.isPending;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name || !formData.default_rate) {
      toast.error('Preencha os campos obrigatórios');
      return;
    }

    try {
      if (isEditing) {
        await updateMutation.mutateAsync({
          id: rule.id,
          data: formData,
        });
      } else {
        await createMutation.mutateAsync(formData);
      }
      onOpenChange(false);
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>
              {isEditing ? 'Editar Regra de Comissão' : 'Nova Regra de Comissão'}
            </DialogTitle>
            <DialogDescription>
              {isEditing
                ? 'Altere os dados da regra de comissão.'
                : 'Crie uma nova regra para cálculo de comissões.'}
            </DialogDescription>
          </DialogHeader>
          
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Nome *</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Ex: Comissão Padrão"
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Descrição</Label>
              <Textarea
                id="description"
                value={formData.description || ''}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Descrição da regra..."
                rows={2}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="type">Tipo *</Label>
                <Select
                  value={formData.type}
                  onValueChange={(value) => 
                    setFormData({ ...formData, type: value as CommissionRuleType })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value={CommissionRuleType.PERCENTUAL}>
                      Percentual (%)
                    </SelectItem>
                    <SelectItem value={CommissionRuleType.FIXO}>
                      Valor Fixo (R$)
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="default_rate">
                  {formData.type === CommissionRuleType.PERCENTUAL ? 'Taxa (%)' : 'Valor (R$)'} *
                </Label>
                <Input
                  id="default_rate"
                  type="number"
                  step="0.01"
                  min="0"
                  value={formData.default_rate}
                  onChange={(e) => setFormData({ ...formData, default_rate: e.target.value })}
                  placeholder={formData.type === CommissionRuleType.PERCENTUAL ? '30' : '50.00'}
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="min_amount">Valor Mínimo (R$)</Label>
                <Input
                  id="min_amount"
                  type="number"
                  step="0.01"
                  min="0"
                  value={formData.min_amount || ''}
                  onChange={(e) => setFormData({ ...formData, min_amount: e.target.value })}
                  placeholder="0.00"
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="max_amount">Valor Máximo (R$)</Label>
                <Input
                  id="max_amount"
                  type="number"
                  step="0.01"
                  min="0"
                  value={formData.max_amount || ''}
                  onChange={(e) => setFormData({ ...formData, max_amount: e.target.value })}
                  placeholder="Sem limite"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="calculation_base">Base de Cálculo</Label>
                <Select
                  value={formData.calculation_base}
                  onValueChange={(value) => 
                    setFormData({ ...formData, calculation_base: value as CalculationBase })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value={CalculationBase.BRUTO}>Valor Bruto</SelectItem>
                    <SelectItem value={CalculationBase.LIQUIDO}>Valor Líquido</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="priority">Prioridade</Label>
                <Input
                  id="priority"
                  type="number"
                  min="0"
                  value={formData.priority || 0}
                  onChange={(e) => setFormData({ ...formData, priority: parseInt(e.target.value) || 0 })}
                  placeholder="0"
                />
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? 'Salvando...' : isEditing ? 'Salvar' : 'Criar'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// DELETE DIALOG
// =============================================================================

interface DeleteDialogProps {
  rule: CommissionRule;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function DeleteDialog({ rule, open, onOpenChange }: DeleteDialogProps) {
  const deleteMutation = useDeleteCommissionRule();

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(rule.id);
      onOpenChange(false);
    } catch {
      // Error handled by mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Confirmar Exclusão</DialogTitle>
          <DialogDescription>
            Tem certeza que deseja excluir a regra &quot;{rule.name}&quot;? Esta ação não pode ser desfeita.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? 'Excluindo...' : 'Excluir'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// MAIN PAGE
// =============================================================================

export default function RegrasComissaoPage() {
  const [showOnlyActive, setShowOnlyActive] = useState(false);
  const { data, isLoading } = useCommissionRules({ active_only: showOnlyActive });
  
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [editingRule, setEditingRule] = useState<CommissionRule | undefined>();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deletingRule, setDeletingRule] = useState<CommissionRule | undefined>();

  const updateMutation = useUpdateCommissionRule();

  const handleEdit = (rule: CommissionRule) => {
    setEditingRule(rule);
    setFormDialogOpen(true);
  };

  const handleDelete = (rule: CommissionRule) => {
    setDeletingRule(rule);
    setDeleteDialogOpen(true);
  };

  const handleToggleActive = async (rule: CommissionRule) => {
    try {
      await updateMutation.mutateAsync({
        id: rule.id,
        data: { is_active: !rule.is_active },
      });
    } catch {
      // Error handled by mutation
    }
  };

  const handleCloseFormDialog = (open: boolean) => {
    setFormDialogOpen(open);
    if (!open) {
      setEditingRule(undefined);
    }
  };

  const handleCloseDeleteDialog = (open: boolean) => {
    setDeleteDialogOpen(open);
    if (!open) {
      setDeletingRule(undefined);
    }
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold">Regras de Comissão</h2>
          <p className="text-sm text-muted-foreground">
            Configure as regras para cálculo automático de comissões
          </p>
        </div>
        <Button onClick={() => setFormDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Regra
        </Button>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center gap-2">
            <Switch
              id="active-only"
              checked={showOnlyActive}
              onCheckedChange={setShowOnlyActive}
            />
            <Label htmlFor="active-only">Mostrar apenas ativas</Label>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="pt-6">
          {isLoading ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Nome</TableHead>
                  <TableHead>Tipo</TableHead>
                  <TableHead>Taxa/Valor</TableHead>
                  <TableHead>Base</TableHead>
                  <TableHead>Prioridade</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-8" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-8" /></TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : data?.data && data.data.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Nome</TableHead>
                  <TableHead>Tipo</TableHead>
                  <TableHead>Taxa/Valor</TableHead>
                  <TableHead>Base</TableHead>
                  <TableHead>Prioridade</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((rule) => (
                  <TableRow key={rule.id} className={cn(!rule.is_active && 'opacity-60')}>
                    <TableCell>
                      <div>
                        <p className="font-medium">{rule.name}</p>
                        {rule.description && (
                          <p className="text-sm text-muted-foreground truncate max-w-[200px]">
                            {rule.description}
                          </p>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {rule.type === CommissionRuleType.PERCENTUAL ? 'Percentual' : 'Fixo'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {rule.type === CommissionRuleType.PERCENTUAL
                        ? formatPercentage(rule.default_rate)
                        : formatCurrency(rule.default_rate)}
                    </TableCell>
                    <TableCell>
                      {rule.calculation_base === CalculationBase.BRUTO ? 'Bruto' : 'Líquido'}
                    </TableCell>
                    <TableCell>{rule.priority || 0}</TableCell>
                    <TableCell>
                      <Badge variant={rule.is_active ? 'default' : 'secondary'}>
                        {rule.is_active ? 'Ativa' : 'Inativa'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => handleEdit(rule)}>
                            <Pencil className="mr-2 h-4 w-4" />
                            Editar
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleToggleActive(rule)}>
                            {rule.is_active ? 'Desativar' : 'Ativar'}
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDelete(rule)}
                            className="text-destructive"
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            Excluir
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <ScrollText className="h-12 w-12 text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium">Nenhuma regra encontrada</h3>
              <p className="text-sm text-muted-foreground mb-4">
                Crie uma regra de comissão para começar
              </p>
              <Button onClick={() => setFormDialogOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Nova Regra
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Dialogs */}
      <RuleFormDialog
        rule={editingRule}
        open={formDialogOpen}
        onOpenChange={handleCloseFormDialog}
      />

      {deletingRule && (
        <DeleteDialog
          rule={deletingRule}
          open={deleteDialogOpen}
          onOpenChange={handleCloseDeleteDialog}
        />
      )}
    </div>
  );
}
