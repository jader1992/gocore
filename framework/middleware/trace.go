package middleware

import (
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/gin"
)

func Trace() gin.HandlerFunc  {
	// 使用函数回调
	return func(c *gin.Context) {
		tracer := c.MustMake(contract.TraceKey).(contract.Trace)
		traceCtx := tracer.ExtractHTTP(c.Request)

		tracer.WithTrace(c, traceCtx)
		// 使用next执行具体的业务逻辑
		c.Next()
	}
}
