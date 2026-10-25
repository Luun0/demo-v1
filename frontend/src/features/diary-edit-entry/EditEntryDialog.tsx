import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap, type DiaryEntry } from '@/shared/api'
import { toLocalInput } from '@/shared/lib'
import { Button, Dialog, useNotify } from '@/shared/ui'
import { diaryKeys, EntryForm, toEntryBody, type EntryValues } from '@/entities/diary-entry'

export function EditEntryDialog({ entry, canChooseVisibility }: { entry: DiaryEntry; canChooseVisibility: boolean }) {
  const { t } = useTranslation('diary')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async (v: EntryValues) =>
      unwrap(await api.PUT('/children/{childId}/diary-entries/{entryId}', { params: { path: { childId: entry.child_id, entryId: entry.id } }, body: toEntryBody(v) })),
    onSuccess: () => {
      notify.ok('diary:updated')
      void qc.invalidateQueries({ queryKey: diaryKeys.all(entry.child_id) })
      setOpen(false)
    },
    onError: notify.fail,
  })
  return (
    <>
      <Button variant="secondary" onClick={() => setOpen(true)}>{t('common:edit')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('common:edit')}>
        <EntryForm
          canChooseVisibility={canChooseVisibility}
          busy={m.isPending}
          defaults={{ type: entry.type, visibility: entry.visibility, content: entry.content, observed_at: toLocalInput(new Date(entry.observed_at)) }}
          onSubmit={(v) => m.mutate(v)}
        />
      </Dialog>
    </>
  )
}
