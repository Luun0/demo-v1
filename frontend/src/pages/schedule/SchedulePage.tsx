import { useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { PageHeader } from '@/shared/ui'
import { activeSpecialists, useTeam } from '@/entities/team'
import { useMe } from '@/entities/viewer'
import { CreateSessionDialog } from '@/features/session-create'
import { ChildHeader } from '@/widgets/child-nav'
import { ScheduleList } from '@/widgets/schedule-list'

export function SchedulePage() {
  const { childId = '' } = useParams()
  const { t } = useTranslation('schedule')
  const me = useMe()
  const team = useTeam(childId)
  const all = activeSpecialists(team.data ?? [])
  const specialists = me.role === 'specialist' ? all.filter((s) => s.user_id === me.id) : all
  return (
    <>
      <PageHeader title={t('title')} actions={specialists.length > 0 && <CreateSessionDialog childId={childId} specialists={specialists} />} />
      <ChildHeader childId={childId} role={me.role} />
      <ScheduleList childId={childId} viewer={me} />
    </>
  )
}
