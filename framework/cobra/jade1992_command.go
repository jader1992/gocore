package cobra

import (
	"github.com/jader1992/gocore/framework"
	"github.com/robfig/cron/v3"
	"log"
)

type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

// SetContainer 设置服务容器的方法是为了在创建根 Command 之后，能将服务容器设置到里面去
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取容器是为了在执行命令的 RunE 函数的时候，能从参数 Command 中获取到服务容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

// set parent
func (c *Command) SetParentNull() {
	c.parent = nil
}

func (c *Command) AddCronCommand(spec string, cmd *Command, args ...string) {
	root := c.Root()
	if root.Cron == nil {
		// 创建一个自定义选项的时间处理器
		newCronParser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		// 覆盖用于解释作业计划的解析器
		cronParse := cron.WithParser(newCronParser)
		// 返回一个新的 Cron 作业运行器，由给定的选项修改
		root.Cron = cron.New(cronParse)
		root.CronSpecs = []CronSpec{}
	}

	// 增加说明信息
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron", // 普通的定时任务
		Cmd:  cmd,
		Spec: spec,
	})

	// 制作一个rootCommand
	// 1. 封装的并不是传递进来的 Command，而是把这个 Command 做了一个副本，并且将其父节点设置为空，让它自身就是一个新的根节点；然后调用这个
	// Command 的 Execute 方法
	// 2. cobra 库在执行任意一个 Command 的 Execute 系列方法时，都会从根 Command 开始，根据参数进行遍历查询。这是因为我们是
	// 通过定时器进行调用的，这个定时器调用并没有真正的控制台，如果希望找到这个 cronCmd，直接调用其 Execute 命令就行。
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	cronCmd.SetContainer(root.GetContainer())

	// 增加调用函数
	root.Cron.AddFunc(spec, func() {
		defer func() {
			// 如果后续的command出现panic，这里要捕获，因为在 cron 中，匿名函数是开启一个 Goroutine 来执行的，而在 Golang 中，
			// 每个 Goroutine 都是平等的，任何一个 Goroutine 出现 panic，都会导致整个进程中止
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			// 打印出error信息
			log.Println(err)
		}
	})
}
