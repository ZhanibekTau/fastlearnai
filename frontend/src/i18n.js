import { createI18n } from 'vue-i18n'
import kk from '@/locales/kk.json'
import ru from '@/locales/ru.json'
import en from '@/locales/en.json'

export const SUPPORTED_LOCALES = ['kk', 'ru', 'en']
const DEFAULT_LOCALE = 'kk'

const saved = localStorage.getItem('locale')
const initialLocale = SUPPORTED_LOCALES.includes(saved) ? saved : DEFAULT_LOCALE

const i18n = createI18n({
  legacy: false,
  locale: initialLocale,
  fallbackLocale: DEFAULT_LOCALE,
  messages: { kk, ru, en },
})

export function setLocale(locale) {
  if (!SUPPORTED_LOCALES.includes(locale)) return
  i18n.global.locale.value = locale
  localStorage.setItem('locale', locale)
  document.documentElement.lang = locale
}

document.documentElement.lang = initialLocale

export default i18n
