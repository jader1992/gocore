package demo

import (
	"github.com/jader1992/gocore/framework"
)

// Service 实现了IService接口
type Service struct {
	// 参数
	container framework.Container
}

func (s *Service) GetAllStudent() []Student {
	return []Student{
		{
			ID:   1,
			Name: "foo1",
		},
		{
			ID:   2,
			Name: "bar2",
		},
	}
}

// NewService 初始化service
func NewService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &Service{container: container}, nil
}
