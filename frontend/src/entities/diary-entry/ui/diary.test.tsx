import { describe, expect, it, vi } from 'vitest'
import { screen } from '@testing-library/react'
import { renderWithProviders } from '@/shared/lib/test/render'
import type { DiaryEntry } from '@/shared/api'
import { EntryCard } from './EntryCard'
import { EntryForm } from './EntryForm'

const entry: DiaryEntry = {
  id: 'e1', child_id: 'c1', author_id: 'u1', type: 'observation', visibility: 'shared',
  content: '<img src=x onerror="alert(1)"><script>alert(2)</script>', observed_at: '2026-10-01T10:00:00Z',
  created_at: '2026-10-01T10:00:00Z', updated_at: '2026-10-01T10:00:00Z',
} as DiaryEntry

describe('diary entry UI', () => {
  it('renders user content as plain text (XSS)', () => {
    const { container } = renderWithProviders(<EntryCard entry={entry} />)
    expect(container.querySelector('img, script')).toBeNull()
    expect(screen.getByText(/onerror/)).toBeInTheDocument()
  })

  it('parents get no visibility choice', () => {
    renderWithProviders(<EntryForm canChooseVisibility={false} busy={false} onSubmit={vi.fn()} />)
    expect(screen.queryByLabelText(/visibility/i)).toBeNull()
  })

  it('specialists can choose specialists_only', () => {
    renderWithProviders(<EntryForm canChooseVisibility busy={false} onSubmit={vi.fn()} />)
    const sel = screen.getByLabelText(/visibility/i)
    expect(sel.querySelectorAll('option')).toHaveLength(2)
  })
})
