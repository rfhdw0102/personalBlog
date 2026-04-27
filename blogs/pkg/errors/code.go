package errors

// 错误码定义
const (
	//  成功
	CodeSuccess = 200

	//  通用错误 4xx
	CodeBadRequest       = 400
	CodeUnauthorized     = 401
	CodeForbidden        = 403
	CodeNotFound         = 404
	CodeMethodNotAllowed = 405

	// 服务器错误 5xx
	CodeInternalError      = 500
	CodeNotImplemented     = 501
	CodeServiceUnavailable = 503
)

// 错误码对应的文本消息
var codeMessages = map[int]string{
	CodeSuccess:            "成功",
	CodeBadRequest:         "请求参数错误",
	CodeUnauthorized:       "未授权",
	CodeForbidden:          "未登录,禁止访问",
	CodeNotFound:           "资源不存在",
	CodeMethodNotAllowed:   "方法不允许",
	CodeInternalError:      "服务器内部错误",
	CodeNotImplemented:     "功能未实现",
	CodeServiceUnavailable: "服务不可用",
}

// GetMessage 获取错误码对应的文本消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
