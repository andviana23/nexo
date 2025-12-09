> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🔐 RBAC - Role-Based Access Control

**Versão:** 1.0
**Última Atualização:** 15/11/2025
**Status:** ✅ Implementado

---

## 📋 Visão Geral

O Barber Analytics Pro implementa **RBAC (Role-Based Access Control)** para controlar o acesso dos usuários a recursos e funcionalidades do sistema. Cada usuário possui uma **role** que define suas permissões.

---

## 👥 Roles Disponíveis

### 🔴 Owner (Proprietário)
**Descrição:** Acesso total ao tenant
**Use Case:** Dono da barbearia

**Permissões:**
- ✅ **Financial:** Criar, visualizar, editar e deletar receitas e despesas
- ✅ **Cashflow:** Visualizar fluxo de caixa
- ✅ **Assinaturas:** Criar, visualizar, editar e deletar assinaturas
- ✅ **Usuários:** Criar, visualizar, editar e deletar usuários
- ✅ **Admin:** Visualizar audit logs, gerenciar feature flags
- ✅ **Dashboard:** Visualizar todos os KPIs

---

### 🟠 Manager (Gerente)
**Descrição:** Gerenciar operações diárias (sem deletar)
**Use Case:** Gerente da barbearia

**Permissões:**
- ✅ **Financial:** Criar, visualizar e editar receitas e despesas
- ❌ **Financial:** ~~Deletar~~ (não permitido)
- ✅ **Cashflow:** Visualizar fluxo de caixa
- ✅ **Assinaturas:** Criar, visualizar e editar assinaturas
- ❌ **Assinaturas:** ~~Deletar~~ (não permitido)
- ✅ **Usuários:** Visualizar apenas (sem criar/editar/deletar)
- ❌ **Admin:** ~~Audit logs, feature flags~~ (não permitido)
- ✅ **Dashboard:** Visualizar todos os KPIs

---

### 🟡 Accountant (Contador)
**Descrição:** Visualizar apenas dados financeiros (somente leitura)
**Use Case:** Contador externo ou assistente administrativo

**Permissões:**
- ✅ **Financial:** Visualizar receitas e despesas (somente leitura)
- ✅ **Cashflow:** Visualizar fluxo de caixa
- ❌ **Financial:** ~~Criar, editar, deletar~~ (não permitido)
- ❌ **Assinaturas:** Sem acesso
- ❌ **Usuários:** Sem acesso
- ❌ **Admin:** Sem acesso
- ✅ **Dashboard:** Visualizar KPIs financeiros

---

### 🟢 Employee (Funcionário)
**Descrição:** Visualizar apenas próprios dados
**Use Case:** Barbeiro/Funcionário

**Permissões:**
- ❌ **Financial:** Sem acesso
- ❌ **Cashflow:** Sem acesso
- ✅ **Assinaturas:** Visualizar apenas próprias assinaturas
- ❌ **Assinaturas:** ~~Criar, editar, deletar~~ (não permitido)
- ❌ **Usuários:** Sem acesso
- ❌ **Admin:** Sem acesso
- ❌ **Dashboard:** Sem acesso

---

## 🔑 Permissões Detalhadas

### Financial

| Permissão | Owner | Manager | Accountant | Employee |
|-----------|-------|---------|------------|----------|
| `receita:create` | ✅ | ✅ | ❌ | ❌ |
| `receita:read` | ✅ | ✅ | ✅ | ❌ |
| `receita:update` | ✅ | ✅ | ❌ | ❌ |
| `receita:delete` | ✅ | ❌ | ❌ | ❌ |
| `despesa:create` | ✅ | ✅ | ❌ | ❌ |
| `despesa:read` | ✅ | ✅ | ✅ | ❌ |
| `despesa:update` | ✅ | ✅ | ❌ | ❌ |
| `despesa:delete` | ✅ | ❌ | ❌ | ❌ |
| `cashflow:read` | ✅ | ✅ | ✅ | ❌ |

### Subscriptions

| Permissão | Owner | Manager | Accountant | Employee |
|-----------|-------|---------|------------|----------|
| `assinatura:create` | ✅ | ✅ | ❌ | ❌ |
| `assinatura:read` | ✅ | ✅ | ❌ | ✅ |
| `assinatura:update` | ✅ | ✅ | ❌ | ❌ |
| `assinatura:delete` | ✅ | ❌ | ❌ | ❌ |

### Users

| Permissão | Owner | Manager | Accountant | Employee |
|-----------|-------|---------|------------|----------|
| `user:create` | ✅ | ❌ | ❌ | ❌ |
| `user:read` | ✅ | ✅ | ❌ | ❌ |
| `user:update` | ✅ | ❌ | ❌ | ❌ |
| `user:delete` | ✅ | ❌ | ❌ | ❌ |

### Admin

| Permissão | Owner | Manager | Accountant | Employee |
|-----------|-------|---------|------------|----------|
| `audit_log:read` | ✅ | ❌ | ❌ | ❌ |
| `feature_flag:read` | ✅ | ❌ | ❌ | ❌ |
| `feature_flag:set` | ✅ | ❌ | ❌ | ❌ |

### Dashboard

| Permissão | Owner | Manager | Accountant | Employee |
|-----------|-------|---------|------------|----------|
| `dashboard:read` | ✅ | ✅ | ✅ | ❌ |

### Appointments

