# Analise Financeiro do Fluxo Agenda → Comanda → Financeiro → Caixa

## Mapa do fluxo real (caminhos observados no codigo)
1) **Agendamento** (`internal/application/usecase/appointment/create_appointment.go`)
   - Valida tenant, profissional, cliente e lista de servicos.
   - Calcula `TotalPrice` e horarios usando os servicos; nenhum vinculo com comanda ou caixa.

2) **Finalizar atendimento** (`appointment_handler.go` → `FinishServiceWithCommandUseCase`)
   - Transiciona status para `AWAITING_PAYMENT` e **cria a comanda** apenas quando a rota `/appointments/{id}/finish` e usada.
   - Comanda gerada com itens dos servicos do agendamento (quantidade fixa = 1, sem profissional/tempo/comissao).
   - Outros caminhos (`/appointments/{id}/complete` ou `UpdateAppointmentStatusUseCase`) permitem marcar `DONE` **sem comanda nem pagamento**.

3) **Comanda** (`command` use cases + mapper)
   - `Create/AddItem`: aceita itens `SERVICO/PRODUTO/PACOTE`, grava apenas `item_id`, descricao e precos em `float64`; nao valida existencia do item nem vincula profissional ou unidade.
   - `AddCommandPayment`: cria `command_payments` calculando liquido a partir das taxas recebidas por parametro; o handler define `taxaPercentual=0` e `taxaFixa=0` (TODO nao implementado). Nenhum impacto em caixa ou financeiro.
   - `CloseCommand`: verifica itens e, se `DeixarSaldoDivida` for falso, exige `TotalRecebido >= Total`. Apenas atualiza status da comanda e, opcionalmente, coloca o agendamento como `DONE`; **nao gera movimentacao financeira, nao registra no caixa, nao processa comissao**.

4) **Financeiro (contas/dre/fluxo)**
   - Use cases de `contas_a_receber/pagar`, DRE e Fluxo de Caixa operam sobre suas proprias entidades e **nao sao chamados em nenhum ponto do ciclo da comanda**. Receita de servicos nao entra no fluxo financeiro oficial.

5) **Caixa Diario**
   - Abertura/fechamento/sangria/reforco funcionam via `caixa` use cases.
   - O tipo `OperacaoCaixa` possui `Tipo=VENDA`, mas nenhuma rota ou use case o cria a partir de pagamentos de comanda; `CaixaDiario.RegistrarEntrada` nunca e chamado. Logo, vendas nao incrementam `TotalEntradas` nem geram rastreabilidade de cliente/profissional.

6) **Comissao**
   - Entidades e regras existem, mas `CreateCommissionItemUseCase` so e exposto por handler proprio; **nenhum gatilho** em add/close da comanda. Descontos da comanda nao recalculam comissao.

## Checklist de conformidade
| Item avaliado | Resultado | Evidencia
| --- | --- | --- |
| Criacao/vinculo de comanda ao finalizar agendamento | ❌ | Somente rota `/finish` cria; `/complete` e `UpdateAppointmentStatus` permitem `DONE` sem comanda.
| Integridade cliente → profissional → horario → servico | ❌ | Comanda nao guarda profissional nem horario; itens aceitam qualquer `item_id` sem validacao.
| Registro de servicos/produtos/tempo/comissao na comanda | ❌ | `command_items` nao armazenam profissional, duracao ou regra de comissao.
| Validacao de meios de pagamento e taxas | ❌ | Handler de pagamento usa taxas 0 (TODO) e nao carrega `MeioPagamento`.
| Finalizacao gera movimentacao financeira obrigatoria | ❌ | `CloseCommandUseCase` encerra sem criar `conta_receber`/`mov_fin`.
| Registro automatico no Caixa Diario com origem COMANDA | ❌ | Nenhum uso de `CreateOperacao` com `Tipo VENDA` ao pagar/comandar.
| Comissao recalculada considerando descontos | ❌ | Nenhum ponto cria ou ajusta `commission_items` a partir da comanda.
| Rastreabilidade de quem finalizou/recebeu e onde caiu o valor | ❌ | `FechadoPor` e `CriadoPor` existem, mas nao ha log/auditoria nem vinculo com caixa/financeiro.

