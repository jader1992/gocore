package contract

import (
	"context"
	"net/http"
)

const TraceKey = "gocore:trace"

const (
	TraceKeyTraceID  = "trace_id"
	TraceKeySpanID   = "span_id"
	TraceKeyCspanID  = "cspan_id"
	TraceKeyParentID = "parent_id"
	TraceKeyMethod   = "method"
	TraceKeyCaller   = "caller"
	TraceKeyTime     = "time"
)

// TraceContext 根据 Google Dapper 定义
type TraceContext struct {
	TraceID  string // 唯一ID
	ParantID string // 父节点spanID
	SpanID   string // 当前节点 SpanID
	CspanID  string // 子节点调用的SpanID，由调用方指定

	Annotation map[string]string // 标记各种信息
}

type Trace interface {
	WithTrace(c context.Context, trace *TraceContext) context.Context // 注册新的跟踪路由
	GetTrace(c context.Context) *TraceContext                         // 从当前上下文中获取traceContext
	NewTrace() *TraceContext                                          // 生成一个新的traceContext

	StartSpan(trace *TraceContext) *TraceContext // 为调用者生成traceContext
	ToMap(trace *TraceContext) map[string]string // 将traceContext 转成map

	ExtractHTTP(req *http.Request) *TraceContext                     // 从http.Request中获取traceContext
	InjectHTTP(req *http.Request, trace *TraceContext) *http.Request // 为http设置traceContext
}
