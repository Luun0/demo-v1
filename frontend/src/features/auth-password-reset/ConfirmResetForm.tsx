import { useEffect, useState } from 'react'
import { Link } from 'react-router'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { api, ApiError } from '@/shared/api'
import { errorKey } from '@/shared/lib'
import { Button, Input } from '@/shared/ui'

const schema = z.object({ password: z.string().min(10, 'too_short').max(256, 'too_long'), confirm: z.string() })
  .refine((v) => v.password === v.confirm, { path: ['confirm'], message: 'mismatch' })
type Values = z.infer<typeof schema>

export function ConfirmResetForm() {
  const { t } = useTranslation('auth')
  const [token] = useState(() => new URLSearchParams(window.location.search).get('token') ?? '')
  const [done, setDone] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { register, handleSubmit, formState } = useForm<Values>({ resolver: zodResolver(schema) })
  useEffect(() => {
    if (token) window.history.replaceState({}, '', window.location.pathname)
  }, [token])
  const submit = handleSubmit(async ({ password }) => {
    setError(null)
    try {
      const r = await api.POST('/auth/password-reset/confirm', { body: { token, password } })
      if (!r.response.ok) setError(r.response.status === 401 ? t('resetInvalid') : t(errorKey(new ApiError(r.response.status, r.error as never))))
      else setDone(true)
    } catch (e) {
      setError(t(errorKey(e)))
    }
  })
  if (!token) return <p role="alert" className="rounded-lg bg-danger/10 p-3 text-danger">{t('resetInvalid')}</p>
  if (done) return <div role="status" className="space-y-3 rounded-xl bg-brand/10 p-4"><p>{t('resetDone')}</p><Link to="/login" className="font-medium text-brand-strong underline">{t('backToLogin')}</Link></div>
  const passwordError = formState.errors.password?.message
  const confirmError = formState.errors.confirm?.message
  return (
    <form onSubmit={submit} noValidate>
      {error && <p role="alert" className="mb-4 rounded-lg bg-danger/10 p-3 text-danger">{error}</p>}
      <Input label={t('password')} type="password" autoComplete="new-password" error={passwordError && t(`errors.${passwordError}`)} {...register('password')} />
      <Input label={t('confirmPassword')} type="password" autoComplete="new-password" error={confirmError && t(confirmError === 'mismatch' ? 'passwordMismatch' : `errors.${confirmError}`)} {...register('confirm')} />
      <Button type="submit" busy={formState.isSubmitting} className="w-full">{t('confirmReset')}</Button>
    </form>
  )
}
