package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/lukedever/api"
)

// Server represents http server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	translator ut.Translator
	config     *api.Config

	// services
	UserService api.UserService
}

// NewServer return server instance
func NewServer(c *api.Config) *Server {
	gin.SetMode(c.Mode)
	trans := initValidatorTrans(c.Locale)
	r := gin.Default()

	svr := &Server{
		router: r,
		httpServer: &http.Server{
			Handler: r,
			Addr:    c.Addr,
		},
		translator: trans,
		config:     c,
	}

	return svr
}

// Shutdown server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen error: %s\n", err)
		}
	}()
	log.Printf("server running on %s", s.config.Addr)

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exiting")
}

func initValidatorTrans(locale string) ut.Translator {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	zhT := zh.New()
	enT := en.New()

	uni := ut.New(enT, zhT)
	trans, ok := uni.GetTranslator(locale)
	if !ok {
		log.Fatalf("uni.GetTranslator(%s) failed", locale)
	}

	var err error
	switch locale {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(v, trans)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(v, trans)
	default:
		err = enTranslations.RegisterDefaultTranslations(v, trans)
	}
	if err != nil {
		panic(err)
	}

	return trans
}
