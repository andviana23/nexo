# NEXO — Components

> Guia de componentes do NEXO. Baseado em **shadcn/ui** (Radix UI + Tailwind).
> Componentes ficam em `src/components/ui`.

---

## 1. Componentes Base (shadcn/ui)

Estes componentes são a base de toda a interface. Eles são acessíveis, responsivos e seguem o tema.

### Button

Botão padrão do sistema.

```tsx
import { Button } from "@/components/ui/button"

<Button variant="default">Salvar</Button>
<Button variant="outline">Cancelar</Button>
<Button variant="ghost">Editar</Button>
<Button variant="destructive">Excluir</Button>
```

### Input & Label

Campos de texto básicos.

```tsx
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

<div className="grid w-full max-w-sm items-center gap-1.5">
  <Label htmlFor="email">Email</Label>
  <Input type="email" id="email" placeholder="Email" />
</div>;
```

### Card

Container padrão para conteúdo agrupado.

```tsx
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

<Card>
  <CardHeader>
    <CardTitle>Título do Card</CardTitle>
    <CardDescription>Descrição opcional.</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Conteúdo principal.</p>
  </CardContent>
  <CardFooter>
    <p>Rodapé.</p>
  </CardFooter>
</Card>;
```

### Dialog (Modal)

Janelas modais para confirmações ou formulários rápidos.

```tsx
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

<Dialog>
  <DialogTrigger>Abrir</DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Você tem certeza?</DialogTitle>
    </DialogHeader>
    {/* Conteúdo */}
  </DialogContent>
</Dialog>;
```

### Table

Tabelas de dados.

```tsx
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Nome</TableHead>
      <TableHead>Status</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow>
      <TableCell>João</TableCell>
      <TableCell>Ativo</TableCell>
    </TableRow>
  </TableBody>
</Table>;
```

---

## 2. Componentes Compostos

Componentes que combinam múltiplos primitivos para funcionalidades específicas.

### DataTable

Uma tabela avançada com ordenação, filtro e paginação (usando `@tanstack/react-table`).

- Local: `src/components/ui/data-table.tsx` (se implementado)
- Uso: Listagens principais (Clientes, Agendamentos).

### Form Wrappers

Componentes que integram `react-hook-form` com os componentes de UI.

```tsx
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';

<FormField
  control={form.control}
  name="username"
  render={({ field }) => (
    <FormItem>
      <FormLabel>Username</FormLabel>
      <FormControl>
        <Input placeholder="shadcn" {...field} />
      </FormControl>
      <FormMessage />
    </FormItem>
  )}
/>;
```

---

## 3. Ícones (Lucide React)

Usamos a biblioteca `lucide-react` para ícones. Eles são leves e consistentes.

```tsx
import { Calendar, CreditCard, User } from "lucide-react"

<Calendar className="mr-2 h-4 w-4" />
<Button>
  <User className="mr-2 h-4 w-4" /> Perfil
</Button>
```

---

## 4. Customização

Para customizar um componente base (ex: mudar o raio da borda do Button):

1. Vá até `src/components/ui/button.tsx`.
2. Edite as classes Tailwind na definição `cva`.

**Não** crie um novo componente `MyButton` se for apenas uma variação de estilo. Use as variantes do componente existente ou adicione uma nova variante no arquivo original.

---

## 5. Toast (Notificações)

Usamos o componente `Sonner` (via shadcn) para notificações.

```tsx
import { toast } from 'sonner';

toast('Evento criado com sucesso', {
  description: 'Domingo, 03 de Dezembro às 9h',
  action: {
    label: 'Desfazer',
    onClick: () => console.log('Undo'),
  },
});
```
