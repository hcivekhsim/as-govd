import { useEffect, useMemo, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { api, ApiError } from './api'
import type { AppealDetail, DictionaryItem } from './types'
import { humanizeFields } from './errors'
import { useLang } from './i18n'
import { CaptchaField, isCaptchaEnabled, type CaptchaFieldHandle } from './CaptchaField'
import { Alert, Field, Button, Loading } from './components'
import { IconPaperclip, IconSend, IconTrash, IconX } from './icons'

interface NewAppealPageProps {
  onCreated: (appeal: AppealDetail) => void
  onCancel: () => void
}

const MAX_FILES = 5
const MAX_FILE_SIZE = 5 * 1024 * 1024
const ALLOWED_EXT = ['doc', 'docx', 'txt', 'rtf', 'odf', 'pdf', 'jpg', 'jpeg', 'tiff', 'png']

interface PickedFile {
  file: File
  error?: string
}

interface FormState {
  applicantCategoryCode: string
  topicCode: string
  organizationName: string
  appealKindCode: string
  appealTypeCode: string
  lastName: string
  firstName: string
  patronymic: string
  identificationNumber: string
  physicalStatusCode: string
  socialCategoryCode: string
  email: string
  citizenshipCode: string
  region: string
  district: string
  locality: string
  street: string
  house: string
  apartment: string
  noteAddress: string
  postalCode: string
  phone: string
  content: string
  plannedRegistrationDate: string
  isAbroad: boolean
}

const emptyForm: FormState = {
  applicantCategoryCode: '',
  topicCode: '',
  organizationName: '',
  appealKindCode: '',
  appealTypeCode: '',
  lastName: '',
  firstName: '',
  patronymic: '',
  identificationNumber: '',
  physicalStatusCode: '',
  socialCategoryCode: '',
  email: '',
  citizenshipCode: '',
  region: '',
  district: '',
  locality: '',
  street: '',
  house: '',
  apartment: '',
  noteAddress: '',
  postalCode: '',
  phone: '',
  content: '',
  plannedRegistrationDate: '',
  isAbroad: false,
}

export function NewAppealPage({ onCreated, onCancel }: NewAppealPageProps) {
  const { t } = useLang()
  const [dictionaries, setDictionaries] = useState<Record<string, DictionaryItem[]>>({})
  const [loadingDicts, setLoadingDicts] = useState(true)
  const [form, setForm] = useState<FormState>(emptyForm)
  const [picked, setPicked] = useState<PickedFile[]>([])
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [busy, setBusy] = useState(false)
  const [captcha, setCaptcha] = useState<string | null>(null)
  const captchaRef = useRef<CaptchaFieldHandle>(null)

  useEffect(() => {
    const kinds = [
      'applicant-categories',
      'topics',
      'appeal-kinds',
      'appeal-types',
      'physical-statuses',
      'social-categories',
      'citizenships',
    ]
    let cancelled = false
    Promise.all(
      kinds.map((kind) =>
        api.dictionary(kind).then((res) => [kind, res.items] as const),
      ),
    )
      .then((entries) => {
        if (cancelled) return
        const map: Record<string, DictionaryItem[]> = {}
        for (const [kind, items] of entries) {
          map[kind] = items
        }
        setDictionaries(map)
        setLoadingDicts(false)
      })
      .catch(() => {
        if (!cancelled) {
          setLoadingDicts(false)
          setError(t('new.dictsFailed'))
        }
      })

    return () => {
      cancelled = true
    }
  }, [])

  const set = <K extends keyof FormState>(key: K, value: FormState[K]) => {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  const onFilesPicked = (list: FileList | null) => {
    if (!list) return
    const room = MAX_FILES - picked.length
    if (room <= 0) return

    const overflow = list.length > room
    const next: PickedFile[] = Array.from(list)
      .slice(0, room)
      .map((file) => {
        const item: PickedFile = { file }
        const ext = (file.name.split('.').pop() ?? '').toLowerCase()
        if (file.size > MAX_FILE_SIZE) {
          item.error = t('new.fileTooLarge')
        } else if (!ALLOWED_EXT.includes(ext)) {
          item.error = t('new.fileWrongType')
        }
        return item
      })
    setPicked((prev) => [...prev, ...next])
    if (overflow) {
      setError(t('new.tooManyFiles'))
    }
  }

  const hasPickedErrors = picked.some((p) => p.error)

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setFieldErrors({})
    try {
      const data = new FormData()
      data.append('applicant_category_code', form.applicantCategoryCode)
      data.append('topic_code', form.topicCode)
      data.append('appeal_kind_code', form.appealKindCode)
      data.append('appeal_type_code', form.appealTypeCode)
      data.append('last_name', form.lastName)
      data.append('first_name', form.firstName)
      data.append('identification_number', form.identificationNumber)
      data.append('email', form.email)
      data.append('locality', form.locality)
      data.append('street', form.street)
      data.append('house', form.house)
      data.append('phone', form.phone)
      data.append('content', form.content)
      data.append('is_abroad', form.isAbroad ? 'true' : 'false')
      data.append('captcha_token', captcha ?? '')

      if (form.applicantCategoryCode === 'legal' && form.organizationName) {
        data.append('organization_name', form.organizationName)
      }
      if (form.patronymic) data.append('patronymic', form.patronymic)
      if (form.physicalStatusCode) data.append('physical_status_code', form.physicalStatusCode)
      if (form.socialCategoryCode) data.append('social_category_code', form.socialCategoryCode)
      if (form.citizenshipCode) data.append('citizenship_code', form.citizenshipCode)
      if (form.region) data.append('region', form.region)
      if (form.district) data.append('district', form.district)
      if (form.apartment) data.append('apartment', form.apartment)
      if (form.noteAddress) data.append('note_address', form.noteAddress)
      if (form.postalCode) data.append('postal_code', form.postalCode)
      if (form.plannedRegistrationDate) {
        data.append('planned_registration_date', form.plannedRegistrationDate)
      }

      for (const item of picked) {
        if (!item.error) {
          data.append('files', item.file)
        }
      }

      const detail = await api.createAppeal(data)
      captchaRef.current?.reset()
      setCaptcha(null)
      onCreated(detail)
    } catch (err) {
      captchaRef.current?.reset()
      setCaptcha(null)
      if (err instanceof ApiError) {
        setError(err.message)
        setFieldErrors(humanizeFields(err.fields))
      } else {
        setError(t('new.submitFailed'))
      }
    } finally {
      setBusy(false)
    }
  }

  const dict = useMemo(() => {
    return {
      categories: dictionaries['applicant-categories'] ?? [],
      topics: dictionaries['topics'] ?? [],
      kinds: dictionaries['appeal-kinds'] ?? [],
      types: dictionaries['appeal-types'] ?? [],
      physical: dictionaries['physical-statuses'] ?? [],
      social: dictionaries['social-categories'] ?? [],
      citizenships: dictionaries['citizenships'] ?? [],
    }
  }, [dictionaries])

  if (loadingDicts) {
    return (
      <div className="container-narrow">
        <Loading text={t('new.loadingDicts')} />
      </div>
    )
  }

  const fl = (key: string) => fieldErrors[t(`fieldErrors.${key}`)]

  return (
    <div className="container-narrow">
      <button type="button" className="back-link" onClick={onCancel}>
        <IconX width={16} height={16} />
        {t('new.back')}
      </button>
      <div className="page-head">
        <div>
          <h1>{t('new.title')}</h1>
          <p className="page-sub">{t('new.subtitle')}</p>
        </div>
      </div>

      {error ? <Alert>{error}</Alert> : null}

      <form onSubmit={submit}>
        <div className="form-section">
          <h2>{t('new.sectionType')}</h2>
          <p className="form-section-note">{t('new.requiredNote')}</p>
          <div className="form-grid">
            <Field label={t('fields.applicantCategory')} required error={fl('ApplicantCategoryCode')}>
              <select
                className={`select ${fl('ApplicantCategoryCode') ? 'is-invalid' : ''}`}
                value={form.applicantCategoryCode}
                onChange={(e) => set('applicantCategoryCode', e.target.value)}
              >
                <option value="">{t('new.chooseCategory')}</option>
                {dict.categories.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
            <Field label={t('fields.topic')} required error={fl('TopicCode')}>
              <select
                className={`select ${fl('TopicCode') ? 'is-invalid' : ''}`}
                value={form.topicCode}
                onChange={(e) => set('topicCode', e.target.value)}
              >
                <option value="">{t('new.chooseTopic')}</option>
                {dict.topics.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
            <Field label={t('fields.appealKind')} required error={fl('AppealKindCode')}>
              <select
                className={`select ${fl('AppealKindCode') ? 'is-invalid' : ''}`}
                value={form.appealKindCode}
                onChange={(e) => set('appealKindCode', e.target.value)}
              >
                <option value="">{t('new.chooseKind')}</option>
                {dict.kinds.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
            <Field label={t('fields.appealType')} required error={fl('AppealTypeCode')}>
              <select
                className={`select ${fl('AppealTypeCode') ? 'is-invalid' : ''}`}
                value={form.appealTypeCode}
                onChange={(e) => set('appealTypeCode', e.target.value)}
              >
                <option value="">{t('new.chooseType')}</option>
                {dict.types.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
          </div>
        </div>

        <div className="form-section">
          <h2>{t('new.sectionApplicant')}</h2>
          <div className="form-grid">
            <Field label={t('fields.lastName')} required error={fl('LastName')}>
              <input
                className={`input ${fl('LastName') ? 'is-invalid' : ''}`}
                value={form.lastName}
                onChange={(e) => set('lastName', e.target.value)}
              />
            </Field>
            <Field label={t('fields.firstName')} required error={fl('FirstName')}>
              <input
                className={`input ${fl('FirstName') ? 'is-invalid' : ''}`}
                value={form.firstName}
                onChange={(e) => set('firstName', e.target.value)}
              />
            </Field>
            <Field label={t('fields.patronymic')} error={fl('Patronymic')}>
              <input
                className={`input ${fl('Patronymic') ? 'is-invalid' : ''}`}
                value={form.patronymic}
                onChange={(e) => set('patronymic', e.target.value)}
              />
            </Field>
            <Field
              label={t('fields.identificationNumber')}
              required
              error={fl('IdentificationNumber')}
            >
              <input
                className={`input ${fl('IdentificationNumber') ? 'is-invalid' : ''}`}
                value={form.identificationNumber}
                onChange={(e) => set('identificationNumber', e.target.value)}
              />
            </Field>
            {form.applicantCategoryCode === 'legal' ? (
              <Field
                label={t('fields.nameOrganization')}
                required
                error={fl('OrganizationName')}
              >
                <input
                  className={`input ${fl('OrganizationName') ? 'is-invalid' : ''}`}
                  value={form.organizationName}
                  onChange={(e) => set('organizationName', e.target.value)}
                />
              </Field>
            ) : null}
            <Field label={t('fields.physicalStatus')} error={fl('PhysicalStatusCode')}>
              <select
                className="select"
                value={form.physicalStatusCode}
                onChange={(e) => set('physicalStatusCode', e.target.value)}
              >
                <option value="">{t('new.notSelected')}</option>
                {dict.physical.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
            <Field label={t('fields.socialCategory')} error={fl('SocialCategoryCode')}>
              <select
                className="select"
                value={form.socialCategoryCode}
                onChange={(e) => set('socialCategoryCode', e.target.value)}
              >
                <option value="">{t('new.notSelected')}</option>
                {dict.social.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
            <Field label={t('fields.email')} required error={fl('Email')}>
              <input
                type="email"
                className={`input ${fl('Email') ? 'is-invalid' : ''}`}
                value={form.email}
                onChange={(e) => set('email', e.target.value)}
                placeholder="name@example.com"
              />
            </Field>
            <Field label={t('fields.phone')} required error={fl('Phone')}>
              <input
                className={`input ${fl('Phone') ? 'is-invalid' : ''}`}
                value={form.phone}
                onChange={(e) => set('phone', e.target.value)}
                placeholder="+37529112233"
              />
            </Field>
          </div>
        </div>

        <div className="form-section">
          <h2>{t('new.sectionAddress')}</h2>
          <div className="form-grid">
            <label className="checkbox span-2">
              <input
                type="checkbox"
                checked={form.isAbroad}
                onChange={(e) => set('isAbroad', e.target.checked)}
              />
              {t('fields.isAbroad')}
            </label>
            <Field label={t('fields.citizenship')} error={fl('CitizenshipCode')}>
              <select
                className="select"
                value={form.citizenshipCode}
                onChange={(e) => set('citizenshipCode', e.target.value)}
              >
                <option value="">{t('new.notSelected')}</option>
                {dict.citizenships.map((d) => (
                  <option key={d.code} value={d.code}>
                    {d.title}
                  </option>
                ))}
              </select>
            </Field>
            <Field label={t('fields.locality')} required error={fl('Locality')}>
              <input
                className={`input ${fl('Locality') ? 'is-invalid' : ''}`}
                value={form.locality}
                onChange={(e) => set('locality', e.target.value)}
              />
            </Field>
            <Field label={t('fields.street')} required error={fl('Street')}>
              <input
                className={`input ${fl('Street') ? 'is-invalid' : ''}`}
                value={form.street}
                onChange={(e) => set('street', e.target.value)}
              />
            </Field>
            <Field label={t('fields.house')} required error={fl('House')}>
              <input
                className={`input ${fl('House') ? 'is-invalid' : ''}`}
                value={form.house}
                onChange={(e) => set('house', e.target.value)}
              />
            </Field>
            <Field label={t('fields.apartment')} error={fl('Apartment')}>
              <input
                className={`input ${fl('Apartment') ? 'is-invalid' : ''}`}
                value={form.apartment}
                onChange={(e) => set('apartment', e.target.value)}
              />
            </Field>
            <Field label={t('fields.region')} error={fl('Region')}>
              <input
                className={`input ${fl('Region') ? 'is-invalid' : ''}`}
                value={form.region}
                onChange={(e) => set('region', e.target.value)}
              />
            </Field>
            <Field label={t('fields.district')} error={fl('District')}>
              <input
                className={`input ${fl('District') ? 'is-invalid' : ''}`}
                value={form.district}
                onChange={(e) => set('district', e.target.value)}
              />
            </Field>
            <Field label={t('fields.postalCode')} error={fl('PostalCode')}>
              <input
                className={`input ${fl('PostalCode') ? 'is-invalid' : ''}`}
                value={form.postalCode}
                onChange={(e) => set('postalCode', e.target.value)}
              />
            </Field>
            <Field label={t('fields.noteAddress')} error={fl('NoteAddress')}>
              <input
                className={`input ${fl('NoteAddress') ? 'is-invalid' : ''}`}
                value={form.noteAddress}
                onChange={(e) => set('noteAddress', e.target.value)}
              />
            </Field>
          </div>
        </div>

        <div className="form-section">
          <h2>{t('new.sectionContent')}</h2>
          <div className="form-grid">
            <Field
              label={t('fields.content')}
              required
              error={fl('Content')}
              className="span-2"
            >
              <textarea
                className={`textarea ${fl('Content') ? 'is-invalid' : ''}`}
                value={form.content}
                onChange={(e) => set('content', e.target.value)}
                placeholder={t('new.contentPlaceholder')}
                rows={6}
              />
            </Field>
            <Field
              label={t('fields.plannedDate')}
              error={fl('PlannedRegistrationDate')}
            >
              <input
                type="date"
                className={`input ${fl('PlannedRegistrationDate') ? 'is-invalid' : ''}`}
                value={form.plannedRegistrationDate}
                onChange={(e) => set('plannedRegistrationDate', e.target.value)}
              />
            </Field>
            <Field
              label={t('fields.attachments')}
              error={hasPickedErrors ? t('new.fixFiles') : undefined}
              className="span-2"
            >
              <div>
                <label className="file-picker">
                  <IconPaperclip />
                  {t('new.filePicker')}
                  <input
                    type="file"
                    multiple
                    accept=".doc,.docx,.txt,.rtf,.odf,.pdf,.jpg,.jpeg,.tiff,.png"
                    onChange={(e) => {
                      onFilesPicked(e.target.files)
                      e.target.value = ''
                    }}
                  />
                </label>
                {picked.map((item, index) => (
                  <div className="picked-file" key={`${item.file.name}-${index}`}>
                    <span className="name" style={item.error ? { color: 'var(--danger)' } : undefined}>
                      {item.file.name}
                    </span>
                    {item.error ? (
                      <span className="size" style={{ color: 'var(--danger)' }}>
                        {item.error}
                      </span>
                    ) : (
                      <span className="size">{Math.ceil(item.file.size / 1024)} КБ</span>
                    )}
                    <Button
                      variant="ghost"
                      className="btn-icon"
                      aria-label={t('files.remove')}
                      onClick={() =>
                        setPicked((prev) => prev.filter((_, i) => i !== index))
                      }
                    >
                      <IconTrash />
                    </Button>
                  </div>
                ))}
              </div>
            </Field>
          </div>
        </div>

        <div className="form-actions">
          <CaptchaField ref={captchaRef} onToken={setCaptcha} />
          <Button variant="ghost" onClick={onCancel}>
            {t('new.cancel')}
          </Button>
          <Button type="submit" disabled={busy || hasPickedErrors || (isCaptchaEnabled() && !captcha)}>
            {busy ? t('new.submitting') : t('new.submit')}
            <IconSend width={16} height={16} />
          </Button>
        </div>
      </form>
    </div>
  )
}
