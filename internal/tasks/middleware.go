package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const currentUserContextKey = "current_user"

func AuthMiddleware(users *UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := AuthorizeUser(c.Request, users)
		if err != nil {
			respondError(c, http.StatusUnauthorized, ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(currentUserContextKey, user)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok || user.Role != RoleAdmin {
			respondError(c, http.StatusForbidden, ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*UserContext, bool) {
	v, ok := c.Get(currentUserContextKey)
	if !ok {
		return nil, false
	}

	user, ok := v.(*UserContext)
	return user, ok
}
