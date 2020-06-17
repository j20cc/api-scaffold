package controllers

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lukedever/gvue-scaffold/app/models"
	"github.com/lukedever/gvue-scaffold/internal/helper"
	"github.com/lukedever/gvue-scaffold/internal/jwt"
	"net/http"
	"time"
)

type User struct {
}

//获取详细信息
func (u *User) GetProfile(c *gin.Context) {
	user := GetUserFromContext(c)
	RespondWithJson(c, http.StatusOK, user)
}

type registerRequest struct {
	Name     string `json:"name" valid:"stringlength(6|15)~用户名应该为6-15个字符"`
	Email    string `json:"email" valid:"required,email~邮箱格式不正确"`
	Password string `json:"password" valid:"stringlength(6|15)~密码应该为6-15个字符"`
}

//浏览器端注册
func (u *User) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	if ok, err := govalidator.ValidateStruct(req); !ok {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	if _, exsist := models.FindUser("name", req.Name); exsist {
		ErrResponse(c, http.StatusUnprocessableEntity, errors.New("该用户已存在"))
		return
	}
	if _, exsist := models.FindUser("email", req.Email); exsist {
		ErrResponse(c, http.StatusUnprocessableEntity, errors.New("该邮箱已被注册"))
		return
	}
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: helper.Md5(req.Password),
	}
	if _, err := user.New(); err != nil {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	go user.SendWelcomeEmail()

	SuccessResponse(c)
}

type loginRequest struct {
	Email    string `json:"email" valid:"required,email~邮箱格式不正确"`
	Password string `json:"password" valid:"stringlength(6|15)~密码应该为6-15个字符"`
}

//登录web端
func (u *User) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	if ok, err := govalidator.ValidateStruct(req); !ok {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	user, exsist := models.FindUser("email", req.Email)
	if !exsist {
		ErrResponse(c, http.StatusUnprocessableEntity, errors.New("该邮箱不存在"))
		return
	}
	if user.Password != helper.Md5(req.Password) {
		ErrResponse(c, http.StatusUnprocessableEntity, errors.New("密码不正确"))
		return
	}
	token, _ := jwt.BuildToken(user.ID)
	user.Token = token
	if err := user.Save(); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	RespondWithJson(c, http.StatusOK, user)
}

//发送重置密码邮件
func (u *User) SendResetEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" valid:"required,email~邮箱格式不正确"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	if ok, err := govalidator.ValidateStruct(req); !ok {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	user, exsist := models.FindUser("email", req.Email)
	if !exsist {
		ErrResponse(c, http.StatusUnprocessableEntity, errors.New("该邮箱未注册"))
		return
	}
	if err := user.SendResetEmail(); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	SuccessResponse(c)
}

type resetPasswordRequest struct {
	Email    string `json:"email" valid:"required,email~邮箱格式不正确"`
	Sign     string `json:"sign" valid:"required~链接错误"`
	Password string `json:"password" valid:"stringlength(6|15)~密码应该为6-15个字符"`
}

//重置密码
func (u *User) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	if ok, err := govalidator.ValidateStruct(req); !ok {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	user, err := models.DecodeSignUrl("reset", req.Sign)
	if err != nil {
		ErrResponse(c, http.StatusBadRequest, err)
		return
	}
	if user.Email != req.Email {
		ErrResponse(c, http.StatusBadRequest, errors.New("链接错误"))
		return
	}
	user.Password = helper.Md5(req.Password)
	if err = user.Save(); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	SuccessResponse(c)
}

//发送确认邮箱的邮件
func (u *User) SendVerifyEmail(c *gin.Context) {
	user := GetUserFromContext(c)
	if err := user.SendVerifyEmail(); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	SuccessResponse(c)
}

//验证邮箱
func (u *User) VerifyEmail(c *gin.Context) {
	sign := c.Query("sign")
	if sign == "" {
		ErrResponse(c, http.StatusBadRequest, errors.New("链接错误"))
		return
	}
	user, err := models.DecodeSignUrl("verify", sign)
	if err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	t := time.Now()
	user.EmailVerifiedAt = &t
	_ = user.Save()

	SuccessResponse(c)
}