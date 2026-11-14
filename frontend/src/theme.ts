export type Theme = 'dark' | 'light'

const STORAGE_KEY = 'theme'

export function getStoredTheme(): Theme | null {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    return saved === 'dark' || saved === 'light' ? saved : null
  } catch {
    return null
  }
}

export function getSystemTheme(): Theme {
  if (typeof window === 'undefined' || !window.matchMedia) return 'dark'
  return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
}

export function applyTheme(theme: Theme): void {
  document.documentElement.dataset.theme = theme
}

export function setTheme(theme: Theme): void {
  applyTheme(theme)
  try {
    localStorage.setItem(STORAGE_KEY, theme)
  } catch {
  }
}

export function resolveInitialTheme(): Theme {
  return getStoredTheme() ?? getSystemTheme()
}
