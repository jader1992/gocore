package demo

const DKey = "demo"

// IService 定义了student接口
type IService interface {
	GetAllStudent() []Student
}

// Student 定义的结构体
type Student struct {
	ID   int
	Name string
}
