# Fluxo de Status de Agendamento — NEXO v1.0

**Versão:** 1.0.0  
**Data de Criação:** 01/12/2025  
**Status:** ✅ Implementado  
**Responsável:** Product + Tech Lead  
**Módulo:** Agendamento  

---

## 📊 Visão Geral

Este documento especifica o **ciclo de vida completo** de um agendamento, desde a criação até a finalização, incluindo todas as transições de status possíveis, cores do card, ações disponíveis via menu de contexto (botão direito) e regras de negócio.

---

## 🎯 Objetivos

1. Definir o **fluxo de estados** de um agendamento
2. Especificar as **cores visuais** de cada status no calendário
3. Detalhar as **ações do menu de contexto** (botão direito)
4. Estabelecer as **regras de transição** entre status
5. Garantir **rastreabilidade** e controle operacional

---

## 🔄 Diagrama de Estados

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CICLO DE VIDA DO AGENDAMENTO                        │
└─────────────────────────────────────────────────────────────────────────────┘

     ┌──────────────┐
     │   CREATED    │ ◄─── Agendamento criado pela recepção
     │  🟡 Amarelo  │      Card aparece no calendário
     └──────┬───────┘      Tamanho proporcional à duração do serviço
            │
            │ [Confirmar]
            ▼
     ┌──────────────┐
     │  CONFIRMED   │ ◄─── Cliente confirmou presença
     │  🟢 Verde    │      (WhatsApp, telefone, etc)
     └──────┬───────┘
            │
            │ [Check-In]
            ▼
     ┌──────────────┐
     │ CHECKED_IN   │ ◄─── Cliente chegou na barbearia
     │  🔵 Azul     │      Notifica profissional (futuro)
     └──────┬───────┘
            │
            │ [Iniciar Atendimento]
            ▼
     ┌──────────────┐
     │ IN_SERVICE   │ ◄─── Profissional está atendendo
     │  🟣 Roxo     │      Serviços sendo executados
     └──────┬───────┘
            │
            │ [Finalizar Atendimento]
            ▼
     ┌──────────────┐
     │ AWAITING_    │ ◄─── Serviços finalizados
     │  PAYMENT     │      Aguardando pagamento
     │  🟠 Laranja  │      Comanda aberta
     └──────┬───────┘
            │
            │ [Fechar Comanda / Concluir]
            ▼
     ┌──────────────┐
     │     DONE     │ ◄─── Agendamento concluído
     │  ⚪ Cinza    │      Pagamento confirmado
     └──────────────┘

┌──────────────────────── FLUXOS ALTERNATIVOS ─────────────────────────┐
│                                                                       │
│  De CREATED ou CONFIRMED:                                            │
│     [Cancelar] ──► CANCELED (⚫ Cinza Escuro)                        │
│                                                                       │
│  De CONFIRMED ou CHECKED_IN:                                         │
│     [Cliente Faltou] ──► NO_SHOW (🔴 Vermelho)                      │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

---

## 📋 Especificação de Status

### Status 1: CREATED (Criado)

**Cor:** 🟡 **Amarelo Dourado** (`#F59E0B` / `bg-amber-500`)

**Descrição:**  
Agendamento criado pela recepção, mas ainda não confirmado pelo cliente.

**Características Visuais:**
- Card no calendário com altura proporcional à duração do serviço
- Cor de fundo: Amarelo dourado suave
- Borda: Amarelo mais escuro
- Texto: Escuro para contraste

**Ações Disponíveis (Menu Contexto):**
1. ✅ **Confirmar Agendamento** → Status: `CONFIRMED`
2. ✏️ **Editar Agendamento** → Abre modal de edição
3. 📋 **Abrir Comanda** → Abre modal de comanda (preview)
4. ❌ **Cancelar Agendamento** → Status: `CANCELED`

**Regras de Transição:**
- Pode ir para: `CONFIRMED`, `CANCELED`
- Não pode pular etapas

**Notificações:**
- 🔔 Envio de confirmação via WhatsApp (futuro)

---

### Status 2: CONFIRMED (Confirmado)

**Cor:** 🟢 **Verde** (`#10B981` / `bg-green-500`)

**Descrição:**  
Cliente confirmou que comparecerá ao agendamento.

**Características Visuais:**
- Card verde para destacar confirmação
- Indicador visual de "confirmado" (ícone de check)

