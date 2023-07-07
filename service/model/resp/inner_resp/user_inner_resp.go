package inner_resp

import (
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/bo"
)

// UserAuthorityInnerResp
type UserAuthorityInnerResp struct {
	// 组织ID
	OrgId int64 `json:"orgId"`
	// 用户id
	UserId int64 `json:"userId"`
	// 是否是外部协作人
	IsOutCollaborator bool `json:"isOutCollaborator"`
	// 是否是组织拥有者
	IsOrgOwner bool `json:"isOrgOwner"`
	// 是否是系统管理员
	IsSysAdmin bool `json:"isSysAdmin"`
	// 是否是子管理员
	IsSubAdmin bool `json:"isSubAdmin"`
	// 组织渠道来源。可能用于后续的权限校验规则。
	OrgSourceChannel string `json:"orgSourceChannel"`
	// 所管理部门ID列表
	RefDeptIds []int64 `json:"refDeptIds"`
	// 所管理角色ID列表
	RefRoleIds []int64 `json:"refRoleIds"`
	// 所拥有操作项
	OptAuth []string `json:"optAuth"`
	// 管理部门的权限
	HasDeptOptAuth bool `json:"hasDeptOptAuth"`
	// 角色操作权限
	HasRoleOptAuth bool `json:"hasRoleOptAuth"`
	// 应用包操作权限
	HasAppPackageOptAuth bool `json:"hasAppPackageOptAuth"`
	// 管理的应用包列表
	ManageAppPackages []int64 `json:"manageAppPackages"`
	ManageApps        []int64 `json:"manageApps"` // ManageApps 管理的应用列表
	// 管理的角色列表
	ManageRoles []int64 `json:"manageRoles"`
	// 管理的部门列表
	ManageDepts []int64 `json:"manageDepts"`
}

type OrgUserPerContext struct {
	auth              *UserAuthorityInnerResp
	manageDepts       map[int64]bool
	manageRoles       map[int64]bool
	manageAppPackages map[int64]bool
}

func NewOrgUserPermissionContext(auth UserAuthorityInnerResp) *OrgUserPerContext {
	manageDepts := make(map[int64]bool)
	manageRoles := make(map[int64]bool)
	manageAppPackages := make(map[int64]bool)
	for _, id := range auth.ManageDepts {
		manageDepts[id] = auth.HasDeptOptAuth
	}
	for _, id := range auth.ManageRoles {
		manageRoles[id] = auth.HasRoleOptAuth
	}
	for _, id := range auth.ManageAppPackages {
		manageAppPackages[id] = auth.HasAppPackageOptAuth
	}
	return &OrgUserPerContext{
		auth:              &auth,
		manageDepts:       manageDepts,
		manageRoles:       manageRoles,
		manageAppPackages: manageAppPackages,
	}
}

// 组织拥有者
func (u *OrgUserPerContext) IsOrgOwner() bool {
	return u.auth.IsOrgOwner
}

// 具备所有权限
func (u *OrgUserPerContext) HasAllPermission() bool {
	return u.auth.IsOrgOwner || u.auth.IsSysAdmin
}

// HasInManageGroup 是否在管理组中
func (u *OrgUserPerContext) HasInManageGroup() bool {
	return u.HasAllPermission() || u.auth.IsSubAdmin
}

// HasCreateDept 是否有添加部门的权限
func (u *OrgUserPerContext) HasCreateDept() bool {
	polarisHasManageDept, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgDeptCreate)
	return u.HasAllPermission() || (u.auth.IsSubAdmin && u.auth.HasDeptOptAuth) || polarisHasManageDept
}

// HasManageDept 是否有管理该部门的权限
func (u *OrgUserPerContext) HasManageDept(deptId int64) bool {
	polarisHasManageDept, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgDeptCreate)
	return u.HasAllPermission() || (u.auth.IsSubAdmin && u.auth.HasDeptOptAuth && (len(u.manageDepts) == 0 || u.manageDepts[deptId])) || polarisHasManageDept
}

// HasManageAllDept 是否有管理所有部门的权限
func (u *OrgUserPerContext) HasManageAllDept() bool {
	polarisHasManageDept, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgDeptCreate)
	return u.HasAllPermission() || (u.auth.IsSubAdmin && u.auth.HasDeptOptAuth && len(u.manageDepts) == 0) || polarisHasManageDept
}

