package dbconfig

import (
	"database/sql"
	"fmt"
)

//在代码中导入驱动
import _ "github.com/go-sql-driver/mysql"

var DB *sql.DB

func init() {
	//告知 数据库用户名、密码、数据库名等信息 告知驱动
	conn := "root:Good8*Finalverse@tcp(127.0.0.1:3306)/herbal?charset=utf8"
	var err error
	//建立链接
	DB, err = sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
		return
	}
	fmt.Println("数据库连接......", DB)
}
