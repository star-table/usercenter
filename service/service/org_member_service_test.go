package service

import (
	"testing"

	"github.com/star-table/usercenter/service/model/req"
)

func TestGetOrgMemberList(t *testing.T) {
	status := 1
	departmentId := int64(96750053334800)
	searchCode := ""
	listReq := req.UserListReq{
		SearchCode:   &searchCode,
		Status:       &status,
		DepartmentId: &departmentId,
	}
	orgMemberList, err := GetOrgMemberList(2797, listReq, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(orgMemberList)
}
