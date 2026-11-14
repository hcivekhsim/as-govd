package tasks

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
)

func AuthorizeUser(r *http.Request, users *UserRepo) (*UserContext, error) {
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		return nil, ErrUnauthorized
	}

	decodedSessionToken, err := url.QueryUnescape(st.Value)
	if err != nil {
		return nil, ErrUnauthorized
	}

	sessionTokenHash := sha256.Sum256([]byte(decodedSessionToken))
	sessionTokenHashHex := base64.URLEncoding.EncodeToString(sessionTokenHash[:])

	userAuth, err := users.getUserAuthBySessionTokenHash(r.Context(), sessionTokenHashHex)
	if err != nil {
		return nil, err
	}

	csrf := r.Header.Get("X-CSRF-Token")
	if csrf == "" {
		return nil, ErrUnauthorized
	}

	decodedCSRF, err := url.QueryUnescape(csrf)
	if err != nil {
		return nil, ErrUnauthorized
	}

	csrfTokenHash := sha256.Sum256([]byte(decodedCSRF))
	csrfTokenHashHex := base64.URLEncoding.EncodeToString(csrfTokenHash[:])
	if csrfTokenHashHex != userAuth.CSRFTokenHash {
		return nil, ErrUnauthorized
	}

	if time.Now().After(userAuth.ExpiresAt) {
		return nil, ErrUnauthorized
	}

	return &UserContext{
		UserID: userAuth.UserID,
		Login:  userAuth.Login,
		Role:   userAuth.Role,
	}, nil
}