## Falhas criticas e pontuacao de risco
| Risco (0-10) | Problema | Arquivo / metodo | Falha logica | Correcao Clean Architecture sugerida |
| --- | --- | --- | --- | --- |
| 10 | Fechamento de comanda nao gera movimentacao financeira nem caixa | `internal/application/usecase/command/close_command.go` `Execute` | Receita permanece apenas em `commands`/`command_payments`; DRE/fluxo/caixa ficam vazios | Criar `FinalizarComandaUseCase` orquestrando pagamento → `conta_receber`/`operacao_caixa` em transacao, depois fechar comanda e atualizar agendamento.
| 9 | Pagamentos ignoram taxas e configuracao do meio de pagamento | `internal/infra/http/handler/command_handler.go` `AddCommandPayment` (TODO) | Liquido/fees incorretos e ausencia de D+ para conciliacao | Carregar `MeioPagamento` pelo ID, aplicar `Taxa/TaxaFixa/DMais`, persistir liquido e data de liquidacao; impedir pagamento se meio inativo.
| 9 | Caixa Diario nao recebe entradas de vendas | Ausencia de chamadas a `CaixaDiarioRepository.CreateOperacao/UpdateTotais` em add pagamento/fechamento | Vendas em dinheiro nao afetam `TotalEntradas`, quebrando conferencias | No novo UC de finalizacao, para meios `DINHEIRO` (ou configurados como caixa), criar `OperacaoCaixa` com origem COMANDA e atualizar totais.
| 8 | Agendamentos podem virar `DONE` sem comanda/pagamento | `internal/infra/http/handler/appointment_handler.go` `CompleteAppointment` e `UpdateAppointmentStatusUseCase` | Bypass do fluxo financeiro gera servicos sem receita registrada | Restringir transicoes finais a UC de finalizacao da comanda; remover/limitar rota `/complete` ou validar existencia de comanda fechada.
| 8 | Comissao inexistente ou incorreta | Nenhum gatilho para `CreateCommissionItemUseCase`; `command_items` nao guardam profissional | Profissionais nao recebem comissao, impossibilitando conferencia | Ao fechar comanda, gerar `commission_items` por item de servico com base no profissional do agendamento e regra aplicavel; recalcular quando descontos mudarem.
| 7 | Multi-unidade inconsistente | Tabela `commands` possui `unit_id` (migrations 036/039) mas entity/repositorio nao populam | Pode falhar RLS `check_tenant_unit_access` e impede segmentacao por unidade | Adicionar `UnitID` em entity/DTO/mapper/repositorio e preencher a partir do agendamento/unidade corrente.
| 6 | Tipos monetarios em `float64` e erro de formatacao | `entity.Command` campos monetarios e helper `formatMoney` retornando rune | Arredondamento impreciso e mensagens de erro incorretas | Migrar valores para `decimal.Decimal` (como `Money`), ajustar helpers e armazenamento; cobrir com testes.
| 5 | Validacao fraca de itens da comanda | `AddCommandItemUseCase` aceita qualquer `item_id` e quantidades duplicadas | Itens inconsistentes/duplicados afetam total e comissao | Ler catalogo de servicos/produtos via porta dedicada, consolidar itens iguais e impedir duplicidade; validar status ativo.

## Correcoes recomendadas (visao operacional)
- Implementar **use case de finalizacao completa**: recebe `command_id`, carrega meios de pagamento, cria `command_payments`, gera `conta_receber` ou baixa imediata, registra `operacao_caixa` conforme meio (dinheiro = entrada, outros = so financeiro), processa `commission_items`, fecha comanda e atualiza agendamento para `DONE`. Tudo em transacao unica.
- Alinhar **multi-unidade**: propagar `unit_id` do agendamento ou contexto para comanda, caixa e financeiro; ajustar DTOs e repositorios.
- Substituir `float64` por `decimal` nas entidades de comanda/pagamento e adicionar testes de arredondamento.
- Fortalecer **validacoes**: verificar existencia/ativo do servico/produto e profissional antes de inserir item; impedir fechamento se `command_payments` nao cobrirem total quando `DeixarSaldoDivida=false`.
- Criar **servico de conciliacao**: usar `MeioPagamento.DMais` para prever liquidacao e alimentar Fluxo de Caixa Diario/DRE.
- Adicionar **eventos de dominio/audit log** (outbox) para comanda criada, item adicionado/removido, pagamento registrado, comanda fechada, operacao de caixa, reprocessamento de comissao.

## Eventos de auditoria sugeridos
- `appointment.finished` com `appointment_id`, `professional_id`, `command_id`, usuario.
- `command.payment_added` com `command_id`, `payment_id`, `meio_pagamento_id`, bruto, liquido, taxas, usuario.
- `command.closed` com totais, troco/saldo_devedor, `fechado_por`, `caixa_id` (quando existir).
- `financial.movement_created` com referencia da comanda/cliente/meio de pagamento.
- `cash.operation_created` para VENDA/SANGRIA/REFORCO com `caixa_id`, `command_id` opcional, usuario.
- `commission.item_created` com `professional_id`, `service_id`, `gross`, `commission_value`, referencia `command_item_id`.

## Pronto para auditoria financeira real?
**Nao.** Falta integracao entre comanda, financeiro e caixa; ausencia de logs auditaveis e uso de `float` inviabilizam conciliacao oficial ou fiscal.

## Comparacao com AppBarber / Trinks / OneBeleza
| Ponto de mercado | Expectativa mercado | Estado NEXO | Gap |
| --- | --- | --- | --- |
| Fechamento obriga movimentacao financeira + registro em caixa | Obrigatorio | Nao implementado | Precisa UC orquestrador e operacao de caixa automatica.
| Comissao automatica por item com descontos | Obrigatorio | Nao existe | Gerar `commission_items` no fechamento com base em regras.
| Rastreabilidade de operador (quem fechou/recebeu) | Obrigatorio com logs | Parcial (campos sem log/auditoria) | Gravar audit log e vincular a caixa/financeiro.
| Controle multi-unidade nas vendas | Suportado | Campo `unit_id` ignorado | Propagar unidade nas entidades.
| Consistencia de valores (sem float) | Usam decimal | Usa float | Migrar para decimal.

## Conclusao
Nivel atual: **alto risco operacional/financeiro**; receita de servicos nao alimenta financeiro nem caixa, e comissao nao e calculada.
Para atingir nivel premium:
- Entregar o UC de finalizacao integrada (pagamento → financeiro → caixa → comissao).
- Corrigir tipos monetarios e multi-unidade.
- Implantar validacoes de itens/meios de pagamento e auditar eventos.
- Somente apos essas entregas o fluxo estara alinhado a AppBarber/Trinks/OneBeleza e apto a auditoria.
