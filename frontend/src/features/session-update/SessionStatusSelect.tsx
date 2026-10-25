import { useId } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap, type SupportSession } from '@/shared/api'
import { useNotify } from '@/shared/ui'
import { sessionKeys } from '@/entities/session'

const statuses = ['scheduled', 'completed', 'cancelled'] as const

export function SessionStatusSelect({ session }: { session: SupportSession }) {
  const { t } = useTranslation('schedule')
  const id = useId()
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async (status: SupportSession['status']) =>
      unwrap(await api.PATCH('/children/{childId}/sessions/{sessionId}', { params: { path: { childId: session.child_id, sessionId: session.id } }, body: { status } })),
    onSuccess: () => {
      notify.ok('schedule:updated')
      void qc.invalidateQueries({ queryKey: sessionKeys.list(session.child_id) })
    },
    onError: notify.fail,
  })
  return (
    <>
      <label htmlFor={id} className="sr-only">{t('statusLabel')}: {session.title}</label>
      <select id={id} value={session.status} disabled={m.isPending} onChange={(e) => m.mutate(e.target.value as SupportSession['status'])} className="min-h-11 rounded-lg border border-line bg-surface px-2">
        {statuses.map((s) => <option key={s} value={s}>{t(`statuses.${s}`)}</option>)}
      </select>
    </>
  )
}
