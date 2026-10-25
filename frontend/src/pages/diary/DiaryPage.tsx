import { useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { PageHeader } from '@/shared/ui'
import { can, useMe } from '@/entities/viewer'
import { CreateEntryDialog } from '@/features/diary-create-entry'
import { ChildHeader } from '@/widgets/child-nav'
import { DiaryFeed } from '@/widgets/diary-feed'

export function DiaryPage() {
  const { childId = '' } = useParams()
  const { t } = useTranslation('diary')
  const me = useMe()
  return (
    <>
      <PageHeader title={t('title')} actions={<CreateEntryDialog childId={childId} canChooseVisibility={can(me.role, 'diary:choose-visibility')} />} />
      <ChildHeader childId={childId} role={me.role} />
      <DiaryFeed childId={childId} viewer={me} />
    </>
  )
}
