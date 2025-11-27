# 📋 Backlog — Módulo de Serviços

> User stories, épicos e tarefas técnicas organizadas por prioridade

---

## 🎯 Épicos

### E-SRV-01: Categorias de Serviço
**Objetivo:** Permitir organização de serviços em categorias

**User Stories:**
- Como gerente, quero criar categorias para organizar meus serviços
- Como gerente, quero editar categorias existentes
- Como gerente, quero deletar categorias não utilizadas
- Como sistema, quero impedir deleção de categorias com serviços vinculados

**Tarefas Técnicas:**
- [ ] T-SRV-001: Criar migration de categorias
- [ ] T-SRV-002: Criar entidade Categoria (domain)
- [ ] T-SRV-003: Criar CategoriaRepository (infra)
- [ ] T-SRV-004: Criar DTOs de Categoria
- [ ] T-SRV-005: Criar Use Cases de Categoria (CRUD)
- [ ] T-SRV-006: Criar Handler de Categoria
- [ ] T-SRV-007: Criar rotas de Categoria
- [ ] T-SRV-008: Criar CategoryService (frontend)
- [ ] T-SRV-009: Criar hook useCategories
- [ ] T-SRV-010: Criar modal de Categoria
- [ ] T-SRV-011: Criar listagem de Categorias
- [ ] T-SRV-012: Testes unitários backend
- [ ] T-SRV-013: Testes integração backend
- [ ] T-SRV-014: Testes E2E frontend

**Estimativa:** 5 dias  
**Prioridade:** 🔴 Alta (bloqueante para serviços)

---

### E-SRV-02: Cadastro Básico de Serviços
**Objetivo:** Permitir criação e gestão de serviços

**User Stories:**
- Como gerente, quero cadastrar um novo serviço com nome, preço e duração
- Como gerente, quero vincular um serviço a uma categoria
- Como gerente, quero editar serviços existentes
- Como gerente, quero inativar serviços temporariamente
- Como gerente, quero deletar serviços não utilizados
- Como sistema, quero impedir deleção de serviços com agendamentos

**Tarefas Técnicas:**
- [ ] T-SRV-015: Criar migration de servicos
- [ ] T-SRV-016: Criar entidade Servico (domain)
- [ ] T-SRV-017: Criar ServicoRepository (infra)
- [ ] T-SRV-018: Criar DTOs de Servico
- [ ] T-SRV-019: Criar Use Cases de Servico (CRUD)
- [ ] T-SRV-020: Criar Handler de Servico
- [ ] T-SRV-021: Criar rotas de Servico
- [ ] T-SRV-022: Criar ServiceService (frontend)
- [ ] T-SRV-023: Criar hook useServices
- [ ] T-SRV-024: Criar página de Serviços
- [ ] T-SRV-025: Criar ServicesList component
- [ ] T-SRV-026: Criar ServiceModal component
- [ ] T-SRV-027: Implementar validações Zod
- [ ] T-SRV-028: Testes unitários backend
- [ ] T-SRV-029: Testes integração backend
- [ ] T-SRV-030: Testes E2E frontend

**Estimativa:** 8 dias  
**Prioridade:** 🔴 Alta

---

### E-SRV-03: Customização por Profissional
**Objetivo:** Permitir valores diferentes por profissional

**User Stories:**
- Como gerente, quero definir quais profissionais executam cada serviço
- Como gerente, quero customizar o preço de um serviço para um profissional específico
- Como gerente, quero customizar a duração de um serviço para um profissional específico
- Como gerente, quero customizar a comissão de um serviço para um profissional específico
- Como sistema, quero usar valores customizados nos agendamentos quando disponíveis
- Como sistema, quero usar valores padrão quando não houver customização

**Tarefas Técnicas:**
- [ ] T-SRV-031: Criar migration de servicos_profissionais
- [ ] T-SRV-032: Criar entidade ServicoProfissional (domain)
- [ ] T-SRV-033: Criar ServicoProfissionalRepository (infra)
- [ ] T-SRV-034: Criar DTOs de customização
- [ ] T-SRV-035: Criar Use Cases de customização
- [ ] T-SRV-036: Atualizar queries com COALESCE
- [ ] T-SRV-037: Criar ProfessionalCustomization component
- [ ] T-SRV-038: Integrar customização no ServiceModal
- [ ] T-SRV-039: Criar lógica de seleção de profissionais
- [ ] T-SRV-040: Implementar validações de customização
- [ ] T-SRV-041: Testes de queries otimizadas
- [ ] T-SRV-042: Testes unitários backend
- [ ] T-SRV-043: Testes integração backend
- [ ] T-SRV-044: Testes E2E frontend