**Ações Disponíveis (Menu Contexto):**
1. ✅ **Fazer Check-In** → Status: `CHECKED_IN`
2. ✏️ **Editar Agendamento** → Abre modal de edição
3. 📋 **Abrir Comanda** → Abre modal de comanda (preview)
4. 🔴 **Cliente Faltou (No-Show)** → Status: `NO_SHOW`
5. ❌ **Cancelar Agendamento** → Status: `CANCELED`

**Regras de Transição:**
- Pode ir para: `CHECKED_IN`, `NO_SHOW`, `CANCELED`
- Se passou do horário e não fez check-in → sugerir marcar como NO_SHOW

**Notificações:**
- 🔔 Lembrete 1 hora antes (WhatsApp - futuro)
- 🔔 Lembrete 1 dia antes (WhatsApp - futuro)

---

### Status 3: CHECKED_IN (Presente)

**Cor:** 🔵 **Azul** (`#3B82F6` / `bg-blue-500`)

**Descrição:**  
Cliente chegou à barbearia e marcou presença.

**Características Visuais:**
- Card azul indicando presença confirmada
- Indicador de "cliente presente" (ícone de usuário com check)

**Ações Disponíveis (Menu Contexto):**
1. ▶️ **Iniciar Atendimento** → Status: `IN_SERVICE`
2. ✏️ **Editar Agendamento** → Abre modal de edição
3. 📋 **Abrir Comanda** → Abre modal de comanda (preview)
4. 🔴 **Cliente Faltou (No-Show)** → Status: `NO_SHOW` (caso tenha saído)
5. ❌ **Cancelar Agendamento** → Status: `CANCELED`

**Regras de Transição:**
- Pode ir para: `IN_SERVICE`, `NO_SHOW`, `CANCELED`
- Registra timestamp de check-in (`checked_in_at`)

**Notificações:**
- 📱 Notifica profissional no app (futuro): "Cliente [Nome] chegou"

---

### Status 4: IN_SERVICE (Em Atendimento)

**Cor:** 🟣 **Roxo** (`#8B5CF6` / `bg-purple-500`)

**Descrição:**  
Profissional iniciou o atendimento ao cliente.

**Características Visuais:**
- Card roxo indicando atendimento em andamento
- Indicador de "em atendimento" (ícone de tesoura ou relógio)

**Ações Disponíveis (Menu Contexto):**
1. ✅ **Finalizar Atendimento** → Status: `AWAITING_PAYMENT`
2. 📋 **Abrir Comanda** → Abre modal de comanda (adicionar itens)
3. ❌ **Cancelar Agendamento** → Status: `CANCELED` (excepcional)

**Regras de Transição:**
- Pode ir para: `AWAITING_PAYMENT`, `CANCELED`
- Registra timestamp de início (`started_at`)
- Ao finalizar, registra timestamp de fim (`finished_at`)

**Comportamento:**
- Timer de duração visível no card (opcional)
- Bloqueia outros agendamentos no mesmo horário/profissional

---

### Status 5: AWAITING_PAYMENT (Aguardando Pagamento)

**Cor:** 🟠 **Laranja** (`#F97316` / `bg-orange-500`)

**Descrição:**  
Atendimento finalizado, aguardando pagamento.

**Características Visuais:**
- Card laranja indicando pendência financeira
- Indicador de "aguardando pagamento" (ícone de cifrão)

**Ações Disponíveis (Menu Contexto):**
1. 💰 **Fechar Comanda** → Abre modal de comanda (pagamento)
2. ✅ **Concluir (Pago)** → Status: `DONE`
3. ❌ **Cancelar Agendamento** → Status: `CANCELED` (excepcional)

**Regras de Transição:**
- Pode ir para: `DONE`, `CANCELED`
- **Obrigatório:** Ter `command_id` vinculado
- Ao clicar no card → Abre `CommandModal` automaticamente

**Comportamento:**
- Comanda criada automaticamente ao entrar neste status
- Vincular `appointment.command_id` ao ID da comanda

**Modal de Comanda:**
- Layout 2 colunas (estilo Trinks)
- Seguir especificação: [ESPECIFICACAO_COMANDA_TRINKS.md](../../Agendamento/ESPECIFICACAO_COMANDA_TRINKS.md)

---

### Status 6: DONE (Concluído)

**Cor:** ⚪ **Cinza Claro** (`#94A3B8` / `bg-slate-400`)

**Descrição:**  
Agendamento concluído com sucesso. Pagamento confirmado.

**Características Visuais:**
- Card cinza claro indicando finalização
- Opacidade reduzida para não poluir visualmente

