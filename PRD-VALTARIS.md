# PRD — NEXO

ERP + CRM para Barbearias Premium

**Versão do PRD:** 1.0
**Status:** Em Desenvolvimento
**Responsável:** Andrey Viana (Product Owner)
**Data:** 21/11/2025

---

## 1. Visão Geral do Produto

### 1.1 Missão do NEXO

Ajudar donos, gerentes e gestores de barbearia a administrar o negócio de forma **profissional e eficiente**, enquanto os barbeiros têm acesso em tempo real ao seu desempenho, comissão e metas, aumentando a consciência financeira e o resultado da operação.

### 1.2 Diferencial Competitivo

NEXO não é apenas um sistema de agendamento.
É um **ERP + CRM exclusivo para barbearias**, desenhado para resolver os principais problemas reais do dono de barbearia:

- Gestão financeira completa (caixa, DRE, fluxo, metas, análise por período).
- Controle real de comissões e desempenho do barbeiro.
- Clube/assinaturas, recorrência e fidelização.
- Gamificação + plano de carreira para equipe.
- KPIs avançados (MRR, ARR, churn, LTV, CAC, ocupação, ticket, retorno).
- Precificação inteligente para produtos e serviços.

### 1.3 Visão de Futuro (5 anos)

Estar entre os **5 ERPs mais usados para barbearias no Brasil**, consolidado em barbearias premium e, após 3–4 anos, expandir para **salões premium** mantendo a mesma tese de gestão orientada a dados.

### 1.4 Público-Alvo

- Barbearias premium com foco em experiência e alta lucratividade.
- Redes de barbearias e operações multi-unidade.
- Donos que querem gestão de verdade, não só agenda.

---

## 2. Personas e Acessos

### 2.1 Personas

- **Dono da Barbearia**
  - Visão estratégica, financeira e de performance.
- **Gerente**
  - Responsável pela operação diária, equipe, metas e resultados.
- **Recepcionista**
  - Focada em agenda, fluxo de clientes e atendimento.
- **Barbeiro**
  - Focado no próprio desempenho, agenda, comissões e evolução.
- **Contador**
  - Olha o financeiro para fins fiscais/contábeis.
- **Cliente Final**
  - Agendamento, histórico, fidelidade, experiência.

### 2.2 Nível de Acesso

- **Dono**

  - Acesso total a todas as unidades.
  - Ver/editar tudo, inclusive financeiro, relatórios e exportações.

- **Gerente**

  - Acesso total às unidades atribuídas.
  - Ver/editar agenda, equipe, comissões, estoque, financeiro e relatórios.
  - Não exporta dados sensíveis (pode ser configurável depois).

- **Recepcionista**

  - Pode:
    - Ver agenda geral da unidade.
    - Criar/editar/mover/cancelar agendamentos.
    - Operar lista da vez.
    - Cadastrar clientes.
    - Ver dados básicos de clientes.
    - Interagir com lista de espera.
    - Ver estoque em nível básico (itens, saldo, alertas).

- **Barbeiro**

  - Pode:
    - Ver **apenas sua própria agenda**.
    - Ver comissões próprias.
    - Ver metas próprias.
    - Ver ranking/gamificação.
    - Ver histórico de seus atendimentos.
  - Não vê dados sensíveis dos clientes (somente nome e serviço).

- **Contador**

  - Acesso **read-only** ao módulo financeiro e relatórios financeiros.

- **Cliente (App)**
  - Pode:
    - Agendar serviços.
    - Ver histórico de agendamentos.
    - Ver histórico de compras.
    - Avaliar atendimentos.
    - Ver saldo de cashback / fidelidade.
    - Receber notificações.
    - Integrar agendamentos com Google Agenda.

### 2.3 Multi-unidade / Rede / Franquia

- **Sim**.
  O NEXO deve suportar multi-unidade:
- Cada unidade com operação própria.
- Consolidação de dados em nível de rede.

---

## 3. Escopo Global — Módulos do Sistema

