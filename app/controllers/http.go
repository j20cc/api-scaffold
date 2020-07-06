package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lukedever/gvue-scaffold/app/validations"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	errModelNotFound = errors.New("model not found")
	defaultPageSize  = 15
)

func getQueryPageSize(c *gin.Context) (int, int) {
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "15")
	p, _ := strconv.Atoi(page)
	s, _ := strconv.Atoi(size)
	if p == 0 {
		p = 1
	}
	if s == 0 {
		s = defaultPageSize
	}

	return p, s
}

func RespondWithJson(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func SuccessResponse(c *gin.Context) {
	RespondWithJson(c, http.StatusOK, gin.H{"message": "success"})
}

func ErrResponse(c *gin.Context, code int, err error) {
	RespondWithJson(c, code, gin.H{"error": err.Error()})
}

func ErrValidateResponse(c *gin.Context, err error, st interface{}) {
	code := http.StatusUnprocessableEntity
	defaultErr := errors.New("验证失败")
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ErrResponse(c, code, defaultErr)
		return
	}
	e := errs[0]
	locale := c.GetHeader("Accept-Language")
	translator := validations.GetTranslator(locale)
	//先使用翻译，但是得到'name是必填的'，想要'姓名是必填的'
	msg := e.Translate(translator)
	t := reflect.TypeOf(st)
	if t.Kind().String() != "struct" {
		ErrResponse(c, code, defaultErr)
		return
	}
	//取出label的tag根据locale替换
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == e.StructField() {
			labels := t.Field(i).Tag.Get("label")
			tags := strings.Split(labels, ",")
			m := make(map[string]string, len(tags))
			for _, tag := range tags {
				i := strings.Split(tag, "=")
				if len(i) == 2 {
					m[i[0]] = i[1]
				}
			}
			if name, ok := m[locale]; ok {
				msg = strings.ReplaceAll(msg, t.Field(i).Tag.Get("json"), name)
				ErrResponse(c, code, errors.New(msg))
				return
			}
			if locale == "en" {
				ErrResponse(c, code, errors.New(msg))
				return
			}
		}
	}
	ErrResponse(c, code, defaultErr)
	return
}