func (u *OrgUserPerContext) HasManageDepts(deptIds []int64) bool {
	polarisHasManageDept, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgDeptCreate)
	if polarisHasManageDept {
		return true
	}
	if u.HasAllPermission() {
		return true
	}
	if u.auth.IsSubAdmin && u.auth.HasDeptOptAuth {
		if len(u.manageDepts) == 0 {
			return true
		}
		for _, id := range deptIds {
			if !u.manageDepts[id] {
				return false
			}
		}
		return true
	}
	return false
}

// HasManageRole 是否有管理该角色的权限
func (u *OrgUserPerContext) HasManageRole(roleId int64) bool {
	polarisHasManageRole, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgAdminGroupCreate)
	return u.HasAllPermission() || (u.auth.IsSubAdmin && u.auth.HasRoleOptAuth && (len(u.manageRoles) == 0 || u.manageRoles[roleId])) || polarisHasManageRole
}

// HasManageAllRole 是否有所有角色的权限
func (u *OrgUserPerContext) HasManageAllRole() bool {
	polarisHasManageRole, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgAdminGroupCreate)
	return u.HasAllPermission() || (u.auth.IsSubAdmin && u.auth.HasRoleOptAuth && len(u.manageRoles) == 0) || polarisHasManageRole
}

func (u *OrgUserPerContext) HasManageRoles(roleIds []int64) bool {
	polarisHasManageRole, _ := slice.Contain(u.auth.OptAuth, consts.OperationOrgAdminGroupCreate)
	if polarisHasManageRole {
		return true
	}
	if u.HasAllPermission() {
		return true
	}
	if u.auth.IsSubAdmin && u.auth.HasRoleOptAuth {
		if len(u.manageRoles) == 0 {
			return true
		}
		for _, id := range roleIds {
			if !u.manageRoles[id] {
				return false
			}
		}
		return true
	}
	return false
}

// HasOpForPolaris 用于极星的是否具备某个权限。如果不是极星来源，则直接通过。
func (u *OrgUserPerContext) HasOpForPolaris(operationCode string) bool {
	if u.HasAllPermission() {
		return true
	}
	return u.HasOperation(operationCode)
}

// CheckOrgSourceChannelIsPolaris 检查组织来源是否是极星的 source_channel
func (u *OrgUserPerContext) CheckOrgSourceChannelIsPolaris() bool {
	return CheckSourceChannelIsPolaris(u.auth.OrgSourceChannel)
}

func (u *OrgUserPerContext) CheckIsOrgSourceChannel(sourceChannel string) bool {
	return u.auth.OrgSourceChannel == sourceChannel
}

// HasOperation 检验是否有某个权限的操作项
func (u *OrgUserPerContext) HasOperation(operationCode string) bool {
	exist, _ := slice.Contain(u.auth.OptAuth, operationCode)
	return exist
}

// GetManageDeptIds 获取管理的所有部门 注意管理全部部门时为空
func (u *OrgUserPerContext) GetManageDeptIds() []int64 {
	return u.auth.ManageDepts
}

// GetManageRoleIds 获取管理的所有角色 注意管理全部角色时为空
func (u *OrgUserPerContext) GetManageRoleIds() []int64 {
	return u.auth.ManageRoles
}

// CheckSourceChannelIsPolaris 校验组织的 source_channel，是否是 `polaris`
func CheckSourceChannelIsPolaris(sourceChannel string) bool {
	polarisEnum := []string{"fs", "web", "ding", "lark-xyjh2019"}
	isPolaris, _ := slice.Contain(polarisEnum, sourceChannel)
	return isPolaris
}

// ObjOptAuthItem
type ObjOptAuthItem struct {
	// ID
	Id int64 `json:"id"`
	// 是否具有管理权限
	HasManage bool `json:"hasManage"`
}

