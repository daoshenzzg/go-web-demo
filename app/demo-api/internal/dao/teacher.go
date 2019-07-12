package dao

import (
	"context"
	"go-web-demo/app/demo-api/internal/model"
	xsql "go-web-demo/library/database/sql"
	"go-web-demo/library/log"
)

const (
	_insertTeacherSQL = "INSERT INTO teacher(teacher_name, create_time, update_time)VALUES(?, ?, ?)"
)

// TxAddTeacher
func (d *Dao) TxAddTeacher(c context.Context, tx *xsql.Tx, teacher *model.Teacher) (lastID int64, err error) {
	res, err := tx.Exec(_insertTeacherSQL, teacher.TeacherName, teacher.CreateTime, teacher.UpdateTime)
	if err != nil {
		log.Error("TxAddTeacher error(%v) teacher: %v", err, teacher)
		return
	}
	lastID, err = res.LastInsertId()
	return
}
