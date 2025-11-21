# 4. Metas Automáticas Baseadas no Faturamento Mínimo

- **Categoria:** METAS
- **Objetivo:** Gerar automaticamente metas mensais baseadas no faturamento mínimo necessário para cobrir despesas fixas + margem de lucro desejada.
- **Escopo:** Backend (engine de cálculo), Frontend (visualização e configuração).

## Plano de Execução (prioridade 4 em Metas)
- **Banco de Dados:** campo `margem_lucro_desejada` em `configuracoes_tenant`; `metas_mensais` com campos `origem` e `status`; índices por tenant/mes.
- **Backend:** job mensal (dia 1) que calcula meta sugerida; endpoints para configurar margem, obter sugerida, aceitar/sobrescrever; auto-aceitar dia 5 se não agir.
- **Frontend:** card “Meta Sugerida” com breakdown (despesas fixas + margem), botões aceitar/definir manual.
- **Cálculos aplicados:** Meta automática usa fórmula (Despesas Fixas / (1 - Margem)) equivalente a Faturamento Mínimo; referenciar `docs/10-calculos/faturamento-minimo-mensal.md` e, opcionalmente, Ponto de Equilíbrio (`ponto-de-equilibrio.md`) se comparar metas.

## Cálculo da Meta

- Sistema calcula meta automática com base em:
  1. **Despesas Fixas do Mês**: Soma de todas as contas a pagar recorrentes (aluguel, internet, etc).
  2. **Margem de Lucro Desejada**: Percentual configurado pelo tenant (ex: 30%).
  3. **Fórmula**: `Meta Automática = Despesas Fixas / (1 - Margem Lucro)`
     - Exemplo: Despesas = R$ 10.000, Margem = 30% → Meta = R$ 10.000 / 0.7 = R$ 14.285,71

## Atualização Automática

- Recalculada no início de cada mês (Cron no dia 1º às 00:01).
- Considera despesas fixas cadastradas como recorrentes.
- Pode ser sobrescrita manualmente pelo Owner/Manager.

## Painel Visual

- Card "Meta Sugerida" no Dashboard.
- Exibição:
  - Meta Automática Calculada (R$)
  - Breakdown: Despesas Fixas + Margem
  - Botão "Aceitar Meta" ou "Definir Manualmente"
  - Status (se meta aceita ou pendente)

## Alertas e Status

- 🔵 **Azul (Pendente)**: Meta calculada mas não aceita/confirmada.
- 🟢 **Verde**: Meta aceita e em uso.
- ⚙️ **Manual**: Owner definiu meta diferente da sugerida.

## Regras

- RN-META-016: Meta automática só é gerada se existirem despesas fixas cadastradas.
- RN-META-017: Margem de lucro deve estar entre 5% e 100%.
- RN-META-018: Meta automática serve como sugestão; Owner pode aceitar, ajustar ou ignorar.
- RN-META-019: Se Owner não aceitar até dia 5, sistema usa meta automática por padrão.
- RN-META-020: Cálculo considera apenas despesas com flag `tipo = FIXA` e `recorrente = true`.

## Dependências

- Módulo Financeiro: Contas a Pagar (despesas fixas).
- Tabela `configuracoes_tenant`: campo `margem_lucro_desejada`.
- Tabela `metas_mensais`: adicionar campo `origem` (MANUAL | AUTOMATICA).
- Cron scheduler para geração no dia 1º.

## Tarefas

1. Adicionar campo `margem_lucro_desejada` na tabela `configuracoes_tenant` (default 30%).
2. Criar endpoint `POST /settings/profit-margin` para configurar margem.
3. Implementar Job `GenerateAutoGoalsJob` executado no dia 1º de cada mês:
   - Calcular despesas fixas do mês.
   - Aplicar fórmula de meta automática.
   - Criar registro em `metas_mensais` com `origem = AUTOMATICA` e `status = PENDENTE`.
4. Criar endpoint `GET /goals/suggested` retornando meta sugerida + breakdown.
5. Criar endpoint `POST /goals/suggested/accept` para aceitar meta automática.
6. Desenvolver UI: card de meta sugerida com botões de ação.
7. Adicionar lógica: se não aceita até dia 5, auto-aceita meta sugerida.
8. Testes: cálculo correto com diferentes margens, cenários sem despesas fixas.

## Critérios de Aceite

- Sistema gera meta automática no dia 1º de cada mês baseada em despesas fixas.
- Fórmula de cálculo está correta (Despesas / (1 - Margem)).
- Owner pode aceitar, ajustar ou ignorar meta sugerida.
- Se não aceita até dia 5, meta automática é aplicada.
- Dashboard exibe breakdown (despesas + margem).
- Margem de lucro configurável entre 5% e 100%.
- Testes cobrem cenários: sem despesas, margens variadas, aceitação/rejeição.
