package domain

import "testing"

func TestGetDeptParentIds(t *testing.T)  {
	resp, err := GetDeptParentIds(2797, []int64{97057403958544})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
