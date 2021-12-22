package trace

import (
	"context"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/gin"
	"net/http"
	"time"
)

type TraceKey string

var ContextKey = TraceKey("trace-key")

type GocoreTraceService struct {
	idService contract.IDService

	traceIDGenerator contract.IDService
	spanIDGenerator  contract.IDService
}

func (t *GocoreTraceService) GetTrace(c context.Context) *contract.TraceContext {
	if ginc, ok := c.(*gin.Context); ok {
		if val, ok2 := ginc.Get(string(ContextKey)); ok2 {
			return val.(*contract.TraceContext)
		}
	}

	if tc, ok := c.Value(ContextKey).(*contract.TraceContext); ok {
		return tc
	}
	return nil
}

func (t *GocoreTraceService) WithTrace(c context.Context, trace *contract.TraceContext) context.Context {
	if ginC, ok := c.(*gin.Context); ok {
		ginC.Set(string(ContextKey), trace)
		return ginC
	} else {
		newC := context.WithValue(c, ContextKey, trace)
		return newC
	}
}

// NewTrace 生成新的跟踪
func (t *GocoreTraceService) NewTrace() *contract.TraceContext {
	var traceID, spanID string
	if t.traceIDGenerator != nil {
		traceID = t.traceIDGenerator.NewID()
	} else {
		traceID = t.idService.NewID()
	}

	if t.spanIDGenerator != nil {
		spanID = t.spanIDGenerator.NewID()
	} else {
		spanID = t.idService.NewID()
	}

	tc := &contract.TraceContext{
		TraceID:    traceID,
		ParentID:   "",
		SpanID:     spanID,
		CSpanID:    "",
		Annotation: map[string]string{},
	}
	return tc
}

func (t *GocoreTraceService) StartSpan(tc *contract.TraceContext) *contract.TraceContext {
	var childSpanId string
	if t.spanIDGenerator != nil {
		childSpanId = t.spanIDGenerator.NewID()
	} else {
		childSpanId = t.idService.NewID()
	}

	childSpan := &contract.TraceContext{
		TraceID:  tc.TraceID,
		ParentID: "",
		SpanID:   tc.SpanID,
		CSpanID:  childSpanId,
		Annotation: map[string]string{
			contract.TraceKeyTime: time.Now().String(),
		},
	}
	return childSpan
}

func (t *GocoreTraceService) ExtractHTTP(req *http.Request) *contract.TraceContext {
	tc := &contract.TraceContext{}
	tc.TraceID = req.Header.Get(contract.TraceKeyTraceID)
	tc.ParentID = req.Header.Get(contract.TraceKeyParentID)
	tc.SpanID = req.Header.Get(contract.TraceKeySpanID)
	tc.CSpanID = ""

	if tc.TraceID == "" {
		tc.TraceID = t.idService.NewID()
	}

	if tc.SpanID == "" {
		tc.SpanID = t.idService.NewID()
	}

	return tc
}

func (t *GocoreTraceService) InjectHTTP(req *http.Request, tc *contract.TraceContext) *http.Request {
	req.Header.Set(contract.TraceKeyTraceID, tc.TraceID)
	req.Header.Set(contract.TraceKeySpanID, tc.SpanID)
	req.Header.Set(contract.TraceKeyCSpanID, tc.CSpanID)
	req.Header.Set(contract.TraceKeyParentID, tc.ParentID)
	return req
}

func (t *GocoreTraceService) ToMap(tc *contract.TraceContext) map[string]string {
	m := map[string]string{}
	if tc == nil {
		return m
	}

	m[contract.TraceKeyTraceID] = tc.TraceID
	m[contract.TraceKeySpanID] = tc.SpanID
	m[contract.TraceKeyCSpanID] = tc.CSpanID
	m[contract.TraceKeyParentID] = tc.ParentID

	if tc.Annotation != nil {
		for k, v := range tc.Annotation {
			m[k] = v
		}
	}
	return m
}

func (t *GocoreTraceService) SetTraceIDService(service contract.IDService)  {
	t.traceIDGenerator = service
}

func (t *GocoreTraceService) SetSpanIDService(service contract.IDService)  {
	t.spanIDGenerator = service
}

func NewGocoreTraceService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	idService := c.MustMake(contract.IDKey).(contract.IDService)
	return &GocoreTraceService{idService: idService}, nil
}
