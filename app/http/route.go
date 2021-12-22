package http

import (
	"github.com/jader1992/gocore/app/http/middleware/cors"
	"github.com/jader1992/gocore/app/http/module/demo"
    "github.com/jader1992/gocore/framework/contract"
    "github.com/jader1992/gocore/framework/gin"
	"github.com/jader1992/gocore/framework/middleware"
    ginSwagger "github.com/jader1992/gocore/framework/middleware/gin-swagger"
    "github.com/jader1992/gocore/framework/middleware/gin-swagger/swaggerFiles"
    "github.com/jader1992/gocore/framework/middleware/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	// 设置静态服务
	//r.Static("/dist/", "./dist/")
    container := r.GetContainer()
    configService := container.MustMake(contract.ConfigKey).(contract.Config)

	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	// 使用链路追踪中间件
	r.Use(middleware.Trace())
	// 使用cors中间件
	r.Use(cors.Default())

    // 如果配置了swagger，则显示swagger的中间件
    if configService.GetBool("app.swagger") == true {
        // 原理：是寻找main下的doc.go, 本项目为什么可以生效： 因为main 包含了 app/http， 而app/http/swagger.go 由包含了./swagger
        r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }

	_ = demo.Register(r) // 注册路由
}
