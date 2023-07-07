package orgvo

import (
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/vo"
)

type CacheUserInfoVo struct {
	vo.Err

	CacheInfo bo.CacheUserInfoBo `json:"data"`
}
