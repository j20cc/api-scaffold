package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lukedever/gvue-scaffold/app/validations"
	"net/http"
	"strconv"
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
	if errs, ok := err.(validator.ValidationErrors); ok {
		err := errs[0]
		locale := c.GetHeader("Accept-Language")
		translator := validations.GetTranslator(locale)
		msg := err.Translate(translator)
		RespondWithJson(c, code, gin.H{"error": msg})
		return
	}
	RespondWithJson(c, code, gin.H{"error": err.Error()})
}
