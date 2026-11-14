import { useEffect, useState } from 'react'
import { api } from './api'
import { LangProvider, useLang } from './i18n'
import { LangSwitch, ThemeToggle } from './controls'
import { AuthPage } from './AuthPage'
import { MyAppealsPage } from './MyAppealsPage'
import { NewAppealPage } from './NewAppealPage'
import { AppealDetailPage } from './AppealDetailPage'
import { ContactsPage } from './ContactsPage'
import { Loading } from './components'
import { IconLogOut } from './icons'

type View =
  | { name: 'list' }
  | { name: 'new' }
  | { name: 'detail'; id: number }
  | { name: 'contacts' }

function Shell() {
  const { t, lang, setLang } = useLang()
  const [checking, setChecking] = useState(true)
  const [authed, setAuthed] = useState(false)
  const [userLogin, setUserLogin] = useState('')
  const [view, setView] = useState<View>({ name: 'list' })

  useEffect(() => {
    let cancelled = false
    api
      .checkSession()
      .then((res) => {
        if (cancelled) return
        const explicit = localStorage.getItem('lang')
        if (
          !explicit &&
          (res.language === 'ru' || res.language === 'be')
        ) {
          setLang(res.language)
        }
        setAuthed(true)
      })
      .catch(() => {
        if (!cancelled) setAuthed(false)
      })
      .finally(() => {
        if (!cancelled) setChecking(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    if (!authed) return
    api.setLanguage(lang).catch(() => {})
  }, [authed, lang])

  const handleAuthed = (login: string) => {
    setUserLogin(login)
    setAuthed(true)
    setView({ name: 'list' })
  }

  const handleLogout = async () => {
    try {
      await api.logout()
    } catch {
    }
    setUserLogin('')
    setAuthed(false)
    setView({ name: 'list' })
  }

  if (checking) {
    return <Loading text={t('app.checkingSession')} />
  }

  if (!authed) {
    return <AuthPage onAuthed={handleAuthed} />
  }

  return (
    <>
      <header className="topbar">
        <div className="wordmark">
          {t('app.wordmark')} <small>{t('org.short')}</small>
        </div>
        <nav className="topnav">
          <button
            type="button"
            className={`topnav-link ${view.name === 'list' ? 'is-active' : ''}`}
            onClick={() => setView({ name: 'list' })}
          >
            {t('app.nav.myAppeals')}
          </button>
          <button
            type="button"
            className={`topnav-link ${view.name === 'new' ? 'is-active' : ''}`}
            onClick={() => setView({ name: 'new' })}
          >
            {t('app.nav.newAppeal')}
          </button>
          <button
            type="button"
            className={`topnav-link ${view.name === 'contacts' ? 'is-active' : ''}`}
            onClick={() => setView({ name: 'contacts' })}
          >
            {t('contacts.nav')}
          </button>
        </nav>
        <div className="topbar-controls">
          <LangSwitch />
          <ThemeToggle />
        </div>
        <div className="topbar-user">
          <span className="topbar-login">{userLogin || t('app.cabinet')}</span>
          <button
            type="button"
            className="btn btn-ghost btn-icon"
            onClick={handleLogout}
            aria-label={t('app.logout')}
            title={t('app.logout')}
          >
            <IconLogOut />
          </button>
        </div>
      </header>

      <main>
        {view.name === 'list' && (
          <MyAppealsPage
            onOpenAppeal={(id) => setView({ name: 'detail', id })}
            onNewAppeal={() => setView({ name: 'new' })}
          />
        )}
        {view.name === 'new' && (
          <NewAppealPage
            onCreated={(appeal) => setView({ name: 'detail', id: appeal.id })}
            onCancel={() => setView({ name: 'list' })}
          />
        )}
        {view.name === 'detail' && (
          <AppealDetailPage appealId={view.id} onBack={() => setView({ name: 'list' })} />
        )}
        {view.name === 'contacts' && <ContactsPage onBack={() => setView({ name: 'list' })} />}
      </main>

      <footer className="footer">
        <div className="footer-inner">
          <span className="footer-org">{t('org.full')}</span>
          <span className="footer-region">{t('org.region')}</span>
          <span className="footer-note">{t('app.footerNote')}</span>
          <button
            type="button"
            className="footer-link"
            onClick={() => setView({ name: 'contacts' })}
          >
            {t('contacts.nav')}
          </button>
        </div>
      </footer>
    </>
  )
}

export default function App() {
  return (
    <LangProvider>
      <Shell />
    </LangProvider>
  )
}
