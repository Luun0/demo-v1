import { useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { formatDate } from '@/shared/lib'
import { ErrorState, Loading, PageHeader } from '@/shared/ui'
import { useChild } from '@/entities/child'
import { useMe } from '@/entities/viewer'
import { EditChildDialog } from '@/features/child-update-profile'
import { ExportButton } from '@/features/child-export-data'
import { DeleteChildButton } from '@/features/child-delete'
import { ChildHeader } from '@/widgets/child-nav'

export function ChildProfilePage() {
  const { childId = '' } = useParams()
  const { t, i18n } = useTranslation('child')
  const me = useMe()
  const q = useChild(childId)
  if (q.isPending) return <Loading />
  if (q.isError) return <ErrorState error={q.error} onRetry={() => void q.refetch()} />
  const c = q.data
  const owner = me.role === 'parent'
  return (
    <>
      <PageHeader title={t('profileTitle')} />
      <ChildHeader childId={childId} role={me.role} />
      <dl className="surface-lift grid gap-4 rounded-2xl border border-line bg-surface p-6 shadow-soft sm:grid-cols-2">
        <div><dt className="text-sm text-ink-soft">{t('name')}</dt><dd className="font-medium">{c.first_name} {c.last_name}</dd></div>
        <div><dt className="text-sm text-ink-soft">{t('birthDate')}</dt><dd className="font-medium">{formatDate(c.birth_date, i18n.language)}</dd></div>
      </dl>
      {owner && (
        <section className="mt-8 space-y-6" aria-labelledby="manage-h">
          <h2 id="manage-h" className="text-lg font-semibold">{t('manage')}</h2>
          <EditChildDialog child={c} />
          <div className="flex flex-wrap gap-3">
            <ExportButton childId={c.id} />
            <DeleteChildButton childId={c.id} />
          </div>
        </section>
      )}
    </>
  )
}
