package demo

import (
	"github.com/jader1992/gocore/framework"
)

// DemoProvider 服务提供方
type DemoProvider struct {
	framework.ServiceProvider

	c framework.Container // 服务容器
}

// Name方法直接将服务对应的字符串凭证返回，在这个例子中就是“jade.demo"
func (sp *DemoProvider) Name() string {
	return DEMO_KEY
}

// Register方法是注册初始化服务实例的方法，我们这里先暂定为NewDemoService
func (sp *DemoProvider) Register(c framework.Container) framework.NewInstance {
	return NewService
}

// IsDefer方法表示是否延迟实例化，我们这里设置为true，将这个服务的实例化延迟到第一次make的时候
func (sp *DemoProvider) IsDefer() bool {
	return false
}

func (sp *DemoProvider) Params(c framework.Container) []interface{} {
	return []interface{}{sp.c}
}

func (sp *DemoProvider) Boot(c framework.Container) error {
	sp.c = c
	return nil
}
