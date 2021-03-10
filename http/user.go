package http

import (
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleLogin handle '/login' route
func (s *Server) HandleLogin(c *gin.Context) {

}

type registerRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"required"`
}

// HandleRegister handle '/register' route
func (s *Server) HandleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.respondWithValidationErr(c, err)
		return
	}
}
