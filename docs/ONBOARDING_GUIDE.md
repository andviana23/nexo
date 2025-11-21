# 🚀 Guia de Onboarding & Primeiro Acesso

**Objetivo:** garantir que um novo cliente consiga criar sua barbearia, configurar preferências iniciais e acessar o dashboard sem intervenção manual.

---

## 🔑 Fluxo Completo (3 passos)

1) **Signup (`/signup` + `POST /auth/signup`)**  
   - Campos obrigatórios: `barberName`, `cnpj` (14 dígitos), `email`, `password`, `name`.  
   - Validações: CNPJ válido, email único, senha forte (8+ chars, maiúscula, minúscula, número e símbolo).  
   - Resultado: cria `tenant` (ativo, plano `free`), cria usuário `OWNER`, retorna `access_token` + `refresh_token` + `user` com `tenant`.

2) **Configuração inicial (`/onboarding` — Step 2)**  
   - Endpoint: `POST /onboarding/configure` com cabeçalhos `Authorization: Bearer <token>` e `X-Tenant-ID`.  
   - Payload:
     ```json
     {
       "business_hours": { "opening_time": "08:00", "closing_time": "18:00", "days_open": ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday"] },
       "financial_settings": { "default_commission_rate": 30, "accepted_payment_methods": ["PIX", "DINHEIRO", "DEBITO", "CREDITO"] },
       "preferences": { "timezone": "America/Sao_Paulo", "default_service_duration": 30 }
     }
     ```
   - Persistência: grava em `tenant_settings` e atualiza `updated_at` do tenant.

3) **Conclusão (`POST /tenants/onboarding/complete`)**  
   - Marca `tenants.onboarding_completed = true`.  
   - Front-end atualiza cache local, cookie `bap.onboarding_completed=true` e redireciona para `/dashboard`.

---

## 🧪 Testes Rápidos (cURL)

```bash
# 1) Signup
curl -X POST https://api.seudominio.com/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"barber_name":"Barbearia Teste","cnpj":"12345678000190","email":"dono@teste.com","password":"Teste@1234","name":"Dono"}'

# 2) Configuração inicial
curl -X POST https://api.seudominio.com/api/v1/onboarding/configure \
  -H "Authorization: Bearer <access_token>" \
  -H "X-Tenant-ID: <tenant_id>" \
  -H "Content-Type: application/json" \
  -d '{"business_hours":{"opening_time":"08:00","closing_time":"18:00","days_open":["monday","tuesday"]},"financial_settings":{"default_commission_rate":25,"accepted_payment_methods":["PIX","DEBITO"]},"preferences":{"timezone":"America/Sao_Paulo","default_service_duration":30}}'

# 3) Concluir onboarding
curl -X POST https://api.seudominio.com/api/v1/tenants/onboarding/complete \
  -H "Authorization: Bearer <access_token>" \
  -H "X-Tenant-ID: <tenant_id>"
```

---

## 📋 Regras de Negócio e Segurança

- **Autenticação:** JWT RS256 obrigatório; headers `Authorization` e `X-Tenant-ID` exigidos em `/onboarding/*`.  
- **Tokens:** `access_token` + `refresh_token` retornados em signup/login; cookies `bap.access_token`, `bap.refresh_token`, `bap.onboarding_completed`, `bap.tenant_id` sincronizados no front.  
- **Redirecionamentos:** middleware força `/onboarding` quando `bap.onboarding_completed=false` mesmo se o usuário tentar acessar rotas privadas.  
- **Idempotência:** `/tenants/onboarding/complete` é seguro para múltiplas chamadas; configurações podem ser atualizadas via POST repetidos.

---

## 🛠️ Troubleshooting

- **422/409 no signup:** verifique CNPJ (14 dígitos válidos) e email único.  
- **401/403 na configuração:** confirme cookies/tokens ainda válidos e se `X-Tenant-ID` está presente.  
- **Ficou preso no dashboard sem concluir onboarding:** limpe cookies `bap.*` e faça login novamente; middleware redireciona para `/onboarding`.  
- **Backend sem JWT carregado:** signup retorna `SERVICE_UNAVAILABLE` até as chaves RSA estarem presentes (`keys/private_key.pem` e `keys/public_key.pem` ou variáveis `JWT_*_PATH`).

---

## ✅ Checklist Rápido de QA

- [ ] Signup cria tenant e usuário OWNER e retorna tokens + tenant.  
- [ ] `/auth/me` retorna `tenant` com `onboarding_completed` correto.  
- [ ] Wizard `/onboarding` salva configurações e atualiza cookie `bap.onboarding_completed`.  
- [ ] Após concluir onboarding, usuário é redirecionado ao dashboard e não volta ao wizard.
