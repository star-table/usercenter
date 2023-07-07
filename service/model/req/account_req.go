package req

// 找回密码请求结构体
type RetrievePasswordReq struct {
	// 账号，可以是邮箱或者手机号
	Username string `json:"username"`
	// 验证码
	AuthCode string `json:"authCode"`
	// 新密码
	NewPassword string `json:"newPassword"`
}

// 重新设置密码请求结构体
type ResetPasswordReq struct {
	// 当前密码
	CurrentPassword string `json:"currentPassword"`
	// 新密码
	NewPassword string `json:"newPassword"`
}

// 修改组织成员状态请求结构体
type UpdateOrgMemberStatusReq struct {
	// 要修改的组织成员列表
	MemberIds []int64 `json:"memberIds"`
	// 状态,  1可用,2禁用,3离职
	Status int `json:"status"`
}

// 移除组织
type RemoveOrgMemberReq struct {
	// 要移除的组织成员列表
	MemberIds []int64 `json:"memberIds"`
}

//清空用户
type EmptyUserReq struct {
	// 状态,  1可用,2禁用
	Status int `json:"status"`
}

//退出登录
type UserQuitReq struct {
	Token string `json:"token"`
}

// 解绑登录方式请求结构体（只剩下一种登录方式的时候不允许解绑）
type UnbindLoginNameReq struct {
	// 地址类型: 1：手机号，2：邮箱
	AddressType int `json:"addressType"`
	// 验证码
	AuthCode string `json:"authCode"`
}

// 绑定手机号或者邮箱请求结构体
type BindLoginNameReq struct {
	// 登录地址，手机号或者邮箱
	Address string `json:"address"`
	// 地址类型: 1：手机号，2：邮箱
	AddressType int `json:"addressType"`
	// 验证码
	AuthCode string `json:"authCode"`
	// 邮箱、手机号换绑时，需要的校验码
	ChangeBindCode string `json:"changeBindCode"`
}

// SetPasswordReq 设置密码请求数据模型
type SetPasswordReq struct {
	UserId int64       `json:"userId"`
	Input  PasswordReq `json:"input"`
}

// PasswordReq 密码数据模型
type PasswordReq struct {
	// 密码
	Password string `json:"password"`
}

// UpdatePasswordReq 修改密码请求数据模型
type UpdatePasswordReq struct {
	OrgId  int64             `json:"orgId"`
	UserId int64             `json:"userId"`
	Input  OldNewPasswordReq `json:"input"`
}

// OldNewPasswordReq 新老密码数据模型
type OldNewPasswordReq struct {
	// 当前密码
	CurrentPassword string `json:"currentPassword"`
	// 新密码
	NewPassword string `json:"newPassword"`
}

// UpdatePasswordByLoginNameReq 修改密码请求数据模型
type UpdatePasswordByLoginNameReq struct {
	Username string            `json:"username"` // Username 用户名，手机号，邮箱
	Input    OldNewPasswordReq `json:"input"`
}
