# 🎯 Plano de Implementação — Módulo de Metas

**Data de Criação:** 27/11/2025  
**Estimativa Total:** 12h  
**Sprint:** 3.1 e 3.2 (28-29/11)  
**Status:** 📋 PLANEJADO

---

## 📊 Análise do Estado Atual

### ✅ Backend — 100% Completo

| Componente | Status | Localização |
|------------|--------|-------------|
| DTOs | ✅ | `backend/internal/application/dto/metas_dto.go` |
| Mapper | ✅ | `backend/internal/application/mapper/metas_mapper.go` |
| Use Cases (15) | ✅ | `backend/internal/application/usecase/metas/` |
| Handler HTTP | ✅ | `backend/internal/infra/http/handler/metas_handler.go` |
| Repository (sqlc) | ✅ | `backend/internal/infra/repository/postgres/` |
| Migrations | ✅ | Tabelas: `metas_mensais`, `metas_barbeiro`, `metas_ticket_medio` |

### ❌ Frontend — 0% Completo

| Componente | Status | Arquivo a Criar |
|------------|--------|-----------------|
| Types/Interfaces | ❌ | `frontend/src/types/metas.ts` |
| Service API | ❌ | `frontend/src/services/metas-service.ts` |
| React Query Hooks | ❌ | `frontend/src/hooks/use-metas.ts` |
| Layout Metas | ❌ | `frontend/src/app/(dashboard)/metas/layout.tsx` |
| Dashboard Metas | ❌ | `frontend/src/app/(dashboard)/metas/page.tsx` |
| Metas Mensais | ❌ | `frontend/src/app/(dashboard)/metas/mensais/page.tsx` |
| Metas Barbeiros | ❌ | `frontend/src/app/(dashboard)/metas/barbeiros/page.tsx` |
| Metas Ticket Médio | ❌ | `frontend/src/app/(dashboard)/metas/ticket/page.tsx` |

---

## 🗂️ Endpoints Backend Disponíveis

### Metas Mensais (Faturamento Geral)
```
POST   /api/v1/metas/monthly      → SetMetaMensal
GET    /api/v1/metas/monthly      → ListMetasMensais
GET    /api/v1/metas/monthly/:id  → GetMetaMensal
PUT    /api/v1/metas/monthly/:id  → UpdateMetaMensal
DELETE /api/v1/metas/monthly/:id  → DeleteMetaMensal
```

### Metas por Barbeiro (Individuais)
```
POST   /api/v1/metas/barbers      → SetMetaBarbeiro
GET    /api/v1/metas/barbers      → ListMetasBarbeiro (?barbeiro_id=uuid)
GET    /api/v1/metas/barbers/:id  → GetMetaBarbeiro
PUT    /api/v1/metas/barbers/:id  → UpdateMetaBarbeiro
DELETE /api/v1/metas/barbers/:id  → DeleteMetaBarbeiro
```

### Metas Ticket Médio
```
POST   /api/v1/metas/ticket      → SetMetaTicket
GET    /api/v1/metas/ticket      → ListMetasTicket
GET    /api/v1/metas/ticket/:id  → GetMetaTicket
PUT    /api/v1/metas/ticket/:id  → UpdateMetaTicket
DELETE /api/v1/metas/ticket/:id  → DeleteMetaTicket
```

---

## 📋 Tarefas Detalhadas

### Fase 1: Tipos e Infraestrutura (1.5h)

#### Task 1.1: Criar `types/metas.ts`
**Tempo:** 30min  
**Arquivo:** `frontend/src/types/metas.ts`

```typescript
// Estrutura esperada:
- MetaMensalResponse (id, mes_ano, meta_faturamento, origem, status, realizado, percentual, timestamps)
- MetaBarbeiroResponse (id, barbeiro_id, mes_ano, metas por categoria, realizados, percentuais)
- MetaTicketMedioResponse (id, mes_ano, tipo, barbeiro_id?, meta_valor, realizado, percentual)
- SetMetaMensalRequest, SetMetaBarbeiroRequest, SetMetaTicketRequest
- Enums: OrigemMeta, StatusMeta, TipoTicketMeta
- Helpers: formatPercentual, calcularProgresso, getStatusColor
```

