import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import type { DiaryEntry } from '@/shared/api'
import { formatDateTime } from '@/shared/lib'
import { Badge } from '@/shared/ui'

/** Diary text is rendered as plain text only (no HTML injection, SEC-21). */
export function EntryCard({ entry, actions }: { entry: DiaryEntry; actions?: ReactNode }) {
  const { t, i18n } = useTranslation('diary')
  return (
    <article className="rounded-xl border border-line bg-surface p-5" aria-labelledby={`entry-${entry.id}`}>
      <header className="mb-2 flex flex-wrap items-center gap-2">
        <h3 id={`entry-${entry.id}`} className="font-semibold text-ink">
          {t(`types.${entry.type}`)}
        </h3>
        <Badge tone={entry.visibility === 'shared' ? 'brand' : 'warn'}>{t(`visibilities.${entry.visibility}`)}</Badge>
        <time className="ml-auto text-sm text-ink-soft" dateTime={entry.observed_at}>
          {formatDateTime(entry.observed_at, i18n.language)}
        </time>
      </header>
      <p className="whitespace-pre-wrap break-words text-ink">{entry.content}</p>
      {actions && <div className="mt-3 flex gap-2">{actions}</div>}
    </article>
  )
}
