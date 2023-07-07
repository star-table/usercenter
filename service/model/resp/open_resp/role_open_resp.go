package open_resp

// RoleInfoResp
type RoleInfoResp struct {
	Id          int64  `json:"id"`
	OrgId       int64  `json:"orgId"`
	RoleGroupId int64  `json:"roleGroupId"`
	Name        string `json:"name"`
}
