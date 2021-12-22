package ssh

import (
    "context"
    "github.com/jader1992/gocore/framework"
    "github.com/jader1992/gocore/framework/contract"
    "golang.org/x/crypto/ssh"
    "sync"
)

type GocoreSSH struct {
    container framework.Container // 服务容器
    clients map[string]*ssh.Client // key为uniqKey, value 为 ssh.Client (连接池)

    lock *sync.RWMutex
}

// GetClient 获取Client实例
func (app *GocoreSSH) GetClient(option ...contract.SSHOption) (*ssh.Client, error) {

    config := GetBaseConfig(app.container)

    for _, opt := range option {
        if err := opt(app.container, config); err != nil {
            return nil, err
        }
    }

    key := config.UniqKey()

    // 判断是否已经实例化
    app.lock.RLock()
    if db, ok := app.clients[key]; ok {
        app.lock.RUnlock()
        return db, nil
    }
    app.lock.RUnlock()

    // 没有实例化
    app.lock.Lock()
    defer app.lock.Unlock()

    addr := config.Host + ":" + config.Port
    client, err := ssh.Dial(config.NetWork, addr, config.ClientConfig)
    if err != nil {
        logService := app.container.MustMake(contract.LogKey).(contract.Log)
        logService.Error(context.Background(), "ssh dial error", map[string]interface{}{
            "err":  err,
            "addr": addr,
        })
    }

    // 挂载到map中, 结束配置
    app.clients[key] = client

    return client, nil
}

// NewGocoreSSH 代表实例化Client
func NewGocoreSSH(params ...interface{}) (interface{}, error) {
    container := params[0].(framework.Container)
    clients := make(map[string]*ssh.Client)
    lock := &sync.RWMutex{}
    return &GocoreSSH{
        container: container,
        clients: clients,
        lock: lock,
    }, nil
}
