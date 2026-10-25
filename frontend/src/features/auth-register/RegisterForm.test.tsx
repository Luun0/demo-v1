import { describe, expect, it } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/shared/lib/test/render'
import { RegisterForm } from './RegisterForm'

describe('RegisterForm', () => {
  it('blocks submit and marks fields invalid when empty', async () => {
    renderWithProviders(<RegisterForm />)
    await userEvent.click(screen.getByRole('button', { name: /^register$/i }))
    const invalid = document.querySelectorAll('[aria-invalid="true"]')
    expect(invalid.length).toBeGreaterThan(0)
    for (const el of invalid) expect(el.getAttribute('aria-describedby')).toBeTruthy()
  })
})
