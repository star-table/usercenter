package open_resp

import "time"

// OrgMemberBaseResp
type OrgMemberBaseResp struct {
	UserID      int64     `json:"userId"`      // UserID
	LoginName   string    `json:"loginName"`   // LoginName
	Name        string    `json:"name"`        // Name 姓名
	NamePy      string    `json:"namePy"`      // NamePy 姓名拼音
	Avatar      string    `json:"avatar"`      // Avatar 用户头像
	Email       string    `json:"email"`       // Email 邮箱
	PhoneRegion string    `json:"phoneRegion"` // PhoneRegion 手机区号
	PhoneNumber string    `json:"phoneNumber"` // PhoneNumber 手机
	EmpNo       string    `json:"empNo"`       // EmpNo 工号
	WeiboIds    string    `json:"weiboIds"`    // WeiboIds 微博ID
	Status      int       `json:"status"`      // Status 状态1启用2禁用
	OrgOwner    int64     `json:"orgOwner"`    // OrgOwner 组织拥有者
	Creator     int64     `json:"creator"`     // Creator 组织成员 创建者
	CreateTime  time.Time `json:"createTime"`  // CreateTime 创建时间
	Updator     int64     `json:"updator"`     // Updator 组织成员 最后修改者ID
	UpdateTime  time.Time `json:"updateTime"`  // UpdateTime 创建时间
	UserBindDeptAndRoleResp
}

// UserDeptBindResp
type UserDeptBindResp struct {
	SimpleMemberInfo
	UserDeptBindData
}

// UserDeptBindData
type UserDeptBindData struct {
	DepartmentId             int64  `json:"departmentId"`             // DepartmentId 部门ID
	DepartmentName           string `json:"departmentName"`           // DepartmentName 部门名称
	OutOrgDepartmentId       string `json:"outOrgDepartmentId"`       // OutOrgDepartmentId
	OutOrgDepartmentCode     string `json:"outOrgDepartmentCode"`     // OutOrgDepartmentCode
	OutOrgDepartmentParentId string `json:"outOrgDepartmentParentId"` // OutOrgDepartmentParentId
	IsLeader                 int    `json:"isLeader"`                 // IsLeader 是否是部门负责人 1是 2否
	PositionId               int64  `json:"positionId"`               // PositionId 注意此处是org内的职级ID
	PositionName             string `json:"positionName"`             // PositionName 职级名称
	PositionLevel            int    `json:"positionLevel"`            // PositionLevel 职级等级
}

// UserRoleBindResp
type UserRoleBindResp struct {
	SimpleMemberInfo
	UserRoleBindData
}

// UserRoleBindData
type UserRoleBindData struct {
	RoleId   int64  `json:"roleId"`
	RoleName string `json:"roleName"`
}

// UserBindDeptAndRoleResp
type UserBindDeptAndRoleResp struct {
	RoleList []UserRoleBindData `json:"roleList"` // RoleList 角色列表
	DeptList []UserDeptBindData `json:"deptList"` // DeptList 部门列表
}

type SimpleMemberInfo struct {
	UserId   int64  `json:"userId"`
	Nickname string `json:"nickname"`
	Status   int    `json:"status"` // Status
}

// OrgUserAuthBaseResp
type OrgUserAuthBaseResp struct {
	UserID     int64   `json:"userId"`     // UserID
	LoginName  string  `json:"loginName"`  // LoginName
	Nickname   string  `json:"nickname"`   // Nickname
	IsOrgOwner bool    `json:"isOrgOwner"` // IsOrgOwner 是否是组织拥有者
	IsSysAdmin bool    `json:"isSysAdmin"` // IsSysAdmin 是否是系统管理员
	IsSubAdmin bool    `json:"isSubAdmin"` // IsSubAdmin 是否是子管理员
	OrgOwner   int64   `json:"orgOwner"`   // OrgOwner 组织拥有者
	DeptIds    []int64 `json:"deptIds"`    // DeptIds 所属部门ID列表
	RoleIds    []int64 `json:"roleIds"`    // RoleIds 所属角色ID列表
}
