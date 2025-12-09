# PRD — Módulo de Agendamento | NEXO v1.0

**Versão do Documento:** 2.0.0  
**Status:** 🟢 COMPLETO (100% Implementado)
**Prioridade:** 🔴 CRÍTICA  
**Data de Criação:** 25/11/2025  
**Última Atualização:** 30/11/2025  
**Responsável:** Andrey Viana (Product Owner)  
**Milestone:** 1.5 (10/12/2025) - ✅ CONCLUÍDO  

---

## 📊 Status de Implementação

| Área | Status | Progresso |
|------|--------|-----------|
| Backend (Go) | ✅ Completo | 100% |
| Frontend Base | ✅ Completo | 100% |
| Comanda/Pagamento | ✅ Completo | 100% |
| Menu de Ações | ✅ Completo | 100% |

### ✅ Implementado (100%)
- ✅ Backend: 6 endpoints, 36 testes
- ✅ Frontend: Calendário FullCalendar, Modal, Selectors
- ✅ React Query com Optimistic Updates
- ✅ Menu de Contexto (Botão Direito) com ações dinâmicas
- ✅ CommandModal (Comanda estilo Trinks) - Layout 2 colunas **COMPLETO**
- ✅ Bloqueios de horário (backend + frontend)
- ✅ Integração Comanda ↔ Agendamento
- ✅ Campo `command_id` em appointments
- ✅ Todos os status: CREATED → CONFIRMED → CHECKED_IN → IN_SERVICE → AWAITING_PAYMENT → DONE
- ✅ Validação de conflitos
- ✅ Multi-tenant + RBAC

### 🚀 Próximas Versões (v2.0+)
- 🔵 Google Calendar Integration (OAuth 2.0)
- 🔵 Notificações WhatsApp (confirmação, lembretes)
- 🔵 Drag & Drop avançado no calendário
- 🔵 Sistema de pontos de fidelidade
- 🔵 Relatórios avançados de ocupação

---

## 1. Executive Summary

### 1.1 Visão Geral

O **Módulo de Agendamento** é o componente central do NEXO, responsável por permitir que barbearias agendem, gerenciem e acompanhem serviços de forma visual, intuitiva e profissional.

**Problema:** Barbearias premium perdem dinheiro com:
- Conflitos de horário (double booking)
- Faltas de cliente sem confirmação (no-show alto)
- Agenda desorganizada (papel, WhatsApp, Excel)
- Falta de visibilidade da ocupação do barbeiro
- Impossibilidade de otimizar agenda

**Solução:** Calendário visual profissional com:
- ✅ Validação automática de conflitos
- ✅ Confirmação de agendamento (reduz no-show)
- ✅ Sincronização com Google Agenda
- ✅ Visão por barbeiro em tempo real
- ✅ CRUD completo (criar, editar, cancelar, reagendar)

---

## 2. Objetivos do Produto

### 2.1 Objetivo Principal

**Permitir que barbearias gerenciem agendamentos de forma visual, sem conflitos e com máxima ocupação dos barbeiros.**

### 2.2 Objetivos Secundários

1. **Reduzir no-show** de clientes (meta: < 10%)
2. **Aumentar ocupação** dos barbeiros (meta: > 80%)
3. **Eliminar conflitos** de horário (meta: 0%)
4. **Melhorar experiência** do cliente (confirmação, lembretes)
5. **Otimizar operação** da recepção (agendar em < 30s)

---

## 3. Métricas de Sucesso (KPIs)

| KPI | Baseline | Meta | Medição |
|-----|----------|------|---------|
| **Taxa de No-Show** | 25% | < 10% | (Agendamentos NO_SHOW / Total) × 100 |
| **Ocupação Média** | 60% | > 80% | Horas agendadas / Horas disponíveis |
| **Conflitos/Mês** | 15 | 0 | Count de conflitos registrados |
| **Tempo de Agendamento** | 3 min | < 30s | Tempo médio para criar agendamento |
| **NPS** | N/A | > 8.0 | Pesquisa de satisfação |

