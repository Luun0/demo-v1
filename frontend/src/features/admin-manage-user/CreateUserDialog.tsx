import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api, unwrap } from '@/shared/api'
import { Button, Dialog, Input, Select, useNotify } from '@/shared/ui'
import { userKeys } from '@/entities/user'

const schema = z.object({
  email: z.email(),
  password: z.string().min(10),
  first_name: z.string().trim().min(1),
  last_name: z.string().trim().min(1),
  role: z.enum(['parent', 'specialist', 'admin']),
})
type Values = z.infer<typeof schema>

export function CreateUserDialog() {
  const { t } = useTranslation('admin')
  const [open, setOpen] = useState(false)
  const qc = useQueryClient()
  const notify = useNotify()
  const { register, handleSubmit, formState, reset } = useForm<Values>({ resolver: zodResolver(schema), defaultValues: { role: 'specialist' } })
  const m = useMutation({
    mutationFn: async (v: Values) => unwrap(await api.POST('/users', { body: { ...v, specialist_approved: v.role === 'specialist' } })),
    onSuccess: () => {
      notify.ok('admin:created')
      void qc.invalidateQueries({ queryKey: userKeys.all })
      reset()
      setOpen(false)
    },
    onError: notify.fail,
  })
  const fe = formState.errors
  return (
    <>
      <Button onClick={() => setOpen(true)}>{t('newUser')}</Button>
      <Dialog open={open} onOpenChange={setOpen} title={t('newUser')}>
        <form onSubmit={handleSubmit((v) => m.mutate(v))} noValidate>
          <Input label={t('auth:firstName')} error={fe.first_name && t('common:required')} {...register('first_name')} />
          <Input label={t('auth:lastName')} error={fe.last_name && t('common:required')} {...register('last_name')} />
          <Input label={t('auth:email')} type="email" autoComplete="off" error={fe.email && t('auth:errors.email')} {...register('email')} />
          <Input label={t('auth:password')} type="password" autoComplete="new-password" error={fe.password && t('auth:errors.too_short')} {...register('password')} />
          <Select label={t('role')} {...register('role')}>
            {(['parent', 'specialist', 'admin'] as const).map((r) => <option key={r} value={r}>{t(`common:roles.${r}`)}</option>)}
          </Select>
          <Button type="submit" busy={m.isPending}>{t('common:create')}</Button>
        </form>
      </Dialog>
    </>
  )
}
