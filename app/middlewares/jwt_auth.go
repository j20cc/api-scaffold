package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lukedever/gvue-scaffold/app/controllers"
	"github.com/lukedever/gvue-scaffold/internal/helper"
	"net/http"
	"strings"
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
