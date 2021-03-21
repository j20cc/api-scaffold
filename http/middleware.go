package http

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens := c.GetHeader("Authorization")
		token := strings.TrimPrefix(tokens, "Bearer ")
		if token == "" {
			s.respondWithAuthErr(c, errInvalidToken)
			c.Abort()
			return
		}
		claims, err := s.parseToken(token)
		if err != nil {
			s.respondWithAuthErr(c, err)
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.UserName)

		c.Next()
	}
}
