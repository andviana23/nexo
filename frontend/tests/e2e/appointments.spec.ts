import { expect, test } from '@playwright/test';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * E2E Tests: Módulo de Agendamento
 * 
 * @description Testes end-to-end para o módulo de Agendamento
 * Valida o fluxo completo de transições de status:
 * CREATED → CONFIRMED → CHECKED_IN → IN_SERVICE → AWAITING_PAYMENT → DONE
 * 
 * Também valida:
 * - Criação de agendamento
 * - Reagendamento
 * - Cancelamento
 * - Bloqueio de horários
 * - Integração com comandas
 */

// Configuração de teste
const TEST_USER = {
  email: 'andrey@tratodebarbados.com',
  password: '@Aa30019258',
};

// Configurar testes para rodarem em série (evitar concorrência)
test.describe.configure({ mode: 'serial' });

let createdAppointmentId: string | null = null;

test.describe('Módulo de Agendamento - Fluxo Completo', () => {
  test.beforeEach(async ({ page }) => {
    // Fazer login antes de cada teste
    await page.goto('/login', { waitUntil: 'domcontentloaded' });
    
    // Aguardar campos estarem disponíveis
    await page.waitForSelector('input[name="email"]', { timeout: 10000 });
    
    await page.fill('input[name="email"]', TEST_USER.email);
    await page.fill('input[name="password"]', TEST_USER.password);
    
    // Clicar no submit
    await page.click('button[type="submit"]');
    
    // Aguardar navegação - pode ir para / ou /dashboard
    await page.waitForURL((url) => url.pathname !== '/login', { timeout: 25000 });
    
    await page.waitForTimeout(500);
    
    // Navegar para a página de agendamentos
    await page.goto('/agendamentos', { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    // Fechar menu mobile se estiver aberto
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);
  });

  test('1. deve exibir o calendário de agendamentos', async ({ page }) => {
    // Aguardar a página carregar
    await page.waitForLoadState('domcontentloaded');
    await page.waitForTimeout(3000);
    
    // Verificar se estamos na página correta
    expect(page.url()).toContain('/agendamentos');
    
    // Verificar se a página tem conteúdo (não é página em branco ou erro)
    const bodyText = await page.locator('body').textContent();
    expect(bodyText).toBeTruthy();
    expect(bodyText!.length).toBeGreaterThan(100);
    
    console.log('✅ Página de agendamentos carregada em:', page.url());
  });

  test('2. deve criar um novo agendamento', async ({ page }) => {
    // Clicar no botão de novo agendamento
    const newButton = page.locator('button').filter({ hasText: /novo agendamento/i }).first();
    await newButton.click();
    
    // Aguardar modal abrir
    await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
    
    // Preencher formulário
    // Selecionar cliente (primeiro da lista)
    const customerSelect = page.locator('select[name="customer_id"]');
    await customerSelect.selectOption({ index: 1 });
    
    await page.waitForTimeout(300);
    
    // Selecionar profissional (primeiro da lista)
    const professionalSelect = page.locator('select[name="professional_id"]');
    await professionalSelect.selectOption({ index: 1 });
    
    await page.waitForTimeout(300);
    
    // Selecionar serviço (primeiro da lista)
    const serviceCheckbox = page.locator('input[type="checkbox"][name="services"]').first();
    await serviceCheckbox.check();
    
    await page.waitForTimeout(300);
    
    // Preencher data (amanhã)
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const dateString = tomorrow.toISOString().split('T')[0];
    
    const dateInput = page.locator('input[type="date"]');
    await dateInput.fill(dateString);
    
    await page.waitForTimeout(300);
    
    // Preencher horário
    const timeInput = page.locator('input[type="time"]');
    await timeInput.fill('14:00');
    
    await page.waitForTimeout(500);
    
    // Submeter formulário
    const submitButton = page.locator('button[type="submit"]').filter({ hasText: /criar/i });
    await submitButton.click();
    
    // Aguardar toast de sucesso
    await page.waitForSelector('text=/agendamento criado/i', { timeout: 10000 });
    
    // Aguardar modal fechar
    await page.waitForTimeout(1000);
    
    // Verificar se o evento apareceu no calendário
    const event = page.locator('.fc-event').first();
    await expect(event).toBeVisible({ timeout: 5000 });
    
    // Capturar ID do agendamento criado (da URL após clicar no evento)
    await event.click();
    await page.waitForURL(/\/agendamentos\/[a-f0-9-]+/i, { timeout: 5000 });
    
    const url = page.url();
    const match = url.match(/\/agendamentos\/([a-f0-9-]+)/i);
    if (match) {
      createdAppointmentId = match[1];
      console.log('✅ Agendamento criado com ID:', createdAppointmentId);
    }
  });

  test('3. deve confirmar o agendamento (CREATED → CONFIRMED)', async ({ page }) => {
    test.skip(!createdAppointmentId, 'Agendamento não foi criado no teste anterior');
    
    // Navegar para a página de detalhes
    await page.goto(`/agendamentos/${createdAppointmentId}`, { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    // Verificar status atual
    const statusBadge = page.locator('[data-testid="appointment-status"]').or(page.locator('text=/status/i').locator('..'));
    await expect(statusBadge).toBeVisible({ timeout: 10000 });
    
    // Abrir menu de ações
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).or(
      page.locator('button[aria-label*="ações"]')
    ).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    // Clicar em confirmar
    const confirmButton = page.locator('text=/confirmar/i').first();
    await confirmButton.click();
    
    // Aguardar toast de sucesso
    await page.waitForSelector('text=/confirmado/i', { timeout: 10000 });
    
    // Aguardar atualização da página
    await page.waitForTimeout(1000);
    
    // Verificar se status mudou para CONFIRMED
    await expect(statusBadge).toContainText(/confirmado/i);
    
    console.log('✅ Agendamento confirmado');
  });

  test('4. deve fazer check-in (CONFIRMED → CHECKED_IN)', async ({ page }) => {
    test.skip(!createdAppointmentId, 'Agendamento não foi criado');
    
    await page.goto(`/agendamentos/${createdAppointmentId}`, { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    // Abrir menu de ações
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    // Clicar em check-in
    const checkinButton = page.locator('text=/check.*in/i').first();
    await checkinButton.click();
    
    // Aguardar toast
    await page.waitForSelector('text=/check.*in/i', { timeout: 10000 });
    
    await page.waitForTimeout(1000);
    
    // Verificar status
    const statusBadge = page.locator('[data-testid="appointment-status"]').or(page.locator('text=/status/i').locator('..'));
    await expect(statusBadge).toContainText(/check.*in/i);
    
    console.log('✅ Check-in realizado');
  });

  test('5. deve iniciar atendimento (CHECKED_IN → IN_SERVICE)', async ({ page }) => {
    test.skip(!createdAppointmentId, 'Agendamento não foi criado');
    
    await page.goto(`/agendamentos/${createdAppointmentId}`, { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    // Abrir menu de ações
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    // Clicar em iniciar atendimento
    const startButton = page.locator('text=/iniciar/i').first();
    await startButton.click();
    
    // Aguardar toast
    await page.waitForSelector('text=/iniciado/i', { timeout: 10000 });
    
    await page.waitForTimeout(1000);
    
    // Verificar status
    const statusBadge = page.locator('[data-testid="appointment-status"]').or(page.locator('text=/status/i').locator('..'));
    await expect(statusBadge).toContainText(/atendimento/i);
    
    console.log('✅ Atendimento iniciado');
  });

  test('6. deve finalizar atendimento (IN_SERVICE → AWAITING_PAYMENT)', async ({ page }) => {
    test.skip(!createdAppointmentId, 'Agendamento não foi criado');
    
    await page.goto(`/agendamentos/${createdAppointmentId}`, { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    // Abrir menu de ações
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    // Clicar em finalizar
    const finishButton = page.locator('text=/finalizar/i').first();
    await finishButton.click();
    
    // Aguardar toast
    await page.waitForSelector('text=/finalizado/i', { timeout: 10000 });
    
    await page.waitForTimeout(1000);
    
    // Verificar status
    const statusBadge = page.locator('[data-testid="appointment-status"]').or(page.locator('text=/status/i').locator('..'));
    await expect(statusBadge).toContainText(/aguardando.*pagamento/i);
    
    console.log('✅ Atendimento finalizado - aguardando pagamento');
  });

  test('7. deve concluir com pagamento (AWAITING_PAYMENT → DONE)', async ({ page }) => {
    test.skip(!createdAppointmentId, 'Agendamento não foi criado');
    
    await page.goto(`/agendamentos/${createdAppointmentId}`, { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    // Abrir menu de ações
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    // Clicar em concluir/finalizar pagamento
    const completeButton = page.locator('text=/concluir|finalizar.*pagamento/i').first();
    await completeButton.click();
    
    // Aguardar toast
    await page.waitForSelector('text=/concluído|pago/i', { timeout: 10000 });
    
    await page.waitForTimeout(1000);
    
    // Verificar status
    const statusBadge = page.locator('[data-testid="appointment-status"]').or(page.locator('text=/status/i').locator('..'));
    await expect(statusBadge).toContainText(/concluído|finalizado/i);
    
    console.log('✅ Agendamento concluído com sucesso');
  });
});

test.describe('Módulo de Agendamento - Bloqueio de Horários', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/login', { waitUntil: 'domcontentloaded' });
    await page.waitForSelector('input[name="email"]', { timeout: 10000 });
    
    await page.fill('input[name="email"]', TEST_USER.email);
    await page.fill('input[name="password"]', TEST_USER.password);
    await page.click('button[type="submit"]');
    
    await page.waitForURL((url) => url.pathname !== '/login', { timeout: 25000 });
    await page.waitForTimeout(500);
    
    // Navegar para agendamentos
    await page.goto('/agendamentos', { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);
  });

  test('8. deve criar bloqueio de horário', async ({ page }) => {
    // Procurar botão de bloquear horário
    const blockButton = page.locator('button').filter({ hasText: /bloquear.*horário/i }).first();
    
    // Se não encontrar, procurar no menu de ações
    const isVisible = await blockButton.isVisible().catch(() => false);
    
    if (!isVisible) {
      // Tentar abrir menu dropdown
      const moreButton = page.locator('button[aria-label*="mais"]').or(
        page.locator('button').filter({ hasText: /\.\.\./i })
      ).first();
      
      if (await moreButton.isVisible().catch(() => false)) {
        await moreButton.click();
        await page.waitForTimeout(500);
      }
    }
    
    await blockButton.click();
    
    // Aguardar modal abrir
    await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
    
    // Preencher formulário
    // Selecionar profissional
    const professionalSelect = page.locator('select[name="professional_id"]');
    await professionalSelect.selectOption({ index: 1 });
    
    await page.waitForTimeout(300);
    
    // Preencher data (amanhã)
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 2);
    const dateString = tomorrow.toISOString().split('T')[0];
    
    const dateInput = page.locator('input[type="date"]');
    await dateInput.fill(dateString);
    
    await page.waitForTimeout(300);
    
    // Preencher horários
    const startTimeInput = page.locator('input[name="start_time"]');
    await startTimeInput.fill('12:00');
    
    const endTimeInput = page.locator('input[name="end_time"]');
    await endTimeInput.fill('13:00');
    
    await page.waitForTimeout(300);
    
    // Preencher motivo
    const reasonInput = page.locator('input[name="reason"]').or(page.locator('textarea[name="reason"]'));
    await reasonInput.fill('Almoço - Teste E2E');
    
    await page.waitForTimeout(500);
    
    // Submeter
    const submitButton = page.locator('button[type="submit"]').filter({ hasText: /bloquear/i });
    await submitButton.click();
    
    // Aguardar toast de sucesso
    await page.waitForSelector('text=/bloqueado|bloqueio criado/i', { timeout: 10000 });
    
    await page.waitForTimeout(1000);
    
    console.log('✅ Bloqueio de horário criado');
  });
});

test.describe('Módulo de Agendamento - Outras Operações', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/login', { waitUntil: 'domcontentloaded' });
    await page.waitForSelector('input[name="email"]', { timeout: 10000 });
    
    await page.fill('input[name="email"]', TEST_USER.email);
    await page.fill('input[name="password"]', TEST_USER.password);
    await page.click('button[type="submit"]');
    
    await page.waitForURL((url) => url.pathname !== '/login', { timeout: 25000 });
    await page.waitForTimeout(500);
    
    await page.goto('/agendamentos', { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);
    
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);
  });

  test('9. deve reagendar um agendamento', async ({ page }) => {
    // Criar agendamento temporário para reagendar
    const newButton = page.locator('button').filter({ hasText: /novo agendamento/i }).first();
    await newButton.click();
    
    await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
    
    // Preencher formulário básico
    const customerSelect = page.locator('select[name="customer_id"]');
    await customerSelect.selectOption({ index: 1 });
    
    const professionalSelect = page.locator('select[name="professional_id"]');
    await professionalSelect.selectOption({ index: 1 });
    
    const serviceCheckbox = page.locator('input[type="checkbox"][name="services"]').first();
    await serviceCheckbox.check();
    
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 3);
    const dateString = tomorrow.toISOString().split('T')[0];
    
    const dateInput = page.locator('input[type="date"]');
    await dateInput.fill(dateString);
    
    const timeInput = page.locator('input[type="time"]');
    await timeInput.fill('10:00');
    
    await page.waitForTimeout(500);
    
    const submitButton = page.locator('button[type="submit"]').filter({ hasText: /criar/i });
    await submitButton.click();
    
    await page.waitForSelector('text=/agendamento criado/i', { timeout: 10000 });
    await page.waitForTimeout(1000);
    
    // Clicar no evento criado
    const event = page.locator('.fc-event').first();
    await event.click();
    
    await page.waitForURL(/\/agendamentos\/[a-f0-9-]+/i, { timeout: 5000 });
    
    // Reagendar
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    const rescheduleButton = page.locator('text=/reagendar/i').first();
    await rescheduleButton.click();
    
    // Aguardar modal
    await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
    
    // Mudar horário
    const newTimeInput = page.locator('input[type="time"]');
    await newTimeInput.fill('11:00');
    
    await page.waitForTimeout(500);
    
    const confirmButton = page.locator('button[type="submit"]').filter({ hasText: /reagendar/i });
    await confirmButton.click();
    
    await page.waitForSelector('text=/reagendado/i', { timeout: 10000 });
    
    console.log('✅ Agendamento reagendado');
  });

  test('10. deve cancelar um agendamento', async ({ page }) => {
    // Criar agendamento temporário para cancelar
    const newButton = page.locator('button').filter({ hasText: /novo agendamento/i }).first();
    await newButton.click();
    
    await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
    
    const customerSelect = page.locator('select[name="customer_id"]');
    await customerSelect.selectOption({ index: 1 });
    
    const professionalSelect = page.locator('select[name="professional_id"]');
    await professionalSelect.selectOption({ index: 1 });
    
    const serviceCheckbox = page.locator('input[type="checkbox"][name="services"]').first();
    await serviceCheckbox.check();
    
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 4);
    const dateString = tomorrow.toISOString().split('T')[0];
    
    const dateInput = page.locator('input[type="date"]');
    await dateInput.fill(dateString);
    
    const timeInput = page.locator('input[type="time"]');
    await timeInput.fill('16:00');
    
    await page.waitForTimeout(500);
    
    const submitButton = page.locator('button[type="submit"]').filter({ hasText: /criar/i });
    await submitButton.click();
    
    await page.waitForSelector('text=/agendamento criado/i', { timeout: 10000 });
    await page.waitForTimeout(1000);
    
    // Clicar no evento
    const event = page.locator('.fc-event').first();
    await event.click();
    
    await page.waitForURL(/\/agendamentos\/[a-f0-9-]+/i, { timeout: 5000 });
    
    // Cancelar
    const actionsButton = page.locator('button').filter({ hasText: /ações/i }).first();
    await actionsButton.click();
    
    await page.waitForTimeout(500);
    
    const cancelButton = page.locator('text=/cancelar/i').first();
    await cancelButton.click();
    
    // Aguardar modal de confirmação
    await page.waitForSelector('[role="dialog"]', { timeout: 5000 });
    
    // Preencher motivo
    const reasonInput = page.locator('textarea[name="reason"]').or(page.locator('input[name="reason"]'));
    await reasonInput.fill('Cliente desistiu - Teste E2E');
    
    await page.waitForTimeout(500);
    
    const confirmCancelButton = page.locator('button[type="submit"]').filter({ hasText: /cancelar/i });
    await confirmCancelButton.click();
    
    await page.waitForSelector('text=/cancelado/i', { timeout: 10000 });
    
    console.log('✅ Agendamento cancelado');
  });
});
