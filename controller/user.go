package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"herbal_demo/dbconfig"
	"herbal_demo/model"
	"herbal_demo/service"
	"herbal_demo/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

func Register(context *gin.Context) {
	//获取用户提交的信息数据
	var user model.UserRegister
	err := context.ShouldBind(&user)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "10001",
			"msg":  "数据解析失败",
		})
		return
	}
	//查验
	isCheck, mapp := service.RegisterCheck(user.User, user.Pwd, user.RepeatPwd)
	if !isCheck {
		context.JSON(http.StatusOK, gin.H{
			"code": mapp["code"],
			"msg":  mapp["msg"],
		})
		return
	}

	uploadFile, header, err := context.Request.FormFile("avatar")
	defer uploadFile.Close()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "30001",
			"msg":  "文件解析失败",
		})
		return
	}

	//请增加完善头像更新逻辑，用户必须上传png/jpg/jpeg三种格式中的一种图片文件
	//如果上传其它类型的文件，则提示”不支持文件格式“
	fileNameArray := strings.Split(header.Filename, ".")
	fileType := fileNameArray[1]
	if fileType != "png" && fileType != "jpg" && fileType != "jpeg" {
		context.JSON(http.StatusOK, gin.H{
			"code": "30002",
			"msg":  "不支持的文件格式，请选择png/jpg/jpeg三种类型的文件上传",
		})
		return
	}
	//2 将文件保存到服务器端 【文件服务器】
	//a. 创建一个新文件，内容为空
	//hash：任意内容输出，固定长度的输出
	dir := "./uploadfile/"
	hash := utils.SHA256(header.Filename)
	path := dir + hash + "." + fileType
	newfile, err := os.Create(path)
	defer newfile.Close()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "30003",
			"msg":  "文件上传失败" + err.Error(),
		})
		return
	}
	//b.将用户上传的内容写到新创建的空文件中
	length, err := io.Copy(newfile, uploadFile)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "30004",
			"msg":  "文件上传失败" + err.Error(),
		})
		return
	}
	fmt.Println(length)

	//向数据库里存储该条记录
	//todo 先进行数据库条件查询。 若数据库已经存 在用户名为即将注册的用户名的记录， 则不进行数据插入操作， 直接提示用户名已被占用
	rows, err := dbconfig.DB.Query("select id from Users where name = ?", user.User)
	//延迟执行
	defer rows.Close()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "10005",
			"msg":  "服务器遇到错误，详细信息：" + err.Error(),
		})
		return
	}
	//如果可以查到数据，则rows.next的值为true
	if rows.Next() {
		context.JSON(http.StatusOK, gin.H{
			"code": "10006",
			"msg":  "用户名已存在，请使用其他用户名!",
		})
		return
	}

	//1.首先调用SHA256函数对密码进行脱敏处理
	hashPwd := utils.SHA256(user.Pwd)
	rs, err := dbconfig.DB.Exec("insert into Users (name , pwd , role) values(?,?,?)", user.User, hashPwd, user.Role)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "10007",
			"msg":  "注册失败，请检查" + err.Error(),
		})
	}
	num, _ := rs.RowsAffected()
	if num != 1 {
		context.JSON(http.StatusOK, gin.H{
			"code": "10008",
			"msg":  "注册异常，请检查",
		})
	}
	//c.将保存的新文件的路径+文件名 信息保存到数据库表中去，高新到指定用户的头像字段

	rs, err = dbconfig.DB.Exec("update Users set avatar = ? where name = ?", path, user.User)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "30005",
			"msg":  "更新头像失败" + err.Error(),
		})
		return
	}

	num, _ = rs.RowsAffected()
	if num != 1 {
		context.JSON(http.StatusOK, gin.H{
			"code": "30006",
			"msg":  "头像更新失败",
		})
		return
	}

	//将处理结果返回前端
	context.JSON(http.StatusOK, gin.H{
		"code": "10000",
		"msg":  "注册成功",
	})

}

func Login(context *gin.Context) {
	//1.解析数据
	var login model.UserLogin
	err := context.ShouldBind(&login)
	fmt.Println(login.Name, login.Pwd)
	if err != nil {
		fmt.Println(err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": "20001",
			"msg":  "数据解析失败",
		})
		return
	}
	//2.参数检查
	isCheck, mapp := service.LoginCheck(login.Name, login.Pwd)
	if !isCheck {
		context.JSON(http.StatusOK, gin.H{
			"code": mapp["code"],
			"msg":  mapp["msg"],
		})
		return
	}
	//3.数据库查询
	err, isExist, user := model.UserQuery(login.Name, login.Pwd, login.Role)
	//数据库遇到错误
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": "20004",
			"msg":  "操作遇到错误，错误信息" + err.Error(),
		})
		return
	}
	//用户登录失败
	if !isExist {
		context.JSON(http.StatusOK, gin.H{
			"code": "20005",
			"msg":  "用户,密码或角色错误，登录失败",
		})
		return
	}

	//需要将用户的信息设置到session中
	utils.SaveUserSession(context, login.Name)

	// 5. 根据用户角色设置跳转路径
	roleToURL := map[string]string{
		"农户":   "/farmer_dashboard.html",
		"用户":   "/home.html",
		"企业":   "/enterprise_dashboard.html",
		"监管机构": "/regulator_dashboard.html",
	}

	redirectURL, exists := roleToURL[login.Role]
	fmt.Println("Redirecting to:", redirectURL)
	if !exists {
		context.JSON(http.StatusForbidden, gin.H{
			"code": "40003",
			"msg":  "无效的用户角色",
		})
		return
	}

	// 6. 返回登录成功提示和跳转路径
	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "登录成功",
		"data": gin.H{
			"user":        user,
			"redirectURL": redirectURL,
		},
	})
}

