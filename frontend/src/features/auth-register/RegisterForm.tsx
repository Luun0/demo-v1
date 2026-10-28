import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { api, ApiError } from '@/shared/api'
import { errorKey } from '@/shared/lib'
import { Button, Input, Select } from '@/shared/ui'

// Mirrors the contract: min 10 chars, role parent|specialist (never admin).
const schema = z.object({
  email: z.email(),
  password: z.string().min(10, 'too_short').max(256, 'too_long'),
  first_name: z.string().trim().min(1).max(100),
  last_name: z.string().trim().min(1).max(100),
  role: z.enum(['parent', 'specialist']),
})
type Values = z.infer<typeof schema>

export function RegisterForm() {
  const { t } = useTranslation('auth')
  const [done, setDone] = useState<Values['role'] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const { register, handleSubmit, watch, formState, setError: setFieldError } = useForm<Values>({ resolver: zodResolver(schema), defaultValues: { role: 'parent' } })
  const fe = formState.errors
  const role = watch('role')

  const onSubmit = handleSubmit(async (v) => {
    setError(null)
    const r = await api.POST('/auth/register', { body: v })
    if (r.response.ok) return setDone(v.role)
    const err = new ApiError(r.response.status, r.error as never)
    for (const [field, code] of Object.entries(err.fieldErrors())) setFieldError(field as keyof Values, { message: code })
    setError(t(errorKey(err)))
  })

  if (done)
    return (
      <div role="status" className="rounded-lg bg-brand/10 p-4">
        <p>{t('registered')}</p>
        {done === 'specialist' && <p className="mt-2">{t('specialistPending')}</p>}
        <Link to="/login" className="mt-3 inline-block font-medium text-brand-ink underline">{t('submitLogin')}</Link>
      </div>
    )

  const msg = (m?: string) => (m ? t(`errors.${['too_short', 'too_long', 'too_common'].includes(m) ? m : 'invalid'}`) : undefined)
  return (
    <form onSubmit={onSubmit} noValidate>
      {error && <p role="alert" className="mb-4 rounded-lg bg-danger/10 p-3 text-danger">{error}</p>}
      <Input label={t('firstName')} autoComplete="given-name" error={msg(fe.first_name?.message)} {...register('first_name')} />
      <Input label={t('lastName')} autoComplete="family-name" error={msg(fe.last_name?.message)} {...register('last_name')} />
      <Input label={t('email')} type="email" autoComplete="email" error={fe.email && t('errors.email')} {...register('email')} />
      <Input label={t('password')} type="password" autoComplete="new-password" hint={t('errors.too_short')} error={msg(fe.password?.message)} {...register('password')} />
      <Select label={t('role')} {...register('role')}>
        <option value="parent">{t('common:roles.parent')}</option>
        <option value="specialist">{t('common:roles.specialist')}</option>
      </Select>
      {role === 'specialist' && <p className="mb-4 -mt-2 rounded-lg bg-muted p-3 text-sm text-ink-soft">{t('registrationRoleHelp')}</p>}
      <Button type="submit" busy={formState.isSubmitting} className="w-full">{t('submitRegister')}</Button>
    </form>
  )
}
