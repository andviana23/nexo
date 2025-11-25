# Plano de Correção de Erros

Arquivo gerado para registrar os problemas encontrados e a ordem de correção recomendada.

---

## 📊 Status Geral: 6/6 Completos (100%) ✅

| Item                    | Status      | Descrição                               |
| ----------------------- | ----------- | --------------------------------------- |
| 1. Backend – compilar   | ✅ COMPLETO | Handlers, helpers, use cases corrigidos |
| 2. Banco de dados       | ✅ COMPLETO | Migrations + SQLC regenerado            |
| 3. Segurança e rotas    | ✅ COMPLETO | JWT + .env + Makefile limpo             |
| 4. Frontend – build     | ✅ COMPLETO | Next.js 14.2.4 + 28 routes OK           |
| 5. Testes/CI            | ✅ COMPLETO | Jest + Go tests passando                |
| 6. Configuração Next.js | ✅ COMPLETO | Webpack + lockfile + docs atualizados   |

**Última atualização:** 24/11/2025
**Status:** 🎉 TODAS AS CORREÇÕES CONCLUÍDAS!

---

## Problemas Identificados

- Backend não compila: `get_dashboard.go` usa ports/métodos inexistentes e value objects como funções.
- Repositórios Postgres com helpers ausentes ou assinaturas erradas (`compensacao_bancaria_repository.go`, `fluxo_caixa_diario_repository.go` e correlatos).
- Handlers de financeiro/estoque/LGPD estão desativados (`*.go.disabled`), mas são usados no `cmd/api/main.go`.
- Migrations reais não cobrem as tabelas/colunas exigidas pelas queries SQLC e entidades.
- Rotas financeiras/pricing expostas sem middleware JWT; JWT usa HS256 hardcoded (diverge do plano RS256).
- Credenciais de banco expostas no `Makefile`.
- Frontend não builda: imports inexistentes (`@/components/ui/button`/`card`), serviços/hooks com tipos e assinaturas quebradas (`params` em apiClient, hooks esperando arrays, services faltantes).
- Script de testes do frontend falha por ausência de testes (`npm test` sem match).
- Configuração Next.js com chaves inválidas (`swcMinify`), lockfiles duplicados e middleware depreciado.

## Checklist em Ordem

1. **Backend – compilar** ✅ COMPLETO
   - [x] Remover/ajustar dependência de `GetDashboardUseCase`: alinhar ports e métodos ou desativar o use case temporariamente.
   - [x] Corrigir helpers ausentes (`pgUUIDToString`, `timestamptzToTime`, `decimalToMoney`, `numericToMoney`, `timestampToTimestamptz`, `int32ToDMais`) ou adequar repositórios para usar os helpers existentes.
   - [x] Recriar/renomear handlers `financial_handler.go` e `stock_handler.go` (e LGPD se necessário) removendo extensão `.disabled` e garantindo que compilam.
2. **Banco de dados** ✅ COMPLETO
   - [x] Revisar migrations: criar/atualizar scripts para todas as tabelas e colunas usadas por `internal/infra/db/schema`/SQLC.
   - [x] Regenerar SQLC se o schema for atualizado.
3. **Segurança e rotas** ✅ COMPLETO
   - [x] Colocar grupos `/financial` e `/pricing` atrás do middleware JWT.
   - [x] Definir estratégia JWT (HS256 vs RS256) e mover segredo/keys para `.env`.
   - [x] Remover `DATABASE_URL` sensível do `Makefile`.
4. **Frontend – build** ✅ COMPLETO
   - [x] Criar/ajustar componentes `components/ui/button` e `components/ui/card` ou alterar imports no `cookie-consent-banner`.
   - [x] Corrigir `apiClient.request` para aceitar `params` ou retirar uso de `params` nos services/hooks; alinhar retornos (hooks que esperam array vs service que devolve item único).
   - [x] Implementar ou remover chamadas para services inexistentes (`listMovimentacoes`, `createContaPagar/Receber`, etc.).
   - [x] Ajustar tipos em hooks de metas/stock/pricing para corresponder às respostas.
   - [x] Rodar `npm run build` e resolver erros de TS restantes.
   - [x] **SOLUÇÃO:** Stack fixada em Next.js 14.2.4 + React/React DOM 18.2.0 + MUI 5.15.21 + Emotion 11.11 para maximizar compatibilidade de SSR e cache (TanStack Query 4).
   - [x] Build passa com sucesso: 28 routes compilados, 0 erros SSR.
5. **Testes/CI** ✅ COMPLETO
   - [x] Ajustar script `npm test` para `--passWithNoTests` ou adicionar testes mínimos.
   - [x] Verificar/go test após correções e adicionar smoke tests básicos.
   - [x] **Implementado:** Jest config para Next.js 14.2.4 + React 18.2.0.
   - [x] **Implementado:** Smoke tests frontend (4 testes passando).
   - [x] **Implementado:** Smoke tests backend Go (2 testes passando).
   - [x] **Corrigido:** financial_handler_integration_test.go (parâmetro dashboard placeholder).
6. **Configuração Next.js** ✅ COMPLETO
   - [x] Remover `swcMinify` inválido do `next.config.js` e decidir lockfile único (pnpm).
   - [x] Atualizar middleware para convenção recomendada (`proxy`) ou confirmar suporte na versão atual.
   - [x] Adicionar `outputFileTracingRoot` para eliminar warnings de lockfile.
   - [x] Configurar webpack aliases para React single instance.
   - [x] Atualizar toda documentação de Next.js 16.0.3 → 14.2.4 (21 arquivos).
