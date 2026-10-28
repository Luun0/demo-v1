import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { RequestResetForm } from '@/features/auth-password-reset'
import { AuthLayout } from '@/widgets/auth-layout'

export function ForgotPasswordPage() {
  const { t } = useTranslation('auth')
  return <AuthLayout title={t('forgotTitle')}>
    <p className="mb-5 text-sm text-ink-soft">{t('forgotText')}</p>
    <p role="note" className="mb-5 rounded-lg border border-warn/50 bg-warn/15 p-3 text-sm text-ink">{t('resetUnavailable')}</p>
    <RequestResetForm />
    <p className="mt-6 text-center"><Link to="/login" className="font-medium text-brand-ink underline">{t('backToLogin')}</Link></p>
  </AuthLayout>
}
