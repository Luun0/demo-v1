import type { Role } from '@/shared/api'

/** UI hints only; the server is the authority (policy.Decide). */
export type Permission = 'child:create' | 'diary:choose-visibility' | 'audit:child' | 'admin' | 'specialist:directory'

const grants: Record<Role, Permission[]> = {
  parent: ['child:create', 'audit:child', 'specialist:directory'],
  specialist: ['diary:choose-visibility'],
  admin: ['admin', 'audit:child', 'specialist:directory'],
}

export function can(role: Role | undefined, p: Permission): boolean {
  return !!role && grants[role].includes(p)
}
