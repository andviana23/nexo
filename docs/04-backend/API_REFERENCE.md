# 📡 API Reference v2.0

> **Atualizado em:** 03/12/2025  
> **Versão:** 2.0  
> **Base URL:** `https://api.seudominio.com/api/v1`  
> **Swagger UI:** `/swagger/index.html`

---

## 📋 Índice

1. [Autenticação](#autenticação)
2. [Agendamentos](#agendamentos)
3. [Clientes (CRM)](#clientes-crm)
4. [Profissionais](#profissionais)
5. [Serviços](#serviços)
6. [Categorias de Serviços](#categorias-de-serviços)
7. [Financeiro](#financeiro)
8. [Metas](#metas)
9. [Estoque](#estoque)
10. [Lista da Vez](#lista-da-vez)
11. [Comandas](#comandas)
12. [Precificação](#precificação)
13. [Meios de Pagamento](#meios-de-pagamento)
14. [Caixa Diário](#caixa-diário)
15. [Códigos de Erro](#códigos-de-erro)

---

## 🔐 Autenticação

### POST /api/v1/auth/login
Realiza login e retorna tokens JWT.

```json
// Request
{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}

// Response 200
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
  "refresh_token": "refresh_eyJ0eXAiOiJKV1QiLCJhbGc...",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "usuario@exemplo.com",
    "nome": "João Silva",
    "tenant_id": "tenant-uuid",
    "role": "owner"
  }
}
```

### POST /api/v1/auth/refresh
Renova o access_token usando refresh_token.

### POST /api/v1/auth/logout
Invalida a sessão atual.

### GET /api/v1/auth/me 🔒
Retorna dados do usuário autenticado.

---

## 📅 Agendamentos

| Método | Endpoint | Descrição | RBAC |
|--------|----------|-----------|------|
| `GET` | `/api/v1/appointments` | Listar agendamentos | Todos |
| `POST` | `/api/v1/appointments` | Criar agendamento | Todos |
| `GET` | `/api/v1/appointments/{id}` | Buscar agendamento | Todos |
| `PATCH` | `/api/v1/appointments/{id}/status` | Atualizar status | Admin |
| `PATCH` | `/api/v1/appointments/{id}/reschedule` | Reagendar | Admin |
| `POST` | `/api/v1/appointments/{id}/confirm` | Confirmar | Todos |
| `POST` | `/api/v1/appointments/{id}/cancel` | Cancelar | Admin |
| `POST` | `/api/v1/appointments/{id}/check-in` | Check-in | Todos |
| `POST` | `/api/v1/appointments/{id}/start` | Iniciar atendimento | Todos |
| `POST` | `/api/v1/appointments/{id}/finish` | Finalizar atendimento | Todos |
| `POST` | `/api/v1/appointments/{id}/complete` | Completar | Admin |
| `POST` | `/api/v1/appointments/{id}/no-show` | Marcar no-show | Owner/Manager |

### Criar Agendamento
```json
// POST /api/v1/appointments
{
  "customer_id": "uuid",
  "professional_id": "uuid",
  "start_time": "2024-12-03T10:00:00Z",
  "services": [
    {"service_id": "uuid", "price": "50.00", "duration": 30}
  ],
  "notes": "Cliente preferencial"
}

// Response 201
{
  "id": "uuid",
  "customer_id": "uuid",
  "customer_name": "João Silva",
  "professional_id": "uuid",
  "professional_name": "Carlos Barbeiro",
  "start_time": "2024-12-03T10:00:00Z",
  "end_time": "2024-12-03T10:30:00Z",
  "duration": 30,
  "total_price": "50.00",
  "status": "CREATED",
  "status_display": "Agendado",
  "status_color": "#3B82F6",
  "services": [...],
  "created_at": "2024-12-03T09:00:00Z"
}
```

### Status do Agendamento
| Status | Display | Cor | Descrição |
|--------|---------|-----|-----------|
| `CREATED` | Agendado | #3B82F6 | Recém criado |
| `CONFIRMED` | Confirmado | #10B981 | Cliente confirmou |
| `CHECKED_IN` | Chegou | #8B5CF6 | Cliente presente |
| `IN_SERVICE` | Em Atendimento | #F59E0B | Serviço em andamento |
| `DONE` | Finalizado | #22C55E | Serviço concluído |
| `NO_SHOW` | Não Compareceu | #EF4444 | Cliente faltou |
| `CANCELED` | Cancelado | #6B7280 | Agendamento cancelado |

---

## 👥 Clientes (CRM)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/customers` | Listar clientes |
| `POST` | `/api/v1/customers` | Criar cliente |
| `GET` | `/api/v1/customers/search` | Buscar clientes |
| `GET` | `/api/v1/customers/stats` | Estatísticas |
| `GET` | `/api/v1/customers/check-phone` | Verificar telefone |
| `GET` | `/api/v1/customers/check-cpf` | Verificar CPF |
| `GET` | `/api/v1/customers/{id}` | Buscar cliente |
| `GET` | `/api/v1/customers/{id}/history` | Histórico de atendimentos |
| `GET` | `/api/v1/customers/{id}/export` | Exportar dados (LGPD) |
| `PUT` | `/api/v1/customers/{id}` | Atualizar cliente |
| `DELETE` | `/api/v1/customers/{id}` | Inativar cliente |

### Criar Cliente
```json
// POST /api/v1/customers
{
  "nome": "João Silva",
  "telefone": "+5511999887766",
  "email": "joao@email.com",
  "cpf": "12345678901",
  "data_nascimento": "1990-05-15",
  "genero": "M",
  "tags": ["VIP"],
  "observacoes": "Cliente preferencial"
}
```

---

## 👨‍💼 Profissionais

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/professionals` | Listar profissionais |
| `POST` | `/api/v1/professionals` | Criar profissional |
| `GET` | `/api/v1/professionals/check-email` | Verificar email |
| `GET` | `/api/v1/professionals/check-cpf` | Verificar CPF |
| `GET` | `/api/v1/professionals/{id}` | Buscar profissional |
| `PUT` | `/api/v1/professionals/{id}` | Atualizar profissional |
| `PUT` | `/api/v1/professionals/{id}/status` | Atualizar status |
| `DELETE` | `/api/v1/professionals/{id}` | Remover profissional |

### Tipos de Profissional
- `BARBEIRO` - Realiza cortes e barbas
- `MANICURE` - Serviços de unha
- `RECEPCIONISTA` - Atendimento
- `GERENTE` - Gestão
- `OUTRO` - Outros

---

## ✂️ Serviços

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/servicos` | Listar serviços |
| `POST` | `/api/v1/servicos` | Criar serviço |
| `GET` | `/api/v1/servicos/stats` | Estatísticas |
| `GET` | `/api/v1/servicos/{id}` | Buscar serviço |
| `PUT` | `/api/v1/servicos/{id}` | Atualizar serviço |
| `DELETE` | `/api/v1/servicos/{id}` | Remover serviço |
| `PATCH` | `/api/v1/servicos/{id}/toggle-status` | Ativar/Desativar |

### Criar Serviço
```json
// POST /api/v1/servicos
{
  "nome": "Corte Masculino",
  "descricao": "Corte tradicional",
  "preco": "50.00",
  "duracao_minutos": 30,
  "comissao_percentual": "30.00",
  "categoria_id": "uuid",
  "cor": "#FF5733",
  "ativo": true
}
```

---

## 📁 Categorias de Serviços

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/categorias-servicos` | Listar categorias |
| `POST` | `/api/v1/categorias-servicos` | Criar categoria |
| `GET` | `/api/v1/categorias-servicos/{id}` | Buscar categoria |
| `PUT` | `/api/v1/categorias-servicos/{id}` | Atualizar categoria |
| `DELETE` | `/api/v1/categorias-servicos/{id}` | Remover categoria |

---

## 💰 Financeiro

### Contas a Pagar

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/financial/payables` | Listar contas a pagar |
| `POST` | `/api/v1/financial/payables` | Criar conta a pagar |
| `GET` | `/api/v1/financial/payables/{id}` | Buscar conta |
| `PUT` | `/api/v1/financial/payables/{id}` | Atualizar conta |
| `DELETE` | `/api/v1/financial/payables/{id}` | Remover conta |
| `POST` | `/api/v1/financial/payables/{id}/payment` | Marcar como pago |

### Contas a Receber

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/financial/receivables` | Listar contas a receber |
| `POST` | `/api/v1/financial/receivables` | Criar conta a receber |
| `GET` | `/api/v1/financial/receivables/{id}` | Buscar conta |
| `PUT` | `/api/v1/financial/receivables/{id}` | Atualizar conta |
| `DELETE` | `/api/v1/financial/receivables/{id}` | Remover conta |
| `POST` | `/api/v1/financial/receivables/{id}/receipt` | Marcar como recebido |

### Relatórios

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/financial/dashboard` | Dashboard consolidado |
| `GET` | `/api/v1/financial/projections` | Projeções financeiras |
| `GET` | `/api/v1/financial/cashflow` | Fluxo de caixa |
| `GET` | `/api/v1/financial/cashflow/{id}` | Fluxo específico |
| `GET` | `/api/v1/financial/dre` | Lista DRE |
| `GET` | `/api/v1/financial/dre/{month}` | DRE mensal |

### Despesas Fixas

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/financial/fixed-expenses` | Listar despesas fixas |
| `POST` | `/api/v1/financial/fixed-expenses` | Criar despesa fixa |
| `GET` | `/api/v1/financial/fixed-expenses/{id}` | Buscar despesa |
| `PUT` | `/api/v1/financial/fixed-expenses/{id}` | Atualizar despesa |
| `DELETE` | `/api/v1/financial/fixed-expenses/{id}` | Remover despesa |
| `PATCH` | `/api/v1/financial/fixed-expenses/{id}/toggle` | Ativar/Desativar |
| `GET` | `/api/v1/financial/fixed-expenses/summary` | Resumo mensal |
| `POST` | `/api/v1/financial/fixed-expenses/generate` | Gerar lançamentos |

### Compensações

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/financial/compensations` | Listar compensações |
| `GET` | `/api/v1/financial/compensations/{id}` | Buscar compensação |
| `DELETE` | `/api/v1/financial/compensations/{id}` | Remover compensação |

---

## 🎯 Metas

### Meta Mensal

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/metas/monthly` | Listar metas mensais |
| `POST` | `/api/v1/metas/monthly` | Criar meta mensal |
| `GET` | `/api/v1/metas/monthly/{id}` | Buscar meta |
| `PUT` | `/api/v1/metas/monthly/{id}` | Atualizar meta |
| `DELETE` | `/api/v1/metas/monthly/{id}` | Remover meta |

### Meta por Barbeiro

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/metas/barbers` | Listar metas barbeiros |
| `POST` | `/api/v1/metas/barbers` | Criar meta barbeiro |
| `GET` | `/api/v1/metas/barbers/{id}` | Buscar meta |
| `PUT` | `/api/v1/metas/barbers/{id}` | Atualizar meta |
| `DELETE` | `/api/v1/metas/barbers/{id}` | Remover meta |

### Meta Ticket Médio

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/metas/ticket` | Listar metas ticket |
| `POST` | `/api/v1/metas/ticket` | Criar meta ticket |
| `GET` | `/api/v1/metas/ticket/{id}` | Buscar meta |
| `PUT` | `/api/v1/metas/ticket/{id}` | Atualizar meta |
| `DELETE` | `/api/v1/metas/ticket/{id}` | Remover meta |

---

## 📦 Estoque

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/stock/items` | Listar produtos |
| `GET` | `/api/v1/stock/items/{id}` | Buscar produto |
| `POST` | `/api/v1/stock/products` | Criar produto |
| `POST` | `/api/v1/stock/entries` | Registrar entrada |
| `POST` | `/api/v1/stock/exit` | Registrar saída |
| `POST` | `/api/v1/stock/adjust` | Ajustar estoque |
| `GET` | `/api/v1/stock/alerts` | Alertas de estoque baixo |

### Registrar Entrada
```json
// POST /api/v1/stock/entries
{
  "produto_id": "uuid",
  "quantidade": "10.00",
  "custo_unitario": "25.00",
  "fornecedor_id": "uuid",
  "nota_fiscal": "NF-12345",
  "observacoes": "Reposição mensal"
}
```

### Registrar Saída
```json
// POST /api/v1/stock/exit
{
  "produto_id": "uuid",
  "quantidade": "2.00",
  "motivo": "VENDA",
  "observacoes": "Venda balcão"
}
```

---

## 🔄 Lista da Vez

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/barber-turn/list` | Listar fila |
| `POST` | `/api/v1/barber-turn/add` | Adicionar à fila |
| `POST` | `/api/v1/barber-turn/record` | Registrar atendimento |
| `PUT` | `/api/v1/barber-turn/{professional_id}/toggle-status` | Pausar/Ativar |
| `DELETE` | `/api/v1/barber-turn/{professional_id}` | Remover da fila |
| `POST` | `/api/v1/barber-turn/reset` | Reset mensal |
| `GET` | `/api/v1/barber-turn/history` | Histórico |
| `GET` | `/api/v1/barber-turn/history/summary` | Resumo histórico |
| `GET` | `/api/v1/barber-turn/available` | Barbeiros disponíveis |

---

## 📝 Comandas

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/commands` | Criar comanda |
| `GET` | `/api/v1/commands/{id}` | Buscar comanda |
| `GET` | `/api/v1/commands/by-appointment/{appointmentId}` | Por agendamento |
| `POST` | `/api/v1/commands/{id}/items` | Adicionar item |
| `DELETE` | `/api/v1/commands/{id}/items/{itemId}` | Remover item |
| `POST` | `/api/v1/commands/{id}/payments` | Adicionar pagamento |
| `DELETE` | `/api/v1/commands/{id}/payments/{paymentId}` | Remover pagamento |
| `POST` | `/api/v1/commands/{id}/close` | Fechar comanda |

---

## 💲 Precificação

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/pricing/config` | Obter configuração |
| `PUT` | `/api/v1/pricing/config` | Atualizar configuração |
| `POST` | `/api/v1/pricing/simulate` | Simular preço |
| `GET` | `/api/v1/pricing/simulations` | Listar simulações |
| `GET` | `/api/v1/pricing/simulations/{id}` | Buscar simulação |
| `DELETE` | `/api/v1/pricing/simulations/{id}` | Remover simulação |

---

## 💳 Meios de Pagamento

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/meios-pagamento` | Listar meios |
| `POST` | `/api/v1/meios-pagamento` | Criar meio |
| `GET` | `/api/v1/meios-pagamento/{id}` | Buscar meio |
| `PUT` | `/api/v1/meios-pagamento/{id}` | Atualizar meio |
| `DELETE` | `/api/v1/meios-pagamento/{id}` | Remover meio |
| `PATCH` | `/api/v1/meios-pagamento/{id}/toggle` | Ativar/Desativar |

---

## 💵 Caixa Diário

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/caixa/status` | Status do caixa |
| `POST` | `/api/v1/caixa/abrir` | Abrir caixa |
| `POST` | `/api/v1/caixa/fechar` | Fechar caixa |
| `GET` | `/api/v1/caixa/atual` | Caixa atual |
| `GET` | `/api/v1/caixa/historico` | Histórico |
| `POST` | `/api/v1/caixa/sangria` | Registrar sangria |
| `POST` | `/api/v1/caixa/reforco` | Registrar reforço |
| `POST` | `/api/v1/caixa/operacao` | Registrar operação |

---

## 🚫 Horários Bloqueados

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/blocked-times` | Listar bloqueios |
| `POST` | `/api/v1/blocked-times` | Criar bloqueio |
| `DELETE` | `/api/v1/blocked-times/{id}` | Remover bloqueio |

---

## ⚠️ Códigos de Erro

| Código | HTTP | Descrição |
|--------|------|-----------|
| `BAD_REQUEST` | 400 | Dados inválidos |
| `UNAUTHORIZED` | 401 | Token inválido/ausente |
| `FORBIDDEN` | 403 | Sem permissão |
| `NOT_FOUND` | 404 | Recurso não encontrado |
| `CONFLICT` | 409 | Conflito (duplicado/horário) |
| `INTERNAL_ERROR` | 500 | Erro interno |

### Formato de Erro
```json
{
  "code": "BAD_REQUEST",
  "message": "Descrição do erro",
  "errors": [
    {"field": "email", "message": "Email inválido"}
  ]
}
```

---

## 🔒 Headers Obrigatórios

```http
Authorization: Bearer {access_token}
Content-Type: application/json
```

---

## 📊 Paginação

Endpoints que retornam listas suportam paginação:

```http
GET /api/v1/customers?page=1&page_size=20

// Response
{
  "items": [...],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

---

## 🏢 Multi-Tenant

Todas as operações são automaticamente filtradas pelo `tenant_id` do usuário autenticado.
- Não é possível acessar dados de outros tenants
- O `tenant_id` é extraído do JWT, nunca do payload

---

## 📖 Swagger UI

Documentação interativa disponível em:
- **Local:** `http://localhost:8080/swagger/index.html`
- **Produção:** `https://api.seudominio.com/swagger/index.html`

---

**Total de Endpoints:** 92+  
**Última Atualização:** 03/12/2025
