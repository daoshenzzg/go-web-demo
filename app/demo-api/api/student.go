package api

type Student struct {
	Id int64 `json:"id"`
	StudName string `json:"stud_name"`
	StudAge int64 `json:"stud_age"`
	StudSex string `json:"stud_sex"`
}