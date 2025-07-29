package models

type Note struct {
	Header  string `json:"header"`
	Content string `json:"content"`
	Id      int64  `json:"id"`
}
