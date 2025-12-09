# PRD — Módulo Financeiro | NEXO v1.0

**Versão do Documento:** 1.0.0  
**Status:** 🔴 PRONTO PARA IMPLEMENTAÇÃO  
**Prioridade:** 🔴 CRÍTICA (MVP)  
**Data de Criação:** 28/11/2025  
**Última Atualização:** 28/11/2025  
**Responsável:** Andrey Viana (Product Owner)  
**Milestone:** MVP v1.0.0  

---

## 📊 Status de Implementação

| Área | Status | Progresso |
|------|--------|-----------|
| Backend - Despesas Fixas | ❌ Não iniciado | 0% |
| Backend - Painel Mensal | ❌ Não iniciado | 0% |
| Backend - Projeções | ❌ Não iniciado | 0% |
| Frontend - Tela Contas Fixas | ❌ Não iniciado | 0% |
| Frontend - Painel Mensal | ❌ Não iniciado | 0% |
| Cron - Gerador Automático | ❌ Não iniciado | 0% |

### ⏳ Pendente
- [ ] Tabela `despesas_fixas`
- [ ] CRUD de Despesas Fixas
- [ ] Gerador automático de contas mensais
- [ ] Endpoint de Painel Mensal
- [ ] Cálculo de projeções
- [ ] Dashboard frontend
- [ ] Tela de Contas Fixas

---

## 1. Executive Summary

### 1.1 Visão Geral

O **Módulo Financeiro** é um diferencial estratégico do NEXO, projetado para transformar dados de receitas, despesas, comissões e assinaturas em **inteligência financeira acionável** para donos de barbearia.

**Problema:** Donos de barbearia NÃO conseguem:
- ❌ Ver quanto precisam faturar por mês para ficar no lucro
- ❌ Saber como está o resultado financeiro atual
- ❌ Controlar despesas fixas de forma organizada
- ❌ Ver quanto falta para fechar o mês no azul
- ❌ Acompanhar previsões de caixa com base em vendas reais
- ❌ Tomar decisões baseadas em dados, não achismo

**Contexto Técnico:**  
O NEXO já possui toda a infraestrutura backend para:
- ✅ Contas a pagar/receber
- ✅ Fluxo de caixa diário
- ✅ DRE (Demonstração do Resultado do Exercício)
- ✅ Snapshot diário de caixa
- ✅ Comissões automáticas
- ✅ Assinaturas recorrentes

**O que falta:** Camada de análise mensal + controle de despesas fixas + projeções inteligentes.

### 1.2 Solução

Criar **3 novos componentes** integrados ao ecossistema financeiro existente:

#### **A) Gestão de Contas Fixas (Recorrentes)**
Tela onde o dono cadastra despesas mensais fixas:
- Aluguel
- Internet, Água, Energia
- Sistemas (NEXO, POS, etc)
- Contador
- Faxina/Limpeza
- Salários e Benefícios
- Outras

✨ **Automação:** Cada despesa gera automaticamente uma conta a pagar todo dia 1º do mês.

#### **B) Painel Financeiro Mensal (Dashboard)**
Dashboard completo com:
- 💰 Total faturado no mês (Serviços + Produtos + Assinaturas)
- 🎯 Meta mensal e % atingida
- 📊 Quanto falta faturar
- 🔴 Despesas Fixas totais
- 🟠 Despesas Variáveis (insumos, comissões, manutenções)
- 🟢 Lucro Operacional até agora

#### **C) Projeção Financeira (Até o Final do Mês)**
O sistema calcula automaticamente:
- 📈 Receita projetada até o último dia
- 💵 Lucro/prejuízo previsto
- 🔮 Probabilidade de bater a meta

**Base de cálculo:**
- Assinaturas confirmadas
- Média diária de faturamento
- Movimento histórico dos últimos 30 dias

---

## 2. Diferencial Competitivo

### 2.1 Comparação com Concorrentes

| Funcionalidade | NEXO | Trinks | AppBarber | BarberSystem |
|----------------|------|--------|-----------|--------------|
| **Painel Financeiro Mensal** | ✅ | ❌ | ❌ | ❌ |
| **Meta Automática Inteligente** | ✅ | ❌ | ❌ | ❌ |
| **Projeção de Lucro** | ✅ | ❌ | ❌ | ❌ |
| **Despesas Fixas Recorrentes** | ✅ | ❌ | 🟡 Parcial | ❌ |
| **Integração Assinaturas + DRE** | ✅ | ❌ | ❌ | ❌ |
| **Análise em Tempo Real** | ✅ | ❌ | ❌ | ❌ |

