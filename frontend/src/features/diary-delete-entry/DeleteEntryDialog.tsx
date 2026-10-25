import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap, type DiaryEntry } from '@/shared/api'
import { Button, Dialog, useNotify } from '@/shared/ui'
import { diaryKeys } from '@/entities/diary-entry'

export function DeleteEntryDialog({ entry }: { entry: DiaryEntry }) {
  const { t } = useTranslation('diary')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const mutation = useMutation({
    mutationFn: async () => unwrap(await api.DELETE('/children/{childId}/diary-entries/{entryId}', {
      params: { path: { childId: entry.child_id, entryId: entry.id } },
    })),
    onSuccess: () => {
      notify.ok('diary:deleted')
      void qc.invalidateQueries({ queryKey: diaryKeys.all(entry.child_id) })
      setOpen(false)
    },
    onError: notify.fail,
  })
  return (
    <>
      <Button variant="danger" onClick={() => setOpen(true)}>{t('deleteEntry')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('deleteTitle')} description={t('deleteText')}>
        <div className="flex justify-end gap-2">
          <Button variant="secondary" onClick={() => setOpen(false)}>{t('common:cancel')}</Button>
          <Button variant="danger" busy={mutation.isPending} onClick={() => mutation.mutate()}>{t('common:delete')}</Button>
        </div>
      </Dialog>
    </>
  )
}
