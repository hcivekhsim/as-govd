import { forwardRef, useImperativeHandle, useRef, useState } from 'react'
import HCaptcha from '@hcaptcha/react-hcaptcha'
import { useLang } from './i18n'
import { getSystemTheme, getStoredTheme } from './theme'

const SITE_KEY = (import.meta.env.VITE_HCAPTCHA_KEY ?? '') as string

export function isCaptchaEnabled(): boolean {
  return SITE_KEY !== ''
}

export interface CaptchaFieldHandle {
  reset: () => void
}

interface CaptchaFieldProps {
  onToken: (token: string | null) => void
}

export const CaptchaField = forwardRef<CaptchaFieldHandle, CaptchaFieldProps>(
  function CaptchaField({ onToken }, ref) {
    const { t, lang } = useLang()
    const widgetRef = useRef<HCaptcha>(null)
    const [failed, setFailed] = useState(false)

    useImperativeHandle(
      ref,
      () => ({
        reset: () => {
          setFailed(false)
          onToken(null)
          widgetRef.current?.resetCaptcha()
        },
      }),
      [onToken],
    )

    if (!SITE_KEY) {
      return null
    }

    const languageOverride = lang === 'be' ? 'be' : 'ru'
    const theme: 'dark' | 'light' = getStoredTheme() ?? getSystemTheme()

    const handleVerify = (token: string) => {
      setFailed(false)
      onToken(token)
    }

    const handleError = () => {
      onToken(null)
      setFailed(true)
    }

    return (
      <div className="captcha-field">
        <HCaptcha
          ref={widgetRef}
          sitekey={SITE_KEY}
          theme={theme}
          languageOverride={languageOverride}
          onVerify={handleVerify}
          onExpire={() => {
            onToken(null)
            widgetRef.current?.resetCaptcha()
          }}
          onError={handleError}
        />
        {failed ? <p className="captcha-error">{t('captcha.failed')}</p> : null}
      </div>
    )
  },
)
