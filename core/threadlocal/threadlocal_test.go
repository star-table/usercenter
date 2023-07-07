package threadlocal

import (
	"testing"

	"github.com/jtolds/gls"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/uuid"
)

func TestSetTraceId(t *testing.T) {

	SetTraceId()
	t.Log(GetTraceId())

	Mgr.SetValues(gls.Values{consts.TraceIdKey: uuid.NewUuid()}, func() {

		t.Log("in ", GetTraceId())
	})
}
