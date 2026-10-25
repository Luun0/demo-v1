export interface Problem {
  type?: string
  title?: string
  status: number
  request_id?: string
  errors?: { field: string; code: string }[]
}

export class ApiError extends Error {
  readonly status: number
  readonly problem: Problem | undefined

  constructor(status: number, problem?: Problem) {
    super(`HTTP ${status}`)
    this.status = status
    this.problem = problem
  }

  fieldErrors(): Record<string, string> {
    const out: Record<string, string> = {}
    for (const e of this.problem?.errors ?? []) out[e.field] = e.code
    return out
  }
}

export function isApiError(e: unknown, status?: number): e is ApiError {
  return e instanceof ApiError && (status === undefined || e.status === status)
}
