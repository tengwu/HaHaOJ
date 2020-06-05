package model

import (
	"time"
)

type User struct {
	ID int64 `json:"id" xorm:"pk autoincr"`
	//基础信息
	CreatedAt    time.Time `json:"created_at" xorm:"created"`                        //创建时间
	Username     string    `json:"username" xorm:"varchar(64) unique index notnull"` //用户名
	Password     string    `json:"password" xorm:"varchar(16) notnull"`              //密码
	School       string    `json:"school" xorm:"varchar(64) notnull index"`          //学校
	Email        string    `json:"email"  xorm:"varchar(32) unique index notnull"`   //邮箱
	Description  string    `json:"description" xorm:"text"`                          //自我描述
	Avatar       string    `json:"avatar"`                                           //头像路径
	Rating       uint      `json:"rating" xorm:"index default 1000"`                 //获得的分数
	IsAdmin      bool      `json:"is_admin"`                                         //是否是管理员,管理员能进后台,然后存在相关权限, 普通用户无任何权限
	IsSuperAdmin bool      `json:"is_super_admin"`                                   //超级管理员,拥有任何权限
	//索引
	Privilege          PrivilegeType     `json:"privilege"`           //权限信息
	Messages           []int64           `json:"messages" `           //私信
	Notifications      []int64           `json:"notifications"`       //通知
	Posts              []int64           `json:"posts"`               //个人发表的帖子
	PostCollections    []int64           `json:"post_collections"`    //个人收藏的帖子
	PassedProblems     map[int64][]int64 `json:"passed_problems"`     //通过的题的id以及对应的所有提交索引id
	FailedProblems     map[int64][]int64 `json:"failed_problems"`     //还未解决的题id和对应的所有提交索引id
	ProblemCollections []int64           `json:"problem_collections"` //收藏的题目
	Follows            []int64           `json:"follows"`             //我关注的人
}
