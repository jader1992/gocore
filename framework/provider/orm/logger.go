package orm

import (
    "context"
    "github.com/jader1992/gocore/framework/contract"
    "gorm.io/gorm/logger"
    "time"
)

// OrmLogger orm的日志实现类, 实现了gorm.Logger.Interface
type OrmLogger struct {
    logger contract.Log // 有一个logger对象存放log服务
}

func (o *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
    return o
}

// Info 对接gocore的Info输出
func (o *OrmLogger) Info(ctx context.Context, s string, i ...interface{})  {
    fields := map[string]interface{}{
        "fields": i,
    }
    o.logger.Info(ctx, s, fields)
}

// Warn 对接gocore的Warn输出
func (o *OrmLogger) Warn(ctx context.Context, s string, i ...interface{})  {
    fields := map[string]interface{}{
        "fields": i,
    }
    o.logger.Warn(ctx, s, fields)
}

// Error 对接gocore的Error输出
func (o *OrmLogger) Error(ctx context.Context, s string, i ...interface{})  {
    fields := map[string]interface{}{
        "fields": i,
    }
    o.logger.Error(ctx, s, fields)
}

// Trace 对接gocore的Trace输出
func (o *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)  {
    sql, rows := fc()
    usedTime := time.Since(begin)
    fields := map[string]interface{}{
        "begin": begin,
        "error": err,
        "sql": sql,
        "rows": rows,
        "time": usedTime,
    }
    s := "orm trace sql"
    o.logger.Trace(ctx, s, fields)
}


// NewOrmLogger 初始化一个Logger
func NewOrmLogger(logger contract.Log) *OrmLogger {
    return &OrmLogger{logger: logger}
}






