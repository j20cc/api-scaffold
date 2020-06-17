package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/lukedever/gvue-scaffold/app/controllers"
	"github.com/lukedever/gvue-scaffold/app/models"
	"github.com/lukedever/gvue-scaffold/internal/jwt"
	"net/http"
	"strconv"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		userId, err := jwt.ParseToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			controllers.ErrResponse(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		if user, exist := models.FindUser("id", strconv.Itoa(userId)); exist {
			c.Set("user", &user)
		} else {
			controllers.ErrResponse(c, http.StatusUnauthorized, controllers.ErrModelNotFound)
			c.Abort()
			return
		}
	}
}
