package middleware

import (
	"gopkg.in/gin-gonic/gin.v1"
	"MediaServer/secret"
	"time"
	"fmt"
)

//校验 uuid
func UUID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//校验 登录 token
func LoginToken() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//单点登录
func SSO() gin.HandlerFunc {
	return func(c *gin.Context) {
		stoken := c.Request.Header.Get("X_STOKEN")
		loginToken, err := secret.ValidateLoginSToken(stoken)
		if err != nil {
			stoken, _ := secret.GenLoginSToken(10, "gxl", time.Now().Add(time.Hour * 24).Unix())
			fmt.Println(stoken)
			c.Abort()
		}

		c.Set("loginToken", loginToken)
		return
	}
}

//校验 csrf token
func CsrfToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		return
		csrfToken := c.Request.Header.Get("HTTP_X_CSRF_TOKEN")
		if csrfToken == "aaa" {
			println("token is ok")
			return
		}

		c.Abort()
		c.JSON(401, gin.H{
			"code":401,
			"msg":"token invalid",
		})
	}
}