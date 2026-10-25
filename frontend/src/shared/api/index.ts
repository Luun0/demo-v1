export { api, unwrap, refreshSession, authFetch, readCookie, API_BASE } from './client'
export { ApiError, isApiError, type Problem } from './errors'
export { getAccessToken, setAccessToken, subscribeToken } from './token'
export type { components, paths } from './schema.gen'

import type { components } from './schema.gen'
type S = components['schemas']
export type User = S['User']
export type Child = S['Child']
export type SupportSession = S['SupportSession']
export type DiaryEntry = S['DiaryEntry']
export type DiaryEntryVersion = S['DiaryEntryVersion']
export type AuditEvent = S['AuditEventView']
export type TeamMember = S['TeamMember']
export type SpecialistSummary = S['SpecialistSummary']
export type Role = User['role']
export type Visibility = DiaryEntry['visibility']