**Ações Disponíveis (Menu Contexto):**
1. 👁️ **Visualizar Detalhes** → Modal readonly
2. 📋 **Ver Comanda** → Modal de comanda em modo leitura
3. 🔄 **Reagendar Cliente** → Criar novo agendamento com mesmos dados

**Regras de Transição:**
- **Estado final** (não pode mudar)
- Permanece no calendário por 7 dias (configurável)

**Comportamento:**
- Comanda fechada (`command.status = CLOSED`)
- Registros financeiros criados
- Comissão calculada para o profissional

---

### Status 7: NO_SHOW (Cliente Faltou)

**Cor:** 🔴 **Vermelho** (`#EF4444` / `bg-red-500`)

**Descrição:**  
Cliente não compareceu ao agendamento.

**Características Visuais:**
- Card vermelho indicando ausência
- Indicador de "faltou" (ícone de usuário com X)

**Ações Disponíveis (Menu Contexto):**
1. 👁️ **Visualizar Detalhes** → Modal readonly
2. 🔄 **Reagendar Cliente** → Criar novo agendamento
3. ⚠️ **Registrar Observação** → Adicionar nota sobre a falta

**Regras de Transição:**
- **Estado final** (não pode mudar)
- Incrementa contador de faltas do cliente

**Comportamento:**
- Registrar no histórico do cliente
- Se cliente tem 3+ faltas → sugerir política de confirmação obrigatória
- Slot de horário liberado automaticamente

---

### Status 8: CANCELED (Cancelado)

**Cor:** ⚫ **Cinza Escuro** (`#475569` / `bg-slate-600`)

**Descrição:**  
Agendamento cancelado (por cliente, recepção ou sistema).

**Características Visuais:**
- Card cinza escuro com riscado
- Opacidade reduzida

**Ações Disponíveis (Menu Contexto):**
1. 👁️ **Visualizar Detalhes** → Modal readonly (motivo do cancelamento)
2. 🔄 **Reagendar Cliente** → Criar novo agendamento

**Regras de Transição:**
- **Estado final** (não pode mudar)
- Pode ser cancelado de qualquer status anterior

**Comportamento:**
- Obrigatório informar motivo do cancelamento
- Registrar quem cancelou (`canceled_by`)
- Slot de horário liberado automaticamente
- Notificar cliente (WhatsApp - futuro)

---

## 🎨 Tabela Resumo de Cores

| Status | Cor | Hex | Tailwind Class |
|--------|-----|-----|----------------|
| CREATED | 🟡 Amarelo Dourado | `#F59E0B` | `bg-amber-500` |
| CONFIRMED | 🟢 Verde | `#10B981` | `bg-green-500` |
| CHECKED_IN | 🔵 Azul | `#3B82F6` | `bg-blue-500` |
| IN_SERVICE | 🟣 Roxo | `#8B5CF6` | `bg-purple-500` |
| AWAITING_PAYMENT | 🟠 Laranja | `#F97316` | `bg-orange-500` |
| DONE | ⚪ Cinza Claro | `#94A3B8` | `bg-slate-400` |
| NO_SHOW | 🔴 Vermelho | `#EF4444` | `bg-red-500` |
| CANCELED | ⚫ Cinza Escuro | `#475569` | `bg-slate-600` |

---

## 🖱️ Menu de Contexto (Botão Direito)

### Estrutura do Menu

O menu de contexto é **dinâmico** e mostra apenas ações válidas para o status atual.

```typescript
interface AppointmentContextMenuProps {
  appointment: AppointmentResponse;
  onEdit?: () => void;
  onConfirm?: () => void;
  onCheckIn?: () => void;
  onStartService?: () => void;
  onFinishService?: () => void;
  onCloseCommand?: () => void;
  onComplete?: () => void;
  onNoShow?: () => void;
  onCancel?: () => void;
}
```

### Ações por Status

#### CREATED
```
┌───────────────────────────────┐
│ ✅ Confirmar Agendamento      │ → CONFIRMED
│ ✏️  Editar Agendamento        │ → Abre modal
│ 📋 Abrir Comanda              │ → Abre comanda (preview)
│ ❌ Cancelar Agendamento       │ → CANCELED
└───────────────────────────────┘
```

#### CONFIRMED
```
┌───────────────────────────────┐
│ ✅ Fazer Check-In             │ → CHECKED_IN
│ ✏️  Editar Agendamento        │ → Abre modal
│ 📋 Abrir Comanda              │ → Abre comanda (preview)
│ 🔴 Cliente Faltou             │ → NO_SHOW
│ ❌ Cancelar                   │ → CANCELED
└───────────────────────────────┘
```

