package demo

// Repository 提供数据，他会与model交互
type Repository struct {

}

// GetUserIds 获取用户id
func (r *Repository) GetUserIds() []int {
	return []int{1, 2}
}

// 获取指定的用户
func (r *Repository) GetUserByIds([]int) []UserModel {
	return []UserModel{
		{
			UserId: 1,
			Name: "foo",
			Age: 11,
		},
		{
			UserId: 2,
			Name: "bar",
			Age: 12,
		},
	}
}

func NewRepository() *Repository {
	return &Repository{}
}
