package ssh

import (
    "github.com/jader1992/gocore/framework"
    "github.com/jader1992/gocore/framework/contract"
)

// GocoreSSHProvider 提供具体实现方法
type GocoreSSHProvider struct {

}

// Register 注册方法
func (h *GocoreSSHProvider) Register(container framework.Container) framework.NewInstance {
    return NewGocoreSSH
}

// Boot 启动调用
func (h *GocoreSSHProvider) Boot(container framework.Container) error {
    return nil
}

// IsDefer 是否延迟初始化
func (h *GocoreSSHProvider) IsDefer() bool {
    return true
}

// Params 获取初始化参数
func (h *GocoreSSHProvider) Params(container framework.Container) []interface{} {
    return []interface{}{container}
}

// Name 获取字符串凭证
func (h *GocoreSSHProvider) Name() string {
    return contract.SSHKey
}

