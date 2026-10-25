import { useParams } from 'react-router'
import { useTranslation } from 'react-i18next'
import { Button, ErrorState, Loading, PageHeader } from '@/shared/ui'
import { AuditTable, useChildAudit } from '@/entities/audit-event'
import { useMe } from '@/entities/viewer'
import { ChildHeader } from '@/widgets/child-nav'

export function AccessLogPage() {
  const { childId = '' } = useParams()
  const { t } = useTranslation('access')
  const me = useMe()
  const q = useChildAudit(childId)
  return (
    <>
      <PageHeader title={t('title')} subtitle={t('subtitle')} />
      <ChildHeader childId={childId} role={me.role} />
      {q.isPending ? <Loading /> : q.isError ? <ErrorState error={q.error} onRetry={() => void q.refetch()} /> : (
        <>
          <AuditTable caption={t('title')} events={q.data.pages.flatMap((p) => p.items)} />
          {q.hasNextPage && <Button variant="secondary" className="mt-4" busy={q.isFetchingNextPage} onClick={() => void q.fetchNextPage()}>{t('common:loadMore')}</Button>}
        </>
      )}
    </>
  )
}
