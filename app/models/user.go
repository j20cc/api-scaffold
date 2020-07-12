package models

import (
	"fmt"
	"github.com/lukedever/gvue-scaffold/internal/mail"
	"github.com/spf13/viper"
	"time"
)

type User struct {
	Model
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Name            string     `json:"name" gorm:"size:50"`
	Email           string     `json:"email" gorm:"size:50"`
	Password        string     `json:"-" gorm:"size:50"`
	Avatar          string     `json:"avatar"`
	Token           string     `json:"token"`
}

func FindUser(key, value string) (*User, bool) {
	var user User
	if MysqlCli.Where(key+" = ?", value).First(&user); user.ID > 0 {
		return &user, true
	} else {
		return nil, false
	}
}

func (u *User) Save() error {
	return MysqlCli.Save(u).Error
}

func (u *User) New() (uint, error) {
	if err := MysqlCli.Create(u).Error; err != nil {
		return 0, err
	} else {
		return u.ID, nil
	}
}

func (u *User) SendWelcomeEmail() {
	link, err := mail.GetSignedUrl(u.Email, mail.Verify)
	if err != nil {
		return
	}
	body := fmt.Sprintf("<h3>%s您好:</h3><p>欢迎注册%s，请点击链接: <a href='%s'>%s</a> 进行确认邮箱</p><p>或者直接复制链接 %s 到浏览器打开</p><p>有效期30分钟</p>", u.Name, viper.GetString("app.name"), link, link, link)
	_ = mail.SendEmail(u.Email, "欢迎注册", body)
}

func (u *User) SendVerifyEmail() error {
	link, err := mail.GetSignedUrl(u.Email, mail.Verify)
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<h3>%s您好:</h3><p>您申请验证邮箱，请点击链接: <a href='%s'>%s</a> 进行确认邮箱</p><p>或者直接复制链接 %s 到浏览器打开</p><p>有效期30分钟</p>", viper.GetString("app.name"), link, link, link)
	return mail.SendEmail(u.Email, "验证邮箱", body)
}

func (u *User) SendResetEmail() error {
	link, err := mail.GetSignedUrl(u.Email, mail.Reset)
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<h3>%s您好:</h3><p>您申请了重置密码，请点击链接: <a href='%s'>%s</a> 进行重置</p><p>或者直接复制链接 %s 到浏览器打开</p><p>有效期30分钟</p>", u.Name, link, link, link)
	return mail.SendEmail(u.Email, "重置密码", body)
}