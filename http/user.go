package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lukedever/api"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleLogin handle '/login' route
func (s *Server) HandleLogin(c *gin.Context) {

}

type registerRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6,max=12"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}

// HandleRegister handle '/register' route
func (s *Server) HandleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.respondWithValidationErr(c, err)
		return
	}

	_, err := s.UserService.FindUserByKV("email", req.Email)
	if err != nil {
		s.respondWithServerErr(c)
		return
	}

	user := api.User{
		Name: strings.TrimFunc(req.Email, func(r rune) bool {
			return string(r) != "@"
		}),
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.UserService.CreateUser(&user); err != nil {

	}
}