🏆 **Só o NEXO oferece:**
1. Meta mensal automática (Fixo + Variável + Margem desejada)
2. Projeção de lucro até o fim do mês
3. Dashboard financeiro em tempo real
4. Conexão de assinaturas + receitas + despesas + comissões

---

## 3. Objetivos do Produto

### 3.1 Objetivo Principal

**Permitir que donos de barbearia entendam a saúde financeira do negócio em uma única tela, com projeções confiáveis e metas inteligentes.**

### 3.2 Objetivos Secundários

1. **Reduzir achismo** nas decisões financeiras (meta: 100% decisões baseadas em dados)
2. **Aumentar consciência de custos** (meta: 80% dos donos sabem seu ponto de equilíbrio)
3. **Melhorar previsibilidade** de caixa (meta: < 10% de desvio na projeção)
4. **Automatizar controle** de despesas fixas (meta: 0% de esquecimento de lançamentos)
5. **Aumentar taxa de permanência** no sistema (donos não cancelam porque veem valor)

---

## 4. Métricas de Sucesso (KPIs)

| KPI | Baseline | Meta | Medição |
|-----|----------|------|---------|
| **Acurácia da Projeção** | N/A | > 90% | (Projetado - Real) / Real × 100 |
| **Uso Diário do Painel** | N/A | > 60% | Sessões diárias com acesso ao painel |
| **Contas Fixas Automatizadas** | 0% | 100% | % de contas geradas automaticamente |
| **Taxa de Permanência (Churn)** | 15% | < 5% | Cancelamentos / Total de clientes |
| **NPS do Módulo Financeiro** | N/A | > 8.5 | Pesquisa de satisfação |

---

## 5. Personas e Necessidades

### 5.1 Persona 1: Dono da Barbearia

**Nome:** Carlos, 38 anos, Dono de 2 barbearias  
**Objetivo:** Maximizar lucro e ter controle financeiro total  

**Necessidades:**
- 🔴 Saber se vai ter lucro no final do mês
- 🔴 Ver quanto precisa faturar para cobrir custos
- 🔴 Controlar despesas fixas sem esquecer nenhuma
- 🟡 Projetar resultado financeiro com base em dados reais
- 🟡 Comparar desempenho entre meses

**Pain Points:**
- "Não sei se estou no lucro ou prejuízo até fechar o mês"
- "Esqueço de lançar aluguel, energia, contador"
- "Não sei quanto preciso vender para pagar tudo"
- "Trabalho muito mas não sobra dinheiro"

**Como o NEXO resolve:**
- ✅ Painel mostra lucro/prejuízo em tempo real
- ✅ Despesas fixas são lançadas automaticamente
- ✅ Meta inteligente calcula quanto precisa faturar
- ✅ Projeção antecipa o resultado do mês

---

### 5.2 Persona 2: Gerente Financeiro

**Nome:** Juliana, 32 anos, Gerente de rede com 4 unidades  
**Objetivo:** Manter todas as unidades no azul  

**Necessidades:**
- 🔴 Comparar desempenho financeiro entre unidades
- 🔴 Identificar quais unidades estão perdendo dinheiro
- 🟡 Acompanhar evolução mensal de custos
- 🟡 Exportar dados para apresentar ao dono

**Pain Points:**
- "Cada unidade tem uma planilha diferente"
- "Não consigo consolidar dados financeiros"
- "Perco tempo fazendo relatórios manuais"

**Como o NEXO resolve:**
- ✅ Dashboard consolidado de todas as unidades
- ✅ Comparação lado a lado do desempenho
- ✅ Exportação automática de relatórios
- ✅ Alertas quando unidade está abaixo da meta

---

### 5.3 Persona 3: Contador

**Nome:** Roberto, 45 anos, Contador de 12 barbearias  
**Objetivo:** Receber dados organizados para fechamento contábil  

