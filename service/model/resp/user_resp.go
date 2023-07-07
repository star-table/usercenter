package resp

import (
	"time"

	"github.com/star-table/usercenter/core/types"
	"github.com/star-table/usercenter/service/model/bo"
)

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
	// 部门/职级信息
	DepartmentList []UserDeptPositionData `json:"departmentList"`
	// 用户所在的管理组
	AdminGroupList []bo.ManageGroupInfoBo `json:"adminGroupList"`
	// 角色信息
	RoleList []UserRoleData `json:"roleList"`
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
	// 是否有设置密码
	HasPassword bool `json:"hasPassword"`
}

// UserDeptPositionData 用户的部门和职级信息
type UserDeptPositionData struct {
	// 部门id
	DeparmentName string `json:"deparmentName"`
	// 部门名称
	DepartmentId int64 `json:"departmentId"`
	// 是否是部门负责人1是2否
	IsLeader int `json:"isLeader"`
	// 职级ID 注意此处是Org内的职级ID
	PositionId int64 `json:"positionId"`
	// 职级名称
	PositionName string `json:"positionName"`
	// 职级级别
	PositionLevel int `json:"positionLevel"`
}

type UserRoleData struct {
	//角色id
	RoleId int64 `json:"roleId"`
	//角色名称
	RoleName string `json:"roleName"`
}

type SearchUserResp struct {
	//搜索结果（1可邀请2已邀请3已添加）
	Status   int               `json:"status"`
	UserInfo *OrgMemberInfoReq `json:"userInfo"`
}

type InviteUserResp struct {
	//成功的邮箱
	SuccessEmail []string `json:"successEmail"`
	//已邀请的邮箱
	InvitedEmail []string `json:"invitedEmail"`
	//已经是用户的邮箱
	IsUserEmail []string `json:"isUserEmail"`
	//不符合规范的邮箱
	InvalidEmail []string `json:"invalidEmail"`
}

type InviteUserListResp struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	List []InviteUserInfo `json:"list"`
}

type InviteUserInfo struct {
	//用户id
	Id int64 `json:"id"`
	//名称
	Name string `json:"name"`
	//邮箱
	Email string `json:"email"`
	//邀请时间
	InviteTime time.Time `json:"inviteTime"`
	//是否24h内已邀请
	IsInvitedRecent bool `json:"isInvitedRecent"`
}

type ExportAddressListResp struct {
	Url string `json:"url"`
}

type GetInviteCodeResp struct {
	InviteCode string `json:"inviteCode"`
	Expire     int    `json:"expire"`
}

type UserStatResp struct {
	//所有成员数量
	AllCount int64 `json:"allCount"`
	//未分配成员数量
	UnallocatedCount int64 `json:"unallocatedCount"`
	//未接受邀请成员数量
	UnreceivedCount int64 `json:"unreceivedCount"`
	//已禁用成员数量
	ForbiddenCount int64 `json:"forbiddenCount"`
	//已离职成员数量
	ResignedCount int64 `json:"resignedCount"`
	//待审核成员数量
	WaitAuditCount int64 `json:"waitAuditCount"`
}

// 个人信息
type PersonalInfo struct {
	ID                 int64      `json:"id"`                 // 主键
	EmplID             *string    `json:"emplId"`             // 工号
	OrgID              int64      `json:"orgId"`              // 组织id
	OrgName            string     `json:"orgName"`            // 组织名称
	OrgCode            string     `json:"orgCode"`            // 组织code
	Name               string     `json:"name"`               // 名称
	ThirdName          string     `json:"thirdName"`          // 第三方名称
	LoginName          string     `json:"loginName"`          // 登录名
	LoginNameEditCount int        `json:"loginNameEditCount"` // 登录名编辑次数
	Email              string     `json:"email"`              // 邮箱
	Mobile             string     `json:"mobile"`             // 电话
	Birthday           types.Time `json:"birthday"`           // 生日
	Sex                int        `json:"sex"`                // 性别
	Rimanente          int        `json:"rimanente"`          // 剩余使用时长
	Level              int        `json:"level"`              // 付费等级
	LevelName          string     `json:"levelName"`          // 付费等级名
	Avatar             string     `json:"avatar"`             // 头像
	SourceChannel      string     `json:"sourceChannel"`      // 来源
	Language           string     `json:"language"`           // 语言
	Motto              string     `json:"motto"`              // 座右铭
	LastLoginIP        string     `json:"lastLoginIp"`        // 上次登录ip
	LastLoginTime      types.Time `json:"lastLoginTime"`      // 上次登录时间
	LoginFailCount     int        `json:"loginFailCount"`     // 登录失败次数
	LastEditPwdTime    types.Time `json:"lastEditPwdTime"`    // LastEditPwdTime
	CreateTime         types.Time `json:"createTime"`         // 创建时间
	UpdateTime         types.Time `json:"updateTime"`         // 更新时间
	PasswordSet        int        `json:"passwordSet"`        // 密码是否设置过(1已设置0未设置)
	RemindBindPhone    int        `json:"remindBindPhone"`    // 是否需要提醒（1需要2不需要）
	IsOrgOwner         bool       `json:"isOrgOwner"`         // IsOrgOwner 是否是组织拥有者
	IsSysAdmin         bool       `json:"isSysAdmin"`         // IsSysAdmin 是否是系统管理员
	IsSubAdmin         bool       `json:"isSubAdmin"`         // IsSubAdmin 是否是子管理员
	IsOutCollaborator  bool       `json:"isOutCollaborator"`  // IsOutCollaborator 是否是外部协作人

	// 部门/职级信息
	DepartmentList []UserDeptPositionData `json:"departmentList"`
	// 角色信息
	RoleList []UserRoleData `json:"roleList"`

	Functions []string `json:"functions"` // 权限
}

// 获取邀请信息响应结构体
type GetInviteInfoResp struct {
	// 组织id
	OrgID int64 `json:"orgId"`
	// 组织名
	OrgName string `json:"orgName"`
	// 邀请人id
	InviterID int64 `json:"inviterId"`
	// 邀请人姓名
	InviterName string `json:"inviterName"`
}

// 解绑成功后返回的数据
type UnbindLoginNameResp struct {
	// 换绑时，第一步请求 unbind 时返回的验证码
	ChangeBindCode string `json:"changeBindCode"`
}
