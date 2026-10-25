import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { ConfirmResetForm } from '@/features/auth-password-reset'
import { AuthLayout } from '@/widgets/auth-layout'

export function ResetPasswordPage() {
  const { t } = useTranslation('auth')
  return <AuthLayout title={t('resetTitle')}>
    <p className="mb-5 text-sm text-ink-soft">{t('resetText')}</p>
    <ConfirmResetForm />
    <p className="mt-6 text-center"><Link to="/login" className="font-medium text-brand-strong underline">{t('backToLogin')}</Link></p>
  </AuthLayout>
}