#### CHECKED_IN
```
┌───────────────────────────────┐
│ ▶️  Iniciar Atendimento       │ → IN_SERVICE
│ ✏️  Editar Agendamento        │ → Abre modal
│ 📋 Abrir Comanda              │ → Abre comanda (preview)
│ 🔴 Cliente Faltou             │ → NO_SHOW
│ ❌ Cancelar                   │ → CANCELED
└───────────────────────────────┘
```

#### IN_SERVICE
```
┌───────────────────────────────┐
│ ✅ Finalizar Atendimento      │ → AWAITING_PAYMENT
│ 📋 Abrir Comanda              │ → Abre comanda (editar)
│ ❌ Cancelar                   │ → CANCELED
└───────────────────────────────┘
```

#### AWAITING_PAYMENT
```
┌───────────────────────────────┐
│ 💰 Fechar Comanda             │ → Abre modal de pagamento
│ ✅ Concluir (Pago)            │ → DONE
│ ❌ Cancelar                   │ → CANCELED
└───────────────────────────────┘
```

#### DONE, NO_SHOW, CANCELED
```
┌───────────────────────────────┐
│ 👁️  Visualizar Detalhes       │ → Modal readonly
│ 🔄 Reagendar Cliente          │ → Novo agendamento
└───────────────────────────────┘
```

---

## 📐 Componentes Visuais do Card

### Estrutura do Card no Calendário

```tsx
<Card className={cn(
  "relative overflow-hidden rounded-lg border-2 shadow-sm transition-all hover:shadow-md cursor-pointer",
  statusColors[appointment.status]
)}>
  {/* Header */}
  <div className="px-3 py-2 border-b">
    <div className="flex items-center justify-between">
      <span className="text-xs font-medium">
        {format(appointment.start_time, 'HH:mm')} - {format(appointment.end_time, 'HH:mm')}
      </span>
      <Badge variant={statusBadgeVariant[appointment.status]} className="text-xs">
        {statusLabels[appointment.status]}
      </Badge>
    </div>
  </div>

  {/* Body */}
  <div className="p-3 space-y-1">
    <h4 className="font-semibold text-sm truncate">
      {appointment.customer_name}
    </h4>
    <p className="text-xs text-muted-foreground truncate">
      {appointment.service_names.join(', ')}
    </p>
    {appointment.professional_name && (
      <p className="text-xs text-muted-foreground flex items-center gap-1">
        <Scissors className="h-3 w-3" />
        {appointment.professional_name}
      </p>
    )}
  </div>

  {/* Footer (se AWAITING_PAYMENT) */}
  {appointment.status === 'AWAITING_PAYMENT' && (
    <div className="px-3 py-2 bg-orange-100 border-t border-orange-200">
      <Button variant="ghost" size="sm" className="w-full text-orange-700">
        <CreditCard className="h-4 w-4 mr-2" />
        Fechar Comanda
      </Button>
    </div>
  )}

  {/* Indicador de Comanda */}
  {appointment.command_id && (
    <div className="absolute top-2 right-2">
      <Receipt className="h-4 w-4 text-orange-600" />
    </div>
  )}
</Card>
```

### Tamanho Proporcional à Duração

O card deve ter altura proporcional à duração do serviço:

```typescript
function calculateCardHeight(durationMinutes: number): string {
  // Cada 15 minutos = 60px (padrão FullCalendar)
  const pixelsPerMinute = 4; // 60px / 15min
  return `${durationMinutes * pixelsPerMinute}px`;
}

// Exemplo:
// 30 min = 120px
// 45 min = 180px
// 60 min = 240px
```

---

## 🔐 Regras de Negócio

### RN-STATUS-001: Transições Permitidas

Matriz de transições válidas:

| De \ Para | CONFIRMED | CHECKED_IN | IN_SERVICE | AWAITING_PAYMENT | DONE | NO_SHOW | CANCELED |
|-----------|-----------|------------|------------|------------------|------|---------|----------|
| CREATED | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| CONFIRMED | - | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ |
| CHECKED_IN | ❌ | - | ✅ | ❌ | ❌ | ✅ | ✅ |
| IN_SERVICE | ❌ | ❌ | - | ✅ | ❌ | ❌ | ✅ |
| AWAITING_PAYMENT | ❌ | ❌ | ❌ | - | ✅ | ❌ | ✅ |

