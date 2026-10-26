import { describe, expect, it } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/shared/lib/test/render'
import { ThemeSwitch } from './ThemeSwitch'

describe('ThemeSwitch', () => {
  it('cycles system, light and dark and persists the selected mode', async () => {
    localStorage.clear()
    const user = userEvent.setup()
    renderWithProviders(<ThemeSwitch />)
    const button = screen.getByRole('button', { name: /Theme: System/ })

    await user.click(button)
    expect(localStorage.getItem('ssp_theme')).toBe('light')
    expect(button).toHaveAccessibleName(/Theme: Light/)
    expect(document.documentElement.dataset.theme).toBe('light')

    await user.click(button)
    expect(localStorage.getItem('ssp_theme')).toBe('dark')
    expect(document.documentElement.dataset.theme).toBe('dark')

    await user.click(button)
    expect(localStorage.getItem('ssp_theme')).toBe('system')
    expect(button).toHaveAccessibleName(/Theme: System/)
  })
})
