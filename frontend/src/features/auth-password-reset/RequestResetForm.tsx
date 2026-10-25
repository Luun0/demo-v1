import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { api, ApiError } from '@/shared/api'
import { errorKey } from '@/shared/lib'
import { Button, Input } from '@/shared/ui'

const schema = z.object({ email: z.email() })
type Values = z.infer<typeof schema>

export function RequestResetForm() {
  const { t } = useTranslation('auth')
  const [sent, setSent] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { register, handleSubmit, formState } = useForm<Values>({ resolver: zodResolver(schema) })
  const submit = handleSubmit(async (body) => {
    setError(null)
    try {
      const r = await api.POST('/auth/password-reset/request', { body })
      if (!r.response.ok) setError(t(errorKey(new ApiError(r.response.status, r.error as never))))
      else setSent(true)
    } catch (e) {
      setError(t(errorKey(e)))
    }
  })
  if (sent) return <div role="status" className="space-y-3 rounded-xl bg-brand/10 p-4"><p>{t('resetRequested')}</p><p className="text-sm">{t('resetLogHint')}</p></div>
  return (
    <form onSubmit={submit} noValidate>
      {error && <p role="alert" className="mb-4 rounded-lg bg-danger/10 p-3 text-danger">{error}</p>}
      <Input label={t('email')} type="email" autoComplete="email" error={formState.errors.email && t('errors.email')} {...register('email')} />
      <Button type="submit" busy={formState.isSubmitting} className="w-full">{t('sendReset')}</Button>
    </form>
  )
}
