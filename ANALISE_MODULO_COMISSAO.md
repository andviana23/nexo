# Análise Completa do Módulo de Comissão

## 1. Mapeamento Geral do Módulo

O módulo de comissão está estruturado seguindo a Clean Architecture do projeto, com separação clara entre entidades, casos de uso e persistência.

### Banco de Dados
*   **`commission_rules`**: Regras de comissão (Global, Unidade, etc.).
*   **`commission_items`**: Tabela principal. Registra cada comissão gerada por item de comanda.
*   **`commission_periods`**: Fechamentos mensais de comissão por profissional.
*   **`comissoes_categoria_profissional`**: Tabela para regras por categoria (existente mas **não utilizada**).

### Backend (Go)
*   **Criação/Cálculo**: `backend/internal/application/usecase/command/finalizar_comanda_integrada.go` (Lógica principal).
*   **Entidades**: `backend/internal/domain/entity/commission_*.go`.
*   **Use Cases**: `backend/internal/application/usecase/commission/`.
*   **Repositórios**: `backend/internal/infra/repository/postgres/commission_*.go`.
*   **Queries**: `backend/internal/infra/db/queries/commission_*.sql`.

### Frontend (Next.js)
*   **Serviços**: `frontend/src/services/commission-service.ts`.
*   **Hooks**: `frontend/src/hooks/use-commissions.ts`.
*   **Tipos**: `frontend/src/types/commission.ts`.

---

## 2. Validação das Regras de Negócio

### Serviços e Produtos
*   **Serviços**: ✅ Implementado. Existe uma hierarquia de 4 níveis para definir a taxa:
    1.  Comissão específica do Serviço (prioridade máxima).
    2.  Comissão do Profissional.
    3.  Regra da Unidade.
    4.  Regra Global do Tenant.
*   **Produtos**: ❌ **NÃO IMPLEMENTADO**. O código em `FinalizarComandaIntegradaUseCase` processa o estoque de produtos, mas **ignora completamente** o cálculo de comissão para itens do tipo `PRODUTO`.
*   **Categorias**: ⚠️ **PARCIAL**. A tabela `comissoes_categoria_profissional` existe no banco, mas a lógica de cálculo (`buscarRegraComissaoHierarquica`) não consulta essa tabela. Regras por categoria de serviço são ignoradas.

### Descontos e Base de Cálculo
*   **Base de Cálculo**: ✅ Suporta `BRUTO` e `LIQUIDO`.
*   **Cálculo Líquido**: ✅ Implementado corretamente. O sistema calcula o valor proporcional do item em relação ao total líquido pago (considerando descontos na comanda).
    *   Fórmula: `(PreçoItem / TotalComanda) * TotalPagamentosLiquidos`.
*   **Cálculo Bruto**: ✅ Usa o `PrecoFinal` do item (que já inclui descontos aplicados diretamente no item, se houver).

### Taxas e Gorjetas
*   **Taxas de Cartão**: ✅ Se a base for `LIQUIDO`, as taxas são deduzidas (pois `TotalPagamentosLiquidos` já desconta a taxa do meio de pagamento). Se for `BRUTO`, ignora as taxas.
*   **Gorjetas**: ❓ O campo `DeixarTrocoGorjeta` existe na comanda, mas não fica claro se entra na base de comissão. Pela lógica atual, se entrar como pagamento, entra no rateio proporcional.

### Regras Especiais
*   **Profissional Diferente**: ⚠️ O sistema usa o profissional do **Agendamento** (`AppointmentID`) para todos os serviços da comanda. Se uma comanda tiver serviços feitos por profissionais diferentes (sem agendamento vinculado item a item), a comissão pode ir para a pessoa errada ou não ser gerada corretamente.
*   **Cancelamento**: ✅ Existe status `CANCELADO` e `ESTORNADO` em `commission_items`.
*   **Data de Competência**: ⚠️ A comissão é registrada com a data do **Fechamento da Comanda** (`time.Now()`), e não a data do agendamento. Isso pode afetar relatórios se a comanda for fechada dias depois.

---

## 3. Verificação de Bugs e Inconsistências

| Gravidade | Local | Problema | Sugestão de Correção |
| :--- | :--- | :--- | :--- |
| 🔴 **CRÍTICO** | `finalizar_comanda_integrada.go` | **Comissão de Produtos inexistente**. Vendas de produtos não geram registro na tabela `commission_items`. | Implementar lógica similar à de serviços para produtos no loop de itens da comanda. |
| 🟠 **ALTO** | `finalizar_comanda_integrada.go` | **Regra por Categoria ignorada**. A tabela `comissoes_categoria_profissional` não é lida. | Adicionar consulta a essa tabela na função `buscarRegraComissaoHierarquica` (entre nível Serviço e Profissional). |
| 🟡 **MÉDIO** | `finalizar_comanda_integrada.go` | **Profissional Único por Comanda**. Assume que todos os itens são do profissional do agendamento. | Permitir vincular profissional a cada item da comanda individualmente, ou validar se itens avulsos têm profissional definido. |
| 🟡 **MÉDIO** | `finalizar_comanda_integrada.go` | **Data de Referência**. Usa data do fechamento, distorcendo relatórios de produtividade real. | Usar `appointment.Date` (se existir) como `ReferenceDate`, ou manter `Now()` apenas para vendas balcão. |

