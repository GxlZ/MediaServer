package secret

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

type Token struct {
	UserId   int64
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	sign = []byte("91u!!T!&#@^#")
)

//生成 签名的 token
//sign 数字签名
//expireTime 过期时间
func GenLoginSToken(userId int64, userName string, expireTime int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Token{
		UserId:userId,
		Username:userName,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt:expireTime,
		},
	})

	sToken, err := token.SignedString(sign)

	return sToken, err
}

//登录 签名token 校验
func ValidateLoginSToken(stoken string) (*Token, error) {

	loginTokenObj := &Token{}

	token, err := jwt.ParseWithClaims(stoken, loginTokenObj, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return sign, nil
	})

	if err != nil {
		return loginTokenObj, err
	}

	if loginToken, ok := token.Claims.(*Token); ok && token.Valid {
		return loginToken, nil
	}

	return loginTokenObj, err
}
