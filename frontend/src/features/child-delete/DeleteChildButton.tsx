import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, ApiError } from '@/shared/api'
import { Button, Dialog, useNotify } from '@/shared/ui'
import { childKeys } from '@/entities/child'

export function DeleteChildButton({ childId }: { childId: string }) {
  const { t } = useTranslation('child')
  const [open, setOpen] = useState(false)
  const nav = useNavigate()
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async () => {
      const r = await api.DELETE('/children/{childId}', { params: { path: { childId } } })
      if (!r.response.ok) throw new ApiError(r.response.status)
    },
    onSuccess: () => {
      notify.ok('child:deleted')
      qc.removeQueries({ queryKey: childKeys.one(childId) })
      void qc.invalidateQueries({ queryKey: childKeys.all })
      void nav('/app', { replace: true })
    },
    onError: notify.fail,
  })
  return (
    <>
      <Button variant="danger" onClick={() => setOpen(true)}>{t('common:delete')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('deleteTitle')} description={t('deleteText')}>
        <div className="flex justify-end gap-3">
          <Button variant="secondary" onClick={() => setOpen(false)}>{t('common:cancel')}</Button>
          <Button variant="danger" busy={m.isPending} onClick={() => m.mutate()}>{t('common:confirm')}</Button>
        </div>
      </Dialog>
    </>
  )
}
