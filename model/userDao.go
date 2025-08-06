package model

import (
	"herbal_demo/dbconfig"
	"herbal_demo/utils"
)

func UserQuery(username string, password string, role string) (error, bool, *User) {
	hashPwd := utils.SHA256(password)
	//rows是结果集
	rows, err := dbconfig.DB.Query("select id, name, role from Users where name = ? and pwd = ? and role = ? ", username, hashPwd, role)
	defer rows.Close()
	//4.判断查询结果
	if err != nil {
		//context.JSON(http.StatusOK, gin.H{
		//	"code": "10003",
		//	"msg":  "服务器遇到错误，详细信息：" + err.Error(),
		//})
		return err, false, nil
	}
	//rows里面没有数据，说明库中没有对应的数据
	existFlag := rows.Next()
	if !existFlag {
		//context.JSON(http.StatusOK, gin.H{
		//	"code": "10006",
		//	"msg":  "登录失败，用户名或密码错误",
		//})
		return nil, false, nil
	}

	//匹配成功，查询到了数据
	//读取查询到的数据
	user := User{}
	rows.Scan(&user.Id, &user.Name)
	return nil, true, &user
}
