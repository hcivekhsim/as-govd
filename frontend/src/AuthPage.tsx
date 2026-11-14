import { useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { api, ApiError, type RegisterPayload } from './api'
import { humanizeFields } from './errors'
import { useLang } from './i18n'
import { LangSwitch, ThemeToggle } from './controls'
import { CaptchaField, isCaptchaEnabled, type CaptchaFieldHandle } from './CaptchaField'
import { Alert, Field, Button } from './components'

interface AuthPageProps {
  onAuthed: (login: string) => void
}

type Mode = 'login' | 'register'

const USER_TYPE_KEYS = ['individual', 'legal', 'foreign', 'stateless'] as const

interface RegisterForm {
  userType: (typeof USER_TYPE_KEYS)[number]
  nameOrganization: string
  login: string
  password: string
  passwordConfirm: string
  email: string
  phone: string
  lastName: string
  firstName: string
  patronymic: string
  isAbroad: boolean
  locality: string
  street: string
  house: string
  apartment: string
  noteAddress: string
  postalCode: string
}

const emptyRegister: RegisterForm = {
  userType: 'individual',
  nameOrganization: '',
  login: '',
  password: '',
  passwordConfirm: '',
  email: '',
  phone: '',
  lastName: '',
  firstName: '',
  patronymic: '',
  isAbroad: false,
  locality: '',
  street: '',
  house: '',
  apartment: '',
  noteAddress: '',
  postalCode: '',
}

export function AuthPage({ onAuthed }: AuthPageProps) {
  const { t } = useLang()
  const [mode, setMode] = useState<Mode>('login')
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [reg, setReg] = useState<RegisterForm>(emptyRegister)
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [busy, setBusy] = useState(false)
  const [loginCaptcha, setLoginCaptcha] = useState<string | null>(null)
  const [regCaptcha, setRegCaptcha] = useState<string | null>(null)
  const loginCaptchaRef = useRef<CaptchaFieldHandle>(null)
  const regCaptchaRef = useRef<CaptchaFieldHandle>(null)

  const resetCaptcha = (isLogin: boolean) => {
    if (isLogin) {
      loginCaptchaRef.current?.reset()
      setLoginCaptcha(null)
    } else {
      regCaptchaRef.current?.reset()
      setRegCaptcha(null)
    }
  }

  const applyError = (e: unknown, isLogin: boolean) => {
    if (e instanceof ApiError) {
      setError(e.message)
      setFieldErrors(humanizeFields(e.fields))
    } else {
      setError(t('auth.requestFailed'))
    }
    resetCaptcha(isLogin)
  }

  const handleLogin = async (e: FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setFieldErrors({})
    try {
      await api.login(login, password, loginCaptcha ?? '')
      resetCaptcha(true)
      onAuthed(login)
    } catch (err) {
      applyError(err, true)
    } finally {
      setBusy(false)
    }
  }

  const handleRegister = async (e: FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setFieldErrors({})
    try {
      const payload: RegisterPayload = {
        user_type: reg.userType,
        login: reg.login,
        password: reg.password,
        password_confirm: reg.passwordConfirm,
        email: reg.email,
        phone: reg.phone,
        last_name: reg.lastName,
        first_name: reg.firstName,
        is_abroad: reg.isAbroad,
        locality: reg.locality,
        street: reg.street,
        house: reg.house,
        captcha_token: regCaptcha ?? '',
      }
      if (reg.userType === 'legal') payload.name_organization = reg.nameOrganization
      if (reg.patronymic) payload.patronymic = reg.patronymic
      if (reg.apartment) payload.apartment = reg.apartment
      if (reg.noteAddress) payload.note_address = reg.noteAddress
      if (reg.postalCode) payload.postal_code = reg.postalCode

      await api.register(payload)
      resetCaptcha(false)
      onAuthed(reg.login)
    } catch (err) {
      applyError(err, false)
    } finally {
      setBusy(false)
    }
  }

  const setRegField = <K extends keyof RegisterForm>(key: K, value: RegisterForm[K]) => {
    setReg((prev) => ({ ...prev, [key]: value }))
  }

  const steps = [t('auth.steps.0'), t('auth.steps.1'), t('auth.steps.2')]

  return (
    <div className="auth">
      <div className="auth-hero">
        <div className="auth-hero-top">
          <div className="wordmark">
            {t('app.wordmark')} <small>{t('org.short')}</small>
          </div>
          <div className="topbar-controls">
            <LangSwitch />
            <ThemeToggle />
          </div>
        </div>
        <h1>{t('auth.heroTitle')}</h1>
        <p className="auth-sub">{t('auth.heroSub')}</p>
        <div className="auth-steps">
          {steps.map((step, index) => (
            <div className="auth-step" key={step}>
              <span className="num">{String(index + 1).padStart(2, '0')}</span>
              <span className="text">{step}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="auth-panel">
        <div className="auth-card">
          <div className="tabs">
            <button
              type="button"
              className={`tab ${mode === 'login' ? 'is-active' : ''}`}
              onClick={() => {
                setMode('login')
                setError(null)
                setFieldErrors({})
              }}
            >
              {t('auth.tabLogin')}
            </button>
            <button
              type="button"
              className={`tab ${mode === 'register' ? 'is-active' : ''}`}
              onClick={() => {
                setMode('register')
                setError(null)
                setFieldErrors({})
              }}
            >
              {t('auth.tabRegister')}
            </button>
          </div>

          {error ? <Alert>{error}</Alert> : null}

          {mode === 'login' ? (
            <form className="form-grid" onSubmit={handleLogin}>
              <Field label={t('fields.login')} required error={fieldErrors[t('fieldErrors.Login')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Login')] ? 'is-invalid' : ''}`}
                  value={login}
                  onChange={(e) => setLogin(e.target.value)}
                  placeholder={t('auth.placeholderLogin')}
                  autoComplete="username"
                />
              </Field>
              <Field label={t('fields.password')} required error={fieldErrors[t('fieldErrors.Password')]}>
                <input
                  type="password"
                  className={`input ${fieldErrors[t('fieldErrors.Password')] ? 'is-invalid' : ''}`}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••••"
                  autoComplete="current-password"
                />
              </Field>
              <CaptchaField ref={loginCaptchaRef} onToken={setLoginCaptcha} />
              <Button type="submit" disabled={busy || (isCaptchaEnabled() && !loginCaptcha)}>
                {busy ? t('auth.loggingIn') : t('auth.login')}
              </Button>
              <p className="auth-note">{t('auth.passwordHint')}</p>
            </form>
          ) : (
            <form className="form-grid" onSubmit={handleRegister}>
              <Field
                label={t('fields.userType')}
                required
                error={fieldErrors[t('fieldErrors.UserType')]}
              >
                <select
                  className="select"
                  value={reg.userType}
                  onChange={(e) => setRegField('userType', e.target.value as RegisterForm['userType'])}
                >
                  {USER_TYPE_KEYS.map((key) => (
                    <option key={key} value={key}>
                      {t(`auth.userTypes.${key}`)}
                    </option>
                  ))}
                </select>
              </Field>

              {reg.userType === 'legal' ? (
                <Field
                  label={t('fields.nameOrganization')}
                  required
                  error={fieldErrors[t('fieldErrors.NameOrganization')]}
                >
                  <input
                    className={`input ${
                      fieldErrors[t('fieldErrors.NameOrganization')] ? 'is-invalid' : ''
                    }`}
                    value={reg.nameOrganization}
                    onChange={(e) => setRegField('nameOrganization', e.target.value)}
                  />
                </Field>
              ) : null}

              <Field label={t('fields.login')} required error={fieldErrors[t('fieldErrors.Login')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Login')] ? 'is-invalid' : ''}`}
                  value={reg.login}
                  onChange={(e) => setRegField('login', e.target.value)}
                  autoComplete="username"
                />
              </Field>
              <Field label={t('fields.password')} required error={fieldErrors[t('fieldErrors.Password')]}>
                <input
                  type="password"
                  className={`input ${fieldErrors[t('fieldErrors.Password')] ? 'is-invalid' : ''}`}
                  value={reg.password}
                  onChange={(e) => setRegField('password', e.target.value)}
                  autoComplete="new-password"
                />
              </Field>
              <Field
                label={t('fields.passwordConfirm')}
                required
                error={fieldErrors[t('fieldErrors.PasswordConfirm')]}
              >
                <input
                  type="password"
                  className={`input ${
                    fieldErrors[t('fieldErrors.PasswordConfirm')] ? 'is-invalid' : ''
                  }`}
                  value={reg.passwordConfirm}
                  onChange={(e) => setRegField('passwordConfirm', e.target.value)}
                  autoComplete="new-password"
                />
              </Field>
              <Field label={t('fields.email')} required error={fieldErrors[t('fieldErrors.Email')]}>
                <input
                  type="email"
                  className={`input ${fieldErrors[t('fieldErrors.Email')] ? 'is-invalid' : ''}`}
                  value={reg.email}
                  onChange={(e) => setRegField('email', e.target.value)}
                  placeholder="name@example.com"
                />
              </Field>
              <Field label={t('fields.phone')} required error={fieldErrors[t('fieldErrors.Phone')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Phone')] ? 'is-invalid' : ''}`}
                  value={reg.phone}
                  onChange={(e) => setRegField('phone', e.target.value)}
                  placeholder="+37529112233"
                />
              </Field>

              <Field label={t('fields.lastName')} required error={fieldErrors[t('fieldErrors.LastName')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.LastName')] ? 'is-invalid' : ''}`}
                  value={reg.lastName}
                  onChange={(e) => setRegField('lastName', e.target.value)}
                />
              </Field>
              <Field label={t('fields.firstName')} required error={fieldErrors[t('fieldErrors.FirstName')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.FirstName')] ? 'is-invalid' : ''}`}
                  value={reg.firstName}
                  onChange={(e) => setRegField('firstName', e.target.value)}
                />
              </Field>
              <Field label={t('fields.patronymic')} error={fieldErrors[t('fieldErrors.Patronymic')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Patronymic')] ? 'is-invalid' : ''}`}
                  value={reg.patronymic}
                  onChange={(e) => setRegField('patronymic', e.target.value)}
                />
              </Field>

              <label className="checkbox span-2">
                <input
                  type="checkbox"
                  checked={reg.isAbroad}
                  onChange={(e) => setRegField('isAbroad', e.target.checked)}
                />
                {t('fields.isAbroad')}
              </label>

              <Field label={t('fields.locality')} required error={fieldErrors[t('fieldErrors.Locality')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Locality')] ? 'is-invalid' : ''}`}
                  value={reg.locality}
                  onChange={(e) => setRegField('locality', e.target.value)}
                />
              </Field>
              <Field label={t('fields.street')} required error={fieldErrors[t('fieldErrors.Street')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Street')] ? 'is-invalid' : ''}`}
                  value={reg.street}
                  onChange={(e) => setRegField('street', e.target.value)}
                />
              </Field>
              <Field label={t('fields.house')} required error={fieldErrors[t('fieldErrors.House')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.House')] ? 'is-invalid' : ''}`}
                  value={reg.house}
                  onChange={(e) => setRegField('house', e.target.value)}
                />
              </Field>
              <Field label={t('fields.apartment')} error={fieldErrors[t('fieldErrors.Apartment')]}>
                <input
                  className={`input ${fieldErrors[t('fieldErrors.Apartment')] ? 'is-invalid' : ''}`}
                  value={reg.apartment}
                  onChange={(e) => setRegField('apartment', e.target.value)}
                />
              </Field>
              <Field label={t('fields.postalCode')} error={fieldErrors[t('fieldErrors.PostalCode')]}>
                <input
                  className={`input ${
                    fieldErrors[t('fieldErrors.PostalCode')] ? 'is-invalid' : ''
                  }`}
                  value={reg.postalCode}
                  onChange={(e) => setRegField('postalCode', e.target.value)}
                />
              </Field>
              <Field label={t('fields.noteAddress')} error={fieldErrors[t('fieldErrors.NoteAddress')]}>
                <input
                  className={`input ${
                    fieldErrors[t('fieldErrors.NoteAddress')] ? 'is-invalid' : ''
                  }`}
                  value={reg.noteAddress}
                  onChange={(e) => setRegField('noteAddress', e.target.value)}
                />
              </Field>

              <CaptchaField ref={regCaptchaRef} onToken={setRegCaptcha} />
              <Button type="submit" disabled={busy || (isCaptchaEnabled() && !regCaptcha)}>
                {busy ? t('auth.registering') : t('auth.register')}
              </Button>
              <p className="auth-note">{t('auth.passwordHintRegister')}</p>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
