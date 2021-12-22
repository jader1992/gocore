package orm

import (
    "context"
    "github.com/jader1992/gocore/framework"
    "github.com/jader1992/gocore/framework/contract"
    "gorm.io/driver/clickhouse"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
    "sync"
    "time"
)

type CoreGorm struct {
	container framework.Container
	dbs        map[string]*gorm.DB // mysql连接池

	lock *sync.RWMutex
}

func (app *CoreGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
    logger := app.container.MustMake(contract.LogKey).(contract.Log)

    dbConfig := GetBaseConfig(app.container)

    // 设置Logger
    OrmLogger := NewOrmLogger(logger)
    dbConfig.Config = &gorm.Config{
        Logger: OrmLogger,
    }

    // option对opt进行修改
    for _, opt := range option { // 这个设计比价好
        if err := opt(app.container, dbConfig); err != nil {
            return nil, err
        }
    }

    // 如果最终的dbLog没有设置dsn,就生成dsn
    if dbConfig.Dsn == "" {
        dsn, err := dbConfig.FormatDSN()
        if err != nil {
            return nil, err
        }
        dbConfig.Dsn = dsn
    }

    // 判断是否已经实例化了gorm.DB
    app.lock.RLock()
    if db, ok := app.dbs[dbConfig.Dsn]; ok {
        app.lock.RUnlock()
        return db, nil
    }
    app.lock.RUnlock()

    // 没有实例化gorm.DB,那么就要进行实例化操作
    app.lock.Lock()
    defer app.lock.Unlock()

    var db *gorm.DB
    var err error

    switch dbConfig.Driver {
    case "mysql":
        db, err = gorm.Open(mysql.Open(dbConfig.Dsn), dbConfig)
    case "postgres":
        db, err = gorm.Open(postgres.Open(dbConfig.Dsn), dbConfig)
    case "sqlite":
        db, err = gorm.Open(sqlite.Open(dbConfig.Dsn), dbConfig)
    case "sqlserver":
        db, err = gorm.Open(sqlserver.Open(dbConfig.Dsn), dbConfig)
    case "clickhouse":
        db, err = gorm.Open(clickhouse.Open(dbConfig.Dsn), dbConfig)
    }

    // 设置对应的连接池配置
    sqlDB, err := db.DB() // 获取的*db.ConnPool.(*sql.DB)
    if err != nil {
        return db, err
    }

    if dbConfig.ConnMaxIdle > 0 {
        sqlDB.SetMaxIdleConns(dbConfig.ConnMaxIdle)
    }

    if dbConfig.ConnMaxOpen > 0 {
        sqlDB.SetMaxOpenConns(dbConfig.ConnMaxOpen)
    }

    if dbConfig.ConnMaxLifetime != "" {
        lifeTime, err := time.ParseDuration(dbConfig.ConnMaxLifetime)
        if err != nil {
            logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
                "err": err,
            })
        } else {
            sqlDB.SetConnMaxLifetime(lifeTime)
        }
    }

    if dbConfig.ConnMaxIdletime != "" {
        idleTime, err := time.ParseDuration(dbConfig.ConnMaxIdletime)
        if err != nil {
            logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
                "err": err,
            })
        } else {
            sqlDB.SetConnMaxIdleTime(idleTime)
        }
    }

    // 挂载到map中，结束配置
    if err != nil {
        app.dbs[dbConfig.Dsn] = db
    }

    return db, nil
}

// NewGorm 初始化gorm容器
func NewGorm(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
    dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}

	return &CoreGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil
}



