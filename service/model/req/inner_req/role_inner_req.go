package inner_req

// RoleListByIdsInnerReq
type RoleListByIdsInnerReq struct {
	// OrgId 组织ID
	OrgId int64 `json:"orgId"`
	// Ids ID列表
	Ids []int64 `json:"ids"`
}

type GetRoleUserIdsReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
}