1. **Agendamento / Agenda**
2. **Lista da Vez**
3. **Sistema de Assinaturas** (via Asaas)
4. **Gestão Financeira**
5. **Comissões de Barbeiros**
6. **Gestão de Estoque**
7. **CRM de Clientes**
8. **Programa de Fidelidade (Cashback)**
9. **Gamificação + Plano de Carreira**
10. **Metas & KPIs**
11. **Precificação Inteligente**
12. **Relatórios Gerenciais**
13. **App do Barbeiro**
14. **App do Cliente**
15. **Integrações Externas**
16. **Multi-unidade / Franquias**
17. **Notas Fiscais (futuro)**

---

## 4. Funcionalidades por Módulo

### 4.1 Módulo de Agendamento

**Objetivo:**
Permitir que recepção/gerência agende serviços de forma visual (calendário estilo AppBarber/Trinks), com controle total de horários, barbeiros, unidades e integração com Google Agenda.

**Regras de Negócio:**

- Não pode agendar com barbeiro inativo.
- Não pode agendar sem cadastro de cliente.
- Intervalo padrão entre horários: **10 minutos**.
- Um agendamento sempre pertence a:
  - 1 unidade
  - 1 barbeiro
  - 1 cliente
  - 1 ou mais serviços
- Recepção pode criar, mover e cancelar agendamentos.
- Status sugeridos de agendamento:
  - `CREATED`
  - `CONFIRMED`
  - `IN_SERVICE`
  - `DONE`
  - `NO_SHOW`
  - `CANCELED`
- **Google Agenda**:
  - Sincronizar:
    - Agendamentos confirmados
    - Cancelamentos
    - Alterações de horário

**Funcionalidades:**

- Visualização em calendário (por barbeiro, por unidade).
- Criação de agendamento com busca de cliente.
- Bloqueio de horário já ocupado.
- Marcação de no-show.
- Reagendamento simples (drag & drop ou formulário).
- Filtros:
  - por barbeiro
  - por serviço
  - por dia/semana.

---

### 4.2 Módulo Lista da Vez

**Objetivo:**
Organizar a ordem de atendimento dos barbeiros de forma **justa e automática**, com histórico e estatísticas.

**Regras de Negócio:**

- Apenas barbeiros **ativos** participam.
- Ordenação base (já existente):
  - `current_points ASC, last_turn_at ASC NULLS FIRST, name ASC`.
- Reset automático mensal (dia 1).
- Histórico preservado para análise.
- A cada atendimento finalizado:
  - o barbeiro vai para o fim da fila.
  - pontos ajustados segundo regras (configuráveis).

**Funcionalidades:**

- Mostrar lista atual da vez por unidade.
- Pausar barbeiro (intervalo).
- Retomar barbeiro pausado.
- Exibir estatísticas:
  - total de atendimentos no período
  - média de atendimentos/dia
  - tempo médio de atendimento.

---

### 4.3 Módulo de Assinaturas (Asaas)

**Objetivo:**
Gerenciar planos recorrentes de clientes (clubes/assinaturas), com cobrança via Asaas (PIX/cartão), controle de limite de uso e suspensão automática por inadimplência.

**Regras de Negócio:**

- Dono pode criar planos personalizados.
- Planos podem ter:
  - preço
  - benefícios
  - limite de uso (ex.: X cortes por mês).
- Não há carência obrigatória.
- Sistema suspende automaticamente benefícios se assinatura não estiver ativa.
- Status da assinatura:
  - `ACTIVE`
  - `PENDING_PAYMENT`
  - `OVERDUE`
  - `CANCELED`
  - `EXPIRED`
- Retry automático de cobrança via Asaas.
- Notificação para gerente/dono/barbeiro quando cliente assinante está inadimplente.

**Funcionalidades:**

- Cadastro de planos.
- Vincular cliente a plano.
- Criar/atualizar cliente no Asaas.
- Gerar link de pagamento (PIX/cartão).
- Receber webhooks do Asaas e atualizar status.
- Demonstrar receita recorrente (MRR/ARR).
- Controle de uso dos benefícios (ex.: quantos serviços já usados no ciclo).

---

