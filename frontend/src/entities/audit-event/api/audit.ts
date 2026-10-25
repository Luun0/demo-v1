import { useInfiniteQuery } from '@tanstack/react-query'
import { api, ApiError, type AuditEvent } from '@/shared/api'

interface Page {
  items: AuditEvent[]
  next: string | null
}

function toPage(r: { data?: AuditEvent[]; response: Response }): Page {
  if (!r.response.ok || !r.data) throw new ApiError(r.response.status)
  return { items: r.data, next: r.response.headers.get('X-Next-Cursor') }
}

export function useChildAudit(childId: string) {
  return useInfiniteQuery({
    queryKey: ['audit', childId],
    initialPageParam: null as string | null,
    getNextPageParam: (p: Page) => p.next,
    queryFn: async ({ pageParam }) =>
      toPage(await api.GET('/children/{childId}/audit-events', { params: { path: { childId }, query: pageParam ? { limit: 50, cursor: pageParam } : { limit: 50 } } })),
  })
}

export function useSystemAudit() {
  return useInfiniteQuery({
    queryKey: ['audit', 'system'],
    initialPageParam: null as string | null,
    getNextPageParam: (p: Page) => p.next,
    queryFn: async ({ pageParam }) => toPage(await api.GET('/audit-events', { params: { query: pageParam ? { limit: 50, cursor: pageParam } : { limit: 50 } } })),
  })
}
