# Análise do Sistema Atual — Barber Analytics Pro / NEXO v2.0

**Contexto:** Análise baseada no PRD-NEXO (21/11/2025) e no código presente neste repositório. Foco em descrever o que já existe e como funciona hoje.

---

## 1. Análise do PRD-NEXO.md
- Visão: ERP/CRM para barbearias premium, multi-unidade, com módulos de agendamento, lista da vez, assinaturas (Asaas), financeiro (caixa/DRE/comissões), estoque, CRM, fidelidade, gamificação, metas/KPIs, precificação, relatórios, apps e integrações externas.
- Regras críticas: multi-tenant obrigatório; bloqueio de benefícios se assinatura inadimplente; comissão só sobre serviços pagos; lista da vez com barbeiros ativos e reset mensal; estoque nunca negativo; privacidade de dados para barbeiro.
- Roadmap: MVP inclui agenda, lista da vez, financeiro básico, comissões, estoque essencial, assinaturas, CRM básico, relatórios mensais simples, permissões. Pós-MVP: fidelidade, gamificação, metas avançadas, precificação, relatórios avançados, apps. Futuro: NF, integrações avançadas, IA.

## 2. Visão geral do sistema realmente implementado
- Backend Go (Clean Architecture) com foco em financeiro, metas e precificação. Apenas rotas POST em handlers; `cmd/api/main.go` não instancia repos/use cases/handlers além de `/health` e `/ping`. Sem auth/RBAC configurados.
- Domain/VO prontos para financeiro/metas/precificação; ports definidos. Repositórios sqlc gerados, implementação parcial (T-CON-003 em 70%).
- Use cases implementados: contas a pagar/receber (criar, marcar pagamento/recebimento), fluxo diário, DRE, compensação bancária, metas (mensal/barbeiro/ticket), precificação (config/simulação/histórico).
- Scheduler criado (`internal/infra/scheduler`) com cron/Prometheus/feature flag, mas sem wiring no main e sem fonte de tenants.
- Frontend Next.js quase vazio; serviços TS criados em `frontend/lib/services` (financial/metas/pricing/stock placeholder) não integrados a UI/hooks.

## 3. Lista das funcionalidades existentes (código)
- Financeiro: Criar Conta a Pagar; Criar Conta a Receber; Marcar Pagamento; Marcar Recebimento; Gerar Fluxo de Caixa Diário; Gerar DRE Mensal; Criar/Marcar Compensação Bancária.
- Metas: Definir Meta Mensal; Definir Meta Barbeiro; Definir Meta Ticket Médio.
- Precificação: Salvar Configuração (margem/markup/impostos/comissão); Simular Preço; Salvar Simulação (use case).
- Scheduler (infra): registro de jobs com métricas/flags, chamados de use cases; placeholders para notificações, estoque mínimo, comissões.
- Frontend services: clientes HTTP com retries/Zod para financeiro, metas, precificação; stock placeholder.

## 4. Fluxo detalhado de cada funcionalidade

### Contas a Pagar — POST `/api/v1/financial/payables`
`internal/infra/http/handler/financial_handler.go`
```
[Início] → Handler recebe JSON
        → c.Get("tenant_id") (assume middleware)
        → Bind + Validate DTO CreateContaPagarRequest
        → Mapper → Money/TipoCusto/Data
        → UseCase CreateContaPagar (porta ContaPagarRepository)
        → Repo persiste
        → Mapper → ContaPagarResponse
        → HTTP 201
[Fim]
```
Entradas: descricao, categoria_id, fornecedor, valor (string), tipo FIXO/VARIAVEL, data_vencimento.  
Saídas: DTO com strings monetárias/datas.  
Validações: required, enum tipo, data válida, valor string parsável.  
Dependências: VO Money/TipoCusto; repos ContaPagar; logger Zap.

### Contas a Receber — POST `/api/v1/financial/receivables`
`financial_handler.go`
```
[Início] → c.Get("tenant_id")
        → Bind + Validate DTO CreateContaReceberRequest
        → Mapper → Money/Data
        → UseCase CreateContaReceber
        → Repo persiste
        → Mapper → ContaReceberResponse
        → HTTP 201
[Fim]
```
Entradas: origem, assinatura_id opcional, descricao_origem, valor (string), data_vencimento.  
Saídas: valores pago/aberto em string.  
Validações: required, uuid opcional, data, valor.

