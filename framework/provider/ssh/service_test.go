package ssh

import (
    "github.com/jader1992/gocore/framework/provider/config"
    "github.com/jader1992/gocore/framework/provider/log"
    tests "github.com/jader1992/gocore/test"
    . "github.com/smartystreets/goconvey/convey"
    "testing"
)


func TestGocoreSSHService_Load(t *testing.T)  {
    container := tests.InitBaseContainer()
    container.Bind(&config.GocoreConfigProvider{})
    container.Bind(&log.GocoreLogServiceProvider{})

    Convey("test get client", t, func() {
        results, err := NewGocoreSSH(container)
        So(err, ShouldBeNil)
        service, _ := results.(*GocoreSSH)
        client, err := service.GetClient(WithConfigPath("ssh.web-01"))
        So(err, ShouldBeNil)
        So(client, ShouldNotBeNil)
        session, err := client.NewSession()
        So(err, ShouldBeNil)
        out, err := session.Output("pwd")
        So(err, ShouldBeNil)
        So(out, ShouldNotBeNil)
        session.Close()
    })
}
