package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-web-demo/app/demo-api/api"
	"go-web-demo/app/demo-api/internal/bcode"
	"go-web-demo/app/demo-api/internal/model"
	"go-web-demo/library/ecode"
	"go-web-demo/library/log"
	xtime "go-web-demo/library/time"
	"time"
)

// ListStudent
func (s *Service) ListStudent(c context.Context, studName string) (studList []*api.Student, err error) {
	list, err := s.dao.ListStudent(c, studName)
	if err != nil {
		return studList, err
	}

	studList = make([]*api.Student, 0)
	for _, item := range list {
		sutd := &api.Student{
			Id:       item.Id,
			StudName: item.StudName,
			StudAge:  item.StudAge,
			StudSex:  item.StudSex,
		}
		studList = append(studList, sutd)
	}
	return studList, err
}

// AddStudent
func (s *Service) AddStudent(c context.Context, stud *model.Student) (id int64, err error) {
	id, err = s.dao.AddStudent(c, stud)
	if err != nil {
		err = bcode.OpStudentErr
		return
	}
	return
}

// UpdateStudentName
func (s *Service) UpdateStudentName(c context.Context, id int64, studName string) (err error) {
	stud := &model.Student{
		Id:       id,
		StudName: studName,
	}
	err = s.dao.UpdateStudentName(c, stud)
	if err != nil {
		err = bcode.OpStudentErr
		return
	}
	return
}

// TxAddTeacherAndStudent
func (s *Service) TxAddTeacherAndStudent(c context.Context) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran() error(%v)", err)
		return
	}
	teacher := &model.Teacher{
		TeacherName: "莫言",
		CreateTime:  xtime.Duration(time.Now().Unix()),
		UpdateTime:  xtime.Duration(time.Now().Unix()),
	}
	_, err = s.dao.TxAddTeacher(c, tx, teacher)
	if err != nil {
		log.Error("s.dao.TxAddTeacher(%v) error(%v)", teacher, err)
		tx.Rollback()
		return
	}
	stud := &model.Student{
		StudName:   "迪丽热巴",
		StudSex:    "女",
		StudAge:    25,
		CreateTime: xtime.Duration(time.Now().Unix()),
		UpdateTime: xtime.Duration(time.Now().Unix()),
	}
	_, err = s.dao.TxAddStudent(c, tx, stud)
	if err != nil {
		log.Error("s.dao.TxAddStudent(%v) error(%v)", stud, err)
		tx.Rollback()
		return
	}
	return tx.Commit()
}

// GetRedisKey
func (s *Service) GetRedisKey(c context.Context, key string) (val string, err error) {
	return s.dao.GetKey(c, key)
}

// SetRedisKey
func (s *Service) SetRedisKey(c context.Context, key, val string, expire int64) (err error) {
	return s.dao.SetKey(c, key, val, expire)
}

// SearchKeyword
func (s *Service) SearchKeyword(c *gin.Context) (keyword string, err error) {
	keyword, err = s.dao.SearchKeyword()
	if err != nil {
		err = ecode.RequestErr
		return
	}
	return
}
