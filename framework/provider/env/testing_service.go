package env

// GocoreTestingEnv 是 Env 的具体实现
type GocoreTestingEnv struct {
	folder string            // 代表.env所在的目录
	maps   map[string]string // 保存所有的环境变量
}

func NewGocoreTestingEnv(params ...interface{}) (interface{}, error) {
	return &GocoreTestingEnv{}, nil
}

// AppEnv 获取表示当前APP环境的变量APP_ENV
func (en *GocoreTestingEnv) AppEnv() string {
	return "testing"
}

// Get 获取某个环境变量，如果没有设置，返回""
func (en *GocoreTestingEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

// IsExist 判断一个环境变量是否有被设置
func (en *GocoreTestingEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// All 获取所有的环境变量，.env和运行环境变量融合后
func (en *GocoreTestingEnv) All() map[string]string {
	return en.maps
}
