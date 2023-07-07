package open_resp

type DeptInfoResp struct {
	Id                       int64  `json:"id"`                       // Id
	OrgId                    int64  `json:"orgId"`                    // OrgId
	Name                     string `json:"name"`                     // Name
	Code                     string `json:"code"`                     // Code
	ParentId                 int64  `json:"parentId"`                 // ParentId
	OutOrgDepartmentId       string `json:"outOrgDepartmentId"`       // OutOrgDepartmentId
	OutOrgDepartmentCode     string `json:"outOrgDepartmentCode"`     // OutOrgDepartmentCode
	OutOrgDepartmentParentId string `json:"outOrgDepartmentParentId"` // OutOrgDepartmentParentId
	Path                     string `json:"path"`                     // Path
	Sort                     int    `json:"sort"`                     // Sort
	IsHide                   int    `json:"isHide"`                   // IsHide
	Status                   int    `json:"status"`                   // Status
}

type DeptMemberListResp struct {
	DeptInfoResp
	UserList []SimpleMemberInfo
}
