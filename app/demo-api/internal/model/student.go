package model

import (
	xtime "go-web-demo/library/time"
)

type Student struct {
	Id int64 `json:"id"`
	StudName string `json:"stud_name"`
	StudAge int64 `json:"stud_age"`
	StudSex string `json:"stud_sex"`
	CreateTime xtime.Duration `json:"create_time"`
	UpdateTime xtime.Duration `json:"update_time"`
}