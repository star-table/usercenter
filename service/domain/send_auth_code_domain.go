package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/slice"
)

func CheckAuthTypeValid(authType int) bool {
	if ok, _ := slice.Contain([]int{
		consts.AuthCodeTypeBind, consts.AuthCodeTypeUnBind, consts.AuthCodeTypeRegister,
		consts.AuthCodeTypeRetrievePwd, consts.AuthCodeTypeChangeSuperAdmin,
	}, authType); ok {
		return true
	}
	return false
}
