import { expect, test } from '@playwright/test';

/**
 * NEXO - Sistema de Gestão para Barbearias
 * E2E Tests: Módulo de Serviços
 * 
 * Sprint 1.4.2 - Serviços Básicos
 * 
 * @description Testes end-to-end para o módulo de Serviços
 * Valida o fluxo completo: Login → Listagem → Criação → Edição → Exclusão
 */

// Configuração de teste
const TEST_USER = {
  email: 'andrey@tratodebarbados.com',
  password: '@Aa30019258',
};

// Configurar testes para rodarem em série (evitar concorrência)
test.describe.configure({ mode: 'serial' });

test.describe('Módulo de Serviços', () => {
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
    
    // Navegar para a página de serviços
    await page.goto('/cadastros/servicos', { waitUntil: 'domcontentloaded' });
    await page.waitForSelector('h1', { timeout: 10000 });
    
    // Fechar menu mobile se estiver aberto
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);
  });

  test('deve exibir a lista de serviços', async ({ page }) => {
    // Verificar título da página
    await expect(page.locator('h1')).toContainText('Serviços');
    
    // Verificar que as estatísticas estão visíveis
    await expect(page.locator('text=Total de Serviços')).toBeVisible();
    await expect(page.locator('text=Preço Médio')).toBeVisible();
    await expect(page.locator('text=Duração Média')).toBeVisible();
    await expect(page.locator('text=Comissão Média')).toBeVisible();
    
    // Verificar que a tabela existe
    await expect(page.locator('table')).toBeVisible();
  });

  test('deve criar um novo serviço', async ({ page }) => {
    // Clicar no botão "Novo Serviço"
    await page.getByRole('button', { name: 'Novo Serviço' }).click();
    
    // Aguardar modal abrir (usando data-slot="dialog-content")
    await expect(page.locator('[data-slot="dialog-content"]')).toBeVisible({ timeout: 10000 });
    await expect(page.getByRole('heading', { name: 'Novo Serviço' })).toBeVisible();
    
    // Preencher formulário
    const timestamp = Date.now();
    await page.fill('input[name="nome"]', `Teste Playwright ${timestamp}`);
    await page.fill('textarea[name="descricao"]', 'Serviço criado via teste automatizado Playwright');
    await page.fill('input[name="preco"]', '55.00');
    await page.fill('input[name="duracao"]', '45');
    await page.fill('input[name="comissao"]', '45');
    
    // Selecionar categoria usando Select do Radix UI
    await page.locator('button[role="combobox"]').click();
    await page.waitForTimeout(300); // Aguardar menu abrir
    const firstOption = page.locator('[role="option"]').first();
    await firstOption.waitFor({ state: 'visible' });
    await firstOption.click();
    
    // Preencher cor (input type="color" requer evaluate para setar valor)
    await page.locator('input[name="cor"]').first().evaluate((el: HTMLInputElement) => {
      el.value = '#3B82F6';
      el.dispatchEvent(new Event('input', { bubbles: true }));
      el.dispatchEvent(new Event('change', { bubbles: true }));
    });
    
    // Submeter formulário
    await page.getByRole('button', { name: 'Salvar' }).click();
    
    // Aguardar toast de sucesso (Sonner)
    const toast = page.locator('[data-sonner-toast]').filter({ hasText: /criado|sucesso/i });
    await expect(toast).toBeVisible({ timeout: 10000 });
    
    // Aguardar modal fechar
    await expect(page.locator('[data-slot="dialog-content"]')).not.toBeVisible({ timeout: 5000 });
    
    // Aguardar tabela recarregar e verificar que o serviço aparece na lista
    await page.waitForTimeout(1000);
    await expect(page.getByText(`Teste Playwright ${timestamp}`)).toBeVisible({ timeout: 5000 });
  });

  test('deve buscar serviços pelo nome', async ({ page }) => {
    // Obter primeiro serviço da lista
    const firstServiceName = await page.locator('table tbody tr').first().locator('td').first().textContent();
    
    if (!firstServiceName) {
      test.skip();
      return;
    }
    
    // Extrair primeiras palavras para buscar
    const searchTerm = firstServiceName.trim().split(' ')[0];
    
    // Digitar no campo de busca
    await page.fill('input[placeholder*="Buscar"]', searchTerm);
    
    // Aguardar debounce e resultado
    await page.waitForTimeout(1000);
    
    // Verificar que existem resultados
    await expect(page.locator('table tbody tr')).not.toHaveCount(0);
  });

  test('deve editar um serviço existente', async ({ page }) => {
    // Clicar no menu de ações do primeiro serviço
    await page.locator('table tbody tr').first().getByRole('button', { name: 'Abrir menu' }).click();
    
    // Aguardar menu aparecer e clicar em "Editar"
    await page.getByRole('menuitem', { name: 'Editar' }).click();
    
    // Aguardar modal abrir
    await expect(page.locator('[data-slot="dialog-content"]')).toBeVisible({ timeout: 5000 });
    await expect(page.getByRole('heading', { name: 'Editar Serviço' })).toBeVisible();
    
    // Modificar nome
    const timestamp = Date.now();
    await page.fill('input[name="nome"]', `Editado ${timestamp}`);
    
    // Modificar preço
    await page.fill('input[name="preco"]', '99.99');
    
    // Salvar
    await page.getByRole('button', { name: 'Salvar' }).click();
    
    // Aguardar toast de sucesso (Sonner)
    const toast = page.locator('[data-sonner-toast]').filter({ hasText: /atualizado|sucesso/i });
    await expect(toast).toBeVisible({ timeout: 10000 });
    
    // Verificar que o modal fechou
    await expect(page.locator('[data-slot="dialog-content"]')).not.toBeVisible({ timeout: 5000 });
  });

  test('deve desativar e reativar um serviço', async ({ page }) => {
    // Encontrar um serviço ativo e pegar seu nome
    const activeRow = page.locator('table tbody tr').filter({ hasText: 'Ativo' }).first();
    const serviceName = await activeRow.locator('td').first().textContent();
    
    // Clicar no menu de ações
    await activeRow.getByRole('button', { name: 'Abrir menu' }).click();
    
    // Aguardar menu abrir e verificar opção disponível
    await page.waitForTimeout(300);
    const menuItem = page.getByRole('menuitem').filter({ hasText: /Desativar|Ativar/ }).first();
    await menuItem.waitFor({ state: 'visible' });
    const action = await menuItem.textContent();
    
    // Clicar na opção
    await menuItem.click();
    
    // Aguardar toast de sucesso (Sonner)
    const toast = page.locator('[data-sonner-toast]');
    await expect(toast).toBeVisible({ timeout: 5000 });
    
    // Aguardar atualização da lista
    await page.waitForTimeout(1500);
    
    // Verificar mudança de status na linha específica do serviço
    if (serviceName) {
      const updatedRow = page.locator('table tbody tr').filter({ hasText: serviceName });
      if (action?.includes('Desativar')) {
        await expect(updatedRow.getByText('Inativo')).toBeVisible({ timeout: 3000 });
      } else {
        await expect(updatedRow.getByText('Ativo')).toBeVisible({ timeout: 3000 });
      }
    }
  });

  test('deve excluir um serviço', async ({ page }) => {
    // Criar um serviço temporário para exclusão
    await page.getByRole('button', { name: 'Novo Serviço' }).click();
    await expect(page.locator('[data-slot="dialog-content"]')).toBeVisible({ timeout: 5000 });
    
    const timestamp = Date.now();
    const serviceName = `Deletar ${timestamp}`;
    
    await page.fill('input[name="nome"]', serviceName);
    await page.fill('input[name="preco"]', '10.00');
    await page.fill('input[name="duracao"]', '10');
    await page.fill('input[name="comissao"]', '10');
    
    // Selecionar categoria
    await page.locator('button[role="combobox"]').click();
    await page.waitForTimeout(300);
    const firstOption = page.locator('[role="option"]').first();
    await firstOption.waitFor({ state: 'visible' });
    await firstOption.click();
    
    await page.getByRole('button', { name: 'Salvar' }).click();
    
    // Aguardar toast de criação
    const createToast = page.locator('[data-sonner-toast]').filter({ hasText: /criado|sucesso/i });
    await expect(createToast).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(500);
    
    // Buscar pelo serviço criado
    await page.fill('input[placeholder*="Buscar"]', serviceName);
    await page.waitForTimeout(1000);
    
    // Clicar no menu de ações da linha com o serviço
    const serviceRow = page.locator('table tbody tr').filter({ hasText: serviceName });
    await serviceRow.getByRole('button', { name: 'Abrir menu' }).click();
    
    // Mock do confirm dialog
    page.on('dialog', async (dialog) => {
      expect(dialog.message()).toContain('excluir');
      await dialog.accept();
    });
    
    // Aguardar menu abrir e clicar em "Excluir"
    await page.waitForTimeout(300);
    await page.getByRole('menuitem', { name: 'Excluir' }).click();
    
    // Aguardar toast de exclusão (Sonner)
    const deleteToast = page.locator('[data-sonner-toast]').filter({ hasText: /excluído|removido|sucesso/i });
    await expect(deleteToast).toBeVisible({ timeout: 10000 });
    
    // Aguardar atualização da lista
    await page.waitForTimeout(1000);
    
    // Verificar que o serviço não aparece mais na lista
    await expect(page.getByText(serviceName).first()).not.toBeVisible({ timeout: 5000 });
  });

  test('deve validar campos obrigatórios no formulário', async ({ page }) => {
    // Clicar no botão "Novo Serviço"
    await page.getByRole('button', { name: 'Novo Serviço' }).click();
    
    // Aguardar modal abrir
    await expect(page.locator('[data-slot="dialog-content"]')).toBeVisible({ timeout: 5000 });
    await page.waitForTimeout(500);
    
    // Tentar submeter sem preencher
    await page.getByRole('button', { name: 'Salvar' }).click();
    
    // Aguardar validação aparecer
    await page.waitForTimeout(500);
    
    // Verificar mensagens de validação (Zod/React Hook Form)
    // Nome, Preço e Duração são obrigatórios
    const errorMessages = page.locator('text=/obrigatório|required/i');
    await expect(errorMessages.first()).toBeVisible({ timeout: 3000 });
    
    // Verificar que o modal não fechou (validação bloqueou submit)
    await expect(page.locator('[data-slot="dialog-content"]')).toBeVisible();
  });

  test('deve exibir estatísticas corretas', async ({ page }) => {
    // Verificar que os cards de estatísticas existem
    await expect(page.locator('text=Total de Serviços')).toBeVisible();
    await expect(page.locator('text=Preço Médio')).toBeVisible();
    await expect(page.locator('text=Duração Média')).toBeVisible();
    await expect(page.locator('text=Comissão Média')).toBeVisible();
  });

  test('deve filtrar serviços por categoria via badge', async ({ page }) => {
    // Aguardar listagem carregar
    await page.waitForTimeout(1000);
    
    // Clicar em um badge de categoria (se existir)
    const categoryBadge = page.locator('table tbody tr').first().locator('span[class*="badge"]').filter({ hasText: /[A-Za-z]/ }).first();
    
    if (await categoryBadge.isVisible()) {
      const categoryName = await categoryBadge.textContent();
      await categoryBadge.click();
      
      // Verificar que todos os resultados são dessa categoria
      const rows = page.locator('table tbody tr');
      const count = await rows.count();
      
      for (let i = 0; i < Math.min(count, 5); i++) {
        await expect(rows.nth(i)).toContainText(categoryName || '');
      }
    }
  });
});
