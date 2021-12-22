package redis

import (
	"context"
    "fmt"
    "github.com/jader1992/gocore/framework/provider/config"
	"github.com/jader1992/gocore/framework/provider/log"
	tests "github.com/jader1992/gocore/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGocoreRedis_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.GocoreConfigProvider{})
	container.Bind(&log.GocoreLogServiceProvider{})

	Convey("test get client", t, func() {
		gocoreRedis, err := NewGocoreRedis(container)
		So(err, ShouldBeNil)
		service, ok := gocoreRedis.(*GocoreRedis)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("redis.write"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		ctx := context.Background()
		err = client.Set(ctx, "foo", "bar", 1*time.Hour).Err()
		So(err, ShouldBeNil)
		val, err := client.Get(ctx, "foo").Result()
		So(err, ShouldBeNil)
		So(val, ShouldEqual, "bar")
		err = client.Del(ctx, "foo").Err()
		So(err, ShouldBeNil)
        err = client.IncrBy(ctx, "nums", 1).Err()
        So(err, ShouldBeNil)
        fmt.Println(client.Get(ctx, "nums"))
	})
}
