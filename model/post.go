package model

import "time"

const (
	Blog = iota
	Puzzle
	Announce
)

// 帖子
type Post struct {
	//基础信息
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Head      string //标题
	Content   string //内容
	Kind      uint   //帖子类型(blog,puzzle,announce)

	//索引
	UserID   uint   //用户ID
	Comments []uint //所有的评
	ProCount uint   //点赞的数量
}

type Comment struct {
	//基础信息
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string

	//索引
	UserID int64   //用户ID
	PostID int64   //帖子ID
	Replys []int64 //索所有的回复

	ProCount uint //赞成的数量
	ConCount uint //反对的数量
}

type Reply struct {
	// 基础信息
	ID        int64
	CreatedAt time.Time
	Content   string
	//索引
	UserID     int64 //用户ID
	CommentID  int64 //评论ID
	PreReplyID int64 //当前所回复的上条回复

	ProCount uint //赞成的数量
	ConCount uint //反对的数量
}
