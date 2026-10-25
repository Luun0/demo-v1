import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap } from '@/shared/api'
import { Button, Dialog, useNotify } from '@/shared/ui'
import { teamKeys } from '@/entities/team'

export function EndAssignmentButton({ childId, specialistId, name }: { childId: string; specialistId: string; name: string }) {
  const { t } = useTranslation('team')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async () =>
      unwrap(await api.PATCH('/children/{childId}/specialist-assignments/{specialistId}', { params: { path: { childId, specialistId } }, body: { ended_at: new Date().toISOString().replace(/\.\d{3}Z$/, 'Z') } })),
    onSuccess: () => {
      notify.ok('team:endedOk')
      setOpen(false)
      void qc.invalidateQueries({ queryKey: teamKeys.team(childId) })
    },
    onError: notify.fail,
  })
  return (
    <>
      <Button variant="secondary" onClick={() => setOpen(true)}>{t('end')}<span className="sr-only">: {name}</span></Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('end')} description={t('endConfirm')}>
        <div className="flex justify-end gap-3">
          <Button variant="secondary" onClick={() => setOpen(false)}>{t('common:cancel')}</Button>
          <Button variant="danger" busy={m.isPending} onClick={() => m.mutate()}>{t('common:confirm')}</Button>
        </div>
      </Dialog>
    </>
  )
}
