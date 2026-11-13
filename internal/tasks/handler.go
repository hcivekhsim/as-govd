package tasks

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	found, err := ah.service.users.getByLogin(c.Request.Context(), user.Login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if found {
		c.JSON(http.StatusBadRequest, gin.H{"user exists": user.Login})
		return
	}

	if errs := ValidateStruct(user); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": SendReportOfStructAuth(errs)})
		return
	}

	id, err := ah.service.RegisterService(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"User created with id": id})
}

func (ah *AuthHandler) LoginHandler(c *gin.Context) {
	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UserID, sessionToken, csrfToken, err := ah.service.LoginService(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, ErrPassowrdMismatch) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", sessionToken, 24*3600, "/", "", false, true)
	c.SetCookie("csrf_token", csrfToken, 24*3600, "/", "", false, false)

	c.String(http.StatusOK, "Login success!")
}

func (ah *AuthHandler) TestHandler(c *gin.Context) {
	if err := Authorize(c.Request, ah.service.users); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.String(http.StatusOK, "Ты авторизован!")
}
