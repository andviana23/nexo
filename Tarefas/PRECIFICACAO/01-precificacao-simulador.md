# CATEGORIA: PRECIFICAÇÃO — Simulador Completo

## Funcionalidade: Calculadora de Preço de Venda (Serviços e Produtos)

- **Objetivo:** Determinar preço de venda ideal considerando custo de insumos, comissão do barbeiro, impostos, markup e margem desejada.

## Plano de Execução (prioridade 1 em Precificação)
- **Banco de Dados:** `precificacao_config` (defaults por tenant) e `precificacao_simulacoes` (histórico de cenários); índices por tenant/item/data. Consumir custo médio de estoque e configurações de comissão/impostos.
- **Backend:** use cases de cálculo e persistência de config/simulações; endpoints `/pricing/config`, `/pricing/simulate`, `/pricing/simulations` e API pública `/api/v1/pricing/calculate` com rate limit e auditoria.
- **Frontend:** telas de configuração e simulador interativo (cenários, waterfall de composição, comparação e exportação).
- **Cálculos aplicados (docs/10-calculos):** Preço de Produto (`preco-produto.md`), Preço de Serviço (`preco-servico.md`), Margem de Lucro (`margem-lucro.md`), Markup (`markup.md`), Custo de Insumo por Serviço (`custo-insumo-servico.md`). Ticket Médio/LTV/CAC são opcionais para análises, não entram no cálculo direto.

### Tarefas

1. **Modelagem de Dados**
   - Criar tabela `precificacao_config` (tenant_id, tipo_item [SERVICO|PRODUTO], margem_desejada, markup_alvo, imposto_percentual, comissao_percentual_default, criado_em, atualizado_em).
   - Criar tabela `precificacao_simulacoes` para histórico (tenant_id, item_id, tipo_item, custo_insumos, comissao, impostos, markup, margem_resultante, preco_sugerido, parametros_json, criado_por, criado_em).
2. **Backend / Use Cases**
   - Implementar `CalculatePrecoVendaUseCase` recebendo inputs (item_id, tipo_item, custo_insumos, comissao_percentual, impostos_percentual, margem_desejada, markup_custom, descontos).
   - Implementar `SavePrecificacaoConfigUseCase` para defaults por tenant.
   - Expor endpoints:
     - `GET /pricing/config` (buscar defaults)
     - `PUT /pricing/config` (atualizar regras)
     - `POST /pricing/simulate` (retornar cálculo e salvar histórico opcional)
     - `GET /pricing/simulations?item_id=&limit=` (listar simulações anteriores)
3. **Frontend / UI**
   - Tela "Configuração de Precificação" com formulário (margem desejada, markup alvo, impostos padrão, comissões padrão) usando RHF + Zod.
   - Tela "Simulador" com inputs dinâmicos por item (auto preencher custo de insumos via estoque, comissões via módulo de comissões, impostos via config) e exibir resultado em cards.
   - Gráfico mini waterfall demonstrando composição do preço (custo → impostos → comissão → margem → preço final).
4. **Integrações**
   - Buscar custo médio dos insumos via módulo `ESTOQUE` (consumo automático por serviço).
   - Buscar comissão padrão do barbeiro via módulo `Comissões Automáticas`.
   - Aplicar impostos configurados por tenant (ex: ISS, ICMS) — armazenar em `precificacao_config`.
5. **Validações e Logs**
   - Logar todas as simulações no `audit_logs` (acao: PRECIFICACAO_SIMULADA, payload: parâmetros e resultado).
   - Garantir multi-tenant em todas as queries (tenant_id obrigatório).
6. **Testes**
   - Unit tests para fórmula principal (Go) cobrindo combinações de markup/margem.
   - Integração: endpoint `/pricing/simulate` retornando os campos corretos.
   - Frontend: testes de formulário (Zod) e snapshot do componente de resultado.

### Fórmulas Matemáticas

- `CustoTotal = CustoInsumos + (CustoInsumos * ImpostoPercentual/100)` _(impostos que incidem sobre custo)_
- `ValorComissao = PrecoVenda * ComissaoPercentual/100`
- `Markup = PrecoVenda / CustoInsumos`
- `Margem = (PrecoVenda - (CustoInsumos + ValorImpostos + ValorComissao)) / PrecoVenda`
- **Preço alvo por margem desejada:**
  - `PrecoVenda = (CustoInsumos + ValorImpostos) / (1 - (ComissaoPercentual/100) - (MargemDesejada/100))`
- **Preço alvo por markup desejado:**
  - `PrecoVenda = CustoInsumos * MarkupAlvo`
