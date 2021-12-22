package console

import (
	"github.com/jader1992/gocore/app/console/command/demo"
	"github.com/jader1992/gocore/app/console/command/foo"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/cobra"
	"github.com/jader1992/gocore/framework/command"
	"time"
)

// RunCommand  初始化根Command并运行
func RunCommand(container framework.Container) error {
	// 根 Command
	var rootCmd = &cobra.Command{
		// 定义根命令的关键字
		Use: "jade1992",
		// 简短介绍
		Short: "gocore 命令",
		Long:  "gocore 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag() // 将默认帮助标志添加到 c
			return cmd.Help()
		},
		// 不需要出现cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	// 为根Command设置服务容器
	rootCmd.SetContainer(container)
	// 绑定框架的命令
	command.AddKernelCommands(rootCmd)
	// 绑定业务的命令
	AddAppCommand(rootCmd)

	// 执行RootCommand
	return rootCmd.Execute()
}

func AddAppCommand(rootCommand *cobra.Command) {
	// demo 例子
	// rootCommand.AddCommand(demo.InitFoo())

	rootCommand.AddCommand(foo.FooCommand)

	// 定时任务
	// rootCommand.AddCronCommand("* * * * * *", demo.initFoo())

	// 定时任务，调度的服务名称为init_func_for_test，每个节点每5s调用一次Foo命令，抢占到了调度任务的节点将抢占锁持续挂载2s才释放
	rootCommand.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}
