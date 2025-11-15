import { createI18n } from 'vue-i18n'

const messages = {
}

const i18n = createI18n({
  legacy: false,
  locale: 'zh', // 默认语言
  fallbackLocale: 'zh', // 回退语言
  messages
})

export default i18n
export { messages };