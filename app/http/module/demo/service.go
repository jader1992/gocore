package demo

// Service 提供User方法
type Service struct {
	repository *Repository // 嵌套了与数据交互的repository
}

func (s *Service) GetUsers() []UserModel {
	// 获取数据
	ids := s.repository.GetUserIds()
	return s.repository.GetUserByIds(ids)
}

func NewService() *Service {
	repository := NewRepository()
	return &Service{
		repository: repository,
	}
}