---

## 4. Personas e Necessidades

### 4.1 Persona 1: Dono da Barbearia

**Nome:** Carlos, 38 anos  
**Objetivo:** Maximizar lucro e otimizar operação  

**Necessidades:**
- 🔴 Ver ocupação de todos os barbeiros
- 🔴 Identificar horários vazios
- 🔴 Acompanhar no-show por barbeiro
- 🟡 Exportar dados para análise

**Pain Points:**
- Não sabe se está perdendo dinheiro com horários vazios
- Não consegue medir performance dos barbeiros
- Perde tempo resolvendo conflitos de agenda

---

### 4.2 Persona 2: Gerente/Recepção

**Nome:** Juliana, 26 anos  
**Objetivo:** Manter agenda organizada e otimizada  

**Necessidades:**
- 🔴 Agendar clientes rapidamente (< 30s)
- 🔴 Ver disponibilidade de barbeiros
- 🔴 Evitar conflitos de horário
- 🟡 Confirmar agendamentos
- 🟡 Remarcar/cancelar facilmente

**Pain Points:**
- Perde tempo checando disponibilidade manualmente
- Conflitos causam estresse e retrabalho
- Cliente insatisfeito com espera

---

### 4.3 Persona 3: Barbeiro

**Nome:** Rafael, 29 anos  
**Objetivo:** Focar no atendimento, sem preocupação com agenda  

**Necessidades:**
- 🔴 Ver apenas seus próprios agendamentos
- 🔴 Saber quem é o próximo cliente
- 🟡 Sincronizar com Google Agenda pessoal
- 🟢 Bloquear horários para almoço/pausa

**Pain Points:**
- Precisa ficar perguntando quem é o próximo
- Não sabe se tem horário livre amanhã
- Agenda do Google desatualizada

---

### 4.4 Persona 4: Cliente Final

**Nome:** Pedro, 32 anos  
**Objetivo:** Agendar corte sem fricção  

**Necessidades:**
- 🔴 Agendar online (futuro - app)
- 🔴 Receber confirmação
- 🟡 Receber lembrete 1h antes
- 🟡 Remarcar facilmente

**Pain Points:**
- Precisa ligar para agendar (inconveniente)
- Esquece do horário marcado
- Dificuldade para remarcar

---

## 5. Requisitos Funcionais (RF)

### 5.1 CRUD de Agendamento

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-001** | Sistema DEVE permitir criar novo agendamento | 🔴 P0 | ⬜ |
| **RF-002** | Sistema DEVE permitir editar agendamento existente | 🔴 P0 | ⬜ |
| **RF-003** | Sistema DEVE permitir cancelar agendamento | 🔴 P0 | ⬜ |
| **RF-004** | Sistema DEVE permitir reagendar (mover data/hora) | 🔴 P0 | ⬜ |
| **RF-005** | Sistema DEVE exibir calendário visual | 🔴 P0 | ⬜ |

### 5.2 Validação e Conflitos

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-006** | Sistema DEVE validar disponibilidade do barbeiro | 🔴 P0 | ⬜ |
| **RF-007** | Sistema DEVE impedir conflitos de horário | 🔴 P0 | ⬜ |
| **RF-008** | Sistema DEVE sugerir horários alternativos | 🟡 P1 | ⬜ |
| **RF-009** | Sistema DEVE validar duração do serviço | 🔴 P0 | ⬜ |
| **RF-010** | Sistema DEVE respeitar intervalo mínimo (10min) | 🟡 P1 | ⬜ |

