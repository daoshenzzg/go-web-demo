package dao

import (
	"context"
	"fmt"
	"go-web-demo/app/demo-api/internal/model"
	xsql "go-web-demo/library/database/sql"
	"go-web-demo/library/log"
	"time"
)

const (
	_updateStudentNameSQL = "UPDATE student SET stud_name=?, create_time=? WHERE id=?"
	_insertStudentSQL     = "INSERT INTO student(stud_name, stud_age, stud_sex, create_time, update_time)VALUES(?, ?, ?, ?, ?)"
)

// Student List
func (d *Dao) ListStudent(c context.Context, studName string) (studList []*model.Student, err error) {
	sql := "SELECT id, stud_name, stud_age, stud_sex "
	sql += "FROM student "
	if len(studName) > 0 {
		sql += fmt.Sprintf(" WHERE stud_name = '%s' ", studName)
	}
	sql += " LIMIT 10"
	studList = make([]*model.Student, 0)
	rows, err := d.db.Query(c, sql)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", sql, err)
		return
	}
	for rows.Next() {
		tmp := new(model.Student)
		err = rows.Scan(&tmp.Id, &tmp.StudName, &tmp.StudAge, &tmp.StudSex)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}
		studList = append(studList, tmp)
	}
	return
}

// AddStudent
func (d *Dao) AddStudent(c context.Context, stud *model.Student) (lastID int64, err error) {
	res, err := d.db.Exec(c, _insertStudentSQL, stud.StudName, stud.StudAge, stud.StudSex, stud.CreateTime, stud.UpdateTime)
	if err != nil {
		log.Error("AddStudent error(%v) student: %v", err, stud)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// TxAddStudent
func (d *Dao) TxAddStudent(c context.Context, tx *xsql.Tx, stud *model.Student) (lastID int64, err error) {
	res, err := tx.Exec(_insertStudentSQL, stud.StudName, stud.StudAge, stud.StudSex, stud.CreateTime, stud.UpdateTime)
	if err != nil {
		log.Error("TxAddStudent error(%v) student: %v", err, stud)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// UpdateStudentName
func (d *Dao) UpdateStudentName(c context.Context, stud *model.Student) (err error) {
	_, err = d.db.Exec(c, _updateStudentNameSQL, stud.StudName, time.Now().Unix(), stud.Id)
	if err != nil {
		log.Error("UpdateStudent error(%v) student: %v", err, stud)
	}
	return
}
