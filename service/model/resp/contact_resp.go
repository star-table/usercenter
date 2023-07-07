package resp

// 通讯录查询响应结构体
type ContactFilterResp struct {
	// 子部门列表
	DepList []ContactDepartment `json:"depList"`
	// 部门层级，顺序返回
	PathNodes []ContactPathNode `json:"pathNodes"`
	// 是否有更多的用户
	UserHasMore bool `json:"userHasMore"`
	// 部门用户总数
	UserTotalNum uint64 `json:"userTotalNum"`
	// 组织总人数
	OrgTotalNum uint64 `json:"orgTotalNum"`
	// 用户列表
	User []ContactUser `json:"user"`
}

// 通讯录搜索响应结构体
type ContactSearchResp struct {
	// 子部门列表， searchType为2时返回
	DepList []ContactDepartment `json:"depList"`
	// 用户列表， searchType为1时返回
	User []ContactUser `json:"user"`
	// 是否有更多
	HasMore bool `json:"hasMore"`
}

// 通讯录层级结点
type ContactPathNode struct {
	// 部门ID
	DepID int64 `json:"depId,string"`
	// 部门名称
	DepName string `json:"depName"`
	// 父部门ID
	ParentID int64 `json:"parentId,string"`
	// 是否可见
	Visible bool `json:"visible"`
}

// 通讯录部门
type ContactDepartment struct {
	// 部门ID
	ID int64 `json:"id,string"`
	// 部门名称
	Name string `json:"name"`
	// 成员总数
	UserCount uint32 `json:"userCount"`
	// 子部门总数
	ChildrenCount uint32 `json:"childrenCount"`
	// 父ID
	ParentID int64 `json:"parentId,string"`
	// 是否可见
	Visible bool `json:"visible"`
	// 部门层级
	PathNodes []ContactPathNode `json:"pathNodes"`
}

// 通讯录用户
type ContactUser struct {
	// 用户ID
	ID int64 `json:"id,string"`
	// 头像
	Avatar string `json:"avatar"`
	// 名称
	Name string `json:"name"`
	// 邮箱
	Email string `json:"email"`
	// 手机号
	Mobile string `json:"mobile"`
	// 登录名
	LoginName string `json:"loginName"`
	// 部门层级
	PathNodesList [][]ContactPathNode `json:"pathNodesList"`
}

type AggregationResp struct {
	// 子部门列表， searchType为2时返回
	DepList []ContactDepartment `json:"depList"`
	// 用户列表， searchType为1时返回
	User []ContactUser `json:"user"`
	// 角色列表
	RoleList []RoleNode `json:"roleList"`
}
