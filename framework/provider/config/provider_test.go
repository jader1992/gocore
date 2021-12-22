package config

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/app"
	"github.com/jader1992/gocore/framework/provider/env"
	tests "github.com/jader1992/gocore/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGocoreConfig_Normal(t *testing.T) {
	Convey("test hade config normal case", t, func() {
		basePath := tests.BasePath
		c := framework.NewGocoreContainer()
		c.Bind(&app.GocoreAppProvider{BaseFolder: basePath})
		c.Bind(&env.GocoreEnvProvider{})
		err := c.Bind(&GocoreConfigProvider{})
		So(err, ShouldBeNil)

		conf := c.MustMake(contract.ConfigKey).(contract.Config)
		So(conf.GetString("database.default.host"), ShouldEqual, "localhost")
		So(conf.GetInt("database.default.port"), ShouldEqual, 3306)
		//So(conf.GetFloat64("database.default.readtime"), ShouldEqual, 2.3)
		//So(conf.GetString("database.mysql.password"), ShouldEqual, "mypassword") // 报错

		maps := conf.GetStringMap("database.default")
		So(maps, ShouldContainKey, "host")
		So(maps["host"], ShouldEqual, "localhost")

		maps2 := conf.GetStringMapString("database.default")
		So(maps2["host"], ShouldEqual, "localhost")

		type Mysql struct {
			Host string `yaml:"host"`
		}
		ms := &Mysql{}
		err = conf.Load("database.default", ms)
		Println(ms)
		So(err, ShouldBeNil)
		So(ms.Host, ShouldEqual, "localhost")
	})
}
