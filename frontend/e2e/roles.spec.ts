import { expect, test } from '@playwright/test'
import { expectAccessible, login, preferEnglish } from './helpers'

test('guest: landing and login are accessible; /app redirects to login', async ({ page }) => {
  await preferEnglish(page)
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1 })).toHaveAccessibleName('Support — together')
  await expectAccessible(page)
  await page.emulateMedia({ colorScheme: 'dark' })
  await expectAccessible(page)
  await page.goto('/login')
  await expectAccessible(page)
  await page.emulateMedia({ colorScheme: 'light' })
  await page.goto('/app')
  await expect(page).toHaveURL(/\/login$/)
  await expectAccessible(page)
  await page.goto('/register')
  await expectAccessible(page)
  await page.goto('/forgot-password')
  await expectAccessible(page)
  await page.goto('/reset-password?token=invalid-demo-token')
  await expectAccessible(page)
})

test('parent: dashboard, diary without specialist-only notes, session restore, no token in storage', async ({ page }) => {
  await login(page, 'parent1@example.test')
  await expect(page.getByRole('heading', { level: 1 })).toContainText('Aida')
  await expectAccessible(page)
  await page.getByRole('link', { name: /Mira/ }).first().click()
  await expectAccessible(page)
  await page.getByRole('navigation', { name: /sections/i }).getByRole('link', { name: 'Diary' }).click()
  await expect(page.getByText('Synthetic note: enjoyed the reading game.')).toBeVisible()
  await expect(page.getByText('Synthetic specialist-only note.')).toHaveCount(0)
  await expect(page.getByLabel('Visibility')).toHaveCount(0)
  await expectAccessible(page)

  const note = `E2E parent note ${Date.now()}`
  await page.getByRole('button', { name: 'New entry' }).click()
  await expect(page.getByRole('dialog').getByLabel('Visibility')).toHaveCount(0)
  await page.getByRole('dialog').getByLabel('Text').fill(note)
  await page.getByRole('dialog').getByRole('button', { name: 'Save' }).click()
  await expect(page.getByText(note)).toBeVisible()

  await page.reload()
  await expect(page.getByText(note)).toBeVisible()
  const storage = await page.evaluate(() => JSON.stringify({ ...localStorage }) + JSON.stringify({ ...sessionStorage }))
  expect(storage).not.toMatch(/eyJ/)

  const ownEntry = page.locator('article').filter({ hasText: note })
  await ownEntry.getByRole('button', { name: 'Delete entry' }).click()
  const confirmation = page.getByRole('dialog')
  await expect(confirmation).toBeVisible()
  await expectAccessible(page)
  await confirmation.getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByText(note)).toHaveCount(0)
  await expectAccessible(page)

  await page.getByRole('link', { name: 'Access log' }).first().click()
  await expect(page.getByRole('table')).toBeVisible()
  await expectAccessible(page)
})

test('specialist: sees specialist-only notes; cannot open access log', async ({ page }) => {
  await login(page, 'specialist1@example.test')
  await page.getByRole('link', { name: /Mira/ }).first().click()
  await page.getByRole('navigation', { name: /sections/i }).getByRole('link', { name: 'Diary' }).click()
  await expect(page.getByText('Synthetic specialist-only note.')).toBeVisible()
  await expect(page.getByLabel('Visibility')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Delete entry' })).toHaveCount(0)
  await expectAccessible(page)
  const url = page.url().replace(/\/diary$/, '/access-log')
  await page.goto(url)
  await expect(page.getByRole('heading', { name: 'Access denied' })).toBeVisible()
  await page.goto(url.replace(/\/access-log$/, '/schedule'))
  await expect(page.getByText('Speech session (synthetic)').first()).toBeVisible()
  await expectAccessible(page)
})

test('admin: users and audit; diary is not reachable', async ({ page }) => {
  await login(page, 'admin@example.test')
  await page.getByRole('link', { name: 'Users' }).first().click()
  await expect(page.getByRole('table')).toContainText('parent1@example.test')
  await expectAccessible(page)
  await page.getByRole('link', { name: 'Audit' }).first().click()
  await expect(page.getByRole('table')).toBeVisible()
  await expectAccessible(page)
  await page.goto('/app')
  await page.getByRole('link', { name: /Mira/ }).first().click()
  await page.goto(page.url() + '/diary')
  await expect(page.getByRole('heading', { name: 'Access denied' })).toBeVisible()
})

test('foreign child id is indistinguishable from missing (404 UI) @idor', async ({ page }) => {
  await login(page, 'parent2@example.test')
  await page.goto('/app/children/00000000-0000-4000-8000-000000000000')
  await expect(page.getByRole('main').getByRole('alert')).toBeVisible()
  await expect(page.getByRole('main')).not.toContainText('Mira')
})

test('mobile guest: support help works in dark mode without overflow @mobile', async ({ page }) => {
  await preferEnglish(page)
  await page.goto('/')
  await expectAccessible(page)
  await page.emulateMedia({ colorScheme: 'dark' })
  await expectAccessible(page)

  await page.getByRole('button', { name: 'Help' }).click()
  await expect(page.getByRole('region', { name: 'SEN Support help' })).toBeVisible()
  await expectAccessible(page)
  await page.getByRole('textbox', { name: 'Your question' }).fill('Who can see the diary?')
  await page.getByRole('button', { name: 'Send question' }).click()
  await expect(page.getByRole('log')).toContainText('Platform administrators cannot see it.')
  await expectAccessible(page)

  await page.emulateMedia({ reducedMotion: 'reduce' })
  const layout = await page.evaluate(() => ({
    viewport: document.documentElement.clientWidth,
    content: document.documentElement.scrollWidth,
    reducedMotion: matchMedia('(prefers-reduced-motion: reduce)').matches,
    panelAnimation: getComputedStyle(document.querySelector('#support-chat-panel')!).animationDuration,
  }))
  expect(layout.content).toBeLessThanOrEqual(layout.viewport)
  expect(layout.reducedMotion).toBe(true)
  expect(Number.parseFloat(layout.panelAnimation)).toBeLessThan(0.001)
})

test('mobile: drawer opens, traps focus and closes on Escape @mobile', async ({ page }) => {
  await login(page, 'parent1@example.test')
  const open = page.getByRole('button', { name: 'More' })
  await open.click()
  const drawer = page.getByRole('dialog')
  await expect(drawer).toBeVisible()
  await expectAccessible(page)
  await page.keyboard.press('Escape')
  await expect(drawer).toBeHidden()
  await expect(open).toBeFocused()
})
