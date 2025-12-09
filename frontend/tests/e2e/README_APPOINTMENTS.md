# Testes E2E - Módulo de Agendamento

## 📋 Visão Geral

Suite completa de testes end-to-end para o módulo de Agendamento do NEXO, cobrindo todos os fluxos principais e transições de status.

## 🧪 Cobertura de Testes

### Fluxo Completo de Status (Testes 1-7)

```
CREATED → CONFIRMED → CHECKED_IN → IN_SERVICE → AWAITING_PAYMENT → DONE
```

1. **Visualização do Calendário**
   - Verifica renderização do FullCalendar
   - Valida presença de botões de navegação
   - Confirma botão de novo agendamento

2. **Criação de Agendamento**
   - Preenche formulário completo
   - Seleciona cliente, profissional, serviço
   - Define data e horário
   - Valida aparição no calendário

3. **Confirmação (CREATED → CONFIRMED)**
   - Abre menu de ações
   - Executa confirmação
   - Valida mudança de status

4. **Check-in (CONFIRMED → CHECKED_IN)**
   - Realiza check-in do cliente
   - Valida timestamp de check-in
   - Confirma status atualizado

5. **Início do Atendimento (CHECKED_IN → IN_SERVICE)**
   - Inicia atendimento
   - Valida timestamp de início
   - Confirma transição de status

6. **Finalização (IN_SERVICE → AWAITING_PAYMENT)**
   - Finaliza atendimento
   - Valida timestamp de finalização
   - Confirma status aguardando pagamento

7. **Conclusão (AWAITING_PAYMENT → DONE)**
   - Conclui com pagamento
   - Valida status final
   - Confirma conclusão do fluxo

### Bloqueio de Horários (Teste 8)

8. **Criação de Bloqueio**
   - Abre modal de bloqueio
   - Preenche profissional, data, horários
   - Define motivo do bloqueio
   - Valida criação com sucesso

### Outras Operações (Testes 9-10)

9. **Reagendamento**
   - Cria agendamento temporário
   - Abre modal de reagendamento
   - Altera data/horário
   - Valida reagendamento

10. **Cancelamento**
    - Cria agendamento temporário
    - Abre modal de cancelamento
    - Informa motivo do cancelamento
    - Valida cancelamento

## 🚀 Como Executar

### Pré-requisitos

1. **Backend rodando**
   ```bash
   cd backend
   make dev
   # ou
   ./start-dev.sh
   ```

2. **Frontend rodando**
   ```bash
   cd frontend
   pnpm dev
   ```

3. **Banco de dados com dados de teste**
   - Tenant configurado
   - Usuário de teste: `andrey@tratodebarbados.com`
   - Pelo menos 1 cliente cadastrado
   - Pelo menos 1 profissional cadastrado
   - Pelo menos 1 serviço cadastrado

### Executar Testes

#### Opção 1: Script Automatizado (Recomendado)

```bash
cd frontend
./run-e2e-appointments.sh
```

#### Opção 2: Comando Direto

```bash
cd frontend

# Instalar navegadores (primeira vez)
npx playwright install --with-deps chromium firefox

# Executar testes
npx playwright test tests/e2e/appointments.spec.ts
```

#### Opção 3: Modo UI Interativo

```bash
cd frontend
npx playwright test tests/e2e/appointments.spec.ts --ui
```

#### Opção 4: Modo Debug

```bash
cd frontend
npx playwright test tests/e2e/appointments.spec.ts --debug
```

### Executar Teste Específico

```bash
# Por número do teste
npx playwright test tests/e2e/appointments.spec.ts:100

# Por nome (grep)
npx playwright test tests/e2e/appointments.spec.ts -g "deve criar um novo agendamento"
```

## 📊 Relatórios

### Ver Relatório HTML

Após executar os testes:

```bash
npx playwright show-report
```

### Screenshots e Vídeos

- Screenshots são capturados automaticamente em falhas
- Localizados em: `test-results/`
- Trace files para debug: `test-results/**/*.zip`

### Ver Trace de Teste Falhado

```bash
npx playwright show-trace test-results/<pasta-do-teste>/trace.zip
```

## 🔧 Configuração

### playwright.config.ts

```typescript
export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
  ],
});
```

### Modo Serial

Os testes de agendamento rodam em **modo serial** (`test.describe.configure({ mode: 'serial' })`) para:

- Evitar conflitos de dados
- Manter ordem do fluxo (CREATED → DONE)
- Compartilhar ID do agendamento criado entre testes

## 🐛 Troubleshooting

### Testes Falhando por Timeout

**Problema:** Testes falham com `timeout exceeded`

**Soluções:**
```bash
# Aumentar timeout global
npx playwright test --timeout=60000

# Ou editar o teste específico
test('nome do teste', async ({ page }) => {
  test.setTimeout(60000);
  // ...
});
```

### Elementos Não Encontrados

**Problema:** `Error: locator.click: Timeout 30000ms exceeded`

**Soluções:**
1. Verificar se o frontend está rodando
2. Confirmar que há dados no banco (cliente, profissional, serviço)
3. Executar em modo debug: `npx playwright test --debug`

### Login Falhando

**Problema:** Não consegue fazer login

**Soluções:**
1. Verificar credenciais em `TEST_USER`
2. Confirmar que o backend está rodando
3. Verificar logs do backend para erros de autenticação

### Modal Não Abre

**Problema:** `Error: waiting for selector "[role="dialog"]"`

**Soluções:**
1. Aumentar `waitForTimeout` antes de clicar no botão
2. Verificar se não há overlay/modal já aberto
3. Pressionar `Escape` para fechar modais anteriores

## 📝 Boas Práticas

### 1. Esperas Explícitas

```typescript
// ❌ Evitar
await page.waitForTimeout(1000);

// ✅ Preferir
await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
```

### 2. Seletores Robustos

```typescript
// ❌ Evitar (frágil)
await page.click('.btn-primary');

// ✅ Preferir (semântico)
await page.locator('button').filter({ hasText: /criar/i }).click();
```

### 3. Validações Claras

```typescript
// ❌ Evitar
expect(await page.textContent('.status')).toBe('Confirmado');

// ✅ Preferir
const statusBadge = page.locator('[data-testid="appointment-status"]');
await expect(statusBadge).toContainText(/confirmado/i);
```

### 4. Limpeza de Estado

```typescript
// Pressionar Escape para fechar menus/modals
await page.keyboard.press('Escape');
await page.waitForTimeout(500);
```

## 🔄 Integração Contínua

### GitHub Actions

```yaml
name: E2E Tests - Appointments

on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'
      
      - name: Install dependencies
        run: cd frontend && pnpm install
      
      - name: Install Playwright
        run: cd frontend && npx playwright install --with-deps
      
      - name: Start Backend
        run: cd backend && make dev &
      
      - name: Start Frontend
        run: cd frontend && pnpm dev &
      
      - name: Wait for services
        run: sleep 10
      
      - name: Run E2E Tests
        run: cd frontend && npx playwright test tests/e2e/appointments.spec.ts
      
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

## 📚 Recursos

- [Playwright Docs](https://playwright.dev)
- [Best Practices](https://playwright.dev/docs/best-practices)
- [Debugging Tests](https://playwright.dev/docs/debug)
- [CI/CD](https://playwright.dev/docs/ci)

## 📞 Suporte

Em caso de dúvidas ou problemas:

1. Verifique os logs dos testes
2. Execute em modo debug
3. Consulte o trace file
4. Abra issue no repositório
