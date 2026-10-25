import { describe, expect, it } from 'vitest'
import { childNav, mainNav } from './nav'

describe('navigation per role mirrors the access matrix', () => {
  it('admin sees no diary/schedule; specialist sees no access log', () => {
    const paths = (r: 'admin' | 'specialist' | 'parent') => childNav('c', r).map((i) => i.to.split('/').pop())
    expect(paths('admin')).not.toContain('diary')
    expect(paths('specialist')).not.toContain('access-log')
    expect(paths('parent')).toEqual(['c', 'schedule', 'diary', 'team', 'access-log'])
    expect(mainNav('parent').map((i) => i.to)).toEqual(['/app'])
    expect(mainNav('admin')).toHaveLength(3)
  })
})
