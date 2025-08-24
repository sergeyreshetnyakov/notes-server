package models

type Note struct {
	Header  string `json:"header" example:"go for a walk"`
	Content string `json:"content" example:"at 3 pm"`
	Id      int64  `json:"id" example:"1"`
}
