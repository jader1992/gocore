package command

import (
	"fmt"
	"github.com/jader1992/gocore/framework/cobra"
	"log"
	"os/exec"
)

// 我们频繁使用到了 go build、npm build 等命令进行前后端编译，在实际生产过程中，这种命令执行肯定会更加频繁
// 可以把需求整理成以下四个命令:
// 编译前端，封装 npm build 命令
// 而编译后端封装 go build 命令
// 同时编译前后端我们同时调用 npm build 和 go build 就行
// 自编译：./gocore build self

// initBuildCommand 初始化命令
func initBuildCommand() *cobra.Command {
	buildCommand.AddCommand(buildSelfCommand)     // 自编译
	buildCommand.AddCommand(buildBackendCommand)  // 编译后端
	buildCommand.AddCommand(buildFrontendCommand) // 编译前端
	buildCommand.AddCommand(buildAllCommand)      // 编译全部
	return buildCommand
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "编译相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var buildSelfCommand = &cobra.Command{
	Use:   "self",
	Short: "编译gocore命令",
	RunE: func(c *cobra.Command, args []string) error {
		// 查找go的路径
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("请先安装go")
		}

		cmd := exec.Command(path, "build", "-o", "gocore", "./")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("go build error:")
			fmt.Println(string(out))
			fmt.Println("-----------------")
			return err
		}
		fmt.Println()
		fmt.Println()
		fmt.Println("后端编译成功，请执行 ./gocore 启动")
		return nil
	},
}

var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "使用go编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(c, args)
	},
}

var buildFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "使用npm编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Fatalln("请安装npm在你的PATH路径下")
		}

		cmd := exec.Command(path, "run", "build") // npm run build
		// 将输出保存在out中
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("=============  前端编译失败 ============")
			fmt.Println(string(out))
			fmt.Println("=============  前端编译失败 ============")
			return err
		}

		// 打印输出
		fmt.Println(string(out))
		fmt.Println("=======前端编译成功=======")
		return nil
	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时编译前端和后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}
