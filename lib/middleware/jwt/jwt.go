package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/webx-top/echo"
)

func New(secret string) *JWT {
	return &JWT{
		Secret: secret,
		CondFn: func(c echo.Context) bool {
			ignore, _ := c.Get(`webx:ignoreJwt`).(bool)
			return !ignore
		},
	}
}

type JWT struct {
	Secret string
	CondFn func(echo.Context) bool
}

func (j *JWT) Validate() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.IsFileServer() || (j.CondFn != nil && j.CondFn(c) == false) {
				return h(c)
			}
			/*//Test
			tokenString, err := j.Response(map[string]interface{}{"uid": "1", "username": "admin"})
			if err == nil {
				println("jwt token:", tokenString)
			}
			//*/
			token, err := jwt.ParseFromRequest(c.Request(), func(token *jwt.Token) (interface{}, error) {
				b := ([]byte(j.Secret))
				return b, nil
			})
			if err != nil {
				return err
			}
			if !token.Valid {
				return errors.New(`Incorrect signature.`)
			}
			c.Set(`webx:jwtClaims`, token.Claims)
			return h(c)
		}
	}
}

func (j *JWT) Claims(c echo.Context) map[string]interface{} {
	r, _ := c.Get(`webx:jwtClaims`).(map[string]interface{})
	return r
}

func (j *JWT) Ignore(on bool, c echo.Context) {
	c.Set(`webx:ignoreJwt`, on)
}

/*
本函数所生成结果的用法
用法一：写入header头的属性“Authorization”中，值设为：前缀BEARER加tokenString的值
用法二：发送post或get参数“access_token”，值设为：tokenString的值
*/
func (j *JWT) Response(values map[string]interface{}) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"
	// Claims
	token.Claims = values
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err = token.SignedString([]byte(j.Secret))
	return
}
