package threadlocal

import (
	"github.com/jtolds/gls"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/model/vo"
)

var (
	Mgr = gls.NewContextManager()
)

func GetHttpContext() *vo.HttpContext {
	if traceId, ok := Mgr.GetValue(consts.HttpContextKey); ok {
		if traceId != nil {
			httpContext := traceId.(vo.HttpContext)
			return &httpContext
		}
	}
	return nil
}

func GetTraceId() string {
	if traceId, ok := Mgr.GetValue(consts.TraceIdKey); ok {
		if traceId != nil {
			return traceId.(string)
		}
	}
	return ""
}

//往threadlocal中设置UUID类型的 traceId
func SetTraceId() {

	Mgr.SetValues(gls.Values{consts.TraceIdKey: uuid.NewUuid()}, func() {})
}

func GetValue(key string) string {
	if traceId, ok := Mgr.GetValue(key); ok {
		if traceId != nil {
			return traceId.(string)
		}
	}
	return ""
}
