package render

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/library/ecode"
	"net/http"
)

type Gin struct {
	c *gin.Context
}

// New create a new gin Context
func New(c *gin.Context) *Gin {
	return &Gin{
		c: c,
	}
}

// JSON common json struct.
type JSON struct {
	// 业务错误码
	Code int `json:"code"`
	// 错误描述
	Msg string `json:"msg"`
	// 响应数据
	Data interface{} `json:"data,omitempty"`
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (g *Gin) JSON(data interface{}, err error) {
	code := ecode.Cause(err)
	g.c.JSON(http.StatusOK, JSON{
		Code: code.Code(),
		Msg:  code.Message(),
		Data: data,
	})
}
