import { useQuery } from '@tanstack/react-query'
import { api, unwrap, type SupportSession } from '@/shared/api'

export const sessionKeys = { list: (child: string) => ['sessions', child] as const }

export function useSessions(childId: string) {
  return useQuery({
    queryKey: sessionKeys.list(childId),
    queryFn: async (): Promise<SupportSession[]> => unwrap(await api.GET('/children/{childId}/sessions', { params: { path: { childId } } })),
  })
}
