import { useTranslation } from 'react-i18next'
import { Button, ErrorState, Loading, PageHeader } from '@/shared/ui'
import { AuditTable, useSystemAudit } from '@/entities/audit-event'

export function AdminAuditPage() {
  const { t } = useTranslation('admin')
  const q = useSystemAudit()
  return (
    <>
      <PageHeader title={t('auditTitle')} subtitle={t('auditSubtitle')} />
      {q.isPending ? <Loading /> : q.isError ? <ErrorState error={q.error} onRetry={() => void q.refetch()} /> : (
        <>
          <AuditTable caption={t('auditTitle')} events={q.data.pages.flatMap((p) => p.items)} />
          {q.hasNextPage && <Button variant="secondary" className="mt-4" busy={q.isFetchingNextPage} onClick={() => void q.fetchNextPage()}>{t('common:loadMore')}</Button>}
        </>
      )}
    </>
  )
}
