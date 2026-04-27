package jwt

import (
	"blogs/pkg/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateToken 生成 JWT Token
func GenerateToken(userID int, username string, role string) (string, error) {
	jwtCfg := config.Get().JWT
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtCfg.ExpireHours * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blogs",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtCfg.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	jwtCfg := config.Get().JWT
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtCfg.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token已过期")
		}
		return nil, errors.New("token无效")
	}
	// 校验并返回自定义载荷
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token无效")
}

// RefreshToken 刷新 Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return GenerateToken(claims.UserID, claims.Username, claims.Roles)
}

// GetUserID 获取用户 ID
func (c *Claims) GetUserID() int {
	return c.UserID
}

// GetUsername 获取用户名
func (c *Claims) GetUsername() string {
	return c.Username
}

// GetRoles 获取用户角色
func (c *Claims) GetRoles() string {
	return c.Roles
}
