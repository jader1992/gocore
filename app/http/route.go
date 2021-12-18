package http

import (
	"github.com/jader1992/gocore/app/http/middleware/cors"
	"github.com/jader1992/gocore/app/http/module/demo"
	"github.com/jader1992/gocore/framework/gin"
	"github.com/jader1992/gocore/framework/middleware"
	"github.com/jader1992/gocore/framework/middleware/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	// 设置静态服务
	//r.Static("/dist/", "./dist/")
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	// 使用链路追踪中间件
	r.Use(middleware.Trace())
	// 使用cors中间件
	r.Use(cors.Default())
	demo.Register(r) // 注册路由
}
