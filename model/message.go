package model

import "time"

//通知消息
type Message struct {
	ID        int64     `json:"id" xorm:"pk autoincr"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	Head      string    `json:"head" xorm:"varchar(64)"` //消息头
	Content   string    `json:"content" xorm:"text"`     //消息内容
}
