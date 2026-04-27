package errors

import (
	"errors"
	"fmt"
)

// BizError 业务错误
type BizError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 创建新的业务错误
func New(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// NewWithErr 创建带原始错误的业务错误
func NewWithErr(code int, message string, err error) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewDefault 创建使用默认消息的业务错误
func NewDefault(code int) *BizError {
	return &BizError{
		Code:    code,
		Message: GetMessage(code),
	}
}

func AsBizError(err error) (*BizError, bool) {
	if err == nil {
		return nil, false
	}

	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return bizErr, true
	}

	return nil, false
}

func Wrap(code int, message string, err error) error {
	if err == nil {
		return nil
	}
	return NewWithErr(code, message, err)
}

// 预定义常用错误
var (
	ErrBadRequest    = NewDefault(CodeBadRequest)
	ErrUnauthorized  = NewDefault(CodeUnauthorized)
	ErrForbidden     = NewDefault(CodeForbidden)
	ErrNotFound      = NewDefault(CodeNotFound)
	ErrInternalError = NewDefault(CodeInternalError)
)
