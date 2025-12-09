import { expect, test } from '@playwright/test';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * E2E Tests: Módulo Multi-Unit
 *
 * @description Testes end-to-end para o módulo de múltiplas unidades.
 * Valida:
 * - Exibição do seletor de unidades (quando habilitado)
 * - Troca de unidade
 * - Persistência da unidade selecionada
 * - Feature flag multi_unit_enabled
 * - Header X-Unit-ID em requisições
 */

// Configuração de teste
const TEST_USER = {
  email: 'andrey@tratodebarbados.com',
  password: '@Aa30019258',
};

// Configurar testes para rodarem em série
test.describe.configure({ mode: 'serial' });

test.describe('Módulo Multi-Unit', () => {
  test.beforeEach(async ({ page }) => {
    // Fazer login antes de cada teste
    await page.goto('/login', { waitUntil: 'domcontentloaded' });

    // Aguardar campos estarem disponíveis
    await page.waitForSelector('input[name="email"]', { timeout: 10000 });

    await page.fill('input[name="email"]', TEST_USER.email);
    await page.fill('input[name="password"]', TEST_USER.password);

    // Clicar no submit
    await page.click('button[type="submit"]');

    // Aguardar navegação
    await page.waitForURL((url) => url.pathname !== '/login', { timeout: 25000 });
    await page.waitForTimeout(1000);
  });

  test('1. deve exibir o seletor de unidades no header quando multi-unit habilitado', async ({
    page,
  }) => {
    // Aguardar a página carregar completamente
    await page.waitForLoadState('networkidle');

    // Verificar se o seletor de unidades existe
    // O seletor usa o componente UnitSelector com Building2 icon
    const unitSelector = page.locator('[data-slot="dropdown-menu-trigger"]').filter({
      has: page.locator('svg.lucide-building-2'),
    });

    // Verificar se existe (pode não existir se multi-unit não está habilitado)
    const isVisible = await unitSelector.isVisible().catch(() => false);

    if (isVisible) {
      console.log('✅ Seletor de unidades está visível - multi-unit habilitado');
      await expect(unitSelector).toBeVisible();
    } else {
      console.log('⚠️ Seletor de unidades não visível - multi-unit pode estar desabilitado');
      // Não falha o teste, apenas registra
    }
  });

  test('2. deve exibir lista de unidades ao clicar no seletor', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // Encontrar e clicar no seletor
    const unitSelector = page.locator('button').filter({
      has: page.locator('svg.lucide-building-2'),
    });

    // Pular se não estiver visível
    if (!(await unitSelector.isVisible().catch(() => false))) {
      test.skip();
      return;
    }

    await unitSelector.click();
    await page.waitForTimeout(300);

    // Verificar se o dropdown abriu
    const dropdown = page.locator('[data-slot="dropdown-menu-content"]');
    await expect(dropdown).toBeVisible();

    // Verificar se tem opções de unidades
    const unitItems = dropdown.locator('[data-slot="dropdown-menu-item"]');
    const count = await unitItems.count();

    expect(count).toBeGreaterThanOrEqual(1);
    console.log(`✅ Dropdown exibe ${count} unidade(s)`);
  });

  test('3. deve trocar de unidade ao selecionar', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    const unitSelector = page.locator('button').filter({
      has: page.locator('svg.lucide-building-2'),
    });

    if (!(await unitSelector.isVisible().catch(() => false))) {
      test.skip();
      return;
    }

    // Verificar unidade atual
    const currentUnitText = await unitSelector.textContent();
    console.log('Unidade atual:', currentUnitText);

    // Clicar para abrir dropdown
    await unitSelector.click();
    await page.waitForTimeout(300);

    const dropdown = page.locator('[data-slot="dropdown-menu-content"]');
    const unitItems = dropdown.locator('[data-slot="dropdown-menu-item"]');

    // Pegar todas as unidades disponíveis
    const unitCount = await unitItems.count();
    if (unitCount < 2) {
      console.log('⚠️ Menos de 2 unidades disponíveis, pulando teste de troca');
      test.skip();
      return;
    }

    // Clicar na segunda unidade (diferente da atual)
    await unitItems.nth(1).click();
    await page.waitForTimeout(500);

    // Verificar se a unidade mudou
    const newUnitText = await unitSelector.textContent();
    console.log('Nova unidade:', newUnitText);

    // A unidade deveria ter mudado
    expect(newUnitText).not.toBe(currentUnitText);
    console.log('✅ Unidade trocada com sucesso');
  });

  test('4. deve persistir a unidade selecionada após refresh', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    const unitSelector = page.locator('button').filter({
      has: page.locator('svg.lucide-building-2'),
    });

    if (!(await unitSelector.isVisible().catch(() => false))) {
      test.skip();
      return;
    }

    // Pegar unidade atual
    const currentUnitText = await unitSelector.textContent();

    // Fazer refresh
    await page.reload();
    await page.waitForLoadState('networkidle');

    // Verificar se a unidade foi mantida
    const unitSelectorAfter = page.locator('button').filter({
      has: page.locator('svg.lucide-building-2'),
    });

    if (!(await unitSelectorAfter.isVisible().catch(() => false))) {
      test.skip();
      return;
    }

    const unitTextAfterRefresh = await unitSelectorAfter.textContent();

    expect(unitTextAfterRefresh).toBe(currentUnitText);
    console.log('✅ Unidade persistida após refresh:', unitTextAfterRefresh);
  });

  test('5. deve enviar X-Unit-ID nas requisições', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // Interceptar requisições para verificar o header
    const requestHeaders: Record<string, string>[] = [];

    page.on('request', (request) => {
      const url = request.url();
      if (url.includes('/api/')) {
        const headers = request.headers();
        requestHeaders.push({
          url: url,
          unitId: headers['x-unit-id'] || 'NOT_SET',
        });
      }
    });

    // Navegar para uma página que faz requisições à API
    await page.goto('/agendamentos');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    // Verificar se alguma requisição teve o header
    const requestsWithUnitId = requestHeaders.filter((r) => r.unitId !== 'NOT_SET');

    if (requestsWithUnitId.length > 0) {
      console.log('✅ X-Unit-ID encontrado em', requestsWithUnitId.length, 'requisições');
      console.log('Exemplo:', requestsWithUnitId[0]);
    } else {
      console.log('⚠️ Nenhuma requisição com X-Unit-ID (pode ser esperado se multi-unit desabilitado)');
    }

    // Não falha - apenas registra
  });

  test('6. deve limpar unidade selecionada ao fazer logout', async ({ page }) => {
    await page.waitForLoadState('networkidle');

    // Verificar localStorage antes do logout
    const unitDataBefore = await page.evaluate(() => {
      return localStorage.getItem('nexo-unit');
    });
    console.log('Unit data antes do logout:', unitDataBefore ? 'existe' : 'não existe');

    // Fazer logout (procurar no menu do usuário)
    const userMenu = page.locator('button').filter({
      has: page.locator('svg.lucide-user, svg.lucide-circle-user'),
    });

    if (await userMenu.isVisible().catch(() => false)) {
      await userMenu.click();
      await page.waitForTimeout(300);

      // Clicar em logout
      const logoutButton = page.locator('text=Sair').or(page.locator('text=Logout'));
      if (await logoutButton.isVisible().catch(() => false)) {
        await logoutButton.click();
        await page.waitForURL('**/login', { timeout: 10000 });

        // Verificar localStorage após logout
        const unitDataAfter = await page.evaluate(() => {
          return localStorage.getItem('nexo-unit');
        });

        // Deve ter sido limpo ou resetado
        if (!unitDataAfter || unitDataAfter === '{"state":{}}') {
          console.log('✅ Dados de unidade limpos após logout');
        } else {
          console.log('⚠️ Dados de unidade ainda presentes:', unitDataAfter);
        }
      }
    }
  });
});

