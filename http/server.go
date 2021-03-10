package http

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"strings"

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

// Run server
func (s *Server) Run() error {
	s.registerRoutes()

	return s.httpServer.ListenAndServe()
}

// Shutdown server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func initValidatorTrans(locale string) ut.Translator {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	// 注册一个获取json tag的自定义方法
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
	// 注册翻译器
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
