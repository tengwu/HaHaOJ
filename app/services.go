package app

import (
	"HappyOnlineJudge/dao"
	"github.com/gin-gonic/gin"
)

//路由处理
func ping(c *gin.Context) {
	c.String(200, "pong")
}

//登陆请求
func login(c *gin.Context) {
	if isLogin(c) {
		deleteSession(c)
	}
	form := new(loginValidtor)
	if err := c.ShouldBind(form); err != nil {
		c.String(403, err.Error())
		return
	}
	if ok, errInfo := form.isOk(); !ok {
		c.String(403, errInfo)
		return
	}
	ud := &dao.UserDao{Username: form.Username}
	id := ud.GetID()
	if id <= 0 {
		c.String(403, "用户名不存在")
		return
	}
	if pwd := dao.OneCol(ud, "password").ToString(); pwd != form.Password {
		c.String(403, "密码错误")
		return
	}
	if dao.IsInRedis(ud) {
		dao.GetSelfAll(ud)
	}
	dao.PutToRedisIfNotIn(ud)
	setSession(c, ud.Username)
	c.String(200, "登陆成功")
}

func logout(c *gin.Context) {
	deleteSession(c)
	c.String(200, "退出成功")
}

//注册请求
func register(c *gin.Context) {
	if isLogin(c) {
		c.String(403, "请退出当前用户")
		return
	}
	form := new(registerValidtor)
	if err := c.ShouldBind(form); err != nil {
		c.String(403, err.Error())
		return
	}
	if ok, errInfo := form.isOk(); !ok {
		c.String(403, errInfo)
		return
	}
	if dao.Count(new(dao.UsersData), []string{"username"}, []interface{}{form.Username}) > 0 {
		c.String(403, "用户名已存在")
		return
	}
	if dao.Count(new(dao.UsersData), []string{"email"}, []interface{}{form.Email}) > 0 {
		c.String(403, "邮箱已被注册")
		return
	}
	ud := &dao.UserDao{
		User: &dao.User{
			Username: form.Username,
			Password: form.Password,
			School:   form.School,
			Email:    form.Email,
		},
	}
	if err := dao.Create(ud); err != nil {
		c.String(500, err.Error())
		return
	}
	setSession(c, form.Username)
	c.String(200, "注册成功")
}

//更新用户的信息
func update(c *gin.Context) {
	form := new(updateValidtor)
	if err := c.ShouldBind(form); err != nil {
		c.String(403, err.Error())
		return
	}
	if ok, errInfo := form.isOk(); !ok {
		c.String(403, errInfo)
		return
	}
	name := getUserName(c)
	mp := make(map[string]interface{}) //要修改的内容
	ud := &dao.UserDao{Username: name}
	if form.Username != "" && form.Username != name {
		if dao.Count(new(dao.UsersData), []string{"username"}, []interface{}{form.Username}) > 0 {
			c.String(403, "用户名已存在")
			return
		}
		mp["username"] = form.Username
	}
	cols := dao.Cols(ud, "password", "email")
	if form.NewPassword != "" {
		if form.OldPassword != cols[0].ToString() {
			c.String(403, "密码错误")
			return
		}
		mp["password"] = form.NewPassword
	}

	if form.Email != "" && form.Email != cols[1].ToString() {
		if dao.Count(new(dao.UsersData), []string{"email"}, []interface{}{form.Email}) > 0 {
			c.String(403, "邮箱已被注册")
			return
		}
		mp["email"] = form.Email
	}

	if form.School != "" {
		mp["school"] = form.School
	}
	if form.Desc != "" {
		mp["description"] = form.Desc
	}
	if len(mp) > 0 {
		if err := ud.Update(mp); err != nil {
			c.String(500, err.Error())
			return
		}
	}
	c.String(200, "修改成功")
}
