package middlewares

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ClaimsJWT 自定义 JWT 载荷结构体
type ClaimsJWT struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成令牌
func GenerateToken(userID uint) (string, error) {

	claims := &ClaimsJWT{
		// ID标识 - 使用传入的实际用户ID
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// 使用 HS256 签名方法创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // 请使用更安全的密钥
	return tokenstring, err
}

// JWT 验证中间件
// gin.HandlerFunc 来定义中间件函数
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "请求未携带令牌"})
			c.Abort()
			return
		}
		// Bearer <token>
		// 分成两个部分
		authBearer := strings.SplitN(authHeader, " ", 2)
		if len(authBearer) != 2 || authBearer[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "请求头格式错误"})
			c.Abort()
			return
		}
		tokenString := authBearer[1]
		claims := &ClaimsJWT{}
		// 解析令牌
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
