package service

import (
	"fmt"
	"testing"

	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/model/req"
)

func TestPersonalInfo(t *testing.T) {
	fmt.Println(conf.Cfg)

	t.Log(PersonalInfo(300287209087438848, 300286823215665152, ""))
}

func TestJson(t *testing.T) {
	//status := int(1)
	//a := req.CreateOrgMemberReq{
	//	PhoneNumber:   "12345679999",
	//	Email:         "2@qq.com",
	//	Name:          "hhhh",
	//	//DepartmentIds: nil,
	//	//RoleIds:       nil,
	//	Status: &status,
	//}
	str := "{\"name\": \"hhhh\", \"email\": \"2@qq.com\", \"phoneNumber\": \"12345679999\", \"status\": 1}"
	res := &req.CreateOrgMemberReq{}
	_ = json.FromJson(str, res)
	fmt.Println(json.ToJsonIgnoreError(res))
}
