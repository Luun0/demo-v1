import { useInfiniteQuery, useQuery } from '@tanstack/react-query'
import { api, ApiError, type DiaryEntry, type DiaryEntryVersion, type Visibility } from '@/shared/api'

export interface DiaryFilters {
  q?: string
  visibility?: Visibility | ''
}

export const diaryKeys = {
  all: (child: string) => ['diary', child] as const,
  list: (child: string, f: DiaryFilters) => ['diary', child, 'list', f] as const,
  count: (child: string) => ['diary', child, 'count'] as const,
  versions: (child: string, entry: string) => ['diary', child, 'versions', entry] as const,
}

interface Page {
  items: DiaryEntry[]
  next: string | null
  total: number
}

async function fetchPage(childId: string, f: DiaryFilters, cursor: string | null, limit = 20): Promise<Page> {
  const query: Record<string, string | number> = { limit }
  if (f.q) query.q = f.q
  if (f.visibility) query.visibility = f.visibility
  if (cursor) query.cursor = cursor
  const r = await api.GET('/children/{childId}/diary-entries', { params: { path: { childId }, query } })
  if (!r.response.ok || !r.data) throw new ApiError(r.response.status)
  return { items: r.data, next: r.response.headers.get('X-Next-Cursor'), total: Number(r.response.headers.get('X-Total-Count') ?? r.data.length) }
}

export function useDiary(childId: string, f: DiaryFilters) {
  return useInfiniteQuery({
    queryKey: diaryKeys.list(childId, f),
    queryFn: ({ pageParam }) => fetchPage(childId, f, pageParam),
    initialPageParam: null as string | null,
    getNextPageParam: (last) => last.next,
  })
}

/** Server-side count of entries visible to the caller (never client-side). */
export function useDiaryCount(childId: string) {
  return useQuery({ queryKey: diaryKeys.count(childId), queryFn: async () => (await fetchPage(childId, {}, null, 1)).total })
}

export function useVersions(childId: string, entryId: string, enabled: boolean) {
  return useQuery({
    queryKey: diaryKeys.versions(childId, entryId),
    enabled,
    queryFn: async (): Promise<DiaryEntryVersion[]> => {
      const r = await api.GET('/children/{childId}/diary-entries/{entryId}/versions', { params: { path: { childId, entryId } } })
      if (!r.response.ok || !r.data) throw new ApiError(r.response.status)
      return r.data
    },
  })
}
