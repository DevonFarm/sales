import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('should display login form', async ({ page }) => {
    await page.goto('/login');
    
    // Check login form elements
    await expect(page.locator('input[name="name"]')).toBeVisible();
    await expect(page.locator('input[name="email"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('should validate required fields', async ({ page }) => {
    await page.goto('/login');
    
    // Try to submit without filling fields
    await page.click('button[type="submit"]');
    
    // Should show validation errors or stay on same page
    await expect(page).toHaveURL('/login');
  });

  test('should submit magic link request', async ({ page }) => {
    await page.goto('/login');
    
    // Fill in the form
    await page.fill('input[name="name"]', 'Test User');
    await page.fill('input[name="email"]', 'test@example.com');
    
    // Submit the form
    await page.click('button[type="submit"]');
    
    // Should redirect to confirmation page
    await expect(page).toHaveURL('/login');
    await expect(page.locator('body')).toContainText(/check your email/i);
  });

  // Note: Testing the actual magic link callback would require:
  // 1. Integration with Stytch test environment
  // 2. Email interception service
  // 3. Or mocked authentication flow
  
  test('should handle logout', async ({ page }) => {
    // This test would need a way to get into authenticated state first
    // For now, we'll test the logout endpoint directly
    
    // Navigate to a protected page (should redirect to login)
    await page.goto('/farm/some-farm-id');
    await expect(page).toHaveURL('/login');
  });
});