### Marcar Pagamento — POST `/api/v1/financial/payables/:id/pay`
`financial_handler.go`
```
[Início] → tenant_id do contexto
        → id path
        → Bind + Validate DTO (data_pagamento)
        → UseCase MarcarPagamento → busca conta, valida estado, marca paga
        → HTTP 200 SuccessResponse
[Fim]
```

### Marcar Recebimento — POST `/api/v1/financial/receivables/:id/receive`
`financial_handler.go`
```
[Início] → tenant_id
        → id path
        → Bind + Validate DTO (valor_pago, data_recebimento)
        → UseCase MarcarRecebimento → busca conta, marca recebida
        → HTTP 200 SuccessResponse
[Fim]
```

### Fluxo de Caixa Diário (use case)
`internal/application/usecase/financial/generate_fluxo_diario.go`
```
[Início] → Validate tenant_id, data (default hoje)
        → Repo.FindByData ou cria entidade
        → Saldo inicial = saldo final do dia anterior (se existir)
        → Entradas confirmadas = sum contas receber pagas no dia
        → Entradas previstas = sum contas receber pendentes no dia
        → Saídas pagas = sum contas pagar pagas no dia
        → Saídas previstas = sum contas pagar pendentes no dia
        → Calcula saldo final
        → Persistir (Create/Update)
[Fim]
```
Depende: FluxoCaixaDiarioRepository, ContaPagar/ReceberRepository.

### DRE Mensal (use case)
`internal/application/usecase/financial/generate_dre.go`
```
[Início] → Validate tenant_id; mes_ano default mês anterior
        → Repo.FindByMesAno ou cria DRE
        → Receitas: sum contas receber pagas no mês (placeholder atribui a serviços)
        → Custos variáveis: zero (placeholder)
        → Despesas: sum contas pagar pagas (placeholder tudo como despesa fixa)
        → dre.Calcular()
        → Persistir (Create/Update)
[Fim]
```
Depende: DREMensalRepository, ContaPagar/ReceberRepository.

### Compensação Bancária
`internal/application/usecase/financial/create_compensacao.go` e `marcar_compensacao.go`
```
Create: valida tenant/ids/datas, cria entidade CompensacaoBancaria, salva via repo.
Marcar (batch/individual): busca pendentes, marca compensado (domínio), atualiza repo.
```

### Metas
`internal/infra/http/handler/metas_handler.go`
- `/api/v1/metas/monthly`:
```
[Início] → tenant_id
        → Bind + Validate SetMetaMensalRequest (mes_ano, meta_faturamento, origem)
        → Mapper → MesAno/Money/OrigemMeta
        → UseCase SetMetaMensal
        → HTTP 200 MetaMensalResponse
[Fim]
```
- `/api/v1/metas/barbers`:
```
[Início] → tenant_id
        → Bind + Validate SetMetaBarbeiroRequest
        → Mapper → MesAno/Money
        → UseCase SetMetaBarbeiro
        → HTTP 200 MetaBarbeiroResponse
[Fim]
```
- `/api/v1/metas/ticket`:
```
[Início] → tenant_id
        → Bind + Validate SetMetaTicketRequest (tipo GERAL/BARBEIRO)
        → Mapper → MesAno/TipoMetaTicket/Money
        → UseCase SetMetaTicket
        → HTTP 200 MetaTicketResponse
[Fim]
```

### Precificação
`internal/infra/http/handler/pricing_handler.go`
- `/api/v1/pricing/config`:
```
[Início] → tenant_id
        → Bind + Validate SaveConfigPrecificacaoRequest
        → Mapper → Percentuais/Decimal
        → UseCase SaveConfigPrecificacao (cria/atualiza por tenant)
        → HTTP 200 PrecificacaoConfigResponse
[Fim]
```
- `/api/v1/pricing/simulate`:
```
[Início] → tenant_id
        → Bind + Validate SimularPrecoRequest
        → Mapper → Money/Percentual + parametros opcionais
        → UseCase SimularPreco (usa config do tenant)
        → HTTP 200 PrecificacaoSimulacaoResponse
[Fim]
```

### Scheduler (infra)
`internal/infra/scheduler/{scheduler.go,jobs.go}`
```
[Início] → Configura cron com schedules/flags via ENV
        → Registra jobs:
            - GenerateDREMonthly (use case, mes anterior)
            - GenerateFluxoDiario (use case, Data=agora)
            - MarcarCompensacoes (use case batch)
            - NotifyPayables (placeholder noop)
            - CheckEstoqueMinimo (placeholder noop)
            - CalculateComissoes (placeholder noop)
        → Métricas Prometheus (duração/erros)
        → Execução por lista de tenants (env)
[Fim]
```
Status: não integrado ao main; sem fonte real de tenants; sem persistência em `cron_run_logs`.