### RN-STATUS-002: Timestamps Obrigatórios

Cada transição deve registrar timestamp:

```typescript
interface AppointmentTimestamps {
  created_at: string;        // Sempre preenchido
  updated_at: string;        // Atualizado em cada mudança
  checked_in_at?: string;    // Quando status = CHECKED_IN
  started_at?: string;       // Quando status = IN_SERVICE
  finished_at?: string;      // Quando status = AWAITING_PAYMENT
  closed_at?: string;        // Quando status = DONE
  canceled_at?: string;      // Quando status = CANCELED
}
```

### RN-STATUS-003: Comanda Obrigatória

- Status `AWAITING_PAYMENT` **exige** `command_id` preenchido
- Ao transitar para `AWAITING_PAYMENT`, criar comanda automaticamente se não existir
- Vincular `appointment.command_id` ao ID da nova comanda

### RN-STATUS-004: Cancelamento com Motivo

Ao cancelar, obrigatório informar:
- `canceled_reason`: string (motivo do cancelamento)
- `canceled_by`: string (ID do usuário que cancelou)

### RN-STATUS-005: No-Show com Registro

Ao marcar como NO_SHOW:
- Incrementar contador de faltas do cliente (`cliente.no_show_count++`)
- Registrar no histórico do cliente
- Se cliente atingir limite de faltas (ex: 3) → bloquear novos agendamentos até contato

---

## 🔔 Notificações (Futuro)

### Gatilhos de Notificação

| Evento | Destinatário | Canal | Template |
|--------|--------------|-------|----------|
| Status → CONFIRMED | Cliente | WhatsApp | "Agendamento confirmado para {data} às {hora}" |
| 24h antes (CONFIRMED) | Cliente | WhatsApp | "Lembrete: amanhã às {hora} você tem agendamento" |
| 1h antes (CONFIRMED) | Cliente | WhatsApp | "Seu horário está próximo! Nos vemos em 1 hora" |
| Status → CHECKED_IN | Profissional | App/Push | "Cliente {nome} chegou" |
| Status → AWAITING_PAYMENT | Financeiro | Sistema | "Comanda {numero} aguardando pagamento" |
| Status → DONE | Cliente | WhatsApp | "Obrigado pela visita! Até breve 😊" |
| Status → NO_SHOW | Gerente | Sistema | "Cliente {nome} faltou ao agendamento" |
| Status → CANCELED | Cliente | WhatsApp | "Agendamento cancelado. Motivo: {motivo}" |

---

## 🧪 Cenários de Teste

### Teste 1: Fluxo Completo Normal

```gherkin
Cenário: Agendamento com fluxo completo bem-sucedido
  Dado que existe um agendamento com status CREATED
  Quando recepcionista confirma o agendamento
  Então status muda para CONFIRMED
  E card fica verde
  
  Quando cliente chega e faz check-in
  Então status muda para CHECKED_IN
  E card fica azul
  E profissional é notificado
  
  Quando profissional inicia atendimento
  Então status muda para IN_SERVICE
  E card fica roxo
  
  Quando profissional finaliza atendimento
  Então status muda para AWAITING_PAYMENT
  E card fica laranja
  E comanda é criada automaticamente
  
  Quando caixa fecha a comanda
  Então status muda para DONE
  E card fica cinza
  E comissão é calculada
```

### Teste 2: Cliente Faltou

```gherkin
Cenário: Cliente confirmou mas não compareceu
  Dado que existe um agendamento com status CONFIRMED
  E horário do agendamento já passou
  Quando recepcionista marca como "Cliente Faltou"
  Então status muda para NO_SHOW
  E card fica vermelho
  E contador de faltas do cliente incrementa
  E horário é liberado
```

### Teste 3: Cancelamento

```gherkin
Cenário: Cliente cancela agendamento
  Dado que existe um agendamento com status CONFIRMED
  Quando recepcionista cancela com motivo "Cliente desistiu"
  Então sistema solicita motivo do cancelamento
  E status muda para CANCELED
  E card fica cinza escuro
  E horário é liberado
  E cliente é notificado
```

### Teste 4: Abertura de Comanda em AWAITING_PAYMENT

```gherkin
Cenário: Clicar em agendamento aguardando pagamento
  Dado que existe um agendamento com status AWAITING_PAYMENT
  E agendamento tem command_id preenchido
  Quando usuário clica no card do agendamento
  Então CommandModal é aberto automaticamente
  E modal mostra dados da comanda
  E layout segue padrão Trinks (2 colunas)
```

