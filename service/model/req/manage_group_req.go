package req

// CreateManageGroup 新建管理组请求参数
type CreateManageGroup struct {
	// 名称
	Name string `json:"name"`
	// 类型 1.系统管理组 2.普通管理组 3.bjx-普通管理员 4.bjx-普通成员,5.bjx-用户创建
	GroupType int `json:"groupType"`
	// 权限项
	OptAuth []string `json:"optAuth"`
}

// UpdateManageGroup 修改管理组请求参数
type UpdateManageGroup struct {
	// ID
	Id int64 `json:"id"`
	// 名称
	Name string `json:"name"`
}

// DeleteManageGroup 删除管理组请求参数
type DeleteManageGroup struct {
	// ID
	Id int64 `json:"id"`
}

// UpdateManageGroupContents 修改管理组内容
type UpdateManageGroupContents struct {
	Id         int64       `json:"id"`         // Id 管理组ID
	Values     []int64     `json:"values"`     // Values ID列表
	ValueIf    interface{} `json:"valueIf"`    // 由于 Values 值是特定的 []int64，为了兼容已经存在的接口，再加一个字段 valueIf（value interface），用于存放其他值。这两者互斥，如果 values 有值，则优先使用 values。
	Key        string      `json:"key"`        // Key 属性名称 user_ids|app_package_ids|dept_ids|role_ids|app_ids
	SourceFrom string      `json:"sourceFrom"` // SourceFrom 调用来源。默认为空，如果是极星，则传 `polaris`
	AuthToken  string      `json:"authToken"`  // AuthToken 校验的 token，做一些特殊操作时，需要带上验证 token，如：更换管理员。
}
