import type { ButtonHTMLAttributes, ReactNode } from 'react'
import type { AppealStatus } from './types'
import { currentLang, useLang } from './i18n'
import {
  IconAlertCircle,
  IconChevronLeft,
  IconChevronRight,
  IconInbox,
} from './icons'

type ButtonVariant = 'primary' | 'secondary' | 'ghost'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
}

export function Button({ variant = 'primary', className, ...props }: ButtonProps) {
  const classes = ['btn', `btn-${variant}`, className].filter(Boolean).join(' ')
  return <button type="button" className={classes} {...props} />
}


interface FieldProps {
  label: string
  required?: boolean
  error?: string
  children: ReactNode
  className?: string
}

export function Field({ label, required, error, children, className }: FieldProps) {
  return (
    <div className={className ? `field ${className}` : 'field'}>
      <label className="field-label">
        {label}
        {required && <span className="req">*</span>}
      </label>
      {children}
      {error ? <span className="field-error">{error}</span> : null}
    </div>
  )
}

const badgeTypes: Record<string, string> = {
  info: 'badge-info',
  progress: 'badge-progress',
  success: 'badge-success',
  warning: 'badge-warning',
}

export function StatusBadge({ status }: { status: AppealStatus }) {
  const cls = badgeTypes[status.group?.badge_type] ?? 'badge-info'
  return (
    <span className={`badge ${cls}`}>
      <span className="dot" />
      {status.title}
    </span>
  )
}


export function Alert({ children }: { children: ReactNode }) {
  return (
    <div className="alert" role="alert">
      <span className="alert-icon">
        <IconAlertCircle />
      </span>
      <div>{children}</div>
    </div>
  )
}


interface PaginationProps {
  page: number
  pageSize: number
  total: number
  onPage: (page: number) => void
}

export function Pagination({ page, pageSize, total, onPage }: PaginationProps) {
  const { t } = useLang()
  const pages = Math.max(1, Math.ceil(total / pageSize))
  const from = total === 0 ? 0 : (page - 1) * pageSize + 1
  const to = Math.min(total, page * pageSize)

  return (
    <div className="pagination">
      <span className="pagination-info">
        {from}–{to} {t('list.of')} {total}
      </span>
      <div className="pagination-controls">
        <Button
          variant="ghost"
          className="btn-icon"
          disabled={page <= 1}
          onClick={() => onPage(page - 1)}
          aria-label={t('list.prevPage')}
        >
          <IconChevronLeft />
        </Button>
        <Button
          variant="ghost"
          className="btn-icon"
          disabled={page >= pages}
          onClick={() => onPage(page + 1)}
          aria-label={t('list.nextPage')}
        >
          <IconChevronRight />
        </Button>
      </div>
    </div>
  )
}

export function EmptyState({ title, note }: { title: string; note?: string }) {
  return (
    <div className="empty-state">
      <span className="empty-icon">
        <IconInbox width={28} height={28} />
      </span>
      <span className="empty-title">{title}</span>
      {note && <span>{note}</span>}
    </div>
  )
}

export function Loading({ text = 'Загрузка' }: { text?: string }) {
  return (
    <div className="loading">
      <span className="spinner" />
      {text}…
    </div>
  )
}

const DATE_LOCALES: Record<string, string> = {
  ru: 'ru-RU',
  be: 'be-BY',
}

function localeTag(): string {
  return DATE_LOCALES[currentLang()] ?? 'ru-RU'
}

function makeDateFmt(): Intl.DateTimeFormat {
  return new Intl.DateTimeFormat(localeTag(), {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  })
}

function makeDateTimeFmt(): Intl.DateTimeFormat {
  return new Intl.DateTimeFormat(localeTag(), {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const fmtCache = new Map<string, Intl.DateTimeFormat>()

function fmt(kind: 'date' | 'datetime'): Intl.DateTimeFormat {
  const key = `${kind}:${localeTag()}`
  let fmtObj = fmtCache.get(key)
  if (!fmtObj) {
    fmtObj = kind === 'date' ? makeDateFmt() : makeDateTimeFmt()
    fmtCache.set(key, fmtObj)
  }
  return fmtObj
}

export function formatDate(value: string | null | undefined): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return fmt('date').format(date)
}

export function formatDateTime(value: string | null | undefined): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return fmt('datetime').format(date)
}

export function formatSize(bytes: number): string {
  if (bytes >= 1024 * 1024) {
    return `${(bytes / (1024 * 1024)).toFixed(1)} МБ`
  }
  if (bytes >= 1024) {
    return `${Math.round(bytes / 1024)} КБ`
  }
  return `${bytes} Б`
}
