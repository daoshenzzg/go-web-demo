package model

import (
	"time"
)

type Teacher struct {
	Id          int64     `json:"id"`
	TeacherName string    `json:"teacher_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}
