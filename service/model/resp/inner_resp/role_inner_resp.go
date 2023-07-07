package inner_resp

// RoleInfoInnerResp
type RoleInfoInnerResp struct {
	Id          int64  `json:"id"`
	OrgId       int64  `json:"orgId"`
	RoleGroupId int64  `json:"roleGroupId"`
	Name        string `json:"name"`
}

type GetRoleUserIdsResp struct {
	Data map[int64][]int64 `json:"data"`
}
