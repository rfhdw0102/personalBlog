package middleware

import (
	"blogs/internal/repository"
	"blogs/pkg/jwt"
	"blogs/pkg/responses"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(redisRepo repository.RedisRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		// 如果前端带有 Bearer 前缀，则去掉
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		if tokenString == "" {
			responses.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 校验 JWT
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			responses.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 校验 Redis
		ctx := context.Background()
		key := fmt.Sprintf("user_token:%d", claims.UserID)

		storedToken, err := redisRepo.GetKey(ctx, key)
		if err != nil {
			responses.Unauthorized(c, "登录已失效")
			c.Abort()
			return
		}

		// 对比 token
		if storedToken != tokenString {
			responses.Unauthorized(c, "账号已在其他设备登录")
			c.Abort()
			return
		}

		// 写入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Roles)

		c.Next()
	}
}

// AdminAuthMiddleware 管理员权限校验中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			responses.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTOptionalMiddleware 可选的 JWT 中间件：如果提供 Token 则解析并设置上下文，不提供则直接通过
func JWTOptionalMiddleware(redisRepo repository.RedisRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		if tokenString == "" {
			c.Next()
			return
		}

		// 校验 JWT
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			// Token 无效或过期，但不中断请求，仅作为游客处理
			c.Next()
			return
		}

		// 校验 Redis
		ctx := context.Background()
		key := fmt.Sprintf("user_token:%d", claims.UserID)

		storedToken, err := redisRepo.GetKey(ctx, key)
		if err != nil || storedToken != tokenString {
			// Token 在 Redis 中不存在或不匹配，不中断请求
			c.Next()
			return
		}

		// 写入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Roles)

		c.Next()
	}
}
