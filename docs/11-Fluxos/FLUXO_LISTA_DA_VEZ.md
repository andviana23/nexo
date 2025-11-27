# Fluxo da Lista da Vez — NEXO

> **Versão:** 2.0  
> **Última atualização:** 26/11/2025  
> **Módulo:** Lista da Vez  
> **Tipo:** Fila giratória manual com reset mensal

---

## 📋 Sumário

1. [Visão Geral](#visão-geral)
2. [Principais Componentes](#principais-componentes)
3. [Lógica Oficial do Módulo](#lógica-oficial-do-módulo)
4. [Fluxo UX Expandido](#fluxo-ux-expandido)
5. [Diagrama Mermaid](#diagrama-mermaid)
6. [Validação do Padrão UX](#validação-do-padrão-ux)

---

## Visão Geral

A **Lista da Vez** é um módulo **totalmente manual**, independente de atendimentos reais, com lógica de **fila giratória** e **reset mensal automático**.

### Características Principais

- ✅ **Manual**: Não depende de agendamentos ou atendimentos reais
- ✅ **Fila Giratória**: Barbeiro atendido vai para o final
- ✅ **Reset Mensal**: Contadores zerados no último dia do mês
- ✅ **Ordem Base**: Sempre respeita ordem de cadastro em empates
- ✅ **Relatórios**: Diários e mensais para análise

---

## Principais Componentes

### 🎭 Atores

| Ator | Descrição |
|------|-----------|
| **Recepção** | Operadora da lista, registra atendimentos manuais |
| **Sistema** | Executa regras internas automaticamente |

### 📦 Entidades

| Entidade | Descrição |
|----------|-----------|
| **Barbeiro** | Profissional com função "Barbeiro" |
| **Lista da Vez** | Fila dinâmica ordenada por contadores |
| **Contador de Atendimentos** | Registro manual por barbeiro |
| **Relatório Diário** | Consolidação de atendimentos do dia |
| **Reset Mensal** | Zeragem automática de contadores |

### ⚡ Eventos Críticos

| Evento | Trigger | Resultado |
|--------|---------|-----------|
| Clique no "+" | Ação manual | Incrementa contador + reordena fila |
| Reordenar lista | Automático após clique | Menor contador → topo |
| Reset mensal | Último dia às 23:59 | Zera contadores, mantém ordem base |
| Visualizar relatório | Ação manual | Exibe dados consolidados |

---

## Lógica Oficial do Módulo

### 🟦 Ações da Recepção (Manuais)

1. Abrir a tela da Lista da Vez
2. Clicar no botão "+" para registrar que o barbeiro atendeu alguém
3. Consultar relatórios diários
4. (Opcional) Filtrar dias/mês/unidade

### 🟩 Ações Automáticas do Sistema

1. Carregar barbeiros e ordená-los pela ordem de cadastro
2. Somar +1 para o barbeiro selecionado
3. Mover barbeiro para o final da fila
4. Reordenar lista automaticamente (menor quantidade → topo)
5. Gerar relatório diário ao final do dia
6. Resetar contadores no último dia do mês às 23:59
7. Registrar histórico no banco de dados

---

## Fluxo UX Expandido

### [1] Início

- Sistema abre o módulo Lista da Vez
- Carrega todos os barbeiros cadastrados com papel: \`Barbeiro\`
- Ordena todos pela **ordem original de cadastro** (nunca muda)

### [2] Exibição da Fila

O sistema exibe para cada barbeiro:

| Campo | Descrição |
|-------|-----------|
| Nome | Nome do barbeiro |
| Contagem | Atendimentos manuais registrados |
| Posição | Posição atual na fila dinâmica |
| Ação | Botão ➕ ao lado |

> ⚠️ **Importante**: Nenhuma relação com agenda ou atendimentos reais.

### [3] Ação: Recepcionista clica em ➕

Quando o botão é clicado, o sistema executa:

1. **Incrementa** o contador daquele barbeiro (+1)
2. **Move** imediatamente o barbeiro para o final da fila
3. **Reordena** a fila:
   - Quem tem menos registros fica no **topo**
   - Em caso de **empate** → volta a ordem original do cadastro

### [4] Relatório Diário

A recepção pode:

- Abrir o relatório daquele dia
- Analisar cores e quantidades
- Ver o total do dia
- Ver totais acumulados no mês

### [5] Reset Mensal Automático

No **último dia do mês**, às **23:59**:

1. Gera relatório geral do mês
2. Zera todos os contadores
3. Mantém a ordem original de cadastro
4. Começa novo ciclo no dia 1

> ✅ **Sem intervenção humana. Sem chance de erro.**

---

## Diagrama Mermaid

\`\`\`mermaid
flowchart TB

%% ================================
%%          INÍCIO DO FLUXO
%% ================================
A([Início do Módulo]) --> B

%% ================================
%%       CARREGAMENTO INICIAL
%% ================================
B[🟩 Sistema carrega barbeiros com função 'Barbeiro'] --> C
C[🟩 Ordenar barbeiros pela ordem de cadastro - ordem base] --> D
D[🟨 Exibir lista dinâmica<br/>• Nome<br/>• Contador manual<br/>• Posição atual<br/>• Botão +] --> E

%% ================================
%%     INTERAÇÃO PRINCIPAL - +
%% ================================
E --> F{🟦 Recepção clicou no botão +?}

F -->|Não| E

F -->|Sim| G[🟩 Incrementar contador do barbeiro selecionado +1]
G --> H[🟩 Mover barbeiro selecionado para o final da lista]
H --> I[�� Reordenar lista:<br/>1. Menor contador → topo<br/>2. Empate → manter ordem de cadastro] 
I --> E

%% ================================
%%   RELATÓRIO DIÁRIO DE USO
%% ================================
E --> J{🟦 Recepção abriu Relatório Diário?}

J -->|Não| E

J -->|Sim| K[🟩 Sistema exibe relatório:<br/>• Atendimentos manuais por barbeiro<br/>• Total do dia<br/>• Histórico colorido<br/>• Totais acumulados] 
K --> E

%% ================================
%%      RESET AUTOMÁTICO MENSAL
%% ================================
E --> L{🟩 Data = último dia do mês às 23:59?}

L -->|Não| E

L -->|Sim| M[🟩 Gerar relatório mensal final]
M --> N[🟩 Zerar todos os contadores da fila]
N --> O[🟩 Restaurar ordem base - ordem de cadastro]
O --> E
\`\`\`

---

## Validação do Padrão UX

### ✅ Por que esse fluxo está no padrão UX correto?

| Critério | Status | Descrição |
|----------|--------|-----------|
| **Estados claros** | ✅ | Cada etapa do sistema é um estado, não um passo solto |
| **Atores definidos** | ✅ | Recepção → ações manuais / Sistema → lógica automatizada |
| **Lógica primária + secundária** | ✅ | Principal: clicar no "+" / Secundário: relatórios / Automático: reset |
| **Decisões UX (diamantes)** | ✅ | Sempre que há escolha humana ou automática, existe um "branch" |
| **Loop contínuo** | ✅ | O sistema sempre retorna para o estado base → Exibir Lista |
| **Single Source of Truth** | ✅ | Fila = contadores manuais / Atendimento real = módulo separado |

### 🎯 Padrões Respeitados

- ✅ Segue o padrão de fluxos de telas **SaaS**
- ✅ **Fila** = contadores manuais (independente)
- ✅ **Atendimento real** = módulo separado (Agendamentos/Financeiro)
- ✅ **Separação de responsabilidades** clara

---

## Regras de Negócio

### RN-LDV-001: Ordem Base
> A ordem base dos barbeiros é definida pela **ordem de cadastro** no sistema. Esta ordem nunca muda e é usada como critério de desempate.

### RN-LDV-002: Incremento de Contador
> Ao clicar no botão "+", o contador do barbeiro é incrementado em 1 e ele é movido para o final da fila.

### RN-LDV-003: Ordenação da Fila
> A fila é ordenada pelo **menor contador** no topo. Em caso de empate, usa-se a **ordem de cadastro**.

### RN-LDV-004: Reset Mensal
> No último dia de cada mês, às 23:59, todos os contadores são zerados automaticamente e um relatório mensal é gerado.

### RN-LDV-005: Independência de Módulos
> A Lista da Vez é **totalmente independente** do módulo de Agendamentos. Os contadores são manuais e não refletem atendimentos reais do sistema.

---

## Histórico de Alterações

| Versão | Data | Autor | Descrição |
|--------|------|-------|-----------|
| 1.0 | - | - | Versão inicial simplificada |
| 2.0 | 26/11/2025 | Sistema | Documentação completa com fluxo UX profissional |