#### Task 1.2: Criar `services/metas-service.ts`
**Tempo:** 30min  
**Arquivo:** `frontend/src/services/metas-service.ts`

```typescript
// Endpoints a mapear:
- Metas Mensais: list, get, create, update, delete
- Metas Barbeiro: list, get, create, update, delete
- Metas Ticket: list, get, create, update, delete
- Helpers: getResumoMetas (agregado)
```

#### Task 1.3: Criar `hooks/use-metas.ts`
**Tempo:** 30min  
**Arquivo:** `frontend/src/hooks/use-metas.ts`

```typescript
// Hooks React Query:
- useMetasMensais() → lista
- useMetaMensal(id) → detalhe
- useCreateMetaMensal() → mutation
- useUpdateMetaMensal() → mutation
- useDeleteMetaMensal() → mutation

// Repetir para Barbeiro e Ticket
- useMetasBarbeiro(barbeiroId?)
- useMetasTicket()
```

---

### Fase 2: Layout e Dashboard (2h)

#### Task 2.1: Criar `metas/layout.tsx`
**Tempo:** 30min  
**Arquivo:** `frontend/src/app/(dashboard)/metas/layout.tsx`

```
Estrutura:
- Header: "Metas"
- Tabs: [Dashboard | Mensais | Por Barbeiro | Ticket Médio]
- Outlet para subpáginas
- Botão "Nova Meta" no header
```

#### Task 2.2: Criar Dashboard `metas/page.tsx`
**Tempo:** 1.5h  
**Arquivo:** `frontend/src/app/(dashboard)/metas/page.tsx`

```
Componentes:
1. KPIs Cards (3-4)
   - Meta Faturamento Atual vs Realizado
   - % Atingimento Geral
   - Barbeiros Acima da Meta
   - Ticket Médio vs Meta

2. Gráfico de Progresso (Progress Bars)
   - Meta Mensal: [###########-------] 67%
   - Por Barbeiro: ranking top 5
   - Ticket Médio: comparativo

3. Timeline de Metas (últimos 6 meses)
   - Mini gráfico de evolução
```

---

### Fase 3: CRUD Metas Mensais (3h)

#### Task 3.1: Página Lista Metas Mensais
**Tempo:** 1.5h  
**Arquivo:** `frontend/src/app/(dashboard)/metas/mensais/page.tsx`

```
Componentes:
1. Filtro: Ano selector (2024, 2025, ...)
2. Tabela DataGrid:
   - Mês/Ano | Meta | Realizado | % | Status | Ações
   - Célula de % com Progress bar inline
   - Status badge (PENDENTE, ACEITA, REJEITADA)
   
3. Modal Criar/Editar:
   - mes_ano (month picker)
   - meta_faturamento (MoneyInput)
   - origem (MANUAL | AUTOMATICA)
   
4. Modal Confirmação Delete
```

#### Task 3.2: Componentes Auxiliares
**Tempo:** 1h  
**Arquivos:**
- `components/metas/MetaMensalForm.tsx`
- `components/metas/MetaProgressBar.tsx`
- `components/metas/MetaStatusBadge.tsx`

#### Task 3.3: Integração e Testes Manuais
**Tempo:** 30min
- Testar CRUD completo
- Validar filtros
- Verificar responsividade

---

### Fase 4: Metas por Barbeiro (3h)

#### Task 4.1: Página Metas por Barbeiro
**Tempo:** 2h  
**Arquivo:** `frontend/src/app/(dashboard)/metas/barbeiros/page.tsx`

