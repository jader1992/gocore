package demo

import (
	"github.com/jader1992/gocore/framework"
)

// TestProvider 服务提供方
type TestProvider struct {
	framework.ServiceProvider

	c framework.Container // 服务容器
}

// Name 方法直接将服务对应的字符串凭证返回，在这个例子中就是“jade.demo"
func (sp *TestProvider) Name() string {
	return DKey
}

// Register 方法是注册初始化服务实例的方法，我们这里先暂定为NewDemoService
func (sp *TestProvider) Register(c framework.Container) framework.NewInstance {
	return NewService
}

// IsDefer 方法表示是否延迟实例化，我们这里设置为true，将这个服务的实例化延迟到第一次make的时候
func (sp *TestProvider) IsDefer() bool {
	return false
}

func (sp *TestProvider) Params(c framework.Container) []interface{} {
	return []interface{}{sp.c}
}

func (sp *TestProvider) Boot(c framework.Container) error {
	sp.c = c
	return nil
}
