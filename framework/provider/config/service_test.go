package config

import (
	"github.com/jader1992/gocore/framework/contract"
	tests "github.com/jader1992/gocore/test"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestGocoreConfig_GetInt(t *testing.T) {
    container := tests.InitBaseContainer()
	Convey("test hade env normal case", t, func() {
		appService := container.MustMake(contract.AppKey).(contract.App)
        envService := container.MustMake(contract.EnvKey).(contract.Env)
		folder := filepath.Join(appService.ConfigFolder(), envService.AppEnv())

		serv, err := NewGocoreConfig(container, folder, map[string]string{})

		So(err, ShouldBeNil)
		conf := serv.(*GocoreConfig)
		timeout := conf.GetString("database.default.timeout")
		So(timeout, ShouldEqual, "10s")
	})
}
