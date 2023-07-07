package resp

//缓存用户登录信息
type CacheUserInfoBo struct {
	OutUserId     string `json:"outUserId"`
	SourceChannel string `json:"sourceChannel"`
	UserId        int64  `json:"userId"`
	CorpId        string `json:"corpId"`
	OrgId         int64  `json:"orgId"`
}

type UserManageAuthResp struct {
	Id int64 `json:"id"`
	// 是否是组织拥有者
	IsOrgOwner bool `json:"isOrgOwner"`
	// 是否是系统管理员
	IsSysAdmin bool `json:"isSysAdmin"`
	// 是否是子管理员
	IsSubAdmin bool `json:"isSubAdmin"`
}
