package utils

import (
	"blogs/pkg/errors"
	"blogs/pkg/responses"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CurrentUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		responses.Unauthorized(c, "未登录")
		return 0, false
	}

	uid, ok := userID.(int)
	if !ok {
		responses.Unauthorized(c, "未登录")
		return 0, false
	}

	return uid, true
}

func CurrentRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		responses.Unauthorized(c, "权限未知")
		return "", false
	}

	value, ok := role.(string)
	if !ok {
		responses.Unauthorized(c, "权限未知")
		return "", false
	}

	return value, true
}

func ParamInt(c *gin.Context, name string) (int, bool) {
	value, err := strconv.Atoi(c.Param(name))
	if err != nil {
		responses.FromError(c, errors.New(errors.CodeBadRequest, "无效的 ID"), errors.CodeBadRequest, "无效的 ID")
		return 0, false
	}

	return value, true
}
