import { useDeferredValue, useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { User, Visibility } from '@/shared/api'
import { Button, EmptyState, ErrorState, Loading } from '@/shared/ui'
import { EntryCard, useDiary } from '@/entities/diary-entry'
import { can } from '@/entities/viewer'
import { EditEntryDialog } from '@/features/diary-edit-entry'
import { HistoryDialog } from '@/features/diary-view-history'
import { DeleteEntryDialog } from '@/features/diary-delete-entry'

export function DiaryFeed({ childId, viewer }: { childId: string; viewer: User }) {
  const { t } = useTranslation('diary')
  const searchId = useId()
  const visId = useId()
  const [q, setQ] = useState('')
  const [vis, setVis] = useState<Visibility | ''>('')
  const deferredQ = useDeferredValue(q.trim())
  const chooseVis = can(viewer.role, 'diary:choose-visibility')
  const feed = useDiary(childId, { q: deferredQ, visibility: vis })
  const total = feed.data?.pages[0]?.total
  return (
    <section aria-label={t('title')}>
      <div className="mb-4 flex flex-wrap items-end gap-3" role="search">
        <div className="grow">
          <label htmlFor={searchId} className="block text-sm font-medium">{t('common:search')}</label>
          <input id={searchId} type="search" value={q} maxLength={100} onChange={(e) => setQ(e.target.value)} placeholder={t('searchPlaceholder')} className="mt-1 min-h-11 w-full rounded-lg border border-line bg-surface px-3" />
        </div>
        {chooseVis && (
          <div>
            <label htmlFor={visId} className="block text-sm font-medium">{t('visibility')}</label>
            <select id={visId} value={vis} onChange={(e) => setVis(e.target.value as Visibility | '')} className="mt-1 min-h-11 rounded-lg border border-line bg-surface px-3">
              <option value="">{t('allVisibilities')}</option>
              <option value="shared">{t('visibilities.shared')}</option>
              <option value="specialists_only">{t('visibilities.specialists_only')}</option>
            </select>
          </div>
        )}
      </div>
      <p className="mb-3 text-sm text-ink-soft" aria-live="polite">{total !== undefined && t('total', { count: total })}</p>
      {feed.isPending ? <Loading /> : feed.isError ? <ErrorState error={feed.error} onRetry={() => void feed.refetch()} /> : total === 0 ? (
        <EmptyState>{t('empty')}</EmptyState>
      ) : (
        <ul className="space-y-3">
          {feed.data.pages.flatMap((p) => p.items).map((e) => (
            <li key={e.id}>
              <EntryCard
                entry={e}
                actions={
                  <>
                    {e.author_id === viewer.id && <EditEntryDialog entry={e} canChooseVisibility={chooseVis} />}
                    {can(viewer.role, 'diary:delete') && e.author_id === viewer.id && <DeleteEntryDialog entry={e} />}
                    <HistoryDialog entry={e} />
                  </>
                }
              />
            </li>
          ))}
        </ul>
      )}
      {feed.hasNextPage && (
        <Button variant="secondary" className="mt-4" busy={feed.isFetchingNextPage} onClick={() => void feed.fetchNextPage()}>{t('common:loadMore')}</Button>
      )}
    </section>
  )
}
