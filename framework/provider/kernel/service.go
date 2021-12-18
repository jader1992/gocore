package kernel

import (
	"github.com/jader1992/gocore/framework/gin"
	"net/http"
)

type GocoreKernelService struct {
	engine *gin.Engine
}

// 返回web引擎
func (s *GocoreKernelService) HttpEngine() http.Handler {
	return s.engine // s.engine是http.Handler的具体实现
}

// NewGocoreKernelService 初始化web引擎服务实例
func NewGocoreKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &GocoreKernelService{engine: httpEngine}, nil
}
