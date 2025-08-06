package model

// 将查询的结果返回终端
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type UserLogin struct {
	Name string `form:"name" json:"name" binding:"required"`
	Pwd  string `form:"pwd" json:"pwd" binding:"required"`
	Role string `form:"role" binding:"required"`
}
type UserRegister struct {
	User      string `form:"user" binding:"required"`
	Pwd       string `form:"pwd" binding:"required"`
	RepeatPwd string `form:"repeatpwd" binding:"required"`
	Role      string `form:"role" binding:"required"`
}
