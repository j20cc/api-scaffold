package models

import (
	"context"
	"errors"
	"fmt"
	"gvue-scaffold/internal/helper"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
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
	if mysqlCli.Where(key+" = ?", value).First(&user); user.ID > 0 {
		return &user, true
	} else {
		return nil, false
	}
}

func (u *User) Save() error {
	return mysqlCli.Save(u).Error
}

func (u *User) New() (uint, error) {
	if err := mysqlCli.Create(u).Error; err != nil {
		return 0, err
	} else {
		return u.ID, nil
	}
}

func (u *User) SendWelcomeEmail() {
	link, err := u.getSignedUrl("verify")
	if err != nil {
		//TODO 记录日志
		return
	}
	body := fmt.Sprintf("<h3>%s您好:</h3><p>欢迎注册%s，请点击链接: <a href='%s'>%s</a> 进行确认邮箱</p><p>或者直接复制链接 %s 到浏览器打开</p><p>有效期30分钟</p>", u.Name, viper.GetString("app.name"), link, link, link)
	_ = helper.SendEmail(u.Email, "欢迎注册", body)
}

func (u *User) SendVerifyEmail() error {
	link, err := u.getSignedUrl("verify")
	if err != nil {
		//TODO 记录日志
		return err
	}
	body := fmt.Sprintf("<h3>%s您好:</h3><p>您申请验证邮箱，请点击链接: <a href='%s'>%s</a> 进行确认邮箱</p><p>或者直接复制链接 %s 到浏览器打开</p><p>有效期30分钟</p>", viper.GetString("app.name"), link, link, link)
	return helper.SendEmail(u.Email, "验证邮箱", body)
}

func (u *User) SendResetEmail() error {
	link, err := u.getSignedUrl("reset")
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<h3>%s您好:</h3><p>您申请了重置密码，请点击链接: <a href='%s'>%s</a> 进行重置</p><p>或者直接复制链接 %s 到浏览器打开</p><p>有效期30分钟</p>", u.Name, link, link, link)
	return helper.SendEmail(u.Email, "重置密码", body)
}

var keySignPrefix = "user:sign:%s:%s"

func (u *User) getSignedUrl(t string) (string, error) {
	host := viper.GetString("app.url")
	//签名
	signstr := helper.Md5(u.Email)
	values := url.Values{}
	//redis key
	key := fmt.Sprintf(keySignPrefix, t, signstr)
	switch t {
	case "reset":
		host += "/password/reset"
		values.Set("email", u.Email)
	case "verify":
		host += "/verification"
	default:
		return "", errors.New("签名失败")
	}
	values.Set("sign", signstr)
	query := values.Encode()
	//构造url
	myurlstr := host + "?" + query
	myurl, err := url.Parse(myurlstr)
	if err != nil {
		return "", err
	}
	//存到redis
	if _, err := redisCli.Set(context.Background(), key, u.Email, time.Minute*30).Result(); err != nil {
		return "", err
	}
	//解析url
	return myurl.String(), nil
}

func DecodeSignUrl(t, sign string) (*User, error) {
	if t != "reset" && t != "verify" {
		return nil, errors.New("链接错误")
	}
	key := fmt.Sprintf(keySignPrefix, t, sign)
	result, err := redisCli.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, errors.New("签名不存在或过期~链接错误")
	} else if err != nil {
		return nil, err
	}
	user, exists := FindUser("email", result)
	if !exists {
		return nil, errors.New("签名不存在或过期~链接错误")
	}
	return user, nil
}
