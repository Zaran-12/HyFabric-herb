package utils

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveUserSession(context *gin.Context, msg string) {
	//需要将用户的信息设置到session中
	session := sessions.Default(context)
	//确保每次的key值唯一
	session.Set("username", msg)
	//sessionvalue := session.Get("username")
	//fmt.Println(sessionvalue)
	session.Save()
	err := session.Save()
	if err != nil {
		fmt.Println("Session 保存失败:", err)
	}

}
