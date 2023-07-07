package inner_req

// UserListByIdsInnerReq
type UserListByIdsInnerReq struct {
	// OrgId 组织ID
	OrgId int64 `json:"orgId"`
	// Ids ID列表
	Ids []int64 `json:"ids"`
}

// UserAuthorityInnerReq
type UserAuthorityInnerReq struct {
	// OrgId 组织ID
	OrgId int64 `json:"orgId"`
	// 用户id
	UserId int64 `json:"userId"`
}

type GetMemberSimpleInfoReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
	// 类型(1成员2部门3角色)
	Type int `json:"type"`
	// 是否需要已经删除的数据
	NeedDelete int `json:"needDelete"`
}

type SimpleReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
}

type GetManageUserReq struct {
	OrgId int64 `json:"orgId"`
	AppId int64 `json:"appId"`
}

type UserListReq struct {
	OrgId int64 `json:"orgId"`
	// 当前用户id，操作人id
	CurUserId int64 `json:"curUserId"`

	//搜索字段
	SearchCode *string `json:"searchCode"`
	//是否已分配部门（1已分配2未分配，默认全部）
	IsAllocate *int `json:"isAllocate"`
	//是否禁用（1启用2禁用，默认全部）
	Status *int `json:"status"`
	// 审核状态。1待审核；2审核通过；3审核不通过
	CheckStatus []int `json:"checkStatus"`
	//角色id
	RoleId *int64 `json:"roleId"`
	//部门id
	DepartmentId *int64 `json:"departmentId"`
	//职级ID
	PositionId *int64 `json:"positionId"`
	// 页码（选填，不填为全部）
	Page int `json:"page"`
	// 数量（选填，不填为全部）
	Size int `json:"size"`
	// 是否开启权限过滤
	AuthFilter *bool `json:"authFilter"`
	// 排序字段 如"create_time desc,id asc"
	Order *string `json:"order"` // Order
}

type MemberSimpleInfoListReq struct {
	OrgId int64 `json:"orgId"`
	//是否已分配部门（1已分配2未分配，默认全部）
	IsAllocate *int `json:"isAllocate"`
	//是否禁用（1启用2禁用，默认全部）
	Status *int `json:"status"`
	// 审核状态。1待审核；2审核通过；3审核不通过
	CheckStatus []int `json:"checkStatus"`
	// 页码（选填，不填为全部）
	Page int `json:"page"`
	// 数量（选填，不填为全部）
	Size int `json:"size"`
	// 排序字段 如"create_time desc,id asc"
	Order   string   `json:"order"`   // Order
	UserIds []string `json:"userIds"` // userIds
}

type GetCommAdminMangeAppsReq struct {
	OrgId int64 `json:"orgId"`
}
