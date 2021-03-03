package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukedever/api"
)

// Server represents http server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server

	Mode string
	Addr string

	// services
	UserService api.UserService
}

// NewServer return server instance
func NewServer(mode, addr string) *Server {
	r := gin.Default()
	gin.SetMode(mode)
	r.GET("/welcome", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome",
		})
	})

	svr := &Server{
		router:     r,
		httpServer: &http.Server{},

		Mode: mode,
		Addr: addr,
	}
	svr.httpServer.Handler = svr.router
	svr.httpServer.Addr = svr.Addr

	return svr
}

// Run server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
