package demo

import (
    "database/sql"
    "github.com/jader1992/gocore/framework/contract"
    "github.com/jader1992/gocore/framework/gin"
    "github.com/jader1992/gocore/framework/provider/orm"
    "time"
)

func (api *Api) DemoOrm(c *gin.Context)  {
    logger := c.MustMakeLog()
    logger.Info(c, "request start", nil)

    // 初始化一个orm.DB
    gormService := c.MustMake(contract.ORMKey).(contract.IORMService)
    db, err := gormService.GetDB(orm.WithConfigPath("database.default"))

    if err != nil {
        logger.Error(c, err.Error(), nil)
        c.AbortWithError(500001, err)
    }
    db.WithContext(c)

    // 将User模型创建到数据库中
    err = db.AutoMigrate(&User{})
    if err != nil {
        c.AbortWithError(500, err)
        return
    }

    logger.Info(c, "migrate ok", nil)

    // 插入一条数据
    email := "foo@gmail.com"
    name := "foo"
    age := uint8(25)
    birthday := time.Date(2001, 1, 1, 1, 1, 1, 1, time.Local)

    user := &User{
        Name: name,
        Email: &email,
        Age: age,
        Birthday: &birthday,
        MemberNumber: sql.NullString{},
        ActivatedAt: sql.NullTime{},
        CreatedAt: time.Now(),
        UpdateAt: time.Now(),
    }

    err = db.Create(user).Error
    logger.Info(c, "insert user", map[string]interface{}{
        "id": user.ID,
        "err": err,
    })

    // 更新一条数据
    user.Name = "bar"
    err = db.Save(user).Error
    logger.Info(c, "update user", map[string]interface{}{
        "id": user.ID,
        "err": err,
    })

    // 查询一条数据

    queryRaw := &User{ID: user.ID}

    err = db.First(queryRaw).Error
    logger.Info(c, "query user", map[string]interface{}{
        "name": queryRaw.Name,
        "err": err,
    })

    c.JSON(200, "ok")
}
