package model

import "time"

type ArticleModel struct {
	Id      int
	Author  string
	Title   string
	Body    string
	Created time.Time
}
