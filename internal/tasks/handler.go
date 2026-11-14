package tasks

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (ah *AuthHandler) RegisterHandler(c *gin.Context) {
	var user Register
	if err := c.ShouldBindJSON(&user); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if err := VerifyCaptchaToken(c.Request.Context(), user.CaptchaToken, c.ClientIP()); err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	found, err := ah.service.users.getByLogin(c.Request.Context(), user.Login)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}
	if found {
		respondValidation(c, map[string]string{"Login": ErrInvalidField.Error()})
		return
	}

	if errs := ValidateStruct(user); errs != nil {
		respondValidation(c, errs)
		return
	}

	id, err := ah.service.RegisterService(c.Request.Context(), user)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	userID, sessionToken, csrfToken, err := ah.service.LoginService(c.Request.Context(), Login{
		Login:        user.Login,
		Password:     user.Password,
		CaptchaToken: user.CaptchaToken,
	})
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	sessionTokenHash := sha256.Sum256([]byte(sessionToken))
	sessionTokenHashHex := base64.URLEncoding.EncodeToString(sessionTokenHash[:])

	csrfTokenHash := sha256.Sum256([]byte(csrfToken))
	csrfTokenHashHex := base64.URLEncoding.EncodeToString(csrfTokenHash[:])

	userAuth := LoginAuth{
		UserID:           userID,
		SessionTokenHash: sessionTokenHashHex,
		CSRFTokenHash:    csrfTokenHashHex,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	err = ah.service.users.createUserAuth(c.Request.Context(), userAuth)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", sessionToken, 24*3600, "/", "", false, true)
	c.SetCookie("csrf_token", csrfToken, 24*3600, "/", "", false, false)

	lang, _ := ah.service.users.getPreferredLanguage(c.Request.Context(), userID)
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Login success!", "language": string(lang)})
}

func (ah *AuthHandler) LoginHandler(c *gin.Context) {
	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if err := VerifyCaptchaToken(c.Request.Context(), user.CaptchaToken, c.ClientIP()); err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	UserID, sessionToken, csrfToken, err := ah.service.LoginService(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, ErrPassowrdMismatch) {
			respondError(c, http.StatusUnauthorized, err)
			return
		}

		respondError(c, statusCodeForError(err), err)
		return
	}

	sessionTokenHash := sha256.Sum256([]byte(sessionToken))
	sessionTokenHashHex := base64.URLEncoding.EncodeToString(sessionTokenHash[:])

	csrfTokenHash := sha256.Sum256([]byte(csrfToken))
	csrfTokenHashHex := base64.URLEncoding.EncodeToString(csrfTokenHash[:])

	userAuth := LoginAuth{
		UserID:           UserID,
		SessionTokenHash: sessionTokenHashHex,
		CSRFTokenHash:    csrfTokenHashHex,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	err = ah.service.users.createUserAuth(c.Request.Context(), userAuth)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", sessionToken, 24*3600, "/", "", false, true)
	c.SetCookie("csrf_token", csrfToken, 24*3600, "/", "", false, false)

	lang, _ := ah.service.users.getPreferredLanguage(c.Request.Context(), UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Login success!", "language": string(lang)})
}

func (ah *AuthHandler) TestHandler(c *gin.Context) {
	user, err := AuthorizeUser(c.Request, ah.service.users)
	if err != nil {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	lang, _ := ah.service.users.getPreferredLanguage(c.Request.Context(), user.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Ты авторизован!", "language": string(lang)})
}

// PreferencesHandler сохраняет язык интерфейса в профиле пользователя.
func (ah *AuthHandler) PreferencesHandler(c *gin.Context) {
	user, ok := CurrentUser(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	var body struct {
		Language string `json:"language"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	lang := Lang(body.Language)
	if !lang.valid() {
		respondValidation(c, map[string]string{"Language": "oneof"})
		return
	}

	if err := ah.service.users.updatePreferredLanguage(c.Request.Context(), user.UserID, lang); err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"language": string(lang)})
}

func (ah *AuthHandler) LogoutHandler(c *gin.Context) {
	if st, err := c.Cookie("session_token"); err == nil && st != "" {
		if decoded, decErr := url.QueryUnescape(st); decErr == nil {
			sessionTokenHash := sha256.Sum256([]byte(decoded))
			sessionTokenHashHex := base64.URLEncoding.EncodeToString(sessionTokenHash[:])
			if err := ah.service.users.deleteUserAuthBySessionTokenHash(c.Request.Context(), sessionTokenHashHex); err != nil {
				respondError(c, statusCodeForError(err), err)
				return
			}
		}
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.SetCookie("csrf_token", "", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "Logout success!"})
}
