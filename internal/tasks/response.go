package tasks

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, status int, err error) {
	if err == nil {
		c.JSON(status, gin.H{"error": localizeError(LangFromRequest(c), errors.New("unknown error"))})
		return
	}

	c.JSON(status, gin.H{"error": localizeError(LangFromRequest(c), err)})
}

func respondValidation(c *gin.Context, errs map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{"errors": SendValidationErrors(errs, LangFromRequest(c))})
}

func respondPaginated(c *gin.Context, items any, total int64, page int, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func statusCodeForError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrRequestTooLarge):
		return http.StatusRequestEntityTooLarge
	case errors.Is(err, ErrFileTooLarge), errors.Is(err, ErrUnsupportedFileType), errors.Is(err, ErrInvalidField), errors.Is(err, ErrInvalidStatusTransition):
		return http.StatusBadRequest
	case errors.Is(err, ErrCaptchaRequired), errors.Is(err, ErrCaptchaFailed):
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
