package demo

const DEMO_KEY = "demo"

// IService 定义了student接口
type IService interface {
	GetAllStudent() []Student
}

// 定义的结构体
type Student struct {
	ID   int
	Name string
}
