package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lukedever/gvue-scaffold/app/models"
	"github.com/lukedever/gvue-scaffold/internal/helper"
	"github.com/lukedever/gvue-scaffold/internal/jwt"
	"net/http"
	"time"
)

type User struct {}

type registerRequest struct {
	Name     string `json:"name" binding:"min=3,max=15"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"min=6,max=15"`
}

//浏览器端注册
func (u *User) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
	//发送欢迎邮件
	go user.SendWelcomeEmail()
	//设置token
	token, _ := jwt.BuildToken(user.ID)
	user.Token = token
	if err := user.Save(); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	RespondWithJson(c, http.StatusOK, user)
}

type loginRequest struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"min=6,max=15"`
}

//登录web端
func (u *User) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	user, exsist := models.FindUser("email", req.Email)
	if !exsist {
		ErrResponse(c, http.StatusUnprocessableEntity, errors.New("该邮箱用户不存在"))
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
		Email string `json:"email" binding:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
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
	Email    string `json:"email" binding:"email"`
	Sign     string `json:"sign" binding:"required"`
	Password string `json:"password" binding:"min=6,max=15"`
}

//重置密码
func (u *User) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
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
	user, exists := models.FindUser("id", c.GetString("userId"))
	if !exists {
		ErrResponse(c, http.StatusUnauthorized, errModelNotFound)
		return
	}
	if err := user.SendVerifyEmail(); err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	SuccessResponse(c)
}

//验证邮箱
func (u *User) VerifyEmail(c *gin.Context) {
	var req struct {
		Sign string `json:"sign" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrResponse(c, http.StatusBadRequest, errors.New("链接错误"))
		return
	}
	user, err := models.DecodeSignUrl("verify", req.Sign)
	if err != nil {
		ErrResponse(c, http.StatusInternalServerError, err)
		return
	}
	t := time.Now()
	user.EmailVerifiedAt = &t
	_ = user.Save()

	SuccessResponse(c)
}

//获取详细信息
func (u *User) GetProfile(c *gin.Context) {
	user, exists := models.FindUser("id", c.GetString("userId"))
	if !exists {
		ErrResponse(c, http.StatusUnauthorized, errModelNotFound)
		return
	}
	RespondWithJson(c, http.StatusOK, user)
}

func (User) Demo(c *gin.Context) {
	var req struct{
		Data string `json:"data" binding:"required,hasUser=1111"`
	}
	if err := c.ShouldBindJSON(&req);err != nil {
		ErrResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
}
