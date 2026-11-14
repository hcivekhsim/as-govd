package tasks

import "time"

type Login struct {
	UserID       int64
	Login        string `json:"login" validate:"required"`
	Password     string `json:"password" validate:"required,gte=10,lte=20,password"`
	CaptchaToken string `json:"captcha_token"`
}

type LoginAuth struct {
	UserID            int64
	Login             string
	Role              string
	PreferredLanguage string
	SessionTokenHash  string
	CSRFTokenHash     string
	ExpiresAt         time.Time
}

type Register struct {
	UserType         string  `json:"user_type" validate:"required,oneof=individual legal foreign stateless"`
	NameOrganization *string `json:"name_organization" validate:"required_if=UserType legal"` // Юридическое лицо = обязательное поле
	Login            string  `json:"login" validate:"required"`                               // validator login
	Password         string  `json:"password" validate:"required,gte=10,lte=20,password"`     // HASH - services
	PasswordConfirm  string  `json:"password_confirm" validate:"required,eqfield=Password"`   // p = pc in model
	Email            string  `json:"email" validate:"required,email"`                         // validator email
	Phone            string  `json:"phone" validate:"required,phone"`                         // validator phone
	LastName         string  `json:"last_name" validate:"required"`                           // validator LastName
	FirstName        string  `json:"first_name" validate:"required"`                          // validator FirstName
	Patronymic       *string `json:"patronymic"`                                              // validator Patronymic
	IsAbroad         bool    `json:"is_abroad"`
	Locality         string  `json:"locality" validate:"required"`
	Street           string  `json:"street" validate:"required"`
	House            string  `json:"house" validate:"required"`
	Apartment        *string `json:"apartment"`
	NoteAddress      *string `json:"note_address"`
	PostalCode       *string `json:"postal_code"`
	CaptchaToken     string  `json:"captcha_token"`
}

type UserContext struct {
	UserID int64
	Login  string
	Role   string
}

type DictionaryItem struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type StatusGroup struct {
	Code      string `json:"code"`
	Title     string `json:"title"`
	BadgeType string `json:"badge_type"`
}

type AppealStatus struct {
	Code  string      `json:"code"`
	Title string      `json:"title"`
	Group StatusGroup `json:"group"`
}

