# Plano de Ajuste Urgente: Isolamento Total de Unidades (Multi-Unit Strict Mode)

**Data:** 11/12/2025
**Prioridade:** CRÍTICA (Bloqueador de Negócio)
**Objetivo:** Garantir que cada Unidade opere como uma entidade isolada (Silo de Dados), impedindo vazamento de informações entre barbearias/filiais distintas dentro do mesmo Tenant.

---

## 1. Contexto e Diagnóstico

A auditoria identificou que o sistema opera corretamente no nível **Multi-Tenant** (isolamento entre empresas/assinantes), mas falha no nível **Multi-Unit** (isolamento entre filiais/unidades da mesma empresa).

**Risco:** Em um modelo de Franquia ou Rede, um franqueado (Unidade A) consegue visualizar profissionais e serviços de outro franqueado (Unidade B), violando a regra de negócio de "100% Separada".

**Causas Raiz:**
1.  **Banco de Dados:** Tabelas `categories_services` e `services` não possuem coluna `unit_id`.
2.  **Backend:** Handlers ignoram o header `X-Unit-ID` e filtram apenas por `tenant_id`.
3.  **Infraestrutura:** Backend precisa de reinício para aplicar regras de CORS (causando erros no frontend).

---

## 2. Plano de Ação Imediato (Correção de Erros)

### 2.1. Infraestrutura (CORS)
*   **Ação:** Reiniciar o serviço backend Go imediatamente.
*   **Comando:** `systemctl restart nexo-backend` (ou equivalente no ambiente de deploy).
*   **Resultado Esperado:** Eliminação dos erros `AxiosError` (CORS) no frontend.

### 2.2. Correção de Dados (Profissionais)
*   **Problema:** Profissionais da unidade "Mangabeiras" estão com `unit_id = NULL`.
*   **Ação:** Executar script SQL de correção.
*   **Script:**
    ```sql
    -- Vincular profissionais órfãos à unidade Mangabeiras (ID da auditoria)
    UPDATE professionals 
    SET unit_id = '5ed7f5b4-5823-443f-b29d-286cc32a02e6' 
    WHERE tenant_id = 'SEU_TENANT_ID' 
      AND unit_id IS NULL;
    ```

---

## 3. Plano de Ação Estrutural (Isolamento Real)

Para atender ao requisito "Cada unidade pode ser uma barbearia diferente", implementaremos o **Isolamento Rígido (Strict Isolation)**.

### 3.1. Alteração de Schema (Banco de Dados)

Adicionar a coluna `unit_id` nas tabelas de catálogo para permitir propriedade exclusiva por unidade.

**Migration SQL (Sugestão):**

```sql
-- 1. Adicionar coluna unit_id em categorias
ALTER TABLE categories_services 
ADD COLUMN unit_id UUID REFERENCES units(id);

-- 2. Adicionar coluna unit_id em serviços
ALTER TABLE services 
ADD COLUMN unit_id UUID REFERENCES units(id);

-- 3. Criar índices para performance
CREATE INDEX idx_categories_unit ON categories_services(unit_id);
CREATE INDEX idx_services_unit ON services(unit_id);
```

> **Nota de Decisão:** Optamos por adicionar `unit_id` diretamente (Opção A da auditoria) em vez de tabela de ligação, pois isso garante que um serviço pertença *exclusivamente* a uma unidade, facilitando a gestão independente de preços e comissões ("Cada unidade é uma barbearia diferente").

### 3.2. Refatoração do Backend (Go)

Alterar os Handlers e Repositórios para exigir e filtrar por `unit_id`.

#### A. Middleware & Contexto
*   Garantir que o `unit_id` extraído do header `X-Unit-ID` seja obrigatório para rotas operacionais.

#### B. Atualização de Queries (SQLC)
Alterar as queries em `internal/infra/db/queries/` para incluir o filtro de unidade.