**Necessidades:**
- 🔴 Exportar DRE mensal automaticamente
- 🔴 Ver todas as despesas lançadas
- 🟡 Categorização correta de custos
- 🟡 Acesso read-only ao financeiro

**Pain Points:**
- "Donos de barbearia não organizam despesas"
- "Recebo dados bagunçados no final do mês"
- "Perco tempo categorizando tudo manualmente"

**Como o NEXO resolve:**
- ✅ Despesas já categorizadas corretamente
- ✅ DRE gerado automaticamente
- ✅ Exportação em formato padronizado
- ✅ Acesso direto via conta de contador

---

## 6. Regras de Negócio (RN)

### 6.1 Despesas Fixas Recorrentes

| ID | Regra | Criticidade |
|----|-------|-------------|
| **RN-FX-001** | Despesas fixas DEVEM gerar lançamentos automáticos todo dia 1º do mês | 🔴 Crítica |
| **RN-FX-002** | Lançamentos gerados podem ser editados individualmente | 🟡 Média |
| **RN-FX-003** | Editar despesa fixa NÃO afeta lançamentos já criados | 🔴 Crítica |
| **RN-FX-004** | Deletar despesa fixa NÃO deleta lançamentos já criados | 🔴 Crítica |
| **RN-FX-005** | Despesa fixa pode ser temporariamente desabilitada | 🟡 Média |
| **RN-FX-006** | Categoria da despesa fixa DEVE ser validada | 🟡 Média |

### 6.2 Painel Mensal

| ID | Regra | Criticidade |
|----|-------|-------------|
| **RN-PNL-001** | Meta mensal pode ser manual OU automática | 🔴 Crítica |
| **RN-PNL-002** | Meta automática = Despesas Fixas + Projeção Variável + Margem Desejada | 🔴 Crítica |
| **RN-PNL-003** | Painel DEVE usar regime de competência, não caixa | 🔴 Crítica |
| **RN-PNL-004** | Comissões entram como despesa operacional | 🔴 Crítica |
| **RN-PNL-005** | Faturamento inclui: Serviços + Produtos + Assinaturas | 🔴 Crítica |
| **RN-PNL-006** | Painel atualiza em tempo real a cada lançamento | 🟡 Média |

### 6.3 Projeções

| ID | Regra | Criticidade |
|----|-------|-------------|
| **RN-PRJ-001** | Projeção DEVE recalcular diariamente às 00:00 | 🔴 Crítica |
| **RN-PRJ-002** | Projeção considera assinaturas confirmadas do mês | 🔴 Crítica |
| **RN-PRJ-003** | Projeção usa média móvel dos últimos 7 dias | 🟡 Média |
| **RN-PRJ-004** | Projeção considera sazonalidade (fim de semana > dias úteis) | 🟢 Baixa |
| **RN-PRJ-005** | Projeção DEVE mostrar cenário otimista e pessimista | 🟢 Baixa |

### 6.4 Multi-Tenant

| ID | Regra | Criticidade |
|----|-------|-------------|
| **RN-MT-001** | Todas as despesas fixas DEVEM ter tenant_id | 🔴 Crítica |
| **RN-MT-002** | Painel mensal filtra apenas dados do tenant ativo | 🔴 Crítica |
| **RN-MT-003** | Unidades diferentes podem ter despesas fixas diferentes | 🟡 Média |

---

## 7. Requisitos Funcionais (RF)

### 7.1 Gestão de Despesas Fixas

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-001** | Sistema DEVE permitir criar despesa fixa recorrente | 🔴 P0 | ❌ |
| **RF-002** | Sistema DEVE permitir editar despesa fixa | 🔴 P0 | ❌ |
| **RF-003** | Sistema DEVE permitir deletar despesa fixa | 🔴 P0 | ❌ |
| **RF-004** | Sistema DEVE listar todas as despesas fixas ativas | 🔴 P0 | ❌ |
| **RF-005** | Sistema DEVE desabilitar despesa fixa temporariamente | 🟡 P1 | ❌ |
| **RF-006** | Sistema DEVE categorizar despesas (predefinido) | 🔴 P0 | ❌ |
| **RF-007** | Sistema DEVE validar dia de vencimento (1-31) | 🟡 P1 | ❌ |

