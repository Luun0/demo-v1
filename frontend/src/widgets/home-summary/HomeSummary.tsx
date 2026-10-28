import { useTranslation } from 'react-i18next'
import type { Child, Role } from '@/shared/api'
import { EmptyState, ErrorState, Loading } from '@/shared/ui'
import { ChildCard, useChildren } from '@/entities/child'
import { useDiaryCount } from '@/entities/diary-entry'

/** Count comes from the server (X-Total-Count), already visibility-filtered (fixes C5). */
function DiaryCount({ child }: { child: Child }) {
  const { t } = useTranslation('dashboard')
  const q = useDiaryCount(child.id)
  if (q.data === undefined) return null
  return <p className="count-pop mt-2 text-sm text-ink">{t('visibleEntries', { count: q.data })}</p>
}

export function HomeSummary({ role }: { role: Role }) {
  const { t } = useTranslation('dashboard')
  const q = useChildren()
  if (q.isPending) return <Loading />
  if (q.isError) return <ErrorState error={q.error} onRetry={() => void q.refetch()} />
  if (q.data.length === 0) return <EmptyState>{t(role === 'specialist' ? 'noChildrenSpecialist' : 'noChildren')}</EmptyState>
  return (
    <ul className="stagger-list grid gap-4 sm:grid-cols-2">
      {q.data.map((c) => (
        <li key={c.id}>
          <ChildCard child={c} role={role} extra={role !== 'admin' && <DiaryCount child={c} />} />
        </li>
      ))}
    </ul>
  )
}
