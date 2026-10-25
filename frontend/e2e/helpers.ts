import { expect, type Page } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

export const PASSWORD = 'Synthetic-Demo-2026'

export async function preferEnglish(page: Page) {
  await page.addInitScript(() => localStorage.setItem('ssp_lang', 'en'))
}

export async function login(page: Page, email: string) {
  await preferEnglish(page)
  await page.goto('/login')
  await page.getByLabel('Email').fill(email)
  await page.getByLabel('Password').fill(PASSWORD)
  await page.getByRole('button', { name: 'Log in' }).click()
  await expect(page).toHaveURL(/\/app$/)
}

/** WCAG 2.1 AA via axe; any violation fails the test. */
export async function expectAccessible(page: Page) {
  const r = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa']).analyze()
  expect(r.violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(', ')}`)).toEqual([])
}