---

## 📊 Métricas e Monitoramento

### KPIs por Status

| Métrica | Cálculo | Meta |
|---------|---------|------|
| **Taxa de Confirmação** | (CONFIRMED / CREATED) × 100 | > 90% |
| **Taxa de Check-In** | (CHECKED_IN / CONFIRMED) × 100 | > 95% |
| **Taxa de No-Show** | (NO_SHOW / CONFIRMED) × 100 | < 10% |
| **Taxa de Conclusão** | (DONE / IN_SERVICE) × 100 | > 98% |
| **Taxa de Cancelamento** | (CANCELED / CREATED) × 100 | < 15% |
| **Tempo Médio em AWAITING_PAYMENT** | AVG(closed_at - finished_at) | < 5 min |

### Alerts

- 🚨 Se taxa de No-Show > 15% → Alerta para gerência
- 🚨 Se agendamento em AWAITING_PAYMENT > 15 min → Alerta para financeiro
- 🚨 Se agendamento em IN_SERVICE > duração prevista + 30min → Alerta para gerente

---

## 🔗 Integrações

### Backend API

Endpoints relacionados ao status:

```typescript
// Confirmar agendamento
POST /api/v1/appointments/:id/confirm
Response: { status: 'CONFIRMED' }

// Check-in
POST /api/v1/appointments/:id/check-in
Response: { status: 'CHECKED_IN', checked_in_at: '...' }

// Iniciar atendimento
POST /api/v1/appointments/:id/start
Response: { status: 'IN_SERVICE', started_at: '...' }

// Finalizar atendimento
POST /api/v1/appointments/:id/finish
Response: { status: 'AWAITING_PAYMENT', command_id: '...' }

// Concluir (marcar como pago)
POST /api/v1/appointments/:id/complete
Response: { status: 'DONE', closed_at: '...' }

// No-Show
POST /api/v1/appointments/:id/no-show
Response: { status: 'NO_SHOW' }

// Cancelar
POST /api/v1/appointments/:id/cancel
Body: { reason: string }
Response: { status: 'CANCELED', canceled_reason: '...' }
```

### Frontend Components

```typescript
// Componentes envolvidos
- AppointmentCard (card visual no calendário)
- AppointmentContextMenu (menu botão direito)
- AppointmentModal (modal de edição)
- CommandModal (modal de comanda)
- StatusBadge (badge de status)
```

---

## ✅ Checklist de Implementação

### Backend ✅ Completo
- [x] Endpoints de transição de status
- [x] Validação de transições permitidas
- [x] Registro de timestamps
- [x] Criação automática de comanda em AWAITING_PAYMENT
- [x] Vinculação appointment ↔ command
- [x] Regras de negócio implementadas
- [x] Testes unitários e integração

### Frontend ✅ Completo
- [x] Cores dos cards por status
- [x] Tamanho proporcional à duração
- [x] Menu de contexto (botão direito)
- [x] Ações dinâmicas por status
- [x] Abertura automática de CommandModal em AWAITING_PAYMENT
- [x] Modal de edição de agendamento
- [x] CommandModal (layout Trinks 2 colunas)
- [x] Badges de status
- [x] Indicadores visuais

### Futuro
- [ ] Notificações WhatsApp
- [ ] Notificação push para profissional
- [ ] Sistema de bloqueio após 3 faltas
- [ ] Dashboard de métricas por status
- [ ] Relatório de no-show por cliente
- [ ] Alerta automático de agendamentos atrasados

---

## 📚 Referências

- [FLUXO_AGENDAMENTO.md](../FLUXO_AGENDAMENTO.md) - Fluxo geral de agendamento
- [ESPECIFICACAO_COMANDA_TRINKS.md](../../Agendamento/ESPECIFICACAO_COMANDA_TRINKS.md) - Spec do modal de comanda
- [PRD_AGENDAMENTO.md](../../Agendamento/PRD_AGENDAMENTO.md) - Product Requirements
- [ARQUITETURA_AGENDAMENTO.md](../../Agendamento/ARQUITETURA_AGENDAMENTO.md) - Arquitetura técnica

---

**Aprovação:**

- [ ] Product Owner
- [ ] Tech Lead
- [ ] UX/UI Designer
- [ ] QA Lead

---

**Última Revisão:** 01/12/2025  
**Próxima Revisão:** 15/12/2025
