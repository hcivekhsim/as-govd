export interface DictionaryItem {
  code: string
  title: string
}

export interface StatusGroup {
  code: string
  title: string
  badge_type: string
}

export interface AppealStatus {
  code: string
  title: string
  group: StatusGroup
}

export interface AppealFile {
  id: number
  file_name: string
  content_type: string
  size_bytes: number
  download_url: string
}

export interface StatusHistoryItem {
  id: number
  from_status?: AppealStatus | null
  to_status: AppealStatus
  comment?: string | null
  changed_by: number
  changed_at: string
}

export interface AppealDetail {
  id: number
  number: string
  created_at: string
  updated_at: string
  submitted_by_user_id: number
  applicant_category: DictionaryItem
  topic: DictionaryItem
  organization_name?: string | null
  appeal_kind: DictionaryItem
  appeal_type: DictionaryItem
  last_name: string
  first_name: string
  patronymic?: string | null
  identification_number: string
  physical_status?: DictionaryItem | null
  social_category?: DictionaryItem | null
  email: string
  is_abroad: boolean
  citizenship?: DictionaryItem | null
  region?: string | null
  district?: string | null
  locality: string
  street: string
  house: string
  apartment?: string | null
  note_address?: string | null
  postal_code?: string | null
  phone: string
  content: string
  planned_registration_date?: string | null
  status: AppealStatus
  allowed_status_transitions?: AppealStatus[]
  files: AppealFile[]
  status_history: StatusHistoryItem[]
}

export interface AppealListItem {
  id: number
  number: string
  created_at: string
  planned_registration_date?: string | null
  applicant_category: DictionaryItem
  topic: DictionaryItem
  appeal_kind: DictionaryItem
  appeal_type: DictionaryItem
  applicant_name: string
  identification_number: string
  citizenship?: DictionaryItem | null
  status: AppealStatus
}

export interface Page<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
