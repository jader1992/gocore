package http

import (
    "github.com/jader1992/gocore/framework"
    "github.com/jader1992/gocore/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	// 设置gin的模式: 设置为Release，为的是默认在启动中不输出调试信息
	gin.SetMode(gin.ReleaseMode)
	//  默认的engine
	//r := gin.Default()
	r := gin.New()
    // 设置了engine
    r.SetContainer(container)

	r.Use(gin.Recovery()) // 返回一个可以从任何恐慌中恢复的中间件

	// 业务绑定路由操作
	Routes(r)
	// 返回绑定路由后的Web引擎
	return r, nil
}
