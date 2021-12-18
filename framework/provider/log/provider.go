package log

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/log/formatter"
	"github.com/jader1992/gocore/framework/provider/log/services"
	"io"
	"strings"
)

// GocoreLogServiceProvider 服务提供者
type GocoreLogServiceProvider struct {
	framework.ServiceProvider

	Driver     string              // Driver
	Level      contract.LogLevel   // 日志级别
	Formatter  contract.Formatter  // 日志输出格式方法
	CtxFielder contract.CtxFielder // 日志context上下文信息获取函数
	Output     io.Writer           // 日志输出信息
}

// Register 注册一个服务实例
func (l *GocoreLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			// 默认使用console
			return services.NewGocoreConsoleLog
		}

		cs := tcs.(contract.Config)
		l.Driver = strings.ToLower(cs.GetString("log.driver"))
	}

	// 根据driver的配置确定
	switch l.Driver {
	case "singlle":
		return services.NewGocoreSingleLog
	case "rotate":
		return services.NewGocoreRotateLog
	case "console":
		return services.NewGocoreConsoleLog
	case "custome":
		return services.NewGocoreCustomLog
	default:
		return services.NewGocoreConsoleLog
	}
}

func (l *GocoreLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

func (l *GocoreLogServiceProvider) IsDefer() bool {
	return false
}

func (l *GocoreLogServiceProvider) Params(c framework.Container) []interface{} {
	// 获取configService
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	// 设置参数formatter
	if l.Formatter == nil {
		l.Formatter = formatter.TextFormatter
		// 从配置中获取日子的格式
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				l.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				l.Formatter = formatter.TextFormatter
			}
		}
	}

	if l.Level == contract.UnknownLevel {
		l.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			l.Level = logLevel(configService.GetString("log.level"))
		}
	}

	// 定义5的参数
	return []interface{}{c, l.Level, l.CtxFielder, l.Formatter, l.Output}
}

// Name 定义对应的服务字符串凭证
func (l *GocoreLogServiceProvider) Name() string {
	return contract.LogKey
}

// logLevel 从字符串中获取日志的级别，switch...case...
func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
