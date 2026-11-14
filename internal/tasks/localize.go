package tasks

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// Поддерживаемые языки интерфейса.
const (
	LangRu Lang = "ru"
	LangBe Lang = "be"
)

type Lang string

func (l Lang) valid() bool {
	return l == LangRu || l == LangBe
}

// LangFromRequest определяет язык ответа: сначала параметр ?lang=,
// затем заголовок Accept-Language (префикс be -> белорусский).
func LangFromRequest(c *gin.Context) Lang {
	if v := strings.TrimSpace(c.Query("lang")); v != "" {
		return normalizeLang(v)
	}
	if al := c.GetHeader("Accept-Language"); al != "" {
		// Accept-Language: be,ru;q=0.9,en;q=0.8 — берём первый тег
		first := strings.SplitN(al, ",", 2)[0]
		return normalizeLang(first)
	}
	return LangRu
}

func normalizeLang(tag string) Lang {
	tag = strings.TrimSpace(strings.ToLower(tag))
	base := strings.SplitN(tag, "-", 2)[0]
	if base == string(LangBe) {
		return LangBe
	}
	return LangRu
}

// localizedMessage — одна строка перевода для sentinel-ошибки.
type localizedMessage struct {
	target error
	ru     string
	be     string
}

// errorMessages — переводы доменных ошибок. Порядок не важен:
// сопоставление идёт через errors.Is.
var errorMessages = []localizedMessage{
	{ErrRequiredField, "Поле обязательное", "Поле абавязковае"},
	{ErrOneOfField, "Недопустимое значение", "Недапушальнае значэнне"},
	{ErrRequiredIfField, "Обязательно для юридического лица", "Абавязкова для юрыдычнай асобы"},
	{ErrGteField, "Слишком короткое значение", "Занадта кароткае значэнне"},
	{ErrLteField, "Слишком длинное значение", "Занадта доўгае значэнне"},
	{ErrPasswordField, "Нужны заглавная, строчная буква и цифра", "Патрэбны вялікая, малая літара і лічба"},
	{ErrEqualField, "Поля не совпадают", "Палі не супадаюць"},
	{ErrInvalidField, "Некорректное значение", "Некарэктнае значэнне"},
	{ErrNotFound, "Ничего не найдено", "Нічога не знойдзена"},
	{ErrPassowrdMismatch, "Пароли не совпадают", "Паролі не супадаюць"},
	{ErrUnauthorized, "Требуется авторизация", "Патрабуецца аўтарызацыя"},
	{ErrForbidden, "Недостаточно прав", "Недастаткова правоў"},
	{ErrFileTooLarge, "Файл превышает допустимый размер", "Файл перавышае дапушчальны памер"},
	{ErrUnsupportedFileType, "Недопустимый тип файла", "Недапушчальны тып файла"},
	{ErrRequestTooLarge, "Размер запроса превышает допустимый лимит", "Памер запыту перавышае дапушчальны ліміт"},
	{ErrInvalidStatusTransition, "Недопустимый переход статуса", "Недапушчальны пераход статусу"},
	{ErrCaptchaRequired, "Подтвердите, что вы не робот", "Пацвердзіце, што вы не робат"},
	{ErrCaptchaFailed, "Проверка не пройдена, попробуйте ещё раз", "Праверка не пройдзена, паспрабуйце яшчэ раз"},
}

// localizeError возвращает текст ошибки на запрошенном языке.
// Для неизвестных ошибок возвращается исходный текст.
func localizeError(lang Lang, err error) string {
	if err == nil {
		return ""
	}
	for _, m := range errorMessages {
		if errors.Is(err, m.target) {
			if lang == LangBe {
				return m.be
			}
			return m.ru
		}
	}
	return err.Error()
}

// validationTagMessages — переводы сообщений валидации по тегам validator.
var validationTagMessages = map[string][2]string{
	"required":    {"Поле обязательное", "Поле абавязковае"},
	"required_if": {"Обязательное поле для выбранного типа", "Абавязковае поле для абранага тыпу"},
	"email":       {"Некорректный email", "Некарэктны email"},
	"phone":       {"Формат +375XXXXXXXXX", "Фармат +375XXXXXXXXX"},
	"password":    {"Нужны заглавная, строчная буква и цифра", "Патрэбны вялікая, малая літара і лічба"},
	"gte":         {"Слишком короткое значение", "Занадта кароткае значэнне"},
	"lte":         {"Слишком длинное значение", "Занадта доўгае значэнне"},
	"eqfield":     {"Поля не совпадают", "Палі не супадаюць"},
	"oneof":       {"Недопустимое значение", "Недапушчальнае значэнне"},
	"invalid":     {"Некорректное значение", "Некарэктнае значэнне"},
}

func validationTagMessage(lang Lang, tag string) string {
	msg, ok := validationTagMessages[tag]
	if !ok {
		return tag
	}
	if lang == LangBe {
		return msg[1]
	}
	return msg[0]
}
