package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukedever/gvue-scaffold/app/models"
	"github.com/lukedever/gvue-scaffold/internal/helper"
	"github.com/lukedever/gvue-scaffold/internal/jwt"
)

// User controller
type User struct{}

type registerRequest struct {
	Name     string `json:"name" binding:"min=3,max=15" label:"zh=用户名"`
	Email    string `json:"email" binding:"email" label:"zh=邮箱"`
	Password string `json:"password" binding:"min=6,max=15" label:"zh=密码"`
}

// Register user action
func (u *User) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrValidateResponse(c, err, req)
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
	RespondWithJSON(c, http.StatusOK, user)
}

type loginRequest struct {
	Email    string `json:"email" binding:"email" label:"zh=邮箱"`
	Password string `json:"password" binding:"min=6,max=15" label:"zh=密码"`
}

// Login action
func (u *User) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrValidateResponse(c, err, req)
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
	RespondWithJSON(c, http.StatusOK, user)
}

// SendResetEmail send reset password email
func (u *User) SendResetEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"email" label:"zh=邮箱"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrValidateResponse(c, err, req)
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
	Email    string `json:"email" binding:"email" label:"zh=邮箱"`
	Sign     string `json:"sign" binding:"required" label:"zh=签名"`
	Password string `json:"password" binding:"min=6,max=15" label:"zh=密码"`
}

// ResetPassword reset password action
func (u *User) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrValidateResponse(c, err, req)
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

// SendVerifyEmail send verify email action
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

// VerifyEmail action
func (u *User) VerifyEmail(c *gin.Context) {
	var req struct {
		Sign string `json:"sign" binding:"required" label:"zh=签名"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrValidateResponse(c, err, req)
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

// GetProfile get login user profile
func (u *User) GetProfile(c *gin.Context) {
	user, exists := models.FindUser("id", c.GetString("userId"))
	if !exists {
		ErrResponse(c, http.StatusUnauthorized, errModelNotFound)
		return
	}
	RespondWithJSON(c, http.StatusOK, user)
}

// Demo action
func (User) Demo(c *gin.Context) {
	var req struct {
		Data string `json:"data" binding:"required,hasUser=1111" label:"zh=数据"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrValidateResponse(c, err, req)
		return
	}
}
