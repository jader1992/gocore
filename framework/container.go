package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provide ServiceProvider) error

	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字获取一个服务
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic
	// 所以在使用这个接口的时候保证服务容器已经为这个关键字凭证绑定了服务提供者
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// GocoreContainer 是服务容器的具体实现，这里的方法会在core.Bind时绑定
type GocoreContainer struct {
	Container

	providers map[string]ServiceProvider // 存储注册的服务提供者，key为字符串凭证
	instances map[string]interface{} // 存储具体的实例，key为字符串凭证

	lock sync.RWMutex
}

// NewHadeContainer 创建一个服务容器
func NewHadeContainer() *GocoreContainer {
	return &GocoreContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock: sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (gocore *GocoreContainer) PrintProvides() []string {
	ret := []string{}
	for _, provide := range gocore.providers {
		name := provide.Name()

		line := fmt.Sprintf(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind 将服务容器和关键字做了绑定
func (gocore *GocoreContainer) Bind(provide ServiceProvider) error {
    gocore.lock.Lock()

	key := provide.Name()
    gocore.providers[key] = provide

    gocore.lock.Unlock()

	// if providers is not defer
	if provide.IsDefer() == false {
		if err := provide.Boot(gocore); err != nil {
			return err // 一些准备工作会报错
		}

		// 实例化方法
		params := provide.Params(gocore)
		method := provide.Register(gocore)
		instance, err := method(params...)
		if err != nil {
			fmt.Println("bind service provider ", key, " error: ", err)
		}
        gocore.instances[key] = instance
	}
	return nil
}

func (gocore *GocoreContainer) IsBind(key string) bool {
	return gocore.findServiceProvider(key) != nil
}

func (gocore *GocoreContainer) findServiceProvider(key string) ServiceProvider {
    gocore.lock.RLock()
	defer gocore.lock.RUnlock()
	if sp, ok := gocore.providers[key]; ok {
		return sp
	}
	return nil
}

func (gocore *GocoreContainer) Make(key string) (interface{}, error) {
	return gocore.make(key, nil, false)
}

func (gocore *GocoreContainer) MustMake(key string) interface{} {
	serv, err := gocore.make(key, nil, false)
	if err != nil {
		panic("container not contain key " + key)
	}
	return serv
}

func (gocore *GocoreContainer) MakeNew(key string, params []interface{}) (interface{}, error)  {
	return gocore.make(key, params, true)
}

// 初始化实例
func (gocore *GocoreContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	err := sp.Boot(gocore);
	if err != nil {
		return nil, err
	}

	if params == nil {
		// 获取默认参数
		params = sp.Params(gocore)
	}

	method := sp.Register(gocore)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// 真正实例化一个服务
func (gocore *GocoreContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error)  {
    gocore.lock.RLock()
	defer gocore.lock.RUnlock()

	sp := gocore.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	// 强制初始化
	if forceNew {
		return gocore.newInstance(sp, params)
	}

	// 不需要强制初始化，如果容器中已经实例化了，那么就直接使用容器中的实例
	if ins, ok := gocore.instances[key]; ok {
		return ins, nil
	}

	// 容器中还未实例化，则进行第一次实例化【单例的使用】
	inst, err := gocore.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	// 保存实例
    gocore.instances[key] = inst
	return inst, nil
}

func (gocore *GocoreContainer) NameList() []string {
    ret := []string{}
    for _, provider := range gocore.providers {
        name := provider.Name()
        ret = append(ret, name)
    }
    return ret
}
