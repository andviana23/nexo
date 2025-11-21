# 1.6 Comissões Automáticas

- **Categoria:** FINANCEIRO
- **Objetivo:** automatizar o cálculo e lançamento de comissões por serviço/barbeiro/venda com múltiplos modelos (fixo, percentual, degraus progressivos), permitindo relatórios por período e exportação em PDF.

## Plano de Execução (quarta prioridade)
- **Banco de Dados:** `comissoes_config` (modelos fixo/percentual/degrau por serviço/barbeiro/tipo) e `comissoes_lancadas` (valores calculados, período, origem); índices por tenant/período/modelo.
- **Backend:** engine de cálculo acionada por eventos/cron após recebimentos confirmados; integração com ProcessarRepasse; exportação PDF; auditoria de configurações.
- **Frontend:** UI de configuração com preview/simulador, relatório por período/barbeiro/modelo, exportação PDF.
- **Cálculos aplicados:** custo variável para DRE e precificação de serviço. Se comissão integra preço do serviço, usar “Cálculo do Preço de Serviço” (`docs/10-calculos/preco-servico.md`) e validar regra de comissão % do PV (ponto marcado como **INFORMAÇÃO AUSENTE** na doc).

## Regras de Negócio

- Modelos suportados:
  - Valor fixo por serviço.
  - Percentual por serviço/barbeiro/tipo de venda (serviço/produto).
  - Degrau progressivo (ex.: 40% até 12k, 45% até 18k, 50% acima disso) calculado mensalmente.
- Comissões só são consideradas quando o recebimento correspondente está `RECEBIDO` (assinaturas) ou `CONFIRMADO` (serviços).
- Configuração feita por tenant e pode ser diferenciada por barbeiro/grupo de serviços.
- Relatórios devem respeitar RBAC (owner/manager) e permitir exportação PDF.

## Dependências Técnicas

- Tabelas `servicos`, `assinatura_invoices`, `receitas`, `comissoes_config`, `comissoes_lancadas`.
- Engine de cálculo acionada por cron ou eventos (processamento de repasse).
- Serviço de exportação PDF e UI de relatórios.

## Riscos

- Configurações complexas podem gerar erros de cálculo (necessário simulador/preview).
- Degraus mal configurados podem se sobrepor (validar ranges).
- Alto volume de lançamentos pode impactar performance (usar batch + indexes).

## Tarefas

1. Modelar `comissoes_config` com suporte a modelos fixo/percentual/degrau e escopos (serviço, barbeiro, tipo de venda).
2. Implementar engine de cálculo que recebe eventos de receita/assinatura e gera lançamentos em `comissoes_lancadas`, incluindo integrações com ProcessarRepasseUseCase.
3. Disponibilizar UI para configuração (forms + pré-visualização) e relatórios com filtros por período/barbeiro/modelo.
4. Criar endpoint de exportação (`POST /financial/comissoes/export`) gerando PDF e logs de auditoria.
5. Adicionar testes unitários (engine), integração (cron + repasse) e e2e (configuração → relatório).
6. Atualizar DRE e dashboard para consumir comissões automaticamente como custo variável.

## Critérios de Aceite

- Configurações aceitam todos os modelos descritos e validam ranges (sem sobreposição de degraus).
- Engine gera lançamentos automaticamente após recebimentos confirmados; dados refletem no fluxo de caixa.
- Relatórios exibem totais por barbeiro/período e permitem exportação PDF.
- Logs de auditoria registram alterações de configuração e geração de relatórios.
- Testes automatizados cobrem cálculos críticos e passam no pipeline CI.
