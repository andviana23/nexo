# 4. Inventário (Ajuste de Estoque)

- **Categoria:** ESTOQUE
- **Objetivo:** Permitir a contagem física e ajuste do saldo sistêmico para refletir a realidade (correção de divergências).
- **Escopo:** Backend, Frontend.

## Plano de Execução (prioridade 4)
- **Banco de Dados:** tabelas `inventarios` (cabeçalho) e `inventario_itens` (detalhe); movimentações de ajuste (ENTRADA/SAIDA com motivo AJUSTE_INVENTARIO).
- **Backend:** endpoint `POST /stock/inventories`, cálculo de divergência e geração das movimentações de ajuste; opção de bloquear movimentações durante inventário.
- **Frontend:** UI de contagem (modo grade), filtros por categoria, relatório de divergências.
- **Cálculos aplicados:** nenhum cálculo financeiro direto; ajustes impactam saldos usados por curva ABC e estoque mínimo.

## Fluxo Operacional

1. Usuário inicia "Novo Inventário".
2. Sistema lista produtos (pode filtrar por categoria).
3. Usuário informa a quantidade contada fisicamente para cada produto.
4. Sistema calcula a diferença (Físico - Sistêmico).
5. Usuário finaliza inventário.
6. Sistema gera movimentações de `AJUSTE` (entrada ou saída) para igualar o saldo.

## Campos Obrigatórios

- `data_inventario`
- Lista de Contagem:
  - `produto_id`
  - `quantidade_contada`

## Comportamentos Esperados

- Bloquear movimentações de estoque durante o inventário (opcional, mas recomendado).
- Registrar quem realizou a contagem.
- Gerar movimentação de ajuste apenas se houver divergência.

## Regras de Baixa/Entrada

- Se Contada < Sistêmica: Gera `SAIDA` (Motivo: AJUSTE_INVENTARIO).
- Se Contada > Sistêmica: Gera `ENTRADA` (Motivo: AJUSTE_INVENTARIO).

## Logs e Auditoria

- Registrar em `audit_logs`:
  - `action`: `CREATE`
  - `resource`: `inventario`
  - `details`: `{ divergencias: [...] }`

## Dependências

- Cadastro de Produtos.

## Tarefas

1. Criar tabela `inventarios` (cabeçalho) e `inventario_itens` (detalhes).
2. Criar endpoint `POST /stock/inventories` para processar a contagem.
3. Lógica de cálculo de divergência e geração de movimentações de ajuste.
4. UI de Inventário (modo grade para digitação rápida).
5. Relatório de divergências de inventário.

## Critérios de Aceite

- Saldo final do produto deve ser igual à quantidade contada.
- Movimentações de ajuste são criadas corretamente para as diferenças.
- Histórico de inventários fica salvo para consulta.
