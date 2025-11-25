# 🎮 Release v1.1.0 — Fidelidade + Gamificação

**Versão:** 1.1.0
**Nome:** Fidelidade e Engajamento
**Status:** ⏳ Planejado
**Data Prevista:** Março 2026
**Dependência:** ✅ v1.0.0 concluído
**Objetivo:** Aumentar retenção de clientes e engajamento de barbeiros

---

## 🎯 Visão Geral

A versão **v1.1.0** adiciona mecanismos de **fidelização de clientes** e **gamificação para barbeiros**, criando loops de engajamento que aumentam lifetime value (LTV) e reduzem churn.

**Principais Módulos:**

- ✅ **Cashback** - Recompensa clientes fiéis
- ✅ **Gamificação** - Engaja e motiva barbeiros
- ✅ **Metas Avançadas** - Tracking automático e alertas

---

## 📋 Funcionalidades Principais

### 1. Programa de Cashback

**Problema que resolve:**
Cliente não volta, concorrência oferece benefícios.

**Solução:**

- Cashback configurável por unidade (ex: 5% do valor gasto)
- Acúmulo automático a cada compra
- Uso parcial/total em próximas compras
- Expiração configurável
- Saldo visível no app do cliente

**Regras de Negócio:**

- Cashback não pode gerar saldo negativo
- Expiração configurada em parâmetros da unidade
- Pode ser usado em serviços e produtos
- Desconto parcial permitido

**Critérios de Aceite:**

- [ ] Cliente acumula cashback automaticamente
- [ ] Cliente vê saldo no app
- [ ] Cashback pode ser usado no checkout
- [ ] Expiração funciona conforme configurado
- [ ] Histórico de movimentações completo

**Implementação Técnica:**
⚪ Planejado para Sprint 14-15 (Fev 2026)

---

### 2. Gamificação de Barbeiros

**Problema que resolve:**
Barbeiro desmotivado, alta rotatividade, falta de evolução clara.

**Solução:**

- Sistema de XP (experiência)
- Níveis: Bronze → Prata → Ouro → Diamante
- XP baseado em:
  - Atendimentos realizados
  - Ticket médio
  - Retenção de clientes
  - Pontualidade (futuro)
- Bônus ao subir de nível
- Plano de carreira com aumento de comissão

**Regras de Negócio:**

- XP calculado automaticamente
- Subida de nível pode dar bônus/comissão maior
- Ranking visível para equipe
- Histórico de evolução preservado

**Critérios de Aceite:**

- [ ] Barbeiro vê XP e nível atual
- [ ] Progressão é calculada corretamente
- [ ] Bônus são aplicados ao atingir nível
- [ ] Ranking atualiza em tempo real
- [ ] Histórico não é perdido

**Implementação Técnica:**
⚪ Planejado para Sprint 15-16 (Fev 2026)

---

### 3. Metas Avançadas

**Problema que resolve:**
Metas simples não motivam, falta feedback em tempo real.

**Solução:**

- Metas automáticas (baseadas em histórico)
- Metas por barbeiro
- Metas de ticket médio
- Alertas de desvio (abaixo de X% da meta)
- Dashboard de progresso em tempo real

**Regras de Negócio:**

- Metas podem ser manuais ou automáticas
- Automáticas usam média dos últimos 3 meses
- Alertas disparam ao desviar >20%
- Progresso atualiza diariamente

**Critérios de Aceite:**

- [ ] Metas automáticas calculam corretamente
- [ ] Alertas funcionam quando há desvio
- [ ] Barbeiro vê progresso da própria meta
- [ ] Manager vê todas as metas da unidade
- [ ] Comparativo meta vs realizado preciso

**Implementação Técnica:**
Ver `/Tarefas/05-METAS/modulo-04-metas-automaticas.md`

---

## 📊 Impacto Esperado

### Métricas de Sucesso

| Métrica                    | Baseline (v1.0) | Meta (v1.1) |
| -------------------------- | --------------- | ----------- |
| **Churn Mensal**           | 15%             | <10%        |
| **LTV**                    | R$ 800          | R$ 1.200    |
| **Frequência Visita**      | 1x/mês          | 1.5x/mês    |
| **Ticket Médio**           | R$ 65           | R$ 75       |
| **Rotatividade Barbeiros** | 30%/ano         | <20%/ano    |
| **NPS Barbeiros**          | 7               | >8          |

---

## 🔗 Implementação Técnica

### Backend

- [ ] Entidade `Cashback` (acúmulo, expiração, uso)
- [ ] Entidade `BarbeiroXP` (pontos, níveis, histórico)
- [ ] Use cases de cashback
- [ ] Use cases de gamificação
- [ ] Cron job de expiração de cashback
- [ ] Cron job de cálculo de XP

### Frontend

- [ ] Tela de configuração de cashback
- [ ] Tela de histórico de cashback (cliente)
- [ ] Dashboard de gamificação (barbeiro)
- [ ] Ranking de barbeiros
- [ ] Alertas de metas

### Mobile

- [ ] App cliente: saldo de cashback
- [ ] App barbeiro: XP e ranking

---

## ✅ Critérios de Conclusão

v1.1.0 estará **PRONTO** quando:

### Funcionalidades

- [ ] Cashback funcionando end-to-end
- [ ] Gamificação ativa para barbeiros
- [ ] Metas avançadas operacionais
- [ ] Apps mobile atualizados

### Qualidade

- [ ] Cobertura de testes >75%
- [ ] Performance mantida (p95 <300ms)
- [ ] UX validado com usuários

### Negócio

- [ ] Churn reduzido >30%
- [ ] LTV aumentado >40%
- [ ] NPS barbeiros >8

---

## 📅 Cronograma

| Milestone           | Data Prevista   | Status       |
| ------------------- | --------------- | ------------ |
| Design UX/UI        | Jan 2026        | ⚪ Planejado |
| Backend Cashback    | Fev 2026        | ⚪ Planejado |
| Backend Gamificação | Fev 2026        | ⚪ Planejado |
| Frontend Web        | Mar 2026        | ⚪ Planejado |
| Mobile Apps         | Mar 2026        | ⚪ Planejado |
| **Release v1.1.0**  | **31 Mar 2026** | ⚪ Planejado |

---

## 📚 Referências

- [PRD - Fidelidade](../../PRD-NEXO.md#48-módulo-de-fidelidade-cashback)
- [PRD - Gamificação](../../PRD-NEXO.md#49-módulo-de-gamificação--plano-de-carreira)
- [PRD - Metas](../../PRD-NEXO.md#410-módulo-de-metas--kpis)

---

**Última Atualização:** 22/11/2025
**Próxima Revisão:** Conclusão de v1.0.0