**Exemplo (Services):**
```sql
-- Antes
SELECT * FROM services WHERE tenant_id = $1;

-- Depois (Isolamento Rígido)
SELECT * FROM services 
WHERE tenant_id = $1 
  AND (unit_id = $2 OR unit_id IS NULL); -- IS NULL permite ver "Padrões da Rede" se desejado, remover se for isolamento total.
```
*Recomendação:* Manter `OR unit_id IS NULL` apenas se houver "Serviços Globais" definidos pela Matriz. Caso contrário, remover para isolamento total.

#### C. Atualização de Handlers
*   **Arquivo:** `internal/api/handlers/professional_handler.go`
*   **Ação:** Ler `unitID` do contexto (`c.Get("unit_id")`) e passar para o UseCase/Repository.
*   **Repetir para:** `category_handler.go`, `service_handler.go`.

---

## 4. Status de Execução (Atualizado)

### ✅ Concluído
1.  **Banco de Dados**:
    *   [x] Adicionada coluna `unit_id` nas tabelas `appointments` e `services`.
    *   [x] Atualizadas queries SQL (`internal/infra/db/queries/`) para filtrar por `unit_id`.
    *   [x] Regenerado código Go com `sqlc generate`.

2.  **Backend (Camada de Domínio e Dados)**:
    *   [x] Atualizadas interfaces (`port`) de `Appointment` e `Servico` para aceitar `unitID`.
    *   [x] Atualizados repositórios (`postgres`) de `Appointment` e `Servico` para implementar a nova assinatura.
    *   [x] Corrigidos erros de compilação e tipos (UUID vs String).

3.  **Backend (Camada de Aplicação)**:
    *   [x] Atualizados UseCases de `Appointment` e `Servico` para propagar `unitID`.
    *   [x] Build do projeto (`go build ./...`) executado com sucesso.

4.  **Backend (Camada HTTP/Handlers + Enforcement)**:
    *   [x] Aplicado `UnitMiddleware` nas rotas críticas (ex.: `/appointments`, `/servicos`, `/professionals`, `/categorias-servicos`).
    *   [x] Atualizado `AppointmentHandler` para exigir `unit_id` do contexto e propagar para os UseCases.
    *   [x] Atualizado `ServicoHandler` para exigir `unit_id` do contexto e sobrescrever `unit_id` do payload (segurança).

5.  **Validação Automatizada (Backend)**:
    *   [x] Suíte de testes do backend passou (`go test ./...`).

### 🚧 Pendente / A Fazer
1.  **Backend (Camada HTTP/Handlers)**:
    *   [ ] Verificar e atualizar `CategoryHandler` e `ProfessionalHandler` (mesmo processo de isolamento: exigir `unit_id` do contexto e filtrar por unidade).

2.  **Infraestrutura**:
    *   [ ] Reiniciar serviço backend para aplicar mudanças.

3.  **Dados**:
    *   [ ] Rodar script SQL para vincular profissionais órfãos às unidades corretas.

4.  **Validação**:
    *   [ ] Teste manual: Criar registro na Unidade A e verificar invisibilidade na Unidade B.
    *   [ ] Teste manual: tentar acessar rotas operacionais sem `X-Unit-ID` e validar bloqueio.

---

## 5. Próximos Passos (Execução)

1.  [x] Criar migration SQL (`alter_tables_add_unit_id`).
2.  [x] Atualizar queries SQLC e regenerar código Go.
3.  [x] Atualizar Handlers para injetar `unit_id` (Appointment + Servicos) e aplicar enforcement via middleware.
4.  [ ] Revisar `CategoryHandler` e `ProfessionalHandler` (garantir filtro por `unit_id` em todas as operações).
5.  [ ] Deploy e Reinício do Backend.
6.  [ ] Rodar script de correção de dados (vincular dados legados).
7.  [ ] Validação manual cross-unit (A não enxerga B e vice-versa).
