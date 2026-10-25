import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { api, ApiError, setAccessToken } from '@/shared/api'
import { errorKey } from '@/shared/lib'
import { Button, Input } from '@/shared/ui'

const schema = z.object({ email: z.email(), password: z.string().min(1) })
type Values = z.infer<typeof schema>

export function LoginForm() {
  const { t } = useTranslation('auth')
  const nav = useNavigate()
  const qc = useQueryClient()
  const [error, setError] = useState<string | null>(null)
  const { register, handleSubmit, formState } = useForm<Values>({ resolver: zodResolver(schema) })

  const onSubmit = handleSubmit(async (v) => {
    setError(null)
    try {
      const r = await api.POST('/auth/login', { body: v })
      if (!r.response.ok || !r.data) {
        setError(r.response.status === 401 ? t('invalidCredentials') : t(errorKey(new ApiError(r.response.status))))
        return
      }
      setAccessToken(r.data.access_token)
      qc.clear()
      void nav('/app', { replace: true })
    } catch (e) {
      setError(t(errorKey(e)))
    }
  })

  return (
    <form onSubmit={onSubmit} noValidate aria-describedby={error ? 'login-error' : undefined}>
      {error && (
        <p id="login-error" role="alert" className="mb-4 rounded-lg bg-danger/10 p-3 text-danger">
          {error}
        </p>
      )}
      <Input label={t('email')} type="email" autoComplete="username" error={formState.errors.email && t('errors.email')} {...register('email')} />
      <Input label={t('password')} type="password" autoComplete="current-password" error={formState.errors.password && t('common:required')} {...register('password')} />
      <Button type="submit" busy={formState.isSubmitting} className="w-full">
        {t('submitLogin')}
      </Button>
    </form>
  )
}
