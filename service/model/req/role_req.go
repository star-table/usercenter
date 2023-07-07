package req

/**
角色相关的响应数据模型
*/

// RoleMemberListReq 角色成员列表查询参数
type RoleMemberListReq struct {
	RoleId int64 `json:"roleId"`
}

// RoleGroupReq 角色组请求结构体
type RoleGroupReq struct {
	// 名称
	Name string `json:"name"`
}

// UpdateRoleGroupReq 修改角色组请求结构体
type UpdateRoleGroupReq struct {
	// ID
	Id int64 `json:"id"`
	RoleGroupReq
}

// CreateRoleReq 创建角色请求结构体
type CreateRoleReq struct {
	// 名称
	Name string `json:"name"`
	// 角色组ID
	RoleGroupId int64 `json:"roleGroupId"`
}

// UpdateRoleReq 更新角色请求结构体
type UpdateRoleReq struct {
	// 角色ID
	Id int64 `json:"id"`
	// 角色名称, 必填
	Name string `json:"name"`
}

// MoveRoleReq 移动角色请求结构体
type MoveRoleReq struct {
	// 角色ID
	Id int64 `json:"id"`
	// 角色名称, 必填
	RoleGroupId int64 `json:"roleGroupId"`
}

// AssignRoleRqe 分配角色请求结构体
type AssignRoleRqe struct {
	// 用户ID
	UserId int64 `json:"userId"`
	// 角色ID列表
	RoleIds []int64 `json:"roleIds"`
}

// UpdateRoleUserRefRqe 修改角色用户请求结构体
type UpdateRoleUserRefRqe struct {
	// 角色ID
	RoleId int64 `json:"roleId"`
	// 用户ID列表
	UserIds []int64 `json:"userIds"`
}

// 角色过滤查询请求结构体
type RoleFilterReq struct {
	// 搜索内容
	Query string `json:"query"`
}
