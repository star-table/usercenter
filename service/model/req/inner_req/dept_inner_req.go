package inner_req

import "github.com/star-table/usercenter/service/model/req"

// DeptListByIdsInnerReq 组织下具体部门列表 请求数据模型
type DeptListByIdsInnerReq struct {
	// OrgId 组织ID
	OrgId int64 `json:"orgId"`
	req.IdsReq
}

// DeptListInnerReq 组织下机构列表 请求数据模型
type DeptListInnerReq struct {
	// OrgId 组织ID
	OrgId int64 `json:"orgId"`
}

//获取部门人数
type GetUserCountByDeptIdsReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//部门id
	DeptIds []int64 `json:"deptIds"`
}

//获取人员所属部门id
type GetUserDeptIdsReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//人员id
	UserId int64 `json:"userId"`
}

//获取人员所属部门id（批量）
type GetUserDeptIdsBatchReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//人员id
	UserIds []int64 `json:"userIds"`
}

type GetUserIdsByDeptIdsReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
	// 部门id
	DeptIds []int64 `json:"deptIds"`
}

type GetLeadersByDeptIdsReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
	// 部门id
	DeptIds []int64 `json:"deptIds"`
}

type GetDeptUserIdsReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
}

type GetFullDeptByIdsReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
	// 部门id
	DeptIds []int64 `json:"deptIds"`
}

type GetDeptListReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
}

type GetUserDeptIdsWithParentIdReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}