### Frontend Services (TS)
`frontend/lib/services/*`
```
client.ts: fetch com Authorization Bearer opcional, retries exponenciais, Zod opcional.
payables/receivables/metas/pricing/dre/fluxo: chamadas REST para rotas /api/v1 existentes; schemas em schemas.ts.
stockService: placeholder para /stock/*.
```

## 5. Comparação PRD vs Código atual
- Coberto parcialmente: Financeiro (contas, fluxo diário, DRE simplificado), Metas (mensal/barbeiro/ticket), Precificação (config/simulação), Compensação bancária.
- Não implementado: Agendamento, Lista da vez, Assinaturas/Asaas, CRM, Fidelidade, Gamificação, Comissões reais, Estoque real, Relatórios avançados, Apps, Integrações externas (Asaas/Google), RBAC, multi-unidade prático (apenas campo tenant_id).
- Roadmap MVP: somente parte do financeiro/metas/precificação começou; demais módulos ausentes.

## 6. Pontos de risco e problemas reais
- Sem wiring/DI: `cmd/api/main.go` não registra repos/use cases/handlers; rotas POST não expostas.
- Auth/RBAC inexistentes: handlers assumem `tenant_id` no contexto; falta middleware JWT/RBAC → risco de vazamento/panic.
- Validator não configurado no Echo → `c.Validate` pode causar panic/500.
- Sem rotas GET/PUT/DELETE para leitura/atualização; apenas POST (criação/ação).
- Repositórios incompletos (T-CON-003 em 70%); sqlc gerado, mas implementações faltantes e não instanciadas.
- Scheduler não iniciado; jobs placeholders e sem tenants provider; sem tabela `cron_run_logs`.
- Frontend sem UI/hooks; serviços apontam para rotas não publicadas (falha de integração).
- Multi-tenant frágil (A5 do diagrama): falta context type-safe, enforcement e isolamento real.
- Integrações externas ausentes (Asaas/Google); requisitos críticos do PRD não atendidos.

## 7. Sugestões de ajustes técnicos
1) Bootstrap/DI no `main`: instanciar repos sqlc, use cases, handlers, registrar rotas `/api/v1/*`, configurar validator e middleware de auth/RBAC (JWT RS256 + roles).
2) Concluir repositórios (T-CON-003) e adicionar testes de integração (tenant isolation, paginação, UNIQUE).
3) Scheduler: integrar no main com graceful shutdown, provider de tenants, migration/repo para `cron_run_logs`, implementar NotifyPayables/CheckEstoqueMinimo/CalculateComissoes de acordo com domínio.
4) Completar endpoints GET/PUT/DELETE para contas, fluxo e DRE; adicionar filtros/paginação e validação.
5) Fortalecer multi-tenant: contexto type-safe de tenant, RBAC por rota, validação de tenant em todos os handlers/repos, evitar panic de `c.Get`.
6) Frontend: provider de auth/token, base URL correta, hooks React Query, páginas mínimas para consumir serviços; alinhar schemas às respostas reais do backend.
7) Planejar e implementar módulos faltantes do PRD (agendamento, lista da vez, assinaturas/Asaas, estoque) com ADRs e docs, seguindo Clean Arch e sqlc.
8) Observabilidade: expor métricas Prometheus, logs estruturados, preparar tracing; validar backups e auditoria mínima.

## 8. Sugestões de fluxos futuros (conforme PRD)
- Agendamento: API/serviço com validação de barbeiro/cliente, bloqueio de horário, sync Google (criar/atualizar/cancelar), drag & drop/reagendamento.
- Lista da vez: serviço com reset mensal automático, histórico, API para pausar/retomar, atualização de ordem ao finalizar atendimento.
- Assinaturas/Asaas: cliente resiliente (retry/backoff/circuit breaker), webhooks, sync diária, bloqueio de benefícios inadimplentes.
- Comissões: cálculo por serviço pago, percentuais por barbeiro, geração de despesas de comissão e relatórios.
- Estoque: entidades/use cases/handlers para entradas/saídas/consumo, ficha técnica por serviço, alerta de estoque mínimo (cron).
- Fidelidade/Gamificação/Metas avançadas: XP, níveis, cashback configurável, KPIs conforme PRD; alertas e cron de métricas.
- Relatórios/export: endpoints com filtros e export CSV/Excel; dashboards por unidade/rede.
- Segurança: RBAC, multi-tenant forte, audit logs, testes de autorização e isolamento.
