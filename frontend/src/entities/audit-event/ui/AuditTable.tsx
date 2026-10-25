import { useTranslation } from 'react-i18next'
import type { AuditEvent } from '@/shared/api'
import { formatDateTime } from '@/shared/lib'
import { Badge } from '@/shared/ui'

const tone = { success: 'brand', denied: 'danger', failure: 'warn' } as const

export function AuditTable({ events, caption }: { events: AuditEvent[]; caption: string }) {
  const { t, i18n } = useTranslation('access')
  return (
    <div className="overflow-x-auto rounded-2xl border border-line bg-surface shadow-soft">
      <table className="w-full text-left">
        <caption className="sr-only">{caption}</caption>
        <thead className="bg-muted text-sm text-ink">
          <tr>
            <th scope="col" className="p-3">{t('when')}</th>
            <th scope="col" className="p-3">{t('action')}</th>
            <th scope="col" className="p-3">{t('outcome')}</th>
          </tr>
        </thead>
        <tbody>
          {events.map((e) => (
            <tr key={e.id} className="border-t border-line">
              <td className="p-3"><time dateTime={e.occurred_at}>{formatDateTime(e.occurred_at, i18n.language)}</time></td>
              <td className="p-3 font-mono text-sm">{e.action}</td>
              <td className="p-3"><Badge tone={tone[e.outcome]}>{t(`outcomes.${e.outcome}`)}</Badge></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
