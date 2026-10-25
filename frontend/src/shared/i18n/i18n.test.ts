import { describe, expect, it } from 'vitest'
import en from './locales/en.json'
import kz from './locales/kz.json'
import ru from './locales/ru.json'

type Tree = { [k: string]: string | Tree }

function keys(t: Tree, prefix = ''): string[] {
  return Object.entries(t).flatMap(([k, v]) => (typeof v === 'string' ? [prefix + k.replace(/_(one|few|many|other)$/, '')] : keys(v, `${prefix}${k}.`)))
}

const uniq = (a: string[]) => [...new Set(a)].sort()

describe('locale completeness (kz / ru / en)', () => {
  it('every key exists in all three locales', () => {
    const base = uniq(keys(en as Tree))
    expect(uniq(keys(ru as Tree))).toEqual(base)
    expect(uniq(keys(kz as Tree))).toEqual(base)
  })

  it('no empty strings and kz/en are actually translated', () => {
    const flat = (t: Tree, p = ''): [string, string][] =>
      Object.entries(t).flatMap(([k, v]) => (typeof v === 'string' ? [[p + k, v] as [string, string]] : flat(v, `${p}${k}.`)))
    const ruMap = new Map(flat(ru as Tree))
    for (const [k, v] of [...flat(kz as Tree), ...flat(en as Tree)]) {
      expect(v.trim(), k).not.toBe('')
      if (!['common.brand'].includes(k) && /[а-яё]/i.test(v) && v === ruMap.get(k)) {
        // kz shares Cyrillic; only identical long strings are suspicious
        expect(v.length < 12, `${k} looks untranslated`).toBe(true)
      }
    }
  })

  it('plural forms required by Intl.PluralRules exist', () => {
    const forms = (loc: string) => new Intl.PluralRules(loc).resolvedOptions().pluralCategories
    for (const [tree, loc] of [[ru, 'ru'], [kz, 'kk'], [en, 'en']] as const) {
      for (const cat of forms(loc)) {
        if (cat === 'zero' || cat === 'two') continue
        expect((tree.dashboard as Tree)[`visibleEntries_${cat}`], `${loc} ${cat}`).toBeTypeOf('string')
      }
    }
  })
})
