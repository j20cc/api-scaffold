package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lukedever/gvue-scaffold/app/validations"
	"github.com/spf13/viper"
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

// RespondWithJSON respond data with json
func RespondWithJSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// SuccessResponse quickly response success message
func SuccessResponse(c *gin.Context) {
	RespondWithJSON(c, http.StatusOK, gin.H{"message": "success"})
}

// ErrResponse return error message
func ErrResponse(c *gin.Context, code int, err error) {
	RespondWithJSON(c, code, gin.H{"error": err.Error()})
}

// ErrValidateResponse return validation message
func ErrValidateResponse(c *gin.Context, err error, st interface{}) {
	code := http.StatusUnprocessableEntity
	defaultErr := errors.New("验证失败")
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ErrResponse(c, code, defaultErr)
		return
	}
	t := reflect.TypeOf(st)
	if t.Kind().String() != "struct" {
		ErrResponse(c, code, defaultErr)
		return
	}

	locale := c.GetHeader("Accept-Language")
	fmt.Println(locale)
	if locale == "" {
		locale = viper.GetString("app.locale")
	}
	translator := validations.GetTranslator(locale)
	errMsg := make(map[string]string, len(errs))
	//先使用翻译，但是得到'name是必填的'，想要'姓名是必填的'
	for _, e := range errs {
		errMsg[e.Field()] = e.Translate(translator)
	}
	//取出label的tag根据locale替换
	pattern := regexp.MustCompile(locale + "=([^,\\s]*)")
	for f, m := range errMsg {
		for i := 0; i < t.NumField(); i++ {
			jsonTag := t.Field(i).Tag.Get("json")
			if jsonTag == f {
				//"en=aaaa,zh=数据,fr=bbbb"
				labels := t.Field(i).Tag.Get("label")
				matches := pattern.FindAllStringSubmatch(labels, -1)
				fmt.Println(labels, matches, locale)
				if len(matches) == 0 || len(matches[0]) < 2 {
					continue
				}
				name := matches[0][1]
				m = strings.ReplaceAll(m, jsonTag, name)
				errMsg[f] = m
			}
		}
	}
	RespondWithJSON(c, http.StatusUnprocessableEntity, gin.H{
		"error": errMsg,
	})
}