### 7.2 Geração Automática de Contas

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-008** | Sistema DEVE gerar contas a pagar automaticamente dia 1º | 🔴 P0 | ❌ |
| **RF-009** | Geração DEVE criar uma conta a pagar para cada despesa fixa ativa | 🔴 P0 | ❌ |
| **RF-010** | Conta gerada DEVE ter status PENDENTE | 🔴 P0 | ❌ |
| **RF-011** | Conta gerada DEVE ter vencimento = dia configurado | 🔴 P0 | ❌ |
| **RF-012** | Sistema DEVE registrar log de geração automática | 🟡 P1 | ❌ |
| **RF-013** | Sistema NÃO DEVE duplicar contas se rodar 2x no mesmo dia | 🔴 P0 | ❌ |

### 7.3 Painel Financeiro Mensal

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-014** | Sistema DEVE exibir total faturado no mês | 🔴 P0 | ❌ |
| **RF-015** | Sistema DEVE exibir meta mensal e % atingida | 🔴 P0 | ❌ |
| **RF-016** | Sistema DEVE exibir quanto falta para atingir meta | 🔴 P0 | ❌ |
| **RF-017** | Sistema DEVE exibir total de despesas fixas | 🔴 P0 | ❌ |
| **RF-018** | Sistema DEVE exibir total de despesas variáveis | 🔴 P0 | ❌ |
| **RF-019** | Sistema DEVE exibir total de comissões | 🔴 P0 | ❌ |
| **RF-020** | Sistema DEVE calcular lucro operacional | 🔴 P0 | ❌ |
| **RF-021** | Sistema DEVE exibir gráfico de faturamento diário | 🟡 P1 | ❌ |
| **RF-022** | Sistema DEVE permitir alternar entre meses | 🔴 P0 | ❌ |

### 7.4 Meta Mensal

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-023** | Sistema DEVE permitir definir meta manual | 🔴 P0 | ❌ |
| **RF-024** | Sistema DEVE calcular meta automática inteligente | 🔴 P0 | ❌ |
| **RF-025** | Meta automática = Fixo + Variável + Margem | 🔴 P0 | ❌ |
| **RF-026** | Sistema DEVE permitir configurar margem desejada (%) | 🔴 P0 | ❌ |
| **RF-027** | Sistema DEVE exibir comparação: meta vs realizado | 🔴 P0 | ❌ |

### 7.5 Projeção Financeira

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-028** | Sistema DEVE calcular receita projetada até fim do mês | 🔴 P0 | ❌ |
| **RF-029** | Sistema DEVE calcular lucro projetado | 🔴 P0 | ❌ |
| **RF-030** | Sistema DEVE considerar assinaturas confirmadas | 🔴 P0 | ❌ |
| **RF-031** | Sistema DEVE usar média móvel de 7 dias | 🟡 P1 | ❌ |
| **RF-032** | Sistema DEVE exibir probabilidade de bater meta | 🟢 P2 | ❌ |
| **RF-033** | Sistema DEVE recalcular projeção diariamente | 🔴 P0 | ❌ |

### 7.6 Exportação e Relatórios

| ID | Requisito | Prioridade | Status |
|----|-----------|------------|--------|
| **RF-034** | Sistema DEVE permitir exportar dados em CSV | 🟡 P1 | ❌ |
| **RF-035** | Sistema DEVE permitir exportar DRE mensal | 🟡 P1 | ❌ |
| **RF-036** | Sistema DEVE permitir imprimir painel mensal | 🟢 P2 | ❌ |

6. 🖥️ Telas do MVP
1) Tela: Contas Fixas

Local: Sidebar → Financeiro → Contas Fixas

Componentes:

Lista de contas fixas

Botão criar

Modal editar

Modal deletar

Toggle “recorrente mensal”

Categoria (dropdown)

Valor

Vencimento

Método de pagamento

2) Tela: Painel Financeiro do Mês

Local: Financeiro → Painel Mensal

Blocos:
🔵 1. Faturamento do Mês

Serviços

Produtos

Assinaturas

Gráfico diário

🟡 2. Meta Mensal

Definir meta

Meta inteligente

Percentual alcançado

Quanto falta

🔴 3. Despesas Fixas

Tabela do mês

Total fixo

🟠 4. Despesas Variáveis

Insumos

Comissões

Manutenção

