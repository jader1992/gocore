package distributed

import (
	"errors"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type LocalDistributedService struct {
	container framework.Container // 服务容器
}

func NewLocalDistrubutedService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("param error")
	}

	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container}, nil
}

// Select 为分布式选择器
func (s LocalDistributedService) Select(serviceName string, appId string, holdTime time.Duration) (selectAppId string, err error) {
	appService := s.container.MustMake(contract.AppKey).(contract.App)
	runtimeFold := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFold, "distribute_"+serviceName)

	// 打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// 尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	// 抢不到文件锁
	if err != nil {
		selectAppIdByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIdByt), err
	}

	// 在一段时间内，选举有效，其它节点在这段时间不能再进行抢占
	go func() {
		defer func() {
			// 释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 释放文件
			lock.Close()
			// 删除文件锁对应的文件
			os.Remove(lockFile)
		}()
		// 创建选举结果有效的计时器
		timer := time.NewTimer(holdTime)
		<-timer.C
	}()

	// 这里已经抢占到了，将抢占到的appId写入文件
	if _, err := lock.WriteString(appId); err != nil {
		return "", err
	}

	return appId, nil
}
