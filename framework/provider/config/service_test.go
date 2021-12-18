package config

import (
	"github.com/jader1992/gocore/framework/contract"
	tests "github.com/jader1992/gocore/test"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestGocoreConfig_GetInt(t *testing.T) {
	Convey("test hade env normal case", t, func() {
		basePath := tests.BASE_PATH
		folder := filepath.Join(basePath, "config")
		serv, err := NewGocoreConfig(folder, map[string]string{}, contract.ConfigKey)
		So(err, ShouldBeNil)
		conf := serv.(*GocoreConfig)
		timeout := conf.GetInt("database.mysql.timeout")
		So(timeout, ShouldEqual, 1)
	})
}