```
Componentes:
1. Filtros:
   - Mês/Ano (month picker)
   - Barbeiro (select - opcional, se vazio mostra todos)

2. Cards por Barbeiro:
   - Avatar + Nome
   - 3 Progress Bars:
     - Serviços Gerais: XX% [========----]
     - Serviços Extras: XX% [======------]
     - Produtos: XX% [============]
   - Total Realizado vs Meta Total
   - Botão Editar/Excluir

3. Modal Criar/Editar Meta Barbeiro:
   - barbeiro_id (select profissionais)
   - mes_ano (month picker)
   - meta_servicos_gerais (MoneyInput)
   - meta_servicos_extras (MoneyInput)
   - meta_produtos (MoneyInput)
   - Preview: Meta Total = soma
   
4. Ranking Visual (se múltiplos barbeiros):
   - Ordenado por % atingimento
   - Medalhas: 🥇 🥈 🥉
```

#### Task 4.2: Componentes
**Tempo:** 1h  
**Arquivos:**
- `components/metas/MetaBarbeiroCard.tsx`
- `components/metas/MetaBarbeiroForm.tsx`
- `components/metas/BarbeiroRanking.tsx`

---

### Fase 5: Metas Ticket Médio (2h)

#### Task 5.1: Página Ticket Médio
**Tempo:** 1.5h  
**Arquivo:** `frontend/src/app/(dashboard)/metas/ticket/page.tsx`

```
Componentes:
1. Seção "Meta Geral da Barbearia"
   - Card grande com gauge/speedometer
   - Ticket Médio Atual vs Meta
   - Variação % mês anterior

2. Seção "Por Barbeiro" (se tipo BARBEIRO)
   - Lista cards menores
   - Cada um com mini gauge

3. Modal Criar/Editar:
   - mes_ano
   - tipo (GERAL | BARBEIRO)
   - barbeiro_id (condicional se tipo=BARBEIRO)
   - meta_valor (MoneyInput)
```

#### Task 5.2: Componentes
**Tempo:** 30min  
**Arquivos:**
- `components/metas/TicketMedioGauge.tsx` (pode usar radix-ui/gauge ou chart)
- `components/metas/MetaTicketForm.tsx`

---

### Fase 6: Integração com Sistema de Bonificação (1.5h)

> **Referência:** `docs/11-Fluxos/FLUXO_METAS.md` — Seção Bonificação Progressiva

#### Task 6.1: Exibir Níveis de Bonificação
**Tempo:** 1h

Adicionar ao `MetaBarbeiroCard.tsx`:
```
Se atingimento >= 100%:
  - Mostrar badge "NÍVEL 1 - Bônus 3%"
  - Se >= 110%: "NÍVEL 2 - Bônus 5%"
  - Se >= 120%: "NÍVEL 3 - Bônus 8%"

Visualização:
[==========|=====|=====|=====] 
  100%     110%  120%  META+20%
   🟢       🟡    🔴    💎
```

#### Task 6.2: Relatório de Bônus Projetado
**Tempo:** 30min

Na página de dashboard, adicionar card:
```
"Bônus Projetado (se mantiver ritmo)"
- Barbeiro X: R$ 150,00 (Nível 2)
- Barbeiro Y: R$ 80,00 (Nível 1)
- Total Projetado: R$ 230,00
```

---

## 📂 Estrutura Final de Arquivos

```
frontend/src/
├── types/
│   └── metas.ts                          # NEW
├── services/
│   └── metas-service.ts                  # NEW
├── hooks/
│   └── use-metas.ts                      # NEW
├── components/
│   └── metas/                            # NEW
│       ├── MetaMensalForm.tsx
│       ├── MetaBarbeiroForm.tsx
│       ├── MetaTicketForm.tsx
│       ├── MetaProgressBar.tsx
│       ├── MetaStatusBadge.tsx
│       ├── MetaBarbeiroCard.tsx
│       ├── BarbeiroRanking.tsx
│       ├── TicketMedioGauge.tsx
│       └── BonificacaoIndicator.tsx
├── app/(dashboard)/
│   └── metas/                            # NEW
│       ├── layout.tsx
│       ├── page.tsx                      # Dashboard
│       ├── mensais/
│       │   └── page.tsx                  # CRUD Metas Mensais
│       ├── barbeiros/
│       │   └── page.tsx                  # CRUD Metas Barbeiro
│       └── ticket/
│           └── page.tsx                  # CRUD Ticket Médio
```