**Estimativa:** 10 dias  
**Prioridade:** 🟡 Média (diferencial competitivo)

---

### E-SRV-04: Recursos Avançados
**Objetivo:** Melhorar UX e produtividade

**User Stories:**
- Como gerente, quero buscar serviços por nome
- Como gerente, quero filtrar serviços por categoria
- Como gerente, quero filtrar por status (ativo/inativo)
- Como gerente, quero duplicar um serviço existente
- Como gerente, quero fazer upload de imagem para o serviço
- Como gerente, quero adicionar tags para busca rápida
- Como sistema, quero ordenar serviços por nome, preço ou categoria

**Tarefas Técnicas:**
- [ ] T-SRV-045: Implementar busca fulltext
- [ ] T-SRV-046: Implementar filtros dinâmicos
- [ ] T-SRV-047: Implementar ordenação customizada
- [ ] T-SRV-048: Criar endpoint de duplicação
- [ ] T-SRV-049: Implementar upload de imagem
- [ ] T-SRV-050: Criar sistema de tags
- [ ] T-SRV-051: Criar componentes de filtro
- [ ] T-SRV-052: Criar SearchBar component
- [ ] T-SRV-053: Implementar debounce na busca
- [ ] T-SRV-054: Testes de performance
- [ ] T-SRV-055: Testes E2E completos

**Estimativa:** 7 dias  
**Prioridade:** 🟢 Baixa (nice to have)

---

## 📊 Priorização (MoSCoW)

### Must Have (Sprint 1.4.1 + 1.4.2)
- ✅ Cadastro de categorias
- ✅ Cadastro de serviços básicos
- ✅ Listagem de serviços
- ✅ Edição de serviços
- ✅ Controle de status (ativo/inativo)
- ✅ Validações de negócio
- ✅ Isolamento multi-tenant

### Should Have (Sprint 1.4.3)
- ✅ Customização por profissional
- ✅ Queries otimizadas com COALESCE
- ✅ UI de seleção de profissionais
- ✅ Validação de valores customizados

### Could Have (Sprint 1.4.4)
- 🔄 Busca por nome
- 🔄 Filtros por categoria e status
- 🔄 Duplicar serviços
- 🔄 Upload de imagens
- 🔄 Sistema de tags

### Won't Have (v2.0+)
- ❌ Pacotes/combos automáticos
- ❌ Preços dinâmicos por horário
- ❌ Desconto por volume
- ❌ Agendamento recorrente de serviços

---

## 🎫 User Stories Detalhadas

### US-SRV-001: Criar Categoria
**Como** gerente da barbearia  
**Quero** criar categorias de serviço  
**Para** organizar meu catálogo de serviços

**Critérios de Aceite:**
- [ ] Posso criar categoria com nome único
- [ ] Posso adicionar descrição opcional
- [ ] Posso escolher cor para visual
- [ ] Sistema valida nome duplicado
- [ ] Sistema me notifica de sucesso/erro
- [ ] Categoria aparece imediatamente na listagem

**Cenários de Teste:**
```gherkin
Cenário: Criar categoria com sucesso
  Dado que estou autenticado como gerente
  Quando acesso "Nova Categoria"
  E preencho nome "Cabelo"
  E escolho cor "#4A90E2"
  E clico "Salvar"
  Então categoria é criada
  E aparece na lista de categorias
  E recebo notificação de sucesso

Cenário: Erro ao duplicar categoria
  Dado que categoria "Barba" já existe
  Quando tento criar categoria "Barba"
  Então recebo erro "Categoria já existe"
  E modal permanece aberto
```

---

### US-SRV-002: Criar Serviço Básico
**Como** gerente da barbearia  
**Quero** cadastrar serviços  
**Para** disponibilizar para agendamento

**Critérios de Aceite:**
- [ ] Posso criar serviço com nome, preço e duração
- [ ] Posso vincular a uma categoria
- [ ] Posso adicionar descrição
- [ ] Sistema valida preço > 0
- [ ] Sistema valida duração >= 5 minutos
- [ ] Sistema valida nome único
- [ ] Serviço criado está ativo por padrão

