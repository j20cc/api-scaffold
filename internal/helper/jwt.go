package helper

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var (
	mySigningKey    = []byte(viper.GetString("jwt.secret"))
	errTokenInvalid = errors.New("token invalid or expired")
)

func getTokenExpireSeconds() int64 {
	expireDay := viper.GetInt("jwt.expire_day")
	return int64(expireDay * 24 * 60 * 60)
}

// BuildToken build a token
func BuildToken(userID uint) (string, error) {
	issuer := viper.GetString("app.name")
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + getTokenExpireSeconds(),
		Issuer:    issuer,
		Id:        strconv.Itoa(int(userID)),
		Subject:   "login",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// ParseToken parse and verify a token
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Id, nil
	}
	return "", errTokenInvalid
}
