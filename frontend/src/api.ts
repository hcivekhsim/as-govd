import type {
  AppealDetail,
  AppealListItem,
  DictionaryItem,
  AppealStatus,
  Page,
} from './types'
import { currentLang } from './i18n'

const API_BASE = '/api'

export class ApiError extends Error {
  status: number
  fields?: Record<string, string>

  constructor(status: number, message: string, fields?: Record<string, string>) {
    super(message)
    this.status = status
    this.fields = fields
  }
}

function getCookie(name: string): string {
  const match = document.cookie.match(new RegExp('(?:^|; )' + name + '=([^;]*)'))
  return match ? decodeURIComponent(match[1]) : ''
}

interface RequestOptions {
  method?: string
  json?: unknown
  form?: FormData
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const method = options.method ?? 'GET'
  const headers: Record<string, string> = {}

  if (options.json !== undefined) {
    headers['Content-Type'] = 'application/json'
  }
  const csrf = getCookie('csrf_token')
  if (csrf) {
    headers['X-CSRF-Token'] = csrf
  }

  headers['Accept-Language'] = currentLang()

  const response = await fetch(API_BASE + path, {
    method,
    headers,
    credentials: 'include',
    body: options.json !== undefined ? JSON.stringify(options.json) : options.form,
  })

  if (response.status === 204) {
    return undefined as T
  }

  let data: unknown = null
  const text = await response.text()
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = null
    }
  }

  if (!response.ok) {
    const payload = (data ?? {}) as { error?: unknown; errors?: unknown }
    let message = 'Не удалось выполнить запрос'

    if (typeof payload.error === 'string') {
      message = payload.error
    } else if (payload.errors && typeof payload.errors === 'object') {
      message = 'Проверьте заполнение полей'
    }

    let fields: Record<string, string> | undefined
    if (payload.errors && typeof payload.errors === 'object') {
      fields = {}
      for (const [key, value] of Object.entries(payload.errors as Record<string, unknown>)) {
        if (typeof value === 'string') {
          fields[key] = value
        }
      }
    }

    throw new ApiError(response.status, message, fields)
  }

  return data as T
}

export interface RegisterPayload {
  user_type: string
  name_organization?: string
  login: string
  password: string
  password_confirm: string
  email: string
  phone: string
  last_name: string
  first_name: string
  patronymic?: string
  is_abroad: boolean
  locality: string
  street: string
  house: string
  apartment?: string
  note_address?: string
  postal_code?: string
  captcha_token?: string
}

export const api = {
  checkSession: () =>
    request<{ message: string; language: 'ru' | 'be' }>('/protected'),

  login: (login: string, password: string, captchaToken: string) =>
    request<{ message: string; language: 'ru' | 'be' }>('/login', {
      method: 'POST',
      json: { login, password, captcha_token: captchaToken },
    }),

  logout: () => request<{ message: string }>('/logout', { method: 'POST' }),

  register: (data: RegisterPayload) =>
    request<{ id: number }>('/register', { method: 'POST', json: data }),

  setLanguage: (language: 'ru' | 'be') =>
    request<{ language: string }>('/preferences', {
      method: 'PUT',
      json: { language },
    }),

  appeals: (params: { page: number; pageSize: number; sort: string; statusCode: string }) => {
    const query = new URLSearchParams({
      page: String(params.page),
      page_size: String(params.pageSize),
      sort: params.sort,
    })
    if (params.statusCode) {
      query.set('status_code', params.statusCode)
    }
    return request<Page<AppealListItem>>(`/appeals?${query.toString()}`)
  },

  appeal: (id: number) => request<AppealDetail>(`/appeals/${id}`),

  dictionary: (kind: string) => request<{ items: DictionaryItem[] }>(`/dictionaries/${kind}`),

  statuses: () => request<{ items: AppealStatus[] }>('/dictionaries/appeal-statuses'),

  createAppeal: (form: FormData) => request<AppealDetail>('/appeals', { method: 'POST', form }),
}
