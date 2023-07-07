package service

import (
	"testing"

	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/vo"
)

// 发送邮箱验证码
func TestSendAuthCode(t *testing.T) {
	var (
		captchaID       = "1"
		captchaPassword = "212091"
	)
	ok, err := SendAuthCode(vo.SendAuthCodeReq{
		AuthType:        6,
		AddressType:     2,
		Address:         "suhanyu@bjx.cloud",
		CaptchaID:       &captchaID,
		CaptchaPassword: &captchaPassword,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v, end", ok)
}

// 测试解绑登录名的逻辑
func TestUnbindLoginName(t *testing.T) {
	// todo config init
	var (
		orgId  = int64(315570826537603072)
		userId = int64(315570713178148864)
		//projectID = int64(0)
		//issueID = int64(0)
		//folderID = int64(0)
	)
	var resp, err = UnbindLoginName(orgId, userId, req.UnbindLoginNameReq{
		AddressType: 2,
		AuthCode:    "123",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

func TestBindLoginName(t *testing.T) {
	var (
		orgId  = int64(315570826537603072)
		userId = int64(315570713178148864)
		//projectID = int64(0)
		//issueID = int64(0)
		//folderID = int64(0)
	)
	err := BindLoginName(orgId, userId, req.BindLoginNameReq{
		Address:        "suhanyujie@qq.com",
		AddressType:    2,
		AuthCode:       "123",
		ChangeBindCode: "202cb962ac59075b964b07152d234b701--",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("end")
}