### 5.3 Visualização

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-011** | Sistema DEVE exibir view diária | 🔴 P0 | ⬜ |
| **RF-012** | Sistema DEVE exibir view semanal | 🔴 P0 | ⬜ |
| **RF-013** | Sistema DEVE exibir view mensal | 🟡 P1 | ⬜ |
| **RF-014** | Sistema DEVE permitir filtrar por barbeiro | 🔴 P0 | ⬜ |
| **RF-015** | Sistema DEVE permitir filtrar por status | 🟡 P1 | ⬜ |
| **RF-016** | Sistema DEVE exibir cores por status | 🔴 P0 | ⬜ |

### 5.4 Status e Lifecycle

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-017** | Sistema DEVE suportar status CREATED | 🔴 P0 | ⬜ |
| **RF-018** | Sistema DEVE suportar status CONFIRMED | 🔴 P0 | ⬜ |
| **RF-019** | Sistema DEVE suportar status IN_SERVICE | 🔴 P0 | ⬜ |
| **RF-020** | Sistema DEVE suportar status DONE | 🔴 P0 | ⬜ |
| **RF-021** | Sistema DEVE suportar status NO_SHOW | 🔴 P0 | ⬜ |
| **RF-022** | Sistema DEVE suportar status CANCELED | 🔴 P0 | ⬜ |
| **RF-023** | Sistema DEVE registrar histórico de mudanças de status | 🟡 P1 | ⬜ |

### 5.5 Integrações

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-024** | Sistema DEVE sincronizar com Google Agenda | 🟡 P1 | ⬜ |
| **RF-025** | Sistema DEVE permitir conectar conta Google (OAuth) | 🟡 P1 | ⬜ |
| **RF-026** | Sistema DEVE atualizar Google Agenda em alterações | 🟡 P1 | ⬜ |
| **RF-027** | Sistema DEVE remover de Google Agenda em cancelamentos | 🟡 P1 | ⬜ |

### 5.6 Notificações (Futuro)

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-028** | Sistema DEVE enviar confirmação via WhatsApp | 🟢 P2 | ⬜ |
| **RF-029** | Sistema DEVE enviar lembrete 1h antes | 🟢 P2 | ⬜ |
| **RF-030** | Sistema DEVE enviar lembrete 24h antes | 🟢 P2 | ⬜ |

---

## 6. Requisitos Não Funcionais (RNF)

### 6.1 Performance

| ID | Requisito | Meta | Medição |
|----|-----------|------|---------|
| **RNF-001** | Tempo de carregamento do calendário | < 1s | P95 |
| **RNF-002** | Latência da API de criação de agendamento | < 200ms | P95 |
| **RNF-003** | Latência da API de validação de conflitos | < 150ms | P95 |
| **RNF-004** | Sincronização Google Calendar | < 500ms | P95 |

### 6.2 Escalabilidade

| ID | Requisito | Meta |
|----|-----------|------|
| **RNF-005** | Suportar 1000 agendamentos/dia por tenant | ✅ |
| **RNF-006** | Suportar 50 barbeiros por tenant | ✅ |
| **RNF-007** | Suportar 10.000 agendamentos históricos | ✅ |

### 6.3 Disponibilidade

| ID | Requisito | Meta |
|----|-----------|------|
| **RNF-008** | Uptime do módulo | > 99.5% |
| **RNF-009** | Tempo de recuperação (MTTR) | < 5 min |

### 6.4 Segurança

| ID | Requisito | Descrição |
|----|-----------|-----------|
| **RNF-010** | Isolamento multi-tenant | TODOS os dados filtrados por `tenant_id` |
| **RNF-011** | Validação de permissões (RBAC) | Barbeiro só vê própria agenda |
| **RNF-012** | Auditoria de ações | Registrar CRUD em `audit_logs` |
| **RNF-013** | Proteção contra CSRF | Tokens CSRF em formulários |

### 6.5 Usabilidade

| ID | Requisito | Descrição |
|----|-----------|-----------|
| **RNF-014** | Responsividade | Mobile, Tablet, Desktop |
| **RNF-015** | Acessibilidade | WCAG 2.1 AA |
| **RNF-016** | Feedback visual | Loading states, toasts, confirmações |

