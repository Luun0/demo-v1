import { useTranslation } from 'react-i18next'
import type { User } from '@/shared/api'
import { formatDate } from '@/shared/lib'
import { Badge, EmptyState, ErrorState, Loading } from '@/shared/ui'
import { activeSpecialists, useTeam } from '@/entities/team'
import { AssignSpecialistForm } from '@/features/team-assign-specialist'
import { EndAssignmentButton } from '@/features/team-end-assignment'

export function TeamPanel({ childId, viewer }: { childId: string; viewer: User }) {
  const { t, i18n } = useTranslation('team')
  const q = useTeam(childId)
  if (q.isPending) return <Loading />
  if (q.isError) return <ErrorState error={q.error} onRetry={() => void q.refetch()} />
  const guardians = q.data.filter((m) => m.role === 'parent')
  const specialists = q.data.filter((m) => m.role === 'specialist')
  const active = activeSpecialists(q.data)
  const manager = viewer.role === 'admin' || guardians.some((g) => g.user_id === viewer.id && g.is_primary)
  return (
    <div className="space-y-8">
      <section aria-labelledby="guardians-h">
        <h2 id="guardians-h" className="mb-3 text-lg font-semibold">{t('guardians')}</h2>
        <ul className="space-y-2">
          {guardians.map((g) => (
            <li key={g.user_id} className="flex items-center gap-3 rounded-lg border border-line bg-surface p-3">
              {g.first_name} {g.last_name} {g.is_primary && <Badge tone="brand">{t('primary')}</Badge>}
            </li>
          ))}
        </ul>
      </section>
      <section aria-labelledby="specialists-h">
        <h2 id="specialists-h" className="mb-3 text-lg font-semibold">{t('specialists')}</h2>
        {specialists.length === 0 ? <EmptyState>{t('empty')}</EmptyState> : (
          <ul className="space-y-2">
            {specialists.map((s) => {
              const isActive = active.includes(s)
              return (
                <li key={s.user_id} className="flex flex-wrap items-center gap-3 rounded-lg border border-line bg-surface p-3">
                  <span className="font-medium">{s.first_name} {s.last_name}</span>
                  <span className="text-sm text-ink-soft">
                    {s.assigned_at && t('assigned', { date: formatDate(s.assigned_at, i18n.language) })}
                    {!isActive && s.ended_at && ` · ${t('ended', { date: formatDate(s.ended_at, i18n.language) })}`}
                  </span>
                  {manager && isActive && <span className="ml-auto"><EndAssignmentButton childId={childId} specialistId={s.user_id} name={`${s.first_name} ${s.last_name}`} /></span>}
                </li>
              )
            })}
          </ul>
        )}
        {manager && <div className="mt-4"><AssignSpecialistForm childId={childId} exclude={active.map((a) => a.user_id)} /></div>}
      </section>
    </div>
  )
}
