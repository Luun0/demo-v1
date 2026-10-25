import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { api, readCookie, setAccessToken } from '@/shared/api'
import { Button, Icon } from '@/shared/ui'

export function LogoutButton({ className }: { className?: string }) {
  const { t } = useTranslation('nav')
  const nav = useNavigate()
  const qc = useQueryClient()
  const logout = async () => {
    try {
      await api.POST('/auth/logout', { params: { header: { 'X-CSRF-Token': readCookie('ssp_csrf') ?? '' } } })
    } finally {
      setAccessToken(null)
      qc.clear()
      void nav('/', { replace: true })
    }
  }
  return (
    <Button variant="ghost" className={className} onClick={() => void logout()}>
      <Icon name="exit" />
      {t('logout')}
    </Button>
  )
}
