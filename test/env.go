package test

import (
    "github.com/jader1992/gocore/framework"
    "github.com/jader1992/gocore/framework/provider/app"
    "github.com/jader1992/gocore/framework/provider/env"
)

const BasePath = "/Users/jade/Desktop/go/coredemo/"

func InitBaseContainer() framework.Container  {
    container := framework.NewGocoreContainer()
    container.Bind(&app.GocoreAppProvider{
        BaseFolder: BasePath,
    })
    container.Bind(&env.GocoreEnvProvider{})
    return container
}
