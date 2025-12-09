# Especificação da Comanda — Estilo Trinks | NEXO v1.0

**Versão:** 2.0
**Última Atualização:** 30/11/2025
**Status:** ✅ Implementado (UI Refatorada)
**Referência:** Print do modal de fechamento de conta do Trinks

---

## 📋 Visão Geral

Este documento especifica o design e comportamento do modal de fechamento de conta (comanda) baseado na interface do Trinks. O objetivo é criar uma experiência de pagamento completa, intuitiva e profissional.

---

## 🎨 Layout do Modal

### Estrutura Geral

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│ [X]                    FECHAMENTO DE CONTA DO DIA 27/11/2025                    │
├───────────────────────────────────────┬─────────────────────────────────────────┤
│                                       │                                         │
│  ┌────────────────────────────────┐   │  ┌─────────────────────────────────┐   │
│  │ 👤 DADOS DO CLIENTE            │   │  │ 💰 FORMAS DE PAGAMENTO          │   │
│  │    Nome: João Silva            │   │  │                                 │   │
│  │    Tel: (11) 99999-9999        │   │  │  > Crédito                      │   │
│  │    Comanda: #1234              │   │  │  > Débito                       │   │
│  │    Cliente desde: 2023         │   │  │  > Transação Bancária (PIX)    │   │
│  │    Pontos: 150 ⭐               │   │  │  > À Vista / Dinheiro          │   │
│  │    [Editar] [Adicionar Pontos] │   │  │  > Outros                       │   │
│  └────────────────────────────────┘   │  │  > Pré-Pago                     │   │
│                                       │  │                                 │   │
│  ┌────────────────────────────────┐   │  └─────────────────────────────────┘   │
│  │ 📑 ABAS DE ITENS               │   │                                         │
│  │ [Serviços] Produtos Pacotes    │   │  ┌─────────────────────────────────┐   │
│  │            Vales    Cupom      │   │  │ 📊 RESUMO                       │   │
│  └────────────────────────────────┘   │  │                                 │   │
│                                       │  │  Total:     R$ 80,00            │   │
│  ┌────────────────────────────────┐   │  │  Recebido:  R$ 0,00             │   │
│  │ 📋 TABELA DE ITENS             │   │  │  Falta:     R$ 80,00            │   │
│  │ ─────────────────────────────  │   │  │                                 │   │
│  │ Item      Preço  Desc.  Pagar  │   │  │  [ ] Deixar como dívida         │   │
│  │ ─────────────────────────────  │   │  │  [ ] Troco como gorjeta         │   │
│  │ ✂️ Corte   R$50   -      R$50  │   │  │                                 │   │
│  │ 🧔 Barba   R$30   -      R$30  │   │  └─────────────────────────────────┘   │
│  │ ─────────────────────────────  │   │                                         │
│  │           Total:       R$80,00 │   │                                         │
│  └────────────────────────────────┘   │                                         │
│                                       │                                         │
├───────────────────────────────────────┴─────────────────────────────────────────┤
│  🔗 Produtos usados nos serviços              [Cancelar]  [🟠 Fechar Conta]     │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 📐 Especificação Detalhada

### 1. Header do Modal

```tsx
<DialogHeader className="bg-muted/50 border-b px-6 py-4">
  <div className="flex items-center justify-between">
    <DialogTitle className="text-lg font-semibold">
      FECHAMENTO DE CONTA DO DIA {format(date, 'dd/MM/yyyy')}
    </DialogTitle>
    <DialogClose asChild>
      <Button variant="ghost" size="icon">
        <X className="h-4 w-4" />
      </Button>
    </DialogClose>
  </div>
</DialogHeader>
```

---

### 2. Card de Dados do Cliente

#### Visual
```tsx
<Card className="mb-4">
  <CardContent className="p-4">
    <div className="flex items-start gap-4">
      {/* Avatar */}
      <Avatar className="h-16 w-16">
        <AvatarImage src={customer.avatar} />
        <AvatarFallback>{customer.initials}</AvatarFallback>
      </Avatar>
      
      {/* Informações */}
      <div className="flex-1">
        <h3 className="font-semibold text-lg">{customer.name}</h3>
        <p className="text-muted-foreground">{customer.phone}</p>
        
        <div className="flex items-center gap-4 mt-2 text-sm">
          <Badge variant="outline">Comanda: #{command.number}</Badge>
          <span>Cliente desde: {customer.memberSince}</span>
        </div>
        
        {/* Pontos de Fidelidade */}
        <div className="flex items-center gap-2 mt-2">
          <Star className="h-4 w-4 text-yellow-500" />
          <span className="text-sm">{customer.points} pontos</span>
          <Button variant="link" size="sm" className="p-0 h-auto">
            Carregar pontos / Ver itens
          </Button>
        </div>
      </div>
      
      {/* Ações */}
      <div className="flex flex-col gap-2">
        <Button variant="outline" size="sm">
          <Pencil className="h-4 w-4 mr-2" />
          Editar cliente
        </Button>
        <Button variant="outline" size="sm">
          <Plus className="h-4 w-4 mr-2" />
          Adicionar pontos
        </Button>
      </div>
    </div>
  </CardContent>
</Card>
```

