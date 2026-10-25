import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap } from '@/shared/api'
import { Button, Dialog, useNotify } from '@/shared/ui'
import { diaryKeys, EntryForm, toEntryBody, type EntryValues } from '@/entities/diary-entry'

export function CreateEntryDialog({ childId, canChooseVisibility }: { childId: string; canChooseVisibility: boolean }) {
  const { t } = useTranslation('diary')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async (v: EntryValues) =>
      unwrap(await api.POST('/children/{childId}/diary-entries', { params: { path: { childId } }, body: toEntryBody(canChooseVisibility ? v : { ...v, visibility: 'shared' }) })),
    onSuccess: () => {
      notify.ok('diary:created')
      void qc.invalidateQueries({ queryKey: diaryKeys.all(childId) })
      setOpen(false)
    },
    onError: notify.fail,
  })
  return (
    <>
      <Button onClick={() => setOpen(true)}>{t('new')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('new')}>
        <EntryForm canChooseVisibility={canChooseVisibility} busy={m.isPending} onSubmit={(v) => m.mutate(v)} />
      </Dialog>
    </>
  )
}
