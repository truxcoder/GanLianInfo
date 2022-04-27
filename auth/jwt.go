package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 这里需要额外记录一个username或ID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

const (
	TokenExpireDuration = time.Hour * 24
)

var (
	MySecret = []byte("为人民服务")
)

// GenToken 生成JWT
func GenToken(id string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		id, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "GanBuGuanLi",                              // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("非法token")
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 通常Token放在Header的Authorization中，并使用Bearer开头
		// 这里要和前端vue配合放在头部的"X-Token"中
		//authHeader := c.Request.Header.Get("Authorization")
		authHeader := c.Request.Header.Get("X-Token")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    2003,
				"message": "鉴权信息为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		//parts := strings.SplitN(authHeader, " ", 2)
		//if !(len(parts) == 2 && parts[0] == "Bearer") {
		//	c.JSON(http.StatusOK, gin.H{
		//		"code":    2004,
		//		"message": "请求头中auth格式有误",
		//	})
		//	c.Abort()
		//	return
		//}
		// parts[1]是获取到的tokenString，使用之前定义好的解析JWT的函数来解析它
		//mc, err := ParseToken(parts[1])
		mc, err := ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    2005,
				"message": "鉴权认证失败",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("userId", mc.ID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
