import { useEffect, useState } from 'react'
import { api, ApiError } from './api'
import type { AppealListItem, AppealStatus, Page } from './types'
import { useLang } from './i18n'
import {
  Alert,
  EmptyState,
  Loading,
  Pagination,
  StatusBadge,
  formatDate,
} from './components'

interface MyAppealsPageProps {
  onOpenAppeal: (id: number) => void
  onNewAppeal: () => void
}

export function MyAppealsPage({ onOpenAppeal, onNewAppeal }: MyAppealsPageProps) {
  const { t } = useLang()
  const [data, setData] = useState<Page<AppealListItem> | null>(null)
  const [statuses, setStatuses] = useState<AppealStatus[]>([])
  const [statusCode, setStatusCode] = useState('')
  const [sort, setSort] = useState('created_at_desc')
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .statuses()
      .then((res) => setStatuses(res.items))
      .catch(() => {
      })
  }, [])

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError(null)
    api
      .appeals({ page, pageSize: 10, sort, statusCode })
      .then((res) => {
        if (!cancelled) setData(res)
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err instanceof ApiError ? err.message : t('list.loadFailed'))
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [page, sort, statusCode])

  return (
    <div className="container">
      <div className="page-head">
        <div>
          <h1>{t('list.title')}</h1>
          <p className="page-sub">{t('list.subtitle')}</p>
        </div>
        <button type="button" className="btn btn-primary" onClick={onNewAppeal}>
          {t('app.nav.newAppeal')}
        </button>
      </div>

      <div className="list-filters">
        <select
          className="select"
          value={statusCode}
          onChange={(e) => {
            setStatusCode(e.target.value)
            setPage(1)
          }}
          aria-label={t('list.filterByStatus')}
        >
          <option value="">{t('list.allStatuses')}</option>
          {statuses.map((s) => (
            <option key={s.code} value={s.code}>
              {s.title}
            </option>
          ))}
        </select>
        <span className="spacer" />
        <select
          className="select"
          value={sort}
          onChange={(e) => setSort(e.target.value)}
          aria-label={t('list.sortBy')}
        >
          <option value="created_at_desc">{t('list.sortNewFirst')}</option>
          <option value="created_at_asc">{t('list.sortOldFirst')}</option>
        </select>
      </div>

      {error ? <Alert>{error}</Alert> : null}
      {loading && !data ? <Loading text={t('list.loading')} /> : null}

      {data &&
        (data.items.length === 0 ? (
          <EmptyState title={t('list.emptyTitle')} note={t('list.emptyNote')} />
        ) : (
          <>
            {data.items.map((item) => (
              <button
                key={item.id}
                type="button"
                className="appeal-row"
                onClick={() => onOpenAppeal(item.id)}
              >
                <span className="appeal-row-number">
                  <span className="mono">{item.number}</span>
                  <span className="date">{formatDate(item.created_at)}</span>
                </span>
                <span className="appeal-row-topic">
                  <span className="title">{item.topic.title}</span>
                  <span className="meta">
                    {item.applicant_category.title} · {item.appeal_kind.title} ·{' '}
                    {item.applicant_name}
                  </span>
                </span>
                <span className="appeal-row-status">
                  <StatusBadge status={item.status} />
                </span>
              </button>
            ))}
            <Pagination
              page={data.page}
              pageSize={data.page_size}
              total={data.total}
              onPage={setPage}
            />
          </>
        ))}
    </div>
  )
}
