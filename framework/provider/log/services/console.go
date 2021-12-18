package services

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"os"
)

// GocoreConsoleLog 代表控制台数输出
type GocoreConsoleLog struct {
	GocoreLog
}

// NewGocoreConsoleLog 实例化GocoreConsoleLog
func NewGocoreConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &GocoreConsoleLog{}

	// 组装数据
	log.SetLevel(level)
	log.SetFormatter(formatter)
	log.SetCtxFielder(ctxFielder)
	// 最重要的内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c

	return log, nil
}
