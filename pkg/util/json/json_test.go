package json

import (
	"testing"
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/types"
)

type TA struct {
	D time.Time
	N string
}

func TestFromJson(t *testing.T) {
	str := "{\"d\":\"2019-01-01 11:11:1000\",\"n\":\"1\"}"
	ta := &TA{}

	FromJson(str, ta)

	t.Log(ta.D)
	t.Log(ToJson(ta))
}

func TestToJson(t *testing.T) {
	ta := TA{
		D: time.Now(),
		N: "abc",
	}
	str, err := ToJson(ta)
	t.Log(err)
	t.Log(str)

	t.Log(consts.BlankTimeObject)
	t.Log(types.Time(consts.BlankTimeObject))
}