#### Dados
| Campo | Tipo | Descrição |
|-------|------|-----------|
| `customer.name` | string | Nome completo do cliente |
| `customer.phone` | string | Telefone formatado |
| `customer.avatar` | string | URL da foto (opcional) |
| `customer.initials` | string | Iniciais para fallback |
| `customer.memberSince` | number | Ano de cadastro |
| `customer.points` | number | Pontos de fidelidade |
| `command.number` | number | Número sequencial da comanda |

---

### 3. Abas de Itens

#### Tabs Disponíveis
```tsx
<Tabs defaultValue="services" className="mb-4">
  <TabsList className="grid w-full grid-cols-5">
    <TabsTrigger value="services">
      <Scissors className="h-4 w-4 mr-2" />
      Serviços
    </TabsTrigger>
    <TabsTrigger value="products" disabled={!hasProducts}>
      <Package className="h-4 w-4 mr-2" />
      Produtos
    </TabsTrigger>
    <TabsTrigger value="packages" disabled>
      <Gift className="h-4 w-4 mr-2" />
      Pacotes
    </TabsTrigger>
    <TabsTrigger value="vouchers" disabled>
      <Ticket className="h-4 w-4 mr-2" />
      Vales
    </TabsTrigger>
    <TabsTrigger value="coupon" disabled>
      <Tag className="h-4 w-4 mr-2" />
      Cupom
    </TabsTrigger>
  </TabsList>
  
  <TabsContent value="services">
    <ItemsTable items={services} />
  </TabsContent>
  
  <TabsContent value="products">
    <ItemsTable items={products} />
  </TabsContent>
</Tabs>
```

#### Status das Abas
| Aba | MVP v1.0 | Futuro |
|-----|----------|--------|
| Serviços | ✅ Implementar | - |
| Produtos | 🚫 Desabilitada | v1.2.0 |
| Pacotes | 🚫 Desabilitada | v2.0.0 |
| Vales | 🚫 Desabilitada | v2.0.0 |
| Cupom | 🚫 Desabilitada | v2.0.0 |

---

### 4. Tabela de Itens

#### Estrutura da Tabela
```tsx
<Table>
  <TableHeader>
    <TableRow>
      <TableHead className="w-[40%]">Item</TableHead>
      <TableHead className="text-right">Preço</TableHead>
      <TableHead className="text-right">Desconto</TableHead>
      <TableHead className="text-right">A pagar</TableHead>
      <TableHead className="w-[40px]"></TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    {items.map((item) => (
      <TableRow key={item.id}>
        <TableCell>
          <div className="flex items-center gap-2">
            {item.type === 'SERVICE' ? (
              <Scissors className="h-4 w-4 text-muted-foreground" />
            ) : (
              <Package className="h-4 w-4 text-muted-foreground" />
            )}
            <span>{item.name}</span>
          </div>
        </TableCell>
        <TableCell className="text-right">
          <Input
            type="text"
            value={formatCurrency(item.price)}
            onChange={(e) => handlePriceChange(item.id, e.target.value)}
            className="w-24 text-right"
          />
        </TableCell>
        <TableCell className="text-right">
          <Input
            type="text"
            value={item.discount ? formatCurrency(item.discount) : '-'}
            onChange={(e) => handleDiscountChange(item.id, e.target.value)}
            className="w-24 text-right"
          />
        </TableCell>
        <TableCell className="text-right font-medium">
          {formatCurrency(item.price - (item.discount || 0))}
        </TableCell>
        <TableCell>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => handleEdit(item)}>
                <Pencil className="h-4 w-4 mr-2" />
                Editar
              </DropdownMenuItem>
              <DropdownMenuItem 
                onClick={() => handleRemove(item)}
                className="text-destructive"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Remover
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </TableCell>
      </TableRow>
    ))}
  </TableBody>
  <TableFooter>
    <TableRow>
      <TableCell colSpan={3}>Total</TableCell>
      <TableCell className="text-right font-bold">
        {formatCurrency(total)}
      </TableCell>
      <TableCell />
    </TableRow>
  </TableFooter>
</Table>
```

#### Funcionalidades
- [x] Preço editável inline
- [x] Desconto editável inline
- [x] Cálculo automático de "A pagar"
- [x] Menu de ações (editar/remover)
- [x] Total atualizado em tempo real

---

