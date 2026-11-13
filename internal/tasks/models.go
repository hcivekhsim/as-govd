package tasks

import "time"

type Login struct {
	UserID   int64
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,gte=10,lte=20,password"`
}

type LoginAuth struct {
	UserID           int64
	Login            string
	SessionTokenHash string
	CSRFTokenHash    string
	ExpiresAt        time.Time
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
}