---

## 4. Integração (Fluxo Completo)

1.  **Agendamento**:
    *   Define o `ProfessionalID`.
    *   Não calcula comissão neste momento.
2.  **Comanda (Abertura/Edição)**:
    *   Adiciona serviços/produtos.
    *   Não pré-calcula comissão (cálculo é feito apenas no fechamento).
3.  **Comanda (Fechamento)**:
    *   O UseCase `FinalizarComandaIntegrada` é acionado.
    *   Verifica se o Caixa está aberto.
    *   Processa pagamentos (Gera `OperacaoCaixa` e `ContaReceber`).
    *   Abate estoque de produtos.
    *   **Gera Comissões (Apenas Serviços)**:
        *   Busca regra hierárquica.
        *   Calcula valor base (Bruto ou Líquido).
        *   Insere registro em `commission_items` com status `PENDENTE`.
4.  **Financeiro (Contas a Pagar)**:
    *   A comissão **não** gera uma `Conta a Pagar` imediata.
    *   Ela fica acumulada em `commission_items`.
    *   É necessário rodar o processo de **Fechamento de Período** (mensal/quinzenal) para agrupar esses itens e gerar uma única `Conta a Pagar` para o profissional.

---

## 5. Análise do Dashboard

O dashboard utiliza as queries do arquivo `commission_items.sql`.

**Por que os dados podem estar incorretos?**
1.  **Falta de Produtos**: Se o usuário espera ver comissão de produtos, o dashboard mostrará valores menores, pois esses registros nunca são criados.
2.  **Filtro de Data**: As queries usam `reference_date`. Se o dashboard filtrar por "Data do Agendamento" mas o backend salvou a "Data do Fechamento" (que pode ser diferente), haverá divergência.
3.  **Status**: O dashboard filtra corretamente `status != 'CANCELADO'` e `status != 'ESTORNADO'`. Porém, se uma comanda for reaberta ou cancelada sem atualizar o status dos itens de comissão (o que parece ser tratado, mas requer atenção), pode haver "sujeira".

**Queries do Dashboard**:
*   `GetCommissionSummaryByService`: Agrupa por serviço.
*   `SumCommissionsByProfessionalAndDateRange`: Agrupa por profissional.

**Conclusão do Dashboard**: As queries estão tecnicamente corretas (SQL), mas a **alimentação dos dados** (no fechamento da comanda) está incompleta (falta produtos) e potencialmente imprecisa (data de referência).

---

## 6. Checklist de Correções Recomendadas

- [x] **Implementar Comissão de Produtos**: Adicionada lógica em `FinalizarComandaIntegradaUseCase` para calcular e salvar comissão de produtos (`buscarRegraComissaoProduto` e `processarComissaoProduto`).
- [x] **Ativar Regras por Categoria**: Integrada a tabela `comissoes_categoria_profissional` na hierarquia de busca de regras (agora são 5 níveis: Serviço → Categoria → Profissional → Unidade → Global).
- [x] **Revisar Data de Referência**: Alterada para usar `appointment.StartTime` (data do agendamento) como `ReferenceDate` em vez da data atual.
- [ ] **Suporte a Múltiplos Profissionais**: Garantir que itens adicionados manualmente na comanda possam ter um profissional diferente do agendamento principal. (Requer alteração de modelo)
- [ ] **Auditoria de Comissões**: Criar um script para recalcular comissões passadas (especialmente de produtos) que não foram geradas.

---

## 7. Alterações Implementadas (10/12/2025)

### Arquivos Modificados

1. **`backend/internal/infra/db/queries/appointments.sql`**
   - Query `GetServiceInfo` agora retorna `categoria_id` para suportar busca de comissão por categoria.

2. **`backend/internal/infra/db/queries/professionals.sql`**
   - Nova query `GetProfessionalCategoryCommission` para buscar comissão específica por categoria de serviço.

3. **`backend/internal/domain/port/appointment_repository.go`**
   - `ServiceInfo` agora inclui `CategoriaID`.
   - `ProfessionalReader` agora inclui método `GetCategoryCommission`.

4. **`backend/internal/infra/repository/postgres/readers.go`**
   - Implementação de `GetCategoryCommission` no `ProfessionalReaderPG`.
   - `FindByID` do `ServiceReaderPG` agora retorna `CategoriaID`.

5. **`backend/internal/application/usecase/command/finalizar_comanda_integrada.go`**
   - Hierarquia de comissões expandida para 5 níveis (incluindo Categoria).
   - Nova função `buscarRegraComissaoProduto` para produtos.
   - Nova função `processarComissaoProduto` para gerar comissão de produtos.
   - Nova função `processarComissaoServicoHierarquicaComData` que usa `appointmentDate`.
   - `ReferenceDate` agora usa a data do agendamento quando disponível.

### Nova Hierarquia de Comissões (5 Níveis)

1. **Serviço**: Comissão específica cadastrada no serviço.
2. **Categoria**: Comissão do profissional para a categoria do serviço (tabela `comissoes_categoria_profissional`).
3. **Profissional**: Comissão padrão do profissional.
4. **Unidade**: Regra de comissão da unidade (tabela `commission_rules`).
5. **Global**: Regra de comissão global do tenant.