| Permissão | Owner | Manager/Receptionist* | Accountant | Employee |
|-----------|-------|-----------------------|------------|----------|
| `appointment:create` | ✅ | ✅ | ❌ | ✅ (somente para si) |
| `appointment:read` | ✅ | ✅ | ❌ | ✅ (somente para si) |
| `appointment:reschedule` | ✅ | ✅ | ❌ | ✅ (somente para si, sem trocar profissional) |
| `appointment:update` | ✅ | ✅ | ❌ | ✅ (somente para si) |
| `appointment:status` | ✅ | ✅ | ❌ | ✅ (somente para si) |
| `appointment:cancel` | ✅ | ✅ | ❌ | ✅ (somente para si) |
| `appointment:availability` | ✅ | ✅ | ❌ | ✅ (apenas disponibilidade do próprio profissional) |

*Recepcionista utiliza o perfil/role `Manager` para acesso aos agendamentos.

---

## 🛠️ Uso no Backend

### Proteger Endpoint com Permissão Específica

```go
import (
    "github.com/andviana23/barber-analytics-backend-v2/internal/domain/entity"
    httpMiddleware "github.com/andviana23/barber-analytics-backend-v2/internal/infrastructure/http/middleware"
)

// Apenas usuários com permissão de deletar receitas
r.Group(func(r chi.Router) {
    r.Use(httpMiddleware.RequirePermission(entity.PermissionReceitaDelete))
    r.Delete("/receitas/{id}", receitaHandler.Delete)
})
```

### Proteger Endpoint com Role Específica

```go
// Apenas Owner
r.Group(func(r chi.Router) {
    r.Use(httpMiddleware.RequireOwner())
    r.Delete("/users/{id}", userHandler.Delete)
})

// Owner OU Manager
r.Group(func(r chi.Router) {
    r.Use(httpMiddleware.RequireOwnerOrManager())
    r.Post("/receitas", receitaHandler.Create)
})

// Múltiplas roles
r.Group(func(r chi.Router) {
    r.Use(httpMiddleware.RequireRole(entity.RoleOwner, entity.RoleManager, entity.RoleAccountant))
    r.Get("/dashboard", dashboardHandler.Get)
})
```

### Obter Role do Usuário

```go
func (h *ReceitaHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    role, err := httpMiddleware.GetRoleFromContext(ctx)
    if err != nil {
        // Tratar erro
        return
    }

    // Usar role para lógica condicional
    if role == entity.RoleOwner {
        // Owner pode fazer coisas extras
    }
}
```

---

## 🔧 Como Adicionar Nova Role

1. **Definir a role** em `internal/domain/entity/role.go`:
```go
const (
    RoleNewRole Role = "new_role"
)
```

2. **Definir permissões** no `RolePermissions`:
```go
RoleNewRole: {
    PermissionReceitaRead,
    PermissionDespesaRead,
},
```

3. **Adicionar descrição** em `GetRoleDescription`:
```go
RoleNewRole: "Descrição da nova role",
```

4. **Atualizar documentação** neste arquivo

---

## 🔧 Como Adicionar Nova Permissão

1. **Definir a permissão** em `internal/domain/entity/role.go`:
```go
const (
    PermissionNewResource Permission = "new_resource:action"
)
```

2. **Adicionar às roles** que devem ter essa permissão:
```go
RoleOwner: {
    // ... permissões existentes
    PermissionNewResource,
},
```

3. **Usar no middleware**:
```go
r.Use(httpMiddleware.RequirePermission(entity.PermissionNewResource))
```

---

## 🧪 Testes

### Testes Unitários

```bash
# Testar permissões de roles
go test ./tests/unit/entity/ -v -run TestRole

# Testes cobrem:
# - Owner tem todas as permissões
# - Manager pode editar mas não deletar
# - Accountant apenas leitura financeira
# - Employee apenas próprias assinaturas
# - Validação de roles inválidas
```

### Testes E2E

```bash
# Testar acesso negado para role sem permissão
curl -H "Authorization: Bearer <token_manager>" \
     -X DELETE http://localhost:8080/api/v1/receitas/uuid
# Esperado: 403 Forbidden

# Testar acesso permitido para role com permissão
curl -H "Authorization: Bearer <token_owner>" \
     -X DELETE http://localhost:8080/api/v1/receitas/uuid
# Esperado: 200 OK
```

---

## 🔒 Segurança

### Validação de Role

- ✅ Role é validada no middleware de autenticação
- ✅ Role inválida retorna 403 Forbidden
- ✅ Role ausente retorna 403 Forbidden
- ✅ Multi-tenant: Role é sempre validada no contexto do tenant

### Bypass Prevention

- ❌ **NÃO** confiar em role enviada pelo cliente (JWT ou header)
- ✅ **SIM** buscar role do banco de dados após autenticação
- ✅ **SIM** validar role a cada requisição (stateless)

---

## 📊 Matriz de Decisão

| Cenário | Role Recomendada |
|---------|------------------|
| Dono da barbearia | Owner |
| Gerente de operações | Manager |
| Contador/Financeiro | Accountant |
| Barbeiro/Funcionário | Employee |
| Recepcionista | Manager (se precisar criar agendamentos) ou Employee |
| Assistente administrativo | Accountant (se apenas visualizar) ou Manager (se editar) |

---

## 🚀 Roadmap Futuro

- [ ] Permissões personalizadas por tenant (RBAC dinâmico)
- [ ] Audit log de mudanças de role
- [ ] UI para gerenciar roles e permissões
- [ ] Roles temporárias (expiração)
- [ ] Permissões por recurso específico (ex: apenas receitas do próprio usuário)

---

**Última Atualização:** 15/11/2025
**Autor:** Andrey Viana
**Status:** ✅ Produção
