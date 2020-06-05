package dao

import (
	"HappyOnlineJudge/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
	"xorm.io/core"
)

type (
	User    = model.User
	Post    = model.Post
	Message = model.Message
	Comment = model.Comment
	Reply   = model.Reply
)

const (
	//user 相关
	USER_REDIS_EXPIRE = time.Hour * 24 * 3 //用户在redis的超时时间
	USER_REDIS_PREFIX = "user_"

	//TIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	superAdminAccount  = "super_admin" //超级管理员账号
	superAdminPassword = "111111"
	superAdminEmail    = "562954019@qq.com"
	engine             *xorm.Engine    //数据库引擎(这里用的mysql)
	rdb                *redis.Client   //redis
	ctx                context.Context //默认值
)

//连接mysql数据库和redis
func connect() error {
	var err error

	//数据库连接
	engine, err = xorm.NewEngine("mysql", "root:root@/hoj_db?charset=utf8")
	if err != nil {
		return err
	}
	err = engine.Ping()
	if err != nil {
		return err
	}
	engine.SetMapper(core.GonicMapper{})
	//redis连接
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx = context.TODO()
	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		return err
	} else {
		fmt.Println(pong, err)
	}
	return nil
}

// mysql表同步和redis初始化
func sync() error {
	if err := engine.Sync2(new(User)); err != nil {
		return err
	}
	//自动设置管理员账号
	ud := &UserDao{Username: superAdminAccount}
	if ud.GetID() == 0 {
		ud.User = &User{
			Username:     superAdminAccount,
			Password:     superAdminPassword,
			IsSuperAdmin: true,
			Email:        superAdminEmail,
		}
		if err := Create(ud); err != nil {
			return err
		}
		fmt.Println("超级管理初始化创建完成!!!")
	}
	return nil
}

func InitDao() error {
	if err := connect(); err != nil {
		return err
	}
	if err := sync(); err != nil {
		return err
	}
	return nil
}
