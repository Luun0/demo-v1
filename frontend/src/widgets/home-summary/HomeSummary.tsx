import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { Child, Role } from '@/shared/api'
import { EmptyState, ErrorState, Loading } from '@/shared/ui'
import { ChildCard, useChildren } from '@/entities/child'
import { useDiary, useDiaryCount } from '@/entities/diary-entry'
import { childNav } from '@/shared/config/nav'
import { formatDateTime } from '@/shared/lib'
import { Link } from 'react-router'

/** Count comes from the server (X-Total-Count), already visibility-filtered (fixes C5). */
function DiaryCount({ child }: { child: Child }) {
  const { t } = useTranslation('dashboard')
  const q = useDiaryCount(child.id)
  if (q.data === undefined) return null
  return <p className="count-pop mt-2 text-sm text-ink">{t('visibleEntries', { count: q.data })}</p>
}

function HomeRecommendations({ child }: { child: Child }) {
  const { t, i18n } = useTranslation(['dashboard', 'diary', 'nav'])
  const diary = useDiary(child.id, {})
  const [visibleCount, setVisibleCount] = useState(2)
  const recommendations = diary.data?.pages.flatMap((page) => page.items)
    .filter((entry) => entry.type === 'recommendation' && entry.visibility === 'shared')
    ?? []
  const visibleEntryCount = diary.data?.pages[0]?.total
  const recent = recommendations.slice(0, visibleCount)
  const diaryPath = childNav(child.id, 'parent').find((item) => item.labelKey === 'nav:diary')?.to
  if (!diaryPath) return null
  if (diary.isPending) return <p className="mt-4 border-t border-line pt-3 text-sm text-ink-soft">{t('recommendationsLoading')}</p>
  if (diary.isError) return <section className="mt-4 border-t border-line pt-3 text-sm text-ink-soft">
    <p>{t('recommendationsError')}</p>
    <Link to={diaryPath} className="mt-2 inline-flex min-h-10 items-center font-medium text-brand-ink underline underline-offset-4">{t('recommendationsOpen')}</Link>
  </section>

  async function showOlder() {
    const target = visibleCount + 2
    let loaded = recommendations.length
    let canFetchNext = diary.hasNextPage
    while (loaded < target && canFetchNext) {
      const result = await diary.fetchNextPage()
      loaded = result.data?.pages.flatMap((page) => page.items)
        .filter((entry) => entry.type === 'recommendation' && entry.visibility === 'shared').length ?? loaded
      canFetchNext = result.hasNextPage && !result.isError
    }
    setVisibleCount(target)
  }

  return (
    <>
      {visibleEntryCount !== undefined && <p className="count-pop mt-2 text-sm text-ink">{t('visibleEntries', { count: visibleEntryCount })}</p>}
      <section className="home-recommendations mt-4 border-t border-line pt-3" aria-label={t('recommendationsTitle')}>
        <h4 className="font-semibold text-brand-ink">{t('recommendationsTitle')}</h4>
        {recent.length ? (
          <ul className="mt-2 space-y-3">
            {recent.map((entry) => <li key={entry.id} className="border-l-2 border-accent/70 pl-3">
              <time className="text-xs text-ink-soft" dateTime={entry.observed_at}>{formatDateTime(entry.observed_at, i18n.language)}</time>
              <p className="mt-0.5 line-clamp-3 whitespace-pre-wrap break-words text-sm text-ink">{entry.content}</p>
            </li>)}
          </ul>
        ) : <p className="mt-1 text-sm text-ink-soft">{t('recommendationsEmpty')}</p>}
        {(recent.length < recommendations.length || diary.hasNextPage) && <button type="button" onClick={() => void showOlder()} disabled={diary.isFetchingNextPage} className="mt-3 min-h-10 font-medium text-brand-ink underline underline-offset-4 disabled:opacity-60">
          {diary.isFetchingNextPage ? t('recommendationsLoadingMore') : t('recommendationsMore')}
        </button>}
        <Link to={diaryPath} className="mt-3 inline-flex min-h-10 items-center font-medium text-brand-ink underline underline-offset-4">{t('recommendationsOpen')}</Link>
      </section>
    </>
  )
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
          <ChildCard child={c} role={role} extra={<>
            {role === 'specialist' && <DiaryCount child={c} />}
            {role === 'parent' && <HomeRecommendations child={c} />}
          </>} />
        </li>
      ))}
    </ul>
  )
}
