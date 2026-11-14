import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import ruCatalog from './ru.json'
import beCatalog from './be.json'

export type Lang = 'ru' | 'be'

export const LANGS: ReadonlyArray<{ value: Lang; label: string }> = [
  { value: 'ru', label: 'Рус' },
  { value: 'be', label: 'Бел' },
]

const STORAGE_KEY = 'lang'

type Catalog = typeof ruCatalog

const catalogs: Record<Lang, Catalog> = {
  ru: ruCatalog,
  be: beCatalog as Catalog,
}

let currentLangValue: Lang = 'ru'

export function currentLang(): Lang {
  return currentLangValue
}

function localeTag(lang: Lang): string {
  return lang === 'be' ? 'be-BY' : 'ru-RU'
}

export function localeOf(lang: Lang): string {
  return localeTag(lang)
}

function detectInitialLang(): Lang {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved === 'be' || saved === 'ru') return saved
  } catch {
  }
  const candidates = navigator.languages ?? [navigator.language]
  for (const candidate of candidates) {
    if (candidate?.toLowerCase().startsWith('be')) return 'be'
  }
  return 'ru'
}

interface LangContextValue {
  lang: Lang
  setLang: (lang: Lang) => void
  t: (key: string, vars?: Record<string, string | number>) => string
  locale: string
}

const LangContext = createContext<LangContextValue | null>(null)

function resolveText(catalog: unknown, path: string): unknown {
  let node: unknown = catalog
  for (const part of path.split('.')) {
    if (node === null || typeof node !== 'object') return undefined
    node = (node as Record<string, unknown>)[part]
  }
  return node
}

export function LangProvider({ children }: { children: ReactNode }) {
  const [lang, setLangState] = useState<Lang>(detectInitialLang)

  useEffect(() => {
    currentLangValue = lang
    document.documentElement.lang = lang
  }, [lang])

  const setLang = useCallback((next: Lang) => {
    setLangState(next)
    try {
      localStorage.setItem(STORAGE_KEY, next)
    } catch {
    }
  }, [])

  const t = useCallback(
    (key: string, vars?: Record<string, string | number>): string => {
      let value = resolveText(catalogs[lang], key)
      if (value === undefined) value = resolveText(catalogs.ru, key)
      let text = typeof value === 'string' ? value : key
      if (vars) {
        for (const [name, replacement] of Object.entries(vars)) {
          text = text.replaceAll(`{{${name}}}`, String(replacement))
        }
      }
      return text
    },
    [lang],
  )

  const value = useMemo<LangContextValue>(
    () => ({ lang, setLang, t, locale: localeTag(lang) }),
    [lang, setLang, t],
  )

  return <LangContext.Provider value={value}>{children}</LangContext.Provider>
}

export function useLang(): LangContextValue {
  const ctx = useContext(LangContext)
  if (!ctx) throw new Error('useLang вне LangProvider')
  return ctx
}
