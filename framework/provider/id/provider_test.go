package id

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/app"
	"github.com/jader1992/gocore/framework/provider/config"
	"github.com/jader1992/gocore/framework/provider/env"
	"github.com/jader1992/gocore/framework/util"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConsolelog_Normal(t *testing.T)  {
	convey.Convey("test gocore console log normal case", t, func() {
		basePath := util.GetExecDirectory()
		c := framework.NewHadeContainer()
		c.Bind(&app.GocoreAppProvider{BaseFolder: basePath})
		c.Bind(&env.GocoreEnvProvider{})
		c.Bind(&config.GocoreConfigProvider{})
		c.Bind(&GocoreIDProvider{})

		idService := c.MustMake(contract.IDKey).(contract.IDService)
		xid := idService.NewID()
		t.Log(xid)
		convey.So(xid, convey.ShouldBeEmpty)
	})
}
