package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lukedever/gvue-scaffold/app/models"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrModelNotFound = errors.New("not found")
	defaultPageSize  = 15
)

func GetUserFromContext(c *gin.Context) *models.User {
	user, exsit := c.Get("user")
	if !exsit {
		ErrResponse(c, http.StatusUnauthorized, ErrModelNotFound)
		return nil
	}
	u, ok := user.(*models.User)
	if !ok {
		ErrResponse(c, http.StatusUnauthorized, ErrModelNotFound)
		return nil
	}
	return u
}

func GetQueryPageSize(c *gin.Context) (int, int) {
	page := c.Query("page")
	size := c.Query("size")
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
	if code == http.StatusUnprocessableEntity {
		msg := strings.Split(err.Error(), ";")
		RespondWithJson(c, code, gin.H{"error": msg[0]})
	} else {
		RespondWithJson(c, code, gin.H{"error": err.Error()})
	}
}
