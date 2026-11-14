import { useLang } from './i18n'
import { CONTACT_EMAIL, DUTY_PHONE_TEL, RECEPTION_PHONE_TEL } from './config/contacts'
import { YandexMap } from './YandexMap'
import { IconArrowLeft } from './icons'

export function ContactsPage({ onBack }: { onBack: () => void }) {
  const { t } = useLang()

  return (
    <div className="container-wide">
      <button type="button" className="back-link" onClick={onBack}>
        <IconArrowLeft width={16} height={16} />
        {t('detail.back')}
      </button>

      <div className="page-head">
        <div>
          <h1>{t('contacts.title')}</h1>
          <p className="page-sub">{t('contacts.subtitle')}</p>
        </div>
      </div>

      {}
      <div className="contacts-grid">
        <div className="card contacts-card">
          <dl className="contacts-list">
            <div className="contact-item">
              <dt>{t('contacts.addressLabel')}</dt>
              <dd>{t('contacts.address')}</dd>
            </div>

            <div className="contact-item">
              <dt>{t('contacts.hoursLabel')}</dt>
              <dd>
                {t('contacts.hours')}
                <span className="contact-extra">
                  {t('contacts.dutyLabel')}:{' '}
                  <a className="contact-link" href={`tel:${DUTY_PHONE_TEL}`}>
                    {t('contacts.dutyPhone')}
                  </a>
                </span>
              </dd>
            </div>

            <div className="contact-item">
              <dt>{t('contacts.receptionLabel')}</dt>
              <dd>
                <a className="contact-link" href={`tel:${RECEPTION_PHONE_TEL}`}>
                  {t('contacts.receptionPhone')}
                </a>
              </dd>
            </div>

            {}
            {CONTACT_EMAIL ? (
              <div className="contact-item">
                <dt>{t('contacts.emailLabel')}</dt>
                <dd>
                  <a className="contact-link" href={`mailto:${CONTACT_EMAIL}`}>
                    {CONTACT_EMAIL}
                  </a>
                </dd>
              </div>
            ) : null}

            <div className="contact-item">
              <dt>{t('contacts.chiefLabel')}</dt>
              <dd>
                {t('contacts.chiefName')}
                <span className="contact-extra">
                  {t('contacts.chiefPhoneLabel')}:{' '}
                  <a className="contact-link" href={`tel:${RECEPTION_PHONE_TEL}`}>
                    {t('contacts.receptionPhone')}
                  </a>
                </span>
                <span className="contact-extra">
                  {t('contacts.scheduleLabel')}: {t('contacts.schedule')}
                </span>
              </dd>
            </div>
          </dl>
        </div>

        <YandexMap />
      </div>
    </div>
  )
}
