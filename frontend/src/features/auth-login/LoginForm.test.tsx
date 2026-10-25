import { describe, expect, it, vi } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { server } from '@/shared/lib/test/msw'
import { renderWithProviders } from '@/shared/lib/test/render'
import { getAccessToken, setAccessToken } from '@/shared/api'
import { LoginForm } from './LoginForm'

describe('LoginForm', () => {
  it('shows a generic error on 401 (no account enumeration)', async () => {
    server.use(http.post('*/api/v1/auth/login', () => HttpResponse.json({ title: 'x', status: 401 }, { status: 401 })))
    renderWithProviders(<LoginForm />)
    await userEvent.type(screen.getByLabelText(/email/i), 'nobody@example.test')
    await userEvent.type(screen.getByLabelText(/password/i), 'wrong-password')
    await userEvent.click(screen.getByRole('button', { name: /log in/i }))
    expect((await screen.findAllByRole('alert')).length).toBeGreaterThan(0)
    expect(getAccessToken()).toBeNull()
  })

  it('validates email before calling the API', async () => {
    renderWithProviders(<LoginForm />)
    await userEvent.type(screen.getByLabelText(/email/i), 'not-an-email')
    await userEvent.click(screen.getByRole('button', { name: /log in/i }))
    expect(screen.getByLabelText(/email/i)).toHaveAttribute('aria-invalid', 'true')
  })

  it('stores the access token in memory only on success', async () => {
    setAccessToken(null)
    server.use(http.post('*/api/v1/auth/login', () => HttpResponse.json({ access_token: 'mem-only', token_type: 'Bearer', expires_in: 900 })))
    renderWithProviders(<LoginForm />)
    await userEvent.type(screen.getByLabelText(/email/i), 'parent1@example.test')
    await userEvent.type(screen.getByLabelText(/password/i), 'Synthetic-Demo-2026')
    await userEvent.click(screen.getByRole('button', { name: /log in/i }))
    await screen.findByRole('button', { name: /log in/i })
    await vi.waitFor(() => expect(getAccessToken()).toBe('mem-only'))
    expect(JSON.stringify({ ...localStorage })).not.toContain('mem-only')
  })
})
