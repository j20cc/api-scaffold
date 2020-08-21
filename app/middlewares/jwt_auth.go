package middlewares

import (
	"errors"
	"gvue-scaffold/app/controllers"
	"gvue-scaffold/internal/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			controllers.ErrResponse(c, http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()
			return
		}
		userId, err := helper.ParseToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil || userId == "" {
			controllers.ErrResponse(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
