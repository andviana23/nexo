# Checklist de Ajustes do Módulo Financeiro (NEXO)

> Este checklist consolida **todos os ajustes** identificados na análise do módulo financeiro.  
> A ordem prioriza risco financeiro, segurança, atomicidade e multi‑tenancy/multi‑unidade.

## Fase 0 — Alinhamento e critérios
- [x] Ler/confirmar PRD e fluxos financeiros aplicáveis (ex.: `docs/07-produto-e-funcionalidades/PRD-NEXO.md`, `docs/11-Fluxos/FLUXO_*`).
- [x] Definir lifecycle canônico de status para:
  - Contas a pagar (hoje: domínio `PENDENTE/PAGO` vs banco `ABERTO/PAGO`).
  - Contas a receber (domínio `PAGO` vs banco `RECEBIDO/CONFIRMADO/ESTORNADO/...`).
  - Compensações (domínio/banco `PREVISTO/CONFIRMADO/COMPENSADO/CANCELADO`).
- [x] Definir estratégia para tabelas legadas `receitas` e `assinaturas` (manter com adaptadores explícitos ou migrar/descontinuar).

## Fase 1 — Correções críticas (bloqueantes para consistência financeira)
1. [x] Corrigir atualização de totais do caixa via webhook Asaas  
   - Onde: `backend/internal/application/usecase/subscription/process_webhook_usecase_v2.go`.  
   - Ação: não zerar `total_sangrias/total_reforcos` ao chamar `UpdateTotais`; preservar valores atuais ou criar método incremental.
2. [x] Normalizar status entre domínio ⇄ sqlc ⇄ schema ⇄ DTO  
   - Onde: `backend/internal/domain/valueobject/enums.go`, schemas `backend/internal/infra/db/schema/contas_a_pagar.sql`, `backend/internal/infra/db/schema/contas_a_receber.sql`, DTOs `backend/internal/application/dto/financial_dto.go`, mappers/repositories.  
   - Ação: escolher nomes canônicos e mapear; garantir que use cases gravem status válidos no banco.
3. [x] Garantir que TODA finalização de comanda gere lançamento financeiro completo  
   - Onde: `backend/internal/application/usecase/command/finalizar_comanda_integrada.go`.  
   - Ações:
     - Criar `ContaReceber` também para pagamentos à vista/PIX (marcar como recebida no ato) com origem correta (`SERVICO`/`PRODUTO`).  
     - Tornar criação idempotente para evitar duplicidade em reprocessamentos.  
     - Persistir vínculo explícito com a comanda/pagamento para rastreabilidade.
4. [x] Implementar estorno financeiro no cancelamento de comanda fechada  
   - Onde: `backend/internal/application/usecase/command/cancel_command.go`.  
   - Ações:
     - Localizar e cancelar/estornar `ContaReceber` vinculada(s).  
     - Criar operação inversa no `operacoes_caixa` ou marcar operações como estornadas.  
     - Recalcular/ajustar totais do caixa.
5. [x] Corrigir DRE por origem (evitar receita “tripla”)  
   - Onde: sqlc `backend/internal/infra/db/queries/contas_a_receber.sql`, repo `backend/internal/infra/repository/postgres/conta_receber_repository.go`, `backend/internal/application/usecase/financial/generate_dre.go`.  
   - Ações:
     - Criar query real `SumContasReceberByOrigem` (origem + período + tenant).  
     - Ajustar `SumByOrigem` para usar a query.  
     - Revalidar cron `GenerateDREMonthly`.
6. [x] Proteger rotas de despesas fixas com RBAC  
   - Onde: `backend/cmd/api/main.go` e/ou `backend/internal/infra/http/handler/despesa_fixa_handler.go`.  
   - Ação: aplicar `RequireOwnerOrManager`/`RequireAdminAccess` em todas as rotas `/financial/fixed-expenses/*`.

## Fase 2 — Completar fluxos financeiros V2
7. [x] Migrar cron diário para `GenerateFluxoDiarioV2`  
   - Onde: `backend/internal/infra/scheduler/jobs.go`.  
   - Ação: usar `received_at` para assinaturas Asaas e manter compatibilidade com contas tradicionais.
