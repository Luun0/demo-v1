import { describe, expect, it } from 'vitest'
import { findFaq, type FaqId } from './faq'

const intents: Array<{ id: FaqId; questions: string[] }> = [
  { id: 'about', questions: ['для чего нужна платформа', 'what is sen support for', 'платформа не үшін қажет', 'что делает sen support', 'how does the platform help', 'sen support не үшін', 'расскажи о платформе', 'tell me about sen support', 'платформа туралы'] },
  { id: 'roles', questions: ['чем отличаются роли', 'what can each role see', 'әр рөл не көре алады', 'что может родитель', 'what does a parent do', 'маманның рөлі қандай', 'права администратора', 'what can an admin access', 'рөлдер туралы'] },
  { id: 'registration', questions: ['как пройти регистрацию', 'how do I register', 'қалай тіркелемін', 'как создать аккаунт', 'how can I sign up', 'тіркелу жолы', 'регистрация специалиста', 'specialist account registration', 'аккаунт ашу'] },
  { id: 'approval', questions: ['почему специалист не выбирается', 'why is the specialist missing', 'неге маман тізімде жоқ', 'почему нет в списке', 'specialist not available to assign', 'маманды таңдай алмаймын', 'администратор не одобрил специалиста', 'specialist approval status', 'маман бекітілмеген'] },
  { id: 'password', questions: ['забыл пароль', 'how do I reset my password', 'құпиясөзді ұмыттым', 'как восстановить пароль', 'forgot password help', 'құпиясөзді қалпына келтіру', 'не приходит ссылка сброса', 'password reset email', 'парольді ауыстыру'] },
  { id: 'profile', questions: ['как добавить профиль ребёнка', 'how do I add a child profile', 'бала профилін қалай қосамын', 'создать карточку ребёнка', 'child profile setup', 'баланың профилі', 'редактировать профиль ребёнка', 'child information profile', 'бала туралы профиль'] },
  { id: 'team', questions: ['как назначить специалиста', 'how to assign a specialist', 'маманды қалай тағайындаймын', 'добавить специалиста ребёнку', 'assign specialist to child', 'маман қосу', 'назначение в команду', 'team assignment help', 'командаға маман қосу'] },
  { id: 'schedule', questions: ['где посмотреть расписание', 'where can I see the schedule', 'кестені қайдан көремін', 'расписание занятий', 'view support sessions', 'кездесу кестесі', 'сессии и календарь', 'appointments schedule', 'сабақ кестесі'] },
  { id: 'diary', questions: ['кто видит дневник', 'who can see the diary', 'күнделікті кім көре алады', 'доступ к записям дневника', 'diary visibility', 'күнделікке кім кіре алады', 'заметки о ребёнке', 'specialist diary notes', 'күнделік жазбалары'] },
  { id: 'privacy', questions: ['как защищён доступ к данным', 'how is data access protected', 'деректерге қолжетімділік қалай қорғалады', 'кто видит мои данные', 'data privacy and access', 'деректердің құпиялылығы', 'защита данных ребёнка', 'privacy of child data', 'баланың деректерін қорғау'] },
  { id: 'medical', questions: ['медицинский совет', 'can I get medical advice', 'медициналық кеңес', 'что делать с диагнозом', 'health advice for my child', 'баланың диагнозы туралы', 'помогите при кризисе', 'medical diagnosis question', 'денсаулық бойынша кеңес'] },
]

describe('FAQ routing for likely support questions', () => {
  it('routes 1,100 synthetic phrasings to their intended topics', () => {
    let routed = 0
    const prefixes = ['Подскажите, пожалуйста', 'У меня вопрос', 'Не понимаю', 'Помогите узнать', 'Could you explain', 'Quick question', 'Please help', 'Айтып жіберіңізші', 'Менде сұрақ бар', 'Көмектесіңіз']
    const suffixes = ['?', ' пожалуйста', ' please', ' рахмет', ' спасибо', ' заранее благодарю', ' көп рақмет', ' thank you', ' помогите', ' түсіндіріңізші']
    for (const intent of intents) {
      for (let variant = 0; variant < 100; variant += 1) {
        const question = intent.questions[variant % intent.questions.length]
        const prefix = prefixes[variant % prefixes.length]
        const suffix = suffixes[Math.floor(variant / prefixes.length)]
        expect(findFaq(`${prefix}: ${question}${suffix}`), `case ${routed + 1}`).toBe(intent.id)
        routed += 1
      }
    }
    expect(routed).toBe(1100)
  })

  it('does not guess for unrelated questions', () => {
    expect(findFaq('Какая сегодня погода в Алматы?')).toBeUndefined()
    expect(findFaq('Напиши мне сочинение про космос')).toBeUndefined()
  })
})
