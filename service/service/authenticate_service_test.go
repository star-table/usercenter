package service

import (
	"testing"

	"github.com/star-table/usercenter/service/model/vo"
)

func TestUserPwdLogin(t *testing.T) {
	res, err := userPwdLogin(vo.UserLoginReq{
		LoginType: 2,
		LoginName: "fuse003",
		Password:  "zZ123456",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
