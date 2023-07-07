package req

type DepartmentListReq struct {
	// 父部门id（选填）
	ParentID *int64 `json:"parentId"`
	// 是否查询最上级部门, 如果是1则为true（选填）
	IsTop *int `json:"isTop"`
	// 是否显示隐藏的部门，如果是1则为true，默认不显示（选填）
	ShowHiding *int `json:"showHiding"`
	// 部门名称（选填）
	Name *string `json:"name"`
	// 页码（选填，不填为全部）
	Page *int64 `json:"page"`
	// 数量（选填，不填为全部）
	Size *int64 `json:"size"`
	// 部门ID范围
	DeptIds *[]int64 `json:"deptIds"`
	// 是否返回部门的人数。0、2不展示，1展示
	IsShowDeptUserCount int `json:"isShowDeptUserCount"`
}

type DepartmentMemberListReq struct {
	// 部门id（选填，不传为全部）
	DepartmentID *int64 `json:"departmentId"`
}

type CreateDepartmentReq struct {
	// 父部门id（根目录为0）
	ParentID int64 `json:"parentId"`
	// 部门名称
	Name string `json:"name"`
	// 部门负责人(选填)
	LeaderIds []int64 `json:"leaderIds"`
}

type UpdateDepartmentReq struct {
	//部门id
	DepartmentId int64 `json:"departmentId"`
	//部门名称（选填）
	Name string `json:"name"`
	// 部门负责人(选填，不传表示不更新，传空数组则表示取消部门负责人)
	LeaderIds []int64 `json:"leaderIds"`
	// 部门下的普通成员(选填，不传表示不更新，传空数组则表示删除所有部门成员)
	UserIds []int64 `json:"userIds"`
}

type DeleteDepartmentReq struct {
	//部门id
	DepartmentId int64 `json:"departmentId"`
}

type AllocateDepartmentReq struct {
	//用户id
	UserIds []int64 `json:"userIds"`
	//部门id
	DepartmentIds []int64 `json:"departmentIds"`
}

// DeptAndPositionReq 部门和职级请求数据模型
type DeptAndPositionReq struct {
	// DepartmentId 部门ID(必填)
	DepartmentId int64 `json:"departmentId"`
	// PositionId 职级ID(必填 注意是Org内的OrgPositionId)
	PositionId int64 `json:"positionId"`
}

// ChangeUserDeptAndPositionReq 修改用户部门和职级请求数据模型
type ChangeUserDeptAndPositionReq struct {
	// UserId 用户ID(必填)
	UserId int64 `json:"userId"`
	// DeptAndPositions 部门和职级ID
	DeptAndPositions []DeptAndPositionReq `json:"deptAndPositions"`
}

// UpdateDeptLeaderReq 修改部门负责人
type UpdateDeptLeaderReq struct {
	//用户id
	UserId int64 `json:"userId"`
	//是否是部门负责人(1是2否)
	IsLeader int `json:"isLeader"`
	//部门id
	DepartmentId int64 `json:"departmentId"`
}

// ChangeDeptSortReq 交换部门排序
type ChangeDeptSortReq struct {
	// 部门ID
	DeptId int64 `json:"deptId"`
	// 部门ID
	Sort int `json:"sort"`
}

// ChangeUserAdminGroupReq 切换用户管理组
type ChangeUserAdminGroupReq struct {
	// UserId 用户 id，为该用户变更管理组
	UserId int64 `json:"userId"`
	// DstAdminGroupId 切换后的管理组 id
	DstAdminGroupId int64 `json:"dstAdminGroupId"`
	// DstAdminGroupIds 切换后的管理组 id 支持多个
	DstAdminGroupIds []int64 `json:"dstAdminGroupIds"`
	// 验证码。如果是更换**超管**，则必须传该字段，通过手机号获取验证码。
	AuthCode string `json:"authCode"`
	// SourceFrom 应用来源：极星应用：`polaris`；
	SourceFrom string `json:"sourceFrom"`
}

// AllocateUserDeptReq 给用户分配部门
type AllocateUserDeptReq struct {
	// UserIds 给这些用户分配部门
	UserIds []int64 `json:"userIds"`
	// DstDeptIds 把用户分配给这些部门
	DstDeptIds []int64 `json:"dstDeptIds"`
}

// DeptRemoveUsersReq 将多个用户移出某个部门
type DeptRemoveUsersReq struct {
	// UserIds 用户 id，将这些用户移出部门
	UserIds []int64 `json:"userIds"`
	// 部门 id
	DeptId int64 `json:"deptId"`
}
