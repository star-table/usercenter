package resp

// 用户登录响应结构体
type UserLoginResp struct {
	// 用户token
	Token string `json:"token"`
	// 用户id
	UserID int64 `json:"userId"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 组织名称
	OrgName string `json:"orgName"`
	// 组织code
	OrgCode string `json:"orgCode"`
	// 用户名称
	Name string `json:"name"`
	// 头像
	Avatar string `json:"avatar"`
	// 是否需要创建组织
	NeedInitOrg bool `json:"needInitOrg"`
}
