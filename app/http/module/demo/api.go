package demo

import (
	demoService "github.com/jader1992/gocore/app/provider/demo"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/gin"
)

func Register(r *gin.Engine) error {
	// 绑定demoProvider提供者
	// r.Bind(&demoService.TestProvider{})

	// 注册路由
	api := NewDemoApi()
	r.GET("/demo/demo", api.Demo)
	r.GET("/demo/demo2", api.Demo2)
	r.POST("/demo/demo_post", api.DemoPost)
    r.GET("/demo/orm", api.DemoOrm)
    r.GET("/demo/cache/redis", api.DemoRedis)
    r.GET("/demo/cache/cache", api.DemoCache)
	return nil
}

// Api 测试api的提供者
type Api struct {
	service *Service // 嵌套了与user方法的service
}

// NewDemoApi 初始化DemoApi
func NewDemoApi() *Api {
	service := NewService()
	return &Api{service: service}
}

var user string

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 string user
// @Router /demo/demo [get]
func (api *Api) Demo(c *gin.Context) {
	//appService := c.MustMake(contract.AppKey).(contract.App) // 获取app服务提供者
	//baseFolder := appService.BaseFolder() 	// 获取项目基础目录
	//users := api.service.GetUsers()
	//UsersDto := UserModelsToUserDTOs(users)
	//c.JSON(200, UsersDto)

	// 测试config
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	password := configService.GetString("database.mysql.password")

	// 测试日志
	logger := c.MustMakeLog()

	// 获取链路追踪信息
	traceService := c.MustMake(contract.TraceKey).(contract.Trace)
	traceContext := traceService.GetTrace(c)
	traceContextMap := traceService.ToMap(traceContext)

	logger.Trace(c, "demo test error", map[string]interface{}{
		"api":   "demo/demo",
		"user":  "jade",
		"trace": traceContextMap,
	})

	c.JSON(200, password + "test")
}

// Demo2 godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDto
// @Router /demo/demo2 [get]
func (api *Api) Demo2(c *gin.Context) {
	// 获取demo服务的提供者
	demoProvide := c.MustMake(demoService.DKey).(demoService.IService)
	students := demoProvide.GetAllStudent()
	usersDto := StudentsToUsersDTOs(students)
	c.JSON(200, usersDto)
}

func (api *Api) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}
	foo := &Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		_ = c.AbortWithError(500, err)
	}
	c.JSON(200, nil)
}
