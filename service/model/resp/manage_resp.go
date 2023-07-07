package resp

import "time"

// SimpleManageGroupInfo 简单管理组信息
type SimpleManageGroupInfo struct {
	Id       int64  `json:"id"`
	OrgId    int64  `json:"orgId"`
	Name     string `json:"name"`
	LangCode string `json:"langCode"`
	OptAuth  []string `json:"optAuth"`
}

// ManageGroupTreeResp 管理组列表信息
type ManageGroupTreeResp struct {
	// 系统管理组
	SysGroup *SimpleManageGroupInfo `json:"sysGroup"`
	// 普通管理组列表
	GeneralGroups []SimpleManageGroupInfo `json:"generalGroups"`
}

// ManageGroupInfo 管理组信息
type ManageGroupInfo struct {
	Id            int64     `json:"id"`
	OrgId         int64     `json:"orgId"`
	Name          string    `json:"name"`
	Type          int       `json:"type"`          // 1 系统管理组 2普通管理组
	LangCode      string    `json:"langCode"`
	UserIds       []int64   `json:"userIds"`       // 成员
	AppPackageIds []int64   `json:"appPackageIds"` // 应用包权限
	AppIds        []int64   `json:"appIds"`        // 应用权限
	UsageIds      []int64   `json:"usageIds"`      // 功能权限  【废弃】
	OptAuth       []string  `json:"optAuth"`        // 功能权限
	DeptIds       []int64   `json:"deptIds"`       // 部门权限
	RoleIds       []int64   `json:"roleIds"`       // 角色权限
	Creator       int64     `json:"creator"`
	CreateTime    time.Time `json:"createTime"`
}

// ManageGroupHeadData 管理组列表附属信息
type ManageGroupHeadData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// ManageGroupHeadDataUser 管理组列表附属信息-用户信息
type ManageGroupHeadDataUser struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	PhoneNumber string   `json:"phoneNumber"`
}

// ManageGroupDetailResp 管理详情信息返回值
type ManageGroupDetailResp struct {
	AdminGroup *ManageGroupInfo `json:"adminGroup"` // AdminGroup 管理组信息
	HeadData   *ManageGroupHead `json:"headData"`   // HeadData 附属信息
}

// ManageGroupHead
type ManageGroupHead struct {
	Users       []ManageGroupHeadDataUser `json:"users"`       // Users
	AppPackages []ManageGroupHeadData `json:"appPackages"` // AppPackages
	Apps        []ManageGroupHeadData `json:"apps"`        // Apps
	Depts       []ManageGroupHeadData `json:"depts"`       // Depts
	Roles       []ManageGroupHeadData `json:"roles"`       // Roles
}

type GetOperationConfigResp struct {
	// OptAuthList 权限项列表
	OptAuthList []GetOperationConfigRespItem `json:"optAuthList"`
	// OptAuthExtraInfo 权限项的额外信息
	OptAuthExtraInfo map[string]interface{} `json:"optAuthExtraInfo"`
}

type GetOperationConfigRespItem struct {
	// Type 权限的分类。功能权限：`FuncPermission`；菜单权限：`MenuPermission`
	Type string `json:"type"`
	// Code 权限项标识
	Code string `json:"code"`
	// Name 权限项名称
	Name string `json:"name"`
	// Group 权限组名称
	Group string `json:"group"`
	// GroupCode 权限组标识
	GroupCode string `json:"groupCode"`
	// ShowStatus 权限项是否展示的状态。1展示；2不展示。
	ShowStatus int `json:"showStatus"`
}
