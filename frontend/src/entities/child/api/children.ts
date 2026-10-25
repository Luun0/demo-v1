import { useQuery } from '@tanstack/react-query'
import { api, unwrap, type Child } from '@/shared/api'

export const childKeys = {
  all: ['children'] as const,
  one: (id: string) => ['children', id] as const,
}

export function useChildren() {
  return useQuery({ queryKey: childKeys.all, queryFn: async (): Promise<Child[]> => unwrap(await api.GET('/children')) })
}

export function useChild(id: string) {
  return useQuery({
    queryKey: childKeys.one(id),
    queryFn: async (): Promise<Child> => unwrap(await api.GET('/children/{childId}', { params: { path: { childId: id } } })),
  })
}