8. [x] Migrar cron mensal para `GenerateDREV2` com regime explícito  
   - Onde: `backend/internal/infra/scheduler/jobs.go`.  
   - Ação: calcular por **COMPETÊNCIA** (`confirmed_at/competencia_mes`) e suportar **CAIXA** quando necessário.
9. [x] Integrar compensações bancárias ao fluxo real  
   - Onde: `finalizar_comanda_integrada.go`, use cases `backend/internal/application/usecase/financial/create_compensacao.go` e `marcar_compensacao.go`, rotas financeiras.  
   - Ações:
     - Gerar `CompensacaoBancaria` automaticamente quando meio de pagamento tiver D+.  
     - Expor endpoints de criação/compensar (ou operar só por cron).  
     - Garantir que o fluxo diário inclua valores previstos/confirmados via compensação.
10. [x] Implementar filtros reais em listagens financeiras  
    - Onde: repositórios `conta_pagar_repository.go`, `conta_receber_repository.go`, handlers de listagem.  
    - Ação: suportar filtros de status/origem/unidade/data conforme DTO.

## Fase 3 — Multi‑unidade real
11. [ ] Introduzir `unit_id` em todas as entidades/tabelas financeiras relevantes  
    - Onde: schemas `contas_a_receber`, `fluxo_caixa_diario`, `dre_mensal`, `caixa_diario`, `operacoes_caixa`.  
    - Ação: migrations + atualização de entidades/DTOs/mappers/repositories.
12. [ ] Propagar `unit_id` desde a comanda/assinatura/estoque  
    - Onde: `Command`, `FinalizarComandaIntegradaUseCase`, webhooks Asaas, estoque.  
    - Ação: ler unidade do contexto (claim/header) e persistir em lançamentos.
13. [ ] Aplicar `UnitMiddleware` nas rotas financeiras e de caixa  
    - Onde: `backend/cmd/api/main.go`.  
    - Ação: bloquear requests sem unidade e filtrar queries por unidade.

## Fase 4 — Integrações adjacentes
14. [ ] Integrar entrada de estoque ao financeiro  
    - Onde: `backend/internal/application/usecase/stock/registrar_entrada.go`.  
    - Ação: quando `GerarFinanceiro=true`, criar `ContaPagar` automática vinculada à entrada.
15. [ ] Desacoplar geração de comissões do fechamento de comanda  
    - Onde: `backend/internal/application/usecase/command/finalizar_comanda_integrada.go` e módulo de comissões.  
    - Ação: mover busca/cálculo para use case de comissões e chamar via serviço/evento.
16. [ ] Considerar assinaturas no cálculo de comissão e origem de receita  
    - Onde: fechamento de comanda e regras de comissão.  
    - Ação: diferenciar serviço “consumido do plano” vs pago avulso para comissão correta.

## Fase 5 — Dívida técnica e alinhamento de contrato
17. [ ] Atualizar Swagger para refletir rotas reais e parâmetros corretos  
    - Onde: `backend/internal/infra/http/handler/financial_handler.go`.
18. [ ] Padronizar domínio usando apenas `domain/port` (ou separar claramente)  
    - Onde: interfaces de comissões vs demais módulos.
19. [ ] Resolver convivência com tabelas legadas  
    - Ação: remover referências a `receitas`/`assinaturas` ou criar adaptadores explícitos.
20. [ ] Alinhar tipos/enums do frontend com backend  
    - Onde: `frontend/src/types/financial.ts`, services/hooks.  
    - Ação: atualizar status/origem/estrutura de DRE para o contrato real.

## Fase 6 — Evolução estrutural (longo prazo)
21. [ ] Centralizar reconhecimento de receitas/despesas no módulo financeiro via eventos de domínio  
    - Ação: comanda/assinatura/estoque/comissão publicam eventos; financeiro consome de forma idempotente.
22. [ ] Implementar outbox/transações para atomicidade entre módulos  
    - Ação: fechamento de comanda vira transação única do ponto de vista financeiro.
23. [ ] Adicionar audit trail completo em lançamentos financeiros  
    - Ação: `user_id`, `origem`, entidade relacionada, timestamps e logs estruturados.