### 5. Painel de Pagamento (Lateral Direita)

#### Formas de Pagamento (Accordion)
```tsx
<Accordion type="single" collapsible className="w-full">
  {paymentMethods.map((method) => (
    <AccordionItem key={method.id} value={method.id}>
      <AccordionTrigger className="hover:no-underline">
        <div className="flex items-center gap-2">
          {method.icon}
          <span>{method.name}</span>
        </div>
      </AccordionTrigger>
      <AccordionContent>
        <div className="space-y-2 pt-2">
          {method.options.map((option) => (
            <Button
              key={option.id}
              variant={selectedPayments.includes(option.id) ? "default" : "outline"}
              className="w-full justify-start"
              onClick={() => handleSelectPayment(option)}
            >
              {option.icon}
              <span className="ml-2">{option.name}</span>
            </Button>
          ))}
        </div>
      </AccordionContent>
    </AccordionItem>
  ))}
</Accordion>
```

#### Formas de Pagamento Disponíveis

| Categoria | Opções | Ícone |
|-----------|--------|-------|
| **Crédito** | Visa, Master, Elo, Amex, Hipercard | `<CreditCard />` |
| **Débito** | Visa, Master, Elo | `<CreditCard />` |
| **Transação Bancária** | PIX, TED, DOC | `<Building />` |
| **À Vista** | Dinheiro | `<Banknote />` |
| **Outros** | Voucher, Vale Funcionário, Permuta | `<Wallet />` |
| **Pré-Pago** | Pacote, Crédito Antecipado | `<Receipt />` |

---

### 6. Resumo Financeiro

```tsx
<Card className="mt-4">
  <CardContent className="p-4">
    <div className="space-y-2">
      <div className="flex justify-between">
        <span className="text-muted-foreground">Total</span>
        <span className="font-medium">{formatCurrency(total)}</span>
      </div>
      <div className="flex justify-between">
        <span className="text-muted-foreground">Recebido</span>
        <span className="font-medium text-green-600">
          {formatCurrency(received)}
        </span>
      </div>
      <Separator />
      <div className="flex justify-between text-lg">
        <span className="font-semibold">
          {remaining > 0 ? 'Falta' : 'Troco'}
        </span>
        <span className={cn(
          "font-bold",
          remaining > 0 ? "text-red-600" : "text-green-600"
        )}>
          {formatCurrency(Math.abs(remaining))}
        </span>
      </div>
    </div>
    
    {/* Opções */}
    <div className="mt-4 space-y-2">
      <div className="flex items-center space-x-2">
        <Checkbox
          id="debt"
          checked={leaveAsDebt}
          onCheckedChange={setLeaveAsDebt}
          disabled={remaining <= 0}
        />
        <label htmlFor="debt" className="text-sm">
          Deixar o que falta como dívida
        </label>
      </div>
      
      <div className="flex items-center space-x-2">
        <Checkbox
          id="tip"
          checked={changeAsTip}
          onCheckedChange={setChangeAsTip}
          disabled={remaining >= 0}
        />
        <label htmlFor="tip" className="text-sm">
          Deixar troco como gorjeta
        </label>
      </div>
    </div>
  </CardContent>
</Card>
```

#### Lógica de Cálculo
```typescript
interface PaymentSummary {
  subtotal: number;      // Soma dos itens
  discount: number;      // Desconto total
  total: number;         // subtotal - discount
  received: number;      // Soma dos pagamentos
  remaining: number;     // total - received (+ = falta, - = troco)
  tip: number;           // Gorjeta (se houver)
}

function calculateSummary(items: CommandItem[], payments: Payment[]): PaymentSummary {
  const subtotal = items.reduce((sum, item) => sum + item.price, 0);
  const discount = items.reduce((sum, item) => sum + (item.discount || 0), 0);
  const total = subtotal - discount;
  const received = payments.reduce((sum, p) => sum + p.amount, 0);
  const remaining = total - received;
  
  return { subtotal, discount, total, received, remaining, tip: 0 };
}
```

---

### 7. Rodapé do Modal

```tsx
<DialogFooter className="border-t px-6 py-4">
  <div className="flex items-center justify-between w-full">
    <Button variant="link" className="text-muted-foreground">
      <Package className="h-4 w-4 mr-2" />
      Produtos usados nos serviços
    </Button>
    
    <div className="flex gap-2">
      <Button variant="outline" onClick={onCancel}>
        Cancelar
      </Button>
      <Button 
        onClick={onClose}
        disabled={remaining > 0 && !leaveAsDebt}
        className="bg-orange-500 hover:bg-orange-600"
      >
        Fechar Conta
      </Button>
    </div>
  </div>
</DialogFooter>
```

---

## 🔄 Fluxo de Interação

