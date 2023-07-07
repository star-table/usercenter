package inner_req

// AddPkgReq
type AddPkgReq struct {
	OrgId  int64 `json:"orgId"`  // OrgId
	PkgId  int64 `json:"pkgId"`  // PkgId
	UserId int64 `json:"userId"` // UserId
}

// DeletePkgReq
type DeletePkgReq struct {
	OrgId  int64 `json:"orgId"`  // OrgId
	PkgId  int64 `json:"pkgId"`  // PkgId
	UserId int64 `json:"userId"` // UserId
}

// AddAppReq
type AddAppReq struct {
	OrgId  int64 `json:"orgId"`  // OrgId
	AppId  int64 `json:"appId"`  // AppId
	UserId int64 `json:"userId"` // UserId
}

// DeleteAppReq
type DeleteAppReq struct {
	OrgId  int64 `json:"orgId"`  // OrgId
	AppId  int64 `json:"appId"`  // AppId
	UserId int64 `json:"userId"` // UserId
}

type GetManagerReq struct {
	OrgId int64 `json:"orgId"`
}

type GetManageGroupTreeReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type UpdateManageGroupContents struct {
	Id             int64   `json:"id"`     // Id 管理组ID
	Values         []int64 `json:"values"` // Values ID列表
	Key            string  `json:"key"`    // Key 属性名称 user_ids|app_package_ids|dept_ids|role_ids|app_ids
	OperatorUserId int64   `json:"operatorUserId"`
	OrgId          int64   `json:"orgId"`
}

type AddUserToSysManageGroupReq struct {
	UserIds []int64 `json:"userIds"`
}

type AddNewMenuToRoleReq struct {
	OrgIds []int64 `json:"orgIds"`
	// 待追加的权限项标识
	PrmItems []string `json:"prmItems"`
}