---

## 7. Regras de Negócio (RN)

### RN-AGE-001: Validação de Barbeiro

**Descrição:** Sistema DEVE validar que o barbeiro está ativo e pertence ao tenant.

**Critérios:**
- ❌ Não pode agendar com barbeiro `ativo = false`
- ❌ Não pode agendar com barbeiro de outro tenant
- ✅ Barbeiro deve ter horário disponível no slot

**Exceção:** `ErrProfessionalInactive` ou `ErrProfessionalNotFound`

---

### RN-AGE-002: Validação de Cliente

**Descrição:** Cliente DEVE existir antes de criar agendamento.

**Critérios:**
- ✅ Cliente com `id` válido
- ✅ Cliente pertence ao mesmo `tenant_id`
- ✅ Cliente `ativo = true`

**Fluxo Alternativo:** Se cliente não existe, sistema DEVE redirecionar para "Cadastrar Cliente".

---

### RN-AGE-003: Intervalo Mínimo

**Descrição:** Deve haver intervalo mínimo de 10 minutos entre agendamentos do mesmo barbeiro.

**Critérios:**
- ✅ `start_time` do novo agendamento >= `end_time` do anterior + 10 min
- ✅ Configurável por tenant (futuro)

**Exceção:** `ErrInsufficientInterval`

---

### RN-AGE-004: Estrutura do Agendamento

**Descrição:** Todo agendamento DEVE ter:

- ✅ 1 tenant (`tenant_id`)
- ✅ 1 barbeiro (`professional_id`)
- ✅ 1 cliente (`customer_id`)
- ✅ 1+ serviços (`service_ids[]`)
- ✅ Data/hora de início (`start_time`)
- ✅ Data/hora de fim (`end_time`)

**Validação:**
- `end_time` > `start_time`
- `service_ids` não pode ser vazio

---

### RN-AGE-005: Status Lifecycle

**Descrição:** Status DEVE seguir transições válidas:

```
CREATED
  ├─> CONFIRMED
  │     ├─> IN_SERVICE
  │     │     ├─> DONE
  │     │     └─> CANCELED
  │     ├─> NO_SHOW
  │     └─> CANCELED
  └─> CANCELED
```

**Transições Proibidas:**
- ❌ `DONE` → `CREATED`
- ❌ `CANCELED` → `CONFIRMED`
- ❌ `NO_SHOW` → `IN_SERVICE`

---

### RN-AGE-006: Conflitos de Horário

**Descrição:** Sistema DEVE impedir conflitos (overlapping) de horário para o mesmo barbeiro.

**Critérios de Conflito:**
```sql
-- Conflito se:
(novo.start_time < existente.end_time) 
AND 
(novo.end_time > existente.start_time)
AND
(existente.status NOT IN ('CANCELED', 'NO_SHOW'))
```

**Exceção:** `ErrTimeSlotConflict`

---

### RN-AGE-007: Sincronização Google Agenda

**Descrição:** Sincronizar APENAS agendamentos com status `CONFIRMED`.

**Regras:**
- ✅ Criar evento no Google ao confirmar
- ✅ Atualizar evento ao reagendar
- ✅ Deletar evento ao cancelar
- ❌ NÃO sincronizar status `CREATED` (pendente)

**Requisitos:**
- Barbeiro deve ter conectado conta Google (OAuth 2.0)
- Armazenar `google_event_id` na tabela `appointments`

---

### RN-AGE-008: Duração do Serviço

**Descrição:** Sistema DEVE calcular `end_time` baseado na soma da duração dos serviços.

**Fórmula:**
```
end_time = start_time + SUM(servicos.duracao_minutos)
```

**Exemplo:**
- Serviço 1: Corte (30 min)
- Serviço 2: Barba (15 min)
- Total: 45 min
- `start_time`: 14:00
- `end_time`: 14:45

---

### RN-AGE-009: Permissões (RBAC)

**Descrição:** Controle de acesso por role.