- **Preço final sugerido:**
  - `PrecoFinal = max(PrecoPorMargem, PrecoPorMarkup)` arredondado para múltiplo configurável (ex.: R$ 5,00).

### Regras

- ✅ RN-PREC-001: Custos de insumos devem vir do estoque (custo médio periódico) para serviços; para produtos usar custo unitário cadastrado.
- ✅ RN-PREC-002: Comissão pode ser fixa ou percentual; se múltiplos barbeiros, usar o maior percentual para garantir margem mínima.
- ✅ RN-PREC-003: Impostos devem permitir configuração por tipo (ISS, ICMS, Simples Nacional) com opção de repasse integral.
- ✅ RN-PREC-004: Markup mínimo default = 1.5; margem desejada default = 30% (editável).
- ✅ RN-PREC-005: O simulador deve permitir comparação de cenários (exibir 3 preços: mínimo, recomendado, premium).
- ✅ RN-PREC-006: Qualquer alteração na configuração precisa de permissão OWNER/MANAGER e gera log em audit trail.
- ✅ RN-PREC-007: API deve bloquear cálculo se inputs obrigatórios ausentes (custo <= 0, margem < 0, markup < 1).
- ✅ RN-PREC-008: Resultado precisa indicar se margem desejada foi atingida; caso não, sugerir ajuste.

### Dependências

- **Estoque:** custo médio dos insumos por serviço/produto (`ESTOQUE/03-consumo-automatico.md`).
- **Comissões Automáticas:** percentual por barbeiro/serviço (`FINANCEIRO/05-comissoes-automaticas.md`).
- **Financeiro:** categorias e impostos cadastrados (`FINANCEIRO/03-contas-a-pagar.md`, `FINANCEIRO/04-contas-a-receber.md`).
- **Metas:** margem desejada pode usar metas financeiras (`METAS/01-meta-geral-mes.md`).
- **RBAC:** somente OWNER/MANAGER podem alterar configurações (ver `docs/06-seguranca/RBAC.md`).

---

## Funcionalidade: Simulador Interativo com Comparação de Cenários

- **Objetivo:** Permitir ao gestor comparar rapidamente variações de preço trocando inputs (margem, comissões, impostos) e visualizar impacto.

### Tarefas

1. **Backend**
   - Endpoint `POST /pricing/scenarios` aceitando array de cenários (até 5) e retornando resultados agregados (preço final, margem, markup, lucro por serviço).
   - Calcular lucro líquido por serviço e por mês (considerando volume informado).
2. **Frontend**
   - UI com tabela comparativa (linhas: custo insumos, comissão, impostos, margem, preço final, lucro por serviço, lucro mensal).
   - Inputs controlados (sliders para margem/markup, campos numéricos para comissões e impostos).
   - Destaque visual do melhor cenário (maior margem sustentada).
3. **Exportação**
   - Botão "Exportar Cenários" gerando PDF/CSV com parâmetros e resultados.
4. **Alertas**
   - Mostrar avisos quando margem < margem mínima configurada ou quando preço sugerido < custo total.

### Fórmulas Matemáticas

- Para cada cenário aplicar as fórmulas da funcionalidade anterior.
- `LucroPorServico = PrecoFinal - (CustoInsumos + ValorComissao + ValorImpostos)`
- `LucroMensal = LucroPorServico * VolumeEstimado`

### Regras

- Limitar a 5 cenários simultâneos para manter performance.
- Permitir salvar "favoritos" (simulações recorrentes) com nome e tags.
- Se volume estimado não informado, usar volume histórico médio (via módulo financeiro).

### Dependências

- Histórico de vendas (para volume médio) via `receitas`.
- Geração de PDF (já introduzido no DRE) para exportar relatórios.

---

## Funcionalidade: API Pública de Precificação

- **Objetivo:** Permitir integrações externas (ex.: app mobile) consultarem preço recomendado via API segura.

### Tarefas

1. Criar endpoint `POST /api/v1/pricing/calculate` (token JWT + scope específico) retornando preço sugerido.
2. Implementar rate limit e auditoria (exigir API Key + tenant).
3. Documentar no `API_REFERENCE.md` com exemplos.
4. Adicionar testes de contrato (Pact) garantindo consistência.

### Fórmulas Matemáticas

- Mesmas da calculadora principal.

### Regras

- Apenas receitas com tenant válido; `tenant_id` vem do token.
- Inputs obrigatórios: item_id ou custo_insumos manual, comissões, impostos ou usar defaults.
- Fail-fast se tenant não configurou precificação (retornar 409 "configuração pendente").

### Dependências

- Auth/JWT, RBAC, audit logs, módulo financeiro para custos default.
