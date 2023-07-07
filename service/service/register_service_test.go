package service

import (
	"log"
	"testing"

	"github.com/star-table/usercenter/service/model/req"
)

func TestCreateUser(t *testing.T) {
	status := 1
	resp, err := CreateOrgMember(1025, 1004, req.CreateOrgMemberReq{
		PhoneNumber:      "15079051001",
		Email:            "15079051001@qq.com",
		Name:             "苏1001",
		DeptAndPositions: nil,
		RoleIds:          []int64{2173},
		Status:           status,
	})
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	log.Printf("%v\n", resp)
}
