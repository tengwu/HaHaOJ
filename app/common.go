package app

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	USER_SESSION_KEY = "who"
	EXPIRE           = time.Hour * 24 * 3
	PREFIX           = "user_"
)

//根据session 验证是否登陆
func isLogin(c *gin.Context) bool {
	return getUserName(c) != ""
}

//根据session获取用户名
func getUserName(c *gin.Context) string {
	session := sessions.Default(c)
	username := session.Get(USER_SESSION_KEY)
	if username == nil {
		return ""
	}
	return username.(string)
}

//设置session
func setSession(c *gin.Context, val string) {
	session := sessions.Default(c)
	session.Set(USER_SESSION_KEY, val)
	session.Save()
}

//删除session
func deleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(USER_SESSION_KEY)
	session.Save()
}
