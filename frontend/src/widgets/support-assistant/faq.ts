export const faqItems = [
  { id: 'about', terms: ['platform', 'платформ', 'платформа', 'sen support', 'не үшін', 'что это', 'what is'] },
  { id: 'roles', terms: ['role', 'roles', 'роли', 'роль', 'рөл', 'рөлі', 'рөлдер', 'кто что видит', 'parent', 'родитель', 'ата ана', 'admin', 'администратор', 'әкімші'] },
  { id: 'registration', terms: ['register', 'registration', 'регистрац', 'зарегистр', 'тіркел', 'создать аккаунт', 'sign up', 'аккаунт', 'жазылу'] },
  { id: 'approval', terms: ['approval', 'одобрен', 'одобр', 'администратор не одобрил', 'мақұлда', 'бекітілмеген', 'маман тізімде', 'специалист в списке', 'specialist missing', 'not available', 'cannot choose', "can't select", 'не выбирается', 'не доступен', 'недоступен', 'нет в списке', 'не могу выбрать', 'таңдай алмаймын', 'cannot select', 'not in the list'] },
  { id: 'password', terms: ['password', 'парол', 'құпиясөз', 'reset', 'сброс', 'восстанов', 'забыл пароль', 'forgot password'] },
  { id: 'profile', terms: ['profile', 'профил', 'ребен', 'child profile', 'child', 'карточка', 'бала', 'add a child'] },
  { id: 'team', terms: ['assign', 'назнач', 'тағайында', 'команда', 'team', 'маманды', 'маман қосу', 'специалиста назнач', 'специалиста ребёнку', 'specialist to child', 'assign a specialist'] },
  { id: 'schedule', terms: ['schedule', 'расписан', 'кесте', 'календар', 'график', 'session', 'сесси', 'сеанс', 'кездесу', 'сабақ', 'appointment'] },
  { id: 'diary', terms: ['diary', 'дневник', 'күнделік', 'запис', 'жазбаны', 'жазба', 'заметки о ребёнке'] },
  { id: 'privacy', terms: ['privacy', 'data', 'access', 'protect', 'данные', 'доступ', 'деректер', 'қолжетімділік', 'құпиялылық', 'құпиялылы', 'защит', 'защита данных', 'безопас', 'пароль', 'құпиясөз'] },
  { id: 'medical', terms: ['medical', 'health', 'advice', 'медицин', 'diagnos', 'диагноз', 'денсаулық', 'здоров', 'дәрігер', 'дәрі', 'совет', 'кеңес', 'кризис'] },
] as const

export type FaqId = (typeof faqItems)[number]['id']

export function findFaq(query: string): FaqId | undefined {
  const normalized = query.toLocaleLowerCase().replace(/[ё]/g, 'е').replace(/[^\p{L}\p{N}\s]/gu, ' ').trim()
  if (!normalized) return undefined
  const words = new Set(normalized.split(/\s+/).filter((word) => word.length > 2))
  let best: { id: FaqId; score: number } | undefined
  for (const item of faqItems) {
    const score = item.terms.reduce((sum, term) => {
      const normalizedTerm = term.toLocaleLowerCase().replace(/[ё]/g, 'е')
      const matches = normalizedTerm.includes(' ')
        ? normalized.includes(normalizedTerm)
        : [...words].some((word) => word === normalizedTerm || (normalizedTerm.length >= 5 && word.startsWith(normalizedTerm)))
      return sum + (matches ? Math.max(1, normalizedTerm.split(/\s+/).length) : 0)
    }, 0)
    if (score > (best?.score ?? 0)) best = { id: item.id, score }
  }
  return best && best.score > 0 ? best.id : undefined
}
