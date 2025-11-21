# 1. Entrada de Estoque

- **Categoria:** ESTOQUE
- **Objetivo:** Registrar a entrada de produtos no estoque, vinculando a fornecedores e atualizando quantidades e custos.
- **Escopo:** Backend (API, validações), Frontend (Formulário de entrada), Banco de Dados (Movimentações, Produtos).

## Plano de Execução (prioridade 1)
- **Banco de Dados:** tabelas de produtos e movimentações; garantir índices por tenant/produto/data; registrar movimentação `ENTRADA`.
- **Backend:** endpoint `POST /stock/entries`, validação de fornecedor/produto do tenant, atualização de saldo, movimentações e (opcional) geração de conta a pagar.
- **Frontend:** formulário com seleção de fornecedor e itens (SKU/nome), cálculo de valor total, confirmação.
- **Cálculos aplicados:** recalcular custo médio se adotado (não há fórmula específica na doc de cálculos; sem dependência direta).

## Fluxo Operacional

1. Usuário acessa tela de "Entrada de Estoque".
2. Seleciona o Fornecedor (ou cadastra um novo).
3. Adiciona produtos à lista de entrada (busca por SKU ou Nome).
4. Informa quantidade, valor unitário de compra e data da entrada.
5. Sistema calcula valor total da entrada.
6. Usuário confirma a operação.
7. Sistema atualiza quantidade do produto (`quantidade_atual += qtd`) e registra movimentação `ENTRADA`.
8. Opcional: Sistema gera conta a pagar no módulo financeiro.

## Campos Obrigatórios

- `fornecedor_id` (UUID)
- `data_entrada` (Date)
- Lista de Itens:
  - `produto_id` (UUID)
  - `quantidade` (Int > 0)
  - `valor_unitario` (Decimal > 0)

## Comportamentos Esperados

- Atualização imediata da quantidade em estoque.
- Recálculo do preço médio de custo do produto (se aplicável).
- Bloqueio de entrada com data futura (configurável).
- Validação de fornecedor ativo.

## Regras de Baixa

- Não se aplica (é uma operação de adição).

## Logs e Auditoria

- Registrar em `audit_logs`:
  - `action`: `CREATE`
  - `resource`: `movimentacao_estoque`
  - `details`: `{ tipo: "ENTRADA", fornecedor_id: "...", itens: [...] }`

## Dependências

- Módulo Financeiro (opcional): Integração para gerar `contas_a_pagar` automaticamente.
- Cadastro de Fornecedores (`fornecedores`).
- Cadastro de Produtos (`produtos`).

## Tarefas

1. Criar endpoint `POST /stock/entries` para registrar entrada em lote.
2. Implementar validação de fornecedor e produtos (pertencem ao tenant).
3. Atualizar tabela `produtos` incrementando quantidade.
4. Registrar movimentações individuais na tabela `movimentacoes` com tipo `ENTRADA`.
5. (Opcional) Criar integração com `CreatePayableUseCase` se flag `gerar_financeiro` for true.
6. Desenvolver interface de entrada com busca de produtos e fornecedores.
7. Testes unitários e de integração cobrindo atualização de saldo e criação de movimentações.

## Critérios de Aceite

- Entrada atualiza corretamente o saldo de todos os produtos listados.
- Movimentações são criadas com tipo `ENTRADA` e vínculo com fornecedor.
- Tentativa de entrada com produto inexistente ou quantidade <= 0 retorna erro.
- Integração financeira (se ativa) cria conta a pagar com valor correto.
