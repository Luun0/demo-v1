import { useQuery } from '@tanstack/react-query'
import { api, unwrap, type User } from '@/shared/api'
import { useAccessToken } from '../model/session'

export const viewerKey = ['viewer'] as const

export function useViewer() {
  const token = useAccessToken()
  return useQuery({
    queryKey: [...viewerKey, token],
    queryFn: async (): Promise<User> => unwrap(await api.GET('/auth/me')),
    enabled: !!token,
    staleTime: 60_000,
  })
}

/** For pages behind RequireAuth: the viewer is guaranteed to be loaded. */
export function useMe(): User {
  const q = useViewer()
  if (!q.data) throw new Error('useMe used outside RequireAuth')
  return q.data
}