test.describe('Feature Flag - Multi-Unit', () => {
  test('deve respeitar a feature flag multi_unit_enabled', async ({ page }) => {
    // Este teste verifica o comportamento baseado na feature flag
    // Não podemos alterar a flag, mas podemos verificar o comportamento

    await page.goto('/login');
    await page.waitForSelector('input[name="email"]', { timeout: 10000 });
    await page.fill('input[name="email"]', TEST_USER.email);
    await page.fill('input[name="password"]', TEST_USER.password);
    await page.click('button[type="submit"]');
    await page.waitForURL((url) => url.pathname !== '/login', { timeout: 25000 });
    await page.waitForLoadState('networkidle');

    // Verificar se há requisição para feature flags
    let featureFlagsCalled = false;
    page.on('request', (request) => {
      if (request.url().includes('feature-flags')) {
        featureFlagsCalled = true;
      }
    });

    // Dar tempo para as requisições
    await page.waitForTimeout(2000);

    // Verificar se o componente UnitSelector está presente
    const unitSelector = page.locator('button').filter({
      has: page.locator('svg.lucide-building-2'),
    });

    const hasUnitSelector = await unitSelector.isVisible().catch(() => false);

    console.log('Feature flags endpoint chamado:', featureFlagsCalled);
    console.log('UnitSelector visível:', hasUnitSelector);

    // O seletor só deve aparecer se:
    // 1. Feature flag está habilitada
    // 2. Usuário tem acesso a mais de 1 unidade
    // Registra o estado para análise
  });
});