---

## ⏱️ Cronograma de Execução

| Data | Horário | Fase | Tarefas | Horas |
|------|---------|------|---------|-------|
| **28/11** | Manhã | 1 | Types + Service + Hooks | 1.5h |
| **28/11** | Manhã | 2 | Layout + Dashboard | 2h |
| **28/11** | Tarde | 3 | CRUD Metas Mensais | 3h |
| **29/11** | Manhã | 4 | Metas por Barbeiro | 3h |
| **29/11** | Tarde | 5 | Metas Ticket Médio | 2h |
| **29/11** | Tarde | 6 | Bonificação + Testes | 1.5h |

**Total:** 13h (margem de 1h para imprevistos)

---

## ✅ Critérios de Aceite

### Funcionalidade
- [ ] CRUD completo Metas Mensais funcionando
- [ ] CRUD completo Metas por Barbeiro funcionando
- [ ] CRUD completo Metas Ticket Médio funcionando
- [ ] Dashboard exibindo KPIs corretos
- [ ] Filtros de período funcionando
- [ ] Progress bars atualizando em tempo real

### UX/UI
- [ ] Design consistente com Design System
- [ ] Responsivo (mobile-first)
- [ ] Feedback visual em ações (loading, success, error)
- [ ] Modais com validação de formulário

### Qualidade
- [ ] TypeScript sem erros
- [ ] ESLint passando
- [ ] Build sem warnings
- [ ] Testado em Chrome e Firefox

### Integração
- [ ] Dados vindos corretamente do backend
- [ ] Tenant isolation respeitado
- [ ] Autenticação funcionando (401 se não logado)

---

## 🔗 Dependências

### Dependências de Projeto
- ✅ Backend Metas completo
- ✅ Autenticação funcionando
- ✅ Componentes shadcn/ui disponíveis
- ✅ React Query configurado

### Dependências de Dados
- Profissionais cadastrados (para select barbeiro)
- Pelo menos 1 tenant com dados

---

## 🚨 Riscos Identificados

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Complexidade visual do gauge/gráficos | Média | Médio | Usar biblioteca pronta (recharts) |
| Cálculo de realizado não implementado no backend | Baixa | Alto | Backend já tem, validar queries |
| Responsividade dos cards em grid | Média | Baixo | Usar CSS Grid + breakpoints |

---

## 📝 Notas de Implementação

### Padrão de Código
Seguir exatamente o padrão estabelecido em `types/financial.ts`:
- Enums com valores string
- Interfaces com snake_case para campos da API
- Helpers para formatação e validação
- Comentários JSDoc

### Componentes
Usar os mesmos componentes do módulo financeiro:
- `Badge` para status
- `Card` para containers
- `DataTable` para listas
- `Sheet/Dialog` para modais
- `Select` para dropdowns
- `Input` para formulários

### Formulários
- Usar Zod para validação
- React Hook Form para controle
- MoneyInput para valores monetários
- MonthPicker para mês/ano

---

## 🔄 Próximos Passos Após Conclusão

1. Atualizar `TAREFAS_MVP_V1.0.0.md` marcando P3.1 e P3.2 como concluídos
2. Integrar metas no cálculo de comissões (`FLUXO_COMISSOES.md`)
3. Adicionar metas nos relatórios gerenciais
4. Criar testes E2E para o módulo

---

**Responsável:** Equipe NEXO  
**Revisão:** 28/11/2025 após Fase 2
