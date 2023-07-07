package open_req

// IdReq Id 请求 数据模型
type IdReq struct {
	Id int64 `json:"id"` // Id
}

// IdReq Id 请求 数据模型
type IdsReq struct {
	Ids []int64 `json:"ids"` // Ids
}

// MemberQueryReq  请求 数据模型
type MemberQueryReq struct {
	UserIds   []int64 `json:"userIds"`   // UserIds 用户ID列表 len == 0 则忽视
	Nickname  string  `json:"nickname"`  // Nickname 昵称 全模糊
	LoginName string  `json:"loginName"` // LoginName 用户名 全模糊
	Mobile    string  `json:"mobile"`    // Mobile 手机号 全模糊
	Email     string  `json:"email"`     // Email 邮箱 全模糊
	Status    int64   `json:"status"`    // Status 状态 有效值(1有效 2无效  3离职)
}
