# ✅ Checklist de Correções — Módulo de Agendamento (Dez/2025)

Escopo: sanear todos os bugs e lacunas mapeadas na análise do módulo de Agendamento. Priorize conforme ordem abaixo (bloqueadores primeiro).

## 1) Fluxo de status e timestamps
- [ ] Persistir `checked_in_at`, `started_at`, `finished_at` ao mudar status (check-in/start/finish) no use case e no repository.
- [ ] Garantir que o update não zere esses campos quando já existentes.
- [ ] Atualizar mapper/response se novos campos forem expostos.
- [ ] Cobrir com testes de integração (check-in → start → finish).

## 2) Escopo de barbeiro (RBAC)
- [ ] Em todas as rotas de ação (`confirm`, `check-in`, `start`, `finish`, `complete`, `cancel`, `no-show`, `reschedule`, `update-status`), validar se `BARBER` só altera agendamentos do próprio `professional_id`.
- [ ] Adicionar testes de integração que garantam 403 quando barbeiro tenta agir em agendamento alheio.

## 3) Validação de profissional no reagendamento
- [ ] Ao trocar `professional_id`, validar existência e status ATIVO (reader).
- [ ] Retornar erro de negócio legível (ex.: `profissional não encontrado` ou `inativo`).
- [ ] Testes cobrindo troca de profissional válido e inválido.

## 4) Códigos HTTP coerentes
- [ ] Devolver 404 para `GET/ PATCH /status /reschedule /cancel` quando agendamento não existe (conforme swagger).
- [ ] Manter 409 para conflitos de horário; 400 apenas para validação.
- [ ] Ajustar testes de integração para os novos códigos.

## 5) RBAC de criação
- [ ] Decidir regra: permitir `BARBER` criar ou não. Alinhar contrato e código.
- [ ] Se restrito, trocar middleware do POST `/appointments` para `RequireAdminAccess` e ajustar frontend se necessário.
- [ ] Atualizar documentação e swagger.

## 6) Disponibilidade (rota faltante)
- [ ] Implementar `GET /appointments/availability` (filtros: profissional, data, range).
- [ ] Consultar conflitos + bloqueios + intervalo mínimo.
- [ ] Adicionar testes unitários/integração e documentar no swagger.

## 7) Atualização/remoção REST
- [ ] Implementar `PUT /appointments/:id` para editar serviços/notas (validação completa de serviços ativos, recalcular preço/duração, conflitos).
- [ ] Implementar `DELETE /appointments/:id` como alias REST para cancelamento (status `CANCELED`), mantendo motivo opcional.
- [ ] Ajustar frontend service se necessário.

## 8) Estatísticas de agendamento
- [ ] Expor endpoint para `AppointmentDailyStats` (ex.: `GET /appointments/stats/daily?date=YYYY-MM-DD`).
- [ ] Validar RBAC (OWNER/MANAGER/RECEPTIONIST).
- [ ] Testes cobrindo contagens e receita.

## 9) Mapper de calendário
- [ ] Preencher `customerName`/`professionalName` em `AppointmentToCalendarEvent` para tooltips completos.
- [ ] Validar integração com FullCalendar (frontend).

## 10) Documentação e contrato
- [ ] Atualizar `docs/Agendamento/API_AGENDAMENTO.md` e `docs/04-backend/API_REFERENCE.md` com rotas reais de workflow e novas rotas (availability, PUT, DELETE, stats).
- [ ] Remover rotas inexistentes ou marcá-las como futuras se não forem implementadas.
- [ ] Sincronizar RBAC descrito vs. aplicado.
- [ ] Regenerar swagger (`backend/docs/swagger.*`) após mudanças.

## 11) Testes E2E/Frontend
- [ ] Ajustar `frontend/tests/e2e/appointments*.spec.ts` para novas rotas/códigos.
- [ ] Adicionar casos para availability e troca de barbeiro.
- [ ] Smoke script (`scripts/smoke_tests_complete.sh`) incluir availability e novo endpoint de stats.

## 12) Migrações e compatibilidade
- [ ] Verificar se novas colunas/índices são necessários para availability ou stats; criar migration se aplicável.
- [ ] Garantir backward compatibility (campos opcionais nos DTOs).

## 13) Observabilidade e logs
- [ ] Incluir contexto (tenant_id, appointment_id, professional_id) nos logs das novas rotas.
- [ ] Adicionar métricas/contadores para transições de status e no-shows.

## 14) Grooming rápido
- [ ] Revisar mensagens de erro retornadas ao frontend (PT-BR, consistentes).
- [ ] Confirmar que `total_price` e preços de serviços seguem formato numérico `"XX.XX"`.

> Conclua os itens em ordem; marque cada checkbox ao finalizar. Prioridade máxima: itens 1–5.
