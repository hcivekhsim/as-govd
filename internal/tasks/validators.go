package tasks

import (
	"errors"
	"regexp"
	"unicode"

	validator "github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	phoneReg = regexp.MustCompile(`^\+[0-9]{10,15}$`)
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("phone", ValidatePhone)
	validate.RegisterValidation("password", ValidatePassword)
}

func ValidateStruct(data any) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var valerrs validator.ValidationErrors
	kverr := make(map[string]string, len(valerrs))
	if errors.As(err, &valerrs) {
		for _, v := range valerrs {
			kverr[v.Field()] = v.Tag()
		}
		return kverr
	}

	return map[string]string{"error": "internal validation error"}
}

// ValidateStruct -> SendReportOfStructAuth
func SendReportOfStructAuth(errs map[string]string) map[string]string {
	return SendValidationErrors(errs, LangRu)
}

// SendValidationErrors локализует отчёт валидации: ключ — имя поля (техническое),
// значение — понятное сообщение на запрошенном языке.
func SendValidationErrors(errs map[string]string, lang Lang) map[string]string {
	outErrs := make(map[string]string, len(errs))
	for f, t := range errs {
		switch t {
		case "required":
			outErrs[f] = validationTagMessage(lang, "required")
		case "required_if":
			outErrs[f] = validationTagMessage(lang, "required_if")
		case "email", "phone":
			outErrs[f] = validationTagMessage(lang, t)
		case "eqfield":
			outErrs[f] = validationTagMessage(lang, "eqfield")
		case "password":
			outErrs[f] = validationTagMessage(lang, "password")
		case "gte":
			outErrs[f] = validationTagMessage(lang, "gte")
		case "lte":
			outErrs[f] = validationTagMessage(lang, "lte")
		case "oneof":
			outErrs[f] = validationTagMessage(lang, "oneof")
		default:
			outErrs[f] = validationTagMessage(lang, "invalid")
		}
	}
	return outErrs
}

// Validators

func ValidatePhone(fl validator.FieldLevel) bool {
	return phoneReg.MatchString(fl.Field().String())
}

// минимум 1 - up, 2 - low, 3 - digit длина от 10 до 20 (gte=10,lte=20)
func ValidatePassword(fl validator.FieldLevel) bool {
	var digit, upper, lower bool
	for _, r := range fl.Field().String() {
		switch {
		case unicode.IsDigit(r):
			digit = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		}
	}
	return digit && upper && lower
}
