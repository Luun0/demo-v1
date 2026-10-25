import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import type { SupportSession } from '@/shared/api'
import { formatDateTime } from '@/shared/lib'
import { Badge } from '@/shared/ui'

const tone = { scheduled: 'brand', completed: 'neutral', cancelled: 'warn' } as const

export function SessionRow({ session, actions }: { session: SupportSession; actions?: ReactNode }) {
  const { t, i18n } = useTranslation('schedule')
  return (
    <li className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-line bg-surface p-4">
      <div>
        <p className="font-semibold text-ink">{session.title}</p>
        <p className="text-ink-soft">
          <time dateTime={session.starts_at}>{formatDateTime(session.starts_at, i18n.language)}</time> · {t(`modes.${session.mode}`)}
        </p>
      </div>
      <div className="flex items-center gap-3">
        <Badge tone={tone[session.status]}>{t(`statuses.${session.status}`)}</Badge>
        {actions}
      </div>
    </li>
  )
}
