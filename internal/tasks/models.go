package tasks

type Login struct{
	Username string `json:"login"`
	Password string `json:"password"`
}

type Register struct {
	UserType         string  `json:"user_type"`
	NameOrganization *string `json:"name_organization"`
	Login            string  `json:"login"`
	Password         string  `json:"password"`
	PasswordConfirm  string  `json:"password_confirm"` 
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	LastName         string  `json:"last_name"`
	FirstName        string  `json:"first_name"`
	Patronymic       *string `json:"patronymic"`
	IsAbroad         bool    `json:"is_abroad"`
	Locality         string  `json:"locality"`
	Street           string  `json:"street"`
	House            string  `json:"house"`
	Apartment        *string `json:"apartment"`
	NoteAddress      *string `json:"note_address"`
	PostalCode       *string `json:"postal_code"`
}