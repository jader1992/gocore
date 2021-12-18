package cobra

import (
	"github.com/jader1992/gocore/framework/contract"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func (c *Command) AddDistributedCronCommand(serviceName string, spec string, cmd *Command, holdTime time.Duration) {
	root := c.Root()

	// 初始化cron
	if root.Cron == nil {
		// 创建一个自定义选项的时间处理器
		newCronParser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		// 覆盖用于解释作业计划的解析器
		cronParse := cron.WithParser(newCronParser)
		// 返回一个新的 Cron 作业运行器，由给定的选项修改
		root.Cron = cron.New(cronParse)
		root.CronSpecs = []CronSpec{}
	}

	// cron命令的注释，这里注意Type为distributed-cron, ServiceName 需要填写
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type:        "distributed-cron",
		Cmd:         cmd,
		Spec:        spec,
		ServiceName: serviceName,
	})

	appService := root.GetContainer().MustMake(contract.AppKey).(contract.App)
	distributeService := root.GetContainer().MustMake(contract.DistributedKey).(contract.Distributed)
	appId := appService.AppId()

	// 复制要执行的command为cronCmd，并且设置为rootCmd
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()

	// cron增加匿名函数
	root.Cron.AddFunc(spec, func() {
		// 防止panic
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		// 节点进行选举，返回选举结果
		selectAppId, err := distributeService.Select(serviceName, appId, holdTime)
		if err != nil {
			return
		}

		// 如果自己没有被选择到，直接返回
		if selectAppId != appId {
			return
		}

		// 如果自己被选择到了，执行这个定时任务
		err = cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}
