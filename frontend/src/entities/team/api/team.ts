import { useQuery } from '@tanstack/react-query'
import { api, unwrap, type SpecialistSummary, type TeamMember } from '@/shared/api'

export const teamKeys = { team: (child: string) => ['team', child] as const, directory: ['specialists'] as const }

export function useTeam(childId: string) {
  return useQuery({
    queryKey: teamKeys.team(childId),
    queryFn: async (): Promise<TeamMember[]> => unwrap(await api.GET('/children/{childId}/team', { params: { path: { childId } } })),
  })
}

export function useSpecialistDirectory(enabled: boolean) {
  return useQuery({ queryKey: teamKeys.directory, enabled, queryFn: async (): Promise<SpecialistSummary[]> => unwrap(await api.GET('/specialists')) })
}

export function activeSpecialists(team: TeamMember[] | undefined, now = Date.now()): TeamMember[] {
  return (team ?? []).filter((m) => m.role === 'specialist' && (!m.ended_at || new Date(m.ended_at).getTime() > now))
}