// UserInfoInnerResp
type UserInfoInnerResp struct {
	Id          int64     `json:"id"`          // Id 用户ID
	LoginName   string    `json:"loginName"`   // LoginName
	Name        string    `json:"name"`        // Name 姓名
	NamePy      string    `json:"namePy"`      // NamePy 姓名拼音
	Avatar      string    `json:"avatar"`      // Avatar 用户头像
	Email       string    `json:"email"`       // Email 邮箱
	PhoneRegion string    `json:"phoneRegion"` // PhoneRegion 手机区号
	PhoneNumber string    `json:"phoneNumber"` // PhoneNumber 手机
	Status      int       `json:"status"`      // Status 状态1启用2禁用3离职
	Creator     int64     `json:"creator"`     // Creator 组织成员 创建者
	CreateTime  time.Time `json:"createTime"`  // CreateTime 创建时间
	Updator     int64     `json:"updator"`     // Updator 组织成员 最后修改者ID
	UpdateTime  time.Time `json:"updateTime"`  // UpdateTime 创建时间
	IsDelete    int       `json:"isDelete"`    // IsDelete
	Type        int       `json:"type"`        // Type 用户类型，1：内部，2：外部
	UserBindDeptAndRoleResp
}

// UserDeptBindResp
type UserDeptBindResp struct {
	SimpleMemberInfo
	UserDeptBindData
}

type SimpleMemberInfo struct {
	UserId   int64  `json:"userId"`
	Nickname string `json:"nickname"`
}

// UserBindDeptAndRoleResp
type UserBindDeptAndRoleResp struct {
	RoleList []UserRoleBindData `json:"roleList"` // RoleList 角色列表
	DeptList []UserDeptBindData `json:"deptList"` // DeptList 部门列表
}

// UserRoleBindData
type UserRoleBindData struct {
	RoleId   int64  `json:"roleId"`
	RoleName string `json:"roleName"`
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

type GetMemberSimpleInfoResp struct {
	Data []SimpleInfo `json:"data"`
}

type SimpleInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
	ParentId int64  `json:"parentId"` //父部门id
}

type RepeatMemberInfo struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	Department []string `json:"department"`
}

type RepeatMemberResp struct {
	User       []RepeatMemberInfo `json:"user"`
	Department []RepeatMemberInfo `json:"department"`
	Role       []RepeatMemberInfo `json:"role"`
}

type GetUsersCouldManageResp struct {
	// 用户列表
	List []bo.SimpleUserInfoBo
}

type UserListResp struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	List []*OrgMemberInfoReq `json:"list"`
}

// 组织成员详情
type OrgMemberInfoReq struct {
	// id
	UserID int64 `json:"userId"`
	// LoginName 用户名
	LoginName string `json:"loginName"`
	// 姓名
	Name string `json:"name"`
	// 姓名拼音
	NamePy string `json:"namePy"`
	// 用户头像
	Avatar string `json:"avatar"`
	// 邮箱
	Email string `json:"email"`
	// 手机区号
	PhoneRegion string `json:"phoneRegion"`
	// 手机号码
	PhoneNumber string `json:"phoneNumber"`
	// 用户所在的管理组
	AdminGroupList []bo.ManageGroupInfoBo `json:"adminGroupList"`
	// 状态1启用2禁用
	Status int `json:"status"`
	// 禁用时间
	StatusChangeTime time.Time `json:"statusChangeTime"`
	// 是否是组织创建人
	IsCreator bool `json:"isCreator"`
	// 是否是组织拥有者
	IsOwner bool `json:"isOwner"`
	// 工号
	EmpNo string `json:"empNo"`
	// 微博ID
	WeiboIds []string `json:"weiboIds"`
	// 组织内用户创建者
	Creator int64 `json:"creator"`
	// 组织内用户创建者名称
	CreatorName int64 `json:"creatorName"`
	// 创建时间
	CreateTime time.Time `json:"createTime"`
}

type MemberSimpleInfoListResp struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	List []*MemberSimpleInfo `json:"list"`
}

// 组织成员详情
type MemberSimpleInfo struct {
	// id
	Id int64 `json:"id"`
	// 姓名
	Name string `json:"name"`
	// 状态1启用2禁用
	Status int `json:"status"`
	// 类型 U表示user D表示部门
	Type string `json:"type"`
	// isDelete = 2
	IsDelete int `json:"isDelete"`
}

type GetCommAdminMangeAppsData struct {
	UserId int64   `json:"userId"`
	AppIds []int64 `json:"appIds"` // [-1] 表示可以管理全部应用
}
