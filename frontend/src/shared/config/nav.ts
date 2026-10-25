import type { Role } from '@/shared/api'

export interface NavItem {
  to: string
  labelKey: string
  end?: boolean
}

export function mainNav(role: Role | undefined): NavItem[] {
  const items: NavItem[] = [{ to: '/app', labelKey: 'nav:dashboard', end: true }]
  if (role === 'admin') items.push({ to: '/app/admin/users', labelKey: 'nav:users' }, { to: '/app/admin/audit', labelKey: 'nav:audit' })
  return items
}

export function childNav(childId: string, role: Role | undefined): NavItem[] {
  const base = `/app/children/${childId}`
  const items: NavItem[] = [{ to: base, labelKey: 'nav:profile', end: true }]
  if (role !== 'admin') items.push({ to: `${base}/schedule`, labelKey: 'nav:schedule' }, { to: `${base}/diary`, labelKey: 'nav:diary' })
  items.push({ to: `${base}/team`, labelKey: 'nav:team' })
  if (role !== 'specialist') items.push({ to: `${base}/access-log`, labelKey: 'nav:accessLog' })
  return items
}
