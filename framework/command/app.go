package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/cobra"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/util"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var appAddress = ""
var appDaemon = false

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVar(&appAddress, "address", ":8888", "设置app启动的地址，默认为:8888")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令, 其包含业务启动，关闭，重启，查询等功能",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 打印帮助文档
		_ = cmd.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 从Command中获取服务容器
		container := cmd.GetContainer()
		// 从服务容器中获取kernel的服务示例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		// 从kernel服务实例中获取引擎
		core := kernelService.HttpEngine()

		configService := container.MustMake(contract.ConfigKey).(contract.Config)

		// 获取appAddress
		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("ADDRESS") != "" {
				appAddress = envService.Get("ADDRESS")
			} else {
				if configService.IsExist("app.address") {
					appAddress = configService.GetString("app.address")
				} else {
					appAddress = ":8888"
				}
			}
		}

		// 创建一个Serve服务
		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)

		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}

		serverPidFile := filepath.Join(pidFolder, "app.pid")
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return err
			}
		}

		// 应用日志
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()

		// daemon模式
		if appDaemon {
			// 创建一个Context
			baseContext := &daemon.Context{
				PidFileName: serverPidFile, // 设置pid文件
				PidFilePerm: 0644,
				LogFileName: serverLogFile, // 设置日志文件
				LogFilePerm: 0640,
				WorkDir:     currentFolder, // 设置工作路径
				Umask:       027,           // 设置所有设置文件的mask，默认为750
				// 子进程的参数，按照这个参数设置，子进程的命令为 ./gocore app start --daemon=true
				Args: []string{"", "app", "start", "--daemon=true"},
			}

			// 启动子进程，d不为空表示当前是父进程，d为空表示是子进程
			d, err := baseContext.Reborn()
			// 把 Reborn 理解成 fork，当调用这个函数的时候，父进程会继续往下走，但是返回值 d 不为空，它的信息是子进程的进程号等信息。
			if err != nil {
				return err
			}

			// 再次进入到 Reborn 函数的时候，返回的 d 就为 nil。所以在 Reborn 的后面，我们让父进程直接 return，
			// 而让子进程继续往后进行操作，这样就达到了 fork 一个子进程的效果了
			if d != nil {
				fmt.Println("app启动成功，pid:", d.Pid)
				fmt.Println("日志文件：", serverLogFile)
				return nil
			}
			defer func(baseContext *daemon.Context) {
				_ = baseContext.Release()
			}(baseContext)

			// 子进程执行真正的app启动操作
			fmt.Println("daemon started")
			gspt.SetProcTitle(strings.TrimPrefix(configService.GetString("app.dev.backend.cmd"), "./") + " app ")

			fmt.Println("app serve url:", appAddress)
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		{
			// 非daemon模式，直接执行
			content := strconv.Itoa(os.Getpid())
			fmt.Println("[PID]", content)
			err := ioutil.WriteFile(serverPidFile, []byte(content), 0644)
			if err != nil {
				return err
			}

			gspt.SetProcTitle(strings.TrimPrefix(configService.GetString("app.dev.backend.cmd"), "./") + " app ")
			fmt.Println("app serve url", appAddress)
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}
	},
}

// 重新启动一个app服务
var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重新启动一个app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		// 如果app.pid 不存在直接以守护进程方式启动apps
		if !util.Exists(serverPidFile) {
			appDaemon = true
			return appStartCommand.RunE(c, args)
		}

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		configService := container.MustMake(contract.ConfigKey).(contract.Config)

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			//由于新、旧进程都是使用同一个端口，所以必须保证旧进程结束，才能启动新的进程。
			if util.CheckProcessExist(pid) {
				// 杀死进程
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}

				closeWait := 5
				if configService.IsExist("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}

				// 确认进程已经关闭，每秒检查1次，最多检测closeWait * 2 秒
				for i := 0; i < closeWait*2; i++ {
					if util.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}

				if util.CheckProcessExist(pid) == true {
					fmt.Println("结束进程失败："+strconv.Itoa(pid), "请查明原因")
					return errors.New("结束进程失败")
				}

				if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}

				fmt.Println("结束进程成功:" + strconv.Itoa(pid))
			}
		}

		appDaemon = true
		return appStartCommand.RunE(c, args)
	},
}

var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := ioutil.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			// 发送SIGTERM命令
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := ioutil.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("停止进程:", pid)
		}
		return nil
	},
}

// 获取启动的app的pid
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		if !util.Exists(serverPidFile) {
			fmt.Println("没有app服务存在")
			return nil
		}
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
				fmt.Println("app服务已经启动, pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}

// startAppServe 启动AppServer, 这个函数会将当前goroutine阻塞
func startAppServe(server *http.Server, c framework.Container) error {
	// 这个goroutine是启动服务的goroutine
	go func() {
		_ = server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量 【主Goroutine监听信号】
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Serve.shutdown graceful结束
	closeWait := 5
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	return nil
}
