# Compliance LGPD — Barber Analytics Pro v2

## Inventário de Dados
- **Usuários:** nome, email, senha (hash), role, ativo, timestamps.
- **Tenants:** nome, CNPJ, plano, status, onboarding.
- **Preferências:** analytics, error tracking e marketing (tabela `user_preferences`).
- **Logs:** audit logs (criação/edição), IP/UA nos access logs, health/metrics.
- **Financeiro:** receitas, despesas, assinaturas (retidos por obrigação fiscal).

## Direitos do Titular
- **Exportação de dados:** `GET /api/v1/me/export` (1 vez a cada 24h) retorna JSON com usuário + tenant. Header `Content-Disposition` enviado.
- **Exclusão/Anonimização:** `DELETE /api/v1/me` exige senha; anonimiza email, limpa hash, desativa conta e registra `deleted_at`.
- **Preferências/Consentimento:** `GET/PUT /api/v1/me/preferences` persiste consentimentos; banner no frontend salva local e sincroniza quando autenticado.

## Consentimento
- Banner/modal com categorias: necessários (sempre ativos), analytics, monitoramento de erros, marketing.
- Persistência em cookie/localStorage e backend (`user_preferences`).
- Serviços opcionais só devem inicializar se preferências permitirem.

## Política de Privacidade
- Publicada em `/privacy` (frontend). Inclui finalidades, bases legais, retenção e canal DPO (dpo@barberanalytics.pro).

## Retenção
- Logs técnicos: 90 dias.
- Dados fiscais: conforme obrigação legal mínima.
- Contas excluídas: desativadas e anonimizadas; podem ser removidas definitivamente após janela de retenção.

## Segurança
- JWT RS256, RBAC, auditoria, backups com teste de restore (vide T-OPS-005).
- Transporte seguro (HTTPS), prevenção de uso de dados sem consentimento (feature flags para opcionais).
