package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s *Server) registerRoutes() {
	r := s.router
	r.GET("/welcome", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome gopher",
		})
	})

	r.POST("/login", s.HandleLogin)
	r.POST("/register", s.HandleRegister)
}

type customJwtClaims struct {
	userID   int
	userName string
	jwt.StandardClaims
}

func createToken(id int, name, secret string) (string, error) {
	claims := customJwtClaims{
		id,
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 100).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func parseToken(tokenString, secret string) {
	token, err := jwt.ParseWithClaims(tokenString, &customJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*customJwtClaims); ok && token.Valid {
		fmt.Printf("%v", claims.userID)
	} else {
		fmt.Println(err)
	}
}

func responseWithData(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func (s *Server) respondWithServerErr(c *gin.Context) {
	responseWithData(c, http.StatusInternalServerError, gin.H{
		"err_msg": "server error",
	})
}

func (s *Server) respondWithValidationErr(c *gin.Context, err error) {
	code := http.StatusUnprocessableEntity
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(code, gin.H{
			"err_msg": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"err_msg": removeTopStruct(errs.Translate(s.translator)),
	})
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
