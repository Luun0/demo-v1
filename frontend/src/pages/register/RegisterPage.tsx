import { Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { RegisterForm } from '@/features/auth-register'
import { AuthLayout } from '@/widgets/auth-layout'

export function RegisterPage() {
  const { t } = useTranslation('auth')
  return (
    <AuthLayout title={t('registerTitle')}>
      <RegisterForm />
      <p className="mt-6 text-center">
        {t('haveAccount')} <Link to="/login" className="font-medium text-brand-strong underline">{t('submitLogin')}</Link>
      </p>
    </AuthLayout>
  )
}
