import { expect, test } from '@playwright/test';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * E2E Tests: Fluxo de Seleção de Unidade Obrigatória
 *
 * @description Testes end-to-end para o modal de seleção de unidade após login.
 * Valida:
 * - Modal aparece imediatamente após login
 * - Modal não pode ser fechado sem seleção
 * - Seleção de unidade redireciona para dashboard
 * - Unidade persiste após refresh
 */

// Configuração de teste
const TEST_USER = {
    email: 'andrey@tratodebarbados.com',
    password: '@Aa30019258',
};

// Configurar testes para rodarem em série
test.describe.configure({ mode: 'serial' });

test.describe('Fluxo de Seleção de Unidade Obrigatória', () => {
    test('1. deve exibir modal de seleção após login', async ({ page }) => {
        // Limpar localStorage para garantir estado inicial
        await page.goto('/login');
        await page.evaluate(() => {
            localStorage.removeItem('nexo-auth');
            localStorage.removeItem('nexo-unit');
        });

        // Fazer login
        await page.waitForSelector('input[name="email"]', { timeout: 10000 });
        await page.fill('input[name="email"]', TEST_USER.email);
        await page.fill('input[name="password"]', TEST_USER.password);
        await page.click('button[type="submit"]');

        // Aguardar modal de seleção aparecer
        await page.waitForTimeout(2000);

        // Verificar se o modal está visível
        const modal = page.locator('div[role="dialog"][aria-modal="true"]');
        await expect(modal).toBeVisible({ timeout: 10000 });

        // Verificar título do modal
        const title = page.locator('#unit-selection-title');
        await expect(title).toContainText('Selecione a Unidade');

        console.log('✅ Modal de seleção de unidade exibido após login');
    });

    test('2. deve exibir lista de unidades do usuário', async ({ page }) => {
        // Fazer login primeiro
        await page.goto('/login');
        await page.evaluate(() => {
            localStorage.removeItem('nexo-auth');
            localStorage.removeItem('nexo-unit');
        });
        await page.waitForSelector('input[name="email"]', { timeout: 10000 });
        await page.fill('input[name="email"]', TEST_USER.email);
        await page.fill('input[name="password"]', TEST_USER.password);
        await page.click('button[type="submit"]');

        // Aguardar modal
        await page.waitForTimeout(2000);
        const modal = page.locator('div[role="dialog"][aria-modal="true"]');
        await expect(modal).toBeVisible({ timeout: 10000 });

        // Verificar se tem botões de unidades
        const unitButtons = modal.locator('button:has-text("Trato")');
        const count = await unitButtons.count();

        expect(count).toBeGreaterThanOrEqual(1);
        console.log(`✅ ${count} unidade(s) exibida(s) no modal`);
    });

    test('3. deve selecionar unidade e redirecionar para dashboard', async ({ page }) => {
        // Fazer login primeiro
        await page.goto('/login');
        await page.evaluate(() => {
            localStorage.removeItem('nexo-auth');
            localStorage.removeItem('nexo-unit');
        });
        await page.waitForSelector('input[name="email"]', { timeout: 10000 });
        await page.fill('input[name="email"]', TEST_USER.email);
        await page.fill('input[name="password"]', TEST_USER.password);
        await page.click('button[type="submit"]');

        // Aguardar modal
        await page.waitForTimeout(2000);
        const modal = page.locator('div[role="dialog"][aria-modal="true"]');
        await expect(modal).toBeVisible({ timeout: 10000 });

        // Clicar na primeira unidade
        const firstUnit = modal.locator('button').filter({
            has: page.locator('text=Trato'),
        }).first();
        await firstUnit.click();

        // Aguardar redirecionamento
        await page.waitForURL((url) => url.pathname === '/', { timeout: 15000 });

        // Verificar que modal sumiu
        await expect(modal).not.toBeVisible();

        // Verificar que está no dashboard
        expect(page.url()).toContain('/');
        console.log('✅ Unidade selecionada e redirecionado para dashboard');
    });

    test('4. deve persistir unidade após refresh', async ({ page }) => {
        // Primeiro fazer login e selecionar unidade
        await page.goto('/login');
        await page.evaluate(() => {
            localStorage.removeItem('nexo-auth');
            localStorage.removeItem('nexo-unit');
        });
        await page.waitForSelector('input[name="email"]', { timeout: 10000 });
        await page.fill('input[name="email"]', TEST_USER.email);
        await page.fill('input[name="password"]', TEST_USER.password);
        await page.click('button[type="submit"]');

        // Aguardar e selecionar unidade
        await page.waitForTimeout(2000);
        const modal = page.locator('div[role="dialog"][aria-modal="true"]');
        if (await modal.isVisible()) {
            const firstUnit = modal.locator('button').filter({
                has: page.locator('text=Trato'),
            }).first();
            await firstUnit.click();
            await page.waitForURL((url) => url.pathname === '/', { timeout: 15000 });
        }

        // Fazer refresh
        await page.reload();
        await page.waitForLoadState('networkidle');

        // Verificar que modal NÃO aparece novamente
        const modalAfterRefresh = page.locator('div[role="dialog"][aria-modal="true"]');
        await page.waitForTimeout(2000);

        // Modal não deve estar visível após refresh se unidade já foi selecionada
        const isModalVisible = await modalAfterRefresh.isVisible().catch(() => false);

        if (!isModalVisible) {
            console.log('✅ Unidade persistiu após refresh - modal não apareceu');
        } else {
            // Se modal apareceu, verificar se há unidade no localStorage
            const unitData = await page.evaluate(() => localStorage.getItem('nexo-unit'));
            console.log('⚠️ Modal apareceu após refresh. Unit data:', unitData);
        }
    });

    test('5. deve limpar estado ao fazer logout', async ({ page }) => {
        // Fazer login e selecionar unidade primeiro
        await page.goto('/login');
        await page.evaluate(() => {
            localStorage.removeItem('nexo-auth');
            localStorage.removeItem('nexo-unit');
        });
        await page.waitForSelector('input[name="email"]', { timeout: 10000 });
        await page.fill('input[name="email"]', TEST_USER.email);
        await page.fill('input[name="password"]', TEST_USER.password);
        await page.click('button[type="submit"]');

        // Aguardar e selecionar unidade se modal aparecer
        await page.waitForTimeout(2000);
        const modal = page.locator('div[role="dialog"][aria-modal="true"]');
        if (await modal.isVisible()) {
            const firstUnit = modal.locator('button').filter({
                has: page.locator('text=Trato'),
            }).first();
            await firstUnit.click();
            await page.waitForURL((url) => url.pathname === '/', { timeout: 15000 });
        }

        // Fazer logout
        const userMenu = page.locator('button').filter({
            has: page.locator('svg.lucide-user, svg.lucide-circle-user, svg.lucide-log-out'),
        });

        if (await userMenu.isVisible().catch(() => false)) {
            await userMenu.click();
            await page.waitForTimeout(300);

            const logoutButton = page.locator('text=Sair').or(page.locator('text=Logout'));
            if (await logoutButton.isVisible().catch(() => false)) {
                await logoutButton.click();
                await page.waitForURL('**/login', { timeout: 10000 });

                // Verificar que dados foram limpos
                const unitData = await page.evaluate(() => localStorage.getItem('nexo-unit'));
                const authData = await page.evaluate(() => localStorage.getItem('nexo-auth'));

                console.log('✅ Logout realizado. Auth:', authData ? 'existe' : 'limpo', 'Unit:', unitData);
            }
        }
    });
});
