import { useTranslation } from 'react-i18next'
import { PageHeader } from '@/shared/ui'
import { can, useMe } from '@/entities/viewer'
import { CreateChildDialog } from '@/features/child-create'
import { HomeSummary } from '@/widgets/home-summary'

export function DashboardPage() {
  const { t } = useTranslation('dashboard')
  const me = useMe()
  return (
    <>
      <PageHeader title={t('greeting', { name: me.first_name })} subtitle={t(`subtitle.${me.role}`)} actions={can(me.role, 'child:create') && <CreateChildDialog />} />
      <HomeSummary role={me.role} />
    </>
  )
}