🟢 5. Resultado Operacional Atual
Receitas totais
- Comissões
- Fixas
- Variáveis
= Lucro/Prejuízo

🟣 6. Projeção Até o Final do Mês

Receita prevista

Lucro previsto

7. 🗄️ Estrutura Técnica (Backend)

O sistema já tem:

Contas a pagar / receber

Fluxo de caixa

DRE

Snapshot diário

Precisamos adicionar:

Nova tabela: despesas_fixas
CREATE TABLE despesas_fixas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    valor DECIMAL(15,2) NOT NULL,
    categoria VARCHAR(100) NOT NULL,
    recorrente BOOLEAN DEFAULT true,
    dia_vencimento INT NOT NULL,
    metodo_pagamento VARCHAR(50),
    unidade_id UUID REFERENCES units(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

Gerador automático no backend

Cron diário:

Se hoje for dia 1:
    Para cada despesa_fixa:
        Criar uma conta a pagar com status PENDENTE

8. 📡 Endpoints Necessários
---

## 8. Requisitos Não Funcionais (RNF)

### 8.1 Performance

| ID | Requisito | Meta | Prioridade |
|----|-----------|------|------------|
| **RNF-001** | Painel mensal DEVE carregar em < 2s | < 2000ms | 🔴 Alta |
| **RNF-002** | Geração automática DEVE processar 1000 despesas em < 30s | < 30s | 🟡 Média |
| **RNF-003** | Projeção DEVE calcular em < 1s | < 1000ms | 🟡 Média |
| **RNF-004** | Exportação CSV DEVE completar em < 5s | < 5000ms | 🟢 Baixa |

### 8.2 Disponibilidade

| ID | Requisito | Meta | Prioridade |
|----|-----------|------|------------|
| **RNF-005** | Sistema DEVE ter uptime > 99.5% | > 99.5% | 🔴 Alta |
| **RNF-006** | Gerador automático DEVE ter retry em caso de falha | 3 tentativas | 🔴 Alta |
| **RNF-007** | Cron DEVE ter monitoramento e alertas | 100% | 🟡 Média |

### 8.3 Segurança

| ID | Requisito | Meta | Prioridade |
|----|-----------|------|------------|
| **RNF-008** | Todas as rotas DEVEM validar tenant_id | 100% | 🔴 Alta |
| **RNF-009** | Dados financeiros DEVEM ser criptografados em trânsito | TLS 1.3 | 🔴 Alta |
| **RNF-010** | Acesso ao painel DEVE ser logado (audit log) | 100% | 🟡 Média |

### 8.4 Usabilidade

| ID | Requisito | Meta | Prioridade |
|----|-----------|------|------------|
| **RNF-011** | Painel DEVE ser responsivo (mobile + desktop) | 100% | 🔴 Alta |
| **RNF-012** | Valores DEVEM ser formatados em BRL | R$ 1.234,56 | 🔴 Alta |
| **RNF-013** | Cores DEVEM seguir Design System | 100% | 🟡 Média |

---

## 9. Arquitetura e Modelo de Dados

### 9.1 Nova Tabela: `despesas_fixas`

```sql
CREATE TABLE despesas_fixas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unidade_id UUID REFERENCES units(id) ON DELETE CASCADE,
    
    -- Dados da despesa
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    valor DECIMAL(15,2) NOT NULL CHECK (valor >= 0),
    categoria VARCHAR(100) NOT NULL,
    
    -- Recorrência
    recorrente BOOLEAN DEFAULT true NOT NULL,
    dia_vencimento INT NOT NULL CHECK (dia_vencimento BETWEEN 1 AND 31),
    
    -- Configurações
    metodo_pagamento VARCHAR(50),
    ativo BOOLEAN DEFAULT true NOT NULL,
    
    -- Auditoria
    criado_em TIMESTAMP DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMP DEFAULT NOW() NOT NULL,
    criado_por UUID REFERENCES users(id),
    atualizado_por UUID REFERENCES users(id),
    
    -- Índices
    CONSTRAINT fk_despesa_fixa_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_despesa_fixa_unidade FOREIGN KEY (unidade_id) REFERENCES units(id)
);

-- Índices para performance
CREATE INDEX idx_despesas_fixas_tenant ON despesas_fixas(tenant_id);
CREATE INDEX idx_despesas_fixas_ativo ON despesas_fixas(tenant_id, ativo);
CREATE INDEX idx_despesas_fixas_unidade ON despesas_fixas(unidade_id);

-- RLS (Row Level Security)
ALTER TABLE despesas_fixas ENABLE ROW LEVEL SECURITY;

CREATE POLICY despesas_fixas_tenant_isolation ON despesas_fixas
    USING (tenant_id = current_setting('app.current_tenant')::uuid);
```

### 9.2 Categorias Predefinidas

```typescript
enum CategoriaDespesaFixa {
  ALUGUEL = 'ALUGUEL',
  CONDOMINIO = 'CONDOMINIO',
  ENERGIA = 'ENERGIA',
  AGUA = 'AGUA',
  INTERNET = 'INTERNET',
  TELEFONE = 'TELEFONE',
  SISTEMAS = 'SISTEMAS',           // NEXO, POS, etc
  CONTADOR = 'CONTADOR',
  LIMPEZA = 'LIMPEZA',
  SEGURANCA = 'SEGURANCA',
  SALARIOS = 'SALARIOS',
  BENEFICIOS = 'BENEFICIOS',       // Vale transporte, alimentação
  MARKETING = 'MARKETING',
  MANUTENCAO = 'MANUTENCAO',
  SEGUROS = 'SEGUROS',
  IMPOSTOS = 'IMPOSTOS',
  OUTRAS = 'OUTRAS'
}
```

### 9.3 Integração com Tabelas Existentes

**Fluxo:**
1. `despesas_fixas` → Cadastro manual pelo dono
2. **Cron (dia 1º)** → Gera `contas_pagar` com status PENDENTE
3. Quando paga → Atualiza `fluxo_caixa_diario`
4. `fluxo_caixa_diario` → Alimenta `dre_mensal`
5. `painel_mensal` → Consome `dre_mensal` + projeções

---

## 10. Endpoints da API

### 10.1 Despesas Fixas

#### `POST /api/v1/financeiro/despesas-fixas`
Criar nova despesa fixa

**Request:**
```json
{
  "nome": "Aluguel Loja Centro",
  "descricao": "Aluguel mensal da unidade centro",
  "valor": "8500.00",
  "categoria": "ALUGUEL",
  "dia_vencimento": 10,
  "metodo_pagamento": "TRANSFERENCIA",
  "unidade_id": "uuid-opcional",
  "ativo": true
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "tenant_id": "uuid",
  "nome": "Aluguel Loja Centro",
  "valor": "8500.00",
  "categoria": "ALUGUEL",
  "dia_vencimento": 10,
  "recorrente": true,
  "ativo": true,
  "criado_em": "2025-11-28T10:00:00Z"
}
```

---

#### `GET /api/v1/financeiro/despesas-fixas`
Listar despesas fixas

**Query Params:**
- `ativo` (boolean): filtrar por status
- `categoria` (string): filtrar por categoria
- `unidade_id` (uuid): filtrar por unidade

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "nome": "Aluguel Loja Centro",
      "valor": "8500.00",
      "categoria": "ALUGUEL",
      "dia_vencimento": 10,
      "ativo": true
    }
  ],
  "total": 12,
  "total_mensal": "24500.00"
}
```

---

#### `PUT /api/v1/financeiro/despesas-fixas/:id`
Atualizar despesa fixa

**Request:**
```json
{
  "nome": "Aluguel Loja Centro - Atualizado",
  "valor": "9000.00"
}
```

**Response:** `200 OK`

---

#### `DELETE /api/v1/financeiro/despesas-fixas/:id`
Deletar despesa fixa

**Response:** `204 No Content`

⚠️ **Importante:** NÃO deleta contas a pagar já geradas

---

### 10.2 Painel Mensal

#### `GET /api/v1/financeiro/painel-mensal/:mes/:ano`
Retorna dashboard completo do mês

**Exemplo:** `/api/v1/financeiro/painel-mensal/11/2025`

**Response:** `200 OK`
```json
{
  "mes": 11,
  "ano": 2025,
  "periodo": "2025-11-01 a 2025-11-30",
  
  "faturamento": {
    "total": "41500.00",
    "servicos": "28000.00",
    "produtos": "8500.00",
    "assinaturas": "5000.00",
    "detalhamento_diario": [
      {"dia": 1, "valor": "1200.00"},
      {"dia": 2, "valor": "1850.00"}
    ]
  },
  
  "meta": {
    "valor": "60000.00",
    "tipo": "AUTOMATICA",
    "porcentagem_atingida": 69.17,
    "falta_faturar": "18500.00",
    "base_calculo": {
      "despesas_fixas": "24500.00",
      "despesas_variaveis_estimadas": "15000.00",
      "margem_desejada": "20500.00"
    }
  },
  
  "despesas": {
    "fixas": {
      "total": "24500.00",
      "itens": [
        {"categoria": "ALUGUEL", "valor": "8500.00", "quantidade": 1},
        {"categoria": "ENERGIA", "valor": "1200.00", "quantidade": 1}
      ]
    },
    "variaveis": {
      "total": "7800.00",
      "insumos": "3200.00",
      "manutencao": "4600.00"
    },
    "comissoes": {
      "total": "12000.00",
      "por_barbeiro": [
        {"barbeiro_id": "uuid", "nome": "João Silva", "valor": "4500.00"}
      ]
    }
  },
  
  "resultado": {
    "lucro_operacional": "18500.00",
    "margem": 44.58,
    "status": "POSITIVO"
  },
  
  "projecao": {
    "receita_projetada": "52700.00",
    "lucro_projetado": "21200.00",
    "probabilidade_bater_meta": 75,
    "dias_restantes": 15,
    "media_diaria_necessaria": "1233.33",
    "media_diaria_atual": "1383.33",
    "cenarios": {
      "otimista": "58000.00",
      "realista": "52700.00",
      "pessimista": "48000.00"
    }
  }
}
```

---

## 11. Diferenciais Competitivos

### 11.1 Comparação com Concorrentes

| Funcionalidade | NEXO | Trinks | AppBarber | BarberSystem |
|----------------|------|--------|-----------|--------------|
| **Painel Financeiro Mensal** | ✅ | ❌ | ❌ | ❌ |
| **Meta Automática Inteligente** | ✅ | ❌ | ❌ | ❌ |
| **Projeção de Lucro** | ✅ | ❌ | ❌ | ❌ |
| **Despesas Fixas Recorrentes** | ✅ | ❌ | 🟡 Parcial | ❌ |
| **Integração Assinaturas + DRE** | ✅ | ❌ | ❌ | ❌ |
| **Análise em Tempo Real** | ✅ | ❌ | ❌ | ❌ |

🏆 **Só o NEXO oferece:**
1. Meta mensal automática (Fixo + Variável + Margem desejada)
2. Projeção de lucro até o fim do mês
3. Dashboard financeiro em tempo real
4. Conexão de assinaturas + receitas + despesas + comissões

---

## 12. Conclusão

O **Módulo Financeiro** do NEXO representa um diferencial competitivo estratégico que nenhum concorrente possui. Ao transformar dados brutos em inteligência financeira acionável, o NEXO se posiciona como o único ERP completo para barbearias premium.

### Impacto Esperado

**Para o Negócio:**
- 📈 Aumento de 50% na taxa de permanência (redução de churn)
- 💰 Aumento do LTV (Lifetime Value) dos clientes
- 🎯 Posicionamento como solução premium diferenciada

**Para os Usuários:**
- 📊 100% dos donos entendem sua saúde financeira
- 💡 Decisões baseadas em dados, não achismo
- ⚡ Economia de 10h/mês em controle financeiro manual
- 🔮 Previsibilidade e controle do resultado mensal

---

**Documento Vivo:** Este PRD será atualizado conforme o desenvolvimento avança.  
**Última Atualização:** 28/11/2025  
**Próxima Revisão:** 05/12/2025

---

## Referências

- [PRD Principal NEXO](../../PRD-VALTARIS.md)
- [Arquitetura Backend](../../docs/04-backend/GUIA_DEV_BACKEND.md)
- [Design System](../../docs/03-frontend/DESIGN_SYSTEM.md)
- [Modelo de Dados](../../docs/02-arquitetura/MODELO_DE_DADOS.md)
- [Fluxos Críticos](../../docs/02-arquitetura/FLUXOS_CRITICOS_SISTEMA.md)