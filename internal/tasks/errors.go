package tasks

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// required, oneof, required_if, gte, lte, password, eqfield, email, phone

var (
	//Register
	ErrRequiredField           = errors.New("Поле обязательное")
	ErrOneOfField              = errors.New("Недопустимое значение")
	ErrRequiredIfField         = errors.New("Обязательно для юридического лица")
	ErrGteField                = errors.New("Слишком короткое значение")
	ErrLteField                = errors.New("Слишком длинное значение")
	ErrPasswordField           = errors.New("Нужны заглавная, строчная буква и цифра")
	ErrEqualField              = errors.New("Поля не совпадают")
	ErrInvalidField            = errors.New("Некорректное значение")
	ErrNotFound                = errors.New("Ничего не найдено")
	ErrPassowrdMismatch        = errors.New("Пароли не совпадают")
	ErrUnauthorized            = errors.New("Требуется авторизация")
	ErrForbidden               = errors.New("Недостаточно прав")
	ErrFileTooLarge            = errors.New("Файл превышает допустимый размер")
	ErrUnsupportedFileType     = errors.New("Недопустимый тип файла")
	ErrRequestTooLarge         = errors.New("Размер запроса превышает допустимый лимит")
	ErrInvalidStatusTransition = errors.New("Недопустимый переход статуса")

	// hCaptcha
	ErrCaptchaRequired = errors.New("Подтвердите, что вы не робот")
	ErrCaptchaFailed   = errors.New("Проверка не пройдена, попробуйте ещё раз")
)

func dbError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("PG ERROR: %s, %w", pgErr.Code, err)
	}
	return err
}