### 4.4 Módulo Financeiro

**Objetivo:**
Ser o **cérebro financeiro** da barbearia: registrar receitas, despesas, DRE, projeções e dar visão real da lucratividade.

**Escopo:**

- Receitas:
  - serviços
  - produtos
  - assinaturas
- Despesas:
  - fixas
  - variáveis
  - recorrentes
- Caixa diário.
- Contas a pagar / receber.
- DRE:
  - mensal
  - comparativo de meses
- Análises:
  - mensal
  - trimestral
  - semestral
  - anual
- Registro de comissões como despesa.
- Integração com MRR/ARR e assinaturas.

**Regras de Negócio (principais):**

- Todo lançamento pertence a:
  - unidade
  - categoria
  - data de competência.
- Comissões são registradas como despesa operacional.
- DRE monta automaticamente com base nos lançamentos.
- Lançamentos podem ser:
  - consolidados
  - conciliados (no futuro com extrato bancário).

**Funcionalidades:**

- Cadastro de categorias de receita/despesa.
- Registro manual de receitas e despesas.
- Marcação de recorrência.
- Visualização de fluxo de caixa.
- Geração de DRE mensal.
- Comparativo DRE mês a mês.
- Exportação CSV/Excel.

---

### 4.5 Módulo de Comissões

**Objetivo:**
Controlar comissões de forma justa, transparente e automática.

**Regras de Negócio:**

- Comissão sempre percentual, configurada por barbeiro.
- Não existe comissão fixa.
- Comissão calculada **por serviço** (unitário).
- Somente serviços pagos geram comissão.
- Comissão nunca pode ultrapassar valor do serviço.
- Bônus podem ser aplicados ao bater metas (gamificação / metas).

**Funcionalidades:**

- Configuração de percentual por barbeiro.
- Configuração de regras de bônus.
- Relatório de comissões:
  - por barbeiro
  - por período
  - por unidade.
- Status da comissão:
  - `PENDING`
  - `PAID`
  - `CANCELED`.

---

### 4.6 Módulo de Estoque

**Objetivo:**
Controlar completamente produtos e insumos, inclusive consumo interno e custo por serviço.

**Regras de Negócio:**

- Não permitir estoque negativo.
- Categorias de estoque personalizadas.
- Produto com:
  - nome
  - categoria
  - unidade
  - custo
  - preço de venda.
- Consumo interno reduz estoque.
- Serviços podem ter “ficha técnica” de insumos para abatimento automático.
- Alerta de baixo estoque com base em estoque mínimo.

**Funcionalidades:**

- Cadastro de produtos/insumos.
- Registro de entradas e saídas.
- Registro de consumo interno (ex.: café, toalha, insumo consumo próprio).
- Relatório de saldo atual.
- Relatório de custo de insumo por serviço.
- Alertas de baixo estoque.

---

### 4.7 Módulo de CRM

**Objetivo:**
Centralizar tudo que é relacionado ao cliente: histórico, comportamento, tags e engajamento.

**Escopo:**

- Cadastro completo de clientes.
- Histórico de agendamentos.
- Histórico de compras (produtos/serviços).
- Origem do cliente (indicação, Instagram, Google, etc.).
- Tags (VIP, Risco, Retenção, Novo, etc.).
- Pontuação de engajamento.
- Preferência de barbeiro.

**Regras de Negócio:**

- Barbeiro vê apenas:
  - nome do cliente
  - serviços realizados
  - sem dados sensíveis (telefone, e-mail, etc.).
- Cliente pode avaliar atendimentos.
- Engajamento pode considerar:
  - frequência de visita
  - ticket médio
  - tempo desde a última visita.

**Funcionalidades:**

- Tela de cliente com linha do tempo.
- Filtros por tags, origem, frequência.
- Relatório de clientes ativos/inativos.
- Base para campanhas de marketing (futuro).

---

### 4.8 Módulo de Fidelidade (Cashback)

**Objetivo:**
Aumentar retenção, fazendo o cliente se sentir recompensado ao consumir serviços/produtos.

**Regras de Negócio:**

