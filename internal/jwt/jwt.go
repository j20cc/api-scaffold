package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var (
	mySigningKey    = []byte(viper.GetString("jwt.secret"))
	tokenInvalidErr = errors.New("token invalid or expired")
)

func GetTokenExpireSeconds() int64 {
	expireDay := viper.GetInt("jwt.expire_day")
	return int64(expireDay * 24 * 60 * 60)
}

func BuildToken(userId uint) (string, error) {
	issuer := viper.GetString("app.name")
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + GetTokenExpireSeconds(),
		Issuer:    issuer,
		Id:        strconv.Itoa(int(userId)),
		Subject:   "login",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ParseToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return strconv.Atoi(claims.Id)
	} else {
		return 0, tokenInvalidErr
	}
}
