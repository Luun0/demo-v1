import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap, type TeamMember } from '@/shared/api'
import { localInputToISO, toLocalInput } from '@/shared/lib'
import { Button, Dialog, Input, Select, useNotify } from '@/shared/ui'
import { sessionKeys } from '@/entities/session'

const schema = z
  .object({ title: z.string().trim().min(1).max(200), specialist_id: z.uuid(), starts_at: z.string().min(1), ends_at: z.string().min(1), mode: z.enum(['online', 'in_person']) })
  .refine((v) => new Date(v.ends_at) > new Date(v.starts_at), { path: ['ends_at'] })
type Values = z.infer<typeof schema>

export function CreateSessionDialog({ childId, specialists }: { childId: string; specialists: TeamMember[] }) {
  const { t } = useTranslation('schedule')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const start = new Date(Date.now() + 86400000)
  start.setMinutes(0, 0, 0)
  const { register, handleSubmit, formState, reset } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { mode: 'online', specialist_id: specialists[0]?.user_id ?? '', starts_at: toLocalInput(start), ends_at: toLocalInput(new Date(start.getTime() + 3600000)) },
  })
  const m = useMutation({
    mutationFn: async (v: Values) =>
      unwrap(await api.POST('/children/{childId}/sessions', { params: { path: { childId } }, body: { ...v, starts_at: localInputToISO(v.starts_at), ends_at: localInputToISO(v.ends_at), status: 'scheduled' } })),
    onSuccess: () => {
      notify.ok('schedule:created')
      void qc.invalidateQueries({ queryKey: sessionKeys.list(childId) })
      reset()
      setOpen(false)
    },
    onError: notify.fail,
  })
  if (specialists.length === 0) return <p className="text-ink-soft">{t('noSpecialist')}</p>
  const fe = formState.errors
  return (
    <>
      <Button onClick={() => setOpen(true)}>{t('new')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('new')}>
        <form onSubmit={handleSubmit((v) => m.mutate(v))} noValidate>
          <Input label={t('sessionTitle')} error={fe.title && t('common:required')} {...register('title')} />
          <Select label={t('specialist')} {...register('specialist_id')}>
            {specialists.map((s) => <option key={s.user_id} value={s.user_id}>{s.first_name} {s.last_name}</option>)}
          </Select>
          <Input label={t('startsAt')} type="datetime-local" error={fe.starts_at && t('common:required')} {...register('starts_at')} />
          <Input label={t('endsAt')} type="datetime-local" error={fe.ends_at && t('auth:errors.invalid')} {...register('ends_at')} />
          <Select label={t('mode')} {...register('mode')}>
            <option value="online">{t('modes.online')}</option>
            <option value="in_person">{t('modes.in_person')}</option>
          </Select>
          <Button type="submit" busy={m.isPending}>{t('common:create')}</Button>
        </form>
      </Dialog>
    </>
  )
}
