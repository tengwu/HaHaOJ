package app

import (
	"github.com/gin-gonic/gin"
)

//中间件

//验证是否登陆
func authLogin(c *gin.Context) {
	if !isLogin(c) {
		c.String(401, "未登陆")
		c.Abort()
		return
	}
}

func authAdmin(c *gin.Context) {

}
