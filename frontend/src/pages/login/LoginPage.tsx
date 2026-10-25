import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { LoginForm } from '@/features/auth-login'
import { AuthLayout } from '@/widgets/auth-layout'

export function LoginPage() {
  const { t } = useTranslation('auth')
  return (
    <AuthLayout title={t('loginTitle')}>
      <LoginForm />
      <p className="mt-6 text-center">
        {t('noAccount')} <Link to="/register" className="font-medium text-brand-strong underline">{t('submitRegister')}</Link>
      </p>
    </AuthLayout>
  )
}
