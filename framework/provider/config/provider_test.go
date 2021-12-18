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
		basePath := tests.BASE_PATH
		c := framework.NewHadeContainer()
		c.Bind(&app.GocoreAppProvider{BaseFolder: basePath})
		c.Bind(&env.GocoreEnvProvider{})
		err := c.Bind(&GocoreConfigProvider{})
		So(err, ShouldBeNil)

		conf := c.MustMake(contract.ConfigKey).(contract.Config)
		So(conf.GetString("database.mysql.hostname"), ShouldEqual, "127.0.0.1")
		So(conf.GetInt("database.mysql.timeout"), ShouldEqual, 1)
		So(conf.GetFloat64("database.mysql.readtime"), ShouldEqual, 2.3)
		//So(conf.GetString("database.mysql.password"), ShouldEqual, "mypassword") // 报错

		maps := conf.GetStringMap("database.mysql")
		So(maps, ShouldContainKey, "hostname")
		So(maps["timeout"], ShouldEqual, 1)

		maps2 := conf.GetStringMapString("database.mysql")
		So(maps2["timeout"], ShouldEqual, "1")

		type Mysql struct {
			Hostname string
			Username string
		}
		ms := &Mysql{}
		err = conf.Load("database.mysql", ms)
		Println(ms)
		So(err, ShouldBeNil)
		So(ms.Hostname, ShouldEqual, "127.0.0.1")
	})
}
