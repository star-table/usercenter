package req

import "time"

// CreateOrgMemberReq 创建组织成员
type CreateOrgMemberReq struct {
	// 账号
	LoginName string `json:"loginName"`
	// 密码
	Password string `json:"password"`
	// 手机区号
	PhoneRegion string `json:"phoneRegion"`
	// 手机号
	PhoneNumber string `json:"phoneNumber"`
	// 邮箱
	Email string `json:"email"`
	// 昵称
	Name string `json:"name"`
	// 部门id
	DeptAndPositions []DeptAndPositionReq `json:"deptAndPositions"`
	// 角色id
	RoleIds []int64 `json:"roleIds"`
	// 状态（1启用2禁用3离职, 默认启用） 忽视0
	Status int `json:"status"`
	// 工号
	EmpNo string `json:"empNo"`
	// 微博
	WeiboIds []string `json:"weiboIds"`
	// 昵称
	Avatar string `json:"avatar"`
}

// UpdateOrgMemberReq 更新组织成员信息
type UpdateOrgMemberReq struct {
	// UserId
	UserId int64 `json:"userId"`
	// 手机区号。可选，不传表示不更新该字段。
	PhoneRegion *string `json:"phoneRegion"`
	// 手机号。可选
	PhoneNumber *string `json:"phoneNumber"`
	// 邮箱。可选
	Email *string `json:"email"`
	// 姓名。可选
	Name *string `json:"name"`
	// 部门/职级列表
	DeptAndPositions []DeptAndPositionReq `json:"deptAndPositions"`
	// 用户所在的部门 ids. polaris 没有职级,因此无法使用 DeptAndPositions 字段
	DepartmentIds []int64 `json:"departmentIds"`
	// 角色id
	RoleIds []int64 `json:"roleIds"`
	// 管理组。一个用户可以在多个管理组中。
	AdminGroup []int64 `json:"adminGroup"`
	// 状态（1启用2禁用3离职 选填）
	Status *int `json:"status"`
	// 工号。可选
	EmpNo *string `json:"empNo"`
	// 微博。
	WeiboIds []string `json:"weiboIds"`
	// 头像。可选
	Avatar *string `json:"avatar"`
}

type UserListReq struct {
	//搜索字段
	SearchCode *string `json:"searchCode"`
	//是否已分配部门（1已分配2未分配，默认全部）
	IsAllocate *int `json:"isAllocate"`
	//是否禁用（1启用2禁用，默认全部）
	Status *int `json:"status"`
	// 审核状态。1待审核；2审核通过；3审核不通过
	CheckStatus []int `json:"checkStatus"`
	//角色id
	RoleId *int64 `json:"roleId"`
	//部门id
	DepartmentId *int64 `json:"departmentId"`
	//职级ID
	PositionId *int64 `json:"positionId"`
	// 页码（选填，不填为全部）
	Page int `json:"page"`
	// 数量（选填，不填为全部）
	Size int `json:"size"`
	// 是否开启权限过滤
	AuthFilter bool `json:"authFilter"`
	// 排序字段 如"create_time desc,id asc"
	Order string `json:"order"` // Order
}

type InviteUserReq struct {
	Data []InviteUserData `json:"data"`
}

type InviteUserData struct {
	//邮箱
	Email string `json:"email"`
	//姓名（再次邀请时不用传了）
	Name string `json:"name"`
}

type SearchUserReq struct {
	Email string `json:"email"`
}

type InviteUserListReq struct {
	// 页码（选填，不填为全部）
	Page int `json:"page"`
	// 数量（选填，不填为全部）
	Size int `json:"size"`
}

type ExportAddressListReq struct {
	//搜索字段
	SearchCode *string `json:"searchCode"`
	//是否已分配部门（1已分配2未分配，默认全部）
	IsAllocate *int `json:"isAllocate"`
	//是否禁用（1启用2禁用，默认全部）
	Status *int `json:"status"`
	//角色id
	RoleId *int64 `json:"roleId"`
	//部门id
	DepartmentId *int64 `json:"departmentId"`
	//导出字段(name,mobile,empNo,email,department,role,isLeader,statusChangeTime,createTime)
	ExportField []string `json:"exportField"`
	// 是否开启权限过滤
	AuthFilter bool `json:"authFilter"`
}

type RemoveInviteUserReq struct {
	//要移除的id
	Ids []int64 `json:"ids"`
	//是否删除全部(1是，其余为否)
	IsAll int `json:"isAll"`
}

// 更改用户个人信息
type UpdateUserInfoReq struct {
	// 姓名
	Name *string `json:"name"`
	// 性别
	Sex *int `json:"sex"`
	// 用户头像
	Avatar *string `json:"avatar"`
	// 生日
	Birthday *time.Time `json:"birthday"`
	// 变动的字段列表
	UpdateFields []string `json:"updateFields"`
}

//获取邀请信息
type GetInviteInfoReq struct {
	InviteCode string `json:"inviteCode"`
}

// 检查用户名是否存在的请求信息
type CheckLoginNameExistReq struct {
	// 用户名，中文
	Name string `json:"name"`
}
