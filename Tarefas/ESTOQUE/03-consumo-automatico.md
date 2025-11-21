# 3. Consumo Automático por Serviço

- **Categoria:** ESTOQUE
- **Objetivo:** Baixar automaticamente do estoque os insumos vinculados a um serviço quando este é realizado/finalizado.
- **Escopo:** Backend (Triggers/Events), Banco de Dados.

## Plano de Execução (prioridade 3)
- **Banco de Dados:** tabela `servico_produtos` (ficha técnica) e movimentações `SAIDA` com motivo `CONSUMO_SERVICO`; garantir índices por serviço/produto.
- **Backend:** listener em evento de finalização de atendimento/serviço; baixa automática conforme ficha técnica; política de saldo insuficiente (alerta ou negativo).
- **Frontend:** UI para configurar ficha técnica (serviço ↔ insumos/quantidades) e visualizar consumo.
- **Cálculos aplicados:** custo de insumos influencia precificação de serviço; referência `docs/10-calculos/custo-insumo-servico.md` para rastrear custo unitário consumido.

## Fluxo Operacional

1. Configuração prévia: Vincular Produtos (Insumos) a Serviços com quantidades padrão (Ficha Técnica).
   - Ex: Serviço "Barba" consome "10ml Creme Barbear" + "1 Lâmina".
2. Barbeiro finaliza um atendimento/serviço.
3. Sistema detecta conclusão do serviço.
4. Sistema consulta ficha técnica do serviço.
5. Sistema gera movimentações de `SAIDA` (Motivo: CONSUMO_SERVICO) para cada insumo.
6. Saldo dos produtos é atualizado.

## Campos Obrigatórios

- `servico_id` (UUID)
- `agendamento_id` / `atendimento_id` (UUID) - Origem do consumo.

## Comportamentos Esperados

- Baixa automática transparente ao usuário no momento da finalização.
- Se saldo insuficiente, registrar a baixa mesmo assim (ficando negativo) ou gerar alerta (configurável).
- Logar qual atendimento gerou o consumo.

## Regras de Baixa

- Decrementa `quantidade_atual` baseado na `ficha_tecnica` do serviço.

## Logs e Auditoria

- Registrar em `audit_logs`:
  - `action`: `CREATE`
  - `resource`: `movimentacao_estoque`
  - `details`: `{ tipo: "SAIDA", motivo: "CONSUMO_SERVICO", servico_id: "...", atendimento_id: "..." }`

## Dependências

- Módulo de Agendamento/Atendimento (Gatilho).
- Cadastro de Serviços (Ficha Técnica: tabela `servico_produtos`).
- Cadastro de Produtos.

## Tarefas

1. Criar tabela `servico_produtos` (servico_id, produto_id, quantidade).
2. Criar endpoint/UI para configurar ficha técnica de serviços.
3. Implementar listener/observer no evento `ServiceCompleted` ou `AtendimentoFinalizado`.
4. Lógica de baixa automática iterando sobre insumos do serviço.
5. Testes de integração: Finalizar serviço -> Verificar se estoque diminuiu.

## Critérios de Aceite

- Ao finalizar serviço com ficha técnica, estoque dos insumos é reduzido.
- Movimentações geradas têm vínculo com o atendimento original.
- Configuração de ficha técnica permite definir quantidades fracionadas (se unidade permitir).
