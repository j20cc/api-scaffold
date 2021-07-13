package http

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
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
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<h1>welcome gopher</h1>"))
	})

	apiR := r.Group("/api")
	{
		// user routes
		apiR.POST("/login", s.HandleLogin)
		apiR.POST("/register", s.HandleRegister)
		apiR.GET("/profile", s.auth(), s.HandleProfile)
	}
}

type customJwtClaims struct {
	UserID   int
	UserName string
	jwt.StandardClaims
}

var errInvalidToken = errors.New("invalid token")

func (s *Server) genToken(id uint, name string) (string, error) {
	c := customJwtClaims{
		int(id),
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(s.config.JWT.TTL) * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *Server) parseToken(tokenString string) (*customJwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customJwtClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*customJwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errInvalidToken
}

func (s *Server) respondWithServerErr(c *gin.Context, err error) {
	log.Printf("error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"err_msg": "server error",
	})
}

func (s *Server) respondWithAuthErr(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"err_msg": err.Error(),
	})
}

func (s *Server) respondWithErr(c *gin.Context, err error) {
	code := http.StatusBadRequest
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(code, gin.H{
			"err_msg": err.Error(),
		})
		return
	}
	code = http.StatusUnprocessableEntity
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

func md5Str(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)
	return fmt.Sprintf("%x", w.Sum(nil))
}
