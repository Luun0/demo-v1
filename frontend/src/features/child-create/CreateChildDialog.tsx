import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap } from '@/shared/api'
import { Button, Dialog, Input, useNotify } from '@/shared/ui'
import { childKeys } from '@/entities/child'

export const POLICY_VERSION = '2026-10'

const schema = z.object({
  first_name: z.string().trim().min(1).max(100),
  last_name: z.string().trim().min(1).max(100),
  birth_date: z.iso.date(),
  consent: z.literal(true),
})
type Values = z.infer<typeof schema>

export function CreateChildDialog() {
  const { t } = useTranslation('child')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const { register, handleSubmit, formState, reset } = useForm<Values>({ resolver: zodResolver(schema) })
  const m = useMutation({
    mutationFn: async (v: Values) =>
      unwrap(await api.POST('/children', { body: { first_name: v.first_name, last_name: v.last_name, birth_date: v.birth_date, policy_version: POLICY_VERSION } })),
    onSuccess: () => {
      notify.ok('child:created')
      void qc.invalidateQueries({ queryKey: childKeys.all })
      reset()
      setOpen(false)
    },
    onError: notify.fail,
  })
  const fe = formState.errors
  return (
    <>
      <Button onClick={() => setOpen(true)}>{t('dashboard:addChild')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('dashboard:addChild')}>
        <form onSubmit={handleSubmit((v) => m.mutate(v))} noValidate>
          <Input label={t('firstName')} error={fe.first_name && t('common:required')} {...register('first_name')} />
          <Input label={t('lastName')} error={fe.last_name && t('common:required')} {...register('last_name')} />
          <Input label={t('birthDate')} type="date" error={fe.birth_date && t('common:required')} {...register('birth_date')} />
          <div className="mb-4 flex items-start gap-3">
            <input id="consent" type="checkbox" className="mt-1 h-6 w-6" aria-invalid={!!fe.consent} aria-describedby={fe.consent ? 'consent-error' : undefined} {...register('consent')} />
            <label htmlFor="consent">{t('consent', { version: POLICY_VERSION })}</label>
          </div>
          {fe.consent && <p id="consent-error" role="alert" className="mb-4 text-danger">{t('consentRequired')}</p>}
          <Button type="submit" busy={m.isPending}>{t('common:create')}</Button>
        </form>
      </Dialog>
    </>
  )
}