type AppealFileResponse struct {
	ID          int64  `json:"id"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
	DownloadURL string `json:"download_url"`
}

type AppealStatusHistoryResponse struct {
	ID         int64         `json:"id"`
	FromStatus *AppealStatus `json:"from_status,omitempty"`
	ToStatus   AppealStatus  `json:"to_status"`
	Comment    *string       `json:"comment,omitempty"`
	ChangedBy  int64         `json:"changed_by"`
	ChangedAt  time.Time     `json:"changed_at"`
}

type AppealListItemResponse struct {
	ID                      int64           `json:"id"`
	Number                  string          `json:"number"`
	CreatedAt               time.Time       `json:"created_at"`
	PlannedRegistrationDate *time.Time      `json:"planned_registration_date,omitempty"`
	ApplicantCategory       DictionaryItem  `json:"applicant_category"`
	Topic                   DictionaryItem  `json:"topic"`
	AppealKind              DictionaryItem  `json:"appeal_kind"`
	AppealType              DictionaryItem  `json:"appeal_type"`
	ApplicantName           string          `json:"applicant_name"`
	IdentificationNumber    string          `json:"identification_number"`
	Citizenship             *DictionaryItem `json:"citizenship,omitempty"`
	Status                  AppealStatus    `json:"status"`
}

type AppealDetailResponse struct {
	ID                       int64                         `json:"id"`
	Number                   string                        `json:"number"`
	CreatedAt                time.Time                     `json:"created_at"`
	UpdatedAt                time.Time                     `json:"updated_at"`
	SubmittedByUserID        int64                         `json:"submitted_by_user_id"`
	ApplicantCategory        DictionaryItem                `json:"applicant_category"`
	Topic                    DictionaryItem                `json:"topic"`
	OrganizationName         *string                       `json:"organization_name,omitempty"`
	AppealKind               DictionaryItem                `json:"appeal_kind"`
	AppealType               DictionaryItem                `json:"appeal_type"`
	LastName                 string                        `json:"last_name"`
	FirstName                string                        `json:"first_name"`
	Patronymic               *string                       `json:"patronymic,omitempty"`
	IdentificationNumber     string                        `json:"identification_number"`
	PhysicalStatus           *DictionaryItem               `json:"physical_status,omitempty"`
	SocialCategory           *DictionaryItem               `json:"social_category,omitempty"`
	Email                    string                        `json:"email"`
	IsAbroad                 bool                          `json:"is_abroad"`
	Citizenship              *DictionaryItem               `json:"citizenship,omitempty"`
	Region                   *string                       `json:"region,omitempty"`
	District                 *string                       `json:"district,omitempty"`
	Locality                 string                        `json:"locality"`
	Street                   string                        `json:"street"`
	House                    string                        `json:"house"`
	Apartment                *string                       `json:"apartment,omitempty"`
	NoteAddress              *string                       `json:"note_address,omitempty"`
	PostalCode               *string                       `json:"postal_code,omitempty"`
	Phone                    string                        `json:"phone"`
	Content                  string                        `json:"content"`
	PlannedRegistrationDate  *time.Time                    `json:"planned_registration_date,omitempty"`
	Status                   AppealStatus                  `json:"status"`
	AllowedStatusTransitions []AppealStatus                `json:"allowed_status_transitions,omitempty"`
	Files                    []AppealFileResponse          `json:"files"`
	StatusHistory            []AppealStatusHistoryResponse `json:"status_history"`
}

type CreateAppealInput struct {
	ApplicantCategoryCode   string     `json:"applicant_category_code" validate:"required"`
	TopicCode               string     `json:"topic_code" validate:"required"`
	OrganizationName        *string    `json:"organization_name"`
	AppealKindCode          string     `json:"appeal_kind_code" validate:"required"`
	AppealTypeCode          string     `json:"appeal_type_code" validate:"required"`
	LastName                string     `json:"last_name" validate:"required"`
	FirstName               string     `json:"first_name" validate:"required"`
	Patronymic              *string    `json:"patronymic"`
	IdentificationNumber    string     `json:"identification_number" validate:"required"`
	PhysicalStatusCode      *string    `json:"physical_status_code"`
	SocialCategoryCode      *string    `json:"social_category_code"`
	Email                   string     `json:"email" validate:"required,email"`
	IsAbroad                bool       `json:"is_abroad"`
	CitizenshipCode         *string    `json:"citizenship_code"`
	Region                  *string    `json:"region"`
	District                *string    `json:"district"`
	Locality                string     `json:"locality" validate:"required"`
	Street                  string     `json:"street" validate:"required"`
	House                   string     `json:"house" validate:"required"`
	Apartment               *string    `json:"apartment"`
	NoteAddress             *string    `json:"note_address"`
	PostalCode              *string    `json:"postal_code"`
	Phone                   string     `json:"phone" validate:"required,phone"`
	Content                 string     `json:"content" validate:"required"`
	PlannedRegistrationDate *time.Time `json:"planned_registration_date"`
}

type ListMeta struct {
	Items    any   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

type AppealListFilter struct {
	Page                  int
	PageSize              int
	Sort                  string
	CitizenshipCode       *string
	Region                *string
	District              *string
	Locality              *string
	Street                *string
	House                 *string
	PostalCode            *string
	Apartment             *string
	IdentificationNumber  *string
	ApplicantName         *string
	TopicCode             *string
	PlannedDateFrom       *time.Time
	PlannedDateTo         *time.Time
	ApplicantCategoryCode *string
	AppealKindCode        *string
	AppealTypeCode        *string
	StatusCode            *string
	StatusGroupCode       *string
	CreatedAtFrom         *time.Time
	CreatedAtTo           *time.Time
	Number                *string
}
