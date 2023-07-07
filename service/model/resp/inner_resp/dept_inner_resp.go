package inner_resp

// DeptInfoInnerResp
type DeptInfoInnerResp struct {
	Id            int64  `json:"id"`
	OrgId         int64  `json:"orgId"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	ParentId      int64  `json:"parentId"`
	Path          string `json:"path"`
	Sort          int    `json:"sort"`
	IsHide        int    `json:"isHide"`
	SourceChannel string `json:"sourceChannel"`
	Status        int    `json:"status"`
	IsDelete      int    `json:"isDelete"` // IsDelete
}

type GetUserCountByDeptIdsResp struct {
	UserCount map[int64]uint64 `json:"userCount"`
}

type GetUserDeptIdsResp struct {
	DeptIds []int64 `json:"deptIds"`
}

type GetUserDeptIdsBatchResp struct {
	Data map[int64][]int64 `json:"data"`
}

type GetUserIdsByDeptIdsResp struct {
	UserIds []int64 `json:"userIds"`
}

type DepartmentLeader struct {
	DepartmentId int64 `json:"departmentId"`
	LeaderId     int64 `json:"leaderId"`
}

type GetLeadersByDeptIdsResp struct {
	Leaders []DepartmentLeader `json:"leaders"`
}

type GetDeptUserIdsResp struct {
	Data map[int64][]int64 `json:"data"`
}

type GetFullDeptByIdsResp struct {
	Data map[int64][]string `json:"data"`
}

type GetUserDeptIdsWithParentIdResp struct {
	DeptIds []int64 `json:"deptIds"`
}