func GetSessionData(context *gin.Context) {
	session := sessions.Default(context)

	// 从 session 中获取用户名
	username := session.Get("username")
	fmt.Println(username)
	if username == nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": "40001",
			"msg":  "未登录或 session 已过期",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "获取成功",
		"data": gin.H{"username": username},
	})
}

func GetAvatar(context *gin.Context) {
	username := context.Query("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "用户名不能为空",
		})
		return
	}

	var avatarPath string
	err := dbconfig.DB.QueryRow("SELECT avatar FROM Users WHERE name = ?", username).Scan(&avatarPath)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"code": "40002",
			"msg":  "用户不存在或头像未设置",
		})
		return
	}
	// 打印文件路径进行调试
	fmt.Println("头像路径:", avatarPath)

	// 打开文件
	file, err := os.Open(avatarPath)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40003",
			"msg":  "无法加载头像文件",
		})
		return
	}
	defer file.Close()

	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "查询成功",
		"data": gin.H{"avatar": avatarPath},
	})
}

func UpdateUserName(context *gin.Context) {
	// 从 session 获取当前用户名
	session := sessions.Default(context)
	oldUsername := session.Get("username")
	//fmt.Println(oldUsername)
	if oldUsername == nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": "40001",
			"msg":  "未登录",
		})
		return
	}

	// 解析前端传来的新用户名
	var form struct {
		NewName string `json:"new_name"`
	}
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40002",
			"msg":  "参数错误",
		})
		return
	}

	// 更新数据库中的用户名
	result, err := dbconfig.DB.Exec("UPDATE Users SET name = ? WHERE name = ?", form.NewName, oldUsername)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "50001",
			"msg":  "更新失败：" + err.Error(),
		})
		return
	}

	// 检查是否成功更新
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40003",
			"msg":  "更新失败，可能用户名不存在",
		})
		return
	}

	// 更新 session 中的用户名
	session.Set("username", form.NewName)
	session.Save()

	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "用户名更新成功",
	})
}

func UpdateAvatar(context *gin.Context) {
	username := context.PostForm("username")
	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "用户名不能为空",
		})
		return
	}

	file, header, err := context.Request.FormFile("avatar")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40002",
			"msg":  "头像上传失败，请选择正确的文件",
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	fileNameArray := strings.Split(header.Filename, ".")
	fileType := fileNameArray[1]
	if fileType != "png" && fileType != "jpg" && fileType != "jpeg" {
		context.JSON(http.StatusOK, gin.H{
			"code": "10009",
			"msg":  "不支持的文件格式，请选择png/jpg/jpeg三种类型的文件上传",
		})
		return
	}

	// 使用哈希生成文件名
	hashFileName := fmt.Sprintf("%s_avatar.%s", utils.SHA256(header.Filename+username), fileType)
	savePath := fmt.Sprintf("./uploadfile/%s", hashFileName)
	fmt.Println(savePath)

	// 保存文件到服务器
	out, err := os.Create(savePath)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40004",
			"msg":  "服务器保存头像失败",
		})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40005",
			"msg":  "头像保存失败",
		})
		return
	}

	// 更新数据库
	_, err = dbconfig.DB.Exec("UPDATE Users SET avatar = ? WHERE name = ?", savePath, username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40006",
			"msg":  "更新数据库失败",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "头像更新成功",
	})
}

func UpdatePassword(context *gin.Context) {
	// 获取用户名和密码数据
	username := context.PostForm("username")
	oldPassword := context.PostForm("old_password")
	newPassword := context.PostForm("new_password")
	confirmPassword := context.PostForm("confirm_password")

	// 检查用户名和密码字段是否为空
	if username == "" || oldPassword == "" || newPassword == "" || confirmPassword == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40001",
			"msg":  "所有字段均为必填",
		})
		return
	}

	// 检查新密码和确认密码是否匹配
	if newPassword != confirmPassword {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "40002",
			"msg":  "新密码与确认密码不匹配",
		})
		return
	}

	// 旧密码哈希化
	hashedOldPassword := utils.SHA256(oldPassword)

	// 从数据库中验证旧密码
	var storedPassword string
	err := dbconfig.DB.QueryRow("SELECT pwd FROM Users WHERE name = ?", username).Scan(&storedPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40003",
			"msg":  "用户不存在或查询失败",
		})
		return
	}

	// 验证旧密码是否正确
	if hashedOldPassword != storedPassword {
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": "40004",
			"msg":  "旧密码错误",
		})
		return
	}

	// 新密码哈希化
	hashedNewPassword := utils.SHA256(newPassword)

	// 更新数据库中的密码
	_, err = dbconfig.DB.Exec("UPDATE Users SET pwd = ? WHERE name = ?", hashedNewPassword, username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": "40005",
			"msg":  "更新密码失败，请稍后重试",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code": "20000",
		"msg":  "密码更新成功",
	})
}
