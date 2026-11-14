import { useEffect, useState } from 'react'
import { api, ApiError } from './api'
import type { AppealDetail } from './types'
import { useLang } from './i18n'
import {
  Alert,
  Loading,
  StatusBadge,
  formatDate,
  formatDateTime,
  formatSize,
} from './components'
import { IconAlertCircle, IconArrowLeft, IconChevronRight, IconDownload, IconFileText } from './icons'

interface AppealDetailPageProps {
  appealId: number
  onBack: () => void
}

export function AppealDetailPage({ appealId, onBack }: AppealDetailPageProps) {
  const { t } = useLang()
  const [detail, setDetail] = useState<AppealDetail | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    api
      .appeal(appealId)
      .then((res) => {
        if (!cancelled) setDetail(res)
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof ApiError
              ? err.status === 403
                ? t('detail.foreignAppeal')
                : err.message
              : t('detail.loadFailed'),
          )
        }
      })
    return () => {
      cancelled = true
    }
  }, [appealId, t])

  if (error) {
    return (
      <div className="container-narrow">
        <button type="button" className="back-link" onClick={onBack}>
          <IconArrowLeft width={16} height={16} />
          {t('detail.back')}
        </button>
        <Alert>{error}</Alert>
      </div>
    )
  }

  if (!detail) {
    return (
      <div className="container-narrow">
        <Loading text={t('detail.loading')} />
      </div>
    )
  }

  const applicantName = [detail.last_name, detail.first_name, detail.patronymic]
    .filter(Boolean)
    .join(' ')

  const address = [
    detail.locality && t('detail.address.locality', { value: detail.locality }),
    detail.street && t('detail.address.street', { value: detail.street }),
    detail.house && t('detail.address.house', { value: detail.house }),
    detail.apartment && t('detail.address.apartment', { value: detail.apartment }),
  ]
    .filter(Boolean)
    .join(', ')

  const meta: Array<{ label: string; value: string; mono?: boolean }> = [
    { label: t('detail.metaApplicantCategory'), value: detail.applicant_category.title },
    { label: t('detail.metaTopic'), value: detail.topic.title },
    { label: t('detail.metaKind'), value: detail.appeal_kind.title },
    { label: t('detail.metaType'), value: detail.appeal_type.title },
    { label: t('detail.metaApplicant'), value: applicantName },
    { label: t('detail.metaIdentificationNumber'), value: detail.identification_number, mono: true },
    { label: t('detail.metaEmail'), value: detail.email },
    { label: t('detail.metaPhone'), value: detail.phone, mono: true },
    { label: t('detail.metaAddress'), value: address },
  ]
  if (detail.organization_name) {
    meta.splice(1, 0, { label: t('detail.metaOrganization'), value: detail.organization_name })
  }
  if (detail.physical_status) {
    meta.push({ label: t('detail.metaPhysicalStatus'), value: detail.physical_status.title })
  }
  if (detail.social_category) {
    meta.push({ label: t('detail.metaSocialCategory'), value: detail.social_category.title })
  }
  if (detail.citizenship) {
    meta.push({ label: t('detail.metaCitizenship'), value: detail.citizenship.title })
  }
  if (detail.is_abroad) {
    meta.push({ label: t('detail.metaLocation'), value: t('detail.locationAbroad') })
  }
  if (detail.postal_code) {
    meta.push({ label: t('detail.metaPostalCode'), value: detail.postal_code, mono: true })
  }
  if (detail.planned_registration_date) {
    meta.push({
      label: t('detail.metaPlannedDate'),
      value: formatDate(detail.planned_registration_date),
    })
  }

  return (
    <div className="container-narrow">
      <button type="button" className="back-link" onClick={onBack}>
        <IconArrowLeft width={16} height={16} />
        {t('detail.back')}
      </button>

      <div className="detail-head">
        <div>
          <h1>
            <span className="mono">{detail.number}</span>
            <StatusBadge status={detail.status} />
          </h1>
          <p className="detail-date">
            {t('detail.submittedAt')} {formatDateTime(detail.created_at)} · {t('detail.updatedAt')}{' '}
            {formatDateTime(detail.updated_at)}
          </p>
        </div>
      </div>

      <div className="card detail-block">
        <div className="meta-grid">
          {meta.map((item) => (
            <div className="meta-item" key={item.label}>
              <dt>{item.label}</dt>
              <dd className={item.mono ? 'mono' : undefined}>{item.value}</dd>
            </div>
          ))}
        </div>
      </div>

      <div className="card detail-block">
        <h2>{t('detail.contentTitle')}</h2>
        <p className="detail-content">{detail.content}</p>
      </div>

      <div className="card detail-block">
        <h2>
          {t('detail.attachmentsTitle')} ({detail.files.length})
        </h2>
        {detail.files.length === 0 ? (
          <p className="detail-date">{t('detail.noFiles')}</p>
        ) : (
          detail.files.map((file) => (
            <div className="file-row" key={file.id}>
              <span className="file-icon">
                <IconFileText />
              </span>
              <span className="file-name">{file.file_name}</span>
              <span className="file-size mono">{formatSize(file.size_bytes)}</span>
              <a
                className="btn btn-ghost btn-icon"
                href={file.download_url}
                aria-label={t('detail.downloadFile', { name: file.file_name })}
              >
                <IconDownload />
              </a>
            </div>
          ))
        )}
      </div>

      <div className="card detail-block">
        <h2>{t('detail.historyTitle')}</h2>
        <ul className="timeline">
          {detail.status_history.map((item) => (
            <li className="timeline-item" key={item.id}>
              <span className="timeline-dot" />
              <div className="timeline-head">
                {item.from_status ? (
                  <>
                    <span style={{ color: 'var(--ink-2)', fontSize: 14 }}>
                      {item.from_status.title}
                    </span>
                    <span className="arrow">
                      <IconChevronRight width={14} height={14} />
                    </span>
                  </>
                ) : null}
                <StatusBadge status={item.to_status} />
                <span className="timeline-date mono">{formatDateTime(item.changed_at)}</span>
              </div>
              {item.comment ? <p className="timeline-comment">{item.comment}</p> : null}
            </li>
          ))}
        </ul>
      </div>

      {detail.status.group.code === 'final_negative' ? (
        <div className="alert">
          <span className="alert-icon">
            <IconAlertCircle />
          </span>
          <div>{t('detail.refusedAlert')}</div>
        </div>
      ) : null}
    </div>
  )
}
