package code

type Error struct {
	code    int
	message string
}

func NewError(code int, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) SetMessage(message string) error {
	e.message = message
	return e
}

var (
	OK                  = NewError(200, "成功")
	BadRequest          = NewError(400, "错误请求")
	Unauthorized        = NewError(401, "未授权")
	Forbidden           = NewError(403, "禁止访问")
	NotFound            = NewError(404, "未找到")
	InternalServerError = NewError(500, "内部服务器错误")
	ServiceUnavailable  = NewError(503, "服务不可用")
	InvalidParameter    = NewError(422, "参数无效")
	DuplicateOperation  = NewError(409, "重复操作")
	StatusNotAvailable  = NewError(1001, "状态不可用")
	StatusException     = NewError(1002, "状态异常")
	MyCustomErrorCode   = NewError(1003, "自定义错误码")
	Expired             = NewError(1004, "已过期")
)
