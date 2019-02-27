package mvc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"sync"

	"person.mgtv.com/framework/logs"
	"person.mgtv.com/framework/resultcode"
)

var mu sync.Mutex
var templateMapping = make(map[string]*template.Template)

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request)
	InitController()
	URLMapping()
	ExistMapping(method string) bool
	HandlerFunc(method string) bool
	Error(code string, err ...error)
	ExecuteTemplate(templatePath string)
	Before()
	After()
}

type ResponseData struct {
	Code    string                 `json:"code"`
	Message string                 `json:"msg"`
	Seq     string                 `json:"seq"`
	Data    map[string]interface{} `json:"data"`
}

type Controller struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	ResponseData   *ResponseData
	Timestamp      time.Time
	Execute        bool
	methodMapping  map[string]func()
}

// 初始化
func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	// 初始化信息
	c.ResponseWriter = w
	c.Request = r
	c.ResponseData = &ResponseData{
		Code: resultcode.SUCCESS,
		Seq:  time.Now().Format("2006-01-02 15:04:05"),
		Data: make(map[string]interface{}),
	}
	c.Timestamp = time.Now()
	c.Execute = true
	c.methodMapping = make(map[string]func())
}

// 方法映射
func (c *Controller) Mapping(method string, fn func()) {
	c.methodMapping[method] = fn
}

func (c *Controller) ExistMapping(method string) bool {
	if _, ok := c.methodMapping[method]; ok {
		return true
	}
	return false
}

func (c *Controller) HandlerFunc(method string) bool {
	if v, ok := c.methodMapping[method]; ok {
		v()
		return true
	}
	return false
}

// 默认成功响应
func (c *Controller) Success() {
	c.Execute = false
	c.ResponseData.Code = resultcode.SUCCESS
	body, _ := json.Marshal(c.ResponseData)
	io.WriteString(c.ResponseWriter, string(body))
}

// 异常响应
func (c *Controller) Error(code string, err ...error) {
	c.Execute = false
	if len(err) > 0 {
		c.ResponseData.Message = err[0].Error()
	} else {
		c.ResponseData.Message = resultcode.ErrorMsg(code)
	}
	c.ResponseData.Code = code
	c.ResponseData.Data = make(map[string]interface{})

	body, _ := json.Marshal(c.ResponseData)
	io.WriteString(c.ResponseWriter, string(body))

	logs.GetLogger("system").Error("Response error! ResponseData=%v", c.ResponseData)
}

// 参数集合
func (c *Controller) Data(key string, value interface{}) {
	c.ResponseData.Data[key] = value
}

// 执行模板方法
func (c *Controller) ExecuteTemplate(templatePath string) {
	if c.Execute {
		t := templateMapping[templatePath]
		if t == nil {
			mu.Lock()
			defer mu.Unlock()
			if t == nil {
				b, _ := ioutil.ReadFile("template" + templatePath)
				t, _ = template.New(templatePath).Parse(string(b))
				templateMapping[templatePath] = t
			}
		}
		t.Execute(c.ResponseWriter, c.ResponseData)
	}
}

// IO流输出
func (c *Controller) Write(body string) {
	c.Execute = false
	io.WriteString(c.ResponseWriter, body)
}

// 解析表单
func (c *Controller) ParseForm() url.Values {
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	return c.Request.Form
}

func (c *Controller) Get(key string) string {
	return c.ParseForm().Get(key)
}

func (c *Controller) GetInt(key string) int {
	v, _ := strconv.Atoi(c.Get(key))
	return v
}

func (c *Controller) GetInt64(key string) int64 {
	v, _ := strconv.ParseInt(c.Get(key), 10, 64)
	return int64(v)
}

func (c *Controller) GetFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(c.Get(key), 64)
	return v
}

// IP地址
func (c *Controller) IP() string {
	return c.Request.Header.Get("X-Forwarded-For")
}

// 逻辑前执行
func (c *Controller) Before() {
	// TODO Do something before...
}

// 逻辑后执行
func (c *Controller) After() {
	dis := (time.Now().Sub(c.Timestamp).Nanoseconds()) / 1e6
	// time|ip|request_uri|const
	access := fmt.Sprintf("%v|%v|%v|", c.IP(), c.Request.URL.Path, dis)
	logs.GetLogger("access").Info(access)
}
