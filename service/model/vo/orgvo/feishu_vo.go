package orgvo

import "github.com/star-table/usercenter/service/model/vo"

type FeiShuAuthRespVo struct {
	vo.Err
	Auth *vo.FeiShuAuthResp `json:"data"`
}

type GetFsAccessTokenReqVo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}

type GetFsAccessTokenRespVo struct {
	AccessToken string `json:"accessToken"`
	vo.Err
}