**Cenários de Teste:**
```gherkin
Cenário: Criar serviço válido
  Dado que estou autenticado como gerente
  E categoria "Cabelo" existe
  Quando acesso "Novo Serviço"
  E preencho:
    | Campo      | Valor           |
    | Nome       | Corte Masculino |
    | Categoria  | Cabelo          |
    | Preço      | 35.00           |
    | Duração    | 30              |
  E clico "Salvar"
  Então serviço é criado
  E está ativo
  E aparece na lista

Cenário: Erro ao criar com preço inválido
  Quando preencho preço "0"
  E clico "Salvar"
  Então recebo erro "Preço deve ser maior que zero"
```

---

### US-SRV-003: Customizar por Profissional
**Como** gerente da barbearia  
**Quero** definir valores diferentes por profissional  
**Para** refletir habilidades e tempo diferentes

**Critérios de Aceite:**
- [ ] Posso marcar quais profissionais executam o serviço
- [ ] Posso customizar preço para um profissional
- [ ] Posso customizar duração para um profissional
- [ ] Posso customizar comissão para um profissional
- [ ] Se não customizado, usa valores padrão
- [ ] Sistema salva customizações corretamente
- [ ] Agendamentos usam valores customizados

**Cenários de Teste:**
```gherkin
Cenário: Customizar serviço para profissional
  Dado que serviço "Barba" existe com:
    | Preço   | 25.00 |
    | Duração | 25    |
  E profissional "Thiago" existe
  Quando edito serviço "Barba"
  E na seção Profissionais:
    | Profissional | Executa | Customizar | Preço | Duração |
    | Thiago       | Sim     | Sim        | 28.00 | 20      |
  E salvo
  Então customização é salva
  E ao buscar serviço para "Thiago"
  Então retorna preço 28.00 e duração 20

Cenário: Usar valores padrão sem customização
  Dado que serviço "Barba" existe
  E profissional "João" NÃO tem customização
  Quando busco serviço para "João"
  Então retorna valores padrão do serviço
```

---

### US-SRV-004: Filtrar Serviços
**Como** gerente da barbearia  
**Quero** filtrar serviços por categoria e status  
**Para** encontrar rapidamente o que preciso

**Critérios de Aceite:**
- [ ] Posso filtrar por categoria
- [ ] Posso filtrar por status (Ativo/Inativo/Todos)
- [ ] Posso buscar por nome
- [ ] Filtros funcionam em combinação
- [ ] Resultados atualizam em tempo real
- [ ] URL reflete os filtros ativos

---

### US-SRV-005: Inativar Serviço
**Como** gerente da barbearia  
**Quero** inativar serviços temporariamente  
**Para** não deletá-los permanentemente

**Critérios de Aceite:**
- [ ] Posso alternar status ativo/inativo
- [ ] Serviço inativo não aparece em agendamentos
- [ ] Serviço inativo mantém agendamentos existentes
- [ ] Posso reativar serviço inativo
- [ ] Histórico de alterações é mantido

---

## 🐛 Bugs Conhecidos / Tech Debt

> A ser preenchido durante desenvolvimento

---

## 📈 Roadmap Futuro (v2.0+)

### Pacotes e Combos
- Criar serviços compostos (ex: Corte + Barba)
- Aplicar desconto em pacotes
- Duração calculada automaticamente

### Preços Dinâmicos
- Preço por horário (pico vs. baixa demanda)
- Preço por dia da semana
- Promoções automáticas

### Agendamento Inteligente
- Sugerir serviços baseado em histórico
- Recomendar profissional ideal
- Prever duração baseada em dados reais

### Analytics
- Serviços mais vendidos
- Rentabilidade por serviço
- Tempo médio real vs. estimado
- Taxa de no-show por serviço

---

## 🔗 Dependências

### Dependências de Entrada (Bloqueantes)
- ✅ Sistema de Autenticação (pronto)
- ✅ Multi-tenant implementado (pronto)
- ✅ Cadastro de Profissionais (pronto)
- ✅ Banco de dados PostgreSQL (pronto)

### Dependências de Saída (Este módulo bloqueia)
- ⏸️ Módulo de Agendamentos (aguardando serviços)
- ⏸️ Módulo Financeiro (cálculo de comissões)
- ⏸️ Relatórios (análise de serviços)

---

**Última atualização:** 26/11/2025  
**Responsável:** Product + Tech Lead
