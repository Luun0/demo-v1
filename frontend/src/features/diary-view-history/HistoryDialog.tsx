import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { DiaryEntry } from '@/shared/api'
import { formatDateTime } from '@/shared/lib'
import { Button, Dialog, ErrorState, Loading } from '@/shared/ui'
import { useVersions } from '@/entities/diary-entry'

export function HistoryDialog({ entry }: { entry: DiaryEntry }) {
  const { t, i18n } = useTranslation('diary')
  const [open, setOpen] = useState(false)
  const q = useVersions(entry.child_id, entry.id, open)
  return (
    <>
      <Button variant="ghost" onClick={() => setOpen(true)}>{t('history')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('history')}>
        {q.isPending ? <Loading /> : q.isError ? <ErrorState error={q.error} onRetry={() => void q.refetch()} /> : (
          <ol className="space-y-3">
            {q.data.map((v) => (
              <li key={v.id} className="rounded-lg border border-line p-3">
                <p className="text-sm text-ink-soft">
                  {t('version', { n: v.version_number })} · <time dateTime={v.changed_at}>{formatDateTime(v.changed_at, i18n.language)}</time> · {t(`visibilities.${v.visibility}`)}
                </p>
                <p className="whitespace-pre-wrap break-words">{v.content}</p>
              </li>
            ))}
          </ol>
        )}
      </Dialog>
    </>
  )
}
