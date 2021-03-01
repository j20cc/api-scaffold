package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lukedever/api"
)

// Server represents http server
type Server struct {
	router *gin.Engine

	Mode string

	// services
	UserService api.UserService
}

// NewServer return server instance
func NewServer() *Server {
	svr := &Server{
		router: gin.Default(),
	}

	return svr
}
