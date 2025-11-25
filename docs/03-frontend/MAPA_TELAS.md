# NEXO — Mapa de Telas

> Estrutura de navegação e principais telas do sistema.

---

## 1. Autenticação `(auth)`

Layout público, centralizado, sem sidebar.

- **/login:** Formulário de login (Email/Senha). Link para recuperação.
- **/register:** Cadastro de novos tenants (Wizard de Onboarding).
- **/forgot-password:** Solicitação de reset de senha.
- **/reset-password:** Definição de nova senha.

---

## 2. Dashboard `(dashboard)`

Layout protegido, com Sidebar lateral e Header.

### 2.1 Visão Geral

- **/dashboard:** KPIs principais (Faturamento, Agendamentos, Novos Clientes). Gráficos de desempenho.

### 2.2 Agenda

- **/agenda:** Calendário interativo (Day/Week/Month view).
  - Modal de Novo Agendamento.
  - Detalhes do Agendamento (Status, Pagamento).

### 2.3 Clientes

- **/clientes:** Listagem de clientes (DataTable).
- **/clientes/novo:** Formulário de cadastro.
- **/clientes/[id]:** Perfil do cliente (Histórico, Dados, Métricas).

### 2.4 Financeiro

- **/financeiro/transacoes:** Lista de receitas e despesas.
- **/financeiro/caixa:** Fluxo de caixa diário/mensal.
- **/financeiro/comissoes:** Relatório e pagamento de comissões.

### 2.5 Serviços & Produtos

- **/servicos:** Catálogo de serviços (Preço, Duração, Profissionais).
- **/produtos:** Controle de estoque e venda de produtos.

### 2.6 Configurações

- **/configuracoes/perfil:** Dados do usuário logado.
- **/configuracoes/empresa:** Dados da barbearia (Logo, Endereço).
- **/configuracoes/equipe:** Gestão de profissionais e permissões.
- **/configuracoes/horarios:** Horário de funcionamento.

---

## 3. Componentes Chave por Tela

| Tela           | Componentes Principais                                   | Dados (Hooks)         |
| -------------- | -------------------------------------------------------- | --------------------- |
| **Login**      | `Card`, `Form`, `Input`, `Button`                        | `useAuth`             |
| **Dashboard**  | `Card` (Metrics), `Chart` (Recharts), `DateRangePicker`  | `useDashboardMetrics` |
| **Agenda**     | `Calendar` (FullCalendar/DayPilot), `Dialog` (New Event) | `useAppointments`     |
| **Clientes**   | `DataTable` (Pagination, Filter), `Sheet` (Quick View)   | `useClients`          |
| **Financeiro** | `DataTable`, `Select` (Period), `Badge` (Status)         | `useTransactions`     |

---

## 4. Fluxos Críticos

1.  **Agendamento:**
    - Cliente/Admin seleciona Serviço -> Profissional -> Horário -> Confirma.
2.  **Checkout:**
    - Finalizar agendamento -> Adicionar produtos (opcional) -> Selecionar forma de pagamento -> Baixa no estoque/Financeiro.
3.  **Onboarding:**
    - Cadastro Empresa -> Configuração Horários -> Cadastro Serviços -> Convite Equipe.
