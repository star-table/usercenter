package resp

import "github.com/star-table/usercenter/core/types"

// 部门列表响应结构体
type DepartmentList struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	List []*Department `json:"list"`
}

// 部门结构体
type Department struct {
	// 主键
	ID int64 `json:"id"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 部门名称
	Name string `json:"name"`
	// 部门标识
	Code string `json:"code"`
	// 父部门id
	ParentID int64 `json:"parentId"`
	// 排序
	Sort int `json:"sort"`
	// 是否隐藏部门,1隐藏,2不隐藏
	IsHide int `json:"isHide"`
	// 部门负责人信息
	LeaderInfo []DeptLeaderUserInfo `json:"leaderInfo"`
	// 部门人数，包含子部门
	DeptUserCount int `json:"deptUserCount"`
	// 来源渠道,
	SourceChannel string `json:"sourceChannel"`
	// 创建人
	Creator int64 `json:"creator"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 是否可编辑
	Editable bool `json:"editable"`
	// 是否可删除
	Deletable bool `json:"deletable"`
}

// DeptLeaderUserInfo
type DeptLeaderUserInfo struct {
	// id
	UserID int64 `json:"userId"`
	// 姓名
	Name string `json:"name"`
	// 姓名拼音
	NamePy string `json:"namePy"`
	// 用户头像
	Avatar string `json:"avatar"`
}

// 部门用户信息
type DepartmentMemberInfo struct {
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
}

type DepartmentListData struct {
	//部门id
	DepartmentId int64 `json:"departmentId"`
	//是否是部门负责人1是2否
	IsLeader int `json:"isLeader"`
	// 部门名称
	Name string `json:"name"`
}
