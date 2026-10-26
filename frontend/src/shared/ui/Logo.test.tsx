import { describe, expect, it } from 'vitest'
import { screen } from '@testing-library/react'
import { renderWithProviders } from '@/shared/lib/test/render'
import { Logo } from './Logo'

describe('Logo', () => {
  it('exposes a named SVG mark and keeps the wordmark as text', () => {
    renderWithProviders(<Logo />)

    expect(screen.getByRole('img', { name: 'SEN Support mark' })).toBeInTheDocument()
    expect(screen.getByText('SEN Support')).toBeInTheDocument()
  })

  it('preserves the inverse treatment for brand surfaces', () => {
    const { container } = renderWithProviders(<Logo inverse />)

    expect(container.firstElementChild).toHaveClass('text-on-brand')
  })
})
