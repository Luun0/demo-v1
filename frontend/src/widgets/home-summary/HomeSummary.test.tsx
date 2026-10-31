import { describe, expect, it } from 'vitest'
import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import type { Child, DiaryEntry } from '@/shared/api'
import { server } from '@/shared/lib/test/msw'
import { renderWithProviders } from '@/shared/lib/test/render'
import { HomeSummary } from './HomeSummary'

const child: Child = {
  id: '10000000-0000-4000-8000-000000000001', first_name: 'Mira', last_name: 'Sample',
  birth_date: '2018-04-12', created_at: '2026-10-23T10:00:00Z', updated_at: '2026-10-23T10:00:00Z',
}

const diary: DiaryEntry[] = [
  { id: 'e1', child_id: child.id, author_id: 'a1', type: 'recommendation', visibility: 'shared', content: 'Practice the puzzle after a short break.', observed_at: '2026-10-23T09:00:00Z', created_at: '2026-10-23T09:01:00Z', updated_at: '2026-10-23T09:01:00Z' },
  { id: 'e2', child_id: child.id, author_id: 'a2', type: 'recommendation', visibility: 'specialists_only', content: 'Private specialist note.', observed_at: '2026-10-22T09:00:00Z', created_at: '2026-10-22T09:01:00Z', updated_at: '2026-10-22T09:01:00Z' },
  { id: 'e3', child_id: child.id, author_id: 'a1', type: 'observation', visibility: 'shared', content: 'Shared observation, not a recommendation.', observed_at: '2026-10-21T09:00:00Z', created_at: '2026-10-21T09:01:00Z', updated_at: '2026-10-21T09:01:00Z' },
]

function serveChildren() {
  server.use(
    http.get('*/api/v1/children', () => HttpResponse.json([child])),
    http.get('*/api/v1/children/:childId/diary-entries', () => HttpResponse.json(diary, { headers: { 'X-Total-Count': '3' } })),
  )
}

describe('HomeSummary recommendations', () => {
  it('shows only family-visible recommendation entries and links to that child diary', async () => {
    let diaryRequests = 0
    server.use(
      http.get('*/api/v1/children', () => HttpResponse.json([child])),
      http.get('*/api/v1/children/:childId/diary-entries', () => {
        diaryRequests += 1
        return HttpResponse.json(diary, { headers: { 'X-Total-Count': '3' } })
      }),
    )
    renderWithProviders(<HomeSummary role="parent" />)

    expect(await screen.findByText('Practice the puzzle after a short break.')).toBeInTheDocument()
    expect(screen.getByText('3 diary entries')).toBeInTheDocument()
    expect(diaryRequests).toBe(1)
    expect(screen.queryByText('Private specialist note.')).not.toBeInTheDocument()
    expect(screen.queryByText('Shared observation, not a recommendation.')).not.toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'All diary entries' })).toHaveAttribute('href', `/app/children/${child.id}/diary`)
  })

  it('does not fetch or display diary content for administrators', async () => {
    let diaryRequests = 0
    serveChildren()
    server.use(http.get('*/api/v1/children/:childId/diary-entries', () => {
      diaryRequests += 1
      return HttpResponse.json(diary)
    }))
    renderWithProviders(<HomeSummary role="admin" />)

    expect(await screen.findByText('Mira Sample')).toBeInTheDocument()
    expect(screen.queryByText('Recent recommendations')).not.toBeInTheDocument()
    expect(diaryRequests).toBe(0)
  })

  it('does not show the parent recommendation panel to specialists', async () => {
    serveChildren()
    renderWithProviders(<HomeSummary role="specialist" />)

    expect(await screen.findByText('Mira Sample')).toBeInTheDocument()
    expect(screen.queryByText('Recent recommendations')).not.toBeInTheDocument()
    expect(screen.queryByText('Practice the puzzle after a short break.')).not.toBeInTheDocument()
  })

  it('shows an honest empty state when the recent page has no recommendations', async () => {
    serveChildren()
    server.use(http.get('*/api/v1/children/:childId/diary-entries', () => HttpResponse.json([diary[2]])))
    renderWithProviders(<HomeSummary role="parent" />)

    expect(await screen.findByText('No recommendations appear in the latest entries.')).toBeInTheDocument()
  })

  it('lets parents load older family-visible recommendations across diary pages', async () => {
    const user = userEvent.setup()
    serveChildren()
    const observations = Array.from({ length: 20 }, (_, index) => ({
      ...diary[2], id: `old-${index}`, content: `Observation ${index}`,
    }))
    server.use(http.get('*/api/v1/children/:childId/diary-entries', ({ request }) => {
      const url = new URL(request.url)
      if (url.searchParams.get('limit') === '1') return HttpResponse.json([], { headers: { 'X-Total-Count': '21' } })
      if (url.searchParams.has('cursor')) return HttpResponse.json([{
        ...diary[0], id: 'older-recommendation', content: 'A recommendation from earlier in the diary.',
      }], { headers: { 'X-Total-Count': '21' } })
      return HttpResponse.json(observations, { headers: { 'X-Total-Count': '21', 'X-Next-Cursor': 'older-page' } })
    }))
    renderWithProviders(<HomeSummary role="parent" />)

    expect(await screen.findByText('No recommendations appear in the latest entries.')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'Show older recommendations' }))
    expect(await screen.findByText('A recommendation from earlier in the diary.')).toBeInTheDocument()
  })

  it('shows a recoverable message if the recommendations request fails', async () => {
    serveChildren()
    server.use(http.get('*/api/v1/children/:childId/diary-entries', () => new HttpResponse(null, { status: 500 })))
    renderWithProviders(<HomeSummary role="parent" />)

    expect(await screen.findByText('Recommendations could not be loaded.')).toBeInTheDocument()
  })
})
