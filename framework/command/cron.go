package command

import (
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/jader1992/gocore/framework/cobra"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/util"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var cronDemon = false

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDemon, "daemon", "d", false, "start serve daemon")
	cronCommand.AddCommand(cronListCommand)
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronRestartCommand)
	cronCommand.AddCommand(cronStopCommand)
	cronCommand.AddCommand(cronStateCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = c.Help()
		}
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(c *cobra.Command, args []string) error {
		cronSpecs := c.Root().CronSpecs
		ps := [][]string{}
		for _, cronSpec := range cronSpecs {
			line := []string{cronSpec.Spec, cronSpec.Cmd.Use, cronSpec.Cmd.Short}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)

		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFolder := appService.RuntimeFolder()
		serverPidFile := filepath.Join(pidFolder, "cron.pid")
		logFolder := appService.LogFolder()
		serverLogFile := filepath.Join(logFolder, "cron.log")
		currentFolder := util.GetExecDirectory()

		// deamon mode
		if cronDemon {
			// 定义守护进程
			cntxt := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0644,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				Umask:       027,
				Args:        []string{"", "cron", "start", "--daemon=true"},
			}
			// 在给定的上下文中运行当前进程的第二个副本
			d, err := cntxt.Reborn() // Reborn 方法启动一个子进程

			if err != nil {
				return err
			}

			if d != nil {
				fmt.Println("cron server started, pid:", d.Pid)
				fmt.Println("log file:", serverLogFile)
				return nil
			}

			// 释放守护进程
			defer func(cntxt *daemon.Context) {
				_ = cntxt.Release()
			}(cntxt)
			fmt.Println("deemon started")
			// 设置进程的名字
			gspt.SetProcTitle("gocore cron")
			c.Root().Cron.Run()
			return nil
		}

		// not daemon mode
		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}

		gspt.SetProcTitle("gocore cron")

		// 直接挂起的逻辑比较简单
		c.Root().Cron.Run()
		return nil
	},
}

var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			// 检查进程是否存在
			if util.CheckProcessExist(pid) {
				// 杀掉进程
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}
				// 检查进程是否被关闭
				for i := 0; i < 10; i++ {
					if util.CheckProcessExist(pid) == false {
						// 检测不存在则退出循环
						break
					}
					// 休眠1秒
					time.Sleep(1 * time.Second)
				}
				fmt.Println("fill process:" + strconv.Itoa(pid))
			}
		}

		// 重新启动
		cronDemon = true
		return cronStartCommand.RunE(c, args)
	},
}

var cronStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取进程pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")
		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}

			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}

			fmt.Println("stop pid:", pid)
		}

		return nil
	},
}

var cronStateCommand = &cobra.Command{
	Use:   "state",
	Short: "cron常驻进程状态",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		appService := container.MustMake(contract.AppKey).(contract.App)
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "cron.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			if util.CheckProcessExist(pid) {
				fmt.Println("cron server started, pid:", pid)
				return nil
			}
		}
		fmt.Println("no cron server start")
		return nil
	},
}
