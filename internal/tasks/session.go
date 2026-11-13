package tasks

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
)

func Authorize(r *http.Request, users *UserRepo) error {
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		// ошибка декодирования cookie сеанса
		return ErrNotFound
	}

	decodedSessionToken, err := url.QueryUnescape(st.Value)
	if err != nil {
		// ошибка декодирования cookie сеанса
		return ErrNotFound
	}

	sessionTokenHash := sha256.Sum256([]byte(decodedSessionToken))
	sessionTokenHashHex := base64.URLEncoding.EncodeToString(sessionTokenHash[:])

	userAuth, err := users.getUserAuthBySessionTokenHash(r.Context(), sessionTokenHashHex)
	if err != nil {
		// ошибка поиска сеанса
		return err
	}

	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" {
		// csrf заголовок отсутствует
		return ErrNotFound
	}

	decodedCSRF, err := url.QueryUnescape(csrf)
	if err != nil {
		// csrf ошибка декодирования
		return ErrNotFound
	}

	csrfTokenHash := sha256.Sum256([]byte(decodedCSRF))
	csrfTokenHashHex := base64.URLEncoding.EncodeToString(csrfTokenHash[:])
	if csrfTokenHashHex != userAuth.CSRFTokenHash {
		// несоответствие хеша csrf
		return ErrNotFound
	}

	if time.Now().After(userAuth.ExpiresAt) {
		// сессия истекла
		return ErrNotFound
	}

	return nil
}
