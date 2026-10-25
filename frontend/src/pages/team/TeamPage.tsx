import { useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { PageHeader } from '@/shared/ui'
import { useMe } from '@/entities/viewer'
import { ChildHeader } from '@/widgets/child-nav'
import { TeamPanel } from '@/widgets/team-panel'

export function TeamPage() {
  const { childId = '' } = useParams()
  const { t } = useTranslation('team')
  const me = useMe()
  return (
    <>
      <PageHeader title={t('title')} />
      <ChildHeader childId={childId} role={me.role} />
      <TeamPanel childId={childId} viewer={me} />
    </>
  )
}