| Role | Ver Todos | Criar | Editar | Cancelar | Ver Próprios |
|------|-----------|-------|--------|----------|--------------|
| **owner** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **manager** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **receptionist** | ✅ (unidade) | ✅ | ✅ | ✅ | ✅ |
| **barbeiro** | ❌ | ❌ | ❌ | ❌ | ✅ |

**Validação:** Middleware RBAC no backend.

---

### RN-AGE-010: Multi-Tenant Enforcement

**Descrição:** NENHUM dado pode cruzar entre tenants.

**Validação Obrigatória:**
```go
// TODAS as queries DEVEM filtrar por tenant_id
SELECT * FROM appointments 
WHERE tenant_id = $1  -- OBRIGATÓRIO
  AND id = $2
```

**Exceção:** `ErrUnauthorizedTenant` (HTTP 403)

---

## 8. Edge Cases e Exceções

### 8.1 Conflito de Criação Simultânea

**Cenário:** Dois usuários agendando o mesmo horário ao mesmo tempo.

**Solução:**
1. Validação de conflito no backend com transaction
2. Lock pessimista ou otimista (versioning)
3. Retornar erro `409 Conflict` para o segundo request
4. Frontend exibe mensagem: "Horário foi reservado. Escolha outro."

---

### 8.2 Barbeiro Inativado Durante Agendamento

**Cenário:** Barbeiro foi inativado enquanto recepção estava criando agendamento.

**Solução:**
1. Backend valida status do barbeiro no momento do `POST`
2. Se inativo, retorna `400 Bad Request`
3. Frontend exibe: "Barbeiro não está mais disponível."

---

### 8.3 Cliente Deletado

**Cenário:** Cliente foi deletado (LGPD) mas possui agendamentos futuros.

**Solução:**
1. Soft delete: `clientes.ativo = false`
2. Agendamentos permanecem (FK com `ON DELETE RESTRICT`)
3. Anonimizar dados: `nome = "Cliente Removido"`

---

### 8.4 Sincronização Google Falhou

**Cenário:** API do Google Calendar retornou erro.

**Solução:**
1. Agendamento é criado localmente (sempre prioritário)
2. Erro é registrado em `audit_logs`
3. Sistema tenta reprocessar em background (retry com exponential backoff)
4. Frontend NÃO bloqueia (exibe warning: "Sincronização pendente")

---

### 8.5 Fuso Horário

**Cenário:** Tenant está em fuso diferente (ex: Manaus vs SP).

**Solução:**
1. Armazenar SEMPRE em UTC no banco
2. Converter para timezone do tenant no frontend
3. Usar `tenant_settings.timezone` (default: `America/Sao_Paulo`)

---

## 9. Critérios de Aceite

### 9.1 Funcionalidades Mínimas (MVP)

- [ ] ✅ Criar agendamento com múltiplos serviços
- [ ] ✅ Editar agendamento existente
- [ ] ✅ Cancelar agendamento
- [ ] ✅ Reagendar (mudar data/hora)
- [ ] ✅ Visualizar calendário semanal
- [ ] ✅ Visualizar calendário diário
- [ ] ✅ Filtrar por barbeiro
- [ ] ✅ Validação de conflitos (tempo real)
- [ ] ✅ Status lifecycle (6 status)
- [ ] ✅ Isolamento multi-tenant (100%)
- [ ] ✅ RBAC (barbeiro read-only)

### 9.2 Integrações (v1.1)

- [ ] 🟡 Sincronização Google Agenda
- [ ] 🟡 OAuth 2.0 (conectar conta Google)
- [ ] 🟢 Notificações WhatsApp (futuro)

### 9.3 UX/UI

- [ ] ✅ Responsivo (mobile, tablet, desktop)
- [ ] ✅ Loading states em todas as ações
- [ ] ✅ Toast de sucesso/erro
- [ ] ✅ Confirmação antes de cancelar
- [ ] ✅ Drag & drop para reagendar (v1.1)

