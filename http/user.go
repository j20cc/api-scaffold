package http

import (
	"errors"
	"net/http"
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

var ErrEmailExists = errors.New("邮箱用户已存在")

// HandleRegister handle '/register' route
func (s *Server) HandleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.respondWithValidationErr(c, err)
		return
	}

	user, _ := s.UserService.FindUserByKV("email", req.Email)
	if user.ID > 0 {
		s.respondWithValidationErr(c, ErrEmailExists)
		return
	}

	ss := strings.Split(req.Email, "@")
	u := api.User{
		Name:     ss[0],
		Email:    req.Email,
		Password: md5Str(req.Password),
	}
	if err := s.UserService.CreateUser(&u); err != nil {
		s.respondWithServerErr(c, err, true)
		return
	}

	s.responseWithData(c, http.StatusOK, u)
}
