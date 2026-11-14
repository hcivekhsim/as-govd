import { useState } from 'react'
import { LANGS, useLang, type Lang } from './i18n'
import { getStoredTheme, getSystemTheme, setTheme, type Theme } from './theme'
import { IconMoon, IconSun } from './icons'

export function LangSwitch() {
  const { lang, setLang, t } = useLang()

  const handle = (next: Lang) => {
    if (next !== lang) setLang(next)
  }

  return (
    <div
      className="lang-switch"
      role="group"
      aria-label={t('lang.switchLabel')}
    >
      {LANGS.map((item) => (
        <button
          key={item.value}
          type="button"
          className={`lang-option ${lang === item.value ? 'is-active' : ''}`}
          onClick={() => handle(item.value)}
          aria-pressed={lang === item.value}
        >
          {item.label}
        </button>
      ))}
    </div>
  )
}

export function ThemeToggle() {
  const [theme, setThemeState] = useState<Theme>(() => getStoredTheme() ?? getSystemTheme())
  const { t } = useLang()
  const isDark = theme === 'dark'

  const toggle = () => {
    const next: Theme = isDark ? 'light' : 'dark'
    setTheme(next)
    setThemeState(next)
  }

  const label = isDark ? t('theme.toLight') : t('theme.toDark')

  return (
    <button
      type="button"
      className="btn btn-ghost btn-icon"
      onClick={toggle}
      aria-label={label}
      title={label}
      data-theme-icon={theme}
    >
      {isDark ? <IconSun /> : <IconMoon />}
    </button>
  )
}
