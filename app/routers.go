package app

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

var r *gin.Engine

//路由
func InitRouters() {
	r = gin.Default()
	store := cookie.NewStore([]byte("secret")) //启用cookie和session
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24 * 3, //3天的过期时间
	})
	r.Use(sessions.Sessions("ginSession", store))
	initUserRouters() //初始化用户相关路由
	if err := r.Run(":8888"); err != nil {
		fmt.Println("路由初始化错误\n", err.Error())
	}
}

func initUserRouters() {

	g0 := r.Group("/") // 无需任何条件的请求
	{
		g0.GET("/ping", ping)
		g0.POST("login", login)
		g0.POST("register", register)
	}
	g1 := r.Group("/") //需要登陆才能进行的请求
	g1.Use(authLogin)  //authLogin 登陆验证中间件
	{
		g1.GET("logout", logout)
		g1.POST("update", update)
	}

	g2 := r.Group("/")
	g2.Use(authAdmin) //需要管理员才能进行的请求
	{

	}
	r.StaticFS("/statics", http.Dir("./statics"))
}
