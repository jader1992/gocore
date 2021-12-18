// Copyright 2021 jade1992.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/jader1992/gocore/app/console"
	"github.com/jader1992/gocore/app/http"
	"github.com/jader1992/gocore/app/provider/demo"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/provider/app"
	"github.com/jader1992/gocore/framework/provider/config"
	"github.com/jader1992/gocore/framework/provider/distributed"
	"github.com/jader1992/gocore/framework/provider/env"
	"github.com/jader1992/gocore/framework/provider/id"
	"github.com/jader1992/gocore/framework/provider/kernel"
	"github.com/jader1992/gocore/framework/provider/log"
	"github.com/jader1992/gocore/framework/provider/trace"
)

func main() {
	// 初始化服务容器
	container := framework.NewHadeContainer()

	// 绑定App服务提供者
	container.Bind(&app.GocoreAppProvider{})
	container.Bind(&demo.DemoProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&distributed.LocalDistributedProvider{}) // 分布式定时任务
	container.Bind(&env.GocoreEnvProvider{})                // ENV相关
	container.Bind(&config.GocoreConfigProvider{})          // config相关
	container.Bind(&log.GocoreLogServiceProvider{})         // 日志文件相关
	container.Bind(&id.GocoreIDProvider{})                  // id生成器
	container.Bind(&trace.GocoreTraceProvider{})            // 链路追踪

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.GocoreKernelProvider{
			HttpEngine: engine,
		})
	}

	console.RunCommand(container)
}
