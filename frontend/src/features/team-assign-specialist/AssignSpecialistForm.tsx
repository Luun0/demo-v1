import { useId, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap } from '@/shared/api'
import { Button, useNotify } from '@/shared/ui'
import { teamKeys, useSpecialistDirectory } from '@/entities/team'

export function AssignSpecialistForm({ childId, exclude }: { childId: string; exclude: string[] }) {
  const { t } = useTranslation('team')
  const id = useId()
  const dir = useSpecialistDirectory(true)
  const [value, setValue] = useState('')
  const qc = useQueryClient()
  const notify = useNotify()
  const m = useMutation({
    mutationFn: async () => unwrap(await api.POST('/children/{childId}/specialist-assignments', { params: { path: { childId } }, body: { specialist_id: value } })),
    onSuccess: () => {
      notify.ok('team:assignedOk')
      setValue('')
      void qc.invalidateQueries({ queryKey: teamKeys.team(childId) })
    },
    onError: notify.fail,
  })
  const options = (dir.data ?? []).filter((s) => !exclude.includes(s.id))
  return (
    <form className="flex flex-wrap items-end gap-3" onSubmit={(e) => { e.preventDefault(); if (value) m.mutate() }}>
      <div>
        <label htmlFor={id} className="block text-sm font-medium">{t('assign')}</label>
        <select id={id} value={value} onChange={(e) => setValue(e.target.value)} className="mt-1 min-h-11 rounded-lg border border-line bg-surface px-3">
          <option value="">{t('chooseSpecialist')}</option>
          {options.map((s) => <option key={s.id} value={s.id}>{s.first_name} {s.last_name}</option>)}
        </select>
      </div>
      <Button type="submit" busy={m.isPending} disabled={!value}>{t('assign')}</Button>
    </form>
  )
}
