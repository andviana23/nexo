# 🛠️ 11 — Módulo de Serviços

**Objetivo:** Implementar o cadastro completo de categorias e serviços com suporte a customização por profissional (preço, duração e comissão personalizados).

**Dependências:**
- Profissionais já cadastrados
- Autenticação e multi-tenant implementados

**Status:** 🟡 Planejado
**Sprint alvo:** Sprint 1.4 (Milestone 1.4)
**Pasta:** `Tarefas/11-SERVICOS/`

---

## 📑 Arquivos deste pacote

- `01-contexto.md` — Visão geral e arquitetura
- `02-backlog.md` — User stories e épicos
- `03-sprint-plan.md` — Plano detalhado de implementação
- `04-checklist-dev.md` — Checklist técnico para dev
- `05-checklist-qa.md` — Critérios de aceite e testes

---

## 🎯 Objetivos Principais

1. ✅ Cadastrar categorias de serviço (Cabelo, Barba, Estética, etc.)
2. ✅ Cadastrar serviços com informações básicas
3. ✅ Precificar serviços com valores padrão
4. ✅ Customizar preço, duração e comissão por profissional
5. ✅ Controlar disponibilidade (ativo/inativo)
6. ✅ Garantir isolamento multi-tenant

---

## 🗂️ Estrutura do Módulo

### Backend (Go)
```
backend/
├── internal/
│   ├── domain/
│   │   └── entity/
│   │       ├── categoria.go
│   │       ├── servico.go
│   │       └── servico_profissional.go
│   ├── application/
│   │   ├── usecase/
│   │   │   ├── categoria/
│   │   │   └── servico/
│   │   └── dto/
│   │       ├── categoria_dto.go
│   │       └── servico_dto.go
│   └── infra/
│       ├── repository/
│       │   └── postgres/
│       │       ├── categoria_repository.go
│       │       └── servico_repository.go
│       └── http/
│           └── handler/
│               ├── categoria_handler.go
│               └── servico_handler.go
└── migrations/
    └── 005_create_categorias_servicos.sql
```

### Frontend (Next.js)
```
frontend/src/
├── app/
│   └── (dashboard)/
│       └── cadastros/
│           └── servicos/
│               ├── page.tsx
│               ├── components/
│               │   ├── ServicesList.tsx
│               │   ├── ServiceModal.tsx
│               │   ├── CategoryModal.tsx
│               │   └── ProfessionalCustomization.tsx
│               └── hooks/
│                   ├── useServices.ts
│                   └── useCategories.ts
├── services/
│   ├── service-service.ts
│   └── category-service.ts
└── types/
    ├── service.ts
    └── category.ts
```

---

## 📊 Fluxo de Implementação

### Fase 1: Categorias (Sprint 1.4.1)
- Migration de categorias
- Backend CRUD de categorias
- Frontend: modal e listagem de categorias
- Testes unitários e integração

### Fase 2: Serviços Básicos (Sprint 1.4.2)
- Migration de serviços
- Backend CRUD de serviços
- Frontend: formulário básico de serviço
- Validações e testes

### Fase 3: Customização por Profissional (Sprint 1.4.3)
- Migration de servicos_profissionais
- Lógica de customização no backend
- UI de seleção e customização de profissionais
- Queries otimizadas com COALESCE
- Testes E2E completos

### Fase 4: Recursos Avançados (Sprint 1.4.4)
- Upload de imagens
- Sistema de tags
- Busca avançada
- Duplicar serviços
- Filtros e ordenação

---

## 🔑 Regras de Negócio Críticas

1. **RN-SRV-001:** Nome de categoria único por tenant
2. **RN-SRV-002:** Preço base deve ser > 0
3. **RN-SRV-003:** Duração mínima de 5 minutos
4. **RN-SRV-004:** Comissão entre 0% e 100%
5. **RN-SRV-005:** Nome de serviço único por tenant
6. **RN-SRV-006:** Customização por profissional é opcional
7. **RN-SRV-007:** Não deletar categoria com serviços vinculados
8. **RN-SRV-008:** Não deletar serviço com agendamentos futuros

---

## 📈 Métricas de Sucesso

- [ ] Tempo médio de cadastro de serviço < 2 minutos
- [ ] Taxa de erro na validação < 5%
- [ ] 100% dos serviços com categoria definida
- [ ] Média de 3+ profissionais por serviço
- [ ] 90% dos serviços com valores customizados

---

## 🔗 Referências

- [FLUXO_CADASTRO_SERVIÇO.md](../../docs/11-Fluxos/FLUXO_CADASTRO_SERVIÇO.md)
- [PRD-VALTARIS.md](../../PRD-VALTARIS.md)
- [MODELO_DE_DADOS.md](../../docs/02-arquitetura/MODELO_DE_DADOS.md)

---

**Responsável:** Tech Lead + Product
**Prazo:** Sprint 1.4 (10/12/2025)
