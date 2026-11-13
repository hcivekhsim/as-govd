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
	outErrs := make(map[string]string, len(errs))
	for f, t := range errs {
		switch t {
		case "required":
			outErrs[f] = ErrRequiredField.Error()
		case "email":
			outErrs[f] = ErrInvalidField.Error()
		case "eqfield":
			outErrs[f] = ErrEqualField.Error()
		case "phone":
			outErrs[f] = ErrInvalidField.Error()
		case "password":
			outErrs[f] = ErrPasswordField.Error()
		case "gte":
			outErrs[f] = ErrGteField.Error()
		case "lte":
			outErrs[f] = ErrLteField.Error()
		case "oneof":
			outErrs[f] = ErrOneOfField.Error()
		default:
			outErrs[f] = ErrInvalidField.Error()
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
