import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap, type Child } from '@/shared/api'
import { Button, Dialog, Input, useNotify } from '@/shared/ui'
import { childKeys } from '@/entities/child'

const schema = z.object({ first_name: z.string().trim().min(1).max(100), last_name: z.string().trim().min(1).max(100), birth_date: z.iso.date() })
type Values = z.infer<typeof schema>

export function EditChildDialog({ child }: { child: Child }) {
  const { t } = useTranslation('child')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const { register, handleSubmit, formState } = useForm<Values>({ resolver: zodResolver(schema), values: { first_name: child.first_name, last_name: child.last_name, birth_date: child.birth_date } })
  const m = useMutation({
    mutationFn: async (v: Values) => unwrap(await api.PATCH('/children/{childId}', { params: { path: { childId: child.id } }, body: v })),
    onSuccess: (c) => {
      qc.setQueryData(childKeys.one(child.id), c)
      void qc.invalidateQueries({ queryKey: childKeys.all })
      notify.ok('common:saved')
      setOpen(false)
    },
    onError: notify.fail,
  })
  return (
    <>
      <Button variant="secondary" onClick={() => setOpen(true)}>{t('editProfile')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('editProfile')}>
        <form onSubmit={handleSubmit((v) => m.mutate(v))} noValidate>
          <Input label={t('firstName')} error={formState.errors.first_name && t('common:required')} {...register('first_name')} />
          <Input label={t('lastName')} error={formState.errors.last_name && t('common:required')} {...register('last_name')} />
          <Input label={t('birthDate')} type="date" error={formState.errors.birth_date && t('common:required')} {...register('birth_date')} />
          <Button type="submit" busy={m.isPending}>{t('common:save')}</Button>
        </form>
      </Dialog>
    </>
  )
}
