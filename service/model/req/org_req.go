package req

// 创建组织请求结构体
type CreateOrgReq struct {
	// 组织名称
	OrgName string `json:"orgName"`
	// 补全个人姓名
	CreatorName *string `json:"creatorName"`
	// 是否要导入示例数据, 1：导入，2：不导入，默认不导入
	ImportSampleData *int `json:"importSampleData"`
	// 来源平台
	SourcePlatform *string `json:"sourcePlatform"`
	// 来源渠道
	SourceChannel *string `json:"sourceChannel"`
	// 所属行业id(选填)
	IndustryId *int64 `json:"industryId"`
	//组织规模（选填）
	Scale *string `json:"scale"`
}

// 转为内部成员请求结构体
type SwitchToInnerMemberReq struct {
	UserIds []int64 `json:"userIds"`
}
