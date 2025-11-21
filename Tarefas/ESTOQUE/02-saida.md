# 2. Saída de Estoque

- **Categoria:** ESTOQUE
- **Objetivo:** Registrar a saída manual de produtos (venda direta, uso interno, perda, avaria), decrementando o saldo.
- **Escopo:** Backend, Frontend, Banco de Dados.

## Plano de Execução (prioridade 2)
- **Banco de Dados:** movimentações `SAIDA` com motivo; locks ou transações para evitar corrida de saldo.
- **Backend:** endpoint `POST /stock/exits`, validação de saldo (com política de negativo opcional), atualização de saldo e registro de motivo; integração opcional com receitas.
- **Frontend:** tela de saída com motivo, seleção de itens, validação de saldo e feedback.
- **Cálculos aplicados:** nenhum cálculo financeiro direto; movimentações alimentam consumo/curva ABC e saldo mínimo.

## Fluxo Operacional

1. Usuário acessa tela de "Saída de Estoque".
2. Seleciona o motivo da saída (Venda, Uso Interno, Perda, Validade, Avaria).
3. Se "Uso Interno", pode vincular a um Profissional/Barbeiro.
4. Adiciona produtos e quantidades.
5. Confirma operação.
6. Sistema valida saldo disponível.
7. Sistema atualiza quantidade (`quantidade_atual -= qtd`) e registra movimentação `SAIDA`.

## Campos Obrigatórios

- `motivo` (Enum: VENDA, USO_INTERNO, PERDA, AVARIA, VALIDADE)
- `data_saida` (Date)
- Lista de Itens:
  - `produto_id` (UUID)
  - `quantidade` (Int > 0)

## Comportamentos Esperados

- Bloqueio de saída se saldo insuficiente (configurável: permitir saldo negativo ou não).
- Validação de motivo.
- Se motivo for VENDA, integração opcional com módulo de Receitas.

## Regras de Baixa

- Decrementa `quantidade_atual` na tabela `produtos`.
- Valida `quantidade_atual >= quantidade` antes de efetivar (exceto se tenant permitir negativo).

## Logs e Auditoria

- Registrar em `audit_logs`:
  - `action`: `CREATE`
  - `resource`: `movimentacao_estoque`
  - `details`: `{ tipo: "SAIDA", motivo: "...", itens: [...] }`

## Dependências

- Cadastro de Produtos.
- Cadastro de Usuários (para vincular profissional em uso interno).
- Módulo Financeiro (se venda gerar receita).

## Tarefas

1. Criar endpoint `POST /stock/exits`.
2. Implementar verificação de saldo disponível (lock otimista ou transação).
3. Registrar movimentações com tipo `SAIDA` e subtipo/motivo.
4. Atualizar saldo em `produtos`.
5. Desenvolver UI de saída com seleção de motivo.
6. Testes de concorrência (duas saídas simultâneas para mesmo produto com pouco saldo).

## Critérios de Aceite

- Saída decrementa saldo corretamente.
- Sistema impede saída maior que saldo (se configurado).
- Motivo da saída é registrado corretamente na movimentação.
