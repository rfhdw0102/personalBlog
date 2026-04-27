package responses

import (
	"blogs/internal/model/dto/response"
	"blogs/pkg/errors"
	"github.com/gin-gonic/gin"
	"math"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` //omitempty 表示：当该字段的值为“空值”时，在 JSON 输出中省略（不包含）这个字段。
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    errors.CodeSuccess,
		Message: errors.GetMessage(errors.CodeSuccess),
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    errors.CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	status := code
	if code < 100 || code > 599 {
		status = 200
	}

	c.JSON(status, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	Error(c, errors.CodeBadRequest, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, errors.CodeUnauthorized, message)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, message string) {
	Error(c, errors.CodeForbidden, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	Error(c, errors.CodeNotFound, message)
}

// InternalError 500 错误
func InternalError(c *gin.Context, message string) {
	Error(c, errors.CodeInternalError, message)
}

func FromError(c *gin.Context, err error, fallbackCode int, fallbackMessage string) {
	if bizErr, ok := errors.AsBizError(err); ok {
		Error(c, bizErr.Code, bizErr.Message)
		return
	}

	Error(c, fallbackCode, fallbackMessage)
}

// NewPageResponse 创建分页响应
func NewPageResponse(list interface{}, total int64, page, pageSize int) *response.PageResponse {
	pages := int(math.Ceil(float64(total) / float64(pageSize)))
	return &response.PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}
