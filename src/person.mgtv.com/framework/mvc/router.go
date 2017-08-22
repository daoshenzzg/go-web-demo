package mvc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"person.mgtv.com/framework/constants"
	"person.mgtv.com/framework/resultcode"
)

type ControllerHandler struct {
	routers map[string]interface{}
}

func NewControllerHandler() *ControllerHandler {
	return &ControllerHandler{
		routers: make(map[string]interface{}),
	}
}

type routerInfo struct {
	name           string
	controllerType reflect.Type
	tplName        string
}

func (handler *ControllerHandler) Add(name string, controller ControllerInterface) {
	reflectType := reflect.Indirect(reflect.ValueOf(controller)).Type()
	router := &routerInfo{}
	router.name = name
	router.controllerType = reflectType

	handler.routers[name] = router
}

func (handler *ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) < 3 {
		http.NotFound(w, r)
		return
	}
	router, method, uri := paths[1], paths[2], r.URL.Path

	// return 404 if router is not exists
	if _, ok := handler.routers[router]; !ok {
		http.NotFound(w, r)
		return
	}

	// recover error
	defer func() {
		if err := recover(); err != nil {
			handler.Error(w, r, err)
			debug.PrintStack()
		}
	}()

	r.ParseForm()

	// reflect view
	execController := handler.NewController(router)
	execController.Init(w, r)
	execController.InitController()
	execController.URLMapping()
	execController.Before()

	// return 404 if method is not exists
	if ok := execController.ExistMapping(method); !ok {
		http.NotFound(w, r)
		return
	}

	exec := execController.HandlerFunc(method)
	if !exec {
		execController.Error(resultcode.ERROR_1000)
		return
	}

	// 执行模板方法
	execController.ExecuteTemplate(uri + ".tpl")
	execController.After()
}

func (handler *ControllerHandler) Error(w http.ResponseWriter, r *http.Request, err interface{}) {
	e, _ := err.(string)

	responseData := &ResponseData{
		Code:    resultcode.ERROR,
		Seq:     time.Now().Format("2006-01-02 15:04:05"),
		Data:    make(map[string]interface{}),
		Message: e,
	}

	data := make(map[string]*ResponseData)
	data["Response"] = responseData

	val, _ := json.Marshal(data)
	io.WriteString(w, string(val))

	// 错误日志
	errmsg := fmt.Sprintf("%s|%s|%s",
		resultcode.ERROR,
		r.FormValue(constants.USER_ID),
		e,
	)
	fmt.Print(errmsg)
}

func (handler *ControllerHandler) NewController(name string) ControllerInterface {
	router := handler.routers[name]
	if router == nil {
		panic(fmt.Sprintf("Router[%s] is not exist", name))
	}

	var reflectType reflect.Type
	if rr, ok := router.(*routerInfo); ok {
		reflectType = rr.controllerType
	} else {
		panic("Router assertion failure")
	}

	controller, ok := reflect.New(reflectType).Interface().(ControllerInterface)
	if !ok {
		panic("Controller assertion failure")
	}
	return controller
}
