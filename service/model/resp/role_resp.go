package resp

import "time"

// TreeRoleGroupNode 角色树节点
type TreeRoleGroupNode struct {
	RoleGroup
	UserIds []int64        `json:"userIds"`
	Roles   []TreeRoleNode `json:"roles"`
}

type TreeRoleNode struct {
	Id          int64     `json:"id"`
	OrgId       int64     `json:"orgId"`
	RoleGroupId int64     `json:"roleGroupId"`
	Name        string    `json:"name"`
	Creator     int64     `json:"creator"`
	CreateTime  time.Time `json:"createTime"`
	Updator     int64     `json:"updator"`
	UpdateTime  time.Time `json:"updateTime"`

	UserIds []int64 `json:"userIds"`
}

type RoleGroup struct {
	Id         int64     `json:"id"`
	OrgId      int64     `json:"orgId"`
	Name       string    `json:"name"`
	Creator    int64     `json:"creator"`
	CreateTime time.Time `json:"createTime"`
	Updator    int64     `json:"updator"`
	UpdateTime time.Time `json:"updateTime"`
}

// 角色用户信息
type RoleMemberInfo struct {
	// id
	UserID int64 `json:"userId"`
	// 姓名
	Name string `json:"name"`
	// 姓名拼音
	NamePy string `json:"namePy"`
	// 用户头像
	Avatar string `json:"avatar"`
	// 工号：企业下唯一
	EmplID string `json:"emplId"`
	// unionId： 开发者账号下唯一
	UnionID string `json:"unionId"`
	// 用户部门信息
	DepartmentList []DepartmentListData `json:"departmentList"`
	// 用户角色信息
	RoleList []RoleListData `json:"roleList"`
}

type RoleListData struct {
	RoleId int64 `json:"roleId"`
}

// RoleListResp 角色列表 响应
type RoleListResp struct {
	// RoleGroups 角色组列表
	RoleGroups []RoleGroupInfo `json:"roleGroups"`
	// Roles 角色列表
	Roles []RoleInfo `json:"roles"`
}

// RoleGroupInfo 角色列表-角色组信息
type RoleGroupInfo struct {
	Id         int64     `json:"id"`
	OrgId      int64     `json:"orgId"`
	Name       string    `json:"name"`
	Creator    int64     `json:"creator"`
	CreateTime time.Time `json:"createTime"`
	Updator    int64     `json:"updator"`
	UpdateTime time.Time `json:"updateTime"`
}

// RoleInfo 角色列表-角色信息
type RoleInfo struct {
	Id          int64     `json:"id"`
	OrgId       int64     `json:"orgId"`
	RoleGroupId int64     `json:"roleGroupId"`
	Name        string    `json:"name"`
	Creator     int64     `json:"creator"`
	CreateTime  time.Time `json:"createTime"`
	Updator     int64     `json:"updator"`
	UpdateTime  time.Time `json:"updateTime"`

	// Editable 是否可编辑成员
	EditableMember bool `json:"editableMember"`
}

// 角色过滤响应结构体
type RoleFilterResp struct {
	// 角色列表
	RoleList []RoleNode `json:"roleList"`
}

// 角色结点信息
type RoleNode struct {
	// 角色名称
	ID int64 `json:"id"`
	// 角色ID
	Name string `json:"name"`
	// 分组ID
	GroupID int64 `json:"groupId"`
	// 分组名称
	GroupName string `json:"groupName"`
	// 成员数
	UserCount uint64 `json:"userCount"`
}
