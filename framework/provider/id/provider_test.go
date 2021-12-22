package id

import (
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/config"
	tests "github.com/jader1992/gocore/test"
	"github.com/smartystreets/goconvey/convey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConsolelog_Normal(t *testing.T) {
	convey.Convey("test gocore console log normal case", t, func() {
		c := tests.InitBaseContainer()
		c.Bind(&config.GocoreConfigProvider{})

		err := c.Bind(&GocoreIDProvider{})
		So(err, ShouldBeNil)

		idService := c.MustMake(contract.IDKey).(contract.IDService)
		xid := idService.NewID()
		t.Log(xid)
		So(xid, convey.ShouldNotBeEmpty)
	})
}
