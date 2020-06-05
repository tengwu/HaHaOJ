package main

import (
	"HappyOnlineJudge/app"
	"HappyOnlineJudge/dao"
	"fmt"
)

func main() {
	if err := dao.InitDao(); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("数据库初始化完成")
	}
	app.InitRouters()
	//if err := dao.Test(); err != nil {
	//	fmt.Println(err.Error())
	//	fmt.Println("测试失败")
	//} else {
	//	fmt.Println("测试通过")
	//}
}
