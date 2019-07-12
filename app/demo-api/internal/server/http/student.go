package http

import (
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"go-web-demo/app/demo-api/internal/model"
	"go-web-demo/library/ecode"
	"go-web-demo/library/render"
	xtime "go-web-demo/library/time"
	"time"
)

// @Summary 学生列表
// @Produce json
// @Param studName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/list [get]
func ListStudent(c *gin.Context) {
	r := render.New(c)

	studName := c.Query("studName")

	studList, err := srv.ListStudent(c, studName)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	r.JSON(studList, ecode.OK)
}

// @Summary 添加学生
// @Produce json
// @Param studName query string true "学生姓名"
// @Param studAge query int true "年龄"
// @Param studSex query string true "性别"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/add [post]
func AddStudent(c *gin.Context) {
	r := render.New(c)

	v := new(struct {
		StudName string `form:"studName" binding:"required,min=1,max=30"`
		StudAge  int64  `form:"studAge" binding:"required,min=1,max=60"`
		StudSex  string `form:"studSex" binding:"required,min=1,max=1"`
	})

	err := c.Bind(v)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}

	stud := &model.Student{
		StudName:   v.StudName,
		StudAge:    v.StudAge,
		StudSex:    v.StudSex,
		CreateTime: xtime.Duration(time.Now().Unix()),
		UpdateTime: xtime.Duration(time.Now().Unix()),
	}
	id, err := srv.AddStudent(c, stud)
	r.JSON(id, err)
}

// @Summary 修改学生
// @Produce json
// @Param id query int true "学生编号"
// @Param StudName query string true "学生姓名"
// @Success 200 {object} render.JSON
// @Router /api/v1/student/update [post]
func UpdateStudentName(c *gin.Context) {
	r := render.New(c)

	v := new(struct {
		Id       int64  `form:"id" binding:"required,gte=1"`
		StudName string `form:"studName" binding:"required,min=1,max=30"`
	})

	err := c.Bind(v)
	if err != nil {
		r.JSON(nil, ecode.RequestErr)
		return
	}
	err = srv.UpdateStudentName(c, v.Id, v.StudName)
	r.JSON(nil, err)
}

func TxAddTeacherAndStudent(c *gin.Context) {
	r := render.New(c)
	err := srv.TxAddTeacherAndStudent(c)
	r.JSON(nil, err)
}

func GetRedisKey(c *gin.Context) {
	r := render.New(c)
	key := c.Query("key")
	if key == "" {
		r.JSON(nil, ecode.RequestErr)
		return
	}
	val, err := srv.GetRedisKey(c, key)
	r.JSON(val, err)
}

func SetRedisKey(c *gin.Context) {
	r := render.New(c)
	key := c.Query("key")
	if len(key) == 0 {
		r.JSON(nil, ecode.RequestErr)
		return
	}
	val := c.Query("val")
	if len(val) == 0 {
		r.JSON(nil, ecode.RequestErr)
		return
	}
	expire := com.StrTo(c.Query("expire")).MustInt64()
	if expire <= 0 || expire > 300 {
		r.JSON(nil, ecode.RequestErr)
		return
	}
	err := srv.SetRedisKey(c, key, val, expire)
	r.JSON(nil, err)
}

func SearchKeyword(c *gin.Context) {
	r := render.New(c)
	keyword, err := srv.SearchKeyword(c)
	r.JSON(keyword, err)
}
