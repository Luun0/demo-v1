import { useTranslation } from 'react-i18next'
import type { User } from '@/shared/api'
import { EmptyState, ErrorState, Loading } from '@/shared/ui'
import { SessionRow, useSessions } from '@/entities/session'
import { SessionStatusSelect } from '@/features/session-update'

export function ScheduleList({ childId, viewer }: { childId: string; viewer: User }) {
  const { t } = useTranslation('schedule')
  const q = useSessions(childId)
  if (q.isPending) return <Loading />
  if (q.isError) return <ErrorState error={q.error} onRetry={() => void q.refetch()} />
  if (q.data.length === 0) return <EmptyState>{t('empty')}</EmptyState>
  return (
    <ul className="space-y-3">
      {q.data.map((s) => (
        <SessionRow key={s.id} session={s} actions={(viewer.role === 'parent' || s.specialist_id === viewer.id) && <SessionStatusSelect session={s} />} />
      ))}
    </ul>
  )
}
