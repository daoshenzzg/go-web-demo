package ecode

// All common ecode
var (
	OK = add(200, "OK")

	NoLogin = add(-101, "账号未登录")

	RequestErr         = add(-400, "请求错误")
	Unauthorized       = add(-401, "未认证")
	AccessDenied       = add(-403, "访问权限不足")
	NotFound           = add(-404, "404")
	MethodNotAllowed   = add(-405, "不支持该方法")
	Conflict           = add(-409, "冲突")
	ServerErr          = add(-500, "服务器错误")
	ServiceUnavailable = add(-503, "过载保护，服务暂时不可用")
	Deadline           = add(-504, "服务调用超时")
	LimitExceed        = add(-509, "超出限制")
)
