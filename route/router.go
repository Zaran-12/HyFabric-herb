package route

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"herbal_demo/controller" // 替换为实际的 controllers 路径
	//	"os"
)

func init() {
	//开始编了。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。。
	engine := gin.Default()

	//引入session redis管理
	cookStore, err := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte("store"))
	cookStore.Options(sessions.Options{
		MaxAge: 315360000, // 设置 session 有效期为 3600 秒
	})
	if err != nil {
		fmt.Println("Error initializing Redis store:", err)
		panic(err)
	}
	engine.Use(sessions.Sessions("mysession", cookStore))

	//用户模块
	userGroup := engine.Group("/user")
	userGroup.POST("/register", controller.Register)
	userGroup.POST("/login", controller.Login)
	userGroup.GET("/get-session", controller.GetSessionData)
	userGroup.POST("/updateUserName", controller.UpdateUserName)
	userGroup.GET("/get-avatar", controller.GetAvatar)
	userGroup.POST("/update-avatar", controller.UpdateAvatar)
	userGroup.POST("/update-password", controller.UpdatePassword)
	userGroup.GET("/combinedTrace", controller.CombinedTraceQuery)

	//农户模块
	farmerGroup := engine.Group("/farmer")
	farmerGroup.POST("/UploadBatchData", controller.UploadBatchData)
	farmerGroup.GET("/getherbalrecords", controller.GetHerbalRecords)
	farmerGroup.DELETE("/deleteherbal/:batchID", controller.DeleteHerbal)

	//企业模块
	enterGroup := engine.Group("/enterprise")
	enterGroup.POST("/UploadGoodsData", controller.UploadGoodsData)
	enterGroup.GET("/getgoodsrecords", controller.GetGoodsRecords)
	enterGroup.DELETE("/deletegoods/:batchID", controller.DeleteGoods)

	//监管部门模块
	regulaGroup := engine.Group("/regulator")
	regulaGroup.GET("/combinedTrace", controller.CombinedTraceQuery)

	//文件
	engine.Static("/uploadfile", "./uploadfile")

	//方式二： 加载 views 目录下的所有html文件。如果还有下级目录，则为 views/**/*
	engine.LoadHTMLGlob("views/*")

	// 设置静态文件路由   将 html 文件中的请求路径【/asset】 映射到 【asset】目录下
	engine.Static("asset", "asset")
	// 将 html 文件中的请求路径【/a/b/c】 映射到 【asset/css】目录下
	engine.Static("/a/b/c", "asset/css")
	//加载图片
	engine.Static("/img", "./img")

	// 设置路由以提供HTML页面
	//用户组
	engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", gin.H{})
	})
	engine.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	engine.GET("/home.html", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})
	engine.GET("/register.html", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})
	engine.GET("/batch-query.html", func(c *gin.Context) {
		c.HTML(200, "batch-query.html", nil)
	})
	engine.GET("/person.html", func(c *gin.Context) {
		c.HTML(200, "person.html", nil)
	})

	//农户组
	engine.GET("/farmer_dashboard.html", func(c *gin.Context) {
		c.HTML(200, "farmer_dashboard.html", nil)
	})
	engine.GET("/records_farmer.html", func(c *gin.Context) {
		c.HTML(200, "records_farmer.html", nil)
	})
	engine.GET("/person_farmer.html", func(c *gin.Context) {
		c.HTML(200, "person_farmer.html", nil)
	})

	//企业组
	engine.GET("/enterprise_dashboard.html", func(c *gin.Context) {
		c.HTML(200, "enterprise_dashboard.html", nil)
	})
	engine.GET("/manage_enterprise.html", func(c *gin.Context) {
		c.HTML(200, "manage_enterprise.html", nil)
	})
	engine.GET("/person_enterprise.html", func(c *gin.Context) {
		c.HTML(200, "person_enterprise.html", nil)
	})

	//监管组
	engine.GET("/regulator_dashboard.html", func(c *gin.Context) {
		c.HTML(200, "regulator_dashboard.html", nil)
	})
	engine.GET("/trace_regulator.html", func(c *gin.Context) {
		c.HTML(200, "trace_regulator.html", nil)
	})
	engine.GET("/exceptions_regulator.html", func(c *gin.Context) {
		c.HTML(200, "exceptions_regulator.html", nil)
	})
	engine.GET("/person_regulator.html", func(c *gin.Context) {
		c.HTML(200, "person_regulator.html", nil)
	})

	engine.Run(":4000")
}