---

## 10. Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Conflitos não detectados | Média | Alto | Validação robusta com transactions |
| Performance com muitos agendamentos | Baixa | Médio | Índices no banco + paginação |
| Sincronização Google falha | Média | Baixo | Retry assíncrono + log de erros |
| Barbeiro edita agenda no Google | Baixa | Médio | Documentar: "NEXO é fonte de verdade" |

---

## 11. Roadmap e Priorização

### v1.0 (MVP) - 10/12/2025

- ✅ CRUD completo
- ✅ Validação de conflitos
- ✅ Calendário visual (FullCalendar)
- ✅ Multi-tenant + RBAC

### v1.1 - 15/01/2026

- 🟡 Google Agenda integration
- 🟡 Drag & drop reagendamento
- 🟡 Notificações por email

### v2.0 - Futuro

- 🟢 App do cliente (agendamento self-service)
- 🟢 Notificações WhatsApp
- 🟢 Bloqueio de horários (férias, almoço)
- 🟢 Agendamento recorrente

---

## 12. Restrições e Observações Técnicas

### 12.1 Licença FullCalendar Scheduler – Modo Avaliação

O NEXO utiliza o **FullCalendar Premium (Scheduler)** durante o período de **avaliação gratuita** para fins exclusivamente de desenvolvimento interno.

**Chave de Licença (Temporária):**

```javascript
schedulerLicenseKey: 'CC-Attribution-NonCommercial-NoDerivatives'
```

**⚠️ Restrições Legais:**

- ❌ **Proibido uso comercial** desta licença.
- ✅ **Permitido apenas** para:
  - Desenvolvimento interno
  - Testes de integração e homologação
  - Demonstrações internas (não para clientes finais)
- ⚠️ **A versão final do NEXO que será usada por barbearias exigirá a compra da licença oficial.**
- 🔄 **Substituir a chave de desenvolvimento pela licença comercial antes do lançamento em produção.**

**Status Atual:**

| Item | Status |
|------|--------|
| Licença de Desenvolvimento | ✅ Ativa (Modo Avaliação) |
| Licença Comercial | ⬜ Pendente (Compra antes da Produção) |
| Ambiente Permitido | Desenvolvimento, Staging |
| Ambiente Bloqueado | Produção (até compra da licença) |

**Referência:** [FullCalendar Pricing](https://fullcalendar.io/pricing)

---

## 13. Dependências Externas

| Dependência | Versão | Propósito |
| **FullCalendar** | 6.x | Calendário visual |
| **Google Calendar API** | v3 | Sincronização |
| **PostgreSQL** | 14+ | Banco de dados |
| **Next.js** | 15.5.6 | Frontend framework |
| **Go** | 1.24 | Backend |

---

## 13. Glossário

| Termo | Definição |
|-------|-----------|
| **Agendamento** | Reserva de horário para cliente com barbeiro específico |
| **Conflito** | Overlapping de horários do mesmo barbeiro |
| **No-Show** | Cliente faltou sem avisar |
| **Slot** | Intervalo de tempo disponível para agendamento |
| **RBAC** | Role-Based Access Control (controle por função) |
| **Multi-Tenant** | Isolamento total de dados entre barbearias |

---

## 14. Anexos

### 14.1 Wireframes

Ver: `docs/Agendamento/DIAGRAMAS_AGENDAMENTO.md`

### 14.2 Fluxos Detalhados

Ver: `docs/11-Fluxos/FLUXO_AGENDAMENTO.md`

### 14.3 Schema de Banco

Ver: `docs/Agendamento/BANCO_AGENDAMENTO.md`

---

**Aprovado por:** Andrey Viana (Product Owner)  
**Data de Aprovação:** 25/11/2025  
**Próxima Revisão:** 10/12/2025  

---

**🚀 Este PRD é a base para a implementação do módulo mais crítico do NEXO! 🚀**