- Cashback configurável por unidade.
- Expiração configurável (parametrizado em tela de parâmetros).
- Cashback pode ser usado em:
  - serviços
  - produtos
  - desconto parcial.
- Não pode gerar saldo negativo.

**Funcionalidades:**

- Configuração de regras de cashback.
- Acúmulo automático a cada venda elegível.
- Uso parcial/total em novas compras.
- Exibição do saldo no aplicativo do cliente.
- Histórico de movimentação de cashback.

---

### 4.9 Módulo de Gamificação + Plano de Carreira

**Objetivo:**
Engajar barbeiros, aumentar performance e criar sensação de evolução e reconhecimento.

**Regras de Negócio:**

- Sistema de XP (pontos de experiência).
- Níveis:
  - Bronze
  - Prata
  - Ouro
  - Diamante
- Ganho de XP baseado em:
  - atendimentos
  - ticket médio
  - retenção de clientes
  - pontualidade (no futuro, se tiver).
- Ao atingir certos níveis:
  - pode haver aumento automático de comissão (plano de carreira).
  - bônus podem ser concedidos.

**Funcionalidades:**

- Painel do barbeiro com:
  - nível atual
  - XP atual
  - XP para próximo nível.
- Painel do gestor com visão geral da equipe.
- Configuração de regras de XP.

---

### 4.10 Módulo de Metas & KPIs

**Objetivo:**
Permitir que donos/gerentes definam metas claras e acompanhem KPIs chave do negócio.

**Metas suportadas:**

- Meta de faturamento.
- Meta de venda de produtos.
- Meta de venda de serviços.
- Meta de despesas (teto).
- Meta de lucro.
- Meta de novos clientes.
- Meta por barbeiro (opcional / futuro).

**KPIs (calculados pelo sistema):**

- MRR & ARR (recorrência via assinaturas).
- Churn.
- LTV.
- CAC (manual ou via integração futura).
- Taxa de ativação (clientes que realmente usam após cadastro).
- % de receita gerada via assinaturas.
- Capacidade operacional (ocupação).
- Tempo médio de atendimento.
- Ticket médio por barbeiro e por unidade.
- Taxa de retorno (cliente que volta em menos de 30 dias).
- Taxa de no-show.

_(Fórmulas detalhadas podem entrar em um anexo depois.)_

---

### 4.11 Módulo de Precificação Inteligente

**Objetivo:**
Ajudar o gestor a definir preços corretos para produtos e serviços, considerando todos os custos e margem desejada.

**Regras de Negócio:**

- Não altera o preço automaticamente (primeira fase).
- Apenas **sugere** preço com base em:
  - custo de compra
  - insumos do serviço
  - comissões
  - impostos (quando parametrizados)
  - taxas de cartão / adquirência
  - margem desejada.

**Funcionalidades:**

- Simulador de preço de produto.
- Simulador de preço de serviço.
- Comparação entre:
  - preço atual
  - preço sugerido
  - margem atual x margem alvo.

---

### 4.12 Módulo de Relatórios

**Objetivo:**
Dar visão gerencial clara em diferentes períodos.

**Períodos suportados:**

- Diário
- Semanal
- Mensal
- Trimestral
- Semestral
- Anual

**Filtros principais:**

- Barbeiro
- Unidade
- Serviço
- Produto
- Ticket médio
- Tipo de cliente (novo/recorrente)

**Relatórios chave:**

- Taxa de ocupação da loja e por barbeiro.
- Ticket médio por barbeiro/unidade.
- Faturamento por período.
- Receita recorrente (assinaturas).
- DRE mensal e comparativo de meses.
- Taxa de fidelização (cliente que volta com o mesmo barbeiro).
- Taxa de retorno em 30 dias.
- Ranking de barbeiros.

**Exportação:**

- CSV
- Excel

---

### 4.13 App do Barbeiro

**Objetivo:**
Ser o painel pessoal do profissional.

**Funcionalidades:**

- Ver agenda própria.
- Ver metas e evolução.
- Ver comissões.
- Ver ranking/nível (gamificação).
- Ver histórico de atendimentos.
- Ver impactos de suas ações (ex.: ticket, retorno).

