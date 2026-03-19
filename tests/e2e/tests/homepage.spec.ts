import { test, expect } from '@playwright/test';

test.describe('Homepage', () => {
  test('should display the homepage correctly', async ({ page }) => {
    await page.goto('/');
    
    // Check that the page loads
    await expect(page).toHaveTitle(/Devon Farm Sales/);
    
    // Check for key elements
    await expect(page.locator('h1')).toContainText('Devon Farm Sales');
    
    // Check that login link is present
    const loginLink = page.locator('a[href="/login"]');
    await expect(loginLink).toBeVisible();
  });

  test('should navigate to login page', async ({ page }) => {
    await page.goto('/');
    
    // Click login link
    await page.click('a[href="/login"]');
    
    // Should be on login page
    await expect(page).toHaveURL('/login');
    await expect(page.locator('h1')).toContainText(/Log in/i);
  });
});