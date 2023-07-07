package req

// 通讯录查询请求结构体
type ContactFilterReq struct {
	// 父部门ID, 传0代表根部门，不需要限制则不传
	ParentID *int64 `json:"parentId"`
	// 成员查询起始位置，默认0
	MemberOffset uint `json:"memberOffset"`
	// 成员查询数量，默认1000
	MemberLimit  uint `json:"memberLimit"`
	// 搜索范围
	Scope ContactScope `json:"scope"`
	// 只查询成员
	OnlyMember bool `json:"onlyMember"`
}

// 通讯录范围
type ContactScope struct {
	// 部门范围
	DepartmentIds *[]int64 `json:"departmentIds"`
	// 用户范围
	UserIds *[]int64 `json:"userIds"`
}

// 通讯录搜索请求结构体
type ContactSearchReq struct {
	// 搜索类型，1：成员，2：部门
	SearchType uint8 `json:"searchType"`
	// 搜索内容
	Query string `json:"query"`
	// Offset, default 0
	Offset uint `json:"offset"`
	// Limit, default 1000
	Limit uint `json:"limit"`
	// 搜索范围
	Scope ContactScope `json:"scope"`
	// 只查询成员
	OnlyMember bool `json:"onlyMember"`
}

type AggregationReq struct {
	// 搜索内容
	Query string `json:"query"`
}