### 1. Abrir Modal
```
Usuário clica em "Finalizar" no agendamento
   ↓
Sistema carrega comanda vinculada
   ↓
Modal abre com dados pré-populados
```

### 2. Editar Itens
```
Usuário clica em valor de preço/desconto
   ↓
Campo fica editável
   ↓
Usuário digita novo valor
   ↓
Total recalcula automaticamente
```

### 3. Selecionar Pagamento
```
Usuário expande accordion de forma de pagamento
   ↓
Seleciona opção específica (ex: PIX)
   ↓
Modal de valor abre
   ↓
Usuário informa valor
   ↓
Pagamento adicionado à lista
   ↓
Resumo atualiza "Recebido"
```

### 4. Fechar Conta
```
"Falta" zerado OU checkbox "dívida" marcado
   ↓
Botão "Fechar Conta" habilitado
   ↓
Usuário clica
   ↓
Confirmação
   ↓
POST /api/commands/:id/close
   ↓
POST /api/payments (para cada pagamento)
   ↓
PUT /api/appointments/:id/status → DONE
   ↓
Modal fecha
   ↓
Toast de sucesso
```

---

## 📦 Componentes Necessários

### Novos Componentes
| Componente | Descrição | Prioridade |
|------------|-----------|------------|
| `CommandModal` | Modal principal | Alta |
| `CustomerCard` | Card de dados do cliente | Alta |
| `ItemsTable` | Tabela de serviços/produtos | Alta |
| `PaymentAccordion` | Accordion de formas de pagamento | Alta |
| `PaymentSummary` | Resumo financeiro | Alta |
| `PaymentValueDialog` | Dialog para informar valor | Alta |

### Dependências (já instaladas)
- `@radix-ui/react-dialog`
- `@radix-ui/react-accordion`
- `@radix-ui/react-checkbox`
- `@radix-ui/react-tabs`
- `@radix-ui/react-dropdown-menu`

---

## 🗄️ Estrutura de Dados

### Command DTO
```typescript
interface CommandResponse {
  id: string;
  command_number: number;
  appointment_id: string | null;
  customer: CustomerSummary;
  professional: ProfessionalSummary;
  items: CommandItem[];
  payments: CommandPayment[];
  status: 'OPEN' | 'AWAITING_PAYMENT' | 'PAID' | 'CANCELED';
  subtotal: string;
  discount: string;
  tip: string;
  total: string;
  opened_at: string;
  closed_at: string | null;
}

interface CommandItem {
  id: string;
  type: 'SERVICE' | 'PRODUCT';
  name: string;
  price: string;
  discount: string | null;
  quantity: number;
  total: string;
}

interface CommandPayment {
  id: string;
  method: PaymentMethod;
  method_detail: string; // "Visa", "PIX", etc.
  amount: string;
  paid_at: string;
}
```

---

## ✅ Checklist de Implementação

### Backend ✅ Completo
- [x] Migração: Criar tabela `commands`
- [x] Migração: Criar tabela `command_items`
- [x] Migração: Criar tabela `command_payments`
- [x] Entity: `Command` com métodos de negócio
- [x] Repository: CRUD + queries específicas
- [x] Use Case: `CreateCommandUseCase`
- [x] Use Case: `AddItemToCommandUseCase`
- [x] Use Case: `RemoveItemFromCommandUseCase`
- [x] Use Case: `AddPaymentToCommandUseCase`
- [x] Use Case: `CloseCommandUseCase`
- [x] Handler: 6 endpoints
- [x] Testes: Unit + Integration

### Frontend ✅ Completo (v2.0 - Layout Trinks)
- [x] Componente: `CommandModal` (refatorado 30/11/2025)
- [x] Layout em 2 colunas responsivo
- [x] Card de cliente com avatar e dados completos
- [x] Tabs de itens (Serviços/Produtos/Pacotes/Vales/Cupom)
- [x] Tabela de itens editável (via `CommandItemsTable`)
- [x] Seletor de formas de pagamento (via `PaymentMethodSelector`)
- [x] Lista de pagamentos selecionados (via `SelectedPaymentsList`)
- [x] Resumo financeiro (via `PaymentSummary`)
- [x] Checkboxes de dívida e gorjeta
- [x] Observações finais
- [x] Botão laranja "🟠 Fechar Conta" (estilo Trinks)
- [x] ScrollArea para longo conteúdo
- [x] Header com data formatada
- [x] Footer com ações
- [x] Hook: `useCommand`
- [x] Hook: `useAddPayment`
- [x] Hook: `useCloseCommand`
- [x] Service: `command-service.ts`
- [x] Types: `command.ts`

---

**Referência Visual:** Print do Trinks fornecido pelo usuário
**Próximo Passo:** Implementar backend (tabelas + endpoints)
**Estimativa:** 3-4 dias de desenvolvimento