---

### 4.14 App do Cliente

**Objetivo:**
Facilitar agendamento e reforçar relacionamento com a barbearia.

**Funcionalidades:**

- Agendar serviços.
- Ver histórico de agendamentos.
- Ver histórico de compras.
- Avaliar atendimentos.
- Ver saldo de cashback.
- Receber lembretes.
- Sincronizar agendamento com Google Agenda.

---

### 4.15 Integrações

**Asaas:**

- Criação de cliente.
- Criação de assinatura.
- Webhooks:
  - pagamento aprovado
  - rejeitado
  - cancelado
- Retry automático.
- Suspensão de benefícios por inadimplência.

**Google Agenda:**

- Enviar agendamentos confirmados.
- Atualizar cancelamentos.
- Atualizar mudanças de horário.

**WhatsApp (futuro):**

- Confirmação.
- Lembretes.
- Lista de espera.

---

### 4.16 Multi-unidade / Franquias

**Objetivo:**
Permitir que donos com várias unidades vejam tudo consolidado e também por unidade.

**Funcionalidades:**

- Painel por unidade.
- Painel consolidado de rede.
- Relatórios por unidade x rede.

---

### 4.17 Notas Fiscais (Futuro)

**Escopo futuro:**

- Emissão de NF integradas.
- Integração com contabilidade (quando for prioridade).

---

## 5. Regras de Negócio Críticas (Resumo)

- **RN-001 — Multi-tenant obrigatório:**
  - Toda operação respeita `tenant_id` (barbearia).
- **RN-002 — Assinatura ativa:**
  - Clientes com assinatura inadimplente têm benefícios bloqueados.
- **RN-003 — Comissões:**
  - Apenas sobre serviços pagos.
- **RN-004 — Lista da vez:**
  - Apenas barbeiros ativos, reset mensal, histórico mantido.
- **RN-005 — Estoque:**
  - Nunca permitir estoque negativo.
- **RN-006 — Privacidade:**
  - Barbeiro não vê dados sensíveis de cliente.

---

## 6. Requisitos Técnicos (Macro)

- Sistema **online** (sem necessidade de modo offline).
- Latência alvo:
  - API < **150ms**
  - Interface < **300ms** em operações principais.
- Logs:
  - Sempre salvar **antes** da alteração (valor antigo).
- Auditoria:
  - Inicialmente sem auditoria completa de banco, mas com logs de ações críticas.
- Backup:
  - Diário.
- Segurança:
  - JWT
  - Multi-tenant forte
  - Separação rígida por barbearia/unidade.

---

## 7. Critérios Gerais de Aceite

Para considerar um módulo “PRONTO” na versão 1.0:

- Funcionalidades principais implementadas conforme regras de negócio.
- Telas responsivas (desktop e mobile).
- Logs de operações críticas funcionando.
- Exportação CSV/Excel nos relatórios principais.
- Perfis de acesso respeitados.
- Erros tratados de forma amigável (sem stacktrace pro usuário).
- KPIs básicos calculando corretamente.

---

## 8. Roadmap (Visão Rápida)

> **Obs.:** detalhamento fino do roadmap pode ficar em outro documento, mas o PRD já aponta direção.

### MVP (1.0 — Core Operacional)

- Agendamento
- Lista da vez
- Financeiro básico (receitas/despesas, caixa, DRE simples)
- Comissões
- Estoque essencial
- Assinaturas (básico Asaas)
- CRM básico
- Relatórios mensais simples
- Permissões

### Pós-MVP (1.1 / 1.2)

- Fidelidade (cashback)
- Gamificação
- Metas & KPIs avançados
- Precificação inteligente
- Relatórios avançados (ocupação, retorno, comparativos trimestrais/anual)
- Apps (barbeiro/cliente)

### Futuro (2.0+)

- Notas fiscais integradas
- Integrações avançadas (bancos, maquininha, etc.)
- Recursos avançados de rede/franquia
- IA para previsão de retorno, ocupação e preços

---
