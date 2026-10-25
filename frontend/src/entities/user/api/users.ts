import { useQuery } from '@tanstack/react-query'
import { api, unwrap, type User } from '@/shared/api'

export const userKeys = { all: ['users'] as const }

export function useUsers() {
  return useQuery({ queryKey: userKeys.all, queryFn: async (): Promise<User[]> => unwrap(await api.GET('/users', { params: { query: { limit: 100 } } })) })
}
