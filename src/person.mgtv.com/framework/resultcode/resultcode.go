package resultcode

var errors = make(map[string]string)

const (
	SUCCESS = "0"

	ERROR_1000 = "1000"
	ERROR_1001 = "1001"

	ERROR = "9999"
)

func init() {
	errors[ERROR_1000] = "方法执行异常"
	errors[ERROR_1001] = "参数缺失"
	errors[ERROR] = "系统异常"
}

func ErrorMsg(code string) string {
	return errors[code]
}
